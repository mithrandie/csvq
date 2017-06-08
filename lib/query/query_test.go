package query

import (
	"path"
	"reflect"
	"testing"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/parser"
)

var executeTests = []struct {
	Input  string
	Result []Result
	Error  string
}{
	{
		Input: "var @var1; @var1 := 1; select @var1 as var1",
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
				Count: 1,
			},
		},
	},
	{
		Input: "var @var1 := 0;",
		Error: "variable @var1 is redeclared",
	},
	{
		Input: "@var2 := 0;",
		Error: "variable @var2 is undefined",
	},
	{
		Input: "select column1 from table1 where column1 = 1 group by column1 having sum(column1) > 0 order by column1 limit 10",
		Result: []Result{
			{
				Type: SELECT,
				View: &View{
					Header: []HeaderField{
						{
							Reference:  "table1",
							Column:     "column1",
							FromTable:  true,
							IsGroupKey: true,
						},
					},
					Records: []Record{
						{
							NewCell(parser.NewString("1")),
						},
					},
					FileInfo: &FileInfo{
						Path:      path.Join(TestDir, "table1.csv"),
						Delimiter: ',',
					},
				},
				Count: 1,
			},
		},
	},
	{
		Input: "select from notexist",
		Error: "syntax error: unexpected FROM",
	},
	{
		Input: "select column1 from notexist",
		Error: "file notexist does not exist",
	},
	{
		Input: "select column1 from table1 where notexist = 1",
		Error: "identifier = notexist: field does not exist",
	},
	{
		Input: "select column1 from table1 group by notexist",
		Error: "identifier = notexist: field does not exist",
	},
	{
		Input: "select column1 from table1 having notexist",
		Error: "identifier = notexist: field does not exist",
	},
	{
		Input: "select column1 from table1 order by notexist",
		Error: "identifier = notexist: field does not exist",
	},
	{
		Input: "select notexist",
		Error: "identifier = notexist: field does not exist",
	},
	{
		Input: "insert into table1 (column1, column2) values (4, 'str4'), (5, 'str5')",
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
						Path:      path.Join(TestDir, "table1.csv"),
						Delimiter: ',',
					},
					OperatedRecords: 2,
				},
				Count: 2,
			},
		},
	},
	{
		Input: "insert into table1 (column1) values (4, 'str4')",
		Error: "field length does not match value length",
	},
	{
		Input: "update table1 set column2 = 'update' where column1 = 2",
		Result: []Result{
			{
				Type: UPDATE,
				View: &View{
					FileInfo: &FileInfo{
						Path:      path.Join(TestDir, "table1.csv"),
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
				Count: 1,
			},
		},
	},
	{
		Input: "update table1 set column2 = 'update' from table1 as t1 join table2 as t2",
		Error: "file table1 is not loaded",
	},
	{
		Input: "delete from table1 where column1 = 2",
		Result: []Result{
			{
				Type: DELETE,
				View: &View{
					FileInfo: &FileInfo{
						Path:      path.Join(TestDir, "table1.csv"),
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
				Count: 1,
			},
		},
	},
	{
		Input: "delete from notexist where column1 = 2",
		Error: "file notexist does not exist",
	},
}

func TestExecute(t *testing.T) {
	Variable = map[string]parser.Primary{}

	tf := cmd.GetFlags()
	tf.Repository = TestDir

	for _, v := range executeTests {
		results, err := Execute(v.Input)
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

		if !reflect.DeepEqual(results, v.Result) {
			t.Errorf("results = %q, want %q for %q", results, v.Result, v.Input)
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
				Path:      path.Join(TestDir, "table1.csv"),
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
				Path:      path.Join(TestDir, "table1.csv"),
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
				Path:      path.Join(TestDir, "table1.csv"),
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
					Path:      path.Join(TestDir, "table1.csv"),
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
					Path:      path.Join(TestDir, "table1.csv"),
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
					Path:      path.Join(TestDir, "table1.csv"),
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
					Path:      path.Join(TestDir, "table1.csv"),
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
