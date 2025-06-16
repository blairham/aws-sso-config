package shared

import (
	"fmt"

	appconfig "github.com/blairham/aws-sso-config/providers/config"
)

// IsValidKey checks if the given key is a valid configuration key
func IsValidKey(key string) bool {
	for _, validKey := range ValidKeys {
		if key == validKey {
			return true
		}
	}
	return false
}

// GetConfigValue gets the value for the specified key from the config
func GetConfigValue(config *appconfig.Config, key string) (string, error) {
	switch key {
	case KeySSOStartURL:
		return config.SSO.StartURL, nil
	case KeySSORegion:
		return config.SSO.Region, nil
	case KeySSORole:
		return config.SSO.Role, nil
	case KeyAWSDefaultRegion:
		return config.AWS.DefaultRegion, nil
	case KeyAWSConfigFile:
		return config.AWS.ConfigFile, nil
	default:
		return "", fmt.Errorf("unknown configuration key: %s", key)
	}
}

// SetConfigValue sets the value for the specified key in the config
func SetConfigValue(config *appconfig.Config, key string, value string) error {
	switch key {
	case KeySSOStartURL:
		config.SSO.StartURL = value
		return nil
	case KeySSORegion:
		config.SSO.Region = value
		return nil
	case KeySSORole:
		config.SSO.Role = value
		return nil
	case KeyAWSDefaultRegion:
		config.AWS.DefaultRegion = value
		return nil
	case KeyAWSConfigFile:
		config.AWS.ConfigFile = value
		return nil
	default:
		return fmt.Errorf("unknown configuration key: %s", key)
	}
}

// PrintAvailableKeys prints all available configuration keys with descriptions to Error
func PrintAvailableKeys(ui interface{ Error(string) }) {
	ui.Error("Available keys:")
	for _, key := range ValidKeys {
		if desc, ok := KeyDescriptions[key]; ok {
			ui.Error(fmt.Sprintf("  %-16s %s", key, desc))
		} else {
			ui.Error(fmt.Sprintf("  %s", key))
		}
	}
}

// OutputAvailableKeys prints all available configuration keys with descriptions to Output
func OutputAvailableKeys(ui interface{ Output(string) }) {
	for _, key := range ValidKeys {
		if desc, ok := KeyDescriptions[key]; ok {
			ui.Output(fmt.Sprintf("  %-16s %s", key, desc))
		} else {
			ui.Output(fmt.Sprintf("  %s", key))
		}
	}
}

// SaveConfigValue saves a configuration value to the configuration file
func SaveConfigValue(configFile string, key string, value string) error {
	cm := appconfig.NewConfigManager(configFile)

	// Load current config or use defaults if file doesn't exist
	config, err := cm.Load()
	if err != nil {
		// If file doesn't exist, create with defaults
		config = appconfig.Default()
	}

	// Update the specific field
	switch key {
	case KeySSOStartURL:
		config.SSO.StartURL = value
		err = cm.SaveProviderConfig("sso", config.SSO)
	case KeySSORegion:
		config.SSO.Region = value
		err = cm.SaveProviderConfig("sso", config.SSO)
	case KeySSORole:
		config.SSO.Role = value
		err = cm.SaveProviderConfig("sso", config.SSO)
	case KeyAWSDefaultRegion:
		config.AWS.DefaultRegion = value
		err = cm.SaveProviderConfig("aws", config.AWS)
	case KeyAWSConfigFile:
		config.AWS.ConfigFile = value
		err = cm.SaveProviderConfig("aws", config.AWS)
	default:
		return fmt.Errorf("unknown configuration key: %s", key)
	}

	return err
}
