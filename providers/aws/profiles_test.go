package aws

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetProfileFromRepoName(t *testing.T) {
	tests := []struct {
		name     string
		repoName string
		expected string
	}{
		{
			name:     "simple repo name",
			repoName: "my-service",
			expected: "my-service",
		},
		{
			name:     "complex repo name",
			repoName: "frontend-dashboard",
			expected: "frontend-dashboard",
		},
		{
			name:     "hyphenated repo name",
			repoName: "user-authentication-service",
			expected: "user-authentication-service",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getProfileFromRepoName(tt.repoName)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConfigFile(t *testing.T) {
	configFile, err := ConfigFile()
	assert.NoError(t, err)

	// Should not be empty
	assert.NotEmpty(t, configFile)

	// Should contain .aws/config
	assert.Contains(t, configFile, ".aws/config")
}

func TestToString(t *testing.T) {
	// Test with non-nil pointer
	str := "test"
	ptr := &str
	result := ToString(ptr)
	assert.Equal(t, "test", result)

	// Test with nil pointer
	var nilPtr *string
	result = ToString(nilPtr)
	assert.Equal(t, "", result)
}

func TestGetProfileWithEnvironmentSet(t *testing.T) {
	// Set AWS_PROFILE environment variable
	os.Setenv(AwsProfile, "test-profile")
	defer os.Unsetenv(AwsProfile)

	profile, err := GetProfile()

	// Should not error and should return empty string (skipped)
	assert.NoError(t, err)
	assert.Equal(t, "", profile)
}

func TestGetProfileNotInGitRepo(t *testing.T) {
	// Save current directory
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)

	// Change to a temporary directory that's not a git repo
	tempDir := t.TempDir()
	os.Chdir(tempDir)

	profile, _ := GetProfile()
	assert.Equal(t, "default", profile)
}

func TestConfigFileFunction(t *testing.T) {
	configFile, err := ConfigFile()
	assert.NoError(t, err)
	assert.Contains(t, configFile, ".aws/config")
	assert.NotEmpty(t, configFile)
}

func TestToStringFunction(t *testing.T) {
	tests := []struct {
		name     string
		input    *string
		expected string
	}{
		{
			name:     "non-nil string",
			input:    stringPointer("test-value"),
			expected: "test-value",
		},
		{
			name:     "empty string",
			input:    stringPointer(""),
			expected: "",
		},
		{
			name:     "nil pointer",
			input:    nil,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToString(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetProfileFromRepoNameEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		repoName string
		expected string
	}{
		{
			name:     "empty repo name",
			repoName: "",
			expected: "",
		},
		{
			name:     "single character",
			repoName: "a",
			expected: "a",
		},
		{
			name:     "numbers and letters",
			repoName: "service123",
			expected: "service123",
		},
		{
			name:     "underscores",
			repoName: "my_service_app",
			expected: "my_service_app",
		},
		{
			name:     "mixed separators",
			repoName: "my-service_app",
			expected: "my-service_app",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getProfileFromRepoName(tt.repoName)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetProfileWithInvalidGitRepo(t *testing.T) {
	// Save current directory and environment
	originalWd, _ := os.Getwd()
	originalProfile := os.Getenv("AWS_PROFILE")
	defer func() {
		os.Chdir(originalWd)
		if originalProfile == "" {
			os.Unsetenv("AWS_PROFILE")
		} else {
			os.Setenv("AWS_PROFILE", originalProfile)
		}
	}()

	// Clear AWS_PROFILE environment variable
	os.Unsetenv("AWS_PROFILE")

	// Create a temporary directory with .git but no remote
	tempDir := t.TempDir()
	os.Chdir(tempDir)

	// Create a .git directory but no remote config
	gitDir := tempDir + "/.git"
	os.Mkdir(gitDir, 0750)

	profile, _ := GetProfile()
	assert.Equal(t, "default", profile)
}

func TestLoadDefaultConfigFunction(t *testing.T) {
	// This is mainly testing that the function doesn't panic
	// In a real environment this would load AWS config
	cfg := LoadDefaultConfig()
	assert.NotNil(t, cfg)
}

// TestValidateAccountID tests the validateAccountID function with various scenarios
func TestValidateAccountID(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()

	// Test with valid terragrunt.hcl file and matching account ID
	validTerragruntContent := `
inputs = {
  account_id = "123456789012"
  region     = "eu-west-1"
}
`
	terragruntPath := filepath.Join(tempDir, "terragrunt.hcl")
	err := os.WriteFile(terragruntPath, []byte(validTerragruntContent), 0600)
	require.NoError(t, err)

	// Should succeed when account ID matches
	err = validateAccountID("123456789012", tempDir)
	assert.NoError(t, err)

	// Should fail when account ID doesn't match
	err = validateAccountID("987654321098", tempDir)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "did not match entry")

	// Test with missing terragrunt file
	invalidDir := filepath.Join(tempDir, "non-existent")
	os.MkdirAll(invalidDir, 0750)
	err = validateAccountID("123456789012", invalidDir)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "could not find terragrunt.hcl")

	// Test with terragrunt file without account_id
	invalidTerragruntContent := `
inputs = {
  region = "eu-west-1"
}
`
	invalidTerragruntPath := filepath.Join(invalidDir, "terragrunt.hcl")
	err = os.WriteFile(invalidTerragruntPath, []byte(invalidTerragruntContent), 0600)
	require.NoError(t, err)

	err = validateAccountID("123456789012", invalidDir)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "could not determine account id")
}

// TestValidateProfile tests the validateProfile function
func TestValidateProfile(t *testing.T) {
	t.Skip("Skipping validateProfile test - it uses hardcoded paths that are difficult to mock")

	// This test is challenging because validateProfile() calls ConfigFile()
	// which constructs a hardcoded path to ~/.aws/config
	// To properly test this, we would need to refactor the code to allow
	// dependency injection of the config file path
}

// Helper function for tests
func stringPointer(s string) *string {
	return &s
}
