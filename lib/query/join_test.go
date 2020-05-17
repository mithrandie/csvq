package query

import (
	"context"
	"reflect"
	"testing"

	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/value"
)

func naturalJoinTestFieldReference(view string, column string) parser.FieldReference {
	return parser.FieldReference{
		BaseExpr: parser.NewBaseExpr(parser.Token{}),
		View:     parser.Identifier{Literal: view},
		Column:   parser.Identifier{BaseExpr: parser.NewBaseExpr(parser.Token{}), Literal: column},
	}
}

func joinUsingTestFieldReference(view string, column string) parser.FieldReference {
	return parser.FieldReference{
		View:   parser.Identifier{Literal: view},
		Column: parser.Identifier{Literal: column},
	}
}

var parseJoinConditionTests = []struct {
	Name          string
	Join          parser.Join
	View          *View
	JoinView      *View
	ResultValue   parser.QueryExpression
	IncludeFields []parser.FieldReference
	ExcludeFields []parser.FieldReference
	Error         string
}{
	{
		Name: "No Condition",
		Join: parser.Join{
			Table:     parser.Table{Alias: parser.Identifier{Literal: "t1"}},
			JoinTable: parser.Table{Alias: parser.Identifier{Literal: "t2"}},
		},
		View:        &View{Header: NewHeaderWithId("table1", []string{"key1", "key2", "key3", "value1", "value2", "value3"})},
		JoinView:    &View{Header: NewHeaderWithId("table2", []string{"key1", "key2", "key3", "value4"})},
		ResultValue: nil,
	},
	{
		Name: "Natural Join",
		Join: parser.Join{
			Table:     parser.Table{Alias: parser.Identifier{Literal: "t1"}},
			JoinTable: parser.Table{Alias: parser.Identifier{Literal: "t2"}},
			Natural:   parser.Token{Token: parser.NATURAL, Literal: "natural"},
		},
		View:     &View{Header: NewHeaderWithId("t1", []string{"key1", "key2", "key3", "value1", "value2", "value3"})},
		JoinView: &View{Header: NewHeaderWithId("t2", []string{"key1", "key2", "key3", "value4"})},
		ResultValue: parser.Logic{
			LHS: parser.Logic{
				LHS: parser.Comparison{
					LHS:      naturalJoinTestFieldReference("t1", "key1"),
					RHS:      naturalJoinTestFieldReference("t2", "key1"),
					Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: "="},
				},
				RHS: parser.Comparison{
					LHS:      naturalJoinTestFieldReference("t1", "key2"),
					RHS:      naturalJoinTestFieldReference("t2", "key2"),
					Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: "="},
				},
				Operator: parser.Token{Token: parser.AND, Literal: "AND"},
			},
			RHS: parser.Comparison{
				LHS:      naturalJoinTestFieldReference("t1", "key3"),
				RHS:      naturalJoinTestFieldReference("t2", "key3"),
				Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: "="},
			},
			Operator: parser.Token{Token: parser.AND, Literal: "AND"},
		},
		IncludeFields: []parser.FieldReference{
			naturalJoinTestFieldReference("t1", "key1"),
			naturalJoinTestFieldReference("t1", "key2"),
			naturalJoinTestFieldReference("t1", "key3"),
		},
		ExcludeFields: []parser.FieldReference{
			naturalJoinTestFieldReference("t2", "key1"),
			naturalJoinTestFieldReference("t2", "key2"),
			naturalJoinTestFieldReference("t2", "key3"),
		},
	},
	{
		Name: "Using Condition",
		Join: parser.Join{
			Table:     parser.Table{Alias: parser.Identifier{Literal: "t1"}},
			JoinTable: parser.Table{Alias: parser.Identifier{Literal: "t2"}},
			Condition: parser.JoinCondition{
				Using: []parser.QueryExpression{
					parser.Identifier{Literal: "key1"},
				},
			},
		},
		View:     &View{Header: NewHeaderWithId("t1", []string{"key1", "key2", "key3", "value1", "value2", "value3"})},
		JoinView: &View{Header: NewHeaderWithId("t2", []string{"key1", "key2", "key3", "value4"})},
		ResultValue: parser.Comparison{
			LHS:      joinUsingTestFieldReference("t1", "key1"),
			RHS:      joinUsingTestFieldReference("t2", "key1"),
			Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: "="},
		},
		IncludeFields: []parser.FieldReference{
			joinUsingTestFieldReference("t1", "key1"),
		},
		ExcludeFields: []parser.FieldReference{
			joinUsingTestFieldReference("t2", "key1"),
		},
	},
	{
		Name: "Right Outer Join Using Condition",
		Join: parser.Join{
			Table:     parser.Table{Alias: parser.Identifier{Literal: "t1"}},
			JoinTable: parser.Table{Alias: parser.Identifier{Literal: "t2"}},
			JoinType:  parser.Token{Token: parser.OUTER, Literal: "outer"},
			Direction: parser.Token{Token: parser.RIGHT, Literal: "right"},
			Condition: parser.JoinCondition{
				Using: []parser.QueryExpression{
					parser.Identifier{Literal: "key1"},
				},
			},
		},
		View:     &View{Header: NewHeaderWithId("t1", []string{"key1", "key2", "key3", "value1", "value2", "value3"})},
		JoinView: &View{Header: NewHeaderWithId("t2", []string{"key1", "key2", "key3", "value4"})},
		ResultValue: parser.Comparison{
			LHS:      joinUsingTestFieldReference("t1", "key1"),
			RHS:      joinUsingTestFieldReference("t2", "key1"),
			Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: "="},
		},
		IncludeFields: []parser.FieldReference{
			joinUsingTestFieldReference("t2", "key1"),
		},
		ExcludeFields: []parser.FieldReference{
			joinUsingTestFieldReference("t1", "key1"),
		},
	},
	{
		Name: "On Condition",
		Join: parser.Join{
			Table:     parser.Table{Alias: parser.Identifier{Literal: "t1"}},
			JoinTable: parser.Table{Alias: parser.Identifier{Literal: "t2"}},
			Condition: parser.JoinCondition{
				On: parser.Comparison{
					LHS:      parser.FieldReference{View: parser.Identifier{Literal: "t1"}, Column: parser.Identifier{Literal: "key1"}},
					RHS:      parser.FieldReference{View: parser.Identifier{Literal: "t2"}, Column: parser.Identifier{Literal: "key1"}},
					Operator: parser.Token{Token: '=', Literal: "="},
				},
			},
		},
		View:     &View{Header: NewHeaderWithId("table1", []string{"key1", "key2", "key3", "value1", "value2", "value3"})},
		JoinView: &View{Header: NewHeaderWithId("table2", []string{"key1", "key2", "key3", "value4"})},
		ResultValue: parser.Comparison{
			LHS:      parser.FieldReference{View: parser.Identifier{Literal: "t1"}, Column: parser.Identifier{Literal: "key1"}},
			RHS:      parser.FieldReference{View: parser.Identifier{Literal: "t2"}, Column: parser.Identifier{Literal: "key1"}},
			Operator: parser.Token{Token: '=', Literal: "="},
		},
	},
	{
		Name: "Natural Join Fields Does Not Duplicate",
		Join: parser.Join{
			Table:     parser.Table{Alias: parser.Identifier{Literal: "t1"}},
			JoinTable: parser.Table{Alias: parser.Identifier{Literal: "t2"}},
			Natural:   parser.Token{Token: parser.NATURAL, Literal: "natural"},
		},
		View:        &View{Header: NewHeaderWithId("table1", []string{"value1", "value2", "value3"})},
		JoinView:    &View{Header: NewHeaderWithId("table2", []string{"value4"})},
		ResultValue: nil,
	},
	{
		Name: "Using Condition View Field Error",
		Join: parser.Join{
			Table:     parser.Table{Alias: parser.Identifier{Literal: "t1"}},
			JoinTable: parser.Table{Alias: parser.Identifier{Literal: "t2"}},
			Condition: parser.JoinCondition{
				Using: []parser.QueryExpression{
					parser.Identifier{Literal: "key1"},
				},
			},
		},
		View:     &View{Header: NewHeaderWithId("t1", []string{"key1", "key2", "key3", "key1", "value1", "value2", "value3"})},
		JoinView: &View{Header: NewHeaderWithId("t2", []string{"key1", "key2", "key3", "value4"})},
		Error:    "field key1 is ambiguous",
	},
	{
		Name: "Using Condition JoinView Field Error",
		Join: parser.Join{
			Table:     parser.Table{Alias: parser.Identifier{Literal: "t1"}},
			JoinTable: parser.Table{Alias: parser.Identifier{Literal: "t2"}},
			Condition: parser.JoinCondition{
				Using: []parser.QueryExpression{
					parser.Identifier{Literal: "key1"},
				},
			},
		},
		View:     &View{Header: NewHeaderWithId("t1", []string{"key1", "key2", "key3", "value1", "value2", "value3"})},
		JoinView: &View{Header: NewHeaderWithId("t2", []string{"key2", "key3", "value4"})},
		Error:    "field key1 does not exist",
	},
}

