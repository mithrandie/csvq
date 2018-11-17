// +build darwin dragonfly freebsd linux netbsd openbsd solaris windows

package cmd

import (
	"os"
	"path/filepath"

	"github.com/chzyer/readline"
	"github.com/mitchellh/go-homedir"
	"github.com/mithrandie/go-text/color"
)

const HistoryFile = ".csvq_history"

type ReadLineTerminal struct {
	terminal *readline.Instance
	fd       int
	palette  *color.Palette
}

func NewTerminal() (VirtualTerminal, error) {
	fd := int(os.Stdin.Fd())

	historyFile, err := historyFilePath()
	if err != nil {
		return nil, err
	}

	p, _ := GetPalette()

	t, err := readline.NewEx(&readline.Config{
		Prompt:                 p.Render(PromptEffect, TerminalPrompt),
		HistoryFile:            historyFile,
		DisableAutoSaveHistory: true,
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

func historyFilePath() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, HistoryFile), nil
}
