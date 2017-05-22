package query

import (
	"reflect"
	"testing"

	"github.com/mithrandie/csvq/lib/parser"
)

var parseJoinConditionTests = []struct {
	Name     string
	Join     parser.Join
	View     *View
	JoinView *View
	Result   parser.Expression
}{
	{
		Name: "No Condition",
		Join: parser.Join{
			Table:     parser.Table{Alias: parser.Identifier{Literal: "t1"}},
			JoinTable: parser.Table{Alias: parser.Identifier{Literal: "t2"}},
		},
		View:     &View{Header: NewHeader("table1", []string{"key1", "key2", "key3", "value1", "value2", "value3"})},
		JoinView: &View{Header: NewHeader("table2", []string{"key1", "key2", "key3", "value4"})},
		Result:   nil,
	},
	{
		Name: "Natural Join",
		Join: parser.Join{
			Table:     parser.Table{Alias: parser.Identifier{Literal: "t1"}},
			JoinTable: parser.Table{Alias: parser.Identifier{Literal: "t2"}},
			Natural:   parser.Token{Token: parser.NATURAL, Literal: "natural"},
		},
		View:     &View{Header: NewHeader("table1", []string{"key1", "key2", "key3", "value1", "value2", "value3"})},
		JoinView: &View{Header: NewHeader("table2", []string{"key1", "key2", "key3", "value4"})},
		Result: parser.Logic{
			LHS: parser.Logic{
				LHS: parser.Comparison{
					LHS:      parser.Identifier{Literal: "t1.key1"},
					RHS:      parser.Identifier{Literal: "t2.key1"},
					Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: "="},
				},
				RHS: parser.Comparison{
					LHS:      parser.Identifier{Literal: "t1.key2"},
					RHS:      parser.Identifier{Literal: "t2.key2"},
					Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: "="},
				},
				Operator: parser.Token{Token: parser.AND, Literal: "AND"},
			},
			RHS: parser.Comparison{
				LHS:      parser.Identifier{Literal: "t1.key3"},
				RHS:      parser.Identifier{Literal: "t2.key3"},
				Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: "="},
			},
			Operator: parser.Token{Token: parser.AND, Literal: "AND"},
		},
	},
	{
		Name: "Using Condition",
		Join: parser.Join{
			Table:     parser.Table{Alias: parser.Identifier{Literal: "t1"}},
			JoinTable: parser.Table{Alias: parser.Identifier{Literal: "t2"}},
			Condition: parser.JoinCondition{
				Using: []parser.Expression{
					parser.Identifier{Literal: "key1"},
				},
			},
		},
		View:     &View{Header: NewHeader("table1", []string{"key1", "key2", "key3", "value1", "value2", "value3"})},
		JoinView: &View{Header: NewHeader("table2", []string{"key1", "key2", "key3", "value4"})},
		Result: parser.Comparison{
			LHS:      parser.Identifier{Literal: "t1.key1"},
			RHS:      parser.Identifier{Literal: "t2.key1"},
			Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: "="},
		},
	},
	{
		Name: "On Condition",
		Join: parser.Join{
			Table:     parser.Table{Alias: parser.Identifier{Literal: "t1"}},
			JoinTable: parser.Table{Alias: parser.Identifier{Literal: "t2"}},
			Condition: parser.JoinCondition{
				On: parser.Comparison{
					LHS:      parser.Identifier{Literal: "t1.key1"},
					RHS:      parser.Identifier{Literal: "t2.key1"},
					Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: "="},
				},
			},
		},
		View:     &View{Header: NewHeader("table1", []string{"key1", "key2", "key3", "value1", "value2", "value3"})},
		JoinView: &View{Header: NewHeader("table2", []string{"key1", "key2", "key3", "value4"})},
		Result: parser.Comparison{
			LHS:      parser.Identifier{Literal: "t1.key1"},
			RHS:      parser.Identifier{Literal: "t2.key1"},
			Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: "="},
		},
	},
	{
		Name: "Natural Join Fields Does Not Duplicate",
		Join: parser.Join{
			Table:     parser.Table{Alias: parser.Identifier{Literal: "t1"}},
			JoinTable: parser.Table{Alias: parser.Identifier{Literal: "t2"}},
			Natural:   parser.Token{Token: parser.NATURAL, Literal: "natural"},
		},
		View:     &View{Header: NewHeader("table1", []string{"value1", "value2", "value3"})},
		JoinView: &View{Header: NewHeader("table2", []string{"value4"})},
		Result:   nil,
	},
}

