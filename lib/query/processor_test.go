package query

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/value"

	"github.com/mithrandie/go-text"
	"github.com/mithrandie/ternary"
)

func TestContextForStoringResults(t *testing.T) {
	defer func() {
		TestTx.SelectedViews = nil
		TestTx.AffectedRows = 0
	}()

	tx := TestTx
	proc := NewProcessor(tx)

	ctx := ContextForStoringResults(context.Background())
	_, err := proc.Execute(ctx, []parser.Statement{
		parser.SelectQuery{
			SelectEntity: parser.SelectEntity{
				SelectClause: parser.SelectClause{
					Fields: []parser.QueryExpression{
						parser.Field{Object: parser.NewIntegerValueFromString("1")},
					},
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error %q", err)
	}

	if len(proc.Tx.SelectedViews) < 1 {
		t.Fatalf("result set is not set to the transaction")
	}
}

var processorExecuteStatementTests = []struct {
	Input            parser.Statement
	UncommittedViews UncommittedViews
	Logs             string
	SelectLogs       []string
	Error            string
	ReturnCode       int
}{
	{
		Input: parser.SetFlag{
			Flag:  parser.Flag{Name: "invalid"},
			Value: parser.NewStringValue("\t"),
		},
		Error:      "@@INVALID is an unknown flag",
		ReturnCode: ReturnCodeApplicationError,
	},
	{
		Input: parser.SetFlag{
			Flag:  parser.Flag{Name: "delimiter"},
			Value: parser.NewStringValue(","),
		},
	},
	{
		Input: parser.Echo{
			Value: parser.Function{
				Name: "trunc_time",
				Args: []parser.QueryExpression{
					parser.Function{
						Name: "datetime",
						Args: []parser.QueryExpression{
							parser.NewStringValue("2001::01::01"),
						},
					},
				},
			},
		},
		Logs: "NULL\n",
	},
	{
		Input: parser.AddFlagElement{
			Flag:  parser.Flag{Name: "datetime_format"},
			Value: parser.NewStringValue("%Y::%m::%d"),
		},
	},
	{
		Input: parser.Echo{
			Value: parser.Function{
				Name: "trunc_time",
				Args: []parser.QueryExpression{
					parser.Function{
						Name: "datetime",
						Args: []parser.QueryExpression{
							parser.NewStringValue("2001::01::01"),
						},
					},
				},
			},
		},
		Logs: "2001-01-01T00:00:00Z\n",
	},
	{
		Input: parser.RemoveFlagElement{
			Flag:  parser.Flag{Name: "datetime_format"},
			Value: parser.NewStringValue("%Y::%m::%d"),
		},
	},
	{
		Input: parser.Echo{
			Value: parser.Function{
				Name: "trunc_time",
				Args: []parser.QueryExpression{
					parser.Function{
						Name: "datetime",
						Args: []parser.QueryExpression{
							parser.NewStringValue("2001::01::01"),
						},
					},
				},
			},
		},
		Logs: "NULL\n",
	},
	{
		Input: parser.ShowFlag{
			Flag: parser.Flag{Name: "repository"},
		},
		Logs: "@@REPOSITORY: " + TestDir + "\n",
	},
	{
		Input: parser.SetEnvVar{
			EnvVar: parser.EnvironmentVariable{Name: "CSVQ_PROC_TEST"},
			Value:  parser.NewStringValue("foo"),
		},
	},
	{
		Input: parser.Echo{
			Value: parser.EnvironmentVariable{Name: "CSVQ_PROC_TEST"},
		},
		Logs: "foo\n",
	},
	{
		Input: parser.Print{
			Value: parser.EnvironmentVariable{Name: "CSVQ_PROC_TEST"},
		},
		Logs: "'foo'\n",
	},
	{
		Input: parser.UnsetEnvVar{
			EnvVar: parser.EnvironmentVariable{Name: "CSVQ_PROC_TEST"},
		},
	},
	{
		Input: parser.Print{
			Value: parser.EnvironmentVariable{Name: "CSVQ_PROC_TEST"},
		},
		Logs: "''\n",
	},
	{
		Input: parser.VariableDeclaration{
			Assignments: []parser.VariableAssignment{
				{
					Variable: parser.Variable{Name: "var1"},
				},
			},
		},
	},
	{
		Input: parser.VariableDeclaration{
			Assignments: []parser.VariableAssignment{
				{
					Variable: parser.Variable{Name: "var2"},
				},
			},
		},
	},
	{
		Input: parser.VariableDeclaration{
			Assignments: []parser.VariableAssignment{
				{
					Variable: parser.Variable{Name: "var3"},
				},
			},
		},
	},
	{
		Input: parser.VariableDeclaration{
			Assignments: []parser.VariableAssignment{
				{
					Variable: parser.Variable{Name: "var4"},
				},
			},
		},
	},
	{
		Input: parser.VariableSubstitution{
			Variable: parser.Variable{Name: "var1"},
			Value:    parser.NewIntegerValueFromString("1"),
		},
	},
	{
		Input: parser.Print{
			Value: parser.Variable{Name: "var1"},
		},
		Logs: "1\n",
	},
	{
		Input: parser.DisposeVariable{
			Variable: parser.Variable{Name: "var4"},
		},
	},
	{
		Input: parser.VariableDeclaration{
			Assignments: []parser.VariableAssignment{
				{
					Variable: parser.Variable{Name: "var4"},
				},
			},
		},
	},
	{
		Input: parser.FunctionDeclaration{
			Name: parser.Identifier{Literal: "userfunc"},
			Parameters: []parser.VariableAssignment{
				{
					Variable: parser.Variable{Name: "arg1"},
				},
			},
			Statements: []parser.Statement{
				parser.Print{
					Value: parser.Variable{Name: "arg1"},
				},
			},
		},
	},
	{
		Input: parser.Function{
			Name: "userfunc",
			Args: []parser.QueryExpression{
				parser.NewIntegerValueFromString("1"),
			},
		},
		Logs: "1\n",
	},
	{
		Input: parser.DisposeFunction{
			Name: parser.Identifier{Literal: "userfunc"},
		},
		Logs: "",
	},
	{
		Input: parser.Function{
			Name: "userfunc",
			Args: []parser.QueryExpression{
				parser.NewIntegerValueFromString("1"),
			},
		},
		Error:      "function userfunc does not exist",
		ReturnCode: ReturnCodeApplicationError,
	},
	{
		Input: parser.CursorDeclaration{
			Cursor: parser.Identifier{Literal: "cur"},
			Query:  selectQueryForCursorTest,
		},
	},
	{
		Input: parser.OpenCursor{
			Cursor: parser.Identifier{Literal: "cur"},
		},
	},
	{
		Input: parser.FetchCursor{
			Cursor: parser.Identifier{Literal: "cur"},
			Position: parser.FetchPosition{
				Position: parser.Token{Token: parser.NEXT, Literal: "next"},
			},
			Variables: []parser.Variable{
				{Name: "var2"},
				{Name: "var3"},
			},
		},
	},
	{
		Input: parser.Print{
			Value: parser.Variable{Name: "var2"},
		},
		Logs: "'1'\n",
	},
	{
		Input: parser.Print{
			Value: parser.Variable{Name: "var3"},
		},
		Logs: "'str1'\n",
	},
	{
		Input: parser.CloseCursor{
			Cursor: parser.Identifier{Literal: "cur"},
		},
	},
	{
		Input: parser.DisposeCursor{
			Cursor: parser.Identifier{Literal: "cur"},
		},
	},
	{
		Input: parser.ViewDeclaration{
			View: parser.Identifier{Literal: "tbl"},
			Fields: []parser.QueryExpression{
				parser.Identifier{Literal: "column1"},
				parser.Identifier{Literal: "column2"},
			},
			Query: parser.SelectQuery{
				SelectEntity: parser.SelectEntity{
					SelectClause: parser.SelectClause{
						Fields: []parser.QueryExpression{
							parser.Field{Object: parser.NewIntegerValueFromString("1")},
							parser.Field{Object: parser.NewIntegerValueFromString("2")},
						},
					},
				},
			},
		},
	},
	{
		Input: parser.StatementPreparation{
			Name:      parser.Identifier{Literal: "stmt"},
			Statement: value.NewString("select 1"),
		},
	},
	{
		Input: parser.ExecuteStatement{
			Name: parser.Identifier{Literal: "invalidstmt"},
		},
		Error:      "statement invalidstmt does not exist",
		ReturnCode: ReturnCodeApplicationError,
	},
	{
		Input: parser.ExecuteStatement{
			Name: parser.Identifier{Literal: "stmt"},
		},
		Logs: "1\n1\n",
	},
	{
		Input: parser.DisposeStatement{
			Name: parser.Identifier{Literal: "stmt"},
		},
	},
	{
		Input: parser.SelectQuery{
			SelectEntity: parser.SelectEntity{
				SelectClause: parser.SelectClause{
					Fields: []parser.QueryExpression{
						parser.Field{
							Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
						},
						parser.Field{
							Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
						},
					},
				},
				FromClause: parser.FromClause{
					Tables: []parser.QueryExpression{
						parser.Table{Object: parser.Identifier{Literal: "tbl"}},
					},
				},
			},
		},
		Logs: "column1,column2\n1,2\n",
	},
	{
		Input: parser.SetFlag{
			Flag:  parser.Flag{Name: "format"},
			Value: parser.NewStringValue("text"),
		},
	},
	{
		Input: parser.SelectQuery{
			SelectEntity: parser.SelectEntity{
				SelectClause: parser.SelectClause{
					Fields: []parser.QueryExpression{
						parser.Field{
							Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
						},
						parser.Field{
							Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
						},
					},
				},
				FromClause: parser.FromClause{
					Tables: []parser.QueryExpression{
						parser.Table{Object: parser.Identifier{Literal: "tbl"}},
					},
				},
				WhereClause: parser.WhereClause{
					BaseExpr: nil,
					Filter:   parser.NewTernaryValueFromString("false"),
				},
			},
		},
		Logs: "Empty RecordSet\n",
	},
	{
		Input: parser.SetFlag{
			Flag:  parser.Flag{Name: "without_header"},
			Value: parser.NewTernaryValueFromString("true"),
		},
	},
	{
		Input: parser.SetFlag{
			Flag:  parser.Flag{Name: "format"},
			Value: parser.NewStringValue("csv"),
		},
	},
	{
		Input: parser.SelectQuery{
			SelectEntity: parser.SelectEntity{
				SelectClause: parser.SelectClause{
					Fields: []parser.QueryExpression{
						parser.Field{
							Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
						},
						parser.Field{
							Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
						},
					},
				},
				FromClause: parser.FromClause{
					Tables: []parser.QueryExpression{
						parser.Table{Object: parser.Identifier{Literal: "tbl"}},
					},
				},
				WhereClause: parser.WhereClause{
					BaseExpr: nil,
					Filter:   parser.NewTernaryValueFromString("false"),
				},
			},
		},
		Logs: "",
	},
	{
		Input: parser.SetFlag{
			Flag:  parser.Flag{Name: "without_header"},
			Value: parser.NewTernaryValueFromString("false"),
		},
	},
	{
		Input: parser.SelectQuery{
			SelectEntity: parser.SelectEntity{
				SelectClause: parser.SelectClause{
					Fields: []parser.QueryExpression{
						parser.Field{
							Object: parser.NewIntegerValueFromString("1234"),
						},
					},
				},
				IntoClause: parser.IntoClause{
					Variables: []parser.Variable{
						{Name: "var4"},
					},
				},
			},
		},
	},
	{
		Input: parser.Print{
			Value: parser.Variable{Name: "var4"},
		},
		Logs: "1234\n",
	},
	{
		Input: parser.DisposeView{
			View: parser.Identifier{Literal: "tbl"},
		},
	},
	{
		Input: parser.AggregateDeclaration{
			Name:   parser.Identifier{Literal: "useraggfunc"},
			Cursor: parser.Identifier{Literal: "list"},
			Statements: []parser.Statement{
				parser.VariableDeclaration{
					Assignments: []parser.VariableAssignment{
						{
							Variable: parser.Variable{Name: "value"},
						},
						{
							Variable: parser.Variable{Name: "fetch"},
						},
					},
				},
				parser.WhileInCursor{
					Variables: []parser.Variable{
						{Name: "fetch"},
					},
					Cursor: parser.Identifier{Literal: "list"},
					Statements: []parser.Statement{
						parser.If{
							Condition: parser.Is{
								LHS: parser.Variable{Name: "fetch"},
								RHS: parser.NewNullValue(),
							},
							Statements: []parser.Statement{
								parser.FlowControl{Token: parser.CONTINUE},
							},
						},
						parser.If{
							Condition: parser.Is{
								LHS: parser.Variable{Name: "value"},
								RHS: parser.NewNullValue(),
							},
							Statements: []parser.Statement{
								parser.VariableSubstitution{
									Variable: parser.Variable{Name: "value"},
									Value:    parser.Variable{Name: "fetch"},
								},
								parser.FlowControl{Token: parser.CONTINUE},
							},
						},
						parser.VariableSubstitution{
							Variable: parser.Variable{Name: "value"},
							Value: parser.Arithmetic{
								LHS:      parser.Variable{Name: "value"},
								RHS:      parser.Variable{Name: "fetch"},
								Operator: parser.Token{Token: '*', Literal: "*"},
							},
						},
					},
				},
				parser.Return{
					Value: parser.Variable{Name: "value"},
				},
			},
		},
	},
	{
		Input: parser.SelectQuery{
			SelectEntity: parser.SelectEntity{
				SelectClause: parser.SelectClause{
					Fields: []parser.QueryExpression{
						parser.Field{
							Object: parser.Function{
								Name: "useraggfunc",
								Args: []parser.QueryExpression{
									parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
								},
							},
							Alias: parser.Identifier{Literal: "multiplication"},
						},
					},
				},
				FromClause: parser.FromClause{
					Tables: []parser.QueryExpression{
						parser.Table{Object: parser.Identifier{Literal: "table1"}},
					},
				},
			},
		},
		Logs: "multiplication\n6\n",
	},
	{
		Input: parser.SelectQuery{
			SelectEntity: parser.SelectEntity{
				SelectClause: parser.SelectClause{
					Fields: []parser.QueryExpression{
						parser.Field{
							Object: parser.Variable{Name: "var1"},
							Alias:  parser.Identifier{Literal: "var1"},
						},
					},
				},
			},
		},
		Logs: "var1\n1\n",
	},
	{
		Input: parser.SetFlag{
			Flag:  parser.Flag{Name: "strip_ending_line_break"},
			Value: parser.NewTernaryValue(ternary.TRUE),
		},
	},
	{
		Input: parser.SelectQuery{
			SelectEntity: parser.SelectEntity{
				SelectClause: parser.SelectClause{
					Fields: []parser.QueryExpression{
						parser.Field{
							Object: parser.Variable{Name: "var1"},
							Alias:  parser.Identifier{Literal: "var1"},
						},
					},
				},
			},
		},
		Logs: "var1\n1",
	},
	{
		Input: parser.SetFlag{
			Flag:  parser.Flag{Name: "strip_ending_line_break"},
			Value: parser.NewTernaryValue(ternary.FALSE),
		},
	},
	{
		Input: parser.VariableDeclaration{
			Assignments: []parser.VariableAssignment{
				{
					Variable: parser.Variable{Name: "var1"},
				},
			},
		},
		Error:      "variable @var1 is redeclared",
		ReturnCode: ReturnCodeApplicationError,
	},
	{
		Input: parser.VariableSubstitution{
			Variable: parser.Variable{Name: "var9"},
			Value:    parser.NewIntegerValueFromString("1"),
		},
		Error:      "variable @var9 is undeclared",
		ReturnCode: ReturnCodeApplicationError,
	},
	{
		Input: parser.InsertQuery{
			Table: parser.Table{Object: parser.Identifier{Literal: "table1"}},
			Fields: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
			},
			ValuesList: []parser.QueryExpression{
				parser.RowValue{
					Value: parser.ValueList{
						Values: []parser.QueryExpression{
							parser.NewIntegerValueFromString("4"),
							parser.NewStringValue("str4"),
						},
					},
				},
				parser.RowValue{
					Value: parser.ValueList{
						Values: []parser.QueryExpression{
							parser.NewIntegerValueFromString("5"),
							parser.NewStringValue("str5"),
						},
					},
				},
			},
		},
		UncommittedViews: UncommittedViews{
			mtx:     &sync.RWMutex{},
			Created: map[string]*FileInfo{},
			Updated: map[string]*FileInfo{
				strings.ToUpper(GetTestFilePath("TABLE1.CSV")): {
					Path:      GetTestFilePath("table1.csv"),
					Delimiter: ',',
					NoHeader:  false,
					Encoding:  text.UTF8,
					LineBreak: text.LF,
					ForUpdate: true,
				},
			},
		},
		Logs: fmt.Sprintf("2 records inserted on %q.\n", GetTestFilePath("table1.csv")),
	},
	{
		Input: parser.UpdateQuery{
			Tables: []parser.QueryExpression{
				parser.Table{Object: parser.Identifier{Literal: "table1"}},
			},
			SetList: []parser.UpdateSet{
				{
					Field: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
					Value: parser.NewStringValue("update"),
				},
			},
			WhereClause: parser.WhereClause{
				Filter: parser.Comparison{
					LHS:      parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
					RHS:      parser.NewIntegerValueFromString("2"),
					Operator: parser.Token{Token: '=', Literal: "="},
				},
			},
		},
		UncommittedViews: UncommittedViews{
			mtx:     &sync.RWMutex{},
			Created: map[string]*FileInfo{},
			Updated: map[string]*FileInfo{
				strings.ToUpper(GetTestFilePath("TABLE1.CSV")): {
					Path:      GetTestFilePath("table1.csv"),
					Delimiter: ',',
					NoHeader:  false,
					Encoding:  text.UTF8,
					LineBreak: text.LF,
					ForUpdate: true,
				},
			},
		},
		Logs: fmt.Sprintf("1 record updated on %q.\n", GetTestFilePath("table1.csv")),
	},
	{
		Input: parser.ReplaceQuery{
			Table: parser.Table{Object: parser.Identifier{Literal: "table1"}},
			Fields: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
			},
			Keys: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
			},
			ValuesList: []parser.QueryExpression{
				parser.RowValue{
					Value: parser.ValueList{
						Values: []parser.QueryExpression{
							parser.NewIntegerValueFromString("4"),
							parser.NewStringValue("str44"),
						},
					},
				},
				parser.RowValue{
					Value: parser.ValueList{
						Values: []parser.QueryExpression{
							parser.NewIntegerValueFromString("6"),
							parser.NewStringValue("str6"),
						},
					},
				},
			},
		},
		UncommittedViews: UncommittedViews{
			mtx:     &sync.RWMutex{},
			Created: map[string]*FileInfo{},
			Updated: map[string]*FileInfo{
				strings.ToUpper(GetTestFilePath("TABLE1.CSV")): {
					Path:      GetTestFilePath("table1.csv"),
					Delimiter: ',',
					NoHeader:  false,
					Encoding:  text.UTF8,
					LineBreak: text.LF,
					ForUpdate: true,
				},
			},
		},
		Logs: fmt.Sprintf("2 records replaced on %q.\n", GetTestFilePath("table1.csv")),
	},
	{
		Input: parser.DeleteQuery{
			FromClause: parser.FromClause{
				Tables: []parser.QueryExpression{
					parser.Table{
						Object: parser.Identifier{Literal: "table1"},
					},
				},
			},
			WhereClause: parser.WhereClause{
				Filter: parser.Comparison{
					LHS:      parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
					RHS:      parser.NewIntegerValueFromString("2"),
					Operator: parser.Token{Token: '=', Literal: "="},
				},
			},
		},
		UncommittedViews: UncommittedViews{
			mtx:     &sync.RWMutex{},
			Created: map[string]*FileInfo{},
			Updated: map[string]*FileInfo{
				strings.ToUpper(GetTestFilePath("TABLE1.CSV")): {
					Path:      GetTestFilePath("table1.csv"),
					Delimiter: ',',
					NoHeader:  false,
					Encoding:  text.UTF8,
					LineBreak: text.LF,
					ForUpdate: true,
				},
			},
		},
		Logs: fmt.Sprintf("1 record deleted on %q.\n", GetTestFilePath("table1.csv")),
	},
	{
		Input: parser.CreateTable{
			Table: parser.Identifier{Literal: "newtable.csv"},
			Fields: []parser.QueryExpression{
				parser.Identifier{Literal: "column1"},
				parser.Identifier{Literal: "column2"},
			},
		},
		UncommittedViews: UncommittedViews{
			mtx: &sync.RWMutex{},
			Created: map[string]*FileInfo{
				strings.ToUpper(GetTestFilePath("NEWTABLE.CSV")): {
					Path:      GetTestFilePath("newtable.csv"),
					Delimiter: ',',
					NoHeader:  false,
					Encoding:  text.UTF8,
					LineBreak: text.LF,
					ForUpdate: true,
				},
			},
			Updated: map[string]*FileInfo{},
		},
		Logs: fmt.Sprintf("file %q is created.\n", GetTestFilePath("newtable.csv")),
	},
	{
		Input: parser.CreateTable{
			Table: parser.Identifier{Literal: "table1.csv"},
			Fields: []parser.QueryExpression{
				parser.Identifier{Literal: "column1"},
				parser.Identifier{Literal: "column2"},
			},
		},
		Error:      fmt.Sprintf("file %s already exists", GetTestFilePath("table1.csv")),
		ReturnCode: ReturnCodeIOError,
	},
	{
		Input: parser.CreateTable{
			Table: parser.Identifier{Literal: "table1.csv"},
			Fields: []parser.QueryExpression{
				parser.Identifier{Literal: "column1"},
				parser.Identifier{Literal: "column2"},
			},
			IfNotExists: true,
		},
		Logs: fmt.Sprintf("file %q already exists.\n", GetTestFilePath("table1.csv")),
	},
	{
		Input: parser.CreateTable{
			Table: parser.Identifier{Literal: "table1.csv"},
			Fields: []parser.QueryExpression{
				parser.Identifier{Literal: "column1"},
				parser.Identifier{Literal: "column2"},
				parser.Identifier{Literal: "column3"},
			},
			IfNotExists: true,
		},
		Error:      "field length does not match",
		ReturnCode: ReturnCodeApplicationError,
	},
	{
		Input: parser.CreateTable{
			Table: parser.Identifier{Literal: "table1.csv"},
			Fields: []parser.QueryExpression{
				parser.Identifier{Literal: "column1"},
				parser.Identifier{Literal: "col"},
			},
			IfNotExists: true,
		},
		Error:      "field col does not exist",
		ReturnCode: ReturnCodeApplicationError,
	},
	{
		Input: parser.AddColumns{
			Table: parser.Identifier{Literal: "table1.csv"},
			Columns: []parser.ColumnDefault{
				{
					Column: parser.Identifier{Literal: "column3"},
				},
			},
		},
		UncommittedViews: UncommittedViews{
			mtx:     &sync.RWMutex{},
			Created: map[string]*FileInfo{},
			Updated: map[string]*FileInfo{
				strings.ToUpper(GetTestFilePath("TABLE1.CSV")): {
					Path:      GetTestFilePath("table1.csv"),
					Delimiter: ',',
					NoHeader:  false,
					Encoding:  text.UTF8,
					LineBreak: text.LF,
					ForUpdate: true,
				},
			},
		},
		Logs: fmt.Sprintf("1 field added on %q.\n", GetTestFilePath("table1.csv")),
	},
	{
		Input: parser.DropColumns{
			Table: parser.Identifier{Literal: "table1"},
			Columns: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
			},
		},
		UncommittedViews: UncommittedViews{
			mtx:     &sync.RWMutex{},
			Created: map[string]*FileInfo{},
			Updated: map[string]*FileInfo{
				strings.ToUpper(GetTestFilePath("TABLE1.CSV")): {
					Path:      GetTestFilePath("table1.csv"),
					Delimiter: ',',
					NoHeader:  false,
					Encoding:  text.UTF8,
					LineBreak: text.LF,
					ForUpdate: true,
				},
			},
		},
		Logs: fmt.Sprintf("1 field dropped on %q.\n", GetTestFilePath("table1.csv")),
	},
	{
		Input: parser.RenameColumn{
			Table: parser.Identifier{Literal: "table1"},
			Old:   parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
			New:   parser.Identifier{Literal: "newcolumn"},
		},
		UncommittedViews: UncommittedViews{
			mtx:     &sync.RWMutex{},
			Created: map[string]*FileInfo{},
			Updated: map[string]*FileInfo{
				strings.ToUpper(GetTestFilePath("TABLE1.CSV")): {
					Path:      GetTestFilePath("table1.csv"),
					Delimiter: ',',
					NoHeader:  false,
					Encoding:  text.UTF8,
					LineBreak: text.LF,
					ForUpdate: true,
				},
			},
		},
		Logs: fmt.Sprintf("1 field renamed on %q.\n", GetTestFilePath("table1.csv")),
	},
	{
		Input: parser.SetTableAttribute{
			Table:     parser.Identifier{Literal: "table1.csv"},
			Attribute: parser.Identifier{Literal: "delimiter"},
			Value:     parser.NewStringValue(","),
		},
		Logs: "Table attributes of " + GetTestFilePath("table1.csv") + " remain unchanged.\n",
	},
	{
		Input: parser.SetTableAttribute{
			Table:     parser.Identifier{Literal: "table1.csv"},
			Attribute: parser.Identifier{Literal: "delimiter"},
			Value:     parser.NewStringValue("\t"),
		},
		UncommittedViews: UncommittedViews{
			mtx:     &sync.RWMutex{},
			Created: map[string]*FileInfo{},
			Updated: map[string]*FileInfo{
				strings.ToUpper(GetTestFilePath("TABLE1.CSV")): {
					Path:      GetTestFilePath("table1.csv"),
					Delimiter: '\t',
					NoHeader:  false,
					Encoding:  text.UTF8,
					LineBreak: text.LF,
					Format:    cmd.TSV,
					ForUpdate: true,
				},
			},
		},
		Logs: "\n" +
			strings.Repeat(" ", (calcShowFieldsWidth("table1.csv", "table1.csv", 22)-(22+len("table1.csv")))/2) + "Attributes Updated in table1.csv\n" +
			strings.Repeat("-", calcShowFieldsWidth("table1.csv", "table1.csv", 22)) + "\n" +
			" Path: " + GetTestFilePath("table1.csv") + "\n" +
			" Format: TSV     Delimiter: '\\t'  Enclose All: false\n" +
			" Encoding: UTF8  LineBreak: LF    Header: true\n" +
			"\n",
	},
	{
		Input: parser.Case{
			When: []parser.CaseWhen{
				{
					Condition: parser.NewTernaryValue(ternary.FALSE),
					Statements: []parser.Statement{
						parser.Print{Value: parser.NewStringValue("1")},
					},
				},
				{
					Condition: parser.NewTernaryValue(ternary.TRUE),
					Statements: []parser.Statement{
						parser.Print{Value: parser.NewStringValue("2")},
					},
				},
			},
		},
		Logs: "'2'\n",
	},
	{
		Input: parser.While{
			Condition: parser.Comparison{
				LHS:      parser.Variable{Name: "while_test"},
				RHS:      parser.NewIntegerValueFromString("3"),
				Operator: parser.Token{Token: '<', Literal: "<"},
			},
			Statements: []parser.Statement{
				parser.VariableSubstitution{
					Variable: parser.Variable{Name: "while_test"},
					Value: parser.Arithmetic{
						LHS:      parser.Variable{Name: "while_test"},
						RHS:      parser.NewIntegerValueFromString("1"),
						Operator: parser.Token{Token: '+', Literal: "+"},
					},
				},
				parser.Print{Value: parser.Variable{Name: "while_test"}},
			},
		},
		Logs: "1\n2\n3\n",
	},
	{
		Input: parser.Exit{
			Code: value.NewInteger(1),
		},
		Error:      ExitMessage,
		ReturnCode: 1,
	},
	{
		Input: parser.Print{
			Value: parser.NewIntegerValue(12345),
		},
		Logs: "12345\n",
	},
	{
		Input: parser.Printf{
			Format: parser.NewStringValue("value: %s"),
			Values: []parser.QueryExpression{
				parser.NewIntegerValue(12345),
			},
		},
		Logs: "value: 12345\n",
	},
	{
		Input: parser.Source{
			FilePath: parser.NewStringValue(GetTestFilePath("source.sql")),
		},
		Logs: "'external executable file'\n",
	},
	{
		Input: parser.Execute{
			BaseExpr:   parser.NewBaseExpr(parser.Token{}),
			Statements: parser.NewStringValue("print 'execute';"),
		},
		Logs: "'execute'\n",
	},
	{
		Input: parser.Trigger{
			Event:   parser.Identifier{Literal: "error"},
			Message: parser.NewStringValue("user error"),
			Code:    value.NewInteger(200),
		},
		Error:      "user error",
		ReturnCode: 200,
	},
	{
		Input: parser.Trigger{
			Event:   parser.Identifier{Literal: "error"},
			Message: parser.NewIntegerValue(200),
		},
		Error:      DefaultUserTriggeredErrorMessage,
		ReturnCode: 200,
	},
	{
		Input: parser.Trigger{
			Event:   parser.Identifier{Literal: "invalid"},
			Message: parser.NewIntegerValue(200),
		},
		Error:      "invalid is an unknown event",
		ReturnCode: ReturnCodeApplicationError,
	},
	{
		Input: parser.ShowObjects{
			Type: parser.Identifier{Literal: "cursors"},
		},
		Logs: "No cursor is declared\n",
	},
	{
		Input: parser.ShowFields{
			Type:  parser.Identifier{Literal: "fields"},
			Table: parser.Identifier{Literal: "table1"},
		},
		Logs: "\n" +
			strings.Repeat(" ", (calcShowFieldsWidth("table1.csv", "table1", 10)-(10+len("table1")))/2) + "Fields in table1\n" +
			strings.Repeat("-", calcShowFieldsWidth("table1.csv", "table1", 10)) + "\n" +
			" Type: Table\n" +
			" Path: " + GetTestFilePath("table1.csv") + "\n" +
			" Format: CSV     Delimiter: ','   Enclose All: false\n" +
			" Encoding: UTF8  LineBreak: LF    Header: true\n" +
			" Status: Fixed\n" +
			" Fields:\n" +
			"   1. column1\n" +
			"   2. column2\n" +
			"\n",
	},
}

