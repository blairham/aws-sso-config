package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSSOConfig(t *testing.T) {
	t.Run("DefaultSSO returns valid defaults", func(t *testing.T) {
		sso := DefaultSSO()
		assert.Equal(t, "https://your-sso-portal.awsapps.com/start", sso.StartURL)
		assert.Equal(t, "us-east-1", sso.Region)
		assert.Equal(t, "AdministratorAccess", sso.Role)
	})

	t.Run("SSO validation passes with valid config", func(t *testing.T) {
		sso := SSOConfig{
			StartURL: "https://test.awsapps.com/start",
			Region:   "us-west-2",
			Role:     "TestRole",
		}
		err := sso.Validate()
		assert.NoError(t, err)
	})

	t.Run("SSO validation fails with missing start URL", func(t *testing.T) {
		sso := SSOConfig{
			Region: "us-west-2",
			Role:   "TestRole",
		}
		err := sso.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "SSO start URL is required")
	})

	t.Run("SSO validation fails with missing region", func(t *testing.T) {
		sso := SSOConfig{
			StartURL: "https://test.awsapps.com/start",
			Role:     "TestRole",
		}
		err := sso.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "SSO region is required")
	})

	t.Run("SSO SetDefaults sets missing values", func(t *testing.T) {
		sso := SSOConfig{}
		sso.SetDefaults()
		assert.Equal(t, "https://your-sso-portal.awsapps.com/start", sso.StartURL)
		assert.Equal(t, "us-east-1", sso.Region)
		assert.Equal(t, "AdministratorAccess", sso.Role)
	})

	t.Run("SSO SetDefaults preserves existing values", func(t *testing.T) {
		sso := SSOConfig{
			StartURL: "https://custom.awsapps.com/start",
			Region:   "eu-west-1",
		}
		sso.SetDefaults()
		assert.Equal(t, "https://custom.awsapps.com/start", sso.StartURL)
		assert.Equal(t, "eu-west-1", sso.Region)
		assert.Equal(t, "AdministratorAccess", sso.Role) // default filled in
	})

	t.Run("SSO GetSectionName returns correct name", func(t *testing.T) {
		sso := SSOConfig{}
		assert.Equal(t, "sso", sso.GetSectionName())
	})

	t.Run("SSO GetDefaultContent returns valid TOML", func(t *testing.T) {
		sso := SSOConfig{}
		content := sso.GetDefaultContent()
		assert.Contains(t, content, "[sso]")
		assert.Contains(t, content, `start_url = "https://your-sso-portal.awsapps.com/start"`)
		assert.Contains(t, content, `region = "us-east-1"`)
		assert.Contains(t, content, `role = "AdministratorAccess"`)
	})
}

func TestAWSConfig(t *testing.T) {
	t.Run("DefaultAWS returns valid defaults", func(t *testing.T) {
		aws := DefaultAWS()
		assert.Equal(t, "us-east-1", aws.DefaultRegion)
		assert.Contains(t, aws.ConfigFile, ".aws/config")
	})

	t.Run("AWS validation passes with valid config", func(t *testing.T) {
		aws := AWSConfig{
			DefaultRegion: "us-west-2",
			ConfigFile:    "/home/user/.aws/config",
		}
		err := aws.Validate()
		assert.NoError(t, err)
	})

	t.Run("AWS validation fails with missing default region", func(t *testing.T) {
		aws := AWSConfig{
			ConfigFile: "/home/user/.aws/config",
		}
		err := aws.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "AWS default region is required")
	})

	t.Run("AWS validation fails with missing config file", func(t *testing.T) {
		aws := AWSConfig{
			DefaultRegion: "us-west-2",
		}
		err := aws.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "AWS config file path is required")
	})

	t.Run("AWS SetDefaults sets missing values", func(t *testing.T) {
		aws := AWSConfig{}
		aws.SetDefaults()
		assert.Equal(t, "us-east-1", aws.DefaultRegion)
		assert.Contains(t, aws.ConfigFile, ".aws/config")
	})

	t.Run("AWS SetDefaults preserves existing values", func(t *testing.T) {
		aws := AWSConfig{
			DefaultRegion: "eu-central-1",
		}
		aws.SetDefaults()
		assert.Equal(t, "eu-central-1", aws.DefaultRegion)
		assert.Contains(t, aws.ConfigFile, ".aws/config") // default filled in
	})

	t.Run("AWS GetSectionName returns correct name", func(t *testing.T) {
		aws := AWSConfig{}
		assert.Equal(t, "aws", aws.GetSectionName())
	})

	t.Run("AWS GetDefaultContent returns valid TOML", func(t *testing.T) {
		aws := AWSConfig{}
		content := aws.GetDefaultContent()
		assert.Contains(t, content, "[aws]")
		assert.Contains(t, content, `default_region = "us-east-1"`)
		assert.Contains(t, content, `config_file = "~/.aws/config"`)
	})
}
