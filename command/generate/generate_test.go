package generate

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sso"
	"github.com/aws/aws-sdk-go-v2/service/sso/types"
	"github.com/mitchellh/cli"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	appconfig "github.com/blairham/aws-sso-config/providers/config"
)

func TestInit(t *testing.T) {
	ui := cli.NewMockUi()
	c := New(ui)

	// Verify the command is properly initialized
	assert.NotNil(t, c.flags)
	assert.NotEmpty(t, c.help)
	assert.Equal(t, synopsis, c.Synopsis())
}

func TestGenerateWithConfigFile(t *testing.T) {
	ui := cli.NewMockUi()

	// Create a temporary directory and files
	tmpDir := t.TempDir()

	// Create AWS config file
	awsConfigFile := filepath.Join(tmpDir, "aws-config")
	err := os.WriteFile(awsConfigFile, []byte("[default]\nregion = us-east-1\n"), 0600)
	require.NoError(t, err)

	// Create app config file
	appConfigFile := filepath.Join(tmpDir, "app-config.toml")
	configContent := `[sso]
start_url = "https://test.awsapps.com/start"
region = "us-west-2"
role = "TestRole"

[aws]
default_region = "eu-west-1"
config_file = "` + awsConfigFile + `"`
	err = os.WriteFile(appConfigFile, []byte(configContent), 0600)
	require.NoError(t, err)

	// Create mock dependencies that simulate successful AWS operations
	mockSSOClient := &MockSSOClient{}
	mockSSOClient.On("ListAccounts", mock.Anything, mock.Anything).Return(
		&sso.ListAccountsOutput{
			AccountList: []types.AccountInfo{
				{
					AccountId:   aws.String("123456789012"),
					AccountName: aws.String("Production Account"),
				},
				{
					AccountId:   aws.String("987654321098"),
					AccountName: aws.String("Development Account"),
				},
			},
		}, nil)

	mockSSOClientFactory := func(cfg aws.Config) SSOClient {
		return mockSSOClient
	}

	token := "mock-access-token"
	mockTokenGenerator := &MockTokenGenerator{
		token: &token,
	}

	mockConfigLoader := func() aws.Config {
		return aws.Config{}
	}

	c := NewWithDependencies(ui, mockSSOClientFactory, mockTokenGenerator, mockConfigLoader)

	// Test successful generation
	exitCode := c.Run([]string{"--config=" + appConfigFile})
	assert.Equal(t, 0, exitCode, "Should successfully generate config")

	// Verify the mock was called
	mockSSOClient.AssertExpectations(t)

	// Verify that the AWS config file was updated
	updatedContent, err := os.ReadFile(awsConfigFile)
	require.NoError(t, err)

	content := string(updatedContent)
	assert.Contains(t, content, "Production Account", "Should contain production account profile")
	assert.Contains(t, content, "Development Account", "Should contain development account profile")
	assert.Contains(t, content, "123456789012", "Should contain production account ID")
	assert.Contains(t, content, "987654321098", "Should contain development account ID")
}

func TestGenerateFlagParsing(t *testing.T) {
	ui := cli.NewMockUi()
	c := New(ui)

	// Test just the flag parsing
	tempDir := t.TempDir()
	appConfigFile := filepath.Join(tempDir, "test-config.yaml")

	// Create a minimal config file that passes validation
	configContent := `sso_start_url: "https://test.awsapps.com/start"
sso_region: "us-west-2"
sso_role: "TestRole"
default_region: "eu-west-1"
config_file: "/tmp/test-aws-config"
`
	err := os.WriteFile(appConfigFile, []byte(configContent), 0600)
	require.NoError(t, err)

	// Just test flag parsing - no need to execute the command fully
	err = c.flags.Parse([]string{"--config=" + appConfigFile, "--diff"})
	require.NoError(t, err)

	// Check that flags were parsed correctly
	assert.Equal(t, appConfigFile, c.configFile)
	assert.True(t, c.diff)
}

func TestGenerateWithoutConfigFile(t *testing.T) {
	t.Skip("Skipping test that requires AWS SSO integration")
}

func TestGenerateInvalidConfigFile(t *testing.T) {
	t.Skip("Skipping test that requires AWS SSO integration")
}

func TestGenerateMalformedConfigFile(t *testing.T) {
	t.Skip("Skipping test that requires AWS SSO integration")
}

