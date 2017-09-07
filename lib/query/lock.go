package query

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mithrandie/csvq/lib/parser"
)

type FileLockContainer struct {
	FileMap       map[string]string
	WaitTimeout   float64 //Seconds
	RetryInterval time.Duration
}

func NewFileLockContainer() *FileLockContainer {
	return &FileLockContainer{
		FileMap:       make(map[string]string),
		WaitTimeout:   0,
		RetryInterval: time.Millisecond,
	}
}

func (c FileLockContainer) LockedByOtherProcess(path string) bool {
	if c.HasLock(path) {
		return false
	}
	_, err := os.Stat(FileLockPath(path))
	return err == nil
}

func (c FileLockContainer) CanRead(path string) bool {
	start := time.Now()

	for {
		if !c.LockedByOtherProcess(path) {
			break
		}

		if time.Since(start).Seconds() > c.WaitTimeout {
			return false
		}
		time.Sleep(c.RetryInterval)
	}

	return true
}

func (c FileLockContainer) HasLock(path string) bool {
	ufpath := strings.ToUpper(path)
	_, ok := c.FileMap[ufpath]
	return ok
}

func (c FileLockContainer) LockWithTimeout(ident parser.Identifier, path string) error {
	var err error
	start := time.Now()

	for {
		if err = c.TryLock(path); err == nil {
			break
		}

		if time.Since(start).Seconds() > c.WaitTimeout {
			return NewFileLockTimeoutError(ident, path)
		}
		time.Sleep(c.RetryInterval)
	}

	return nil
}

func (c FileLockContainer) TryLock(path string) error {
	lockFile := FileLockPath(path)
	fp, err := os.OpenFile(lockFile, os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return err
	}
	fp.Close()

	c.FileMap[strings.ToUpper(path)] = lockFile
	return nil
}

func (c FileLockContainer) Unlock(path string) error {
	if !c.HasLock(path) {
		return nil
	}
	ufpath := strings.ToUpper(path)
	err := os.Remove(c.FileMap[ufpath])
	delete(c.FileMap, ufpath)
	return err
}

func (c FileLockContainer) UnlockAll() error {
	for key := range c.FileMap {
		if err := c.Unlock(key); err != nil {
			return err
		}
	}
	return nil
}

func FileLockPath(path string) string {
	dir := filepath.Dir(path)
	basename := filepath.Base(path)
	return filepath.Join(dir, "."+basename+".lock")
}
