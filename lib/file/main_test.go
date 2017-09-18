package file

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func GetTestFilePath(filename string) string {
	return filepath.Join(TestDir, filename)
}

var TestDir = filepath.Join(os.TempDir(), "csvq_file_test")

func TestMain(m *testing.M) {
	os.Exit(run(m))
}

func run(m *testing.M) int {
	defer teardown()

	setup()
	return m.Run()
}

func setup() {
	if _, err := os.Stat(TestDir); err == nil {
		os.RemoveAll(TestDir)
	}

	if _, err := os.Stat(TestDir); os.IsNotExist(err) {
		os.Mkdir(TestDir, 0755)
	}

	fp, _ := os.Create(LockFilePath(GetTestFilePath("locked_by_other.txt")))
	fp.Close()

	fp, _ = os.Create(GetTestFilePath("open.txt"))
	fp.Close()

	fp, _ = os.Create(GetTestFilePath("update.txt"))
	fp.Close()

	UpdateWaitTimeout(0.1, 10*time.Millisecond)
	LockFiles = make(LockFileContainer)
}

func teardown() {
	UnlockAll()

	if _, err := os.Stat(TestDir); err == nil {
		os.RemoveAll(TestDir)
	}
}
