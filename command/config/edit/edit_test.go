package edit

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mitchellh/cli"
)

func TestCmd_Run(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "aws-sso-config-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	configFile := filepath.Join(tmpDir, "test-config")

	ui := cli.NewMockUi()
	cmd := New(ui)

	// Test with EDITOR not set and no available editors
	// This should fail gracefully
	originalEditor := os.Getenv("EDITOR")
	originalPath := os.Getenv("PATH")

	// Clear EDITOR and PATH to simulate no available editors
	os.Setenv("EDITOR", "")
	os.Setenv("PATH", "")

	// Restore environment after test
	defer func() {
		if originalEditor != "" {
			os.Setenv("EDITOR", originalEditor)
		} else {
			os.Unsetenv("EDITOR")
		}
		os.Setenv("PATH", originalPath)
	}()

	args := []string{configFile}
	code := cmd.Run(args)

	if code == 0 {
		t.Errorf("Expected non-zero exit code when no editor is available, got %d", code)
	}

	if !contains(ui.ErrorWriter.String(), "No editor found") {
		t.Errorf("Expected error message about no editor found, got: %s", ui.ErrorWriter.String())
	}
}

func TestCmd_ensureConfigFileExists(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "aws-sso-config-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	configFile := filepath.Join(tmpDir, "test-config")

	ui := cli.NewMockUi()
	cmd := New(ui)

	// Test creating config file
	err = cmd.ensureConfigFileExists(configFile)
	if err != nil {
		t.Errorf("ensureConfigFileExists failed: %v", err)
	}

	// Check if file was created
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		t.Errorf("Config file was not created")
	}

	// Test with existing file
	err = cmd.ensureConfigFileExists(configFile)
	if err != nil {
		t.Errorf("ensureConfigFileExists failed with existing file: %v", err)
	}
}

func TestCmd_Help(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := New(ui)

	help := cmd.Help()
	if help == "" {
		t.Error("Help should not be empty")
	}

	if !contains(help, "aws-sso-config config edit") {
		t.Error("Help should contain usage information")
	}

	if !contains(help, "EDITOR") {
		t.Error("Help should mention EDITOR environment variable")
	}
}

func TestCmd_Synopsis(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := New(ui)

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Error("Synopsis should not be empty")
	}

	if !contains(synopsis, "editor") {
		t.Error("Synopsis should mention editor")
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && (s[:len(substr)] == substr ||
			s[len(s)-len(substr):] == substr ||
			containsSubstring(s, substr))))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
