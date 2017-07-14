package query

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/ternary"
)

var executeTests = []struct {
	Name       string
	Input      string
	Output     string
	UpdateFile string
	Content    string
	Error      string
}{
	{
		Name:  "Select Query",
		Input: "select 1 from dual",
		Output: "+---+\n" +
			"| 1 |\n" +
			"+---+\n" +
			"| 1 |\n" +
			"+---+\n",
	},
	{
		Name:  "Insert Query",
		Input: "insert into insert_query values (4, 'str4'), (5, 'str5')",
		Output: fmt.Sprintf("%d records inserted on %q\n", 2, GetTestFilePath("insert_query.csv")) +
			fmt.Sprintf("Commit: file %q is updated.\n", GetTestFilePath("insert_query.csv")),
		UpdateFile: GetTestFilePath("insert_query.csv"),
		Content: "\"column1\",\"column2\"\n" +
			"\"1\",\"str1\"\n" +
			"\"2\",\"str2\"\n" +
			"\"3\",\"str3\"\n" +
			"4,\"str4\"\n" +
			"5,\"str5\"",
	},
	{
		Name:  "Update Query",
		Input: "update update_query set column2 = 'update' where column1 = 2",
		Output: fmt.Sprintf("%d record updated on %q\n", 1, GetTestFilePath("update_query.csv")) +
			fmt.Sprintf("Commit: file %q is updated.\n", GetTestFilePath("update_query.csv")),
		UpdateFile: GetTestFilePath("update_query.csv"),
		Content: "\"column1\",\"column2\"\n" +
			"\"1\",\"str1\"\n" +
			"\"2\",\"update\"\n" +
			"\"3\",\"str3\"",
	},
	{
		Name:   "Update Query No Record Updated",
		Input:  "update update_query set column2 = 'update' where false",
		Output: fmt.Sprintf("no record updated on %q\n", GetTestFilePath("update_query.csv")),
	},
	{
		Name:  "Delete Query",
		Input: "delete from delete_query where column1 = 2",
		Output: fmt.Sprintf("%d record deleted on %q\n", 1, GetTestFilePath("delete_query.csv")) +
			fmt.Sprintf("Commit: file %q is updated.\n", GetTestFilePath("delete_query.csv")),
		UpdateFile: GetTestFilePath("delete_query.csv"),
		Content: "\"column1\",\"column2\"\n" +
			"\"1\",\"str1\"\n" +
			"\"3\",\"str3\"",
	},
	{
		Name:   "Delete Query No Record Deleted",
		Input:  "delete from delete_query where false",
		Output: fmt.Sprintf("no record deleted on %q\n", GetTestFilePath("delete_query.csv")),
	},
	{
		Name:  "Create Table",
		Input: "create table `create_table.csv` (column1, column2)",
		Output: fmt.Sprintf("file %q is created\n", GetTestFilePath("create_table.csv")) +
			fmt.Sprintf("Commit: file %q is created.\n", GetTestFilePath("create_table.csv")),
		UpdateFile: GetTestFilePath("create_table.csv"),
		Content:    "\"column1\",\"column2\"\n",
	},
	{
		Name:  "Add Columns",
		Input: "alter table add_columns add column3",
		Output: fmt.Sprintf("%d field added on %q\n", 1, GetTestFilePath("add_columns.csv")) +
			fmt.Sprintf("Commit: file %q is updated.\n", GetTestFilePath("add_columns.csv")),
		UpdateFile: GetTestFilePath("add_columns.csv"),
		Content: "\"column1\",\"column2\",\"column3\"\n" +
			"\"1\",\"str1\",\n" +
			"\"2\",\"str2\",\n" +
			"\"3\",\"str3\",",
	},
	{
		Name:  "Drop Columns",
		Input: "alter table drop_columns drop column1",
		Output: fmt.Sprintf("%d field dropped on %q\n", 1, GetTestFilePath("drop_columns.csv")) +
			fmt.Sprintf("Commit: file %q is updated.\n", GetTestFilePath("drop_columns.csv")),
		UpdateFile: GetTestFilePath("drop_columns.csv"),
		Content: "\"column2\"\n" +
			"\"str1\"\n" +
			"\"str2\"\n" +
			"\"str3\"",
	},
	{
		Name:  "Rename Column",
		Input: "alter table rename_column rename column1 to newcolumn",
		Output: fmt.Sprintf("%d field renamed on %q\n", 1, GetTestFilePath("rename_column.csv")) +
			fmt.Sprintf("Commit: file %q is updated.\n", GetTestFilePath("rename_column.csv")),
		UpdateFile: GetTestFilePath("rename_column.csv"),
		Content: "\"newcolumn\",\"column2\"\n" +
			"\"1\",\"str1\"\n" +
			"\"2\",\"str2\"\n" +
			"\"3\",\"str3\"",
	},
	{
		Name:   "Print",
		Input:  "var @a := 1; print @a;",
		Output: "1\n",
	},
	{
		Name:  "Query Execution Error",
		Input: "select from",
		Error: "syntax error: unexpected from [L:1 C:8]",
	},
}

func TestExecute(t *testing.T) {
	tf := cmd.GetFlags()
	tf.Format = cmd.TEXT
	tf.Repository = TestDir

	for _, v := range executeTests {
		out, err := Execute(v.Input)

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

		if out != v.Output {
			t.Errorf("%s: output = %q, want %q", v.Name, out, v.Output)
		}

		if 0 < len(v.UpdateFile) {
			fp, _ := os.Open(v.UpdateFile)
			buf, _ := ioutil.ReadAll(fp)
			if string(buf) != v.Content {
				t.Errorf("%s: content = %q, want %q", v.Name, string(buf), v.Content)
			}
		}
	}
}

var executeStatementTests = []struct {
	Input  parser.Statement
	Result []Result
	Error  string
}{
	{
		Input: parser.SetFlag{
			Name:  "@@invalid",
			Value: parser.NewString("\t"),
		},
		Error: "invalid flag name: @@invalid",
	},
	{
		Input: parser.VariableDeclaration{
			Assignments: []parser.Expression{
				parser.VariableAssignment{
					Name: "@var1",
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
					Name: "@var1",
				},
			},
		},
		Error: "variable @var1 is redeclared",
	},
	{
		Input: parser.VariableSubstitution{
			Variable: parser.Variable{Name: "@var2"},
			Value:    parser.NewInteger(1),
		},
		Error: "variable @var2 is undefined",
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
				Log:           fmt.Sprintf("2 records inserted on %q", GetTestFilePath("table1.csv")),
			},
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
				Log:           fmt.Sprintf("1 record updated on %q", GetTestFilePath("table1.csv")),
			},
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
				Log:           fmt.Sprintf("1 record deleted on %q", GetTestFilePath("table1.csv")),
			},
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
				Log: fmt.Sprintf("file %q is created", GetTestFilePath("newtable.csv")),
			},
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
				Log:           fmt.Sprintf("1 field added on %q", GetTestFilePath("table1.csv")),
			},
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
				Log:           fmt.Sprintf("1 field dropped on %q", GetTestFilePath("table1.csv")),
			},
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
				Log:           fmt.Sprintf("1 field renamed on %q", GetTestFilePath("table1.csv")),
			},
		},
	},
	{
		Input: parser.Print{
			Value: parser.NewInteger(12345),
		},
		Result: []Result{
			{
				Type: PRINT,
				Log:  "12345",
			},
		},
	},
}

func TestExecuteStatement(t *testing.T) {
	GlobalVars = map[string]parser.Primary{}

	tf := cmd.GetFlags()
	tf.Repository = TestDir

	for _, v := range executeStatementTests {
		ViewCache.Clear()
		ResultSet = []Result{}

		_, _, err := ExecuteStatement(v.Input)
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

		if !reflect.DeepEqual(ResultSet, v.Result) {
			t.Errorf("results = %q, want %q for %q", ResultSet, v.Result, v.Input)
		}
	}
}

var ifStmtTests = []struct {
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
		Error: "field notexist does not exist",
	},
}

func TestIfStmt(t *testing.T) {
	for _, v := range ifStmtTests {
		Rollback()

		flow, result, err := IfStmt(v.Stmt)
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
		if result != v.Result {
			t.Errorf("%s: result = %q, want %q", v.Name, result, v.Result)
		}
	}
}

var whileTests = []struct {
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
		Error: "field notexist does not exist",
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
		Error: "field notexist does not exist",
	},
}