func TestParseJoinCondition(t *testing.T) {
	for _, v := range parseJoinConditionTests {
		r := ParseJoinCondition(v.Join, v.View, v.JoinView)
		if !reflect.DeepEqual(r, v.Result) {
			t.Errorf("%s: condition = %q, want %q", v.Name, r, v.Result)
		}
	}
}

func TestCrossJoin(t *testing.T) {
	view := &View{
		Header: NewHeader("table1", []string{"column1", "column2"}),
		Records: []Record{
			NewRecord([]parser.Primary{
				parser.NewInteger(1),
				parser.NewString("str1"),
			}),
			NewRecord([]parser.Primary{
				parser.NewInteger(2),
				parser.NewString("str2"),
			}),
		},
	}
	joinView := &View{
		Header: NewHeader("table2", []string{"column3", "column4"}),
		Records: []Record{
			NewRecord([]parser.Primary{
				parser.NewInteger(3),
				parser.NewString("str3"),
			}),
			NewRecord([]parser.Primary{
				parser.NewInteger(4),
				parser.NewString("str4"),
			}),
		},
	}
	expect := &View{
		Header: []HeaderField{
			{Reference: "table1", Column: "column1", FromTable: true},
			{Reference: "table1", Column: "column2", FromTable: true},
			{Reference: "table2", Column: "column3", FromTable: true},
			{Reference: "table2", Column: "column4", FromTable: true},
		},
		Records: []Record{
			NewRecord([]parser.Primary{
				parser.NewInteger(1),
				parser.NewString("str1"),
				parser.NewInteger(3),
				parser.NewString("str3"),
			}),
			NewRecord([]parser.Primary{
				parser.NewInteger(1),
				parser.NewString("str1"),
				parser.NewInteger(4),
				parser.NewString("str4"),
			}),
			NewRecord([]parser.Primary{
				parser.NewInteger(2),
				parser.NewString("str2"),
				parser.NewInteger(3),
				parser.NewString("str3"),
			}),
			NewRecord([]parser.Primary{
				parser.NewInteger(2),
				parser.NewString("str2"),
				parser.NewInteger(4),
				parser.NewString("str4"),
			}),
		},
	}

	joined := CrossJoin(view, joinView)
	if !reflect.DeepEqual(joined, expect) {
		t.Errorf("cross join: result = %q, want %q", joined, expect)
	}
}

