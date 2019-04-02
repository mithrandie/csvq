package query

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/mithrandie/csvq/lib/file"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/value"

	"github.com/mitchellh/go-homedir"
	"github.com/mithrandie/go-text"
	"github.com/mithrandie/go-text/json"
)

type SyncMapStruct interface {
	SortedKeys() []string
	load(string) (interface{}, bool)
}

func SyncMapEqual(m1 SyncMapStruct, m2 SyncMapStruct) bool {
	mkeys := m1.SortedKeys()
	vlist := make([]interface{}, 0, len(mkeys))
	for _, key := range mkeys {
		v, _ := m1.load(key)
		vlist = append(vlist, v)
	}

	m2keys := m2.SortedKeys()
	vlist2 := make([]interface{}, 0, len(m2keys))
	for _, key := range m2keys {
		v, _ := m2.load(key)
		vlist2 = append(vlist2, v)
	}
	return reflect.DeepEqual(vlist, vlist2)
}

func SyncMapListEqual(l1 []SyncMapStruct, l2 []SyncMapStruct) bool {
	if len(l1) != len(l2) {
		return false
	}
	for i := 0; i < len(l1); i++ {
		if !SyncMapEqual(l1[i], l2[i]) {
			return false
		}
	}
	return true
}

func TempViewScopesToSyncMapList(scopes TemporaryViewScopes) []SyncMapStruct {
	l := make([]SyncMapStruct, 0, len(scopes))
	for _, v := range scopes {
		l = append(l, v)
	}
	return l
}
func VariableScopesToSyncMapList(scopes VariableScopes) []SyncMapStruct {
	l := make([]SyncMapStruct, 0, len(scopes))
	for _, v := range scopes {
		l = append(l, v)
	}
	return l
}

func GetTestFilePath(filename string) string {
	return filepath.Join(TestDir, filename)
}

var tempdir, _ = filepath.Abs(os.TempDir())
var TestDir = filepath.Join(tempdir, "csvq_query_test")
var TestDataDir string
var CompletionTestDir = filepath.Join(TestDir, "completion")
var CompletionTestSubDir = filepath.Join(TestDir, "completion", "sub")
var TestLocation = "UTC"
var NowForTest = time.Date(2012, 2, 3, 9, 18, 15, 0, GetTestLocation())
var HomeDir string

var TestTx, _ = NewTransaction(context.Background(), file.DefaultWaitTimeout, file.DefaultRetryDelay, NewSession())

func GetTestLocation() *time.Location {
	l, _ := time.LoadLocation(TestLocation)
	return l
}

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

