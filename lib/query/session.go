package query

import (
	"bytes"
	"io"
	"os"

	"github.com/mithrandie/csvq/lib/cmd"
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

type Input struct {
	reader io.Reader
}

func NewInput(r io.Reader) *Input {
	return &Input{reader: r}
}

func (r *Input) Read(p []byte) (int, error) {
	return r.reader.Read(p)
}

func (r *Input) Close() error {
	if rc, ok := r.reader.(io.ReadCloser); ok {
		return rc.Close()
	}
	return nil
}

type Output struct {
	*bytes.Buffer
}

func NewOutput() *Output {
	return &Output{new(bytes.Buffer)}
}

func (w *Output) Close() error {
	return nil
}

type Session struct {
	ScreenFd uintptr
	Stdin    io.ReadCloser
	Stdout   io.WriteCloser
	Stderr   io.WriteCloser
	OutFile  io.Writer
	Terminal VirtualTerminal
}

func NewSession() *Session {
	return &Session{
		ScreenFd: os.Stdin.Fd(),
		Stdin:    os.Stdin,
		Stdout:   os.Stdout,
		Stderr:   os.Stderr,
		OutFile:  nil,
		Terminal: nil,
	}
}

func (sess *Session) CanReadStdin() bool {
	if sess.Stdin == nil {
		return false
	}
	if f, ok := sess.Stdin.(*os.File); ok {
		return cmd.IsReadableFromPipeOrRedirection(f)
	}
	return true
}

func (sess *Session) Log(log string, quiet bool) {
	if !quiet {
		if err := sess.WriteToStdoutWithLineBreak(log); err != nil {
			println(err.Error())
		}
	}
}

func (sess *Session) LogNotice(log string, quiet bool) {
	if !quiet {
		if err := sess.WriteToStdoutWithLineBreak(cmd.Notice(log)); err != nil {
			println(err.Error())
		}
	}
}

func (sess *Session) LogWarn(log string, quiet bool) {
	if !quiet {
		if err := sess.WriteToStdoutWithLineBreak(cmd.Warn(log)); err != nil {
			println(err.Error())
		}
	}
}

func (sess *Session) LogError(log string) {
	if err := sess.WriteToStderrWithLineBreak(cmd.Error(log)); err != nil {
		println(err.Error())
	}
}

func (sess *Session) WriteToStdout(s string) error {
	if sess.Terminal != nil {
		return sess.Terminal.Write(s)
	}

	_, err := sess.Stdout.Write([]byte(s))
	return err
}

func (sess *Session) WriteToStdoutWithLineBreak(s string) error {
	if 0 < len(s) && s[len(s)-1] != '\n' {
		s = s + "\n"
	}
	return sess.WriteToStdout(s)
}

func (sess *Session) WriteToStderr(s string) error {
	if sess.Terminal != nil {
		return sess.Terminal.WriteError(s)
	}

	_, err := sess.Stderr.Write([]byte(s))
	return err
}

func (sess *Session) WriteToStderrWithLineBreak(s string) error {
	if 0 < len(s) && s[len(s)-1] != '\n' {
		s = s + "\n"
	}
	return sess.WriteToStderr(s)
}
