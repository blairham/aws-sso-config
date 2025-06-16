package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfigManager(t *testing.T) {
	t.Run("with custom config file", func(t *testing.T) {
		customPath := "/tmp/custom-config"
		cm := NewConfigManager(customPath)
		assert.Equal(t, customPath, cm.configFile)
	})

	t.Run("with empty config file defaults to home directory", func(t *testing.T) {
		cm := NewConfigManager("")
		assert.Contains(t, cm.configFile, ".awsssoconfig")
	})
}

func TestConfigManagerLoad(t *testing.T) {
	t.Run("load non-existent config creates default file and returns config", func(t *testing.T) {
		tempDir := t.TempDir()
		configFile := filepath.Join(tempDir, "test-config")
		cm := NewConfigManager(configFile)

		// First verify the file doesn't exist
		assert.NoFileExists(t, configFile)

		config, err := cm.Load()

		// If error occurs, print for debugging
		if err != nil {
			t.Logf("Load error: %v", err)
		}

		require.NoError(t, err)
		assert.NotNil(t, config)

		// Check that defaults are set
		assert.Equal(t, "https://your-sso-portal.awsapps.com/start", config.SSO.StartURL)
		assert.Equal(t, "us-east-1", config.SSO.Region)
		assert.Equal(t, "AdministratorAccess", config.SSO.Role)
		assert.Equal(t, "us-east-1", config.AWS.DefaultRegion)
		assert.Contains(t, config.AWS.ConfigFile, ".aws/config")

		// Check that config file was created
		assert.FileExists(t, configFile)
	})

	t.Run("load existing config", func(t *testing.T) {
		tempDir := t.TempDir()
		configFile := filepath.Join(tempDir, "test-config")

		// Create a test config file
		content := `[sso]
start_url = "https://test.awsapps.com/start"
region = "us-west-2"
role = "TestRole"

[aws]
default_region = "eu-central-1"
config_file = "/custom/aws/config"
`
		err := os.WriteFile(configFile, []byte(content), 0600)
		require.NoError(t, err)

		cm := NewConfigManager(configFile)
		config, err := cm.Load()
		require.NoError(t, err)

		assert.Equal(t, "https://test.awsapps.com/start", config.SSO.StartURL)
		assert.Equal(t, "us-west-2", config.SSO.Region)
		assert.Equal(t, "TestRole", config.SSO.Role)
		assert.Equal(t, "eu-central-1", config.AWS.DefaultRegion)
		assert.Equal(t, "/custom/aws/config", config.AWS.ConfigFile)
	})

	t.Run("load config with missing sections", func(t *testing.T) {
		tempDir := t.TempDir()
		configFile := filepath.Join(tempDir, "test-config")

		// Create a config file with only SSO section
		content := `[sso]
start_url = "https://test.awsapps.com/start"
region = "us-west-2"
`
		err := os.WriteFile(configFile, []byte(content), 0600)
		require.NoError(t, err)

		cm := NewConfigManager(configFile)
		config, err := cm.Load()
		require.NoError(t, err)

		// SSO values should be loaded
		assert.Equal(t, "https://test.awsapps.com/start", config.SSO.StartURL)
		assert.Equal(t, "us-west-2", config.SSO.Region)
		assert.Equal(t, "AdministratorAccess", config.SSO.Role) // default

		// AWS values should be defaults
		assert.Equal(t, "us-east-1", config.AWS.DefaultRegion)
		assert.Contains(t, config.AWS.ConfigFile, ".aws/config")
	})
}

