package cmd

import (
	"bufio"
	"io"
	"os"

	"github.com/mithrandie/go-file"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	TERMINAL_PROMPT            string = "csvq > "
	TERMINAL_CONTINUOUS_PROMPT string = "     > "
)

var Term *Terminal

func ToStdout(s string) error {
	if Term != nil {
		return Term.Write(s)
	}
	return CreateFile("", s)
}

func CreateFile(filename string, s string) error {
	var fp *os.File
	var err error

	if len(filename) < 1 {
		fp = os.Stdout
	} else {
		fp, err = file.Create(filename)
		if err != nil {
			return err
		}
		defer file.Close(fp)
	}

	w := bufio.NewWriter(fp)
	_, err = w.WriteString(s)
	if err != nil {
		return err
	}
	w.Flush()

	return nil
}

func UpdateFile(fp *os.File, s string) error {
	defer file.Close(fp)

	fp.Truncate(0)
	fp.Seek(0, 0)

	w := bufio.NewWriter(fp)
	if _, err := w.WriteString(s); err != nil {
		return err
	}
	w.Flush()

	return nil
}

func TryCreateFile(filename string) error {
	fp, err := os.OpenFile(filename, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}

	fp.Close()
	os.Remove(filename)

	return nil
}

type Terminal struct {
	terminal *terminal.Terminal
	oldFd    int
	oldState *terminal.State
}

func NewTerminal() (*Terminal, error) {
	oldFd := int(os.Stdin.Fd())
	oldState, err := terminal.MakeRaw(oldFd)
	if err != nil {
		return nil, err
	}

	return &Terminal{
		terminal: terminal.NewTerminal(NewStdIO(), TERMINAL_PROMPT),
		oldFd:    oldFd,
		oldState: oldState,
	}, nil
}

func (t *Terminal) Restore() {
	terminal.Restore(t.oldFd, t.oldState)
}

func (t *Terminal) ReadLine() (string, error) {
	return t.terminal.ReadLine()
}

func (t *Terminal) Write(s string) error {
	_, err := t.terminal.Write([]byte(s))
	return err
}

func (t *Terminal) SetPrompt() {
	t.terminal.SetPrompt(TERMINAL_PROMPT)
}

func (t *Terminal) SetContinuousPrompt() {
	t.terminal.SetPrompt(TERMINAL_CONTINUOUS_PROMPT)
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
		reader: os.Stdin,
		writer: os.Stdout,
	}
}
