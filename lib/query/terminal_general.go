//go:build !darwin && !dragonfly && !freebsd && !linux && !netbsd && !openbsd && !solaris && !windows

package query

import (
	"context"
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
	tx        *Transaction
}

func NewTerminal(ctx context.Context, scope *ReferenceScope) (VirtualTerminal, error) {
	stdin := int(scope.Tx.Session.ScreenFd())
	origState, err := terminal.MakeRaw(stdin)
	if err != nil {
		return nil, err
	}

	rawState, err := terminal.GetState(stdin)
	if err != nil {
		return nil, err
	}

	prompt := NewPrompt(scope)
	if err = prompt.LoadConfig(); err != nil {
		return nil, err
	}

	t := SSHTerminal{
		terminal:  terminal.NewTerminal(NewStdIO(scope.Tx.Session), scope.Tx.Palette.Render(cmd.PromptEffect, TerminalPrompt)),
		stdin:     stdin,
		origState: origState,
		rawState:  rawState,
		prompt:    prompt,
		tx:        scope.Tx,
	}

	_ = t.RestoreOriginalMode()
	t.SetPrompt(ctx)
	_ = t.RestoreRawMode()
	return t, nil
}

func (t SSHTerminal) Teardown() error {
	return t.RestoreOriginalMode()
}

func (t SSHTerminal) RestoreRawMode() error {
	return terminal.Restore(t.stdin, t.rawState)
}

func (t SSHTerminal) RestoreOriginalMode() error {
	return terminal.Restore(t.stdin, t.origState)
}

func (t SSHTerminal) ReadLine() (string, error) {
	if w, h, err := terminal.GetSize(t.stdin); err == nil {
		if err = t.terminal.SetSize(w, h); err != nil {
			return "", err
		}
	}

	_ = t.RestoreRawMode()
	s, err := t.terminal.ReadLine()
	_ = t.RestoreOriginalMode()
	return s, err
}

func (t SSHTerminal) Write(s string) error {
	_, err := t.terminal.Write([]byte(s))
	return err
}

func (t SSHTerminal) WriteError(s string) error {
	_, err := t.tx.Session.stderr.Write([]byte(s))
	return err
}

func (t SSHTerminal) SetPrompt(ctx context.Context) {
	str, err := t.prompt.RenderPrompt(ctx)
	if err != nil {
		t.tx.LogError(err.Error())
	}
	t.terminal.SetPrompt(str)
}

func (t SSHTerminal) SetContinuousPrompt(ctx context.Context) {
	str, err := t.prompt.RenderContinuousPrompt(ctx)
	if err != nil {
		t.tx.LogError(err.Error())
	}
	t.terminal.SetPrompt(str)
}

func (t SSHTerminal) SaveHistory(s string) error {
	return nil
}

func (t SSHTerminal) GetSize() (int, int, error) {
	return terminal.GetSize(t.stdin)
}

func (t SSHTerminal) ReloadConfig() error {
	return t.prompt.LoadConfig()
}

type StdIO struct {
	reader io.Reader
	writer io.Writer
}

func (t SSHTerminal) UpdateCompleter() {
	//Do Nothing
}

func (sh *StdIO) Read(p []byte) (n int, err error) {
	return sh.reader.Read(p)
}

func (sh *StdIO) Write(p []byte) (n int, err error) {
	return sh.writer.Write(p)
}

func NewStdIO(sess *Session) *StdIO {
	return &StdIO{
		reader: sess.stdin,
		writer: sess.stdout,
	}
}
