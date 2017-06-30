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
		Name: "Parentheses",
		Expr: parser.Parentheses{
			Expr: parser.NewString("str"),
		},
		Result: parser.NewString("str"),
	},
	{
		Name: "FieldReference",
		Filter: []FilterRecord{
			{
				View: &View{
					Header: NewHeader("table1", []string{"column1", "column2"}),
					Records: []Record{
						NewRecord(1, []parser.Primary{
							parser.NewInteger(1),
							parser.NewString("str"),
						}),
						NewRecord(2, []parser.Primary{
							parser.NewInteger(2),
							parser.NewString("strstr"),
						}),
					},
				},
				RecordIndex: 1,
			},
		},
		Expr:   parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
		Result: parser.NewString("strstr"),
	},
	{
		Name: "FieldReference ColumnNotExist Error",
		Filter: []FilterRecord{
			{
				View: &View{
					Header: NewHeader("table1", []string{"column1", "column2"}),
					Records: []Record{
						NewRecord(1, []parser.Primary{
							parser.NewInteger(1),
							parser.NewString("str"),
						}),
						NewRecord(2, []parser.Primary{
							parser.NewInteger(2),
							parser.NewString("strstr"),
						}),
					},
				},
				RecordIndex: 1,
			},
		},
		Expr:  parser.FieldReference{Column: parser.Identifier{Literal: "column3"}},
		Error: "field column3 does not exist",
	},
	{
		Name: "FieldReference FieldAmbigous Error",
		Filter: []FilterRecord{
			{
				View: &View{
					Header: NewHeader("table1", []string{"column1", "column1"}),
					Records: []Record{
						NewRecord(1, []parser.Primary{
							parser.NewInteger(1),
							parser.NewString("str"),
						}),
						NewRecord(2, []parser.Primary{
							parser.NewInteger(2),
							parser.NewString("strstr"),
						}),
					},
				},
				RecordIndex: 1,
			},
		},
		Expr:  parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
		Error: "field column1 is ambiguous",
	},
	{
		Name: "FieldReference Not Group Key Error",
		Filter: []FilterRecord{
			{
				View: &View{
					Header: []HeaderField{
						{
							Reference: "table1",
							Column:    "column1",
						},
						{
							Reference: "table1",
							Column:    "column2",
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
		Expr:  parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
		Error: "field column1 is not a group key",
	},
	{
		Name: "FieldReference Fields Ambiguous Error with Multiple Tables",
		Filter: []FilterRecord{
			{
				View: &View{
					Header: NewHeader("table1", []string{"column1", "column2"}),
					Records: []Record{
						NewRecord(1, []parser.Primary{
							parser.NewInteger(1),
							parser.NewString("str"),
						}),
						NewRecord(2, []parser.Primary{
							parser.NewInteger(2),
							parser.NewString("strstr"),
						}),
					},
				},
				RecordIndex: 1,
			},
			{
				View: &View{
					Header: NewHeader("table2", []string{"column1", "column2"}),
					Records: []Record{
						NewRecord(1, []parser.Primary{
							parser.NewInteger(1),
							parser.NewString("str"),
						}),
						NewRecord(2, []parser.Primary{
							parser.NewInteger(2),
							parser.NewString("strstr"),
						}),
					},
				},
				RecordIndex: 1,
			},
		},
		Expr:  parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
		Error: "field column1 is ambiguous",
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
		Error: "field notexist does not exist",
	},
	{
		Name: "Arithmetic RHS Error",
		Expr: parser.Arithmetic{
			LHS:      parser.NewInteger(1),
			RHS:      parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			Operator: '+',
		},
		Error: "field notexist does not exist",
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
		Error: "field notexist does not exist",
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
			Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: "="},
		},
		Result: parser.NewTernary(ternary.FALSE),
	},
	{
		Name: "Comparison LHS Error",
		Expr: parser.Comparison{
			LHS:      parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			RHS:      parser.NewInteger(2),
			Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: "="},
		},
		Error: "field notexist does not exist",
	},
	{
		Name: "Comparison RHS Error",
		Expr: parser.Comparison{
			LHS:      parser.NewInteger(1),
			RHS:      parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: "="},
		},
		Error: "field notexist does not exist",
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
			Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: "="},
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
									Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: "="},
								},
							},
						},
					},
				},
			},
			Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: "="},
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
			Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: "="},
		},
		Error: "field notexist does not exist",
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
			Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: "="},
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
			Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: "="},
		},
		Error: "field notexist does not exist",
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
			Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: "="},
		},
		Error: "subquery returns too many records, should be only one record",
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
			Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: "="},
		},
		Error: "row value length does not match",
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
		Error: "field notexist does not exist",
	},
	{
		Name: "Is RHS Error",
		Expr: parser.Is{
			LHS:      parser.NewInteger(1),
			RHS:      parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			Negation: parser.Token{Token: parser.NOT, Literal: "not"},
		},
		Error: "field notexist does not exist",
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
		Error: "field notexist does not exist",
	},
	{
		Name: "Between Low Error",
		Expr: parser.Between{
			LHS:      parser.NewInteger(2),
			Low:      parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			High:     parser.NewInteger(3),
			Negation: parser.Token{Token: parser.NOT, Literal: "not"},
		},
		Error: "field notexist does not exist",
	},
	{
		Name: "Between High Error",
		Expr: parser.Between{
			LHS:      parser.NewInteger(2),
			Low:      parser.NewInteger(1),
			High:     parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			Negation: parser.Token{Token: parser.NOT, Literal: "not"},
		},
		Error: "field notexist does not exist",
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
		Error: "field notexist does not exist",
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
		Error: "field notexist does not exist",
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
		Error: "field notexist does not exist",
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
		Error: "row value length does not match",
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
		Error: "row value length does not match",
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
		Error: "field notexist does not exist",
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
		Error: "field notexist does not exist",
	},
	{
		Name: "In Subquery",
		Filter: []FilterRecord{
			{
				View: &View{
					Header: NewHeader("table2", []string{"column3", "column4"}),
					Records: []Record{
						NewRecord(1, []parser.Primary{
							parser.NewInteger(1),
							parser.NewString("str2"),
						}),
					},
				},
				RecordIndex: 0,
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
									Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: "="},
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
									parser.Field{Object: parser.Identifier{Literal: "column1"}},
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
									Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: "="},
								},
							},
						},
					},
				},
			},
			Negation: parser.Token{Token: parser.NOT, Literal: "not"},
		},
		Error: "field notexist does not exist",
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
		Error: "field notexist does not exist",
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
		Error: "field notexist does not exist",
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
		Error: "row value length does not match",
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
			Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: "<>"},
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
			Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: "<>"},
		},
		Error: "field notexist does not exist",
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
									Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: "="},
								},
							},
						},
					},
				},
			},
			Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: "<>"},
		},
		Error: "field notexist does not exist",
	},
	{
		Name: "Any Row Value Length Not Match Error",
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
			Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: "<>"},
		},
		Error: "row value length does not match",
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
			Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: ">"},
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
			Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: ">"},
		},
		Error: "field notexist does not exist",
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
									Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: "="},
								},
							},
						},
					},
				},
			},
			Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: ">"},
		},
		Error: "field notexist does not exist",
	},
	{
		Name: "All Row Value Length Not Match Error",
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
			Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: ">"},
		},
		Error: "row value length does not match",
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
		Error: "field notexist does not exist",
	},
	{
		Name: "Like Pattern Error",
		Expr: parser.Like{
			LHS:      parser.NewString("abcdefg"),
			Pattern:  parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			Negation: parser.Token{Token: parser.NOT, Literal: "not"},
		},
		Error: "field notexist does not exist",
	},
	{
		Name: "Exists",
		Filter: []FilterRecord{
			{
				View: &View{
					Header: NewHeader("table2", []string{"column3", "column4"}),
					Records: []Record{
						NewRecord(1, []parser.Primary{
							parser.NewInteger(1),
							parser.NewString("str2"),
						}),
					},
				},
				RecordIndex: 0,
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
								Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: "="},
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
								Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: "="},
							},
						},
					},
				},
			},
		},
		Error: "field notexist does not exist",
	},
	{
		Name: "Subquery",
		Filter: []FilterRecord{
			{
				View: &View{
					Header: NewHeader("table2", []string{"column3", "column4"}),
					Records: []Record{
						NewRecord(1, []parser.Primary{
							parser.NewInteger(1),
							parser.NewString("str2"),
						}),
					},
				},
				RecordIndex: 0,
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
							Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: "="},
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
		Error: "field notexist does not exist",
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
		Error: "subquery returns too many records, should be only one record",
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
		Error: "subquery returns too many fields, should be only one field",
	},
	{
		Name: "Function",
		Expr: parser.Function{
			Name: "coalesce",
			Option: parser.Option{
				Args: []parser.Expression{
					parser.NewNull(),
					parser.NewString("str"),
				},
			},
		},
		Result: parser.NewString("str"),
	},
	{
		Name: "Function Is Not Exist",
		Expr: parser.Function{
			Name: "notexist",
			Option: parser.Option{
				Args: []parser.Expression{
					parser.NewNull(),
					parser.NewString("str"),
				},
			},
		},
		Error: "function notexist is not exist",
	},
	{
		Name: "Function Option Error",
		Expr: parser.Function{
			Name: "coalesce",
			Option: parser.Option{
				Distinct: parser.Token{Token: parser.DISTINCT, Literal: "distinct"},
				Args: []parser.Expression{
					parser.NewNull(),
					parser.NewString("str"),
				},
			},
		},
		Error: "syntax error: unexpected distinct",
	},
	{
		Name: "Function Argument Error",
		Expr: parser.Function{
			Name: "coalesce",
			Option: parser.Option{
				Args: []parser.Expression{
					parser.AllColumns{},
				},
			},
		},
		Error: "syntax error: unexpected *",
	},
	{
		Name: "Aggregate Function",
		Filter: []FilterRecord{
			{
				View: &View{
					Header: NewHeader("table1", []string{"column1", "column2"}),
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
		Expr: parser.Function{
			Name: "avg",
			Option: parser.Option{
				Args: []parser.Expression{
					parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
				},
			},
		},
		Result: parser.NewInteger(2),
	},
	{
		Name: "Aggregate Function Not Grouped Error",
		Filter: []FilterRecord{
			{
				View: &View{
					Header: NewHeader("table1", []string{"column1", "column2"}),
					Records: []Record{
						NewRecord(1, []parser.Primary{
							parser.NewInteger(1),
							parser.NewString("str2"),
						}),
					},
				},
				RecordIndex: 0,
			},
		},
		Expr: parser.Function{
			Name: "avg",
			Option: parser.Option{
				Args: []parser.Expression{
					parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
				},
			},
		},
		Error: "function avg: records are not grouped",
	},
	{
		Name: "Aggregate Function No Argument Error",
		Filter: []FilterRecord{
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
		Expr: parser.Function{
			Name: "avg",
			Option: parser.Option{
				Args: []parser.Expression{},
			},
		},
		Error: "function avg requires 1 argument",
	},
	{
		Name: "Aggregate Function Too Many Arguments Error",
		Filter: []FilterRecord{
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
		Expr: parser.Function{
			Name: "avg",
			Option: parser.Option{
				Args: []parser.Expression{
					parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
					parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
				},
			},
		},
		Error: "function avg has too many arguments",
	},
	{
		Name: "Aggregate Function Unpermitted AllColumns",
		Filter: []FilterRecord{
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
		Expr: parser.Function{
			Name: "avg",
			Option: parser.Option{
				Args: []parser.Expression{
					parser.AllColumns{},
				},
			},
		},
		Error: "syntax error: avg(*)",
	},
	{
		Name: "Aggregate Function Duplicate Error",
		Filter: []FilterRecord{
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
		Expr: parser.Function{
			Name: "avg",
			Option: parser.Option{
				Args: []parser.Expression{
					parser.Function{
						Name: "avg",
						Option: parser.Option{
							Args: []parser.Expression{
								parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
							},
						},
					},
				},
			},
		},
		Error: "syntax error: avg(avg(column1))",
	},
	{
		Name: "Aggregate Function Count With AllColumns",
		Filter: []FilterRecord{
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
		Expr: parser.Function{
			Name: "count",
			Option: parser.Option{
				Args: []parser.Expression{
					parser.AllColumns{},
				},
			},
		},
		Result: parser.NewInteger(3),
	},
	{
		Name: "GroupConcat Function",
		Filter: []FilterRecord{
			{
				View: &View{
					Header: NewHeader("table1", []string{"column1", "column2"}),
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
		Expr: parser.GroupConcat{
			GroupConcat: "group_concat",
			Option: parser.Option{
				Distinct: parser.Token{Token: parser.DISTINCT, Literal: "distinct"},
				Args: []parser.Expression{
					parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
				},
			},
			OrderBy: parser.OrderByClause{
				Items: []parser.Expression{
					parser.OrderItem{Item: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
				},
			},
			Separator: ",",
		},
		Result: parser.NewString("str1,str2"),
	},
	{
		Name: "GroupConcat Function Null",
		Filter: []FilterRecord{
			{
				View: &View{
					Header: NewHeader("table1", []string{"column1", "column2"}),
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
		Expr: parser.GroupConcat{
			Option: parser.Option{
				Distinct: parser.Token{Token: parser.DISTINCT, Literal: "distinct"},
				Args: []parser.Expression{
					parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
				},
			},
		},
		Result: parser.NewNull(),
	},
	{
		Name: "GroupConcat Function Not Grouped Error",
		Filter: []FilterRecord{
			{
				View: &View{
					Header: NewHeader("table1", []string{"column1", "column2"}),
					Records: []Record{
						NewRecord(1, []parser.Primary{
							parser.NewInteger(1),
							parser.NewString("str2"),
						}),
					},
				},
				RecordIndex: 0,
			},
		},
		Expr: parser.GroupConcat{
			GroupConcat: "group_concat",
			Option: parser.Option{
				Args: []parser.Expression{},
			},
		},
		Error: "function group_concat: records are not grouped",
	},
	{
		Name: "GroupConcat Function Argument Error",
		Filter: []FilterRecord{
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
		Expr: parser.GroupConcat{
			GroupConcat: "group_concat",
			Option: parser.Option{
				Args: []parser.Expression{},
			},
		},
		Error: "function group_concat takes 1 argument",
	},
	{
		Name: "GroupConcat Function AllColumns",
		Filter: []FilterRecord{
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
		Expr: parser.GroupConcat{
			GroupConcat: "group_concat",
			Option: parser.Option{
				Args: []parser.Expression{
					parser.AllColumns{},
				},
			},
		},
		Error: "syntax error: group_concat(*)",
	},
	{
		Name: "GroupConcat Function Identification Error",
		Filter: []FilterRecord{
			{
				View: &View{
					Header: NewHeader("table1", []string{"column1", "column2"}),
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
		Expr: parser.GroupConcat{
			GroupConcat: "group_concat",
			Option: parser.Option{
				Distinct: parser.Token{Token: parser.DISTINCT, Literal: "distinct"},
				Args: []parser.Expression{
					parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
				},
			},
			OrderBy: parser.OrderByClause{
				Items: []parser.Expression{
					parser.OrderItem{Item: parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}}},
				},
			},
			Separator: ",",
		},
		Error: "field notexist does not exist",
	},
	{
		Name: "GroupConcat Function Duplicate Error",
		Filter: []FilterRecord{
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
		Expr: parser.GroupConcat{
			GroupConcat: "group_concat",
			Option: parser.Option{
				Args: []parser.Expression{
					parser.Function{
						Name: "avg",
						Option: parser.Option{
							Args: []parser.Expression{
								parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
							},
						},
					},
				},
			},
		},
		Error: "syntax error: group_concat(avg(column1))",
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
						Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: "="},
					},
					Result: parser.NewString("A"),
				},
				parser.CaseWhen{
					Condition: parser.Comparison{
						LHS:      parser.NewInteger(2),
						RHS:      parser.NewInteger(2),
						Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: "="},
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
		Error: "field notexist does not exist",
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
		Error: "field notexist does not exist",
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
		Error: "field notexist does not exist",
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
		Error: "field notexist does not exist",
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
		Name: "Logic NOT",
		Expr: parser.Logic{
			RHS:      parser.NewTernary(ternary.FALSE),
			Operator: parser.Token{Token: parser.NOT, Literal: "not"},
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
		Error: "field notexist does not exist",
	},
	{
		Name: "Logic RHS Error",
		Expr: parser.Logic{
			LHS:      parser.NewTernary(ternary.FALSE),
			RHS:      parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			Operator: parser.Token{Token: parser.AND, Literal: "and"},
		},
		Error: "field notexist does not exist",
	},
	{
		Name: "Variable",
		Expr: parser.Variable{
			Name: "var1",
		},
		Result: parser.NewInteger(1),
	},
	{
		Name: "Variable Undefined Error",
		Expr: parser.Variable{
			Name: "undefined",
		},
		Error: "variable undefined is undefined",
	},
	{
		Name: "Variable Substitution",
		Expr: parser.VariableSubstitution{
			Variable: parser.Variable{Name: "var1"},
			Value:    parser.NewInteger(2),
		},
		Result: parser.NewInteger(2),
	},
	{
		Name: "Variable Substitution Undefined Error",
		Expr: parser.VariableSubstitution{
			Variable: parser.Variable{Name: "undefined"},
			Value:    parser.NewInteger(2),
		},
		Error: "variable undefined is undefined",
	},
}

func TestFilter_Evaluate(t *testing.T) {
	GlobalVars = map[string]parser.Primary{
		"var1": parser.NewInteger(1),
	}

	tf := cmd.GetFlags()
	tf.Repository = TestDataDir

	for _, v := range filterEvaluateTests {
		ViewCache.Clear()
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
