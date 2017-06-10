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
		Name:       "Insert Query",
		Input:      "insert into insert_query values (4, 'str4'), (5, 'str5')",
		Output:     fmt.Sprintf("%d records inserted on %q\n", 2, GetTestFilePath("insert_query.csv")),
		UpdateFile: GetTestFilePath("insert_query.csv"),
		Content: "\"column1\",\"column2\"\n" +
			"\"1\",\"str1\"\n" +
			"\"2\",\"str2\"\n" +
			"\"3\",\"str3\"\n" +
			"4,\"str4\"\n" +
			"5,\"str5\"",
	},
	{
		Name:       "Update Query",
		Input:      "update update_query set column2 = 'update' where column1 = 2",
		Output:     fmt.Sprintf("%d record updated on %q\n", 1, GetTestFilePath("update_query.csv")),
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
		Name:       "Delete Query",
		Input:      "delete from delete_query where column1 = 2",
		Output:     fmt.Sprintf("%d record deleted on %q\n", 1, GetTestFilePath("delete_query.csv")),
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
		Name:       "Create Table",
		Input:      "create table create_table.csv (column1, column2)",
		Output:     fmt.Sprintf("file %q is created\n", GetTestFilePath("create_table.csv")),
		UpdateFile: GetTestFilePath("create_table.csv"),
		Content:    "\"column1\",\"column2\"\n",
	},
	{
		Name:       "Add Columns",
		Input:      "alter table add_columns add column3",
		Output:     fmt.Sprintf("%d field added on %q\n", 1, GetTestFilePath("add_columns.csv")),
		UpdateFile: GetTestFilePath("add_columns.csv"),
		Content: "\"column1\",\"column2\",\"column3\"\n" +
			"\"1\",\"str1\",\n" +
			"\"2\",\"str2\",\n" +
			"\"3\",\"str3\",",
	},
	{
		Name:       "Drop Columns",
		Input:      "alter table drop_columns drop column1",
		Output:     fmt.Sprintf("%d field dropped on %q\n", 1, GetTestFilePath("drop_columns.csv")),
		UpdateFile: GetTestFilePath("drop_columns.csv"),
		Content: "\"column2\"\n" +
			"\"str1\"\n" +
			"\"str2\"\n" +
			"\"str3\"",
	},
	{
		Name:       "Rename Column",
		Input:      "alter table rename_column rename column1 to newcolumn",
		Output:     fmt.Sprintf("%d field renamed on %q\n", 1, GetTestFilePath("rename_column.csv")),
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
		Error: "syntax error: unexpected FROM",
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
			SelectClause: parser.SelectClause{
				Fields: []parser.Expression{
					parser.Field{
						Object: parser.Variable{Name: "@var1"},
						Alias:  parser.Identifier{Literal: "var1"},
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
							Alias: "var1",
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
				parser.Identifier{Literal: "column1"},
				parser.Identifier{Literal: "column2"},
			},
			ValuesList: []parser.Expression{
				parser.InsertValues{
					Values: []parser.Expression{
						parser.NewInteger(4),
						parser.NewString("str4"),
					},
				},
				parser.InsertValues{
					Values: []parser.Expression{
						parser.NewInteger(5),
						parser.NewString("str5"),
					},
				},
			},
		},
		Result: []Result{
			{
				Type: INSERT,
				View: &View{
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
					FileInfo: &FileInfo{
						Path:      GetTestFilePath("table1.csv"),
						Delimiter: ',',
					},
					OperatedRecords: 2,
				},
				Log: fmt.Sprintf("2 records inserted on %q", GetTestFilePath("table1.csv")),
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
					Field: parser.Identifier{Literal: "column2"},
					Value: parser.NewString("update"),
				},
			},
			WhereClause: parser.WhereClause{
				Filter: parser.Comparison{
					LHS:      parser.Identifier{Literal: "column1"},
					RHS:      parser.NewInteger(2),
					Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: "="},
				},
			},
		},
		Result: []Result{
			{
				Type: UPDATE,
				View: &View{
					FileInfo: &FileInfo{
						Path:      GetTestFilePath("table1.csv"),
						Delimiter: ',',
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
				Log: fmt.Sprintf("1 record updated on %q", GetTestFilePath("table1.csv")),
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
					LHS:      parser.Identifier{Literal: "column1"},
					RHS:      parser.NewInteger(2),
					Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: "="},
				},
			},
		},
		Result: []Result{
			{
				Type: DELETE,
				View: &View{
					FileInfo: &FileInfo{
						Path:      GetTestFilePath("table1.csv"),
						Delimiter: ',',
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
				Log: fmt.Sprintf("1 record deleted on %q", GetTestFilePath("table1.csv")),
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
				View: &View{
					FileInfo: &FileInfo{
						Path:      GetTestFilePath("newtable.csv"),
						Delimiter: ',',
					},
					Header: NewHeaderWithoutId("newtable", []string{"column1", "column2"}),
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
				View: &View{
					FileInfo: &FileInfo{
						Path:      GetTestFilePath("table1.csv"),
						Delimiter: ',',
					},
					Header: NewHeaderWithoutId("table1", []string{"column1", "column2", "column3"}),
					Records: []Record{
						NewRecordWithoutId([]parser.Primary{
							parser.NewString("1"),
							parser.NewString("str1"),
							parser.NewNull(),
						}),
						NewRecordWithoutId([]parser.Primary{
							parser.NewString("2"),
							parser.NewString("str2"),
							parser.NewNull(),
						}),
						NewRecordWithoutId([]parser.Primary{
							parser.NewString("3"),
							parser.NewString("str3"),
							parser.NewNull(),
						}),
					},
					OperatedFields: 1,
				},
				Log: fmt.Sprintf("1 field added on %q", GetTestFilePath("table1.csv")),
			},
		},
	},
	{
		Input: parser.DropColumns{
			Table: parser.Identifier{Literal: "table1"},
			Columns: []parser.Expression{
				parser.Identifier{Literal: "column1"},
			},
		},
		Result: []Result{
			{
				Type: DROP_COLUMNS,
				View: &View{
					FileInfo: &FileInfo{
						Path:      GetTestFilePath("table1.csv"),
						Delimiter: ',',
					},
					Header: NewHeaderWithoutId("table1", []string{"column2"}),
					Records: []Record{
						NewRecordWithoutId([]parser.Primary{
							parser.NewString("str1"),
						}),
						NewRecordWithoutId([]parser.Primary{
							parser.NewString("str2"),
						}),
						NewRecordWithoutId([]parser.Primary{
							parser.NewString("str3"),
						}),
					},
					OperatedFields: 1,
				},
				Log: fmt.Sprintf("1 field dropped on %q", GetTestFilePath("table1.csv")),
			},
		},
	},
	{
		Input: parser.RenameColumn{
			Table: parser.Identifier{Literal: "table1"},
			Old:   parser.Identifier{Literal: "column1"},
			New:   parser.Identifier{Literal: "newcolumn"},
		},
		Result: []Result{
			{
				Type: RENAME_COLUMN,
				View: &View{
					FileInfo: &FileInfo{
						Path:      GetTestFilePath("table1.csv"),
						Delimiter: ',',
					},
					Header: NewHeaderWithoutId("table1", []string{"newcolumn", "column2"}),
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
				Log: fmt.Sprintf("1 field renamed on %q", GetTestFilePath("table1.csv")),
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
	Variable = map[string]parser.Primary{}

	tf := cmd.GetFlags()
	tf.Repository = TestDir

	for _, v := range executeStatementTests {
		ViewCache.Clear()
		ResultSet = []Result{}

		_, err := ExecuteStatement(v.Input)
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
				parser.Identifier{Literal: "column1"},
			},
			Values: "values",
			ValuesList: []parser.Expression{
				parser.InsertValues{
					Values: []parser.Expression{
						parser.NewInteger(4),
					},
				},
				parser.InsertValues{
					Values: []parser.Expression{
						parser.NewInteger(5),
					},
				},
			},
		},
		Result: &View{
			FileInfo: &FileInfo{
				Path:      GetTestFilePath("table1.csv"),
				Delimiter: ',',
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
				parser.InsertValues{
					Values: []parser.Expression{
						parser.NewInteger(4),
						parser.NewString("str4"),
					},
				},
				parser.InsertValues{
					Values: []parser.Expression{
						parser.NewInteger(5),
						parser.NewString("str5"),
					},
				},
			},
		},
		Result: &View{
			FileInfo: &FileInfo{
				Path:      GetTestFilePath("table1.csv"),
				Delimiter: ',',
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
				parser.Identifier{Literal: "column1"},
			},
			Values: "values",
			ValuesList: []parser.Expression{
				parser.InsertValues{
					Values: []parser.Expression{
						parser.NewInteger(4),
					},
				},
				parser.InsertValues{
					Values: []parser.Expression{
						parser.NewInteger(5),
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
				parser.Identifier{Literal: "notexist"},
			},
			Values: "values",
			ValuesList: []parser.Expression{
				parser.InsertValues{
					Values: []parser.Expression{
						parser.NewInteger(4),
					},
				},
				parser.InsertValues{
					Values: []parser.Expression{
						parser.NewInteger(5),
					},
				},
			},
		},
		Error: "identifier = notexist: field does not exist",
	},
	{
		Name: "Insert Select Query",
		Query: parser.InsertQuery{
			Insert: "insert",
			Into:   "into",
			Table:  parser.Identifier{Literal: "table1"},
			Fields: []parser.Expression{
				parser.Identifier{Literal: "column1"},
				parser.Identifier{Literal: "column2"},
			},
			Query: parser.SelectQuery{
				SelectClause: parser.SelectClause{
					Fields: []parser.Expression{
						parser.Field{Object: parser.Identifier{Literal: "column3"}},
						parser.Field{Object: parser.Identifier{Literal: "column4"}},
					},
				},
				FromClause: parser.FromClause{
					Tables: []parser.Expression{
						parser.Table{Object: parser.Identifier{Literal: "table2"}},
					},
				},
			},
		},
		Result: &View{
			FileInfo: &FileInfo{
				Path:      GetTestFilePath("table1.csv"),
				Delimiter: ',',
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
		Name: "Insert Query Field Does Not Exist Error",
		Query: parser.InsertQuery{
			Insert: "insert",
			Into:   "into",
			Table:  parser.Identifier{Literal: "table1"},
			Fields: []parser.Expression{
				parser.Identifier{Literal: "notexist"},
			},
			Values: "values",
			ValuesList: []parser.Expression{
				parser.InsertValues{
					Values: []parser.Expression{
						parser.NewInteger(4),
					},
				},
				parser.InsertValues{
					Values: []parser.Expression{
						parser.NewInteger(5),
					},
				},
			},
		},
		Error: "identifier = notexist: field does not exist",
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
					Field: parser.Identifier{Literal: "column2"},
					Value: parser.NewString("update"),
				},
			},
			WhereClause: parser.WhereClause{
				Filter: parser.Comparison{
					LHS:      parser.Identifier{Literal: "column1"},
					RHS:      parser.NewInteger(2),
					Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: "="},
				},
			},
		},
		Result: []*View{
			{
				FileInfo: &FileInfo{
					Path:      GetTestFilePath("table1.csv"),
					Delimiter: ',',
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
					Field: parser.Identifier{Literal: "column2"},
					Value: parser.Identifier{Literal: "column4"},
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
								LHS:      parser.Identifier{Literal: "column1"},
								RHS:      parser.Identifier{Literal: "column3"},
								Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: "="},
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
					Field: parser.Identifier{Literal: "column2"},
					Value: parser.NewString("update"),
				},
			},
			WhereClause: parser.WhereClause{
				Filter: parser.Comparison{
					LHS:      parser.Identifier{Literal: "column1"},
					RHS:      parser.NewInteger(2),
					Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: "="},
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
					Field: parser.Identifier{Literal: "column1"},
					Value: parser.NewString("update"),
				},
			},
			WhereClause: parser.WhereClause{
				Filter: parser.Comparison{
					LHS:      parser.Identifier{Literal: "notexist"},
					RHS:      parser.NewInteger(2),
					Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: "="},
				},
			},
		},
		Error: "identifier = notexist: field does not exist",
	},
	{
		Name: "Update Query File Is Not Loaded Error",
		Query: parser.UpdateQuery{
			Update: "update",
			Tables: []parser.Expression{
				parser.Table{Object: parser.Identifier{Literal: "table1"}},
			},
			Set: "set",
			SetList: []parser.Expression{
				parser.UpdateSet{
					Field: parser.Identifier{Literal: "column2"},
					Value: parser.Identifier{Literal: "column4"},
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
								LHS:      parser.Identifier{Literal: "column1"},
								RHS:      parser.Identifier{Literal: "column3"},
								Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: "="},
							},
						},
					}},
				},
			},
		},
		Error: "file table1 is not loaded",
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
					Field: parser.Identifier{Literal: "notexist"},
					Value: parser.NewString("update"),
				},
			},
			WhereClause: parser.WhereClause{
				Filter: parser.Comparison{
					LHS:      parser.Identifier{Literal: "column1"},
					RHS:      parser.NewInteger(2),
					Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: "="},
				},
			},
		},
		Error: "identifier = notexist: field does not exist",
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
					Field: parser.Identifier{Literal: "column1"},
					Value: parser.Identifier{Literal: "notexist"},
				},
			},
			WhereClause: parser.WhereClause{
				Filter: parser.Comparison{
					LHS:      parser.Identifier{Literal: "column1"},
					RHS:      parser.NewInteger(2),
					Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: "="},
				},
			},
		},
		Error: "identifier = notexist: field does not exist",
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
					Field: parser.Identifier{Literal: "column2"},
					Value: parser.Identifier{Literal: "column4"},
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
					LHS:      parser.Identifier{Literal: "column1"},
					RHS:      parser.NewInteger(2),
					Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: "="},
				},
			},
		},
		Result: []*View{
			{
				FileInfo: &FileInfo{
					Path:      GetTestFilePath("table1.csv"),
					Delimiter: ',',
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
								LHS:      parser.Identifier{Literal: "column1"},
								RHS:      parser.Identifier{Literal: "column3"},
								Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: "="},
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
								LHS:      parser.Identifier{Literal: "column1"},
								RHS:      parser.Identifier{Literal: "column3"},
								Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: "="},
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
					LHS:      parser.Identifier{Literal: "column1"},
					RHS:      parser.NewInteger(2),
					Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: "="},
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
					LHS:      parser.Identifier{Literal: "column1"},
					RHS:      parser.Identifier{Literal: "notexist"},
					Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: "="},
				},
			},
		},
		Error: "identifier = notexist: field does not exist",
	},
	{
		Name: "Delete Query File Is Not Loaded Error",
		Query: parser.DeleteQuery{
			Delete: "delete",
			Tables: []parser.Expression{
				parser.Table{Object: parser.Identifier{Literal: "table1"}},
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
								LHS:      parser.Identifier{Literal: "column1"},
								RHS:      parser.Identifier{Literal: "column3"},
								Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: "="},
							},
						},
					}},
				},
			},
		},
		Error: "file table1 is not loaded",
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
				Column:   parser.Identifier{Literal: "column1"},
			},
		},
		Result: &View{
			FileInfo: &FileInfo{
				Path:      GetTestFilePath("table1.csv"),
				Delimiter: ',',
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
				Column:   parser.Identifier{Literal: "column2"},
			},
		},
		Result: &View{
			FileInfo: &FileInfo{
				Path:      GetTestFilePath("table1.csv"),
				Delimiter: ',',
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
				Column:   parser.Identifier{Literal: "notexist"},
			},
		},
		Error: "identifier = notexist: field does not exist",
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
					Value:  parser.Identifier{Literal: "notexist.column1"},
				},
			},
		},
		Error: "identifier = notexist.column1: field does not exist",
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
				parser.Identifier{Literal: "column2"},
			},
		},
		Result: &View{
			FileInfo: &FileInfo{
				Path:      GetTestFilePath("table1.csv"),
				Delimiter: ',',
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
				parser.Identifier{Literal: "column2"},
			},
		},
		Error: "file notexist does not exist",
	},
	{
		Name: "Drop Columns Field Does Not Exist Error",
		Query: parser.DropColumns{
			Table: parser.Identifier{Literal: "table1"},
			Columns: []parser.Expression{
				parser.Identifier{Literal: "notexist"},
			},
		},
		Error: "identifier = notexist: field does not exist",
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
			Old:   parser.Identifier{Literal: "column2"},
			New:   parser.Identifier{Literal: "newcolumn"},
		},
		Result: &View{
			FileInfo: &FileInfo{
				Path:      GetTestFilePath("table1.csv"),
				Delimiter: ',',
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
			Old:   parser.Identifier{Literal: "column2"},
			New:   parser.Identifier{Literal: "newcolumn"},
		},
		Error: "file notexist does not exist",
	},
	{
		Name: "Rename Column Field Duplicate Error",
		Query: parser.RenameColumn{
			Table: parser.Identifier{Literal: "table1"},
			Old:   parser.Identifier{Literal: "column2"},
			New:   parser.Identifier{Literal: "column1"},
		},
		Error: "field column1 is duplicate",
	},
	{
		Name: "Rename Column Field Does Not Exist Error",
		Query: parser.RenameColumn{
			Table: parser.Identifier{Literal: "table1"},
			Old:   parser.Identifier{Literal: "notexist"},
			New:   parser.Identifier{Literal: "newcolumn"},
		},
		Error: "identifier = notexist: field does not exist",
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
		if !reflect.DeepEqual(result, v.Result) {
			t.Errorf("%s: result = %q, want %q", v.Name, result, v.Result)
		}
	}
}
