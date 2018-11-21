package excmd

import (
	"os/exec"

	"github.com/mithrandie/csvq/lib/cmd"
)

func Exec(args []string) error {
	if len(args) < 1 {
		return nil
	}

	if cmd.Terminal != nil {
		cmd.Terminal.RestoreOriginalMode()
		defer cmd.Terminal.RestoreRawMode()
	}

	c := exec.Command(args[0], args[1:]...)
	c.Stdin = cmd.Stdin
	c.Stdout = cmd.Stdout
	c.Stderr = cmd.Stderr

	return c.Run()
}
