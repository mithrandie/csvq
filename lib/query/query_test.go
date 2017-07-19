package query

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/parser"
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
		Output: fmt.Sprintf("%d records inserted on %q.\n", 2, GetTestFilePath("insert_query.csv")) +
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
		Output: fmt.Sprintf("%d record updated on %q.\n", 1, GetTestFilePath("update_query.csv")) +
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
		Output: fmt.Sprintf("no record updated on %q.\n", GetTestFilePath("update_query.csv")),
	},
	{
		Name:  "Delete Query",
		Input: "delete from delete_query where column1 = 2",
		Output: fmt.Sprintf("%d record deleted on %q.\n", 1, GetTestFilePath("delete_query.csv")) +
			fmt.Sprintf("Commit: file %q is updated.\n", GetTestFilePath("delete_query.csv")),
		UpdateFile: GetTestFilePath("delete_query.csv"),
		Content: "\"column1\",\"column2\"\n" +
			"\"1\",\"str1\"\n" +
			"\"3\",\"str3\"",
	},
	{
		Name:   "Delete Query No Record Deleted",
		Input:  "delete from delete_query where false",
		Output: fmt.Sprintf("no record deleted on %q.\n", GetTestFilePath("delete_query.csv")),
	},
	{
		Name:  "Create Table",
		Input: "create table `create_table.csv` (column1, column2)",
		Output: fmt.Sprintf("file %q is created.\n", GetTestFilePath("create_table.csv")) +
			fmt.Sprintf("Commit: file %q is created.\n", GetTestFilePath("create_table.csv")),
		UpdateFile: GetTestFilePath("create_table.csv"),
		Content:    "\"column1\",\"column2\"\n",
	},
	{
		Name:  "Add Columns",
		Input: "alter table add_columns add column3",
		Output: fmt.Sprintf("%d field added on %q.\n", 1, GetTestFilePath("add_columns.csv")) +
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
		Output: fmt.Sprintf("%d field dropped on %q.\n", 1, GetTestFilePath("drop_columns.csv")) +
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
		Output: fmt.Sprintf("%d field renamed on %q.\n", 1, GetTestFilePath("rename_column.csv")) +
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
		Error: "[L:1 C:8] syntax error: unexpected from",
	},
}

