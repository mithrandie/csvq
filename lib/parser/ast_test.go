package parser

import (
	"reflect"
	"testing"
	"time"

	"github.com/mithrandie/ternary"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/value"
)

func TestBaseExpr_Line(t *testing.T) {
	e := BaseExpr{
		line:       3,
		char:       5,
		sourceFile: "source.sql",
	}

	expect := 3
	if e.Line() != expect {
		t.Errorf("line = %d, want %d for %#v", e.Line(), expect, e)
	}
}

func TestBaseExpr_Char(t *testing.T) {
	e := BaseExpr{
		line:       3,
		char:       5,
		sourceFile: "source.sql",
	}

	expect := 5
	if e.Char() != expect {
		t.Errorf("line = %d, want %d for %#v", e.Char(), expect, e)
	}
}

func TestBaseExpr_SourceFile(t *testing.T) {
	e := BaseExpr{
		line:       3,
		char:       5,
		sourceFile: "source.sql",
	}

	expect := "source.sql"
	if e.SourceFile() != expect {
		t.Errorf("line = %q, want %q for %#v", e.SourceFile(), expect, e)
	}
}

func TestBaseExpr_HasParseInfo(t *testing.T) {
	expr := &BaseExpr{}
	if !expr.HasParseInfo() {
		t.Errorf("has parse info = %t, want %t for %#v", expr.HasParseInfo(), true, expr)
	}

	queryExpr := NewNullValue()
	if queryExpr.HasParseInfo() {
		t.Errorf("has parse info = %t, want %t for %#v", expr.HasParseInfo(), false, queryExpr)
	}
}

func TestPrimitiveType_String(t *testing.T) {
	e := NewTernaryValueFromString("true")
	expect := "TRUE"
	if e.String() != expect {
		t.Errorf("result = %q, want %q for %q ", e.String(), expect, e)
	}

	e = NewTernaryValue(ternary.FALSE)
	expect = "FALSE"
	if e.String() != expect {
		t.Errorf("result = %q, want %q for %q ", e.String(), expect, e)
	}

	e = NewStringValue("str")
	expect = "'str'"
	if e.String() != expect {
		t.Errorf("result = %q, want %q for %q ", e.String(), expect, e)
	}

	e = NewIntegerValue(1)
	expect = "1"
	if e.String() != expect {
		t.Errorf("result = %q, want %q for %q ", e.String(), expect, e)
	}

	e = NewFloatValue(1.234)
	expect = "1.234"
	if e.String() != expect {
		t.Errorf("result = %q, want %q for %q ", e.String(), expect, e)
	}

	e = NewNullValue()
	expect = "NULL"
	if e.String() != expect {
		t.Errorf("result = %q, want %q for %q ", e.String(), expect, e)
	}
}

func TestPrimitiveType_IsInteger(t *testing.T) {
	e := NewDatetimeValue(time.Date(2012, 2, 4, 9, 18, 15, 0, time.Local))
	if e.IsInteger() != false {
		t.Errorf("result = %t, want %t for %q ", e.IsInteger(), false, e)
	}

	e = NewIntegerValue(1)
	if e.IsInteger() != true {
		t.Errorf("result = %t, want %t for %q ", e.IsInteger(), true, e)
	}
}

func TestPlaceHolder_String(t *testing.T) {
	s := "?"
	ordinal := 3
	e := Placeholder{Literal: s, Ordinal: ordinal, Name: ""}
	expect := "?{3}"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}

	s = ":foo"
	ordinal = 5
	e = Placeholder{Literal: s, Ordinal: ordinal, Name: "foo"}
	expect = ":foo"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestIdentifier_String(t *testing.T) {
	s := "abcde"
	e := Identifier{Literal: s}
	if e.String() != s {
		t.Errorf("string = %q, want %q for %#v", e.String(), s, e)
	}

	s = "abcde"
	e = Identifier{Literal: s, Quoted: true}
	if e.String() != cmd.QuoteIdentifier(s) {
		t.Errorf("string = %q, want %q for %#v", e.String(), cmd.QuoteIdentifier(s), e)
	}
}

