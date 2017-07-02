package query

import (
	"io"
	"os"
	"path"
	"testing"
	"time"

	"github.com/mithrandie/csvq/lib/cmd"
)

func GetTestFilePath(filename string) string {
	return path.Join(TestDir, filename)
}

var TestDir = path.Join(os.TempDir(), "csvq_query_test")
var TestDataDir string
var TestLocation = "America/Los_Angeles"

func GetTestLocation() *time.Location {
	l, _ := time.LoadLocation(TestLocation)
	return l
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
	flags := cmd.GetFlags()
	flags.Location = TestLocation
	flags.Now = "2012-02-03 09:18:15"

	wdir, _ := os.Getwd()
	TestDataDir = path.Join(wdir, "..", "..", "testdata", "csv")

	r, _ := os.Open(path.Join(TestDataDir, "empty.txt"))
	os.Stdin = r

	if _, err := os.Stat(TestDir); os.IsNotExist(err) {
		os.Mkdir(TestDir, 0755)
	}

	copyfile(path.Join(TestDir, "table1.csv"), path.Join(TestDataDir, "table1.csv"))
	copyfile(path.Join(TestDir, "table2.csv"), path.Join(TestDataDir, "table2.csv"))
	copyfile(path.Join(TestDir, "table4.csv"), path.Join(TestDataDir, "table4.csv"))
	copyfile(path.Join(TestDir, "group_table.csv"), path.Join(TestDataDir, "group_table.csv"))
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
