package query

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/mithrandie/csvq/lib/file"
	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/value"
)

func TestTransaction_Commit(t *testing.T) {
	defer func() {
		_ = TestTx.ReleaseResources()
		_ = TestTx.cachedViews.Clean(TestTx.FileContainer)
		TestTx.Session.SetStdout(NewDiscard())
		initFlag(TestTx.Flags)
	}()

	TestTx.Flags.SetQuiet(false)

	ch, _ := file.NewHandlerForCreate(TestTx.FileContainer, GetTestFilePath("create_file.csv"))
	uh, _ := file.NewHandlerForUpdate(context.Background(), TestTx.FileContainer, GetTestFilePath("updated_file_1.csv"), TestTx.WaitTimeout, TestTx.RetryDelay)

	TestTx.cachedViews = GenerateViewMap([]*View{
		{
			Header:    NewHeader("created_file", []string{"column1", "column2"}),
			RecordSet: RecordSet{},
			FileInfo: &FileInfo{
				Path:    GetTestFilePath("created_file.csv"),
				Handler: ch,
			},
		},
		{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecord([]value.Primary{
					value.NewString("update1"),
					value.NewString("update2"),
				}),
				NewRecord([]value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
				}),
			},
			FileInfo: &FileInfo{
				Path:    GetTestFilePath("updated_file_1.csv"),
				Handler: uh,
			},
		},
	})

	TestTx.uncommittedViews = UncommittedViews{
		mtx: &sync.RWMutex{},
		Created: map[string]*FileInfo{
			strings.ToUpper(GetTestFilePath("created_file.csv")): {
				Path:    GetTestFilePath("created_file.csv"),
				Handler: ch,
			},
		},
		Updated: map[string]*FileInfo{
			strings.ToUpper(GetTestFilePath("updated_file_1.csv")): {
				Path:    GetTestFilePath("updated_file_1.csv"),
				Handler: uh,
			},
		},
	}

	expect := fmt.Sprintf("Commit: file %q is created.\nCommit: file %q is updated.\n", GetTestFilePath("created_file.csv"), GetTestFilePath("updated_file_1.csv"))

	tx := TestTx

	out := NewOutput()
	tx.Session.SetStdout(out)

	_ = TestTx.Commit(context.Background(), NewFilter(tx), parser.TransactionControl{Token: parser.COMMIT})

	log := out.String()

	if string(log) != expect {
		t.Errorf("Commit: log = %q, want %q", string(log), expect)
	}
}

func TestTransaction_Rollback(t *testing.T) {
	defer func() {
		_ = TestTx.ReleaseResources()
		_ = TestTx.cachedViews.Clean(TestTx.FileContainer)
		TestTx.Session.SetStdout(NewDiscard())
		initFlag(TestTx.Flags)
	}()

	TestTx.Flags.SetQuiet(false)

	TestTx.uncommittedViews = UncommittedViews{
		mtx: &sync.RWMutex{},
		Created: map[string]*FileInfo{
			strings.ToUpper(GetTestFilePath("created_file.csv")): {
				Path: GetTestFilePath("created_file.csv"),
			},
		},
		Updated: map[string]*FileInfo{
			strings.ToUpper(GetTestFilePath("updated_file_1.csv")): {
				Path: GetTestFilePath("updated_file_1.csv"),
			},
		},
	}

	expect := fmt.Sprintf("Rollback: file %q is deleted.\nRollback: file %q is restored.\n", GetTestFilePath("created_file.csv"), GetTestFilePath("updated_file_1.csv"))

	tx := TestTx

	out := NewOutput()
	tx.Session.SetStdout(out)

	_ = TestTx.Rollback(NewFilter(tx), nil)

	log := out.String()

	if string(log) != expect {
		t.Errorf("Rollback: log = %q, want %q", string(log), expect)
	}
}
