package query

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/ternary"
)

var procedureExecuteStatementTests = []struct {
	Input  parser.Statement
	Result []Result
	Logs   []string
	Error  string
}{
	{
		Input: parser.SetFlag{
			Name:  "@@invalid",
			Value: parser.NewString("\t"),
		},
		Error: "[L:- C:-] SET: flag name @@invalid is invalid",
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
		Input: parser.VariableSubstitution{
			Variable: parser.Variable{Name: "@var1"},
			Value:    parser.NewInteger(1),
		},
		Result: []Result{},
	},
	{
		Input: parser.FunctionDeclaration{
			Name: parser.Identifier{Literal: "userfunc"},
			Parameters: []parser.Variable{
				{Name: "@arg1"},
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
			Args: []parser.Expression{
				parser.NewInteger(1),
			},
		},
		Logs: []string{
			"1",
		},
	},
	{
		Input: parser.SelectQuery{
			SelectEntity: parser.SelectEntity{
				SelectClause: parser.SelectClause{
					Fields: []parser.Expression{
						parser.Field{
							Object: parser.Variable{Name: "@var1"},
							Alias:  parser.Identifier{Literal: "var1"},
						},
					},
				},
			},
		},
		Result: []Result{
			{
				Type: SELECT,
				View: &View{
					Header: []HeaderField{
						{
							Column:    "@var1",
							Alias:     "var1",
							FromTable: true,
						},
					},
					Records: []Record{
						{
							NewCell(parser.NewInteger(1)),
						},
					},
				},
			},
		},
	},
	{
		Input: parser.VariableDeclaration{
			Assignments: []parser.Expression{
				parser.VariableAssignment{
					Variable: parser.Variable{Name: "@var1"},
				},
			},
		},
		Error: "[L:- C:-] variable @var1 is redeclared",
	},
	{
		Input: parser.VariableSubstitution{
			Variable: parser.Variable{Name: "@var2"},
			Value:    parser.NewInteger(1),
		},
		Error: "[L:- C:-] variable @var2 is undefined",
	},
	{
		Input: parser.InsertQuery{
			Table: parser.Identifier{Literal: "table1"},
			Fields: []parser.Expression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
			},
			ValuesList: []parser.Expression{
				parser.RowValue{
					Value: parser.ValueList{
						Values: []parser.Expression{
							parser.NewInteger(4),
							parser.NewString("str4"),
						},
					},
				},
				parser.RowValue{
					Value: parser.ValueList{
						Values: []parser.Expression{
							parser.NewInteger(5),
							parser.NewString("str5"),
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
		Logs: []string{
			fmt.Sprintf("2 records inserted on %q.", GetTestFilePath("table1.csv")),
		},
	},
	{
		Input: parser.UpdateQuery{
			Tables: []parser.Expression{
				parser.Table{Object: parser.Identifier{Literal: "table1"}},
			},
			SetList: []parser.Expression{
				parser.UpdateSet{
					Field: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
					Value: parser.NewString("update"),
				},
			},
			WhereClause: parser.WhereClause{
				Filter: parser.Comparison{
					LHS:      parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
					RHS:      parser.NewInteger(2),
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
		Logs: []string{
			fmt.Sprintf("1 record updated on %q.", GetTestFilePath("table1.csv")),
		},
	},
	{
		Input: parser.DeleteQuery{
			FromClause: parser.FromClause{
				Tables: []parser.Expression{
					parser.Table{
						Object: parser.Identifier{Literal: "table1"},
					},
				},
			},
			WhereClause: parser.WhereClause{
				Filter: parser.Comparison{
					LHS:      parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
					RHS:      parser.NewInteger(2),
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
		Logs: []string{
			fmt.Sprintf("1 record deleted on %q.", GetTestFilePath("table1.csv")),
		},
	},
	{
		Input: parser.CreateTable{
			Table: parser.Identifier{Literal: "newtable.csv"},
			Fields: []parser.Expression{
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
		Logs: []string{
			fmt.Sprintf("file %q is created.", GetTestFilePath("newtable.csv")),
		},
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
		Logs: []string{
			fmt.Sprintf("1 field added on %q.", GetTestFilePath("table1.csv")),
		},
	},
	{
		Input: parser.DropColumns{
			Table: parser.Identifier{Literal: "table1"},
			Columns: []parser.Expression{
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
		Logs: []string{
			fmt.Sprintf("1 field dropped on %q.", GetTestFilePath("table1.csv")),
		},
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
		Logs: []string{
			fmt.Sprintf("1 field renamed on %q.", GetTestFilePath("table1.csv")),
		},
	},
	{
		Input: parser.Print{
			Value: parser.NewInteger(12345),
		},
		Logs: []string{
			"12345",
		},
	},
	{
		Input: parser.Printf{
			Values: []parser.Expression{
				parser.NewString("value: %s"),
				parser.NewInteger(12345),
			},
		},
		Logs: []string{
			"value: 12345",
		},
	},
	{
		Input: parser.Source{
			FilePath: GetTestFilePath("source.sql"),
		},
		Logs: []string{
			"'external executable file'",
		},
	},
}

func TestProcedure_ExecuteStatement(t *testing.T) {
	tf := cmd.GetFlags()
	tf.Repository = TestDir

	UserFunctions = UserDefinedFunctionMap{}
	proc := NewProcedure()

	for _, v := range procedureExecuteStatementTests {
		ViewCache.Clear()
		Results = []Result{}
		Logs = []string{}

		_, err := proc.ExecuteStatement(v.Input)
		if err != nil {
			if len(v.Error) < 1 {
				t.Errorf("unexpected error %q for %q", err, v.Input)
			} else if err.Error() != v.Error {
				t.Errorf("error %q, want error %q for %q", err, v.Error, v.Input)
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
		if v.Logs != nil {
			if !reflect.DeepEqual(Logs, v.Logs) {
				t.Errorf("logs = %s, want %s for %q", Logs, v.Logs, v.Input)
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
			Condition: parser.NewTernary(ternary.TRUE),
			Statements: []parser.Statement{
				parser.Print{Value: parser.NewString("1")},
				parser.TransactionControl{Token: parser.COMMIT},
			},
		},
		ResultFlow: TERMINATE,
		Result:     "'1'\n",
	},
	{
		Name: "If Statement Execute Nothing",
		Stmt: parser.If{
			Condition: parser.NewTernary(ternary.FALSE),
			Statements: []parser.Statement{
				parser.Print{Value: parser.NewString("1")},
				parser.TransactionControl{Token: parser.COMMIT},
			},
		},
		ResultFlow: TERMINATE,
		Result:     "",
	},
	{
		Name: "If Statement Execute ElseIf",
		Stmt: parser.If{
			Condition: parser.NewTernary(ternary.FALSE),
			Statements: []parser.Statement{
				parser.Print{Value: parser.NewString("1")},
				parser.TransactionControl{Token: parser.COMMIT},
			},
			ElseIf: []parser.ProcExpr{
				parser.ElseIf{
					Condition: parser.NewTernary(ternary.TRUE),
					Statements: []parser.Statement{
						parser.Print{Value: parser.NewString("2")},
						parser.TransactionControl{Token: parser.COMMIT},
					},
				},
				parser.ElseIf{
					Condition: parser.NewTernary(ternary.FALSE),
					Statements: []parser.Statement{
						parser.Print{Value: parser.NewString("3")},
						parser.TransactionControl{Token: parser.COMMIT},
					},
				},
			},
			Else: parser.Else{
				Statements: []parser.Statement{
					parser.Print{Value: parser.NewString("4")},
					parser.TransactionControl{Token: parser.COMMIT},
				},
			},
		},
		ResultFlow: TERMINATE,
		Result:     "'2'\n",
	},
	{
		Name: "If Statement Execute Else",
		Stmt: parser.If{
			Condition: parser.NewTernary(ternary.FALSE),
			Statements: []parser.Statement{
				parser.Print{Value: parser.NewString("1")},
				parser.TransactionControl{Token: parser.COMMIT},
			},
			ElseIf: []parser.ProcExpr{
				parser.ElseIf{
					Condition: parser.NewTernary(ternary.FALSE),
					Statements: []parser.Statement{
						parser.Print{Value: parser.NewString("2")},
						parser.TransactionControl{Token: parser.COMMIT},
					},
				},
				parser.ElseIf{
					Condition: parser.NewTernary(ternary.FALSE),
					Statements: []parser.Statement{
						parser.Print{Value: parser.NewString("3")},
						parser.TransactionControl{Token: parser.COMMIT},
					},
				},
			},
			Else: parser.Else{
				Statements: []parser.Statement{
					parser.Print{Value: parser.NewString("4")},
					parser.TransactionControl{Token: parser.COMMIT},
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
				parser.Print{Value: parser.NewString("1")},
				parser.TransactionControl{Token: parser.COMMIT},
			},
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
}

func TestProcedure_IfStmt(t *testing.T) {
	proc := NewProcedure()

	for _, v := range procedureIfStmtTests {
		proc.Rollback()
		Logs = []string{}

		flow, err := proc.IfStmt(v.Stmt)
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
		if ReadLog() != v.Result {
			t.Errorf("%s: result = %q, want %q", v.Name, ReadLog(), v.Result)
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
				RHS:      parser.NewInteger(3),
				Operator: "<",
			},
			Statements: []parser.Statement{
				parser.VariableSubstitution{
					Variable: parser.Variable{Name: "@while_test"},
					Value: parser.Arithmetic{
						LHS:      parser.Variable{Name: "@while_test"},
						RHS:      parser.NewInteger(1),
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
				RHS:      parser.NewInteger(3),
				Operator: "<",
			},
			Statements: []parser.Statement{
				parser.VariableSubstitution{
					Variable: parser.Variable{Name: "@while_test_count"},
					Value: parser.Arithmetic{
						LHS:      parser.Variable{Name: "@while_test_count"},
						RHS:      parser.NewInteger(1),
						Operator: '+',
					},
				},
				parser.VariableSubstitution{
					Variable: parser.Variable{Name: "@while_test"},
					Value: parser.Arithmetic{
						LHS:      parser.Variable{Name: "@while_test"},
						RHS:      parser.NewInteger(1),
						Operator: '+',
					},
				},
				parser.If{
					Condition: parser.Comparison{
						LHS:      parser.Variable{Name: "@while_test_count"},
						RHS:      parser.NewInteger(2),
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
				RHS:      parser.NewInteger(3),
				Operator: "<",
			},
			Statements: []parser.Statement{
				parser.VariableSubstitution{
					Variable: parser.Variable{Name: "@while_test_count"},
					Value: parser.Arithmetic{
						LHS:      parser.Variable{Name: "@while_test_count"},
						RHS:      parser.NewInteger(1),
						Operator: '+',
					},
				},
				parser.VariableSubstitution{
					Variable: parser.Variable{Name: "@while_test"},
					Value: parser.Arithmetic{
						LHS:      parser.Variable{Name: "@while_test"},
						RHS:      parser.NewInteger(1),
						Operator: '+',
					},
				},
				parser.If{
					Condition: parser.Comparison{
						LHS:      parser.Variable{Name: "@while_test_count"},
						RHS:      parser.NewInteger(2),
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
				RHS:      parser.NewInteger(3),
				Operator: "<",
			},
			Statements: []parser.Statement{
				parser.VariableSubstitution{
					Variable: parser.Variable{Name: "@while_test_count"},
					Value: parser.Arithmetic{
						LHS:      parser.Variable{Name: "@while_test_count"},
						RHS:      parser.NewInteger(1),
						Operator: '+',
					},
				},
				parser.VariableSubstitution{
					Variable: parser.Variable{Name: "@while_test"},
					Value: parser.Arithmetic{
						LHS:      parser.Variable{Name: "@while_test"},
						RHS:      parser.NewInteger(1),
						Operator: '+',
					},
				},
				parser.If{
					Condition: parser.Comparison{
						LHS:      parser.Variable{Name: "@while_test_count"},
						RHS:      parser.NewInteger(2),
						Operator: "=",
					},
					Statements: []parser.Statement{
						parser.FlowControl{Token: parser.EXIT},
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
						RHS:      parser.NewInteger(1),
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
				RHS:      parser.NewInteger(3),
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
		proc.Rollback()
		Logs = []string{}

		if _, err := proc.Filter.VariablesList[0].Get(parser.Variable{Name: "@while_test"}); err != nil {
			proc.Filter.VariablesList[0].Add(parser.Variable{Name: "@while_test"}, parser.NewInteger(0))
		}
		proc.Filter.VariablesList[0].Set(parser.Variable{Name: "@while_test"}, parser.NewInteger(0))

		if _, err := proc.Filter.VariablesList[0].Get(parser.Variable{Name: "@while_test_count"}); err != nil {
			proc.Filter.VariablesList[0].Add(parser.Variable{Name: "@while_test_count"}, parser.NewInteger(0))
		}
		proc.Filter.VariablesList[0].Set(parser.Variable{Name: "@while_test_count"}, parser.NewInteger(0))

		flow, err := proc.While(v.Stmt)
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
		if ReadLog() != v.Result {
			t.Errorf("%s: result = %q, want %q", v.Name, ReadLog(), v.Result)
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
						RHS:      parser.NewInteger(2),
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
						RHS:      parser.NewInteger(2),
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
		Name: "While In Cursor Exit",
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
						RHS:      parser.NewInteger(2),
						Operator: "=",
					},
					Statements: []parser.Statement{
						parser.FlowControl{Token: parser.EXIT},
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
		Logs = []string{}

		proc.Filter.VariablesList[0] = Variables{
			"@var1": parser.NewNull(),
			"@var2": parser.NewNull(),
		}
		proc.Filter.CursorsList[0] = CursorMap{
			"CUR": &Cursor{
				query: selectQueryForCursorTest,
			},
		}
		ViewCache.Clear()
		proc.Filter.CursorsList.Open(parser.Identifier{Literal: "cur"}, proc.Filter)

		flow, err := proc.WhileInCursor(v.Stmt)
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
		if ReadLog() != v.Result {
			t.Errorf("%s: result = %q, want %q", v.Name, ReadLog(), v.Result)
		}
	}
}
