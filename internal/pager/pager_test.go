package pager

import (
	"bytes"
	"os"
	"testing"

	"github.com/mitchellh/cli"
)

func TestNew(t *testing.T) {
	ui := &cli.BasicUi{
		Writer:      &bytes.Buffer{},
		ErrorWriter: &bytes.Buffer{},
	}

	pager := New(ui)
	if pager == nil {
		t.Fatal("Expected pager to be created")
	}

	if pager.ui != ui {
		t.Error("Expected pager to have the provided UI")
	}

	if pager.threshold <= 0 {
		t.Error("Expected pager to have a positive threshold")
	}
}

func TestNewWithThreshold(t *testing.T) {
	ui := &cli.BasicUi{
		Writer:      &bytes.Buffer{},
		ErrorWriter: &bytes.Buffer{},
	}

	threshold := 10
	pager := NewWithThreshold(ui, threshold)
	if pager.threshold != threshold {
		t.Errorf("Expected threshold %d, got %d", threshold, pager.threshold)
	}
}

func TestSetForceEnabled(t *testing.T) {
	ui := &cli.BasicUi{
		Writer:      &bytes.Buffer{},
		ErrorWriter: &bytes.Buffer{},
	}

	pager := New(ui)
	pager.SetForceEnabled(true)
	if !pager.forceEnabled {
		t.Error("Expected forceEnabled to be true")
	}

	pager.SetForceEnabled(false)
	if pager.forceEnabled {
		t.Error("Expected forceEnabled to be false")
	}
}

func TestShouldPage(t *testing.T) {
	ui := &cli.BasicUi{
		Writer:      &bytes.Buffer{},
		ErrorWriter: &bytes.Buffer{},
	}

	pager := NewWithThreshold(ui, 5)

	// Test with small output
	smallLines := []string{"line1", "line2", "line3"}
	if pager.shouldPage(smallLines) {
		t.Error("Expected small output not to be paged")
	}

	// Test with large output
	largeLines := make([]string, 10)
	for i := range largeLines {
		largeLines[i] = "line"
	}

	// This test depends on terminal detection, so we'll test force enabled
	pager.SetForceEnabled(true)
	if !pager.shouldPage(smallLines) {
		t.Error("Expected forced paging to work with small output")
	}

	// Test NO_PAGER environment variable
	os.Setenv("NO_PAGER", "1")
	defer os.Unsetenv("NO_PAGER")

	pager.SetForceEnabled(false)
	if pager.shouldPage(largeLines) {
		t.Error("Expected NO_PAGER to disable paging")
	}
}

func TestOutput(t *testing.T) {
	var stdout bytes.Buffer
	ui := &cli.BasicUi{
		Writer:      &stdout,
		ErrorWriter: &bytes.Buffer{},
	}

	pager := NewWithThreshold(ui, 10) // High threshold to avoid paging
	lines := []string{"line1", "line2", "line3"}

	pager.Output(lines)

	output := stdout.String()
	for _, line := range lines {
		if !contains(output, line) {
			t.Errorf("Expected output to contain %s", line)
		}
	}
}

func TestGetPagerCommand(t *testing.T) {
	// Clean environment
	os.Unsetenv("NO_PAGER")
	os.Unsetenv("PAGER")
	os.Unsetenv("AWS_SSO_CONFIG_PAGER")

	// Test default behavior
	cmd := getPagerCommand()
	if cmd == "" {
		t.Error("Expected to find a pager command")
	}

	// Test custom pager via environment (should be rejected if not in allowlist)
	os.Setenv("PAGER", "custom-pager")
	defer os.Unsetenv("PAGER")

	cmd = getPagerCommand()
	// custom-pager is not in allowlist, so should fall back to default
	if cmd == "custom-pager" {
		t.Errorf("Expected custom-pager to be rejected, but got %s", cmd)
	}

	// Test with a valid pager
	os.Setenv("PAGER", "less")
	cmd = getPagerCommand()
	if cmd != "less" {
		t.Errorf("Expected less, got %s", cmd)
	}

	// Test AWS_SSO_CONFIG_PAGER with valid pager
	os.Setenv("AWS_SSO_CONFIG_PAGER", "more")
	defer os.Unsetenv("AWS_SSO_CONFIG_PAGER")

	cmd = getPagerCommand()
	if cmd != "less" { // PAGER should take precedence
		t.Errorf("Expected less to take precedence, got %s", cmd)
	}

	os.Unsetenv("PAGER")
	cmd = getPagerCommand()
	if cmd != "more" {
		t.Errorf("Expected more, got %s", cmd)
	}
}

func TestGetTerminalHeight(t *testing.T) {
	height := getTerminalHeight()
	if height <= 0 {
		t.Error("Expected positive terminal height")
	}

	// Test with LINES environment variable
	os.Setenv("LINES", "50")
	defer os.Unsetenv("LINES")

	height = getTerminalHeight()
	if height != 50 {
		t.Errorf("Expected height 50, got %d", height)
	}

	// Test with invalid LINES
	os.Setenv("LINES", "invalid")
	height = getTerminalHeight()
	if height <= 0 {
		t.Error("Expected fallback to work with invalid LINES")
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && (s[:len(substr)] == substr ||
			s[len(s)-len(substr):] == substr ||
			len(s) > len(substr) && someContains(s, substr))))
}

func someContains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
