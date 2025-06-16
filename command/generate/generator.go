package generate

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sso"
	"github.com/aws/aws-sdk-go-v2/service/sso/types"
	"github.com/bigkevmcd/go-configparser"

	appconfig "github.com/blairham/aws-sso-config/providers/config"
)

// Interface for AWS SSO operations to allow mocking in tests
type SSOClient interface {
	GetRoleCredentials(ctx context.Context, params *sso.GetRoleCredentialsInput, optFns ...func(*sso.Options)) (*sso.GetRoleCredentialsOutput, error)
	ListAccounts(ctx context.Context, params *sso.ListAccountsInput, optFns ...func(*sso.Options)) (*sso.ListAccountsOutput, error)
	ListAccountRoles(ctx context.Context, params *sso.ListAccountRolesInput, optFns ...func(*sso.Options)) (*sso.ListAccountRolesOutput, error)
}

// ConfigGeneratorIface defines the interface for config generator operations
type ConfigGeneratorIface interface {
	GenerateConfigFile(ssoClient SSOClient, token *string, configFile string, showDiff bool, appCfg *appconfig.Config) error
	ListAccountsWithClient(ssoClient SSOClient, token *string) ([]types.AccountInfo, error)
	GetAccountRolesWithClient(ssoClient SSOClient, token *string, accountID string) ([]types.RoleInfo, error)
	WriteSectionToConfig(configParser *configparser.ConfigParser, sectionName string, values map[string]string) error
}

// ConfigGenerator holds the config generator implementation
type ConfigGenerator struct {
	SSOStartURL   string
	SSORegion     string
	DefaultRegion string
}

// Returns a new config generator with the given parameters
func NewConfigGenerator(ssoStartURL, ssoRegion, defaultRegion string) *ConfigGenerator {
	return &ConfigGenerator{
		SSOStartURL:   ssoStartURL,
		SSORegion:     ssoRegion,
		DefaultRegion: defaultRegion,
	}
}

// ListAccountsWithClient lists AWS accounts with the provided SSO client
func (g *ConfigGenerator) ListAccountsWithClient(ssoClient SSOClient, token *string) ([]types.AccountInfo, error) {
	if token == nil {
		return nil, errors.New("no SSO token provided")
	}

	// List accounts
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	accountsOutput, err := ssoClient.ListAccounts(ctx, &sso.ListAccountsInput{
		AccessToken: token,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list accounts: %w", err)
	}

	return accountsOutput.AccountList, nil
}

// GetAccountRolesWithClient gets AWS account roles using the provided SSO client
func (g *ConfigGenerator) GetAccountRolesWithClient(ssoClient SSOClient, token *string, accountID string) ([]types.RoleInfo, error) {
	if token == nil {
		return nil, errors.New("no SSO token provided")
	}

	// List roles for this account
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	rolesOutput, err := ssoClient.ListAccountRoles(ctx, &sso.ListAccountRolesInput{
		AccessToken: token,
		AccountId:   aws.String(accountID),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list roles for account %s: %w", accountID, err)
	}

	return rolesOutput.RoleList, nil
}

// WriteSectionToConfig writes a section to the config parser
func (g *ConfigGenerator) WriteSectionToConfig(configParser *configparser.ConfigParser, sectionName string, values map[string]string) error {
	// Add section if it doesn't exist
	if !configParser.HasSection(sectionName) {
		err := configParser.AddSection(sectionName)
		if err != nil {
			return fmt.Errorf("failed to add section %s: %w", sectionName, err)
		}
	}

	// Write values
	for key, value := range values {
		err := configParser.Set(sectionName, key, value)
		if err != nil {
			return fmt.Errorf("failed to set %s=%s in section %s: %w", key, value, sectionName, err)
		}
	}

	return nil
}
