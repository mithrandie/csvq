package query

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/value"
)

func TestNewViewFromGroupedRecord(t *testing.T) {
	fr := ReferenceRecord{
		view: &View{
			Header: NewHeaderWithId("table1", []string{"column1", "column2"}),
			RecordSet: []Record{
				{
					NewGroupCell([]value.Primary{value.NewInteger(1), value.NewInteger(2), value.NewInteger(3)}),
					NewGroupCell([]value.Primary{value.NewInteger(1), value.NewInteger(2), value.NewInteger(3)}),
					NewGroupCell([]value.Primary{value.NewString("str1"), value.NewString("str2"), value.NewString("str3")}),
				},
			},
		},
		recordIndex: 0,
		cache:       NewFieldIndexCache(10, LimitToUseFieldIndexSliceChache),
	}
	expect := &View{
		Header: NewHeaderWithId("table1", []string{"column1", "column2"}),
		RecordSet: []Record{
			{NewCell(value.NewInteger(1)), NewCell(value.NewInteger(1)), NewCell(value.NewString("str1"))},
			{NewCell(value.NewInteger(2)), NewCell(value.NewInteger(2)), NewCell(value.NewString("str2"))},
			{NewCell(value.NewInteger(3)), NewCell(value.NewInteger(3)), NewCell(value.NewString("str3"))},
		},
	}

	result, _ := NewViewFromGroupedRecord(context.Background(), TestTx.Flags, fr)
	if !reflect.DeepEqual(result, expect) {
		t.Errorf("result = %v, want %v", result, expect)
	}
}

var viewWhereTests = []struct {
	Name   string
	CPU    int
	View   *View
	Where  parser.WhereClause
	Result RecordSet
	Error  string
}{
	{
		Name: "Where",
		View: &View{
			Header: NewHeaderWithId("table1", []string{"column1", "column2"}),
			RecordSet: RecordSet{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
				}),
			},
		},
		Where: parser.WhereClause{
			Filter: parser.Comparison{
				LHS:      parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
				RHS:      parser.NewIntegerValueFromString("2"),
				Operator: parser.Token{Token: '=', Literal: "="},
			},
		},
		Result: RecordSet{
			NewRecordWithId(2, []value.Primary{
				value.NewString("2"),
				value.NewString("str2"),
			}),
		},
	},
	{
		Name: "Where in Multi Threading",
		CPU:  3,
		View: &View{
			Header: NewHeaderWithId("table1", []string{"column1", "column2"}),
			RecordSet: RecordSet{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
				}),
			},
		},
		Where: parser.WhereClause{
			Filter: parser.Comparison{
				LHS:      parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
				RHS:      parser.NewIntegerValueFromString("2"),
				Operator: parser.Token{Token: '=', Literal: "="},
			},
		},
		Result: RecordSet{
			NewRecordWithId(2, []value.Primary{
				value.NewString("2"),
				value.NewString("str2"),
			}),
		},
	},
	{
		Name: "Where Filter Error",
		View: &View{
			Header: NewHeaderWithId("table1", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
				}),
			},
		},
		Where: parser.WhereClause{
			Filter: parser.Comparison{
				LHS:      parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
				RHS:      parser.NewIntegerValueFromString("2"),
				Operator: parser.Token{Token: '=', Literal: "="},
			},
		},
		Error: "field notexist does not exist",
	},
}

func TestView_Where(t *testing.T) {
	defer initFlag(TestTx.Flags)

	scope := NewReferenceScope(TestTx)
	ctx := context.Background()
	for _, v := range viewWhereTests {
		TestTx.Flags.CPU = 1
		if v.CPU != 0 {
			TestTx.Flags.CPU = v.CPU
		}

		err := v.View.Where(ctx, scope, v.Where)
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
		if !reflect.DeepEqual(v.View.RecordSet, v.Result) {
			t.Errorf("%s: result = %s, want %s", v.Name, v.View.RecordSet, v.Result)
		}
	}
}

var viewGroupByTests = []struct {
	Name       string
	View       *View
	GroupBy    parser.GroupByClause
	Result     *View
	IsGrouped  bool
	GroupItems []string
	Error      string
}{
	{
		Name: "Group By",
		View: &View{
			Header: NewHeaderWithId("table1", []string{"column1", "column2", "column3"}),
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
					value.NewString("group1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
					value.NewString("group2"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
					value.NewString("group1"),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("4"),
					value.NewString("str4"),
					value.NewString("group2"),
				}),
			},
		},
		GroupBy: parser.GroupByClause{
			Items: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column3"}},
			},
		},
		Result: &View{
			Header: []HeaderField{
				{
					View:   "table1",
					Column: InternalIdColumn,
				},
				{
					View:        "table1",
					Column:      "column1",
					Number:      1,
					IsFromTable: true,
				},
				{
					View:        "table1",
					Column:      "column2",
					Number:      2,
					IsFromTable: true,
				},
				{
					View:        "table1",
					Column:      "column3",
					Number:      3,
					IsFromTable: true,
					IsGroupKey:  true,
				},
			},
			RecordSet: []Record{
				{
					NewGroupCell([]value.Primary{value.NewInteger(1), value.NewInteger(3)}),
					NewGroupCell([]value.Primary{value.NewString("1"), value.NewString("3")}),
					NewGroupCell([]value.Primary{value.NewString("str1"), value.NewString("str3")}),
					NewGroupCell([]value.Primary{value.NewString("group1"), value.NewString("group1")}),
				},
				{
					NewGroupCell([]value.Primary{value.NewInteger(2), value.NewInteger(4)}),
					NewGroupCell([]value.Primary{value.NewString("2"), value.NewString("4")}),
					NewGroupCell([]value.Primary{value.NewString("str2"), value.NewString("str4")}),
					NewGroupCell([]value.Primary{value.NewString("group2"), value.NewString("group2")}),
				},
			},
			isGrouped: true,
		},
	},
	{
		Name: "Group By With ColumnNumber",
		View: &View{
			Header: NewHeaderWithId("table1", []string{"column1", "column2", "column3"}),
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
					value.NewString("group1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
					value.NewString("group2"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
					value.NewString("group1"),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("4"),
					value.NewString("str4"),
					value.NewString("group2"),
				}),
			},
		},
		GroupBy: parser.GroupByClause{
			Items: []parser.QueryExpression{
				parser.ColumnNumber{View: parser.Identifier{Literal: "table1"}, Number: value.NewInteger(3)},
			},
		},
		Result: &View{
			Header: []HeaderField{
				{
					View:   "table1",
					Column: InternalIdColumn,
				},
				{
					View:        "table1",
					Column:      "column1",
					Number:      1,
					IsFromTable: true,
				},
				{
					View:        "table1",
					Column:      "column2",
					Number:      2,
					IsFromTable: true,
				},
				{
					View:        "table1",
					Column:      "column3",
					Number:      3,
					IsFromTable: true,
					IsGroupKey:  true,
				},
			},
			RecordSet: []Record{
				{
					NewGroupCell([]value.Primary{value.NewInteger(1), value.NewInteger(3)}),
					NewGroupCell([]value.Primary{value.NewString("1"), value.NewString("3")}),
					NewGroupCell([]value.Primary{value.NewString("str1"), value.NewString("str3")}),
					NewGroupCell([]value.Primary{value.NewString("group1"), value.NewString("group1")}),
				},
				{
					NewGroupCell([]value.Primary{value.NewInteger(2), value.NewInteger(4)}),
					NewGroupCell([]value.Primary{value.NewString("2"), value.NewString("4")}),
					NewGroupCell([]value.Primary{value.NewString("str2"), value.NewString("str4")}),
					NewGroupCell([]value.Primary{value.NewString("group2"), value.NewString("group2")}),
				},
			},
			isGrouped: true,
		},
	},
	{
		Name: "Group By Evaluation Error",
		View: &View{
			Header: NewHeaderWithId("table1", []string{"column1", "column2", "column3"}),
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
					value.NewString("group1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
					value.NewString("group2"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
					value.NewString("group1"),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("4"),
					value.NewString("str4"),
					value.NewString("group2"),
				}),
			},
		},
		GroupBy: parser.GroupByClause{
			Items: []parser.QueryExpression{
				parser.ColumnNumber{View: parser.Identifier{Literal: "table1"}, Number: value.NewInteger(0)},
			},
		},
		Error: "field table1.0 does not exist",
	},
	{
		Name: "Group By Empty Record",
		View: &View{
			Header:    NewHeaderWithId("table1", []string{"column1", "column2", "column3"}),
			RecordSet: []Record{},
		},
		GroupBy: parser.GroupByClause{
			Items: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column3"}},
			},
		},
		Result: &View{
			Header: []HeaderField{
				{
					View:   "table1",
					Column: InternalIdColumn,
				},
				{
					View:        "table1",
					Column:      "column1",
					Number:      1,
					IsFromTable: true,
				},
				{
					View:        "table1",
					Column:      "column2",
					Number:      2,
					IsFromTable: true,
				},
				{
					View:        "table1",
					Column:      "column3",
					Number:      3,
					IsFromTable: true,
					IsGroupKey:  true,
				},
			},
			RecordSet: []Record{},
			isGrouped: true,
		},
	},
	{
		Name: "Group By Empty Record with No Condition",
		View: &View{
			Header:    NewHeaderWithId("table1", []string{"column1", "column2", "column3"}),
			RecordSet: []Record{},
		},
		GroupBy: parser.GroupByClause{
			Items: nil,
		},
		Result: &View{
			Header: []HeaderField{
				{
					View:   "table1",
					Column: InternalIdColumn,
				},
				{
					View:        "table1",
					Column:      "column1",
					Number:      1,
					IsFromTable: true,
				},
				{
					View:        "table1",
					Column:      "column2",
					Number:      2,
					IsFromTable: true,
				},
				{
					View:        "table1",
					Column:      "column3",
					Number:      3,
					IsFromTable: true,
				},
			},
			RecordSet: []Record{},
			isGrouped: true,
		},
	},
}

