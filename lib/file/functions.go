package file

import (
	"context"
	"os"
	"path/filepath"
	"time"
)

func GetTimeoutContext(ctx context.Context, waitTimeOut time.Duration) (context.Context, context.CancelFunc) {
	if _, ok := ctx.Deadline(); ok {
		return ctx, func() {
			// Dummy function
		}
	}

	return context.WithTimeout(context.Background(), waitTimeOut)
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