var innerJoinTests = []struct {
	Name      string
	View      *View
	JoinView  *View
	Condition parser.Expression
	Filter    Filter
	Result    *View
	Error     string
}{
	{
		Name: "Inner Join",
		View: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			Records: []Record{
				NewRecord([]parser.Primary{
					parser.NewInteger(1),
					parser.NewString("str1"),
				}),
				NewRecord([]parser.Primary{
					parser.NewInteger(2),
					parser.NewString("str2"),
				}),
				NewRecord([]parser.Primary{
					parser.NewInteger(3),
					parser.NewString("str3"),
				}),
			},
		},
		JoinView: &View{
			Header: NewHeader("table2", []string{"column1", "column3"}),
			Records: []Record{
				NewRecord([]parser.Primary{
					parser.NewInteger(1),
					parser.NewString("str1"),
				}),
				NewRecord([]parser.Primary{
					parser.NewInteger(2),
					parser.NewString("str22"),
				}),
			},
		},
		Condition: parser.Comparison{
			LHS:      parser.Identifier{Literal: "table1.column1"},
			RHS:      parser.Identifier{Literal: "table2.column1"},
			Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: "="},
		},
		Result: &View{
			Header: []HeaderField{
				{Reference: "table1", Column: "column1", FromTable: true},
				{Reference: "table1", Column: "column2", FromTable: true},
				{Reference: "table2", Column: "column1", FromTable: true},
				{Reference: "table2", Column: "column3", FromTable: true},
			},
			Records: []Record{
				NewRecord([]parser.Primary{
					parser.NewInteger(1),
					parser.NewString("str1"),
					parser.NewInteger(1),
					parser.NewString("str1"),
				}),
				NewRecord([]parser.Primary{
					parser.NewInteger(2),
					parser.NewString("str2"),
					parser.NewInteger(2),
					parser.NewString("str22"),
				}),
			},
		},
	},
	{
		Name: "Inner Join Filter Error",
		View: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			Records: []Record{
				NewRecord([]parser.Primary{
					parser.NewInteger(1),
					parser.NewString("str1"),
				}),
			},
		},
		JoinView: &View{
			Header: NewHeader("table2", []string{"column1", "column3"}),
			Records: []Record{
				NewRecord([]parser.Primary{
					parser.NewInteger(1),
					parser.NewString("str1"),
				}),
			},
		},
		Condition: parser.Comparison{
			LHS:      parser.Identifier{Literal: "table1.column1"},
			RHS:      parser.Identifier{Literal: "table2.notexist"},
			Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: "="},
		},
		Error: "identifier = table2.notexist: field does not exist",
	},
}