func TestProcessor_ExecuteStatement(t *testing.T) {
	defer func() {
		_ = TestTx.ReleaseResources()
		TestTx.uncommittedViews.Clean()
		TestTx.Session.SetStdout(NewDiscard())
		initFlag(TestTx.Flags)
	}()

	TestTx.Flags.Repository = TestDir
	TestTx.Flags.ExportOptions.Format = cmd.CSV

	tx := TestTx
	proc := NewProcessor(tx)
	_ = proc.ReferenceScope.DeclareVariableDirectly(parser.Variable{Name: "while_test"}, value.NewInteger(0))
	ctx := context.Background()

	for _, v := range processorExecuteStatementTests {
		_ = TestTx.ReleaseResources()
		TestTx.uncommittedViews = NewUncommittedViews()

		out := NewOutput()
		tx.Session.SetStdout(out)
		_, err := proc.ExecuteStatement(ctx, v.Input)
		log := out.String()

		if err != nil {
			var code int
			if apperr, ok := err.(Error); ok {
				if len(v.Error) < 1 {
					t.Errorf("unexpected error %q for %q", err, v.Input)
				} else if err.Error() != v.Error {
					t.Errorf("error %q, want error %q for %q", err, v.Error, v.Input)
				}

				code = apperr.Code()
			}
			if code != v.ReturnCode {
				t.Errorf("error code %d, want error code %d for %q", code, v.ReturnCode, v.Input)
			}
			continue
		}
		if 0 < len(v.Error) {
			t.Errorf("no error, want error %q for %q", v.Error, v.Input)
			continue
		}

		if v.UncommittedViews.mtx != nil {
			for _, r := range TestTx.uncommittedViews.Created {
				if r.Handler != nil {
					if r.Path != r.Handler.Path() {
						t.Errorf("file pointer = %q, want %q for %q", r.Handler.Path(), r.Path, v.Input)
					}
					_ = TestTx.FileContainer.Close(r.Handler)
					r.Handler = nil
				}
			}
			for _, r := range TestTx.uncommittedViews.Updated {
				if r.Handler != nil {
					if r.Path != r.Handler.Path() {
						t.Errorf("file pointer = %q, want %q for %q", r.Handler.Path(), r.Path, v.Input)
					}
					_ = TestTx.FileContainer.Close(r.Handler)
					r.Handler = nil
				}
			}

			if !reflect.DeepEqual(TestTx.uncommittedViews, v.UncommittedViews) {
				t.Errorf("uncomitted views = %v, want %v for %q", TestTx.uncommittedViews, v.UncommittedViews, v.Input)
			}
		}
		if 0 < len(v.Logs) {
			if log != v.Logs {
				t.Errorf("logs = %s, want %s for %q", log, v.Logs, v.Input)
			}
		}
		if v.SelectLogs != nil {
			selectLog := log
			if !reflect.DeepEqual(selectLog, v.SelectLogs) {
				t.Errorf("select logs = %s, want %s for %q", selectLog, v.SelectLogs, v.Input)
			}
		}
	}
}

