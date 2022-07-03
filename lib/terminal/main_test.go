package terminal

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/mithrandie/csvq/lib/file"
	"github.com/mithrandie/csvq/lib/option"
	"github.com/mithrandie/csvq/lib/query"

	"github.com/mitchellh/go-homedir"
)

var tempdir, _ = filepath.Abs(os.TempDir())
var TestDir = filepath.Join(tempdir, "csvq_terminal_test")
var CompletionTestDir = filepath.Join(TestDir, "completion")
var CompletionTestSubDir = filepath.Join(TestDir, "completion", "sub")
var TestLocation = "UTC"
var TestTx, _ = query.NewTransaction(context.Background(), file.DefaultWaitTimeout, file.DefaultRetryDelay, query.NewSession())
var TestDataDir string
var HomeDir string

func GetWD() string {
	wdir, _ := os.Getwd()
	return wdir
}

func TestMain(m *testing.M) {
	os.Exit(run(m))
}

func run(m *testing.M) int {
	defer teardown()

	setup()
	return m.Run()
}

func initFlag(flags *option.Flags) {
	cpu := runtime.NumCPU() / 2
	if cpu < 1 {
		cpu = 1
	}

	flags.Repository = "."
	flags.DatetimeFormat = []string{}
	flags.AnsiQuotes = false
	flags.StrictEqual = false
	flags.WaitTimeout = 15
	flags.ImportOptions = option.NewImportOptions()
	flags.ExportOptions = option.NewExportOptions()
	flags.Quiet = false
	flags.LimitRecursion = 5
	flags.CPU = cpu
	flags.Stats = false
	flags.SetColor(false)
	_ = flags.SetLocation(TestLocation)
}

func copyfile(dstfile string, srcfile string) error {
	src, err := os.Open(srcfile)
	if err != nil {
		return err
	}
	defer func() { _ = src.Close() }()

	dst, err := os.Create(dstfile)
	if err != nil {
		return err
	}
	defer func() { _ = dst.Close() }()

	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}

	return nil
}

func setup() {
	if _, err := os.Stat(TestDir); err == nil {
		_ = os.RemoveAll(TestDir)
	}
	if _, err := os.Stat(TestDir); os.IsNotExist(err) {
		_ = os.Mkdir(TestDir, 0755)
	}
	if _, err := os.Stat(CompletionTestDir); os.IsNotExist(err) {
		_ = os.Mkdir(CompletionTestDir, 0755)
	}
	if _, err := os.Stat(CompletionTestSubDir); os.IsNotExist(err) {
		_ = os.Mkdir(CompletionTestSubDir, 0755)
	}

	HomeDir, _ = homedir.Dir()
	TestDataDir = filepath.Join(GetWD(), "..", "..", "testdata", "csv")

	_ = copyfile(filepath.Join(CompletionTestDir, "table1.csv"), filepath.Join(TestDataDir, "table1.csv"))
	_ = copyfile(filepath.Join(CompletionTestDir, ".table1.csv"), filepath.Join(TestDataDir, "table1.csv"))
	_ = copyfile(filepath.Join(CompletionTestDir, "source.sql"), filepath.Join(filepath.Join(GetWD(), "..", "..", "testdata"), "source.sql"))
	_ = copyfile(filepath.Join(CompletionTestSubDir, "table2.csv"), filepath.Join(TestDataDir, "table2.csv"))

	_ = os.Setenv("CSVQ_TEST_ENV", "foo")
	query.Version = "v1.0.0"
	TestTx.Session.SetStdout(query.NewDiscard())
	TestTx.Session.SetStderr(query.NewDiscard())
	initFlag(TestTx.Flags)
}

func teardown() {
	if _, err := os.Stat(TestDir); err == nil {
		_ = os.RemoveAll(TestDir)
	}
}
