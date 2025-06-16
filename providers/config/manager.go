package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

// ConfigManager handles loading from a single configuration file with multiple sections
type ConfigManager struct {
	configFile string
}

// NewConfigManager creates a new configuration manager
func NewConfigManager(configFile string) *ConfigManager {
	if configFile == "" {
		home, _ := homedir.Dir()
		configFile = filepath.Join(home, ".awsssoconfig")
	}
	return &ConfigManager{configFile: configFile}
}

// Load loads configuration from a single file with multiple sections
func (cm *ConfigManager) Load() (*Config, error) {
	v := viper.New()

	// Set up configuration file
	v.SetConfigFile(cm.configFile)
	v.SetConfigType("toml")

	// Environment variable configuration
	v.SetEnvPrefix("AWS_SSO_CONFIG")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Read config file (create if not found)
	if err := v.ReadInConfig(); err != nil {
		// Check for "file not found" error more broadly
		if strings.Contains(err.Error(), "no such file or directory") ||
			strings.Contains(err.Error(), "cannot find the file") {
			// Create default config file
			if createErr := cm.createDefaultConfig(); createErr == nil {
				// Try reading again after creation
				if readErr := v.ReadInConfig(); readErr != nil {
					return nil, fmt.Errorf("error reading newly created config file: %w", readErr)
				}
			} else {
				return nil, fmt.Errorf("error creating config file: %w", createErr)
			}
		} else {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	config := &Config{}

	// Load SSO section
	if ssoData := v.Sub("sso"); ssoData != nil {
		if err := ssoData.Unmarshal(&config.SSO); err != nil {
			return nil, fmt.Errorf("error unmarshaling SSO config: %w", err)
		}
	}

	// Load AWS section
	if awsData := v.Sub("aws"); awsData != nil {
		if err := awsData.Unmarshal(&config.AWS); err != nil {
			return nil, fmt.Errorf("error unmarshaling AWS config: %w", err)
		}
	}

	// Set defaults for any missing values
	config.SetDefaults()

	return config, nil
}

// createDefaultConfig creates a default configuration file with all sections
func (cm *ConfigManager) createDefaultConfig() error {
	// Ensure the directory exists
	configDir := filepath.Dir(cm.configFile)
	if err := os.MkdirAll(configDir, 0750); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Combine default content from all sections
	sso := DefaultSSO()
	aws := DefaultAWS()

	content := sso.GetDefaultContent() + "\n" + aws.GetDefaultContent()

	return os.WriteFile(cm.configFile, []byte(content), 0600)
}

// SaveProviderConfig saves a provider configuration section to the single config file
func (cm *ConfigManager) SaveProviderConfig(provider string, data interface{}) error {
	v := viper.New()
	v.SetConfigFile(cm.configFile)
	v.SetConfigType("toml")

	// Try to read existing config first
	_ = v.ReadInConfig() // Ignore error if file doesn't exist

	// Set the provider section data
	switch provider {
	case "sso":
		if ssoData, ok := data.(SSOConfig); ok {
			if ssoData.StartURL != "" {
				v.Set("sso.start_url", ssoData.StartURL)
			}
			if ssoData.Region != "" {
				v.Set("sso.region", ssoData.Region)
			}
			if ssoData.Role != "" {
				v.Set("sso.role", ssoData.Role)
			}
		}
	case "aws":
		if awsData, ok := data.(AWSConfig); ok {
			if awsData.DefaultRegion != "" {
				v.Set("aws.default_region", awsData.DefaultRegion)
			}
			if awsData.ConfigFile != "" {
				v.Set("aws.config_file", awsData.ConfigFile)
			}
		}
	default:
		return fmt.Errorf("unknown provider: %s", provider)
	}

	// Ensure the directory exists
	configDir := filepath.Dir(cm.configFile)
	if err := os.MkdirAll(configDir, 0750); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	return v.WriteConfig()
}

// LoadConfigForKey loads only the configuration needed for a specific key
func LoadConfigForKey(configFile string, key string) (*Config, error) {
	var cm *ConfigManager

	if configFile != "" {
		cm = NewConfigManager(configFile)
	} else {
		cm = NewConfigManager("")
	}

	config := &Config{}

	// Load the full config since it's in a single file anyway
	fullConfig, err := cm.Load()
	if err != nil {
		// If loading fails, create defaults for the specific provider
		if strings.HasPrefix(key, "sso.") {
			config.SSO = DefaultSSO()
		} else if strings.HasPrefix(key, "aws.") {
			config.AWS = DefaultAWS()
		} else {
			return nil, fmt.Errorf("unknown key prefix for key: %s", key)
		}
		return config, nil
	}

	return fullConfig, nil
}

// Load loads configuration using the default configuration manager
func Load(configPath string) (*Config, error) {
	var cm *ConfigManager

	if configPath != "" {
		cm = NewConfigManager(configPath)
	} else {
		cm = NewConfigManager("")
	}

	return cm.Load()
}