func TestView_GroupBy(t *testing.T) {
	scope := NewReferenceScope(TestTx)
	ctx := context.Background()
	for _, v := range viewGroupByTests {
		err := v.View.GroupBy(ctx, scope, v.GroupBy)
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
		if !reflect.DeepEqual(v.View, v.Result) {
			t.Errorf("%s: result = %v, want %v", v.Name, v.View, v.Result)
		}
	}
}

var viewHavingTests = []struct {
	Name   string
	View   *View
	Having parser.HavingClause
	Result RecordSet
	Error  string
}{
	{
		Name: "Having",
		View: &View{
			Header: []HeaderField{
				{
					View:   "table1",
					Column: InternalIdColumn,
				},
				{
					View:        "table1",
					Column:      "column1",
					IsFromTable: true,
				},
				{
					View:        "table1",
					Column:      "column2",
					IsFromTable: true,
				},
				{
					View:        "table1",
					Column:      "column3",
					IsFromTable: true,
					IsGroupKey:  true,
				},
			},
			isGrouped: true,
			RecordSet: RecordSet{
				{
					NewGroupCell([]value.Primary{value.NewString("1"), value.NewString("3")}),
					NewGroupCell([]value.Primary{value.NewString("1"), value.NewString("3")}),
					NewGroupCell([]value.Primary{value.NewString("str1"), value.NewString("str3")}),
					NewGroupCell([]value.Primary{value.NewString("group1"), value.NewString("group1")}),
				},
				{
					NewGroupCell([]value.Primary{value.NewString("2"), value.NewString("4")}),
					NewGroupCell([]value.Primary{value.NewString("2"), value.NewString("4")}),
					NewGroupCell([]value.Primary{value.NewString("str2"), value.NewString("str4")}),
					NewGroupCell([]value.Primary{value.NewString("group2"), value.NewString("group2")}),
				},
			},
		},
		Having: parser.HavingClause{
			Filter: parser.Comparison{
				LHS: parser.AggregateFunction{
					Name:     "sum",
					Distinct: parser.Token{},
					Args: []parser.QueryExpression{
						parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
					},
				},
				RHS:      parser.NewIntegerValueFromString("5"),
				Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: ">"},
			},
		},
		Result: RecordSet{
			{
				NewGroupCell([]value.Primary{value.NewString("2"), value.NewString("4")}),
				NewGroupCell([]value.Primary{value.NewString("2"), value.NewString("4")}),
				NewGroupCell([]value.Primary{value.NewString("str2"), value.NewString("str4")}),
				NewGroupCell([]value.Primary{value.NewString("group2"), value.NewString("group2")}),
			},
		},
	},
	{
		Name: "Having Filter Error",
		View: &View{
			Header: []HeaderField{
				{
					View:   "table1",
					Column: InternalIdColumn,
				},
				{
					View:        "table1",
					Column:      "column1",
					IsFromTable: true,
				},
				{
					View:        "table1",
					Column:      "column2",
					IsFromTable: true,
				},
				{
					View:        "table1",
					Column:      "column3",
					IsFromTable: true,
					IsGroupKey:  true,
				},
			},
			isGrouped: true,
			RecordSet: RecordSet{
				{
					NewGroupCell([]value.Primary{value.NewString("1"), value.NewString("3")}),
					NewGroupCell([]value.Primary{value.NewString("1"), value.NewString("3")}),
					NewGroupCell([]value.Primary{value.NewString("str1"), value.NewString("str3")}),
					NewGroupCell([]value.Primary{value.NewString("group1"), value.NewString("group1")}),
				},
				{
					NewGroupCell([]value.Primary{value.NewString("2"), value.NewString("4")}),
					NewGroupCell([]value.Primary{value.NewString("2"), value.NewString("4")}),
					NewGroupCell([]value.Primary{value.NewString("str2"), value.NewString("str4")}),
					NewGroupCell([]value.Primary{value.NewString("group2"), value.NewString("group2")}),
				},
			},
		},
		Having: parser.HavingClause{
			Filter: parser.Comparison{
				LHS: parser.AggregateFunction{
					Name:     "sum",
					Distinct: parser.Token{},
					Args: []parser.QueryExpression{
						parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
					},
				},
				RHS:      parser.NewIntegerValueFromString("5"),
				Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: ">"},
			},
		},
		Error: "field notexist does not exist",
	},
	{
		Name: "Having Not Grouped",
		View: &View{
			Header: NewHeaderWithId("table1", []string{"column1", "column2", "column3"}),
			RecordSet: RecordSet{
				NewRecordWithId(1, []value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
					value.NewString("group2"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("4"),
					value.NewString("str4"),
					value.NewString("group2"),
				}),
			},
		},
		Having: parser.HavingClause{
			Filter: parser.Comparison{
				LHS: parser.AggregateFunction{
					Name:     "sum",
					Distinct: parser.Token{},
					Args: []parser.QueryExpression{
						parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
					},
				},
				RHS:      parser.NewIntegerValueFromString("5"),
				Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: ">"},
			},
		},
		Result: RecordSet{
			{
				NewGroupCell([]value.Primary{value.NewInteger(1), value.NewInteger(2)}),
				NewGroupCell([]value.Primary{value.NewString("2"), value.NewString("4")}),
				NewGroupCell([]value.Primary{value.NewString("str2"), value.NewString("str4")}),
				NewGroupCell([]value.Primary{value.NewString("group2"), value.NewString("group2")}),
			},
		},
	},
	{
		Name: "Having All RecordSet Filter Error",
		View: &View{
			Header: NewHeaderWithId("table1", []string{"column1", "column2", "column3"}),
			RecordSet: RecordSet{
				NewRecordWithId(1, []value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
					value.NewString("group2"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("4"),
					value.NewString("str4"),
					value.NewString("group2"),
				}),
			},
		},
		Having: parser.HavingClause{
			Filter: parser.Comparison{
				LHS: parser.AggregateFunction{
					Name:     "sum",
					Distinct: parser.Token{},
					Args: []parser.QueryExpression{
						parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
					},
				},
				RHS:      parser.NewIntegerValueFromString("5"),
				Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: ">"},
			},
		},
		Error: "field notexist does not exist",
	},
}

func TestView_Having(t *testing.T) {
	scope := NewReferenceScope(TestTx)
	ctx := context.Background()
	for _, v := range viewHavingTests {
		err := v.View.Having(ctx, scope, v.Having)
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
		if !reflect.DeepEqual(v.View.RecordSet, v.Result) {
			t.Errorf("%s: result = %s, want %s", v.Name, v.View.RecordSet, v.Result)
		}
	}
}

