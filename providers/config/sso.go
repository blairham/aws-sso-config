package config

import "fmt"

// SSOConfig holds SSO-specific configuration
type SSOConfig struct {
	StartURL string `mapstructure:"start_url" toml:"start_url"`
	Region   string `mapstructure:"region" toml:"region"`
	Role     string `mapstructure:"role" toml:"role"`
}

// DefaultSSO returns the default SSO configuration
func DefaultSSO() SSOConfig {
	return SSOConfig{
		StartURL: "https://your-sso-portal.awsapps.com/start",
		Region:   "us-east-1",
		Role:     "AdministratorAccess",
	}
}

// Validate validates the SSO configuration
func (s *SSOConfig) Validate() error {
	if s.StartURL == "" {
		return fmt.Errorf("SSO start URL is required")
	}
	if s.Region == "" {
		return fmt.Errorf("SSO region is required")
	}
	return nil
}

// SetDefaults sets default values for any missing SSO configuration
func (s *SSOConfig) SetDefaults() {
	if s.StartURL == "" {
		s.StartURL = "https://your-sso-portal.awsapps.com/start"
	}
	if s.Region == "" {
		s.Region = "us-east-1"
	}
	if s.Role == "" {
		s.Role = "AdministratorAccess"
	}
}

// GetSectionName returns the TOML section name for SSO configuration
func (s *SSOConfig) GetSectionName() string {
	return "sso"
}

// GetDefaultContent returns the default TOML content for SSO section
func (s *SSOConfig) GetDefaultContent() string {
	return `# AWS SSO Configuration
[sso]
start_url = "https://your-sso-portal.awsapps.com/start"
region = "us-east-1"
role = "AdministratorAccess"
`
}
