package parser

import (
	"reflect"
	"testing"
	"time"

	"github.com/mithrandie/csvq/lib/ternary"
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
		t.Errorf("line = %d, want %d for %#v", e.SourceFile(), expect, e)
	}
}

func TestIsNull(t *testing.T) {
	var p Primary

	p = NewInteger(1)
	if IsNull(p) {
		t.Errorf("value %#p is evaluated as is a null, but it is not so", p)
	}

	p = NewNull()
	if !IsNull(p) {
		t.Errorf("value %#p is evaluated as is not a null, but it is so", p)
	}
}

func TestBaseExpr_HasParseInfo(t *testing.T) {
	var expr *BaseExpr

	if expr.HasParseInfo() {
		t.Errorf("has parse info = %t, want %t for %#v", expr.HasParseInfo(), false, expr)
	}

	expr = &BaseExpr{}
	if !expr.HasParseInfo() {
		t.Errorf("has parse info = %t, want %t for %#v", expr.HasParseInfo(), true, expr)
	}
}

func TestString_String(t *testing.T) {
	s := "abcde"
	p := NewString(s)
	expect := "'" + s + "'"
	if p.String() != expect {
		t.Errorf("string = %q, want %q for %#v", p.String(), expect, p)
	}
}

func TestString_Value(t *testing.T) {
	s := "abcde"
	p := NewString(s)
	if p.Value() != s {
		t.Errorf("value = %q, want %q for %#v", p.Value(), s, p)
	}
}

func TestString_Ternary(t *testing.T) {
	s := "1"
	p := NewString(s)
	if p.Ternary() != ternary.TRUE {
		t.Errorf("ternary = %s, want %s for %#v", p.Ternary(), ternary.TRUE, p)
	}

	s = "0"
	p = NewString(s)
	if p.Ternary() != ternary.FALSE {
		t.Errorf("ternary = %s, want %s for %#v", p.Ternary(), ternary.FALSE, p)
	}
	s = "unknown"
	p = NewString(s)
	if p.Ternary() != ternary.UNKNOWN {
		t.Errorf("ternary = %s, want %s for %#v", p.Ternary(), ternary.UNKNOWN, p)
	}
}

func TestInteger_String(t *testing.T) {
	s := "1"
	p := NewInteger(1)
	if p.String() != s {
		t.Errorf("string = %q, want %q for %#v", p.String(), s, p)
	}
}

func TestInteger_Value(t *testing.T) {
	i := NewInteger(1)
	expect := int64(1)

	if i.Value() != expect {
		t.Errorf("value = %d, want %d for %#v", i.Value(), expect, i)
	}
}

func TestInteger_Ternary(t *testing.T) {
	p := NewInteger(1)
	if p.Ternary() != ternary.TRUE {
		t.Errorf("ternary = %s, want %s for %#v", p.Ternary(), ternary.TRUE, p)
	}
	p = NewInteger(0)
	if p.Ternary() != ternary.FALSE {
		t.Errorf("ternary = %s, want %s for %#v", p.Ternary(), ternary.FALSE, p)
	}
	p = NewInteger(2)
	if p.Ternary() != ternary.UNKNOWN {
		t.Errorf("ternary = %s, want %s for %#v", p.Ternary(), ternary.UNKNOWN, p)
	}
}

func TestFloat_String(t *testing.T) {
	s := "1.234"
	p := NewFloat(1.234)
	if p.String() != s {
		t.Errorf("string = %q, want %q for %#v", p.String(), s, p)
	}
}

func TestFloat_Value(t *testing.T) {
	f := NewFloat(1.234)
	expect := float64(1.234)

	if f.Value() != expect {
		t.Errorf("value = %f, want %f for %#v", f.Value(), expect, f)
	}
}

func TestFloat_Ternary(t *testing.T) {
	p := NewFloat(1)
	if p.Ternary() != ternary.TRUE {
		t.Errorf("ternary = %s, want %s for %#v", p.Ternary(), ternary.TRUE, p)
	}
	p = NewFloat(0)
	if p.Ternary() != ternary.FALSE {
		t.Errorf("ternary = %s, want %s for %#v", p.Ternary(), ternary.FALSE, p)
	}
	p = NewFloat(2)
	if p.Ternary() != ternary.UNKNOWN {
		t.Errorf("ternary = %s, want %s for %#v", p.Ternary(), ternary.UNKNOWN, p)
	}
}

func TestBoolean_String(t *testing.T) {
	s := "true"
	p := NewBoolean(true)
	if p.String() != s {
		t.Errorf("string = %q, want %q for %#v", p.String(), s, p)
	}
}

func TestBoolean_Value(t *testing.T) {
	p := NewBoolean(true)
	if p.Value() != true {
		t.Errorf("bool = %t, want %t for %#v", p.Value(), true, p)
	}
}

