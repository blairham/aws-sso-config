package get

import (
	"os"
	"testing"

	"github.com/mitchellh/cli"
	"github.com/stretchr/testify/assert"
)

func TestGetCommand(t *testing.T) {
	ui := cli.NewMockUi()
	c := New(ui)

	// Verify the command is properly initialized
	assert.NotNil(t, c)
	assert.Equal(t, "Get a configuration value", c.Synopsis())
}

func TestGetNoArgs(t *testing.T) {
	ui := cli.NewMockUi()
	c := New(ui)

	exitCode := c.Run([]string{})
	assert.Equal(t, 1, exitCode)

	errorOutput := ui.ErrorWriter.String()
	assert.Contains(t, errorOutput, "Usage: aws-sso-config config get <key>")
}

func TestGetInvalidKey(t *testing.T) {
	ui := cli.NewMockUi()
	c := New(ui)

	exitCode := c.Run([]string{"invalid_key"})
	assert.Equal(t, 1, exitCode)

	errorOutput := ui.ErrorWriter.String()
	assert.Contains(t, errorOutput, "Invalid configuration key: invalid_key")
}

func TestGetValidKey(t *testing.T) {
	// Setup temporary home directory with default config
	tempHome := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempHome)
	defer os.Setenv("HOME", originalHome)

	ui := cli.NewMockUi()
	c := New(ui)

	// Getting a valid key should work (will use defaults)
	exitCode := c.Run([]string{"sso.start_url"})
	assert.Equal(t, 0, exitCode)

	output := ui.OutputWriter.String()
	assert.Contains(t, output, "https://your-sso-portal.awsapps.com/start")
}

func TestGetHelp(t *testing.T) {
	ui := cli.NewMockUi()
	c := New(ui)

	help := c.Help()
	assert.Contains(t, help, "Usage: aws-sso-config config get <key>")
	assert.Contains(t, help, "Get a configuration value")
	assert.Contains(t, help, "sso.start_url")
	assert.Contains(t, help, "Examples:")
}
