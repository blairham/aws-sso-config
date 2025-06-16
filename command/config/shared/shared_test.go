package shared

import (
	"testing"

	"github.com/stretchr/testify/assert"

	appconfig "github.com/blairham/aws-sso-config/providers/config"
)

func TestIsValidKey(t *testing.T) {
	// Test valid keys
	validKeys := []string{
		"sso.start_url",
		"sso.region",
		"sso.role",
		"aws.default_region",
		"aws.config_file",
	}

	for _, key := range validKeys {
		assert.True(t, IsValidKey(key), "Key %s should be valid", key)
	}

	// Test invalid keys
	invalidKeys := []string{
		"invalid_key",
		"",
		"sso_start_url_invalid",
		"non_existent",
	}

	for _, key := range invalidKeys {
		assert.False(t, IsValidKey(key), "Key %s should be invalid", key)
	}
}

func TestValidKeysConstant(t *testing.T) {
	// Verify that ValidKeys has the expected number of keys
	assert.Len(t, ValidKeys, 5, "ValidKeys should contain 5 keys")

	// Verify all expected keys are present
	expectedKeys := map[string]bool{
		"sso.start_url":      true,
		"sso.region":         true,
		"sso.role":           true,
		"aws.default_region": true,
		"aws.config_file":    true,
	}

	for _, key := range ValidKeys {
		assert.True(t, expectedKeys[key], "Unexpected key in ValidKeys: %s", key)
	}
}

func TestKeyDescriptions(t *testing.T) {
	// Verify that each valid key has a description
	for _, key := range ValidKeys {
		desc, exists := KeyDescriptions[key]
		assert.True(t, exists, "Key %s should have a description", key)
		assert.NotEmpty(t, desc, "Description for key %s should not be empty", key)
	}
}

func TestGetConfigValue(t *testing.T) {
	// Create a test config
	config := &appconfig.Config{}
	config.SSO.StartURL = "https://test.awsapps.com/start"
	config.SSO.Region = "us-west-2"
	config.SSO.Role = "TestRole"
	config.AWS.DefaultRegion = "us-east-1"
	config.AWS.ConfigFile = "/test/config"

	tests := []struct {
		key      string
		expected string
	}{
		{KeySSOStartURL, "https://test.awsapps.com/start"},
		{KeySSORegion, "us-west-2"},
		{KeySSORole, "TestRole"},
		{KeyAWSDefaultRegion, "us-east-1"},
		{KeyAWSConfigFile, "/test/config"},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			value, err := GetConfigValue(config, tt.key)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, value)
		})
	}

	// Test invalid key
	_, err := GetConfigValue(config, "invalid_key")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown configuration key")
}

func TestSetConfigValue(t *testing.T) {
	config := &appconfig.Config{}

	tests := []struct {
		key   string
		value string
	}{
		{KeySSOStartURL, "https://new.awsapps.com/start"},
		{KeySSORegion, "eu-west-1"},
		{KeySSORole, "NewRole"},
		{KeyAWSDefaultRegion, "ap-south-1"},
		{KeyAWSConfigFile, "/new/config"},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			err := SetConfigValue(config, tt.key, tt.value)
			assert.NoError(t, err)

			// Verify the value was set
			value, err := GetConfigValue(config, tt.key)
			assert.NoError(t, err)
			assert.Equal(t, tt.value, value)
		})
	}

	// Test invalid key
	err := SetConfigValue(config, "invalid_key", "test_value")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown configuration key")
}

func TestAllValidKeysHaveConstants(t *testing.T) {
	// Ensure all valid keys have corresponding constants
	expectedConstants := map[string]string{
		"sso.start_url":      KeySSOStartURL,
		"sso.region":         KeySSORegion,
		"sso.role":           KeySSORole,
		"aws.default_region": KeyAWSDefaultRegion,
		"aws.config_file":    KeyAWSConfigFile,
	}

	for validKey, expectedConstant := range expectedConstants {
		assert.Contains(t, ValidKeys, validKey)
		assert.Equal(t, validKey, expectedConstant)
	}
}

func TestKeyConstantsMatch(t *testing.T) {
	// Test that key constants match their string values
	assert.Equal(t, "sso.start_url", KeySSOStartURL)
	assert.Equal(t, "sso.region", KeySSORegion)
	assert.Equal(t, "sso.role", KeySSORole)
	assert.Equal(t, "aws.default_region", KeyAWSDefaultRegion)
	assert.Equal(t, "aws.config_file", KeyAWSConfigFile)
}
