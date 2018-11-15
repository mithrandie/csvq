package file

import (
	"os"
	"path/filepath"
	"time"

	"github.com/mithrandie/go-file"
)

func UpdateWaitTimeout(waitTimeout float64, retryInterval time.Duration) {
	WaitTimeout = waitTimeout
	RetryInterval = retryInterval
	file.WaitTimeout = waitTimeout
	file.RetryInterval = retryInterval
}

func LockFilePath(path string) string {
	dir := filepath.Dir(path)
	basename := filepath.Base(path)
	return filepath.Join(dir, "."+basename+LockFileSuffix)
}

func TempFilePath(path string) string {
	dir := filepath.Dir(path)
	basename := filepath.Base(path)
	return filepath.Join(dir, "."+basename+TempFileSuffix)
}

func Exists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}
