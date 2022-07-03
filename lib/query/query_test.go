package query

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/mithrandie/csvq/lib/option"
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
	ResultScopes  *ReferenceScope
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
		ResultScopes: GenerateReferenceScope([]map[string]map[string]interface{}{
			{
				scopeNameVariables: {
					"var1": value.NewString("1"),
					"var2": value.NewString("str1"),
				},
			},
		}, nil, time.Time{}, nil),
	},
	{
		Name:    "Fetch Cursor Second Time",
		CurName: parser.Identifier{Literal: "cur"},
		Variables: []parser.Variable{
			{Name: "var1"},
			{Name: "var2"},
		},
		Success: true,
		ResultScopes: GenerateReferenceScope([]map[string]map[string]interface{}{
			{
				scopeNameVariables: {
					"var1": value.NewString("2"),
					"var2": value.NewString("str2"),
				},
			},
		}, nil, time.Time{}, nil),
	},
	{
		Name:    "Fetch Cursor Third Time",
		CurName: parser.Identifier{Literal: "cur"},
		Variables: []parser.Variable{
			{Name: "var1"},
			{Name: "var2"},
		},
		Success: true,
		ResultScopes: GenerateReferenceScope([]map[string]map[string]interface{}{
			{
				scopeNameVariables: {
					"var1": value.NewString("3"),
					"var2": value.NewString("str3"),
				},
			},
		}, nil, time.Time{}, nil),
	},
	{
		Name:    "Fetch Cursor Forth Time",
		CurName: parser.Identifier{Literal: "cur"},
		Variables: []parser.Variable{
			{Name: "var1"},
			{Name: "var2"},
		},
		Success: false,
		ResultScopes: GenerateReferenceScope([]map[string]map[string]interface{}{
			{
				scopeNameVariables: {
					"var1": value.NewString("3"),
					"var2": value.NewString("str3"),
				},
			},
		}, nil, time.Time{}, nil),
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
		ResultScopes: GenerateReferenceScope([]map[string]map[string]interface{}{
			{
				scopeNameVariables: {
					"var1": value.NewString("2"),
					"var2": value.NewString("str2"),
				},
			},
		}, nil, time.Time{}, nil),
	},
	{
		Name:    "Fetch Cursor Fetch Error",
		CurName: parser.Identifier{Literal: "notexist"},
		Variables: []parser.Variable{
			{Name: "var1"},
			{Name: "var2"},
		},
		Error: "cursor notexist is undeclared",
	},
	{
		Name:    "Fetch Cursor Not Match Number Error",
		CurName: parser.Identifier{Literal: "cur2"},
		Variables: []parser.Variable{
			{Name: "var1"},
		},
		Error: "fetching from cursor cur2 returns 2 values",
	},
	{
		Name:    "Fetch Cursor Substitution Error",
		CurName: parser.Identifier{Literal: "cur2"},
		Variables: []parser.Variable{
			{Name: "var1"},
			{Name: "notexist"},
		},
		Error: "variable @notexist is undeclared",
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
		Error: "field notexist does not exist",
	},
	{
		Name:    "Fetch Cursor Number Not Integer Error",
		CurName: parser.Identifier{Literal: "cur"},
		FetchPosition: parser.FetchPosition{
			Position: parser.Token{Token: parser.ABSOLUTE, Literal: "absolute"},
			Number:   parser.NewNullValue(),
		},
		Variables: []parser.Variable{
			{Name: "var1"},
			{Name: "var2"},
		},
		Error: "fetching position NULL is not an integer value",
	},
}

