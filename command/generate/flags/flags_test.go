package flags

import (
	"testing"
)

func TestNewDiffFlag(t *testing.T) {
	flag := NewDiffFlag()

	if flag == nil {
		t.Fatal("NewDiffFlag() returned nil")
	}

	if flag.GetFlagName() != "diff" {
		t.Errorf("GetFlagName() = %q, expected %q", flag.GetFlagName(), "diff")
	}

	if flag.GetShortFlag() != "d" {
		t.Errorf("GetShortFlag() = %q, expected %q", flag.GetShortFlag(), "d")
	}

	if flag.GetDescription() == "" {
		t.Error("GetDescription() returned empty string")
	}

	if flag.GetUsage() == "" {
		t.Error("GetUsage() returned empty string")
	}
}

func TestNewConfigFlag(t *testing.T) {
	flag := NewConfigFlag()

	if flag == nil {
		t.Fatal("NewConfigFlag() returned nil")
	}

	if flag.GetFlagName() != "config" {
		t.Errorf("GetFlagName() = %q, expected %q", flag.GetFlagName(), "config")
	}

	if flag.GetShortFlag() != "c" {
		t.Errorf("GetShortFlag() = %q, expected %q", flag.GetShortFlag(), "c")
	}

	if flag.GetDescription() == "" {
		t.Error("GetDescription() returned empty string")
	}

	if flag.GetUsage() == "" {
		t.Error("GetUsage() returned empty string")
	}
}

func TestFlagRegistry_NewFlagRegistry(t *testing.T) {
	registry := NewFlagRegistry()

	if registry == nil {
		t.Fatal("NewFlagRegistry() returned nil")
	}

	flags := registry.GetAllFlags()
	if len(flags) == 0 {
		t.Error("GetAllFlags() returned no flags")
	}

	// Should have at least diff and config flags
	expectedFlags := []string{"diff", "config"}
	for _, expectedFlag := range expectedFlags {
		found := false
		for _, flag := range flags {
			if flag.GetFlagName() == expectedFlag {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected flag %q not found in registry", expectedFlag)
		}
	}
}

func TestFlagRegistry_GetFlagByName(t *testing.T) {
	registry := NewFlagRegistry()

	tests := []struct {
		name        string
		expectedNil bool
	}{
		{"diff", false},
		{"config", false},
		{"nonexistent", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag := registry.GetFlagByName(tt.name)
			if tt.expectedNil && flag != nil {
				t.Errorf("GetFlagByName(%q) returned %v, expected nil", tt.name, flag)
			}
			if !tt.expectedNil && flag == nil {
				t.Errorf("GetFlagByName(%q) returned nil, expected a flag", tt.name)
			}
			if !tt.expectedNil && flag != nil && flag.GetFlagName() != tt.name {
				t.Errorf("GetFlagByName(%q) returned flag with name %q", tt.name, flag.GetFlagName())
			}
		})
	}
}

func TestDiffFlag_Interface(t *testing.T) {
	flag := NewDiffFlag()

	// Test that DiffFlag implements the Flag interface
	var _ Flag = flag

	if flag.GetDescription() == "" {
		t.Error("DiffFlag GetDescription() returned empty string")
	}

	expectedSubstring := "diff"
	if !containsIgnoreCase(flag.GetDescription(), expectedSubstring) {
		t.Errorf("DiffFlag GetDescription() = %q, expected to contain %q", flag.GetDescription(), expectedSubstring)
	}
}

func TestConfigFlag_Interface(t *testing.T) {
	flag := NewConfigFlag()

	// Test that ConfigFlag implements the Flag interface
	var _ Flag = flag

	if flag.GetDescription() == "" {
		t.Error("ConfigFlag GetDescription() returned empty string")
	}

	expectedSubstring := "config"
	if !containsIgnoreCase(flag.GetDescription(), expectedSubstring) {
		t.Errorf("ConfigFlag GetDescription() = %q, expected to contain %q", flag.GetDescription(), expectedSubstring)
	}
}

// Helper function to check if a string contains a substring (case-insensitive)
func containsIgnoreCase(s, substr string) bool {
	s = toLower(s)
	substr = toLower(substr)
	return contains(s, substr)
}

func toLower(s string) string {
	result := make([]byte, len(s))
	for i, c := range []byte(s) {
		if c >= 'A' && c <= 'Z' {
			result[i] = c + 32
		} else {
			result[i] = c
		}
	}
	return string(result)
}

func contains(s, substr string) bool {
	if len(substr) > len(s) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
