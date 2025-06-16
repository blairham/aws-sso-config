package shared

// Configuration key constants
const (
	KeySSOStartURL      = "sso.start_url"
	KeySSORegion        = "sso.region"
	KeySSORole          = "sso.role"
	KeyAWSDefaultRegion = "aws.default_region"
	KeyAWSConfigFile    = "aws.config_file"
)

// ValidKeys contains all valid configuration keys
var ValidKeys = []string{
	KeySSOStartURL,
	KeySSORegion,
	KeySSORole,
	KeyAWSDefaultRegion,
	KeyAWSConfigFile,
}

// KeyDescriptions maps configuration keys to their descriptions
var KeyDescriptions = map[string]string{
	KeySSOStartURL:      "Your AWS SSO start URL",
	KeySSORegion:        "AWS region for SSO (e.g., us-east-1)",
	KeySSORole:          "SSO role name (e.g., AdministratorAccess)",
	KeyAWSDefaultRegion: "Default AWS region for profiles",
	KeyAWSConfigFile:    "Path to AWS config file",
}
