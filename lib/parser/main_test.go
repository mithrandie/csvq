package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func GetTestFilePath(filename string) string {
	return filepath.Join(TestDir, filename)
}

var TestDir = filepath.Join(os.TempDir(), "csvq_parser_test")

func TestMain(m *testing.M) {
	os.Exit(run(m))
}

func run(m *testing.M) int {
	defer teardown()

	setup()
	return m.Run()
}

func setup() {
	fp, _ := os.Create(GetTestFilePath("dummy.sql"))
	defer fp.Close()
}

func teardown() {
	if _, err := os.Stat(TestDir); err == nil {
		os.RemoveAll(TestDir)
	}
}
