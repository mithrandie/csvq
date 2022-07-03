package query

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"
	"sync"
	"testing"

	"github.com/mithrandie/csvq/lib/option"

	"github.com/mithrandie/go-text"

	"github.com/mithrandie/csvq/lib/file"
	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/value"
)

func TestTransaction_Commit(t *testing.T) {
	defer func() {
		_ = TestTx.ReleaseResources()
		_ = TestTx.CachedViews.Clean(TestTx.FileContainer)
		TestTx.Session.SetStdout(NewDiscard())
		initFlag(TestTx.Flags)
	}()

	TestTx.Flags.SetQuiet(false)

	ch, _ := file.NewHandlerForCreate(TestTx.FileContainer, GetTestFilePath("created_file.csv"))
	uh, _ := file.NewHandlerForUpdate(context.Background(), TestTx.FileContainer, GetTestFilePath("updated_file_1.csv"), TestTx.WaitTimeout, TestTx.RetryDelay)

	TestTx.CachedViews = GenerateViewMap([]*View{
		{
			Header:    NewHeader("created_file", []string{"column1", "column2"}),
			RecordSet: RecordSet{},
			FileInfo: &FileInfo{
				Path:      GetTestFilePath("created_file.csv"),
				Handler:   ch,
				Encoding:  text.UTF8,
				Format:    option.CSV,
				Delimiter: ',',
				LineBreak: text.LF,
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
				Path:      GetTestFilePath("updated_file_1.csv"),
				Handler:   uh,
				Encoding:  text.UTF8,
				Format:    option.CSV,
				Delimiter: ',',
				LineBreak: text.LF,
			},
		},
	})

	TestTx.UncommittedViews = UncommittedViews{
		mtx: &sync.RWMutex{},
		Created: map[string]*FileInfo{
			strings.ToUpper(GetTestFilePath("created_file.csv")): {
				Path:      GetTestFilePath("created_file.csv"),
				Handler:   ch,
				Encoding:  text.UTF8,
				Format:    option.CSV,
				Delimiter: ',',
				LineBreak: text.LF,
			},
		},
		Updated: map[string]*FileInfo{
			strings.ToUpper(GetTestFilePath("updated_file_1.csv")): {
				Path:      GetTestFilePath("updated_file_1.csv"),
				Handler:   uh,
				Encoding:  text.UTF8,
				Format:    option.CSV,
				Delimiter: ',',
				LineBreak: text.LF,
			},
		},
	}

	expect := fmt.Sprintf("Commit: file %q is created.\nCommit: file %q is updated.\n", GetTestFilePath("created_file.csv"), GetTestFilePath("updated_file_1.csv"))

	tx := TestTx

	out := NewOutput()
	tx.Session.SetStdout(out)

	err := TestTx.Commit(context.Background(), NewReferenceScope(tx), parser.TransactionControl{Token: parser.COMMIT})
	if err != nil {
		t.Fatalf("unexpected error %q", err.Error())
	}

	log := out.String()

	if log != expect {
		t.Errorf("Commit: log = %q, want %q", log, expect)
	}

	expectedCreatedContents := "column1,column2\n"
	createdContents, err := ioutil.ReadFile(GetTestFilePath("created_file.csv"))
	if err != nil {
		t.Fatalf("unexpected error %q", err.Error())
	}

	if expectedCreatedContents != string(createdContents) {
		t.Errorf("created contents = %q, want %q", string(createdContents), expectedCreatedContents)
	}

	expectedUpdatedContents := "column1,column2\n1,str1\nupdate1,update2\n3,str3\n"
	updatedContents, err := ioutil.ReadFile(GetTestFilePath("updated_file_1.csv"))
	if err != nil {
		t.Fatalf("unexpected error %q", err.Error())
	}

	if expectedUpdatedContents != string(updatedContents) {
		t.Errorf("updated contents = %q, want %q", string(updatedContents), expectedUpdatedContents)
	}

	// Flags.StripEndingLineBreak = true
	TestTx.Flags.ExportOptions.StripEndingLineBreak = true
	ch, _ = file.NewHandlerForCreate(TestTx.FileContainer, GetTestFilePath("created_file_1.csv"))
	uh, _ = file.NewHandlerForUpdate(context.Background(), TestTx.FileContainer, GetTestFilePath("updated_file_1.csv"), TestTx.WaitTimeout, TestTx.RetryDelay)
	TestTx.CachedViews = GenerateViewMap([]*View{
		{
			Header:    NewHeader("created_file_1", []string{"column1", "column2"}),
			RecordSet: RecordSet{},
			FileInfo: &FileInfo{
				Path:      GetTestFilePath("created_file_1.csv"),
				Handler:   ch,
				Encoding:  text.UTF8,
				Format:    option.CSV,
				Delimiter: ',',
				LineBreak: text.LF,
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
				Path:      GetTestFilePath("updated_file_1.csv"),
				Handler:   uh,
				Encoding:  text.UTF8,
				Format:    option.CSV,
				Delimiter: ',',
				LineBreak: text.LF,
			},
		},
	})

	TestTx.UncommittedViews = UncommittedViews{
		mtx: &sync.RWMutex{},
		Created: map[string]*FileInfo{
			strings.ToUpper(GetTestFilePath("created_file_1.csv")): {
				Path:      GetTestFilePath("created_file_1.csv"),
				Handler:   ch,
				Encoding:  text.UTF8,
				Format:    option.CSV,
				Delimiter: ',',
				LineBreak: text.LF,
			},
		},
		Updated: map[string]*FileInfo{
			strings.ToUpper(GetTestFilePath("updated_file_1.csv")): {
				Path:      GetTestFilePath("updated_file_1.csv"),
				Handler:   uh,
				Encoding:  text.UTF8,
				Format:    option.CSV,
				Delimiter: ',',
				LineBreak: text.LF,
			},
		},
	}

	err = TestTx.Commit(context.Background(), NewReferenceScope(tx), parser.TransactionControl{Token: parser.COMMIT})
	if err != nil {
		t.Fatalf("unexpected error %q", err.Error())
	}

	expectedCreatedContents = "column1,column2"
	createdContents, err = ioutil.ReadFile(GetTestFilePath("created_file_1.csv"))
	if err != nil {
		t.Fatalf("unexpected error %q", err.Error())
	}

	if expectedCreatedContents != string(createdContents) {
		t.Errorf("created contents = %q, want %q", string(createdContents), expectedCreatedContents)
	}

	expectedUpdatedContents = "column1,column2\n1,str1\nupdate1,update2\n3,str3"
	updatedContents, err = ioutil.ReadFile(GetTestFilePath("updated_file_1.csv"))
	if err != nil {
		t.Fatalf("unexpected error %q", err.Error())
	}

	if expectedUpdatedContents != string(updatedContents) {
		t.Errorf("updated contents = %q, want %q", string(updatedContents), expectedUpdatedContents)
	}
}

func TestTransaction_Rollback(t *testing.T) {
	defer func() {
		_ = TestTx.ReleaseResources()
		_ = TestTx.CachedViews.Clean(TestTx.FileContainer)
		TestTx.Session.SetStdout(NewDiscard())
		initFlag(TestTx.Flags)
	}()

	TestTx.Flags.SetQuiet(false)

	TestTx.UncommittedViews = UncommittedViews{
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

	_ = TestTx.Rollback(NewReferenceScope(tx), nil)

	log := out.String()

	if log != expect {
		t.Errorf("Rollback: log = %q, want %q", log, expect)
	}
}
