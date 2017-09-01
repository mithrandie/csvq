package query

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/ternary"
	"github.com/mithrandie/csvq/lib/value"
)

var procedureExecuteStatementTests = []struct {
	Input      parser.Statement
	Result     []Result
	Logs       string
	SelectLogs []string
	Error      string
	ErrorCode  int
}{
	{
		Input: parser.SetFlag{
			Name:  "@@invalid",
			Value: value.NewString("\t"),
		},
		Error:     "[L:- C:-] SET: flag name @@invalid is invalid",
		ErrorCode: 1,
	},
	{
		Input: parser.VariableDeclaration{
			Assignments: []parser.Expression{
				parser.VariableAssignment{
					Variable: parser.Variable{Name: "@var1"},
				},
			},
		},
		Result: []Result{},
	},
	{
		Input: parser.VariableDeclaration{
			Assignments: []parser.Expression{
				parser.VariableAssignment{
					Variable: parser.Variable{Name: "@var2"},
				},
			},
		},
		Result: []Result{},
	},
	{
		Input: parser.VariableDeclaration{
			Assignments: []parser.Expression{
				parser.VariableAssignment{
					Variable: parser.Variable{Name: "@var3"},
				},
			},
		},
		Result: []Result{},
	},
	{
		Input: parser.VariableDeclaration{
			Assignments: []parser.Expression{
				parser.VariableAssignment{
					Variable: parser.Variable{Name: "@var4"},
				},
			},
		},
		Result: []Result{},
	},
	{
		Input: parser.VariableSubstitution{
			Variable: parser.Variable{Name: "@var1"},
			Value:    parser.NewIntegerValueFromString("1"),
		},
		Result: []Result{},
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
			Assignments: []parser.Expression{
				parser.VariableAssignment{
					Variable: parser.Variable{Name: "@var4"},
				},
			},
		},
	},
	{
		Input: parser.FunctionDeclaration{
			Name: parser.Identifier{Literal: "userfunc"},
			Parameters: []parser.Expression{
				parser.VariableAssignment{
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
		Logs: "'1'\n",
	},
	{
		Input: parser.Print{
			Value: parser.Variable{Name: "@var3"},
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
		Input: parser.TableDeclaration{
			Table: parser.Identifier{Literal: "tbl"},
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
		Input: parser.DisposeTable{
			Table: parser.Identifier{Literal: "tbl"},
		},
	},
	{
		Input: parser.AggregateDeclaration{
			Name:   parser.Identifier{Literal: "useraggfunc"},
			Cursor: parser.Identifier{Literal: "list"},
			Statements: []parser.Statement{
				parser.VariableDeclaration{
					Assignments: []parser.Expression{
						parser.VariableAssignment{
							Variable: parser.Variable{Name: "@value"},
						},
						parser.VariableAssignment{
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
			Assignments: []parser.Expression{
				parser.VariableAssignment{
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
		Error:     "[L:- C:-] variable @var9 is undefined",
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
		Result: []Result{
			{
				Type: INSERT,
				FileInfo: &FileInfo{
					Path:      GetTestFilePath("table1.csv"),
					Delimiter: ',',
					NoHeader:  false,
					Encoding:  cmd.UTF8,
					LineBreak: cmd.LF,
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
			SetList: []parser.Expression{
				parser.UpdateSet{
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
		Result: []Result{
			{
				Type: UPDATE,
				FileInfo: &FileInfo{
					Path:      GetTestFilePath("table1.csv"),
					Delimiter: ',',
					NoHeader:  false,
					Encoding:  cmd.UTF8,
					LineBreak: cmd.LF,
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
		Result: []Result{
			{
				Type: DELETE,
				FileInfo: &FileInfo{
					Path:      GetTestFilePath("table1.csv"),
					Delimiter: ',',
					NoHeader:  false,
					Encoding:  cmd.UTF8,
					LineBreak: cmd.LF,
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
		Result: []Result{
			{
				Type: CREATE_TABLE,
				FileInfo: &FileInfo{
					Path:      GetTestFilePath("newtable.csv"),
					Delimiter: ',',
					NoHeader:  false,
					Encoding:  cmd.UTF8,
					LineBreak: cmd.LF,
				},
			},
		},
		Logs: fmt.Sprintf("file %q is created.\n", GetTestFilePath("newtable.csv")),
	},
	{
		Input: parser.AddColumns{
			Table: parser.Identifier{Literal: "table1.csv"},
			Columns: []parser.Expression{
				parser.ColumnDefault{
					Column: parser.Identifier{Literal: "column3"},
				},
			},
		},
		Result: []Result{
			{
				Type: ADD_COLUMNS,
				FileInfo: &FileInfo{
					Path:      GetTestFilePath("table1.csv"),
					Delimiter: ',',
					NoHeader:  false,
					Encoding:  cmd.UTF8,
					LineBreak: cmd.LF,
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
		Result: []Result{
			{
				Type: DROP_COLUMNS,
				FileInfo: &FileInfo{
					Path:      GetTestFilePath("table1.csv"),
					Delimiter: ',',
					NoHeader:  false,
					Encoding:  cmd.UTF8,
					LineBreak: cmd.LF,
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
		Result: []Result{
			{
				Type: RENAME_COLUMN,
				FileInfo: &FileInfo{
					Path:      GetTestFilePath("table1.csv"),
					Delimiter: ',',
					NoHeader:  false,
					Encoding:  cmd.UTF8,
					LineBreak: cmd.LF,
				},
				OperatedCount: 1,
			},
		},
		Logs: fmt.Sprintf("1 field renamed on %q.\n", GetTestFilePath("table1.csv")),
	},
	{
		Input: parser.Case{
			When: []parser.Expression{
				parser.CaseWhen{
					Condition: parser.NewTernaryValue(ternary.FALSE),
					Statements: []parser.Statement{
						parser.Print{Value: parser.NewStringValue("1")},
					},
				},
				parser.CaseWhen{
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
			Format: "value: %s",
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
		Input: parser.Trigger{
			Token:   parser.ERROR,
			Message: parser.NewStringValue("user error"),
			Code:    value.NewInteger(200),
		},
		Error:     "[L:- C:-] user error",
		ErrorCode: 200,
	},
	{
		Input: parser.Trigger{
			Token:   parser.ERROR,
			Message: parser.NewIntegerValue(200),
		},
		Error:     "[L:- C:-] ",
		ErrorCode: 200,
	},
}

func TestProcedure_ExecuteStatement(t *testing.T) {
	initFlag()
	tf := cmd.GetFlags()
	tf.Repository = TestDir
	tf.Format = cmd.CSV

	proc := NewProcedure()
	proc.Filter.VariablesList[0].Add(parser.Variable{Name: "@while_test"}, value.NewInteger(0))

	for _, v := range procedureExecuteStatementTests {
		ViewCache.Clear()
		Results = []Result{}
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
			} else if ex, ok := err.(*Exit); ok {
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
			if !reflect.DeepEqual(Results, v.Result) {
				t.Errorf("results = %q, want %q for %q", Results, v.Result, v.Input)
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
		ResultFlow: TERMINATE,
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
		ResultFlow: TERMINATE,
		Result:     "",
	},
	{
		Name: "If Statement Execute ElseIf",
		Stmt: parser.If{
			Condition: parser.NewTernaryValue(ternary.FALSE),
			Statements: []parser.Statement{
				parser.Print{Value: parser.NewStringValue("1")},
			},
			ElseIf: []parser.Expression{
				parser.ElseIf{
					Condition: parser.NewTernaryValue(ternary.TRUE),
					Statements: []parser.Statement{
						parser.Print{Value: parser.NewStringValue("2")},
					},
				},
				parser.ElseIf{
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
		ResultFlow: TERMINATE,
		Result:     "'2'\n",
	},
	{
		Name: "If Statement Execute Else",
		Stmt: parser.If{
			Condition: parser.NewTernaryValue(ternary.FALSE),
			Statements: []parser.Statement{
				parser.Print{Value: parser.NewStringValue("1")},
			},
			ElseIf: []parser.Expression{
				parser.ElseIf{
					Condition: parser.NewTernaryValue(ternary.FALSE),
					Statements: []parser.Statement{
						parser.Print{Value: parser.NewStringValue("2")},
					},
				},
				parser.ElseIf{
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
		ResultFlow: TERMINATE,
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
		Error: "[L:- C:-] field notexist does not exist",
	},
}

func TestProcedure_IfStmt(t *testing.T) {
	proc := NewProcedure()

	for _, v := range procedureIfStmtTests {
		oldStdout := os.Stdout

		r, w, _ := os.Pipe()
		os.Stdout = w
		proc.Rollback()
		w.Close()

		r, w, _ = os.Pipe()
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
			When: []parser.Expression{
				parser.CaseWhen{
					Condition: parser.NewTernaryValue(ternary.FALSE),
					Statements: []parser.Statement{
						parser.Print{Value: parser.NewStringValue("1")},
					},
				},
				parser.CaseWhen{
					Condition: parser.NewTernaryValue(ternary.TRUE),
					Statements: []parser.Statement{
						parser.Print{Value: parser.NewStringValue("2")},
					},
				},
			},
		},
		ResultFlow: TERMINATE,
		Result:     "'2'\n",
	},
	{
		Name: "Case Comparison",
		Stmt: parser.Case{
			Value: parser.NewIntegerValue(2),
			When: []parser.Expression{
				parser.CaseWhen{
					Condition: parser.NewIntegerValue(1),
					Statements: []parser.Statement{
						parser.Print{Value: parser.NewStringValue("1")},
					},
				},
				parser.CaseWhen{
					Condition: parser.NewIntegerValue(2),
					Statements: []parser.Statement{
						parser.Print{Value: parser.NewStringValue("2")},
					},
				},
			},
		},
		ResultFlow: TERMINATE,
		Result:     "'2'\n",
	},
	{
		Name: "Case Else",
		Stmt: parser.Case{
			When: []parser.Expression{
				parser.CaseWhen{
					Condition: parser.NewTernaryValue(ternary.FALSE),
					Statements: []parser.Statement{
						parser.Print{Value: parser.NewStringValue("1")},
					},
				},
				parser.CaseWhen{
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
		ResultFlow: TERMINATE,
		Result:     "'3'\n",
	},
	{
		Name: "Case No Match",
		Stmt: parser.Case{
			When: []parser.Expression{
				parser.CaseWhen{
					Condition: parser.NewTernaryValue(ternary.FALSE),
					Statements: []parser.Statement{
						parser.Print{Value: parser.NewStringValue("1")},
					},
				},
				parser.CaseWhen{
					Condition: parser.NewTernaryValue(ternary.FALSE),
					Statements: []parser.Statement{
						parser.Print{Value: parser.NewStringValue("2")},
					},
				},
			},
		},
		ResultFlow: TERMINATE,
		Result:     "",
	},
	{
		Name: "Case Comparison Value Error",
		Stmt: parser.Case{
			Value: parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			When: []parser.Expression{
				parser.CaseWhen{
					Condition: parser.NewIntegerValue(1),
					Statements: []parser.Statement{
						parser.Print{Value: parser.NewStringValue("1")},
					},
				},
				parser.CaseWhen{
					Condition: parser.NewIntegerValue(2),
					Statements: []parser.Statement{
						parser.Print{Value: parser.NewStringValue("2")},
					},
				},
			},
		},
		ResultFlow: ERROR,
		Error:      "[L:- C:-] field notexist does not exist",
	},
	{
		Name: "Case Condition Error",
		Stmt: parser.Case{
			When: []parser.Expression{
				parser.CaseWhen{
					Condition: parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
					Statements: []parser.Statement{
						parser.Print{Value: parser.NewStringValue("1")},
					},
				},
			},
		},
		ResultFlow: ERROR,
		Error:      "[L:- C:-] field notexist does not exist",
	},
}

func TestProcedure_Case(t *testing.T) {
	proc := NewProcedure()

	for _, v := range procedureCaseStmtTests {
		oldStdout := os.Stdout

		r, w, _ := os.Pipe()
		os.Stdout = w
		proc.Rollback()
		w.Close()

		r, w, _ = os.Pipe()
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
		ResultFlow: TERMINATE,
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
		ResultFlow: TERMINATE,
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
		ResultFlow: TERMINATE,
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
		ResultFlow: EXIT,
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
	proc := NewProcedure()

	for _, v := range procedureWhileTests {
		if _, err := proc.Filter.VariablesList[0].Get(parser.Variable{Name: "@while_test"}); err != nil {
			proc.Filter.VariablesList[0].Add(parser.Variable{Name: "@while_test"}, value.NewInteger(0))
		}
		proc.Filter.VariablesList[0].Set(parser.Variable{Name: "@while_test"}, value.NewInteger(0))

		if _, err := proc.Filter.VariablesList[0].Get(parser.Variable{Name: "@while_test_count"}); err != nil {
			proc.Filter.VariablesList[0].Add(parser.Variable{Name: "@while_test_count"}, value.NewInteger(0))
		}
		proc.Filter.VariablesList[0].Set(parser.Variable{Name: "@while_test_count"}, value.NewInteger(0))

		oldStdout := os.Stdout

		r, w, _ := os.Pipe()
		os.Stdout = w
		proc.Rollback()
		w.Close()

		r, w, _ = os.Pipe()
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
		ResultFlow: TERMINATE,
		Result:     "'1'\n'2'\n'3'\n",
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
		ResultFlow: TERMINATE,
		Result:     "'1'\n'3'\n",
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
		ResultFlow: TERMINATE,
		Result:     "'1'\n",
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
		ResultFlow: EXIT,
		Result:     "'1'\n",
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
		Error: "[L:- C:-] variable @var3 is undefined",
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
		proc.Filter.VariablesList[0] = Variables{
			"@var1": value.NewNull(),
			"@var2": value.NewNull(),
		}
		proc.Filter.CursorsList[0] = CursorMap{
			"CUR": &Cursor{
				query: selectQueryForCursorTest,
			},
		}
		ViewCache.Clear()
		proc.Filter.CursorsList.Open(parser.Identifier{Literal: "cur"}, proc.Filter)

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