var processorIfStmtTests = []struct {
	Name        string
	Stmt        parser.If
	ResultFlow  StatementFlow
	ReturnValue value.Primary
	Result      string
	Error       string
}{
	{
		Name: "If Statement",
		Stmt: parser.If{
			Condition: parser.NewTernaryValue(ternary.TRUE),
			Statements: []parser.Statement{
				parser.Print{Value: parser.NewStringValue("1")},
			},
		},
		ResultFlow: Terminate,
		Result:     "'1'\n",
	},
	{
		Name: "If Statement Execute Nothing",
		Stmt: parser.If{
			Condition: parser.NewTernaryValue(ternary.FALSE),
			Statements: []parser.Statement{
				parser.Print{Value: parser.NewStringValue("1")},
			},
		},
		ResultFlow: Terminate,
		Result:     "",
	},
	{
		Name: "If Statement Execute ElseIf",
		Stmt: parser.If{
			Condition: parser.NewTernaryValue(ternary.FALSE),
			Statements: []parser.Statement{
				parser.Print{Value: parser.NewStringValue("1")},
			},
			ElseIf: []parser.ElseIf{
				{
					Condition: parser.NewTernaryValue(ternary.TRUE),
					Statements: []parser.Statement{
						parser.Print{Value: parser.NewStringValue("2")},
					},
				},
				{
					Condition: parser.NewTernaryValue(ternary.FALSE),
					Statements: []parser.Statement{
						parser.Print{Value: parser.NewStringValue("3")},
					},
				},
			},
			Else: parser.Else{
				Statements: []parser.Statement{
					parser.Print{Value: parser.NewStringValue("4")},
				},
			},
		},
		ResultFlow: Terminate,
		Result:     "'2'\n",
	},
	{
		Name: "If Statement Execute Else",
		Stmt: parser.If{
			Condition: parser.NewTernaryValue(ternary.FALSE),
			Statements: []parser.Statement{
				parser.Print{Value: parser.NewStringValue("1")},
			},
			ElseIf: []parser.ElseIf{
				{
					Condition: parser.NewTernaryValue(ternary.FALSE),
					Statements: []parser.Statement{
						parser.Print{Value: parser.NewStringValue("2")},
					},
				},
				{
					Condition: parser.NewTernaryValue(ternary.FALSE),
					Statements: []parser.Statement{
						parser.Print{Value: parser.NewStringValue("3")},
					},
				},
			},
			Else: parser.Else{
				Statements: []parser.Statement{
					parser.Print{Value: parser.NewStringValue("4")},
				},
			},
		},
		ResultFlow: Terminate,
		Result:     "'4'\n",
	},
	{
		Name: "If Statement Filter Error",
		Stmt: parser.If{
			Condition: parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			Statements: []parser.Statement{
				parser.Print{Value: parser.NewStringValue("1")},
			},
		},
		Error: "field notexist does not exist",
	},
	{
		Name: "If Statement Return Value",
		Stmt: parser.If{
			Condition: parser.NewTernaryValue(ternary.TRUE),
			Statements: []parser.Statement{
				parser.Return{Value: parser.NewStringValue("1")},
			},
		},
		ResultFlow:  Return,
		ReturnValue: value.NewString("1"),
		Result:      "",
	},
}