func TestWhile(t *testing.T) {
	for _, v := range whileTests {
		Rollback()
		if _, err := GlobalVars.Get("@while_test"); err != nil {
			GlobalVars.Add("@while_test", parser.NewInteger(0))
		}
		GlobalVars.Set("@while_test", parser.NewInteger(0))

		if _, err := GlobalVars.Get("@while_test_count"); err != nil {
			GlobalVars.Add("@while_test_count", parser.NewInteger(0))
		}
		GlobalVars.Set("@while_test_count", parser.NewInteger(0))

		flow, result, err := While(v.Stmt)
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
		if result != v.Result {
			t.Errorf("%s: result = %q, want %q", v.Name, result, v.Result)
		}
	}
}

var whileInCursorTests = []struct {
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
		Error: "variable @var3 is undefined",
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
		Error: "field notexist does not exist",
	},
}

func TestWhileInCursor(t *testing.T) {
	tf := cmd.GetFlags()
	tf.Repository = TestDir

	for _, v := range whileInCursorTests {
		Cursors = CursorMap{
			"cur": &Cursor{
				name:  "cur",
				query: selectQueryForCursorTest,
			},
		}
		Cursors.Open("cur")

		GlobalVars = Variables{
			"@var1": parser.NewNull(),
			"@var2": parser.NewNull(),
		}

		flow, result, err := WhileInCursor(v.Stmt)
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
		if result != v.Result {
			t.Errorf("%s: result = %q, want %q", v.Name, result, v.Result)
		}
	}
}

var fetchCursorTests = []struct {
	Name          string
	CurName       string
	FetchPosition parser.Expression
	Variables     []parser.Variable
	Success       bool
	ResultVars    Variables
	Error         string
}{
	{
		Name:    "Fetch Cursor First Time",
		CurName: "cur",
		Variables: []parser.Variable{
			{Name: "@var1"},
			{Name: "@var2"},
		},
		Success: true,
		ResultVars: Variables{
			"@var1": parser.NewString("1"),
			"@var2": parser.NewString("str1"),
		},
	},
	{
		Name:    "Fetch Cursor Second Time",
		CurName: "cur",
		Variables: []parser.Variable{
			{Name: "@var1"},
			{Name: "@var2"},
		},
		Success: true,
		ResultVars: Variables{
			"@var1": parser.NewString("2"),
			"@var2": parser.NewString("str2"),
		},
	},
	{
		Name:    "Fetch Cursor Third Time",
		CurName: "cur",
		Variables: []parser.Variable{
			{Name: "@var1"},
			{Name: "@var2"},
		},
		Success: true,
		ResultVars: Variables{
			"@var1": parser.NewString("3"),
			"@var2": parser.NewString("str3"),
		},
	},
	{
		Name:    "Fetch Cursor Forth Time",
		CurName: "cur",
		Variables: []parser.Variable{
			{Name: "@var1"},
			{Name: "@var2"},
		},
		Success: false,
		ResultVars: Variables{
			"@var1": parser.NewString("3"),
			"@var2": parser.NewString("str3"),
		},
	},
	{
		Name:    "Fetch Cursor Absolute",
		CurName: "cur",
		FetchPosition: parser.FetchPosition{
			Position: parser.Token{Token: parser.ABSOLUTE, Literal: "absolute"},
			Number:   parser.NewInteger(1),
		},
		Variables: []parser.Variable{
			{Name: "@var1"},
			{Name: "@var2"},
		},
		Success: true,
		ResultVars: Variables{
			"@var1": parser.NewString("2"),
			"@var2": parser.NewString("str2"),
		},
	},
	{
		Name:    "Fetch Cursor Fetch Error",
		CurName: "notexist",
		Variables: []parser.Variable{
			{Name: "@var1"},
			{Name: "@var2"},
		},
		Error: "cursor notexist does not exist",
	},
	{
		Name:    "Fetch Cursor Not Match Number Error",
		CurName: "cur2",
		Variables: []parser.Variable{
			{Name: "@var1"},
		},
		Error: "cursor cur2 field length does not match variables number",
	},
	{
		Name:    "Fetch Cursor Substitution Error",
		CurName: "cur2",
		Variables: []parser.Variable{
			{Name: "@var1"},
			{Name: "@notexist"},
		},
		Error: "variable @notexist is undefined",
	},
	{
		Name:    "Fetch Cursor Number Value Error",
		CurName: "cur",
		FetchPosition: parser.FetchPosition{
			Position: parser.Token{Token: parser.ABSOLUTE, Literal: "absolute"},
			Number:   parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
		},
		Variables: []parser.Variable{
			{Name: "@var1"},
			{Name: "@var2"},
		},
		Error: "field notexist does not exist",
	},
	{
		Name:    "Fetch Cursor Number Not Integer Error",
		CurName: "cur",
		FetchPosition: parser.FetchPosition{
			Position: parser.Token{Token: parser.ABSOLUTE, Literal: "absolute"},
			Number:   parser.NewNull(),
		},
		Variables: []parser.Variable{
			{Name: "@var1"},
			{Name: "@var2"},
		},
		Error: "fetch position NULL is not a integer",
	},
}

func TestFetchCursor(t *testing.T) {
	tf := cmd.GetFlags()
	tf.Repository = TestDir

	Cursors = CursorMap{
		"cur": &Cursor{
			name:  "cur",
			query: selectQueryForCursorTest,
		},
		"cur2": &Cursor{
			name:  "cur2",
			query: selectQueryForCursorTest,
		},
	}
	Cursors.Open("cur")
	Cursors.Open("cur2")

	GlobalVars = Variables{
		"@var1": parser.NewNull(),
		"@var2": parser.NewNull(),
	}

	for _, v := range fetchCursorTests {
		success, err := FetchCursor(v.CurName, v.FetchPosition, v.Variables)
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
		if success != v.Success {
			t.Errorf("%s: success = %t, want %t", v.Name, success, v.Success)
		}
		if !reflect.DeepEqual(GlobalVars, v.ResultVars) {
			t.Errorf("%s: global vars = %q, want %q", v.Name, GlobalVars, v.ResultVars)
		}
	}
}

