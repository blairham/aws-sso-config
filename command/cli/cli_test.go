package cli_test

import (
	"bytes"
	"testing"

	mcli "github.com/mitchellh/cli"
	"github.com/stretchr/testify/assert"

	"github.com/blairham/aws-sso-config/command/cli"
)

func TestBasicUI(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	ui := &cli.BasicUI{
		BasicUi: mcli.BasicUi{
			Writer:      &stdout,
			ErrorWriter: &stderr,
		},
	}

	ui.Output("Hello, world!")
	ui.Error("Oops, something went wrong.")

	assert.Equal(t, "Hello, world!\n", stdout.String())
	assert.Equal(t, "Oops, something went wrong.\n", stderr.String())
}

func TestBasicUIInfo(t *testing.T) {
	var stdout bytes.Buffer

	ui := &cli.BasicUI{
		BasicUi: mcli.BasicUi{
			Writer: &stdout,
		},
	}

	ui.Info("This is an info message")
	assert.Equal(t, "This is an info message\n", stdout.String())
}

func TestBasicUIWarn(t *testing.T) {
	var stderr bytes.Buffer

	ui := &cli.BasicUI{
		BasicUi: mcli.BasicUi{
			ErrorWriter: &stderr,
		},
	}

	ui.Warn("This is a warning message")
	assert.Equal(t, "This is a warning message\n", stderr.String())
}

func TestBasicUIMultipleMessages(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	ui := &cli.BasicUI{
		BasicUi: mcli.BasicUi{
			Writer:      &stdout,
			ErrorWriter: &stderr,
		},
	}

	ui.Output("First output")
	ui.Output("Second output")
	ui.Error("First error")
	ui.Info("Info message")
	ui.Warn("Warning message")

	expectedOutput := "First output\nSecond output\nInfo message\n"
	expectedError := "First error\nWarning message\n"

	assert.Equal(t, expectedOutput, stdout.String())
	assert.Equal(t, expectedError, stderr.String())
}

func TestBasicUIEmptyMessages(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	ui := &cli.BasicUI{
		BasicUi: mcli.BasicUi{
			Writer:      &stdout,
			ErrorWriter: &stderr,
		},
	}

	ui.Output("")
	ui.Error("")
	ui.Info("")
	ui.Warn("")

	// Each call should add a newline even for empty strings
	assert.Equal(t, "\n\n", stdout.String())
	assert.Equal(t, "\n\n", stderr.String())
}

func TestBasicUIStdout(t *testing.T) {
	var stdout bytes.Buffer

	ui := &cli.BasicUI{
		BasicUi: mcli.BasicUi{
			Writer: &stdout,
		},
	}

	// Test that Stdout() returns the correct writer
	writer := ui.Stdout()
	assert.NotNil(t, writer)

	// Write directly to the writer and verify
	writer.Write([]byte("Direct write to stdout"))
	assert.Equal(t, "Direct write to stdout", stdout.String())
}

func TestBasicUIStderr(t *testing.T) {
	var stderr bytes.Buffer

	ui := &cli.BasicUI{
		BasicUi: mcli.BasicUi{
			ErrorWriter: &stderr,
		},
	}

	// Test that Stderr() returns the correct writer
	writer := ui.Stderr()
	assert.NotNil(t, writer)

	// Write directly to the writer and verify
	writer.Write([]byte("Direct write to stderr"))
	assert.Equal(t, "Direct write to stderr", stderr.String())
}

func TestStdoutAndStderrNilSafety(t *testing.T) {
	// Test that Stdout and Stderr don't panic with nil values
	ui := &cli.BasicUI{}

	// These should not panic, but will return nil
	stdoutWriter := ui.Stdout()
	stderrWriter := ui.Stderr()

	assert.Nil(t, stdoutWriter)
	assert.Nil(t, stderrWriter)
}

func TestBasicUIInterface(t *testing.T) {
	// Verify that BasicUI implements the UI interface
	var _ cli.UI = &cli.BasicUI{}
}
