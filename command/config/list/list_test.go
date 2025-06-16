package list

import (
	"bytes"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestListCommand(t *testing.T) {
	var stdout, stderr bytes.Buffer
	ui := &cli.BasicUi{
		Writer:      &stdout,
		ErrorWriter: &stderr,
	}

	cmd := New(ui)
	if cmd == nil {
		t.Fatal("New() returned nil")
	}
}

func TestListNoArgs(t *testing.T) {
	var stdout, stderr bytes.Buffer
	ui := &cli.BasicUi{
		Writer:      &stdout,
		ErrorWriter: &stderr,
	}

	cmd := New(ui)
	code := cmd.Run([]string{})

	if code != 0 {
		t.Errorf("Expected exit code 0, got %d", code)
	}

	// The output should be in key=value format, even if values are empty
	// We can't reliably test specific values since they depend on the user's config
	// But we can test that the output format is correct
	output := stdout.String()

	// If there's any output, it should be in key=value format
	lines := strings.Split(strings.TrimSpace(output), "\n")
	for _, line := range lines {
		if line != "" && !strings.Contains(line, "=") {
			t.Errorf("Expected line to be in key=value format, got %s", line)
		}
	}
}

func TestListWithArgs(t *testing.T) {
	var stdout, stderr bytes.Buffer
	ui := &cli.BasicUi{
		Writer:      &stdout,
		ErrorWriter: &stderr,
	}

	cmd := New(ui)
	code := cmd.Run([]string{"extra", "args"})

	if code != 1 {
		t.Errorf("Expected exit code 1, got %d", code)
	}

	errorOutput := stderr.String()
	if !strings.Contains(errorOutput, "Usage: aws-sso-config config list") {
		t.Errorf("Expected error to contain usage message, got %s", errorOutput)
	}

	if !strings.Contains(errorOutput, "This command takes no arguments except optional flags.") {
		t.Errorf("Expected error to contain no arguments message, got %s", errorOutput)
	}
}

func TestListHelp(t *testing.T) {
	var stdout, stderr bytes.Buffer
	ui := &cli.BasicUi{
		Writer:      &stdout,
		ErrorWriter: &stderr,
	}

	cmd := New(ui)
	help := cmd.Help()

	if !strings.Contains(help, "Usage: aws-sso-config config list") {
		t.Errorf("Expected help to contain usage, got %s", help)
	}

	if !strings.Contains(help, "List all configuration variables") {
		t.Errorf("Expected help to contain description, got %s", help)
	}
}

func TestListSynopsis(t *testing.T) {
	var stdout, stderr bytes.Buffer
	ui := &cli.BasicUi{
		Writer:      &stdout,
		ErrorWriter: &stderr,
	}

	cmd := New(ui)
	synopsis := cmd.Synopsis()

	expected := "List all configuration variables and their values"
	if synopsis != expected {
		t.Errorf("Expected synopsis %q, got %q", expected, synopsis)
	}
}