var viewSelectTests = []struct {
	Name   string
	View   *View
	Scope  *ReferenceScope
	Select parser.SelectClause
	Result *View
	Error  string
}{
	{
		Name: "Select",
		View: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", Number: 1, IsFromTable: true},
				{View: "table1", Column: "column2", Number: 2, IsFromTable: true},
				{View: "table2", Column: InternalIdColumn},
				{View: "table2", Column: "column3", Number: 1, IsFromTable: true},
				{View: "table2", Column: "column4", Number: 2, IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewInteger(1),
					value.NewString("1"),
					value.NewString("str1"),
					value.NewInteger(1),
					value.NewString("2"),
					value.NewString("str22"),
				}),
				NewRecord([]value.Primary{
					value.NewInteger(1),
					value.NewString("1"),
					value.NewString("str1"),
					value.NewInteger(2),
					value.NewString("3"),
					value.NewString("str33"),
				}),
				NewRecord([]value.Primary{
					value.NewInteger(1),
					value.NewString("1"),
					value.NewString("str1"),
					value.NewInteger(3),
					value.NewString("1"),
					value.NewString("str44"),
				}),
				NewRecord([]value.Primary{
					value.NewInteger(2),
					value.NewString("2"),
					value.NewString("str2"),
					value.NewInteger(1),
					value.NewString("2"),
					value.NewString("str22"),
				}),
			},
		},
		Select: parser.SelectClause{
			Fields: []parser.QueryExpression{
				parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}, Alias: parser.Identifier{Literal: "c2"}},
				parser.Field{Object: parser.AllColumns{}},
				parser.Field{Object: parser.NewIntegerValueFromString("1"), Alias: parser.Identifier{Literal: "a"}},
				parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}, Alias: parser.Identifier{Literal: "c2a"}},
				parser.Field{Object: parser.ColumnNumber{View: parser.Identifier{Literal: "table2"}, Number: value.NewInteger(1)}, Alias: parser.Identifier{Literal: "t21"}},
				parser.Field{Object: parser.ColumnNumber{View: parser.Identifier{Literal: "table2"}, Number: value.NewInteger(1)}, Alias: parser.Identifier{Literal: "t21a"}},
				parser.Field{Object: parser.PrimitiveType{
					Literal: "2012-01-01",
					Value:   value.NewDatetime(time.Date(2012, 1, 1, 0, 0, 0, 0, GetTestLocation())),
				}},
			},
		},
		Result: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", Number: 1, IsFromTable: true},
				{View: "table1", Column: "column2", Aliases: []string{"c2", "c2a"}, Number: 2, IsFromTable: true},
				{View: "table2", Column: InternalIdColumn},
				{View: "table2", Column: "column3", Aliases: []string{"t21", "t21a"}, Number: 1, IsFromTable: true},
				{View: "table2", Column: "column4", Number: 2, IsFromTable: true},
				{Identifier: "@__PT:I:1", Column: "1", Aliases: []string{"a"}},
				{Identifier: "@__PT:D:2012-01-01T00:00:00Z", Column: "2012-01-01T00:00:00Z"},
			},
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewInteger(1),
					value.NewString("1"),
					value.NewString("str1"),
					value.NewInteger(1),
					value.NewString("2"),
					value.NewString("str22"),
					value.NewInteger(1),
					value.NewDatetime(time.Date(2012, 1, 1, 0, 0, 0, 0, GetTestLocation())),
				}),
				NewRecord([]value.Primary{
					value.NewInteger(1),
					value.NewString("1"),
					value.NewString("str1"),
					value.NewInteger(2),
					value.NewString("3"),
					value.NewString("str33"),
					value.NewInteger(1),
					value.NewDatetime(time.Date(2012, 1, 1, 0, 0, 0, 0, GetTestLocation())),
				}),
				NewRecord([]value.Primary{
					value.NewInteger(1),
					value.NewString("1"),
					value.NewString("str1"),
					value.NewInteger(3),
					value.NewString("1"),
					value.NewString("str44"),
					value.NewInteger(1),
					value.NewDatetime(time.Date(2012, 1, 1, 0, 0, 0, 0, GetTestLocation())),
				}),
				NewRecord([]value.Primary{
					value.NewInteger(2),
					value.NewString("2"),
					value.NewString("str2"),
					value.NewInteger(1),
					value.NewString("2"),
					value.NewString("str22"),
					value.NewInteger(1),
					value.NewDatetime(time.Date(2012, 1, 1, 0, 0, 0, 0, GetTestLocation())),
				}),
			},
			selectFields: []int{2, 1, 2, 4, 5, 6, 2, 4, 4, 7},
		},
	},
	{
		Name: "Select using Table Wildcard",
		View: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", Number: 1, IsFromTable: true},
				{View: "table1", Column: "column2", Number: 2, IsFromTable: true},
				{View: "table2", Column: InternalIdColumn},
				{View: "table2", Column: "column3", Number: 1, IsFromTable: true},
				{View: "table2", Column: "column4", Number: 2, IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewInteger(1),
					value.NewString("1"),
					value.NewString("str1"),
					value.NewInteger(1),
					value.NewString("2"),
					value.NewString("str22"),
				}),
				NewRecord([]value.Primary{
					value.NewInteger(1),
					value.NewString("1"),
					value.NewString("str1"),
					value.NewInteger(2),
					value.NewString("3"),
					value.NewString("str33"),
				}),
				NewRecord([]value.Primary{
					value.NewInteger(1),
					value.NewString("1"),
					value.NewString("str1"),
					value.NewInteger(3),
					value.NewString("1"),
					value.NewString("str44"),
				}),
				NewRecord([]value.Primary{
					value.NewInteger(2),
					value.NewString("2"),
					value.NewString("str2"),
					value.NewInteger(1),
					value.NewString("2"),
					value.NewString("str22"),
				}),
			},
		},
		Select: parser.SelectClause{
			Fields: []parser.QueryExpression{
				parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}, Alias: parser.Identifier{Literal: "c2"}},
				parser.Field{Object: parser.FieldReference{View: parser.Identifier{Literal: "table2"}, Column: parser.AllColumns{}}},
				parser.Field{Object: parser.NewIntegerValueFromString("1"), Alias: parser.Identifier{Literal: "a"}},
				parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}, Alias: parser.Identifier{Literal: "c2a"}},
				parser.Field{Object: parser.ColumnNumber{View: parser.Identifier{Literal: "table2"}, Number: value.NewInteger(1)}, Alias: parser.Identifier{Literal: "t21"}},
				parser.Field{Object: parser.ColumnNumber{View: parser.Identifier{Literal: "table2"}, Number: value.NewInteger(1)}, Alias: parser.Identifier{Literal: "t21a"}},
				parser.Field{Object: parser.PrimitiveType{
					Literal: "2012-01-01",
					Value:   value.NewDatetime(time.Date(2012, 1, 1, 0, 0, 0, 0, GetTestLocation())),
				}},
			},
		},
		Result: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", Number: 1, IsFromTable: true},
				{View: "table1", Column: "column2", Aliases: []string{"c2", "c2a"}, Number: 2, IsFromTable: true},
				{View: "table2", Column: InternalIdColumn},
				{View: "table2", Column: "column3", Aliases: []string{"t21", "t21a"}, Number: 1, IsFromTable: true},
				{View: "table2", Column: "column4", Number: 2, IsFromTable: true},
				{Identifier: "@__PT:I:1", Column: "1", Aliases: []string{"a"}},
				{Identifier: "@__PT:D:2012-01-01T00:00:00Z", Column: "2012-01-01T00:00:00Z"},
			},
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewInteger(1),
					value.NewString("1"),
					value.NewString("str1"),
					value.NewInteger(1),
					value.NewString("2"),
					value.NewString("str22"),
					value.NewInteger(1),
					value.NewDatetime(time.Date(2012, 1, 1, 0, 0, 0, 0, GetTestLocation())),
				}),
				NewRecord([]value.Primary{
					value.NewInteger(1),
					value.NewString("1"),
					value.NewString("str1"),
					value.NewInteger(2),
					value.NewString("3"),
					value.NewString("str33"),
					value.NewInteger(1),
					value.NewDatetime(time.Date(2012, 1, 1, 0, 0, 0, 0, GetTestLocation())),
				}),
				NewRecord([]value.Primary{
					value.NewInteger(1),
					value.NewString("1"),
					value.NewString("str1"),
					value.NewInteger(3),
					value.NewString("1"),
					value.NewString("str44"),
					value.NewInteger(1),
					value.NewDatetime(time.Date(2012, 1, 1, 0, 0, 0, 0, GetTestLocation())),
				}),
				NewRecord([]value.Primary{
					value.NewInteger(2),
					value.NewString("2"),
					value.NewString("str2"),
					value.NewInteger(1),
					value.NewString("2"),
					value.NewString("str22"),
					value.NewInteger(1),
					value.NewDatetime(time.Date(2012, 1, 1, 0, 0, 0, 0, GetTestLocation())),
				}),
			},
			selectFields: []int{2, 4, 5, 6, 2, 4, 4, 7},
		},
	},
	{
		Name: "Select Distinct",
		View: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
				{View: "table2", Column: InternalIdColumn},
				{View: "table2", Column: "column3", IsFromTable: true},
				{View: "table2", Column: "column4", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewInteger(1),
					value.NewString("1"),
					value.NewString("str1"),
					value.NewInteger(1),
					value.NewString("2"),
					value.NewString("str22"),
				}),
				NewRecord([]value.Primary{
					value.NewInteger(1),
					value.NewString("1"),
					value.NewString("str1"),
					value.NewInteger(2),
					value.NewString("3"),
					value.NewString("str33"),
				}),
				NewRecord([]value.Primary{
					value.NewInteger(1),
					value.NewString("1"),
					value.NewString("str1"),
					value.NewInteger(3),
					value.NewString("4"),
					value.NewString("str44"),
				}),
				NewRecord([]value.Primary{
					value.NewInteger(2),
					value.NewString("2"),
					value.NewString("str2"),
					value.NewInteger(1),
					value.NewString("2"),
					value.NewString("str22"),
				}),
			},
		},
		Select: parser.SelectClause{
			Distinct: parser.Token{Token: parser.DISTINCT, Literal: "distinct"},
			Fields: []parser.QueryExpression{
				parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
				parser.Field{Object: parser.NewIntegerValueFromString("1"), Alias: parser.Identifier{Literal: "a"}},
			},
		},
		Result: &View{
			Header: []HeaderField{
				{View: "table1", Column: "column1", IsFromTable: true},
				{Identifier: "@__PT:I:1", Column: "1", Aliases: []string{"a"}},
			},
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("1"),
					value.NewInteger(1),
				}),
				NewRecord([]value.Primary{
					value.NewString("2"),
					value.NewInteger(1),
				}),
			},
			selectFields: []int{0, 1},
		},
	},
	{
		Name: "Select Aggregate Function",
		View: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
				}),
			},
		},
		Select: parser.SelectClause{
			Fields: []parser.QueryExpression{
				parser.Field{
					Object: parser.AggregateFunction{
						Name:     "sum",
						Distinct: parser.Token{},
						Args: []parser.QueryExpression{
							parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
						},
					},
				},
			},
		},
		Result: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
				{Identifier: "SUM(column1)", Column: "SUM(column1)"},
			},
			RecordSet: []Record{
				{
					NewGroupCell([]value.Primary{value.NewInteger(1), value.NewInteger(2)}),
					NewGroupCell([]value.Primary{value.NewString("1"), value.NewString("2")}),
					NewGroupCell([]value.Primary{value.NewString("str1"), value.NewString("str2")}),
					NewCell(value.NewFloat(3)),
				},
			},
			selectFields: []int{3},
		},
	},
	{
		Name: "Select Aggregate Function Not Group Key Error",
		View: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
				}),
			},
		},
		Select: parser.SelectClause{
			Fields: []parser.QueryExpression{
				parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
				parser.Field{
					Object: parser.AggregateFunction{
						Name:     "sum",
						Distinct: parser.Token{},
						Args: []parser.QueryExpression{
							parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
						},
					},
				},
			},
		},
		Error: "field column2 is not a group key",
	},
	{
		Name: "Select Aggregate Function All RecordSet Lazy Evaluation",
		View: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
				}),
			},
		},
		Select: parser.SelectClause{
			Fields: []parser.QueryExpression{
				parser.Field{Object: parser.NewIntegerValueFromString("1")},
				parser.Field{
					Object: parser.Arithmetic{
						LHS: parser.AggregateFunction{
							Name:     "sum",
							Distinct: parser.Token{},
							Args: []parser.QueryExpression{
								parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
							},
						},
						RHS:      parser.NewIntegerValueFromString("1"),
						Operator: parser.Token{Token: '+', Literal: "+"},
					},
				},
			},
		},
		Result: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
				{Identifier: "@__PT:I:1", Column: "1"},
				{Identifier: "SUM(column1) + 1", Column: "SUM(column1) + 1"},
			},
			RecordSet: []Record{
				{
					NewGroupCell([]value.Primary{value.NewInteger(1), value.NewInteger(2)}),
					NewGroupCell([]value.Primary{value.NewString("1"), value.NewString("2")}),
					NewGroupCell([]value.Primary{value.NewString("str1"), value.NewString("str2")}),
					NewCell(value.NewInteger(1)),
					NewCell(value.NewFloat(4)),
				},
			},
			selectFields: []int{3, 4},
		},
	},
	{
		Name: "Select Analytic Function",
		View: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(2),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(3),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(5),
				}),
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(1),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(4),
				}),
			},
		},
		Select: parser.SelectClause{
			Fields: []parser.QueryExpression{
				parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
				parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
				parser.Field{
					Object: parser.AnalyticFunction{
						Name: "row_number",
						AnalyticClause: parser.AnalyticClause{
							PartitionClause: parser.PartitionClause{
								Values: []parser.QueryExpression{
									parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
								},
							},
							OrderByClause: parser.OrderByClause{
								Items: []parser.QueryExpression{
									parser.OrderItem{
										Value: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
									},
								},
							},
						},
					},
					Alias: parser.Identifier{Literal: "rownum"},
				},
			},
		},
		Result: &View{
			Header: []HeaderField{
				{View: "table1", Column: "column1", Number: 1, IsFromTable: true},
				{View: "table1", Column: "column2", Number: 2, IsFromTable: true},
				{Identifier: "ROW_NUMBER() OVER (PARTITION BY column1 ORDER BY column2)", Column: "ROW_NUMBER() OVER (PARTITION BY column1 ORDER BY column2)", Aliases: []string{"rownum"}},
			},
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(1),
					value.NewInteger(1),
				}),
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(2),
					value.NewInteger(2),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(3),
					value.NewInteger(1),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(4),
					value.NewInteger(2),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(5),
					value.NewInteger(3),
				}),
			},
			selectFields: []int{0, 1, 2},
		},
	},
	{
		Name: "Select Analytic Function Not Exist Error",
		View: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(2),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(3),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(5),
				}),
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(1),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(4),
				}),
			},
		},
		Select: parser.SelectClause{
			Fields: []parser.QueryExpression{
				parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
				parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
				parser.Field{
					Object: parser.AnalyticFunction{
						Name: "notexist",
						AnalyticClause: parser.AnalyticClause{
							PartitionClause: parser.PartitionClause{
								Values: []parser.QueryExpression{
									parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
								},
							},
							OrderByClause: parser.OrderByClause{
								Items: []parser.QueryExpression{
									parser.OrderItem{
										Value: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
									},
								},
							},
						},
					},
					Alias: parser.Identifier{Literal: "rownum"},
				},
			},
		},
		Error: "function notexist does not exist",
	},
	{
		Name: "Select Analytic Function Partition Error",
		View: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(2),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(3),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(5),
				}),
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(1),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(4),
				}),
			},
		},
		Select: parser.SelectClause{
			Fields: []parser.QueryExpression{
				parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
				parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
				parser.Field{
					Object: parser.AnalyticFunction{
						Name: "row_number",
						AnalyticClause: parser.AnalyticClause{
							PartitionClause: parser.PartitionClause{
								Values: []parser.QueryExpression{
									parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
								},
							},
							OrderByClause: parser.OrderByClause{
								Items: []parser.QueryExpression{
									parser.OrderItem{
										Value: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
									},
								},
							},
						},
					},
					Alias: parser.Identifier{Literal: "rownum"},
				},
			},
		},
		Error: "field notexist does not exist",
	},
	{
		Name: "Select Analytic Function Order Error",
		View: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(2),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(3),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(5),
				}),
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(1),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(4),
				}),
			},
		},
		Select: parser.SelectClause{
			Fields: []parser.QueryExpression{
				parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
				parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
				parser.Field{
					Object: parser.AnalyticFunction{
						Name: "row_number",
						AnalyticClause: parser.AnalyticClause{
							PartitionClause: parser.PartitionClause{
								Values: []parser.QueryExpression{
									parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
								},
							},
							OrderByClause: parser.OrderByClause{
								Items: []parser.QueryExpression{
									parser.OrderItem{
										Value: parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
									},
								},
							},
						},
					},
					Alias: parser.Identifier{Literal: "rownum"},
				},
			},
		},
		Error: "field notexist does not exist",
	},
	{
		Name: "Select User Defined Analytic Function",
		View: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(2),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(3),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(5),
				}),
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(1),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(4),
				}),
			},
		},
		Scope: GenerateReferenceScope([]map[string]map[string]interface{}{
			{
				scopeNameFunctions: {
					"USERAGGFUNC": &UserDefinedFunction{
						Name:         parser.Identifier{Literal: "useraggfunc"},
						IsAggregate:  true,
						Cursor:       parser.Identifier{Literal: "list"},
						RequiredArgs: 0,
						Statements: []parser.Statement{
							parser.VariableDeclaration{
								Assignments: []parser.VariableAssignment{
									{
										Variable: parser.Variable{Name: "value"},
									},
									{
										Variable: parser.Variable{Name: "fetch"},
									},
								},
							},
							parser.WhileInCursor{
								Variables: []parser.Variable{
									{Name: "fetch"},
								},
								Cursor: parser.Identifier{Literal: "list"},
								Statements: []parser.Statement{
									parser.If{
										Condition: parser.Is{
											LHS: parser.Variable{Name: "value"},
											RHS: parser.NewNullValue(),
										},
										Statements: []parser.Statement{
											parser.VariableSubstitution{
												Variable: parser.Variable{Name: "value"},
												Value:    parser.Variable{Name: "fetch"},
											},
											parser.FlowControl{Token: parser.CONTINUE},
										},
									},
									parser.VariableSubstitution{
										Variable: parser.Variable{Name: "value"},
										Value: parser.Arithmetic{
											LHS:      parser.Variable{Name: "value"},
											RHS:      parser.Variable{Name: "fetch"},
											Operator: parser.Token{Token: '+', Literal: "+"},
										},
									},
								},
							},

							parser.Return{
								Value: parser.Variable{Name: "value"},
							},
						},
					},
				},
			},
		}, nil, time.Time{}, nil),
		Select: parser.SelectClause{
			Fields: []parser.QueryExpression{
				parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
				parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
				parser.Field{
					Object: parser.AnalyticFunction{
						Name: "useraggfunc",
						Args: []parser.QueryExpression{
							parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
						},
						AnalyticClause: parser.AnalyticClause{},
					},
				},
			},
		},
		Result: &View{
			Header: []HeaderField{
				{View: "table1", Column: "column1", Number: 1, IsFromTable: true},
				{View: "table1", Column: "column2", Number: 2, IsFromTable: true},
				{Identifier: "USERAGGFUNC(column2) OVER ()", Column: "USERAGGFUNC(column2) OVER ()"},
			},
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(2),
					value.NewInteger(15),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(3),
					value.NewInteger(15),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(5),
					value.NewInteger(15),
				}),
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(1),
					value.NewInteger(15),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(4),
					value.NewInteger(15),
				}),
			},
			selectFields: []int{0, 1, 2},
		},
	},
	{
		Name: "Select Aggregate Empty Rows",
		View: &View{
			Header: []HeaderField{
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{},
		},
		Select: parser.SelectClause{
			Fields: []parser.QueryExpression{
				parser.Field{Object: parser.NewIntegerValue(1)},
				parser.Field{Object: parser.AggregateFunction{Name: "count", Args: []parser.QueryExpression{parser.AllColumns{}}}},
				parser.Field{Object: parser.AggregateFunction{Name: "sum", Args: []parser.QueryExpression{parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}}}},
			},
		},
		Result: &View{
			Header: []HeaderField{
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
				{Identifier: "@__PT:I:1", Column: "1"},
				{Identifier: "COUNT(*)", Column: "COUNT(*)"},
				{Identifier: "SUM(column1)", Column: "SUM(column1)"},
			},
			RecordSet: []Record{
				NewRecord([]value.Primary{
					nil,
					nil,
					value.NewInteger(1),
					value.NewInteger(0),
					value.NewNull(),
				}),
			},
			selectFields: []int{2, 3, 4},
		},
	},
	{
		Name: "Select Compound Function with Aggregate Empty Rows",
		View: &View{
			Header: []HeaderField{
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{},
		},
		Select: parser.SelectClause{
			Fields: []parser.QueryExpression{
				parser.Field{
					Object: parser.Function{Name: "coalesce",
						Args: []parser.QueryExpression{
							parser.AggregateFunction{Name: "sum", Args: []parser.QueryExpression{parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}}},
							parser.NewIntegerValue(0),
						},
					},
				},
			},
		},
		Result: &View{
			Header: []HeaderField{
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
				{Identifier: "COALESCE(SUM(column1), 0)", Column: "COALESCE(SUM(column1), 0)"},
			},
			RecordSet: []Record{
				NewRecord([]value.Primary{
					nil,
					nil,
					value.NewInteger(0),
				}),
			},
			selectFields: []int{2},
		},
	},
}

