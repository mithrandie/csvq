// +build darwin dragonfly freebsd linux netbsd openbsd solaris windows

package cmd

import (
	"bufio"
	"os"
	"path/filepath"

	"github.com/chzyer/readline"
	"github.com/mitchellh/go-homedir"
)

const HISTORY_FILE = ".csvq_history"

type ReadLineTerminal struct {
	terminal *readline.Instance
}

func NewTerminal() (VirtualTerminal, error) {
	historyFile, err := historyFilePath()
	if err != nil {
		return nil, err
	}

	t, err := readline.NewEx(&readline.Config{
		Prompt:                 TERMINAL_PROMPT,
		HistoryFile:            historyFile,
		DisableAutoSaveHistory: true,
	})
	if err != nil {
		return nil, err
	}

	return ReadLineTerminal{
		terminal: t,
	}, nil
}

func (t ReadLineTerminal) Teardown() {
	t.terminal.Close()
}

func (t ReadLineTerminal) ReadLine() (string, error) {
	return t.terminal.Readline()
}

func (t ReadLineTerminal) Write(s string) error {
	w := bufio.NewWriter(os.Stdout)
	_, err := w.WriteString(s)
	if err != nil {
		return err
	}
	w.Flush()
	return nil
}

func (t ReadLineTerminal) SetPrompt() {
	t.terminal.SetPrompt(TERMINAL_PROMPT)
}

func (t ReadLineTerminal) SetContinuousPrompt() {
	t.terminal.SetPrompt(TERMINAL_CONTINUOUS_PROMPT)
}

func (t ReadLineTerminal) SaveHistory(s string) {
	t.terminal.SaveHistory(s)
}

func historyFilePath() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	home, err = homedir.Expand(home)
	if err != nil {
		return "", err
	}
	return filepath.Join(home, HISTORY_FILE), nil
}
