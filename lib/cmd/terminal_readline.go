// +build darwin dragonfly freebsd linux netbsd openbsd solaris windows

package cmd

import (
	"bufio"
	"github.com/mithrandie/csvq/lib/color"
	"os"
	"path/filepath"

	"github.com/chzyer/readline"
	"github.com/mitchellh/go-homedir"
)

const HistoryFile = ".csvq_history"

type ReadLineTerminal struct {
	terminal *readline.Instance
	fd       int
}

func NewTerminal() (VirtualTerminal, error) {
	fd := int(os.Stdin.Fd())

	historyFile, err := historyFilePath()
	if err != nil {
		return nil, err
	}

	t, err := readline.NewEx(&readline.Config{
		Prompt:                 color.Blue(TerminalPrompt),
		HistoryFile:            historyFile,
		DisableAutoSaveHistory: true,
	})
	if err != nil {
		return nil, err
	}

	return ReadLineTerminal{
		terminal: t,
		fd:       fd,
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
	t.terminal.SetPrompt(color.Blue(TerminalPrompt))
}

func (t ReadLineTerminal) SetContinuousPrompt() {
	t.terminal.SetPrompt(color.Blue(TerminalContinuousPrompt))
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
	home, err = homedir.Expand(home)
	if err != nil {
		return "", err
	}
	return filepath.Join(home, HistoryFile), nil
}
