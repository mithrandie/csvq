package file

import (
	"os"
	"strings"

	"github.com/mithrandie/go-file"
)

type LockFileContainer map[string]*LockFile

type LockFile struct {
	Path string
	Fp   *os.File
}

func (l *LockFile) Close() error {
	if err := file.Close(l.Fp); err != nil {
		return ParseError(err)
	}
	if err := os.Remove(l.Path); err != nil {
		return NewIOError(err.Error())
	}
	return nil
}

func (l LockFileContainer) Exists(path string) bool {
	ufpath := strings.ToUpper(path)
	_, ok := LockFiles[ufpath]
	return ok
}
