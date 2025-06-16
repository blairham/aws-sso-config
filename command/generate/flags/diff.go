package flags

// DiffFlag represents the diff flag configuration
type DiffFlag struct {
	BaseFlag
}

// NewDiffFlag creates a new diff flag configuration
func NewDiffFlag() *DiffFlag {
	return &DiffFlag{
		BaseFlag: BaseFlag{
			Name:        "diff",
			ShortFlag:   "d",
			Description: "Enable diff output to show changes before applying them",
			Usage:       "Show differences that would be made to the AWS config file without applying them",
		},
	}
}