var selectTests = []struct {
	Name   string
	Query  parser.SelectQuery
	Result *View
	Error  string
}{
	{
		Name: "Select",
		Query: parser.SelectQuery{
			SelectEntity: parser.SelectEntity{
				SelectClause: parser.SelectClause{
					Fields: []parser.Expression{
						parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
						parser.Field{Object: parser.Function{Name: "count", Option: parser.Option{Args: []parser.Expression{parser.AllColumns{}}}}},
					},
				},
				FromClause: parser.FromClause{
					Tables: []parser.Expression{
						parser.Table{Object: parser.Identifier{Literal: "group_table"}},
					},
				},
				WhereClause: parser.WhereClause{
					Filter: parser.Comparison{
						LHS:      parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
						RHS:      parser.NewInteger(3),
						Operator: "<",
					},
				},
				GroupByClause: parser.GroupByClause{
					Items: []parser.Expression{
						parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
					},
				},
				HavingClause: parser.HavingClause{
					Filter: parser.Comparison{
						LHS:      parser.Function{Name: "count", Option: parser.Option{Args: []parser.Expression{parser.AllColumns{}}}},
						RHS:      parser.NewInteger(1),
						Operator: ">",
					},
				},
			},
			OrderByClause: parser.OrderByClause{
				Items: []parser.Expression{
					parser.OrderItem{Value: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
				},
			},
			LimitClause: parser.LimitClause{
				Value: parser.NewInteger(5),
			},
			OffsetClause: parser.OffsetClause{
				Value: parser.NewInteger(0),
			},
		},
		Result: &View{
			FileInfo: &FileInfo{
				Path:      GetTestFilePath("group_table.csv"),
				Delimiter: ',',
				NoHeader:  false,
				Encoding:  cmd.UTF8,
				LineBreak: cmd.LF,
			},
			Header: []HeaderField{
				{
					Reference: "group_table",
					Column:    "column1",
					FromTable: true,
				},
				{
					Column:    "count(*)",
					Alias:     "count(*)",
					FromTable: true,
				},
			},
			Records: []Record{
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("1"),
					parser.NewInteger(2),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("2"),
					parser.NewInteger(2),
				}),
			},
		},
	},
	{
		Name: "Union",
		Query: parser.SelectQuery{
			SelectEntity: parser.SelectSet{
				LHS: parser.SelectEntity{
					SelectClause: parser.SelectClause{
						Fields: []parser.Expression{
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
						},
					},
					FromClause: parser.FromClause{
						Tables: []parser.Expression{
							parser.Table{Object: parser.Identifier{Literal: "table1"}},
						},
					},
				},
				Operator: parser.Token{Token: parser.UNION, Literal: "union"},
				RHS: parser.SelectEntity{
					SelectClause: parser.SelectClause{
						Fields: []parser.Expression{
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column3"}}},
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column4"}}},
						},
					},
					FromClause: parser.FromClause{
						Tables: []parser.Expression{
							parser.Table{Object: parser.Identifier{Literal: "table4"}},
						},
					},
				},
			},
		},
		Result: &View{
			Header: []HeaderField{
				{
					Reference: "table1",
					Column:    "column1",
					FromTable: true,
				},
				{
					Reference: "table1",
					Column:    "column2",
					FromTable: true,
				},
			},
			Records: []Record{
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("2"),
					parser.NewString("str2"),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("3"),
					parser.NewString("str3"),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("4"),
					parser.NewString("str4"),
				}),
			},
		},
	},
	{
		Name: "Intersect",
		Query: parser.SelectQuery{
			SelectEntity: parser.SelectSet{
				LHS: parser.SelectEntity{
					SelectClause: parser.SelectClause{
						Fields: []parser.Expression{
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
						},
					},
					FromClause: parser.FromClause{
						Tables: []parser.Expression{
							parser.Table{Object: parser.Identifier{Literal: "table1"}},
						},
					},
				},
				Operator: parser.Token{Token: parser.INTERSECT, Literal: "intersect"},
				RHS: parser.SelectEntity{
					SelectClause: parser.SelectClause{
						Fields: []parser.Expression{
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column3"}}},
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column4"}}},
						},
					},
					FromClause: parser.FromClause{
						Tables: []parser.Expression{
							parser.Table{Object: parser.Identifier{Literal: "table4"}},
						},
					},
				},
			},
		},
		Result: &View{
			Header: []HeaderField{
				{
					Reference: "table1",
					Column:    "column1",
					FromTable: true,
				},
				{
					Reference: "table1",
					Column:    "column2",
					FromTable: true,
				},
			},
			Records: []Record{
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("2"),
					parser.NewString("str2"),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("3"),
					parser.NewString("str3"),
				}),
			},
		},
	},
	{
		Name: "Except",
		Query: parser.SelectQuery{
			SelectEntity: parser.SelectSet{
				LHS: parser.SelectEntity{
					SelectClause: parser.SelectClause{
						Fields: []parser.Expression{
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
						},
					},
					FromClause: parser.FromClause{
						Tables: []parser.Expression{
							parser.Table{Object: parser.Identifier{Literal: "table1"}},
						},
					},
				},
				Operator: parser.Token{Token: parser.EXCEPT, Literal: "except"},
				RHS: parser.SelectEntity{
					SelectClause: parser.SelectClause{
						Fields: []parser.Expression{
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column3"}}},
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column4"}}},
						},
					},
					FromClause: parser.FromClause{
						Tables: []parser.Expression{
							parser.Table{Object: parser.Identifier{Literal: "table4"}},
						},
					},
				},
			},
		},
		Result: &View{
			Header: []HeaderField{
				{
					Reference: "table1",
					Column:    "column1",
					FromTable: true,
				},
				{
					Reference: "table1",
					Column:    "column2",
					FromTable: true,
				},
			},
			Records: []Record{
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
			},
		},
	},
	{
		Name: "Union with SubQuery",
		Query: parser.SelectQuery{
			SelectEntity: parser.SelectSet{
				LHS: parser.Subquery{
					Query: parser.SelectQuery{
						SelectEntity: parser.SelectEntity{
							SelectClause: parser.SelectClause{
								Fields: []parser.Expression{
									parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
									parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
								},
							},
							FromClause: parser.FromClause{
								Tables: []parser.Expression{
									parser.Table{Object: parser.Identifier{Literal: "table1"}},
								},
							},
						},
					},
				},
				Operator: parser.Token{Token: parser.UNION, Literal: "union"},
				RHS: parser.SelectEntity{
					SelectClause: parser.SelectClause{
						Fields: []parser.Expression{
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column3"}}},
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column4"}}},
						},
					},
					FromClause: parser.FromClause{
						Tables: []parser.Expression{
							parser.Table{Object: parser.Identifier{Literal: "table4"}},
						},
					},
				},
			},
		},
		Result: &View{
			Header: []HeaderField{
				{
					Reference: "table1",
					Column:    "column1",
					FromTable: true,
				},
				{
					Reference: "table1",
					Column:    "column2",
					FromTable: true,
				},
			},
			Records: []Record{
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("2"),
					parser.NewString("str2"),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("3"),
					parser.NewString("str3"),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("4"),
					parser.NewString("str4"),
				}),
			},
		},
	},
	{
		Name: "Union Field Length Error",
		Query: parser.SelectQuery{
			SelectEntity: parser.SelectSet{
				LHS: parser.SelectEntity{
					SelectClause: parser.SelectClause{
						Fields: []parser.Expression{
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
						},
					},
					FromClause: parser.FromClause{
						Tables: []parser.Expression{
							parser.Table{Object: parser.Identifier{Literal: "table1"}},
						},
					},
				},
				Operator: parser.Token{Token: parser.UNION, Literal: "union"},
				RHS: parser.SelectEntity{
					SelectClause: parser.SelectClause{
						Fields: []parser.Expression{
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column3"}}},
						},
					},
					FromClause: parser.FromClause{
						Tables: []parser.Expression{
							parser.Table{Object: parser.Identifier{Literal: "table4"}},
						},
					},
				},
			},
		},
		Error: "UNION: field length does not match",
	},
	{
		Name: "Union LHS Error",
		Query: parser.SelectQuery{
			SelectEntity: parser.SelectSet{
				LHS: parser.SelectEntity{
					SelectClause: parser.SelectClause{
						Fields: []parser.Expression{
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}}},
						},
					},
					FromClause: parser.FromClause{
						Tables: []parser.Expression{
							parser.Table{Object: parser.Identifier{Literal: "table1"}},
						},
					},
				},
				Operator: parser.Token{Token: parser.UNION, Literal: "union"},
				RHS: parser.SelectEntity{
					SelectClause: parser.SelectClause{
						Fields: []parser.Expression{
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column3"}}},
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column4"}}},
						},
					},
					FromClause: parser.FromClause{
						Tables: []parser.Expression{
							parser.Table{Object: parser.Identifier{Literal: "table4"}},
						},
					},
				},
			},
		},
		Error: "field notexist does not exist",
	},
	{
		Name: "Union RHS Error",
		Query: parser.SelectQuery{
			SelectEntity: parser.SelectSet{
				LHS: parser.SelectEntity{
					SelectClause: parser.SelectClause{
						Fields: []parser.Expression{
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
						},
					},
					FromClause: parser.FromClause{
						Tables: []parser.Expression{
							parser.Table{Object: parser.Identifier{Literal: "table1"}},
						},
					},
				},
				Operator: parser.Token{Token: parser.UNION, Literal: "union"},
				RHS: parser.SelectEntity{
					SelectClause: parser.SelectClause{
						Fields: []parser.Expression{
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column3"}}},
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}}},
						},
					},
					FromClause: parser.FromClause{
						Tables: []parser.Expression{
							parser.Table{Object: parser.Identifier{Literal: "table4"}},
						},
					},
				},
			},
		},
		Error: "field notexist does not exist",
	},
	{
		Name: "Common Tables",
		Query: parser.SelectQuery{
			WithClause: parser.WithClause{
				With: "with",
				InlineTables: []parser.Expression{
					parser.InlineTable{
						Name: parser.Identifier{Literal: "it"},
						Columns: []parser.Expression{
							parser.Identifier{Literal: "c1"},
						},
						As: "as",
						Query: parser.SelectQuery{
							SelectEntity: parser.SelectEntity{
								SelectClause: parser.SelectClause{
									Select: "select",
									Fields: []parser.Expression{
										parser.Field{Object: parser.NewInteger(2)},
									},
								},
							},
						},
					},
				},
			},
			SelectEntity: parser.SelectEntity{
				SelectClause: parser.SelectClause{
					Fields: []parser.Expression{
						parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "c1"}}},
					},
				},
				FromClause: parser.FromClause{
					Tables: []parser.Expression{
						parser.Table{Object: parser.Identifier{Literal: "it"}},
					},
				},
			},
		},
		Result: &View{
			Header: []HeaderField{
				{
					Reference: "it",
					Column:    "c1",
					FromTable: true,
				},
			},
			Records: []Record{
				NewRecordWithoutId([]parser.Primary{
					parser.NewInteger(2),
				}),
			},
		},
	},
	{
		Name: "Common Tables Recursion",
		Query: parser.SelectQuery{
			WithClause: parser.WithClause{
				With: "with",
				InlineTables: []parser.Expression{
					parser.InlineTable{
						Recursive: parser.Token{Token: parser.RECURSIVE, Literal: "recursive"},
						Name:      parser.Identifier{Literal: "it"},
						Columns: []parser.Expression{
							parser.Identifier{Literal: "n"},
						},
						As: "as",
						Query: parser.SelectQuery{
							SelectEntity: parser.SelectSet{
								LHS: parser.SelectEntity{
									SelectClause: parser.SelectClause{
										Select: "select",
										Fields: []parser.Expression{
											parser.Field{Object: parser.NewInteger(1)},
										},
									},
								},
								Operator: parser.Token{Token: parser.UNION, Literal: "union"},
								RHS: parser.SelectEntity{
									SelectClause: parser.SelectClause{
										Select: "select",
										Fields: []parser.Expression{
											parser.Field{
												Object: parser.Arithmetic{
													LHS:      parser.FieldReference{Column: parser.Identifier{Literal: "n"}},
													RHS:      parser.NewInteger(1),
													Operator: '+',
												},
											},
										},
									},
									FromClause: parser.FromClause{
										Tables: []parser.Expression{
											parser.Table{Object: parser.Identifier{Literal: "it"}},
										},
									},
									WhereClause: parser.WhereClause{
										Filter: parser.Comparison{
											LHS:      parser.FieldReference{Column: parser.Identifier{Literal: "n"}},
											RHS:      parser.NewInteger(3),
											Operator: "<",
										},
									},
								},
							},
						},
					},
				},
			},
			SelectEntity: parser.SelectEntity{
				SelectClause: parser.SelectClause{
					Fields: []parser.Expression{
						parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "n"}}},
					},
				},
				FromClause: parser.FromClause{
					Tables: []parser.Expression{
						parser.Table{Object: parser.Identifier{Literal: "it"}},
					},
				},
			},
		},
		Result: &View{
			Header: []HeaderField{
				{
					Reference: "it",
					Column:    "n",
					FromTable: true,
				},
			},
			Records: []Record{
				NewRecordWithoutId([]parser.Primary{
					parser.NewInteger(1),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewInteger(2),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewInteger(3),
				}),
			},
		},
	},
	{
		Name: "Common Tables Recursion Field Length Error",
		Query: parser.SelectQuery{
			WithClause: parser.WithClause{
				With: "with",
				InlineTables: []parser.Expression{
					parser.InlineTable{
						Recursive: parser.Token{Token: parser.RECURSIVE, Literal: "recursive"},
						Name:      parser.Identifier{Literal: "it"},
						Columns: []parser.Expression{
							parser.Identifier{Literal: "n"},
						},
						As: "as",
						Query: parser.SelectQuery{
							SelectEntity: parser.SelectSet{
								LHS: parser.SelectEntity{
									SelectClause: parser.SelectClause{
										Select: "select",
										Fields: []parser.Expression{
											parser.Field{Object: parser.NewInteger(1)},
										},
									},
								},
								Operator: parser.Token{Token: parser.UNION, Literal: "union"},
								RHS: parser.SelectEntity{
									SelectClause: parser.SelectClause{
										Select: "select",
										Fields: []parser.Expression{
											parser.Field{
												Object: parser.Arithmetic{
													LHS:      parser.FieldReference{Column: parser.Identifier{Literal: "n"}},
													RHS:      parser.NewInteger(1),
													Operator: '+',
												},
											},
											parser.Field{Object: parser.NewInteger(2)},
										},
									},
									FromClause: parser.FromClause{
										Tables: []parser.Expression{
											parser.Table{Object: parser.Identifier{Literal: "it"}},
										},
									},
									WhereClause: parser.WhereClause{
										Filter: parser.Comparison{
											LHS:      parser.FieldReference{Column: parser.Identifier{Literal: "n"}},
											RHS:      parser.NewInteger(3),
											Operator: "<",
										},
									},
								},
							},
						},
					},
				},
			},
			SelectEntity: parser.SelectEntity{
				SelectClause: parser.SelectClause{
					Fields: []parser.Expression{
						parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "n"}}},
					},
				},
				FromClause: parser.FromClause{
					Tables: []parser.Expression{
						parser.Table{Object: parser.Identifier{Literal: "it"}},
					},
				},
			},
		},
		Error: "UNION: field length does not match",
	},
}

