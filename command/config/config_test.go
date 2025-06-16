package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
	"github.com/stretchr/testify/assert"

	"github.com/blairham/aws-sso-config/command/config/shared"
)

func TestConfig(t *testing.T) {
	ui := cli.NewMockUi()
	c := New(ui)

	// Verify the command is properly initialized
	assert.NotNil(t, c)
	assert.Equal(t, "Read and write configuration values", c.Synopsis())
}

func TestConfigNoArgs(t *testing.T) {
	ui := cli.NewMockUi()
	c := New(ui)

	exitCode := c.Run([]string{})
	assert.Equal(t, 1, exitCode)

	errorOutput := ui.ErrorWriter.String()
	assert.Contains(t, errorOutput, "Usage: aws-sso-config config <subcommand>")
	assert.Contains(t, errorOutput, "get <key>")
	assert.Contains(t, errorOutput, "set <key> <value>")
	assert.Contains(t, errorOutput, "list")
}

func TestConfigInvalidSubcommand(t *testing.T) {
	ui := cli.NewMockUi()
	c := New(ui)

	exitCode := c.Run([]string{"invalid"})
	assert.Equal(t, 1, exitCode)

	errorOutput := ui.ErrorWriter.String()
	assert.Contains(t, errorOutput, "Unknown subcommand: invalid")
}

func TestConfigGetInvalidArgs(t *testing.T) {
	ui := cli.NewMockUi()
	c := New(ui)

	// No key provided
	exitCode := c.Run([]string{"get"})
	assert.Equal(t, 1, exitCode)

	errorOutput := ui.ErrorWriter.String()
	assert.Contains(t, errorOutput, "Usage: aws-sso-config config get <key>")
}

func TestConfigSetInvalidArgs(t *testing.T) {
	ui := cli.NewMockUi()
	c := New(ui)

	// No key or value provided
	exitCode := c.Run([]string{"set"})
	assert.Equal(t, 1, exitCode)

	errorOutput := ui.ErrorWriter.String()
	assert.Contains(t, errorOutput, "Usage: aws-sso-config config set <key> <value>")

	// Only key provided
	ui = cli.NewMockUi()
	c = New(ui)
	exitCode = c.Run([]string{"set", "sso_start_url"})
	assert.Equal(t, 1, exitCode)
}

func TestConfigGetSet(t *testing.T) {
	// Setup temporary home directory
	tempHome := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempHome)
	defer os.Setenv("HOME", originalHome)

	configPath := filepath.Join(tempHome, ".awsssoconfig")

	// Test set
	ui := cli.NewMockUi()
	c := New(ui)

	exitCode := c.Run([]string{"set", "sso.start_url", "https://test.awsapps.com/start"})
	assert.Equal(t, 0, exitCode)

	output := ui.OutputWriter.String()
	assert.Contains(t, output, "Updated sso.start_url = https://test.awsapps.com/start")

	// Verify file was created
	_, err := os.Stat(configPath)
	assert.NoError(t, err)

	// Test get
	ui = cli.NewMockUi()
	c = New(ui)

	exitCode = c.Run([]string{"get", "sso.start_url"})
	assert.Equal(t, 0, exitCode)

	output = ui.OutputWriter.String()
	assert.Contains(t, output, "https://test.awsapps.com/start")
}

func TestConfigSetInvalidKey(t *testing.T) {
	ui := cli.NewMockUi()
	c := New(ui)

	exitCode := c.Run([]string{"set", "invalid_key", "value"})
	assert.Equal(t, 1, exitCode)

	errorOutput := ui.ErrorWriter.String()
	assert.Contains(t, errorOutput, "Invalid configuration key: invalid_key")
}

func TestConfigGetInvalidKey(t *testing.T) {
	// Setup temporary home directory
	tempHome := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempHome)
	defer os.Setenv("HOME", originalHome)

	ui := cli.NewMockUi()
	c := New(ui)

	exitCode := c.Run([]string{"get", "invalid_key"})
	assert.Equal(t, 1, exitCode)

	errorOutput := ui.ErrorWriter.String()
	assert.Contains(t, errorOutput, "Invalid configuration key: invalid_key")
}

func TestConfigHelp(t *testing.T) {
	ui := cli.NewMockUi()
	c := New(ui)

	help := c.Help()
	assert.Contains(t, help, "Usage: aws-sso-config config <subcommand>")
	assert.Contains(t, help, "get <key>")
	assert.Contains(t, help, "set <key> <value>")
	assert.Contains(t, help, "list")
	assert.Contains(t, help, "sso.start_url")
	assert.Contains(t, help, "Examples:")

	// The --list flag should NOT appear in help since it's a hidden/secret flag
	assert.NotContains(t, help, "-l, --list")
	assert.NotContains(t, help, "[-l|--list]")
}