func TestView_Select(t *testing.T) {
	ctx := context.Background()
	for _, v := range viewSelectTests {
		if v.Scope == nil {
			v.Scope = NewReferenceScope(TestTx)
		}
		err := v.View.Select(ctx, v.Scope, v.Select)
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
		if !reflect.DeepEqual(v.View.Header, v.Result.Header) {
			t.Errorf("%s: header = %v, want %v", v.Name, v.View.Header, v.Result.Header)
		}
		if !reflect.DeepEqual(v.View.RecordSet, v.Result.RecordSet) {
			t.Errorf("%s: records = %s, want %s", v.Name, v.View.RecordSet, v.Result.RecordSet)
		}
		if !reflect.DeepEqual(v.View.selectFields, v.Result.selectFields) {
			t.Errorf("%s: select indices = %v, want %v", v.Name, v.View.selectFields, v.Result.selectFields)
		}
	}
}

var viewOrderByTests = []struct {
	Name    string
	View    *View
	OrderBy parser.OrderByClause
	Result  *View
	Error   string
}{
	{
		Name: "Order By",
		View: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
				{View: "table1", Column: "column3", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("3"),
					value.NewString("2"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("4"),
					value.NewString("3"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewString("1"),
					value.NewString("4"),
					value.NewString("3"),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("1"),
					value.NewString("3"),
					value.NewNull(),
				}),
				NewRecordWithId(5, []value.Primary{
					value.NewNull(),
					value.NewString("2"),
					value.NewString("4"),
				}),
			},
		},
		OrderBy: parser.OrderByClause{
			Items: []parser.QueryExpression{
				parser.OrderItem{
					Value: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
				},
				parser.OrderItem{
					Value:     parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
					Direction: parser.Token{Token: parser.DESC, Literal: "desc"},
				},
				parser.OrderItem{
					Value: parser.FieldReference{Column: parser.Identifier{Literal: "column3"}},
				},
				parser.OrderItem{
					Value: parser.NewIntegerValueFromString("1"),
				},
			},
		},
		Result: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
				{View: "table1", Column: "column3", IsFromTable: true},
				{Identifier: "@__PT:I:1", Column: "1"},
			},
			RecordSet: []Record{
				NewRecordWithId(5, []value.Primary{
					value.NewNull(),
					value.NewString("2"),
					value.NewString("4"),
					value.NewInteger(1),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("4"),
					value.NewString("3"),
					value.NewInteger(1),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewString("1"),
					value.NewString("4"),
					value.NewString("3"),
					value.NewInteger(1),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("1"),
					value.NewString("3"),
					value.NewNull(),
					value.NewInteger(1),
				}),
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("3"),
					value.NewString("2"),
					value.NewInteger(1),
				}),
			},
		},
	},
	{
		Name: "Order By with Cached SortValues",
		View: &View{
			Header: NewHeaderWithId("table1", []string{"column1", "column2", "column3"}),
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("3"),
					value.NewString("2"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("4"),
					value.NewString("3"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewString("1"),
					value.NewString("4"),
					value.NewString("3"),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("1"),
					value.NewString("3"),
					value.NewNull(),
				}),
				NewRecordWithId(5, []value.Primary{
					value.NewNull(),
					value.NewString("2"),
					value.NewString("4"),
				}),
			},
			sortValuesInEachCell: [][]*SortValue{
				{nil, nil, NewSortValue(value.NewString("3"), TestTx.Flags), nil},
				{nil, nil, NewSortValue(value.NewString("4"), TestTx.Flags), nil},
				{nil, nil, NewSortValue(value.NewString("4"), TestTx.Flags), nil},
				{nil, nil, NewSortValue(value.NewString("3"), TestTx.Flags), nil},
				{nil, nil, NewSortValue(value.NewString("2"), TestTx.Flags), nil},
			},
		},
		OrderBy: parser.OrderByClause{
			Items: []parser.QueryExpression{
				parser.OrderItem{
					Value: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
				},
			},
		},
		Result: &View{
			Header: NewHeaderWithId("table1", []string{"column1", "column2", "column3"}),
			RecordSet: []Record{
				NewRecordWithId(5, []value.Primary{
					value.NewNull(),
					value.NewString("2"),
					value.NewString("4"),
				}),
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("3"),
					value.NewString("2"),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("1"),
					value.NewString("3"),
					value.NewNull(),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("4"),
					value.NewString("3"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewString("1"),
					value.NewString("4"),
					value.NewString("3"),
				}),
			},
			sortValuesInEachCell: [][]*SortValue{
				{nil, nil, NewSortValue(value.NewString("2"), TestTx.Flags), nil},
				{nil, nil, NewSortValue(value.NewString("3"), TestTx.Flags), nil},
				{nil, nil, NewSortValue(value.NewString("3"), TestTx.Flags), nil},
				{nil, nil, NewSortValue(value.NewString("4"), TestTx.Flags), nil},
				{nil, nil, NewSortValue(value.NewString("4"), TestTx.Flags), nil},
			},
		},
	},
	{
		Name: "Order By With Null Positions",
		View: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewNull(),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("2"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewNull(),
					value.NewString("2"),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("1"),
					value.NewNull(),
				}),
				NewRecordWithId(5, []value.Primary{
					value.NewString("1"),
					value.NewString("2"),
				}),
			},
		},
		OrderBy: parser.OrderByClause{
			Items: []parser.QueryExpression{
				parser.OrderItem{
					Value:         parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
					NullsPosition: parser.Token{Token: parser.LAST, Literal: "last"},
				},
				parser.OrderItem{
					Value:         parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
					NullsPosition: parser.Token{Token: parser.FIRST, Literal: "first"},
				},
			},
		},
		Result: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewNull(),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("1"),
					value.NewNull(),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("2"),
				}),
				NewRecordWithId(5, []value.Primary{
					value.NewString("1"),
					value.NewString("2"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewNull(),
					value.NewString("2"),
				}),
			},
		},
	},
	{
		Name: "Order By Record Extend Error",
		View: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewNull(),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("2"),
				}),
			},
		},
		OrderBy: parser.OrderByClause{
			Items: []parser.QueryExpression{
				parser.OrderItem{
					Value: parser.AggregateFunction{
						Name:     "sum",
						Distinct: parser.Token{},
						Args: []parser.QueryExpression{
							parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
						},
					},
				},
			},
		},
		Error: "function sum cannot aggregate not grouping records",
	},
}

