package set

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mitchellh/cli"
	"github.com/stretchr/testify/assert"
)

func TestSetCommand(t *testing.T) {
	ui := cli.NewMockUi()
	c := New(ui)

	// Verify the command is properly initialized
	assert.NotNil(t, c)
	assert.Equal(t, "Set a configuration value", c.Synopsis())
}

func TestSetNoArgs(t *testing.T) {
	ui := cli.NewMockUi()
	c := New(ui)

	exitCode := c.Run([]string{})
	assert.Equal(t, 1, exitCode)

	errorOutput := ui.ErrorWriter.String()
	assert.Contains(t, errorOutput, "Usage: aws-sso-config config set <key> <value>")
}

func TestSetInsufficientArgs(t *testing.T) {
	ui := cli.NewMockUi()
	c := New(ui)

	exitCode := c.Run([]string{"sso_start_url"})
	assert.Equal(t, 1, exitCode)

	errorOutput := ui.ErrorWriter.String()
	assert.Contains(t, errorOutput, "Usage: aws-sso-config config set <key> <value>")
}

func TestSetInvalidKey(t *testing.T) {
	ui := cli.NewMockUi()
	c := New(ui)

	exitCode := c.Run([]string{"invalid_key", "value"})
	assert.Equal(t, 1, exitCode)

	errorOutput := ui.ErrorWriter.String()
	assert.Contains(t, errorOutput, "Invalid configuration key: invalid_key")
}

func TestSetValidKeyValue(t *testing.T) {
	// Setup temporary home directory
	tempHome := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempHome)
	defer os.Setenv("HOME", originalHome)

	configPath := filepath.Join(tempHome, ".awsssoconfig")

	ui := cli.NewMockUi()
	c := New(ui)

	exitCode := c.Run([]string{"sso.start_url", "https://test.awsapps.com/start"})
	assert.Equal(t, 0, exitCode)

	output := ui.OutputWriter.String()
	assert.Contains(t, output, "Updated sso.start_url = https://test.awsapps.com/start")

	// Verify file was created
	_, err := os.Stat(configPath)
	assert.NoError(t, err)
}

func TestSetHelp(t *testing.T) {
	ui := cli.NewMockUi()
	c := New(ui)

	help := c.Help()
	assert.Contains(t, help, "Usage: aws-sso-config config set <key> <value>")
	assert.Contains(t, help, "Set a configuration value")
	assert.Contains(t, help, "sso.start_url")
	assert.Contains(t, help, "Examples:")
}

func TestSetMultipleWordValue(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "aws-sso-config-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Set up HOME environment variable to use temp directory
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	ui := cli.NewMockUi()
	c := New(ui)

	// Test setting a value with multiple words (without quotes)
	exitCode := c.Run([]string{"sso.role", "Administrator", "Access", "Role"})
	assert.Equal(t, 0, exitCode)

	output := ui.OutputWriter.String()
	assert.Contains(t, output, "Updated sso.role = Administrator Access Role")
}

func TestSetSingleWordValue(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "aws-sso-config-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Set up HOME environment variable to use temp directory
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	ui := cli.NewMockUi()
	c := New(ui)

	// Test setting a value with a single word
	exitCode := c.Run([]string{"aws.default_region", "us-west-2"})
	assert.Equal(t, 0, exitCode)

	output := ui.OutputWriter.String()
	assert.Contains(t, output, "Updated aws.default_region = us-west-2")
}