func TestBoolean_Ternary(t *testing.T) {
	p := NewBoolean(true)
	if p.Ternary() != ternary.TRUE {
		t.Errorf("ternary = %s, want %s for %#v", p.Ternary(), ternary.TRUE, p)
	}
}

func TestTernary_String(t *testing.T) {
	s := "TRUE"
	p := NewTernary(ternary.TRUE)
	if p.String() != s {
		t.Errorf("string = %q, want %q for %#v", p.String(), s, p)
	}

	s = "true"
	p = NewTernaryFromString("true")
	if p.String() != s {
		t.Errorf("string = %q, want %q for %#v", p.String(), s, p)
	}
}

func TestTernary_Ternary(t *testing.T) {
	p := NewTernary(ternary.TRUE)
	if p.Ternary() != ternary.TRUE {
		t.Errorf("ternary = %s, want %s for %#v", p.Ternary(), ternary.TRUE, p)
	}
}

func TestDatetime_String(t *testing.T) {
	s := "2012-01-01 12:34:56"
	p := NewDatetimeFromString(s)

	expect := "'" + s + "'"
	if p.String() != expect {
		t.Errorf("string = %q, want %q for %#v", p.String(), expect, p)
	}

	s = "2012-01-01T12:34:56-08:00"
	tm, _ := StrToTime(s)
	p = Datetime{
		value: tm,
	}
	expect = "'" + s + "'"
	if p.String() != expect {
		t.Errorf("string = %q, want %q for %#v", p.String(), expect, p)
	}
}

func TestDatetime_Value(t *testing.T) {
	d := NewDatetimeFromString("2012-01-01 12:34:56")
	expect := time.Date(2012, time.January, 1, 12, 34, 56, 0, time.Local)

	if d.Value() != expect {
		t.Errorf("value = %q, want %t for %#v", d.Value(), expect, d)
	}

	d = NewDatetimeFromString("2012-01-01T12:34:56-08:00")
	l, _ := time.LoadLocation("America/Los_Angeles")
	expect = time.Date(2012, time.January, 1, 12, 34, 56, 0, l)

	if d.Value().Sub(expect).Seconds() != 0 {
		t.Errorf("value = %q, want %t for %#v", d.Value(), expect, d)
	}
}

func TestDatetime_Ternary(t *testing.T) {
	p := NewDatetimeFromString("2012-01-01T12:34:56-08:00")
	if p.Ternary() != ternary.UNKNOWN {
		t.Errorf("ternary = %s, want %s for %#v", p.Ternary(), ternary.UNKNOWN, p)
	}
}

func TestDatetime_Format(t *testing.T) {
	dtstring := "2012-08-01T04:03:05.123-08:00"
	dt := NewDatetimeFromString(dtstring)
	expect := "2012-08-01T04:03:05-08:00"
	if dt.Format(time.RFC3339) != expect {
		t.Errorf("result = %q, want %q for %q ", dt.Format(time.RFC3339), expect, dtstring)
	}
}

func TestNull_String(t *testing.T) {
	s := "null"
	p := NewNullFromString(s)
	if p.String() != s {
		t.Errorf("string = %q, want %q for %#v", p.String(), s, p)
	}

	p = NewNull()
	if p.String() != "NULL" {
		t.Errorf("string = %q, want %q for %#v", p.String(), "NULL", p)
	}
}

func TestNull_Ternary(t *testing.T) {
	p := NewNull()
	if p.Ternary() != ternary.UNKNOWN {
		t.Errorf("ternary = %s, want %s for %#v", p.Ternary(), ternary.UNKNOWN, p)
	}
}

func TestPrimitiveType_String(t *testing.T) {
	e := NewTernaryValue(ternary.TRUE)
	expect := "TRUE"
	if e.String() != expect {
		t.Errorf("result = %q, want %q for %q ", e.String(), expect, e)
	}
}

