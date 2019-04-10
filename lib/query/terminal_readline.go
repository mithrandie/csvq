// +build darwin dragonfly freebsd linux netbsd openbsd solaris windows

package query

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mithrandie/csvq/lib/cmd"

	"github.com/mitchellh/go-homedir"
	"github.com/mithrandie/readline-csvq"
)

type ReadLineTerminal struct {
	terminal  *readline.Instance
	fd        int
	prompt    *Prompt
	env       *cmd.Environment
	completer *Completer
	tx        *Transaction
}

func NewTerminal(ctx context.Context, filter *Filter) (VirtualTerminal, error) {
	fd := int(filter.tx.Session.ScreenFd())

	limit := *filter.tx.Environment.InteractiveShell.HistoryLimit
	historyFile, err := HistoryFilePath(filter.tx.Environment.InteractiveShell.HistoryFile)
	if err != nil {
		filter.tx.LogWarn(fmt.Sprintf("cannot detect filepath: %q", filter.tx.Environment.InteractiveShell.HistoryFile), false)
		limit = -1
	}

	prompt := NewPrompt(filter)
	completer := NewCompleter(filter)

	t, err := readline.NewEx(&readline.Config{
		HistoryFile:            historyFile,
		DisableAutoSaveHistory: true,
		HistoryLimit:           limit,
		HistorySearchFold:      true,
		Listener:               new(ReadlineListener),
		Stdin:                  filter.tx.Session.Stdin(),
		Stdout:                 filter.tx.Session.Stdout(),
		Stderr:                 filter.tx.Session.Stderr(),
	})
	if err != nil {
		return nil, err
	}

	terminal := ReadLineTerminal{
		terminal:  t,
		fd:        fd,
		prompt:    prompt,
		env:       filter.tx.Environment,
		completer: completer,
		tx:        filter.tx,
	}

	terminal.setCompleter()
	terminal.setKillWholeLine()
	terminal.setViMode()
	if err = prompt.LoadConfig(); err != nil {
		return nil, err
	}

	terminal.SetPrompt(ctx)
	return terminal, nil
}

func (t ReadLineTerminal) Teardown() error {
	return t.terminal.Close()
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

func (t ReadLineTerminal) SetPrompt(ctx context.Context) {
	str, err := t.prompt.RenderPrompt(ctx)
	if err != nil {
		t.tx.LogError(err.Error())
	}
	t.terminal.SetPrompt(str)
}

func (t ReadLineTerminal) SetContinuousPrompt(ctx context.Context) {
	str, err := t.prompt.RenderContinuousPrompt(ctx)
	if err != nil {
		t.tx.LogError(err.Error())
	}
	t.terminal.SetPrompt(str)
}

func (t ReadLineTerminal) SaveHistory(s string) error {
	return t.terminal.SaveHistory(s)
}

func (t ReadLineTerminal) GetSize() (int, int, error) {
	return readline.GetSize(t.fd)
}

func (t ReadLineTerminal) ReloadConfig() error {
	t.setCompleter()
	t.setKillWholeLine()
	t.setViMode()
	return t.prompt.LoadConfig()
}

func (t ReadLineTerminal) UpdateCompleter() {
	if t.completer != nil {
		t.completer.Update()
	}
}

func (t ReadLineTerminal) setCompleter() {
	if *t.env.InteractiveShell.Completion {
		t.terminal.Config.AutoComplete = t.completer
	} else {
		t.terminal.Config.AutoComplete = nil
	}
}

func (t ReadLineTerminal) setKillWholeLine() {
	if *t.env.InteractiveShell.KillWholeLine {
		t.terminal.EnableKillWholeLine()
	} else {
		t.terminal.DisableKillWholeLine()
	}
}

func (t ReadLineTerminal) setViMode() {
	t.terminal.SetVimMode(*t.env.InteractiveShell.ViMode)
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