func TestView_OrderBy(t *testing.T) {
	scope := NewReferenceScope(TestTx)
	ctx := context.Background()
	for _, v := range viewOrderByTests {
		err := v.View.OrderBy(ctx, scope, v.OrderBy)
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
		if !reflect.DeepEqual(v.View.Header, v.Result.Header) {
			t.Errorf("%s: header = %v, want %v", v.Name, v.View.Header, v.Result.Header)
		}
		if !reflect.DeepEqual(v.View.RecordSet, v.Result.RecordSet) {
			t.Errorf("%s: records = %s, want %s", v.Name, v.View.RecordSet, v.Result.RecordSet)
		}
	}
}

var viewExtendRecordCapacityTests = []struct {
	Name   string
	View   *View
	Scope  *ReferenceScope
	Exprs  []parser.QueryExpression
	Result int
	Error  string
}{
	{
		Name: "ExtendRecordCapacity",
		View: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			RecordSet: RecordSet{
				NewRecord([]value.Primary{
					value.NewInteger(1),
					value.NewInteger(2),
				}),
			},
			isGrouped: true,
		},
		Scope: GenerateReferenceScope([]map[string]map[string]interface{}{
			{
				scopeNameFunctions: {
					"USERFUNC": &UserDefinedFunction{
						Name: parser.Identifier{Literal: "userfunc"},
						Parameters: []parser.Variable{
							{Name: "arg1"},
						},
						RequiredArgs: 1,
						Statements: []parser.Statement{
							parser.Return{Value: parser.Variable{Name: "arg1"}},
						},
						IsAggregate: true,
					},
				},
			},
		}, nil, time.Time{}, nil),
		Exprs: []parser.QueryExpression{
			parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
			parser.ColumnNumber{View: parser.Identifier{Literal: "table1"}, Number: value.NewInteger(2)},
			parser.Function{
				Name: "userfunc",
				Args: []parser.QueryExpression{
					parser.NewIntegerValueFromString("1"),
				},
			},
			parser.AggregateFunction{
				Name:     "avg",
				Distinct: parser.Token{},
				Args: []parser.QueryExpression{
					parser.AggregateFunction{
						Name: "avg",
						Args: []parser.QueryExpression{
							parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
						},
					},
				},
			},
			parser.AnalyticFunction{
				Name: "rank",
				AnalyticClause: parser.AnalyticClause{
					PartitionClause: parser.PartitionClause{
						Values: []parser.QueryExpression{
							parser.Arithmetic{
								LHS:      parser.NewIntegerValueFromString("1"),
								RHS:      parser.NewIntegerValueFromString("2"),
								Operator: parser.Token{Token: '+', Literal: "+"},
							},
						},
					},
					OrderByClause: parser.OrderByClause{
						Items: []parser.QueryExpression{
							parser.OrderItem{
								Value: parser.Arithmetic{
									LHS:      parser.NewIntegerValueFromString("3"),
									RHS:      parser.NewIntegerValueFromString("4"),
									Operator: parser.Token{Token: '+', Literal: "+"},
								},
							},
						},
					},
				},
			},
			parser.Arithmetic{
				LHS:      parser.NewIntegerValueFromString("5"),
				RHS:      parser.NewIntegerValueFromString("6"),
				Operator: parser.Token{Token: '+', Literal: "+"},
			},
		},
		Result: 8,
	},
}

func TestView_ExtendRecordCapacity(t *testing.T) {
	ctx := context.Background()
	for _, v := range viewExtendRecordCapacityTests {
		if v.Scope == nil {
			v.Scope = NewReferenceScope(TestTx)
		}

		err := v.View.ExtendRecordCapacity(ctx, v.Scope, v.Exprs, nil)
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
		if cap(v.View.RecordSet[0]) != v.Result {
			t.Errorf("%s: record capacity = %d, want %d", v.Name, cap(v.View.RecordSet[0]), v.Result)
		}
	}
}

