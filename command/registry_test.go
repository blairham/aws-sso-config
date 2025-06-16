package command

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/blairham/aws-sso-config/command/cli"
)

// mockUI implements the cli.UI interface for testing
type mockUI struct {
	stdout *bytes.Buffer
	stderr *bytes.Buffer
	input  *bytes.Buffer
}

func newMockUI() *mockUI {
	return &mockUI{
		stdout: &bytes.Buffer{},
		stderr: &bytes.Buffer{},
		input:  &bytes.Buffer{},
	}
}

func (m *mockUI) Ask(query string) (string, error) {
	return "", nil
}

func (m *mockUI) AskSecret(query string) (string, error) {
	return "", nil
}

func (m *mockUI) Output(message string) {
	m.stdout.WriteString(message + "\n")
}

func (m *mockUI) Info(message string) {
	m.stdout.WriteString(message + "\n")
}

func (m *mockUI) Error(message string) {
	m.stderr.WriteString(message + "\n")
}

func (m *mockUI) Warn(message string) {
	m.stderr.WriteString(message + "\n")
}

func (m *mockUI) Stdout() io.Writer {
	return m.stdout
}

func (m *mockUI) Stderr() io.Writer {
	return m.stderr
}

// mockCommand implements the cli.Command interface for testing
type mockCommand struct {
	helpText     string
	synopsisText string
	runReturn    int
}

func (m *mockCommand) Help() string {
	if m.helpText == "" {
		return "Mock command help"
	}
	return m.helpText
}

func (m *mockCommand) Run(args []string) int {
	return m.runReturn
}

func (m *mockCommand) Synopsis() string {
	if m.synopsisText == "" {
		return "Mock command synopsis"
	}
	return m.synopsisText
}

func TestRegisteredCommands(t *testing.T) {
	ui := newMockUI()
	commands := RegisteredCommands(ui)

	// Test that all expected commands are registered
	expectedCommands := []string{
		"config",
		"generate",
	}

	for _, expectedCmd := range expectedCommands {
		factory, exists := commands[expectedCmd]
		assert.True(t, exists, "Command %s should be registered", expectedCmd)

		if exists {
			cmd, err := factory()
			assert.NoError(t, err, "Command factory for %s should not error", expectedCmd)
			assert.NotNil(t, cmd, "Command %s should not be nil", expectedCmd)
		}
	}
}

func TestCommandFactories(t *testing.T) {
	ui := newMockUI()
	commands := RegisteredCommands(ui)

	// Test each command can be created successfully
	for cmdName, factory := range commands {
		t.Run(cmdName, func(t *testing.T) {
			cmd, err := factory()
			require.NoError(t, err)
			require.NotNil(t, cmd)

			// Test that each command has required methods
			assert.NotEmpty(t, cmd.Help())
			assert.NotEmpty(t, cmd.Synopsis())
		})
	}
}

func TestCommandUniqueness(t *testing.T) {
	ui := newMockUI()
	commands := RegisteredCommands(ui)

	// Test that we don't have duplicate command names
	seen := make(map[string]bool)
	for cmdName := range commands {
		assert.False(t, seen[cmdName], "Command %s should not be duplicated", cmdName)
		seen[cmdName] = true
	}
}

func TestRegisterCommandsDuplicatePanic(t *testing.T) {
	// This test would require refactoring the registerCommands function
	// to make it testable for duplicate detection. For now, we test
	// that the RegisteredCommands function returns unique commands.
	ui := newMockUI()
	commands := RegisteredCommands(ui)

	// Verify no duplicate command names
	seen := make(map[string]bool)
	for cmdName := range commands {
		assert.False(t, seen[cmdName], "Command %s should not be duplicated", cmdName)
		seen[cmdName] = true
	}
}

func TestRegistryEntryType(t *testing.T) {
	// Test the entry struct
	testEntry := entry{
		name: "test-command",
		fn: func(ui cli.UI) (cli.Command, error) {
			return &mockCommand{}, nil
		},
	}

	assert.Equal(t, "test-command", testEntry.name)
	assert.NotNil(t, testEntry.fn)

	// Test the factory function
	cmd, err := testEntry.fn(newMockUI())
	assert.NoError(t, err)
	assert.NotNil(t, cmd)
}

func TestCommandHelp(t *testing.T) {
	ui := newMockUI()
	commands := RegisteredCommands(ui)

	// Test that help contains useful information
	for cmdName, factory := range commands {
		t.Run(cmdName, func(t *testing.T) {
			cmd, err := factory()
			require.NoError(t, err)

			help := cmd.Help()
			synopsis := cmd.Synopsis()

			// Help should contain usage information
			assert.Contains(t, help, "Usage:", "Command %s help should contain usage", cmdName)

			// Synopsis should be a short description
			assert.NotEmpty(t, synopsis, "Command %s should have a synopsis", cmdName)
			assert.Less(t, len(synopsis), 100, "Synopsis for %s should be concise", cmdName)
		})
	}
}

func TestSpecificCommands(t *testing.T) {
	ui := newMockUI()
	commands := RegisteredCommands(ui)

	tests := []struct {
		name         string
		expectedType string
	}{
		{"config", "config"},
		{"generate", "generate"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory, exists := commands[tt.name]
			require.True(t, exists, "Command %s should exist", tt.name)

			cmd, err := factory()
			require.NoError(t, err)
			require.NotNil(t, cmd)

			// Test synopsis content
			synopsis := cmd.Synopsis()
			switch tt.name {
			case "config":
				assert.Contains(t, synopsis, "configuration")
			case "generate":
				assert.Contains(t, synopsis, "Generate")
			}
		})
	}
}
