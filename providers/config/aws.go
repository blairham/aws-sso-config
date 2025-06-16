package config

import (
	"fmt"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

// AWSConfig holds AWS-specific configuration
type AWSConfig struct {
	DefaultRegion string `mapstructure:"default_region" toml:"default_region"`
	ConfigFile    string `mapstructure:"config_file" toml:"config_file"`
}

// DefaultAWS returns the default AWS configuration
func DefaultAWS() AWSConfig {
	home, _ := homedir.Dir()
	return AWSConfig{
		DefaultRegion: "us-east-1",
		ConfigFile:    filepath.Join(home, ".aws", "config"),
	}
}

// Validate validates the AWS configuration
func (a *AWSConfig) Validate() error {
	if a.DefaultRegion == "" {
		return fmt.Errorf("AWS default region is required")
	}
	if a.ConfigFile == "" {
		return fmt.Errorf("AWS config file path is required")
	}
	return nil
}

// SetDefaults sets default values for any missing AWS configuration
func (a *AWSConfig) SetDefaults() {
	if a.DefaultRegion == "" {
		a.DefaultRegion = "us-east-1"
	}
	if a.ConfigFile == "" {
		home, _ := homedir.Dir()
		a.ConfigFile = filepath.Join(home, ".aws", "config")
	}
}

// GetSectionName returns the TOML section name for AWS configuration
func (a *AWSConfig) GetSectionName() string {
	return "aws"
}

// GetDefaultContent returns the default TOML content for AWS section
func (a *AWSConfig) GetDefaultContent() string {
	return `# AWS Configuration
[aws]
default_region = "us-east-1"
config_file = "~/.aws/config"
`
}
