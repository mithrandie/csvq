// +build windows

package query

import (
	"github.com/mithrandie/readline-csvq"
)

func init() {
	stdin = readline.NewRawReader()
}