var viewLimitTests = []struct {
	Name   string
	View   *View
	Limit  parser.LimitClause
	Result *View
	Error  string
}{
	{
		Name: "Limit",
		View: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
				}),
			},
		},
		Limit: parser.LimitClause{Value: parser.NewIntegerValueFromString("2")},
		Result: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
			},
		},
	},
	{
		Name: "Limit With Ties",
		View: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
				}),
			},
			sortValuesInEachRecord: []SortValues{
				{
					&SortValue{Type: IntegerType, Integer: 1},
					&SortValue{Type: StringType, String: "str1"},
				},
				{
					&SortValue{Type: IntegerType, Integer: 1},
					&SortValue{Type: StringType, String: "str1"},
				},
				{
					&SortValue{Type: IntegerType, Integer: 1},
					&SortValue{Type: StringType, String: "str1"},
				},
				{
					&SortValue{Type: IntegerType, Integer: 2},
					&SortValue{Type: StringType, String: "str2"},
				},
				{
					&SortValue{Type: IntegerType, Integer: 3},
					&SortValue{Type: StringType, String: "str3"},
				},
			},
		},
		Limit: parser.LimitClause{Value: parser.NewIntegerValueFromString("2"), Restriction: parser.Token{Token: parser.TIES}},
		Result: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
			},
			sortValuesInEachRecord: []SortValues{
				{
					&SortValue{Type: IntegerType, Integer: 1},
					&SortValue{Type: StringType, String: "str1"},
				},
				{
					&SortValue{Type: IntegerType, Integer: 1},
					&SortValue{Type: StringType, String: "str1"},
				},
				{
					&SortValue{Type: IntegerType, Integer: 1},
					&SortValue{Type: StringType, String: "str1"},
				},
				{
					&SortValue{Type: IntegerType, Integer: 2},
					&SortValue{Type: StringType, String: "str2"},
				},
				{
					&SortValue{Type: IntegerType, Integer: 3},
					&SortValue{Type: StringType, String: "str3"},
				},
			},
		},
	},
	{
		Name: "Limit By Percentage",
		View: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
				}),
			},
			offset: 1,
		},
		Limit: parser.LimitClause{Value: parser.NewFloatValue(50.5), Unit: parser.Token{Token: parser.PERCENT}},
		Result: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
				}),
			},
			offset: 1,
		},
	},
	{
		Name: "Limit By Over 100 Percentage",
		View: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
				}),
			},
		},
		Limit: parser.LimitClause{Value: parser.NewFloatValue(150), Unit: parser.Token{Token: parser.PERCENT}},
		Result: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
				}),
			},
		},
	},
	{
		Name: "Limit By Negative Percentage",
		View: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
				}),
			},
		},
		Limit: parser.LimitClause{Value: parser.NewFloatValue(-10), Unit: parser.Token{Token: parser.PERCENT}},
		Result: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{},
		},
	},
	{
		Name: "Limit Greater Than RecordSet",
		View: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
			},
		},
		Limit: parser.LimitClause{Value: parser.NewIntegerValueFromString("5")},
		Result: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
			},
		},
	},
	{
		Name: "Limit Evaluate Error",
		View: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
			},
		},
		Limit: parser.LimitClause{Value: parser.Variable{Name: "notexist"}},
		Error: "variable @notexist is undeclared",
	},
	{
		Name: "Limit Value Error",
		View: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
			},
		},
		Limit: parser.LimitClause{Value: parser.NewStringValue("str")},
		Error: "limit number of records 'str' is not an integer value",
	},
	{
		Name: "Limit Negative Value",
		View: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
			},
		},
		Limit: parser.LimitClause{Value: parser.NewIntegerValue(-1)},
		Result: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{},
		},
	},
	{
		Name: "Limit By Percentage Value Error",
		View: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
			},
		},
		Limit: parser.LimitClause{Value: parser.NewStringValue("str"), Unit: parser.Token{Token: parser.PERCENT}},
		Error: "limit percentage 'str' is not a float value",
	},
}

func TestView_Limit(t *testing.T) {
	scope := NewReferenceScope(TestTx)
	ctx := context.Background()
	for _, v := range viewLimitTests {
		err := v.View.Limit(ctx, scope, v.Limit)
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
		if !reflect.DeepEqual(v.View, v.Result) {
			t.Errorf("%s: view = %v, want %v", v.Name, v.View, v.Result)
		}
	}
}

var viewOffsetTests = []struct {
	Name   string
	View   *View
	Offset parser.OffsetClause
	Result *View
	Error  string
}{
	{
		Name: "Offset",
		View: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
				}),
			},
		},
		Offset: parser.OffsetClause{Value: parser.NewIntegerValueFromString("3")},
		Result: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(4, []value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
				}),
			},
			offset: 3,
		},
	},
	{
		Name: "Offset Equal To Record Length",
		View: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
				}),
			},
		},
		Offset: parser.OffsetClause{Value: parser.NewIntegerValueFromString("4")},
		Result: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{},
			offset:    4,
		},
	},
	{
		Name: "Offset Evaluate Error",
		View: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
				}),
			},
		},
		Offset: parser.OffsetClause{Value: parser.Variable{Name: "notexist"}},
		Error:  "variable @notexist is undeclared",
	},
	{
		Name: "Offset Value Error",
		View: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
				}),
			},
		},
		Offset: parser.OffsetClause{Value: parser.NewStringValue("str")},
		Error:  "offset number 'str' is not an integer value",
	},
	{
		Name: "Offset Negative Number",
		View: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
			},
		},
		Offset: parser.OffsetClause{Value: parser.NewIntegerValue(-3)},
		Result: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
			},
			offset: 0,
		},
	},
}

func TestView_Offset(t *testing.T) {
	scope := NewReferenceScope(TestTx)
	ctx := context.Background()
	for _, v := range viewOffsetTests {
		err := v.View.Offset(ctx, scope, v.Offset)
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
		if !reflect.DeepEqual(v.View, v.Result) {
			t.Errorf("%s: view = %v, want %v", v.Name, v.View, v.Result)
		}
	}
}

var viewInsertValuesTests = []struct {
	Name        string
	Fields      []parser.QueryExpression
	ValuesList  []parser.QueryExpression
	Result      *View
	UpdateCount int
	Error       string
}{
	{
		Name: "InsertValues",
		Fields: []parser.QueryExpression{
			parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
		},
		ValuesList: []parser.QueryExpression{
			parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewIntegerValueFromString("3"),
					},
				},
			},
			parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewIntegerValueFromString("4"),
					},
				},
			},
		},
		Result: &View{
			Header: NewHeaderWithId("table1", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewInteger(1),
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecord([]value.Primary{
					value.NewInteger(2),
					value.NewString("2"),
					value.NewString("str2"),
				}),
				NewRecord([]value.Primary{
					value.NewNull(),
					value.NewInteger(3),
					value.NewNull(),
				}),
				NewRecord([]value.Primary{
					value.NewNull(),
					value.NewInteger(4),
					value.NewNull(),
				}),
			},
		},
		UpdateCount: 2,
	},
	{
		Name: "InsertValues Field Length Does Not Match Error",
		Fields: []parser.QueryExpression{
			parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
			parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
		},
		ValuesList: []parser.QueryExpression{
			parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewIntegerValueFromString("3"),
					},
				},
			},
		},
		Error: "row value should contain exactly 2 values",
	},
	{
		Name: "InsertValues Value Evaluation Error",
		Fields: []parser.QueryExpression{
			parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
		},
		ValuesList: []parser.QueryExpression{
			parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
					},
				},
			},
		},
		Error: "field notexist does not exist",
	},
	{
		Name: "InsertValues Field Does Not Exist Error",
		Fields: []parser.QueryExpression{
			parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
		},
		ValuesList: []parser.QueryExpression{
			parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewIntegerValueFromString("3"),
					},
				},
			},
		},
		Error: "field notexist does not exist",
	},
}

func TestView_InsertValues(t *testing.T) {
	view := &View{
		Header: NewHeaderWithId("table1", []string{"column1", "column2"}),
		RecordSet: []Record{
			NewRecordWithId(1, []value.Primary{
				value.NewString("1"),
				value.NewString("str1"),
			}),
			NewRecordWithId(2, []value.Primary{
				value.NewString("2"),
				value.NewString("str2"),
			}),
		},
	}

	scope := NewReferenceScope(TestTx)
	ctx := context.Background()
	for _, v := range viewInsertValuesTests {
		cnt, err := view.InsertValues(ctx, scope, v.Fields, v.ValuesList)
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
		if !reflect.DeepEqual(view, v.Result) {
			t.Errorf("%s: result = %v, want %v", v.Name, view, v.Result)
		}
		if cnt != v.UpdateCount {
			t.Errorf("%s: update count = %d, want %d", v.Name, cnt, v.UpdateCount)
		}

	}
}

var viewInsertFromQueryTests = []struct {
	Name        string
	Fields      []parser.QueryExpression
	Query       parser.SelectQuery
	Result      *View
	UpdateCount int
	Error       string
}{
	{
		Name: "InsertFromQuery",
		Fields: []parser.QueryExpression{
			parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
		},
		Query: parser.SelectQuery{
			SelectEntity: parser.SelectEntity{
				SelectClause: parser.SelectClause{
					Fields: []parser.QueryExpression{
						parser.Field{Object: parser.NewIntegerValueFromString("3")},
					},
				},
			},
		},
		Result: &View{
			Header: NewHeaderWithId("table1", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewInteger(1),
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecord([]value.Primary{
					value.NewInteger(2),
					value.NewString("2"),
					value.NewString("str2"),
				}),
				NewRecord([]value.Primary{
					value.NewNull(),
					value.NewInteger(3),
					value.NewNull(),
				}),
			},
		},
		UpdateCount: 1,
	},
	{
		Name: "InsertFromQuery Field Lenght Does Not Match Error",
		Fields: []parser.QueryExpression{
			parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
			parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
		},
		Query: parser.SelectQuery{
			SelectEntity: parser.SelectEntity{
				SelectClause: parser.SelectClause{
					Fields: []parser.QueryExpression{
						parser.Field{Object: parser.NewIntegerValueFromString("3")},
					},
				},
			},
		},
		Error: "select query should return exactly 2 fields",
	},
	{
		Name: "InsertFromQuery Exuecution Error",
		Fields: []parser.QueryExpression{
			parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
		},
		Query: parser.SelectQuery{
			SelectEntity: parser.SelectEntity{
				SelectClause: parser.SelectClause{
					Fields: []parser.QueryExpression{
						parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}}},
					},
				},
			},
		},
		Error: "field notexist does not exist",
	},
}

func TestView_InsertFromQuery(t *testing.T) {
	view := &View{
		Header: NewHeaderWithId("table1", []string{"column1", "column2"}),
		RecordSet: []Record{
			NewRecordWithId(1, []value.Primary{
				value.NewString("1"),
				value.NewString("str1"),
			}),
			NewRecordWithId(2, []value.Primary{
				value.NewString("2"),
				value.NewString("str2"),
			}),
		},
	}

	scope := NewReferenceScope(TestTx)
	ctx := context.Background()
	for _, v := range viewInsertFromQueryTests {
		cnt, err := view.InsertFromQuery(ctx, scope, v.Fields, v.Query)
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
		if !reflect.DeepEqual(view, v.Result) {
			t.Errorf("%s: result = %v, want %v", v.Name, view, v.Result)
		}
		if cnt != v.UpdateCount {
			t.Errorf("%s: update count = %d, want %d", v.Name, cnt, v.UpdateCount)
		}
	}
}

