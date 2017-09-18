package file

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mithrandie/go-file"
)

var WaitTimeout float64 = 30.0
var RetryInterval time.Duration = 50 * time.Millisecond

var LockFiles = make(LockFileContainer)

func UpdateWaitTimeout(waitTimeout float64, retryInterval time.Duration) {
	WaitTimeout = waitTimeout
	RetryInterval = retryInterval
	file.WaitTimeout = waitTimeout
	file.RetryInterval = retryInterval
}

func OpenToRead(path string) (*os.File, error) {
	if !CanRead(path) {
		return nil, NewTimeoutError(path)
	}

	fp, err := file.OpenToReadWithTimeout(path)
	if err != nil {
		return nil, ParseError(err)
	}
	return fp, nil
}

func OpenToUpdate(path string) (*os.File, error) {
	if err := LockWithTimeout(path); err != nil {
		return nil, err
	}

	fp, err := file.OpenToUpdateWithTimeout(path)
	if err != nil {
		Unlock(path)
		return nil, ParseError(err)
	}
	return fp, nil
}

func Create(path string) (*os.File, error) {
	if !LockFiles.Exists(path) {
		if err := TryLock(path); err != nil {
			return nil, err
		}
	}

	fp, err := file.Create(path)
	if err != nil {
		Unlock(path)
		return nil, ParseError(err)
	}
	return fp, nil
}

func Close(fp *os.File) error {
	err := file.Close(fp)
	Unlock(fp.Name())
	return err
}

func Unlock(path string) {
	ufpath := strings.ToUpper(path)
	if lockFile, ok := LockFiles[ufpath]; ok {
		lockFile.Close()
		delete(LockFiles, ufpath)
	}
}

func UnlockAll() {
	for key, lockFile := range LockFiles {
		lockFile.Close()
		delete(LockFiles, key)
	}
}

func IsLockedByOtherProcess(path string) bool {
	if LockFiles.Exists(path) {
		return false
	}
	_, err := os.Stat(LockFilePath(path))
	return err == nil
}

func CanRead(path string) bool {
	var start time.Time

	for {
		if start.IsZero() {
			start = time.Now()
		} else if time.Since(start).Seconds() > WaitTimeout {
			return false
			break
		}

		if !IsLockedByOtherProcess(path) {
			break
		}

		time.Sleep(RetryInterval)
	}

	return true
}

func LockWithTimeout(path string) error {
	var err error
	var start time.Time

	for {
		if start.IsZero() {
			start = time.Now()
		} else if time.Since(start).Seconds() > WaitTimeout {
			err = NewTimeoutError(path)
			break
		}

		if err = TryLock(path); err == nil {
			break
		}
		time.Sleep(RetryInterval)
	}

	return err
}

func TryLock(path string) error {
	lockFilePath := LockFilePath(path)
	fp, err := file.Create(lockFilePath)
	if err != nil {
		return NewLockError(fmt.Sprintf("unable to create lock file for %q", path))
	}

	LockFiles[strings.ToUpper(path)] = &LockFile{
		Path: lockFilePath,
		Fp:   fp,
	}
	return nil
}

func LockFilePath(path string) string {
	dir := filepath.Dir(path)
	basename := filepath.Base(path)
	return filepath.Join(dir, "."+basename+".lock")
}