func TestSelect(t *testing.T) {
	tf := cmd.GetFlags()
	tf.Repository = TestDir

	for _, v := range selectTests {
		ViewCache.Clear()
		result, err := Select(v.Query)
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
		if !reflect.DeepEqual(result, v.Result) {
			t.Errorf("%s: result = %q, want %q", v.Name, result, v.Result)
		}
	}
}

var insertTests = []struct {
	Name   string
	Query  parser.InsertQuery
	Result *View
	Error  string
}{
	{
		Name: "Insert Query",
		Query: parser.InsertQuery{
			Insert: "insert",
			Into:   "into",
			Table:  parser.Identifier{Literal: "table1"},
			Fields: []parser.Expression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
			},
			Values: "values",
			ValuesList: []parser.Expression{
				parser.RowValue{
					Value: parser.ValueList{
						Values: []parser.Expression{
							parser.NewInteger(4),
						},
					},
				},
				parser.RowValue{
					Value: parser.ValueList{
						Values: []parser.Expression{
							parser.NewInteger(5),
						},
					},
				},
			},
		},
		Result: &View{
			FileInfo: &FileInfo{
				Path:      GetTestFilePath("table1.csv"),
				Delimiter: ',',
				NoHeader:  false,
				Encoding:  cmd.UTF8,
				LineBreak: cmd.LF,
			},
			Header: NewHeaderWithoutId("table1", []string{"column1", "column2"}),
			Records: []Record{
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("2"),
					parser.NewString("str2"),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("3"),
					parser.NewString("str3"),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewInteger(4),
					parser.NewNull(),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewInteger(5),
					parser.NewNull(),
				}),
			},
			OperatedRecords: 2,
		},
	},
	{
		Name: "Insert Query All Columns",
		Query: parser.InsertQuery{
			Insert: "insert",
			Into:   "into",
			Table:  parser.Identifier{Literal: "table1"},
			Values: "values",
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
		Result: &View{
			FileInfo: &FileInfo{
				Path:      GetTestFilePath("table1.csv"),
				Delimiter: ',',
				NoHeader:  false,
				Encoding:  cmd.UTF8,
				LineBreak: cmd.LF,
			},
			Header: NewHeaderWithoutId("table1", []string{"column1", "column2"}),
			Records: []Record{
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("2"),
					parser.NewString("str2"),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("3"),
					parser.NewString("str3"),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewInteger(4),
					parser.NewString("str4"),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewInteger(5),
					parser.NewString("str5"),
				}),
			},
			OperatedRecords: 2,
		},
	},
	{
		Name: "Insert Query File Does Not Exist Error",
		Query: parser.InsertQuery{
			Insert: "insert",
			Into:   "into",
			Table:  parser.Identifier{Literal: "notexist"},
			Fields: []parser.Expression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
			},
			Values: "values",
			ValuesList: []parser.Expression{
				parser.RowValue{
					Value: parser.ValueList{
						Values: []parser.Expression{
							parser.NewInteger(4),
						},
					},
				},
				parser.RowValue{
					Value: parser.ValueList{
						Values: []parser.Expression{
							parser.NewInteger(5),
						},
					},
				},
			},
		},
		Error: "file notexist does not exist",
	},
	{
		Name: "Insert Query Field Does Not Exist Error",
		Query: parser.InsertQuery{
			Insert: "insert",
			Into:   "into",
			Table:  parser.Identifier{Literal: "table1"},
			Fields: []parser.Expression{
				parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			},
			Values: "values",
			ValuesList: []parser.Expression{
				parser.RowValue{
					Value: parser.ValueList{
						Values: []parser.Expression{
							parser.NewInteger(4),
						},
					},
				},
				parser.RowValue{
					Value: parser.ValueList{
						Values: []parser.Expression{
							parser.NewInteger(5),
						},
					},
				},
			},
		},
		Error: "field notexist does not exist",
	},
	{
		Name: "Insert Select Query",
		Query: parser.InsertQuery{
			Insert: "insert",
			Into:   "into",
			Table:  parser.Identifier{Literal: "table1"},
			Fields: []parser.Expression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
			},
			Query: parser.SelectQuery{
				SelectEntity: parser.SelectEntity{
					SelectClause: parser.SelectClause{
						Fields: []parser.Expression{
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column3"}}},
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column4"}}},
						},
					},
					FromClause: parser.FromClause{
						Tables: []parser.Expression{
							parser.Table{Object: parser.Identifier{Literal: "table2"}},
						},
					},
				},
			},
		},
		Result: &View{
			FileInfo: &FileInfo{
				Path:      GetTestFilePath("table1.csv"),
				Delimiter: ',',
				NoHeader:  false,
				Encoding:  cmd.UTF8,
				LineBreak: cmd.LF,
			},
			Header: NewHeaderWithoutId("table1", []string{"column1", "column2"}),
			Records: []Record{
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("2"),
					parser.NewString("str2"),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("3"),
					parser.NewString("str3"),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("2"),
					parser.NewString("str22"),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("3"),
					parser.NewString("str33"),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("4"),
					parser.NewString("str44"),
				}),
			},
			OperatedRecords: 3,
		},
	},
	{
		Name: "Insert Select Query Field Does Not Exist Error",
		Query: parser.InsertQuery{
			Insert: "insert",
			Into:   "into",
			Table:  parser.Identifier{Literal: "table1"},
			Fields: []parser.Expression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
			},
			Query: parser.SelectQuery{
				SelectEntity: parser.SelectEntity{
					SelectClause: parser.SelectClause{
						Fields: []parser.Expression{
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column3"}}},
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column4"}}},
						},
					},
					FromClause: parser.FromClause{
						Tables: []parser.Expression{
							parser.Table{Object: parser.Identifier{Literal: "table2"}},
						},
					},
				},
			},
		},
		Error: "field length does not match value length",
	},
}

