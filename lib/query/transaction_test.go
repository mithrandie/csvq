package query

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/mithrandie/csvq/lib/file"
	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/value"
)

func TestTransaction_Commit(t *testing.T) {
	defer func() {
		_ = TestTx.ReleaseResources()
		_ = TestTx.cachedViews.Clean(TestTx.FileContainer)
		initFlag(TestTx.Flags)
	}()

	TestTx.Flags.SetQuiet(false)

	ch, _ := file.NewHandlerForCreate(TestTx.FileContainer, GetTestFilePath("create_file.csv"))
	uh, _ := file.NewHandlerForUpdate(context.Background(), TestTx.FileContainer, GetTestFilePath("updated_file_1.csv"), TestTx.WaitTimeout, TestTx.RetryDelay)

	TestTx.cachedViews = ViewMap{
		strings.ToUpper(GetTestFilePath("created_file.csv")): &View{
			Header:    NewHeader("created_file", []string{"column1", "column2"}),
			RecordSet: RecordSet{},
			FileInfo: &FileInfo{
				Path:    GetTestFilePath("created_file.csv"),
				Handler: ch,
			},
		},
		strings.ToUpper(GetTestFilePath("updated_file_1.csv")): &View{
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
	}

	TestTx.uncommittedViews = &UncommittedViews{
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

	r, w, _ := os.Pipe()
	tx.Session.Stdout = w

	_ = TestTx.Commit(NewFilter(tx), parser.TransactionControl{Token: parser.COMMIT})

	_ = w.Close()
	log, _ := ioutil.ReadAll(r)

	if string(log) != expect {
		t.Errorf("Commit: log = %q, want %q", string(log), expect)
	}
}

func TestTransaction_Rollback(t *testing.T) {
	defer func() {
		_ = TestTx.ReleaseResources()
		_ = TestTx.cachedViews.Clean(TestTx.FileContainer)
		initFlag(TestTx.Flags)
	}()

	TestTx.Flags.SetQuiet(false)

	TestTx.uncommittedViews = &UncommittedViews{
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

	r, w, _ := os.Pipe()
	tx.Session.Stdout = w

	_ = TestTx.Rollback(NewFilter(tx), nil)

	_ = w.Close()
	log, _ := ioutil.ReadAll(r)

	if string(log) != expect {
		t.Errorf("Rollback: log = %q, want %q", string(log), expect)
	}
}
