//go:build !windows

package query

import (
	"io"
)

func GetStdinForREPL() io.ReadCloser {
	return stdin
}
