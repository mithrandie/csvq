package query

import (
	"io"
	"os"

	"github.com/mithrandie/csvq/lib/cmd"
)

var (
	ScreenFd                = os.Stdin.Fd()
	Stdin    io.ReadCloser  = os.Stdin
	Stdout   io.WriteCloser = os.Stdout
	Stderr   io.WriteCloser = os.Stderr
	OutFile  io.Writer
	Terminal VirtualTerminal
)

type Discard struct {
}

func NewDiscard() *Discard {
	return &Discard{}
}

func (d Discard) Write(p []byte) (int, error) {
	return len(p), nil
}

func (d Discard) Close() error {
	return nil
}

func Log(log string, quiet bool) {
	if !quiet {
		if err := WriteToStdoutWithLineBreak(log); err != nil {
			println(err.Error())
		}
	}
}

func LogNotice(log string, quiet bool) {
	if !quiet {
		if err := WriteToStdoutWithLineBreak(cmd.Notice(log)); err != nil {
			println(err.Error())
		}
	}
}

func LogWarn(log string, quiet bool) {
	if !quiet {
		if err := WriteToStdoutWithLineBreak(cmd.Warn(log)); err != nil {
			println(err.Error())
		}
	}
}

func LogError(log string) {
	if err := WriteToStderrWithLineBreak(cmd.Error(log)); err != nil {
		println(err.Error())
	}
}

func WriteToStdout(s string) error {
	if Terminal != nil {
		return Terminal.Write(s)
	}

	_, err := Stdout.Write([]byte(s))
	return err
}

func WriteToStdoutWithLineBreak(s string) error {
	if 0 < len(s) && s[len(s)-1] != '\n' {
		s = s + "\n"
	}
	return WriteToStdout(s)
}

func WriteToStderr(s string) error {
	if Terminal != nil {
		return Terminal.WriteError(s)
	}

	_, err := Stderr.Write([]byte(s))
	return err
}

func WriteToStderrWithLineBreak(s string) error {
	if 0 < len(s) && s[len(s)-1] != '\n' {
		s = s + "\n"
	}
	return WriteToStderr(s)
}
