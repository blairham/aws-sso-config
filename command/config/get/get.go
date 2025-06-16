package get

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
		c.UI.Error("Usage: aws-sso-config config get <key>")
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

	// Load current config (only what's needed for this key)
	config, err := appconfig.LoadConfigForKey("", key)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error loading config: %v", err))
		return 1
	}

	// Get the value for the specified key
	value, err := shared.GetConfigValue(config, key)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	c.UI.Output(value)
	return 0
}

func (c *cmd) Help() string {
	return `Usage: aws-sso-config config get <key>

  Get a configuration value.

Available configuration keys:
  sso.start_url        Your AWS SSO start URL
  sso.region          AWS region for SSO (e.g., us-east-1)
  sso.role            SSO role name (e.g., AdministratorAccess)
  aws.default_region  Default AWS region for profiles
  aws.config_file     Path to AWS config file

Examples:
  # Get the SSO start URL
  aws-sso-config config get sso.start_url

  # Get the default region
  aws-sso-config config get aws.default_region
`
}

func (c *cmd) Synopsis() string {
	return "Get a configuration value"
}