var viewReplaceValuesTests = []struct {
	Name        string
	Fields      []parser.QueryExpression
	Keys        []parser.QueryExpression
	ValuesList  []parser.QueryExpression
	Result      *View
	UpdateCount int
	Error       string
}{
	{
		Name: "ReplaceValues",
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
						parser.NewIntegerValueFromString("1"),
						parser.NewStringValue("str3"),
					},
				},
			},
			parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewIntegerValueFromString("4"),
						parser.NewStringValue("str4"),
					},
				},
			},
		},
		Result: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("1"),
					value.NewString("str3"),
				}),
				NewRecord([]value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
				}),
				NewRecord([]value.Primary{
					value.NewInteger(4),
					value.NewString("str4"),
				}),
			},
		},
		UpdateCount: 2,
	},
	{
		Name: "ReplaceValues Field Length Does Not Match Error",
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
						parser.NewIntegerValueFromString("3"),
					},
				},
			},
		},
		Error: "row value should contain exactly 2 values",
	},
	{
		Name: "ReplaceValues Value Evaluation Error",
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
						parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
					},
				},
			},
		},
		Error: "field notexist does not exist",
	},
	{
		Name: "ReplaceValues Field Does Not Exist Error",
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
						parser.NewIntegerValueFromString("3"),
					},
				},
			},
		},
		Error: "field notexist does not exist",
	},
	{
		Name: "ReplaceValues Key Does Not Exist Error",
		Fields: []parser.QueryExpression{
			parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
		},
		Keys: []parser.QueryExpression{
			parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
		},
		ValuesList: []parser.QueryExpression{
			parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewIntegerValueFromString("3"),
					},
				},
			},
		},
		Error: "field notexist does not exist",
	},
	{
		Name: "ReplaceValues Key Not Set Error",
		Fields: []parser.QueryExpression{
			parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
		},
		Keys: []parser.QueryExpression{
			parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
		},
		ValuesList: []parser.QueryExpression{
			parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewIntegerValueFromString("3"),
					},
				},
			},
		},
		Error: "replace Key column2 is not set",
	},
}

func TestView_ReplaceValues(t *testing.T) {
	view := &View{
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
		},
	}

	scope := NewReferenceScope(TestTx)
	ctx := context.Background()
	for _, v := range viewReplaceValuesTests {
		cnt, err := view.ReplaceValues(ctx, scope, v.Fields, v.ValuesList, v.Keys)
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
		if !reflect.DeepEqual(view, v.Result) {
			t.Errorf("%s: result = %v, want %v", v.Name, view, v.Result)
		}
		if cnt != v.UpdateCount {
			t.Errorf("%s: update count = %d, want %d", v.Name, cnt, v.UpdateCount)
		}

	}
}

var viewReplaceFromQueryTests = []struct {
	Name        string
	Fields      []parser.QueryExpression
	Keys        []parser.QueryExpression
	Query       parser.SelectQuery
	Result      *View
	UpdateCount int
	Error       string
}{
	{
		Name: "ReplaceFromQuery",
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
						parser.Field{Object: parser.NewIntegerValueFromString("1")},
						parser.Field{Object: parser.NewStringValue("str3")},
					},
				},
			},
		},
		Result: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("1"),
					value.NewString("str3"),
				}),
				NewRecord([]value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
				}),
			},
		},
		UpdateCount: 1,
	},
	{
		Name: "ReplaceFromQuery Field Lenght Does Not Match Error",
		Fields: []parser.QueryExpression{
			parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
			parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
		},
		Query: parser.SelectQuery{
			SelectEntity: parser.SelectEntity{
				SelectClause: parser.SelectClause{
					Fields: []parser.QueryExpression{
						parser.Field{Object: parser.NewIntegerValueFromString("3")},
					},
				},
			},
		},
		Error: "select query should return exactly 2 fields",
	},
}

func TestView_ReplaceFromQuery(t *testing.T) {
	view := &View{
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
		},
	}

	scope := NewReferenceScope(TestTx)
	ctx := context.Background()
	for _, v := range viewReplaceFromQueryTests {
		cnt, err := view.ReplaceFromQuery(ctx, scope, v.Fields, v.Query, v.Keys)
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
		if !reflect.DeepEqual(view, v.Result) {
			t.Errorf("%s: result = %v, want %v", v.Name, view, v.Result)
		}
		if cnt != v.UpdateCount {
			t.Errorf("%s: update count = %d, want %d", v.Name, cnt, v.UpdateCount)
		}
	}
}

func TestView_Fix(t *testing.T) {
	view := &View{
		Header: []HeaderField{
			{View: "table1", Column: InternalIdColumn},
			{View: "table1", Column: "column1", IsFromTable: true},
			{View: "table1", Column: "column2", IsFromTable: true},
		},
		RecordSet: []Record{
			NewRecordWithId(1, []value.Primary{
				value.NewString("1"),
				value.NewString("str1"),
			}),
			NewRecordWithId(2, []value.Primary{
				value.NewString("1"),
				value.NewString("str1"),
			}),
		},
		selectFields: []int{2},
	}
	expect := &View{
		Header: NewHeader("table1", []string{"column2"}),
		RecordSet: []Record{
			NewRecord([]value.Primary{
				value.NewString("str1"),
			}),
			NewRecord([]value.Primary{
				value.NewString("str1"),
			}),
		},
		selectFields: []int(nil),
	}

	_ = view.Fix(context.Background(), TestTx.Flags)
	if !reflect.DeepEqual(view, expect) {
		t.Errorf("fix: view = %v, want %v", view, expect)
	}

	view = &View{
		Header: []HeaderField{
			{View: "table1", Column: InternalIdColumn},
			{View: "table1", Column: "column1", IsFromTable: true},
			{View: "table1", Column: "column2", IsFromTable: true},
		},
		RecordSet: []Record{
			NewRecordWithId(1, []value.Primary{
				value.NewString("1"),
				value.NewString("str1"),
			}),
			NewRecordWithId(2, []value.Primary{
				value.NewString("1"),
				value.NewString("str1"),
			}),
		},
		selectFields: []int{2, 2, 2, 2},
	}
	expect = &View{
		Header: NewHeader("table1", []string{"column2", "column2", "column2", "column2"}),
		RecordSet: []Record{
			NewRecord([]value.Primary{
				value.NewString("str1"),
				value.NewString("str1"),
				value.NewString("str1"),
				value.NewString("str1"),
			}),
			NewRecord([]value.Primary{
				value.NewString("str1"),
				value.NewString("str1"),
				value.NewString("str1"),
				value.NewString("str1"),
			}),
		},
		selectFields: []int(nil),
	}

	_ = view.Fix(context.Background(), TestTx.Flags)
	if !reflect.DeepEqual(view, expect) {
		t.Errorf("fix: view = %v, want %v", view, expect)
	}
}

