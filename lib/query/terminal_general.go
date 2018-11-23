// +build !darwin,!dragonfly,!freebsd,!linux,!netbsd,!openbsd,!solaris,!windows

package query

import (
	"io"

	"github.com/mithrandie/go-text/color"

	"golang.org/x/crypto/ssh/terminal"
)

type SSHTerminal struct {
	terminal  *terminal.Terminal
	stdin     int
	origState *terminal.State
	rawState  *terminal.State
	palette   *color.Palette
}

func NewTerminal() (VirtualTerminal, error) {
	stdin := int(ScreenFd)
	origState, err := terminal.MakeRaw(stdin)
	if err != nil {
		return nil, err
	}

	rawState, err := terminal.GetState(stdin)
	if err != nil {
		return nil, err
	}

	p, _ := GetPalette()

	return SSHTerminal{
		terminal:  terminal.NewTerminal(NewStdIO(), p.Render(PromptEffect, TerminalPrompt)),
		stdin:     stdin,
		origState: origState,
		rawState:  rawState,
		palette:   p,
	}, nil
}

func (t SSHTerminal) Teardown() {
	t.RestoreOriginalMode()
}

func (t SSHTerminal) RestoreRawMode() error {
	return terminal.Restore(t.stdin, t.rawState)
}

func (t SSHTerminal) RestoreOriginalMode() error {
	return terminal.Restore(t.stdin, t.origState)
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
	return terminal.GetSize(t.stdin)
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
		reader: Stdin,
		writer: Stdout,
	}
}
