package query

import (
	"io"
	"os"
)

var (
	ScreenFd                = os.Stdin.Fd()
	Stdin    io.ReadCloser  = os.Stdin
	Stdout   io.WriteCloser = os.Stdout
	Stderr   io.WriteCloser = os.Stderr
	OutFile  io.Writer
	Terminal VirtualTerminal
)

func Log(log string, quiet bool) {
	if !quiet {
		WriteToStdoutWithLineBreak(log)
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
		return Terminal.Write(s)
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
