package action

import (
	"io"
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

	copyfile(path.Join(TestDir, "insert_query.csv"), path.Join(TestDataDir, "table1.csv"))
	copyfile(path.Join(TestDir, "update_query.csv"), path.Join(TestDataDir, "table1.csv"))
	copyfile(path.Join(TestDir, "delete_query.csv"), path.Join(TestDataDir, "table1.csv"))
	copyfile(path.Join(TestDir, "add_columns.csv"), path.Join(TestDataDir, "table1.csv"))
	copyfile(path.Join(TestDir, "drop_columns.csv"), path.Join(TestDataDir, "table1.csv"))
	copyfile(path.Join(TestDir, "rename_column.csv"), path.Join(TestDataDir, "table1.csv"))
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