func TestGenerateInvalidFlags(t *testing.T) {
	t.Skip("Skipping test that requires AWS SSO integration")
}

func TestGenerateHelpOutput(t *testing.T) {
	ui := cli.NewMockUi()
	c := New(ui)

	help := c.Help()
	assert.Contains(t, help, "Usage: aws-sso-config generate")
	assert.Contains(t, help, "-diff")
	assert.Contains(t, help, "--config string")
	assert.Contains(t, help, "Examples:")
	assert.Contains(t, help, "Enable diff output")
	assert.Contains(t, help, "Path to configuration file")
}

func TestGenerateSynopsis(t *testing.T) {
	ui := cli.NewMockUi()
	c := New(ui)

	assert.Equal(t, synopsis, c.Synopsis())
}

func TestGenerateConfigValidation(t *testing.T) {
	t.Skip("Skipping test that requires AWS SSO integration")
}

func TestGenerateFlagParsingTable(t *testing.T) {
	t.Skip("Skipping test that requires AWS SSO integration")

	tests := []struct {
		name     string
		args     []string
		wantDiff bool
		wantFile string
	}{
		{
			name:     "no flags",
			args:     []string{},
			wantDiff: false,
			wantFile: "",
		},
		{
			name:     "diff flag only",
			args:     []string{"--diff"},
			wantDiff: true,
			wantFile: "",
		},
		{
			name:     "config flag only",
			args:     []string{"--config=test.yaml"},
			wantDiff: false,
			wantFile: "test.yaml",
		},
		{
			name:     "both flags",
			args:     []string{"--diff", "--config=my-config.yaml"},
			wantDiff: true,
			wantFile: "my-config.yaml",
		},
		{
			name:     "flags in different order",
			args:     []string{"--config=another.yaml", "--diff"},
			wantDiff: true,
			wantFile: "another.yaml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ui := cli.NewMockUi()
			c := New(ui)

			// We expect these to fail due to authentication, but flags should parse
			c.Run(tt.args)

			assert.Equal(t, tt.wantDiff, c.diff)
			assert.Equal(t, tt.wantFile, c.configFile)
		})
	}
}

// TestGenerateAwsConfigFile tests the generateAwsConfigFile function
func TestGenerateAwsConfigFile(t *testing.T) {
	t.Skip("Skipping test that requires mocking AWS SSO pagination")
}

// TestRunError tests the Run function error handling with mocks
func TestRunError(t *testing.T) {
	ui := cli.NewMockUi()

	// Create mock dependencies
	mockSSOClientFactory := func(cfg aws.Config) SSOClient {
		return &MockSSOClient{}
	}

	mockTokenGenerator := &MockTokenGenerator{
		shouldFail: true, // This will cause the token generation to fail
	}

	mockConfigLoader := func() aws.Config {
		return aws.Config{}
	}

	c := NewWithDependencies(ui, mockSSOClientFactory, mockTokenGenerator, mockConfigLoader)

	// Create a temporary directory and files
	tmpDir := t.TempDir()

	// Create an invalid TOML config file to trigger the error
	invalidConfigFile := filepath.Join(tmpDir, "invalid.toml")
	err := os.WriteFile(invalidConfigFile, []byte("invalid toml content: ["), 0600)
	require.NoError(t, err)

	// Test with invalid config file - this should fail during config loading
	exitCode := c.Run([]string{"--config=" + invalidConfigFile})
	assert.NotEqual(t, 0, exitCode, "Should return non-zero exit code for invalid config file")
}

// TestRunConfigFileError tests Run with config file error
func TestRunConfigFileError(t *testing.T) {
	ui := cli.NewMockUi()
	c := New(ui)

	// Create a temporary invalid YAML file
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "invalid.yaml")
	err := os.WriteFile(configFile, []byte("invalid yaml content: ["), 0600)
	require.NoError(t, err)

	// Test with invalid YAML config
	exitCode := c.Run([]string{"--config=" + configFile})
	assert.NotEqual(t, 0, exitCode, "Should return non-zero exit code for invalid YAML")
}

// TestRunDefaultConfig tests Run with default configuration
func TestRunDefaultConfig(t *testing.T) {
	t.Skip("Skipping test that requires AWS SSO authentication")

	ui := cli.NewMockUi()
	c := New(ui)

	// This test will likely fail due to authentication requirements, but it exercises the default config path
	exitCode := c.Run([]string{})
	// We expect a non-zero exit code due to authentication failure
	assert.NotEqual(t, 0, exitCode, "Should return non-zero exit code due to authentication failure")
}