func TestPrimitiveType_IsInteger(t *testing.T) {
	e := NewDatetimeValue(time.Date(2012, 2, 4, 9, 18, 15, 0, time.Local))
	if e.IsInteger() != false {
		t.Errorf("result = %t, want %t for %q ", e.IsInteger(), false, e)
	}

	e = NewPrimitiveType(NewInteger(1))
	if e.IsInteger() != true {
		t.Errorf("result = %t, want %t for %q ", e.IsInteger(), true, e)
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
	if e.String() != quoteIdentifier(s) {
		t.Errorf("string = %q, want %q for %#v", e.String(), quoteIdentifier(s), e)
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
		Number: NewInteger(3),
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
			Values: []Expression{
				NewIntegerValue(1),
				NewIntegerValue(2),
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
		Values: []Expression{
			NewIntegerValue(1),
			NewIntegerValue(2),
		},
	}
	expect := "(1, 2)"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestRowValueList_String(t *testing.T) {
	e := RowValueList{
		RowValues: []Expression{
			ValueList{
				Values: []Expression{
					NewIntegerValue(1),
					NewIntegerValue(2),
				},
			},
			ValueList{
				Values: []Expression{
					NewIntegerValue(3),
					NewIntegerValue(4),
				},
			},
		},
	}
	expect := "((1, 2), (3, 4))"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestSelectQuery_String(t *testing.T) {
	e := SelectQuery{
		WithClause: WithClause{
			With: "with",
			InlineTables: []Expression{
				InlineTable{
					Name: Identifier{Literal: "ct"},
					As:   "as",
					Query: SelectQuery{
						SelectEntity: SelectEntity{
							SelectClause: SelectClause{
								Select: "select",
								Fields: []Expression{
									Field{Object: NewIntegerValue(1)},
								},
							},
						},
					},
				},
			},
		},
		SelectEntity: SelectEntity{
			SelectClause: SelectClause{
				Select: "select",
				Fields: []Expression{Field{Object: Identifier{Literal: "column"}}},
			},
			FromClause: FromClause{
				From:   "from",
				Tables: []Expression{Table{Object: Identifier{Literal: "table"}}},
			},
		},
		OrderByClause: OrderByClause{
			OrderBy: "order by",
			Items: []Expression{
				OrderItem{
					Value: Identifier{Literal: "column"},
				},
			},
		},
		LimitClause: LimitClause{
			Limit: "limit",
			Value: NewIntegerValue(10),
		},
		OffsetClause: OffsetClause{
			Offset: "offset",
			Value:  NewIntegerValue(10),
		},
	}
	expect := "with ct as (select 1) select column from table order by column limit 10 offset 10"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestSelectSet_String(t *testing.T) {
	e := SelectSet{
		LHS: SelectEntity{
			SelectClause: SelectClause{
				Select: "select",
				Fields: []Expression{Field{Object: NewIntegerValue(1)}},
			},
		},
		Operator: Token{Token: UNION, Literal: "union"},
		All:      Token{Token: ALL, Literal: "all"},
		RHS: SelectEntity{
			SelectClause: SelectClause{
				Select: "select",
				Fields: []Expression{Field{Object: NewIntegerValue(2)}},
			},
		},
	}
	expect := "select 1 union all select 2"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestSelectEntity_String(t *testing.T) {
	e := SelectEntity{
		SelectClause: SelectClause{
			Select: "select",
			Fields: []Expression{Field{Object: Identifier{Literal: "column"}}},
		},
		FromClause: FromClause{
			From:   "from",
			Tables: []Expression{Table{Object: Identifier{Literal: "table"}}},
		},
		WhereClause: WhereClause{
			Where: "where",
			Filter: Comparison{
				LHS:      Identifier{Literal: "column"},
				Operator: ">",
				RHS:      NewIntegerValue(1),
			},
		},
		GroupByClause: GroupByClause{
			GroupBy: "group by",
			Items: []Expression{
				Identifier{Literal: "column1"},
			},
		},
		HavingClause: HavingClause{
			Having: "having",
			Filter: Comparison{
				LHS:      Identifier{Literal: "column"},
				Operator: ">",
				RHS:      NewIntegerValue(1),
			},
		},
	}

	expect := "select column from table where column > 1 group by column1 having column > 1"
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
		Select:   "select",
		Distinct: Token{Token: DISTINCT, Literal: "distinct"},
		Fields: []Expression{
			Field{
				Object: Identifier{Literal: "column1"},
			},
			Field{
				Object: Identifier{Literal: "column2"},
				As:     "as",
				Alias:  Identifier{Literal: "alias"},
			},
		},
	}
	expect := "select distinct column1, column2 as alias"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestFromClause_String(t *testing.T) {
	e := FromClause{
		From: "from",
		Tables: []Expression{
			Table{Object: Identifier{Literal: "table1"}},
			Table{Object: Identifier{Literal: "table2"}},
		},
	}
	expect := "from table1, table2"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestWhereClause_String(t *testing.T) {
	e := WhereClause{
		Where: "where",
		Filter: Comparison{
			LHS:      Identifier{Literal: "column"},
			Operator: ">",
			RHS:      NewIntegerValue(1),
		},
	}
	expect := "where column > 1"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestGroupByClause_String(t *testing.T) {
	e := GroupByClause{
		GroupBy: "group by",
		Items: []Expression{
			Identifier{Literal: "column1"},
			Identifier{Literal: "column2"},
		},
	}
	expect := "group by column1, column2"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestHavingClause_String(t *testing.T) {
	e := HavingClause{
		Having: "having",
		Filter: Comparison{
			LHS:      Identifier{Literal: "column"},
			Operator: ">",
			RHS:      NewIntegerValue(1),
		},
	}
	expect := "having column > 1"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestOrderByClause_String(t *testing.T) {
	e := OrderByClause{
		OrderBy: "order by",
		Items: []Expression{
			OrderItem{
				Value: Identifier{Literal: "column1"},
			},
			OrderItem{
				Value:     Identifier{Literal: "column2"},
				Direction: Token{Token: ASC, Literal: "asc"},
			},
		},
	}
	expect := "order by column1, column2 asc"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestLimitClause_String(t *testing.T) {
	e := LimitClause{Limit: "limit", Value: NewIntegerValue(10), With: LimitWith{With: "with", Type: Token{Token: TIES, Literal: "ties"}}}
	expect := "limit 10 with ties"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}

	e = LimitClause{Limit: "limit", Value: NewIntegerValue(10), Percent: "percent"}
	expect = "limit 10 percent"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestLimitClause_IsPercentage(t *testing.T) {
	e := LimitClause{Limit: "limit", Value: NewIntegerValue(10)}
	if e.IsPercentage() {
		t.Errorf("percentage = %t, want %t for %#v", e.IsPercentage(), false, e)
	}

	e = LimitClause{Limit: "limit", Value: NewIntegerValue(10), Percent: "percent"}
	if !e.IsPercentage() {
		t.Errorf("percentage = %t, want %t for %#v", e.IsPercentage(), true, e)
	}
}

func TestLimitClause_IsWithTies(t *testing.T) {
	e := LimitClause{Limit: "limit", Value: NewIntegerValue(10)}
	if e.IsWithTies() {
		t.Errorf("with ties = %t, want %t for %#v", e.IsWithTies(), false, e)
	}

	e = LimitClause{Limit: "limit", Value: NewIntegerValue(10), With: LimitWith{With: "with", Type: Token{Token: TIES, Literal: "ties"}}}
	if !e.IsWithTies() {
		t.Errorf("with ties = %t, want %t for %#v", e.IsWithTies(), true, e)
	}
}

func TestLimitWith_String(t *testing.T) {
	e := LimitWith{With: "with", Type: Token{Token: TIES, Literal: "ties"}}
	expect := "with ties"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestOffsetClause_String(t *testing.T) {
	e := OffsetClause{Offset: "offset", Value: NewIntegerValue(10)}
	expect := "offset 10"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestWithClause_String(t *testing.T) {
	e := WithClause{
		With: "with",
		InlineTables: []Expression{
			InlineTable{
				Name: Identifier{Literal: "alias1"},
				As:   "as",
				Query: SelectQuery{
					SelectEntity: SelectEntity{
						SelectClause: SelectClause{
							Select: "select",
							Fields: []Expression{
								NewIntegerValue(1),
							},
						},
					},
				},
			},
			InlineTable{
				Recursive: Token{Token: RECURSIVE, Literal: "recursive"},
				Name:      Identifier{Literal: "alias2"},
				As:        "as",
				Query: SelectQuery{
					SelectEntity: SelectEntity{
						SelectClause: SelectClause{
							Select: "select",
							Fields: []Expression{
								NewIntegerValue(2),
							},
						},
					},
				},
			},
		},
	}
	expect := "with alias1 as (select 1), recursive alias2 as (select 2)"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestInlineTable_String(t *testing.T) {
	e := InlineTable{
		Recursive: Token{Token: RECURSIVE, Literal: "recursive"},
		Name:      Identifier{Literal: "alias"},
		Fields: []Expression{
			Identifier{Literal: "column1"},
		},
		As: "as",
		Query: SelectQuery{
			SelectEntity: SelectEntity{
				SelectClause: SelectClause{
					Select: "select",
					Fields: []Expression{
						NewIntegerValue(1),
					},
				},
			},
		},
	}
	expect := "recursive alias (column1) as (select 1)"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestInlineTable_IsRecursive(t *testing.T) {
	e := InlineTable{
		Recursive: Token{Token: RECURSIVE, Literal: "recursive"},
		Name:      Identifier{Literal: "alias"},
		As:        "as",
		Query: SelectQuery{
			SelectEntity: SelectEntity{
				SelectClause: SelectClause{
					Select: "select",
					Fields: []Expression{
						NewIntegerValue(1),
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
		As:   "as",
		Query: SelectQuery{
			SelectEntity: SelectEntity{
				SelectClause: SelectClause{
					Select: "select",
					Fields: []Expression{
						NewIntegerValue(1),
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
					Select: "select",
					Fields: []Expression{
						NewIntegerValue(1),
					},
				},
				FromClause: FromClause{
					From: "from",
					Tables: []Expression{
						Dual{Dual: "dual"},
					},
				},
			},
		},
	}
	expect := "(select 1 from dual)"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestComparison_String(t *testing.T) {
	e := Comparison{
		LHS:      Identifier{Literal: "column"},
		Operator: ">",
		RHS:      NewIntegerValue(1),
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
		Is:       "is",
		LHS:      Identifier{Literal: "column"},
		RHS:      NewNullValueFromString("null"),
		Negation: Token{Token: NOT, Literal: "not"},
	}
	expect := "column is not null"
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
		Between:  "between",
		And:      "and",
		LHS:      Identifier{Literal: "column"},
		Low:      NewIntegerValue(-10),
		High:     NewIntegerValue(10),
		Negation: Token{Token: NOT, Literal: "not"},
	}
	expect := "column not between -10 and 10"
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
		In:  "in",
		LHS: Identifier{Literal: "column"},
		Values: RowValue{
			Value: ValueList{
				Values: []Expression{
					NewIntegerValue(1),
					NewIntegerValue(2),
				},
			},
		},
		Negation: Token{Token: NOT, Literal: "not"},
	}
	expect := "column not in (1, 2)"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestAll_String(t *testing.T) {
	e := All{
		All: "all",
		LHS: RowValue{
			Value: ValueList{
				Values: []Expression{
					Identifier{Literal: "column1"},
					Identifier{Literal: "column2"},
				},
			},
		},
		Operator: ">",
		Values: Subquery{
			Query: SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						Select: "select",
						Fields: []Expression{
							NewIntegerValue(1),
						},
					},
					FromClause: FromClause{
						From: "from",
						Tables: []Expression{
							Dual{Dual: "dual"},
						},
					},
				},
			},
		},
	}
	expect := "(column1, column2) > all (select 1 from dual)"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestAny_String(t *testing.T) {
	e := Any{
		Any: "any",
		LHS: RowValue{
			Value: ValueList{
				Values: []Expression{
					Identifier{Literal: "column1"},
					Identifier{Literal: "column2"},
				},
			},
		},
		Operator: ">",
		Values: Subquery{
			Query: SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						Select: "select",
						Fields: []Expression{
							NewIntegerValue(1),
						},
					},
					FromClause: FromClause{
						From: "from",
						Tables: []Expression{
							Dual{Dual: "dual"},
						},
					},
				},
			},
		},
	}
	expect := "(column1, column2) > any (select 1 from dual)"
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
		Like:     "like",
		LHS:      Identifier{Literal: "column"},
		Pattern:  NewStringValue("pattern"),
		Negation: Token{Token: NOT, Literal: "not"},
	}
	expect := "column not like 'pattern'"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestExists_String(t *testing.T) {
	e := Exists{
		Exists: "exists",
		Query: Subquery{
			Query: SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						Select: "select",
						Fields: []Expression{
							NewIntegerValue(1),
						},
					},
					FromClause: FromClause{
						From: "from",
						Tables: []Expression{
							Dual{Dual: "dual"},
						},
					},
				},
			},
		},
	}
	expect := "exists (select 1 from dual)"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestArithmetic_String(t *testing.T) {
	e := Arithmetic{
		LHS:      Identifier{Literal: "column"},
		Operator: int('+'),
		RHS:      NewIntegerValue(2),
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
	expect := "true and false"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestUnaryLogic_String(t *testing.T) {
	e := UnaryLogic{
		Operator: Token{Token: NOT, Literal: "not"},
		Operand:  NewTernaryValueFromString("false"),
	}
	expect := "not false"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}

	e = UnaryLogic{
		Operator: Token{Token: '!', Literal: "!"},
		Operand:  NewTernaryValueFromString("false"),
	}
	expect = "!false"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestConcat_String(t *testing.T) {
	e := Concat{
		Items: []Expression{
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
		Args: []Expression{
			Identifier{Literal: "column"},
		},
	}
	expect := "sum(column)"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestAggregateFunction_String(t *testing.T) {
	e := AggregateFunction{
		Name:     "sum",
		Distinct: Token{Token: DISTINCT, Literal: "distinct"},
		Args: []Expression{
			FieldReference{Column: Identifier{Literal: "column"}},
		},
	}
	expect := "sum(distinct column)"
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
		As:     "as",
		Alias:  Identifier{Literal: "alias"},
	}
	expect := "table as alias"
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
}

func TestTable_Name(t *testing.T) {
	e := Table{
		Object: Identifier{Literal: "table.csv"},
		As:     "as",
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
		Object: Subquery{
			Query: SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						Select: "select",
						Fields: []Expression{
							NewIntegerValue(1),
						},
					},
					FromClause: FromClause{
						From: "from",
						Tables: []Expression{
							Dual{Dual: "dual"},
						},
					},
				},
			},
		},
	}
	expect = Identifier{Literal: "(select 1 from dual)"}
	if !reflect.DeepEqual(e.Name(), expect) {
		t.Errorf("name = %q, want %q for %#v", e.Name(), expect, e)
	}
}

func TestJoin_String(t *testing.T) {
	e := Join{
		Join:      "join",
		Table:     Table{Object: Identifier{Literal: "table1"}},
		JoinTable: Table{Object: Identifier{Literal: "table2"}},
		Natural:   Token{Token: NATURAL, Literal: "natural"},
	}
	expect := "table1 natural join table2"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}

	e = Join{
		Join:      "join",
		Table:     Table{Object: Identifier{Literal: "table1"}},
		JoinTable: Table{Object: Identifier{Literal: "table2"}},
		JoinType:  Token{Token: OUTER, Literal: "outer"},
		Direction: Token{Token: LEFT, Literal: "left"},
		Condition: JoinCondition{
			Literal: "using",
			Using: []Expression{
				Identifier{Literal: "column"},
			},
		},
	}
	expect = "table1 left outer join table2 using (column)"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestJoinCondition_String(t *testing.T) {
	e := JoinCondition{
		Literal: "on",
		On: Comparison{
			LHS:      Identifier{Literal: "column"},
			Operator: ">",
			RHS:      NewIntegerValue(1),
		},
	}
	expect := "on column > 1"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}

	e = JoinCondition{
		Literal: "using",
		Using: []Expression{
			Identifier{Literal: "column1"},
			Identifier{Literal: "column2"},
		},
	}
	expect = "using (column1, column2)"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestField_String(t *testing.T) {
	e := Field{
		Object: Identifier{Literal: "column"},
		As:     "as",
		Alias:  Identifier{Literal: "alias"},
	}
	expect := "column as alias"
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
	e := Field{
		Object: Identifier{Literal: "column"},
		As:     "as",
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
		Object: NewDatetimeValueFromString("2012-01-01"),
	}
	expect = "2012-01-01"
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
	s := "dual"
	e := Dual{Dual: s}
	if e.String() != s {
		t.Errorf("string = %q, want %q for %#v", e.String(), s, e)
	}
}

func TestStdin_String(t *testing.T) {
	s := "stdin"
	e := Stdin{Stdin: s}
	if e.String() != s {
		t.Errorf("string = %q, want %q for %#v", e.String(), s, e)
	}
}

func TestOrderItem_String(t *testing.T) {
	e := OrderItem{
		Value:     Identifier{Literal: "column"},
		Direction: Token{Token: DESC, Literal: "desc"},
	}
	expect := "column desc"
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
		Value:    Identifier{Literal: "column"},
		Nulls:    "nulls",
		Position: Token{Token: FIRST, Literal: "first"},
	}
	expect = "column nulls first"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestCase_String(t *testing.T) {
	e := Case{
		Case:  "case",
		End:   "end",
		Value: Identifier{Literal: "column"},
		When: []Expression{
			CaseWhen{
				When:      "when",
				Then:      "then",
				Condition: NewIntegerValue(1),
				Result:    NewStringValue("A"),
			},
			CaseWhen{
				When:      "when",
				Then:      "then",
				Condition: NewIntegerValue(2),
				Result:    NewStringValue("B"),
			},
		},
		Else: CaseElse{Else: "else", Result: NewStringValue("C")},
	}
	expect := "case column when 1 then 'A' when 2 then 'B' else 'C' end"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}

	e = Case{
		Case: "case",
		End:  "end",
		When: []Expression{
			CaseWhen{
				When: "when",
				Then: "then",
				Condition: Comparison{
					LHS:      Identifier{Literal: "column"},
					Operator: ">",
					RHS:      NewIntegerValue(1),
				},
				Result: NewStringValue("A"),
			},
		},
	}
	expect = "case when column > 1 then 'A' end"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestCaseWhen_String(t *testing.T) {
	e := CaseWhen{
		When: "when",
		Then: "then",
		Condition: Comparison{
			LHS:      Identifier{Literal: "column"},
			Operator: ">",
			RHS:      NewIntegerValue(1),
		},
		Result: NewStringValue("abcde"),
	}
	expect := "when column > 1 then 'abcde'"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestCaseElse_String(t *testing.T) {
	e := CaseElse{Else: "else", Result: NewStringValue("abcde")}
	expect := "else 'abcde'"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestListAgg_String(t *testing.T) {
	e := ListAgg{
		ListAgg:  "listagg",
		Distinct: Token{Token: DISTINCT, Literal: "distinct"},
		Args: []Expression{
			Identifier{Literal: "column1"},
			NewStringValue(","),
		},
		WithinGroup: "within group",
		OrderBy: OrderByClause{
			OrderBy: "order by",
			Items:   []Expression{Identifier{Literal: "column1"}},
		},
	}
	expect := "listagg(distinct column1, ',') within group (order by column1)"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}

	e = ListAgg{
		ListAgg:  "listagg",
		Distinct: Token{Token: DISTINCT, Literal: "distinct"},
		Args: []Expression{
			Identifier{Literal: "column1"},
			NewStringValue(","),
		},
		WithinGroup: "within group",
	}
	expect = "listagg(distinct column1, ',') within group ()"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestListAgg_IsDistinct(t *testing.T) {
	e := ListAgg{}
	if e.IsDistinct() == true {
		t.Errorf("distinct = %t, want %t for %#v", e.IsDistinct(), false, e)
	}

	e = ListAgg{Distinct: Token{Token: DISTINCT, Literal: "distinct"}}
	if e.IsDistinct() == false {
		t.Errorf("distinct = %t, want %t for %#v", e.IsDistinct(), true, e)
	}
}

func TestAnalyticFunction_String(t *testing.T) {
	e := AnalyticFunction{
		Name:     "avg",
		Distinct: Token{Token: DISTINCT, Literal: "distinct"},
		Args: []Expression{
			Identifier{Literal: "column4"},
		},
		IgnoreNulls:    true,
		IgnoreNullsLit: "ignore nulls",
		Over:           "over",
		AnalyticClause: AnalyticClause{
			Partition: Partition{
				PartitionBy: "partition by",
				Values: []Expression{
					Identifier{Literal: "column1"},
					Identifier{Literal: "column2"},
				},
			},
			OrderByClause: OrderByClause{
				OrderBy: "order by",
				Items: []Expression{
					OrderItem{Value: Identifier{Literal: "column3"}},
				},
			},
		},
	}
	expect := "avg(distinct column4 ignore nulls) over (partition by column1, column2 order by column3)"
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

func TestAnalyticClause_String(t *testing.T) {
	e := AnalyticClause{
		Partition: Partition{
			PartitionBy: "partition by",
			Values: []Expression{
				Identifier{Literal: "column1"},
				Identifier{Literal: "column2"},
			},
		},
		OrderByClause: OrderByClause{
			OrderBy: "order by",
			Items: []Expression{
				OrderItem{Value: Identifier{Literal: "column3"}},
			},
		},
	}
	expect := "partition by column1, column2 order by column3"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestAnalyticClause_PartitionValues(t *testing.T) {
	e := AnalyticClause{
		Partition: Partition{
			PartitionBy: "partition by",
			Values: []Expression{
				Identifier{Literal: "column1"},
				Identifier{Literal: "column2"},
			},
		},
	}
	expect := []Expression{
		Identifier{Literal: "column1"},
		Identifier{Literal: "column2"},
	}
	if !reflect.DeepEqual(e.PartitionValues(), expect) {
		t.Errorf("partition values = %q, want %q for %#v", e.PartitionValues(), expect, e)
	}

	e = AnalyticClause{}
	expect = []Expression(nil)
	if !reflect.DeepEqual(e.PartitionValues(), expect) {
		t.Errorf("partition values = %q, want %q for %#v", e.PartitionValues(), expect, e)
	}
}

func TestAnalyticClause_OrderValues(t *testing.T) {
	e := AnalyticClause{
		OrderByClause: OrderByClause{
			OrderBy: "order by",
			Items: []Expression{
				OrderItem{Value: Identifier{Literal: "column3"}},
			},
		},
	}
	expect := []Expression{
		Identifier{Literal: "column3"},
	}
	if !reflect.DeepEqual(e.OrderValues(), expect) {
		t.Errorf("order values = %q, want %q for %#v", e.OrderValues(), expect, e)
	}

	e = AnalyticClause{}
	expect = []Expression(nil)
	if !reflect.DeepEqual(e.OrderValues(), expect) {
		t.Errorf("order values = %q, want %q for %#v", e.OrderValues(), expect, e)
	}
}

func TestPartition_String(t *testing.T) {
	e := Partition{
		PartitionBy: "partition by",
		Values: []Expression{
			Identifier{Literal: "column1"},
			Identifier{Literal: "column2"},
		},
	}
	expect := "partition by column1, column2"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestVariable_String(t *testing.T) {
	e := Variable{
		Name: "@var",
	}
	expect := "@var"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestVariableSubstitution_String(t *testing.T) {
	e := VariableSubstitution{
		Variable: Variable{
			Name: "@var",
		},
		Value: NewIntegerValue(1),
	}
	expect := "@var := 1"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestVariableAssignment_String(t *testing.T) {
	e := VariableAssignment{
		Variable: Variable{Name: "@var"},
		Value:    NewIntegerValue(1),
	}
	expect := "@var := 1"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}

	e = VariableAssignment{
		Variable: Variable{Name: "@var"},
		Value:    nil,
	}
	expect = "@var"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestInsertQuery_String(t *testing.T) {
	e := InsertQuery{
		WithClause: WithClause{
			With: "with",
			InlineTables: []Expression{
				InlineTable{
					Name: Identifier{Literal: "ct"},
					As:   "as",
					Query: SelectQuery{
						SelectEntity: SelectEntity{
							SelectClause: SelectClause{
								Select: "select",
								Fields: []Expression{
									Field{Object: NewIntegerValue(1)},
								},
							},
						},
					},
				},
			},
		},
		Insert: "insert",
		Into:   "into",
		Table:  Identifier{Literal: "table1"},
		Fields: []Expression{
			Identifier{Literal: "column1"},
			Identifier{Literal: "column2"},
		},
		Values: "values",
		ValuesList: []Expression{
			RowValue{
				Value: ValueList{
					Values: []Expression{
						NewIntegerValue(1),
						NewIntegerValue(2),
					},
				},
			},
			RowValue{
				Value: ValueList{
					Values: []Expression{
						NewIntegerValue(3),
						NewIntegerValue(4),
					},
				},
			},
		},
	}
	expect := "with ct as (select 1) insert into table1 (column1, column2) values (1, 2), (3, 4)"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}

	e = InsertQuery{
		Insert: "insert",
		Into:   "into",
		Table:  Identifier{Literal: "table1"},
		Fields: []Expression{
			Identifier{Literal: "column1"},
			Identifier{Literal: "column2"},
		},
		Query: SelectQuery{
			SelectEntity: SelectEntity{
				SelectClause: SelectClause{
					Select: "select",
					Fields: []Expression{
						NewIntegerValue(1),
						NewIntegerValue(2),
					},
				},
			},
		},
	}
	expect = "insert into table1 (column1, column2) select 1, 2"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestUpdateQuery_String(t *testing.T) {
	e := UpdateQuery{
		WithClause: WithClause{
			With: "with",
			InlineTables: []Expression{
				InlineTable{
					Name: Identifier{Literal: "ct"},
					As:   "as",
					Query: SelectQuery{
						SelectEntity: SelectEntity{
							SelectClause: SelectClause{
								Select: "select",
								Fields: []Expression{
									Field{Object: NewIntegerValue(1)},
								},
							},
						},
					},
				},
			},
		},
		Update: "update",
		Tables: []Expression{
			Table{
				Object: Identifier{Literal: "table1"},
			},
		},
		Set: "set",
		SetList: []Expression{
			UpdateSet{
				Field: FieldReference{Column: Identifier{Literal: "column1"}},
				Value: NewIntegerValue(1),
			},
			UpdateSet{
				Field: FieldReference{Column: Identifier{Literal: "column2"}},
				Value: NewIntegerValue(2),
			},
		},
		FromClause: FromClause{
			From: "from",
			Tables: []Expression{
				Table{
					Object: Identifier{Literal: "table1"},
				},
			},
		},
		WhereClause: WhereClause{
			Where: "where",
			Filter: Comparison{
				LHS:      Identifier{Literal: "column3"},
				Operator: "=",
				RHS:      NewIntegerValue(3),
			},
		},
	}
	expect := "with ct as (select 1) update table1 set column1 = 1, column2 = 2 from table1 where column3 = 3"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestUpdateSet_String(t *testing.T) {
	e := UpdateSet{
		Field: FieldReference{Column: Identifier{Literal: "column1"}},
		Value: NewIntegerValue(1),
	}
	expect := "column1 = 1"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestDeleteQuery_String(t *testing.T) {
	e := DeleteQuery{
		WithClause: WithClause{
			With: "with",
			InlineTables: []Expression{
				InlineTable{
					Name: Identifier{Literal: "ct"},
					As:   "as",
					Query: SelectQuery{
						SelectEntity: SelectEntity{
							SelectClause: SelectClause{
								Select: "select",
								Fields: []Expression{
									Field{Object: NewIntegerValue(1)},
								},
							},
						},
					},
				},
			},
		},
		Delete: "delete",
		Tables: []Expression{
			Table{
				Object: Identifier{Literal: "table1"},
			},
		},
		FromClause: FromClause{
			From: "from",
			Tables: []Expression{
				Table{
					Object: Identifier{Literal: "table1"},
				},
			},
		},
		WhereClause: WhereClause{
			Where: "where",
			Filter: Comparison{
				LHS:      Identifier{Literal: "column1"},
				Operator: "=",
				RHS:      NewIntegerValue(1),
			},
		},
	}
	expect := "with ct as (select 1) delete table1 from table1 where column1 = 1"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestColumnDefault_String(t *testing.T) {
	e := ColumnDefault{
		Column:  Identifier{Literal: "column1"},
		Default: "default",
		Value:   NewNullValueFromString("null"),
	}
	expect := "column1 default null"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestColumnPosition_String(t *testing.T) {
	e := ColumnPosition{
		Position: Token{Token: AFTER, Literal: "after"},
		Column:   Identifier{Literal: "column1"},
	}
	expect := "after column1"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestFetchPosition_String(t *testing.T) {
	e := FetchPosition{
		Position: Token{Token: ABSOLUTE, Literal: "absolute"},
		Number:   NewIntegerValue(1),
	}
	expect := "absolute 1"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestCursorStatus_String(t *testing.T) {
	e := CursorStatus{
		CursorLit: "cursor",
		Cursor:    Identifier{Literal: "cur"},
		Is:        "is",
		Negation:  Token{Token: NOT, Literal: "not"},
		Type:      RANGE,
		TypeLit:   "in range",
	}
	expect := "cursor cur is not in range"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestCursorAttrebute_String(t *testing.T) {
	e := CursorAttrebute{
		CursorLit: "cursor",
		Cursor:    Identifier{Literal: "cur"},
		Attrebute: Token{Token: COUNT, Literal: "count"},
	}
	expect := "cursor cur count"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}
