package set

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/mitchellh/go-homedir"

	"github.com/blairham/aws-sso-config/command/config/shared"
)

type cmd struct {
	UI cli.Ui
}

func New(ui cli.Ui) *cmd {
	return &cmd{UI: ui}
}

func (c *cmd) Run(args []string) int {
	if len(args) < 2 {
		c.UI.Error("Usage: aws-sso-config config set <key> <value>")
		c.UI.Error("")
		shared.PrintAvailableKeys(c.UI)
		return 1
	}

	key := args[0]
	// Join all remaining arguments as the value to support values with spaces
	value := strings.Join(args[1:], " ")

	// Validate the key
	if !shared.IsValidKey(key) {
		c.UI.Error(fmt.Sprintf("Invalid configuration key: %s", key))
		c.UI.Error("")
		shared.PrintAvailableKeys(c.UI)
		return 1
	}

	// Get the config file path
	home, err := homedir.Dir()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error getting home directory: %v", err))
		return 1
	}
	configFile := filepath.Join(home, ".awsssoconfig")

	// Use the single-file configuration saving
	if err := shared.SaveConfigValue(configFile, key, value); err != nil {
		c.UI.Error(fmt.Sprintf("Error updating config: %v", err))
		return 1
	}

	c.UI.Output(fmt.Sprintf("Updated %s = %s", key, value))
	return 0
}

func (c *cmd) Help() string {
	return `Usage: aws-sso-config config set <key> <value>

  Set a configuration value.

  The value can be provided with or without quotes. Multiple words
  will be joined with spaces to form the complete value.

Available configuration keys:
  sso.start_url        Your AWS SSO start URL
  sso.region          AWS region for SSO (e.g., us-east-1)
  sso.role            SSO role name (e.g., AdministratorAccess)
  aws.default_region  Default AWS region for profiles
  aws.config_file     Path to AWS config file

Examples:
  # Set the SSO start URL (no quotes needed)
  aws-sso-config config set sso.start_url https://mycompany.awsapps.com/start

  # Set the default region
  aws-sso-config config set aws.default_region us-west-2

  # Set the AWS config file path
  aws-sso-config config set aws.config_file ~/.aws/config

  # Values with spaces work without quotes
  aws-sso-config config set sso.role Administrator Access Role

  # Quotes still work if preferred
  aws-sso-config config set sso.start_url "https://mycompany.awsapps.com/start"
`
}

func (c *cmd) Synopsis() string {
	return "Set a configuration value"
}
