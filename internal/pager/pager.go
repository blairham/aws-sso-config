package pager

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/mitchellh/cli"
)

// Pager handles pagination of output
type Pager struct {
	ui           cli.Ui
	threshold    int  // Number of lines before paging kicks in
	forceEnabled bool // Force paging even if output fits on screen
}

// New creates a new pager with default settings
func New(ui cli.Ui) *Pager {
	return &Pager{
		ui:        ui,
		threshold: getTerminalHeight() - 2, // Leave some space for prompt
	}
}

// NewWithThreshold creates a new pager with a custom threshold
func NewWithThreshold(ui cli.Ui, threshold int) *Pager {
	return &Pager{
		ui:        ui,
		threshold: threshold,
	}
}

// SetForceEnabled forces paging to be enabled regardless of output size
func (p *Pager) SetForceEnabled(force bool) {
	p.forceEnabled = force
}

// Output handles output with automatic paging
func (p *Pager) Output(lines []string) {
	if !p.shouldPage(lines) {
		// Output directly without paging
		for _, line := range lines {
			p.ui.Output(line)
		}
		return
	}

	// Use paging
	if err := p.pageOutput(lines); err != nil {
		// Fallback to direct output if paging fails
		p.ui.Error(fmt.Sprintf("Paging failed: %v", err))
		for _, line := range lines {
			p.ui.Output(line)
		}
	}
}

// shouldPage determines if output should be paged
func (p *Pager) shouldPage(lines []string) bool {
	if p.forceEnabled {
		return true
	}

	// Don't page if output is small
	if len(lines) <= p.threshold {
		return false
	}

	// Don't page if we're not in a terminal
	if !isTerminal() {
		return false
	}

	// Don't page if NO_PAGER environment variable is set
	if os.Getenv("NO_PAGER") != "" {
		return false
	}

	return true
}

// pageOutput sends output through a pager using temporary file approach for better interactivity
func (p *Pager) pageOutput(lines []string) error {
	pagerCmd := getPagerCommand()
	if pagerCmd == "" {
		return fmt.Errorf("no suitable pager found")
	}

	// Create a temporary file to hold the content
	tmpFile, err := os.CreateTemp("", "aws-sso-config-*.txt")
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer os.Remove(tmpFile.Name()) // Clean up

	// Write content to the temporary file
	writer := bufio.NewWriter(tmpFile)
	for _, line := range lines {
		if _, err := writer.WriteString(line + "\n"); err != nil {
			tmpFile.Close()
			return fmt.Errorf("failed to write to temporary file: %w", err)
		}
	}
	writer.Flush()
	tmpFile.Close()

	// Split the pager command into command and args
	cmdParts := strings.Fields(pagerCmd)
	if len(cmdParts) == 0 {
		return fmt.Errorf("invalid pager command")
	}

	// Validate the pager command for security
	if !isValidPagerCommand(pagerCmd) {
		return fmt.Errorf("invalid pager command")
	}

	// Add the temporary file as an argument
	args := append(cmdParts[1:], tmpFile.Name())
	// #nosec G204 - pagerCmd is validated against allowlist in isValidPagerCommand
	cmd := exec.Command(cmdParts[0], args...)

	// Connect the pager directly to the terminal for full interactivity
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the pager and wait for it to finish
	return cmd.Run()
}

// getPagerCommand returns the preferred pager command with security validation
func getPagerCommand() string {
	// Check environment variables in order of preference, but validate them
	if pager := os.Getenv("PAGER"); pager != "" {
		if isValidPagerCommand(pager) {
			return pager
		}
	}

	if pager := os.Getenv("AWS_SSO_CONFIG_PAGER"); pager != "" {
		if isValidPagerCommand(pager) {
			return pager
		}
	}

	// Try common pagers in order of preference
	pagers := []string{
		"less -R", // -R: raw control chars for colors, allows full interactivity
		"less",    // Basic less
		"more",    // Traditional more
		"cat",     // Fallback
	}

	for _, pager := range pagers {
		cmdParts := strings.Fields(pager)
		if len(cmdParts) > 0 {
			if _, err := exec.LookPath(cmdParts[0]); err == nil {
				return pager
			}
		}
	}

	return ""
}

// isValidPagerCommand validates that a pager command is safe to execute
func isValidPagerCommand(pagerCmd string) bool {
	// Parse the command
	cmdParts := strings.Fields(pagerCmd)
	if len(cmdParts) == 0 {
		return false
	}

	// Check that the base command exists and is executable
	baseProg := cmdParts[0]
	if _, err := exec.LookPath(baseProg); err != nil {
		return false
	}

	// Define allowed pager programs for security
	allowedPagers := map[string]bool{
		"less":  true,
		"more":  true,
		"most":  true,
		"cat":   true,
		"bat":   true,
		"pg":    true,
		"pager": true,
	}

	// Extract just the program name (remove path)
	progName := baseProg
	if lastSlash := strings.LastIndex(baseProg, "/"); lastSlash >= 0 {
		progName = baseProg[lastSlash+1:]
	}

	// Only allow known safe pager programs
	if !allowedPagers[progName] {
		return false
	}

	// Validate arguments don't contain dangerous characters
	for _, arg := range cmdParts[1:] {
		if containsDangerousChars(arg) {
			return false
		}
	}

	return true
}

// containsDangerousChars checks if a string contains characters that could be used for command injection
func containsDangerousChars(s string) bool {
	dangerous := []string{
		";", "&", "|", "`", "$", "(", ")", "{", "}",
		"<", ">", "!", "*", "?", "[", "]", "~", "^",
	}

	for _, char := range dangerous {
		if strings.Contains(s, char) {
			return true
		}
	}

	return false
}

// getTerminalHeight returns the height of the terminal
func getTerminalHeight() int {
	// Try to get terminal size from environment first
	if lines := os.Getenv("LINES"); lines != "" {
		if height, err := strconv.Atoi(lines); err == nil && height > 0 {
			return height
		}
	}

	// Try to get terminal size using system calls
	if runtime.GOOS != "windows" {
		if height := getUnixTerminalHeight(); height > 0 {
			return height
		}
	}

	// Default fallback
	return 24
}

// getUnixTerminalHeight gets terminal height on Unix systems
func getUnixTerminalHeight() int {
	// Try using stty
	if cmd := exec.Command("stty", "size"); cmd.Err == nil {
		if output, err := cmd.Output(); err == nil {
			parts := strings.Fields(string(output))
			if len(parts) >= 2 {
				if height, err := strconv.Atoi(parts[0]); err == nil {
					return height
				}
			}
		}
	}

	// Try using tput
	if cmd := exec.Command("tput", "lines"); cmd.Err == nil {
		if output, err := cmd.Output(); err == nil {
			if height, err := strconv.Atoi(strings.TrimSpace(string(output))); err == nil {
				return height
			}
		}
	}

	return 0
}

// isTerminal checks if we're running in a terminal
func isTerminal() bool {
	// On Unix-like systems, check if it's a character device
	if stat, err := os.Stdout.Stat(); err == nil {
		return (stat.Mode() & os.ModeCharDevice) != 0
	}

	return false
}
