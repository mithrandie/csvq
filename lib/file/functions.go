package file

import (
	"context"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func GetTimeoutContext(ctx context.Context) (context.Context, context.CancelFunc) {
	if _, ok := ctx.Deadline(); ok {
		return ctx, func() {
			// Dummy function
		}
	}

	return context.WithTimeout(context.Background(), WaitTimeout)
}

func UpdateWaitTimeout(waitTimeout float64, retryDelay time.Duration) error {
	d, err := time.ParseDuration(strconv.FormatFloat(waitTimeout, 'f', -1, 64) + "s")
	if err != nil {
		return err
	}

	WaitTimeout = d
	RetryDelay = retryDelay
	return nil
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
