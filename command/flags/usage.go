package flags

import (
	"bytes"
	"strings"

	"github.com/spf13/pflag"
)

func Usage(txt string, flags *pflag.FlagSet) string {
	u := &Usager{
		Usage: txt,
		Flags: flags,
	}

	return u.String()
}

type Usager struct {
	Usage string
	Flags *pflag.FlagSet
}

func (u *Usager) String() string {
	out := new(bytes.Buffer)
	out.WriteString(strings.TrimSpace(u.Usage))
	out.WriteString("\n")
	out.WriteString("\n")

	return strings.TrimRight(out.String(), "\n")
}