func TestParseJoinCondition(t *testing.T) {
	for _, v := range parseJoinConditionTests {
		r, ifields, xfields, err := ParseJoinCondition(v.Join, v.View, v.JoinView)
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
		if !reflect.DeepEqual(r, v.ResultValue) {
			t.Errorf("%s: condition = %q, want %q", v.Name, r, v.ResultValue)
		}
		if !reflect.DeepEqual(ifields, v.IncludeFields) {
			t.Errorf("%s: include fields = %q, want %q", v.Name, ifields, v.IncludeFields)
		}
		if !reflect.DeepEqual(xfields, v.ExcludeFields) {
			t.Errorf("%s: exclude fields = %q, want %q", v.Name, xfields, v.ExcludeFields)
		}
	}
}

func TestCrossJoin(t *testing.T) {
	view := &View{
		Header: NewHeaderWithId("table1", []string{"column1", "column2"}),
		RecordSet: []Record{
			NewRecordWithId(1, []value.Primary{
				value.NewInteger(1),
				value.NewString("str1"),
			}),
			NewRecordWithId(2, []value.Primary{
				value.NewInteger(2),
				value.NewString("str2"),
			}),
		},
	}
	joinView := &View{
		Header: NewHeaderWithId("table2", []string{"column3", "column4"}),
		RecordSet: []Record{
			NewRecordWithId(1, []value.Primary{
				value.NewInteger(3),
				value.NewString("str3"),
			}),
			NewRecordWithId(2, []value.Primary{
				value.NewInteger(4),
				value.NewString("str4"),
			}),
		},
	}
	expect := &View{
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
				value.NewInteger(1),
				value.NewString("str1"),
				value.NewInteger(1),
				value.NewInteger(3),
				value.NewString("str3"),
			}),
			NewRecord([]value.Primary{
				value.NewInteger(1),
				value.NewInteger(1),
				value.NewString("str1"),
				value.NewInteger(2),
				value.NewInteger(4),
				value.NewString("str4"),
			}),
			NewRecord([]value.Primary{
				value.NewInteger(2),
				value.NewInteger(2),
				value.NewString("str2"),
				value.NewInteger(1),
				value.NewInteger(3),
				value.NewString("str3"),
			}),
			NewRecord([]value.Primary{
				value.NewInteger(2),
				value.NewInteger(2),
				value.NewString("str2"),
				value.NewInteger(2),
				value.NewInteger(4),
				value.NewString("str4"),
			}),
		},
	}

	_ = CrossJoin(context.Background(), NewReferenceScope(TestTx), view, joinView)
	if !reflect.DeepEqual(view, expect) {
		t.Errorf("Cross Join: result = %v, want %v", view, expect)
	}
}

