package parser

import (
	"testing"
	"time"

	"github.com/mithrandie/csvq/lib/ternary"
)

func TestIsPrimary(t *testing.T) {
	var e Expression

	e = Identifier{Literal: "foo"}
	if IsPrimary(e) {
		t.Errorf("expression %#v is evaluated as is a implementation of primary, but it is not so", e)
	}

	e = Integer{literal: "1"}
	if !IsPrimary(e) {
		t.Errorf("expression %#v is evaluated as is not a implementation of primary, but it is so", e)
	}

	if IsPrimary(nil) {
		t.Error("nil is evaluated as is a implementation of primary,, want empty string for nil")
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

func TestString_Bool(t *testing.T) {
	s := "true"
	p := NewString(s)
	if p.Bool() != true {
		t.Errorf("bool = %t, want %t for %#v", p.Bool(), true, p)
	}

	s = "false"
	p = NewString(s)
	if p.Bool() != false {
		t.Errorf("bool = %t, want %t for %#v", p.Bool(), false, p)
	}

	s = "error"
	p = NewString(s)
	if p.Bool() != false {
		t.Errorf("bool = %t, want %t for %#v", p.Bool(), false, p)
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

func TestInteger_Bool(t *testing.T) {
	p := NewInteger(1)
	if p.Bool() != true {
		t.Errorf("bool = %t, want %t for %#v", p.Bool(), true, p)
	}

	p = NewInteger(0)
	if p.Bool() != false {
		t.Errorf("bool = %t, want %t for %#v", p.Bool(), false, p)
	}
}

func TestInteger_Ternary(t *testing.T) {
	p := NewInteger(1)
	if p.Ternary() != ternary.TRUE {
		t.Errorf("ternary = %s, want %s for %#v", p.Ternary(), ternary.TRUE, p)
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

func TestFloat_Bool(t *testing.T) {
	p := NewFloat(1)
	if p.Bool() != true {
		t.Errorf("bool = %t, want %t for %#v", p.Bool(), true, p)
	}

	p = NewFloat(0)
	if p.Bool() != false {
		t.Errorf("bool = %t, want %t for %#v", p.Bool(), false, p)
	}
}

func TestFloat_Ternary(t *testing.T) {
	p := NewFloat(1)
	if p.Ternary() != ternary.TRUE {
		t.Errorf("ternary = %s, want %s for %#v", p.Ternary(), ternary.TRUE, p)
	}
}

func TestBoolean_String(t *testing.T) {
	s := "true"
	p := NewBoolean(true)
	if p.String() != s {
		t.Errorf("string = %q, want %q for %#v", p.String(), s, p)
	}
}

func TestBoolean_Bool(t *testing.T) {
	p := NewBoolean(true)
	if p.Bool() != true {
		t.Errorf("bool = %t, want %t for %#v", p.Bool(), true, p)
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
}

func TestTernary_Bool(t *testing.T) {
	p := NewTernary(ternary.TRUE)
	if p.Bool() != true {
		t.Errorf("bool = %t, want %t for %#v", p.Bool(), true, p)
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

func TestDatetime_Bool(t *testing.T) {
	d := NewDatetime(time.Time{})

	if d.Bool() != false {
		t.Errorf("bool = %t, want %t for %#v", d.Bool(), false, d)
	}

	d = NewDatetimeFromString("2000-01-01 00:00:00")

	if d.Bool() != true {
		t.Errorf("bool = %t, want %t for %#v", d.Bool(), true, d)
	}
}

func TestDatetime_Ternary(t *testing.T) {
	p := NewDatetime(time.Time{})
	if p.Ternary() != ternary.FALSE {
		t.Errorf("ternary = %s, want %s for %#v", p.Ternary(), ternary.FALSE, p)
	}
}

func TestDatetime_Format(t *testing.T) {
	dtstring := "2012-08-01T04:03:05.123-08:00"
	dt := NewDatetimeFromString(dtstring)
	expect := "2012-08-01 04:03:05.123"
	if dt.Format() != expect {
		t.Errorf("result = %q, want %q for %q ", dt.Format(), expect, dtstring)
	}
}

var datetimeSetFormatTests = []struct {
	Datetime string
	Format   string
	Result   string
}{
	{
		Datetime: "2012-08-01T04:03:05.123-08:00",
		Format:   "datetime: %Y-%m-%d %H:%i:%s %% %g",
		Result:   "datetime: 2012-08-01 04:03:05 % g",
	},
	{
		Datetime: "2012-08-01T04:03:05.123-08:00",
		Format:   "%a",
		Result:   "Wed",
	},
	{
		Datetime: "2012-08-01T04:03:05.123-08:00",
		Format:   "%b",
		Result:   "Aug",
	},
	{
		Datetime: "2012-08-01T04:03:05.123-08:00",
		Format:   "%c",
		Result:   "8",
	},
	{
		Datetime: "2012-08-01T04:03:05.123-08:00",
		Format:   "%E",
		Result:   " 1",
	},
	{
		Datetime: "2012-08-01T04:03:05.123-08:00",
		Format:   "%e",
		Result:   "1",
	},
	{
		Datetime: "2012-08-01T04:03:05.123-08:00",
		Format:   "%F",
		Result:   ".123",
	},
	{
		Datetime: "2012-08-01T04:03:05.123-08:00",
		Format:   "%f",
		Result:   ".123000",
	},
	{
		Datetime: "2012-08-01T04:03:05.123-08:00",
		Format:   "%h",
		Result:   "04",
	},
	{
		Datetime: "2012-08-01T04:03:05.123-08:00",
		Format:   "%l",
		Result:   "4",
	},
	{
		Datetime: "2012-08-01T04:03:05.123-08:00",
		Format:   "%M",
		Result:   "August",
	},
	{
		Datetime: "2012-08-01T04:03:05.123-08:00",
		Format:   "%N",
		Result:   ".123",
	},
	{
		Datetime: "2012-08-01T04:03:05.123-08:00",
		Format:   "%n",
		Result:   ".123000000",
	},
	{
		Datetime: "2012-08-01T04:03:05.123-08:00",
		Format:   "%p",
		Result:   "AM",
	},
	{
		Datetime: "2012-08-01T04:03:05.123-08:00",
		Format:   "%r",
		Result:   "04:03:05 AM",
	},
	{
		Datetime: "2012-08-01T04:03:05.123-08:00",
		Format:   "%T",
		Result:   "04:03:05",
	},
	{
		Datetime: "2012-08-01T04:03:05.123-08:00",
		Format:   "%W",
		Result:   "Wednesday",
	},
	{
		Datetime: "2012-08-01T04:03:05.123-08:00",
		Format:   "%y",
		Result:   "12",
	},
	{
		Datetime: "2012-08-01T04:03:05.123-08:00",
		Format:   "%Z",
		Result:   "-08:00",
	},
	{
		Datetime: "2012-08-01 04:03:05.123 PDT",
		Format:   "%z",
		Result:   "PDT",
	},
}

func TestDatetime_SetFormat(t *testing.T) {
	for _, v := range datetimeSetFormatTests {
		dt := NewDatetimeFromString(v.Datetime)
		dt.SetFormat(v.Format)
		if dt.Format() != v.Result {
			t.Errorf("result = %q, want %q for %q with %q", dt.Format(), v.Result, v.Format, v.Datetime)
		}
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

func TestNull_Bool(t *testing.T) {
	p := NewNull()
	if p.Bool() != false {
		t.Errorf("bool = %t, want %t for %#v", p.Bool(), false, p)
	}
}

func TestNull_Ternary(t *testing.T) {
	p := NewNull()
	if p.Ternary() != ternary.UNKNOWN {
		t.Errorf("ternary = %s, want %s for %#v", p.Ternary(), ternary.UNKNOWN, p)
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

func TestParentheses_String(t *testing.T) {
	s := "abcde"
	e := Parentheses{Expr: String{literal: s}}
	expect := "('abcde')"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestSelectQuery_String(t *testing.T) {
	e := SelectQuery{
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
					Item: Identifier{Literal: "column"},
				},
			},
		},
		LimitClause: LimitClause{
			Limit:  "limit",
			Number: 10,
		},
	}
	expect := "select column from table order by column limit 10"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestSelectSet_String(t *testing.T) {
	e := SelectSet{
		LHS: SelectEntity{
			SelectClause: SelectClause{
				Select: "select",
				Fields: []Expression{Field{Object: NewInteger(1)}},
			},
		},
		Operator: Token{Token: UNION, Literal: "union"},
		All:      Token{Token: ALL, Literal: "all"},
		RHS: SelectEntity{
			SelectClause: SelectClause{
				Select: "select",
				Fields: []Expression{Field{Object: NewInteger(2)}},
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
				Operator: Token{Token: COMPARISON_OP, Literal: ">"},
				RHS:      Integer{literal: "1"},
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
				Operator: Token{Token: COMPARISON_OP, Literal: ">"},
				RHS:      Integer{literal: "1"},
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
				As:     Token{Token: AS, Literal: "as"},
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
			Operator: Token{Token: COMPARISON_OP, Literal: ">"},
			RHS:      Integer{literal: "1"},
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
			Operator: Token{Token: COMPARISON_OP, Literal: ">"},
			RHS:      Integer{literal: "1"},
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
				Item: Identifier{Literal: "column1"},
			},
			OrderItem{
				Item:      Identifier{Literal: "column2"},
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
	e := LimitClause{Limit: "limit", Number: 10}
	expect := "limit 10"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestSubquery_String(t *testing.T) {
	e := Subquery{
		Query: SelectQuery{
			SelectEntity: SelectEntity{
				SelectClause: SelectClause{
					Select: "select",
					Fields: []Expression{
						Integer{literal: "1"},
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
		Operator: Token{Token: COMPARISON_OP, Literal: ">"},
		RHS:      Integer{literal: "1"},
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
		RHS:      Null{literal: "null"},
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
		Low:      Integer{literal: "-10"},
		High:     Integer{literal: "10"},
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
		List: []Expression{
			Integer{literal: "1"},
			Integer{literal: "2"},
		},
		Negation: Token{Token: NOT, Literal: "not"},
	}
	expect := "column not in (1, 2)"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}

	e = In{
		In:  "in",
		LHS: Identifier{Literal: "column"},
		Query: Subquery{
			Query: SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						Select: "select",
						Fields: []Expression{
							Integer{literal: "1"},
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
	expect = "column in (select 1 from dual)"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestAll_String(t *testing.T) {
	e := All{
		All:      "all",
		LHS:      Identifier{Literal: "column"},
		Operator: Token{Token: COMPARISON_OP, Literal: ">"},
		Query: Subquery{
			Query: SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						Select: "select",
						Fields: []Expression{
							Integer{literal: "1"},
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
	expect := "column > all (select 1 from dual)"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestAny_String(t *testing.T) {
	e := Any{
		Any:      "any",
		LHS:      Identifier{Literal: "column"},
		Operator: Token{Token: COMPARISON_OP, Literal: ">"},
		Query: Subquery{
			Query: SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						Select: "select",
						Fields: []Expression{
							Integer{literal: "1"},
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
	expect := "column > any (select 1 from dual)"
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
		Pattern:  String{literal: "pattern"},
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
							Integer{literal: "1"},
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
		RHS:      Integer{literal: "2"},
	}
	expect := "column + 2"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestLogic_String(t *testing.T) {
	e := Logic{
		LHS:      Boolean{literal: "true"},
		Operator: Token{Token: AND, Literal: "and"},
		RHS:      Boolean{literal: "false"},
	}
	expect := "true and false"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}

	e = Logic{
		Operator: Token{Token: NOT, Literal: "not"},
		RHS:      Boolean{literal: "false"},
	}
	expect = "not false"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestConcat_String(t *testing.T) {
	e := Concat{
		Items: []Expression{
			Identifier{Literal: "column"},
			String{literal: "a"},
		},
	}
	expect := "column || 'a'"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestOption_IsDistinct(t *testing.T) {
	e := Option{}
	if e.IsDistinct() == true {
		t.Errorf("distinct = %t, want %t for %#v", e.IsDistinct(), false, e)
	}

	e = Option{Distinct: Token{Token: DISTINCT, Literal: "distinct"}}
	if e.IsDistinct() == false {
		t.Errorf("distinct = %t, want %t for %#v", e.IsDistinct(), true, e)
	}
}

func TestOption_String(t *testing.T) {
	e := Option{
		Distinct: Token{Token: DISTINCT, Literal: "distinct"},
		Args: []Expression{
			Identifier{Literal: "column"},
			String{literal: "a"},
		},
	}
	expect := "distinct column, 'a'"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}

	e = Option{}
	expect = ""
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestFunction_String(t *testing.T) {
	e := Function{
		Name: "sum",
		Option: Option{
			Args: []Expression{
				Identifier{Literal: "column"},
			},
		},
	}
	expect := "sum(column)"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestTable_String(t *testing.T) {
	e := Table{
		Object: Identifier{Literal: "table"},
		As:     Token{Token: AS, Literal: "as"},
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
		As:     Token{Token: AS, Literal: "as"},
		Alias:  Identifier{Literal: "alias"},
	}
	expect := "alias"
	if e.Name() != expect {
		t.Errorf("name = %q, want %q for %#v", e.Name(), expect, e)
	}

	e = Table{
		Object: Identifier{Literal: "/path/to/table.csv"},
	}
	expect = "table"
	if e.Name() != expect {
		t.Errorf("name = %q, want %q for %#v", e.Name(), expect, e)
	}

	e = Table{
		Object: Subquery{
			Query: SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						Select: "select",
						Fields: []Expression{
							Integer{literal: "1"},
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
	expect = "(select 1 from dual)"
	if e.Name() != expect {
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
			Operator: Token{Token: COMPARISON_OP, Literal: ">"},
			RHS:      Integer{literal: "1"},
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
		As:     Token{Token: AS, Literal: "as"},
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
		Object: String{literal: "foo"},
	}
	expect = "foo"
	if e.Name() != expect {
		t.Errorf("name = %q, want %q for %#v", e.Name(), expect, e)
	}

	e = Field{
		Object: NewDatetimeFromString("2012-01-01"),
	}
	expect = "2012-01-01"
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
		Item:      Identifier{Literal: "column"},
		Direction: Token{Token: DESC, Literal: "desc"},
	}
	expect := "column desc"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}

	e = OrderItem{
		Item: Identifier{Literal: "column"},
	}
	expect = "column"
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
				Condition: Integer{literal: "1"},
				Result:    String{literal: "A"},
			},
			CaseWhen{
				When:      "when",
				Then:      "then",
				Condition: Integer{literal: "2"},
				Result:    String{literal: "B"},
			},
		},
		Else: CaseElse{Else: "else", Result: String{literal: "C"}},
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
					Operator: Token{Token: COMPARISON_OP, Literal: ">"},
					RHS:      Integer{literal: "1"},
				},
				Result: String{literal: "A"},
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
			Operator: Token{Token: COMPARISON_OP, Literal: ">"},
			RHS:      Integer{literal: "1"},
		},
		Result: String{literal: "abcde"},
	}
	expect := "when column > 1 then 'abcde'"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestCaseElse_String(t *testing.T) {
	e := CaseElse{Else: "else", Result: String{literal: "abcde"}}
	expect := "else 'abcde'"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestGroupConcat_String(t *testing.T) {
	e := GroupConcat{
		GroupConcat: "group_concat",
		Option:      Option{Args: []Expression{Identifier{Literal: "column1"}}},
		OrderBy: OrderByClause{
			OrderBy: "order by",
			Items:   []Expression{Identifier{Literal: "column1"}},
		},
		SeparatorLit: "separator",
		Separator:    ",",
	}
	expect := "group_concat(column1 order by column1 separator ',')"
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
		Value: NewInteger(1),
	}
	expect := "@var := 1"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestVariableAssignment_String(t *testing.T) {
	e := VariableAssignment{
		Name:  "@var",
		Value: NewInteger(1),
	}
	expect := "@var = 1"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}

	e = VariableAssignment{
		Name:  "@var",
		Value: nil,
	}
	expect = "@var"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestInsertQuery_String(t *testing.T) {
	e := InsertQuery{
		Insert: "insert",
		Into:   "into",
		Table:  Identifier{Literal: "table1"},
		Fields: []Expression{
			Identifier{Literal: "column1"},
			Identifier{Literal: "column2"},
		},
		Values: "values",
		ValuesList: []Expression{
			InsertValues{
				Values: []Expression{
					NewInteger(1),
					NewInteger(2),
				},
			},
			InsertValues{
				Values: []Expression{
					NewInteger(3),
					NewInteger(4),
				},
			},
		},
	}
	expect := "insert into table1 (column1, column2) values (1, 2), (3, 4)"
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
						NewInteger(1),
						NewInteger(2),
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

func TestInsertValues_String(t *testing.T) {
	e := InsertValues{
		Values: []Expression{
			NewInteger(1),
			NewInteger(2),
		},
	}
	expect := "(1, 2)"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestUpdateQuery_String(t *testing.T) {
	e := UpdateQuery{
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
				Value: NewInteger(1),
			},
			UpdateSet{
				Field: FieldReference{Column: Identifier{Literal: "column2"}},
				Value: NewInteger(2),
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
				Operator: Token{Token: COMPARISON_OP, Literal: "="},
				RHS:      NewInteger(3),
			},
		},
	}
	expect := "update table1 set column1 = 1, column2 = 2 from table1 where column3 = 3"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestUpdateSet_String(t *testing.T) {
	e := UpdateSet{
		Field: FieldReference{Column: Identifier{Literal: "column1"}},
		Value: NewInteger(1),
	}
	expect := "column1 = 1"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestDeleteQuery_String(t *testing.T) {
	e := DeleteQuery{
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
				Operator: Token{Token: COMPARISON_OP, Literal: "="},
				RHS:      NewInteger(1),
			},
		},
	}
	expect := "delete table1 from table1 where column1 = 1"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestCreateTable_String(t *testing.T) {
	e := CreateTable{
		CreateTable: "create table",
		Table:       Identifier{Literal: "newtable"},
		Fields: []Expression{
			Identifier{Literal: "column1"},
			Identifier{Literal: "column2"},
		},
	}
	expect := "create table newtable (column1, column2)"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestAddColumns_String(t *testing.T) {
	e := AddColumns{
		AlterTable: "alter table",
		Table:      Identifier{Literal: "table1"},
		Add:        "add",
		Columns: []Expression{
			ColumnDefault{
				Column: Identifier{Literal: "newcolumn"},
			},
			ColumnDefault{
				Column:  Identifier{Literal: "newcolumn2"},
				Default: "default",
				Value:   Null{literal: "null"},
			},
		},
		Position: ColumnPosition{
			Position: Token{Token: FIRST, Literal: "first"},
		},
	}
	expect := "alter table table1 add (newcolumn, newcolumn2 default null) first"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestColumnDefault_String(t *testing.T) {
	e := ColumnDefault{
		Column:  Identifier{Literal: "column1"},
		Default: "default",
		Value:   Null{literal: "null"},
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

func TestDropColumns_String(t *testing.T) {
	e := DropColumns{
		AlterTable: "alter table",
		Table:      Identifier{Literal: "table1"},
		Drop:       "drop",
		Columns: []Expression{
			Identifier{Literal: "column1"},
			Identifier{Literal: "column2"},
		},
	}
	expect := "alter table table1 drop (column1, column2)"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}

func TestRenameColumn_String(t *testing.T) {
	e := RenameColumn{
		AlterTable: "alter table",
		Table:      Identifier{Literal: "table1"},
		Rename:     "rename",
		Old:        FieldReference{Column: Identifier{Literal: "oldcolumn"}},
		To:         "to",
		New:        Identifier{Literal: "newcolumn"},
	}
	expect := "alter table table1 rename oldcolumn to newcolumn"
	if e.String() != expect {
		t.Errorf("string = %q, want %q for %#v", e.String(), expect, e)
	}
}
