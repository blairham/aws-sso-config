package aws

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssooidc"
	"github.com/pkg/browser"

	appconfig "github.com/blairham/aws-sso-config/providers/config"
)

// Interface for SSO OIDC operations to allow mocking in tests
type SSOOIDCClient interface {
	RegisterClient(ctx context.Context, params *ssooidc.RegisterClientInput, optFns ...func(*ssooidc.Options)) (*ssooidc.RegisterClientOutput, error)
	StartDeviceAuthorization(
		ctx context.Context,
		params *ssooidc.StartDeviceAuthorizationInput,
		optFns ...func(*ssooidc.Options),
	) (*ssooidc.StartDeviceAuthorizationOutput, error)
	CreateToken(ctx context.Context, params *ssooidc.CreateTokenInput, optFns ...func(*ssooidc.Options)) (*ssooidc.CreateTokenOutput, error)
}

// Wrapper around AWS functions to facilitate testing
type AWSProvider struct {
	SSOOIDCClient SSOOIDCClient
	BrowserOpener func(string) error
	TokenPoller   func(client SSOOIDCClient, register *ssooidc.RegisterClientOutput, deviceAuth *ssooidc.StartDeviceAuthorizationOutput) *string
	Cfg           aws.Config
}

// Creates a new default AWS provider
func NewDefaultAWSProvider() *AWSProvider {
	cfg := LoadDefaultConfig()
	ssooidcClient := ssooidc.NewFromConfig(cfg)

	return &AWSProvider{
		SSOOIDCClient: ssooidcClient,
		BrowserOpener: browser.OpenURL,
		TokenPoller: func(client SSOOIDCClient, register *ssooidc.RegisterClientOutput, deviceAuth *ssooidc.StartDeviceAuthorizationOutput) *string {
			// We need to cast the interface to the concrete type
			if ssoClient, ok := client.(*ssooidc.Client); ok {
				return pollForToken(ssoClient, register, deviceAuth)
			}
			return nil
		},
		Cfg: cfg,
	}
}

// Generate token with provider's configuration
func (p *AWSProvider) GenerateToken(appCfg *appconfig.Config) *string {
	// create sso oidc client to trigger login flow
	ssooidcClient := p.SSOOIDCClient

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// register your client which is triggering the login flow
	register, err := ssooidcClient.RegisterClient(ctx, &ssooidc.RegisterClientInput{
		ClientName: aws.String("aws-sso-config-cli"),
		ClientType: aws.String("public"),
		Scopes:     []string{"sso-portal:*"},
	})
	if err != nil {
		return nil
	}

	// authorize your device using the client registration response
	deviceAuth, err := ssooidcClient.StartDeviceAuthorization(ctx, &ssooidc.StartDeviceAuthorizationInput{
		ClientId:     register.ClientId,
		ClientSecret: register.ClientSecret,
		StartUrl:     aws.String(appCfg.SSOStartURL()),
	})
	if err != nil {
		return nil
	}

	// trigger OIDC login. open browser to login and wait for authorization
	url := aws.ToString(deviceAuth.VerificationUriComplete)
	_ = p.BrowserOpener(url)

	return p.TokenPoller(ssooidcClient, register, deviceAuth)
}