func TestInsert(t *testing.T) {
	tf := cmd.GetFlags()
	tf.Repository = TestDir

	for _, v := range insertTests {
		ViewCache.Clear()
		result, err := Insert(v.Query)
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
		if !reflect.DeepEqual(result, v.Result) {
			t.Errorf("%s: result = %q, want %q", v.Name, result, v.Result)
		}
	}
}

var updateTests = []struct {
	Name   string
	Query  parser.UpdateQuery
	Result []*View
	Error  string
}{
	{
		Name: "Update Query",
		Query: parser.UpdateQuery{
			Update: "update",
			Tables: []parser.Expression{
				parser.Table{Object: parser.Identifier{Literal: "table1"}},
			},
			Set: "set",
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
		Result: []*View{
			{
				FileInfo: &FileInfo{
					Path:      GetTestFilePath("table1.csv"),
					Delimiter: ',',
					NoHeader:  false,
					Encoding:  cmd.UTF8,
					LineBreak: cmd.LF,
				},
				Header: NewHeaderWithoutId("table1", []string{"column1", "column2"}),
				Records: []Record{
					NewRecordWithoutId([]parser.Primary{
						parser.NewString("1"),
						parser.NewString("str1"),
					}),
					NewRecordWithoutId([]parser.Primary{
						parser.NewString("2"),
						parser.NewString("update"),
					}),
					NewRecordWithoutId([]parser.Primary{
						parser.NewString("3"),
						parser.NewString("str3"),
					}),
				},
				OperatedRecords: 1,
			},
		},
	},
	{
		Name: "Update Query Multiple Table",
		Query: parser.UpdateQuery{
			Update: "update",
			Tables: []parser.Expression{
				parser.Table{Object: parser.Identifier{Literal: "t1"}},
			},
			Set: "set",
			SetList: []parser.Expression{
				parser.UpdateSet{
					Field: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
					Value: parser.FieldReference{Column: parser.Identifier{Literal: "column4"}},
				},
			},
			FromClause: parser.FromClause{
				Tables: []parser.Expression{
					parser.Table{Object: parser.Join{
						Table: parser.Table{
							Object: parser.Identifier{Literal: "table1"},
							Alias:  parser.Identifier{Literal: "t1"},
						},
						JoinTable: parser.Table{
							Object: parser.Identifier{Literal: "table2"},
							Alias:  parser.Identifier{Literal: "t2"},
						},
						Condition: parser.JoinCondition{
							On: parser.Comparison{
								LHS:      parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
								RHS:      parser.FieldReference{Column: parser.Identifier{Literal: "column3"}},
								Operator: "=",
							},
						},
					}},
				},
			},
		},
		Result: []*View{
			{
				FileInfo: &FileInfo{
					Path:      GetTestFilePath("table1.csv"),
					Delimiter: ',',
					NoHeader:  false,
					Encoding:  cmd.UTF8,
					LineBreak: cmd.LF,
				},
				Header: NewHeaderWithoutId("t1", []string{"column1", "column2"}),
				Records: []Record{
					NewRecordWithoutId([]parser.Primary{
						parser.NewString("1"),
						parser.NewString("str1"),
					}),
					NewRecordWithoutId([]parser.Primary{
						parser.NewString("2"),
						parser.NewString("str22"),
					}),
					NewRecordWithoutId([]parser.Primary{
						parser.NewString("3"),
						parser.NewString("str33"),
					}),
				},
				OperatedRecords: 2,
			},
		},
	},
	{
		Name: "Update Query File Does Not Exist Error",
		Query: parser.UpdateQuery{
			Update: "update",
			Tables: []parser.Expression{
				parser.Table{Object: parser.Identifier{Literal: "notexist"}},
			},
			Set: "set",
			SetList: []parser.Expression{
				parser.UpdateSet{
					Field: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
					Value: parser.NewString("update"),
				},
			},
			WhereClause: parser.WhereClause{
				Filter: parser.Comparison{
					LHS:      parser.Identifier{Literal: "column1"},
					RHS:      parser.NewInteger(2),
					Operator: "=",
				},
			},
		},
		Error: "file notexist does not exist",
	},
	{
		Name: "Update Query Filter Error",
		Query: parser.UpdateQuery{
			Update: "update",
			Tables: []parser.Expression{
				parser.Table{Object: parser.Identifier{Literal: "table1"}},
			},
			Set: "set",
			SetList: []parser.Expression{
				parser.UpdateSet{
					Field: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
					Value: parser.NewString("update"),
				},
			},
			WhereClause: parser.WhereClause{
				Filter: parser.Comparison{
					LHS:      parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
					RHS:      parser.NewInteger(2),
					Operator: "=",
				},
			},
		},
		Error: "field notexist does not exist",
	},
	{
		Name: "Update Query File Is Not Loaded Error",
		Query: parser.UpdateQuery{
			Update: "update",
			Tables: []parser.Expression{
				parser.Table{Object: parser.Identifier{Literal: "notexist"}},
			},
			Set: "set",
			SetList: []parser.Expression{
				parser.UpdateSet{
					Field: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
					Value: parser.FieldReference{Column: parser.Identifier{Literal: "column4"}},
				},
			},
			FromClause: parser.FromClause{
				Tables: []parser.Expression{
					parser.Table{Object: parser.Join{
						Table: parser.Table{
							Object: parser.Identifier{Literal: "table1"},
							Alias:  parser.Identifier{Literal: "t1"},
						},
						JoinTable: parser.Table{
							Object: parser.Identifier{Literal: "table2"},
							Alias:  parser.Identifier{Literal: "t2"},
						},
						Condition: parser.JoinCondition{
							On: parser.Comparison{
								LHS:      parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
								RHS:      parser.FieldReference{Column: parser.Identifier{Literal: "column3"}},
								Operator: "=",
							},
						},
					}},
				},
			},
		},
		Error: "file notexist is not loaded",
	},
	{
		Name: "Update Query Update Table Is Not Specified Error",
		Query: parser.UpdateQuery{
			Update: "update",
			Tables: []parser.Expression{
				parser.Table{Object: parser.Identifier{Literal: "table1"}},
			},
			Set: "set",
			SetList: []parser.Expression{
				parser.UpdateSet{
					Field: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
					Value: parser.FieldReference{Column: parser.Identifier{Literal: "column4"}},
				},
			},
			FromClause: parser.FromClause{
				Tables: []parser.Expression{
					parser.Table{Object: parser.Join{
						Table: parser.Table{
							Object: parser.Identifier{Literal: "table1"},
							Alias:  parser.Identifier{Literal: "t1"},
						},
						JoinTable: parser.Table{
							Object: parser.Identifier{Literal: "table2"},
							Alias:  parser.Identifier{Literal: "t2"},
						},
						Condition: parser.JoinCondition{
							On: parser.Comparison{
								LHS:      parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
								RHS:      parser.FieldReference{Column: parser.Identifier{Literal: "column3"}},
								Operator: "=",
							},
						},
					}},
				},
			},
		},
		Error: "table t1 is not specified in tables to update",
	},
	{
		Name: "Update Query Update Field Error",
		Query: parser.UpdateQuery{
			Update: "update",
			Tables: []parser.Expression{
				parser.Table{Object: parser.Identifier{Literal: "table1"}},
			},
			Set: "set",
			SetList: []parser.Expression{
				parser.UpdateSet{
					Field: parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
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
		Error: "field notexist does not exist",
	},
	{
		Name: "Update Query Update Value Error",
		Query: parser.UpdateQuery{
			Update: "update",
			Tables: []parser.Expression{
				parser.Table{Object: parser.Identifier{Literal: "table1"}},
			},
			Set: "set",
			SetList: []parser.Expression{
				parser.UpdateSet{
					Field: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
					Value: parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
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
		Error: "field notexist does not exist",
	},
	{
		Name: "Update Query Record Is Ambiguous Error",
		Query: parser.UpdateQuery{
			Update: "update",
			Tables: []parser.Expression{
				parser.Table{Object: parser.Identifier{Literal: "t1"}},
			},
			Set: "set",
			SetList: []parser.Expression{
				parser.UpdateSet{
					Field: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
					Value: parser.FieldReference{Column: parser.Identifier{Literal: "column4"}},
				},
			},
			FromClause: parser.FromClause{
				Tables: []parser.Expression{
					parser.Table{Object: parser.Join{
						Table: parser.Table{
							Object: parser.Identifier{Literal: "table1"},
							Alias:  parser.Identifier{Literal: "t1"},
						},
						JoinTable: parser.Table{
							Object: parser.Identifier{Literal: "table2"},
							Alias:  parser.Identifier{Literal: "t2"},
						},
						JoinType: parser.Token{Token: parser.CROSS, Literal: "cross"},
					}},
				},
			},
		},
		Error: "record to update is ambiguous",
	},
}

func TestUpdate(t *testing.T) {
	tf := cmd.GetFlags()
	tf.Repository = TestDir

	for _, v := range updateTests {
		ViewCache.Clear()
		result, err := Update(v.Query)
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
		if !reflect.DeepEqual(result, v.Result) {
			t.Errorf("%s: result = %q, want %q", v.Name, result, v.Result)
		}
	}
}

var deleteTests = []struct {
	Name   string
	Query  parser.DeleteQuery
	Result []*View
	Error  string
}{
	{
		Name: "Delete Query",
		Query: parser.DeleteQuery{
			Delete: "delete",
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
		Result: []*View{
			{
				FileInfo: &FileInfo{
					Path:      GetTestFilePath("table1.csv"),
					Delimiter: ',',
					NoHeader:  false,
					Encoding:  cmd.UTF8,
					LineBreak: cmd.LF,
				},
				Header: NewHeaderWithoutId("table1", []string{"column1", "column2"}),
				Records: []Record{
					NewRecordWithoutId([]parser.Primary{
						parser.NewString("1"),
						parser.NewString("str1"),
					}),
					NewRecordWithoutId([]parser.Primary{
						parser.NewString("3"),
						parser.NewString("str3"),
					}),
				},
				OperatedRecords: 1,
			},
		},
	},
	{
		Name: "Delete Query Multiple Table",
		Query: parser.DeleteQuery{
			Delete: "delete",
			Tables: []parser.Expression{
				parser.Table{Object: parser.Identifier{Literal: "t1"}},
			},
			FromClause: parser.FromClause{
				Tables: []parser.Expression{
					parser.Table{Object: parser.Join{
						Table: parser.Table{
							Object: parser.Identifier{Literal: "table1"},
							Alias:  parser.Identifier{Literal: "t1"},
						},
						JoinTable: parser.Table{
							Object: parser.Identifier{Literal: "table2"},
							Alias:  parser.Identifier{Literal: "t2"},
						},
						Condition: parser.JoinCondition{
							On: parser.Comparison{
								LHS:      parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
								RHS:      parser.FieldReference{Column: parser.Identifier{Literal: "column3"}},
								Operator: "=",
							},
						},
					}},
				},
			},
		},
		Result: []*View{
			{
				FileInfo: &FileInfo{
					Path:      GetTestFilePath("table1.csv"),
					Delimiter: ',',
					NoHeader:  false,
					Encoding:  cmd.UTF8,
					LineBreak: cmd.LF,
				},
				Header: NewHeaderWithoutId("t1", []string{"column1", "column2"}),
				Records: []Record{
					NewRecordWithoutId([]parser.Primary{
						parser.NewString("1"),
						parser.NewString("str1"),
					}),
				},
				OperatedRecords: 2,
			},
		},
	},
	{
		Name: "Delete Query File Is Not Specified Error",
		Query: parser.DeleteQuery{
			Delete: "delete",
			FromClause: parser.FromClause{
				Tables: []parser.Expression{
					parser.Table{Object: parser.Join{
						Table: parser.Table{
							Object: parser.Identifier{Literal: "table1"},
							Alias:  parser.Identifier{Literal: "t1"},
						},
						JoinTable: parser.Table{
							Object: parser.Identifier{Literal: "table2"},
							Alias:  parser.Identifier{Literal: "t2"},
						},
						Condition: parser.JoinCondition{
							On: parser.Comparison{
								LHS:      parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
								RHS:      parser.FieldReference{Column: parser.Identifier{Literal: "column3"}},
								Operator: "=",
							},
						},
					}},
				},
			},
		},
		Error: "update file is not specified",
	},
	{
		Name: "Delete Query File Does Not Exist Error",
		Query: parser.DeleteQuery{
			Delete: "delete",
			FromClause: parser.FromClause{
				Tables: []parser.Expression{
					parser.Table{
						Object: parser.Identifier{Literal: "notexist"},
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
		Error: "file notexist does not exist",
	},
	{
		Name: "Delete Query Filter Error",
		Query: parser.DeleteQuery{
			Delete: "delete",
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
					RHS:      parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
					Operator: "=",
				},
			},
		},
		Error: "field notexist does not exist",
	},
	{
		Name: "Delete Query File Is Not Loaded Error",
		Query: parser.DeleteQuery{
			Delete: "delete",
			Tables: []parser.Expression{
				parser.Table{Object: parser.Identifier{Literal: "notexist"}},
			},
			FromClause: parser.FromClause{
				Tables: []parser.Expression{
					parser.Table{Object: parser.Join{
						Table: parser.Table{
							Object: parser.Identifier{Literal: "table1"},
							Alias:  parser.Identifier{Literal: "t1"},
						},
						JoinTable: parser.Table{
							Object: parser.Identifier{Literal: "table2"},
							Alias:  parser.Identifier{Literal: "t2"},
						},
						Condition: parser.JoinCondition{
							On: parser.Comparison{
								LHS:      parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
								RHS:      parser.FieldReference{Column: parser.Identifier{Literal: "column3"}},
								Operator: "=",
							},
						},
					}},
				},
			},
		},
		Error: "file notexist is not loaded",
	},
}

func TestDelete(t *testing.T) {
	tf := cmd.GetFlags()
	tf.Repository = TestDir

	for _, v := range deleteTests {
		ViewCache.Clear()
		result, err := Delete(v.Query)
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
		if !reflect.DeepEqual(result, v.Result) {
			t.Errorf("%s: result = %q, want %q", v.Name, result, v.Result)
		}
	}
}

var createTableTests = []struct {
	Name   string
	Query  parser.CreateTable
	Result *View
	Error  string
}{
	{
		Name: "Create Table",
		Query: parser.CreateTable{
			Table: parser.Identifier{Literal: "create_table.csv"},
			Fields: []parser.Expression{
				parser.Identifier{Literal: "column1"},
				parser.Identifier{Literal: "column2"},
			},
		},
		Result: &View{
			FileInfo: &FileInfo{
				Path:      GetTestFilePath("create_table.csv"),
				Delimiter: ',',
				NoHeader:  false,
				Encoding:  cmd.UTF8,
				LineBreak: cmd.LF,
			},
			Header: NewHeaderWithoutId("create_table", []string{"column1", "column2"}),
		},
	},
	{
		Name: "Create Table Field Duplicate Error",
		Query: parser.CreateTable{
			Table: parser.Identifier{Literal: "create_table.csv"},
			Fields: []parser.Expression{
				parser.Identifier{Literal: "column1"},
				parser.Identifier{Literal: "column1"},
			},
		},
		Error: "field column1 is duplicate",
	},
}

func TestCreateTable(t *testing.T) {
	tf := cmd.GetFlags()
	tf.Repository = TestDir

	for _, v := range createTableTests {
		ViewCache.Clear()
		result, err := CreateTable(v.Query)
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
		if !reflect.DeepEqual(result, v.Result) {
			t.Errorf("%s: result = %q, want %q", v.Name, result, v.Result)
		}
	}
}

var addColumnsTests = []struct {
	Name   string
	Query  parser.AddColumns
	Result *View
	Error  string
}{
	{
		Name: "Add Columns",
		Query: parser.AddColumns{
			Table: parser.Identifier{Literal: "table1.csv"},
			Columns: []parser.Expression{
				parser.ColumnDefault{
					Column: parser.Identifier{Literal: "column3"},
				},
				parser.ColumnDefault{
					Column: parser.Identifier{Literal: "column4"},
				},
			},
		},
		Result: &View{
			FileInfo: &FileInfo{
				Path:      GetTestFilePath("table1.csv"),
				Delimiter: ',',
				NoHeader:  false,
				Encoding:  cmd.UTF8,
				LineBreak: cmd.LF,
			},
			Header: NewHeaderWithoutId("table1", []string{"column1", "column2", "column3", "column4"}),
			Records: []Record{
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
					parser.NewNull(),
					parser.NewNull(),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("2"),
					parser.NewString("str2"),
					parser.NewNull(),
					parser.NewNull(),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("3"),
					parser.NewString("str3"),
					parser.NewNull(),
					parser.NewNull(),
				}),
			},
			OperatedFields: 2,
		},
	},
	{
		Name: "Add Columns First",
		Query: parser.AddColumns{
			Table: parser.Identifier{Literal: "table1.csv"},
			Columns: []parser.Expression{
				parser.ColumnDefault{
					Column: parser.Identifier{Literal: "column3"},
					Value:  parser.Function{Name: "auto_increment"},
				},
				parser.ColumnDefault{
					Column: parser.Identifier{Literal: "column4"},
					Value:  parser.NewInteger(1),
				},
			},
			Position: parser.ColumnPosition{
				Position: parser.Token{Token: parser.FIRST},
			},
		},
		Result: &View{
			FileInfo: &FileInfo{
				Path:      GetTestFilePath("table1.csv"),
				Delimiter: ',',
				NoHeader:  false,
				Encoding:  cmd.UTF8,
				LineBreak: cmd.LF,
			},
			Header: NewHeaderWithoutId("table1", []string{"column3", "column4", "column1", "column2"}),
			Records: []Record{
				NewRecordWithoutId([]parser.Primary{
					parser.NewInteger(1),
					parser.NewInteger(1),
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewInteger(2),
					parser.NewInteger(1),
					parser.NewString("2"),
					parser.NewString("str2"),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewInteger(3),
					parser.NewInteger(1),
					parser.NewString("3"),
					parser.NewString("str3"),
				}),
			},
			OperatedFields: 2,
		},
	},
	{
		Name: "Add Columns After",
		Query: parser.AddColumns{
			Table: parser.Identifier{Literal: "table1.csv"},
			Columns: []parser.Expression{
				parser.ColumnDefault{
					Column: parser.Identifier{Literal: "column3"},
				},
				parser.ColumnDefault{
					Column: parser.Identifier{Literal: "column4"},
					Value:  parser.NewInteger(1),
				},
			},
			Position: parser.ColumnPosition{
				Position: parser.Token{Token: parser.AFTER},
				Column:   parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
			},
		},
		Result: &View{
			FileInfo: &FileInfo{
				Path:      GetTestFilePath("table1.csv"),
				Delimiter: ',',
				NoHeader:  false,
				Encoding:  cmd.UTF8,
				LineBreak: cmd.LF,
			},
			Header: NewHeaderWithoutId("table1", []string{"column1", "column3", "column4", "column2"}),
			Records: []Record{
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("1"),
					parser.NewNull(),
					parser.NewInteger(1),
					parser.NewString("str1"),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("2"),
					parser.NewNull(),
					parser.NewInteger(1),
					parser.NewString("str2"),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("3"),
					parser.NewNull(),
					parser.NewInteger(1),
					parser.NewString("str3"),
				}),
			},
			OperatedFields: 2,
		},
	},
	{
		Name: "Add Columns Before",
		Query: parser.AddColumns{
			Table: parser.Identifier{Literal: "table1.csv"},
			Columns: []parser.Expression{
				parser.ColumnDefault{
					Column: parser.Identifier{Literal: "column3"},
				},
				parser.ColumnDefault{
					Column: parser.Identifier{Literal: "column4"},
					Value:  parser.NewInteger(1),
				},
			},
			Position: parser.ColumnPosition{
				Position: parser.Token{Token: parser.BEFORE},
				Column:   parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
			},
		},
		Result: &View{
			FileInfo: &FileInfo{
				Path:      GetTestFilePath("table1.csv"),
				Delimiter: ',',
				NoHeader:  false,
				Encoding:  cmd.UTF8,
				LineBreak: cmd.LF,
			},
			Header: NewHeaderWithoutId("table1", []string{"column1", "column3", "column4", "column2"}),
			Records: []Record{
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("1"),
					parser.NewNull(),
					parser.NewInteger(1),
					parser.NewString("str1"),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("2"),
					parser.NewNull(),
					parser.NewInteger(1),
					parser.NewString("str2"),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("3"),
					parser.NewNull(),
					parser.NewInteger(1),
					parser.NewString("str3"),
				}),
			},
			OperatedFields: 2,
		},
	},
	{
		Name: "Add Columns Load Error",
		Query: parser.AddColumns{
			Table: parser.Identifier{Literal: "notexist"},
			Columns: []parser.Expression{
				parser.ColumnDefault{
					Column: parser.Identifier{Literal: "column3"},
				},
				parser.ColumnDefault{
					Column: parser.Identifier{Literal: "column4"},
				},
			},
		},
		Error: "file notexist does not exist",
	},
	{
		Name: "Add Columns Position Column Does Not Exist Error",
		Query: parser.AddColumns{
			Table: parser.Identifier{Literal: "table1.csv"},
			Columns: []parser.Expression{
				parser.ColumnDefault{
					Column: parser.Identifier{Literal: "column3"},
				},
				parser.ColumnDefault{
					Column: parser.Identifier{Literal: "column2"},
					Value:  parser.NewInteger(1),
				},
			},
			Position: parser.ColumnPosition{
				Position: parser.Token{Token: parser.BEFORE},
				Column:   parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			},
		},
		Error: "field notexist does not exist",
	},
	{
		Name: "Add Columns Field Duplicate Error",
		Query: parser.AddColumns{
			Table: parser.Identifier{Literal: "table1.csv"},
			Columns: []parser.Expression{
				parser.ColumnDefault{
					Column: parser.Identifier{Literal: "column3"},
				},
				parser.ColumnDefault{
					Column: parser.Identifier{Literal: "column1"},
					Value:  parser.NewInteger(1),
				},
			},
		},
		Error: "field column1 is duplicate",
	},
	{
		Name: "Add Columns Default Value Error",
		Query: parser.AddColumns{
			Table: parser.Identifier{Literal: "table1.csv"},
			Columns: []parser.Expression{
				parser.ColumnDefault{
					Column: parser.Identifier{Literal: "column3"},
				},
				parser.ColumnDefault{
					Column: parser.Identifier{Literal: "column4"},
					Value:  parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
				},
			},
		},
		Error: "field notexist does not exist",
	},
}

func TestAddColumns(t *testing.T) {
	tf := cmd.GetFlags()
	tf.Repository = TestDir

	for _, v := range addColumnsTests {
		ViewCache.Clear()
		result, err := AddColumns(v.Query)
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
		if !reflect.DeepEqual(result, v.Result) {
			t.Errorf("%s: result = %q, want %q", v.Name, result, v.Result)
		}
	}
}

var dropColumnsTests = []struct {
	Name   string
	Query  parser.DropColumns
	Result *View
	Error  string
}{
	{
		Name: "Drop Columns",
		Query: parser.DropColumns{
			Table: parser.Identifier{Literal: "table1"},
			Columns: []parser.Expression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
			},
		},
		Result: &View{
			FileInfo: &FileInfo{
				Path:      GetTestFilePath("table1.csv"),
				Delimiter: ',',
				NoHeader:  false,
				Encoding:  cmd.UTF8,
				LineBreak: cmd.LF,
			},
			Header: NewHeaderWithoutId("table1", []string{"column1"}),
			Records: []Record{
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("1"),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("2"),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("3"),
				}),
			},
			OperatedFields: 1,
		},
	},
	{
		Name: "Drop Columns Load Error",
		Query: parser.DropColumns{
			Table: parser.Identifier{Literal: "notexist"},
			Columns: []parser.Expression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
			},
		},
		Error: "file notexist does not exist",
	},
	{
		Name: "Drop Columns Field Does Not Exist Error",
		Query: parser.DropColumns{
			Table: parser.Identifier{Literal: "table1"},
			Columns: []parser.Expression{
				parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			},
		},
		Error: "field notexist does not exist",
	},
}