var innerJoinTests = []struct {
	Name      string
	CPU       int
	View      *View
	JoinView  *View
	Condition parser.QueryExpression
	Scope     *ReferenceScope
	Result    *View
	Error     string
}{
	{
		Name: "Inner Join",
		View: &View{
			Header: NewHeaderWithId("table1", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewInteger(1),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewInteger(2),
					value.NewString("str2"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewInteger(3),
					value.NewString("str3"),
				}),
			},
		},
		JoinView: &View{
			Header: NewHeaderWithId("table2", []string{"column1", "column3"}),
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewInteger(1),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewInteger(2),
					value.NewString("str22"),
				}),
			},
		},
		Condition: parser.Comparison{
			LHS:      parser.FieldReference{View: parser.Identifier{Literal: "table1"}, Column: parser.Identifier{Literal: "column1"}},
			RHS:      parser.FieldReference{View: parser.Identifier{Literal: "table2"}, Column: parser.Identifier{Literal: "column1"}},
			Operator: parser.Token{Token: '=', Literal: "="},
		},
		Result: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", Number: 1, IsFromTable: true},
				{View: "table1", Column: "column2", Number: 2, IsFromTable: true},
				{View: "table2", Column: InternalIdColumn},
				{View: "table2", Column: "column1", Number: 1, IsFromTable: true},
				{View: "table2", Column: "column3", Number: 2, IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewInteger(1),
					value.NewInteger(1),
					value.NewString("str1"),
					value.NewInteger(1),
					value.NewInteger(1),
					value.NewString("str1"),
				}),
				NewRecord([]value.Primary{
					value.NewInteger(2),
					value.NewInteger(2),
					value.NewString("str2"),
					value.NewInteger(2),
					value.NewInteger(2),
					value.NewString("str22"),
				}),
			},
		},
	},
	{
		Name: "Inner Join in Multi Threading",
		CPU:  2,
		View: &View{
			Header: NewHeaderWithId("table1", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewInteger(1),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewInteger(2),
					value.NewString("str2"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewInteger(3),
					value.NewString("str3"),
				}),
			},
		},
		JoinView: &View{
			Header: NewHeaderWithId("table2", []string{"column1", "column3"}),
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewInteger(1),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewInteger(2),
					value.NewString("str22"),
				}),
			},
		},
		Condition: parser.Comparison{
			LHS:      parser.FieldReference{View: parser.Identifier{Literal: "table1"}, Column: parser.Identifier{Literal: "column1"}},
			RHS:      parser.FieldReference{View: parser.Identifier{Literal: "table2"}, Column: parser.Identifier{Literal: "column1"}},
			Operator: parser.Token{Token: '=', Literal: "="},
		},
		Result: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", Number: 1, IsFromTable: true},
				{View: "table1", Column: "column2", Number: 2, IsFromTable: true},
				{View: "table2", Column: InternalIdColumn},
				{View: "table2", Column: "column1", Number: 1, IsFromTable: true},
				{View: "table2", Column: "column3", Number: 2, IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewInteger(1),
					value.NewInteger(1),
					value.NewString("str1"),
					value.NewInteger(1),
					value.NewInteger(1),
					value.NewString("str1"),
				}),
				NewRecord([]value.Primary{
					value.NewInteger(2),
					value.NewInteger(2),
					value.NewString("str2"),
					value.NewInteger(2),
					value.NewInteger(2),
					value.NewString("str22"),
				}),
			},
		},
	},
	{
		Name: "Inner Join With No Condition",
		View: &View{
			Header: NewHeaderWithId("table1", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewInteger(1),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewInteger(2),
					value.NewString("str2"),
				}),
			},
		},
		JoinView: &View{
			Header: NewHeaderWithId("table2", []string{"column3", "column4"}),
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewInteger(3),
					value.NewString("str3"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewInteger(4),
					value.NewString("str4"),
				}),
			},
		},
		Condition: nil,
		Result: &View{
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
					value.NewInteger(1),
					value.NewString("str1"),
					value.NewInteger(1),
					value.NewInteger(3),
					value.NewString("str3"),
				}),
				NewRecord([]value.Primary{
					value.NewInteger(1),
					value.NewInteger(1),
					value.NewString("str1"),
					value.NewInteger(2),
					value.NewInteger(4),
					value.NewString("str4"),
				}),
				NewRecord([]value.Primary{
					value.NewInteger(2),
					value.NewInteger(2),
					value.NewString("str2"),
					value.NewInteger(1),
					value.NewInteger(3),
					value.NewString("str3"),
				}),
				NewRecord([]value.Primary{
					value.NewInteger(2),
					value.NewInteger(2),
					value.NewString("str2"),
					value.NewInteger(2),
					value.NewInteger(4),
					value.NewString("str4"),
				}),
			},
		},
	},
	{
		Name: "Inner Join Filter Error",
		View: &View{
			Header: NewHeaderWithId("table1", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewInteger(1),
					value.NewString("str1"),
				}),
				NewRecordWithId(1, []value.Primary{
					value.NewInteger(2),
					value.NewString("str2"),
				}),
				NewRecordWithId(1, []value.Primary{
					value.NewInteger(3),
					value.NewString("str3"),
				}),
				NewRecordWithId(1, []value.Primary{
					value.NewInteger(4),
					value.NewString("str4"),
				}),
			},
		},
		JoinView: &View{
			Header: NewHeaderWithId("table2", []string{"column1", "column3"}),
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewInteger(1),
					value.NewString("str1"),
				}),
			},
		},
		Condition: parser.Comparison{
			LHS:      parser.FieldReference{View: parser.Identifier{Literal: "table1"}, Column: parser.Identifier{Literal: "column1"}},
			RHS:      parser.FieldReference{View: parser.Identifier{Literal: "table2"}, Column: parser.Identifier{Literal: "notexist"}},
			Operator: parser.Token{Token: '=', Literal: "="},
		},
		Error: "field table2.notexist does not exist",
	},
}