func setup() {
	if _, err := os.Stat(TestDir); err == nil {
		_ = os.RemoveAll(TestDir)
	}

	cmd.TestTime = NowForTest

	TestDataDir = filepath.Join(GetWD(), "..", "..", "testdata", "csv")

	r, _ := os.Open(filepath.Join(TestDataDir, "empty.txt"))
	os.Stdin = r

	if _, err := os.Stat(TestDir); os.IsNotExist(err) {
		_ = os.Mkdir(TestDir, 0755)
	}

	_ = copyfile(filepath.Join(TestDir, "table_sjis.csv"), filepath.Join(TestDataDir, "table_sjis.csv"))
	_ = copyfile(filepath.Join(TestDir, "table_noheader.csv"), filepath.Join(TestDataDir, "table_noheader.csv"))
	_ = copyfile(filepath.Join(TestDir, "table_broken.csv"), filepath.Join(TestDataDir, "table_broken.csv"))
	_ = copyfile(filepath.Join(TestDir, "table1.csv"), filepath.Join(TestDataDir, "table1.csv"))
	_ = copyfile(filepath.Join(TestDir, "table1_bom.csv"), filepath.Join(TestDataDir, "table1_bom.csv"))
	_ = copyfile(filepath.Join(TestDir, "table1b.csv"), filepath.Join(TestDataDir, "table1b.csv"))
	_ = copyfile(filepath.Join(TestDir, "table2.csv"), filepath.Join(TestDataDir, "table2.csv"))
	_ = copyfile(filepath.Join(TestDir, "table4.csv"), filepath.Join(TestDataDir, "table4.csv"))
	_ = copyfile(filepath.Join(TestDir, "table5.csv"), filepath.Join(TestDataDir, "table5.csv"))
	_ = copyfile(filepath.Join(TestDir, "group_table.csv"), filepath.Join(TestDataDir, "group_table.csv"))
	_ = copyfile(filepath.Join(TestDir, "insert_query.csv"), filepath.Join(TestDataDir, "table1.csv"))
	_ = copyfile(filepath.Join(TestDir, "update_query.csv"), filepath.Join(TestDataDir, "table1.csv"))
	_ = copyfile(filepath.Join(TestDir, "delete_query.csv"), filepath.Join(TestDataDir, "table1.csv"))
	_ = copyfile(filepath.Join(TestDir, "add_columns.csv"), filepath.Join(TestDataDir, "table1.csv"))
	_ = copyfile(filepath.Join(TestDir, "drop_columns.csv"), filepath.Join(TestDataDir, "table1.csv"))
	_ = copyfile(filepath.Join(TestDir, "rename_column.csv"), filepath.Join(TestDataDir, "table1.csv"))
	_ = copyfile(filepath.Join(TestDir, "updated_file_1.csv"), filepath.Join(TestDataDir, "table1.csv"))
	_ = copyfile(filepath.Join(TestDir, "dup_name.csv"), filepath.Join(TestDataDir, "dup_name.csv"))

	_ = copyfile(filepath.Join(TestDir, "table3.tsv"), filepath.Join(TestDataDir, "table3.tsv"))
	_ = copyfile(filepath.Join(TestDir, "dup_name.tsv"), filepath.Join(TestDataDir, "dup_name.tsv"))

	_ = copyfile(filepath.Join(TestDir, "table.json"), filepath.Join(TestDataDir, "table.json"))
	_ = copyfile(filepath.Join(TestDir, "table_h.json"), filepath.Join(TestDataDir, "table_h.json"))
	_ = copyfile(filepath.Join(TestDir, "table_a.json"), filepath.Join(TestDataDir, "table_a.json"))

	_ = copyfile(filepath.Join(TestDir, "table6.ltsv"), filepath.Join(TestDataDir, "table6.ltsv"))
	_ = copyfile(filepath.Join(TestDir, "table6_bom.ltsv"), filepath.Join(TestDataDir, "table6_bom.ltsv"))

	_ = copyfile(filepath.Join(TestDir, "fixed_length.txt"), filepath.Join(TestDataDir, "fixed_length.txt"))
	_ = copyfile(filepath.Join(TestDir, "fixed_length_bom.txt"), filepath.Join(TestDataDir, "fixed_length_bom.txt"))
	_ = copyfile(filepath.Join(TestDir, "fixed_length_sl.txt"), filepath.Join(TestDataDir, "fixed_length_sl.txt"))

	_ = copyfile(filepath.Join(TestDir, "autoselect"), filepath.Join(TestDataDir, "autoselect"))

	_ = copyfile(filepath.Join(TestDir, "source.sql"), filepath.Join(filepath.Join(GetWD(), "..", "..", "testdata"), "source.sql"))
	_ = copyfile(filepath.Join(TestDir, "source_syntaxerror.sql"), filepath.Join(filepath.Join(GetWD(), "..", "..", "testdata"), "source_syntaxerror.sql"))

	_ = os.Setenv("CSVQ_TEST_ENV", "foo")

	if _, err := os.Stat(CompletionTestDir); os.IsNotExist(err) {
		_ = os.Mkdir(CompletionTestDir, 0755)
	}
	if _, err := os.Stat(CompletionTestSubDir); os.IsNotExist(err) {
		_ = os.Mkdir(CompletionTestSubDir, 0755)
	}
	_ = copyfile(filepath.Join(CompletionTestDir, "table1.csv"), filepath.Join(TestDataDir, "table1.csv"))
	_ = copyfile(filepath.Join(CompletionTestDir, ".table1.csv"), filepath.Join(TestDataDir, "table1.csv"))
	_ = copyfile(filepath.Join(CompletionTestDir, "source.sql"), filepath.Join(filepath.Join(GetWD(), "..", "..", "testdata"), "source.sql"))
	_ = copyfile(filepath.Join(CompletionTestSubDir, "table2.csv"), filepath.Join(TestDataDir, "table2.csv"))

	Version = "v1.0.0"
	HomeDir, _ = homedir.Dir()
	TestTx.Session.Stdout = NewDiscard()
	TestTx.Session.Stderr = NewDiscard()
}

