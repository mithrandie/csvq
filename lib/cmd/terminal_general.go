// +build !darwin,!dragonfly,!freebsd,!linux,!netbsd,!openbsd,!solaris,!windows

package cmd

import (
	"io"
	"os"

	"github.com/mithrandie/go-text/color"

	"golang.org/x/crypto/ssh/terminal"
)

type SSHTerminal struct {
	terminal *terminal.Terminal
	oldFd    int
	oldState *terminal.State
	palette  *color.Palette
}

func NewTerminal() (VirtualTerminal, error) {
	oldFd := int(os.Stdin.Fd())
	oldState, err := terminal.MakeRaw(oldFd)
	if err != nil {
		return nil, err
	}

	p := GetPalette()

	return SSHTerminal{
		terminal: terminal.NewTerminal(NewStdIO(), p.Render(PromptEffect, TerminalPrompt)),
		oldFd:    oldFd,
		oldState: oldState,
		palette:  p,
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

func (t SSHTerminal) WriteError(s string) error {
	_, err := t.terminal.Write([]byte(s))
	return err
}

func (t SSHTerminal) SetPrompt() {
	t.terminal.SetPrompt(t.palette.Render(PromptEffect, TerminalPrompt))
}

func (t SSHTerminal) SetContinuousPrompt() {
	t.terminal.SetPrompt(t.palette.Render(PromptEffect, TerminalContinuousPrompt))
}

func (t SSHTerminal) SaveHistory(s string) {
	return
}

func (t SSHTerminal) GetSize() (int, int, error) {
	return terminal.GetSize(t.oldFd)
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
