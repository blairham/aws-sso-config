package edit

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/mitchellh/cli"
	"github.com/mitchellh/go-homedir"

	appconfig "github.com/blairham/aws-sso-config/providers/config"
)

type cmd struct {
	UI cli.Ui
}

func New(ui cli.Ui) *cmd {
	return &cmd{UI: ui}
}

func (c *cmd) Run(args []string) int {
	// Determine config file path
	var configFile string
	if len(args) > 0 {
		configFile = args[0]
	} else {
		// Use default config file location
		home, err := homedir.Dir()
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error getting home directory: %v", err))
			return 1
		}
		configFile = filepath.Join(home, ".awsssoconfig")
	}

	// Ensure config file exists
	if err := c.ensureConfigFileExists(configFile); err != nil {
		c.UI.Error(fmt.Sprintf("Error ensuring config file exists: %v", err))
		return 1
	}

	// Get editor from environment or use default
	editor := os.Getenv("EDITOR")
	if editor == "" {
		// Default to vim on Unix-like systems, or notepad on Windows
		if _, err := exec.LookPath("vim"); err == nil {
			editor = "vim"
		} else if _, err := exec.LookPath("nano"); err == nil {
			editor = "nano"
		} else if _, err := exec.LookPath("vi"); err == nil {
			editor = "vi"
		} else {
			c.UI.Error("No editor found. Please set the EDITOR environment variable.")
			c.UI.Error("For example: export EDITOR=vim")
			return 1
		}
	}

	// Validate that the editor exists and is executable
	editorPath, err := exec.LookPath(editor)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Editor '%s' not found in PATH", editor))
		return 1
	}

	// Open the config file in the editor
	// #nosec G204 - This is intentional: we want to launch the user's preferred editor
	// The editor path is validated above using exec.LookPath
	cmd := exec.Command(editorPath, configFile)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		c.UI.Error(fmt.Sprintf("Error running editor: %v", err))
		return 1
	}

	// Validate the config file after editing
	if err := c.validateConfigFile(configFile); err != nil {
		c.UI.Error(fmt.Sprintf("Warning: Configuration file validation failed: %v", err))
		c.UI.Error("Please check your configuration file for syntax errors.")
		return 1
	}

	c.UI.Output("Configuration file edited successfully.")
	return 0
}

// ensureConfigFileExists creates the config file if it doesn't exist
func (c *cmd) ensureConfigFileExists(configFile string) error {
	// Check if file exists
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		// Create the directory if it doesn't exist
		dir := filepath.Dir(configFile)
		if err := os.MkdirAll(dir, 0750); err != nil {
			return fmt.Errorf("failed to create config directory: %w", err)
		}

		// Create a default config file
		cm := appconfig.NewConfigManager(configFile)
		config, err := cm.Load() // This will create the file with defaults
		if err != nil {
			return fmt.Errorf("failed to create default config file: %w", err)
		}

		// Validate the created config
		if err := config.Validate(); err != nil {
			return fmt.Errorf("default config validation failed: %w", err)
		}

		c.UI.Output(fmt.Sprintf("Created default configuration file at: %s", configFile))
	}

	return nil
}

// validateConfigFile validates the configuration file after editing
func (c *cmd) validateConfigFile(configFile string) error {
	cm := appconfig.NewConfigManager(configFile)
	config, err := cm.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	return config.Validate()
}

func (c *cmd) Help() string {
	return `Usage: aws-sso-config config edit [config-file]

  Open the configuration file in an editor for manual editing.

  The editor used is determined by the EDITOR environment variable.
  If EDITOR is not set, the command will try to use vim, nano, or vi
  in that order.

Arguments:
  config-file    Optional path to the configuration file.
                 If not provided, uses the default location (~/.awsssoconfig)

Examples:
  # Edit the default configuration file
  aws-sso-config config edit

  # Edit a specific configuration file
  aws-sso-config config edit /path/to/config

  # Set your preferred editor
  export EDITOR=nano
  aws-sso-config config edit

Environment Variables:
  EDITOR         The editor to use for editing the configuration file
`
}

func (c *cmd) Synopsis() string {
	return "Open configuration file in an editor"
}
