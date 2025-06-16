package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/blairham/aws-sso-config/command/cli"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected int
	}{
		{
			name:     "no arguments",
			args:     []string{},
			expected: 127, // CLI returns 127 when no subcommand is provided
		},
		{
			name:     "help flag",
			args:     []string{"--help"},
			expected: 0,
		},
		{
			name:     "version flag",
			args:     []string{"--version"},
			expected: 0,
		},
		{
			name:     "invalid command",
			args:     []string{"invalid-command"},
			expected: 127, // CLI returns 127 for invalid commands
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Run(tt.args)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCreateCLI(t *testing.T) {
	// Create a basic UI for testing
	ui := &cli.BasicUI{}

	// Test CLI creation
	cliInstance := createCLI(ui, []string{})
	assert.NotNil(t, cliInstance)

	// Test that commands are registered
	commands := cliInstance.Commands
	assert.NotEmpty(t, commands)

	// Check for expected commands
	expectedCommands := []string{"config", "generate"}
	for _, expectedCmd := range expectedCommands {
		_, exists := commands[expectedCmd]
		assert.True(t, exists, "Expected command %s to be registered", expectedCmd)
	}
}

func TestVersionInfo(t *testing.T) {
	// Test that version variables exist (they're set by build flags)
	assert.NotEmpty(t, version)
	assert.NotEmpty(t, commit)
	assert.NotEmpty(t, buildTime)
}

func TestMainFunction(t *testing.T) {
	// Test that main function doesn't panic
	// We can't easily test os.Exit, but we can test that main runs without panic

	// Save original osExit
	originalOsExit := osExit
	defer func() { osExit = originalOsExit }()

	var exitCode int
	osExit = func(code int) { exitCode = code }

	// Save original args
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	// Test with help args
	os.Args = []string{"aws-sso-config", "--help"}

	// This should not panic
	assert.NotPanics(t, func() {
		main()
	})

	// Should exit with 0 for help
	assert.Equal(t, 0, exitCode)
}
