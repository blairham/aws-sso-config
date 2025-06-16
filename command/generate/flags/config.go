package flags

// ConfigFlag represents the config flag configuration
type ConfigFlag struct {
	BaseFlag
}

// NewConfigFlag creates a new config flag configuration
func NewConfigFlag() *ConfigFlag {
	return &ConfigFlag{
		BaseFlag: BaseFlag{
			Name:        "config",
			ShortFlag:   "c",
			Description: "Path to configuration file",
			Usage:       "Path to configuration file. If not specified, uses environment variables and defaults.",
		},
	}
}
