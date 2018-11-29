// +build !darwin,!dragonfly,!freebsd,!linux,!netbsd,!openbsd,!solaris,!windows

package query

import (
	"io"

	"github.com/mithrandie/csvq/lib/cmd"

	"golang.org/x/crypto/ssh/terminal"
)

type SSHTerminal struct {
	terminal  *terminal.Terminal
	stdin     int
	origState *terminal.State
	rawState  *terminal.State
	prompt    *Prompt
}

func NewTerminal(filter *Filter) (VirtualTerminal, error) {
	stdin := int(ScreenFd)
	origState, err := terminal.MakeRaw(stdin)
	if err != nil {
		return nil, err
	}

	rawState, err := terminal.GetState(stdin)
	if err != nil {
		return nil, err
	}

	p, _ := cmd.GetPalette()
	prompt := NewPrompt(filter, p)
	prompt.LoadConfig()

	t := SSHTerminal{
		terminal:  terminal.NewTerminal(NewStdIO(), p.Render(cmd.PromptEffect, TerminalPrompt)),
		stdin:     stdin,
		origState: origState,
		rawState:  rawState,
		prompt:    prompt,
	}

	t.RestoreOriginalMode()
	t.SetPrompt()
	t.RestoreRawMode()
	return t, nil
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
	if w, h, err := terminal.GetSize(t.stdin); err == nil {
		t.terminal.SetSize(w, h)
	}

	t.RestoreRawMode()
	s, e := t.terminal.ReadLine()
	t.RestoreOriginalMode()
	return s, e
}

func (t SSHTerminal) Write(s string) error {
	_, err := t.terminal.Write([]byte(s))
	return err
}

func (t SSHTerminal) WriteError(s string) error {
	_, err := Stderr.Write([]byte(s))
	return err
}

func (t SSHTerminal) SetPrompt() {
	str, err := t.prompt.RenderPrompt()
	if err != nil {
		LogError(err.Error())
	}
	t.terminal.SetPrompt(str)
}

func (t SSHTerminal) SetContinuousPrompt() {
	str, err := t.prompt.RenderContinuousPrompt()
	if err != nil {
		LogError(err.Error())
	}
	t.terminal.SetPrompt(str)
}

func (t SSHTerminal) SaveHistory(s string) {
	return
}

func (t SSHTerminal) GetSize() (int, int, error) {
	return terminal.GetSize(t.stdin)
}

func (t SSHTerminal) ReloadPromptConfig() error {
	return t.prompt.LoadConfig()
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