func TestExecute(t *testing.T) {
	tf := cmd.GetFlags()
	tf.Format = cmd.TEXT
	tf.Repository = TestDir

	for _, v := range executeTests {
		Logs = []string{}
		out, err := Execute(v.Input, "")

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

var fetchCursorTests = []struct {
	Name          string
	CurName       parser.Identifier
	FetchPosition parser.Expression
	Variables     []parser.Variable
	Success       bool
	ResultVars    Variables
	Error         string
}{
	{
		Name:    "Fetch Cursor First Time",
		CurName: parser.Identifier{Literal: "cur"},
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
		CurName: parser.Identifier{Literal: "cur"},
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
		CurName: parser.Identifier{Literal: "cur"},
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
		CurName: parser.Identifier{Literal: "cur"},
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
		CurName: parser.Identifier{Literal: "cur"},
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
		CurName: parser.Identifier{Literal: "notexist"},
		Variables: []parser.Variable{
			{Name: "@var1"},
			{Name: "@var2"},
		},
		Error: "[L:- C:-] cursor notexist is undefined",
	},
	{
		Name:    "Fetch Cursor Not Match Number Error",
		CurName: parser.Identifier{Literal: "cur2"},
		Variables: []parser.Variable{
			{Name: "@var1"},
		},
		Error: "[L:- C:-] fetching from cursor cur2 returns 2 values",
	},
	{
		Name:    "Fetch Cursor Substitution Error",
		CurName: parser.Identifier{Literal: "cur2"},
		Variables: []parser.Variable{
			{Name: "@var1"},
			{Name: "@notexist"},
		},
		Error: "[L:- C:-] variable @notexist is undefined",
	},
	{
		Name:    "Fetch Cursor Number Value Error",
		CurName: parser.Identifier{Literal: "cur"},
		FetchPosition: parser.FetchPosition{
			Position: parser.Token{Token: parser.ABSOLUTE, Literal: "absolute"},
			Number:   parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
		},
		Variables: []parser.Variable{
			{Name: "@var1"},
			{Name: "@var2"},
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
	{
		Name:    "Fetch Cursor Number Not Integer Error",
		CurName: parser.Identifier{Literal: "cur"},
		FetchPosition: parser.FetchPosition{
			Position: parser.Token{Token: parser.ABSOLUTE, Literal: "absolute"},
			Number:   parser.NewNull(),
		},
		Variables: []parser.Variable{
			{Name: "@var1"},
			{Name: "@var2"},
		},
		Error: "[L:- C:-] fetching position NULL is not an integer value",
	},
}

func TestFetchCursor(t *testing.T) {
	tf := cmd.GetFlags()
	tf.Repository = TestDir

	filter := NewFilter([]Variables{
		{
			"@var1": parser.NewNull(),
			"@var2": parser.NewNull(),
		},
	})

	Cursors = CursorMap{
		"CUR": &Cursor{
			query: selectQueryForCursorTest,
		},
		"CUR2": &Cursor{
			query: selectQueryForCursorTest,
		},
	}
	ViewCache.Clear()
	Cursors.Open(parser.Identifier{Literal: "cur"}, filter)
	ViewCache.Clear()
	Cursors.Open(parser.Identifier{Literal: "cur2"}, filter)

	for _, v := range fetchCursorTests {
		success, err := FetchCursor(v.CurName, v.FetchPosition, v.Variables, filter)
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
		if !reflect.DeepEqual(filter.VariablesList[0], v.ResultVars) {
			t.Errorf("%s: global vars = %q, want %q", v.Name, filter.VariablesList[0], v.ResultVars)
		}
	}
}

var declareTableTests = []struct {
	Name    string
	ViewMap ViewMap
	Expr    parser.TableDeclaration
	Result  ViewMap
	Error   string
}{
	{
		Name: "Declare Table",
		Expr: parser.TableDeclaration{
			Table: parser.Identifier{Literal: "tbl"},
			Fields: []parser.Expression{
				parser.Identifier{Literal: "column1"},
				parser.Identifier{Literal: "column2"},
			},
		},
		Result: ViewMap{
			"TBL": {
				FileInfo: &FileInfo{
					Path:      "tbl",
					Temporary: true,
				},
				Header: []HeaderField{
					{
						Reference: "tbl",
						Column:    "column1",
						FromTable: true,
					},
					{
						Reference: "tbl",
						Column:    "column2",
						FromTable: true,
					},
				},
				Records: Records{},
			},
		},
	},
	{
		Name: "Declare Table Field Duplicate Error",
		Expr: parser.TableDeclaration{
			Table: parser.Identifier{Literal: "tbl"},
			Fields: []parser.Expression{
				parser.Identifier{Literal: "column1"},
				parser.Identifier{Literal: "column1"},
			},
		},
		Error: "[L:- C:-] field name column1 is a duplicate",
	},
	{
		Name: "Declare Table From Query",
		Expr: parser.TableDeclaration{
			Table: parser.Identifier{Literal: "tbl"},
			Fields: []parser.Expression{
				parser.Identifier{Literal: "column1"},
				parser.Identifier{Literal: "column2"},
			},
			Query: parser.SelectQuery{
				SelectEntity: parser.SelectEntity{
					SelectClause: parser.SelectClause{
						Fields: []parser.Expression{
							parser.Field{Object: parser.NewInteger(1)},
							parser.Field{Object: parser.NewInteger(2)},
						},
					},
				},
			},
		},
		Result: ViewMap{
			"TBL": {
				FileInfo: &FileInfo{
					Path:      "tbl",
					Temporary: true,
				},
				Header: []HeaderField{
					{
						Reference: "tbl",
						Column:    "column1",
						FromTable: true,
					},
					{
						Reference: "tbl",
						Column:    "column2",
						FromTable: true,
					},
				},
				Records: Records{
					NewRecordWithoutId([]parser.Primary{
						parser.NewInteger(1),
						parser.NewInteger(2),
					}),
				},
			},
		},
	},
	{
		Name: "Declare Table From Query Query Error",
		Expr: parser.TableDeclaration{
			Table: parser.Identifier{Literal: "tbl"},
			Fields: []parser.Expression{
				parser.Identifier{Literal: "column1"},
				parser.Identifier{Literal: "column2"},
			},
			Query: parser.SelectQuery{
				SelectEntity: parser.SelectEntity{
					SelectClause: parser.SelectClause{
						Fields: []parser.Expression{
							parser.Field{Object: parser.NewInteger(1)},
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}}},
						},
					},
				},
			},
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
	{
		Name: "Declare Table From Query Field Update Error",
		Expr: parser.TableDeclaration{
			Table: parser.Identifier{Literal: "tbl"},
			Fields: []parser.Expression{
				parser.Identifier{Literal: "column1"},
			},
			Query: parser.SelectQuery{
				SelectEntity: parser.SelectEntity{
					SelectClause: parser.SelectClause{
						Fields: []parser.Expression{
							parser.Field{Object: parser.NewInteger(1)},
							parser.Field{Object: parser.NewInteger(2)},
						},
					},
				},
			},
		},
		Error: "[L:- C:-] select query should return exactly 1 field for temporary table tbl",
	},
	{
		Name: "Declare Table Redeclaration Error",
		ViewMap: ViewMap{
			"TBL": {
				FileInfo: &FileInfo{
					Path:      "tbl",
					Temporary: true,
				},
			},
		},
		Expr: parser.TableDeclaration{
			Table: parser.Identifier{Literal: "tbl"},
			Fields: []parser.Expression{
				parser.Identifier{Literal: "column1"},
				parser.Identifier{Literal: "column2"},
			},
		},
		Error: "[L:- C:-] temporary table tbl is redeclared",
	},
}

func TestDeclareTable(t *testing.T) {
	filter := NewEmptyFilter()

	for _, v := range declareTableTests {
		if v.ViewMap == nil {
			ViewCache.Clear()
		} else {
			ViewCache = v.ViewMap
		}

		err := DeclareTable(v.Expr, filter)
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
		if !reflect.DeepEqual(ViewCache, v.Result) {
			t.Errorf("%s: view cache = %q, want %q", v.Name, ViewCache, v.Result)
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
						parser.Field{Object: parser.AggregateFunction{Name: "count", Option: parser.AggregateOption{Args: []parser.Expression{parser.AllColumns{}}}}},
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
						LHS:      parser.AggregateFunction{Name: "count", Option: parser.AggregateOption{Args: []parser.Expression{parser.AllColumns{}}}},
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
		Error: "[L:- C:-] result set to be combined should contain exactly 2 fields",
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
		Error: "[L:- C:-] field notexist does not exist",
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
		Error: "[L:- C:-] field notexist does not exist",
	},
	{
		Name: "Inline Tables",
		Query: parser.SelectQuery{
			WithClause: parser.WithClause{
				With: "with",
				InlineTables: []parser.Expression{
					parser.InlineTable{
						Name: parser.Identifier{Literal: "it"},
						Fields: []parser.Expression{
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
		Name: "Inline Tables Recursion",
		Query: parser.SelectQuery{
			WithClause: parser.WithClause{
				With: "with",
				InlineTables: []parser.Expression{
					parser.InlineTable{
						Recursive: parser.Token{Token: parser.RECURSIVE, Literal: "recursive"},
						Name:      parser.Identifier{Literal: "it"},
						Fields: []parser.Expression{
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
		Name: "Inline Tables Recursion Field Length Error",
		Query: parser.SelectQuery{
			WithClause: parser.WithClause{
				With: "with",
				InlineTables: []parser.Expression{
					parser.InlineTable{
						Recursive: parser.Token{Token: parser.RECURSIVE, Literal: "recursive"},
						Name:      parser.Identifier{Literal: "it"},
						Fields: []parser.Expression{
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
		Error: "[L:- C:-] result set to be combined should contain exactly 1 field",
	},
}

func TestSelect(t *testing.T) {
	tf := cmd.GetFlags()
	tf.Repository = TestDir

	filter := NewEmptyFilter()

	for _, v := range selectTests {
		ViewCache.Clear()
		result, err := Select(v.Query, filter)
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
		Error: "[L:- C:-] file notexist does not exist",
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
		Error: "[L:- C:-] field notexist does not exist",
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
		Error: "[L:- C:-] select query should return exactly 1 field",
	},
}

func TestInsert(t *testing.T) {
	tf := cmd.GetFlags()
	tf.Repository = TestDir

	filter := NewEmptyFilter()

	for _, v := range insertTests {
		ViewCache.Clear()
		result, err := Insert(v.Query, filter)
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
		Error: "[L:- C:-] file notexist does not exist",
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
		Error: "[L:- C:-] field notexist does not exist",
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
		Error: "[L:- C:-] table notexist is not loaded",
	},
	{
		Name: "Update Query Update Table Is Not Specified Error",
		Query: parser.UpdateQuery{
			Update: "update",
			Tables: []parser.Expression{
				parser.Table{Object: parser.Identifier{Literal: "t2"}},
			},
			Set: "set",
			SetList: []parser.Expression{
				parser.UpdateSet{
					Field: parser.FieldReference{View: parser.Identifier{Literal: "t1"}, Column: parser.Identifier{Literal: "column2"}},
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
		Error: "[L:- C:-] field t1.column2 does not exist in the tables to update",
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
		Error: "[L:- C:-] field notexist does not exist",
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
		Error: "[L:- C:-] field notexist does not exist",
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
		Error: "[L:- C:-] value column4 to set in the field column2 is ambiguous",
	},
}

func TestUpdate(t *testing.T) {
	tf := cmd.GetFlags()
	tf.Repository = TestDir

	filter := NewEmptyFilter()

	for _, v := range updateTests {
		ViewCache.Clear()
		result, err := Update(v.Query, filter)
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
		Name: "Delete Query Tables Not Specified Error",
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
		Error: "[L:- C:-] tables to delete records are not specified",
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
		Error: "[L:- C:-] file notexist does not exist",
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
		Error: "[L:- C:-] field notexist does not exist",
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
		Error: "[L:- C:-] table notexist is not loaded",
	},
}

func TestDelete(t *testing.T) {
	tf := cmd.GetFlags()
	tf.Repository = TestDir

	filter := NewEmptyFilter()

	for _, v := range deleteTests {
		ViewCache.Clear()
		result, err := Delete(v.Query, filter)
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
		Error: "[L:- C:-] field name column1 is a duplicate",
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
					Value:  parser.NewInteger(2),
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
					parser.NewInteger(2),
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
					parser.NewInteger(2),
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
		Error: "[L:- C:-] file notexist does not exist",
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
		Error: "[L:- C:-] field notexist does not exist",
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
		Error: "[L:- C:-] field name column1 is a duplicate",
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
		Error: "[L:- C:-] field notexist does not exist",
	},
}

func TestAddColumns(t *testing.T) {
	tf := cmd.GetFlags()
	tf.Repository = TestDir

	filter := NewEmptyFilter()

	for _, v := range addColumnsTests {
		ViewCache.Clear()
		result, err := AddColumns(v.Query, filter)
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
		Error: "[L:- C:-] file notexist does not exist",
	},
	{
		Name: "Drop Columns Field Does Not Exist Error",
		Query: parser.DropColumns{
			Table: parser.Identifier{Literal: "table1"},
			Columns: []parser.Expression{
				parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			},
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
}

func TestDropColumns(t *testing.T) {
	tf := cmd.GetFlags()
	tf.Repository = TestDir

	for _, v := range dropColumnsTests {
		ViewCache.Clear()
		result, err := DropColumns(v.Query, NewEmptyFilter())
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
		Error: "[L:- C:-] file notexist does not exist",
	},
	{
		Name: "Rename Column Field Duplicate Error",
		Query: parser.RenameColumn{
			Table: parser.Identifier{Literal: "table1"},
			Old:   parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
			New:   parser.Identifier{Literal: "column1"},
		},
		Error: "[L:- C:-] field name column1 is a duplicate",
	},
	{
		Name: "Rename Column Field Does Not Exist Error",
		Query: parser.RenameColumn{
			Table: parser.Identifier{Literal: "table1"},
			Old:   parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			New:   parser.Identifier{Literal: "newcolumn"},
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
}

func TestRenameColumn(t *testing.T) {
	tf := cmd.GetFlags()
	tf.Repository = TestDir

	for _, v := range renameColumnTests {
		ViewCache.Clear()
		result, err := RenameColumn(v.Query, NewEmptyFilter())
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
