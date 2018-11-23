// +build darwin dragonfly freebsd linux netbsd openbsd solaris windows

package query

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/chzyer/readline"
	"github.com/mitchellh/go-homedir"
	"github.com/mithrandie/csvq/lib/cmd"
)

type ReadLineTerminal struct {
	terminal *readline.Instance
	fd       int
	prompt   *Prompt
}

func NewTerminal(filter *Filter) (VirtualTerminal, error) {
	fd := int(ScreenFd)

	p, _ := cmd.GetPalette()
	env, _ := cmd.GetEnvironment()

	limit := env.InteractiveShell.HistoryLimit
	historyFile, err := HistoryFilePath(env.InteractiveShell.HistoryFile)
	if err != nil {
		WriteToStderrWithLineBreak(fmt.Sprintf("cannot detect filepath: %q", env.InteractiveShell.HistoryFile))
		limit = -1
	}

	prompt := NewPrompt(filter, p)
	prompt.LoadConfig()

	t, err := readline.NewEx(&readline.Config{
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

	terminal := ReadLineTerminal{
		terminal: t,
		fd:       fd,
		prompt:   prompt,
	}

	terminal.SetPrompt()
	return terminal, nil
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
	str, err := t.prompt.RenderPrompt()
	if err != nil {
		WriteToStderrWithLineBreak(cmd.Error(err.Error()))
	}
	t.terminal.SetPrompt(str)
}

func (t ReadLineTerminal) SetContinuousPrompt() {
	str, err := t.prompt.RenderContinuousPrompt()
	if err != nil {
		WriteToStderrWithLineBreak(cmd.Error(err.Error()))
	}
	t.terminal.SetPrompt(str)
}

func (t ReadLineTerminal) SaveHistory(s string) {
	t.terminal.SaveHistory(s)
}

func (t ReadLineTerminal) GetSize() (int, int, error) {
	return readline.GetSize(t.fd)
}

func (t ReadLineTerminal) ReloadPromptConfig() error {
	return t.prompt.LoadConfig()
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