func TestDropColumns(t *testing.T) {
	tf := cmd.GetFlags()
	tf.Repository = TestDir

	for _, v := range dropColumnsTests {
		ViewCache.Clear()
		result, err := DropColumns(v.Query)
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
		if !reflect.DeepEqual(result, v.Result) {
			t.Errorf("%s: result = %q, want %q", v.Name, result, v.Result)
		}
	}
}

var renameColumnTests = []struct {
	Name   string
	Query  parser.RenameColumn
	Result *View
	Error  string
}{
	{
		Name: "Rename Column",
		Query: parser.RenameColumn{
			Table: parser.Identifier{Literal: "table1"},
			Old:   parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
			New:   parser.Identifier{Literal: "newcolumn"},
		},
		Result: &View{
			FileInfo: &FileInfo{
				Path:      GetTestFilePath("table1.csv"),
				Delimiter: ',',
				NoHeader:  false,
				Encoding:  cmd.UTF8,
				LineBreak: cmd.LF,
			},
			Header: NewHeaderWithoutId("table1", []string{"column1", "newcolumn"}),
			Records: []Record{
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("2"),
					parser.NewString("str2"),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("3"),
					parser.NewString("str3"),
				}),
			},
			OperatedFields: 1,
		},
	},
	{
		Name: "Rename Column Load Error",
		Query: parser.RenameColumn{
			Table: parser.Identifier{Literal: "notexist"},
			Old:   parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
			New:   parser.Identifier{Literal: "newcolumn"},
		},
		Error: "file notexist does not exist",
	},
	{
		Name: "Rename Column Field Duplicate Error",
		Query: parser.RenameColumn{
			Table: parser.Identifier{Literal: "table1"},
			Old:   parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
			New:   parser.Identifier{Literal: "column1"},
		},
		Error: "field column1 is duplicate",
	},
	{
		Name: "Rename Column Field Does Not Exist Error",
		Query: parser.RenameColumn{
			Table: parser.Identifier{Literal: "table1"},
			Old:   parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			New:   parser.Identifier{Literal: "newcolumn"},
		},
		Error: "field notexist does not exist",
	},
}

