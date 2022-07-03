//go:build darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris || windows

package terminal

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mithrandie/csvq/lib/option"
	"github.com/mithrandie/csvq/lib/query"

	"github.com/mitchellh/go-homedir"
	"github.com/mithrandie/readline-csvq"
)

type ReadLineTerminal struct {
	terminal  *readline.Instance
	fd        int
	prompt    *Prompt
	env       *option.Environment
	completer *Completer
	tx        *query.Transaction
}

func NewTerminal(ctx context.Context, scope *query.ReferenceScope) (query.VirtualTerminal, error) {
	fd := int(scope.Tx.Session.ScreenFd())

	limit := *scope.Tx.Environment.InteractiveShell.HistoryLimit
	historyFile, err := HistoryFilePath(scope.Tx.Environment.InteractiveShell.HistoryFile)
	if err != nil {
		scope.Tx.LogWarn(fmt.Sprintf("cannot detect filepath: %q", scope.Tx.Environment.InteractiveShell.HistoryFile), false)
		limit = -1
	}

	prompt := NewPrompt(scope)
	completer := NewCompleter(scope)

	t, err := readline.NewEx(&readline.Config{
		HistoryFile:            historyFile,
		DisableAutoSaveHistory: true,
		HistoryLimit:           limit,
		HistorySearchFold:      true,
		Listener:               new(ReadlineListener),
		Stdin:                  readline.NewCancelableStdin(scope.Tx.Session.Stdin()),
		Stdout:                 scope.Tx.Session.Stdout(),
		Stderr:                 scope.Tx.Session.Stderr(),
	})
	if err != nil {
		return nil, err
	}

	terminal := ReadLineTerminal{
		terminal:  t,
		fd:        fd,
		prompt:    prompt,
		env:       scope.Tx.Environment,
		completer: completer,
		tx:        scope.Tx,
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