func teardown() {
	if _, err := os.Stat(TestDir); err == nil {
		_ = os.RemoveAll(TestDir)
	}
}

func initFlag(flags *cmd.Flags) {
	cpu := runtime.NumCPU() / 2
	if cpu < 1 {
		cpu = 1
	}

	flags.Repository = "."
	flags.Location = TestLocation
	flags.DatetimeFormat = []string{}
	flags.WaitTimeout = 15
	flags.ImportFormat = cmd.CSV
	flags.Delimiter = ','
	flags.DelimiterPositions = nil
	flags.SingleLine = false
	flags.JsonQuery = ""
	flags.Encoding = text.UTF8
	flags.NoHeader = false
	flags.WithoutNull = false
	flags.Format = cmd.TEXT
	flags.WriteEncoding = text.UTF8
	flags.WriteDelimiter = ','
	flags.WriteDelimiterPositions = nil
	flags.WriteAsSingleLine = false
	flags.WithoutHeader = false
	flags.LineBreak = text.LF
	flags.EncloseAll = false
	flags.JsonEscape = json.Backslash
	flags.PrettyPrint = false
	flags.EastAsianEncoding = false
	flags.CountDiacriticalSign = false
	flags.CountFormatCode = false
	flags.Quiet = false
	flags.CPU = cpu
	flags.Stats = false
	flags.SetColor(false)
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

func GenerateVariableMap(values map[string]value.Primary) VariableMap {
	m := NewVariableMap()
	for k, v := range values {
		m.Store(k, v)
	}
	return m
}

func GenerateViewMap(values []*View) ViewMap {
	m := NewViewMap()
	for _, v := range values {
		m.Store(v.FileInfo.Path, v)
	}
	return m
}

func GenerateStatementMap(values []*PreparedStatement) PreparedStatementMap {
	m := NewPreparedStatementMap()
	for _, v := range values {
		m.Store(v.Name, v)
	}
	return m
}

func GenerateBenchGroupedViewFilter() Filter {
	primaries := make([]value.Primary, 10000)
	for i := 0; i < 10000; i++ {
		primaries[i] = value.NewInteger(int64(i))
	}

	view := &View{
		Header: NewHeader("table1", []string{"c1"}),
		RecordSet: []Record{
			{
				NewGroupCell(primaries),
			},
		},
		isGrouped: true,
	}

	tx, _ := NewTransaction(context.Background(), file.DefaultWaitTimeout, file.DefaultRetryDelay, NewSession())

	return Filter{
		records: []filterRecord{
			{view: view},
		},
		tx: tx,
	}
}

func GenerateBenchView(tableName string, records int) *View {
	view := &View{
		Header:    NewHeader(tableName, []string{"c1"}),
		RecordSet: make([]Record, records),
	}

	for i := 0; i < records; i++ {
		view.RecordSet[i] = NewRecord([]value.Primary{value.NewInteger(int64(i))})
	}

	return view
}
