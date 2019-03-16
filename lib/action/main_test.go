package action

import (
	"os"
	"path/filepath"
	"testing"
)

var TestDir = filepath.Join(os.TempDir(), "csvq_action_test")
var TestDataDir string

func GetTestFilePath(filename string) string {
	return filepath.Join(TestDir, filename)
}

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
		_ = os.RemoveAll(TestDir)
	}

	wdir, _ := os.Getwd()
	TestDataDir = filepath.Join(wdir, "..", "..", "testdata", "csv")

	if _, err := os.Stat(TestDir); os.IsNotExist(err) {
		_ = os.Mkdir(TestDir, 0755)
	}

	r, _ := os.Open(filepath.Join(TestDataDir, "empty.txt"))
	os.Stdin = r
}

func teardown() {
	if _, err := os.Stat(TestDir); err == nil {
		_ = os.RemoveAll(TestDir)
	}
}