func TestFieldReference_String(t *testing.T) {
	e := FieldReference{
		View:   Identifier{Literal: "table1"},
		Column: Identifier{Literal: "column1"},
	}
	expect := "table1.column1"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}

	e = FieldReference{
		Column: Identifier{Literal: "column1"},
	}
	expect = "column1"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestColumnNumber_String(t *testing.T) {
	e := ColumnNumber{
		View:   Identifier{Literal: "table1"},
		Number: value.NewInteger(3),
	}
	expect := "table1.3"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestParentheses_String(t *testing.T) {
	s := "abcde"
	e := Parentheses{Expr: NewStringValue(s)}
	expect := "('abcde')"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestRowValue_String(t *testing.T) {
	e := RowValue{
		Value: ValueList{
			Values: []QueryExpression{
				NewIntegerValueFromString("1"),
				NewIntegerValueFromString("2"),
			},
		},
	}
	expect := "(1, 2)"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestValueList_String(t *testing.T) {
	e := ValueList{
		Values: []QueryExpression{
			NewIntegerValueFromString("1"),
			NewIntegerValueFromString("2"),
		},
	}
	expect := "(1, 2)"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestRowValueList_String(t *testing.T) {
	e := RowValueList{
		RowValues: []QueryExpression{
			ValueList{
				Values: []QueryExpression{
					NewIntegerValueFromString("1"),
					NewIntegerValueFromString("2"),
				},
			},
			ValueList{
				Values: []QueryExpression{
					NewIntegerValueFromString("3"),
					NewIntegerValueFromString("4"),
				},
			},
		},
	}
	expect := "((1, 2), (3, 4))"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestSelectQuery_IsForUpdate(t *testing.T) {
	e := SelectQuery{
		SelectEntity: SelectEntity{
			SelectClause: SelectClause{
				Fields: []QueryExpression{Field{Object: Identifier{Literal: "column"}}},
			},
		},
		Context: Token{Token: UPDATE, Literal: "update"},
	}
	if !e.IsForUpdate() {
		t.Errorf("IsForUpdate() = %t, want %t for %#v", e.IsForUpdate(), true, e)
	}

	e = SelectQuery{
		SelectEntity: SelectEntity{
			SelectClause: SelectClause{
				Fields: []QueryExpression{Field{Object: Identifier{Literal: "column"}}},
			},
		},
	}
	if e.IsForUpdate() {
		t.Errorf("IsForUpdate() = %t, want %t for %#v", e.IsForUpdate(), false, e)
	}
}

func TestSelectQuery_String(t *testing.T) {
	e := SelectQuery{
		WithClause: WithClause{
			InlineTables: []QueryExpression{
				InlineTable{
					Name: Identifier{Literal: "ct"},
					Query: SelectQuery{
						SelectEntity: SelectEntity{
							SelectClause: SelectClause{
								Fields: []QueryExpression{
									Field{Object: NewIntegerValueFromString("1")},
								},
							},
						},
					},
				},
			},
		},
		SelectEntity: SelectEntity{
			SelectClause: SelectClause{
				Fields: []QueryExpression{Field{Object: Identifier{Literal: "column"}}},
			},
			FromClause: FromClause{
				Tables: []QueryExpression{
					Table{Object: Identifier{Literal: "table"}},
				},
			},
		},
		OrderByClause: OrderByClause{
			Items: []QueryExpression{
				OrderItem{
					Value: Identifier{Literal: "column"},
				},
			},
		},
		LimitClause: LimitClause{
			Type:  Token{Token: LIMIT, Literal: "limit"},
			Value: NewIntegerValueFromString("10"),
			OffsetClause: OffsetClause{
				Value: NewIntegerValueFromString("10"),
			},
		},
		Context: Token{Token: UPDATE, Literal: "update"},
	}
	expect := "WITH ct AS (SELECT 1) SELECT column FROM table ORDER BY column LIMIT 10 OFFSET 10 FOR UPDATE"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestSelectSet_String(t *testing.T) {
	e := SelectSet{
		LHS: SelectEntity{
			SelectClause: SelectClause{
				Fields: []QueryExpression{Field{Object: NewIntegerValueFromString("1")}},
			},
		},
		Operator: Token{Token: UNION, Literal: "union"},
		All:      Token{Token: ALL, Literal: "all"},
		RHS: SelectEntity{
			SelectClause: SelectClause{
				Fields: []QueryExpression{Field{Object: NewIntegerValueFromString("2")}},
			},
		},
	}
	expect := "SELECT 1 UNION ALL SELECT 2"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestSelectEntity_String(t *testing.T) {
	e := SelectEntity{
		SelectClause: SelectClause{
			Fields: []QueryExpression{Field{Object: Identifier{Literal: "column"}}},
		},
		IntoClause: IntoClause{
			Variables: []Variable{
				{Name: "var1"},
				{Name: "var2"},
			},
		},
		FromClause: FromClause{
			Tables: []QueryExpression{
				Table{Object: Identifier{Literal: "table"}},
			},
		},
		WhereClause: WhereClause{
			Filter: Comparison{
				LHS:      Identifier{Literal: "column"},
				Operator: Token{Token: '>', Literal: ">"},
				RHS:      NewIntegerValueFromString("1"),
			},
		},
		GroupByClause: GroupByClause{
			Items: []QueryExpression{
				Identifier{Literal: "column1"},
			},
		},
		HavingClause: HavingClause{
			Filter: Comparison{
				LHS:      Identifier{Literal: "column"},
				Operator: Token{Token: '>', Literal: ">"},
				RHS:      NewIntegerValueFromString("1"),
			},
		},
	}

	expect := "SELECT column INTO @var1, @var2 FROM table WHERE column > 1 GROUP BY column1 HAVING column > 1"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestSelectClause_IsDistinct(t *testing.T) {
	e := SelectClause{}
	if e.IsDistinct() == true {
		t.Errorf("distinct = %t, want %t for %#v", e.IsDistinct(), false, e)
	}

	e = SelectClause{Distinct: Token{Token: DISTINCT, Literal: "distinct"}}
	if e.IsDistinct() == false {
		t.Errorf("distinct = %t, want %t for %#v", e.IsDistinct(), true, e)
	}
}

func TestSelectClause_String(t *testing.T) {
	e := SelectClause{
		Distinct: Token{Token: DISTINCT, Literal: "distinct"},
		Fields: []QueryExpression{
			Field{
				Object: Identifier{Literal: "column1"},
			},
			Field{
				Object: Identifier{Literal: "column2"},
				As:     Token{Token: AS, Literal: "as"},
				Alias:  Identifier{Literal: "alias"},
			},
		},
	}
	expect := "SELECT DISTINCT column1, column2 AS alias"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestIntoClause_String(t *testing.T) {
	e := IntoClause{
		Variables: []Variable{
			{Name: "var1"},
			{Name: "var2"},
		},
	}
	expect := "INTO @var1, @var2"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestFromClause_String(t *testing.T) {
	e := FromClause{
		Tables: []QueryExpression{
			Table{Object: Identifier{Literal: "table1"}},
			Table{Object: Identifier{Literal: "table2"}},
		},
	}
	expect := "FROM table1, table2"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestWhereClause_String(t *testing.T) {
	e := WhereClause{
		Filter: Comparison{
			LHS:      Identifier{Literal: "column"},
			Operator: Token{Token: '>', Literal: ">"},
			RHS:      NewIntegerValueFromString("1"),
		},
	}
	expect := "WHERE column > 1"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestGroupByClause_String(t *testing.T) {
	e := GroupByClause{
		Items: []QueryExpression{
			Identifier{Literal: "column1"},
			Identifier{Literal: "column2"},
		},
	}
	expect := "GROUP BY column1, column2"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestHavingClause_String(t *testing.T) {
	e := HavingClause{
		Filter: Comparison{
			LHS:      Identifier{Literal: "column"},
			Operator: Token{Token: '>', Literal: ">"},
			RHS:      NewIntegerValueFromString("1"),
		},
	}
	expect := "HAVING column > 1"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestOrderByClause_String(t *testing.T) {
	e := OrderByClause{
		Items: []QueryExpression{
			OrderItem{
				Value: Identifier{Literal: "column1"},
			},
			OrderItem{
				Value:     Identifier{Literal: "column2"},
				Direction: Token{Token: ASC, Literal: "asc"},
			},
		},
	}
	expect := "ORDER BY column1, column2 ASC"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestLimitClause_String(t *testing.T) {
	e := LimitClause{
		OffsetClause: OffsetClause{Value: NewIntegerValueFromString("10")},
	}
	expect := "OFFSET 10"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}

	e = LimitClause{
		Type:         Token{Token: LIMIT, Literal: "limit"},
		Value:        NewIntegerValueFromString("10"),
		Unit:         Token{Token: ROWS, Literal: "rows"},
		Restriction:  Token{Token: TIES, Literal: "with ties"},
		OffsetClause: OffsetClause{Value: NewIntegerValueFromString("10")},
	}
	expect = "LIMIT 10 ROWS WITH TIES OFFSET 10"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}

	e = LimitClause{
		Type:         Token{Token: FETCH, Literal: "fetch"},
		Position:     Token{Token: NEXT, Literal: "next"},
		Value:        NewIntegerValueFromString("10"),
		Unit:         Token{Token: ROWS, Literal: "rows"},
		Restriction:  Token{Token: TIES, Literal: "with ties"},
		OffsetClause: OffsetClause{Value: NewIntegerValueFromString("10")},
	}
	expect = "OFFSET 10 FETCH NEXT 10 ROWS WITH TIES"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestLimitClause_Percentage(t *testing.T) {
	e := LimitClause{Type: Token{Token: LIMIT, Literal: "limit"}, Value: NewIntegerValue(10)}
	if e.Percentage() {
		t.Errorf("percentage = %t, want %t for %#v", e.Percentage(), false, e)
	}

	e = LimitClause{Type: Token{Token: LIMIT, Literal: "limit"}, Value: NewIntegerValue(10), Unit: Token{Token: PERCENT, Literal: "percent"}}
	if !e.Percentage() {
		t.Errorf("percentage = %t, want %t for %#v", e.Percentage(), true, e)
	}
}

func TestLimitClause_WithTies(t *testing.T) {
	e := LimitClause{Type: Token{Token: LIMIT, Literal: "limit"}, Value: NewIntegerValue(10)}
	if e.WithTies() {
		t.Errorf("with ties = %t, want %t for %#v", e.WithTies(), false, e)
	}

	e = LimitClause{Type: Token{Token: LIMIT, Literal: "limit"}, Value: NewIntegerValue(10), Restriction: Token{Token: TIES, Literal: "with ties"}}
	if !e.WithTies() {
		t.Errorf("with ties = %t, want %t for %#v", e.WithTies(), true, e)
	}
}

func TestOffsetClause_String(t *testing.T) {
	e := OffsetClause{Value: NewIntegerValueFromString("10")}
	expect := "OFFSET 10"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}

	e = OffsetClause{Value: NewIntegerValueFromString("10"), Unit: Token{Token: ROWS, Literal: "rows"}}
	expect = "OFFSET 10 ROWS"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestWithClause_String(t *testing.T) {
	e := WithClause{
		InlineTables: []QueryExpression{
			InlineTable{
				Name: Identifier{Literal: "alias1"},
				Query: SelectQuery{
					SelectEntity: SelectEntity{
						SelectClause: SelectClause{
							Fields: []QueryExpression{
								NewIntegerValueFromString("1"),
							},
						},
					},
				},
			},
			InlineTable{
				Recursive: Token{Token: RECURSIVE, Literal: "recursive"},
				Name:      Identifier{Literal: "alias2"},
				Query: SelectQuery{
					SelectEntity: SelectEntity{
						SelectClause: SelectClause{
							Fields: []QueryExpression{
								NewIntegerValueFromString("2"),
							},
						},
					},
				},
			},
		},
	}
	expect := "WITH alias1 AS (SELECT 1), RECURSIVE alias2 AS (SELECT 2)"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestInlineTable_String(t *testing.T) {
	e := InlineTable{
		Recursive: Token{Token: RECURSIVE, Literal: "recursive"},
		Name:      Identifier{Literal: "it"},
		Fields: []QueryExpression{
			Identifier{Literal: "column1"},
		},
		Query: SelectQuery{
			SelectEntity: SelectEntity{
				SelectClause: SelectClause{
					Fields: []QueryExpression{
						NewIntegerValueFromString("1"),
					},
				},
			},
		},
	}
	expect := "RECURSIVE it (column1) AS (SELECT 1)"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestInlineTable_IsRecursive(t *testing.T) {
	e := InlineTable{
		Recursive: Token{Token: RECURSIVE, Literal: "recursive"},
		Name:      Identifier{Literal: "alias"},
		Query: SelectQuery{
			SelectEntity: SelectEntity{
				SelectClause: SelectClause{
					Fields: []QueryExpression{
						NewIntegerValueFromString("1"),
					},
				},
			},
		},
	}
	if e.IsRecursive() != true {
		t.Errorf("IsRecursive = %t, want %t for %#v", e.IsRecursive(), true, e)
	}

	e = InlineTable{
		Name: Identifier{Literal: "alias"},
		Query: SelectQuery{
			SelectEntity: SelectEntity{
				SelectClause: SelectClause{
					Fields: []QueryExpression{
						NewIntegerValueFromString("1"),
					},
				},
			},
		},
	}
	if e.IsRecursive() != false {
		t.Errorf("IsRecursive = %t, want %t for %#v", e.IsRecursive(), false, e)
	}
}

func TestSubquery_String(t *testing.T) {
	e := Subquery{
		Query: SelectQuery{
			SelectEntity: SelectEntity{
				SelectClause: SelectClause{
					Fields: []QueryExpression{
						NewIntegerValueFromString("1"),
					},
				},
				FromClause: FromClause{
					Tables: []QueryExpression{Dual{}},
				},
			},
		},
	}
	expect := "(SELECT 1 FROM DUAL)"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestTableObject_String(t *testing.T) {
	e := TableObject{
		Type:          Token{Token: FIXED, Literal: "fixed"},
		FormatElement: NewStringValue("[1, 2, 3]"),
		Path:          Identifier{Literal: "fixed_length.dat", Quoted: true},
		Args:          []QueryExpression{NewStringValue("utf8")},
	}
	expect := "FIXED('[1, 2, 3]', `fixed_length.dat`, 'utf8')"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}

	e = TableObject{
		Type:          Token{Token: FIXED, Literal: "fixed"},
		FormatElement: NewStringValue("[1, 2, 3]"),
		Path:          Identifier{Literal: "fixed_length.dat", Quoted: true},
		Args:          nil,
	}
	expect = "FIXED('[1, 2, 3]', `fixed_length.dat`)"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}

	e = TableObject{
		Type: Token{Token: LTSV, Literal: "ltsv"},
		Path: Identifier{Literal: "table.ltsv", Quoted: true},
		Args: []QueryExpression{NewStringValue("utf8")},
	}
	expect = "LTSV(`table.ltsv`, 'utf8')"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestJsonQuery_String(t *testing.T) {
	e := JsonQuery{
		JsonQuery: Token{Token: JSON_ROW, Literal: "json_array"},
		Query:     NewStringValue("key"),
		JsonText:  NewStringValue("{\"key\":1}"),
	}
	expect := "JSON_ROW('key', '{\"key\":1}')"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestComparison_String(t *testing.T) {
	e := Comparison{
		LHS:      Identifier{Literal: "column"},
		Operator: Token{Token: '>', Literal: ">"},
		RHS:      NewIntegerValueFromString("1"),
	}
	expect := "column > 1"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestIs_IsNegated(t *testing.T) {
	e := Is{}
	if e.IsNegated() == true {
		t.Errorf("negation = %t, want %t for %#v", e.IsNegated(), false, e)
	}

	e = Is{Negation: Token{Token: NOT, Literal: "not"}}
	if e.IsNegated() == false {
		t.Errorf("negation = %t, want %t for %#v", e.IsNegated(), true, e)
	}
}

func TestIs_String(t *testing.T) {
	e := Is{
		LHS:      Identifier{Literal: "column"},
		RHS:      NewNullValue(),
		Negation: Token{Token: NOT, Literal: "not"},
	}
	expect := "column IS NOT NULL"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestBetween_IsNegated(t *testing.T) {
	e := Between{}
	if e.IsNegated() == true {
		t.Errorf("negation = %t, want %t for %#v", e.IsNegated(), false, e)
	}

	e = Between{Negation: Token{Token: NOT, Literal: "not"}}
	if e.IsNegated() == false {
		t.Errorf("negation = %t, want %t for %#v", e.IsNegated(), true, e)
	}
}

func TestBetween_String(t *testing.T) {
	e := Between{
		LHS:      Identifier{Literal: "column"},
		Low:      NewIntegerValueFromString("-10"),
		High:     NewIntegerValueFromString("10"),
		Negation: Token{Token: NOT, Literal: "not"},
	}
	expect := "column NOT BETWEEN -10 AND 10"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestIn_IsNegated(t *testing.T) {
	e := In{}
	if e.IsNegated() == true {
		t.Errorf("negation = %t, want %t for %#v", e.IsNegated(), false, e)
	}

	e = In{Negation: Token{Token: NOT, Literal: "not"}}
	if e.IsNegated() == false {
		t.Errorf("negation = %t, want %t for %#v", e.IsNegated(), true, e)
	}
}

func TestIn_String(t *testing.T) {
	e := In{
		LHS: Identifier{Literal: "column"},
		Values: RowValue{
			Value: ValueList{
				Values: []QueryExpression{
					NewIntegerValueFromString("1"),
					NewIntegerValueFromString("2"),
				},
			},
		},
		Negation: Token{Token: NOT, Literal: "not"},
	}
	expect := "column NOT IN (1, 2)"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestAll_String(t *testing.T) {
	e := All{
		LHS: RowValue{
			Value: ValueList{
				Values: []QueryExpression{
					Identifier{Literal: "column1"},
					Identifier{Literal: "column2"},
				},
			},
		},
		Operator: Token{Token: '>', Literal: ">"},
		Values: Subquery{
			Query: SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						Fields: []QueryExpression{
							NewIntegerValueFromString("1"),
						},
					},
					FromClause: FromClause{
						Tables: []QueryExpression{Dual{}},
					},
				},
			},
		},
	}
	expect := "(column1, column2) > ALL (SELECT 1 FROM DUAL)"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestAny_String(t *testing.T) {
	e := Any{
		LHS: RowValue{
			Value: ValueList{
				Values: []QueryExpression{
					Identifier{Literal: "column1"},
					Identifier{Literal: "column2"},
				},
			},
		},
		Operator: Token{Token: '>', Literal: ">"},
		Values: Subquery{
			Query: SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						Fields: []QueryExpression{
							NewIntegerValueFromString("1"),
						},
					},
					FromClause: FromClause{
						Tables: []QueryExpression{Dual{}},
					},
				},
			},
		},
	}
	expect := "(column1, column2) > ANY (SELECT 1 FROM DUAL)"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestLike_IsNegated(t *testing.T) {
	e := Like{}
	if e.IsNegated() == true {
		t.Errorf("negation = %t, want %t for %#v", e.IsNegated(), false, e)
	}

	e = Like{Negation: Token{Token: NOT, Literal: "not"}}
	if e.IsNegated() == false {
		t.Errorf("negation = %t, want %t for %#v", e.IsNegated(), true, e)
	}
}

func TestLike_String(t *testing.T) {
	e := Like{
		LHS:      Identifier{Literal: "column"},
		Pattern:  NewStringValue("pattern"),
		Negation: Token{Token: NOT, Literal: "not"},
	}
	expect := "column NOT LIKE 'pattern'"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestExists_String(t *testing.T) {
	e := Exists{
		Query: Subquery{
			Query: SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						Fields: []QueryExpression{
							NewIntegerValueFromString("1"),
						},
					},
					FromClause: FromClause{
						Tables: []QueryExpression{Dual{}},
					},
				},
			},
		},
	}
	expect := "EXISTS (SELECT 1 FROM DUAL)"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestArithmetic_String(t *testing.T) {
	e := Arithmetic{
		LHS:      Identifier{Literal: "column"},
		Operator: Token{Token: '+', Literal: "+"},
		RHS:      NewIntegerValueFromString("2"),
	}
	expect := "column + 2"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestUnaryArithmetic_String(t *testing.T) {
	e := UnaryArithmetic{
		Operand:  Identifier{Literal: "column"},
		Operator: Token{Token: '-', Literal: "-"},
	}
	expect := "-column"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestLogic_String(t *testing.T) {
	e := Logic{
		LHS:      NewTernaryValueFromString("true"),
		Operator: Token{Token: AND, Literal: "and"},
		RHS:      NewTernaryValueFromString("false"),
	}
	expect := "TRUE AND FALSE"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestUnaryLogic_String(t *testing.T) {
	e := UnaryLogic{
		Operator: Token{Token: NOT, Literal: "not"},
		Operand:  NewTernaryValueFromString("false"),
	}
	expect := "NOT FALSE"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}

	e = UnaryLogic{
		Operator: Token{Token: '!', Literal: "!"},
		Operand:  NewTernaryValueFromString("false"),
	}
	expect = "!FALSE"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestConcat_String(t *testing.T) {
	e := Concat{
		Items: []QueryExpression{
			Identifier{Literal: "column"},
			NewStringValue("a"),
		},
	}
	expect := "column || 'a'"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestFunction_String(t *testing.T) {
	e := Function{
		Name: "sum",
		Args: []QueryExpression{
			Identifier{Literal: "column"},
		},
	}
	expect := "SUM(column)"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}

	e = Function{
		Name: "substring",
		Args: []QueryExpression{
			Identifier{Literal: "column"},
			NewIntegerValue(2),
			NewIntegerValue(5),
		},
	}
	expect = "SUBSTRING(column, 2, 5)"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}

	e = Function{
		Name: "substring",
		Args: []QueryExpression{
			Identifier{Literal: "column"},
			NewIntegerValue(2),
		},
		From: Token{Token: FROM, Literal: "from"},
	}
	expect = "SUBSTRING(column FROM 2)"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}

	e = Function{
		Name: "substring",
		Args: []QueryExpression{
			Identifier{Literal: "column"},
			NewIntegerValue(2),
			NewIntegerValue(5),
		},
		From: Token{Token: FROM, Literal: "from"},
		For:  Token{Token: FOR, Literal: "for"},
	}
	expect = "SUBSTRING(column FROM 2 FOR 5)"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestAggregateFunction_String(t *testing.T) {
	e := AggregateFunction{
		Name:     "sum",
		Distinct: Token{Token: DISTINCT, Literal: "distinct"},
		Args: []QueryExpression{
			FieldReference{Column: Identifier{Literal: "column"}},
		},
	}
	expect := "SUM(DISTINCT column)"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestAggregateFunction_IsDistinct(t *testing.T) {
	e := AggregateFunction{}
	if e.IsDistinct() == true {
		t.Errorf("distinct = %t, want %t for %#v", e.IsDistinct(), false, e)
	}

	e = AggregateFunction{Distinct: Token{Token: DISTINCT, Literal: "distinct"}}
	if e.IsDistinct() == false {
		t.Errorf("distinct = %t, want %t for %#v", e.IsDistinct(), true, e)
	}
}

func TestTable_String(t *testing.T) {
	e := Table{
		Object: Identifier{Literal: "table"},
		As:     Token{Token: AS, Literal: "as"},
		Alias:  Identifier{Literal: "alias"},
	}
	expect := "table AS alias"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}

	e = Table{
		Object: Identifier{Literal: "table"},
	}
	expect = "table"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}

	e = Table{
		Lateral: Token{Token: LATERAL, Literal: "lateral"},
		Object:  Identifier{Literal: "table"},
	}
	expect = "LATERAL table"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestTable_Name(t *testing.T) {
	e := Table{
		Object: Identifier{Literal: "table.csv"},
		As:     Token{Token: AS, Literal: "as"},
		Alias:  Identifier{Literal: "alias"},
	}
	expect := Identifier{Literal: "alias"}
	if !reflect.DeepEqual(e.Name(), expect) {
		t.Errorf("name = %q, want %q for %#v", e.Name(), expect, e)
	}

	e = Table{
		Object: Identifier{Literal: "/path/to/table.csv"},
	}
	expect = Identifier{Literal: "table"}
	if !reflect.DeepEqual(e.Name(), expect) {
		t.Errorf("name = %q, want %q for %#v", e.Name(), expect, e)
	}

	e = Table{
		Object: Stdin{},
	}
	expect = Identifier{Literal: "STDIN"}
	if !reflect.DeepEqual(e.Name(), expect) {
		t.Errorf("name = %q, want %q for %#v", e.Name(), expect, e)
	}

	e = Table{
		Object: Subquery{
			Query: SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						Fields: []QueryExpression{
							NewIntegerValueFromString("1"),
						},
					},
					FromClause: FromClause{
						Tables: []QueryExpression{Dual{}},
					},
				},
			},
		},
	}
	expect = Identifier{}
	if !reflect.DeepEqual(e.Name(), expect) {
		t.Errorf("name = %q, want %q for %#v", e.Name(), expect, e)
	}

	e = Table{
		Object: TableObject{
			Type:          Token{Token: FIXED, Literal: "fixed"},
			FormatElement: NewStringValue("[1, 2, 3]"),
			Path:          Identifier{Literal: "fixed_length.dat", Quoted: true},
			Args:          nil,
		},
	}
	expect = Identifier{Literal: "fixed_length"}
	if !reflect.DeepEqual(e.Name(), expect) {
		t.Errorf("name = %q, want %q for %#v", e.Name(), expect, e)
	}
}

func TestJoin_String(t *testing.T) {
	e := Join{
		Table:     Table{Object: Identifier{Literal: "table1"}},
		JoinTable: Table{Object: Identifier{Literal: "table2"}},
		Natural:   Token{Token: NATURAL, Literal: "natural"},
	}
	expect := "table1 NATURAL JOIN table2"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}

	e = Join{
		Table:     Table{Object: Identifier{Literal: "table1"}},
		JoinTable: Table{Object: Identifier{Literal: "table2"}},
		JoinType:  Token{Token: OUTER, Literal: "outer"},
		Direction: Token{Token: LEFT, Literal: "left"},
		Condition: JoinCondition{
			Using: []QueryExpression{
				Identifier{Literal: "column"},
			},
		},
	}
	expect = "table1 LEFT OUTER JOIN table2 USING (column)"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestJoinCondition_String(t *testing.T) {
	e := JoinCondition{
		On: Comparison{
			LHS:      Identifier{Literal: "column"},
			Operator: Token{Token: '>', Literal: ">"},
			RHS:      NewIntegerValueFromString("1"),
		},
	}
	expect := "ON column > 1"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}

	e = JoinCondition{
		Using: []QueryExpression{
			Identifier{Literal: "column1"},
			Identifier{Literal: "column2"},
		},
	}
	expect = "USING (column1, column2)"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestField_String(t *testing.T) {
	e := Field{
		Object: Identifier{Literal: "column"},
		As:     Token{Token: AS, Literal: "as"},
		Alias:  Identifier{Literal: "alias"},
	}
	expect := "column AS alias"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}

	e = Field{
		Object: Identifier{Literal: "column"},
	}
	expect = "column"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestField_Name(t *testing.T) {
	location, _ := time.LoadLocation("UTC")

	e := Field{
		Object: Identifier{Literal: "column"},
		As:     Token{Token: AS, Literal: "as"},
		Alias:  Identifier{Literal: "alias"},
	}
	expect := "alias"
	if e.Name() != expect {
		t.Errorf("name = %q, want %q for %#v", e.Name(), expect, e)
	}

	e = Field{
		Object: Identifier{Literal: "column"},
	}
	expect = "column"
	if e.Name() != expect {
		t.Errorf("name = %q, want %q for %#v", e.Name(), expect, e)
	}

	e = Field{
		Object: NewStringValue("foo"),
	}
	expect = "foo"
	if e.Name() != expect {
		t.Errorf("name = %q, want %q for %#v", e.Name(), expect, e)
	}

	e = Field{
		Object: NewDatetimeValueFromString("2012-01-01 00:00:00 +00:00", nil, location),
	}
	expect = "2012-01-01 00:00:00 +00:00"
	if e.Name() != expect {
		t.Errorf("name = %q, want %q for %#v", e.Name(), expect, e)
	}

	e = Field{
		Object: FieldReference{
			View:   Identifier{Literal: "tbl"},
			Column: Identifier{Literal: "column1"},
		},
	}
	expect = "column1"
	if e.Name() != expect {
		t.Errorf("name = %q, want %q for %#v", e.Name(), expect, e)
	}
}

func TestAllColumns_String(t *testing.T) {
	e := AllColumns{}
	if e.String() != "*" {
		t.Errorf("string = %q, want %q for %#v", e.String(), "*", e)
	}
}

func TestDual_String(t *testing.T) {
	s := "DUAL"
	e := Dual{}
	if e.String() != s {
		t.Errorf("string = %q, want %q for %#v", e.String(), s, e)
	}
}

func TestStdin_String(t *testing.T) {
	s := "STDIN"
	e := Stdin{}
	if e.String() != s {
		t.Errorf("string = %q, want %q for %#v", e.String(), s, e)
	}
}

func TestOrderItem_String(t *testing.T) {
	e := OrderItem{
		Value:     Identifier{Literal: "column"},
		Direction: Token{Token: DESC, Literal: "desc"},
	}
	expect := "column DESC"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}

	e = OrderItem{
		Value: Identifier{Literal: "column"},
	}
	expect = "column"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}

	e = OrderItem{
		Value:         Identifier{Literal: "column"},
		NullsPosition: Token{Token: FIRST, Literal: "first"},
	}
	expect = "column NULLS FIRST"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestCase_String(t *testing.T) {
	e := CaseExpr{
		Value: Identifier{Literal: "column"},
		When: []QueryExpression{
			CaseExprWhen{
				Condition: NewIntegerValueFromString("1"),
				Result:    NewStringValue("A"),
			},
			CaseExprWhen{
				Condition: NewIntegerValueFromString("2"),
				Result:    NewStringValue("B"),
			},
		},
		Else: CaseExprElse{Result: NewStringValue("C")},
	}
	expect := "CASE column WHEN 1 THEN 'A' WHEN 2 THEN 'B' ELSE 'C' END"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}

	e = CaseExpr{
		When: []QueryExpression{
			CaseExprWhen{
				Condition: Comparison{
					LHS:      Identifier{Literal: "column"},
					Operator: Token{Token: '>', Literal: ">"},
					RHS:      NewIntegerValueFromString("1"),
				},
				Result: NewStringValue("A"),
			},
		},
	}
	expect = "CASE WHEN column > 1 THEN 'A' END"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestCaseWhen_String(t *testing.T) {
	e := CaseExprWhen{
		Condition: Comparison{
			LHS:      Identifier{Literal: "column"},
			Operator: Token{Token: '>', Literal: ">"},
			RHS:      NewIntegerValueFromString("1"),
		},
		Result: NewStringValue("abcde"),
	}
	expect := "WHEN column > 1 THEN 'abcde'"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestCaseElse_String(t *testing.T) {
	e := CaseExprElse{Result: NewStringValue("abcde")}
	expect := "ELSE 'abcde'"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestListFunction_String(t *testing.T) {
	e := ListFunction{
		Name:     "listagg",
		Distinct: Token{Token: DISTINCT, Literal: "distinct"},
		Args: []QueryExpression{
			Identifier{Literal: "column1"},
			NewStringValue(","),
		},
		OrderBy: OrderByClause{
			Items: []QueryExpression{Identifier{Literal: "column1"}},
		},
	}
	expect := "LISTAGG(DISTINCT column1, ',') WITHIN GROUP (ORDER BY column1)"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}

	e = ListFunction{
		Name:     "listagg",
		Distinct: Token{Token: DISTINCT, Literal: "distinct"},
		Args: []QueryExpression{
			Identifier{Literal: "column1"},
			NewStringValue(","),
		},
	}
	expect = "LISTAGG(DISTINCT column1, ',')"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestListFunction_IsDistinct(t *testing.T) {
	e := ListFunction{}
	if e.IsDistinct() == true {
		t.Errorf("distinct = %t, want %t for %#v", e.IsDistinct(), false, e)
	}

	e = ListFunction{Distinct: Token{Token: DISTINCT, Literal: "distinct"}}
	if e.IsDistinct() == false {
		t.Errorf("distinct = %t, want %t for %#v", e.IsDistinct(), true, e)
	}
}

func TestAnalyticFunction_String(t *testing.T) {
	e := AnalyticFunction{
		Name:     "avg",
		Distinct: Token{Token: DISTINCT, Literal: "distinct"},
		Args: []QueryExpression{
			Identifier{Literal: "column4"},
		},
		IgnoreType: Token{Token: NULLS, Literal: "nulls"},
		AnalyticClause: AnalyticClause{
			PartitionClause: PartitionClause{
				Values: []QueryExpression{
					Identifier{Literal: "column1"},
					Identifier{Literal: "column2"},
				},
			},
			OrderByClause: OrderByClause{
				Items: []QueryExpression{
					OrderItem{Value: Identifier{Literal: "column3"}},
				},
			},
		},
	}
	expect := "AVG(DISTINCT column4 IGNORE NULLS) OVER (PARTITION BY column1, column2 ORDER BY column3)"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestAnalyticFunction_IsDistinct(t *testing.T) {
	e := AnalyticFunction{}
	if e.IsDistinct() == true {
		t.Errorf("distinct = %t, want %t for %#v", e.IsDistinct(), false, e)
	}

	e = AnalyticFunction{Distinct: Token{Token: DISTINCT, Literal: "distinct"}}
	if e.IsDistinct() == false {
		t.Errorf("distinct = %t, want %t for %#v", e.IsDistinct(), true, e)
	}
}

func TestAnalyticFunction_IgnoreNulls(t *testing.T) {
	e := AnalyticFunction{}
	if e.IgnoreNulls() == true {
		t.Errorf("IgnoreNulls() = %t, want %t for %#v", e.IgnoreNulls(), false, e)
	}

	e = AnalyticFunction{IgnoreType: Token{Token: NULLS, Literal: "nulls"}}
	if e.IgnoreNulls() == false {
		t.Errorf("IgnoreNulls() = %t, want %t for %#v", e.IgnoreNulls(), true, e)
	}
}

func TestAnalyticClause_String(t *testing.T) {
	e := AnalyticClause{
		PartitionClause: PartitionClause{
			Values: []QueryExpression{
				Identifier{Literal: "column1"},
				Identifier{Literal: "column2"},
			},
		},
		OrderByClause: OrderByClause{
			Items: []QueryExpression{
				OrderItem{Value: Identifier{Literal: "column3"}},
			},
		},
		WindowingClause: WindowingClause{
			FrameLow: WindowFramePosition{
				Direction: Token{Token: CURRENT, Literal: "current"},
			},
		},
	}
	expect := "PARTITION BY column1, column2 ORDER BY column3 ROWS CURRENT ROW"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestAnalyticClause_PartitionValues(t *testing.T) {
	e := AnalyticClause{
		PartitionClause: PartitionClause{
			Values: []QueryExpression{
				Identifier{Literal: "column1"},
				Identifier{Literal: "column2"},
			},
		},
	}
	expect := []QueryExpression{
		Identifier{Literal: "column1"},
		Identifier{Literal: "column2"},
	}
	if !reflect.DeepEqual(e.PartitionValues(), expect) {
		t.Errorf("partition values = %q, want %q for %#v", e.PartitionValues(), expect, e)
	}

	e = AnalyticClause{}
	expect = []QueryExpression(nil)
	if !reflect.DeepEqual(e.PartitionValues(), expect) {
		t.Errorf("partition values = %q, want %q for %#v", e.PartitionValues(), expect, e)
	}
}

func TestPartition_String(t *testing.T) {
	e := PartitionClause{
		Values: []QueryExpression{
			Identifier{Literal: "column1"},
			Identifier{Literal: "column2"},
		},
	}
	expect := "PARTITION BY column1, column2"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestWindowingClause_String(t *testing.T) {
	e := WindowingClause{
		FrameLow: WindowFramePosition{
			Direction: Token{Token: CURRENT, Literal: "current"},
		},
	}
	expect := "ROWS CURRENT ROW"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}

	e = WindowingClause{
		FrameLow: WindowFramePosition{
			Direction: Token{Token: PRECEDING, Literal: "preceding"},
			Offset:    1,
		},
		FrameHigh: WindowFramePosition{
			Direction: Token{Token: FOLLOWING, Literal: "following"},
			Unbounded: Token{Token: UNBOUNDED, Literal: "unbounded"},
		},
	}
	expect = "ROWS BETWEEN 1 PRECEDING AND UNBOUNDED FOLLOWING"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestVariable_String(t *testing.T) {
	e := Variable{
		Name: "var",
	}
	expect := "@var"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestVariableSubstitution_String(t *testing.T) {
	e := VariableSubstitution{
		Variable: Variable{
			Name: "var",
		},
		Value: NewIntegerValueFromString("1"),
	}
	expect := "@var := 1"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestEnvironmentVariable_String(t *testing.T) {
	e := EnvironmentVariable{
		Name: "envvar",
	}
	expect := "@%envvar"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}

	e = EnvironmentVariable{
		Name:   "envvar",
		Quoted: true,
	}
	expect = "@%`envvar`"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestRuntimeInformation_String(t *testing.T) {
	e := RuntimeInformation{
		Name: "ri",
	}
	expect := "@#RI"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestFlag_String(t *testing.T) {
	e := Flag{
		Name: "flag",
	}
	expect := "@@FLAG"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestCursorStatus_String(t *testing.T) {
	e := CursorStatus{
		Cursor:   Identifier{Literal: "cur"},
		Negation: Token{Token: NOT, Literal: "not"},
		Type:     Token{Token: RANGE, Literal: "range"},
	}
	expect := "CURSOR cur IS NOT IN RANGE"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}

	e = CursorStatus{
		Cursor: Identifier{Literal: "cur"},
		Type:   Token{Token: OPEN, Literal: "open"},
	}
	expect = "CURSOR cur IS OPEN"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestCursorAttrebute_String(t *testing.T) {
	e := CursorAttrebute{
		Cursor:    Identifier{Literal: "cur"},
		Attrebute: Token{Token: COUNT, Literal: "count"},
	}
	expect := "CURSOR cur COUNT"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}