func TestInnerJoin(t *testing.T) {
	for _, v := range innerJoinTests {
		result, err := InnerJoin(v.View, v.JoinView, v.Condition, v.Filter)
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

var outerJoinTests = []struct {
	Name      string
	View      *View
	JoinView  *View
	Condition parser.Expression
	Direction int
	Filter    Filter
	Result    *View
	Error     string
}{
	{
		Name: "Left Outer Join",
		View: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			Records: []Record{
				NewRecord([]parser.Primary{
					parser.NewInteger(1),
					parser.NewString("str1"),
				}),
				NewRecord([]parser.Primary{
					parser.NewInteger(2),
					parser.NewString("str2"),
				}),
				NewRecord([]parser.Primary{
					parser.NewInteger(3),
					parser.NewString("str3"),
				}),
			},
		},
		JoinView: &View{
			Header: NewHeader("table2", []string{"column1", "column3"}),
			Records: []Record{
				NewRecord([]parser.Primary{
					parser.NewInteger(2),
					parser.NewString("str22"),
				}),
				NewRecord([]parser.Primary{
					parser.NewInteger(3),
					parser.NewString("str33"),
				}),
				NewRecord([]parser.Primary{
					parser.NewInteger(4),
					parser.NewString("str44"),
				}),
			},
		},
		Condition: parser.Comparison{
			LHS:      parser.Identifier{Literal: "table1.column1"},
			RHS:      parser.Identifier{Literal: "table2.column1"},
			Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: "="},
		},
		Direction: parser.LEFT,
		Result: &View{
			Header: []HeaderField{
				{Reference: "table1", Column: "column1", FromTable: true},
				{Reference: "table1", Column: "column2", FromTable: true},
				{Reference: "table2", Column: "column1", FromTable: true},
				{Reference: "table2", Column: "column3", FromTable: true},
			},
			Records: []Record{
				NewRecord([]parser.Primary{
					parser.NewInteger(1),
					parser.NewString("str1"),
					parser.NewNull(),
					parser.NewNull(),
				}),
				NewRecord([]parser.Primary{
					parser.NewInteger(2),
					parser.NewString("str2"),
					parser.NewInteger(2),
					parser.NewString("str22"),
				}),
				NewRecord([]parser.Primary{
					parser.NewInteger(3),
					parser.NewString("str3"),
					parser.NewInteger(3),
					parser.NewString("str33"),
				}),
			},
		},
	},
	{
		Name: "Right Outer Join",
		View: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			Records: []Record{
				NewRecord([]parser.Primary{
					parser.NewInteger(1),
					parser.NewString("str1"),
				}),
				NewRecord([]parser.Primary{
					parser.NewInteger(2),
					parser.NewString("str2"),
				}),
				NewRecord([]parser.Primary{
					parser.NewInteger(3),
					parser.NewString("str3"),
				}),
			},
		},
		JoinView: &View{
			Header: NewHeader("table2", []string{"column1", "column3"}),
			Records: []Record{
				NewRecord([]parser.Primary{
					parser.NewInteger(2),
					parser.NewString("str22"),
				}),
				NewRecord([]parser.Primary{
					parser.NewInteger(3),
					parser.NewString("str33"),
				}),
				NewRecord([]parser.Primary{
					parser.NewInteger(4),
					parser.NewString("str44"),
				}),
			},
		},
		Condition: parser.Comparison{
			LHS:      parser.Identifier{Literal: "table1.column1"},
			RHS:      parser.Identifier{Literal: "table2.column1"},
			Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: "="},
		},
		Direction: parser.RIGHT,
		Result: &View{
			Header: []HeaderField{
				{Reference: "table1", Column: "column1", FromTable: true},
				{Reference: "table1", Column: "column2", FromTable: true},
				{Reference: "table2", Column: "column1", FromTable: true},
				{Reference: "table2", Column: "column3", FromTable: true},
			},
			Records: []Record{
				NewRecord([]parser.Primary{
					parser.NewInteger(2),
					parser.NewString("str2"),
					parser.NewInteger(2),
					parser.NewString("str22"),
				}),
				NewRecord([]parser.Primary{
					parser.NewInteger(3),
					parser.NewString("str3"),
					parser.NewInteger(3),
					parser.NewString("str33"),
				}),
				NewRecord([]parser.Primary{
					parser.NewNull(),
					parser.NewNull(),
					parser.NewInteger(4),
					parser.NewString("str44"),
				}),
			},
		},
	},
	{
		Name: "Full Outer Join",
		View: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			Records: []Record{
				NewRecord([]parser.Primary{
					parser.NewInteger(1),
					parser.NewString("str1"),
				}),
				NewRecord([]parser.Primary{
					parser.NewInteger(2),
					parser.NewString("str2"),
				}),
				NewRecord([]parser.Primary{
					parser.NewInteger(3),
					parser.NewString("str3"),
				}),
			},
		},
		JoinView: &View{
			Header: NewHeader("table2", []string{"column1", "column3"}),
			Records: []Record{
				NewRecord([]parser.Primary{
					parser.NewInteger(2),
					parser.NewString("str22"),
				}),
				NewRecord([]parser.Primary{
					parser.NewInteger(3),
					parser.NewString("str33"),
				}),
				NewRecord([]parser.Primary{
					parser.NewInteger(4),
					parser.NewString("str44"),
				}),
			},
		},
		Condition: parser.Comparison{
			LHS:      parser.Identifier{Literal: "table1.column1"},
			RHS:      parser.Identifier{Literal: "table2.column1"},
			Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: "="},
		},
		Direction: parser.FULL,
		Result: &View{
			Header: []HeaderField{
				{Reference: "table1", Column: "column1", FromTable: true},
				{Reference: "table1", Column: "column2", FromTable: true},
				{Reference: "table2", Column: "column1", FromTable: true},
				{Reference: "table2", Column: "column3", FromTable: true},
			},
			Records: []Record{
				NewRecord([]parser.Primary{
					parser.NewInteger(1),
					parser.NewString("str1"),
					parser.NewNull(),
					parser.NewNull(),
				}),
				NewRecord([]parser.Primary{
					parser.NewInteger(2),
					parser.NewString("str2"),
					parser.NewInteger(2),
					parser.NewString("str22"),
				}),
				NewRecord([]parser.Primary{
					parser.NewInteger(3),
					parser.NewString("str3"),
					parser.NewInteger(3),
					parser.NewString("str33"),
				}),
				NewRecord([]parser.Primary{
					parser.NewNull(),
					parser.NewNull(),
					parser.NewInteger(4),
					parser.NewString("str44"),
				}),
			},
		},
	},
	{
		Name: "Left Outer Join Filter Error",
		View: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			Records: []Record{
				NewRecord([]parser.Primary{
					parser.NewInteger(1),
					parser.NewString("str1"),
				}),
				NewRecord([]parser.Primary{
					parser.NewInteger(2),
					parser.NewString("str2"),
				}),
				NewRecord([]parser.Primary{
					parser.NewInteger(3),
					parser.NewString("str3"),
				}),
			},
		},
		JoinView: &View{
			Header: NewHeader("table2", []string{"column1", "column3"}),
			Records: []Record{
				NewRecord([]parser.Primary{
					parser.NewInteger(2),
					parser.NewString("str22"),
				}),
				NewRecord([]parser.Primary{
					parser.NewInteger(3),
					parser.NewString("str33"),
				}),
				NewRecord([]parser.Primary{
					parser.NewInteger(4),
					parser.NewString("str44"),
				}),
			},
		},
		Condition: parser.Comparison{
			LHS:      parser.Identifier{Literal: "table1.notexist"},
			RHS:      parser.Identifier{Literal: "table2.column1"},
			Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: "="},
		},
		Direction: parser.LEFT,
		Error:     "identifier = table1.notexist: field does not exist",
	},
	{
		Name: "Outer Join Direction Undefined",
		View: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			Records: []Record{
				NewRecord([]parser.Primary{
					parser.NewInteger(1),
					parser.NewString("str1"),
				}),
				NewRecord([]parser.Primary{
					parser.NewInteger(2),
					parser.NewString("str2"),
				}),
				NewRecord([]parser.Primary{
					parser.NewInteger(3),
					parser.NewString("str3"),
				}),
			},
		},
		JoinView: &View{
			Header: NewHeader("table2", []string{"column1", "column3"}),
			Records: []Record{
				NewRecord([]parser.Primary{
					parser.NewInteger(2),
					parser.NewString("str22"),
				}),
				NewRecord([]parser.Primary{
					parser.NewInteger(3),
					parser.NewString("str33"),
				}),
				NewRecord([]parser.Primary{
					parser.NewInteger(4),
					parser.NewString("str44"),
				}),
			},
		},
		Condition: parser.Comparison{
			LHS:      parser.Identifier{Literal: "table1.column1"},
			RHS:      parser.Identifier{Literal: "table2.column1"},
			Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: "="},
		},
		Direction: parser.TOKEN_UNDEFINED,
		Result: &View{
			Header: []HeaderField{
				{Reference: "table1", Column: "column1", FromTable: true},
				{Reference: "table1", Column: "column2", FromTable: true},
				{Reference: "table2", Column: "column1", FromTable: true},
				{Reference: "table2", Column: "column3", FromTable: true},
			},
			Records: []Record{
				NewRecord([]parser.Primary{
					parser.NewInteger(1),
					parser.NewString("str1"),
					parser.NewNull(),
					parser.NewNull(),
				}),
				NewRecord([]parser.Primary{
					parser.NewInteger(2),
					parser.NewString("str2"),
					parser.NewInteger(2),
					parser.NewString("str22"),
				}),
				NewRecord([]parser.Primary{
					parser.NewInteger(3),
					parser.NewString("str3"),
					parser.NewInteger(3),
					parser.NewString("str33"),
				}),
			},
		},
	},
}

func TestOuterJoin(t *testing.T) {
	for _, v := range outerJoinTests {
		result, err := OuterJoin(v.View, v.JoinView, v.Condition, v.Direction, v.Filter)
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