func TestProcessor_IfStmt(t *testing.T) {
	defer func() {
		TestTx.Session.SetStdout(NewDiscard())
		initFlag(TestTx.Flags)
	}()

	TestTx.Flags.SetQuiet(true)
	tx := TestTx
	proc := NewProcessor(tx)

	for _, v := range processorIfStmtTests {
		out := NewOutput()
		tx.Session.SetStdout(out)

		proc.returnVal = nil
		flow, err := proc.IfStmt(context.Background(), v.Stmt)
		log := out.String()

		if err != nil {
			if len(v.Error) < 1 {
				t.Errorf("%s: unexpected error %q", v.Name, err)
			} else if err.Error() != v.Error {
				t.Errorf("%s: error %q, want error %q", v.Name, err.Error(), v.Error)
			}
			continue
		}
		if 0 < len(v.Error) {
			t.Errorf("%s: no error, want error %q", v.Name, v.Error)
			continue
		}
		if flow != v.ResultFlow {
			t.Errorf("%s: result flow = %q, want %q", v.Name, flow, v.ResultFlow)
		}
		if !reflect.DeepEqual(proc.returnVal, v.ReturnValue) {
			t.Errorf("%s: return = %t, want %t", v.Name, proc.returnVal, v.ReturnValue)
		}
		if log != v.Result {
			t.Errorf("%s: result = %q, want %q", v.Name, log, v.Result)
		}
	}
}

