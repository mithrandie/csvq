//go:build !windows

package terminal

import (
	"io"
	"os"
)

func GetStdinForREPL() io.ReadCloser {
	return os.Stdin
}
