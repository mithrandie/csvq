//go:build windows

package terminal

import (
	"io"

	"github.com/mithrandie/readline-csvq"
)

func GetStdinForREPL() io.ReadCloser {
	return readline.NewRawReader()
}
