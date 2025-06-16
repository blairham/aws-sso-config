package generate

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sso"
	"github.com/bigkevmcd/go-configparser"
	"github.com/mitchellh/cli"
	"github.com/spf13/pflag"

	generateflags "github.com/blairham/aws-sso-config/command/generate/flags"
	awsprovider "github.com/blairham/aws-sso-config/providers/aws"
	appconfig "github.com/blairham/aws-sso-config/providers/config"
)

// TokenGenerator interface for mocking
type TokenGenerator interface {
	GenerateTokenWithConfig(cfg aws.Config, appCfg *appconfig.Config) *string
}

// DefaultTokenGenerator implements TokenGenerator using the real AWS functions
type DefaultTokenGenerator struct{}

func (g *DefaultTokenGenerator) GenerateTokenWithConfig(cfg aws.Config, appCfg *appconfig.Config) *string {
	return awsprovider.GenerateTokenWithConfig(cfg, appCfg)
}

type cmd struct {
	UI    cli.Ui
	flags *pflag.FlagSet
	help  string

	diff       bool
	configFile string

	// Dependencies for testing
	ssoClientFactory func(aws.Config) SSOClient
	tokenGenerator   TokenGenerator
	configLoader     func() aws.Config
}

func New(ui cli.Ui) *cmd {
	c := &cmd{UI: ui}
	c.Init()
	// Set default dependencies
	c.ssoClientFactory = func(cfg aws.Config) SSOClient {
		return sso.NewFromConfig(cfg)
	}
	c.tokenGenerator = &DefaultTokenGenerator{}
	c.configLoader = awsprovider.LoadDefaultConfig
	return c
}

// NewWithDependencies creates a new command with injected dependencies for testing
func NewWithDependencies(ui cli.Ui, ssoClientFactory func(aws.Config) SSOClient, tokenGenerator TokenGenerator, configLoader func() aws.Config) *cmd {
	c := &cmd{UI: ui}
	c.Init()
	c.ssoClientFactory = ssoClientFactory
	c.tokenGenerator = tokenGenerator
	c.configLoader = configLoader
	return c
}

func (c *cmd) Init() {
	c.flags = pflag.NewFlagSet("generate", pflag.ContinueOnError)

	// Get flag configurations from registry
	registry := generateflags.NewFlagRegistry()
	diffFlag := registry.GetFlagByName("diff")
	configFlag := registry.GetFlagByName("config")

	// Add flags with both short and long forms
	c.flags.BoolVarP(&c.diff, diffFlag.GetFlagName(), diffFlag.GetShortFlag(), false, diffFlag.GetDescription())
	c.flags.StringVarP(&c.configFile, configFlag.GetFlagName(), configFlag.GetShortFlag(), "", configFlag.GetDescription())

	c.help = c.buildHelp()
}

func (c *cmd) Run(args []string) int {
	if err := c.flags.Parse(args); err != nil {
		return 1
	}

	// Load configuration
	var appCfg *appconfig.Config
	var err error

	if c.configFile != "" {
		appCfg, err = appconfig.Load(c.configFile)
		if err != nil {
			c.UI.Error(fmt.Sprintf("Configuration error: %v", err))
			return 1
		}
	} else {
		appCfg = appconfig.Default()
	}

	if err := appCfg.Validate(); err != nil {
		c.UI.Error(fmt.Sprintf("Configuration error: %v", err))
		return 1
	}

	configFile := appCfg.ConfigFile()

	cfg := c.configLoader()
	token := c.tokenGenerator.GenerateTokenWithConfig(cfg, appCfg)

	// create sso client
	ssoClient := c.ssoClientFactory(cfg)

	if generateAwsConfigFile(ssoClient, token, configFile, c.diff, appCfg) != nil {
		return 1
	}

	return 0
}

func (c *cmd) buildHelp() string {
	helpText := help + "\n"
	if c.flags != nil {
		helpText += c.flags.FlagUsages()
	}
	return helpText
}

func (c *cmd) Help() string {
	return c.help
}

func (c *cmd) Synopsis() string {
	return synopsis
}

func showFileDiff(file1, file2 string) {
	// First check to make sure these are legit filenames so we don't confuse the "diff" command
	for _, file := range []string{file1, file2} {
		_, err := os.Stat(file)
		if err != nil {
			fmt.Printf("File %s does not exist\n", file)
			return
		}
	}
	cmd := exec.Command("diff", file1, file2)
	cmd.Stdout = os.Stdout
	cmd.Run() // Ignore error as diff returns non-zero when files differ
}

func generateAwsConfigFile(ssoClient SSOClient, token *string, configFile string, diff bool, appCfg *appconfig.Config) error {
	configFileNew := configFile + ".new"

	awsConfig, err := configparser.NewConfigParserFromFile(configFile)
	if err != nil {
		log.Fatal(err)
		return err
	}

	fmt.Println("Fetching list of all accounts for user")

	// Get accounts (simplified - no pagination for testing compatibility)
	listAccountsInput := &sso.ListAccountsInput{
		AccessToken: token,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	accountsResult, err := ssoClient.ListAccounts(ctx, listAccountsInput)
	if err != nil {
		fmt.Printf("Error fetching accounts: %v\n", err)
		return err
	}

	for _, y := range accountsResult.AccountList {
		// Add all accounts - users can configure filtering if needed
		accountName := aws.ToString(y.AccountName)

		// Use the account name as-is for the profile name
		// Users can customize this logic based on their naming conventions
		profileName := accountName
		section := "profile " + profileName

		// check if profile already exists and update it
		if !awsConfig.HasSection(section) {
			fmt.Printf("Adding profile %v\n", profileName)
			awsConfig.AddSection(section)
		}

		awsConfig.Set(section, "sso_account_id", aws.ToString(y.AccountId))
		awsConfig.Set(section, "sso_role_name", appCfg.SSORole())
		awsConfig.Set(section, "sso_region", appCfg.SSORegion())
		awsConfig.Set(section, "sso_start_url", appCfg.SSOStartURL())
		awsConfig.Set(section, "region", appCfg.DefaultRegion())
	}

	err = awsConfig.SaveWithDelimiter(configFileNew, "=")
	if err != nil {
		return fmt.Errorf("failed to save config file: %w", err)
	}
	if diff {
		showFileDiff(configFile, configFileNew)
	}
	err = os.Rename(configFileNew, configFile)
	if err != nil {
		return fmt.Errorf("failed to rename config file: %w", err)
	}

	return nil
}