func TestFetchCursor(t *testing.T) {
	defer func() {
		_ = TestTx.CachedViews.Clean(TestTx.FileContainer)
		initFlag(TestTx.Flags)
	}()

	TestTx.Flags.Repository = TestDir

	scope := GenerateReferenceScope([]map[string]map[string]interface{}{
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
				"CUR2": &Cursor{
					query: selectQueryForCursorTest,
					mtx:   &sync.Mutex{},
				},
			},
		},
	}, nil, time.Time{}, nil)

	ctx := context.Background()
	_ = TestTx.CachedViews.Clean(TestTx.FileContainer)
	_ = scope.OpenCursor(ctx, parser.Identifier{Literal: "cur"}, nil)
	_ = TestTx.CachedViews.Clean(TestTx.FileContainer)
	_ = scope.OpenCursor(ctx, parser.Identifier{Literal: "cur2"}, nil)

	for _, v := range fetchCursorTests {
		success, err := FetchCursor(ctx, scope, v.CurName, v.FetchPosition, v.Variables)
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

		if !SyncMapEqual(scope.Blocks[0].Variables, v.ResultScopes.Blocks[0].Variables) {
			t.Errorf("%s: variables = %v, want %v", v.Name, scope.Blocks, v.ResultScopes.Blocks)
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
		Result: GenerateViewMap([]*View{
			{
				FileInfo: &FileInfo{
					Path:                  "tbl",
					ViewType:              ViewTypeTemporaryTable,
					restorePointHeader:    NewHeader("tbl", []string{"column1", "column2"}),
					restorePointRecordSet: RecordSet{},
				},
				Header:    NewHeader("tbl", []string{"column1", "column2"}),
				RecordSet: RecordSet{},
			},
		}),
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
		Error: "field name column1 is a duplicate",
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
		Result: GenerateViewMap([]*View{
			{
				FileInfo: &FileInfo{
					Path:               "tbl",
					ViewType:           ViewTypeTemporaryTable,
					restorePointHeader: NewHeader("tbl", []string{"column1", "column2"}),
					restorePointRecordSet: RecordSet{
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
		}),
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
		Error: "field notexist does not exist",
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
		Error: "select query should return exactly 1 field for view tbl",
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
		Error: "field name column1 is a duplicate",
	},
	{
		Name: "Declare View Redeclaration Error",
		ViewMap: GenerateViewMap([]*View{
			{
				FileInfo: &FileInfo{
					Path:     "tbl",
					ViewType: ViewTypeTemporaryTable,
				},
			},
		}),
		Expr: parser.ViewDeclaration{
			View: parser.Identifier{Literal: "tbl"},
			Fields: []parser.QueryExpression{
				parser.Identifier{Literal: "column1"},
				parser.Identifier{Literal: "column2"},
			},
		},
		Error: "view tbl is redeclared",
	},
}

func TestDeclareView(t *testing.T) {
	scope := NewReferenceScope(TestTx)
	ctx := context.Background()

	for _, v := range declareViewTests {
		if v.ViewMap.SyncMap == nil {
			scope.Blocks[0].TemporaryTables = NewViewMap()
		} else {
			scope.Blocks[0].TemporaryTables = v.ViewMap
		}

		err := DeclareView(ctx, scope, v.Expr)
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
		if !SyncMapEqual(scope.Blocks[0].TemporaryTables, v.Result) {
			t.Errorf("%s: view cache = %v, want %v", v.Name, TestTx.CachedViews, v.Result)
		}
	}
}

var selectTests = []struct {
	Name         string
	Query        parser.SelectQuery
	Result       *View
	SetVariables map[parser.Variable]value.Primary
	Error        string
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
						Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: "<"},
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
						Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: ">"},
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
				OffsetClause: parser.OffsetClause{
					Value: parser.NewIntegerValue(0),
				},
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
					Column:      "COUNT(*)",
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
				Type:  parser.Token{Token: parser.FETCH},
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
		Name: "Separate scopes on both sides of a set operator",
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
				All:      parser.Token{Token: parser.ALL, Literal: "all"},
				RHS: parser.SelectEntity{
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
		Error: "result set to be combined should contain exactly 2 fields",
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
		Error: "field notexist does not exist",
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
		Error: "field notexist does not exist",
	},
	{
		Name: "Inline Tables",
		Query: parser.SelectQuery{
			WithClause: parser.WithClause{
				InlineTables: []parser.QueryExpression{
					parser.InlineTable{
						Name: parser.Identifier{Literal: "it"},
						Fields: []parser.QueryExpression{
							parser.Identifier{Literal: "c1"},
						},
						Query: parser.SelectQuery{
							SelectEntity: parser.SelectEntity{
								SelectClause: parser.SelectClause{
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
				InlineTables: []parser.QueryExpression{
					parser.InlineTable{
						Name: parser.Identifier{Literal: "it"},
						Fields: []parser.QueryExpression{
							parser.Identifier{Literal: "c1"},
						},
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
		Error: "select query should return exactly 1 field for inline table it",
	},
	{
		Name: "Inline Tables Recursion",
		Query: parser.SelectQuery{
			WithClause: parser.WithClause{
				InlineTables: []parser.QueryExpression{
					parser.InlineTable{
						Recursive: parser.Token{Token: parser.RECURSIVE, Literal: "recursive"},
						Name:      parser.Identifier{Literal: "it"},
						Fields: []parser.QueryExpression{
							parser.Identifier{Literal: "n"},
						},
						Query: parser.SelectQuery{
							SelectEntity: parser.SelectSet{
								LHS: parser.SelectEntity{
									SelectClause: parser.SelectClause{
										Fields: []parser.QueryExpression{
											parser.Field{Object: parser.NewIntegerValueFromString("1")},
										},
									},
								},
								Operator: parser.Token{Token: parser.UNION, Literal: "union"},
								RHS: parser.SelectEntity{
									SelectClause: parser.SelectClause{
										Fields: []parser.QueryExpression{
											parser.Field{
												Object: parser.Arithmetic{
													LHS:      parser.FieldReference{Column: parser.Identifier{Literal: "n"}},
													RHS:      parser.NewIntegerValueFromString("1"),
													Operator: parser.Token{Token: '+', Literal: "+"},
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
											Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: "<"},
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
				InlineTables: []parser.QueryExpression{
					parser.InlineTable{
						Recursive: parser.Token{Token: parser.RECURSIVE, Literal: "recursive"},
						Name:      parser.Identifier{Literal: "it"},
						Fields: []parser.QueryExpression{
							parser.Identifier{Literal: "n"},
						},
						Query: parser.SelectQuery{
							SelectEntity: parser.SelectSet{
								LHS: parser.SelectEntity{
									SelectClause: parser.SelectClause{
										Fields: []parser.QueryExpression{
											parser.Field{Object: parser.NewIntegerValueFromString("1")},
										},
									},
								},
								Operator: parser.Token{Token: parser.UNION, Literal: "union"},
								RHS: parser.SelectEntity{
									SelectClause: parser.SelectClause{
										Fields: []parser.QueryExpression{
											parser.Field{
												Object: parser.Arithmetic{
													LHS:      parser.FieldReference{Column: parser.Identifier{Literal: "n"}},
													RHS:      parser.NewIntegerValueFromString("1"),
													Operator: parser.Token{Token: '+', Literal: "+"},
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
											Operator: parser.Token{Token: '<', Literal: "<"},
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
		Error: "result set to be combined should contain exactly 1 field",
	},
	{
		Name: "Inline Tables Recursion Recursion Limit Exceeded Error",
		Query: parser.SelectQuery{
			WithClause: parser.WithClause{
				InlineTables: []parser.QueryExpression{
					parser.InlineTable{
						Recursive: parser.Token{Token: parser.RECURSIVE, Literal: "recursive"},
						Name:      parser.Identifier{Literal: "it"},
						Fields: []parser.QueryExpression{
							parser.Identifier{Literal: "n"},
						},
						Query: parser.SelectQuery{
							SelectEntity: parser.SelectSet{
								LHS: parser.SelectEntity{
									SelectClause: parser.SelectClause{
										Fields: []parser.QueryExpression{
											parser.Field{Object: parser.NewIntegerValueFromString("1")},
										},
									},
								},
								Operator: parser.Token{Token: parser.UNION, Literal: "union"},
								RHS: parser.SelectEntity{
									SelectClause: parser.SelectClause{
										Fields: []parser.QueryExpression{
											parser.Field{
												Object: parser.Arithmetic{
													LHS:      parser.FieldReference{Column: parser.Identifier{Literal: "n"}},
													RHS:      parser.NewIntegerValueFromString("1"),
													Operator: parser.Token{Token: '+', Literal: "+"},
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
											RHS:      parser.NewIntegerValueFromString("10"),
											Operator: parser.Token{Token: '<', Literal: "<"},
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
		Error: "iteration of recursive query exceeded the limit 5",
	},
	{
		Name: "Select Into Variables",
		Query: parser.SelectQuery{
			SelectEntity: parser.SelectEntity{
				SelectClause: parser.SelectClause{
					Fields: []parser.QueryExpression{
						parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
						parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
					},
				},
				IntoClause: parser.IntoClause{
					Variables: []parser.Variable{
						{Name: "var1"},
						{Name: "var2"},
					},
				},
				FromClause: parser.FromClause{
					Tables: []parser.QueryExpression{
						parser.Table{Object: parser.Identifier{Literal: "table1"}},
					},
				},
			},
			LimitClause: parser.LimitClause{
				Type:  parser.Token{Token: parser.LIMIT},
				Value: parser.NewIntegerValueFromString("1"),
			},
		},
		SetVariables: map[parser.Variable]value.Primary{
			parser.Variable{Name: "var1"}: value.NewString("1"),
			parser.Variable{Name: "var2"}: value.NewString("str1"),
		},
	},
	{
		Name: "Select Into Variables Empty Result Set",
		Query: parser.SelectQuery{
			SelectEntity: parser.SelectEntity{
				SelectClause: parser.SelectClause{
					Fields: []parser.QueryExpression{
						parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
						parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
					},
				},
				IntoClause: parser.IntoClause{
					Variables: []parser.Variable{
						{Name: "var1"},
						{Name: "var2"},
					},
				},
				FromClause: parser.FromClause{
					Tables: []parser.QueryExpression{
						parser.Table{Object: parser.Identifier{Literal: "table1"}},
					},
				},
				WhereClause: parser.WhereClause{
					Filter: parser.NewTernaryValueFromString("false"),
				},
			},
			LimitClause: parser.LimitClause{
				Value: parser.NewIntegerValueFromString("1"),
			},
		},
		SetVariables: map[parser.Variable]value.Primary{
			parser.Variable{Name: "var1"}: value.NewNull(),
			parser.Variable{Name: "var2"}: value.NewNull(),
		},
	},
	{
		Name: "Select Into Variables Too Many Records",
		Query: parser.SelectQuery{
			SelectEntity: parser.SelectEntity{
				SelectClause: parser.SelectClause{
					Fields: []parser.QueryExpression{
						parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
						parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
					},
				},
				IntoClause: parser.IntoClause{
					Variables: []parser.Variable{
						{Name: "var1"},
						{Name: "var2"},
					},
				},
				FromClause: parser.FromClause{
					Tables: []parser.QueryExpression{
						parser.Table{Object: parser.Identifier{Literal: "table1"}},
					},
				},
			},
		},
		Error: "select into query returns too many records, should return only one record",
	},
	{
		Name: "Select Into Variables Field Length Not Match",
		Query: parser.SelectQuery{
			SelectEntity: parser.SelectEntity{
				SelectClause: parser.SelectClause{
					Fields: []parser.QueryExpression{
						parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
						parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
					},
				},
				IntoClause: parser.IntoClause{
					Variables: []parser.Variable{
						{Name: "var1"},
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
		Error: "select into query should return exactly 1 field",
	},
	{
		Name: "Select Into Variables Undeclared Variable",
		Query: parser.SelectQuery{
			SelectEntity: parser.SelectEntity{
				SelectClause: parser.SelectClause{
					Fields: []parser.QueryExpression{
						parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
						parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
					},
				},
				IntoClause: parser.IntoClause{
					Variables: []parser.Variable{
						{Name: "var1"},
						{Name: "undeclared"},
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
		Error: "variable @undeclared is undeclared",
	},
}

func TestSelect(t *testing.T) {
	defer func() {
		_ = TestTx.CachedViews.Clean(TestTx.FileContainer)
		initFlag(TestTx.Flags)
	}()

	TestTx.Flags.Repository = TestDir

	scope := NewReferenceScope(TestTx)
	ctx := context.Background()
	_ = scope.DeclareVariable(ctx, parser.VariableDeclaration{Assignments: []parser.VariableAssignment{
		{Variable: parser.Variable{Name: "var1"}},
		{Variable: parser.Variable{Name: "var2"}, Value: parser.NewIntegerValueFromString("2")},
	}})

	for _, v := range selectTests {
		_ = TestTx.CachedViews.Clean(TestTx.FileContainer)
		result, err := Select(ctx, scope, v.Query)
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
		if v.Result != nil {
			if !reflect.DeepEqual(result, v.Result) {
				t.Errorf("%s: result = %v, want %v", v.Name, result, v.Result)
			}
		}
		if 0 < len(v.SetVariables) {
			for variable, expectValue := range v.SetVariables {
				val, _ := scope.GetVariable(variable)
				if !reflect.DeepEqual(val, expectValue) {
					t.Errorf("%s: variable %s = %v, want %v", v.Name, variable, val, expectValue)
				}
			}
		}
	}
}

var insertTests = []struct {
	Name         string
	Query        parser.InsertQuery
	ResultFile   *FileInfo
	UpdateCount  int
	ViewCache    ViewMap
	ResultScopes *ReferenceScope
	Error        string
}{
	{
		Name: "Insert Query",
		Query: parser.InsertQuery{
			WithClause: parser.WithClause{
				InlineTables: []parser.QueryExpression{
					parser.InlineTable{
						Name: parser.Identifier{Literal: "it"},
						Fields: []parser.QueryExpression{
							parser.Identifier{Literal: "c1"},
						},
						Query: parser.SelectQuery{
							SelectEntity: parser.SelectEntity{
								SelectClause: parser.SelectClause{
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
			ForUpdate: true,
		},
		UpdateCount: 2,
		ViewCache: GenerateViewMap([]*View{
			{
				FileInfo: &FileInfo{
					Path:      GetTestFilePath("table1.csv"),
					Delimiter: ',',
					NoHeader:  false,
					Encoding:  text.UTF8,
					LineBreak: text.LF,
					ForUpdate: true,
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
			},
		}),
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
			Path:      "tmpview",
			Delimiter: ',',
			ViewType:  ViewTypeTemporaryTable,
		},
		UpdateCount: 2,
		ResultScopes: GenerateReferenceScope([]map[string]map[string]interface{}{
			{
				scopeNameTempTables: {
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
							Path:      "tmpview",
							Delimiter: ',',
							ViewType:  ViewTypeTemporaryTable,
						},
					},
				},
			},
		}, nil, time.Time{}, nil),
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
			ForUpdate: true,
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
		Error: "file notexist does not exist",
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
		Error: "field notexist does not exist",
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
			ForUpdate: true,
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
		Error: "select query should return exactly 1 field",
	},
}

func TestInsert(t *testing.T) {
	defer func() {
		_ = TestTx.ReleaseResources()
		_ = TestTx.CachedViews.Clean(TestTx.FileContainer)
		initFlag(TestTx.Flags)
	}()

	TestTx.Flags.Repository = TestDir
	TestTx.Flags.Quiet = false

	scope := GenerateReferenceScope([]map[string]map[string]interface{}{
		{
			scopeNameTempTables: {
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
						Path:      "tmpview",
						Delimiter: ',',
						ViewType:  ViewTypeTemporaryTable,
					},
				},
			},
		},
	}, nil, time.Time{}, nil)

	ctx := context.Background()
	for _, v := range insertTests {
		_ = TestTx.ReleaseResources()
		result, cnt, err := Insert(ctx, scope, v.Query)
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

		TestTx.CachedViews.Range(func(key, value interface{}) bool {
			view := value.(*View)
			if view.FileInfo.Handler != nil {
				if view.FileInfo.Path != view.FileInfo.Handler.Path() {
					t.Errorf("file pointer = %q, want %q for %q", view.FileInfo.Handler.Path(), view.FileInfo.Path, v.Name)
				}
				_ = TestTx.FileContainer.Close(view.FileInfo.Handler)
				view.FileInfo.Handler = nil
			}
			return true
		})

		if !reflect.DeepEqual(result, v.ResultFile) {
			t.Errorf("%s: fileinfo = %v, want %v", v.Name, result, v.ResultFile)
		}

		if !reflect.DeepEqual(cnt, v.UpdateCount) {
			t.Errorf("%s: update count = %d, want %d", v.Name, cnt, v.UpdateCount)
		}

		if v.ViewCache.SyncMap != nil {
			if !SyncMapEqual(TestTx.CachedViews, v.ViewCache) {
				t.Errorf("%s: view cache = %v, want %v", v.Name, TestTx.CachedViews, v.ViewCache)
			}
		}
		if v.ResultScopes != nil {
			if !SyncMapEqual(scope.Blocks[0].TemporaryTables, v.ResultScopes.Blocks[0].TemporaryTables) {
				t.Errorf("%s: temporary views list = %v, want %v", v.Name, scope.Blocks[0].TemporaryTables, v.ResultScopes.Blocks[0].TemporaryTables)
			}
		}
	}
}

var updateTests = []struct {
	Name         string
	Query        parser.UpdateQuery
	ResultFiles  []*FileInfo
	UpdateCounts []int
	ViewCache    ViewMap
	ResultScopes *ReferenceScope
	Error        string
}{
	{
		Name: "Update Query",
		Query: parser.UpdateQuery{
			WithClause: parser.WithClause{
				InlineTables: []parser.QueryExpression{
					parser.InlineTable{
						Name: parser.Identifier{Literal: "it"},
						Fields: []parser.QueryExpression{
							parser.Identifier{Literal: "c1"},
						},
						Query: parser.SelectQuery{
							SelectEntity: parser.SelectEntity{
								SelectClause: parser.SelectClause{
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
					Operator: parser.Token{Token: '=', Literal: "="},
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
				ForUpdate: true,
			},
		},
		UpdateCounts: []int{1},
		ViewCache: GenerateViewMap([]*View{
			{
				FileInfo: &FileInfo{
					Path:      GetTestFilePath("table1.csv"),
					Delimiter: ',',
					NoHeader:  false,
					Encoding:  text.UTF8,
					LineBreak: text.LF,
					ForUpdate: true,
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
			},
		}),
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
				Path:      "tmpview",
				Delimiter: ',',
				ViewType:  ViewTypeTemporaryTable,
			},
		},
		UpdateCounts: []int{2},
		ResultScopes: GenerateReferenceScope([]map[string]map[string]interface{}{
			{
				scopeNameTempTables: {
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
							Path:      "tmpview",
							Delimiter: ',',
							ViewType:  ViewTypeTemporaryTable,
						},
					},
				},
			},
		}, nil, time.Time{}, nil),
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
								Operator: parser.Token{Token: '=', Literal: "="},
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
				ForUpdate: true,
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
					Operator: parser.Token{Token: '=', Literal: "="},
				},
			},
		},
		Error: "file notexist does not exist",
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
					Operator: parser.Token{Token: '=', Literal: "="},
				},
			},
		},
		Error: "field notexist does not exist",
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
								Operator: parser.Token{Token: '=', Literal: "="},
							},
						},
					}},
				},
			},
		},
		Error: "table notexist is not loaded",
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
								Operator: parser.Token{Token: '=', Literal: "="},
							},
						},
					}},
				},
			},
		},
		Error: "field t1.column2 does not exist in the tables to update",
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
					Operator: parser.Token{Token: '=', Literal: "="},
				},
			},
		},
		Error: "field notexist does not exist",
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
					Operator: parser.Token{Token: '=', Literal: "="},
				},
			},
		},
		Error: "field notexist does not exist",
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
		Error: "value column4 to set in the field column2 is ambiguous",
	},
}

func TestUpdate(t *testing.T) {
	defer func() {
		_ = TestTx.ReleaseResources()
		_ = TestTx.CachedViews.Clean(TestTx.FileContainer)
		initFlag(TestTx.Flags)
	}()

	TestTx.Flags.Repository = TestDir
	TestTx.Flags.Quiet = false

	scope := GenerateReferenceScope([]map[string]map[string]interface{}{
		{
			scopeNameTempTables: {
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
						Path:      "tmpview",
						Delimiter: ',',
						ViewType:  ViewTypeTemporaryTable,
					},
				},
			},
		},
	}, nil, time.Time{}, nil)

	ctx := context.Background()
	for _, v := range updateTests {
		_ = TestTx.ReleaseResources()
		files, cnt, err := Update(ctx, scope, v.Query)
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

		TestTx.CachedViews.Range(func(key, value interface{}) bool {
			view := value.(*View)
			if view.FileInfo.Handler != nil {
				if view.FileInfo.Path != view.FileInfo.Handler.Path() {
					t.Errorf("file pointer = %q, want %q for %q", view.FileInfo.Handler.Path(), view.FileInfo.Path, v.Name)
				}
				_ = TestTx.FileContainer.Close(view.FileInfo.Handler)
				view.FileInfo.Handler = nil
			}
			return true
		})

		if !reflect.DeepEqual(files, v.ResultFiles) {
			t.Errorf("%s: fileinfo = %v, want %v", v.Name, files, v.ResultFiles)
		}

		if !reflect.DeepEqual(cnt, v.UpdateCounts) {
			t.Errorf("%s: update count = %v, want %v", v.Name, cnt, v.UpdateCounts)
		}

		if v.ViewCache.SyncMap != nil {
			if !SyncMapEqual(TestTx.CachedViews, v.ViewCache) {
				t.Errorf("%s: view cache = %v, want %v", v.Name, TestTx.CachedViews, v.ViewCache)
			}
		}
		if v.ResultScopes != nil {
			if !SyncMapEqual(scope.Blocks[0].TemporaryTables, v.ResultScopes.Blocks[0].TemporaryTables) {
				t.Errorf("%s: temporary views list = %v, want %v", v.Name, scope.Blocks[0].TemporaryTables, v.ResultScopes.Blocks[0].TemporaryTables)
			}
		}
	}
}

var replaceTests = []struct {
	Name         string
	Query        parser.ReplaceQuery
	ResultFile   *FileInfo
	UpdateCount  int
	ViewCache    ViewMap
	ResultScopes *ReferenceScope
	Error        string
}{
	{
		Name: "Replace Query",
		Query: parser.ReplaceQuery{
			WithClause: parser.WithClause{
				InlineTables: []parser.QueryExpression{
					parser.InlineTable{
						Name: parser.Identifier{Literal: "it"},
						Fields: []parser.QueryExpression{
							parser.Identifier{Literal: "c1"},
							parser.Identifier{Literal: "c2"},
						},
						Query: parser.SelectQuery{
							SelectEntity: parser.SelectEntity{
								SelectClause: parser.SelectClause{
									Fields: []parser.QueryExpression{
										parser.Field{Object: parser.NewIntegerValueFromString("2")},
										parser.Field{Object: parser.NewStringValue("str3")},
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
							parser.NewStringValue("str4"),
						},
					},
				},
				parser.RowValue{
					Value: parser.Subquery{
						Query: parser.SelectQuery{
							SelectEntity: parser.SelectEntity{
								SelectClause: parser.SelectClause{
									Fields: []parser.QueryExpression{
										parser.Field{Object: parser.FieldReference{View: parser.Identifier{Literal: "it"}, Column: parser.Identifier{Literal: "c1"}}},
										parser.Field{Object: parser.FieldReference{View: parser.Identifier{Literal: "it"}, Column: parser.Identifier{Literal: "c2"}}},
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
		ResultFile: &FileInfo{
			Path:      GetTestFilePath("table1.csv"),
			Delimiter: ',',
			NoHeader:  false,
			Encoding:  text.UTF8,
			LineBreak: text.LF,
			ForUpdate: true,
		},
		UpdateCount: 2,
		ViewCache: GenerateViewMap([]*View{
			{
				FileInfo: &FileInfo{
					Path:      GetTestFilePath("table1.csv"),
					Delimiter: ',',
					NoHeader:  false,
					Encoding:  text.UTF8,
					LineBreak: text.LF,
					ForUpdate: true,
				},
				Header: NewHeader("table1", []string{"column1", "column2"}),
				RecordSet: []Record{
					NewRecord([]value.Primary{
						value.NewString("1"),
						value.NewString("str1"),
					}),
					NewRecord([]value.Primary{
						value.NewString("2"),
						value.NewString("str3"),
					}),
					NewRecord([]value.Primary{
						value.NewString("3"),
						value.NewString("str3"),
					}),
					NewRecord([]value.Primary{
						value.NewInteger(4),
						value.NewString("str4"),
					}),
				},
			},
		}),
	},
	{
		Name: "Replace Query to Empty Table",
		Query: parser.ReplaceQuery{
			Table: parser.Table{Object: parser.Identifier{Literal: "table_empty"}},
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
							parser.NewStringValue("str4"),
						},
					},
				},
			},
		},
		ResultFile: &FileInfo{
			Path:      GetTestFilePath("table_empty.csv"),
			Delimiter: ',',
			NoHeader:  false,
			Encoding:  text.UTF8,
			LineBreak: text.LF,
			ForUpdate: true,
		},
		UpdateCount: 1,
		ViewCache: GenerateViewMap([]*View{
			{
				FileInfo: &FileInfo{
					Path:      GetTestFilePath("table_empty.csv"),
					Delimiter: ',',
					NoHeader:  false,
					Encoding:  text.UTF8,
					LineBreak: text.LF,
					ForUpdate: true,
				},
				Header: NewHeader("table_empty", []string{"column1", "column2"}),
				RecordSet: []Record{
					NewRecord([]value.Primary{
						value.NewInteger(4),
						value.NewString("str4"),
					}),
				},
			},
		}),
	},
	{
		Name: "Replace Query For Temporary View",
		Query: parser.ReplaceQuery{
			Table: parser.Table{Object: parser.Identifier{Literal: "tmpview"}, Alias: parser.Identifier{Literal: "t"}},
			Fields: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
			},
			Keys: []parser.QueryExpression{
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
			Path:      "tmpview",
			Delimiter: ',',
			ViewType:  ViewTypeTemporaryTable,
		},
		UpdateCount: 2,
		ResultScopes: GenerateReferenceScope([]map[string]map[string]interface{}{
			{
				scopeNameTempTables: {
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
						},
						FileInfo: &FileInfo{
							Path:      "tmpview",
							Delimiter: ',',
							ViewType:  ViewTypeTemporaryTable,
						},
					},
				},
			},
		}, nil, time.Time{}, nil),
	},
	{
		Name: "Replace Query All Fields",
		Query: parser.ReplaceQuery{
			Table: parser.Table{Object: parser.Identifier{Literal: "table1"}},
			Keys: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
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
		ResultFile: &FileInfo{
			Path:      GetTestFilePath("table1.csv"),
			Delimiter: ',',
			NoHeader:  false,
			Encoding:  text.UTF8,
			LineBreak: text.LF,
			ForUpdate: true,
		},
		UpdateCount: 2,
	},
	{
		Name: "Replace Query File Does Not Exist Error",
		Query: parser.ReplaceQuery{
			Table: parser.Table{Object: parser.Identifier{Literal: "notexist"}},
			Fields: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
			},
			Keys: []parser.QueryExpression{
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
		Error: "file notexist does not exist",
	},
	{
		Name: "Replace Query Field Does Not Exist Error",
		Query: parser.ReplaceQuery{
			Table: parser.Table{Object: parser.Identifier{Literal: "table1"}},
			Fields: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			},
			Keys: []parser.QueryExpression{
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
		Error: "field notexist does not exist",
	},
	{
		Name: "Replace Select Query",
		Query: parser.ReplaceQuery{
			Table: parser.Table{Object: parser.Identifier{Literal: "table1"}},
			Fields: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
			},
			Keys: []parser.QueryExpression{
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
		ResultFile: &FileInfo{
			Path:      GetTestFilePath("table1.csv"),
			Delimiter: ',',
			NoHeader:  false,
			Encoding:  text.UTF8,
			LineBreak: text.LF,
			ForUpdate: true,
		},
		UpdateCount: 3,
	},
	{
		Name: "Replace Select Query Field Does Not Exist Error",
		Query: parser.ReplaceQuery{
			Table: parser.Table{Object: parser.Identifier{Literal: "table1"}},
			Fields: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
			},
			Keys: []parser.QueryExpression{
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
		Error: "select query should return exactly 1 field",
	},
}

func TestReplace(t *testing.T) {
	defer func() {
		_ = TestTx.ReleaseResources()
		_ = TestTx.CachedViews.Clean(TestTx.FileContainer)
		initFlag(TestTx.Flags)
	}()

	TestTx.Flags.Repository = TestDir
	TestTx.Flags.Quiet = false

	scope := GenerateReferenceScope([]map[string]map[string]interface{}{
		{
			scopeNameTempTables: {
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
						Path:      "tmpview",
						Delimiter: ',',
						ViewType:  ViewTypeTemporaryTable,
					},
				},
			},
		},
	}, nil, time.Time{}, nil)

	ctx := context.Background()
	for _, v := range replaceTests {
		_ = TestTx.ReleaseResources()
		result, cnt, err := Replace(ctx, scope, v.Query)
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

		TestTx.CachedViews.Range(func(key, value interface{}) bool {
			view := value.(*View)
			if view.FileInfo.Handler != nil {
				if view.FileInfo.Path != view.FileInfo.Handler.Path() {
					t.Errorf("file pointer = %q, want %q for %q", view.FileInfo.Handler.Path(), view.FileInfo.Path, v.Name)
				}
				_ = TestTx.FileContainer.Close(view.FileInfo.Handler)
				view.FileInfo.Handler = nil
			}
			return true
		})

		if !reflect.DeepEqual(result, v.ResultFile) {
			t.Errorf("%s: fileinfo = %v, want %v", v.Name, result, v.ResultFile)
		}

		if !reflect.DeepEqual(cnt, v.UpdateCount) {
			t.Errorf("%s: update count = %d, want %d", v.Name, cnt, v.UpdateCount)
		}

		if v.ViewCache.SyncMap != nil {
			if !SyncMapEqual(TestTx.CachedViews, v.ViewCache) {
				t.Errorf("%s: view cache = %v, want %v", v.Name, TestTx.CachedViews, v.ViewCache)
			}
		}
		if v.ResultScopes != nil {
			if !SyncMapEqual(scope.Blocks[0].TemporaryTables, v.ResultScopes.Blocks[0].TemporaryTables) {
				t.Errorf("%s: temporary views list = %v, want %v", v.Name, scope.Blocks[0].TemporaryTables, v.ResultScopes.Blocks[0].TemporaryTables)
			}
		}
	}
}

var deleteTests = []struct {
	Name         string
	Query        parser.DeleteQuery
	ResultFiles  []*FileInfo
	UpdateCounts []int
	ViewCache    ViewMap
	ResultScopes *ReferenceScope
	Error        string
}{
	{
		Name: "Delete Query",
		Query: parser.DeleteQuery{
			WithClause: parser.WithClause{
				InlineTables: []parser.QueryExpression{
					parser.InlineTable{
						Name: parser.Identifier{Literal: "it"},
						Fields: []parser.QueryExpression{
							parser.Identifier{Literal: "c1"},
						},
						Query: parser.SelectQuery{
							SelectEntity: parser.SelectEntity{
								SelectClause: parser.SelectClause{
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
					Operator: parser.Token{Token: '=', Literal: "="},
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
				ForUpdate: true,
			},
		},
		UpdateCounts: []int{1},
		ViewCache: GenerateViewMap([]*View{
			{
				FileInfo: &FileInfo{
					Path:      GetTestFilePath("table1.csv"),
					Delimiter: ',',
					NoHeader:  false,
					Encoding:  text.UTF8,
					LineBreak: text.LF,
					ForUpdate: true,
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
			},
		}),
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
					Operator: parser.Token{Token: '=', Literal: "="},
				},
			},
		},
		ResultFiles: []*FileInfo{
			{
				Path:      "tmpview",
				Delimiter: ',',
				ViewType:  ViewTypeTemporaryTable,
			},
		},
		UpdateCounts: []int{1},
		ResultScopes: GenerateReferenceScope([]map[string]map[string]interface{}{
			{
				scopeNameTempTables: {
					"TMPVIEW": &View{
						Header: NewHeader("tmpview", []string{"column1", "column2"}),
						RecordSet: []Record{
							NewRecord([]value.Primary{
								value.NewString("1"),
								value.NewString("str1"),
							}),
						},
						FileInfo: &FileInfo{
							Path:      "tmpview",
							Delimiter: ',',
							ViewType:  ViewTypeTemporaryTable,
						},
					},
				},
			},
		}, nil, time.Time{}, nil),
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
								Operator: parser.Token{Token: '=', Literal: "="},
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
				ForUpdate: true,
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
								Operator: parser.Token{Token: '=', Literal: "="},
							},
						},
					}},
				},
			},
		},
		Error: "tables to delete records are not specified",
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
					Operator: parser.Token{Token: '=', Literal: "="},
				},
			},
		},
		Error: "file notexist does not exist",
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
					Operator: parser.Token{Token: '=', Literal: "="},
				},
			},
		},
		Error: "field notexist does not exist",
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
								Operator: parser.Token{Token: '=', Literal: "="},
							},
						},
					}},
				},
			},
		},
		Error: "table notexist is not loaded",
	},
}

func TestDelete(t *testing.T) {
	defer func() {
		_ = TestTx.ReleaseResources()
		_ = TestTx.CachedViews.Clean(TestTx.FileContainer)
		initFlag(TestTx.Flags)
	}()

	TestTx.Flags.Repository = TestDir
	TestTx.Flags.Quiet = false

	scope := GenerateReferenceScope([]map[string]map[string]interface{}{
		{
			scopeNameTempTables: {
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
						Path:      "tmpview",
						Delimiter: ',',
						ViewType:  ViewTypeTemporaryTable,
					},
				},
			},
		},
	}, nil, time.Time{}, nil)

	ctx := context.Background()
	for _, v := range deleteTests {
		_ = TestTx.ReleaseResources()
		files, cnt, err := Delete(ctx, scope, v.Query)
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

		TestTx.CachedViews.Range(func(key, value interface{}) bool {
			view := value.(*View)
			if view.FileInfo.Handler != nil {
				if view.FileInfo.Path != view.FileInfo.Handler.Path() {
					t.Errorf("file pointer = %q, want %q for %q", view.FileInfo.Handler.Path(), view.FileInfo.Path, v.Name)
				}
				_ = TestTx.FileContainer.Close(view.FileInfo.Handler)
				view.FileInfo.Handler = nil
			}
			return true
		})

		if !reflect.DeepEqual(files, v.ResultFiles) {
			t.Errorf("%s: fileinfo = %v, want %v", v.Name, files, v.ResultFiles)
		}

		if !reflect.DeepEqual(cnt, v.UpdateCounts) {
			t.Errorf("%s: update count = %v, want %v", v.Name, cnt, v.UpdateCounts)
		}

		if v.ViewCache.SyncMap != nil {
			if !SyncMapEqual(TestTx.CachedViews, v.ViewCache) {
				t.Errorf("%s: view cache = %v, want %v", v.Name, TestTx.CachedViews, v.ViewCache)
			}
		}
		if v.ResultScopes != nil {
			if !SyncMapEqual(scope.Blocks[0].TemporaryTables, v.ResultScopes.Blocks[0].TemporaryTables) {
				t.Errorf("%s: temporary views list = %v, want %v", v.Name, scope.Blocks[0].TemporaryTables, v.ResultScopes.Blocks[0].TemporaryTables)
			}
		}
	}
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
			ForUpdate: true,
		},
		ViewCache: GenerateViewMap([]*View{
			{
				FileInfo: &FileInfo{
					Path:      GetTestFilePath("create_table_1.csv"),
					Delimiter: ',',
					NoHeader:  false,
					Encoding:  text.UTF8,
					LineBreak: text.LF,
					ForUpdate: true,
				},
				Header:    NewHeader("create_table_1", []string{"column1", "column2"}),
				RecordSet: RecordSet{},
			},
		}),
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
			ForUpdate: true,
		},
		ViewCache: GenerateViewMap([]*View{
			{
				FileInfo: &FileInfo{
					Path:      GetTestFilePath("create_table_1.csv"),
					Delimiter: ',',
					NoHeader:  false,
					Encoding:  text.UTF8,
					LineBreak: text.LF,
					ForUpdate: true,
				},
				Header: NewHeader("create_table_1", []string{"column1", "column2"}),
				RecordSet: RecordSet{
					NewRecord([]value.Primary{
						value.NewInteger(1),
						value.NewInteger(2),
					}),
				},
			},
		}),
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
		Error: fmt.Sprintf("file %s already exists", GetTestFilePath("table1.csv")),
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
		Error: "field name column1 is a duplicate",
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
		Error: "field notexist does not exist",
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
		Error: "select query should return exactly 1 field for table create_table_1.csv",
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
		Error: "field name column1 is a duplicate",
	},
}

func TestCreateTable(t *testing.T) {
	defer func() {
		_ = TestTx.ReleaseResources()
		_ = TestTx.CachedViews.Clean(TestTx.FileContainer)
		initFlag(TestTx.Flags)
	}()

	TestTx.Flags.Repository = TestDir
	TestTx.Flags.Quiet = false

	scope := NewReferenceScope(TestTx)
	ctx := context.Background()
	for _, v := range createTableTests {
		_ = TestTx.ReleaseResources()

		result, err := CreateTable(ctx, scope, v.Query)

		if result != nil {
			_ = TestTx.FileContainer.Close(result.Handler)
			result.Handler = nil
		}
		TestTx.CachedViews.Range(func(key, value interface{}) bool {
			view := value.(*View)
			if view.FileInfo != nil {
				_ = TestTx.FileContainer.Close(view.FileInfo.Handler)
				view.FileInfo.Handler = nil
			}
			return true
		})

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

		if v.ViewCache.SyncMap != nil {
			if !SyncMapEqual(TestTx.CachedViews, v.ViewCache) {
				t.Errorf("%s: view cache = %v, want %v", v.Name, TestTx.CachedViews, v.ViewCache)
			}
		}
	}
}

var addColumnsTests = []struct {
	Name         string
	Query        parser.AddColumns
	ResultFile   *FileInfo
	UpdateCount  int
	ViewCache    ViewMap
	ResultScopes *ReferenceScope
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
			ForUpdate: true,
		},
		UpdateCount: 2,
		ViewCache: GenerateViewMap([]*View{
			{
				FileInfo: &FileInfo{
					Path:      GetTestFilePath("table1.csv"),
					Delimiter: ',',
					NoHeader:  false,
					Encoding:  text.UTF8,
					LineBreak: text.LF,
					ForUpdate: true,
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
			},
		}),
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
			Path:      "tmpview",
			Delimiter: ',',
			ViewType:  ViewTypeTemporaryTable,
		},
		UpdateCount: 2,
		ResultScopes: GenerateReferenceScope([]map[string]map[string]interface{}{
			{
				scopeNameTempTables: {
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
							Path:      "tmpview",
							Delimiter: ',',
							ViewType:  ViewTypeTemporaryTable,
						},
					},
				},
			},
		}, nil, time.Time{}, nil),
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
			ForUpdate: true,
		},
		UpdateCount: 2,
		ViewCache: GenerateViewMap([]*View{
			{
				FileInfo: &FileInfo{
					Path:      GetTestFilePath("table1.csv"),
					Delimiter: ',',
					NoHeader:  false,
					Encoding:  text.UTF8,
					LineBreak: text.LF,
					ForUpdate: true,
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
			},
		}),
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
			ForUpdate: true,
		},
		UpdateCount: 2,
		ViewCache: GenerateViewMap([]*View{
			{
				FileInfo: &FileInfo{
					Path:      GetTestFilePath("table1.csv"),
					Delimiter: ',',
					NoHeader:  false,
					Encoding:  text.UTF8,
					LineBreak: text.LF,
					ForUpdate: true,
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
			},
		}),
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
			ForUpdate: true,
		},
		UpdateCount: 2,
		ViewCache: GenerateViewMap([]*View{
			{
				FileInfo: &FileInfo{
					Path:      GetTestFilePath("table1.csv"),
					Delimiter: ',',
					NoHeader:  false,
					Encoding:  text.UTF8,
					LineBreak: text.LF,
					ForUpdate: true,
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
			},
		}),
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
		Error: "file notexist does not exist",
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
		Error: "field notexist does not exist",
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
		Error: "field name column1 is a duplicate",
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
		Error: "field notexist does not exist",
	},
}

func TestAddColumns(t *testing.T) {
	defer func() {
		_ = TestTx.ReleaseResources()
		_ = TestTx.CachedViews.Clean(TestTx.FileContainer)
		initFlag(TestTx.Flags)
	}()

	TestTx.Flags.Repository = TestDir
	TestTx.Flags.Quiet = false

	scope := GenerateReferenceScope([]map[string]map[string]interface{}{
		{
			scopeNameTempTables: {
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
						Path:      "tmpview",
						Delimiter: ',',
						ViewType:  ViewTypeTemporaryTable,
					},
				},
			},
		},
	}, nil, time.Time{}, nil)

	ctx := context.Background()
	for _, v := range addColumnsTests {
		_ = TestTx.ReleaseResources()
		result, cnt, err := AddColumns(ctx, scope, v.Query)
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

		TestTx.CachedViews.Range(func(key, value interface{}) bool {
			view := value.(*View)
			if view.FileInfo.Handler != nil {
				if view.FileInfo.Path != view.FileInfo.Handler.Path() {
					t.Errorf("file pointer = %q, want %q for %q", view.FileInfo.Handler.Path(), view.FileInfo.Path, v.Name)
				}
				_ = TestTx.FileContainer.Close(view.FileInfo.Handler)
				view.FileInfo.Handler = nil
			}
			return true
		})

		if !reflect.DeepEqual(result, v.ResultFile) {
			t.Errorf("%s: fileinfo = %v, want %v", v.Name, result, v.ResultFile)
		}

		if cnt != v.UpdateCount {
			t.Errorf("%s: update count = %d, want %d", v.Name, cnt, v.UpdateCount)
		}

		if v.ViewCache.SyncMap != nil {
			if !SyncMapEqual(TestTx.CachedViews, v.ViewCache) {
				t.Errorf("%s: view cache = %v, want %v", v.Name, TestTx.CachedViews, v.ViewCache)
			}
		}
		if v.ResultScopes != nil {
			if !SyncMapEqual(scope.Blocks[0].TemporaryTables, v.ResultScopes.Blocks[0].TemporaryTables) {
				t.Errorf("%s: temporary views list = %v, want %v", v.Name, scope.Blocks[0].TemporaryTables, v.ResultScopes.Blocks[0].TemporaryTables)
			}
		}
	}
}

var dropColumnsTests = []struct {
	Name         string
	Query        parser.DropColumns
	Result       *FileInfo
	UpdateCount  int
	ViewCache    ViewMap
	ResultScopes *ReferenceScope
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
			ForUpdate: true,
		},
		UpdateCount: 1,
		ViewCache: GenerateViewMap([]*View{
			{
				FileInfo: &FileInfo{
					Path:      GetTestFilePath("table1.csv"),
					Delimiter: ',',
					NoHeader:  false,
					Encoding:  text.UTF8,
					LineBreak: text.LF,
					ForUpdate: true,
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
			},
		}),
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
			Path:      "tmpview",
			Delimiter: ',',
			ViewType:  ViewTypeTemporaryTable,
		},
		UpdateCount: 1,
		ResultScopes: GenerateReferenceScope([]map[string]map[string]interface{}{
			{
				scopeNameTempTables: {
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
							Path:      "tmpview",
							Delimiter: ',',
							ViewType:  ViewTypeTemporaryTable,
						},
					},
				},
			},
		}, nil, time.Time{}, nil),
	},
	{
		Name: "Drop Fields Load Error",
		Query: parser.DropColumns{
			Table: parser.Identifier{Literal: "notexist"},
			Columns: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
			},
		},
		Error: "file notexist does not exist",
	},
	{
		Name: "Drop Fields Field Does Not Exist Error",
		Query: parser.DropColumns{
			Table: parser.Identifier{Literal: "table1"},
			Columns: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			},
		},
		Error: "field notexist does not exist",
	},
}

