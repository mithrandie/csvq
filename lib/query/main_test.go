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

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/file"
	"github.com/mithrandie/csvq/lib/value"

	"github.com/mitchellh/go-homedir"
)

type syncMapStruct interface {
	SortedKeys() []string
	LoadDirect(string) (interface{}, bool)
}

func SyncMapEqual(m1 syncMapStruct, m2 syncMapStruct) bool {
	mkeys := m1.SortedKeys()
	vlist := make([]interface{}, 0, len(mkeys))
	for _, key := range mkeys {
		v, _ := m1.LoadDirect(key)
		vlist = append(vlist, v)
	}

	m2keys := m2.SortedKeys()
	vlist2 := make([]interface{}, 0, len(m2keys))
	for _, key := range m2keys {
		v, _ := m2.LoadDirect(key)
		vlist2 = append(vlist2, v)
	}
	return reflect.DeepEqual(vlist, vlist2)
}

func BlockScopeListEqual(s1 []BlockScope, s2 []BlockScope) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := range s1 {
		if !SyncMapEqual(s1[i].variables, s2[i].variables) {
			return false
		}
		if !SyncMapEqual(s1[i].temporaryTables, s2[i].temporaryTables) {
			return false
		}
		if !SyncMapEqual(s1[i].cursors, s2[i].cursors) {
			return false
		}
		if !reflect.DeepEqual(s1[i].functions, s2[i].functions) {
			return false
		}
	}
	return true
}

func NodeScopeListEqual(s1 []NodeScope, s2 []NodeScope) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := range s1 {
		if !reflect.DeepEqual(s1[i].inlineTables, s2[i].inlineTables) {
			return false
		}
		if !reflect.DeepEqual(s1[i].aliases, s2[i].aliases) {
			return false
		}
	}
	return true
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

	_ = copyfile(filepath.Join(TestDir, "table7.jsonl"), filepath.Join(TestDataDir, "table7.jsonl"))

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
	TestTx.Session.SetStdout(NewDiscard())
	TestTx.Session.SetStderr(NewDiscard())
	initFlag(TestTx.Flags)
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
	flags.DatetimeFormat = []string{}
	flags.AnsiQuotes = false
	flags.StrictEqual = false
	flags.WaitTimeout = 15
	flags.ImportOptions = cmd.NewImportOptions()
	flags.ExportOptions = cmd.NewExportOptions()
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

const (
	scopeNameVariables    = "v"
	scopeNameTempTables   = "t"
	scopeNameCursors      = "c"
	scopeNameFunctions    = "f"
	scopeNameInlineTables = "i"
	scopeNameAliases      = "a"
)

func GenerateReferenceScope(blocks []map[string]map[string]interface{}, nodes []map[string]map[string]interface{}, now time.Time, records []ReferenceRecord) *ReferenceScope {
	rs := NewReferenceScope(TestTx)
	for i := 1; i < len(blocks); i++ {
		rs = rs.CreateChild()
	}

	for i := range blocks {
		for n := range blocks[i] {
			for k, v := range blocks[i][n] {
				switch n {
				case scopeNameVariables:
					rs.blocks[i].variables.Store(k, v.(value.Primary))
				case scopeNameTempTables:
					rs.blocks[i].temporaryTables.Store(k, v.(*View))
				case scopeNameCursors:
					rs.blocks[i].cursors.Store(k, v.(*Cursor))
				case scopeNameFunctions:
					rs.blocks[i].functions.Store(k, v.(*UserDefinedFunction))
				}
			}
		}
	}

	if nodes != nil {
		ns := make([]NodeScope, len(nodes))
		for i := range nodes {
			ns[i] = GetNodeScope()
			for n := range nodes[i] {
				for k, v := range nodes[i][n] {
					switch n {
					case scopeNameInlineTables:
						ns[i].inlineTables[k] = v.(*View)
					case scopeNameAliases:
						ns[i].aliases[k] = v.(string)
					}
				}
			}
		}
		rs.nodes = ns
	}

	if !now.IsZero() {
		rs.now = now
	}

	if records != nil {
		rs.Records = records
	}

	return rs
}

func GenerateViewMap(values []*View) ViewMap {
	m := NewViewMap()
	for _, v := range values {
		m.Store(v.FileInfo.Path, v)
	}
	return m
}

func GenerateCursorMap(values []*Cursor) CursorMap {
	m := NewCursorMap()
	for _, v := range values {
		m.Store(v.name, v)
	}
	return m
}

func GenerateUserDefinedFunctionMap(values []*UserDefinedFunction) UserDefinedFunctionMap {
	m := NewUserDefinedFunctionMap()
	for _, v := range values {
		m.Store(v.Name.Literal, v)
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

var (
	testLetterRunes = []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func randomStr(length int) string {
	s := make([]rune, length)
	for i := 0; i < length; i++ {
		s[i] = testLetterRunes[cmd.GetRand().Intn(len(testLetterRunes))]
	}
	return string(s)
}
