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
