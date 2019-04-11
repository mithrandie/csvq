package query

import (
	"context"
	"os"
	"reflect"
	"sync"
	"testing"

	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/value"

	"github.com/mithrandie/ternary"
)

var filterEvaluateTests = []struct {
	Name          string
	Filter        *Filter
	Expr          parser.QueryExpression
	ReplaceValues *ReplaceValues
	Result        value.Primary
	Error         string
}{
	{
		Name:   "nil",
		Expr:   nil,
		Result: value.NewTernary(ternary.TRUE),
	},
	{
		Name:   "PrimitiveType",
		Expr:   parser.NewStringValue("str"),
		Result: value.NewString("str"),
	},
	{
		Name:  "Invalid Value Error",
		Expr:  parser.AllColumns{},
		Error: "*: cannot evaluate as a value",
	},
	{
		Name: "Parentheses",
		Expr: parser.Parentheses{
			Expr: parser.NewStringValue("str"),
		},
		Result: value.NewString("str"),
	},
	{
		Name: "FieldReference",
		Filter: &Filter{
			records: []filterRecord{
				{
					view: &View{
						Header: NewHeaderWithId("table1", []string{"column1", "column2"}),
						RecordSet: []Record{
							NewRecordWithId(1, []value.Primary{
								value.NewInteger(1),
								value.NewString("str"),
							}),
							NewRecordWithId(2, []value.Primary{
								value.NewInteger(2),
								value.NewString("strstr"),
							}),
						},
						Tx: TestTx,
					},
					recordIndex: 1,
				},
			},
		},
		Expr:   parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
		Result: value.NewString("strstr"),
	},
	{
		Name: "FieldReference ColumnNotExist Error",
		Filter: &Filter{
			records: []filterRecord{
				{
					view: &View{
						Header: NewHeaderWithId("table1", []string{"column1", "column2"}),
						RecordSet: []Record{
							NewRecordWithId(1, []value.Primary{
								value.NewInteger(1),
								value.NewString("str"),
							}),
							NewRecordWithId(2, []value.Primary{
								value.NewInteger(2),
								value.NewString("strstr"),
							}),
						},
						Tx: TestTx,
					},
					recordIndex: 1,
				},
			},
		},
		Expr:  parser.FieldReference{Column: parser.Identifier{Literal: "column3"}},
		Error: "field column3 does not exist",
	},
	{
		Name: "FieldReference FieldAmbigous Error",
		Filter: &Filter{
			records: []filterRecord{
				{
					view: &View{
						Header: NewHeaderWithId("table1", []string{"column1", "column1"}),
						RecordSet: []Record{
							NewRecordWithId(1, []value.Primary{
								value.NewInteger(1),
								value.NewString("str"),
							}),
							NewRecordWithId(2, []value.Primary{
								value.NewInteger(2),
								value.NewString("strstr"),
							}),
						},
						Tx: TestTx,
					},
					recordIndex: 1,
				},
			},
		},
		Expr:  parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
		Error: "field column1 is ambiguous",
	},
	{
		Name: "FieldReference Not Group Key Error",
		Filter: &Filter{
			records: []filterRecord{
				{
					view: &View{
						Header: []HeaderField{
							{
								View:        "table1",
								Column:      "column1",
								IsFromTable: true,
							},
							{
								View:   "table1",
								Column: "column2",
							},
						},
						RecordSet: []Record{
							{
								NewGroupCell([]value.Primary{
									value.NewInteger(1),
									value.NewInteger(2),
								}),
							},
							{
								NewGroupCell([]value.Primary{
									value.NewString("str1"),
									value.NewString("str2"),
								}),
							},
						},
						isGrouped: true,
						Tx:        TestTx,
					},
					recordIndex: 0,
				},
			},
		},
		Expr:  parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
		Error: "field column1 is not a group key",
	},
	{
		Name: "ColumnNumber",
		Filter: &Filter{
			records: []filterRecord{
				{
					view: &View{
						Header: NewHeaderWithId("table1", []string{"column1", "column2"}),
						RecordSet: []Record{
							NewRecordWithId(1, []value.Primary{
								value.NewInteger(1),
								value.NewString("str"),
							}),
							NewRecordWithId(2, []value.Primary{
								value.NewInteger(2),
								value.NewString("strstr"),
							}),
						},
						Tx: TestTx,
					},
					recordIndex: 1,
				},
			},
		},
		Expr:   parser.ColumnNumber{View: parser.Identifier{Literal: "table1"}, Number: value.NewInteger(2)},
		Result: value.NewString("strstr"),
	},
	{
		Name: "ColumnNumber Not Exist Error",
		Filter: &Filter{
			records: []filterRecord{
				{
					view: &View{
						Header: NewHeaderWithId("table1", []string{"column1", "column2"}),
						RecordSet: []Record{
							NewRecordWithId(1, []value.Primary{
								value.NewInteger(1),
								value.NewString("str"),
							}),
							NewRecordWithId(2, []value.Primary{
								value.NewInteger(2),
								value.NewString("strstr"),
							}),
						},
						Tx: TestTx,
					},
					recordIndex: 1,
				},
			},
		},
		Expr:  parser.ColumnNumber{View: parser.Identifier{Literal: "table1"}, Number: value.NewInteger(9)},
		Error: "field table1.9 does not exist",
	},
	{
		Name: "ColumnNumber Not Group Key Error",
		Filter: &Filter{
			records: []filterRecord{
				{
					view: &View{
						Header: []HeaderField{
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
								IsGroupKey:  true,
							},
						},
						RecordSet: []Record{
							{
								NewGroupCell([]value.Primary{
									value.NewInteger(1),
									value.NewInteger(2),
								}),
							},
							{
								NewGroupCell([]value.Primary{
									value.NewString("str1"),
									value.NewString("str2"),
								}),
							},
						},
						isGrouped: true,
						Tx:        TestTx,
					},
					recordIndex: 0,
				},
			},
		},
		Expr:  parser.ColumnNumber{View: parser.Identifier{Literal: "table1"}, Number: value.NewInteger(1)},
		Error: "field table1.1 is not a group key",
	},
	{
		Name: "Arithmetic",
		Expr: parser.Arithmetic{
			LHS:      parser.NewIntegerValue(1),
			RHS:      parser.NewIntegerValue(2),
			Operator: '+',
		},
		Result: value.NewInteger(3),
	},
	{
		Name: "Arithmetic LHS Error",
		Expr: parser.Arithmetic{
			LHS:      parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			RHS:      parser.NewIntegerValue(2),
			Operator: '+',
		},
		Error: "field notexist does not exist",
	},
	{
		Name: "Arithmetic LHS Is Null",
		Expr: parser.Arithmetic{
			LHS:      parser.NewNullValue(),
			RHS:      parser.NewIntegerValue(2),
			Operator: '+',
		},
		Result: value.NewNull(),
	},
	{
		Name: "Arithmetic RHS Error",
		Expr: parser.Arithmetic{
			LHS:      parser.NewIntegerValue(1),
			RHS:      parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			Operator: '+',
		},
		Error: "field notexist does not exist",
	},
	{
		Name: "UnaryArithmetic Integer",
		Expr: parser.UnaryArithmetic{
			Operand:  parser.NewIntegerValue(1),
			Operator: parser.Token{Token: '-', Literal: "-"},
		},
		Result: value.NewInteger(-1),
	},
	{
		Name: "UnaryArithmetic Float",
		Expr: parser.UnaryArithmetic{
			Operand:  parser.NewFloatValue(1.234),
			Operator: parser.Token{Token: '-', Literal: "-"},
		},
		Result: value.NewFloat(-1.234),
	},
	{
		Name: "UnaryArithmetic Operand Error",
		Expr: parser.UnaryArithmetic{
			Operand:  parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			Operator: parser.Token{Token: '-', Literal: "-"},
		},
		Error: "field notexist does not exist",
	},
	{
		Name: "UnaryArithmetic Cast Failure",
		Expr: parser.UnaryArithmetic{
			Operand:  parser.NewStringValue("str"),
			Operator: parser.Token{Token: '-', Literal: "-"},
		},
		Result: value.NewNull(),
	},
	{
		Name: "Concat",
		Expr: parser.Concat{
			Items: []parser.QueryExpression{
				parser.NewStringValue("a"),
				parser.NewStringValue("b"),
				parser.NewStringValue("c"),
			},
		},
		Result: value.NewString("abc"),
	},
	{
		Name: "Concat FieldNotExist Error",
		Expr: parser.Concat{
			Items: []parser.QueryExpression{
				parser.NewStringValue("a"),
				parser.NewStringValue("b"),
				parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			},
		},
		Error: "field notexist does not exist",
	},
	{
		Name: "Concat Including Null",
		Expr: parser.Concat{
			Items: []parser.QueryExpression{
				parser.NewStringValue("a"),
				parser.NewNullValue(),
				parser.NewStringValue("c"),
			},
		},
		Result: value.NewNull(),
	},
	{
		Name: "Comparison",
		Expr: parser.Comparison{
			LHS:      parser.NewIntegerValue(1),
			RHS:      parser.NewIntegerValue(2),
			Operator: "=",
		},
		Result: value.NewTernary(ternary.FALSE),
	},
	{
		Name: "Comparison LHS Error",
		Expr: parser.Comparison{
			LHS:      parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			RHS:      parser.NewIntegerValue(2),
			Operator: "=",
		},
		Error: "field notexist does not exist",
	},
	{
		Name: "Comparison LHS Is Null",
		Expr: parser.Comparison{
			LHS:      parser.NewNullValue(),
			RHS:      parser.NewIntegerValue(2),
			Operator: "=",
		},
		Result: value.NewTernary(ternary.UNKNOWN),
	},
	{
		Name: "Comparison RHS Error",
		Expr: parser.Comparison{
			LHS:      parser.NewIntegerValue(1),
			RHS:      parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			Operator: "=",
		},
		Error: "field notexist does not exist",
	},
	{
		Name: "Comparison with Row Values",
		Expr: parser.Comparison{
			LHS: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewIntegerValue(1),
						parser.NewIntegerValue(2),
					},
				},
			},
			RHS: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewIntegerValue(1),
						parser.NewIntegerValue(2),
					},
				},
			},
			Operator: "=",
		},
		Result: value.NewTernary(ternary.TRUE),
	},
	{
		Name: "Comparison with Row Value and Subquery",
		Expr: parser.Comparison{
			LHS: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewStringValue("1"),
						parser.NewStringValue("str1"),
					},
				},
			},
			RHS: parser.RowValue{
				Value: parser.Subquery{
					Query: parser.SelectQuery{
						SelectEntity: parser.SelectEntity{
							SelectClause: parser.SelectClause{
								Select: "select",
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
							WhereClause: parser.WhereClause{
								Filter: parser.Comparison{
									LHS:      parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
									RHS:      parser.NewIntegerValue(1),
									Operator: "=",
								},
							},
						},
					},
				},
			},
			Operator: "=",
		},
		Result: value.NewTernary(ternary.TRUE),
	},
	{
		Name: "Comparison with Row Values LHS Error",
		Expr: parser.Comparison{
			LHS: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
						parser.NewIntegerValue(2),
					},
				},
			},
			RHS: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewIntegerValue(1),
						parser.NewIntegerValue(2),
					},
				},
			},
			Operator: "=",
		},
		Error: "field notexist does not exist",
	},
	{
		Name: "Comparison with Row Value and Subquery Returns No Record",
		Expr: parser.Comparison{
			LHS: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewStringValue("1"),
						parser.NewStringValue("str1"),
					},
				},
			},
			RHS: parser.RowValue{
				Value: parser.Subquery{
					Query: parser.SelectQuery{
						SelectEntity: parser.SelectEntity{
							SelectClause: parser.SelectClause{
								Select: "select",
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
							WhereClause: parser.WhereClause{
								Filter: parser.NewTernaryValue(ternary.FALSE),
							},
						},
					},
				},
			},
			Operator: "=",
		},
		Result: value.NewTernary(ternary.UNKNOWN),
	},
	{
		Name: "Comparison with Row Value and LHS Subquery Returns No Record",
		Expr: parser.Comparison{
			LHS: parser.RowValue{
				Value: parser.Subquery{
					Query: parser.SelectQuery{
						SelectEntity: parser.SelectEntity{
							SelectClause: parser.SelectClause{
								Select: "select",
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
							WhereClause: parser.WhereClause{
								Filter: parser.NewTernaryValue(ternary.FALSE),
							},
						},
					},
				},
			},
			RHS: parser.RowValue{
				Value: parser.Subquery{
					Query: parser.SelectQuery{
						SelectEntity: parser.SelectEntity{
							SelectClause: parser.SelectClause{
								Select: "select",
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
							WhereClause: parser.WhereClause{
								Filter: parser.NewTernaryValue(ternary.FALSE),
							},
						},
					},
				},
			},
			Operator: "=",
		},
		Result: value.NewTernary(ternary.UNKNOWN),
	},
	{
		Name: "Comparison with Row Value and Subquery Query Error",
		Expr: parser.Comparison{
			LHS: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewStringValue("1"),
						parser.NewStringValue("str1"),
					},
				},
			},
			RHS: parser.RowValue{
				Value: parser.Subquery{
					Query: parser.SelectQuery{
						SelectEntity: parser.SelectEntity{
							SelectClause: parser.SelectClause{
								Select: "select",
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
					},
				},
			},
			Operator: "=",
		},
		Error: "field notexist does not exist",
	},
	{
		Name: "Comparison with Row Values RHS Error",
		Expr: parser.Comparison{
			LHS: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewIntegerValue(1),
						parser.NewIntegerValue(2),
					},
				},
			},
			RHS: parser.RowValue{
				Value: parser.Subquery{
					Query: parser.SelectQuery{
						SelectEntity: parser.SelectEntity{
							SelectClause: parser.SelectClause{
								Select: "select",
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
			},
			Operator: "=",
		},
		Error: "subquery returns too many records, should return only one record",
	},
	{
		Name: "Comparison with Row Values Value Length Not Match Error",
		Expr: parser.Comparison{
			LHS: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewIntegerValue(1),
						parser.NewIntegerValue(2),
					},
				},
			},
			RHS: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewIntegerValue(1),
					},
				},
			},
			Operator: "=",
		},
		Error: "row value should contain exactly 2 values",
	},
	{
		Name: "Comparison with Row Value and JsonQuery",
		Expr: parser.Comparison{
			LHS: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewStringValue("1"),
						parser.NewStringValue("str1"),
					},
				},
			},
			RHS: parser.RowValue{
				Value: parser.JsonQuery{
					Query:    parser.NewStringValue("key{key2, key3}"),
					JsonText: parser.NewStringValue("{\"key\": {\"key2\": 1, \"key3\": \"str1\"}}"),
				},
			},
			Operator: "=",
		},
		Result: value.NewTernary(ternary.TRUE),
	},
	{
		Name: "Comparison with Row Value and JsonQuery Query Evaluation Error",
		Expr: parser.Comparison{
			LHS: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewStringValue("1"),
						parser.NewStringValue("str1"),
					},
				},
			},
			RHS: parser.RowValue{
				Value: parser.JsonQuery{
					Query:    parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
					JsonText: parser.NewStringValue("{\"key\": {\"key2\": 1, \"key3\": \"str1\"}}"),
				},
			},
			Operator: "=",
		},
		Error: "field notexist does not exist",
	},
	{
		Name: "Comparison with Row Value and JsonQuery Query is Null",
		Expr: parser.Comparison{
			LHS: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewStringValue("1"),
						parser.NewStringValue("str1"),
					},
				},
			},
			RHS: parser.RowValue{
				Value: parser.JsonQuery{
					Query:    parser.NewNullValue(),
					JsonText: parser.NewStringValue("{\"key\": {\"key2\": 1, \"key3\": \"str1\"}}"),
				},
			},
			Operator: "=",
		},
		Result: value.NewTernary(ternary.UNKNOWN),
	},
	{
		Name: "Comparison with Row Value and JsonQuery Type Error",
		Expr: parser.Comparison{
			LHS: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewStringValue("1"),
						parser.NewStringValue("str1"),
					},
				},
			},
			RHS: parser.RowValue{
				Value: parser.JsonQuery{
					Query:    parser.NewStringValue("key[]"),
					JsonText: parser.NewStringValue("{\"key\": {\"key2\": 1, \"key3\": \"str1\"}}"),
				},
			},
			Operator: "=",
		},
		Error: "json loading error: json value must be an array",
	},
	{
		Name: "Comparison with Row Value and JsonQuery Empty Result Set",
		Expr: parser.Comparison{
			LHS: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewStringValue("1"),
						parser.NewStringValue("str1"),
					},
				},
			},
			RHS: parser.RowValue{
				Value: parser.JsonQuery{
					Query:    parser.NewStringValue("key{}"),
					JsonText: parser.NewStringValue("{\"key\": []}"),
				},
			},
			Operator: "=",
		},
		Result: value.NewTernary(ternary.UNKNOWN),
	},
	{
		Name: "Comparison with Row Value and JsonQuery Too Many Records",
		Expr: parser.Comparison{
			LHS: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewStringValue("1"),
						parser.NewStringValue("str1"),
					},
				},
			},
			RHS: parser.RowValue{
				Value: parser.JsonQuery{
					Query:    parser.NewStringValue("key{key2, key3}"),
					JsonText: parser.NewStringValue("{\"key\": [{\"key2\": 1, \"key3\": \"str1\"}, {\"key2\": 1, \"key3\": \"str1\"}] }"),
				},
			},
			Operator: "=",
		},
		Error: "json query returns too many records, should return only one record",
	},
	{
		Name: "Is",
		Expr: parser.Is{
			LHS:      parser.NewIntegerValue(1),
			RHS:      parser.NewNullValue(),
			Negation: parser.Token{Token: parser.NOT, Literal: "not"},
		},
		Result: value.NewTernary(ternary.TRUE),
	},
	{
		Name: "Is LHS Error",
		Expr: parser.Is{
			LHS:      parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			RHS:      parser.NewNullValue(),
			Negation: parser.Token{Token: parser.NOT, Literal: "not"},
		},
		Error: "field notexist does not exist",
	},
	{
		Name: "Is RHS Error",
		Expr: parser.Is{
			LHS:      parser.NewIntegerValue(1),
			RHS:      parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			Negation: parser.Token{Token: parser.NOT, Literal: "not"},
		},
		Error: "field notexist does not exist",
	},
	{
		Name: "Between",
		Expr: parser.Between{
			LHS:      parser.NewIntegerValue(2),
			Low:      parser.NewIntegerValue(1),
			High:     parser.NewIntegerValue(3),
			Negation: parser.Token{Token: parser.NOT, Literal: "not"},
		},
		Result: value.NewTernary(ternary.FALSE),
	},
	{
		Name: "Between LHS Error",
		Expr: parser.Between{
			LHS:      parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			Low:      parser.NewIntegerValue(1),
			High:     parser.NewIntegerValue(3),
			Negation: parser.Token{Token: parser.NOT, Literal: "not"},
		},
		Error: "field notexist does not exist",
	},
	{
		Name: "Between LHS Is Null",
		Expr: parser.Between{
			LHS:      parser.NewNullValue(),
			Low:      parser.NewIntegerValue(1),
			High:     parser.NewIntegerValue(3),
			Negation: parser.Token{Token: parser.NOT, Literal: "not"},
		},
		Result: value.NewTernary(ternary.UNKNOWN),
	},
	{
		Name: "Between Low Error",
		Expr: parser.Between{
			LHS:      parser.NewIntegerValue(2),
			Low:      parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			High:     parser.NewIntegerValue(3),
			Negation: parser.Token{Token: parser.NOT, Literal: "not"},
		},
		Error: "field notexist does not exist",
	},
	{
		Name: "Between Low Comparison False",
		Expr: parser.Between{
			LHS:      parser.NewIntegerValue(2),
			Low:      parser.NewIntegerValue(3),
			High:     parser.NewIntegerValue(5),
			Negation: parser.Token{Token: parser.NOT, Literal: "not"},
		},
		Result: value.NewTernary(ternary.TRUE),
	},
	{
		Name: "Between High Error",
		Expr: parser.Between{
			LHS:      parser.NewIntegerValue(2),
			Low:      parser.NewIntegerValue(1),
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
					Values: []parser.QueryExpression{
						parser.NewIntegerValue(1),
						parser.NewIntegerValue(2),
					},
				},
			},
			Low: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewIntegerValue(1),
						parser.NewIntegerValue(1),
					},
				},
			},
			High: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewIntegerValue(1),
						parser.NewIntegerValue(3),
					},
				},
			},
		},
		Result: value.NewTernary(ternary.TRUE),
	},
	{
		Name: "Between with LHS Subquery Returns No Records",
		Expr: parser.Between{
			LHS: parser.RowValue{
				Value: parser.Subquery{
					Query: parser.SelectQuery{
						SelectEntity: parser.SelectEntity{
							SelectClause: parser.SelectClause{
								Select: "select",
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
							WhereClause: parser.WhereClause{
								Filter: parser.NewTernaryValue(ternary.FALSE),
							},
						},
					},
				},
			},
			Low: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewIntegerValue(1),
						parser.NewIntegerValue(1),
					},
				},
			},
			High: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewIntegerValue(1),
						parser.NewIntegerValue(3),
					},
				},
			},
		},
		Result: value.NewTernary(ternary.UNKNOWN),
	},
	{
		Name: "Between with Row Values LHS Error",
		Expr: parser.Between{
			LHS: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
						parser.NewIntegerValue(2),
					},
				},
			},
			Low: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewIntegerValue(1),
						parser.NewIntegerValue(1),
					},
				},
			},
			High: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewIntegerValue(1),
						parser.NewIntegerValue(3),
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
					Values: []parser.QueryExpression{
						parser.NewIntegerValue(1),
						parser.NewIntegerValue(2),
					},
				},
			},
			Low: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
						parser.NewIntegerValue(1),
					},
				},
			},
			High: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewIntegerValue(1),
						parser.NewIntegerValue(3),
					},
				},
			},
		},
		Error: "field notexist does not exist",
	},
	{
		Name: "Between with Row Values Low Comparison False",
		Expr: parser.Between{
			LHS: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewIntegerValue(1),
						parser.NewIntegerValue(2),
					},
				},
			},
			Low: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewIntegerValue(1),
						parser.NewIntegerValue(3),
					},
				},
			},
			High: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewIntegerValue(1),
						parser.NewIntegerValue(5),
					},
				},
			},
		},
		Result: value.NewTernary(ternary.FALSE),
	},
	{
		Name: "Between with Row Values High Error",
		Expr: parser.Between{
			LHS: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewIntegerValue(1),
						parser.NewIntegerValue(2),
					},
				},
			},
			Low: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewIntegerValue(1),
						parser.NewIntegerValue(1),
					},
				},
			},
			High: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
						parser.NewIntegerValue(3),
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
					Values: []parser.QueryExpression{
						parser.NewIntegerValue(1),
						parser.NewIntegerValue(2),
					},
				},
			},
			Low: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewIntegerValue(1),
					},
				},
			},
			High: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewIntegerValue(1),
						parser.NewIntegerValue(3),
					},
				},
			},
		},
		Error: "row value should contain exactly 2 values",
	},
	{
		Name: "Between with Row Values High Comparison Error",
		Expr: parser.Between{
			LHS: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewIntegerValue(1),
						parser.NewIntegerValue(2),
					},
				},
			},
			Low: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewIntegerValue(1),
						parser.NewIntegerValue(1),
					},
				},
			},
			High: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewIntegerValue(3),
					},
				},
			},
		},
		Error: "row value should contain exactly 2 values",
	},
	{
		Name: "In",
		Expr: parser.In{
			LHS: parser.NewIntegerValue(2),
			Values: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewIntegerValue(1),
						parser.NewIntegerValue(2),
						parser.NewIntegerValue(3),
					},
				},
			},
			Negation: parser.Token{Token: parser.NOT, Literal: "not"},
		},
		Result: value.NewTernary(ternary.FALSE),
	},
	{
		Name: "In LHS Error",
		Expr: parser.In{
			LHS: parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			Values: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewIntegerValue(1),
						parser.NewIntegerValue(2),
						parser.NewIntegerValue(3),
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
			LHS: parser.NewIntegerValue(2),
			Values: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewIntegerValue(1),
						parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
						parser.NewIntegerValue(3),
					},
				},
			},
			Negation: parser.Token{Token: parser.NOT, Literal: "not"},
		},
		Error: "field notexist does not exist",
	},
	{
		Name: "In Subquery",
		Filter: &Filter{
			records: []filterRecord{
				{
					view: &View{
						Header: NewHeaderWithId("table2", []string{"column3", "column4"}),
						RecordSet: []Record{
							NewRecordWithId(1, []value.Primary{
								value.NewInteger(1),
								value.NewString("str2"),
							}),
						},
					},
					recordIndex: 0,
				},
			},
		},
		Expr: parser.In{
			LHS: parser.NewIntegerValue(2),
			Values: parser.RowValue{
				Value: parser.Subquery{
					Query: parser.SelectQuery{
						SelectEntity: parser.SelectEntity{
							SelectClause: parser.SelectClause{
								Select: "select",
								Fields: []parser.QueryExpression{
									parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
								},
							},
							FromClause: parser.FromClause{
								Tables: []parser.QueryExpression{
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
		Result: value.NewTernary(ternary.FALSE),
	},
	{
		Name: "In Subquery Execution Error",
		Expr: parser.In{
			LHS: parser.NewIntegerValue(2),
			Values: parser.RowValue{
				Value: parser.Subquery{
					Query: parser.SelectQuery{
						SelectEntity: parser.SelectEntity{
							SelectClause: parser.SelectClause{
								Select: "select",
								Fields: []parser.QueryExpression{
									parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
								},
							},
							FromClause: parser.FromClause{
								Tables: []parser.QueryExpression{
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
		Error: "field notexist does not exist",
	},
	{
		Name: "In Subquery Too Many Field Error",
		Expr: parser.In{
			LHS: parser.NewIntegerValue(2),
			Values: parser.RowValue{
				Value: parser.Subquery{
					Query: parser.SelectQuery{
						SelectEntity: parser.SelectEntity{
							SelectClause: parser.SelectClause{
								Select: "select",
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
			},
			Negation: parser.Token{Token: parser.NOT, Literal: "not"},
		},
		Error: "subquery returns too many fields, should return only one field",
	},
	{
		Name: "In Subquery Returns No Record",
		Expr: parser.In{
			LHS: parser.NewIntegerValue(2),
			Values: parser.RowValue{
				Value: parser.Subquery{
					Query: parser.SelectQuery{
						SelectEntity: parser.SelectEntity{
							SelectClause: parser.SelectClause{
								Select: "select",
								Fields: []parser.QueryExpression{
									parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
								},
							},
							FromClause: parser.FromClause{
								Tables: []parser.QueryExpression{
									parser.Table{Object: parser.Identifier{Literal: "table1"}},
								},
							},
							WhereClause: parser.WhereClause{
								Filter: parser.NewTernaryValue(ternary.FALSE),
							},
						},
					},
				},
			},
		},
		Result: value.NewTernary(ternary.FALSE),
	},
	{
		Name: "In JsonArray",
		Filter: &Filter{
			records: []filterRecord{
				{
					view: &View{
						Header: NewHeaderWithId("table2", []string{"column3", "column4"}),
						RecordSet: []Record{
							NewRecordWithId(1, []value.Primary{
								value.NewInteger(1),
								value.NewString("str2"),
							}),
						},
						Tx: TestTx,
					},
					recordIndex: 0,
				},
			},
		},
		Expr: parser.In{
			LHS: parser.NewIntegerValue(2),
			Values: parser.RowValue{
				Value: parser.JsonQuery{
					Query:    parser.NewStringValue("key[]"),
					JsonText: parser.NewStringValue("{\"key\":[2, 3]}"),
				},
			},
			Negation: parser.Token{Token: parser.NOT, Literal: "not"},
		},
		Result: value.NewTernary(ternary.FALSE),
	},
	{
		Name: "In JsonArray Query is Null",
		Filter: &Filter{
			records: []filterRecord{
				{
					view: &View{
						Header: NewHeaderWithId("table2", []string{"column3", "column4"}),
						RecordSet: []Record{
							NewRecordWithId(1, []value.Primary{
								value.NewInteger(1),
								value.NewString("str2"),
							}),
						},
						Tx: TestTx,
					},
					recordIndex: 0,
				},
			},
		},
		Expr: parser.In{
			LHS: parser.NewIntegerValue(2),
			Values: parser.RowValue{
				Value: parser.JsonQuery{
					Query:    parser.NewNullValue(),
					JsonText: parser.NewStringValue("{\"key\":[2, 3]}"),
				},
			},
			Negation: parser.Token{Token: parser.NOT, Literal: "not"},
		},
		Result: value.NewTernary(ternary.TRUE),
	},
	{
		Name: "In JsonArray Query Evaluation Error",
		Filter: &Filter{
			records: []filterRecord{
				{
					view: &View{
						Header: NewHeaderWithId("table2", []string{"column3", "column4"}),
						RecordSet: []Record{
							NewRecordWithId(1, []value.Primary{
								value.NewInteger(1),
								value.NewString("str2"),
							}),
						},
					},
					recordIndex: 0,
				},
			},
		},
		Expr: parser.In{
			LHS: parser.NewIntegerValue(2),
			Values: parser.RowValue{
				Value: parser.JsonQuery{
					Query:    parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
					JsonText: parser.NewStringValue("{\"key\":[2, 3]}"),
				},
			},
			Negation: parser.Token{Token: parser.NOT, Literal: "not"},
		},
		Error: "field notexist does not exist",
	},
	{
		Name: "In JsonArray JsonText Evaluation Error",
		Filter: &Filter{
			records: []filterRecord{
				{
					view: &View{
						Header: NewHeaderWithId("table2", []string{"column3", "column4"}),
						RecordSet: []Record{
							NewRecordWithId(1, []value.Primary{
								value.NewInteger(1),
								value.NewString("str2"),
							}),
						},
						Tx: TestTx,
					},
					recordIndex: 0,
				},
			},
		},
		Expr: parser.In{
			LHS: parser.NewIntegerValue(2),
			Values: parser.RowValue{
				Value: parser.JsonQuery{
					Query:    parser.NewStringValue("key"),
					JsonText: parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
				},
			},
			Negation: parser.Token{Token: parser.NOT, Literal: "not"},
		},
		Error: "field notexist does not exist",
	},
	{
		Name: "In JsonArray Query Load Error",
		Filter: &Filter{
			records: []filterRecord{
				{
					view: &View{
						Header: NewHeaderWithId("table2", []string{"column3", "column4"}),
						RecordSet: []Record{
							NewRecordWithId(1, []value.Primary{
								value.NewInteger(1),
								value.NewString("str2"),
							}),
						},
						Tx: TestTx,
					},
					recordIndex: 0,
				},
			},
		},
		Expr: parser.In{
			LHS: parser.NewIntegerValue(2),
			Values: parser.RowValue{
				Value: parser.JsonQuery{
					Query:    parser.NewStringValue("'key"),
					JsonText: parser.NewStringValue("{\"key\":[2, 3]}"),
				},
			},
			Negation: parser.Token{Token: parser.NOT, Literal: "not"},
		},
		Error: "json loading error: column 4: string not terminated",
	},
	{
		Name: "In JsonArray Empty Result Set",
		Filter: &Filter{
			records: []filterRecord{
				{
					view: &View{
						Header: NewHeaderWithId("table2", []string{"column3", "column4"}),
						RecordSet: []Record{
							NewRecordWithId(1, []value.Primary{
								value.NewInteger(1),
								value.NewString("str2"),
							}),
						},
					},
					recordIndex: 0,
				},
			},
		},
		Expr: parser.In{
			LHS: parser.NewIntegerValue(2),
			Values: parser.RowValue{
				Value: parser.JsonQuery{
					Query:    parser.NewStringValue("key[]"),
					JsonText: parser.NewStringValue("{\"key\":[]}"),
				},
			},
			Negation: parser.Token{Token: parser.NOT, Literal: "not"},
		},
		Result: value.NewTernary(ternary.TRUE),
	},
	{
		Name: "In with Row Values",
		Expr: parser.In{
			LHS: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewIntegerValue(1),
						parser.NewIntegerValue(2),
					},
				},
			},
			Values: parser.RowValueList{
				RowValues: []parser.QueryExpression{
					parser.RowValue{
						Value: parser.ValueList{
							Values: []parser.QueryExpression{
								parser.NewIntegerValue(1),
								parser.NewIntegerValue(1),
							},
						},
					},
					parser.RowValue{
						Value: parser.ValueList{
							Values: []parser.QueryExpression{
								parser.NewIntegerValue(1),
								parser.NewIntegerValue(2),
							},
						},
					},
				},
			},
		},
		Result: value.NewTernary(ternary.TRUE),
	},
	{
		Name: "In with Row Value and Subquery",
		Expr: parser.In{
			LHS: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewStringValue("2"),
						parser.NewStringValue("str2"),
					},
				},
			},
			Values: parser.Subquery{
				Query: parser.SelectQuery{
					SelectEntity: parser.SelectEntity{
						SelectClause: parser.SelectClause{
							Select: "select",
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
		},
		Result: value.NewTernary(ternary.TRUE),
	},
	{
		Name: "In with Row Value and JsonQuery",
		Expr: parser.In{
			LHS: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewStringValue("2"),
						parser.NewStringValue("str2"),
					},
				},
			},
			Values: parser.JsonQuery{
				Query:    parser.NewStringValue("key{key2, key3}"),
				JsonText: parser.NewStringValue("{\"key\":{\"key2\":2, \"key3\":\"str2\"}}"),
			},
		},
		Result: value.NewTernary(ternary.TRUE),
	},
	{
		Name: "In with Row Values LHS Error",
		Expr: parser.In{
			LHS: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
						parser.NewIntegerValue(2),
					},
				},
			},
			Values: parser.RowValueList{
				RowValues: []parser.QueryExpression{
					parser.RowValue{
						Value: parser.ValueList{
							Values: []parser.QueryExpression{
								parser.NewIntegerValue(1),
								parser.NewIntegerValue(1),
							},
						},
					},
					parser.RowValue{
						Value: parser.ValueList{
							Values: []parser.QueryExpression{
								parser.NewIntegerValue(1),
								parser.NewIntegerValue(2),
							},
						},
					},
				},
			},
		},
		Error: "field notexist does not exist",
	},
	{
		Name: "In with Row Value and Subquery Query Error",
		Expr: parser.In{
			LHS: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewStringValue("2"),
						parser.NewStringValue("str2"),
					},
				},
			},
			Values: parser.Subquery{
				Query: parser.SelectQuery{
					SelectEntity: parser.SelectEntity{
						SelectClause: parser.SelectClause{
							Select: "select",
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
				},
			},
		},
		Error: "field notexist does not exist",
	},
	{
		Name: "In with Row Value and Subquery Empty Result Set",
		Expr: parser.In{
			LHS: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewStringValue("2"),
						parser.NewStringValue("str2"),
					},
				},
			},
			Values: parser.Subquery{
				Query: parser.SelectQuery{
					SelectEntity: parser.SelectEntity{
						SelectClause: parser.SelectClause{
							Select: "select",
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
						WhereClause: parser.WhereClause{
							Filter: parser.NewTernaryValue(ternary.FALSE),
						},
					},
				},
			},
		},
		Result: value.NewTernary(ternary.FALSE),
	},
	{
		Name: "In with Row Value and JsonQuery Query Evaluation Error",
		Expr: parser.In{
			LHS: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewStringValue("2"),
						parser.NewStringValue("str2"),
					},
				},
			},
			Values: parser.JsonQuery{
				Query:    parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
				JsonText: parser.NewStringValue("{\"key\":{\"key2\":2, \"key3\":\"str2\"}}"),
			},
		},
		Error: "field notexist does not exist",
	},
	{
		Name: "In with Row Value and JsonQuery Query is Null",
		Expr: parser.In{
			LHS: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewStringValue("2"),
						parser.NewStringValue("str2"),
					},
				},
			},
			Values: parser.JsonQuery{
				Query:    parser.NewNullValue(),
				JsonText: parser.NewStringValue("{\"key\":{\"key2\":2, \"key3\":\"str2\"}}"),
			},
		},
		Result: value.NewTernary(ternary.FALSE),
	},
	{
		Name: "In with Row Value and JsonQuery Loading Error",
		Expr: parser.In{
			LHS: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewStringValue("2"),
						parser.NewStringValue("str2"),
					},
				},
			},
			Values: parser.JsonQuery{
				Query:    parser.NewStringValue("key{"),
				JsonText: parser.NewStringValue("{\"key\":[]}"),
			},
		},
		Error: "json loading error: column 4: unexpected termination",
	},
	{
		Name: "In with Row Value and JsonQuery Empty Result Set",
		Expr: parser.In{
			LHS: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewStringValue("2"),
						parser.NewStringValue("str2"),
					},
				},
			},
			Values: parser.JsonQuery{
				Query:    parser.NewStringValue("key{}"),
				JsonText: parser.NewStringValue("{\"key\":[]}"),
			},
		},
		Result: value.NewTernary(ternary.FALSE),
	},
	{
		Name: "In with Row Value and JsonQuery Field Length Error",
		Expr: parser.In{
			LHS: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewStringValue("2"),
						parser.NewStringValue("str2"),
					},
				},
			},
			Values: parser.JsonQuery{
				Query:    parser.NewStringValue("key{}"),
				JsonText: parser.NewStringValue("{\"key\":{}}"),
			},
		},
		Error: "row value should contain exactly 2 values",
	},
	{
		Name: "In with Row Values Values Error",
		Expr: parser.In{
			LHS: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewIntegerValue(1),
						parser.NewIntegerValue(2),
					},
				},
			},
			Values: parser.RowValueList{
				RowValues: []parser.QueryExpression{
					parser.RowValue{
						Value: parser.ValueList{
							Values: []parser.QueryExpression{
								parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
								parser.NewIntegerValue(1),
							},
						},
					},
					parser.RowValue{
						Value: parser.ValueList{
							Values: []parser.QueryExpression{
								parser.NewIntegerValue(1),
								parser.NewIntegerValue(2),
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
					Values: []parser.QueryExpression{
						parser.NewIntegerValue(1),
						parser.NewIntegerValue(2),
					},
				},
			},
			Values: parser.RowValueList{
				RowValues: []parser.QueryExpression{
					parser.RowValue{
						Value: parser.ValueList{
							Values: []parser.QueryExpression{
								parser.NewIntegerValue(1),
								parser.NewIntegerValue(1),
							},
						},
					},
					parser.RowValue{
						Value: parser.ValueList{
							Values: []parser.QueryExpression{
								parser.NewIntegerValue(2),
							},
						},
					},
				},
			},
		},
		Error: "row value should contain exactly 2 values",
	},
	{
		Name: "In with Row Value and Subquery Length Not Match Error",
		Expr: parser.In{
			LHS: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewStringValue("2"),
						parser.NewStringValue("str2"),
					},
				},
			},
			Values: parser.Subquery{
				Query: parser.SelectQuery{
					SelectEntity: parser.SelectEntity{
						SelectClause: parser.SelectClause{
							Select: "select",
							Fields: []parser.QueryExpression{
								parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
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
		},
		Error: "select query should return exactly 2 fields",
	},
	{
		Name: "Any",
		Expr: parser.Any{
			LHS: parser.NewIntegerValue(5),
			Values: parser.RowValue{
				Value: parser.Subquery{
					Query: parser.SelectQuery{
						SelectEntity: parser.SelectEntity{
							SelectClause: parser.SelectClause{
								Select: "select",
								Fields: []parser.QueryExpression{
									parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
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
			},
			Operator: "<>",
		},
		Result: value.NewTernary(ternary.TRUE),
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
								Fields: []parser.QueryExpression{
									parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
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
			},
			Operator: "<>",
		},
		Error: "field notexist does not exist",
	},
	{
		Name: "Any Query Execution Error",
		Expr: parser.Any{
			LHS: parser.NewIntegerValue(2),
			Values: parser.RowValue{
				Value: parser.Subquery{
					Query: parser.SelectQuery{
						SelectEntity: parser.SelectEntity{
							SelectClause: parser.SelectClause{
								Select: "select",
								Fields: []parser.QueryExpression{
									parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
								},
							},
							FromClause: parser.FromClause{
								Tables: []parser.QueryExpression{
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
		Error: "field notexist does not exist",
	},
	{
		Name: "Any Row Value Select Field Not Match Error",
		Expr: parser.Any{
			LHS: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewIntegerValue(1),
						parser.NewIntegerValue(2),
					},
				},
			},
			Values: parser.Subquery{
				Query: parser.SelectQuery{
					SelectEntity: parser.SelectEntity{
						SelectClause: parser.SelectClause{
							Select: "select",
							Fields: []parser.QueryExpression{
								parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
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
			Operator: "<>",
		},
		Error: "select query should return exactly 2 fields",
	},
	{
		Name: "Any Row Value Length Not Match Error",
		Expr: parser.Any{
			LHS: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewIntegerValue(1),
						parser.NewIntegerValue(2),
					},
				},
			},
			Values: parser.RowValueList{
				RowValues: []parser.QueryExpression{
					parser.RowValue{
						Value: parser.ValueList{
							Values: []parser.QueryExpression{
								parser.NewIntegerValue(1),
								parser.NewIntegerValue(2),
							},
						},
					},
					parser.RowValue{
						Value: parser.ValueList{
							Values: []parser.QueryExpression{
								parser.NewIntegerValue(1),
							},
						},
					},
				},
			},
			Operator: "<>",
		},
		Error: "row value should contain exactly 2 values",
	},
	{
		Name: "Any with Row Value and JsonQuery Field Length Error",
		Expr: parser.Any{
			LHS: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewStringValue("2"),
						parser.NewStringValue("str2"),
					},
				},
			},
			Values: parser.JsonQuery{
				Query:    parser.NewStringValue("key{}"),
				JsonText: parser.NewStringValue("{\"key\":{}}"),
			},
		},
		Error: "row value should contain exactly 2 values",
	},
	{
		Name: "All",
		Expr: parser.All{
			LHS: parser.NewIntegerValue(5),
			Values: parser.RowValue{
				Value: parser.Subquery{
					Query: parser.SelectQuery{
						SelectEntity: parser.SelectEntity{
							SelectClause: parser.SelectClause{
								Select: "select",
								Fields: []parser.QueryExpression{
									parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
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
			},
			Operator: ">",
		},
		Result: value.NewTernary(ternary.TRUE),
	},
	{
		Name: "All False",
		Expr: parser.All{
			LHS: parser.NewIntegerValue(-99),
			Values: parser.RowValue{
				Value: parser.Subquery{
					Query: parser.SelectQuery{
						SelectEntity: parser.SelectEntity{
							SelectClause: parser.SelectClause{
								Select: "select",
								Fields: []parser.QueryExpression{
									parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
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
			},
			Operator: ">",
		},
		Result: value.NewTernary(ternary.FALSE),
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
								Fields: []parser.QueryExpression{
									parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
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
			},
			Operator: ">",
		},
		Error: "field notexist does not exist",
	},
	{
		Name: "All Query Execution Error",
		Expr: parser.All{
			LHS: parser.NewIntegerValue(5),
			Values: parser.RowValue{
				Value: parser.Subquery{
					Query: parser.SelectQuery{
						SelectEntity: parser.SelectEntity{
							SelectClause: parser.SelectClause{
								Select: "select",
								Fields: []parser.QueryExpression{
									parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
								},
							},
							FromClause: parser.FromClause{
								Tables: []parser.QueryExpression{
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
		Error: "field notexist does not exist",
	},
	{
		Name: "All Row Value Select Field Not Match Error",
		Expr: parser.All{
			LHS: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewIntegerValue(1),
						parser.NewIntegerValue(2),
					},
				},
			},
			Values: parser.Subquery{
				Query: parser.SelectQuery{
					SelectEntity: parser.SelectEntity{
						SelectClause: parser.SelectClause{
							Select: "select",
							Fields: []parser.QueryExpression{
								parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
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
			Operator: ">",
		},
		Error: "select query should return exactly 2 fields",
	},
	{
		Name: "All Row Value Length Not Match Error",
		Expr: parser.All{
			LHS: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewIntegerValue(1),
						parser.NewIntegerValue(2),
					},
				},
			},
			Values: parser.RowValueList{
				RowValues: []parser.QueryExpression{
					parser.RowValue{
						Value: parser.ValueList{
							Values: []parser.QueryExpression{
								parser.NewIntegerValue(1),
								parser.NewIntegerValue(2),
							},
						},
					},
					parser.RowValue{
						Value: parser.ValueList{
							Values: []parser.QueryExpression{
								parser.NewIntegerValue(1),
							},
						},
					},
				},
			},
			Operator: "=",
		},
		Error: "row value should contain exactly 2 values",
	},
	{
		Name: "All with Row Value and JsonQuery Field Length Error",
		Expr: parser.All{
			LHS: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewStringValue("2"),
						parser.NewStringValue("str2"),
					},
				},
			},
			Values: parser.JsonQuery{
				Query:    parser.NewStringValue("key{}"),
				JsonText: parser.NewStringValue("{\"key\":{}}"),
			},
		},
		Error: "row value should contain exactly 2 values",
	},
	{
		Name: "Like",
		Expr: parser.Like{
			LHS:      parser.NewStringValue("abcdefg"),
			Pattern:  parser.NewStringValue("_bc%"),
			Negation: parser.Token{Token: parser.NOT, Literal: "not"},
		},
		Result: value.NewTernary(ternary.FALSE),
	},
	{
		Name: "Like LHS Error",
		Expr: parser.Like{
			LHS:      parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			Pattern:  parser.NewStringValue("_bc%"),
			Negation: parser.Token{Token: parser.NOT, Literal: "not"},
		},
		Error: "field notexist does not exist",
	},
	{
		Name: "Like Pattern Error",
		Expr: parser.Like{
			LHS:      parser.NewStringValue("abcdefg"),
			Pattern:  parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			Negation: parser.Token{Token: parser.NOT, Literal: "not"},
		},
		Error: "field notexist does not exist",
	},
	{
		Name: "Exists",
		Filter: &Filter{
			records: []filterRecord{
				{
					view: &View{
						Header: NewHeaderWithId("table2", []string{"column3", "column4"}),
						RecordSet: []Record{
							NewRecordWithId(1, []value.Primary{
								value.NewInteger(1),
								value.NewString("str2"),
							}),
						},
						Tx: TestTx,
					},
					recordIndex: 0,
				},
			},
		},
		Expr: parser.Exists{
			Query: parser.Subquery{
				Query: parser.SelectQuery{
					SelectEntity: parser.SelectEntity{
						SelectClause: parser.SelectClause{
							Select: "select",
							Fields: []parser.QueryExpression{
								parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
							},
						},
						FromClause: parser.FromClause{
							Tables: []parser.QueryExpression{
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
		Result: value.NewTernary(ternary.TRUE),
	},
	{
		Name: "Exists No Record",
		Expr: parser.Exists{
			Query: parser.Subquery{
				Query: parser.SelectQuery{
					SelectEntity: parser.SelectEntity{
						SelectClause: parser.SelectClause{
							Select: "select",
							Fields: []parser.QueryExpression{
								parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
							},
						},
						FromClause: parser.FromClause{
							Tables: []parser.QueryExpression{
								parser.Table{Object: parser.Identifier{Literal: "table1"}},
							},
						},
						WhereClause: parser.WhereClause{
							Filter: parser.NewTernaryValue(ternary.FALSE),
						},
					},
				},
			},
		},
		Result: value.NewTernary(ternary.FALSE),
	},
	{
		Name: "Exists Query Execution Error",
		Expr: parser.Exists{
			Query: parser.Subquery{
				Query: parser.SelectQuery{
					SelectEntity: parser.SelectEntity{
						SelectClause: parser.SelectClause{
							Select: "select",
							Fields: []parser.QueryExpression{
								parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
							},
						},
						FromClause: parser.FromClause{
							Tables: []parser.QueryExpression{
								parser.Table{Object: parser.Identifier{Literal: "table1"}},
							},
						},
						WhereClause: parser.WhereClause{
							Filter: parser.Comparison{
								LHS:      parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
								RHS:      parser.NewStringValue("str2"),
								Operator: "=",
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
		Filter: &Filter{
			records: []filterRecord{
				{
					view: &View{
						Header: NewHeaderWithId("table2", []string{"column3", "column4"}),
						RecordSet: []Record{
							NewRecordWithId(1, []value.Primary{
								value.NewInteger(1),
								value.NewString("str2"),
							}),
						},
						Tx: TestTx,
					},
					recordIndex: 0,
				},
			},
		},
		Expr: parser.Subquery{
			Query: parser.SelectQuery{
				SelectEntity: parser.SelectEntity{
					SelectClause: parser.SelectClause{
						Select: "select",
						Fields: []parser.QueryExpression{
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
						},
					},
					FromClause: parser.FromClause{
						Tables: []parser.QueryExpression{
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
					Value: parser.NewIntegerValue(1),
				},
			},
		},
		Result: value.NewString("2"),
	},
	{
		Name: "Subquery No Record",
		Expr: parser.Subquery{
			Query: parser.SelectQuery{
				SelectEntity: parser.SelectEntity{
					SelectClause: parser.SelectClause{
						Select: "select",
						Fields: []parser.QueryExpression{
							parser.Field{Object: parser.NewIntegerValue(1)},
						},
					},
					FromClause: parser.FromClause{
						Tables: []parser.QueryExpression{
							parser.Table{Object: parser.Identifier{Literal: "table1"}},
						},
					},
					WhereClause: parser.WhereClause{
						Filter: parser.NewTernaryValue(ternary.FALSE),
					},
				},
				LimitClause: parser.LimitClause{
					Value: parser.NewIntegerValue(1),
				},
			},
		},
		Result: value.NewNull(),
	},
	{
		Name: "Subquery Execution Error",
		Expr: parser.Subquery{
			Query: parser.SelectQuery{
				SelectEntity: parser.SelectEntity{
					SelectClause: parser.SelectClause{
						Select: "select",
						Fields: []parser.QueryExpression{
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}}},
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
		Error: "field notexist does not exist",
	},
	{
		Name: "Subquery Too Many RecordSet Error",
		Expr: parser.Subquery{
			Query: parser.SelectQuery{
				SelectEntity: parser.SelectEntity{
					SelectClause: parser.SelectClause{
						Select: "select",
						Fields: []parser.QueryExpression{
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
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
		Error: "subquery returns too many records, should return only one record",
	},
	{
		Name: "Subquery Too Many Fields Error",
		Expr: parser.Subquery{
			Query: parser.SelectQuery{
				SelectEntity: parser.SelectEntity{
					SelectClause: parser.SelectClause{
						Select: "select",
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
				LimitClause: parser.LimitClause{
					Value: parser.NewIntegerValue(1),
				},
			},
		},
		Error: "subquery returns too many fields, should return only one field",
	},
	{
		Name: "Function",
		Expr: parser.Function{
			Name: "coalesce",
			Args: []parser.QueryExpression{
				parser.NewNullValue(),
				parser.NewStringValue("str"),
			},
		},
		Result: value.NewString("str"),
	},
	{
		Name: "Function Now",
		Expr: parser.Function{
			Name: "now",
		},
		Result: value.NewDatetime(NowForTest),
	},
	{
		Name: "User Defined Function",
		Filter: &Filter{
			functions: UserDefinedFunctionScopes{
				UserDefinedFunctionMap{
					"USERFUNC": &UserDefinedFunction{
						Name: parser.Identifier{Literal: "userfunc"},
						Parameters: []parser.Variable{
							{Name: "arg1"},
						},
						RequiredArgs: 1,
						Statements: []parser.Statement{
							parser.Return{Value: parser.Variable{Name: "arg1"}},
						},
					},
				},
			},
		},
		Expr: parser.Function{
			Name: "userfunc",
			Args: []parser.QueryExpression{
				parser.NewIntegerValue(1),
			},
		},
		Result: value.NewInteger(1),
	},
	{
		Name: "User Defined Function Argument Length Error",
		Filter: &Filter{
			functions: UserDefinedFunctionScopes{
				UserDefinedFunctionMap{
					"USERFUNC": &UserDefinedFunction{
						Name: parser.Identifier{Literal: "userfunc"},
						Parameters: []parser.Variable{
							{Name: "arg1"},
						},
						RequiredArgs: 1,
						Statements: []parser.Statement{
							parser.Return{Value: parser.Variable{Name: "arg1"}},
						},
					},
				},
			},
		},
		Expr: parser.Function{
			Name: "userfunc",
			Args: []parser.QueryExpression{},
		},
		Error: "function userfunc takes exactly 1 argument",
	},
	{
		Name: "Function Not Exist Error",
		Expr: parser.Function{
			Name: "notexist",
			Args: []parser.QueryExpression{
				parser.NewNullValue(),
				parser.NewStringValue("str"),
			},
		},
		Error: "function notexist does not exist",
	},
	{
		Name: "Function Evaluate Error",
		Expr: parser.Function{
			Name: "coalesce",
			Args: []parser.QueryExpression{
				parser.NewNullValue(),
				parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			},
		},
		Error: "field notexist does not exist",
	},
	{
		Name: "Aggregate Function",
		Filter: &Filter{
			records: []filterRecord{
				{
					view: &View{
						Header: NewHeaderWithId("table1", []string{"column1", "column2"}),
						RecordSet: []Record{
							{
								NewGroupCell([]value.Primary{
									value.NewInteger(1),
									value.NewInteger(2),
									value.NewInteger(3),
								}),
								NewGroupCell([]value.Primary{
									value.NewInteger(1),
									value.NewNull(),
									value.NewInteger(3),
								}),
								NewGroupCell([]value.Primary{
									value.NewString("str1"),
									value.NewString("str2"),
									value.NewString("str3"),
								}),
							},
						},
						isGrouped: true,
						Tx:        TestTx,
					},
					recordIndex: 0,
				},
			},
		},
		Expr: parser.AggregateFunction{
			Name:     "avg",
			Distinct: parser.Token{},
			Args: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
			},
		},
		Result: value.NewInteger(2),
	},
	{
		Name: "Aggregate Function Argument Length Error",
		Filter: &Filter{
			records: []filterRecord{
				{
					view: &View{
						Header: NewHeader("table1", []string{"column1", "column2"}),
						RecordSet: []Record{
							{
								NewGroupCell([]value.Primary{
									value.NewInteger(1),
									value.NewNull(),
									value.NewInteger(3),
								}),
								NewGroupCell([]value.Primary{
									value.NewString("str1"),
									value.NewString("str2"),
									value.NewString("str3"),
								}),
							},
						},
						isGrouped: true,
						Tx:        TestTx,
					},
					recordIndex: 0,
				},
			},
		},
		Expr: parser.AggregateFunction{
			Name:     "avg",
			Distinct: parser.Token{},
		},
		Error: "function avg takes exactly 1 argument",
	},
	{
		Name: "Aggregate Function Not Grouped Error",
		Filter: &Filter{
			records: []filterRecord{
				{
					view: &View{
						Header: NewHeader("table1", []string{"column1", "column2"}),
						RecordSet: []Record{
							NewRecordWithId(1, []value.Primary{
								value.NewInteger(1),
								value.NewString("str2"),
							}),
						},
						Tx: TestTx,
					},
					recordIndex: 0,
				},
			},
		},
		Expr: parser.AggregateFunction{
			Name:     "avg",
			Distinct: parser.Token{},
			Args: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
			},
		},
		Error: "function avg cannot aggregate not grouping records",
	},
	{
		Name: "Aggregate Function Nested Error",
		Filter: &Filter{
			records: []filterRecord{
				{
					view: &View{
						Header: NewHeader("table1", []string{"column1", "column2"}),
						RecordSet: []Record{
							{
								NewGroupCell([]value.Primary{
									value.NewInteger(1),
									value.NewNull(),
									value.NewInteger(3),
								}),
								NewGroupCell([]value.Primary{
									value.NewString("str1"),
									value.NewString("str2"),
									value.NewString("str3"),
								}),
							},
						},
						isGrouped: true,
						Tx:        TestTx,
					},
					recordIndex: 0,
				},
			},
		},
		Expr: parser.AggregateFunction{
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
		Error: "aggregate functions are nested at avg(avg(column1))",
	},
	{
		Name: "Aggregate Function Count With AllColumns",
		Filter: &Filter{
			records: []filterRecord{
				{
					view: &View{
						Header: NewHeader("table1", []string{"column1", "column2"}),
						RecordSet: []Record{
							{
								NewGroupCell([]value.Primary{
									value.NewInteger(1),
									value.NewNull(),
									value.NewInteger(3),
								}),
								NewGroupCell([]value.Primary{
									value.NewString("str1"),
									value.NewString("str2"),
									value.NewString("str3"),
								}),
							},
						},
						isGrouped: true,
						Tx:        TestTx,
					},
					recordIndex: 0,
				},
			},
		},
		Expr: parser.AggregateFunction{
			Name:     "count",
			Distinct: parser.Token{},
			Args: []parser.QueryExpression{
				parser.AllColumns{},
			},
		},
		Result: value.NewInteger(3),
	},
	{
		Name:   "Aggregate Function As a Statement Error",
		Filter: &Filter{},
		Expr: parser.AggregateFunction{
			Name:     "avg",
			Distinct: parser.Token{},
			Args: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
			},
		},
		Result: value.NewNull(),
	},
	{
		Name: "Aggregate Function User Defined",
		Filter: &Filter{
			records: []filterRecord{
				{
					view: &View{
						Header: NewHeader("table1", []string{"column1", "column2"}),
						RecordSet: []Record{
							{
								NewGroupCell([]value.Primary{
									value.NewInteger(1),
									value.NewNull(),
									value.NewInteger(3),
									value.NewInteger(3),
								}),
								NewGroupCell([]value.Primary{
									value.NewString("str1"),
									value.NewString("str2"),
									value.NewString("str3"),
									value.NewString("str4"),
								}),
							},
						},
						isGrouped: true,
						Tx:        TestTx,
					},
					recordIndex: 0,
				},
			},
			functions: UserDefinedFunctionScopes{
				{
					"USERAGGFUNC": &UserDefinedFunction{
						Name:        parser.Identifier{Literal: "useraggfunc"},
						IsAggregate: true,
						Cursor:      parser.Identifier{Literal: "column1"},
						Parameters: []parser.Variable{
							{Name: "default"},
						},
						RequiredArgs: 1,
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
								Cursor: parser.Identifier{Literal: "column1"},
								Statements: []parser.Statement{
									parser.If{
										Condition: parser.Is{
											LHS: parser.Variable{Name: "fetch"},
											RHS: parser.NewNullValue(),
										},
										Statements: []parser.Statement{
											parser.FlowControl{Token: parser.CONTINUE},
										},
									},
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
											Operator: '*',
										},
									},
								},
							},

							parser.If{
								Condition: parser.Is{
									LHS: parser.Variable{Name: "value"},
									RHS: parser.NewNullValue(),
								},
								Statements: []parser.Statement{
									parser.VariableSubstitution{
										Variable: parser.Variable{Name: "value"},
										Value:    parser.Variable{Name: "default"},
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
		},
		Expr: parser.AggregateFunction{
			Name:     "useraggfunc",
			Distinct: parser.Token{Token: parser.DISTINCT, Literal: "distinct"},
			Args: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
				parser.NewIntegerValue(0),
			},
		},
		Result: value.NewInteger(3),
	},
	{
		Name: "Aggregate Function User Defined Argument Length Error",
		Filter: &Filter{
			records: []filterRecord{
				{
					view: &View{
						Header: NewHeader("table1", []string{"column1", "column2"}),
						RecordSet: []Record{
							{
								NewGroupCell([]value.Primary{
									value.NewInteger(1),
									value.NewNull(),
									value.NewInteger(3),
									value.NewInteger(3),
								}),
								NewGroupCell([]value.Primary{
									value.NewString("str1"),
									value.NewString("str2"),
									value.NewString("str3"),
									value.NewString("str4"),
								}),
							},
						},
						isGrouped: true,
						Tx:        TestTx,
					},
					recordIndex: 0,
				},
			},
			functions: UserDefinedFunctionScopes{
				{
					"USERAGGFUNC": &UserDefinedFunction{
						Name:        parser.Identifier{Literal: "useraggfunc"},
						IsAggregate: true,
						Cursor:      parser.Identifier{Literal: "column1"},
						Parameters: []parser.Variable{
							{Name: "default"},
						},
						Statements: []parser.Statement{
							parser.Return{
								Value: parser.Variable{Name: "value"},
							},
						},
					},
				},
			},
		},
		Expr: parser.AggregateFunction{
			Name:     "useraggfunc",
			Distinct: parser.Token{Token: parser.DISTINCT, Literal: "distinct"},
			Args: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
			},
		},
		Error: "function useraggfunc takes exactly 2 arguments",
	},
	{
		Name: "Aggregate Function User Defined Argument Evaluation Error",
		Filter: &Filter{
			records: []filterRecord{
				{
					view: &View{
						Header: NewHeader("table1", []string{"column1", "column2"}),
						RecordSet: []Record{
							{
								NewGroupCell([]value.Primary{
									value.NewInteger(1),
									value.NewNull(),
									value.NewInteger(3),
									value.NewInteger(3),
								}),
								NewGroupCell([]value.Primary{
									value.NewString("str1"),
									value.NewString("str2"),
									value.NewString("str3"),
									value.NewString("str4"),
								}),
							},
						},
						isGrouped: true,
						Tx:        TestTx,
					},
					recordIndex: 0,
				},
			},
			functions: UserDefinedFunctionScopes{
				{
					"USERAGGFUNC": &UserDefinedFunction{
						Name:        parser.Identifier{Literal: "useraggfunc"},
						IsAggregate: true,
						Cursor:      parser.Identifier{Literal: "column1"},
						Parameters: []parser.Variable{
							{Name: "default"},
						},
						Statements: []parser.Statement{
							parser.Return{
								Value: parser.Variable{Name: "value"},
							},
						},
					},
				},
			},
		},
		Expr: parser.AggregateFunction{
			Name:     "useraggfunc",
			Distinct: parser.Token{Token: parser.DISTINCT, Literal: "distinct"},
			Args: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
				parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			},
		},
		Error: "field notexist does not exist",
	},
	{
		Name: "Aggregate Function Execute User Defined Aggregate Function Passed As Scala Function",
		Filter: &Filter{
			records: []filterRecord{
				{
					view: &View{
						Header: NewHeader("table1", []string{"column1", "column2"}),
						RecordSet: []Record{
							{
								NewGroupCell([]value.Primary{
									value.NewInteger(1),
									value.NewNull(),
									value.NewInteger(3),
								}),
								NewGroupCell([]value.Primary{
									value.NewString("str1"),
									value.NewString("str2"),
									value.NewString("str3"),
								}),
							},
						},
						isGrouped: true,
						Tx:        TestTx,
					},
					recordIndex: 0,
				},
			},
			functions: UserDefinedFunctionScopes{
				{
					"USERAGGFUNC": &UserDefinedFunction{
						Name:        parser.Identifier{Literal: "useraggfunc"},
						IsAggregate: true,
						Cursor:      parser.Identifier{Literal: "column1"},
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
								Cursor: parser.Identifier{Literal: "column1"},
								Statements: []parser.Statement{
									parser.If{
										Condition: parser.Is{
											LHS: parser.Variable{Name: "fetch"},
											RHS: parser.NewNullValue(),
										},
										Statements: []parser.Statement{
											parser.FlowControl{Token: parser.CONTINUE},
										},
									},
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
											Operator: '*',
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
		},
		Expr: parser.Function{
			Name: "useraggfunc",
			Args: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
			},
		},
		Result: value.NewInteger(3),
	},
	{
		Name: "Aggregate Function Execute User Defined Aggregate Function Undeclared Error",
		Filter: &Filter{
			records: []filterRecord{
				{
					view: &View{
						Header: NewHeader("table1", []string{"column1", "column2"}),
						RecordSet: []Record{
							{
								NewGroupCell([]value.Primary{
									value.NewInteger(1),
									value.NewNull(),
									value.NewInteger(3),
								}),
								NewGroupCell([]value.Primary{
									value.NewString("str1"),
									value.NewString("str2"),
									value.NewString("str3"),
								}),
							},
						},
						isGrouped: true,
						Tx:        TestTx,
					},
					recordIndex: 0,
				},
			},
			functions: UserDefinedFunctionScopes{
				{
					"USERAGGFUNC": &UserDefinedFunction{
						Name:        parser.Identifier{Literal: "useraggfunc"},
						IsAggregate: true,
						Cursor:      parser.Identifier{Literal: "column1"},
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
								Cursor: parser.Identifier{Literal: "column1"},
								Statements: []parser.Statement{
									parser.If{
										Condition: parser.Is{
											LHS: parser.Variable{Name: "fetch"},
											RHS: parser.NewNullValue(),
										},
										Statements: []parser.Statement{
											parser.FlowControl{Token: parser.CONTINUE},
										},
									},
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
											Operator: '*',
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
		},
		Expr: parser.AggregateFunction{
			Name: "undefined",
			Args: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
			},
		},
		Error: "function undefined does not exist",
	},
	{
		Name: "ListAgg Function",
		Filter: &Filter{
			records: []filterRecord{
				{
					view: &View{
						Header: NewHeaderWithId("table1", []string{"column1", "column2"}),
						RecordSet: []Record{
							{
								NewGroupCell([]value.Primary{
									value.NewInteger(1),
									value.NewInteger(2),
									value.NewInteger(3),
									value.NewInteger(4),
								}),
								NewGroupCell([]value.Primary{
									value.NewInteger(1),
									value.NewInteger(2),
									value.NewInteger(3),
									value.NewInteger(4),
								}),
								NewGroupCell([]value.Primary{
									value.NewString("str2"),
									value.NewString("str1"),
									value.NewNull(),
									value.NewString("str2"),
								}),
							},
						},
						Filter:    NewFilter(TestTx),
						isGrouped: true,
						Tx:        TestTx,
					},
					recordIndex: 0,
				},
			},
		},
		Expr: parser.ListFunction{
			Name:     "listagg",
			Distinct: parser.Token{Token: parser.DISTINCT, Literal: "distinct"},
			Args: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
				parser.NewStringValue(","),
			},
			OrderBy: parser.OrderByClause{
				Items: []parser.QueryExpression{
					parser.OrderItem{Value: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
				},
			},
		},
		Result: value.NewString("str1,str2"),
	},
	{
		Name: "ListAgg Function Null",
		Filter: &Filter{
			records: []filterRecord{
				{
					view: &View{
						Header: NewHeaderWithId("table1", []string{"column1", "column2"}),
						RecordSet: []Record{
							{
								NewGroupCell([]value.Primary{
									value.NewInteger(1),
									value.NewInteger(2),
									value.NewInteger(3),
									value.NewInteger(4),
								}),
								NewGroupCell([]value.Primary{
									value.NewInteger(1),
									value.NewInteger(2),
									value.NewInteger(3),
									value.NewInteger(4),
								}),
								NewGroupCell([]value.Primary{
									value.NewNull(),
									value.NewNull(),
									value.NewNull(),
									value.NewNull(),
								}),
							},
						},
						Filter:    NewFilter(TestTx),
						isGrouped: true,
						Tx:        TestTx,
					},
					recordIndex: 0,
				},
			},
		},
		Expr: parser.ListFunction{
			Distinct: parser.Token{Token: parser.DISTINCT, Literal: "distinct"},
			Args: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
				parser.NewStringValue(","),
			},
		},
		Result: value.NewNull(),
	},
	{
		Name: "ListAgg Function Argument Length Error",
		Filter: &Filter{
			records: []filterRecord{
				{
					view: &View{
						Header: NewHeaderWithId("table1", []string{"column1", "column2"}),
						RecordSet: []Record{
							NewRecordWithId(1, []value.Primary{
								value.NewInteger(1),
								value.NewString("str2"),
							}),
						},
						Filter: NewFilter(TestTx),
						Tx:     TestTx,
					},
					recordIndex: 0,
				},
			},
		},
		Expr: parser.ListFunction{
			Name:     "listagg",
			Distinct: parser.Token{Token: parser.DISTINCT, Literal: "distinct"},
			OrderBy: parser.OrderByClause{
				Items: []parser.QueryExpression{
					parser.OrderItem{Value: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
				},
			},
		},
		Error: "function listagg takes 1 or 2 arguments",
	},
	{
		Name: "ListAgg Function Not Grouped Error",
		Filter: &Filter{
			records: []filterRecord{
				{
					view: &View{
						Header: NewHeaderWithId("table1", []string{"column1", "column2"}),
						RecordSet: []Record{
							NewRecordWithId(1, []value.Primary{
								value.NewInteger(1),
								value.NewString("str2"),
							}),
						},
						Filter: NewFilter(TestTx),
						Tx:     TestTx,
					},
					recordIndex: 0,
				},
			},
		},
		Expr: parser.ListFunction{
			Name:     "listagg",
			Distinct: parser.Token{Token: parser.DISTINCT, Literal: "distinct"},
			Args: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
				parser.NewStringValue(","),
			},
			OrderBy: parser.OrderByClause{
				Items: []parser.QueryExpression{
					parser.OrderItem{Value: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
				},
			},
		},
		Error: "function listagg cannot aggregate not grouping records",
	},
	{
		Name: "ListAgg Function Sort Error",
		Filter: &Filter{
			records: []filterRecord{
				{
					view: &View{
						Header: NewHeaderWithId("table1", []string{"column1", "column2"}),
						RecordSet: []Record{
							{
								NewGroupCell([]value.Primary{
									value.NewInteger(1),
									value.NewInteger(2),
									value.NewInteger(3),
									value.NewInteger(4),
								}),
								NewGroupCell([]value.Primary{
									value.NewString("str2"),
									value.NewString("str1"),
									value.NewNull(),
									value.NewString("str2"),
								}),
							},
						},
						Filter:    NewFilter(TestTx),
						isGrouped: true,
						Tx:        TestTx,
					},
					recordIndex: 0,
				},
			},
		},
		Expr: parser.ListFunction{
			Name:     "listagg",
			Distinct: parser.Token{Token: parser.DISTINCT, Literal: "distinct"},
			Args: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
				parser.NewStringValue(","),
			},
			OrderBy: parser.OrderByClause{
				Items: []parser.QueryExpression{
					parser.OrderItem{Value: parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}}},
				},
			},
		},
		Error: "field notexist does not exist",
	},
	{
		Name: "ListAgg Function Nested Error",
		Filter: &Filter{
			records: []filterRecord{
				{
					view: &View{
						Header: NewHeaderWithId("table1", []string{"column1", "column2"}),
						RecordSet: []Record{
							{
								NewGroupCell([]value.Primary{
									value.NewInteger(1),
									value.NewNull(),
									value.NewInteger(3),
								}),
								NewGroupCell([]value.Primary{
									value.NewString("str1"),
									value.NewString("str2"),
									value.NewString("str3"),
								}),
							},
						},
						Filter:    NewFilter(TestTx),
						isGrouped: true,
						Tx:        TestTx,
					},
					recordIndex: 0,
				},
			},
		},
		Expr: parser.ListFunction{
			Name: "listagg",
			Args: []parser.QueryExpression{
				parser.AggregateFunction{
					Name:     "avg",
					Distinct: parser.Token{},
					Args: []parser.QueryExpression{
						parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
					},
				},
			},
		},
		Error: "aggregate functions are nested at listagg(avg(column1))",
	},
	{
		Name: "ListAgg Function Second Argument Evaluation Error",
		Filter: &Filter{
			records: []filterRecord{
				{
					view: &View{
						Header: NewHeaderWithId("table1", []string{"column1", "column2"}),
						RecordSet: []Record{
							{
								NewGroupCell([]value.Primary{
									value.NewInteger(1),
									value.NewNull(),
									value.NewInteger(3),
								}),
								NewGroupCell([]value.Primary{
									value.NewString("str1"),
									value.NewString("str2"),
									value.NewString("str3"),
								}),
							},
						},
						Filter:    NewFilter(TestTx),
						isGrouped: true,
						Tx:        TestTx,
					},
					recordIndex: 0,
				},
			},
		},
		Expr: parser.ListFunction{
			Name: "listagg",
			Args: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
			},
		},
		Error: "the second argument must be a string for function listagg",
	},
	{
		Name: "ListAgg Function Second Argument Not String Error",
		Filter: &Filter{
			records: []filterRecord{
				{
					view: &View{
						Header: NewHeaderWithId("table1", []string{"column1", "column2"}),
						RecordSet: []Record{
							{
								NewGroupCell([]value.Primary{
									value.NewInteger(1),
									value.NewNull(),
									value.NewInteger(3),
								}),
								NewGroupCell([]value.Primary{
									value.NewString("str1"),
									value.NewString("str2"),
									value.NewString("str3"),
								}),
							},
						},
						Filter:    NewFilter(TestTx),
						isGrouped: true,
						Tx:        TestTx,
					},
					recordIndex: 0,
				},
			},
		},
		Expr: parser.ListFunction{
			Name: "listagg",
			Args: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
				parser.NewNullValue(),
			},
		},
		Error: "the second argument must be a string for function listagg",
	},
	{
		Name:   "ListAgg Function As a Statement Error",
		Filter: &Filter{},
		Expr: parser.ListFunction{
			Name:     "listagg",
			Distinct: parser.Token{Token: parser.DISTINCT, Literal: "distinct"},
			Args: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
				parser.NewStringValue(","),
			},
			OrderBy: parser.OrderByClause{
				Items: []parser.QueryExpression{
					parser.OrderItem{Value: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
				},
			},
		},
		Result: value.NewNull(),
	},
	{
		Name: "JsonAgg Function",
		Filter: &Filter{
			records: []filterRecord{
				{
					view: &View{
						Header: NewHeaderWithId("table1", []string{"column1", "column2"}),
						RecordSet: []Record{
							{
								NewGroupCell([]value.Primary{
									value.NewInteger(1),
									value.NewInteger(2),
									value.NewInteger(3),
									value.NewInteger(4),
								}),
								NewGroupCell([]value.Primary{
									value.NewInteger(1),
									value.NewInteger(2),
									value.NewInteger(3),
									value.NewInteger(4),
								}),
								NewGroupCell([]value.Primary{
									value.NewString("str2"),
									value.NewString("str1"),
									value.NewNull(),
									value.NewString("str2"),
								}),
							},
						},
						Filter:    NewFilter(TestTx),
						isGrouped: true,
						Tx:        TestTx,
					},
					recordIndex: 0,
				},
			},
		},
		Expr: parser.ListFunction{
			Name:     "json_agg",
			Distinct: parser.Token{Token: parser.DISTINCT, Literal: "distinct"},
			Args: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
			},
			OrderBy: parser.OrderByClause{
				Items: []parser.QueryExpression{
					parser.OrderItem{Value: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
				},
			},
		},
		Result: value.NewString("[null,\"str1\",\"str2\"]"),
	},
	{
		Name: "JsonAgg Function Arguments Error",
		Filter: &Filter{
			records: []filterRecord{
				{
					view: &View{
						Header: NewHeaderWithId("table1", []string{"column1", "column2"}),
						RecordSet: []Record{
							{
								NewGroupCell([]value.Primary{
									value.NewInteger(1),
									value.NewInteger(2),
								}),
								NewGroupCell([]value.Primary{
									value.NewInteger(1),
									value.NewInteger(2),
								}),
								NewGroupCell([]value.Primary{
									value.NewString("str2"),
									value.NewString("str1"),
								}),
							},
						},
						Filter:    NewFilter(TestTx),
						isGrouped: true,
						Tx:        TestTx,
					},
					recordIndex: 0,
				},
			},
		},
		Expr: parser.ListFunction{
			Name:     "json_agg",
			Distinct: parser.Token{Token: parser.DISTINCT, Literal: "distinct"},
			Args:     []parser.QueryExpression{},
			OrderBy: parser.OrderByClause{
				Items: []parser.QueryExpression{
					parser.OrderItem{Value: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
				},
			},
		},
		Error: "function json_agg takes exactly 1 argument",
	},
	{
		Name: "CaseExpr Comparison",
		Expr: parser.CaseExpr{
			Value: parser.NewIntegerValue(2),
			When: []parser.QueryExpression{
				parser.CaseExprWhen{
					Condition: parser.NewIntegerValue(1),
					Result:    parser.NewStringValue("A"),
				},
				parser.CaseExprWhen{
					Condition: parser.NewIntegerValue(2),
					Result:    parser.NewStringValue("B"),
				},
			},
		},
		Result: value.NewString("B"),
	},
	{
		Name: "CaseExpr Filter",
		Expr: parser.CaseExpr{
			When: []parser.QueryExpression{
				parser.CaseExprWhen{
					Condition: parser.Comparison{
						LHS:      parser.NewIntegerValue(2),
						RHS:      parser.NewIntegerValue(1),
						Operator: "=",
					},
					Result: parser.NewStringValue("A"),
				},
				parser.CaseExprWhen{
					Condition: parser.Comparison{
						LHS:      parser.NewIntegerValue(2),
						RHS:      parser.NewIntegerValue(2),
						Operator: "=",
					},
					Result: parser.NewStringValue("B"),
				},
			},
		},
		Result: value.NewString("B"),
	},
	{
		Name: "CaseExpr Else",
		Expr: parser.CaseExpr{
			Value: parser.NewIntegerValue(0),
			When: []parser.QueryExpression{
				parser.CaseExprWhen{
					Condition: parser.NewIntegerValue(1),
					Result:    parser.NewStringValue("A"),
				},
				parser.CaseExprWhen{
					Condition: parser.NewIntegerValue(2),
					Result:    parser.NewStringValue("B"),
				},
			},
			Else: parser.CaseExprElse{
				Result: parser.NewStringValue("C"),
			},
		},
		Result: value.NewString("C"),
	},
	{
		Name: "CaseExpr No Else",
		Expr: parser.CaseExpr{
			Value: parser.NewIntegerValue(0),
			When: []parser.QueryExpression{
				parser.CaseExprWhen{
					Condition: parser.NewIntegerValue(1),
					Result:    parser.NewStringValue("A"),
				},
				parser.CaseExprWhen{
					Condition: parser.NewIntegerValue(2),
					Result:    parser.NewStringValue("B"),
				},
			},
		},
		Result: value.NewNull(),
	},
	{
		Name: "CaseExpr Value Error",
		Expr: parser.CaseExpr{
			Value: parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			When: []parser.QueryExpression{
				parser.CaseExprWhen{
					Condition: parser.NewIntegerValue(1),
					Result:    parser.NewStringValue("A"),
				},
				parser.CaseExprWhen{
					Condition: parser.NewIntegerValue(2),
					Result:    parser.NewStringValue("B"),
				},
			},
		},
		Error: "field notexist does not exist",
	},
	{
		Name: "CaseExpr When Condition Error",
		Expr: parser.CaseExpr{
			Value: parser.NewIntegerValue(2),
			When: []parser.QueryExpression{
				parser.CaseExprWhen{
					Condition: parser.NewIntegerValue(1),
					Result:    parser.NewStringValue("A"),
				},
				parser.CaseExprWhen{
					Condition: parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
					Result:    parser.NewStringValue("B"),
				},
			},
		},
		Error: "field notexist does not exist",
	},
	{
		Name: "CaseExpr When Result Error",
		Expr: parser.CaseExpr{
			Value: parser.NewIntegerValue(2),
			When: []parser.QueryExpression{
				parser.CaseExprWhen{
					Condition: parser.NewIntegerValue(1),
					Result:    parser.NewStringValue("A"),
				},
				parser.CaseExprWhen{
					Condition: parser.NewIntegerValue(2),
					Result:    parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
				},
			},
		},
		Error: "field notexist does not exist",
	},
	{
		Name: "CaseExpr Else Result Error",
		Expr: parser.CaseExpr{
			Value: parser.NewIntegerValue(0),
			When: []parser.QueryExpression{
				parser.CaseExprWhen{
					Condition: parser.NewIntegerValue(1),
					Result:    parser.NewStringValue("A"),
				},
				parser.CaseExprWhen{
					Condition: parser.NewIntegerValue(2),
					Result:    parser.NewStringValue("B"),
				},
			},
			Else: parser.CaseExprElse{
				Result: parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			},
		},
		Error: "field notexist does not exist",
	},
	{
		Name: "Logic AND",
		Expr: parser.Logic{
			LHS:      parser.NewTernaryValue(ternary.TRUE),
			RHS:      parser.NewTernaryValue(ternary.FALSE),
			Operator: parser.Token{Token: parser.AND, Literal: "and"},
		},
		Result: value.NewTernary(ternary.FALSE),
	},
	{
		Name: "Logic AND Decided with LHS",
		Expr: parser.Logic{
			LHS:      parser.NewTernaryValue(ternary.FALSE),
			RHS:      parser.NewTernaryValue(ternary.FALSE),
			Operator: parser.Token{Token: parser.AND, Literal: "and"},
		},
		Result: value.NewTernary(ternary.FALSE),
	},
	{
		Name: "Logic OR",
		Expr: parser.Logic{
			LHS:      parser.NewTernaryValue(ternary.FALSE),
			RHS:      parser.NewTernaryValue(ternary.TRUE),
			Operator: parser.Token{Token: parser.OR, Literal: "or"},
		},
		Result: value.NewTernary(ternary.TRUE),
	},
	{
		Name: "Logic OR Decided with LHS",
		Expr: parser.Logic{
			LHS:      parser.NewTernaryValue(ternary.TRUE),
			RHS:      parser.NewTernaryValue(ternary.FALSE),
			Operator: parser.Token{Token: parser.OR, Literal: "or"},
		},
		Result: value.NewTernary(ternary.TRUE),
	},
	{
		Name: "Logic LHS Error",
		Expr: parser.Logic{
			LHS:      parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			RHS:      parser.NewTernaryValue(ternary.FALSE),
			Operator: parser.Token{Token: parser.AND, Literal: "and"},
		},
		Error: "field notexist does not exist",
	},
	{
		Name: "Logic RHS Error",
		Expr: parser.Logic{
			LHS:      parser.NewTernaryValue(ternary.UNKNOWN),
			RHS:      parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			Operator: parser.Token{Token: parser.AND, Literal: "and"},
		},
		Error: "field notexist does not exist",
	},
	{
		Name: "UnaryLogic",
		Expr: parser.UnaryLogic{
			Operand:  parser.NewTernaryValue(ternary.FALSE),
			Operator: parser.Token{Token: parser.NOT, Literal: "not"},
		},
		Result: value.NewTernary(ternary.TRUE),
	},
	{
		Name: "UnaryLogic Operand Error",
		Expr: parser.UnaryLogic{
			Operand:  parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			Operator: parser.Token{Token: parser.NOT, Literal: "not"},
		},
		Error: "field notexist does not exist",
	},
	{
		Name: "Variable",
		Filter: NewFilterWithScopes(
			TestTx,
			[]VariableMap{
				GenerateVariableMap(map[string]value.Primary{
					"var1": value.NewInteger(1),
				}),
			},
			[]ViewMap{NewViewMap()},
			[]CursorMap{{}},
			[]UserDefinedFunctionMap{{}},
		),
		Expr: parser.Variable{
			Name: "var1",
		},
		Result: value.NewInteger(1),
	},
	{
		Name: "Environment Variable",
		Expr: parser.EnvironmentVariable{
			Name: "CSVQ_TEST_ENV",
		},
		Result: value.NewString("foo"),
	},
	{
		Name: "Runtime Information",
		Expr: parser.RuntimeInformation{
			Name: "version",
		},
		Result: value.NewString("v1.0.0"),
	},
	{
		Name: "Flag",
		Expr: parser.Flag{
			Name: "json_escape",
		},
		Result: value.NewString("BACKSLASH"),
	},
	{
		Name: "Flag Ivalid Flag Name Error",
		Expr: parser.Flag{
			Name: "invalid",
		},
		Error: "@@invalid is an unknown flag",
	},
	{
		Name: "Variable Undeclared Error",
		Expr: parser.Variable{
			Name: "undefined",
		},
		Error: "variable @undefined is undeclared",
	},
	{
		Name: "Variable Substitution",
		Filter: NewFilterWithScopes(
			TestTx,
			[]VariableMap{
				GenerateVariableMap(map[string]value.Primary{
					"var1": value.NewInteger(1),
				}),
			},
			[]ViewMap{NewViewMap()},
			[]CursorMap{{}},
			[]UserDefinedFunctionMap{{}},
		),
		Expr: parser.VariableSubstitution{
			Variable: parser.Variable{Name: "var1"},
			Value:    parser.NewIntegerValue(2),
		},
		Result: value.NewInteger(2),
	},
	{
		Name: "Variable Substitution Undeclared Error",
		Expr: parser.VariableSubstitution{
			Variable: parser.Variable{Name: "undefined"},
			Value:    parser.NewIntegerValue(2),
		},
		Error: "variable @undefined is undeclared",
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
		Result: value.NewTernary(ternary.FALSE),
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
		Result: value.NewTernary(ternary.TRUE),
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
		Error: "cursor notexist is undeclared",
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
		Error: "cursor notexist is undeclared",
	},
	{
		Name: "Cursor Attribute Count",
		Expr: parser.CursorAttrebute{
			CursorLit: "cursor",
			Cursor:    parser.Identifier{Literal: "cur"},
			Attrebute: parser.Token{Token: parser.COUNT, Literal: "count"},
		},
		Result: value.NewInteger(3),
	},
	{
		Name: "Cursor Attribute Count Error",
		Expr: parser.CursorAttrebute{
			CursorLit: "cursor",
			Cursor:    parser.Identifier{Literal: "notexist"},
			Attrebute: parser.Token{Token: parser.COUNT, Literal: "count"},
		},
		Error: "cursor notexist is undeclared",
	},
	{
		Name: "Placeholder Ordinal",
		Expr: parser.Placeholder{
			Literal: "?",
			Ordinal: 1,
		},
		ReplaceValues: &ReplaceValues{
			Values: []parser.QueryExpression{parser.NewIntegerValueFromString("1")},
			Names:  map[string]int{},
		},
		Result: value.NewInteger(1),
	},
	{
		Name: "Placeholder Named",
		Expr: parser.Placeholder{
			Literal: ":val",
			Ordinal: 1,
			Name:    "val",
		},
		ReplaceValues: &ReplaceValues{
			Values: []parser.QueryExpression{parser.NewIntegerValueFromString("1")},
			Names:  map[string]int{"val": 0},
		},
		Result: value.NewInteger(1),
	},
	{
		Name: "Placeholder Replace Values Not Exist Error",
		Expr: parser.Placeholder{
			Literal: "?",
			Ordinal: 1,
		},
		Error: "replace value for ?{1} is not specified",
	},
	{
		Name: "Placeholder Ordinal Replace Value Not Exist Error",
		Expr: parser.Placeholder{
			Literal: "?",
			Ordinal: 10,
		},
		ReplaceValues: &ReplaceValues{
			Values: []parser.QueryExpression{parser.NewIntegerValueFromString("1")},
		},
		Error: "replace value for ?{10} is not specified",
	},
	{
		Name: "Placeholder Named Replace Value Not Exist Error",
		Expr: parser.Placeholder{
			Literal: ":notexist",
			Ordinal: 1,
			Name:    "notexist",
		},
		ReplaceValues: &ReplaceValues{
			Values: []parser.QueryExpression{parser.NewIntegerValueFromString("1")},
			Names:  map[string]int{"val": 0},
		},
		Error: "replace value for :notexist is not specified",
	},
}

func TestFilter_Evaluate(t *testing.T) {
	defer func() {
		_ = TestTx.cachedViews.Clean(TestTx.FileContainer)
		initFlag(TestTx.Flags)
	}()

	TestTx.Flags.Repository = TestDataDir

	cursors := CursorMap{
		"CUR": &Cursor{
			query: selectQueryForCursorTest,
			mtx:   &sync.Mutex{},
		},
	}

	ctx := context.Background()
	_ = cursors.Open(ctx, NewFilter(TestTx), parser.Identifier{Literal: "cur"}, nil)
	_, _ = cursors.Fetch(parser.Identifier{Literal: "cur"}, parser.NEXT, 0)

	for _, v := range filterEvaluateTests {
		_ = TestTx.cachedViews.Clean(TestTx.FileContainer)

		if v.Filter == nil {
			v.Filter = NewFilter(TestTx)
		} else {
			v.Filter.tx = TestTx
		}
		for i := range v.Filter.records {
			v.Filter.records[i].view.Filter = v.Filter
		}

		v.Filter.cursors = append(v.Filter.cursors, cursors)

		evalCtx := ctx
		if v.ReplaceValues != nil {
			evalCtx = ContextForPreparedStatement(ctx, v.ReplaceValues)
		}
		result, err := v.Filter.Evaluate(evalCtx, v.Expr)
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

var filterEvaluateEmbeddedStringTests = []struct {
	Input  string
	Expect string
	Error  string
}{
	{
		Input:  "str",
		Expect: "str",
	},
	{
		Input:  "@var",
		Expect: "1",
	},
	{
		Input:  "@%CSVQ_TEST_FILTER",
		Expect: "FILTER_TEST",
	},
	{
		Input:  "@#version",
		Expect: "v1.0.0",
	},
	{
		Input:  "abc${}def",
		Expect: "abcdef",
	},
	{
		Input:  "abc${@var}def",
		Expect: "abc1def",
	},
	{
		Input: "@notexist",
		Error: "variable @notexist is undeclared",
	},
	{
		Input: "@#notexist",
		Error: "@#notexist is an unknown runtime information",
	},
	{
		Input: "abc${invalid expr}def",
		Error: "invalid expr [L:1 C:9] syntax error: unexpected token \"expr\"",
	},
	{
		Input: "abc${print 1;}def",
		Error: "'print 1;': cannot evaluate as a value",
	},
	{
		Input: "abc${print 1;print2;}def",
		Error: "print 1;print2; [L:1 C:15] syntax error: unexpected token \";\"",
	},
	{
		Input: "abc${@notexist}def",
		Error: "@notexist [L:1 C:1] variable @notexist is undeclared",
	},
}

func TestFilter_EvaluateEmbeddedString(t *testing.T) {
	filter := NewFilter(TestTx)
	_ = filter.variables[0].Add(parser.Variable{Name: "var"}, value.NewInteger(1))
	_ = os.Setenv("CSVQ_TEST_FILTER", "FILTER_TEST")

	for _, v := range filterEvaluateEmbeddedStringTests {
		result, err := filter.EvaluateEmbeddedString(context.Background(), v.Input)

		if err != nil {
			if len(v.Error) < 1 {
				t.Errorf("unexpected error %q for %q", err, v.Input)
			} else if err.Error() != v.Error {
				t.Errorf("error %q, want error %q for %q", err.Error(), v.Error, v.Input)
			}
			continue
		}
		if 0 < len(v.Error) {
			t.Errorf("no error, want error %q for %q", v.Error, v.Input)
			continue
		}

		if result != v.Expect {
			t.Errorf("result = %q, want %q for %q", result, v.Expect, v.Input)
		}
	}
}

func BenchmarkFilter_EvaluateCountAllColumns(b *testing.B) {
	ctx := context.Background()
	filter := GenerateBenchGroupedViewFilter()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = filter.Evaluate(ctx, parser.AggregateFunction{
			Name:     "count",
			Distinct: parser.Token{},
			Args: []parser.QueryExpression{
				parser.AllColumns{},
			},
		})
	}
}

func BenchmarkFilter_EvaluateCount(b *testing.B) {
	ctx := context.Background()
	filter := GenerateBenchGroupedViewFilter()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = filter.Evaluate(ctx, parser.AggregateFunction{
			Name:     "count",
			Distinct: parser.Token{},
			Args: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "c1"}},
			},
		})
	}
}

func BenchmarkFilter_EvaluateSingleThread(b *testing.B) {
	ctx := context.Background()

	for i := 0; i < b.N; i++ {
		filter := NewFilter(TestTx)

		for j := 0; j < 150; j++ {
			_, _ = filter.Evaluate(ctx, parser.Comparison{
				LHS:      parser.NewIntegerValue(1),
				RHS:      parser.NewStringValue("1"),
				Operator: "=",
			})
		}
	}
}

func BenchmarkFilter_EvaluateMultiThread(b *testing.B) {
	ctx := context.Background()

	for i := 0; i < b.N; i++ {
		wg := sync.WaitGroup{}
		for i := 0; i < 3; i++ {
			wg.Add(1)
			go func(thIdx int) {
				filter := NewFilter(TestTx)

				for j := 0; j < 50; j++ {
					_, _ = filter.Evaluate(ctx, parser.Comparison{
						LHS:      parser.NewIntegerValue(1),
						RHS:      parser.NewStringValue("1"),
						Operator: "=",
					})
				}
				wg.Done()
			}(i)
		}
		wg.Wait()
	}
}

var filterEvaluateFieldReferenceBenchFilter = &Filter{
	records: []filterRecord{
		{
			view: &View{
				Header: NewHeader("table1", []string{"column1", "column2", "column3"}),
				RecordSet: RecordSet{
					NewRecord([]value.Primary{
						value.NewInteger(1),
						value.NewInteger(1),
						value.NewInteger(1),
					}),
				},
			},
			recordIndex: 0,
		},
	},
}

var filterEvaluateFieldReferenceWithIndexCacheBenchFilter = &Filter{
	records: []filterRecord{
		{
			view: &View{
				Header: NewHeader("table1", []string{"column1", "column2", "column3"}),
				RecordSet: RecordSet{
					NewRecord([]value.Primary{
						value.NewInteger(1),
						value.NewInteger(1),
						value.NewInteger(1),
					}),
				},
			},
			recordIndex:           0,
			fieldReferenceIndices: make(map[string]int),
		},
	},
}

var filterEvaluateFieldReferenceBenchExpr = parser.FieldReference{
	Column: parser.Identifier{Literal: "column3"},
}

func BenchmarkFilter_EvaluateFieldReference(b *testing.B) {
	ctx := context.Background()
	filter := filterEvaluateFieldReferenceBenchFilter
	expr := filterEvaluateFieldReferenceBenchExpr
	for i := 0; i < b.N; i++ {
		_, _ = filter.Evaluate(ctx, expr)
	}
}

func BenchmarkFilter_EvaluateFieldReferenceWithIndexCache(b *testing.B) {
	ctx := context.Background()
	filter := filterEvaluateFieldReferenceWithIndexCacheBenchFilter
	expr := filterEvaluateFieldReferenceBenchExpr
	for i := 0; i < b.N; i++ {
		_, _ = filter.Evaluate(ctx, expr)
	}
}