func TestRenameColumn(t *testing.T) {
	tf := cmd.GetFlags()
	tf.Repository = TestDir

	for _, v := range renameColumnTests {
		ViewCache.Clear()
		result, err := RenameColumn(v.Query)
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
		if !reflect.DeepEqual(result, v.Result) {
			t.Errorf("%s: result = %q, want %q", v.Name, result, v.Result)
		}
	}
}

var printTests = []struct {
	Name   string
	Query  parser.Print
	Result string
	Error  string
}{
	{
		Name: "Print",
		Query: parser.Print{
			Value: parser.NewString("foo"),
		},
		Result: "'foo'",
	},
	{
		Name: "Print Error",
		Query: parser.Print{
			Value: parser.Variable{
				Name: "var",
			},
		},
		Error: "variable var is undefined",
	},
}

func TestPrint(t *testing.T) {
	for _, v := range printTests {
		result, err := Print(v.Query)
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
		if result != v.Result {
			t.Errorf("%s: result = %q, want %q", v.Name, result, v.Result)
		}
	}
}

var setFlagTests = []struct {
	Name            string
	Query           parser.SetFlag
	ResultFlag      string
	ResultStlValue  string
	ResultBoolValue bool
	Error           string
}{
	{
		Name: "Set Delimiter",
		Query: parser.SetFlag{
			Name:  "@@delimiter",
			Value: parser.NewString("\t"),
		},
		ResultFlag:     "delimiter",
		ResultStlValue: "\t",
	},
	{
		Name: "Set Encoding",
		Query: parser.SetFlag{
			Name:  "@@encoding",
			Value: parser.NewString("SJIS"),
		},
		ResultFlag:     "encoding",
		ResultStlValue: "SJIS",
	},
	{
		Name: "Set Repository",
		Query: parser.SetFlag{
			Name:  "@@repository",
			Value: parser.NewString(TestDir),
		},
		ResultFlag:     "repository",
		ResultStlValue: TestDir,
	},
	{
		Name: "Set NoHeader",
		Query: parser.SetFlag{
			Name:  "@@no_header",
			Value: parser.NewBoolean(true),
		},
		ResultFlag:      "no_header",
		ResultBoolValue: true,
	},
	{
		Name: "Set WithoutNull",
		Query: parser.SetFlag{
			Name:  "@@without_null",
			Value: parser.NewBoolean(true),
		},
		ResultFlag:      "without_null",
		ResultBoolValue: true,
	},
	{
		Name: "Set Delimiter Value Error",
		Query: parser.SetFlag{
			Name:  "@@delimiter",
			Value: parser.NewBoolean(true),
		},
		Error: "invalid flag value: @@delimiter = true",
	},
	{
		Name: "Set WithoutNull Value Error",
		Query: parser.SetFlag{
			Name:  "@@without_null",
			Value: parser.NewString("string"),
		},
		Error: "invalid flag value: @@without_null = 'string'",
	},
	{
		Name: "Invalid Flag Error",
		Query: parser.SetFlag{
			Name:  "@@invalid",
			Value: parser.NewString("string"),
		},
		Error: "invalid flag name: @@invalid",
	},
}