var processorCaseStmtTests = []struct {
	Name       string
	Stmt       parser.Case
	ResultFlow StatementFlow
	Result     string
	Error      string
}{
	{
		Name: "Case",
		Stmt: parser.Case{
			When: []parser.CaseWhen{
				{
					Condition: parser.NewTernaryValue(ternary.FALSE),
					Statements: []parser.Statement{
						parser.Print{Value: parser.NewStringValue("1")},
					},
				},
				{
					Condition: parser.NewTernaryValue(ternary.TRUE),
					Statements: []parser.Statement{
						parser.Print{Value: parser.NewStringValue("2")},
					},
				},
			},
		},
		ResultFlow: Terminate,
		Result:     "'2'\n",
	},
	{
		Name: "Case Comparison",
		Stmt: parser.Case{
			Value: parser.NewIntegerValue(2),
			When: []parser.CaseWhen{
				{
					Condition: parser.NewIntegerValue(1),
					Statements: []parser.Statement{
						parser.Print{Value: parser.NewStringValue("1")},
					},
				},
				{
					Condition: parser.NewIntegerValue(2),
					Statements: []parser.Statement{
						parser.Print{Value: parser.NewStringValue("2")},
					},
				},
			},
		},
		ResultFlow: Terminate,
		Result:     "'2'\n",
	},
	{
		Name: "Case Else",
		Stmt: parser.Case{
			When: []parser.CaseWhen{
				{
					Condition: parser.NewTernaryValue(ternary.FALSE),
					Statements: []parser.Statement{
						parser.Print{Value: parser.NewStringValue("1")},
					},
				},
				{
					Condition: parser.NewTernaryValue(ternary.FALSE),
					Statements: []parser.Statement{
						parser.Print{Value: parser.NewStringValue("2")},
					},
				},
			},
			Else: parser.CaseElse{
				Statements: []parser.Statement{
					parser.Print{Value: parser.NewStringValue("3")},
				},
			},
		},
		ResultFlow: Terminate,
		Result:     "'3'\n",
	},
	{
		Name: "Case No Match",
		Stmt: parser.Case{
			When: []parser.CaseWhen{
				{
					Condition: parser.NewTernaryValue(ternary.FALSE),
					Statements: []parser.Statement{
						parser.Print{Value: parser.NewStringValue("1")},
					},
				},
				{
					Condition: parser.NewTernaryValue(ternary.FALSE),
					Statements: []parser.Statement{
						parser.Print{Value: parser.NewStringValue("2")},
					},
				},
			},
		},
		ResultFlow: Terminate,
		Result:     "",
	},
	{
		Name: "Case Comparison Value Error",
		Stmt: parser.Case{
			Value: parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			When: []parser.CaseWhen{
				{
					Condition: parser.NewIntegerValue(1),
					Statements: []parser.Statement{
						parser.Print{Value: parser.NewStringValue("1")},
					},
				},
				{
					Condition: parser.NewIntegerValue(2),
					Statements: []parser.Statement{
						parser.Print{Value: parser.NewStringValue("2")},
					},
				},
			},
		},
		ResultFlow: TerminateWithError,
		Error:      "field notexist does not exist",
	},
	{
		Name: "Case Condition Error",
		Stmt: parser.Case{
			When: []parser.CaseWhen{
				{
					Condition: parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
					Statements: []parser.Statement{
						parser.Print{Value: parser.NewStringValue("1")},
					},
				},
			},
		},
		ResultFlow: TerminateWithError,
		Error:      "field notexist does not exist",
	},
}

