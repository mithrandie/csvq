package query

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/file"
	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/value"

	"github.com/mithrandie/go-text"
	"github.com/mithrandie/go-text/json"
)

var fetchCursorTests = []struct {
	Name          string
	CurName       parser.Identifier
	FetchPosition parser.FetchPosition
	Variables     []parser.Variable
	Success       bool
	ResultVars    VariableMap
	Error         string
}{
	{
		Name:    "Fetch Cursor First Time",
		CurName: parser.Identifier{Literal: "cur"},
		Variables: []parser.Variable{
			{Name: "var1"},
			{Name: "var2"},
		},
		Success: true,
		ResultVars: VariableMap{
			"var1": value.NewString("1"),
			"var2": value.NewString("str1"),
		},
	},
	{
		Name:    "Fetch Cursor Second Time",
		CurName: parser.Identifier{Literal: "cur"},
		Variables: []parser.Variable{
			{Name: "var1"},
			{Name: "var2"},
		},
		Success: true,
		ResultVars: VariableMap{
			"var1": value.NewString("2"),
			"var2": value.NewString("str2"),
		},
	},
	{
		Name:    "Fetch Cursor Third Time",
		CurName: parser.Identifier{Literal: "cur"},
		Variables: []parser.Variable{
			{Name: "var1"},
			{Name: "var2"},
		},
		Success: true,
		ResultVars: VariableMap{
			"var1": value.NewString("3"),
			"var2": value.NewString("str3"),
		},
	},
	{
		Name:    "Fetch Cursor Forth Time",
		CurName: parser.Identifier{Literal: "cur"},
		Variables: []parser.Variable{
			{Name: "var1"},
			{Name: "var2"},
		},
		Success: false,
		ResultVars: VariableMap{
			"var1": value.NewString("3"),
			"var2": value.NewString("str3"),
		},
	},
	{
		Name:    "Fetch Cursor Absolute",
		CurName: parser.Identifier{Literal: "cur"},
		FetchPosition: parser.FetchPosition{
			Position: parser.Token{Token: parser.ABSOLUTE, Literal: "absolute"},
			Number:   parser.NewIntegerValueFromString("1"),
		},
		Variables: []parser.Variable{
			{Name: "var1"},
			{Name: "var2"},
		},
		Success: true,
		ResultVars: VariableMap{
			"var1": value.NewString("2"),
			"var2": value.NewString("str2"),
		},
	},
	{
		Name:    "Fetch Cursor Fetch Error",
		CurName: parser.Identifier{Literal: "notexist"},
		Variables: []parser.Variable{
			{Name: "var1"},
			{Name: "var2"},
		},
		Error: "[L:- C:-] cursor notexist is undeclared",
	},
	{
		Name:    "Fetch Cursor Not Match Number Error",
		CurName: parser.Identifier{Literal: "cur2"},
		Variables: []parser.Variable{
			{Name: "var1"},
		},
		Error: "[L:- C:-] fetching from cursor cur2 returns 2 values",
	},
	{
		Name:    "Fetch Cursor Substitution Error",
		CurName: parser.Identifier{Literal: "cur2"},
		Variables: []parser.Variable{
			{Name: "var1"},
			{Name: "notexist"},
		},
		Error: "[L:- C:-] variable @notexist is undeclared",
	},
	{
		Name:    "Fetch Cursor Number Value Error",
		CurName: parser.Identifier{Literal: "cur"},
		FetchPosition: parser.FetchPosition{
			Position: parser.Token{Token: parser.ABSOLUTE, Literal: "absolute"},
			Number:   parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
		},
		Variables: []parser.Variable{
			{Name: "var1"},
			{Name: "var2"},
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
	{
		Name:    "Fetch Cursor Number Not Integer Error",
		CurName: parser.Identifier{Literal: "cur"},
		FetchPosition: parser.FetchPosition{
			Position: parser.Token{Token: parser.ABSOLUTE, Literal: "absolute"},
			Number:   parser.NewNullValueFromString("null"),
		},
		Variables: []parser.Variable{
			{Name: "var1"},
			{Name: "var2"},
		},
		Error: "[L:- C:-] fetching position null is not an integer value",
	},
}

func TestFetchCursor(t *testing.T) {
	tf := cmd.GetFlags()
	tf.Repository = TestDir

	filter := NewFilter(
		[]VariableMap{
			{
				"var1": value.NewNull(),
				"var2": value.NewNull(),
			},
		},
		[]ViewMap{{}},
		[]CursorMap{
			{
				"CUR": &Cursor{
					query: selectQueryForCursorTest,
				},
				"CUR2": &Cursor{
					query: selectQueryForCursorTest,
				},
			},
		},
		[]UserDefinedFunctionMap{{}},
	)

	ViewCache.Clean()
	filter.Cursors.Open(parser.Identifier{Literal: "cur"}, filter)
	ViewCache.Clean()
	filter.Cursors.Open(parser.Identifier{Literal: "cur2"}, filter)

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
		if !reflect.DeepEqual(filter.Variables[0], v.ResultVars) {
			t.Errorf("%s: global vars = %q, want %q", v.Name, filter.Variables[0], v.ResultVars)
		}
	}
}

var declareViewTests = []struct {
	Name    string
	ViewMap ViewMap
	Expr    parser.ViewDeclaration
	Result  ViewMap
	Error   string
}{
	{
		Name: "Declare View",
		Expr: parser.ViewDeclaration{
			View: parser.Identifier{Literal: "tbl"},
			Fields: []parser.QueryExpression{
				parser.Identifier{Literal: "column1"},
				parser.Identifier{Literal: "column2"},
			},
		},
		Result: ViewMap{
			"TBL": {
				FileInfo: &FileInfo{
					Path:             "tbl",
					IsTemporary:      true,
					InitialHeader:    NewHeader("tbl", []string{"column1", "column2"}),
					InitialRecordSet: RecordSet{},
				},
				Header:    NewHeader("tbl", []string{"column1", "column2"}),
				RecordSet: RecordSet{},
			},
		},
	},
	{
		Name: "Declare View Field Duplicate Error",
		Expr: parser.ViewDeclaration{
			View: parser.Identifier{Literal: "tbl"},
			Fields: []parser.QueryExpression{
				parser.Identifier{Literal: "column1"},
				parser.Identifier{Literal: "column1"},
			},
		},
		Error: "[L:- C:-] field name column1 is a duplicate",
	},
	{
		Name: "Declare View From Query",
		Expr: parser.ViewDeclaration{
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
		Result: ViewMap{
			"TBL": {
				FileInfo: &FileInfo{
					Path:          "tbl",
					IsTemporary:   true,
					InitialHeader: NewHeader("tbl", []string{"column1", "column2"}),
					InitialRecordSet: RecordSet{
						NewRecord([]value.Primary{
							value.NewInteger(1),
							value.NewInteger(2),
						}),
					},
				},
				Header: NewHeader("tbl", []string{"column1", "column2"}),
				RecordSet: RecordSet{
					NewRecord([]value.Primary{
						value.NewInteger(1),
						value.NewInteger(2),
					}),
				},
			},
		},
	},
	{
		Name: "Declare View From Query Query Error",
		Expr: parser.ViewDeclaration{
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
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}}},
						},
					},
				},
			},
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
	{
		Name: "Declare View From Query Field Update Error",
		Expr: parser.ViewDeclaration{
			View: parser.Identifier{Literal: "tbl"},
			Fields: []parser.QueryExpression{
				parser.Identifier{Literal: "column1"},
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
		Error: "[L:- C:-] select query should return exactly 1 field for view tbl",
	},
	{
		Name: "Declare View  From Query Field Duplicate Error",
		Expr: parser.ViewDeclaration{
			View: parser.Identifier{Literal: "tbl"},
			Fields: []parser.QueryExpression{
				parser.Identifier{Literal: "column1"},
				parser.Identifier{Literal: "column1"},
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
		Error: "[L:- C:-] field name column1 is a duplicate",
	},
	{
		Name: "Declare View Redeclaration Error",
		ViewMap: ViewMap{
			"TBL": {
				FileInfo: &FileInfo{
					Path:        "tbl",
					IsTemporary: true,
				},
			},
		},
		Expr: parser.ViewDeclaration{
			View: parser.Identifier{Literal: "tbl"},
			Fields: []parser.QueryExpression{
				parser.Identifier{Literal: "column1"},
				parser.Identifier{Literal: "column2"},
			},
		},
		Error: "[L:- C:-] view tbl is redeclared",
	},
}

func TestDeclareView(t *testing.T) {
	filter := NewEmptyFilter()

	for _, v := range declareViewTests {
		if v.ViewMap == nil {
			filter.TempViews = []ViewMap{{}}
		} else {
			filter.TempViews = []ViewMap{v.ViewMap}
		}

		err := DeclareView(v.Expr, filter)
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
		if !reflect.DeepEqual(filter.TempViews[0], v.Result) {
			t.Errorf("%s: view cache = %v, want %v", v.Name, ViewCache, v.Result)
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
					Fields: []parser.QueryExpression{
						parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
						parser.Field{Object: parser.AggregateFunction{Name: "count", Args: []parser.QueryExpression{parser.AllColumns{}}}},
					},
				},
				FromClause: parser.FromClause{
					Tables: []parser.QueryExpression{
						parser.Table{Object: parser.Identifier{Literal: "group_table"}},
					},
				},
				WhereClause: parser.WhereClause{
					Filter: parser.Comparison{
						LHS:      parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
						RHS:      parser.NewIntegerValueFromString("3"),
						Operator: "<",
					},
				},
				GroupByClause: parser.GroupByClause{
					Items: []parser.QueryExpression{
						parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
					},
				},
				HavingClause: parser.HavingClause{
					Filter: parser.Comparison{
						LHS:      parser.AggregateFunction{Name: "count", Args: []parser.QueryExpression{parser.AllColumns{}}},
						RHS:      parser.NewIntegerValueFromString("1"),
						Operator: ">",
					},
				},
			},
			OrderByClause: parser.OrderByClause{
				Items: []parser.QueryExpression{
					parser.OrderItem{Value: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
				},
			},
			LimitClause: parser.LimitClause{
				Value: parser.NewIntegerValueFromString("5"),
			},
			OffsetClause: parser.OffsetClause{
				Value: parser.NewIntegerValue(0),
			},
		},
		Result: &View{
			FileInfo: &FileInfo{
				Path:      GetTestFilePath("group_table.csv"),
				Delimiter: ',',
				NoHeader:  false,
				Encoding:  text.UTF8,
				LineBreak: text.LF,
			},
			Header: []HeaderField{
				{
					View:        "group_table",
					Column:      "column1",
					Number:      1,
					IsFromTable: true,
				},
				{
					Column:      "count(*)",
					Number:      2,
					IsFromTable: true,
				},
			},
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("1"),
					value.NewInteger(2),
				}),
				NewRecord([]value.Primary{
					value.NewString("2"),
					value.NewInteger(2),
				}),
			},
		},
	},
	{
		Name: "Select Replace Fields",
		Query: parser.SelectQuery{
			SelectEntity: parser.SelectEntity{
				SelectClause: parser.SelectClause{
					Fields: []parser.QueryExpression{
						parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
						parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
					},
				},
				FromClause: parser.FromClause{
					Tables: []parser.QueryExpression{
						parser.Table{Object: parser.Identifier{Literal: "table1"}},
					},
				},
			},
			LimitClause: parser.LimitClause{
				Value: parser.NewIntegerValueFromString("1"),
			},
		},
		Result: &View{
			FileInfo: &FileInfo{
				Path:      GetTestFilePath("table1.csv"),
				Delimiter: ',',
				NoHeader:  false,
				Encoding:  text.UTF8,
				LineBreak: text.LF,
			},
			Header: NewHeader("table1", []string{"column2", "column1"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("str1"),
					value.NewString("1"),
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
						Fields: []parser.QueryExpression{
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
						},
					},
					FromClause: parser.FromClause{
						Tables: []parser.QueryExpression{
							parser.Table{Object: parser.Identifier{Literal: "table1"}},
						},
					},
				},
				Operator: parser.Token{Token: parser.UNION, Literal: "union"},
				RHS: parser.SelectEntity{
					SelectClause: parser.SelectClause{
						Fields: []parser.QueryExpression{
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column3"}}},
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column4"}}},
						},
					},
					FromClause: parser.FromClause{
						Tables: []parser.QueryExpression{
							parser.Table{Object: parser.Identifier{Literal: "table4"}},
						},
					},
				},
			},
		},
		Result: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecord([]value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
				}),
				NewRecord([]value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
				}),
				NewRecord([]value.Primary{
					value.NewString("4"),
					value.NewString("str4"),
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
						Fields: []parser.QueryExpression{
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
						},
					},
					FromClause: parser.FromClause{
						Tables: []parser.QueryExpression{
							parser.Table{Object: parser.Identifier{Literal: "table1"}},
						},
					},
				},
				Operator: parser.Token{Token: parser.INTERSECT, Literal: "intersect"},
				RHS: parser.SelectEntity{
					SelectClause: parser.SelectClause{
						Fields: []parser.QueryExpression{
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column3"}}},
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column4"}}},
						},
					},
					FromClause: parser.FromClause{
						Tables: []parser.QueryExpression{
							parser.Table{Object: parser.Identifier{Literal: "table4"}},
						},
					},
				},
			},
		},
		Result: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
				}),
				NewRecord([]value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
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
						Fields: []parser.QueryExpression{
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
						},
					},
					FromClause: parser.FromClause{
						Tables: []parser.QueryExpression{
							parser.Table{Object: parser.Identifier{Literal: "table1"}},
						},
					},
				},
				Operator: parser.Token{Token: parser.EXCEPT, Literal: "except"},
				RHS: parser.SelectEntity{
					SelectClause: parser.SelectClause{
						Fields: []parser.QueryExpression{
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column3"}}},
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column4"}}},
						},
					},
					FromClause: parser.FromClause{
						Tables: []parser.QueryExpression{
							parser.Table{Object: parser.Identifier{Literal: "table4"}},
						},
					},
				},
			},
		},
		Result: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
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
								Fields: []parser.QueryExpression{
									parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
									parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
								},
							},
							FromClause: parser.FromClause{
								Tables: []parser.QueryExpression{
									parser.Table{Object: parser.Identifier{Literal: "table1"}},
								},
							},
						},
					},
				},
				Operator: parser.Token{Token: parser.UNION, Literal: "union"},
				RHS: parser.SelectEntity{
					SelectClause: parser.SelectClause{
						Fields: []parser.QueryExpression{
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column3"}}},
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column4"}}},
						},
					},
					FromClause: parser.FromClause{
						Tables: []parser.QueryExpression{
							parser.Table{Object: parser.Identifier{Literal: "table4"}},
						},
					},
				},
			},
		},
		Result: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecord([]value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
				}),
				NewRecord([]value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
				}),
				NewRecord([]value.Primary{
					value.NewString("4"),
					value.NewString("str4"),
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
						Fields: []parser.QueryExpression{
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
						},
					},
					FromClause: parser.FromClause{
						Tables: []parser.QueryExpression{
							parser.Table{Object: parser.Identifier{Literal: "table1"}},
						},
					},
				},
				Operator: parser.Token{Token: parser.UNION, Literal: "union"},
				RHS: parser.SelectEntity{
					SelectClause: parser.SelectClause{
						Fields: []parser.QueryExpression{
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column3"}}},
						},
					},
					FromClause: parser.FromClause{
						Tables: []parser.QueryExpression{
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
						Fields: []parser.QueryExpression{
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}}},
						},
					},
					FromClause: parser.FromClause{
						Tables: []parser.QueryExpression{
							parser.Table{Object: parser.Identifier{Literal: "table1"}},
						},
					},
				},
				Operator: parser.Token{Token: parser.UNION, Literal: "union"},
				RHS: parser.SelectEntity{
					SelectClause: parser.SelectClause{
						Fields: []parser.QueryExpression{
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column3"}}},
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column4"}}},
						},
					},
					FromClause: parser.FromClause{
						Tables: []parser.QueryExpression{
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
						Fields: []parser.QueryExpression{
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
						},
					},
					FromClause: parser.FromClause{
						Tables: []parser.QueryExpression{
							parser.Table{Object: parser.Identifier{Literal: "table1"}},
						},
					},
				},
				Operator: parser.Token{Token: parser.UNION, Literal: "union"},
				RHS: parser.SelectEntity{
					SelectClause: parser.SelectClause{
						Fields: []parser.QueryExpression{
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column3"}}},
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}}},
						},
					},
					FromClause: parser.FromClause{
						Tables: []parser.QueryExpression{
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
				InlineTables: []parser.QueryExpression{
					parser.InlineTable{
						Name: parser.Identifier{Literal: "it"},
						Fields: []parser.QueryExpression{
							parser.Identifier{Literal: "c1"},
						},
						As: "as",
						Query: parser.SelectQuery{
							SelectEntity: parser.SelectEntity{
								SelectClause: parser.SelectClause{
									Select: "select",
									Fields: []parser.QueryExpression{
										parser.Field{Object: parser.NewIntegerValueFromString("2")},
									},
								},
							},
						},
					},
				},
			},
			SelectEntity: parser.SelectEntity{
				SelectClause: parser.SelectClause{
					Fields: []parser.QueryExpression{
						parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "c1"}}},
					},
				},
				FromClause: parser.FromClause{
					Tables: []parser.QueryExpression{
						parser.Table{Object: parser.Identifier{Literal: "it"}},
					},
				},
			},
		},
		Result: &View{
			Header: NewHeader("it", []string{"c1"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewInteger(2),
				}),
			},
		},
	},
	{
		Name: "Inline Tables Field Length Error",
		Query: parser.SelectQuery{
			WithClause: parser.WithClause{
				With: "with",
				InlineTables: []parser.QueryExpression{
					parser.InlineTable{
						Name: parser.Identifier{Literal: "it"},
						Fields: []parser.QueryExpression{
							parser.Identifier{Literal: "c1"},
						},
						As: "as",
						Query: parser.SelectQuery{
							SelectEntity: parser.SelectSet{
								LHS: parser.SelectEntity{
									SelectClause: parser.SelectClause{
										Fields: []parser.QueryExpression{
											parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
											parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
										},
									},
									FromClause: parser.FromClause{
										Tables: []parser.QueryExpression{
											parser.Table{Object: parser.Identifier{Literal: "table1"}},
										},
									},
								},
								Operator: parser.Token{Token: parser.UNION, Literal: "union"},
								RHS: parser.SelectEntity{
									SelectClause: parser.SelectClause{
										Fields: []parser.QueryExpression{
											parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column3"}}},
											parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column4"}}},
										},
									},
									FromClause: parser.FromClause{
										Tables: []parser.QueryExpression{
											parser.Table{Object: parser.Identifier{Literal: "table4"}},
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
					Fields: []parser.QueryExpression{
						parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "c1"}}},
					},
				},
				FromClause: parser.FromClause{
					Tables: []parser.QueryExpression{
						parser.Table{Object: parser.Identifier{Literal: "it"}},
					},
				},
			},
		},
		Error: "[L:- C:-] select query should return exactly 1 field for inline table it",
	},
	{
		Name: "Inline Tables Recursion",
		Query: parser.SelectQuery{
			WithClause: parser.WithClause{
				With: "with",
				InlineTables: []parser.QueryExpression{
					parser.InlineTable{
						Recursive: parser.Token{Token: parser.RECURSIVE, Literal: "recursive"},
						Name:      parser.Identifier{Literal: "it"},
						Fields: []parser.QueryExpression{
							parser.Identifier{Literal: "n"},
						},
						As: "as",
						Query: parser.SelectQuery{
							SelectEntity: parser.SelectSet{
								LHS: parser.SelectEntity{
									SelectClause: parser.SelectClause{
										Select: "select",
										Fields: []parser.QueryExpression{
											parser.Field{Object: parser.NewIntegerValueFromString("1")},
										},
									},
								},
								Operator: parser.Token{Token: parser.UNION, Literal: "union"},
								RHS: parser.SelectEntity{
									SelectClause: parser.SelectClause{
										Select: "select",
										Fields: []parser.QueryExpression{
											parser.Field{
												Object: parser.Arithmetic{
													LHS:      parser.FieldReference{Column: parser.Identifier{Literal: "n"}},
													RHS:      parser.NewIntegerValueFromString("1"),
													Operator: '+',
												},
											},
										},
									},
									FromClause: parser.FromClause{
										Tables: []parser.QueryExpression{
											parser.Table{Object: parser.Identifier{Literal: "it"}},
										},
									},
									WhereClause: parser.WhereClause{
										Filter: parser.Comparison{
											LHS:      parser.FieldReference{Column: parser.Identifier{Literal: "n"}},
											RHS:      parser.NewIntegerValueFromString("3"),
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
					Fields: []parser.QueryExpression{
						parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "n"}}},
					},
				},
				FromClause: parser.FromClause{
					Tables: []parser.QueryExpression{
						parser.Table{Object: parser.Identifier{Literal: "it"}},
					},
				},
			},
		},
		Result: &View{
			Header: []HeaderField{
				{
					View:        "it",
					Column:      "n",
					Number:      1,
					IsFromTable: true,
				},
			},
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewInteger(1),
				}),
				NewRecord([]value.Primary{
					value.NewInteger(2),
				}),
				NewRecord([]value.Primary{
					value.NewInteger(3),
				}),
			},
		},
	},
	{
		Name: "Inline Tables Recursion Field Length Error",
		Query: parser.SelectQuery{
			WithClause: parser.WithClause{
				With: "with",
				InlineTables: []parser.QueryExpression{
					parser.InlineTable{
						Recursive: parser.Token{Token: parser.RECURSIVE, Literal: "recursive"},
						Name:      parser.Identifier{Literal: "it"},
						Fields: []parser.QueryExpression{
							parser.Identifier{Literal: "n"},
						},
						As: "as",
						Query: parser.SelectQuery{
							SelectEntity: parser.SelectSet{
								LHS: parser.SelectEntity{
									SelectClause: parser.SelectClause{
										Select: "select",
										Fields: []parser.QueryExpression{
											parser.Field{Object: parser.NewIntegerValueFromString("1")},
										},
									},
								},
								Operator: parser.Token{Token: parser.UNION, Literal: "union"},
								RHS: parser.SelectEntity{
									SelectClause: parser.SelectClause{
										Select: "select",
										Fields: []parser.QueryExpression{
											parser.Field{
												Object: parser.Arithmetic{
													LHS:      parser.FieldReference{Column: parser.Identifier{Literal: "n"}},
													RHS:      parser.NewIntegerValueFromString("1"),
													Operator: '+',
												},
											},
											parser.Field{Object: parser.NewIntegerValueFromString("2")},
										},
									},
									FromClause: parser.FromClause{
										Tables: []parser.QueryExpression{
											parser.Table{Object: parser.Identifier{Literal: "it"}},
										},
									},
									WhereClause: parser.WhereClause{
										Filter: parser.Comparison{
											LHS:      parser.FieldReference{Column: parser.Identifier{Literal: "n"}},
											RHS:      parser.NewIntegerValueFromString("3"),
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
					Fields: []parser.QueryExpression{
						parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "n"}}},
					},
				},
				FromClause: parser.FromClause{
					Tables: []parser.QueryExpression{
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
		ViewCache.Clean()
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
			t.Errorf("%s: result = %v, want %v", v.Name, result, v.Result)
		}
	}
}

var insertTests = []struct {
	Name         string
	Query        parser.InsertQuery
	ResultFile   *FileInfo
	UpdateCount  int
	ViewCache    ViewMap
	TempViewList TemporaryViewScopes
	Error        string
}{
	{
		Name: "Insert Query",
		Query: parser.InsertQuery{
			WithClause: parser.WithClause{
				With: "with",
				InlineTables: []parser.QueryExpression{
					parser.InlineTable{
						Name: parser.Identifier{Literal: "it"},
						Fields: []parser.QueryExpression{
							parser.Identifier{Literal: "c1"},
						},
						As: "as",
						Query: parser.SelectQuery{
							SelectEntity: parser.SelectEntity{
								SelectClause: parser.SelectClause{
									Select: "select",
									Fields: []parser.QueryExpression{
										parser.Field{Object: parser.NewIntegerValueFromString("2")},
									},
								},
							},
						},
					},
				},
			},
			Table: parser.Table{Object: parser.Identifier{Literal: "table1"}},
			Fields: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
			},
			ValuesList: []parser.QueryExpression{
				parser.RowValue{
					Value: parser.ValueList{
						Values: []parser.QueryExpression{
							parser.NewIntegerValueFromString("4"),
						},
					},
				},
				parser.RowValue{
					Value: parser.ValueList{
						Values: []parser.QueryExpression{
							parser.Subquery{
								Query: parser.SelectQuery{
									SelectEntity: parser.SelectEntity{
										SelectClause: parser.SelectClause{
											Select: "select",
											Fields: []parser.QueryExpression{
												parser.Field{Object: parser.FieldReference{View: parser.Identifier{Literal: "it"}, Column: parser.Identifier{Literal: "c1"}}},
											},
										},
										FromClause: parser.FromClause{
											Tables: []parser.QueryExpression{
												parser.Table{Object: parser.Identifier{Literal: "it"}},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		ResultFile: &FileInfo{
			Path:      GetTestFilePath("table1.csv"),
			Delimiter: ',',
			NoHeader:  false,
			Encoding:  text.UTF8,
			LineBreak: text.LF,
		},
		UpdateCount: 2,
		ViewCache: ViewMap{
			strings.ToUpper(GetTestFilePath("table1.csv")): &View{
				FileInfo: &FileInfo{
					Path:      GetTestFilePath("table1.csv"),
					Delimiter: ',',
					NoHeader:  false,
					Encoding:  text.UTF8,
					LineBreak: text.LF,
				},
				Header: NewHeader("table1", []string{"column1", "column2"}),
				RecordSet: []Record{
					NewRecord([]value.Primary{
						value.NewString("1"),
						value.NewString("str1"),
					}),
					NewRecord([]value.Primary{
						value.NewString("2"),
						value.NewString("str2"),
					}),
					NewRecord([]value.Primary{
						value.NewString("3"),
						value.NewString("str3"),
					}),
					NewRecord([]value.Primary{
						value.NewInteger(4),
						value.NewNull(),
					}),
					NewRecord([]value.Primary{
						value.NewInteger(2),
						value.NewNull(),
					}),
				},
				ForUpdate: true,
			},
		},
	},
	{
		Name: "Insert Query For Temporary View",
		Query: parser.InsertQuery{
			Table: parser.Table{Object: parser.Identifier{Literal: "tmpview"}, Alias: parser.Identifier{Literal: "t"}},
			Fields: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
			},
			ValuesList: []parser.QueryExpression{
				parser.RowValue{
					Value: parser.ValueList{
						Values: []parser.QueryExpression{
							parser.NewIntegerValueFromString("4"),
						},
					},
				},
				parser.RowValue{
					Value: parser.ValueList{
						Values: []parser.QueryExpression{
							parser.NewIntegerValueFromString("2"),
						},
					},
				},
			},
		},
		ResultFile: &FileInfo{
			Path:        "tmpview",
			Delimiter:   ',',
			IsTemporary: true,
		},
		UpdateCount: 2,
		TempViewList: TemporaryViewScopes{
			ViewMap{
				"TMPVIEW": &View{
					Header: NewHeader("tmpview", []string{"column1", "column2"}),
					RecordSet: []Record{
						NewRecord([]value.Primary{
							value.NewString("1"),
							value.NewString("str1"),
						}),
						NewRecord([]value.Primary{
							value.NewString("2"),
							value.NewString("str2"),
						}),
						NewRecord([]value.Primary{
							value.NewInteger(4),
							value.NewNull(),
						}),
						NewRecord([]value.Primary{
							value.NewInteger(2),
							value.NewNull(),
						}),
					},
					FileInfo: &FileInfo{
						Path:        "tmpview",
						Delimiter:   ',',
						IsTemporary: true,
					},
					ForUpdate: true,
				},
			},
		},
	},
	{
		Name: "Insert Query All Fields",
		Query: parser.InsertQuery{
			Table: parser.Table{Object: parser.Identifier{Literal: "table1"}},
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
		ResultFile: &FileInfo{
			Path:      GetTestFilePath("table1.csv"),
			Delimiter: ',',
			NoHeader:  false,
			Encoding:  text.UTF8,
			LineBreak: text.LF,
		},
		UpdateCount: 2,
	},
	{
		Name: "Insert Query File Does Not Exist Error",
		Query: parser.InsertQuery{
			Table: parser.Table{Object: parser.Identifier{Literal: "notexist"}},
			Fields: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
			},
			ValuesList: []parser.QueryExpression{
				parser.RowValue{
					Value: parser.ValueList{
						Values: []parser.QueryExpression{
							parser.NewIntegerValueFromString("4"),
						},
					},
				},
				parser.RowValue{
					Value: parser.ValueList{
						Values: []parser.QueryExpression{
							parser.NewIntegerValueFromString("5"),
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
			Table: parser.Table{Object: parser.Identifier{Literal: "table1"}},
			Fields: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			},
			ValuesList: []parser.QueryExpression{
				parser.RowValue{
					Value: parser.ValueList{
						Values: []parser.QueryExpression{
							parser.NewIntegerValueFromString("4"),
						},
					},
				},
				parser.RowValue{
					Value: parser.ValueList{
						Values: []parser.QueryExpression{
							parser.NewIntegerValueFromString("5"),
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
			Table: parser.Table{Object: parser.Identifier{Literal: "table1"}},
			Fields: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
			},
			Query: parser.SelectQuery{
				SelectEntity: parser.SelectEntity{
					SelectClause: parser.SelectClause{
						Fields: []parser.QueryExpression{
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column3"}}},
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column4"}}},
						},
					},
					FromClause: parser.FromClause{
						Tables: []parser.QueryExpression{
							parser.Table{Object: parser.Identifier{Literal: "table2"}},
						},
					},
				},
			},
		},
		ResultFile: &FileInfo{
			Path:      GetTestFilePath("table1.csv"),
			Delimiter: ',',
			NoHeader:  false,
			Encoding:  text.UTF8,
			LineBreak: text.LF,
		},
		UpdateCount: 3,
	},
	{
		Name: "Insert Select Query Field Does Not Exist Error",
		Query: parser.InsertQuery{
			Table: parser.Table{Object: parser.Identifier{Literal: "table1"}},
			Fields: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
			},
			Query: parser.SelectQuery{
				SelectEntity: parser.SelectEntity{
					SelectClause: parser.SelectClause{
						Fields: []parser.QueryExpression{
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column3"}}},
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column4"}}},
						},
					},
					FromClause: parser.FromClause{
						Tables: []parser.QueryExpression{
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
	tf.Quiet = false

	filter := NewEmptyFilter()
	filter.TempViews = TemporaryViewScopes{
		ViewMap{
			"TMPVIEW": &View{
				Header: NewHeader("tmpview", []string{"column1", "column2"}),
				RecordSet: []Record{
					NewRecord([]value.Primary{
						value.NewString("1"),
						value.NewString("str1"),
					}),
					NewRecord([]value.Primary{
						value.NewString("2"),
						value.NewString("str2"),
					}),
				},
				FileInfo: &FileInfo{
					Path:        "tmpview",
					Delimiter:   ',',
					IsTemporary: true,
				},
			},
		},
	}

	for _, v := range insertTests {
		ReleaseResources()
		result, cnt, err := Insert(v.Query, filter)
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

		for _, v2 := range ViewCache {
			if v2.FileInfo.Handler != nil {
				if v2.FileInfo.Path != v2.FileInfo.Handler.Path() {
					t.Errorf("file pointer = %q, want %q for %q", v2.FileInfo.Handler.Path(), v2.FileInfo.Path, v.Name)
				}
				v2.FileInfo.Close()
				v2.FileInfo.Handler = nil
			}
		}

		if !reflect.DeepEqual(result, v.ResultFile) {
			t.Errorf("%s: fileinfo = %v, want %v", v.Name, result, v.ResultFile)
		}

		if !reflect.DeepEqual(cnt, v.UpdateCount) {
			t.Errorf("%s: update count = %d, want %d", v.Name, cnt, v.UpdateCount)
		}

		if v.ViewCache != nil {
			if !reflect.DeepEqual(ViewCache, v.ViewCache) {
				t.Errorf("%s: view cache = %v, want %v", v.Name, ViewCache, v.ViewCache)
			}
		}
		if v.TempViewList != nil {
			if !reflect.DeepEqual(filter.TempViews, v.TempViewList) {
				t.Errorf("%s: temporary views list = %v, want %v", v.Name, filter.TempViews, v.TempViewList)
			}
		}
	}
	ReleaseResources()
}

var updateTests = []struct {
	Name         string
	Query        parser.UpdateQuery
	ResultFiles  []*FileInfo
	UpdateCounts []int
	ViewCache    ViewMap
	TempViewList TemporaryViewScopes
	Error        string
}{
	{
		Name: "Update Query",
		Query: parser.UpdateQuery{
			WithClause: parser.WithClause{
				With: "with",
				InlineTables: []parser.QueryExpression{
					parser.InlineTable{
						Name: parser.Identifier{Literal: "it"},
						Fields: []parser.QueryExpression{
							parser.Identifier{Literal: "c1"},
						},
						As: "as",
						Query: parser.SelectQuery{
							SelectEntity: parser.SelectEntity{
								SelectClause: parser.SelectClause{
									Select: "select",
									Fields: []parser.QueryExpression{
										parser.Field{Object: parser.NewIntegerValueFromString("2")},
									},
								},
							},
						},
					},
				},
			},
			Tables: []parser.QueryExpression{
				parser.Table{Object: parser.Identifier{Literal: "table1"}},
			},
			SetList: []parser.UpdateSet{
				{
					Field: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
					Value: parser.NewStringValue("update1"),
				},
				{
					Field: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
					Value: parser.NewStringValue("update2"),
				},
			},
			WhereClause: parser.WhereClause{
				Filter: parser.Comparison{
					LHS: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
					RHS: parser.Subquery{
						Query: parser.SelectQuery{
							SelectEntity: parser.SelectEntity{
								SelectClause: parser.SelectClause{
									Select: "select",
									Fields: []parser.QueryExpression{
										parser.Field{Object: parser.FieldReference{View: parser.Identifier{Literal: "it"}, Column: parser.Identifier{Literal: "c1"}}},
									},
								},
								FromClause: parser.FromClause{
									Tables: []parser.QueryExpression{
										parser.Table{Object: parser.Identifier{Literal: "it"}},
									},
								},
							},
						},
					},
					Operator: "=",
				},
			},
		},
		ResultFiles: []*FileInfo{
			{
				Path:      GetTestFilePath("table1.csv"),
				Delimiter: ',',
				NoHeader:  false,
				Encoding:  text.UTF8,
				LineBreak: text.LF,
			},
		},
		UpdateCounts: []int{1},
		ViewCache: ViewMap{
			strings.ToUpper(GetTestFilePath("table1.csv")): &View{
				FileInfo: &FileInfo{
					Path:      GetTestFilePath("table1.csv"),
					Delimiter: ',',
					NoHeader:  false,
					Encoding:  text.UTF8,
					LineBreak: text.LF,
				},
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
				ForUpdate: true,
			},
		},
	},
	{
		Name: "Update Query For Temporary View",
		Query: parser.UpdateQuery{
			Tables: []parser.QueryExpression{
				parser.Table{Object: parser.Identifier{Literal: "tmpview"}, Alias: parser.Identifier{Literal: "t1"}},
			},
			SetList: []parser.UpdateSet{
				{
					Field: parser.ColumnNumber{View: parser.Identifier{Literal: "t1"}, Number: value.NewInteger(2)},
					Value: parser.NewStringValue("update"),
				},
			},
		},
		ResultFiles: []*FileInfo{
			{
				Path:        "tmpview",
				Delimiter:   ',',
				IsTemporary: true,
			},
		},
		UpdateCounts: []int{2},
		TempViewList: TemporaryViewScopes{
			ViewMap{
				"TMPVIEW": &View{
					Header: NewHeader("tmpview", []string{"column1", "column2"}),
					RecordSet: []Record{
						NewRecord([]value.Primary{
							value.NewString("1"),
							value.NewString("update"),
						}),
						NewRecord([]value.Primary{
							value.NewString("2"),
							value.NewString("update"),
						}),
					},
					FileInfo: &FileInfo{
						Path:        "tmpview",
						Delimiter:   ',',
						IsTemporary: true,
					},
				},
			},
		},
	},
	{
		Name: "Update Query Multiple Table",
		Query: parser.UpdateQuery{
			Tables: []parser.QueryExpression{
				parser.Table{Object: parser.Identifier{Literal: "t1"}},
			},
			SetList: []parser.UpdateSet{
				{
					Field: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
					Value: parser.FieldReference{Column: parser.Identifier{Literal: "column4"}},
				},
			},
			FromClause: parser.FromClause{
				Tables: []parser.QueryExpression{
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
		ResultFiles: []*FileInfo{
			{
				Path:      GetTestFilePath("table1.csv"),
				Delimiter: ',',
				NoHeader:  false,
				Encoding:  text.UTF8,
				LineBreak: text.LF,
			},
		},
		UpdateCounts: []int{2},
	},
	{
		Name: "Update Query File Does Not Exist Error",
		Query: parser.UpdateQuery{
			Tables: []parser.QueryExpression{
				parser.Table{Object: parser.Identifier{Literal: "notexist"}},
			},
			SetList: []parser.UpdateSet{
				{
					Field: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
					Value: parser.NewStringValue("update"),
				},
			},
			WhereClause: parser.WhereClause{
				Filter: parser.Comparison{
					LHS:      parser.Identifier{Literal: "column1"},
					RHS:      parser.NewIntegerValueFromString("2"),
					Operator: "=",
				},
			},
		},
		Error: "[L:- C:-] file notexist does not exist",
	},
	{
		Name: "Update Query Filter Error",
		Query: parser.UpdateQuery{
			Tables: []parser.QueryExpression{
				parser.Table{Object: parser.Identifier{Literal: "table1"}},
			},
			SetList: []parser.UpdateSet{
				{
					Field: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
					Value: parser.NewStringValue("update"),
				},
			},
			WhereClause: parser.WhereClause{
				Filter: parser.Comparison{
					LHS:      parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
					RHS:      parser.NewIntegerValueFromString("2"),
					Operator: "=",
				},
			},
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
	{
		Name: "Update Query File Is Not Loaded Error",
		Query: parser.UpdateQuery{
			Tables: []parser.QueryExpression{
				parser.Table{Object: parser.Identifier{Literal: "notexist"}},
			},
			SetList: []parser.UpdateSet{
				{
					Field: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
					Value: parser.FieldReference{Column: parser.Identifier{Literal: "column4"}},
				},
			},
			FromClause: parser.FromClause{
				Tables: []parser.QueryExpression{
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
			Tables: []parser.QueryExpression{
				parser.Table{Object: parser.Identifier{Literal: "t2"}},
			},
			SetList: []parser.UpdateSet{
				{
					Field: parser.FieldReference{View: parser.Identifier{Literal: "t1"}, Column: parser.Identifier{Literal: "column2"}},
					Value: parser.FieldReference{Column: parser.Identifier{Literal: "column4"}},
				},
			},
			FromClause: parser.FromClause{
				Tables: []parser.QueryExpression{
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
			Tables: []parser.QueryExpression{
				parser.Table{Object: parser.Identifier{Literal: "table1"}},
			},
			SetList: []parser.UpdateSet{
				{
					Field: parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
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
		Error: "[L:- C:-] field notexist does not exist",
	},
	{
		Name: "Update Query Update Value Error",
		Query: parser.UpdateQuery{
			Tables: []parser.QueryExpression{
				parser.Table{Object: parser.Identifier{Literal: "table1"}},
			},
			SetList: []parser.UpdateSet{
				{
					Field: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
					Value: parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
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
		Error: "[L:- C:-] field notexist does not exist",
	},
	{
		Name: "Update Query Record Is Ambiguous Error",
		Query: parser.UpdateQuery{
			Tables: []parser.QueryExpression{
				parser.Table{Object: parser.Identifier{Literal: "t1"}},
			},
			SetList: []parser.UpdateSet{
				{
					Field: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
					Value: parser.FieldReference{Column: parser.Identifier{Literal: "column4"}},
				},
			},
			FromClause: parser.FromClause{
				Tables: []parser.QueryExpression{
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
	tf.Quiet = false

	filter := NewEmptyFilter()
	filter.TempViews = TemporaryViewScopes{
		ViewMap{
			"TMPVIEW": &View{
				Header: NewHeader("tmpview", []string{"column1", "column2"}),
				RecordSet: []Record{
					NewRecord([]value.Primary{
						value.NewString("1"),
						value.NewString("str1"),
					}),
					NewRecord([]value.Primary{
						value.NewString("2"),
						value.NewString("str2"),
					}),
				},
				FileInfo: &FileInfo{
					Path:        "tmpview",
					Delimiter:   ',',
					IsTemporary: true,
				},
			},
		},
	}

	for _, v := range updateTests {
		ReleaseResources()
		files, cnt, err := Update(v.Query, filter)
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

		for _, v2 := range ViewCache {
			if v2.FileInfo.Handler != nil {
				if v2.FileInfo.Path != v2.FileInfo.Handler.Path() {
					t.Errorf("file pointer = %q, want %q for %q", v2.FileInfo.Handler.Path(), v2.FileInfo.Path, v.Name)
				}
				v2.FileInfo.Close()
				v2.FileInfo.Handler = nil
			}
		}

		if !reflect.DeepEqual(files, v.ResultFiles) {
			t.Errorf("%s: fileinfo = %v, want %v", v.Name, files, v.ResultFiles)
		}

		if !reflect.DeepEqual(cnt, v.UpdateCounts) {
			t.Errorf("%s: update count = %v, want %v", v.Name, cnt, v.UpdateCounts)
		}

		if v.ViewCache != nil {
			if !reflect.DeepEqual(ViewCache, v.ViewCache) {
				t.Errorf("%s: view cache = %v, want %v", v.Name, ViewCache, v.ViewCache)
			}
		}
		if v.TempViewList != nil {
			if !reflect.DeepEqual(filter.TempViews, v.TempViewList) {
				t.Errorf("%s: temporary views list = %v, want %v", v.Name, filter.TempViews, v.TempViewList)
			}
		}
	}
	ReleaseResources()
}

var deleteTests = []struct {
	Name         string
	Query        parser.DeleteQuery
	ResultFiles  []*FileInfo
	UpdateCounts []int
	ViewCache    ViewMap
	TempViewList TemporaryViewScopes
	Error        string
}{
	{
		Name: "Delete Query",
		Query: parser.DeleteQuery{
			WithClause: parser.WithClause{
				With: "with",
				InlineTables: []parser.QueryExpression{
					parser.InlineTable{
						Name: parser.Identifier{Literal: "it"},
						Fields: []parser.QueryExpression{
							parser.Identifier{Literal: "c1"},
						},
						As: "as",
						Query: parser.SelectQuery{
							SelectEntity: parser.SelectEntity{
								SelectClause: parser.SelectClause{
									Select: "select",
									Fields: []parser.QueryExpression{
										parser.Field{Object: parser.NewIntegerValueFromString("2")},
									},
								},
							},
						},
					},
				},
			},
			FromClause: parser.FromClause{
				Tables: []parser.QueryExpression{
					parser.Table{
						Object: parser.Identifier{Literal: "table1"},
					},
				},
			},
			WhereClause: parser.WhereClause{
				Filter: parser.Comparison{
					LHS: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
					RHS: parser.Subquery{
						Query: parser.SelectQuery{
							SelectEntity: parser.SelectEntity{
								SelectClause: parser.SelectClause{
									Select: "select",
									Fields: []parser.QueryExpression{
										parser.Field{Object: parser.FieldReference{View: parser.Identifier{Literal: "it"}, Column: parser.Identifier{Literal: "c1"}}},
									},
								},
								FromClause: parser.FromClause{
									Tables: []parser.QueryExpression{
										parser.Table{Object: parser.Identifier{Literal: "it"}},
									},
								},
							},
						},
					},
					Operator: "=",
				},
			},
		},
		ResultFiles: []*FileInfo{
			{
				Path:      GetTestFilePath("table1.csv"),
				Delimiter: ',',
				NoHeader:  false,
				Encoding:  text.UTF8,
				LineBreak: text.LF,
			},
		},
		UpdateCounts: []int{1},
		ViewCache: ViewMap{
			strings.ToUpper(GetTestFilePath("table1.csv")): &View{
				FileInfo: &FileInfo{
					Path:      GetTestFilePath("table1.csv"),
					Delimiter: ',',
					NoHeader:  false,
					Encoding:  text.UTF8,
					LineBreak: text.LF,
				},
				Header: NewHeader("table1", []string{"column1", "column2"}),
				RecordSet: []Record{
					NewRecord([]value.Primary{
						value.NewString("1"),
						value.NewString("str1"),
					}),
					NewRecord([]value.Primary{
						value.NewString("3"),
						value.NewString("str3"),
					}),
				},
				ForUpdate: true,
			},
		},
	},
	{
		Name: "Delete Query For Temporary View",
		Query: parser.DeleteQuery{
			FromClause: parser.FromClause{
				Tables: []parser.QueryExpression{
					parser.Table{
						Object: parser.Identifier{Literal: "tmpview"},
						Alias:  parser.Identifier{Literal: "t1"},
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
		ResultFiles: []*FileInfo{
			{
				Path:        "tmpview",
				Delimiter:   ',',
				IsTemporary: true,
			},
		},
		UpdateCounts: []int{1},
		TempViewList: TemporaryViewScopes{
			ViewMap{
				"TMPVIEW": &View{
					Header: NewHeader("tmpview", []string{"column1", "column2"}),
					RecordSet: []Record{
						NewRecord([]value.Primary{
							value.NewString("1"),
							value.NewString("str1"),
						}),
					},
					FileInfo: &FileInfo{
						Path:        "tmpview",
						Delimiter:   ',',
						IsTemporary: true,
					},
				},
			},
		},
	},
	{
		Name: "Delete Query Multiple Table",
		Query: parser.DeleteQuery{
			Tables: []parser.QueryExpression{
				parser.Table{Object: parser.Identifier{Literal: "t1"}},
			},
			FromClause: parser.FromClause{
				Tables: []parser.QueryExpression{
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
		ResultFiles: []*FileInfo{
			{
				Path:      GetTestFilePath("table1.csv"),
				Delimiter: ',',
				NoHeader:  false,
				Encoding:  text.UTF8,
				LineBreak: text.LF,
			},
		},
		UpdateCounts: []int{2},
	},
	{
		Name: "Delete Query Tables Not Specified Error",
		Query: parser.DeleteQuery{
			FromClause: parser.FromClause{
				Tables: []parser.QueryExpression{
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
			FromClause: parser.FromClause{
				Tables: []parser.QueryExpression{
					parser.Table{
						Object: parser.Identifier{Literal: "notexist"},
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
		Error: "[L:- C:-] file notexist does not exist",
	},
	{
		Name: "Delete Query Filter Error",
		Query: parser.DeleteQuery{
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
			Tables: []parser.QueryExpression{
				parser.Table{Object: parser.Identifier{Literal: "notexist"}},
			},
			FromClause: parser.FromClause{
				Tables: []parser.QueryExpression{
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
	tf.Quiet = false

	filter := NewEmptyFilter()
	filter.TempViews = TemporaryViewScopes{
		ViewMap{
			"TMPVIEW": &View{
				Header: NewHeader("tmpview", []string{"column1", "column2"}),
				RecordSet: []Record{
					NewRecord([]value.Primary{
						value.NewString("1"),
						value.NewString("str1"),
					}),
					NewRecord([]value.Primary{
						value.NewString("2"),
						value.NewString("str2"),
					}),
				},
				FileInfo: &FileInfo{
					Path:        "tmpview",
					Delimiter:   ',',
					IsTemporary: true,
				},
			},
		},
	}

	for _, v := range deleteTests {
		ReleaseResources()
		files, cnt, err := Delete(v.Query, filter)
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

		for _, v2 := range ViewCache {
			if v2.FileInfo.Handler != nil {
				if v2.FileInfo.Path != v2.FileInfo.Handler.Path() {
					t.Errorf("file pointer = %q, want %q for %q", v2.FileInfo.Handler.Path(), v2.FileInfo.Path, v.Name)
				}
				v2.FileInfo.Close()
				v2.FileInfo.Handler = nil
			}
		}

		if !reflect.DeepEqual(files, v.ResultFiles) {
			t.Errorf("%s: fileinfo = %v, want %v", v.Name, files, v.ResultFiles)
		}

		if !reflect.DeepEqual(cnt, v.UpdateCounts) {
			t.Errorf("%s: update count = %v, want %v", v.Name, cnt, v.UpdateCounts)
		}

		if v.ViewCache != nil {
			if !reflect.DeepEqual(ViewCache, v.ViewCache) {
				t.Errorf("%s: view cache = %v, want %v", v.Name, ViewCache, v.ViewCache)
			}
		}
		if v.TempViewList != nil {
			if !reflect.DeepEqual(filter.TempViews, v.TempViewList) {
				t.Errorf("%s: temporary views list = %v, want %v", v.Name, filter.TempViews, v.TempViewList)
			}
		}
	}
	ReleaseResources()
}

var createTableTests = []struct {
	Name       string
	Query      parser.CreateTable
	ResultFile *FileInfo
	ViewCache  ViewMap
	Error      string
}{
	{
		Name: "Create Table",
		Query: parser.CreateTable{
			Table: parser.Identifier{Literal: "create_table_1.csv"},
			Fields: []parser.QueryExpression{
				parser.Identifier{Literal: "column1"},
				parser.Identifier{Literal: "column2"},
			},
		},
		ResultFile: &FileInfo{
			Path:      GetTestFilePath("create_table_1.csv"),
			Delimiter: ',',
			NoHeader:  false,
			Encoding:  text.UTF8,
			LineBreak: text.LF,
		},
		ViewCache: ViewMap{
			strings.ToUpper(GetTestFilePath("create_table_1.csv")): &View{
				FileInfo: &FileInfo{
					Path:      GetTestFilePath("create_table_1.csv"),
					Delimiter: ',',
					NoHeader:  false,
					Encoding:  text.UTF8,
					LineBreak: text.LF,
				},
				Header:    NewHeader("create_table_1", []string{"column1", "column2"}),
				RecordSet: RecordSet{},
				ForUpdate: true,
			},
		},
	},
	{
		Name: "Create Table From Select Query",
		Query: parser.CreateTable{
			Table: parser.Identifier{Literal: "create_table_1.csv"},
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
		ResultFile: &FileInfo{
			Path:      GetTestFilePath("create_table_1.csv"),
			Delimiter: ',',
			NoHeader:  false,
			Encoding:  text.UTF8,
			LineBreak: text.LF,
		},
		ViewCache: ViewMap{
			strings.ToUpper(GetTestFilePath("create_table_1.csv")): &View{
				FileInfo: &FileInfo{
					Path:      GetTestFilePath("create_table_1.csv"),
					Delimiter: ',',
					NoHeader:  false,
					Encoding:  text.UTF8,
					LineBreak: text.LF,
				},
				Header: NewHeader("create_table_1", []string{"column1", "column2"}),
				RecordSet: RecordSet{
					NewRecord([]value.Primary{
						value.NewInteger(1),
						value.NewInteger(2),
					}),
				},
				ForUpdate: true,
			},
		},
	},
	{
		Name: "Create Table File Already Exist Error",
		Query: parser.CreateTable{
			Table: parser.Identifier{Literal: "table1.csv"},
			Fields: []parser.QueryExpression{
				parser.Identifier{Literal: "column1"},
				parser.Identifier{Literal: "column2"},
			},
		},
		Error: "[L:- C:-] file table1.csv already exists",
	},
	{
		Name: "Create Table Field Duplicate Error",
		Query: parser.CreateTable{
			Table: parser.Identifier{Literal: "create_table_1.csv"},
			Fields: []parser.QueryExpression{
				parser.Identifier{Literal: "column1"},
				parser.Identifier{Literal: "column1"},
			},
		},
		Error: "[L:- C:-] field name column1 is a duplicate",
	},
	{
		Name: "Create Table Select Query Execution Error",
		Query: parser.CreateTable{
			Table: parser.Identifier{Literal: "create_table_1.csv"},
			Fields: []parser.QueryExpression{
				parser.Identifier{Literal: "column1"},
				parser.Identifier{Literal: "column2"},
			},
			Query: parser.SelectQuery{
				SelectEntity: parser.SelectEntity{
					SelectClause: parser.SelectClause{
						Fields: []parser.QueryExpression{
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}}},
							parser.Field{Object: parser.NewIntegerValueFromString("2")},
						},
					},
				},
			},
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
	{
		Name: "Create Table From Select Query Field Length Not Match Error",
		Query: parser.CreateTable{
			Table: parser.Identifier{Literal: "create_table_1.csv"},
			Fields: []parser.QueryExpression{
				parser.Identifier{Literal: "column1"},
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
		Error: "[L:- C:-] select query should return exactly 1 field for table create_table_1.csv",
	},
	{
		Name: "Create Table From Select Query Field Name Duplicate Error",
		Query: parser.CreateTable{
			Table: parser.Identifier{Literal: "create_table_1.csv"},
			Fields: []parser.QueryExpression{
				parser.Identifier{Literal: "column1"},
				parser.Identifier{Literal: "column1"},
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
		Error: "[L:- C:-] field name column1 is a duplicate",
	},
}

func TestCreateTable(t *testing.T) {
	tf := cmd.GetFlags()
	tf.Repository = TestDir
	tf.Quiet = false

	for _, v := range createTableTests {
		ReleaseResources()

		result, err := CreateTable(v.Query, NewEmptyFilter())

		if result != nil {
			result.Close()
			result.Handler = nil
		}
		for _, view := range ViewCache {
			if view.FileInfo != nil {
				view.FileInfo.Close()
				view.FileInfo.Handler = nil
			}
		}

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

		if !reflect.DeepEqual(result, v.ResultFile) {
			t.Errorf("%s: result = %v, want %v", v.Name, result, v.ResultFile)
		}

		if v.ViewCache != nil {
			if !reflect.DeepEqual(ViewCache, v.ViewCache) {
				t.Errorf("%s: view cache = %v, want %v", v.Name, ViewCache, v.ViewCache)
			}
		}
	}
	ReleaseResources()
}

var addColumnsTests = []struct {
	Name         string
	Query        parser.AddColumns
	ResultFile   *FileInfo
	UpdateCount  int
	ViewCache    ViewMap
	TempViewList TemporaryViewScopes
	Error        string
}{
	{
		Name: "Add Fields",
		Query: parser.AddColumns{
			Table: parser.Identifier{Literal: "table1.csv"},
			Columns: []parser.ColumnDefault{
				{
					Column: parser.Identifier{Literal: "column3"},
				},
				{
					Column: parser.Identifier{Literal: "column4"},
				},
			},
		},
		ResultFile: &FileInfo{
			Path:      GetTestFilePath("table1.csv"),
			Delimiter: ',',
			NoHeader:  false,
			Encoding:  text.UTF8,
			LineBreak: text.LF,
		},
		UpdateCount: 2,
		ViewCache: ViewMap{
			strings.ToUpper(GetTestFilePath("table1.csv")): &View{
				FileInfo: &FileInfo{
					Path:      GetTestFilePath("table1.csv"),
					Delimiter: ',',
					NoHeader:  false,
					Encoding:  text.UTF8,
					LineBreak: text.LF,
				},
				Header: NewHeader("table1", []string{"column1", "column2", "column3", "column4"}),
				RecordSet: []Record{
					NewRecord([]value.Primary{
						value.NewString("1"),
						value.NewString("str1"),
						value.NewNull(),
						value.NewNull(),
					}),
					NewRecord([]value.Primary{
						value.NewString("2"),
						value.NewString("str2"),
						value.NewNull(),
						value.NewNull(),
					}),
					NewRecord([]value.Primary{
						value.NewString("3"),
						value.NewString("str3"),
						value.NewNull(),
						value.NewNull(),
					}),
				},
				ForUpdate: true,
			},
		},
	},
	{
		Name: "Add Fields For Temporary View",
		Query: parser.AddColumns{
			Table: parser.Identifier{Literal: "tmpview"},
			Columns: []parser.ColumnDefault{
				{
					Column: parser.Identifier{Literal: "column3"},
				},
				{
					Column: parser.Identifier{Literal: "column4"},
				},
			},
		},
		ResultFile: &FileInfo{
			Path:        "tmpview",
			Delimiter:   ',',
			IsTemporary: true,
		},
		UpdateCount: 2,
		TempViewList: TemporaryViewScopes{
			ViewMap{
				"TMPVIEW": &View{
					Header: NewHeader("tmpview", []string{"column1", "column2", "column3", "column4"}),
					RecordSet: []Record{
						NewRecord([]value.Primary{
							value.NewString("1"),
							value.NewString("str1"),
							value.NewNull(),
							value.NewNull(),
						}),
						NewRecord([]value.Primary{
							value.NewString("2"),
							value.NewString("str2"),
							value.NewNull(),
							value.NewNull(),
						}),
					},
					FileInfo: &FileInfo{
						Path:        "tmpview",
						Delimiter:   ',',
						IsTemporary: true,
					},
					ForUpdate: true,
				},
			},
		},
	},
	{
		Name: "Add Fields First",
		Query: parser.AddColumns{
			Table: parser.Identifier{Literal: "table1.csv"},
			Columns: []parser.ColumnDefault{
				{
					Column: parser.Identifier{Literal: "column3"},
					Value:  parser.NewIntegerValueFromString("2"),
				},
				{
					Column: parser.Identifier{Literal: "column4"},
					Value:  parser.NewIntegerValueFromString("1"),
				},
			},
			Position: parser.ColumnPosition{
				Position: parser.Token{Token: parser.FIRST},
			},
		},
		ResultFile: &FileInfo{
			Path:      GetTestFilePath("table1.csv"),
			Delimiter: ',',
			NoHeader:  false,
			Encoding:  text.UTF8,
			LineBreak: text.LF,
		},
		UpdateCount: 2,
		ViewCache: ViewMap{
			strings.ToUpper(GetTestFilePath("table1.csv")): &View{
				FileInfo: &FileInfo{
					Path:      GetTestFilePath("table1.csv"),
					Delimiter: ',',
					NoHeader:  false,
					Encoding:  text.UTF8,
					LineBreak: text.LF,
				},
				Header: NewHeader("table1", []string{"column3", "column4", "column1", "column2"}),
				RecordSet: []Record{
					NewRecord([]value.Primary{
						value.NewInteger(2),
						value.NewInteger(1),
						value.NewString("1"),
						value.NewString("str1"),
					}),
					NewRecord([]value.Primary{
						value.NewInteger(2),
						value.NewInteger(1),
						value.NewString("2"),
						value.NewString("str2"),
					}),
					NewRecord([]value.Primary{
						value.NewInteger(2),
						value.NewInteger(1),
						value.NewString("3"),
						value.NewString("str3"),
					}),
				},
				ForUpdate: true,
			},
		},
	},
	{
		Name: "Add Fields After",
		Query: parser.AddColumns{
			Table: parser.Identifier{Literal: "table1.csv"},
			Columns: []parser.ColumnDefault{
				{
					Column: parser.Identifier{Literal: "column3"},
				},
				{
					Column: parser.Identifier{Literal: "column4"},
					Value:  parser.NewIntegerValueFromString("1"),
				},
			},
			Position: parser.ColumnPosition{
				Position: parser.Token{Token: parser.AFTER},
				Column:   parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
			},
		},
		ResultFile: &FileInfo{
			Path:      GetTestFilePath("table1.csv"),
			Delimiter: ',',
			NoHeader:  false,
			Encoding:  text.UTF8,
			LineBreak: text.LF,
		},
		UpdateCount: 2,
		ViewCache: ViewMap{
			strings.ToUpper(GetTestFilePath("table1.csv")): &View{
				FileInfo: &FileInfo{
					Path:      GetTestFilePath("table1.csv"),
					Delimiter: ',',
					NoHeader:  false,
					Encoding:  text.UTF8,
					LineBreak: text.LF,
				},
				Header: NewHeader("table1", []string{"column1", "column3", "column4", "column2"}),
				RecordSet: []Record{
					NewRecord([]value.Primary{
						value.NewString("1"),
						value.NewNull(),
						value.NewInteger(1),
						value.NewString("str1"),
					}),
					NewRecord([]value.Primary{
						value.NewString("2"),
						value.NewNull(),
						value.NewInteger(1),
						value.NewString("str2"),
					}),
					NewRecord([]value.Primary{
						value.NewString("3"),
						value.NewNull(),
						value.NewInteger(1),
						value.NewString("str3"),
					}),
				},
				ForUpdate: true,
			},
		},
	},
	{
		Name: "Add Fields Before",
		Query: parser.AddColumns{
			Table: parser.Identifier{Literal: "table1.csv"},
			Columns: []parser.ColumnDefault{
				{
					Column: parser.Identifier{Literal: "column3"},
				},
				{
					Column: parser.Identifier{Literal: "column4"},
					Value:  parser.NewIntegerValueFromString("1"),
				},
			},
			Position: parser.ColumnPosition{
				Position: parser.Token{Token: parser.BEFORE},
				Column:   parser.ColumnNumber{View: parser.Identifier{Literal: "table1"}, Number: value.NewInteger(2)},
			},
		},
		ResultFile: &FileInfo{
			Path:      GetTestFilePath("table1.csv"),
			Delimiter: ',',
			NoHeader:  false,
			Encoding:  text.UTF8,
			LineBreak: text.LF,
		},
		UpdateCount: 2,
		ViewCache: ViewMap{
			strings.ToUpper(GetTestFilePath("table1.csv")): &View{
				FileInfo: &FileInfo{
					Path:      GetTestFilePath("table1.csv"),
					Delimiter: ',',
					NoHeader:  false,
					Encoding:  text.UTF8,
					LineBreak: text.LF,
				},
				Header: NewHeader("table1", []string{"column1", "column3", "column4", "column2"}),
				RecordSet: []Record{
					NewRecord([]value.Primary{
						value.NewString("1"),
						value.NewNull(),
						value.NewInteger(1),
						value.NewString("str1"),
					}),
					NewRecord([]value.Primary{
						value.NewString("2"),
						value.NewNull(),
						value.NewInteger(1),
						value.NewString("str2"),
					}),
					NewRecord([]value.Primary{
						value.NewString("3"),
						value.NewNull(),
						value.NewInteger(1),
						value.NewString("str3"),
					}),
				},
				ForUpdate: true,
			},
		},
	},
	{
		Name: "Add Fields Load Error",
		Query: parser.AddColumns{
			Table: parser.Identifier{Literal: "notexist"},
			Columns: []parser.ColumnDefault{
				{
					Column: parser.Identifier{Literal: "column3"},
				},
				{
					Column: parser.Identifier{Literal: "column4"},
				},
			},
		},
		Error: "[L:- C:-] file notexist does not exist",
	},
	{
		Name: "Add Fields Position Column Does Not Exist Error",
		Query: parser.AddColumns{
			Table: parser.Identifier{Literal: "table1.csv"},
			Columns: []parser.ColumnDefault{
				{
					Column: parser.Identifier{Literal: "column3"},
				},
				{
					Column: parser.Identifier{Literal: "column2"},
					Value:  parser.NewIntegerValueFromString("1"),
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
		Name: "Add Fields Field Duplicate Error",
		Query: parser.AddColumns{
			Table: parser.Identifier{Literal: "table1.csv"},
			Columns: []parser.ColumnDefault{
				{
					Column: parser.Identifier{Literal: "column3"},
				},
				{
					Column: parser.Identifier{Literal: "column1"},
					Value:  parser.NewIntegerValueFromString("1"),
				},
			},
		},
		Error: "[L:- C:-] field name column1 is a duplicate",
	},
	{
		Name: "Add Fields Default Value Error",
		Query: parser.AddColumns{
			Table: parser.Identifier{Literal: "table1.csv"},
			Columns: []parser.ColumnDefault{
				{
					Column: parser.Identifier{Literal: "column3"},
				},
				{
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
	tf.Quiet = false

	filter := NewEmptyFilter()
	filter.TempViews = TemporaryViewScopes{
		ViewMap{
			"TMPVIEW": &View{
				Header: NewHeader("tmpview", []string{"column1", "column2"}),
				RecordSet: []Record{
					NewRecord([]value.Primary{
						value.NewString("1"),
						value.NewString("str1"),
					}),
					NewRecord([]value.Primary{
						value.NewString("2"),
						value.NewString("str2"),
					}),
				},
				FileInfo: &FileInfo{
					Path:        "tmpview",
					Delimiter:   ',',
					IsTemporary: true,
				},
			},
		},
	}
	for _, v := range addColumnsTests {
		ReleaseResources()
		result, cnt, err := AddColumns(v.Query, filter)
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

		for _, v2 := range ViewCache {
			if v2.FileInfo.Handler != nil {
				if v2.FileInfo.Path != v2.FileInfo.Handler.Path() {
					t.Errorf("file pointer = %q, want %q for %q", v2.FileInfo.Handler.Path(), v2.FileInfo.Path, v.Name)
				}
				v2.FileInfo.Close()
				v2.FileInfo.Handler = nil
			}
		}

		if !reflect.DeepEqual(result, v.ResultFile) {
			t.Errorf("%s: fileinfo = %v, want %v", v.Name, result, v.ResultFile)
		}

		if cnt != v.UpdateCount {
			t.Errorf("%s: update count = %d, want %d", v.Name, cnt, v.UpdateCount)
		}

		if v.ViewCache != nil {
			if !reflect.DeepEqual(ViewCache, v.ViewCache) {
				t.Errorf("%s: view cache = %v, want %v", v.Name, ViewCache, v.ViewCache)
			}
		}
		if v.TempViewList != nil {
			if !reflect.DeepEqual(filter.TempViews, v.TempViewList) {
				t.Errorf("%s: temporary views list = %v, want %v", v.Name, filter.TempViews, v.TempViewList)
			}
		}
	}
	ReleaseResources()
}

var dropColumnsTests = []struct {
	Name         string
	Query        parser.DropColumns
	Result       *FileInfo
	UpdateCount  int
	ViewCache    ViewMap
	TempViewList TemporaryViewScopes
	Error        string
}{
	{
		Name: "Drop Fields",
		Query: parser.DropColumns{
			Table: parser.Identifier{Literal: "table1"},
			Columns: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
			},
		},
		Result: &FileInfo{
			Path:      GetTestFilePath("table1.csv"),
			Delimiter: ',',
			NoHeader:  false,
			Encoding:  text.UTF8,
			LineBreak: text.LF,
		},
		UpdateCount: 1,
		ViewCache: ViewMap{
			strings.ToUpper(GetTestFilePath("table1.csv")): &View{
				FileInfo: &FileInfo{
					Path:      GetTestFilePath("table1.csv"),
					Delimiter: ',',
					NoHeader:  false,
					Encoding:  text.UTF8,
					LineBreak: text.LF,
				},
				Header: NewHeader("table1", []string{"column1"}),
				RecordSet: []Record{
					NewRecord([]value.Primary{
						value.NewString("1"),
					}),
					NewRecord([]value.Primary{
						value.NewString("2"),
					}),
					NewRecord([]value.Primary{
						value.NewString("3"),
					}),
				},
				ForUpdate: true,
			},
		},
	},
	{
		Name: "Drop Fields For Temporary View",
		Query: parser.DropColumns{
			Table: parser.Identifier{Literal: "tmpview"},
			Columns: []parser.QueryExpression{
				parser.ColumnNumber{View: parser.Identifier{Literal: "tmpview"}, Number: value.NewInteger(2)},
			},
		},
		Result: &FileInfo{
			Path:        "tmpview",
			Delimiter:   ',',
			IsTemporary: true,
		},
		UpdateCount: 1,
		TempViewList: TemporaryViewScopes{
			ViewMap{
				"TMPVIEW": &View{
					Header: NewHeader("tmpview", []string{"column1"}),
					RecordSet: []Record{
						NewRecord([]value.Primary{
							value.NewString("1"),
						}),
						NewRecord([]value.Primary{
							value.NewString("2"),
						}),
					},
					FileInfo: &FileInfo{
						Path:        "tmpview",
						Delimiter:   ',',
						IsTemporary: true,
					},
					ForUpdate: true,
				},
			},
		},
	},
	{
		Name: "Drop Fields Load Error",
		Query: parser.DropColumns{
			Table: parser.Identifier{Literal: "notexist"},
			Columns: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
			},
		},
		Error: "[L:- C:-] file notexist does not exist",
	},
	{
		Name: "Drop Fields Field Does Not Exist Error",
		Query: parser.DropColumns{
			Table: parser.Identifier{Literal: "table1"},
			Columns: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			},
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
}

func TestDropColumns(t *testing.T) {
	tf := cmd.GetFlags()
	tf.Repository = TestDir
	tf.Quiet = false

	filter := NewEmptyFilter()
	filter.TempViews = TemporaryViewScopes{
		ViewMap{
			"TMPVIEW": &View{
				Header: NewHeader("tmpview", []string{"column1", "column2"}),
				RecordSet: []Record{
					NewRecord([]value.Primary{
						value.NewString("1"),
						value.NewString("str1"),
					}),
					NewRecord([]value.Primary{
						value.NewString("2"),
						value.NewString("str2"),
					}),
				},
				FileInfo: &FileInfo{
					Path:        "tmpview",
					Delimiter:   ',',
					IsTemporary: true,
				},
			},
		},
	}

	for _, v := range dropColumnsTests {
		ReleaseResources()
		result, cnt, err := DropColumns(v.Query, filter)
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

		for _, v2 := range ViewCache {
			if v2.FileInfo.Handler != nil {
				if v2.FileInfo.Path != v2.FileInfo.Handler.Path() {
					t.Errorf("file pointer = %q, want %q for %q", v2.FileInfo.Handler.Path(), v2.FileInfo.Path, v.Name)
				}
				v2.FileInfo.Close()
				v2.FileInfo.Handler = nil
			}
		}

		if !reflect.DeepEqual(result, v.Result) {
			t.Errorf("%s: fileinfo = %v, want %v", v.Name, result, v.Result)
		}

		if cnt != v.UpdateCount {
			t.Errorf("%s: update count = %d, want %d", v.Name, cnt, v.UpdateCount)
		}

		if v.ViewCache != nil {
			if !reflect.DeepEqual(ViewCache, v.ViewCache) {
				t.Errorf("%s: view cache = %v, want %v", v.Name, ViewCache, v.ViewCache)
			}
		}
		if v.TempViewList != nil {
			if !reflect.DeepEqual(filter.TempViews, v.TempViewList) {
				t.Errorf("%s: temporary views list = %v, want %v", v.Name, filter.TempViews, v.TempViewList)
			}
		}
	}
	ReleaseResources()
}

var renameColumnTests = []struct {
	Name         string
	Query        parser.RenameColumn
	Result       *FileInfo
	ViewCache    ViewMap
	TempViewList TemporaryViewScopes
	Error        string
}{
	{
		Name: "Rename Column",
		Query: parser.RenameColumn{
			Table: parser.Identifier{Literal: "table1"},
			Old:   parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
			New:   parser.Identifier{Literal: "newcolumn"},
		},
		Result: &FileInfo{
			Path:      GetTestFilePath("table1.csv"),
			Delimiter: ',',
			NoHeader:  false,
			Encoding:  text.UTF8,
			LineBreak: text.LF,
		},
		ViewCache: ViewMap{
			strings.ToUpper(GetTestFilePath("table1.csv")): &View{
				FileInfo: &FileInfo{
					Path:      GetTestFilePath("table1.csv"),
					Delimiter: ',',
					NoHeader:  false,
					Encoding:  text.UTF8,
					LineBreak: text.LF,
				},
				Header: NewHeader("table1", []string{"column1", "newcolumn"}),
				RecordSet: []Record{
					NewRecord([]value.Primary{
						value.NewString("1"),
						value.NewString("str1"),
					}),
					NewRecord([]value.Primary{
						value.NewString("2"),
						value.NewString("str2"),
					}),
					NewRecord([]value.Primary{
						value.NewString("3"),
						value.NewString("str3"),
					}),
				},
				ForUpdate: true,
			},
		},
	},
	{
		Name: "Rename Column For Temporary View",
		Query: parser.RenameColumn{
			Table: parser.Identifier{Literal: "tmpview"},
			Old:   parser.ColumnNumber{View: parser.Identifier{Literal: "tmpview"}, Number: value.NewInteger(2)},
			New:   parser.Identifier{Literal: "newcolumn"},
		},
		Result: &FileInfo{
			Path:        "tmpview",
			Delimiter:   ',',
			IsTemporary: true,
		},
		TempViewList: TemporaryViewScopes{
			ViewMap{
				"TMPVIEW": &View{
					Header: NewHeader("tmpview", []string{"column1", "newcolumn"}),
					RecordSet: []Record{
						NewRecord([]value.Primary{
							value.NewString("1"),
							value.NewString("str1"),
						}),
						NewRecord([]value.Primary{
							value.NewString("2"),
							value.NewString("str2"),
						}),
					},
					FileInfo: &FileInfo{
						Path:        "tmpview",
						Delimiter:   ',',
						IsTemporary: true,
					},
					ForUpdate: true,
				},
			},
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
	tf.Quiet = false

	filter := NewEmptyFilter()
	filter.TempViews = TemporaryViewScopes{
		ViewMap{
			"TMPVIEW": &View{
				Header: NewHeader("tmpview", []string{"column1", "column2"}),
				RecordSet: []Record{
					NewRecord([]value.Primary{
						value.NewString("1"),
						value.NewString("str1"),
					}),
					NewRecord([]value.Primary{
						value.NewString("2"),
						value.NewString("str2"),
					}),
				},
				FileInfo: &FileInfo{
					Path:        "tmpview",
					Delimiter:   ',',
					IsTemporary: true,
				},
			},
		},
	}

	for _, v := range renameColumnTests {
		ReleaseResources()
		result, err := RenameColumn(v.Query, filter)
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

		for _, v2 := range ViewCache {
			if v2.FileInfo.Handler != nil {
				if v2.FileInfo.Path != v2.FileInfo.Handler.Path() {
					t.Errorf("file pointer = %q, want %q for %q", v2.FileInfo.Handler.Path(), v2.FileInfo.Path, v.Name)
				}
				v2.FileInfo.Close()
				v2.FileInfo.Handler = nil
			}
		}

		if !reflect.DeepEqual(result, v.Result) {
			t.Errorf("%s: result = %v, want %v", v.Name, result, v.Result)
		}

		if v.ViewCache != nil {
			if !reflect.DeepEqual(ViewCache, v.ViewCache) {
				t.Errorf("%s: view cache = %v, want %v", v.Name, ViewCache, v.ViewCache)
			}
		}
		if v.TempViewList != nil {
			if !reflect.DeepEqual(filter.TempViews, v.TempViewList) {
				t.Errorf("%s: temporary views list = %v, want %v", v.Name, filter.TempViews, v.TempViewList)
			}
		}
	}
	ReleaseResources()
}

var setTableAttributeTests = []struct {
	Name   string
	Query  parser.SetTableAttribute
	Expect *FileInfo
	Error  string
}{
	{
		Name: "Set Delimiter to CSV",
		Query: parser.SetTableAttribute{
			Table:     parser.Identifier{Literal: "table1.csv"},
			Attribute: parser.Identifier{Literal: "delimiter"},
			Value:     parser.NewStringValue(";"),
		},
		Expect: &FileInfo{
			Path:      GetTestFilePath("table1.csv"),
			Delimiter: ';',
			Format:    cmd.CSV,
			Encoding:  text.UTF8,
			LineBreak: text.LF,
		},
	},
	{
		Name: "Set Delimiter to TSV",
		Query: parser.SetTableAttribute{
			Table:     parser.Identifier{Literal: "table1.csv"},
			Attribute: parser.Identifier{Literal: "delimiter"},
			Value:     parser.NewStringValue("\t"),
		},
		Expect: &FileInfo{
			Path:      GetTestFilePath("table1.csv"),
			Delimiter: '\t',
			Format:    cmd.TSV,
			Encoding:  text.UTF8,
			LineBreak: text.LF,
		},
	},
	{
		Name: "Set Delimiter to FIXED",
		Query: parser.SetTableAttribute{
			Table:     parser.Identifier{Literal: "table1.csv"},
			Attribute: parser.Identifier{Literal: "delimiter"},
			Value:     parser.NewStringValue("SPACES"),
		},
		Expect: &FileInfo{
			Path:      GetTestFilePath("table1.csv"),
			Delimiter: ',',
			Format:    cmd.FIXED,
			Encoding:  text.UTF8,
			LineBreak: text.LF,
		},
	},
	{
		Name: "Set Delimiter Error",
		Query: parser.SetTableAttribute{
			Table:     parser.Identifier{Literal: "table1.csv"},
			Attribute: parser.Identifier{Literal: "delimiter"},
			Value:     parser.NewStringValue("aa"),
		},
		Error: "[L:- C:-] delimiter must be one character, \"SPACES\" or JSON array of integers",
	},
	{
		Name: "Set Delimiter Not Allowed Value",
		Query: parser.SetTableAttribute{
			Table:     parser.Identifier{Literal: "table1.csv"},
			Attribute: parser.Identifier{Literal: "delimiter"},
			Value:     parser.NewNullValueFromString("null"),
		},
		Error: "[L:- C:-] null for delimiter is not allowed",
	},
	{
		Name: "Set Format to Text",
		Query: parser.SetTableAttribute{
			Table:     parser.Identifier{Literal: "table1.csv"},
			Attribute: parser.Identifier{Literal: "format"},
			Value:     parser.NewStringValue("text"),
		},
		Expect: &FileInfo{
			Path:      GetTestFilePath("table1.csv"),
			Delimiter: ',',
			Format:    cmd.TEXT,
			Encoding:  text.UTF8,
			LineBreak: text.LF,
		},
	},
	{
		Name: "Set Format to JSON",
		Query: parser.SetTableAttribute{
			Table:     parser.Identifier{Literal: "table1.csv"},
			Attribute: parser.Identifier{Literal: "format"},
			Value:     parser.NewStringValue("json"),
		},
		Expect: &FileInfo{
			Path:      GetTestFilePath("table1.csv"),
			Delimiter: ',',
			Format:    cmd.JSON,
			Encoding:  text.UTF8,
			LineBreak: text.LF,
		},
	},
	{
		Name: "Set Format to TSV",
		Query: parser.SetTableAttribute{
			Table:     parser.Identifier{Literal: "table1.csv"},
			Attribute: parser.Identifier{Literal: "format"},
			Value:     parser.NewStringValue("tsv"),
		},
		Expect: &FileInfo{
			Path:      GetTestFilePath("table1.csv"),
			Delimiter: '\t',
			Format:    cmd.TSV,
			Encoding:  text.UTF8,
			LineBreak: text.LF,
		},
	},
	{
		Name: "Set Format Error",
		Query: parser.SetTableAttribute{
			Table:     parser.Identifier{Literal: "table1.csv"},
			Attribute: parser.Identifier{Literal: "format"},
			Value:     parser.NewStringValue("invalid"),
		},
		Error: "[L:- C:-] format must be one of CSV|TSV|FIXED|JSON|GFM|ORG|TEXT|JSONH|JSONA",
	},
	{
		Name: "Set Encoding to SJIS",
		Query: parser.SetTableAttribute{
			Table:     parser.Identifier{Literal: "table1.csv"},
			Attribute: parser.Identifier{Literal: "encoding"},
			Value:     parser.NewStringValue("sjis"),
		},
		Expect: &FileInfo{
			Path:      GetTestFilePath("table1.csv"),
			Delimiter: ',',
			Format:    cmd.CSV,
			Encoding:  text.SJIS,
			LineBreak: text.LF,
		},
	},
	{
		Name: "Set Encoding to SJIS with Identifier",
		Query: parser.SetTableAttribute{
			Table:     parser.Identifier{Literal: "table1.csv"},
			Attribute: parser.Identifier{Literal: "encoding"},
			Value:     parser.Identifier{Literal: "sjis"},
		},
		Expect: &FileInfo{
			Path:      GetTestFilePath("table1.csv"),
			Delimiter: ',',
			Format:    cmd.CSV,
			Encoding:  text.SJIS,
			LineBreak: text.LF,
		},
	},
	{
		Name: "Set Encoding Error",
		Query: parser.SetTableAttribute{
			Table:     parser.Identifier{Literal: "table1.csv"},
			Attribute: parser.Identifier{Literal: "encoding"},
			Value:     parser.NewStringValue("invalid"),
		},
		Error: "[L:- C:-] encoding must be one of UTF8|SJIS",
	},
	{
		Name: "Set Encoding Error in JSON Format",
		Query: parser.SetTableAttribute{
			Table:     parser.Identifier{Literal: "table.json"},
			Attribute: parser.Identifier{Literal: "encoding"},
			Value:     parser.NewStringValue("sjis"),
		},
		Error: "[L:- C:-] json format is supported only UTF8",
	},
	{
		Name: "Set LineBreak to CRLF",
		Query: parser.SetTableAttribute{
			Table:     parser.Identifier{Literal: "table1.csv"},
			Attribute: parser.Identifier{Literal: "line_break"},
			Value:     parser.NewStringValue("crlf"),
		},
		Expect: &FileInfo{
			Path:      GetTestFilePath("table1.csv"),
			Delimiter: ',',
			Format:    cmd.CSV,
			Encoding:  text.UTF8,
			LineBreak: text.CRLF,
		},
	},
	{
		Name: "Set LineBreak Error",
		Query: parser.SetTableAttribute{
			Table:     parser.Identifier{Literal: "table1.csv"},
			Attribute: parser.Identifier{Literal: "line_break"},
			Value:     parser.NewStringValue("invalid"),
		},
		Error: "[L:- C:-] line-break must be one of CRLF|LF|CR",
	},
	{
		Name: "Set NoHeader to true",
		Query: parser.SetTableAttribute{
			Table:     parser.Identifier{Literal: "table1.csv"},
			Attribute: parser.Identifier{Literal: "header"},
			Value:     parser.NewTernaryValueFromString("false"),
		},
		Expect: &FileInfo{
			Path:      GetTestFilePath("table1.csv"),
			Delimiter: ',',
			Format:    cmd.CSV,
			Encoding:  text.UTF8,
			LineBreak: text.LF,
			NoHeader:  true,
		},
	},
	{
		Name: "Set NoHeader Not Allowed Value",
		Query: parser.SetTableAttribute{
			Table:     parser.Identifier{Literal: "table1.csv"},
			Attribute: parser.Identifier{Literal: "header"},
			Value:     parser.NewNullValueFromString("null"),
		},
		Error: "[L:- C:-] null for header is not allowed",
	},
	{
		Name: "Set EncloseAll to true",
		Query: parser.SetTableAttribute{
			Table:     parser.Identifier{Literal: "table1.csv"},
			Attribute: parser.Identifier{Literal: "enclose_all"},
			Value:     parser.NewStringValue("true"),
		},
		Expect: &FileInfo{
			Path:       GetTestFilePath("table1.csv"),
			Delimiter:  ',',
			Format:     cmd.CSV,
			Encoding:   text.UTF8,
			LineBreak:  text.LF,
			EncloseAll: true,
		},
	},
	{
		Name: "Set JsonEscape to HEX",
		Query: parser.SetTableAttribute{
			Table:     parser.Identifier{Literal: "table.json"},
			Attribute: parser.Identifier{Literal: "json_escape"},
			Value:     parser.NewStringValue("hex"),
		},
		Expect: &FileInfo{
			Path:        GetTestFilePath("table.json"),
			Delimiter:   ',',
			Format:      cmd.JSON,
			Encoding:    text.UTF8,
			LineBreak:   text.LF,
			JsonEscape:  json.HexDigits,
			PrettyPrint: false,
		},
	},
	{
		Name: "Set JsonEscape Error",
		Query: parser.SetTableAttribute{
			Table:     parser.Identifier{Literal: "table.json"},
			Attribute: parser.Identifier{Literal: "json_escape"},
			Value:     parser.NewStringValue("invalid"),
		},
		Error: "[L:- C:-] json-escape must be one of BACKSLASH|HEX|HEXALL",
	},
	{
		Name: "Set PrettyPring to true",
		Query: parser.SetTableAttribute{
			Table:     parser.Identifier{Literal: "table.json"},
			Attribute: parser.Identifier{Literal: "pretty_print"},
			Value:     parser.NewTernaryValueFromString("true"),
		},
		Expect: &FileInfo{
			Path:        GetTestFilePath("table.json"),
			Delimiter:   ',',
			Format:      cmd.JSON,
			Encoding:    text.UTF8,
			LineBreak:   text.LF,
			PrettyPrint: true,
		},
	},
	{
		Name: "Not Exist Table Error",
		Query: parser.SetTableAttribute{
			Table:     parser.Identifier{Literal: "notexist.csv"},
			Attribute: parser.Identifier{Literal: "delimiter"},
			Value:     parser.NewStringValue(","),
		},
		Error: "[L:- C:-] file notexist.csv does not exist",
	},
	{
		Name: "Temporary View Error",
		Query: parser.SetTableAttribute{
			Table:     parser.Identifier{Literal: "tmpview"},
			Attribute: parser.Identifier{Literal: "delimiter"},
			Value:     parser.NewStringValue(","),
		},
		Error: "[L:- C:-] tmpview is not a table that has attributes",
	},
	{
		Name: "Value Evaluation Error",
		Query: parser.SetTableAttribute{
			Table:     parser.Identifier{Literal: "table1.csv"},
			Attribute: parser.Identifier{Literal: "delimiter"},
			Value:     parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
	{
		Name: "Not Exist Attribute Error",
		Query: parser.SetTableAttribute{
			Table:     parser.Identifier{Literal: "table1.csv"},
			Attribute: parser.Identifier{Literal: "notexist"},
			Value:     parser.NewStringValue(","),
		},
		Error: "[L:- C:-] table attribute notexist does not exist",
	},
}

func TestSetTableAttribute(t *testing.T) {
	tf := cmd.GetFlags()
	tf.Repository = TestDir
	tf.Quiet = false

	filter := NewEmptyFilter()
	filter.TempViews = TemporaryViewScopes{
		ViewMap{
			"TMPVIEW": &View{
				Header: NewHeader("tmpview", []string{"column1", "column2"}),
				RecordSet: []Record{
					NewRecord([]value.Primary{
						value.NewString("1"),
						value.NewString("str1"),
					}),
					NewRecord([]value.Primary{
						value.NewString("2"),
						value.NewString("str2"),
					}),
				},
				FileInfo: &FileInfo{
					Path:        "tmpview",
					Delimiter:   ',',
					IsTemporary: true,
				},
			},
		},
	}

	for _, v := range setTableAttributeTests {
		ReleaseResources()

		_, _, err := SetTableAttribute(v.Query, filter)
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

		for _, v2 := range ViewCache {
			if v2.FileInfo.Handler != nil {
				if v2.FileInfo.Path != v2.FileInfo.Handler.Path() {
					t.Errorf("file pointer = %q, want %q for %q", v2.FileInfo.Handler.Path(), v2.FileInfo.Path, v.Name)
				}
				v2.FileInfo.Close()
				v2.FileInfo.Handler = nil
			}
		}

		view := NewView()
		view.LoadFromTableIdentifier(v.Query.Table, filter.CreateNode())

		if !reflect.DeepEqual(view.FileInfo, v.Expect) {
			t.Errorf("%s: result = %v, want %v", v.Name, view.FileInfo, v.Expect)
		}

		_, _, err = SetTableAttribute(v.Query, filter)
		if err == nil {
			t.Errorf("%s: no error, want TableAttributeUnchangedError for duplicate set", v.Name)
		} else if _, ok := err.(*TableAttributeUnchangedError); !ok {
			t.Errorf("%s: error = %T, want TableAttributeUnchangedError for duplicate set", v.Name, err)
		}
	}
	ReleaseResources()
}

func TestCommit(t *testing.T) {
	cmd.GetFlags().SetQuiet(false)

	ch, _ := file.NewHandlerForCreate(GetTestFilePath("create_file.csv"))
	uh, _ := file.NewHandlerForUpdate(GetTestFilePath("updated_file_1.csv"))

	ViewCache = ViewMap{
		strings.ToUpper(GetTestFilePath("created_file.csv")): &View{
			Header:    NewHeader("created_file", []string{"column1", "column2"}),
			RecordSet: RecordSet{},
			FileInfo: &FileInfo{
				Path:    GetTestFilePath("created_file.csv"),
				Handler: ch,
			},
		},
		strings.ToUpper(GetTestFilePath("updated_file_1.csv")): &View{
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
				Path:    GetTestFilePath("updated_file_1.csv"),
				Handler: uh,
			},
		},
	}

	UncommittedViews = &UncommittedViewMap{
		Created: map[string]*FileInfo{
			strings.ToUpper(GetTestFilePath("created_file.csv")): {
				Path:    GetTestFilePath("created_file.csv"),
				Handler: ch,
			},
		},
		Updated: map[string]*FileInfo{
			strings.ToUpper(GetTestFilePath("updated_file_1.csv")): {
				Path:    GetTestFilePath("updated_file_1.csv"),
				Handler: uh,
			},
		},
	}

	expect := fmt.Sprintf("Commit: file %q is created.\nCommit: file %q is updated.\n", GetTestFilePath("created_file.csv"), GetTestFilePath("updated_file_1.csv"))

	oldStdout := Stdout
	r, w, _ := os.Pipe()
	Stdout = w

	Commit(parser.TransactionControl{Token: parser.COMMIT}, NewEmptyFilter())

	w.Close()
	Stdout = oldStdout
	log, _ := ioutil.ReadAll(r)

	if string(log) != expect {
		t.Errorf("Commit: log = %q, want %q", string(log), expect)
	}
}

func TestRollback(t *testing.T) {
	cmd.GetFlags().SetQuiet(false)

	UncommittedViews = &UncommittedViewMap{
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

	oldStdout := Stdout
	r, w, _ := os.Pipe()
	Stdout = w

	Rollback(nil, NewEmptyFilter())

	w.Close()
	Stdout = oldStdout
	log, _ := ioutil.ReadAll(r)

	if string(log) != expect {
		t.Errorf("Rollback: log = %q, want %q", string(log), expect)
	}
}
