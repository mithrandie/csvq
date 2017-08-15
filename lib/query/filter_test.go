package query

import (
	"reflect"
	"testing"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/ternary"
)

var filterEvaluateTests = []struct {
	Name   string
	Filter Filter
	Expr   parser.Expression
	Result parser.Primary
	Error  string
}{
	{
		Name:   "nil",
		Expr:   nil,
		Result: parser.NewTernary(ternary.TRUE),
	},
	{
		Name:   "Primary",
		Expr:   parser.NewString("str"),
		Result: parser.NewString("str"),
	},
	{
		Name:  "Syntax Error",
		Expr:  parser.AllColumns{},
		Error: "[L:- C:-] syntax error: unexpected *",
	},
	{
		Name: "Parentheses",
		Expr: parser.Parentheses{
			Expr: parser.NewString("str"),
		},
		Result: parser.NewString("str"),
	},
	{
		Name: "FieldReference",
		Filter: Filter{
			Records: []FilterRecord{
				{
					View: &View{
						Header: NewHeaderWithId("table1", []string{"column1", "column2"}),
						Records: []Record{
							NewRecordWithId(1, []parser.Primary{
								parser.NewInteger(1),
								parser.NewString("str"),
							}),
							NewRecordWithId(2, []parser.Primary{
								parser.NewInteger(2),
								parser.NewString("strstr"),
							}),
						},
					},
					RecordIndex: 1,
				},
			},
		},
		Expr:   parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
		Result: parser.NewString("strstr"),
	},
	{
		Name: "FieldReference ColumnNotExist Error",
		Filter: Filter{
			Records: []FilterRecord{
				{
					View: &View{
						Header: NewHeaderWithId("table1", []string{"column1", "column2"}),
						Records: []Record{
							NewRecordWithId(1, []parser.Primary{
								parser.NewInteger(1),
								parser.NewString("str"),
							}),
							NewRecordWithId(2, []parser.Primary{
								parser.NewInteger(2),
								parser.NewString("strstr"),
							}),
						},
					},
					RecordIndex: 1,
				},
			},
		},
		Expr:  parser.FieldReference{Column: parser.Identifier{Literal: "column3"}},
		Error: "[L:- C:-] field column3 does not exist",
	},
	{
		Name: "FieldReference FieldAmbigous Error",
		Filter: Filter{
			Records: []FilterRecord{
				{
					View: &View{
						Header: NewHeaderWithId("table1", []string{"column1", "column1"}),
						Records: []Record{
							NewRecordWithId(1, []parser.Primary{
								parser.NewInteger(1),
								parser.NewString("str"),
							}),
							NewRecordWithId(2, []parser.Primary{
								parser.NewInteger(2),
								parser.NewString("strstr"),
							}),
						},
					},
					RecordIndex: 1,
				},
			},
		},
		Expr:  parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
		Error: "[L:- C:-] field column1 is ambiguous",
	},
	{
		Name: "FieldReference Not Group Key Error",
		Filter: Filter{
			Records: []FilterRecord{
				{
					View: &View{
						Header: []HeaderField{
							{
								View:      "table1",
								Column:    "column1",
								FromTable: true,
							},
							{
								View:   "table1",
								Column: "column2",
							},
						},
						Records: []Record{
							{
								NewGroupCell([]parser.Primary{
									parser.NewInteger(1),
									parser.NewInteger(2),
								}),
							},
							{
								NewGroupCell([]parser.Primary{
									parser.NewString("str1"),
									parser.NewString("str2"),
								}),
							},
						},
						isGrouped: true,
					},
					RecordIndex: 0,
				},
			},
		},
		Expr:  parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
		Error: "[L:- C:-] field column1 is not a group key",
	},
	{
		Name: "ColumnNumber",
		Filter: Filter{
			Records: []FilterRecord{
				{
					View: &View{
						Header: NewHeaderWithId("table1", []string{"column1", "column2"}),
						Records: []Record{
							NewRecordWithId(1, []parser.Primary{
								parser.NewInteger(1),
								parser.NewString("str"),
							}),
							NewRecordWithId(2, []parser.Primary{
								parser.NewInteger(2),
								parser.NewString("strstr"),
							}),
						},
					},
					RecordIndex: 1,
				},
			},
		},
		Expr:   parser.ColumnNumber{View: parser.Identifier{Literal: "table1"}, Number: parser.NewInteger(2)},
		Result: parser.NewString("strstr"),
	},
	{
		Name: "ColumnNumber Not Exist Error",
		Filter: Filter{
			Records: []FilterRecord{
				{
					View: &View{
						Header: NewHeaderWithId("table1", []string{"column1", "column2"}),
						Records: []Record{
							NewRecordWithId(1, []parser.Primary{
								parser.NewInteger(1),
								parser.NewString("str"),
							}),
							NewRecordWithId(2, []parser.Primary{
								parser.NewInteger(2),
								parser.NewString("strstr"),
							}),
						},
					},
					RecordIndex: 1,
				},
			},
		},
		Expr:  parser.ColumnNumber{View: parser.Identifier{Literal: "table1"}, Number: parser.NewInteger(9)},
		Error: "[L:- C:-] field number table1.9 does not exist",
	},
	{
		Name: "ColumnNumber Not Group Key Error",
		Filter: Filter{
			Records: []FilterRecord{
				{
					View: &View{
						Header: []HeaderField{
							{
								View:      "table1",
								Column:    "column1",
								Number:    1,
								FromTable: true,
							},
							{
								View:       "table1",
								Column:     "column2",
								Number:     2,
								FromTable:  true,
								IsGroupKey: true,
							},
						},
						Records: []Record{
							{
								NewGroupCell([]parser.Primary{
									parser.NewInteger(1),
									parser.NewInteger(2),
								}),
							},
							{
								NewGroupCell([]parser.Primary{
									parser.NewString("str1"),
									parser.NewString("str2"),
								}),
							},
						},
						isGrouped: true,
					},
					RecordIndex: 0,
				},
			},
		},
		Expr:  parser.ColumnNumber{View: parser.Identifier{Literal: "table1"}, Number: parser.NewInteger(1)},
		Error: "[L:- C:-] field number table1.1 is not a group key",
	},
	{
		Name: "Arithmetic",
		Expr: parser.Arithmetic{
			LHS:      parser.NewInteger(1),
			RHS:      parser.NewInteger(2),
			Operator: '+',
		},
		Result: parser.NewInteger(3),
	},
	{
		Name: "Arithmetic LHS Error",
		Expr: parser.Arithmetic{
			LHS:      parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			RHS:      parser.NewInteger(2),
			Operator: '+',
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
	{
		Name: "Arithmetic RHS Error",
		Expr: parser.Arithmetic{
			LHS:      parser.NewInteger(1),
			RHS:      parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			Operator: '+',
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
	{
		Name: "UnaryArithmetic",
		Expr: parser.UnaryArithmetic{
			Operand:  parser.NewInteger(1),
			Operator: parser.Token{Token: '-', Literal: "-"},
		},
		Result: parser.NewInteger(-1),
	},
	{
		Name: "UnaryArithmetic Operand Error",
		Expr: parser.UnaryArithmetic{
			Operand:  parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			Operator: parser.Token{Token: '-', Literal: "-"},
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
	{
		Name: "UnaryArithmetic Cast Failure",
		Expr: parser.UnaryArithmetic{
			Operand:  parser.NewString("str"),
			Operator: parser.Token{Token: '-', Literal: "-"},
		},
		Result: parser.NewNull(),
	},
	{
		Name: "Concat",
		Expr: parser.Concat{
			Items: []parser.Expression{
				parser.NewString("a"),
				parser.NewString("b"),
				parser.NewString("c"),
			},
		},
		Result: parser.NewString("abc"),
	},
	{
		Name: "Concat FieldNotExist Error",
		Expr: parser.Concat{
			Items: []parser.Expression{
				parser.NewString("a"),
				parser.NewString("b"),
				parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			},
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
	{
		Name: "Concat Including Null",
		Expr: parser.Concat{
			Items: []parser.Expression{
				parser.NewString("a"),
				parser.NewNull(),
				parser.NewString("c"),
			},
		},
		Result: parser.NewNull(),
	},
	{
		Name: "Comparison",
		Expr: parser.Comparison{
			LHS:      parser.NewInteger(1),
			RHS:      parser.NewInteger(2),
			Operator: "=",
		},
		Result: parser.NewTernary(ternary.FALSE),
	},
	{
		Name: "Comparison LHS Error",
		Expr: parser.Comparison{
			LHS:      parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			RHS:      parser.NewInteger(2),
			Operator: "=",
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
	{
		Name: "Comparison RHS Error",
		Expr: parser.Comparison{
			LHS:      parser.NewInteger(1),
			RHS:      parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			Operator: "=",
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
	{
		Name: "Comparison with Row Values",
		Expr: parser.Comparison{
			LHS: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.Expression{
						parser.NewInteger(1),
						parser.NewInteger(2),
					},
				},
			},
			RHS: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.Expression{
						parser.NewInteger(1),
						parser.NewInteger(2),
					},
				},
			},
			Operator: "=",
		},
		Result: parser.NewTernary(ternary.TRUE),
	},
	{
		Name: "Comparison with Row Value and Subquery",
		Expr: parser.Comparison{
			LHS: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.Expression{
						parser.NewString("1"),
						parser.NewString("str1"),
					},
				},
			},
			RHS: parser.RowValue{
				Value: parser.Subquery{
					Query: parser.SelectQuery{
						SelectEntity: parser.SelectEntity{
							SelectClause: parser.SelectClause{
								Select: "select",
								Fields: []parser.Expression{
									parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
									parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
								},
							},
							FromClause: parser.FromClause{
								Tables: []parser.Expression{
									parser.Table{Object: parser.Identifier{Literal: "table1"}},
								},
							},
							WhereClause: parser.WhereClause{
								Filter: parser.Comparison{
									LHS:      parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
									RHS:      parser.NewInteger(1),
									Operator: "=",
								},
							},
						},
					},
				},
			},
			Operator: "=",
		},
		Result: parser.NewTernary(ternary.TRUE),
	},
	{
		Name: "Comparison with Row Values LHS Error",
		Expr: parser.Comparison{
			LHS: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.Expression{
						parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
						parser.NewInteger(2),
					},
				},
			},
			RHS: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.Expression{
						parser.NewInteger(1),
						parser.NewInteger(2),
					},
				},
			},
			Operator: "=",
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
	{
		Name: "Comparison with Row Value and Subquery Returns No Record",
		Expr: parser.Comparison{
			LHS: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.Expression{
						parser.NewString("1"),
						parser.NewString("str1"),
					},
				},
			},
			RHS: parser.RowValue{
				Value: parser.Subquery{
					Query: parser.SelectQuery{
						SelectEntity: parser.SelectEntity{
							SelectClause: parser.SelectClause{
								Select: "select",
								Fields: []parser.Expression{
									parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
									parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
								},
							},
							FromClause: parser.FromClause{
								Tables: []parser.Expression{
									parser.Table{Object: parser.Identifier{Literal: "table1"}},
								},
							},
							WhereClause: parser.WhereClause{
								Filter: parser.NewTernary(ternary.FALSE),
							},
						},
					},
				},
			},
			Operator: "=",
		},
		Result: parser.NewTernary(ternary.UNKNOWN),
	},
	{
		Name: "Comparison with Row Value and Subquery Query Error",
		Expr: parser.Comparison{
			LHS: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.Expression{
						parser.NewString("1"),
						parser.NewString("str1"),
					},
				},
			},
			RHS: parser.RowValue{
				Value: parser.Subquery{
					Query: parser.SelectQuery{
						SelectEntity: parser.SelectEntity{
							SelectClause: parser.SelectClause{
								Select: "select",
								Fields: []parser.Expression{
									parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
									parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}}},
								},
							},
							FromClause: parser.FromClause{
								Tables: []parser.Expression{
									parser.Table{Object: parser.Identifier{Literal: "table1"}},
								},
							},
						},
					},
				},
			},
			Operator: "=",
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
	{
		Name: "Comparison with Row Values RHS Error",
		Expr: parser.Comparison{
			LHS: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.Expression{
						parser.NewInteger(1),
						parser.NewInteger(2),
					},
				},
			},
			RHS: parser.RowValue{
				Value: parser.Subquery{
					Query: parser.SelectQuery{
						SelectEntity: parser.SelectEntity{
							SelectClause: parser.SelectClause{
								Select: "select",
								Fields: []parser.Expression{
									parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
									parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
								},
							},
							FromClause: parser.FromClause{
								Tables: []parser.Expression{
									parser.Table{Object: parser.Identifier{Literal: "table1"}},
								},
							},
						},
					},
				},
			},
			Operator: "=",
		},
		Error: "[L:- C:-] subquery returns too many records, should return only one record",
	},
	{
		Name: "Comparison with Row Values Value Length Not Match Error",
		Expr: parser.Comparison{
			LHS: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.Expression{
						parser.NewInteger(1),
						parser.NewInteger(2),
					},
				},
			},
			RHS: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.Expression{
						parser.NewInteger(1),
					},
				},
			},
			Operator: "=",
		},
		Error: "[L:- C:-] row value should contain exactly 2 values",
	}, {
		Name: "Is",
		Expr: parser.Is{
			LHS:      parser.NewInteger(1),
			RHS:      parser.NewNull(),
			Negation: parser.Token{Token: parser.NOT, Literal: "not"},
		},
		Result: parser.NewTernary(ternary.TRUE),
	},
	{
		Name: "Is LHS Error",
		Expr: parser.Is{
			LHS:      parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			RHS:      parser.NewNull(),
			Negation: parser.Token{Token: parser.NOT, Literal: "not"},
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
	{
		Name: "Is RHS Error",
		Expr: parser.Is{
			LHS:      parser.NewInteger(1),
			RHS:      parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			Negation: parser.Token{Token: parser.NOT, Literal: "not"},
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
	{
		Name: "Between",
		Expr: parser.Between{
			LHS:      parser.NewInteger(2),
			Low:      parser.NewInteger(1),
			High:     parser.NewInteger(3),
			Negation: parser.Token{Token: parser.NOT, Literal: "not"},
		},
		Result: parser.NewTernary(ternary.FALSE),
	},
	{
		Name: "Between LHS Error",
		Expr: parser.Between{
			LHS:      parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			Low:      parser.NewInteger(1),
			High:     parser.NewInteger(3),
			Negation: parser.Token{Token: parser.NOT, Literal: "not"},
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
	{
		Name: "Between Low Error",
		Expr: parser.Between{
			LHS:      parser.NewInteger(2),
			Low:      parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			High:     parser.NewInteger(3),
			Negation: parser.Token{Token: parser.NOT, Literal: "not"},
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
	{
		Name: "Between High Error",
		Expr: parser.Between{
			LHS:      parser.NewInteger(2),
			Low:      parser.NewInteger(1),
			High:     parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			Negation: parser.Token{Token: parser.NOT, Literal: "not"},
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
	{
		Name: "Between with Row Values",
		Expr: parser.Between{
			LHS: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.Expression{
						parser.NewInteger(1),
						parser.NewInteger(2),
					},
				},
			},
			Low: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.Expression{
						parser.NewInteger(1),
						parser.NewInteger(1),
					},
				},
			},
			High: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.Expression{
						parser.NewInteger(1),
						parser.NewInteger(3),
					},
				},
			},
		},
		Result: parser.NewTernary(ternary.TRUE),
	},
	{
		Name: "Between with Row Values LHS Error",
		Expr: parser.Between{
			LHS: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.Expression{
						parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
						parser.NewInteger(2),
					},
				},
			},
			Low: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.Expression{
						parser.NewInteger(1),
						parser.NewInteger(1),
					},
				},
			},
			High: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.Expression{
						parser.NewInteger(1),
						parser.NewInteger(3),
					},
				},
			},
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
	{
		Name: "Between with Row Values Low Error",
		Expr: parser.Between{
			LHS: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.Expression{
						parser.NewInteger(1),
						parser.NewInteger(2),
					},
				},
			},
			Low: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.Expression{
						parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
						parser.NewInteger(1),
					},
				},
			},
			High: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.Expression{
						parser.NewInteger(1),
						parser.NewInteger(3),
					},
				},
			},
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
	{
		Name: "Between with Row Values High Error",
		Expr: parser.Between{
			LHS: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.Expression{
						parser.NewInteger(1),
						parser.NewInteger(2),
					},
				},
			},
			Low: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.Expression{
						parser.NewInteger(1),
						parser.NewInteger(1),
					},
				},
			},
			High: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.Expression{
						parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
						parser.NewInteger(3),
					},
				},
			},
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
	{
		Name: "Between with Row Values Low Comparison Error",
		Expr: parser.Between{
			LHS: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.Expression{
						parser.NewInteger(1),
						parser.NewInteger(2),
					},
				},
			},
			Low: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.Expression{
						parser.NewInteger(1),
					},
				},
			},
			High: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.Expression{
						parser.NewInteger(1),
						parser.NewInteger(3),
					},
				},
			},
		},
		Error: "[L:- C:-] row value should contain exactly 2 values",
	},
	{
		Name: "Between with Row Values High Comparison Error",
		Expr: parser.Between{
			LHS: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.Expression{
						parser.NewInteger(1),
						parser.NewInteger(2),
					},
				},
			},
			Low: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.Expression{
						parser.NewInteger(1),
						parser.NewInteger(1),
					},
				},
			},
			High: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.Expression{
						parser.NewInteger(3),
					},
				},
			},
		},
		Error: "[L:- C:-] row value should contain exactly 2 values",
	},
	{
		Name: "In",
		Expr: parser.In{
			LHS: parser.NewInteger(2),
			Values: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.Expression{
						parser.NewInteger(1),
						parser.NewInteger(2),
						parser.NewInteger(3),
					},
				},
			},
			Negation: parser.Token{Token: parser.NOT, Literal: "not"},
		},
		Result: parser.NewTernary(ternary.FALSE),
	},
	{
		Name: "In LHS Error",
		Expr: parser.In{
			LHS: parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			Values: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.Expression{
						parser.NewInteger(1),
						parser.NewInteger(2),
						parser.NewInteger(3),
					},
				},
			},
			Negation: parser.Token{Token: parser.NOT, Literal: "not"},
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
	{
		Name: "In List Error",
		Expr: parser.In{
			LHS: parser.NewInteger(2),
			Values: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.Expression{
						parser.NewInteger(1),
						parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
						parser.NewInteger(3),
					},
				},
			},
			Negation: parser.Token{Token: parser.NOT, Literal: "not"},
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
	{
		Name: "In Subquery",
		Filter: Filter{
			Records: []FilterRecord{
				{
					View: &View{
						Header: NewHeaderWithId("table2", []string{"column3", "column4"}),
						Records: []Record{
							NewRecordWithId(1, []parser.Primary{
								parser.NewInteger(1),
								parser.NewString("str2"),
							}),
						},
					},
					RecordIndex: 0,
				},
			},
		},
		Expr: parser.In{
			LHS: parser.NewInteger(2),
			Values: parser.RowValue{
				Value: parser.Subquery{
					Query: parser.SelectQuery{
						SelectEntity: parser.SelectEntity{
							SelectClause: parser.SelectClause{
								Select: "select",
								Fields: []parser.Expression{
									parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
								},
							},
							FromClause: parser.FromClause{
								Tables: []parser.Expression{
									parser.Table{Object: parser.Identifier{Literal: "table1"}},
								},
							},
							WhereClause: parser.WhereClause{
								Filter: parser.Comparison{
									LHS:      parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
									RHS:      parser.FieldReference{View: parser.Identifier{Literal: "table2"}, Column: parser.Identifier{Literal: "column4"}},
									Operator: "=",
								},
							},
						},
					},
				},
			},
			Negation: parser.Token{Token: parser.NOT, Literal: "not"},
		},
		Result: parser.NewTernary(ternary.FALSE),
	},
	{
		Name: "In Subquery Execution Error",
		Expr: parser.In{
			LHS: parser.NewInteger(2),
			Values: parser.RowValue{
				Value: parser.Subquery{
					Query: parser.SelectQuery{
						SelectEntity: parser.SelectEntity{
							SelectClause: parser.SelectClause{
								Select: "select",
								Fields: []parser.Expression{
									parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
								},
							},
							FromClause: parser.FromClause{
								Tables: []parser.Expression{
									parser.Table{Object: parser.Identifier{Literal: "table1"}},
								},
							},
							WhereClause: parser.WhereClause{
								Filter: parser.Comparison{
									LHS:      parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
									RHS:      parser.FieldReference{View: parser.Identifier{Literal: "table2"}, Column: parser.Identifier{Literal: "column4"}},
									Operator: "=",
								},
							},
						},
					},
				},
			},
			Negation: parser.Token{Token: parser.NOT, Literal: "not"},
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
	{
		Name: "In Subquery Too Many Field Error",
		Expr: parser.In{
			LHS: parser.NewInteger(2),
			Values: parser.RowValue{
				Value: parser.Subquery{
					Query: parser.SelectQuery{
						SelectEntity: parser.SelectEntity{
							SelectClause: parser.SelectClause{
								Select: "select",
								Fields: []parser.Expression{
									parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
									parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
								},
							},
							FromClause: parser.FromClause{
								Tables: []parser.Expression{
									parser.Table{Object: parser.Identifier{Literal: "table1"}},
								},
							},
						},
					},
				},
			},
			Negation: parser.Token{Token: parser.NOT, Literal: "not"},
		},
		Error: "[L:- C:-] subquery returns too many fields, should return only one field",
	},
	{
		Name: "In Subquery Returns No Record",
		Expr: parser.In{
			LHS: parser.NewInteger(2),
			Values: parser.RowValue{
				Value: parser.Subquery{
					Query: parser.SelectQuery{
						SelectEntity: parser.SelectEntity{
							SelectClause: parser.SelectClause{
								Select: "select",
								Fields: []parser.Expression{
									parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
								},
							},
							FromClause: parser.FromClause{
								Tables: []parser.Expression{
									parser.Table{Object: parser.Identifier{Literal: "table1"}},
								},
							},
							WhereClause: parser.WhereClause{
								Filter: parser.NewTernary(ternary.FALSE),
							},
						},
					},
				},
			},
		},
		Result: parser.NewTernary(ternary.FALSE),
	},
	{
		Name: "In with Row Values",
		Expr: parser.In{
			LHS: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.Expression{
						parser.NewInteger(1),
						parser.NewInteger(2),
					},
				},
			},
			Values: parser.RowValueList{
				RowValues: []parser.Expression{
					parser.RowValue{
						Value: parser.ValueList{
							Values: []parser.Expression{
								parser.NewInteger(1),
								parser.NewInteger(1),
							},
						},
					},
					parser.RowValue{
						Value: parser.ValueList{
							Values: []parser.Expression{
								parser.NewInteger(1),
								parser.NewInteger(2),
							},
						},
					},
				},
			},
		},
		Result: parser.NewTernary(ternary.TRUE),
	},
	{
		Name: "In with Row Value and Subquery",
		Expr: parser.In{
			LHS: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.Expression{
						parser.NewString("2"),
						parser.NewString("str2"),
					},
				},
			},
			Values: parser.Subquery{
				Query: parser.SelectQuery{
					SelectEntity: parser.SelectEntity{
						SelectClause: parser.SelectClause{
							Select: "select",
							Fields: []parser.Expression{
								parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
								parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
							},
						},
						FromClause: parser.FromClause{
							Tables: []parser.Expression{
								parser.Table{Object: parser.Identifier{Literal: "table1"}},
							},
						},
					},
				},
			},
		},
		Result: parser.NewTernary(ternary.TRUE),
	},
	{
		Name: "In with Row Values LHS Error",
		Expr: parser.In{
			LHS: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.Expression{
						parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
						parser.NewInteger(2),
					},
				},
			},
			Values: parser.RowValueList{
				RowValues: []parser.Expression{
					parser.RowValue{
						Value: parser.ValueList{
							Values: []parser.Expression{
								parser.NewInteger(1),
								parser.NewInteger(1),
							},
						},
					},
					parser.RowValue{
						Value: parser.ValueList{
							Values: []parser.Expression{
								parser.NewInteger(1),
								parser.NewInteger(2),
							},
						},
					},
				},
			},
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
	{
		Name: "In with Row Value and Subquery Query Error",
		Expr: parser.In{
			LHS: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.Expression{
						parser.NewString("2"),
						parser.NewString("str2"),
					},
				},
			},
			Values: parser.Subquery{
				Query: parser.SelectQuery{
					SelectEntity: parser.SelectEntity{
						SelectClause: parser.SelectClause{
							Select: "select",
							Fields: []parser.Expression{
								parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
								parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}}},
							},
						},
						FromClause: parser.FromClause{
							Tables: []parser.Expression{
								parser.Table{Object: parser.Identifier{Literal: "table1"}},
							},
						},
					},
				},
			},
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
	{
		Name: "In with Row Value and Subquery Empty Result Set",
		Expr: parser.In{
			LHS: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.Expression{
						parser.NewString("2"),
						parser.NewString("str2"),
					},
				},
			},
			Values: parser.Subquery{
				Query: parser.SelectQuery{
					SelectEntity: parser.SelectEntity{
						SelectClause: parser.SelectClause{
							Select: "select",
							Fields: []parser.Expression{
								parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
								parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
							},
						},
						FromClause: parser.FromClause{
							Tables: []parser.Expression{
								parser.Table{Object: parser.Identifier{Literal: "table1"}},
							},
						},
						WhereClause: parser.WhereClause{
							Filter: parser.NewTernary(ternary.FALSE),
						},
					},
				},
			},
		},
		Result: parser.NewTernary(ternary.FALSE),
	},
	{
		Name: "In with Row Values Values Error",
		Expr: parser.In{
			LHS: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.Expression{
						parser.NewInteger(1),
						parser.NewInteger(2),
					},
				},
			},
			Values: parser.RowValueList{
				RowValues: []parser.Expression{
					parser.RowValue{
						Value: parser.ValueList{
							Values: []parser.Expression{
								parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
								parser.NewInteger(1),
							},
						},
					},
					parser.RowValue{
						Value: parser.ValueList{
							Values: []parser.Expression{
								parser.NewInteger(1),
								parser.NewInteger(2),
							},
						},
					},
				},
			},
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
	{
		Name: "In with Row Values Length Not Match Error ",
		Expr: parser.In{
			LHS: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.Expression{
						parser.NewInteger(1),
						parser.NewInteger(2),
					},
				},
			},
			Values: parser.RowValueList{
				RowValues: []parser.Expression{
					parser.RowValue{
						Value: parser.ValueList{
							Values: []parser.Expression{
								parser.NewInteger(1),
								parser.NewInteger(1),
							},
						},
					},
					parser.RowValue{
						Value: parser.ValueList{
							Values: []parser.Expression{
								parser.NewInteger(2),
							},
						},
					},
				},
			},
		},
		Error: "[L:- C:-] row value should contain exactly 2 values",
	},
	{
		Name: "In with Row Value and Subquery Length Not Match Error",
		Expr: parser.In{
			LHS: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.Expression{
						parser.NewString("2"),
						parser.NewString("str2"),
					},
				},
			},
			Values: parser.Subquery{
				Query: parser.SelectQuery{
					SelectEntity: parser.SelectEntity{
						SelectClause: parser.SelectClause{
							Select: "select",
							Fields: []parser.Expression{
								parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
							},
						},
						FromClause: parser.FromClause{
							Tables: []parser.Expression{
								parser.Table{Object: parser.Identifier{Literal: "table1"}},
							},
						},
					},
				},
			},
		},
		Error: "[L:- C:-] select query should return exactly 2 fields",
	},
	{
		Name: "Any",
		Expr: parser.Any{
			LHS: parser.NewInteger(5),
			Values: parser.RowValue{
				Value: parser.Subquery{
					Query: parser.SelectQuery{
						SelectEntity: parser.SelectEntity{
							SelectClause: parser.SelectClause{
								Select: "select",
								Fields: []parser.Expression{
									parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
								},
							},
							FromClause: parser.FromClause{
								Tables: []parser.Expression{
									parser.Table{Object: parser.Identifier{Literal: "table1"}},
								},
							},
						},
					},
				},
			},
			Operator: "<>",
		},
		Result: parser.NewTernary(ternary.TRUE),
	},
	{
		Name: "Any LHS Error",
		Expr: parser.Any{
			LHS: parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			Values: parser.RowValue{
				Value: parser.Subquery{
					Query: parser.SelectQuery{
						SelectEntity: parser.SelectEntity{
							SelectClause: parser.SelectClause{
								Select: "select",
								Fields: []parser.Expression{
									parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
								},
							},
							FromClause: parser.FromClause{
								Tables: []parser.Expression{
									parser.Table{Object: parser.Identifier{Literal: "table1"}},
								},
							},
						},
					},
				},
			},
			Operator: "<>",
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
	{
		Name: "Any Query Execution Error",
		Expr: parser.Any{
			LHS: parser.NewInteger(2),
			Values: parser.RowValue{
				Value: parser.Subquery{
					Query: parser.SelectQuery{
						SelectEntity: parser.SelectEntity{
							SelectClause: parser.SelectClause{
								Select: "select",
								Fields: []parser.Expression{
									parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
								},
							},
							FromClause: parser.FromClause{
								Tables: []parser.Expression{
									parser.Table{Object: parser.Identifier{Literal: "table1"}},
								},
							},
							WhereClause: parser.WhereClause{
								Filter: parser.Comparison{
									LHS:      parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
									RHS:      parser.FieldReference{View: parser.Identifier{Literal: "table2"}, Column: parser.Identifier{Literal: "column4"}},
									Operator: "=",
								},
							},
						},
					},
				},
			},
			Operator: "<>",
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
	{
		Name: "Any Row Value Select Field Not Match Error",
		Expr: parser.Any{
			LHS: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.Expression{
						parser.NewInteger(1),
						parser.NewInteger(2),
					},
				},
			},
			Values: parser.Subquery{
				Query: parser.SelectQuery{
					SelectEntity: parser.SelectEntity{
						SelectClause: parser.SelectClause{
							Select: "select",
							Fields: []parser.Expression{
								parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
							},
						},
						FromClause: parser.FromClause{
							Tables: []parser.Expression{
								parser.Table{Object: parser.Identifier{Literal: "table1"}},
							},
						},
					},
				},
			},
			Operator: "<>",
		},
		Error: "[L:- C:-] select query should return exactly 2 fields",
	},
	{
		Name: "Any Row Value Length Not Match Error",
		Expr: parser.Any{
			LHS: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.Expression{
						parser.NewInteger(1),
						parser.NewInteger(2),
					},
				},
			},
			Values: parser.RowValueList{
				RowValues: []parser.Expression{
					parser.RowValue{
						Value: parser.ValueList{
							Values: []parser.Expression{
								parser.NewInteger(1),
								parser.NewInteger(1),
							},
						},
					},
					parser.RowValue{
						Value: parser.ValueList{
							Values: []parser.Expression{
								parser.NewInteger(1),
							},
						},
					},
				},
			},
			Operator: "<>",
		},
		Error: "[L:- C:-] row value should contain exactly 2 values",
	},
	{
		Name: "All",
		Expr: parser.All{
			LHS: parser.NewInteger(5),
			Values: parser.RowValue{
				Value: parser.Subquery{
					Query: parser.SelectQuery{
						SelectEntity: parser.SelectEntity{
							SelectClause: parser.SelectClause{
								Select: "select",
								Fields: []parser.Expression{
									parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
								},
							},
							FromClause: parser.FromClause{
								Tables: []parser.Expression{
									parser.Table{Object: parser.Identifier{Literal: "table1"}},
								},
							},
						},
					},
				},
			},
			Operator: ">",
		},
		Result: parser.NewTernary(ternary.TRUE),
	},
	{
		Name: "All LHS Error",
		Expr: parser.All{
			LHS: parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			Values: parser.RowValue{
				Value: parser.Subquery{
					Query: parser.SelectQuery{
						SelectEntity: parser.SelectEntity{
							SelectClause: parser.SelectClause{
								Select: "select",
								Fields: []parser.Expression{
									parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
								},
							},
							FromClause: parser.FromClause{
								Tables: []parser.Expression{
									parser.Table{Object: parser.Identifier{Literal: "table1"}},
								},
							},
						},
					},
				},
			},
			Operator: ">",
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
	{
		Name: "All Query Execution Error",
		Expr: parser.All{
			LHS: parser.NewInteger(5),
			Values: parser.RowValue{
				Value: parser.Subquery{
					Query: parser.SelectQuery{
						SelectEntity: parser.SelectEntity{
							SelectClause: parser.SelectClause{
								Select: "select",
								Fields: []parser.Expression{
									parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
								},
							},
							FromClause: parser.FromClause{
								Tables: []parser.Expression{
									parser.Table{Object: parser.Identifier{Literal: "table1"}},
								},
							},
							WhereClause: parser.WhereClause{
								Filter: parser.Comparison{
									LHS:      parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
									RHS:      parser.FieldReference{View: parser.Identifier{Literal: "table2"}, Column: parser.Identifier{Literal: "column4"}},
									Operator: "=",
								},
							},
						},
					},
				},
			},
			Operator: ">",
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
	{
		Name: "All Row Value Select Field Not Match Error",
		Expr: parser.All{
			LHS: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.Expression{
						parser.NewInteger(1),
						parser.NewInteger(2),
					},
				},
			},
			Values: parser.Subquery{
				Query: parser.SelectQuery{
					SelectEntity: parser.SelectEntity{
						SelectClause: parser.SelectClause{
							Select: "select",
							Fields: []parser.Expression{
								parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
							},
						},
						FromClause: parser.FromClause{
							Tables: []parser.Expression{
								parser.Table{Object: parser.Identifier{Literal: "table1"}},
							},
						},
					},
				},
			},
			Operator: ">",
		},
		Error: "[L:- C:-] select query should return exactly 2 fields",
	},
	{
		Name: "All Row Value Length Not Match Error",
		Expr: parser.All{
			LHS: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.Expression{
						parser.NewInteger(1),
						parser.NewInteger(2),
					},
				},
			},
			Values: parser.RowValueList{
				RowValues: []parser.Expression{
					parser.RowValue{
						Value: parser.ValueList{
							Values: []parser.Expression{
								parser.NewInteger(1),
								parser.NewInteger(1),
							},
						},
					},
					parser.RowValue{
						Value: parser.ValueList{
							Values: []parser.Expression{
								parser.NewInteger(1),
							},
						},
					},
				},
			},
			Operator: ">",
		},
		Error: "[L:- C:-] row value should contain exactly 2 values",
	},
	{
		Name: "Like",
		Expr: parser.Like{
			LHS:      parser.NewString("abcdefg"),
			Pattern:  parser.NewString("_bc%"),
			Negation: parser.Token{Token: parser.NOT, Literal: "not"},
		},
		Result: parser.NewTernary(ternary.FALSE),
	},
	{
		Name: "Like LHS Error",
		Expr: parser.Like{
			LHS:      parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			Pattern:  parser.NewString("_bc%"),
			Negation: parser.Token{Token: parser.NOT, Literal: "not"},
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
	{
		Name: "Like Pattern Error",
		Expr: parser.Like{
			LHS:      parser.NewString("abcdefg"),
			Pattern:  parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			Negation: parser.Token{Token: parser.NOT, Literal: "not"},
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
	{
		Name: "Exists",
		Filter: Filter{
			Records: []FilterRecord{
				{
					View: &View{
						Header: NewHeaderWithId("table2", []string{"column3", "column4"}),
						Records: []Record{
							NewRecordWithId(1, []parser.Primary{
								parser.NewInteger(1),
								parser.NewString("str2"),
							}),
						},
					},
					RecordIndex: 0,
				},
			},
		},
		Expr: parser.Exists{
			Query: parser.Subquery{
				Query: parser.SelectQuery{
					SelectEntity: parser.SelectEntity{
						SelectClause: parser.SelectClause{
							Select: "select",
							Fields: []parser.Expression{
								parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
							},
						},
						FromClause: parser.FromClause{
							Tables: []parser.Expression{
								parser.Table{Object: parser.Identifier{Literal: "table1"}},
							},
						},
						WhereClause: parser.WhereClause{
							Filter: parser.Comparison{
								LHS:      parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
								RHS:      parser.FieldReference{View: parser.Identifier{Literal: "table2"}, Column: parser.Identifier{Literal: "column4"}},
								Operator: "=",
							},
						},
					},
				},
			},
		},
		Result: parser.NewTernary(ternary.TRUE),
	},
	{
		Name: "Exists No Record",
		Expr: parser.Exists{
			Query: parser.Subquery{
				Query: parser.SelectQuery{
					SelectEntity: parser.SelectEntity{
						SelectClause: parser.SelectClause{
							Select: "select",
							Fields: []parser.Expression{
								parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
							},
						},
						FromClause: parser.FromClause{
							Tables: []parser.Expression{
								parser.Table{Object: parser.Identifier{Literal: "table1"}},
							},
						},
						WhereClause: parser.WhereClause{
							Filter: parser.NewTernary(ternary.FALSE),
						},
					},
				},
			},
		},
		Result: parser.NewTernary(ternary.FALSE),
	},
	{
		Name: "Exists Query Execution Error",
		Expr: parser.Exists{
			Query: parser.Subquery{
				Query: parser.SelectQuery{
					SelectEntity: parser.SelectEntity{
						SelectClause: parser.SelectClause{
							Select: "select",
							Fields: []parser.Expression{
								parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
							},
						},
						FromClause: parser.FromClause{
							Tables: []parser.Expression{
								parser.Table{Object: parser.Identifier{Literal: "table1"}},
							},
						},
						WhereClause: parser.WhereClause{
							Filter: parser.Comparison{
								LHS:      parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
								RHS:      parser.NewString("str2"),
								Operator: "=",
							},
						},
					},
				},
			},
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
	{
		Name: "Subquery",
		Filter: Filter{
			Records: []FilterRecord{
				{
					View: &View{
						Header: NewHeaderWithId("table2", []string{"column3", "column4"}),
						Records: []Record{
							NewRecordWithId(1, []parser.Primary{
								parser.NewInteger(1),
								parser.NewString("str2"),
							}),
						},
					},
					RecordIndex: 0,
				},
			},
		},
		Expr: parser.Subquery{
			Query: parser.SelectQuery{
				SelectEntity: parser.SelectEntity{
					SelectClause: parser.SelectClause{
						Select: "select",
						Fields: []parser.Expression{
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
						},
					},
					FromClause: parser.FromClause{
						Tables: []parser.Expression{
							parser.Table{Object: parser.Identifier{Literal: "table1"}},
						},
					},
					WhereClause: parser.WhereClause{
						Filter: parser.Comparison{
							LHS:      parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
							RHS:      parser.FieldReference{View: parser.Identifier{Literal: "table2"}, Column: parser.Identifier{Literal: "column4"}},
							Operator: "=",
						},
					},
				},
				LimitClause: parser.LimitClause{
					Value: parser.NewInteger(1),
				},
			},
		},
		Result: parser.NewString("2"),
	},
	{
		Name: "Subquery No Record",
		Expr: parser.Subquery{
			Query: parser.SelectQuery{
				SelectEntity: parser.SelectEntity{
					SelectClause: parser.SelectClause{
						Select: "select",
						Fields: []parser.Expression{
							parser.Field{Object: parser.NewInteger(1)},
						},
					},
					FromClause: parser.FromClause{
						Tables: []parser.Expression{
							parser.Table{Object: parser.Identifier{Literal: "table1"}},
						},
					},
					WhereClause: parser.WhereClause{
						Filter: parser.NewTernary(ternary.FALSE),
					},
				},
				LimitClause: parser.LimitClause{
					Value: parser.NewInteger(1),
				},
			},
		},
		Result: parser.NewNull(),
	},
	{
		Name: "Subquery Execution Error",
		Expr: parser.Subquery{
			Query: parser.SelectQuery{
				SelectEntity: parser.SelectEntity{
					SelectClause: parser.SelectClause{
						Select: "select",
						Fields: []parser.Expression{
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}}},
						},
					},
					FromClause: parser.FromClause{
						Tables: []parser.Expression{
							parser.Table{Object: parser.Identifier{Literal: "table1"}},
						},
					},
				},
			},
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
	{
		Name: "Subquery Too Many Records Error",
		Expr: parser.Subquery{
			Query: parser.SelectQuery{
				SelectEntity: parser.SelectEntity{
					SelectClause: parser.SelectClause{
						Select: "select",
						Fields: []parser.Expression{
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
						},
					},
					FromClause: parser.FromClause{
						Tables: []parser.Expression{
							parser.Table{Object: parser.Identifier{Literal: "table1"}},
						},
					},
				},
			},
		},
		Error: "[L:- C:-] subquery returns too many records, should return only one record",
	},
	{
		Name: "Subquery Too Many Fields Error",
		Expr: parser.Subquery{
			Query: parser.SelectQuery{
				SelectEntity: parser.SelectEntity{
					SelectClause: parser.SelectClause{
						Select: "select",
						Fields: []parser.Expression{
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
						},
					},
					FromClause: parser.FromClause{
						Tables: []parser.Expression{
							parser.Table{Object: parser.Identifier{Literal: "table1"}},
						},
					},
				},
				LimitClause: parser.LimitClause{
					Value: parser.NewInteger(1),
				},
			},
		},
		Error: "[L:- C:-] subquery returns too many fields, should return only one field",
	},
	{
		Name: "Function",
		Expr: parser.Function{
			Name: "coalesce",
			Args: []parser.Expression{
				parser.NewNull(),
				parser.NewString("str"),
			},
		},
		Result: parser.NewString("str"),
	},
	{
		Name: "User Defined Function",
		Filter: Filter{
			FunctionsList: UserDefinedFunctionsList{
				UserDefinedFunctionMap{
					"USERFUNC": &UserDefinedFunction{
						Name: parser.Identifier{Literal: "userfunc"},
						Parameters: []parser.Variable{
							{Name: "@arg1"},
						},
						RequiredArgs: 1,
						Statements: []parser.Statement{
							parser.Return{Value: parser.Variable{Name: "@arg1"}},
						},
					},
				},
			},
		},
		Expr: parser.Function{
			Name: "userfunc",
			Args: []parser.Expression{
				parser.NewInteger(1),
			},
		},
		Result: parser.NewInteger(1),
	},
	{
		Name: "User Defined Function Argument Length Error",
		Filter: Filter{
			FunctionsList: UserDefinedFunctionsList{
				UserDefinedFunctionMap{
					"USERFUNC": &UserDefinedFunction{
						Name: parser.Identifier{Literal: "userfunc"},
						Parameters: []parser.Variable{
							{Name: "@arg1"},
						},
						RequiredArgs: 1,
						Statements: []parser.Statement{
							parser.Return{Value: parser.Variable{Name: "@arg1"}},
						},
					},
				},
			},
		},
		Expr: parser.Function{
			Name: "userfunc",
			Args: []parser.Expression{},
		},
		Error: "[L:- C:-] function userfunc takes exactly 1 argument",
	},
	{
		Name: "Function Not Exist Error",
		Expr: parser.Function{
			Name: "notexist",
			Args: []parser.Expression{
				parser.NewNull(),
				parser.NewString("str"),
			},
		},
		Error: "[L:- C:-] function notexist does not exist",
	},
	{
		Name: "Function Evaluate Error",
		Expr: parser.Function{
			Name: "coalesce",
			Args: []parser.Expression{
				parser.NewNull(),
				parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			},
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
	{
		Name: "Aggregate Function",
		Filter: Filter{
			Records: []FilterRecord{
				{
					View: &View{
						Header: NewHeaderWithId("table1", []string{"column1", "column2"}),
						Records: []Record{
							{
								NewGroupCell([]parser.Primary{
									parser.NewInteger(1),
									parser.NewInteger(2),
									parser.NewInteger(3),
								}),
								NewGroupCell([]parser.Primary{
									parser.NewInteger(1),
									parser.NewNull(),
									parser.NewInteger(3),
								}),
								NewGroupCell([]parser.Primary{
									parser.NewString("str1"),
									parser.NewString("str2"),
									parser.NewString("str3"),
								}),
							},
						},
						isGrouped: true,
					},
					RecordIndex: 0,
				},
			},
		},
		Expr: parser.AggregateFunction{
			Name:     "avg",
			Distinct: parser.Token{},
			Args: []parser.Expression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
			},
		},
		Result: parser.NewInteger(2),
	},
	{
		Name: "Aggregate Function Argument Length Error",
		Filter: Filter{
			Records: []FilterRecord{
				{
					View: &View{
						Header: NewHeader("table1", []string{"column1", "column2"}),
						Records: []Record{
							{
								NewGroupCell([]parser.Primary{
									parser.NewInteger(1),
									parser.NewNull(),
									parser.NewInteger(3),
								}),
								NewGroupCell([]parser.Primary{
									parser.NewString("str1"),
									parser.NewString("str2"),
									parser.NewString("str3"),
								}),
							},
						},
						isGrouped: true,
					},
					RecordIndex: 0,
				},
			},
		},
		Expr: parser.AggregateFunction{
			Name:     "avg",
			Distinct: parser.Token{},
		},
		Error: "[L:- C:-] function avg takes exactly 1 argument",
	},
	{
		Name: "Aggregate Function Not Grouped Error",
		Filter: Filter{
			Records: []FilterRecord{
				{
					View: &View{
						Header: NewHeader("table1", []string{"column1", "column2"}),
						Records: []Record{
							NewRecordWithId(1, []parser.Primary{
								parser.NewInteger(1),
								parser.NewString("str2"),
							}),
						},
					},
					RecordIndex: 0,
				},
			},
		},
		Expr: parser.AggregateFunction{
			Name:     "avg",
			Distinct: parser.Token{},
			Args: []parser.Expression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
			},
		},
		Error: "[L:- C:-] function avg cannot aggregate not grouping records",
	},
	{
		Name: "Aggregate Function Nested Error",
		Filter: Filter{
			Records: []FilterRecord{
				{
					View: &View{
						Header: NewHeader("table1", []string{"column1", "column2"}),
						Records: []Record{
							{
								NewGroupCell([]parser.Primary{
									parser.NewInteger(1),
									parser.NewNull(),
									parser.NewInteger(3),
								}),
								NewGroupCell([]parser.Primary{
									parser.NewString("str1"),
									parser.NewString("str2"),
									parser.NewString("str3"),
								}),
							},
						},
						isGrouped: true,
					},
					RecordIndex: 0,
				},
			},
		},
		Expr: parser.AggregateFunction{
			Name:     "avg",
			Distinct: parser.Token{},
			Args: []parser.Expression{
				parser.AggregateFunction{
					Name: "avg",
					Args: []parser.Expression{
						parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
					},
				},
			},
		},
		Error: "[L:- C:-] aggregate functions are nested at avg(avg(column1))",
	},
	{
		Name: "Aggregate Function Count With AllColumns",
		Filter: Filter{
			Records: []FilterRecord{
				{
					View: &View{
						Header: NewHeader("table1", []string{"column1", "column2"}),
						Records: []Record{
							{
								NewGroupCell([]parser.Primary{
									parser.NewInteger(1),
									parser.NewNull(),
									parser.NewInteger(3),
								}),
								NewGroupCell([]parser.Primary{
									parser.NewString("str1"),
									parser.NewString("str2"),
									parser.NewString("str3"),
								}),
							},
						},
						isGrouped: true,
					},
					RecordIndex: 0,
				},
			},
		},
		Expr: parser.AggregateFunction{
			Name:     "count",
			Distinct: parser.Token{},
			Args: []parser.Expression{
				parser.AllColumns{},
			},
		},
		Result: parser.NewInteger(3),
	},
	{
		Name:   "Aggregate Function As a Statement Error",
		Filter: Filter{},
		Expr: parser.AggregateFunction{
			Name:     "avg",
			Distinct: parser.Token{},
			Args: []parser.Expression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
			},
		},
		Error: "[L:- C:-] function avg cannot be used as a statement",
	},
	{
		Name: "Aggregate Function User Defined",
		Filter: Filter{
			Records: []FilterRecord{
				{
					View: &View{
						Header: NewHeader("table1", []string{"column1", "column2"}),
						Records: []Record{
							{
								NewGroupCell([]parser.Primary{
									parser.NewInteger(1),
									parser.NewNull(),
									parser.NewInteger(3),
									parser.NewInteger(3),
								}),
								NewGroupCell([]parser.Primary{
									parser.NewString("str1"),
									parser.NewString("str2"),
									parser.NewString("str3"),
									parser.NewString("str4"),
								}),
							},
						},
						isGrouped: true,
					},
					RecordIndex: 0,
				},
			},
			FunctionsList: UserDefinedFunctionsList{
				{
					"USERAGGFUNC": &UserDefinedFunction{
						Name:        parser.Identifier{Literal: "useraggfunc"},
						IsAggregate: true,
						Cursor:      parser.Identifier{Literal: "column1"},
						Parameters: []parser.Variable{
							{Name: "@default"},
						},
						RequiredArgs: 1,
						Statements: []parser.Statement{
							parser.VariableDeclaration{
								Assignments: []parser.Expression{
									parser.VariableAssignment{
										Variable: parser.Variable{Name: "@value"},
									},
									parser.VariableAssignment{
										Variable: parser.Variable{Name: "@fetch"},
									},
								},
							},
							parser.WhileInCursor{
								Variables: []parser.Variable{
									{Name: "@fetch"},
								},
								Cursor: parser.Identifier{Literal: "column1"},
								Statements: []parser.Statement{
									parser.If{
										Condition: parser.Is{
											LHS: parser.Variable{Name: "@fetch"},
											RHS: parser.NewNull(),
										},
										Statements: []parser.Statement{
											parser.FlowControl{Token: parser.CONTINUE},
										},
									},
									parser.If{
										Condition: parser.Is{
											LHS: parser.Variable{Name: "@value"},
											RHS: parser.NewNull(),
										},
										Statements: []parser.Statement{
											parser.VariableSubstitution{
												Variable: parser.Variable{Name: "@value"},
												Value:    parser.Variable{Name: "@fetch"},
											},
											parser.FlowControl{Token: parser.CONTINUE},
										},
									},
									parser.VariableSubstitution{
										Variable: parser.Variable{Name: "@value"},
										Value: parser.Arithmetic{
											LHS:      parser.Variable{Name: "@value"},
											RHS:      parser.Variable{Name: "@fetch"},
											Operator: '*',
										},
									},
								},
							},

							parser.If{
								Condition: parser.Is{
									LHS: parser.Variable{Name: "@value"},
									RHS: parser.NewNull(),
								},
								Statements: []parser.Statement{
									parser.VariableSubstitution{
										Variable: parser.Variable{Name: "@value"},
										Value:    parser.Variable{Name: "@default"},
									},
								},
							},

							parser.Return{
								Value: parser.Variable{Name: "@value"},
							},
						},
					},
				},
			},
		},
		Expr: parser.AggregateFunction{
			Name:     "useraggfunc",
			Distinct: parser.Token{Token: parser.DISTINCT, Literal: "distinct"},
			Args: []parser.Expression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
				parser.NewInteger(0),
			},
		},
		Result: parser.NewInteger(3),
	},
	{
		Name: "Aggregate Function User Defined Argument Length Error",
		Filter: Filter{
			Records: []FilterRecord{
				{
					View: &View{
						Header: NewHeader("table1", []string{"column1", "column2"}),
						Records: []Record{
							{
								NewGroupCell([]parser.Primary{
									parser.NewInteger(1),
									parser.NewNull(),
									parser.NewInteger(3),
									parser.NewInteger(3),
								}),
								NewGroupCell([]parser.Primary{
									parser.NewString("str1"),
									parser.NewString("str2"),
									parser.NewString("str3"),
									parser.NewString("str4"),
								}),
							},
						},
						isGrouped: true,
					},
					RecordIndex: 0,
				},
			},
			FunctionsList: UserDefinedFunctionsList{
				{
					"USERAGGFUNC": &UserDefinedFunction{
						Name:        parser.Identifier{Literal: "useraggfunc"},
						IsAggregate: true,
						Cursor:      parser.Identifier{Literal: "column1"},
						Parameters: []parser.Variable{
							{Name: "@default"},
						},
						Statements: []parser.Statement{
							parser.Return{
								Value: parser.Variable{Name: "@value"},
							},
						},
					},
				},
			},
		},
		Expr: parser.AggregateFunction{
			Name:     "useraggfunc",
			Distinct: parser.Token{Token: parser.DISTINCT, Literal: "distinct"},
			Args: []parser.Expression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
			},
		},
		Error: "[L:- C:-] function useraggfunc takes exactly 2 arguments",
	},
	{
		Name: "Aggregate Function User Defined Argument Evaluation Error",
		Filter: Filter{
			Records: []FilterRecord{
				{
					View: &View{
						Header: NewHeader("table1", []string{"column1", "column2"}),
						Records: []Record{
							{
								NewGroupCell([]parser.Primary{
									parser.NewInteger(1),
									parser.NewNull(),
									parser.NewInteger(3),
									parser.NewInteger(3),
								}),
								NewGroupCell([]parser.Primary{
									parser.NewString("str1"),
									parser.NewString("str2"),
									parser.NewString("str3"),
									parser.NewString("str4"),
								}),
							},
						},
						isGrouped: true,
					},
					RecordIndex: 0,
				},
			},
			FunctionsList: UserDefinedFunctionsList{
				{
					"USERAGGFUNC": &UserDefinedFunction{
						Name:        parser.Identifier{Literal: "useraggfunc"},
						IsAggregate: true,
						Cursor:      parser.Identifier{Literal: "column1"},
						Parameters: []parser.Variable{
							{Name: "@default"},
						},
						Statements: []parser.Statement{
							parser.Return{
								Value: parser.Variable{Name: "@value"},
							},
						},
					},
				},
			},
		},
		Expr: parser.AggregateFunction{
			Name:     "useraggfunc",
			Distinct: parser.Token{Token: parser.DISTINCT, Literal: "distinct"},
			Args: []parser.Expression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
				parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			},
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
	{
		Name: "Aggregate Function Execute User Defined Aggregate Function Passed As Scala Function",
		Filter: Filter{
			Records: []FilterRecord{
				{
					View: &View{
						Header: NewHeader("table1", []string{"column1", "column2"}),
						Records: []Record{
							{
								NewGroupCell([]parser.Primary{
									parser.NewInteger(1),
									parser.NewNull(),
									parser.NewInteger(3),
								}),
								NewGroupCell([]parser.Primary{
									parser.NewString("str1"),
									parser.NewString("str2"),
									parser.NewString("str3"),
								}),
							},
						},
						isGrouped: true,
					},
					RecordIndex: 0,
				},
			},
			FunctionsList: UserDefinedFunctionsList{
				{
					"USERAGGFUNC": &UserDefinedFunction{
						Name:        parser.Identifier{Literal: "useraggfunc"},
						IsAggregate: true,
						Cursor:      parser.Identifier{Literal: "column1"},
						Statements: []parser.Statement{
							parser.VariableDeclaration{
								Assignments: []parser.Expression{
									parser.VariableAssignment{
										Variable: parser.Variable{Name: "@value"},
									},
									parser.VariableAssignment{
										Variable: parser.Variable{Name: "@fetch"},
									},
								},
							},
							parser.WhileInCursor{
								Variables: []parser.Variable{
									{Name: "@fetch"},
								},
								Cursor: parser.Identifier{Literal: "column1"},
								Statements: []parser.Statement{
									parser.If{
										Condition: parser.Is{
											LHS: parser.Variable{Name: "@fetch"},
											RHS: parser.NewNull(),
										},
										Statements: []parser.Statement{
											parser.FlowControl{Token: parser.CONTINUE},
										},
									},
									parser.If{
										Condition: parser.Is{
											LHS: parser.Variable{Name: "@value"},
											RHS: parser.NewNull(),
										},
										Statements: []parser.Statement{
											parser.VariableSubstitution{
												Variable: parser.Variable{Name: "@value"},
												Value:    parser.Variable{Name: "@fetch"},
											},
											parser.FlowControl{Token: parser.CONTINUE},
										},
									},
									parser.VariableSubstitution{
										Variable: parser.Variable{Name: "@value"},
										Value: parser.Arithmetic{
											LHS:      parser.Variable{Name: "@value"},
											RHS:      parser.Variable{Name: "@fetch"},
											Operator: '*',
										},
									},
								},
							},
							parser.Return{
								Value: parser.Variable{Name: "@value"},
							},
						},
					},
				},
			},
		},
		Expr: parser.Function{
			Name: "useraggfunc",
			Args: []parser.Expression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
			},
		},
		Result: parser.NewInteger(3),
	},
	{
		Name: "Aggregate Function Execute User Defined Aggregate Function Undefined Error",
		Filter: Filter{
			Records: []FilterRecord{
				{
					View: &View{
						Header: NewHeader("table1", []string{"column1", "column2"}),
						Records: []Record{
							{
								NewGroupCell([]parser.Primary{
									parser.NewInteger(1),
									parser.NewNull(),
									parser.NewInteger(3),
								}),
								NewGroupCell([]parser.Primary{
									parser.NewString("str1"),
									parser.NewString("str2"),
									parser.NewString("str3"),
								}),
							},
						},
						isGrouped: true,
					},
					RecordIndex: 0,
				},
			},
			FunctionsList: UserDefinedFunctionsList{
				{
					"USERAGGFUNC": &UserDefinedFunction{
						Name:        parser.Identifier{Literal: "useraggfunc"},
						IsAggregate: true,
						Cursor:      parser.Identifier{Literal: "column1"},
						Statements: []parser.Statement{
							parser.VariableDeclaration{
								Assignments: []parser.Expression{
									parser.VariableAssignment{
										Variable: parser.Variable{Name: "@value"},
									},
									parser.VariableAssignment{
										Variable: parser.Variable{Name: "@fetch"},
									},
								},
							},
							parser.WhileInCursor{
								Variables: []parser.Variable{
									{Name: "@fetch"},
								},
								Cursor: parser.Identifier{Literal: "column1"},
								Statements: []parser.Statement{
									parser.If{
										Condition: parser.Is{
											LHS: parser.Variable{Name: "@fetch"},
											RHS: parser.NewNull(),
										},
										Statements: []parser.Statement{
											parser.FlowControl{Token: parser.CONTINUE},
										},
									},
									parser.If{
										Condition: parser.Is{
											LHS: parser.Variable{Name: "@value"},
											RHS: parser.NewNull(),
										},
										Statements: []parser.Statement{
											parser.VariableSubstitution{
												Variable: parser.Variable{Name: "@value"},
												Value:    parser.Variable{Name: "@fetch"},
											},
											parser.FlowControl{Token: parser.CONTINUE},
										},
									},
									parser.VariableSubstitution{
										Variable: parser.Variable{Name: "@value"},
										Value: parser.Arithmetic{
											LHS:      parser.Variable{Name: "@value"},
											RHS:      parser.Variable{Name: "@fetch"},
											Operator: '*',
										},
									},
								},
							},
							parser.Return{
								Value: parser.Variable{Name: "@value"},
							},
						},
					},
				},
			},
		},
		Expr: parser.AggregateFunction{
			Name: "undefined",
			Args: []parser.Expression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
			},
		},
		Error: "[L:- C:-] function undefined does not exist",
	},
	{
		Name: "ListAgg Function",
		Filter: Filter{
			Records: []FilterRecord{
				{
					View: &View{
						Header: NewHeaderWithId("table1", []string{"column1", "column2"}),
						Records: []Record{
							{
								NewGroupCell([]parser.Primary{
									parser.NewInteger(1),
									parser.NewInteger(2),
									parser.NewInteger(3),
									parser.NewInteger(4),
								}),
								NewGroupCell([]parser.Primary{
									parser.NewInteger(1),
									parser.NewInteger(2),
									parser.NewInteger(3),
									parser.NewInteger(4),
								}),
								NewGroupCell([]parser.Primary{
									parser.NewString("str2"),
									parser.NewString("str1"),
									parser.NewNull(),
									parser.NewString("str2"),
								}),
							},
						},
						isGrouped: true,
					},
					RecordIndex: 0,
				},
			},
		},
		Expr: parser.ListAgg{
			ListAgg:  "listagg",
			Distinct: parser.Token{Token: parser.DISTINCT, Literal: "distinct"},
			Args: []parser.Expression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
				parser.NewString(","),
			},
			OrderBy: parser.OrderByClause{
				Items: []parser.Expression{
					parser.OrderItem{Value: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
				},
			},
		},
		Result: parser.NewString("str1,str2"),
	},
	{
		Name: "ListAgg Function Null",
		Filter: Filter{
			Records: []FilterRecord{
				{
					View: &View{
						Header: NewHeaderWithId("table1", []string{"column1", "column2"}),
						Records: []Record{
							{
								NewGroupCell([]parser.Primary{
									parser.NewInteger(1),
									parser.NewInteger(2),
									parser.NewInteger(3),
									parser.NewInteger(4),
								}),
								NewGroupCell([]parser.Primary{
									parser.NewInteger(1),
									parser.NewInteger(2),
									parser.NewInteger(3),
									parser.NewInteger(4),
								}),
								NewGroupCell([]parser.Primary{
									parser.NewNull(),
									parser.NewNull(),
									parser.NewNull(),
									parser.NewNull(),
								}),
							},
						},
						isGrouped: true,
					},
					RecordIndex: 0,
				},
			},
		},
		Expr: parser.ListAgg{
			Distinct: parser.Token{Token: parser.DISTINCT, Literal: "distinct"},
			Args: []parser.Expression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
				parser.NewString(","),
			},
		},
		Result: parser.NewNull(),
	},
	{
		Name: "ListAgg Function Argument Length Error",
		Filter: Filter{
			Records: []FilterRecord{
				{
					View: &View{
						Header: NewHeaderWithId("table1", []string{"column1", "column2"}),
						Records: []Record{
							NewRecordWithId(1, []parser.Primary{
								parser.NewInteger(1),
								parser.NewString("str2"),
							}),
						},
					},
					RecordIndex: 0,
				},
			},
		},
		Expr: parser.ListAgg{
			ListAgg:  "listagg",
			Distinct: parser.Token{Token: parser.DISTINCT, Literal: "distinct"},
			OrderBy: parser.OrderByClause{
				Items: []parser.Expression{
					parser.OrderItem{Value: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
				},
			},
		},
		Error: "[L:- C:-] function listagg takes 1 or 2 arguments",
	},
	{
		Name: "ListAgg Function Not Grouped Error",
		Filter: Filter{
			Records: []FilterRecord{
				{
					View: &View{
						Header: NewHeaderWithId("table1", []string{"column1", "column2"}),
						Records: []Record{
							NewRecordWithId(1, []parser.Primary{
								parser.NewInteger(1),
								parser.NewString("str2"),
							}),
						},
					},
					RecordIndex: 0,
				},
			},
		},
		Expr: parser.ListAgg{
			ListAgg:  "listagg",
			Distinct: parser.Token{Token: parser.DISTINCT, Literal: "distinct"},
			Args: []parser.Expression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
				parser.NewString(","),
			},
			OrderBy: parser.OrderByClause{
				Items: []parser.Expression{
					parser.OrderItem{Value: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
				},
			},
		},
		Error: "[L:- C:-] function listagg cannot aggregate not grouping records",
	},
	{
		Name: "ListAgg Function Sort Error",
		Filter: Filter{
			Records: []FilterRecord{
				{
					View: &View{
						Header: NewHeaderWithId("table1", []string{"column1", "column2"}),
						Records: []Record{
							{
								NewGroupCell([]parser.Primary{
									parser.NewInteger(1),
									parser.NewInteger(2),
									parser.NewInteger(3),
									parser.NewInteger(4),
								}),
								NewGroupCell([]parser.Primary{
									parser.NewString("str2"),
									parser.NewString("str1"),
									parser.NewNull(),
									parser.NewString("str2"),
								}),
							},
						},
						isGrouped: true,
					},
					RecordIndex: 0,
				},
			},
		},
		Expr: parser.ListAgg{
			ListAgg:  "listagg",
			Distinct: parser.Token{Token: parser.DISTINCT, Literal: "distinct"},
			Args: []parser.Expression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
				parser.NewString(","),
			},
			OrderBy: parser.OrderByClause{
				Items: []parser.Expression{
					parser.OrderItem{Value: parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}}},
				},
			},
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
	{
		Name: "ListAgg Function Nested Error",
		Filter: Filter{
			Records: []FilterRecord{
				{
					View: &View{
						Header: NewHeaderWithId("table1", []string{"column1", "column2"}),
						Records: []Record{
							{
								NewGroupCell([]parser.Primary{
									parser.NewInteger(1),
									parser.NewNull(),
									parser.NewInteger(3),
								}),
								NewGroupCell([]parser.Primary{
									parser.NewString("str1"),
									parser.NewString("str2"),
									parser.NewString("str3"),
								}),
							},
						},
						isGrouped: true,
					},
					RecordIndex: 0,
				},
			},
		},
		Expr: parser.ListAgg{
			ListAgg: "listagg",
			Args: []parser.Expression{
				parser.AggregateFunction{
					Name:     "avg",
					Distinct: parser.Token{},
					Args: []parser.Expression{
						parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
					},
				},
			},
		},
		Error: "[L:- C:-] aggregate functions are nested at listagg(avg(column1))",
	},
	{
		Name: "ListAgg Function Second Argument Evaluation Error",
		Filter: Filter{
			Records: []FilterRecord{
				{
					View: &View{
						Header: NewHeaderWithId("table1", []string{"column1", "column2"}),
						Records: []Record{
							{
								NewGroupCell([]parser.Primary{
									parser.NewInteger(1),
									parser.NewNull(),
									parser.NewInteger(3),
								}),
								NewGroupCell([]parser.Primary{
									parser.NewString("str1"),
									parser.NewString("str2"),
									parser.NewString("str3"),
								}),
							},
						},
						isGrouped: true,
					},
					RecordIndex: 0,
				},
			},
		},
		Expr: parser.ListAgg{
			ListAgg: "listagg",
			Args: []parser.Expression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
			},
		},
		Error: "[L:- C:-] the second argument must be a string for function listagg",
	},
	{
		Name: "ListAgg Function Second Argument Not String Error",
		Filter: Filter{
			Records: []FilterRecord{
				{
					View: &View{
						Header: NewHeaderWithId("table1", []string{"column1", "column2"}),
						Records: []Record{
							{
								NewGroupCell([]parser.Primary{
									parser.NewInteger(1),
									parser.NewNull(),
									parser.NewInteger(3),
								}),
								NewGroupCell([]parser.Primary{
									parser.NewString("str1"),
									parser.NewString("str2"),
									parser.NewString("str3"),
								}),
							},
						},
						isGrouped: true,
					},
					RecordIndex: 0,
				},
			},
		},
		Expr: parser.ListAgg{
			ListAgg: "listagg",
			Args: []parser.Expression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
				parser.NewNull(),
			},
		},
		Error: "[L:- C:-] the second argument must be a string for function listagg",
	},
	{
		Name:   "ListAgg Function As a Statement Error",
		Filter: Filter{},
		Expr: parser.ListAgg{
			ListAgg:  "listagg",
			Distinct: parser.Token{Token: parser.DISTINCT, Literal: "distinct"},
			Args: []parser.Expression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
				parser.NewString(","),
			},
			OrderBy: parser.OrderByClause{
				Items: []parser.Expression{
					parser.OrderItem{Value: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
				},
			},
		},
		Error: "[L:- C:-] function listagg cannot be used as a statement",
	},
	{
		Name: "Case Comparison",
		Expr: parser.Case{
			Value: parser.NewInteger(2),
			When: []parser.Expression{
				parser.CaseWhen{
					Condition: parser.NewInteger(1),
					Result:    parser.NewString("A"),
				},
				parser.CaseWhen{
					Condition: parser.NewInteger(2),
					Result:    parser.NewString("B"),
				},
			},
		},
		Result: parser.NewString("B"),
	},
	{
		Name: "Case Filter",
		Expr: parser.Case{
			When: []parser.Expression{
				parser.CaseWhen{
					Condition: parser.Comparison{
						LHS:      parser.NewInteger(2),
						RHS:      parser.NewInteger(1),
						Operator: "=",
					},
					Result: parser.NewString("A"),
				},
				parser.CaseWhen{
					Condition: parser.Comparison{
						LHS:      parser.NewInteger(2),
						RHS:      parser.NewInteger(2),
						Operator: "=",
					},
					Result: parser.NewString("B"),
				},
			},
		},
		Result: parser.NewString("B"),
	},
	{
		Name: "Case Else",
		Expr: parser.Case{
			Value: parser.NewInteger(0),
			When: []parser.Expression{
				parser.CaseWhen{
					Condition: parser.NewInteger(1),
					Result:    parser.NewString("A"),
				},
				parser.CaseWhen{
					Condition: parser.NewInteger(2),
					Result:    parser.NewString("B"),
				},
			},
			Else: parser.CaseElse{
				Result: parser.NewString("C"),
			},
		},
		Result: parser.NewString("C"),
	},
	{
		Name: "Case No Else",
		Expr: parser.Case{
			Value: parser.NewInteger(0),
			When: []parser.Expression{
				parser.CaseWhen{
					Condition: parser.NewInteger(1),
					Result:    parser.NewString("A"),
				},
				parser.CaseWhen{
					Condition: parser.NewInteger(2),
					Result:    parser.NewString("B"),
				},
			},
		},
		Result: parser.NewNull(),
	},
	{
		Name: "Case Value Error",
		Expr: parser.Case{
			Value: parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			When: []parser.Expression{
				parser.CaseWhen{
					Condition: parser.NewInteger(1),
					Result:    parser.NewString("A"),
				},
				parser.CaseWhen{
					Condition: parser.NewInteger(2),
					Result:    parser.NewString("B"),
				},
			},
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
	{
		Name: "Case When Condition Error",
		Expr: parser.Case{
			Value: parser.NewInteger(2),
			When: []parser.Expression{
				parser.CaseWhen{
					Condition: parser.NewInteger(1),
					Result:    parser.NewString("A"),
				},
				parser.CaseWhen{
					Condition: parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
					Result:    parser.NewString("B"),
				},
			},
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
	{
		Name: "Case When Result Error",
		Expr: parser.Case{
			Value: parser.NewInteger(2),
			When: []parser.Expression{
				parser.CaseWhen{
					Condition: parser.NewInteger(1),
					Result:    parser.NewString("A"),
				},
				parser.CaseWhen{
					Condition: parser.NewInteger(2),
					Result:    parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
				},
			},
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
	{
		Name: "Case Else Result Error",
		Expr: parser.Case{
			Value: parser.NewInteger(0),
			When: []parser.Expression{
				parser.CaseWhen{
					Condition: parser.NewInteger(1),
					Result:    parser.NewString("A"),
				},
				parser.CaseWhen{
					Condition: parser.NewInteger(2),
					Result:    parser.NewString("B"),
				},
			},
			Else: parser.CaseElse{
				Result: parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			},
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
	{
		Name: "Logic AND",
		Expr: parser.Logic{
			LHS:      parser.NewTernary(ternary.TRUE),
			RHS:      parser.NewTernary(ternary.FALSE),
			Operator: parser.Token{Token: parser.AND, Literal: "and"},
		},
		Result: parser.NewTernary(ternary.FALSE),
	},
	{
		Name: "Logic OR",
		Expr: parser.Logic{
			LHS:      parser.NewTernary(ternary.TRUE),
			RHS:      parser.NewTernary(ternary.FALSE),
			Operator: parser.Token{Token: parser.OR, Literal: "or"},
		},
		Result: parser.NewTernary(ternary.TRUE),
	},
	{
		Name: "Logic LHS Error",
		Expr: parser.Logic{
			LHS:      parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			RHS:      parser.NewTernary(ternary.FALSE),
			Operator: parser.Token{Token: parser.AND, Literal: "and"},
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
	{
		Name: "Logic RHS Error",
		Expr: parser.Logic{
			LHS:      parser.NewTernary(ternary.FALSE),
			RHS:      parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			Operator: parser.Token{Token: parser.AND, Literal: "and"},
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
	{
		Name: "UnaryLogic",
		Expr: parser.UnaryLogic{
			Operand:  parser.NewTernary(ternary.FALSE),
			Operator: parser.Token{Token: parser.NOT, Literal: "not"},
		},
		Result: parser.NewTernary(ternary.TRUE),
	},
	{
		Name: "UnaryLogic Operand Error",
		Expr: parser.UnaryLogic{
			Operand:  parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			Operator: parser.Token{Token: parser.NOT, Literal: "not"},
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
	{
		Name: "Variable",
		Filter: NewFilter(
			[]Variables{{
				"@var1": parser.NewInteger(1),
			}},
			[]ViewMap{{}},
			[]CursorMap{{}},
			[]UserDefinedFunctionMap{{}},
		),
		Expr: parser.Variable{
			Name: "@var1",
		},
		Result: parser.NewInteger(1),
	},
	{
		Name:   "Variable Undefined Error",
		Filter: NewEmptyFilter(),
		Expr: parser.Variable{
			Name: "@undefined",
		},
		Error: "[L:- C:-] variable @undefined is undefined",
	},
	{
		Name: "Variable Substitution",
		Filter: NewFilter(
			[]Variables{{
				"@var1": parser.NewInteger(1),
			}},
			[]ViewMap{{}},
			[]CursorMap{{}},
			[]UserDefinedFunctionMap{{}},
		),
		Expr: parser.VariableSubstitution{
			Variable: parser.Variable{Name: "@var1"},
			Value:    parser.NewInteger(2),
		},
		Result: parser.NewInteger(2),
	},
	{
		Name:   "Variable Substitution Undefined Error",
		Filter: NewEmptyFilter(),
		Expr: parser.VariableSubstitution{
			Variable: parser.Variable{Name: "@undefined"},
			Value:    parser.NewInteger(2),
		},
		Error: "[L:- C:-] variable @undefined is undefined",
	},
	{
		Name: "Cursor Status Is Not Open",
		Expr: parser.CursorStatus{
			CursorLit: "cursor",
			Cursor:    parser.Identifier{Literal: "cur"},
			Is:        "is",
			Negation:  parser.Token{Token: parser.NOT, Literal: "not"},
			Type:      parser.OPEN,
			TypeLit:   "open",
		},
		Result: parser.NewTernary(ternary.FALSE),
	},
	{
		Name: "Cursor Status Is In Range",
		Expr: parser.CursorStatus{
			CursorLit: "cursor",
			Cursor:    parser.Identifier{Literal: "cur"},
			Is:        "is",
			Type:      parser.RANGE,
			TypeLit:   "in range",
		},
		Result: parser.NewTernary(ternary.TRUE),
	},
	{
		Name: "Cursor Status Open Error",
		Expr: parser.CursorStatus{
			CursorLit: "cursor",
			Cursor:    parser.Identifier{Literal: "notexist"},
			Is:        "is",
			Type:      parser.OPEN,
			TypeLit:   "open",
		},
		Error: "[L:- C:-] cursor notexist is undefined",
	},
	{
		Name: "Cursor Status In Range Error",
		Expr: parser.CursorStatus{
			CursorLit: "cursor",
			Cursor:    parser.Identifier{Literal: "notexist"},
			Is:        "is",
			Type:      parser.RANGE,
			TypeLit:   "in range",
		},
		Error: "[L:- C:-] cursor notexist is undefined",
	},
	{
		Name: "Cursor Attribute Count",
		Expr: parser.CursorAttrebute{
			CursorLit: "cursor",
			Cursor:    parser.Identifier{Literal: "cur"},
			Attrebute: parser.Token{Token: parser.COUNT, Literal: "count"},
		},
		Result: parser.NewInteger(3),
	},
	{
		Name: "Cursor Attribute Count Error",
		Expr: parser.CursorAttrebute{
			CursorLit: "cursor",
			Cursor:    parser.Identifier{Literal: "notexist"},
			Attrebute: parser.Token{Token: parser.COUNT, Literal: "count"},
		},
		Error: "[L:- C:-] cursor notexist is undefined",
	},
}

func TestFilter_Evaluate(t *testing.T) {
	initFlag()
	tf := cmd.GetFlags()
	tf.Repository = TestDataDir

	cursors := CursorMap{
		"CUR": &Cursor{
			query: selectQueryForCursorTest,
		},
	}
	ViewCache.Clear()
	cursors.Open(parser.Identifier{Literal: "cur"}, NewEmptyFilter())
	cursors.Fetch(parser.Identifier{Literal: "cur"}, parser.NEXT, 0)

	for _, v := range filterEvaluateTests {
		ViewCache.Clear()
		v.Filter.CursorsList = append(v.Filter.CursorsList, cursors)
		result, err := v.Filter.Evaluate(v.Expr)
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

func BenchmarkFilter_Evaluate1(b *testing.B) {
	filter := GenerateBenchGroupedViewFilter()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		filter.Evaluate(parser.AggregateFunction{
			Name:     "count",
			Distinct: parser.Token{},
			Args: []parser.Expression{
				parser.AllColumns{},
			},
		})
	}
}

func BenchmarkFilter_Evaluate2(b *testing.B) {
	filter := GenerateBenchGroupedViewFilter()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		filter.Evaluate(parser.AggregateFunction{
			Name:     "count",
			Distinct: parser.Token{},
			Args: []parser.Expression{
				parser.FieldReference{Column: parser.Identifier{Literal: "c1"}},
			},
		})
	}
}