func TestConfigManagerSaveProviderConfig(t *testing.T) {
	t.Run("save SSO config", func(t *testing.T) {
		tempDir := t.TempDir()
		configFile := filepath.Join(tempDir, "test-config")
		cm := NewConfigManager(configFile)

		ssoConfig := SSOConfig{
			StartURL: "https://example.awsapps.com/start",
			Region:   "eu-west-1",
			Role:     "MyRole",
		}

		err := cm.SaveProviderConfig("sso", ssoConfig)
		require.NoError(t, err)

		// Verify file was created and contains correct content
		assert.FileExists(t, configFile)

		// Load config and verify values
		config, err := cm.Load()
		require.NoError(t, err)
		assert.Equal(t, "https://example.awsapps.com/start", config.SSO.StartURL)
		assert.Equal(t, "eu-west-1", config.SSO.Region)
		assert.Equal(t, "MyRole", config.SSO.Role)
	})

	t.Run("save AWS config", func(t *testing.T) {
		tempDir := t.TempDir()
		configFile := filepath.Join(tempDir, "test-config")
		cm := NewConfigManager(configFile)

		awsConfig := AWSConfig{
			DefaultRegion: "ap-southeast-1",
			ConfigFile:    "/custom/aws/config",
		}

		err := cm.SaveProviderConfig("aws", awsConfig)
		require.NoError(t, err)

		// Load config and verify values
		config, err := cm.Load()
		require.NoError(t, err)
		assert.Equal(t, "ap-southeast-1", config.AWS.DefaultRegion)
		assert.Equal(t, "/custom/aws/config", config.AWS.ConfigFile)
	})

	t.Run("save to existing config preserves other sections", func(t *testing.T) {
		tempDir := t.TempDir()
		configFile := filepath.Join(tempDir, "test-config")
		cm := NewConfigManager(configFile)

		// First save SSO config
		ssoConfig := SSOConfig{
			StartURL: "https://example.awsapps.com/start",
			Region:   "us-east-1",
			Role:     "SSORole",
		}
		err := cm.SaveProviderConfig("sso", ssoConfig)
		require.NoError(t, err)

		// Then save AWS config
		awsConfig := AWSConfig{
			DefaultRegion: "us-west-2",
			ConfigFile:    "/custom/config",
		}
		err = cm.SaveProviderConfig("aws", awsConfig)
		require.NoError(t, err)

		// Verify both sections are present
		config, err := cm.Load()
		require.NoError(t, err)
		assert.Equal(t, "https://example.awsapps.com/start", config.SSO.StartURL)
		assert.Equal(t, "us-east-1", config.SSO.Region)
		assert.Equal(t, "SSORole", config.SSO.Role)
		assert.Equal(t, "us-west-2", config.AWS.DefaultRegion)
		assert.Equal(t, "/custom/config", config.AWS.ConfigFile)
	})

	t.Run("invalid provider returns error", func(t *testing.T) {
		tempDir := t.TempDir()
		configFile := filepath.Join(tempDir, "test-config")
		cm := NewConfigManager(configFile)

		err := cm.SaveProviderConfig("invalid", "some data")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unknown provider: invalid")
	})
}

func TestLoadConfigForKey(t *testing.T) {
	t.Run("load config for SSO key", func(t *testing.T) {
		tempDir := t.TempDir()
		configFile := filepath.Join(tempDir, "test-config")

		// Create a test config file
		content := `[sso]
start_url = "https://test.awsapps.com/start"
region = "us-west-2"
role = "TestRole"
`
		err := os.WriteFile(configFile, []byte(content), 0600)
		require.NoError(t, err)

		config, err := LoadConfigForKey(configFile, "sso.start_url")
		require.NoError(t, err)
		assert.Equal(t, "https://test.awsapps.com/start", config.SSO.StartURL)
	})

	t.Run("load config for AWS key", func(t *testing.T) {
		tempDir := t.TempDir()
		configFile := filepath.Join(tempDir, "test-config")

		// Create a test config file
		content := `[aws]
default_region = "eu-central-1"
config_file = "/custom/aws/config"
`
		err := os.WriteFile(configFile, []byte(content), 0600)
		require.NoError(t, err)

		config, err := LoadConfigForKey(configFile, "aws.default_region")
		require.NoError(t, err)
		assert.Equal(t, "eu-central-1", config.AWS.DefaultRegion)
	})

	t.Run("load config for key with missing file returns defaults", func(t *testing.T) {
		config, err := LoadConfigForKey("/nonexistent/config", "sso.start_url")
		require.NoError(t, err)
		assert.Equal(t, "https://your-sso-portal.awsapps.com/start", config.SSO.StartURL)
	})

	t.Run("invalid key prefix returns error", func(t *testing.T) {
		config, err := LoadConfigForKey("/nonexistent/config", "invalid.key")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unknown key prefix for key: invalid.key")
		assert.Nil(t, config)
	})
}