func TestProcessor_Case(t *testing.T) {
	defer func() {
		TestTx.Session.SetStdout(NewDiscard())
		initFlag(TestTx.Flags)
	}()

	TestTx.Flags.SetQuiet(true)
	tx := TestTx
	proc := NewProcessor(tx)

	for _, v := range processorCaseStmtTests {
		out := NewOutput()
		tx.Session.SetStdout(out)
		flow, err := proc.Case(context.Background(), v.Stmt)
		log := out.String()

		if err != nil {
			if len(v.Error) < 1 {
				t.Errorf("%s: unexpected error %q", v.Name, err)
			} else if err.Error() != v.Error {
				t.Errorf("%s: error %q, want error %q", v.Name, err.Error(), v.Error)
			}
			continue
		}
		if 0 < len(v.Error) {
			t.Errorf("%s: no error, want error %q", v.Name, v.Error)
			continue
		}
		if flow != v.ResultFlow {
			t.Errorf("%s: result flow = %q, want %q", v.Name, flow, v.ResultFlow)
		}
		if log != v.Result {
			t.Errorf("%s: result = %q, want %q", v.Name, log, v.Result)
		}
	}
}

var processorWhileTests = []struct {
	Name        string
	Stmt        parser.While
	ResultFlow  StatementFlow
	ReturnValue value.Primary
	Result      string
	Error       string
}{
	{
		Name: "While Statement",
		Stmt: parser.While{
			Condition: parser.Comparison{
				LHS:      parser.Variable{Name: "while_test"},
				RHS:      parser.NewIntegerValueFromString("3"),
				Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: "<"},
			},
			Statements: []parser.Statement{
				parser.VariableSubstitution{
					Variable: parser.Variable{Name: "while_test"},
					Value: parser.Arithmetic{
						LHS:      parser.Variable{Name: "while_test"},
						RHS:      parser.NewIntegerValueFromString("1"),
						Operator: parser.Token{Token: '+', Literal: "+"},
					},
				},
				parser.Print{Value: parser.Variable{Name: "while_test"}},
				parser.TransactionControl{Token: parser.COMMIT},
			},
		},
		ResultFlow: Terminate,
		Result:     "1\n2\n3\n",
	},
	{
		Name: "While Statement Continue",
		Stmt: parser.While{
			Condition: parser.Comparison{
				LHS:      parser.Variable{Name: "while_test_count"},
				RHS:      parser.NewIntegerValueFromString("3"),
				Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: "<"},
			},
			Statements: []parser.Statement{
				parser.VariableSubstitution{
					Variable: parser.Variable{Name: "while_test_count"},
					Value: parser.Arithmetic{
						LHS:      parser.Variable{Name: "while_test_count"},
						RHS:      parser.NewIntegerValueFromString("1"),
						Operator: parser.Token{Token: '+', Literal: "+"},
					},
				},
				parser.VariableSubstitution{
					Variable: parser.Variable{Name: "while_test"},
					Value: parser.Arithmetic{
						LHS:      parser.Variable{Name: "while_test"},
						RHS:      parser.NewIntegerValueFromString("1"),
						Operator: parser.Token{Token: '+', Literal: "+"},
					},
				},
				parser.If{
					Condition: parser.Comparison{
						LHS:      parser.Variable{Name: "while_test_count"},
						RHS:      parser.NewIntegerValueFromString("2"),
						Operator: parser.Token{Token: '=', Literal: "="},
					},
					Statements: []parser.Statement{
						parser.FlowControl{Token: parser.CONTINUE},
					},
				},
				parser.Print{Value: parser.Variable{Name: "while_test"}},
				parser.TransactionControl{Token: parser.COMMIT},
			},
		},
		ResultFlow: Terminate,
		Result:     "1\n3\n",
	},
	{
		Name: "While Statement Break",
		Stmt: parser.While{
			Condition: parser.Comparison{
				LHS:      parser.Variable{Name: "while_test_count"},
				RHS:      parser.NewIntegerValueFromString("3"),
				Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: "<"},
			},
			Statements: []parser.Statement{
				parser.VariableSubstitution{
					Variable: parser.Variable{Name: "while_test_count"},
					Value: parser.Arithmetic{
						LHS:      parser.Variable{Name: "while_test_count"},
						RHS:      parser.NewIntegerValueFromString("1"),
						Operator: parser.Token{Token: '+', Literal: "+"},
					},
				},
				parser.VariableSubstitution{
					Variable: parser.Variable{Name: "while_test"},
					Value: parser.Arithmetic{
						LHS:      parser.Variable{Name: "while_test"},
						RHS:      parser.NewIntegerValueFromString("1"),
						Operator: parser.Token{Token: '+', Literal: "+"},
					},
				},
				parser.If{
					Condition: parser.Comparison{
						LHS:      parser.Variable{Name: "while_test_count"},
						RHS:      parser.NewIntegerValueFromString("2"),
						Operator: parser.Token{Token: '=', Literal: "="},
					},
					Statements: []parser.Statement{
						parser.FlowControl{Token: parser.BREAK},
					},
				},
				parser.Print{Value: parser.Variable{Name: "while_test"}},
				parser.TransactionControl{Token: parser.COMMIT},
			},
		},
		ResultFlow: Terminate,
		Result:     "1\n",
	},
	{
		Name: "While Statement Exit",
		Stmt: parser.While{
			Condition: parser.Comparison{
				LHS:      parser.Variable{Name: "while_test_count"},
				RHS:      parser.NewIntegerValueFromString("3"),
				Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: "<"},
			},
			Statements: []parser.Statement{
				parser.VariableSubstitution{
					Variable: parser.Variable{Name: "while_test_count"},
					Value: parser.Arithmetic{
						LHS:      parser.Variable{Name: "while_test_count"},
						RHS:      parser.NewIntegerValueFromString("1"),
						Operator: parser.Token{Token: '+', Literal: "+"},
					},
				},
				parser.VariableSubstitution{
					Variable: parser.Variable{Name: "while_test"},
					Value: parser.Arithmetic{
						LHS:      parser.Variable{Name: "while_test"},
						RHS:      parser.NewIntegerValueFromString("1"),
						Operator: parser.Token{Token: '+', Literal: "+"},
					},
				},
				parser.If{
					Condition: parser.Comparison{
						LHS:      parser.Variable{Name: "while_test_count"},
						RHS:      parser.NewIntegerValueFromString("2"),
						Operator: parser.Token{Token: '=', Literal: "="},
					},
					Statements: []parser.Statement{
						parser.Exit{},
					},
				},
				parser.Print{Value: parser.Variable{Name: "while_test"}},
				parser.TransactionControl{Token: parser.COMMIT},
			},
		},
		ResultFlow: Exit,
		Result:     "1\n",
	},
	{
		Name: "While Statement Filter Error",
		Stmt: parser.While{
			Condition: parser.Comparison{
				LHS:      parser.Variable{Name: "while_test"},
				RHS:      parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
				Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: "<"},
			},
			Statements: []parser.Statement{
				parser.VariableSubstitution{
					Variable: parser.Variable{Name: "while_test"},
					Value: parser.Arithmetic{
						LHS:      parser.Variable{Name: "while_test"},
						RHS:      parser.NewIntegerValueFromString("1"),
						Operator: parser.Token{Token: '+', Literal: "+"},
					},
				},
				parser.Print{Value: parser.Variable{Name: "while_test"}},
				parser.TransactionControl{Token: parser.COMMIT},
			},
		},
		Error: "field notexist does not exist",
	},
	{
		Name: "While Statement Execution Error",
		Stmt: parser.While{
			Condition: parser.Comparison{
				LHS:      parser.Variable{Name: "while_test"},
				RHS:      parser.NewIntegerValueFromString("3"),
				Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: "<"},
			},
			Statements: []parser.Statement{
				parser.VariableSubstitution{
					Variable: parser.Variable{Name: "while_test"},
					Value: parser.Arithmetic{
						LHS:      parser.Variable{Name: "while_test"},
						RHS:      parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
						Operator: parser.Token{Token: '+', Literal: "+"},
					},
				},
				parser.Print{Value: parser.Variable{Name: "while_test"}},
				parser.TransactionControl{Token: parser.COMMIT},
			},
		},
		Error: "field notexist does not exist",
	},
	{
		Name: "While Statement Return Value",
		Stmt: parser.While{
			Condition: parser.Comparison{
				LHS:      parser.Variable{Name: "while_test"},
				RHS:      parser.NewIntegerValueFromString("3"),
				Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: "<"},
			},
			Statements: []parser.Statement{
				parser.Return{Value: parser.NewStringValue("1")},
				parser.VariableSubstitution{
					Variable: parser.Variable{Name: "while_test"},
					Value: parser.Arithmetic{
						LHS:      parser.Variable{Name: "while_test"},
						RHS:      parser.NewIntegerValueFromString("1"),
						Operator: parser.Token{Token: '+', Literal: "+"},
					},
				},
				parser.Print{Value: parser.Variable{Name: "while_test"}},
				parser.TransactionControl{Token: parser.COMMIT},
			},
		},
		ResultFlow:  Return,
		ReturnValue: value.NewString("1"),
	},
}

