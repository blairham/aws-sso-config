# Generate Command Flags

This package contains flag definitions for the `generate` command using [spf13/pflag](https://github.com/spf13/pflag). Each flag is implemented in its own file, making it easy to add or remove flags.

## Architecture

The flag system uses a modular approach with pflag where:

1. **Each flag is in its own file** (e.g., `diff.go`, `config.go`)
2. **Common interface**: All flags implement the `Flag` interface:
   - `GetFlagName() string` - Returns the long flag name (without dashes)
   - `GetShortFlag() string` - Returns the short flag name (single character)
   - `GetDescription() string` - Returns a description for help text
   - `GetUsage() string` - Returns detailed usage information

3. **Central registry**: The `FlagRegistry` manages all flags and provides:
   - `GetAllFlags()` - Returns all registered flags
   - `GetFlagByName(name)` - Returns a specific flag by name

4. **pflag integration**: Uses `pflag.FlagSet` for enhanced flag parsing with:
   - Long flags (e.g., `--diff`, `--config`)
   - Short flags (e.g., `-d`, `-c`)
   - POSIX-style flag parsing
   - Better error messages

## Available Flags

### --diff, -d

**File:** `diff.go`

The diff flag enables diff output to show changes before applying them to the AWS config file.

**Usage:**
```bash
aws-sso-config generate --diff
aws-sso-config generate -d
aws-sso-config generate --config my-config.toml --diff
aws-sso-config generate -c my-config.toml -d
```

### --config, -c

**File:** `config.go`

The config flag specifies a custom configuration file path.

**Usage:**
```bash
aws-sso-config generate --config my-config.toml
aws-sso-config generate -c my-config.toml
aws-sso-config generate --config /path/to/config.toml --diff
aws-sso-config generate -c /path/to/config.toml -d
```

## Adding New Flags

To add a new flag:

1. **Create a new file** (e.g., `verbose.go`) with:
   ```go
   package flags

   type VerboseFlag struct {
       Name        string
       ShortFlag   string
       Description string
       Usage       string
   }

   func NewVerboseFlag() *VerboseFlag {
       return &VerboseFlag{
           Name:        "verbose",
           ShortFlag:   "v",
           Description: "Enable verbose output",
           Usage:       "Enable verbose output for debugging",
       }
   }

   func (f *VerboseFlag) GetFlagName() string { return f.Name }
   func (f *VerboseFlag) GetShortFlag() string { return f.ShortFlag }
   func (f *VerboseFlag) GetDescription() string { return f.Description }
   func (f *VerboseFlag) GetUsage() string { return f.Usage }
   ```

2. **Register the flag** in `flags.go` by adding it to `NewFlagRegistry()`:
   ```go
   flags: []Flag{
       NewDiffFlag(),
       NewConfigFlag(),
       NewVerboseFlag(), // Add your new flag here
   }
   ```

3. **Update the generate command** in `generate.go` to use the new flag:
   ```go
   verboseFlag := registry.GetFlagByName("verbose")
   c.flags.BoolVarP(&c.verbose, verboseFlag.GetFlagName(), verboseFlag.GetShortFlag(), false, verboseFlag.GetDescription())
   ```

4. **Add tests** in `flags_test.go`
5. **Update this README**

## Removing Flags

To remove a flag:

1. Delete the flag's file (e.g., `diff.go`)
2. Remove it from the registry in `flags.go`
3. Remove its usage from `generate.go`
4. Update tests and documentation

## pflag Benefits

Using pflag provides several advantages over the standard flag package:

- **POSIX-style flags**: Support for both `--flag` and `-f` syntax
- **Better error messages**: More user-friendly error reporting
- **Flag combinations**: Support for combining short flags (e.g., `-dc`)
- **Consistent behavior**: Matches GNU-style flag conventions
- **No flag reordering**: Flags can appear anywhere in the command line

## Testing

Run tests for this package:

```bash
go test ./command/generate/flags/
```

## File Structure

```
command/generate/flags/
├── README.md          # This documentation
├── flags.go          # Flag interface and registry (pflag integration)
├── flags_test.go     # Tests for the flag system
├── diff.go           # Diff flag implementation
└── config.go         # Config flag implementation
```
