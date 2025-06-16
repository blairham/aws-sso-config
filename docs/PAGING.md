# Interactive Paging Implementation

The AWS SSO Config tool now includes full interactive paging support for the `config list` command, providing a `less`-like experience with complete keyboard navigation.

## Interactive Features

### Full Navigation Support
- **Arrow Keys**: ↑/↓ move up/down one line
- **Page Navigation**: Space bar for page down, 'b' for page up
- **Search**: '/' to search forward, '?' to search backward
- **Quit**: 'q' to exit the pager
- **Home/End**: 'g' to go to beginning, 'G' to go to end
- **Help**: 'h' for help (in less)

### Smart Activation
- Automatically detects when output would exceed the terminal height
- Only activates paging when running in a terminal (not when piped or redirected)
- Respects common Unix conventions for paging behavior

### Environment Variable Support
- `NO_PAGER`: Disables paging entirely when set to any value
- `PAGER`: Specifies which pager program to use (default: `less -FRX`)
- `AWS_SSO_CONFIG_PAGER`: Tool-specific pager override
- `LINES`: Override terminal height detection

### Pager Selection
The tool automatically selects the best available pager in this order:
1. `$AWS_SSO_CONFIG_PAGER` (if set)
2. `$PAGER` (if set)
3. `less -R` (if available) - Provides full interactive navigation with color support
4. `less` (if available) - Basic less functionality
5. `more` (if available) - Traditional more pager
6. `cat` (fallback) - Direct output

### Command Options
- `--force-paging`: Force paging even for short output (useful for testing)

## Usage Examples

```bash
# Normal usage - pages automatically if needed
aws-sso-config config list

# Disable paging
NO_PAGER=1 aws-sso-config config list

# Use a specific pager
PAGER="more" aws-sso-config config list

# Use tool-specific pager
AWS_SSO_CONFIG_PAGER="cat -n" aws-sso-config config list

# Force paging for testing
aws-sso-config config list --force-paging
```

## Implementation Details

### Interactive Terminal Connection
- **Temporary File Approach**: Content is written to a temporary file and opened in the pager
- **Full TTY Access**: Pager has complete control over the terminal for interactive features
- **Direct Connection**: Stdin, stdout, and stderr are connected directly to the terminal
- **Clean Cleanup**: Temporary files are automatically removed after paging

### Enhanced User Experience
- **Full Keyboard Support**: All standard pager keybindings work (arrows, space, search, etc.)
- **Color Preservation**: The `-R` flag preserves colors when available
- **Search Functionality**: Built-in search with '/' and '?' commands
- **Help Access**: 'h' key shows help in supported pagers

### Testing
- Comprehensive test suite for pager functionality
- Mockable UI interface for testing
- Environment variable testing
- Edge case handling

## Future Enhancements

When more configuration options are added to the tool, the paging system will automatically handle larger outputs. The current threshold is set to terminal height minus 2 lines for the prompt.

Potential future enhancements:
- Color support in pagers
- Custom paging thresholds
- Search functionality within paged output
- Mouse support for scroll wheels

## Backwards Compatibility

The paging implementation is fully backwards compatible:
- Existing scripts and automation continue to work unchanged
- No breaking changes to command-line interface
- Respects existing Unix conventions (NO_PAGER, PAGER variables)
- Graceful fallback when paging is not available or appropriate
