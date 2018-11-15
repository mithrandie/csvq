package cmd

import (
	"bufio"
	"os"

	"github.com/mithrandie/go-text/color"
)

func WriteToStdout(s string) error {
	if Terminal != nil {
		return Terminal.Write(s)
	}

	w := bufio.NewWriter(os.Stdout)
	_, err := w.WriteString(s)
	if err != nil {
		return err
	}
	w.Flush()

	return nil
}

func WriteToStdErr(s string) error {
	if Terminal != nil {
		return Terminal.Write(s)
	}

	w := bufio.NewWriter(os.Stderr)
	_, err := w.WriteString(color.Error(s))
	if err != nil {
		return err
	}
	w.Flush()

	return nil
}
