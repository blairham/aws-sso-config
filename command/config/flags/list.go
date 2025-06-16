package flags

// ListFlag represents the list flag configuration
type ListFlag struct {
	BaseFlag
}

// NewListFlag creates a new list flag configuration
func NewListFlag() *ListFlag {
	return &ListFlag{
		BaseFlag: BaseFlag{
			Name:        "list",
			ShortFlag:   "l",
			Description: "List all configuration variables and their values",
			Usage:       "List all configuration variables and their values",
		},
	}
}
