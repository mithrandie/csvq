// +build windows

package query

import (
	"io"

	"github.com/mithrandie/readline-csvq"
)

func GetStdinForREPL() io.ReadCloser {
	return readline.NewRawReader()
}
