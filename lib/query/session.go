package query

import (
	"bytes"
	"io"
	"os"
	"sync"

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
	screenFd uintptr
	stdin    io.ReadCloser
	stdout   io.WriteCloser
	stderr   io.WriteCloser
	outFile  io.Writer
	terminal VirtualTerminal

	mtx *sync.Mutex
}

func NewSession() *Session {
	return &Session{
		screenFd: os.Stdin.Fd(),
		stdin:    os.Stdin,
		stdout:   os.Stdout,
		stderr:   os.Stderr,
		outFile:  nil,
		terminal: nil,

		mtx: &sync.Mutex{},
	}
}

func (sess *Session) ScreenFd() uintptr {
	return sess.screenFd
}

func (sess *Session) Stdin() io.ReadCloser {
	return sess.stdin
}

func (sess *Session) Stdout() io.WriteCloser {
	return sess.stdout
}

func (sess *Session) Stderr() io.WriteCloser {
	return sess.stderr
}

func (sess *Session) OutFile() io.Writer {
	return sess.outFile
}

func (sess *Session) Terminal() VirtualTerminal {
	return sess.terminal
}

func (sess *Session) SetStdin(r io.ReadCloser) {
	sess.mtx.Lock()
	sess.stdin = r
	sess.mtx.Unlock()
}

func (sess *Session) SetStdout(w io.WriteCloser) {
	sess.mtx.Lock()
	sess.stdout = w
	sess.mtx.Unlock()
}

func (sess *Session) SetStderr(w io.WriteCloser) {
	sess.mtx.Lock()
	sess.stderr = w
	sess.mtx.Unlock()
}

func (sess *Session) SetOutFile(w io.Writer) {
	sess.mtx.Lock()
	sess.outFile = w
	sess.mtx.Unlock()
}

func (sess *Session) SetTerminal(t VirtualTerminal) {
	sess.mtx.Lock()
	sess.terminal = t
	sess.mtx.Unlock()
}

func (sess *Session) CanReadStdin() bool {
	if sess.stdin == nil {
		return false
	}
	if f, ok := sess.stdin.(*os.File); ok {
		return cmd.IsReadableFromPipeOrRedirection(f)
	}
	return true
}

func (sess *Session) WriteToStdout(s string) (err error) {
	sess.mtx.Lock()
	if sess.terminal != nil {
		err = sess.terminal.Write(s)
	} else if sess.stdout != nil {
		_, err = sess.stdout.Write([]byte(s))
	}
	sess.mtx.Unlock()
	return
}

func (sess *Session) WriteToStdoutWithLineBreak(s string) error {
	if 0 < len(s) && s[len(s)-1] != '\n' {
		s = s + "\n"
	}
	return sess.WriteToStdout(s)
}

func (sess *Session) WriteToStderr(s string) (err error) {
	sess.mtx.Lock()
	if sess.terminal != nil {
		err = sess.terminal.WriteError(s)
	} else if sess.stderr != nil {
		_, err = sess.stderr.Write([]byte(s))
	}
	sess.mtx.Unlock()
	return
}

func (sess *Session) WriteToStderrWithLineBreak(s string) error {
	if 0 < len(s) && s[len(s)-1] != '\n' {
		s = s + "\n"
	}
	return sess.WriteToStderr(s)
}
