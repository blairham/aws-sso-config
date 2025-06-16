package cli

import (
	"io"

	mcli "github.com/mitchellh/cli"
)

// UI implements the mitchellh/cli.Ui interface, while exposing the underlying
// io.Writer used for stdout and stderr.
type UI interface {
	mcli.Ui
	Stdout() io.Writer
	Stderr() io.Writer
}

// Command is an alias to reduce the diff. It can be removed at any time.
type Command mcli.Command

// BasicUI augments mitchellh/cli.BasicUi by exposing the underlying io.Writer.
type BasicUI struct {
	mcli.BasicUi
}

func (b *BasicUI) Stdout() io.Writer {
	return b.Writer
}

func (b *BasicUI) Stderr() io.Writer {
	return b.ErrorWriter
}
