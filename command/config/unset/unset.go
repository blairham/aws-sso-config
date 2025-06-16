package unset

import (
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/blairham/aws-sso-config/command/config/shared"
	appconfig "github.com/blairham/aws-sso-config/providers/config"
)

type cmd struct {
	UI cli.Ui
}

func New(ui cli.Ui) *cmd {
	return &cmd{UI: ui}
}

func (c *cmd) Run(args []string) int {
	if len(args) != 1 {
		c.UI.Error("Usage: aws-sso-config config unset <key>")
		c.UI.Error("")
		shared.PrintAvailableKeys(c.UI)
		return 1
	}

	key := args[0]

	// Validate the key
	if !shared.IsValidKey(key) {
		c.UI.Error(fmt.Sprintf("Invalid configuration key: %s", key))
		c.UI.Error("")
		shared.PrintAvailableKeys(c.UI)
		return 1
	}

	// Load current config
	config, err := appconfig.LoadConfigForKey("", key)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error loading config: %v", err))
		return 1
	}

	// Get the default value for this key
	defaultValue, err := c.getDefaultValue(key)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error getting default value: %v", err))
		return 1
	}

	// Reset the value to default
	if err := shared.SetConfigValue(config, key, defaultValue); err != nil {
		c.UI.Error(fmt.Sprintf("Error setting config value: %v", err))
		return 1
	}

	// Save the updated configuration
	if err := shared.SaveConfigValue("", key, defaultValue); err != nil {
		c.UI.Error(fmt.Sprintf("Error saving config: %v", err))
		return 1
	}

	c.UI.Output(fmt.Sprintf("Reset %s to default value: %s", key, defaultValue))
	return 0
}

// getDefaultValue returns the default value for a given configuration key
func (c *cmd) getDefaultValue(key string) (string, error) {
	switch key {
	case shared.KeySSOStartURL:
		return appconfig.DefaultSSO().StartURL, nil
	case shared.KeySSORegion:
		return appconfig.DefaultSSO().Region, nil
	case shared.KeySSORole:
		return appconfig.DefaultSSO().Role, nil
	case shared.KeyAWSDefaultRegion:
		return appconfig.DefaultAWS().DefaultRegion, nil
	case shared.KeyAWSConfigFile:
		return appconfig.DefaultAWS().ConfigFile, nil
	default:
		return "", fmt.Errorf("unknown configuration key: %s", key)
	}
}

func (c *cmd) Help() string {
	return `Usage: aws-sso-config config unset <key>

  Reset a configuration value to its default.

  This command removes any custom configuration for the specified key
  and restores it to the default value.

Available configuration keys:
  sso.start_url        Your AWS SSO start URL
  sso.region          AWS region for SSO (e.g., us-east-1)
  sso.role            SSO role name (e.g., AdministratorAccess)
  aws.default_region  Default AWS region for profiles
  aws.config_file     Path to AWS config file

Examples:
  # Reset SSO start URL to default
  aws-sso-config config unset sso.start_url

  # Reset default region to default
  aws-sso-config config unset aws.default_region

  # Reset AWS config file path to default
  aws-sso-config config unset aws.config_file
`
}

func (c *cmd) Synopsis() string {
	return "Reset a configuration value to its default"
}
