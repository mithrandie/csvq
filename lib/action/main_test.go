package action

import (
	"os"
	"path"
	"testing"
)

var TestDir = path.Join(os.TempDir(), "csvq_action_test")
var TestDataDir string

func GetTestFilePath(filename string) string {
	return path.Join(TestDir, filename)
}

func TestMain(m *testing.M) {
	setup()
	r := m.Run()
	teardown()
	os.Exit(r)
}

func setup() {
	wdir, _ := os.Getwd()
	TestDataDir = path.Join(wdir, "..", "..", "testdata", "csv")

	if _, err := os.Stat(TestDir); os.IsNotExist(err) {
		os.Mkdir(TestDir, 0755)
	}

	r, _ := os.Open(path.Join(TestDataDir, "empty.txt"))
	os.Stdin = r
}

func teardown() {
	if _, err := os.Stat(TestDir); err == nil {
		os.RemoveAll(TestDir)
	}
}
