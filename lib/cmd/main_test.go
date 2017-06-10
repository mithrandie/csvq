package cmd

import (
	"os"
	"path"
	"testing"
)

var TestDir = path.Join(os.TempDir(), "csvq_cmd_test")

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
	if _, err := os.Stat(TestDir); os.IsNotExist(err) {
		os.Mkdir(TestDir, 0755)
	}
}

func teardown() {
	if _, err := os.Stat(TestDir); err == nil {
		os.RemoveAll(TestDir)
	}
}