func TestInnerJoin(t *testing.T) {
	defer initFlag(TestTx.Flags)

	for _, v := range innerJoinTests {
		TestTx.Flags.CPU = 1
		if v.CPU != 0 {
			TestTx.Flags.CPU = v.CPU
		}

		if v.Scope == nil {
			v.Scope = NewReferenceScope(TestTx)
		}

		err := InnerJoin(context.Background(), v.Scope, v.View, v.JoinView, v.Condition)
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

var outerJoinTests = []struct {
	Name      string
	View      *View
	JoinView  *View
	Condition parser.QueryExpression
	Direction int
	Scope     *ReferenceScope
	Result    *View
	Error     string
}{
	{
		Name: "Left Outer Join",
		View: &View{
			Header: NewHeaderWithId("table1", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewInteger(1),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewInteger(2),
					value.NewString("str2"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewInteger(3),
					value.NewString("str3"),
				}),
			},
		},
		JoinView: &View{
			Header: NewHeaderWithId("table2", []string{"column1", "column3"}),
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewInteger(2),
					value.NewString("str22"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewInteger(3),
					value.NewString("str33"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewInteger(4),
					value.NewString("str44"),
				}),
			},
		},
		Condition: parser.Comparison{
			LHS:      parser.FieldReference{View: parser.Identifier{Literal: "table1"}, Column: parser.Identifier{Literal: "column1"}},
			RHS:      parser.FieldReference{View: parser.Identifier{Literal: "table2"}, Column: parser.Identifier{Literal: "column1"}},
			Operator: parser.Token{Token: '=', Literal: "="},
		},
		Direction: parser.LEFT,
		Result: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", Number: 1, IsFromTable: true},
				{View: "table1", Column: "column2", Number: 2, IsFromTable: true},
				{View: "table2", Column: InternalIdColumn},
				{View: "table2", Column: "column1", Number: 1, IsFromTable: true},
				{View: "table2", Column: "column3", Number: 2, IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewInteger(1),
					value.NewInteger(1),
					value.NewString("str1"),
					value.NewNull(),
					value.NewNull(),
					value.NewNull(),
				}),
				NewRecord([]value.Primary{
					value.NewInteger(2),
					value.NewInteger(2),
					value.NewString("str2"),
					value.NewInteger(1),
					value.NewInteger(2),
					value.NewString("str22"),
				}),
				NewRecord([]value.Primary{
					value.NewInteger(3),
					value.NewInteger(3),
					value.NewString("str3"),
					value.NewInteger(2),
					value.NewInteger(3),
					value.NewString("str33"),
				}),
			},
		},
	},
	{
		Name: "Right Outer Join",
		View: &View{
			Header: NewHeaderWithId("table1", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewInteger(1),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewInteger(2),
					value.NewString("str2"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewInteger(3),
					value.NewString("str3"),
				}),
			},
		},
		JoinView: &View{
			Header: NewHeaderWithId("table2", []string{"column1", "column3"}),
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewInteger(2),
					value.NewString("str22"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewInteger(3),
					value.NewString("str33"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewInteger(4),
					value.NewString("str44"),
				}),
			},
		},
		Condition: parser.Comparison{
			LHS:      parser.FieldReference{View: parser.Identifier{Literal: "table1"}, Column: parser.Identifier{Literal: "column1"}},
			RHS:      parser.FieldReference{View: parser.Identifier{Literal: "table2"}, Column: parser.Identifier{Literal: "column1"}},
			Operator: parser.Token{Token: '=', Literal: "="},
		},
		Direction: parser.RIGHT,
		Result: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", Number: 1, IsFromTable: true},
				{View: "table1", Column: "column2", Number: 2, IsFromTable: true},
				{View: "table2", Column: InternalIdColumn},
				{View: "table2", Column: "column1", Number: 1, IsFromTable: true},
				{View: "table2", Column: "column3", Number: 2, IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewInteger(2),
					value.NewInteger(2),
					value.NewString("str2"),
					value.NewInteger(1),
					value.NewInteger(2),
					value.NewString("str22"),
				}),
				NewRecord([]value.Primary{
					value.NewInteger(3),
					value.NewInteger(3),
					value.NewString("str3"),
					value.NewInteger(2),
					value.NewInteger(3),
					value.NewString("str33"),
				}),
				NewRecord([]value.Primary{
					value.NewNull(),
					value.NewNull(),
					value.NewNull(),
					value.NewInteger(3),
					value.NewInteger(4),
					value.NewString("str44"),
				}),
			},
		},
	},
	{
		Name: "Full Outer Join",
		View: &View{
			Header: NewHeaderWithId("table1", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewInteger(1),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewInteger(2),
					value.NewString("str2"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewInteger(3),
					value.NewString("str3"),
				}),
			},
		},
		JoinView: &View{
			Header: NewHeaderWithId("table2", []string{"column1", "column3"}),
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewInteger(2),
					value.NewString("str22"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewInteger(3),
					value.NewString("str33"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewInteger(4),
					value.NewString("str44"),
				}),
			},
		},
		Condition: parser.Comparison{
			LHS:      parser.FieldReference{View: parser.Identifier{Literal: "table1"}, Column: parser.Identifier{Literal: "column1"}},
			RHS:      parser.FieldReference{View: parser.Identifier{Literal: "table2"}, Column: parser.Identifier{Literal: "column1"}},
			Operator: parser.Token{Token: '=', Literal: "="},
		},
		Direction: parser.FULL,
		Result: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", Number: 1, IsFromTable: true},
				{View: "table1", Column: "column2", Number: 2, IsFromTable: true},
				{View: "table2", Column: InternalIdColumn},
				{View: "table2", Column: "column1", Number: 1, IsFromTable: true},
				{View: "table2", Column: "column3", Number: 2, IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewInteger(1),
					value.NewInteger(1),
					value.NewString("str1"),
					value.NewNull(),
					value.NewNull(),
					value.NewNull(),
				}),
				NewRecord([]value.Primary{
					value.NewInteger(2),
					value.NewInteger(2),
					value.NewString("str2"),
					value.NewInteger(1),
					value.NewInteger(2),
					value.NewString("str22"),
				}),
				NewRecord([]value.Primary{
					value.NewInteger(3),
					value.NewInteger(3),
					value.NewString("str3"),
					value.NewInteger(2),
					value.NewInteger(3),
					value.NewString("str33"),
				}),
				NewRecord([]value.Primary{
					value.NewNull(),
					value.NewNull(),
					value.NewNull(),
					value.NewInteger(3),
					value.NewInteger(4),
					value.NewString("str44"),
				}),
			},
		},
	},
	{
		Name: "Left Outer Join Filter Error",
		View: &View{
			Header: NewHeaderWithId("table1", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewInteger(1),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewInteger(2),
					value.NewString("str2"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewInteger(3),
					value.NewString("str3"),
				}),
			},
		},
		JoinView: &View{
			Header: NewHeaderWithId("table2", []string{"column1", "column3"}),
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewInteger(2),
					value.NewString("str22"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewInteger(3),
					value.NewString("str33"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewInteger(4),
					value.NewString("str44"),
				}),
			},
		},
		Condition: parser.Comparison{
			LHS:      parser.FieldReference{View: parser.Identifier{Literal: "table1"}, Column: parser.Identifier{Literal: "notexist"}},
			RHS:      parser.FieldReference{View: parser.Identifier{Literal: "table2"}, Column: parser.Identifier{Literal: "column1"}},
			Operator: parser.Token{Token: '=', Literal: "="},
		},
		Direction: parser.LEFT,
		Error:     "field table1.notexist does not exist",
	},
	{
		Name: "Outer Join Direction Undefined",
		View: &View{
			Header: NewHeaderWithId("table1", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewInteger(1),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewInteger(2),
					value.NewString("str2"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewInteger(3),
					value.NewString("str3"),
				}),
			},
		},
		JoinView: &View{
			Header: NewHeaderWithId("table2", []string{"column1", "column3"}),
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewInteger(2),
					value.NewString("str22"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewInteger(3),
					value.NewString("str33"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewInteger(4),
					value.NewString("str44"),
				}),
			},
		},
		Condition: parser.Comparison{
			LHS:      parser.FieldReference{View: parser.Identifier{Literal: "table1"}, Column: parser.Identifier{Literal: "column1"}},
			RHS:      parser.FieldReference{View: parser.Identifier{Literal: "table2"}, Column: parser.Identifier{Literal: "column1"}},
			Operator: parser.Token{Token: '=', Literal: "="},
		},
		Direction: parser.TokenUndefined,
		Result: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", Number: 1, IsFromTable: true},
				{View: "table1", Column: "column2", Number: 2, IsFromTable: true},
				{View: "table2", Column: InternalIdColumn},
				{View: "table2", Column: "column1", Number: 1, IsFromTable: true},
				{View: "table2", Column: "column3", Number: 2, IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewInteger(1),
					value.NewInteger(1),
					value.NewString("str1"),
					value.NewNull(),
					value.NewNull(),
					value.NewNull(),
				}),
				NewRecord([]value.Primary{
					value.NewInteger(2),
					value.NewInteger(2),
					value.NewString("str2"),
					value.NewInteger(1),
					value.NewInteger(2),
					value.NewString("str22"),
				}),
				NewRecord([]value.Primary{
					value.NewInteger(3),
					value.NewInteger(3),
					value.NewString("str3"),
					value.NewInteger(2),
					value.NewInteger(3),
					value.NewString("str33"),
				}),
			},
		},
	},
}

