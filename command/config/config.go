package config

import (
	"fmt"

	"github.com/mitchellh/cli"
	"github.com/spf13/pflag"

	"github.com/blairham/aws-sso-config/command/config/edit"
	configflags "github.com/blairham/aws-sso-config/command/config/flags"
	"github.com/blairham/aws-sso-config/command/config/get"
	"github.com/blairham/aws-sso-config/command/config/list"
	"github.com/blairham/aws-sso-config/command/config/set"
	"github.com/blairham/aws-sso-config/command/config/unset"
)

type cmd struct {
	UI    cli.Ui
	flags *pflag.FlagSet

	// Flag variables
	list bool
}

func New(ui cli.Ui) *cmd {
	c := &cmd{UI: ui}
	c.init()
	return c
}

func (c *cmd) init() {
	c.flags = pflag.NewFlagSet("config", pflag.ContinueOnError)

	// Get flag configurations from registry
	registry := configflags.NewFlagRegistry()
	listFlag := registry.GetFlagByName("list")

	// Add the list flag as a secret flag (won't show in help)
	c.flags.BoolVarP(&c.list, listFlag.GetFlagName(), listFlag.GetShortFlag(), false, listFlag.GetDescription())
	c.flags.MarkHidden("list") // Make it a secret flag
}

func (c *cmd) Run(args []string) int {
	// Parse flags
	if err := c.flags.Parse(args); err != nil {
		return 1
	}

	// Get remaining args after flag parsing
	remainingArgs := c.flags.Args()

	// Check if list flag was used
	if c.list {
		listCmd := list.New(c.UI)
		return listCmd.Run([]string{})
	}

	if len(remainingArgs) == 0 {
		c.showUsage()
		return 1
	}

	subcommand := remainingArgs[0]
	subArgs := remainingArgs[1:]

	switch subcommand {
	case "get":
		getCmd := get.New(c.UI)
		return getCmd.Run(subArgs)
	case "set":
		setCmd := set.New(c.UI)
		return setCmd.Run(subArgs)
	case "unset":
		unsetCmd := unset.New(c.UI)
		return unsetCmd.Run(subArgs)
	case "list":
		listCmd := list.New(c.UI)
		return listCmd.Run(subArgs)
	case "edit":
		editCmd := edit.New(c.UI)
		return editCmd.Run(subArgs)
	default:
		c.UI.Error(fmt.Sprintf("Unknown subcommand: %s", subcommand))
		c.UI.Error("")
		c.showUsage()
		return 1
	}
}

// showUsage displays the usage information including available flags
func (c *cmd) showUsage() {
	c.UI.Error("Usage: aws-sso-config config <subcommand>")
	c.UI.Error("")
	c.UI.Error("Available subcommands:")
	c.UI.Error("  get <key>             Get a configuration value")
	c.UI.Error("  set <key> <value>     Set a configuration value")
	c.UI.Error("  unset <key>           Reset a configuration value to its default")
	c.UI.Error("  list                  List all configuration variables and their values")
	c.UI.Error("  edit [config-file]    Open configuration file in an editor")
	c.UI.Error("")

	// Show flags (note: hidden flags like --list won't appear here)
	if c.flags.HasFlags() {
		c.UI.Error("Flags:")
		c.UI.Error(c.flags.FlagUsages())
	}
}

func (c *cmd) Help() string {
	help := `Usage: aws-sso-config config <subcommand>

  Manage configuration settings for aws-sso-config.

Subcommands:
  get <key>            Get a configuration value
  set <key> <value>    Set a configuration value
  unset <key>          Reset a configuration value to its default
  list                 List all configuration variables and their values
  edit [config-file]   Open configuration file in an editor

Available configuration keys:
  sso.start_url        Your AWS SSO start URL
  sso.region          AWS region for SSO (e.g., us-east-1)
  sso.role            SSO role name (e.g., AdministratorAccess)
  aws.default_region  Default AWS region for profiles
  aws.config_file     Path to AWS config file

Examples:
  # Get the SSO start URL
  aws-sso-config config get sso.start_url

  # Set the SSO start URL (no quotes needed)
  aws-sso-config config set sso.start_url https://mycompany.awsapps.com/start

  # Reset SSO start URL to default
  aws-sso-config config unset sso.start_url

  # Set the default region
  aws-sso-config config set aws.default_region us-west-2

  # Reset default region to default
  aws-sso-config config unset aws.default_region

  # Set the AWS config file path
  aws-sso-config config set aws.config_file ~/.aws/config

  # List all configuration variables and their values
  aws-sso-config config list

  # Edit the configuration file
  aws-sso-config config edit

  # Edit a specific configuration file
  aws-sso-config config edit /path/to/config
`

	// Add pflag usage information (excluding hidden flags)
	if c.flags.HasAvailableFlags() {
		help += "\nFlags:\n" + c.flags.FlagUsages()
	}

	return help
}

func (c *cmd) Synopsis() string {
	return "Read and write configuration values"
}