// TestRunParseError tests Run with parse error
func TestRunParseError(t *testing.T) {
	ui := cli.NewMockUi()
	c := New(ui)

	// Test with invalid flag
	exitCode := c.Run([]string{"-invalid-flag"})
	assert.Equal(t, 1, exitCode, "Should return exit code 1 for flag parse error")
}

// TestRunWithMocks tests the Run function with successful mocks
func TestRunWithMocks(t *testing.T) {
	ui := cli.NewMockUi()

	// Create a temporary directory and files
	tmpDir := t.TempDir()

	// Create AWS config file
	awsConfigFile := filepath.Join(tmpDir, "aws-config")
	err := os.WriteFile(awsConfigFile, []byte("[default]\nregion = us-east-1\n"), 0600)
	require.NoError(t, err)

	// Create app config file
	appConfigFile := filepath.Join(tmpDir, "app-config.toml")
	configContent := `[sso]
start_url = "https://test.awsapps.com/start"
region = "us-west-2"
role = "TestRole"

[aws]
default_region = "eu-west-1"
config_file = "` + awsConfigFile + `"`
	err = os.WriteFile(appConfigFile, []byte(configContent), 0600)
	require.NoError(t, err)

	// Create mock dependencies
	mockSSOClient := &MockSSOClient{}
	mockSSOClient.On("ListAccounts", mock.Anything, mock.Anything).Return(
		&sso.ListAccountsOutput{
			AccountList: []types.AccountInfo{
				{
					AccountId:   aws.String("123456789012"),
					AccountName: aws.String("Test Account"),
				},
			},
		}, nil)

	mockSSOClientFactory := func(cfg aws.Config) SSOClient {
		return mockSSOClient
	}

	token := "mock-access-token"
	mockTokenGenerator := &MockTokenGenerator{
		token: &token,
	}

	mockConfigLoader := func() aws.Config {
		return aws.Config{}
	}

	c := NewWithDependencies(ui, mockSSOClientFactory, mockTokenGenerator, mockConfigLoader)

	// Test with valid config
	exitCode := c.Run([]string{"--config=" + appConfigFile})
	assert.Equal(t, 0, exitCode, "Should return zero exit code for successful run")

	// Verify the mock was called
	mockSSOClient.AssertExpectations(t)
}

// MockTokenGenerator implements TokenGenerator for testing
type MockTokenGenerator struct {
	shouldFail bool
	token      *string
}

func (m *MockTokenGenerator) GenerateTokenWithConfig(cfg aws.Config, appCfg *appconfig.Config) *string {
	if m.shouldFail {
		return nil
	}
	if m.token != nil {
		return m.token
	}
	token := "mock-token"
	return &token
}

// TestFlagFormatDocumentation documents correct and incorrect flag usage patterns
func TestFlagFormatDocumentation(t *testing.T) {
	ui := cli.NewMockUi()

	tests := []struct {
		name           string
		args           []string
		expectError    bool
		expectedConfig string
		description    string
	}{
		{
			name:           "correct long form config flag",
			args:           []string{"--config=/tmp/test"},
			expectError:    false,
			expectedConfig: "/tmp/test",
			description:    "This is the recommended way to use the config flag with long form",
		},
		{
			name:           "correct short form config flag with equals",
			args:           []string{"-c=/tmp/test"},
			expectError:    false,
			expectedConfig: "/tmp/test",
			description:    "This is the recommended way to use the config flag with short form",
		},
		{
			name:           "correct short form config flag with space",
			args:           []string{"-c", "/tmp/test"},
			expectError:    false,
			expectedConfig: "/tmp/test",
			description:    "This is also valid short form usage",
		},
		{
			name:           "misleading but technically valid: -config=value",
			args:           []string{"-config=/tmp/test"},
			expectError:    false,
			expectedConfig: "onfig=/tmp/test",
			description:    "This looks wrong but pflag parses -config=/tmp/test as -c onfig=/tmp/test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New(ui)
			err := c.flags.Parse(tt.args)

			if tt.expectError {
				assert.Error(t, err, tt.description)
			} else {
				assert.NoError(t, err, tt.description)
				assert.Equal(t, tt.expectedConfig, c.configFile,
					"Config value should match expected for: %s", tt.description)
			}
		})
	}
}
