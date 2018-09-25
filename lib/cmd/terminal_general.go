// +build !darwin,!dragonfly,!freebsd,!linux,!netbsd,!openbsd,!solaris,!windows

package cmd

import (
	"io"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

type SSHTerminal struct {
	terminal *terminal.Terminal
	oldFd    int
	oldState *terminal.State
}

func NewTerminal() (VirtualTerminal, error) {
	oldFd := int(os.Stdin.Fd())
	oldState, err := terminal.MakeRaw(oldFd)
	if err != nil {
		return nil, err
	}

	return SSHTerminal{
		terminal: terminal.NewTerminal(NewStdIO(), TerminalPrompt),
		oldFd:    oldFd,
		oldState: oldState,
	}, nil
}

func (t SSHTerminal) Teardown() {
	terminal.Restore(t.oldFd, t.oldState)
}

func (t SSHTerminal) ReadLine() (string, error) {
	return t.terminal.ReadLine()
}

func (t SSHTerminal) Write(s string) error {
	_, err := t.terminal.Write([]byte(s))
	return err
}

func (t SSHTerminal) SetPrompt() {
	t.terminal.SetPrompt(TerminalPrompt)
}

func (t SSHTerminal) SetContinuousPrompt() {
	t.terminal.SetPrompt(TerminalContinuousPrompt)
}

func (t SSHTerminal) SaveHistory(s string) {
	return
}

type StdIO struct {
	reader io.Reader
	writer io.Writer
}

func (sh *StdIO) Read(p []byte) (n int, err error) {
	return sh.reader.Read(p)
}

func (sh *StdIO) Write(p []byte) (n int, err error) {
	return sh.writer.Write(p)
}

func NewStdIO() *StdIO {
	return &StdIO{
		reader: os.Stdin,
		writer: os.Stdout,
	}
}
