package cmd

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/mitchellh/go-homedir"
)

var TestDir = filepath.Join(os.TempDir(), "csvq_cmd_test")
var TestDataDir string
var HomeDir string

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

	wdir, _ := os.Getwd()
	TestDataDir = filepath.Join(wdir, "..", "..", "testdata", "csv")

	HomeDir, _ = homedir.Dir()

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