func TestOuterJoin(t *testing.T) {
	for _, v := range outerJoinTests {
		if v.Scope == nil {
			v.Scope = NewReferenceScope(TestTx)
		}

		err := OuterJoin(context.Background(), v.Scope, v.View, v.JoinView, v.Condition, v.Direction)
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
			t.Log(v.View.RecordSet)
			t.Log(v.Result.RecordSet)
		}
	}
}

var calcMinimumRequiredTests = []struct {
	Int1    int
	Int2    int
	Default int
	Expect  int
}{
	{
		Int1:    13,
		Int2:    20,
		Default: 80,
		Expect:  5,
	},
	{
		Int1:    1,
		Int2:    200,
		Default: 80,
		Expect:  1,
	},
	{
		Int1:    199,
		Int2:    1,
		Default: 80,
		Expect:  100,
	},
	{
		Int1:    1,
		Int2:    0,
		Default: 80,
		Expect:  80,
	},
	{
		Int1:    1,
		Int2:    1,
		Default: 80,
		Expect:  80,
	},
}

func TestCalcMinimumRequired(t *testing.T) {
	for _, v := range calcMinimumRequiredTests {
		result := CalcMinimumRequired(v.Int1, v.Int2, v.Default)
		if result != v.Expect {
			t.Errorf("result = %d, want %d for %d, %d, %d", result, v.Expect, v.Int1, v.Int2, v.Default)
		}
	}
}