func TestDropColumns(t *testing.T) {
	defer func() {
		_ = TestTx.ReleaseResources()
		_ = TestTx.CachedViews.Clean(TestTx.FileContainer)
		initFlag(TestTx.Flags)
	}()

	TestTx.Flags.Repository = TestDir
	TestTx.Flags.Quiet = false

	scope := GenerateReferenceScope([]map[string]map[string]interface{}{
		{
			scopeNameTempTables: {
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
						Path:      "tmpview",
						Delimiter: ',',
						ViewType:  ViewTypeTemporaryTable,
					},
				},
			},
		},
	}, nil, time.Time{}, nil)

	ctx := context.Background()
	for _, v := range dropColumnsTests {
		_ = TestTx.ReleaseResources()
		result, cnt, err := DropColumns(ctx, scope, v.Query)
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

		TestTx.CachedViews.Range(func(key, value interface{}) bool {
			view := value.(*View)
			if view.FileInfo.Handler != nil {
				if view.FileInfo.Path != view.FileInfo.Handler.Path() {
					t.Errorf("file pointer = %q, want %q for %q", view.FileInfo.Handler.Path(), view.FileInfo.Path, v.Name)
				}
				_ = TestTx.FileContainer.Close(view.FileInfo.Handler)
				view.FileInfo.Handler = nil
			}
			return true
		})

		if !reflect.DeepEqual(result, v.Result) {
			t.Errorf("%s: fileinfo = %v, want %v", v.Name, result, v.Result)
		}

		if cnt != v.UpdateCount {
			t.Errorf("%s: update count = %d, want %d", v.Name, cnt, v.UpdateCount)
		}

		if v.ViewCache.SyncMap != nil {
			if !SyncMapEqual(TestTx.CachedViews, v.ViewCache) {
				t.Errorf("%s: view cache = %v, want %v", v.Name, TestTx.CachedViews, v.ViewCache)
			}
		}
		if v.ResultScopes != nil {
			if !SyncMapEqual(scope.Blocks[0].TemporaryTables, v.ResultScopes.Blocks[0].TemporaryTables) {
				t.Errorf("%s: temporary views list = %v, want %v", v.Name, scope.Blocks[0].TemporaryTables, v.ResultScopes.Blocks[0].TemporaryTables)
			}
		}
	}
}