func TestConfigList(t *testing.T) {
	ui := cli.NewMockUi()
	c := New(ui)

	exitCode := c.Run([]string{"list"})
	assert.Equal(t, 0, exitCode)

	output := ui.OutputWriter.String()
	// The output should be in key=value format
	lines := strings.Split(strings.TrimSpace(output), "\n")
	for _, line := range lines {
		if line != "" && !strings.Contains(line, "=") {
			t.Errorf("Expected line to be in key=value format, got %s", line)
		}
	}
}

func TestConfigListFlag(t *testing.T) {
	ui := cli.NewMockUi()
	c := New(ui)

	exitCode := c.Run([]string{"-l"})
	assert.Equal(t, 0, exitCode)

	output := ui.OutputWriter.String()
	// The output should be in key=value format
	lines := strings.Split(strings.TrimSpace(output), "\n")
	for _, line := range lines {
		if line != "" && !strings.Contains(line, "=") {
			t.Errorf("Expected line to be in key=value format, got %s", line)
		}
	}
}

func TestConfigListLongFlag(t *testing.T) {
	ui := cli.NewMockUi()
	c := New(ui)

	exitCode := c.Run([]string{"--list"})
	assert.Equal(t, 0, exitCode)

	output := ui.OutputWriter.String()
	// The output should be in key=value format
	lines := strings.Split(strings.TrimSpace(output), "\n")
	for _, line := range lines {
		if line != "" && !strings.Contains(line, "=") {
			t.Errorf("Expected line to be in key=value format, got %s", line)
		}
	}
}

func TestConfigListFlagWithOtherArgs(t *testing.T) {
	ui := cli.NewMockUi()
	c := New(ui)

	// Test that -l flag takes precedence and ignores other args
	exitCode := c.Run([]string{"-l", "some", "extra", "args"})
	assert.Equal(t, 0, exitCode)

	output := ui.OutputWriter.String()
	// The output should be in key=value format
	lines := strings.Split(strings.TrimSpace(output), "\n")
	for _, line := range lines {
		if line != "" && !strings.Contains(line, "=") {
			t.Errorf("Expected line to be in key=value format, got %s", line)
		}
	}
}

func TestConfigFlagParsing(t *testing.T) {
	ui := cli.NewMockUi()
	c := New(ui)

	// Test that the flag system is properly initialized
	assert.NotNil(t, c.flags)

	// Test that the list flag is available and hidden
	flag := c.flags.Lookup("list")
	assert.NotNil(t, flag, "Expected list flag to be available")
	assert.True(t, flag.Hidden, "Expected list flag to be hidden")
}

func TestConfigUnsetSubcommand(t *testing.T) {
	ui := cli.NewMockUi()
	c := New(ui)

	// Test with no arguments to unset
	exitCode := c.Run([]string{"unset"})
	assert.Equal(t, 1, exitCode)

	errorOutput := ui.ErrorWriter.String()
	assert.Contains(t, errorOutput, "Usage: aws-sso-config config unset <key>")
}

func TestConfigUnsetHelp(t *testing.T) {
	ui := cli.NewMockUi()
	c := New(ui)

	help := c.Help()
	assert.Contains(t, help, "unset <key>")
	assert.Contains(t, help, "Reset a configuration value to its default")
	assert.Contains(t, help, "aws-sso-config config unset")
}

func TestValidKey(t *testing.T) {
	// Test valid keys
	validKeys := []string{
		"sso.start_url",
		"sso.region",
		"sso.role",
		"aws.default_region",
		"aws.config_file",
	}

	for _, key := range validKeys {
		assert.True(t, shared.IsValidKey(key), "Key %s should be valid", key)
	}

	// Test invalid keys
	invalidKeys := []string{
		"invalid_key",
		"",
		"sso_start_url_invalid",
	}

	for _, key := range invalidKeys {
		assert.False(t, shared.IsValidKey(key), "Key %s should be invalid", key)
	}
}

func TestConfigEditSubcommand(t *testing.T) {
	ui := cli.NewMockUi()
	c := New(ui)

	// Test with EDITOR not set and no available editors
	// This should fail gracefully
	originalEditor := os.Getenv("EDITOR")
	originalPath := os.Getenv("PATH")

	// Clear EDITOR and PATH to simulate no available editors
	os.Setenv("EDITOR", "")
	os.Setenv("PATH", "")

	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "aws-sso-config-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	configFile := filepath.Join(tmpDir, "test-config")

	// Restore environment after test
	defer func() {
		if originalEditor != "" {
			os.Setenv("EDITOR", originalEditor)
		} else {
			os.Unsetenv("EDITOR")
		}
		os.Setenv("PATH", originalPath)
	}()

	exitCode := c.Run([]string{"edit", configFile})
	assert.Equal(t, 1, exitCode)

	errorOutput := ui.ErrorWriter.String()
	assert.Contains(t, errorOutput, "No editor found")
}

func TestConfigEditHelp(t *testing.T) {
	ui := cli.NewMockUi()
	c := New(ui)

	help := c.Help()
	assert.Contains(t, help, "edit [config-file]")
	assert.Contains(t, help, "Open configuration file in an editor")
	assert.Contains(t, help, "aws-sso-config config edit")
}