func TestProcessor_While(t *testing.T) {
	defer func() {
		TestTx.Session.SetStdout(NewDiscard())
		initFlag(TestTx.Flags)
	}()

	TestTx.Flags.SetQuiet(true)
	tx := TestTx
	proc := NewProcessor(tx)

	for _, v := range processorWhileTests {
		proc.returnVal = nil
		if _, ok := proc.ReferenceScope.CurrentBlock().variables.Get(parser.Variable{Name: "while_test"}); !ok {
			_ = proc.ReferenceScope.DeclareVariableDirectly(parser.Variable{Name: "while_test"}, value.NewInteger(0))
		}
		_ = proc.ReferenceScope.CurrentBlock().variables.Set(parser.Variable{Name: "while_test"}, value.NewInteger(0))

		if _, ok := proc.ReferenceScope.CurrentBlock().variables.Get(parser.Variable{Name: "while_test_count"}); !ok {
			_ = proc.ReferenceScope.DeclareVariableDirectly(parser.Variable{Name: "while_test_count"}, value.NewInteger(0))
		}
		_ = proc.ReferenceScope.CurrentBlock().variables.Set(parser.Variable{Name: "while_test_count"}, value.NewInteger(0))

		out := NewOutput()
		tx.Session.SetStdout(out)
		flow, err := proc.While(context.Background(), v.Stmt)
		log := out.String()

		if err != nil {
			if len(v.Error) < 1 {
				t.Errorf("%s: unexpected error %q", v.Name, err)
			} else if err.Error() != v.Error {
				t.Errorf("%s: error %q, want error %q", v.Name, err.Error(), v.Error)
			}
			continue
		}
		if 0 < len(v.Error) {
			t.Errorf("%s: no error, want error %q", v.Name, v.Error)
			continue
		}
		if flow != v.ResultFlow {
			t.Errorf("%s: result flow = %q, want %q", v.Name, flow, v.ResultFlow)
		}
		if !reflect.DeepEqual(proc.returnVal, v.ReturnValue) {
			t.Errorf("%s: return = %t, want %t", v.Name, proc.returnVal, v.ReturnValue)
		}
		if log != v.Result {
			t.Errorf("%s: result = %q, want %q", v.Name, log, v.Result)
		}
	}
}