var renameColumnTests = []struct {
	Name         string
	Query        parser.RenameColumn
	Result       *FileInfo
	ViewCache    ViewMap
	ResultScopes *ReferenceScope
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
			ForUpdate: true,
		},
		ViewCache: GenerateViewMap([]*View{
			{
				FileInfo: &FileInfo{
					Path:      GetTestFilePath("table1.csv"),
					Delimiter: ',',
					NoHeader:  false,
					Encoding:  text.UTF8,
					LineBreak: text.LF,
					ForUpdate: true,
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
			},
		}),
	},
	{
		Name: "Rename Column For Temporary View",
		Query: parser.RenameColumn{
			Table: parser.Identifier{Literal: "tmpview"},
			Old:   parser.ColumnNumber{View: parser.Identifier{Literal: "tmpview"}, Number: value.NewInteger(2)},
			New:   parser.Identifier{Literal: "newcolumn"},
		},
		Result: &FileInfo{
			Path:      "tmpview",
			Delimiter: ',',
			ViewType:  ViewTypeTemporaryTable,
		},
		ResultScopes: GenerateReferenceScope([]map[string]map[string]interface{}{
			{
				scopeNameTempTables: {
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
							Path:      "tmpview",
							Delimiter: ',',
							ViewType:  ViewTypeTemporaryTable,
						},
					},
				},
			},
		}, nil, time.Time{}, nil),
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
		Error: "field name column1 is a duplicate",
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
	defer func() {
		_ = TestTx.ReleaseResources()
		_ = TestTx.CachedViews.Clean(TestTx.FileContainer)
		initFlag(TestTx.Flags)
	}()

	TestTx.Flags.Repository = TestDir
	TestTx.Flags.Quiet = false

	scope := GenerateReferenceScope([]map[string]map[string]interface{}{
		{
			scopeNameTempTables: {
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
						Path:      "tmpview",
						Delimiter: ',',
						ViewType:  ViewTypeTemporaryTable,
					},
				},
			},
		},
	}, nil, time.Time{}, nil)

	ctx := context.Background()
	for _, v := range renameColumnTests {
		_ = TestTx.ReleaseResources()
		result, err := RenameColumn(ctx, scope, v.Query)
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

		TestTx.CachedViews.Range(func(key, value interface{}) bool {
			view := value.(*View)
			if view.FileInfo.Handler != nil {
				if view.FileInfo.Path != view.FileInfo.Handler.Path() {
					t.Errorf("file pointer = %q, want %q for %q", view.FileInfo.Handler.Path(), view.FileInfo.Path, v.Name)
				}
				_ = TestTx.FileContainer.Close(view.FileInfo.Handler)
				view.FileInfo.Handler = nil
			}
			return true
		})

		if !reflect.DeepEqual(result, v.Result) {
			t.Errorf("%s: result = %v, want %v", v.Name, result, v.Result)
		}

		if v.ViewCache.SyncMap != nil {
			if !SyncMapEqual(TestTx.CachedViews, v.ViewCache) {
				t.Errorf("%s: view cache = %v, want %v", v.Name, TestTx.CachedViews, v.ViewCache)
			}
		}
		if v.ResultScopes != nil {
			if !SyncMapEqual(scope.Blocks[0].TemporaryTables, v.ResultScopes.Blocks[0].TemporaryTables) {
				t.Errorf("%s: temporary views list = %v, want %v", v.Name, scope.Blocks[0].TemporaryTables, v.ResultScopes.Blocks[0].TemporaryTables)
			}
		}
	}
}

