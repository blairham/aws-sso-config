# Config Flags Package

This package manages all flags for the `config` command, providing a clean and modular way to add and manage command-line flags using [spf13/pflag](https://github.com/spf13/pflag).

## Architecture

The flags package uses an interface-based approach where each flag type implements the `Flag` interface:

```go
type Flag interface {
    GetFlagName() string      // Returns the flag name (e.g., "list")
    GetShortFlag() string     // Returns the short flag (e.g., "l")
    GetDescription() string   // Returns description for help text
    GetUsage() string        // Returns usage information
}
```

The `FlagRegistry` manages all available flags and provides methods to:
- Get all flags: `GetAllFlags()`
- Get a specific flag by name: `GetFlagByName(name)`

## Adding a New Flag

To add a new flag, follow these steps:

### 1. Create a new flag definition file

Create a new file in this directory (e.g., `newflag.go`) with the following structure:

```go
package flags

// NewFlag represents the new flag configuration
type NewFlag struct {
    Name        string
    ShortFlag   string
    Description string
    Usage       string
}

// NewNewFlag creates a new flag configuration
func NewNewFlag() *NewFlag {
    return &NewFlag{
        Name:        "newflag",
        ShortFlag:   "n",
        Description: "Description of what the new flag does",
        Usage:       "Usage information for the new flag",
    }
}

// GetFlagName returns the flag name
func (f *NewFlag) GetFlagName() string {
    return f.Name
}

// GetShortFlag returns the short flag
func (f *NewFlag) GetShortFlag() string {
    return f.ShortFlag
}

// GetDescription returns the flag description
func (f *NewFlag) GetDescription() string {
    return f.Description
}

// GetUsage returns the flag usage information
func (f *NewFlag) GetUsage() string {
    return f.Usage
}
```

### 2. Register the flag

Add your new flag to the `NewFlagRegistry()` function in `flags.go`:

```go
func NewFlagRegistry() *FlagRegistry {
    return &FlagRegistry{
        flags: []Flag{
            NewListFlag(),
            NewNewFlag(),  // Add your new flag here
        },
    }
}
```

### 3. Add flag handling in the config command

In `command/config/config.go`, add the flag variable and registration:

```go
type cmd struct {
    UI    cli.Ui
    flags *pflag.FlagSet

    // Flag variables
    list    bool
    newflag bool  // Add your new flag variable
}

func (c *cmd) init() {
    c.flags = pflag.NewFlagSet("config", pflag.ContinueOnError)

    // Get flag configurations from registry
    registry := configflags.NewFlagRegistry()

    // Register list flag
    listFlag := registry.GetFlagByName("list")
    c.flags.BoolVarP(&c.list, listFlag.GetFlagName(), listFlag.GetShortFlag(), false, listFlag.GetDescription())
    c.flags.MarkHidden("list") // Make it a secret flag

    // Register your new flag
    newFlag := registry.GetFlagByName("newflag")
    c.flags.BoolVarP(&c.newflag, newFlag.GetFlagName(), newFlag.GetShortFlag(), false, newFlag.GetDescription())
}

func (c *cmd) Run(args []string) int {
    // Parse flags
    if err := c.flags.Parse(args); err != nil {
        return 1
    }

    // Get remaining args after flag parsing
    remainingArgs := c.flags.Args()

    // Check if your new flag was used
    if c.newflag {
        // Handle your new flag logic here
        c.UI.Output("New flag executed!")
        return 0
    }

    // ... rest of the Run method
}
### 4. Add tests

Create tests for your new flag in the test file to ensure it works correctly.

## Hidden/Secret Flags

To make a flag hidden (not shown in help output), use pflag's `MarkHidden` method:

```go
c.flags.MarkHidden("flagname")
```

This is useful for maintaining backward compatibility or providing undocumented functionality.

## Existing Flags

### List Flag (`-l`, `--list`)
- **File**: `list.go`
- **Purpose**: Lists all configuration variables and their values
- **Flags**: `-l`, `--list`
- **Behavior**: Calls the `config list` subcommand
- **Status**: Hidden/secret flag (not shown in help)

## Benefits of This Architecture

1. **Modularity**: Each flag is self-contained in its own definition
2. **Testability**: Each flag can be tested independently  
3. **Extensibility**: Adding new flags doesn't require modifying existing code
4. **Clean Separation**: Flag definitions are separated from command logic
5. **pflag Integration**: Full support for both short and long flag forms
6. **Consistency**: All flags follow the same interface pattern
7. **Hidden Flags**: Support for secret/hidden flags that don't appear in help

## Implementation Details

- Uses spf13/pflag for advanced flag parsing capabilities
- Flags are processed before subcommands in the main config command
- If a flag is found, its logic executes before subcommand processing
- Help text is automatically generated using pflag's FlagUsages()
- Hidden flags can be used but don't appear in help output
- Both short (`-l`) and long (`--list`) flag forms are supported
