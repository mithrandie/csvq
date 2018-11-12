package query

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/mithrandie/go-text"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/file"
	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/value"

	"github.com/mithrandie/ternary"
)

var procedureExecuteStatementTests = []struct {
	Input      parser.Statement
	Result     []ExecResult
	Logs       string
	SelectLogs []string
	Error      string
	ErrorCode  int
}{
	{
		Input: parser.SetFlag{
			Name:  "@@invalid",
			Value: parser.NewStringValue("\t"),
		},
		Error:     "[L:- C:-] flag @@invalid does not exist",
		ErrorCode: 1,
	},
	{
		Input: parser.SetFlag{
			Name:  "@@delimiter",
			Value: parser.NewStringValue(","),
		},
		Logs: " @@DELIMITER: ',' | SPACES\n",
	},
	{
		Input: parser.ShowFlag{
			Name: "@@repository",
		},
		Logs: " @@REPOSITORY: " + TestDir + "\n",
	},
	{
		Input: parser.VariableDeclaration{
			Assignments: []parser.VariableAssignment{
				{
					Variable: parser.Variable{Name: "@var1"},
				},
			},
		},
		Result: []ExecResult{},
	},
	{
		Input: parser.VariableDeclaration{
			Assignments: []parser.VariableAssignment{
				{
					Variable: parser.Variable{Name: "@var2"},
				},
			},
		},
		Result: []ExecResult{},
	},
	{
		Input: parser.VariableDeclaration{
			Assignments: []parser.VariableAssignment{
				{
					Variable: parser.Variable{Name: "@var3"},
				},
			},
		},
		Result: []ExecResult{},
	},
	{
		Input: parser.VariableDeclaration{
			Assignments: []parser.VariableAssignment{
				{
					Variable: parser.Variable{Name: "@var4"},
				},
			},
		},
		Result: []ExecResult{},
	},
	{
		Input: parser.VariableSubstitution{
			Variable: parser.Variable{Name: "@var1"},
			Value:    parser.NewIntegerValueFromString("1"),
		},
		Result: []ExecResult{},
	},
	{
		Input: parser.Print{
			Value: parser.Variable{Name: "@var1"},
		},
		Logs: "1\n",
	},
	{
		Input: parser.DisposeVariable{
			Variable: parser.Variable{Name: "@var4"},
		},
	},
	{
		Input: parser.VariableDeclaration{
			Assignments: []parser.VariableAssignment{
				{
					Variable: parser.Variable{Name: "@var4"},
				},
			},
		},
	},
	{
		Input: parser.FunctionDeclaration{
			Name: parser.Identifier{Literal: "userfunc"},
			Parameters: []parser.VariableAssignment{
				{
					Variable: parser.Variable{Name: "@arg1"},
				},
			},
			Statements: []parser.Statement{
				parser.Print{
					Value: parser.Variable{Name: "@arg1"},
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
		Error:     "[L:- C:-] function userfunc does not exist",
		ErrorCode: 1,
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
				{Name: "@var2"},
				{Name: "@var3"},
			},
		},
	},
	{
		Input: parser.Print{
			Value: parser.Variable{Name: "@var2"},
		},
		Logs: "\"1\"\n",
	},
	{
		Input: parser.Print{
			Value: parser.Variable{Name: "@var3"},
		},
		Logs: "\"str1\"\n",
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
		Logs: "\"column1\",\"column2\"\n1,2\n",
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
							Variable: parser.Variable{Name: "@value"},
						},
						{
							Variable: parser.Variable{Name: "@fetch"},
						},
					},
				},
				parser.WhileInCursor{
					Variables: []parser.Variable{
						{Name: "@fetch"},
					},
					Cursor: parser.Identifier{Literal: "list"},
					Statements: []parser.Statement{
						parser.If{
							Condition: parser.Is{
								LHS: parser.Variable{Name: "@fetch"},
								RHS: parser.NewNullValue(),
							},
							Statements: []parser.Statement{
								parser.FlowControl{Token: parser.CONTINUE},
							},
						},
						parser.If{
							Condition: parser.Is{
								LHS: parser.Variable{Name: "@value"},
								RHS: parser.NewNullValue(),
							},
							Statements: []parser.Statement{
								parser.VariableSubstitution{
									Variable: parser.Variable{Name: "@value"},
									Value:    parser.Variable{Name: "@fetch"},
								},
								parser.FlowControl{Token: parser.CONTINUE},
							},
						},
						parser.VariableSubstitution{
							Variable: parser.Variable{Name: "@value"},
							Value: parser.Arithmetic{
								LHS:      parser.Variable{Name: "@value"},
								RHS:      parser.Variable{Name: "@fetch"},
								Operator: '*',
							},
						},
					},
				},
				parser.Return{
					Value: parser.Variable{Name: "@value"},
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
		Logs: "\"multiplication\"\n6\n",
	},
	{
		Input: parser.SelectQuery{
			SelectEntity: parser.SelectEntity{
				SelectClause: parser.SelectClause{
					Fields: []parser.QueryExpression{
						parser.Field{
							Object: parser.Variable{Name: "@var1"},
							Alias:  parser.Identifier{Literal: "var1"},
						},
					},
				},
			},
		},
		Logs: "\"var1\"\n1\n",
	},
	{
		Input: parser.VariableDeclaration{
			Assignments: []parser.VariableAssignment{
				{
					Variable: parser.Variable{Name: "@var1"},
				},
			},
		},
		Error:     "[L:- C:-] variable @var1 is redeclared",
		ErrorCode: 1,
	},
	{
		Input: parser.VariableSubstitution{
			Variable: parser.Variable{Name: "@var9"},
			Value:    parser.NewIntegerValueFromString("1"),
		},
		Error:     "[L:- C:-] variable @var9 is undeclared",
		ErrorCode: 1,
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
		Result: []ExecResult{
			{
				Type: InsertQuery,
				FileInfo: &FileInfo{
					Path:      GetTestFilePath("table1.csv"),
					Delimiter: ',',
					NoHeader:  false,
					Encoding:  text.UTF8,
					LineBreak: text.LF,
				},
				OperatedCount: 2,
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
					Operator: "=",
				},
			},
		},
		Result: []ExecResult{
			{
				Type: UpdateQuery,
				FileInfo: &FileInfo{
					Path:      GetTestFilePath("table1.csv"),
					Delimiter: ',',
					NoHeader:  false,
					Encoding:  text.UTF8,
					LineBreak: text.LF,
				},
				OperatedCount: 1,
			},
		},
		Logs: fmt.Sprintf("1 record updated on %q.\n", GetTestFilePath("table1.csv")),
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
					Operator: "=",
				},
			},
		},
		Result: []ExecResult{
			{
				Type: DeleteQuery,
				FileInfo: &FileInfo{
					Path:      GetTestFilePath("table1.csv"),
					Delimiter: ',',
					NoHeader:  false,
					Encoding:  text.UTF8,
					LineBreak: text.LF,
				},
				OperatedCount: 1,
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
		Result: []ExecResult{
			{
				Type: CreateTableQuery,
				FileInfo: &FileInfo{
					Path:      GetTestFilePath("newtable.csv"),
					Delimiter: ',',
					NoHeader:  false,
					Encoding:  text.UTF8,
					LineBreak: text.LF,
				},
			},
		},
		Logs: fmt.Sprintf("file %q is created.\n", GetTestFilePath("newtable.csv")),
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
		Result: []ExecResult{
			{
				Type: AddColumnsQuery,
				FileInfo: &FileInfo{
					Path:      GetTestFilePath("table1.csv"),
					Delimiter: ',',
					NoHeader:  false,
					Encoding:  text.UTF8,
					LineBreak: text.LF,
				},
				OperatedCount: 1,
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
		Result: []ExecResult{
			{
				Type: DropColumnsQuery,
				FileInfo: &FileInfo{
					Path:      GetTestFilePath("table1.csv"),
					Delimiter: ',',
					NoHeader:  false,
					Encoding:  text.UTF8,
					LineBreak: text.LF,
				},
				OperatedCount: 1,
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
		Result: []ExecResult{
			{
				Type: RenameColumnQuery,
				FileInfo: &FileInfo{
					Path:      GetTestFilePath("table1.csv"),
					Delimiter: ',',
					NoHeader:  false,
					Encoding:  text.UTF8,
					LineBreak: text.LF,
				},
				OperatedCount: 1,
			},
		},
		Logs: fmt.Sprintf("1 field renamed on %q.\n", GetTestFilePath("table1.csv")),
	},
	{
		Input: parser.SetTableAttribute{
			Table:     parser.Identifier{Literal: "table1.csv"},
			Attribute: parser.Identifier{Literal: "delimiter"},
			Value:     parser.NewStringValue("\t"),
		},
		Logs: "\n" +
			strings.Repeat(" ", (calcShowFieldsWidth("table1.csv", "table1.csv", 22)-(22+len("table1.csv")))/2) + "Attributes Updated in table1.csv\n" +
			strings.Repeat("-", calcShowFieldsWidth("table1.csv", "table1.csv", 22)) + "\n" +
			" Path: " + GetTestFilePath("table1.csv") + "\n" +
			" Format: TSV     Delimiter: '\\t'\n" +
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
		Logs: "\"2\"\n",
	},
	{
		Input: parser.While{
			Condition: parser.Comparison{
				LHS:      parser.Variable{Name: "@while_test"},
				RHS:      parser.NewIntegerValueFromString("3"),
				Operator: "<",
			},
			Statements: []parser.Statement{
				parser.VariableSubstitution{
					Variable: parser.Variable{Name: "@while_test"},
					Value: parser.Arithmetic{
						LHS:      parser.Variable{Name: "@while_test"},
						RHS:      parser.NewIntegerValueFromString("1"),
						Operator: '+',
					},
				},
				parser.Print{Value: parser.Variable{Name: "@while_test"}},
			},
		},
		Logs: "1\n2\n3\n",
	},
	{
		Input: parser.Exit{
			Code: value.NewInteger(1),
		},
		Error:     "",
		ErrorCode: 1,
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
		Logs: "\"external executable file\"\n",
	},
	{
		Input: parser.Execute{
			BaseExpr:   parser.NewBaseExpr(parser.Token{}),
			Statements: parser.NewStringValue("print 'execute';"),
		},
		Logs: "\"execute\"\n",
	},
	{
		Input: parser.Trigger{
			Event:   parser.Identifier{Literal: "error"},
			Message: parser.NewStringValue("user error"),
			Code:    value.NewInteger(200),
		},
		Error:     "[L:- C:-] user error",
		ErrorCode: 200,
	},
	{
		Input: parser.Trigger{
			Event:   parser.Identifier{Literal: "error"},
			Message: parser.NewIntegerValue(200),
		},
		Error:     "[L:- C:-] ",
		ErrorCode: 200,
	},
	{
		Input: parser.Trigger{
			Event:   parser.Identifier{Literal: "invalid"},
			Message: parser.NewIntegerValue(200),
		},
		Error:     "[L:- C:-] invalid is an unknown event",
		ErrorCode: 1,
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
			" Format: CSV     Delimiter: ','\n" +
			" Encoding: UTF8  LineBreak: LF    Header: true\n" +
			" Status: Fixed\n" +
			" Fields:\n" +
			"   1. column1\n" +
			"   2. column2\n" +
			"\n",
	},
}

func TestProcedure_ExecuteStatement(t *testing.T) {
	initFlag()
	tf := cmd.GetFlags()
	tf.Repository = TestDir
	tf.Format = cmd.CSV

	proc := NewProcedure()
	proc.Filter.Variables[0].Add(parser.Variable{Name: "@while_test"}, value.NewInteger(0))

	for _, v := range procedureExecuteStatementTests {
		ReleaseResources()
		ExecResults = []ExecResult{}
		SelectLogs = []string{}

		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		_, err := proc.ExecuteStatement(v.Input)

		w.Close()
		os.Stdout = oldStdout

		log, _ := ioutil.ReadAll(r)

		if err != nil {
			var code int
			if apperr, ok := err.(AppError); ok {
				if len(v.Error) < 1 {
					t.Errorf("unexpected error %q for %q", err, v.Input)
				} else if err.Error() != v.Error {
					t.Errorf("error %q, want error %q for %q", err, v.Error, v.Input)
				}

				code = apperr.GetCode()
			} else if ex, ok := err.(*ForcedExit); ok {
				code = ex.GetCode()
			}
			if code != v.ErrorCode {
				t.Errorf("error code %d, want error code %d for %q", code, v.ErrorCode, v.Input)
			}
			continue
		}
		if 0 < len(v.Error) {
			t.Errorf("no error, want error %q for %q", v.Error, v.Input)
			continue
		}

		if v.Result != nil {
			for _, r := range ExecResults {
				if r.FileInfo.File != nil {
					if r.FileInfo.Path != r.FileInfo.File.Name() {
						t.Errorf("file pointer = %q, want %q for %q", r.FileInfo.File.Name(), r.FileInfo.Path, v.Input)
					}
					file.Close(r.FileInfo.File)
					r.FileInfo.File = nil
				}
			}

			if !reflect.DeepEqual(ExecResults, v.Result) {
				t.Errorf("results = %v, want %v for %q", ExecResults, v.Result, v.Input)
			}
		}
		if 0 < len(v.Logs) {
			if string(log) != v.Logs {
				t.Errorf("logs = %s, want %s for %q", string(log), v.Logs, v.Input)
			}
		}
		if v.SelectLogs != nil {
			if !reflect.DeepEqual(SelectLogs, v.SelectLogs) {
				t.Errorf("select logs = %s, want %s for %q", SelectLogs, v.SelectLogs, v.Input)
			}
		}
	}

	ReleaseResources()
	ExecResults = []ExecResult{}
	SelectLogs = []string{}
}

var procedureIfStmtTests = []struct {
	Name       string
	Stmt       parser.If
	ResultFlow StatementFlow
	Result     string
	Error      string
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
		Result:     "\"1\"\n",
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
		Result:     "\"2\"\n",
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
		Result:     "\"4\"\n",
	},
	{
		Name: "If Statement Filter Error",
		Stmt: parser.If{
			Condition: parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			Statements: []parser.Statement{
				parser.Print{Value: parser.NewStringValue("1")},
			},
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
}

func TestProcedure_IfStmt(t *testing.T) {
	cmd.SetQuiet(true)
	proc := NewProcedure()

	for _, v := range procedureIfStmtTests {
		oldStdout := os.Stdout

		r, w, _ := os.Pipe()
		os.Stdout = w

		flow, err := proc.IfStmt(v.Stmt)

		w.Close()
		os.Stdout = oldStdout

		log, _ := ioutil.ReadAll(r)

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
		if string(log) != v.Result {
			t.Errorf("%s: result = %q, want %q", v.Name, string(log), v.Result)
		}
	}
}

var procedureCaseStmtTests = []struct {
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
		Result:     "\"2\"\n",
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
		Result:     "\"2\"\n",
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
		Result:     "\"3\"\n",
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
		ResultFlow: Error,
		Error:      "[L:- C:-] field notexist does not exist",
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
		ResultFlow: Error,
		Error:      "[L:- C:-] field notexist does not exist",
	},
}

func TestProcedure_Case(t *testing.T) {
	cmd.SetQuiet(true)
	proc := NewProcedure()

	for _, v := range procedureCaseStmtTests {
		oldStdout := os.Stdout

		r, w, _ := os.Pipe()
		os.Stdout = w

		flow, err := proc.Case(v.Stmt)

		w.Close()
		os.Stdout = oldStdout

		log, _ := ioutil.ReadAll(r)

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
		if string(log) != v.Result {
			t.Errorf("%s: result = %q, want %q", v.Name, string(log), v.Result)
		}
	}
}

var procedureWhileTests = []struct {
	Name       string
	Stmt       parser.While
	ResultFlow StatementFlow
	Result     string
	Error      string
}{
	{
		Name: "While Statement",
		Stmt: parser.While{
			Condition: parser.Comparison{
				LHS:      parser.Variable{Name: "@while_test"},
				RHS:      parser.NewIntegerValueFromString("3"),
				Operator: "<",
			},
			Statements: []parser.Statement{
				parser.VariableSubstitution{
					Variable: parser.Variable{Name: "@while_test"},
					Value: parser.Arithmetic{
						LHS:      parser.Variable{Name: "@while_test"},
						RHS:      parser.NewIntegerValueFromString("1"),
						Operator: '+',
					},
				},
				parser.Print{Value: parser.Variable{Name: "@while_test"}},
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
				LHS:      parser.Variable{Name: "@while_test_count"},
				RHS:      parser.NewIntegerValueFromString("3"),
				Operator: "<",
			},
			Statements: []parser.Statement{
				parser.VariableSubstitution{
					Variable: parser.Variable{Name: "@while_test_count"},
					Value: parser.Arithmetic{
						LHS:      parser.Variable{Name: "@while_test_count"},
						RHS:      parser.NewIntegerValueFromString("1"),
						Operator: '+',
					},
				},
				parser.VariableSubstitution{
					Variable: parser.Variable{Name: "@while_test"},
					Value: parser.Arithmetic{
						LHS:      parser.Variable{Name: "@while_test"},
						RHS:      parser.NewIntegerValueFromString("1"),
						Operator: '+',
					},
				},
				parser.If{
					Condition: parser.Comparison{
						LHS:      parser.Variable{Name: "@while_test_count"},
						RHS:      parser.NewIntegerValueFromString("2"),
						Operator: "=",
					},
					Statements: []parser.Statement{
						parser.FlowControl{Token: parser.CONTINUE},
					},
				},
				parser.Print{Value: parser.Variable{Name: "@while_test"}},
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
				LHS:      parser.Variable{Name: "@while_test_count"},
				RHS:      parser.NewIntegerValueFromString("3"),
				Operator: "<",
			},
			Statements: []parser.Statement{
				parser.VariableSubstitution{
					Variable: parser.Variable{Name: "@while_test_count"},
					Value: parser.Arithmetic{
						LHS:      parser.Variable{Name: "@while_test_count"},
						RHS:      parser.NewIntegerValueFromString("1"),
						Operator: '+',
					},
				},
				parser.VariableSubstitution{
					Variable: parser.Variable{Name: "@while_test"},
					Value: parser.Arithmetic{
						LHS:      parser.Variable{Name: "@while_test"},
						RHS:      parser.NewIntegerValueFromString("1"),
						Operator: '+',
					},
				},
				parser.If{
					Condition: parser.Comparison{
						LHS:      parser.Variable{Name: "@while_test_count"},
						RHS:      parser.NewIntegerValueFromString("2"),
						Operator: "=",
					},
					Statements: []parser.Statement{
						parser.FlowControl{Token: parser.BREAK},
					},
				},
				parser.Print{Value: parser.Variable{Name: "@while_test"}},
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
				LHS:      parser.Variable{Name: "@while_test_count"},
				RHS:      parser.NewIntegerValueFromString("3"),
				Operator: "<",
			},
			Statements: []parser.Statement{
				parser.VariableSubstitution{
					Variable: parser.Variable{Name: "@while_test_count"},
					Value: parser.Arithmetic{
						LHS:      parser.Variable{Name: "@while_test_count"},
						RHS:      parser.NewIntegerValueFromString("1"),
						Operator: '+',
					},
				},
				parser.VariableSubstitution{
					Variable: parser.Variable{Name: "@while_test"},
					Value: parser.Arithmetic{
						LHS:      parser.Variable{Name: "@while_test"},
						RHS:      parser.NewIntegerValueFromString("1"),
						Operator: '+',
					},
				},
				parser.If{
					Condition: parser.Comparison{
						LHS:      parser.Variable{Name: "@while_test_count"},
						RHS:      parser.NewIntegerValueFromString("2"),
						Operator: "=",
					},
					Statements: []parser.Statement{
						parser.Exit{},
					},
				},
				parser.Print{Value: parser.Variable{Name: "@while_test"}},
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
				LHS:      parser.Variable{Name: "@while_test"},
				RHS:      parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
				Operator: "<",
			},
			Statements: []parser.Statement{
				parser.VariableSubstitution{
					Variable: parser.Variable{Name: "@while_test"},
					Value: parser.Arithmetic{
						LHS:      parser.Variable{Name: "@while_test"},
						RHS:      parser.NewIntegerValueFromString("1"),
						Operator: '+',
					},
				},
				parser.Print{Value: parser.Variable{Name: "@while_test"}},
				parser.TransactionControl{Token: parser.COMMIT},
			},
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
	{
		Name: "While Statement Execution Error",
		Stmt: parser.While{
			Condition: parser.Comparison{
				LHS:      parser.Variable{Name: "@while_test"},
				RHS:      parser.NewIntegerValueFromString("3"),
				Operator: "<",
			},
			Statements: []parser.Statement{
				parser.VariableSubstitution{
					Variable: parser.Variable{Name: "@while_test"},
					Value: parser.Arithmetic{
						LHS:      parser.Variable{Name: "@while_test"},
						RHS:      parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
						Operator: '+',
					},
				},
				parser.Print{Value: parser.Variable{Name: "@while_test"}},
				parser.TransactionControl{Token: parser.COMMIT},
			},
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
}

func TestProcedure_While(t *testing.T) {
	cmd.SetQuiet(true)
	proc := NewProcedure()

	for _, v := range procedureWhileTests {
		if _, err := proc.Filter.Variables[0].Get(parser.Variable{Name: "@while_test"}); err != nil {
			proc.Filter.Variables[0].Add(parser.Variable{Name: "@while_test"}, value.NewInteger(0))
		}
		proc.Filter.Variables[0].Set(parser.Variable{Name: "@while_test"}, value.NewInteger(0))

		if _, err := proc.Filter.Variables[0].Get(parser.Variable{Name: "@while_test_count"}); err != nil {
			proc.Filter.Variables[0].Add(parser.Variable{Name: "@while_test_count"}, value.NewInteger(0))
		}
		proc.Filter.Variables[0].Set(parser.Variable{Name: "@while_test_count"}, value.NewInteger(0))

		oldStdout := os.Stdout

		r, w, _ := os.Pipe()
		os.Stdout = w

		flow, err := proc.While(v.Stmt)

		w.Close()
		os.Stdout = oldStdout

		log, _ := ioutil.ReadAll(r)

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
		if string(log) != v.Result {
			t.Errorf("%s: result = %q, want %q", v.Name, string(log), v.Result)
		}
	}
}

var procedureWhileInCursorTests = []struct {
	Name       string
	Stmt       parser.WhileInCursor
	ResultFlow StatementFlow
	Result     string
	Error      string
}{
	{
		Name: "While In Cursor",
		Stmt: parser.WhileInCursor{
			Variables: []parser.Variable{
				{Name: "@var1"},
				{Name: "@var2"},
			},
			Cursor: parser.Identifier{Literal: "cur"},
			Statements: []parser.Statement{
				parser.Print{Value: parser.Variable{Name: "@var1"}},
				parser.TransactionControl{Token: parser.COMMIT},
			},
		},
		ResultFlow: Terminate,
		Result:     "\"1\"\n\"2\"\n\"3\"\n",
	},
	{
		Name: "While In Cursor With Declaration",
		Stmt: parser.WhileInCursor{
			WithDeclaration: true,
			Variables: []parser.Variable{
				{Name: "@declvar1"},
				{Name: "@declvar2"},
			},
			Cursor: parser.Identifier{Literal: "cur"},
			Statements: []parser.Statement{
				parser.Print{Value: parser.Variable{Name: "@declvar1"}},
				parser.TransactionControl{Token: parser.COMMIT},
			},
		},
		ResultFlow: Terminate,
		Result:     "\"1\"\n\"2\"\n\"3\"\n",
	},
	{
		Name: "While In Cursor Continue",
		Stmt: parser.WhileInCursor{
			Variables: []parser.Variable{
				{Name: "@var1"},
				{Name: "@var2"},
			},
			Cursor: parser.Identifier{Literal: "cur"},
			Statements: []parser.Statement{
				parser.If{
					Condition: parser.Comparison{
						LHS:      parser.Variable{Name: "@var1"},
						RHS:      parser.NewIntegerValueFromString("2"),
						Operator: "=",
					},
					Statements: []parser.Statement{
						parser.FlowControl{Token: parser.CONTINUE},
					},
				},
				parser.Print{Value: parser.Variable{Name: "@var1"}},
				parser.TransactionControl{Token: parser.COMMIT},
			},
		},
		ResultFlow: Terminate,
		Result:     "\"1\"\n\"3\"\n",
	},
	{
		Name: "While In Cursor Break",
		Stmt: parser.WhileInCursor{
			Variables: []parser.Variable{
				{Name: "@var1"},
				{Name: "@var2"},
			},
			Cursor: parser.Identifier{Literal: "cur"},
			Statements: []parser.Statement{
				parser.If{
					Condition: parser.Comparison{
						LHS:      parser.Variable{Name: "@var1"},
						RHS:      parser.NewIntegerValueFromString("2"),
						Operator: "=",
					},
					Statements: []parser.Statement{
						parser.FlowControl{Token: parser.BREAK},
					},
				},
				parser.Print{Value: parser.Variable{Name: "@var1"}},
				parser.TransactionControl{Token: parser.COMMIT},
			},
		},
		ResultFlow: Terminate,
		Result:     "\"1\"\n",
	},
	{
		Name: "While In Cursor Exit With Code",
		Stmt: parser.WhileInCursor{
			Variables: []parser.Variable{
				{Name: "@var1"},
				{Name: "@var2"},
			},
			Cursor: parser.Identifier{Literal: "cur"},
			Statements: []parser.Statement{
				parser.If{
					Condition: parser.Comparison{
						LHS:      parser.Variable{Name: "@var1"},
						RHS:      parser.NewIntegerValueFromString("2"),
						Operator: "=",
					},
					Statements: []parser.Statement{
						parser.Exit{},
					},
				},
				parser.Print{Value: parser.Variable{Name: "@var1"}},
				parser.TransactionControl{Token: parser.COMMIT},
			},
		},
		ResultFlow: Exit,
		Result:     "\"1\"\n",
	},
	{
		Name: "While In Cursor Fetch Error",
		Stmt: parser.WhileInCursor{
			Variables: []parser.Variable{
				{Name: "@var1"},
				{Name: "@var3"},
			},
			Cursor: parser.Identifier{Literal: "cur"},
			Statements: []parser.Statement{
				parser.Print{Value: parser.Variable{Name: "@var1"}},
				parser.TransactionControl{Token: parser.COMMIT},
			},
		},
		Error: "[L:- C:-] variable @var3 is undeclared",
	},
	{
		Name: "While In Cursor Statement Execution Error",
		Stmt: parser.WhileInCursor{
			Variables: []parser.Variable{
				{Name: "@var1"},
				{Name: "@var2"},
			},
			Cursor: parser.Identifier{Literal: "cur"},
			Statements: []parser.Statement{
				parser.If{
					Condition: parser.Comparison{
						LHS:      parser.Variable{Name: "@var1"},
						RHS:      parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
						Operator: "=",
					},
					Statements: []parser.Statement{
						parser.FlowControl{Token: parser.BREAK},
					},
				},
				parser.Print{Value: parser.Variable{Name: "@var1"}},
				parser.TransactionControl{Token: parser.COMMIT},
			},
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
}

func TestProcedure_WhileInCursor(t *testing.T) {
	tf := cmd.GetFlags()
	tf.Repository = TestDir

	proc := NewProcedure()

	for _, v := range procedureWhileInCursorTests {
		proc.Filter.Variables[0] = VariableMap{
			"@var1": value.NewNull(),
			"@var2": value.NewNull(),
		}
		proc.Filter.Cursors[0] = CursorMap{
			"CUR": &Cursor{
				query: selectQueryForCursorTest,
			},
		}
		ViewCache.Clean()
		proc.Filter.Cursors.Open(parser.Identifier{Literal: "cur"}, proc.Filter)

		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		flow, err := proc.WhileInCursor(v.Stmt)

		w.Close()
		os.Stdout = oldStdout

		log, _ := ioutil.ReadAll(r)

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
		if string(log) != v.Result {
			t.Errorf("%s: result = %q, want %q", v.Name, string(log), v.Result)
		}
	}
}