var setTableAttributeTests = []struct {
	Name   string
	Query  parser.SetTableAttribute
	Expect *FileInfo
	Error  string
}{
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
			Format:    option.TSV,
			Encoding:  text.UTF8,
			LineBreak: text.LF,
			ForUpdate: true,
		},
	},
	{
		Name: "Set Delimiter to TSV with TableObject",
		Query: parser.SetTableAttribute{
			Table: parser.TableObject{
				Type:          parser.Token{Token: parser.CSV, Literal: "csv"},
				Path:          parser.Identifier{Literal: "table1.csv"},
				FormatElement: parser.NewStringValue(","),
			},
			Attribute: parser.Identifier{Literal: "delimiter"},
			Value:     parser.NewStringValue("\t"),
		},
		Expect: &FileInfo{
			Path:      GetTestFilePath("table1.csv"),
			Delimiter: '\t',
			Format:    option.TSV,
			Encoding:  text.UTF8,
			LineBreak: text.LF,
			ForUpdate: true,
		},
	},
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
			Format:    option.CSV,
			Encoding:  text.UTF8,
			LineBreak: text.LF,
			ForUpdate: true,
		},
	},
	{
		Name: "Set Delimiter Error",
		Query: parser.SetTableAttribute{
			Table:     parser.Identifier{Literal: "table1.csv"},
			Attribute: parser.Identifier{Literal: "delimiter"},
			Value:     parser.NewStringValue("aa"),
		},
		Error: "delimiter must be one character",
	},
	{
		Name: "Set Delimiter Not Allowed Value",
		Query: parser.SetTableAttribute{
			Table:     parser.Identifier{Literal: "table1.csv"},
			Attribute: parser.Identifier{Literal: "delimiter"},
			Value:     parser.NewNullValue(),
		},
		Error: "NULL for delimiter is not allowed",
	},
	{
		Name: "Set DelimiterPositions",
		Query: parser.SetTableAttribute{
			Table:     parser.Identifier{Literal: "table1.csv"},
			Attribute: parser.Identifier{Literal: "delimiter_positions"},
			Value:     parser.NewStringValue("S[2, 5, 10]"),
		},
		Expect: &FileInfo{
			Path:               GetTestFilePath("table1.csv"),
			Delimiter:          ',',
			DelimiterPositions: []int{2, 5, 10},
			Format:             option.FIXED,
			Encoding:           text.UTF8,
			SingleLine:         true,
			LineBreak:          text.LF,
			ForUpdate:          true,
		},
	},
	{
		Name: "Set DelimiterPositions Error",
		Query: parser.SetTableAttribute{
			Table:     parser.Identifier{Literal: "table1.csv"},
			Attribute: parser.Identifier{Literal: "delimiter_positions"},
			Value:     parser.NewStringValue("invalid"),
		},
		Error: "delimiter positions must be \"SPACES\" or a JSON array of integers",
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
			Format:    option.TEXT,
			Encoding:  text.UTF8,
			LineBreak: text.LF,
			ForUpdate: true,
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
			Format:    option.JSON,
			Encoding:  text.UTF8,
			LineBreak: text.LF,
			ForUpdate: true,
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
			Format:    option.TSV,
			Encoding:  text.UTF8,
			LineBreak: text.LF,
			ForUpdate: true,
		},
	},
	{
		Name: "Set Format Error",
		Query: parser.SetTableAttribute{
			Table:     parser.Identifier{Literal: "table1.csv"},
			Attribute: parser.Identifier{Literal: "format"},
			Value:     parser.NewStringValue("invalid"),
		},
		Error: "format must be one of CSV|TSV|FIXED|JSON|JSONL|LTSV|GFM|ORG|BOX|TEXT",
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
			Format:    option.CSV,
			Encoding:  text.SJIS,
			LineBreak: text.LF,
			ForUpdate: true,
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
			Format:    option.CSV,
			Encoding:  text.SJIS,
			LineBreak: text.LF,
			ForUpdate: true,
		},
	},
	{
		Name: "Set Encoding Error",
		Query: parser.SetTableAttribute{
			Table:     parser.Identifier{Literal: "table1.csv"},
			Attribute: parser.Identifier{Literal: "encoding"},
			Value:     parser.NewStringValue("invalid"),
		},
		Error: "encoding must be one of UTF8|UTF8M|UTF16|UTF16BE|UTF16LE|UTF16BEM|UTF16LEM|SJIS",
	},
	{
		Name: "Set Encoding Error in JSON Format",
		Query: parser.SetTableAttribute{
			Table:     parser.Identifier{Literal: "table.json"},
			Attribute: parser.Identifier{Literal: "encoding"},
			Value:     parser.NewStringValue("sjis"),
		},
		Error: "json format is supported only UTF8",
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
			Format:    option.CSV,
			Encoding:  text.UTF8,
			LineBreak: text.CRLF,
			ForUpdate: true,
		},
	},
	{
		Name: "Set LineBreak Error",
		Query: parser.SetTableAttribute{
			Table:     parser.Identifier{Literal: "table1.csv"},
			Attribute: parser.Identifier{Literal: "line_break"},
			Value:     parser.NewStringValue("invalid"),
		},
		Error: "line-break must be one of CRLF|CR|LF",
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
			Format:    option.CSV,
			Encoding:  text.UTF8,
			LineBreak: text.LF,
			NoHeader:  true,
			ForUpdate: true,
		},
	},
	{
		Name: "Set NoHeader Not Allowed Value",
		Query: parser.SetTableAttribute{
			Table:     parser.Identifier{Literal: "table1.csv"},
			Attribute: parser.Identifier{Literal: "header"},
			Value:     parser.NewNullValue(),
		},
		Error: "NULL for header is not allowed",
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
			Format:     option.CSV,
			Encoding:   text.UTF8,
			LineBreak:  text.LF,
			EncloseAll: true,
			ForUpdate:  true,
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
			Format:      option.JSON,
			Encoding:    text.UTF8,
			LineBreak:   text.LF,
			JsonEscape:  json.HexDigits,
			PrettyPrint: false,
			ForUpdate:   true,
		},
	},
	{
		Name: "Set JsonEscape Error",
		Query: parser.SetTableAttribute{
			Table:     parser.Identifier{Literal: "table.json"},
			Attribute: parser.Identifier{Literal: "json_escape"},
			Value:     parser.NewStringValue("invalid"),
		},
		Error: "json escape type must be one of BACKSLASH|HEX|HEXALL",
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
			Format:      option.JSON,
			Encoding:    text.UTF8,
			LineBreak:   text.LF,
			PrettyPrint: true,
			ForUpdate:   true,
		},
	},
	{
		Name: "Not Exist Table Error",
		Query: parser.SetTableAttribute{
			Table:     parser.Identifier{Literal: "notexist.csv"},
			Attribute: parser.Identifier{Literal: "delimiter"},
			Value:     parser.NewStringValue(","),
		},
		Error: "file notexist.csv does not exist",
	},
	{
		Name: "Temporary View Error",
		Query: parser.SetTableAttribute{
			Table:     parser.Identifier{Literal: "tmpview"},
			Attribute: parser.Identifier{Literal: "delimiter"},
			Value:     parser.NewStringValue(","),
		},
		Error: "view has no attributes",
	},
	{
		Name: "Value Evaluation Error",
		Query: parser.SetTableAttribute{
			Table:     parser.Identifier{Literal: "table1.csv"},
			Attribute: parser.Identifier{Literal: "delimiter"},
			Value:     parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
		},
		Error: "field notexist does not exist",
	},
	{
		Name: "Not Exist Attribute Error",
		Query: parser.SetTableAttribute{
			Table:     parser.Identifier{Literal: "table1.csv"},
			Attribute: parser.Identifier{Literal: "notexist"},
			Value:     parser.NewStringValue(","),
		},
		Error: "table attribute notexist does not exist",
	},
}

