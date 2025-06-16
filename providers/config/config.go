package config

// Config holds the application configuration
type Config struct {
	// Provider configurations
	SSO SSOConfig `mapstructure:"sso" toml:"sso"`
	AWS AWSConfig `mapstructure:"aws" toml:"aws"`
}

// Backward compatibility getters
func (c *Config) SSOStartURL() string {
	return c.SSO.StartURL
}

func (c *Config) SSORegion() string {
	return c.SSO.Region
}

func (c *Config) SSORole() string {
	return c.SSO.Role
}

// AWS configuration getters
func (c *Config) DefaultRegion() string {
	return c.AWS.DefaultRegion
}

func (c *Config) ConfigFile() string {
	return c.AWS.ConfigFile
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if err := c.SSO.Validate(); err != nil {
		return err
	}
	if err := c.AWS.Validate(); err != nil {
		return err
	}
	return nil
}

// Default returns a default configuration
func Default() *Config {
	return &Config{
		SSO: DefaultSSO(),
		AWS: DefaultAWS(),
	}
}

// SetDefaults sets default values for any missing configuration
func (c *Config) SetDefaults() {
	c.SSO.SetDefaults()
	c.AWS.SetDefaults()
}
