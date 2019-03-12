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
}

func NewTerminal(ctx context.Context, filter *Filter) (VirtualTerminal, error) {
	fd := int(ScreenFd)

	p, _ := cmd.GetPalette()
	env, _ := cmd.GetEnvironment()

	limit := *env.InteractiveShell.HistoryLimit
	historyFile, err := HistoryFilePath(env.InteractiveShell.HistoryFile)
	if err != nil {
		LogError(fmt.Sprintf("cannot detect filepath: %q", env.InteractiveShell.HistoryFile))
		limit = -1
	}

	prompt := NewPrompt(filter, p)
	completer := NewCompleter(filter)

	t, err := readline.NewEx(&readline.Config{
		HistoryFile:            historyFile,
		DisableAutoSaveHistory: true,
		HistoryLimit:           limit,
		HistorySearchFold:      true,
		Listener:               new(ReadlineListener),
		Stdin:                  Stdin,
		Stdout:                 Stdout,
		Stderr:                 Stderr,
	})
	if err != nil {
		return nil, err
	}

	terminal := ReadLineTerminal{
		terminal:  t,
		fd:        fd,
		prompt:    prompt,
		env:       env,
		completer: completer,
	}

	terminal.setCompleter()
	terminal.setKillWholeLine()
	terminal.setViMode()
	prompt.LoadConfig()

	terminal.SetPrompt(ctx)
	return terminal, nil
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

func (t ReadLineTerminal) SetPrompt(ctx context.Context) {
	str, err := t.prompt.RenderPrompt(ctx)
	if err != nil {
		LogError(err.Error())
	}
	t.terminal.SetPrompt(str)
}

func (t ReadLineTerminal) SetContinuousPrompt(ctx context.Context) {
	str, err := t.prompt.RenderContinuousPrompt(ctx)
	if err != nil {
		LogError(err.Error())
	}
	t.terminal.SetPrompt(str)
}

func (t ReadLineTerminal) SaveHistory(s string) {
	t.terminal.SaveHistory(s)
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