func TestLoad(t *testing.T) {
	t.Run("load with custom config path", func(t *testing.T) {
		tempDir := t.TempDir()
		configFile := filepath.Join(tempDir, "test-config")

		// Create a test config file
		content := `[sso]
start_url = "https://test.awsapps.com/start"
region = "us-west-2"
`
		err := os.WriteFile(configFile, []byte(content), 0600)
		require.NoError(t, err)

		config, err := Load(configFile)
		require.NoError(t, err)
		assert.Equal(t, "https://test.awsapps.com/start", config.SSO.StartURL)
		assert.Equal(t, "us-west-2", config.SSO.Region)
	})

	t.Run("load with empty config path uses default location", func(t *testing.T) {
		// Clean up any existing config file
		homeConfigFile := filepath.Join(os.Getenv("HOME"), ".awsssoconfig")
		originalContent := ""
		configExists := false

		// Back up existing config if it exists
		if data, err := os.ReadFile(homeConfigFile); err == nil {
			originalContent = string(data)
			configExists = true
		}

		// Remove the config file for this test
		os.Remove(homeConfigFile)

		// Restore after test
		defer func() {
			if configExists {
				os.WriteFile(homeConfigFile, []byte(originalContent), 0600)
			} else {
				os.Remove(homeConfigFile)
			}
		}()

		config, err := Load("")
		require.NoError(t, err)
		assert.NotNil(t, config)
		// Should have default values
		assert.Equal(t, "https://your-sso-portal.awsapps.com/start", config.SSO.StartURL)
	})
}

func TestConfigValidate(t *testing.T) {
	t.Run("valid config passes validation", func(t *testing.T) {
		config := &Config{
			SSO: SSOConfig{
				StartURL: "https://test.awsapps.com/start",
				Region:   "us-east-1",
			},
			AWS: AWSConfig{
				DefaultRegion: "us-east-1",
				ConfigFile:    "/home/user/.aws/config",
			},
		}

		err := config.Validate()
		assert.NoError(t, err)
	})

	t.Run("missing SSO start URL fails validation", func(t *testing.T) {
		config := &Config{
			SSO: SSOConfig{
				Region: "us-east-1",
			},
			AWS: AWSConfig{
				DefaultRegion: "us-east-1",
				ConfigFile:    "/home/user/.aws/config",
			},
		}

		err := config.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "SSO start URL is required")
	})

	t.Run("missing SSO region fails validation", func(t *testing.T) {
		config := &Config{
			SSO: SSOConfig{
				StartURL: "https://test.awsapps.com/start",
			},
			AWS: AWSConfig{
				DefaultRegion: "us-east-1",
				ConfigFile:    "/home/user/.aws/config",
			},
		}

		err := config.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "SSO region is required")
	})

	t.Run("missing AWS default region fails validation", func(t *testing.T) {
		config := &Config{
			SSO: SSOConfig{
				StartURL: "https://test.awsapps.com/start",
				Region:   "us-east-1",
			},
			AWS: AWSConfig{
				ConfigFile: "/home/user/.aws/config",
			},
		}

		err := config.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "AWS default region is required")
	})

	t.Run("missing AWS config file fails validation", func(t *testing.T) {
		config := &Config{
			SSO: SSOConfig{
				StartURL: "https://test.awsapps.com/start",
				Region:   "us-east-1",
			},
			AWS: AWSConfig{
				DefaultRegion: "us-east-1",
			},
		}

		err := config.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "AWS config file path is required")
	})
}

func TestDefault(t *testing.T) {
	config := Default()
	assert.NotNil(t, config)
	assert.Equal(t, "https://your-sso-portal.awsapps.com/start", config.SSO.StartURL)
	assert.Equal(t, "us-east-1", config.SSO.Region)
	assert.Equal(t, "AdministratorAccess", config.SSO.Role)
	assert.Equal(t, "us-east-1", config.AWS.DefaultRegion)
	assert.Contains(t, config.AWS.ConfigFile, ".aws/config")
}

func TestConfigBackwardCompatibilityGetters(t *testing.T) {
	config := &Config{
		SSO: SSOConfig{
			StartURL: "https://test.awsapps.com/start",
			Region:   "us-west-2",
			Role:     "TestRole",
		},
		AWS: AWSConfig{
			DefaultRegion: "eu-central-1",
			ConfigFile:    "/custom/config",
		},
	}

	// Test SSO getters
	assert.Equal(t, "https://test.awsapps.com/start", config.SSOStartURL())
	assert.Equal(t, "us-west-2", config.SSORegion())
	assert.Equal(t, "TestRole", config.SSORole())

	// Test AWS getters
	assert.Equal(t, "eu-central-1", config.DefaultRegion())
	assert.Equal(t, "/custom/config", config.ConfigFile())
}
