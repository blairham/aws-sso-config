package unset

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mitchellh/cli"

	"github.com/blairham/aws-sso-config/command/config/shared"
	appconfig "github.com/blairham/aws-sso-config/providers/config"
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

	// Test with no arguments
	code := cmd.Run([]string{})
	if code == 0 {
		t.Errorf("Expected non-zero exit code with no arguments, got %d", code)
	}
	if !contains(ui.ErrorWriter.String(), "Usage: aws-sso-config config unset <key>") {
		t.Errorf("Expected usage message, got: %s", ui.ErrorWriter.String())
	}

	// Reset UI for next test
	ui = cli.NewMockUi()
	cmd = New(ui)

	// Test with invalid key
	code = cmd.Run([]string{"invalid.key"})
	if code == 0 {
		t.Errorf("Expected non-zero exit code with invalid key, got %d", code)
	}
	if !contains(ui.ErrorWriter.String(), "Invalid configuration key") {
		t.Errorf("Expected invalid key message, got: %s", ui.ErrorWriter.String())
	}

	// Reset UI for next test
	ui = cli.NewMockUi()
	cmd = New(ui)

	// Test with valid key - first create a config file with custom values
	cm := appconfig.NewConfigManager(configFile)
	config := appconfig.Default()
	config.SSO.StartURL = "https://custom.example.com/start"
	err = cm.SaveProviderConfig("sso", config.SSO)
	if err != nil {
		t.Fatalf("Failed to save test config: %v", err)
	}

	// Now test unsetting the value
	// We need to set up the config file path for the test
	// Since SaveConfigValue uses default path, we'll use a different approach
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	// Create .awsssoconfig in temp home
	defaultConfigFile := filepath.Join(tmpDir, ".awsssoconfig")
	cm2 := appconfig.NewConfigManager(defaultConfigFile)
	config2 := appconfig.Default()
	config2.SSO.StartURL = "https://custom.example.com/start"
	err = cm2.SaveProviderConfig("sso", config2.SSO)
	if err != nil {
		t.Fatalf("Failed to save test config: %v", err)
	}

	code = cmd.Run([]string{"sso.start_url"})
	if code != 0 {
		t.Errorf("Expected zero exit code with valid key, got %d. Error: %s", code, ui.ErrorWriter.String())
	}

	// Verify the value was reset to default
	output := ui.OutputWriter.String()
	defaultSSO := appconfig.DefaultSSO()
	if !contains(output, defaultSSO.StartURL) {
		t.Errorf("Expected output to contain default start URL %s, got: %s", defaultSSO.StartURL, output)
	}
}

func TestCmd_getDefaultValue(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := New(ui)

	testCases := []struct {
		key           string
		expectedValue string
		expectError   bool
	}{
		{shared.KeySSOStartURL, appconfig.DefaultSSO().StartURL, false},
		{shared.KeySSORegion, appconfig.DefaultSSO().Region, false},
		{shared.KeySSORole, appconfig.DefaultSSO().Role, false},
		{shared.KeyAWSDefaultRegion, appconfig.DefaultAWS().DefaultRegion, false},
		{shared.KeyAWSConfigFile, appconfig.DefaultAWS().ConfigFile, false},
		{"invalid.key", "", true},
	}

	for _, tc := range testCases {
		value, err := cmd.getDefaultValue(tc.key)
		if tc.expectError {
			if err == nil {
				t.Errorf("Expected error for key %s, but got none", tc.key)
			}
		} else {
			if err != nil {
				t.Errorf("Unexpected error for key %s: %v", tc.key, err)
			}
			if value != tc.expectedValue {
				t.Errorf("For key %s, expected value %s, got %s", tc.key, tc.expectedValue, value)
			}
		}
	}
}

func TestCmd_Help(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := New(ui)

	help := cmd.Help()
	if help == "" {
		t.Error("Help should not be empty")
	}

	if !contains(help, "aws-sso-config config unset") {
		t.Error("Help should contain usage information")
	}

	if !contains(help, "Reset a configuration value to its default") {
		t.Error("Help should contain description")
	}
}

func TestCmd_Synopsis(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := New(ui)

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Error("Synopsis should not be empty")
	}

	if !contains(synopsis, "Reset") {
		t.Error("Synopsis should mention reset")
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
