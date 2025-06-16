package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssooidc"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/browser"

	appconfig "github.com/blairham/aws-sso-config/providers/config"
)

type SSOCacheEntry struct {
	AccessToken string    `json:"accessToken"`
	ExpiresAt   time.Time `json:"expiresAt"`
}

func ToString(p *string) string {
	return aws.ToString(p)
}

func ConfigFile() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	return home + "/.aws/config", nil
}

func LoadDefaultConfig() aws.Config {
	// load default aws config
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion("us-east-1"))
	if err != nil {
		log.Fatal("error: ", err)
	}

	return cfg
}

func generateToken(cfg aws.Config) *string {
	// create sso oidc client to trigger login flow
	ssooidcClient := ssooidc.NewFromConfig(cfg)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// register your client which is triggering the login flow
	register, err := ssooidcClient.RegisterClient(ctx, &ssooidc.RegisterClientInput{
		ClientName: aws.String("sample-client-name"),
		ClientType: aws.String("public"),
		Scopes:     []string{"sso-portal:*"},
	})
	if err != nil {
		fmt.Printf("Failed to register client: %v\n", err)
		return nil
	}

	// authorize your device using the client registration response
	deviceAuth, err := ssooidcClient.StartDeviceAuthorization(ctx, &ssooidc.StartDeviceAuthorizationInput{
		ClientId:     register.ClientId,
		ClientSecret: register.ClientSecret,
		StartUrl:     aws.String("https://your-sso-portal.awsapps.com/start"), // Replace with your SSO URL
	})
	if err != nil {
		fmt.Printf("Failed to start device authorization: %v\n", err)
		return nil
	}

	// trigger OIDC login. open browser to login and wait for authorization
	url := aws.ToString(deviceAuth.VerificationUriComplete)
	fmt.Printf("Opening browser for AWS SSO login...\n%v\n", url)
	err = browser.OpenURL(url)
	if err != nil {
		fmt.Printf("Failed to open browser automatically. Please manually open: %v\n", url)
	}

	fmt.Println("Waiting for authorization... (this may take a few moments)")

	return pollForToken(ssooidcClient, register, deviceAuth)
}

func getCurrentToken() *string {
	// Best effort attempt to get token from sso cache.
	// If you can't for whatever reason, return nil, and the code will walk the user through generating a token

	usr, err := user.Current()
	if err != nil {
		return nil
	}
	dir := filepath.Join(usr.HomeDir, ".aws/sso/cache")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		filename := filepath.Join(dir, entry.Name())
		byteValue, err := os.ReadFile(filename)
		if err != nil {
			continue
		}
		var cacheEntry SSOCacheEntry
		err = json.Unmarshal(byteValue, &cacheEntry)
		// If the file ends in .json it should probably have valid json, but meh
		if err != nil {
			continue
		}
		// If it didn't fill in the fields, it is probably not a cache entry, some other random json
		if cacheEntry.AccessToken == "" || cacheEntry.ExpiresAt.IsZero() {
			continue
		}
		// Make sure it's not already expired
		if time.Now().After(cacheEntry.ExpiresAt) {
			continue
		}

		return &cacheEntry.AccessToken
	}

	return nil
}

func GetToken(cfg aws.Config) *string {
	token := getCurrentToken()
	if token != nil {
		return token
	}

	return generateToken(cfg)
}

func pollForToken(ssooidcClient *ssooidc.Client, register *ssooidc.RegisterClientOutput, deviceAuth *ssooidc.StartDeviceAuthorizationOutput) *string {
	// Poll for token creation with exponential backoff
	var token *ssooidc.CreateTokenOutput
	var err error
	maxAttempts := 30 // About 5 minutes total
	interval := time.Second * 5

	tokenCtx, tokenCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer tokenCancel()

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		token, err = ssooidcClient.CreateToken(tokenCtx, &ssooidc.CreateTokenInput{
			ClientId:     register.ClientId,
			ClientSecret: register.ClientSecret,
			DeviceCode:   deviceAuth.DeviceCode,
			GrantType:    aws.String("urn:ietf:params:oauth:grant-type:device_code"),
		})

		if err == nil {
			fmt.Println("✓ Authorization successful!")
			break
		}

		// Check if this is an authorization pending error (expected while waiting)
		if strings.Contains(err.Error(), "authorization_pending") || strings.Contains(err.Error(), "slow_down") {
			if attempt%6 == 0 { // Print status every 30 seconds
				fmt.Printf("Still waiting for authorization... (attempt %d/%d)\n", attempt, maxAttempts)
			}
			time.Sleep(interval)
			continue
		}

		// For other errors, break immediately
		fmt.Printf("Authorization error: %v\n", err)
		break
	}

	if err != nil {
		if strings.Contains(err.Error(), "authorization_pending") {
			fmt.Println("Authorization timeout. Please try again.")
		}
		return nil
	}

	return token.AccessToken
}

func GenerateTokenWithConfig(cfg aws.Config, appCfg *appconfig.Config) *string {
	// create sso oidc client to trigger login flow
	ssooidcClient := ssooidc.NewFromConfig(cfg)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// register your client which is triggering the login flow
	register, err := ssooidcClient.RegisterClient(ctx, &ssooidc.RegisterClientInput{
		ClientName: aws.String("aws-sso-config-cli"),
		ClientType: aws.String("public"),
		Scopes:     []string{"sso-portal:*"},
	})
	if err != nil {
		fmt.Printf("Failed to register client: %v\n", err)
		return nil
	}

	// authorize your device using the client registration response
	deviceAuth, err := ssooidcClient.StartDeviceAuthorization(ctx, &ssooidc.StartDeviceAuthorizationInput{
		ClientId:     register.ClientId,
		ClientSecret: register.ClientSecret,
		StartUrl:     aws.String(appCfg.SSOStartURL()),
	})
	if err != nil {
		fmt.Printf("Failed to start device authorization: %v\n", err)
		return nil
	}

	// trigger OIDC login. open browser to login and wait for authorization
	url := aws.ToString(deviceAuth.VerificationUriComplete)
	fmt.Printf("Opening browser for AWS SSO login...\n%v\n", url)
	err = browser.OpenURL(url)
	if err != nil {
		fmt.Printf("Failed to open browser automatically. Please manually open: %v\n", url)
	}

	fmt.Println("Waiting for authorization... (this may take a few moments)")

	return pollForToken(ssooidcClient, register, deviceAuth)
}