var processorWhileInCursorTests = []struct {
	Name        string
	Stmt        parser.WhileInCursor
	ResultFlow  StatementFlow
	ReturnValue value.Primary
	Result      string
	Error       string
}{
	{
		Name: "While In Cursor",
		Stmt: parser.WhileInCursor{
			Variables: []parser.Variable{
				{Name: "var1"},
				{Name: "var2"},
			},
			Cursor: parser.Identifier{Literal: "cur"},
			Statements: []parser.Statement{
				parser.Print{Value: parser.Variable{Name: "var1"}},
				parser.TransactionControl{Token: parser.COMMIT},
			},
		},
		ResultFlow: Terminate,
		Result:     "'1'\n'2'\n'3'\n",
	},
	{
		Name: "While In Cursor With Declaration",
		Stmt: parser.WhileInCursor{
			WithDeclaration: true,
			Variables: []parser.Variable{
				{Name: "declvar1"},
				{Name: "declvar2"},
			},
			Cursor: parser.Identifier{Literal: "cur"},
			Statements: []parser.Statement{
				parser.Print{Value: parser.Variable{Name: "declvar1"}},
				parser.TransactionControl{Token: parser.COMMIT},
			},
		},
		ResultFlow: Terminate,
		Result:     "'1'\n'2'\n'3'\n",
	},
	{
		Name: "While In Cursor Continue",
		Stmt: parser.WhileInCursor{
			Variables: []parser.Variable{
				{Name: "var1"},
				{Name: "var2"},
			},
			Cursor: parser.Identifier{Literal: "cur"},
			Statements: []parser.Statement{
				parser.If{
					Condition: parser.Comparison{
						LHS:      parser.Variable{Name: "var1"},
						RHS:      parser.NewIntegerValueFromString("2"),
						Operator: parser.Token{Token: '=', Literal: "="},
					},
					Statements: []parser.Statement{
						parser.FlowControl{Token: parser.CONTINUE},
					},
				},
				parser.Print{Value: parser.Variable{Name: "var1"}},
				parser.TransactionControl{Token: parser.COMMIT},
			},
		},
		ResultFlow: Terminate,
		Result:     "'1'\n'3'\n",
	},
	{
		Name: "While In Cursor Break",
		Stmt: parser.WhileInCursor{
			Variables: []parser.Variable{
				{Name: "var1"},
				{Name: "var2"},
			},
			Cursor: parser.Identifier{Literal: "cur"},
			Statements: []parser.Statement{
				parser.If{
					Condition: parser.Comparison{
						LHS:      parser.Variable{Name: "var1"},
						RHS:      parser.NewIntegerValueFromString("2"),
						Operator: parser.Token{Token: '=', Literal: "="},
					},
					Statements: []parser.Statement{
						parser.FlowControl{Token: parser.BREAK},
					},
				},
				parser.Print{Value: parser.Variable{Name: "var1"}},
				parser.TransactionControl{Token: parser.COMMIT},
			},
		},
		ResultFlow: Terminate,
		Result:     "'1'\n",
	},
	{
		Name: "While In Cursor Exit With Code",
		Stmt: parser.WhileInCursor{
			Variables: []parser.Variable{
				{Name: "var1"},
				{Name: "var2"},
			},
			Cursor: parser.Identifier{Literal: "cur"},
			Statements: []parser.Statement{
				parser.If{
					Condition: parser.Comparison{
						LHS:      parser.Variable{Name: "var1"},
						RHS:      parser.NewIntegerValueFromString("2"),
						Operator: parser.Token{Token: '=', Literal: "="},
					},
					Statements: []parser.Statement{
						parser.Exit{},
					},
				},
				parser.Print{Value: parser.Variable{Name: "var1"}},
				parser.TransactionControl{Token: parser.COMMIT},
			},
		},
		ResultFlow: Exit,
		Result:     "'1'\n",
	},
	{
		Name: "While In Cursor Fetch Error",
		Stmt: parser.WhileInCursor{
			Variables: []parser.Variable{
				{Name: "var1"},
				{Name: "var3"},
			},
			Cursor: parser.Identifier{Literal: "cur"},
			Statements: []parser.Statement{
				parser.Print{Value: parser.Variable{Name: "var1"}},
				parser.TransactionControl{Token: parser.COMMIT},
			},
		},
		Error: "variable @var3 is undeclared",
	},
	{
		Name: "While In Cursor Statement Execution Error",
		Stmt: parser.WhileInCursor{
			Variables: []parser.Variable{
				{Name: "var1"},
				{Name: "var2"},
			},
			Cursor: parser.Identifier{Literal: "cur"},
			Statements: []parser.Statement{
				parser.If{
					Condition: parser.Comparison{
						LHS:      parser.Variable{Name: "var1"},
						RHS:      parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
						Operator: parser.Token{Token: '=', Literal: "="},
					},
					Statements: []parser.Statement{
						parser.FlowControl{Token: parser.BREAK},
					},
				},
				parser.Print{Value: parser.Variable{Name: "var1"}},
				parser.TransactionControl{Token: parser.COMMIT},
			},
		},
		Error: "field notexist does not exist",
	},
	{
		Name: "While In Cursor Return Value",
		Stmt: parser.WhileInCursor{
			Variables: []parser.Variable{
				{Name: "var1"},
				{Name: "var2"},
			},
			Cursor: parser.Identifier{Literal: "cur"},
			Statements: []parser.Statement{
				parser.Return{Value: parser.NewStringValue("1")},
				parser.Print{Value: parser.Variable{Name: "var1"}},
				parser.TransactionControl{Token: parser.COMMIT},
			},
		},
		ResultFlow:  Return,
		ReturnValue: value.NewString("1"),
	},
}

func TestProcessor_WhileInCursor(t *testing.T) {
	defer func() {
		_ = TestTx.cachedViews.Clean(TestTx.FileContainer)
		TestTx.Session.SetStdout(NewDiscard())
		initFlag(TestTx.Flags)
	}()

	TestTx.Flags.Repository = TestDir

	tx := TestTx
	proc := NewProcessor(tx)
	ctx := context.Background()

	for _, v := range processorWhileInCursorTests {
		proc.ReferenceScope = GenerateReferenceScope([]map[string]map[string]interface{}{
			{
				scopeNameVariables: {
					"var1": value.NewNull(),
					"var2": value.NewNull(),
				},
				scopeNameCursors: {
					"CUR": &Cursor{
						query: selectQueryForCursorTest,
						mtx:   &sync.Mutex{},
					},
				},
			},
		}, nil, time.Time{}, nil)
		_ = TestTx.cachedViews.Clean(TestTx.FileContainer)
		_ = proc.ReferenceScope.OpenCursor(ctx, parser.Identifier{Literal: "cur"}, nil)

		out := NewOutput()
		tx.Session.SetStdout(out)
		flow, err := proc.WhileInCursor(context.Background(), v.Stmt)
		log := out.String()

		if err != nil {
			if len(v.Error) < 1 {
				t.Errorf("%s: unexpected error %q", v.Name, err)
			} else if err.Error() != v.Error {
				t.Errorf("%s: error %q, want error %q", v.Name, err.Error(), v.Error)
			}
			continue
		}
		if 0 < len(v.Error) {
			t.Errorf("%s: no error, want error %q", v.Name, v.Error)
			continue
		}
		if flow != v.ResultFlow {
			t.Errorf("%s: result flow = %q, want %q", v.Name, flow, v.ResultFlow)
		}
		if !reflect.DeepEqual(proc.returnVal, v.ReturnValue) {
			t.Errorf("%s: return = %t, want %t", v.Name, proc.returnVal, v.ReturnValue)
		}
		if log != v.Result {
			t.Errorf("%s: result = %q, want %q", v.Name, log, v.Result)
		}
	}
}

var processorExecExternalCommand = []struct {
	Name  string
	Stmt  parser.ExternalCommand
	Error string
}{
	{
		Name: "Error in Splitting Arguments",
		Stmt: parser.ExternalCommand{
			Command: "cmd arg 'arg",
		},
		Error: "external command: string not terminated",
	},
	{
		Name: "Error in Scanning Argument",
		Stmt: parser.ExternalCommand{
			Command: "cmd 'arg arg@'",
		},
		Error: "external command: invalid variable symbol",
	},
	{
		Name: "Error in Evaluation of Variable",
		Stmt: parser.ExternalCommand{
			Command: "cmd @__not_exist__",
		},
		Error: "external command: variable @__not_exist__ is undeclared",
	},
}

func TestProcessor_ExecExternalCommand(t *testing.T) {
	proc := NewProcessor(TestTx)

	for _, v := range processorExecExternalCommand {
		err := proc.ExecExternalCommand(context.Background(), v.Stmt)

		if err != nil {
			if len(v.Error) < 1 {
				t.Errorf("%s: unexpected error %q", v.Name, err)
			} else if err.Error() != v.Error {
				t.Errorf("%s: error %q, want error %q", v.Name, err.Error(), v.Error)
			}
			continue
		}
		if 0 < len(v.Error) {
			t.Errorf("%s: no error, want error %q", v.Name, v.Error)
			continue
		}
	}
}
