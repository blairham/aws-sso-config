package list

import (
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/blairham/aws-sso-config/command/config/shared"
	"github.com/blairham/aws-sso-config/internal/pager"
	appconfig "github.com/blairham/aws-sso-config/providers/config"
)

type cmd struct {
	UI cli.Ui
}

func New(ui cli.Ui) *cmd {
	return &cmd{UI: ui}
}

func (c *cmd) Run(args []string) int {
	// Parse arguments for paging options
	forcePaging := false
	var filteredArgs []string

	for _, arg := range args {
		if arg == "--force-paging" {
			forcePaging = true
		} else {
			filteredArgs = append(filteredArgs, arg)
		}
	}

	if len(filteredArgs) != 0 {
		c.UI.Error("Usage: aws-sso-config config list [--force-paging]")
		c.UI.Error("")
		c.UI.Error("This command takes no arguments except optional flags.")
		return 1
	}

	// Load current configuration
	config, err := appconfig.Load("")
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error loading config: %v", err))
		return 1
	}

	// Collect all output lines for paging
	var outputLines []string

	// Collect all configuration values in key=value format like git config --list
	for _, key := range shared.ValidKeys {
		value, err := shared.GetConfigValue(config, key)
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error getting value for %s: %v", key, err))
			continue
		}
		// Only output non-empty values
		if value != "" {
			outputLines = append(outputLines, fmt.Sprintf("%s=%s", key, value))
		}
	}

	// Use pager for output
	p := pager.New(c.UI)
	if forcePaging {
		p.SetForceEnabled(true)
	}
	p.Output(outputLines)

	return 0
}

func (c *cmd) Help() string {
	return `Usage: aws-sso-config config list [--force-paging]

  List all configuration variables set in the config file with their values.
  Output format is key=value, one per line, similar to 'git config --list'.

  The output will automatically use an interactive pager (like 'less') when the
  output would be too long for the terminal screen. The pager provides full
  navigation with arrow keys, search functionality, and more.

Interactive Pager Controls:
  ↑/↓ Arrow keys      Navigate up/down one line
  Space / Page Down   Navigate down one page
  b / Page Up         Navigate up one page
  g / Home            Go to beginning
  G / End             Go to end
  / text              Search forward for 'text'
  ? text              Search backward for 'text'
  n                   Next search result
  N                   Previous search result
  q                   Quit the pager
  h                   Show help (in less)

Flags:
  --force-paging       Force paging even for short output (for testing)

Environment Variables:
  NO_PAGER            Disable paging entirely
  PAGER               Specify which pager to use (default: less -R)
  AWS_SSO_CONFIG_PAGER  Specify pager for this tool specifically
  LINES               Override terminal height detection

Examples:
  # List all configuration values (auto-paging)
  aws-sso-config config list

  # Force interactive paging (useful for testing)
  aws-sso-config config list --force-paging

  # Disable paging completely
  NO_PAGER=1 aws-sso-config config list

  # Use a specific pager
  PAGER="more" aws-sso-config config list

  # Use less with specific options
  PAGER="less -S" aws-sso-config config list
`
}

func (c *cmd) Synopsis() string {
	return "List all configuration variables and their values"
}
