package cmd

import (
	"io"
	"os"
)

var Terminal VirtualTerminal

const (
	TerminalPrompt           string = "csvq > "
	TerminalContinuousPrompt string = "     > "
)

var (
	Stdin  io.ReadCloser  = os.Stdin
	Stdout io.WriteCloser = os.Stdout
	Stderr io.WriteCloser = os.Stderr
)

type VirtualTerminal interface {
	ReadLine() (string, error)
	Write(string) error
	WriteError(string) error
	SetPrompt()
	SetContinuousPrompt()
	SaveHistory(string)
	Teardown()
	GetSize() (int, int, error)
	RestoreRawMode() error
	RestoreOriginalMode() error
}
