package main

import (
	"fmt"
	"os"

	mcli "github.com/mitchellh/cli"

	"github.com/blairham/aws-sso-config/command"
	"github.com/blairham/aws-sso-config/command/cli"
)

// Version information - set by build flags
var (
	version   = "dev"
	commit    = "unknown"
	buildTime = "unknown"
)

// Variables for testing
var osExit = os.Exit

func main() {
	osExit(Run(os.Args[1:]))
}

// For testing purposes
var createCLI = func(ui *cli.BasicUI, args []string) *mcli.CLI {
	cmds := command.RegisteredCommands(ui)
	var names []string
	for c := range cmds {
		names = append(names, c)
	}

	return &mcli.CLI{
		Name:                       "aws-sso-config",
		Version:                    fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, buildTime),
		Args:                       args,
		Commands:                   cmds,
		Autocomplete:               true,
		AutocompleteNoDefaultFlags: true,
		HelpFunc:                   mcli.FilteredHelpFunc(names, mcli.BasicHelpFunc("aws-sso-config")),
		HelpWriter:                 os.Stdout,
	}
}

func Run(args []string) int {
	ui := &cli.BasicUI{
		BasicUi: mcli.BasicUi{
			Reader:      os.Stdin,
			Writer:      os.Stdout,
			ErrorWriter: os.Stderr,
		},
	}

	cliInstance := createCLI(ui, args)

	exitCode, err := cliInstance.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing CLI: %v\n", err)
		return 1
	}

	return exitCode
}