func TestSetFlag(t *testing.T) {
	flags := cmd.GetFlags()

	for _, v := range setFlagTests {
		err := SetFlag(v.Query)
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

		switch strings.ToUpper(v.ResultFlag) {
		case "DELIMITER":
			if string(flags.Delimiter) != v.ResultStlValue {
				t.Errorf("%s: delimiter = %q, want %q", v.Name, string(flags.Delimiter), v.ResultStlValue)
			}
		case "ENCODING":
			if flags.Encoding.String() != v.ResultStlValue {
				t.Errorf("%s: encoding = %q, want %q", v.Name, flags.Encoding.String(), v.ResultStlValue)
			}
		case "REPOSITORY":
			if flags.Repository != v.ResultStlValue {
				t.Errorf("%s: repository = %q, want %q", v.Name, flags.Repository, v.ResultStlValue)
			}
		case "NO-HEADER":
			if flags.NoHeader != v.ResultBoolValue {
				t.Errorf("%s: no-header = %t, want %t", v.Name, flags.NoHeader, v.ResultBoolValue)
			}
		case "WITHOUT-NULL":
			if flags.WithoutNull != v.ResultBoolValue {
				t.Errorf("%s: without-null = %t, want %t", v.Name, flags.WithoutNull, v.ResultBoolValue)
			}
		}
	}
}
