// +build darwin dragonfly freebsd linux netbsd openbsd solaris windows

package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/chzyer/readline"
	"github.com/mitchellh/go-homedir"
	"github.com/mithrandie/go-text/color"
)

type ReadLineTerminal struct {
	terminal *readline.Instance
	fd       int
	palette  *color.Palette
}

func NewTerminal() (VirtualTerminal, error) {
	fd := int(os.Stdin.Fd())

	p, _ := GetPalette()
	env, _ := GetEnvironment()

	limit := env.InteractiveShell.HistoryLimit
	historyFile, err := HistoryFilePath(env.InteractiveShell.HistoryFile)
	if err != nil {
		WriteToStderrWithLineBreak(fmt.Sprintf("cannot detect filepath: %q", env.InteractiveShell.HistoryFile))
		limit = -1
	}

	t, err := readline.NewEx(&readline.Config{
		Prompt:                 p.Render(PromptEffect, TerminalPrompt),
		HistoryFile:            historyFile,
		DisableAutoSaveHistory: true,
		HistoryLimit:           limit,
		Stdin:                  Stdin,
		Stdout:                 Stdout,
		Stderr:                 Stderr,
	})
	if err != nil {
		return nil, err
	}

	return ReadLineTerminal{
		terminal: t,
		fd:       fd,
		palette:  p,
	}, nil
}

func (t ReadLineTerminal) Teardown() {
	t.terminal.Close()
}

func (t ReadLineTerminal) RestoreRawMode() error {
	return t.terminal.Terminal.EnterRawMode()
}

func (t ReadLineTerminal) RestoreOriginalMode() error {
	return t.terminal.Terminal.ExitRawMode()
}

func (t ReadLineTerminal) ReadLine() (string, error) {
	return t.terminal.Readline()
}

func (t ReadLineTerminal) Write(s string) error {
	_, err := t.terminal.Write([]byte(s))
	return err
}

func (t ReadLineTerminal) WriteError(s string) error {
	_, err := t.terminal.Stderr().Write([]byte(s))
	return err
}

func (t ReadLineTerminal) SetPrompt() {
	t.terminal.SetPrompt(t.palette.Render(PromptEffect, TerminalPrompt))
}

func (t ReadLineTerminal) SetContinuousPrompt() {
	t.terminal.SetPrompt(t.palette.Render(PromptEffect, TerminalContinuousPrompt))
}

func (t ReadLineTerminal) SaveHistory(s string) {
	t.terminal.SaveHistory(s)
}

func (t ReadLineTerminal) GetSize() (int, int, error) {
	return readline.GetSize(t.fd)
}

func HistoryFilePath(filename string) (string, error) {
	if filename[0] == '~' {
		if fpath, err := homedir.Expand(filename); err == nil {
			return fpath, nil
		}
	}

	fpath := os.ExpandEnv(filename)

	if filepath.IsAbs(fpath) {
		return fpath, nil
	}

	home, err := homedir.Dir()
	if err != nil {
		return filename, err
	}
	return filepath.Join(home, fpath), nil
}
