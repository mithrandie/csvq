package cmd

import (
	"io"
	"os"
	"path/filepath"
	"testing"
)

var TestDir = filepath.Join(os.TempDir(), "csvq_cmd_test")
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
	wdir, _ := os.Getwd()
	TestDataDir = filepath.Join(wdir, "..", "..", "testdata", "csv")

	if _, err := os.Stat(TestDir); os.IsNotExist(err) {
		os.Mkdir(TestDir, 0755)
	}

	copyfile(filepath.Join(TestDir, "table1.csv"), filepath.Join(TestDataDir, "table1.csv"))
}

func teardown() {
	if _, err := os.Stat(TestDir); err == nil {
		os.RemoveAll(TestDir)
	}
}

func copyfile(dstfile string, srcfile string) error {
	src, err := os.Open(srcfile)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(dstfile)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}

	return nil
}
