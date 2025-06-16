package flags

// Flag represents a common interface for all flags
type Flag interface {
	GetFlagName() string
	GetShortFlag() string
	GetDescription() string
	GetUsage() string
}

// BaseFlag provides a common implementation for all flags
type BaseFlag struct {
	Name        string
	ShortFlag   string
	Description string
	Usage       string
}

// GetFlagName returns the flag name
func (f *BaseFlag) GetFlagName() string {
	return f.Name
}

// GetShortFlag returns the short flag
func (f *BaseFlag) GetShortFlag() string {
	return f.ShortFlag
}

// GetDescription returns the flag description
func (f *BaseFlag) GetDescription() string {
	return f.Description
}

// GetUsage returns the flag usage information
func (f *BaseFlag) GetUsage() string {
	return f.Usage
}

// FlagRegistry manages all available flags for the config command
type FlagRegistry struct {
	flags []Flag
}

// NewFlagRegistry creates a new flag registry with all available flags
func NewFlagRegistry() *FlagRegistry {
	return &FlagRegistry{
		flags: []Flag{
			NewListFlag(),
		},
	}
}

// GetAllFlags returns all registered flags
func (r *FlagRegistry) GetAllFlags() []Flag {
	return r.flags
}

// GetFlagByName returns a flag by its name, or nil if not found
func (r *FlagRegistry) GetFlagByName(name string) Flag {
	for _, flag := range r.flags {
		if flag.GetFlagName() == name {
			return flag
		}
	}
	return nil
}