func TestSetTableAttribute(t *testing.T) {
	defer func() {
		_ = TestTx.ReleaseResources()
		_ = TestTx.CachedViews.Clean(TestTx.FileContainer)
		initFlag(TestTx.Flags)
	}()

	TestTx.Flags.Repository = TestDir
	TestTx.Flags.Quiet = false

	scope := GenerateReferenceScope([]map[string]map[string]interface{}{
		{
			scopeNameTempTables: {
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
						Path:      "tmpview",
						Delimiter: ',',
						ViewType:  ViewTypeTemporaryTable,
					},
				},
			},
		},
	}, nil, time.Time{}, nil)

	ctx := context.Background()
	for _, v := range setTableAttributeTests {
		_ = TestTx.ReleaseResources()

		_, _, err := SetTableAttribute(ctx, scope, v.Query)
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

		TestTx.CachedViews.Range(func(key, value interface{}) bool {
			view := value.(*View)
			if view.FileInfo.Handler != nil {
				if view.FileInfo.Path != view.FileInfo.Handler.Path() {
					t.Errorf("file pointer = %q, want %q for %q", view.FileInfo.Handler.Path(), view.FileInfo.Path, v.Name)
				}
				_ = TestTx.FileContainer.Close(view.FileInfo.Handler)
				view.FileInfo.Handler = nil
			}
			return true
		})

		view, _ := LoadViewFromTableIdentifier(ctx, scope.CreateNode(), v.Query.Table, false, false)

		if !reflect.DeepEqual(view.FileInfo, v.Expect) {
			t.Errorf("%s: result = %v, want %v", v.Name, view.FileInfo, v.Expect)
		}

		_, _, err = SetTableAttribute(ctx, scope, v.Query)
		if err == nil {
			t.Errorf("%s: no error, want TableAttributeUnchangedError for duplicate set", v.Name)
		} else if _, ok := err.(*TableAttributeUnchangedError); !ok {
			t.Errorf("%s: error = %T, want TableAttributeUnchangedError for duplicate set", v.Name, err)
		}
	}
}