func TestView_Union(t *testing.T) {
	view := &View{
		Header: []HeaderField{
			{View: "table1", Column: "column1", IsFromTable: true},
			{View: "table1", Column: "column2", IsFromTable: true},
		},
		RecordSet: RecordSet{
			NewRecord([]value.Primary{
				value.NewString("1"),
				value.NewString("str1"),
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
				value.NewString("2"),
				value.NewString("str2"),
			}),
		},
	}

	calcView := &View{
		Header: []HeaderField{
			{View: "table2", Column: "column3", IsFromTable: true},
			{View: "table2", Column: "column4", IsFromTable: true},
		},
		RecordSet: RecordSet{
			NewRecord([]value.Primary{
				value.NewString("2"),
				value.NewString("str2"),
			}),
			NewRecord([]value.Primary{
				value.NewString("3"),
				value.NewString("str3"),
			}),
		},
	}

	expect := &View{
		Header: []HeaderField{
			{View: "table1", Column: "column1", IsFromTable: true},
			{View: "table1", Column: "column2", IsFromTable: true},
		},
		RecordSet: RecordSet{
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
	}

	ctx := context.Background()
	err := view.Union(ctx, TestTx.Flags, calcView, false)
	if err != nil {
		t.Errorf("unexpected error %q", err)
	}
	if !reflect.DeepEqual(view, expect) {
		t.Errorf("union: view = %v, want %v", view, expect)
	}

	view = &View{
		Header: []HeaderField{
			{View: "table1", Column: "column1", IsFromTable: true},
			{View: "table1", Column: "column2", IsFromTable: true},
		},
		RecordSet: RecordSet{
			NewRecord([]value.Primary{
				value.NewString("1"),
				value.NewString("str1"),
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
				value.NewString("2"),
				value.NewString("str2"),
			}),
		},
	}

	expect = &View{
		Header: []HeaderField{
			{View: "table1", Column: "column1", IsFromTable: true},
			{View: "table1", Column: "column2", IsFromTable: true},
		},
		RecordSet: RecordSet{
			NewRecord([]value.Primary{
				value.NewString("1"),
				value.NewString("str1"),
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
				value.NewString("2"),
				value.NewString("str2"),
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
	}

	err = view.Union(ctx, TestTx.Flags, calcView, true)
	if err != nil {
		t.Errorf("unexpected error %q", err)
	}
	if !reflect.DeepEqual(view, expect) {
		t.Errorf("union all: view = %v, want %v", view, expect)
	}
}

func TestView_Except(t *testing.T) {
	view := &View{
		Header: []HeaderField{
			{View: "table1", Column: "column1", IsFromTable: true},
			{View: "table1", Column: "column2", IsFromTable: true},
		},
		RecordSet: RecordSet{
			NewRecord([]value.Primary{
				value.NewString("1"),
				value.NewString("str1"),
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
				value.NewString("2"),
				value.NewString("str2"),
			}),
		},
	}

	calcView := &View{
		Header: []HeaderField{
			{View: "table2", Column: "column3", IsFromTable: true},
			{View: "table2", Column: "column4", IsFromTable: true},
		},
		RecordSet: RecordSet{
			NewRecord([]value.Primary{
				value.NewString("2"),
				value.NewString("str2"),
			}),
			NewRecord([]value.Primary{
				value.NewString("3"),
				value.NewString("str3"),
			}),
		},
	}

	expect := &View{
		Header: []HeaderField{
			{View: "table1", Column: "column1", IsFromTable: true},
			{View: "table1", Column: "column2", IsFromTable: true},
		},
		RecordSet: RecordSet{
			NewRecord([]value.Primary{
				value.NewString("1"),
				value.NewString("str1"),
			}),
		},
	}

	ctx := context.Background()
	err := view.Except(ctx, TestTx.Flags, calcView, false)
	if err != nil {
		t.Errorf("unexpected error %q", err)
	}
	if !reflect.DeepEqual(view, expect) {
		t.Errorf("except: view = %v, want %v", view, expect)
	}

	view = &View{
		Header: []HeaderField{
			{View: "table1", Column: "column1", IsFromTable: true},
			{View: "table1", Column: "column2", IsFromTable: true},
		},
		RecordSet: RecordSet{
			NewRecord([]value.Primary{
				value.NewString("1"),
				value.NewString("str1"),
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
				value.NewString("2"),
				value.NewString("str2"),
			}),
		},
	}

	expect = &View{
		Header: []HeaderField{
			{View: "table1", Column: "column1", IsFromTable: true},
			{View: "table1", Column: "column2", IsFromTable: true},
		},
		RecordSet: RecordSet{
			NewRecord([]value.Primary{
				value.NewString("1"),
				value.NewString("str1"),
			}),
			NewRecord([]value.Primary{
				value.NewString("1"),
				value.NewString("str1"),
			}),
		},
	}

	err = view.Except(ctx, TestTx.Flags, calcView, true)
	if err != nil {
		t.Errorf("unexpected error %q", err)
	}
	if !reflect.DeepEqual(view, expect) {
		t.Errorf("except all: view = %v, want %v", view, expect)
	}
}

func TestView_Intersect(t *testing.T) {
	view := &View{
		Header: []HeaderField{
			{View: "table1", Column: "column1", IsFromTable: true},
			{View: "table1", Column: "column2", IsFromTable: true},
		},
		RecordSet: RecordSet{
			NewRecord([]value.Primary{
				value.NewString("1"),
				value.NewString("str1"),
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
				value.NewString("2"),
				value.NewString("str2"),
			}),
		},
	}

	calcView := &View{
		Header: []HeaderField{
			{View: "table2", Column: "column3", IsFromTable: true},
			{View: "table2", Column: "column4", IsFromTable: true},
		},
		RecordSet: RecordSet{
			NewRecord([]value.Primary{
				value.NewString("2"),
				value.NewString("str2"),
			}),
			NewRecord([]value.Primary{
				value.NewString("3"),
				value.NewString("str3"),
			}),
		},
	}

	expect := &View{
		Header: []HeaderField{
			{View: "table1", Column: "column1", IsFromTable: true},
			{View: "table1", Column: "column2", IsFromTable: true},
		},
		RecordSet: RecordSet{
			NewRecord([]value.Primary{
				value.NewString("2"),
				value.NewString("str2"),
			}),
		},
	}

	ctx := context.Background()
	err := view.Intersect(ctx, TestTx.Flags, calcView, false)
	if err != nil {
		t.Errorf("unexpected error %q", err)
	}
	if !reflect.DeepEqual(view, expect) {
		t.Errorf("intersect: view = %v, want %v", view, expect)
	}

	view = &View{
		Header: []HeaderField{
			{View: "table1", Column: "column1", IsFromTable: true},
			{View: "table1", Column: "column2", IsFromTable: true},
		},
		RecordSet: RecordSet{
			NewRecord([]value.Primary{
				value.NewString("1"),
				value.NewString("str1"),
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
				value.NewString("2"),
				value.NewString("str2"),
			}),
		},
	}

	expect = &View{
		Header: []HeaderField{
			{View: "table1", Column: "column1", IsFromTable: true},
			{View: "table1", Column: "column2", IsFromTable: true},
		},
		RecordSet: RecordSet{
			NewRecord([]value.Primary{
				value.NewString("2"),
				value.NewString("str2"),
			}),
			NewRecord([]value.Primary{
				value.NewString("2"),
				value.NewString("str2"),
			}),
		},
	}

	err = view.Intersect(ctx, TestTx.Flags, calcView, true)
	if err != nil {
		t.Errorf("unexpected error %q", err)
	}
	if !reflect.DeepEqual(view, expect) {
		t.Errorf("intersect all: view = %v, want %v", view, expect)
	}
}

func TestView_FieldIndex(t *testing.T) {
	view := &View{
		Header: []HeaderField{
			{View: "table1", Column: "column1", Number: 1, IsFromTable: true},
			{View: "table1", Column: "column2", Number: 2, IsFromTable: true},
		},
	}
	fieldRef := parser.FieldReference{
		Column: parser.Identifier{Literal: "column1"},
	}
	expect := 0

	idx, _ := view.FieldIndex(fieldRef)
	if idx != expect {
		t.Errorf("field index = %d, want %d", idx, expect)
	}

	columnNum := parser.ColumnNumber{
		View:   parser.Identifier{Literal: "table1"},
		Number: value.NewInteger(2),
	}
	expect = 1

	idx, _ = view.FieldIndex(columnNum)
	if idx != expect {
		t.Errorf("field index = %d, want %d", idx, expect)
	}
}

func TestView_FieldIndices(t *testing.T) {
	view := &View{
		Header: []HeaderField{
			{View: "table1", Column: "column1", IsFromTable: true},
			{View: "table1", Column: "column2", IsFromTable: true},
		},
	}
	fields := []parser.QueryExpression{
		parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
		parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
	}
	expect := []int{1, 0}

	indices, _ := view.FieldIndices(fields)
	if !reflect.DeepEqual(indices, expect) {
		t.Errorf("field indices = %v, want %v", indices, expect)
	}

	fields = []parser.QueryExpression{
		parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
		parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
	}
	expectErr := "field notexist does not exist"
	_, err := view.FieldIndices(fields)
	if err == nil {
		t.Errorf("no error, want error %s", expectErr)
	} else if err.Error() != expectErr {
		t.Errorf("error = %s, want %s", err, expectErr)
	}
}

func TestView_FieldViewName(t *testing.T) {
	view := &View{
		Header: []HeaderField{
			{View: "table1", Column: "column1", IsFromTable: true},
			{View: "table2", Column: "column2", IsFromTable: true},
		},
	}
	fieldRef := parser.FieldReference{
		Column: parser.Identifier{Literal: "column1"},
	}
	expect := "table1"

	ref, _ := view.FieldViewName(fieldRef)
	if ref != expect {
		t.Errorf("field reference = %s, want %s", ref, expect)
	}

	fieldRef = parser.FieldReference{
		Column: parser.Identifier{Literal: "notexist"},
	}
	expectErr := "field notexist does not exist"
	_, err := view.FieldViewName(fieldRef)
	if err == nil {
		t.Errorf("no error, want error %s", expectErr)
	} else if err.Error() != expectErr {
		t.Errorf("error = %s, want %s", err, expectErr)
	}
}

func TestView_InternalRecordId(t *testing.T) {
	view := &View{
		Header: NewHeaderWithId("table1", []string{"column1", "column2"}),
		RecordSet: []Record{
			NewRecordWithId(0, []value.Primary{value.NewInteger(1), value.NewString("str1")}),
			NewRecordWithId(1, []value.Primary{value.NewInteger(2), value.NewString("str2")}),
			NewRecordWithId(2, []value.Primary{value.NewInteger(3), value.NewString("str3")}),
		},
	}
	ref := "table1"
	recordIndex := 1
	expect := 1

	id, _ := view.InternalRecordId(ref, recordIndex)
	if id != expect {
		t.Errorf("field internal id = %d, want %d", id, expect)
	}

	view.RecordSet[1][0] = NewCell(value.NewNull())
	expectErr := "internal record id is empty"
	_, err := view.InternalRecordId(ref, recordIndex)
	if err == nil {
		t.Errorf("no error, want error %s", expectErr)
	} else if err.Error() != expectErr {
		t.Errorf("error = %s, want %s", err, expectErr)
	}

	view = &View{
		Header: []HeaderField{
			{View: "table1", Column: "column1", IsFromTable: true},
			{View: "table2", Column: "column2", IsFromTable: true},
		},
	}
	expectErr = "internal record id does not exist"
	_, err = view.InternalRecordId(ref, recordIndex)
	if err == nil {
		t.Errorf("no error, want error %s", expectErr)
	} else if err.Error() != expectErr {
		t.Errorf("error = %s, want %s", err, expectErr)
	}
}

func BenchmarkView_GroupBy(b *testing.B) {
	view := &View{
		Header:    NewHeader("t", []string{"c1", "c2", "c3"}),
		RecordSet: make(RecordSet, 10000),
	}
	for i := int64(0); i < 10000; i++ {
		view.RecordSet[i] = NewRecord([]value.Primary{
			value.NewInteger(i),
			value.NewString(randomStr(1)),
			value.NewString(randomStr(1)),
		})
	}

	ctx := context.Background()
	scope := NewReferenceScope(TestTx)
	clause := parser.GroupByClause{
		Items: []parser.QueryExpression{
			parser.FieldReference{Column: parser.Identifier{Literal: "c2"}},
			parser.FieldReference{Column: parser.Identifier{Literal: "c3"}},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v := &View{
			Header:    view.Header.Copy(),
			RecordSet: view.RecordSet.Copy(),
		}
		_ = v.GroupBy(ctx, scope, clause)
	}
}

func BenchmarkView_SelectDistinct(b *testing.B) {
	view := &View{
		Header:    NewHeader("t", []string{"c1", "c2", "c3"}),
		RecordSet: make(RecordSet, 10000),
	}
	for i := int64(0); i < 10000; i++ {
		view.RecordSet[i] = NewRecord([]value.Primary{
			value.NewInteger(i),
			value.NewString(randomStr(1)),
			value.NewString(randomStr(1)),
		})
	}

	ctx := context.Background()
	scope := NewReferenceScope(TestTx)
	clause := parser.SelectClause{
		Distinct: parser.Token{Token: parser.DISTINCT},
		Fields: []parser.QueryExpression{
			parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "c1"}}},
			parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "c2"}}},
			parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "c3"}}},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v := &View{
			Header:    view.Header.Copy(),
			RecordSet: view.RecordSet.Copy(),
		}
		_ = v.Select(ctx, scope, clause)
	}
}