func GenerateBenchView(tableName string, records int, startIdx int) *View {
	view := &View{
		Header:    NewHeader(tableName, []string{"c1"}),
		RecordSet: make(RecordSet, records),
	}

	for i := 0; i < records; i++ {
		view.RecordSet[i] = NewRecord([]value.Primary{value.NewInteger(int64(i + startIdx))})
	}

	return view
}

func BenchmarkCrossJoin(b *testing.B) {
	ctx := context.Background()
	scope := NewReferenceScope(TestTx)

	for i := 0; i < b.N; i++ {
		view := GenerateBenchView("t1", 100, 0)
		joinView := GenerateBenchView("t2", 100, 50)

		_ = CrossJoin(ctx, scope, view, joinView)
	}
}

func BenchmarkInnerJoin(b *testing.B) {
	condition := parser.Comparison{
		LHS:      parser.FieldReference{View: parser.Identifier{Literal: "t1"}, Column: parser.Identifier{Literal: "c1"}},
		RHS:      parser.FieldReference{View: parser.Identifier{Literal: "t2"}, Column: parser.Identifier{Literal: "c1"}},
		Operator: parser.Token{Token: '=', Literal: "="},
	}

	ctx := context.Background()
	scope := NewReferenceScope(TestTx)

	for i := 0; i < b.N; i++ {
		view := GenerateBenchView("t1", 100, 0)
		joinView := GenerateBenchView("t2", 100, 50)

		_ = InnerJoin(ctx, scope, view, joinView, condition)
	}
}

func BenchmarkOuterJoin(b *testing.B) {
	condition := parser.Comparison{
		LHS:      parser.FieldReference{View: parser.Identifier{Literal: "t1"}, Column: parser.Identifier{Literal: "c1"}},
		RHS:      parser.FieldReference{View: parser.Identifier{Literal: "t2"}, Column: parser.Identifier{Literal: "c1"}},
		Operator: parser.Token{Token: '=', Literal: "="},
	}

	ctx := context.Background()
	scope := NewReferenceScope(TestTx)

	for i := 0; i < b.N; i++ {
		view := GenerateBenchView("t1", 100, 0)
		joinView := GenerateBenchView("t2", 100, 50)

		_ = OuterJoin(ctx, scope, view, joinView, condition, parser.LEFT)
	}
}
