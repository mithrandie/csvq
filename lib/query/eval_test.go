package query

import (
	"context"
	"os"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/value"
	"github.com/mithrandie/ternary"
)

var evaluateTests = []struct {
	Name          string
	Scope         *ReferenceScope
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
		Scope: GenerateReferenceScope(nil, nil, time.Time{}, []ReferenceRecord{
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
				},
				recordIndex: 1,
				cache:       NewFieldIndexCache(10, LimitToUseFieldIndexSliceChache),
			},
		}),
		Expr:   parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
		Result: value.NewString("strstr"),
	},
	{
		Name: "FieldReference ColumnNotExist Error",
		Scope: GenerateReferenceScope(nil, nil, time.Time{}, []ReferenceRecord{
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
				},
				recordIndex: 1,
				cache:       NewFieldIndexCache(10, LimitToUseFieldIndexSliceChache),
			},
		}),
		Expr:  parser.FieldReference{Column: parser.Identifier{Literal: "column3"}},
		Error: "field column3 does not exist",
	},
	{
		Name: "FieldReference FieldAmbigous Error",
		Scope: GenerateReferenceScope(nil, nil, time.Time{}, []ReferenceRecord{
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
				},
				recordIndex: 1,
				cache:       NewFieldIndexCache(10, LimitToUseFieldIndexSliceChache),
			},
		}),
		Expr:  parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
		Error: "field column1 is ambiguous",
	},
	{
		Name: "FieldReference Not Group Key Error",
		Scope: GenerateReferenceScope(nil, nil, time.Time{}, []ReferenceRecord{
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
				},
				recordIndex: 0,
				cache:       NewFieldIndexCache(10, LimitToUseFieldIndexSliceChache),
			},
		}),
		Expr:  parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
		Error: "field column1 is not a group key",
	},
	{
		Name: "ColumnNumber",
		Scope: GenerateReferenceScope(nil, nil, time.Time{}, []ReferenceRecord{
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
				},
				recordIndex: 1,
				cache:       NewFieldIndexCache(10, LimitToUseFieldIndexSliceChache),
			},
		}),
		Expr:   parser.ColumnNumber{View: parser.Identifier{Literal: "table1"}, Number: value.NewInteger(2)},
		Result: value.NewString("strstr"),
	},
	{
		Name: "ColumnNumber Not Exist Error",
		Scope: GenerateReferenceScope(nil, nil, time.Time{}, []ReferenceRecord{
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
				},
				recordIndex: 1,
				cache:       NewFieldIndexCache(10, LimitToUseFieldIndexSliceChache),
			},
		}),
		Expr:  parser.ColumnNumber{View: parser.Identifier{Literal: "table1"}, Number: value.NewInteger(9)},
		Error: "field table1.9 does not exist",
	},
	{
		Name: "ColumnNumber Not Group Key Error",
		Scope: GenerateReferenceScope(nil, nil, time.Time{}, []ReferenceRecord{
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
				},
				recordIndex: 0,
				cache:       NewFieldIndexCache(10, LimitToUseFieldIndexSliceChache),
			},
		}),
		Expr:  parser.ColumnNumber{View: parser.Identifier{Literal: "table1"}, Number: value.NewInteger(1)},
		Error: "field table1.1 is not a group key",
	},
	{
		Name: "Arithmetic",
		Expr: parser.Arithmetic{
			LHS:      parser.NewIntegerValue(1),
			RHS:      parser.NewIntegerValue(2),
			Operator: parser.Token{Token: '+', Literal: "+"},
		},
		Result: value.NewInteger(3),
	},
	{
		Name: "Arithmetic LHS Error",
		Expr: parser.Arithmetic{
			LHS:      parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			RHS:      parser.NewIntegerValue(2),
			Operator: parser.Token{Token: '+', Literal: "+"},
		},
		Error: "field notexist does not exist",
	},
	{
		Name: "Arithmetic LHS Is Null",
		Expr: parser.Arithmetic{
			LHS:      parser.NewNullValue(),
			RHS:      parser.NewIntegerValue(2),
			Operator: parser.Token{Token: '+', Literal: "+"},
		},
		Result: value.NewNull(),
	},
	{
		Name: "Arithmetic RHS Error",
		Expr: parser.Arithmetic{
			LHS:      parser.NewIntegerValue(1),
			RHS:      parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			Operator: parser.Token{Token: '+', Literal: "+"},
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
			Operator: parser.Token{Token: '=', Literal: "="},
		},
		Result: value.NewTernary(ternary.FALSE),
	},
	{
		Name: "Comparison With Single RowValue",
		Expr: parser.Comparison{
			LHS:      parser.RowValue{Value: parser.ValueList{Values: []parser.QueryExpression{parser.NewStringValue("2")}}},
			RHS:      parser.NewIntegerValue(2),
			Operator: parser.Token{Token: '=', Literal: "="},
		},
		Result: value.NewTernary(ternary.TRUE),
	},
	{
		Name: "Comparison LHS Error",
		Expr: parser.Comparison{
			LHS:      parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			RHS:      parser.NewIntegerValue(2),
			Operator: parser.Token{Token: '=', Literal: "="},
		},
		Error: "field notexist does not exist",
	},
	{
		Name: "Comparison LHS Is Null",
		Expr: parser.Comparison{
			LHS:      parser.NewNullValue(),
			RHS:      parser.NewIntegerValue(2),
			Operator: parser.Token{Token: '=', Literal: "="},
		},
		Result: value.NewTernary(ternary.UNKNOWN),
	},
	{
		Name: "Comparison RHS Error",
		Expr: parser.Comparison{
			LHS:      parser.NewIntegerValue(1),
			RHS:      parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			Operator: parser.Token{Token: '=', Literal: "="},
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
			Operator: parser.Token{Token: '=', Literal: "="},
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
									Operator: parser.Token{Token: '=', Literal: "="},
								},
							},
						},
					},
				},
			},
			Operator: parser.Token{Token: '=', Literal: "="},
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
			Operator: parser.Token{Token: '=', Literal: "="},
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
			Operator: parser.Token{Token: '=', Literal: "="},
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
			Operator: parser.Token{Token: '=', Literal: "="},
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
			Operator: parser.Token{Token: '=', Literal: "="},
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
			Operator: parser.Token{Token: '=', Literal: "="},
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
			Operator: parser.Token{Token: '=', Literal: "="},
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
			Operator: parser.Token{Token: '=', Literal: "="},
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
			Operator: parser.Token{Token: '=', Literal: "="},
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
			Operator: parser.Token{Token: '=', Literal: "="},
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
			Operator: parser.Token{Token: '=', Literal: "="},
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
			Operator: parser.Token{Token: '=', Literal: "="},
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
			Operator: parser.Token{Token: '=', Literal: "="},
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
		},
		Result: value.NewTernary(ternary.TRUE),
	},
	{
		Name: "Not In",
		Expr: parser.In{
			LHS: parser.NewIntegerValue(2),
			Values: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewIntegerValue(1),
						parser.NewIntegerValue(2),
						parser.NewNullValue(),
					},
				},
			},
			Negation: parser.Token{Token: parser.NOT, Literal: "not"},
		},
		Result: value.NewTernary(ternary.FALSE),
	},
	{
		Name: "Not In UNKNOWN",
		Expr: parser.In{
			LHS: parser.NewIntegerValue(2),
			Values: parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewIntegerValue(3),
						parser.NewIntegerValue(4),
						parser.NewNullValue(),
					},
				},
			},
			Negation: parser.Token{Token: parser.NOT, Literal: "not"},
		},
		Result: value.NewTernary(ternary.UNKNOWN),
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
		Scope: GenerateReferenceScope(nil, nil, time.Time{}, []ReferenceRecord{
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
				cache:       NewFieldIndexCache(10, LimitToUseFieldIndexSliceChache),
			},
		}),
		Expr: parser.In{
			LHS: parser.NewIntegerValue(2),
			Values: parser.RowValue{
				Value: parser.Subquery{
					Query: parser.SelectQuery{
						SelectEntity: parser.SelectEntity{
							SelectClause: parser.SelectClause{
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
									Operator: parser.Token{Token: '=', Literal: "="},
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
									Operator: parser.Token{Token: '=', Literal: "="},
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
		Scope: GenerateReferenceScope(nil, nil, time.Time{}, []ReferenceRecord{
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
				cache:       NewFieldIndexCache(10, LimitToUseFieldIndexSliceChache),
			},
		}),
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
		Scope: GenerateReferenceScope(nil, nil, time.Time{}, []ReferenceRecord{
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
				cache:       NewFieldIndexCache(10, LimitToUseFieldIndexSliceChache),
			},
		}),
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
		Scope: GenerateReferenceScope(nil, nil, time.Time{}, []ReferenceRecord{
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
				cache:       NewFieldIndexCache(10, LimitToUseFieldIndexSliceChache),
			},
		}),
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
		Scope: GenerateReferenceScope(nil, nil, time.Time{}, []ReferenceRecord{
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
				cache:       NewFieldIndexCache(10, LimitToUseFieldIndexSliceChache),
			},
		}),
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
		Scope: GenerateReferenceScope(nil, nil, time.Time{}, []ReferenceRecord{
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
				cache:       NewFieldIndexCache(10, LimitToUseFieldIndexSliceChache),
			},
		}),
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
		Scope: GenerateReferenceScope(nil, nil, time.Time{}, []ReferenceRecord{
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
				cache:       NewFieldIndexCache(10, LimitToUseFieldIndexSliceChache),
			},
		}),
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
			Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: "<>"},
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
			Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: "<>"},
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
									Operator: parser.Token{Token: '=', Literal: "="},
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
			Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: "<>"},
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
			Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: "<>"},
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
			Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: ">"},
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
			Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: ">"},
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
			Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: ">"},
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
									Operator: parser.Token{Token: '=', Literal: "="},
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
			Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: ">"},
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
			Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: "="},
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
		Scope: GenerateReferenceScope(nil, nil, time.Time{}, []ReferenceRecord{
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
				cache:       NewFieldIndexCache(10, LimitToUseFieldIndexSliceChache),
			},
		}),
		Expr: parser.Exists{
			Query: parser.Subquery{
				Query: parser.SelectQuery{
					SelectEntity: parser.SelectEntity{
						SelectClause: parser.SelectClause{
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
								Operator: parser.Token{Token: '=', Literal: "="},
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
								Operator: parser.Token{Token: '=', Literal: "="},
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
		Scope: GenerateReferenceScope(nil, nil, time.Time{}, []ReferenceRecord{
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
				cache:       NewFieldIndexCache(10, LimitToUseFieldIndexSliceChache),
			},
		}),
		Expr: parser.Subquery{
			Query: parser.SelectQuery{
				SelectEntity: parser.SelectEntity{
					SelectClause: parser.SelectClause{
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
							Operator: parser.Token{Token: '=', Literal: "="},
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
					},
				},
			},
		}, nil, time.Time{}, nil),
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
					},
				},
			},
		}, nil, time.Time{}, nil),
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
		Scope: GenerateReferenceScope(nil, nil, time.Time{}, []ReferenceRecord{
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
				},
				recordIndex: 0,
				cache:       NewFieldIndexCache(10, LimitToUseFieldIndexSliceChache),
			},
		}),
		Expr: parser.AggregateFunction{
			Name:     "avg",
			Distinct: parser.Token{},
			Args: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
			},
		},
		Result: value.NewFloat(2),
	},
	{
		Name: "Aggregate Function Argument Length Error",
		Scope: GenerateReferenceScope(nil, nil, time.Time{}, []ReferenceRecord{
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
				},
				recordIndex: 0,
				cache:       NewFieldIndexCache(10, LimitToUseFieldIndexSliceChache),
			},
		}),
		Expr: parser.AggregateFunction{
			Name:     "avg",
			Distinct: parser.Token{},
		},
		Error: "function avg takes exactly 1 argument",
	},
	{
		Name: "Aggregate Function Not Grouped Error",
		Scope: GenerateReferenceScope(nil, nil, time.Time{}, []ReferenceRecord{
			{
				view: &View{
					Header: NewHeader("table1", []string{"column1", "column2"}),
					RecordSet: []Record{
						NewRecordWithId(1, []value.Primary{
							value.NewInteger(1),
							value.NewString("str2"),
						}),
					},
				},
				recordIndex: 0,
				cache:       NewFieldIndexCache(10, LimitToUseFieldIndexSliceChache),
			},
		}),
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
		Scope: GenerateReferenceScope(nil, nil, time.Time{}, []ReferenceRecord{
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
				},
				recordIndex: 0,
				cache:       NewFieldIndexCache(10, LimitToUseFieldIndexSliceChache),
			},
		}),
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
		Error: "aggregate functions are nested at AVG(AVG(column1))",
	},
	{
		Name: "Aggregate Function Count With AllColumns",
		Scope: GenerateReferenceScope(nil, nil, time.Time{}, []ReferenceRecord{
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
				},
				recordIndex: 0,
				cache:       NewFieldIndexCache(10, LimitToUseFieldIndexSliceChache),
			},
		}),
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
		Name: "Aggregate Function Count With Null",
		Scope: GenerateReferenceScope(nil, nil, time.Time{}, []ReferenceRecord{
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
				},
				recordIndex: 0,
				cache:       NewFieldIndexCache(10, LimitToUseFieldIndexSliceChache),
			},
		}),
		Expr: parser.AggregateFunction{
			Name:     "count",
			Distinct: parser.Token{},
			Args: []parser.QueryExpression{
				parser.PrimitiveType{Value: value.NewNull()},
			},
		},
		Result: value.NewInteger(0),
	},
	{
		Name: "Aggregate Function As a Statement Error",
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
		Scope: GenerateReferenceScope([]map[string]map[string]interface{}{
			{
				scopeNameFunctions: {
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
											Operator: parser.Token{Token: '*', Literal: "*"},
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
		}, nil, time.Time{}, []ReferenceRecord{
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
				},
				recordIndex: 0,
				cache:       NewFieldIndexCache(10, LimitToUseFieldIndexSliceChache),
			},
		}),
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
		Scope: GenerateReferenceScope([]map[string]map[string]interface{}{
			{
				scopeNameFunctions: {
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
		}, nil, time.Time{}, []ReferenceRecord{
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
				},
				recordIndex: 0,
				cache:       NewFieldIndexCache(10, LimitToUseFieldIndexSliceChache),
			},
		}),
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
		Scope: GenerateReferenceScope([]map[string]map[string]interface{}{
			{
				scopeNameFunctions: {
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
		}, nil, time.Time{}, []ReferenceRecord{
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
				},
				recordIndex: 0,
				cache:       NewFieldIndexCache(10, LimitToUseFieldIndexSliceChache),
			},
		}),
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
		Name: "Aggregate Function Execute User Defined Aggregate Function Passed As Scalar Function",
		Scope: GenerateReferenceScope([]map[string]map[string]interface{}{
			{
				scopeNameFunctions: {
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
											Operator: parser.Token{Token: '*', Literal: "*"},
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
		}, nil, time.Time{}, []ReferenceRecord{
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
				},
				recordIndex: 0,
				cache:       NewFieldIndexCache(10, LimitToUseFieldIndexSliceChache),
			},
		}),
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
		Scope: GenerateReferenceScope([]map[string]map[string]interface{}{
			{
				scopeNameFunctions: {
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
											Operator: parser.Token{Token: '*', Literal: "*"},
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
		}, nil, time.Time{}, []ReferenceRecord{
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
				},
				recordIndex: 0,
				cache:       NewFieldIndexCache(10, LimitToUseFieldIndexSliceChache),
			},
		}),
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
		Scope: GenerateReferenceScope(nil, nil, time.Time{}, []ReferenceRecord{
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
					isGrouped: true,
				},
				recordIndex: 0,
				cache:       NewFieldIndexCache(10, LimitToUseFieldIndexSliceChache),
			},
		}),
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
		Scope: GenerateReferenceScope(nil, nil, time.Time{}, []ReferenceRecord{
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
					isGrouped: true,
				},
				recordIndex: 0,
				cache:       NewFieldIndexCache(10, LimitToUseFieldIndexSliceChache),
			},
		}),
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
		Scope: GenerateReferenceScope(nil, nil, time.Time{}, []ReferenceRecord{
			{
				view: &View{
					Header: NewHeaderWithId("table1", []string{"column1", "column2"}),
					RecordSet: []Record{
						NewRecordWithId(1, []value.Primary{
							value.NewInteger(1),
							value.NewString("str2"),
						}),
					},
				},
				recordIndex: 0,
				cache:       NewFieldIndexCache(10, LimitToUseFieldIndexSliceChache),
			},
		}),
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
		Scope: GenerateReferenceScope(nil, nil, time.Time{}, []ReferenceRecord{
			{
				view: &View{
					Header: NewHeaderWithId("table1", []string{"column1", "column2"}),
					RecordSet: []Record{
						NewRecordWithId(1, []value.Primary{
							value.NewInteger(1),
							value.NewString("str2"),
						}),
					},
				},
				recordIndex: 0,
				cache:       NewFieldIndexCache(10, LimitToUseFieldIndexSliceChache),
			},
		}),
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
		Scope: GenerateReferenceScope(nil, nil, time.Time{}, []ReferenceRecord{
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
					isGrouped: true,
				},
				recordIndex: 0,
				cache:       NewFieldIndexCache(10, LimitToUseFieldIndexSliceChache),
			},
		}),
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
		Scope: GenerateReferenceScope(nil, nil, time.Time{}, []ReferenceRecord{
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
					isGrouped: true,
				},
				recordIndex: 0,
				cache:       NewFieldIndexCache(10, LimitToUseFieldIndexSliceChache),
			},
		}),
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
		Error: "aggregate functions are nested at LISTAGG(AVG(column1))",
	},
	{
		Name: "ListAgg Function Second Argument Evaluation Error",
		Scope: GenerateReferenceScope(nil, nil, time.Time{}, []ReferenceRecord{
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
					isGrouped: true,
				},
				recordIndex: 0,
				cache:       NewFieldIndexCache(10, LimitToUseFieldIndexSliceChache),
			},
		}),
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
		Scope: GenerateReferenceScope(nil, nil, time.Time{}, []ReferenceRecord{
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
					isGrouped: true,
				},
				recordIndex: 0,
				cache:       NewFieldIndexCache(10, LimitToUseFieldIndexSliceChache),
			},
		}),
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
		Name: "ListAgg Function As a Statement Error",
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
		Scope: GenerateReferenceScope(nil, nil, time.Time{}, []ReferenceRecord{
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
					isGrouped: true,
				},
				recordIndex: 0,
				cache:       NewFieldIndexCache(10, LimitToUseFieldIndexSliceChache),
			},
		}),
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
		Scope: GenerateReferenceScope(nil, nil, time.Time{}, []ReferenceRecord{
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
					isGrouped: true,
				},
				recordIndex: 0,
				cache:       NewFieldIndexCache(10, LimitToUseFieldIndexSliceChache),
			},
		}),
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
		Name: "Analytic Function",
		Scope: GenerateReferenceScope(nil, nil, time.Time{}, []ReferenceRecord{
			{
				view: &View{
					Header: append(
						NewHeader("table1", []string{"column1", "column2"}),
						HeaderField{Identifier: "RANK() OVER (PARTITION BY column1 ORDER BY column2)", Column: "RANK() OVER (PARTITION BY column1 ORDER BY column2)"},
					),
					RecordSet: []Record{
						NewRecord([]value.Primary{
							value.NewString("a"),
							value.NewInteger(11),
							value.NewInteger(1),
						}),
						NewRecord([]value.Primary{
							value.NewString("b"),
							value.NewInteger(22),
							value.NewInteger(2),
						}),
						NewRecord([]value.Primary{
							value.NewString("c"),
							value.NewInteger(33),
							value.NewInteger(3),
						}),
					},
				},
				recordIndex: 1,
				cache:       NewFieldIndexCache(10, LimitToUseFieldIndexSliceChache),
			},
		}),
		Expr: parser.AnalyticFunction{
			Name: "rank",
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
		Result: value.NewInteger(2),
	},
	{
		Name: "Analytic Function Not Allowed Error",
		Scope: GenerateReferenceScope(nil, nil, time.Time{}, []ReferenceRecord{
			{
				view: &View{
					Header: NewHeader("table1", []string{"column1", "column2"}),
					RecordSet: []Record{
						NewRecord([]value.Primary{
							value.NewString("a"),
							value.NewInteger(11),
						}),
						NewRecord([]value.Primary{
							value.NewString("b"),
							value.NewInteger(22),
						}),
						NewRecord([]value.Primary{
							value.NewString("c"),
							value.NewInteger(33),
						}),
					},
				},
				recordIndex: 1,
				cache:       NewFieldIndexCache(10, LimitToUseFieldIndexSliceChache),
			},
		}),
		Expr: parser.AnalyticFunction{
			Name: "rank",
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
		Error: "analytic function rank is only available in select clause or order by clause",
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
						Operator: parser.Token{Token: '=', Literal: "="},
					},
					Result: parser.NewStringValue("A"),
				},
				parser.CaseExprWhen{
					Condition: parser.Comparison{
						LHS:      parser.NewIntegerValue(2),
						RHS:      parser.NewIntegerValue(2),
						Operator: parser.Token{Token: '=', Literal: "="},
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
		Scope: GenerateReferenceScope([]map[string]map[string]interface{}{
			{
				scopeNameVariables: {
					"var1": value.NewInteger(1),
				},
			},
		}, nil, time.Time{}, nil),
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
		Error: "@@INVALID is an unknown flag",
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
		Scope: GenerateReferenceScope([]map[string]map[string]interface{}{
			{
				scopeNameVariables: {
					"var1": value.NewInteger(1),
				},
			},
		}, nil, time.Time{}, nil),
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
			Cursor:   parser.Identifier{Literal: "cur"},
			Negation: parser.Token{Token: parser.NOT, Literal: "not"},
			Type:     parser.Token{Token: parser.OPEN, Literal: "open"},
		},
		Result: value.NewTernary(ternary.FALSE),
	},
	{
		Name: "Cursor Status Is In Range",
		Expr: parser.CursorStatus{
			Cursor: parser.Identifier{Literal: "cur"},
			Type:   parser.Token{Token: parser.RANGE, Literal: "range"},
		},
		Result: value.NewTernary(ternary.TRUE),
	},
	{
		Name: "Cursor Status Open Error",
		Expr: parser.CursorStatus{
			Cursor: parser.Identifier{Literal: "notexist"},
			Type:   parser.Token{Token: parser.OPEN, Literal: "open"},
		},
		Error: "cursor notexist is undeclared",
	},
	{
		Name: "Cursor Status In Range Error",
		Expr: parser.CursorStatus{
			Cursor: parser.Identifier{Literal: "notexist"},
			Type:   parser.Token{Token: parser.RANGE, Literal: "range"},
		},
		Error: "cursor notexist is undeclared",
	},
	{
		Name: "Cursor Attribute Count",
		Expr: parser.CursorAttrebute{
			Cursor:    parser.Identifier{Literal: "cur"},
			Attrebute: parser.Token{Token: parser.COUNT, Literal: "count"},
		},
		Result: value.NewInteger(3),
	},
	{
		Name: "Cursor Attribute Count Error",
		Expr: parser.CursorAttrebute{
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

func TestEvaluate(t *testing.T) {
	defer func() {
		_ = TestTx.cachedViews.Clean(TestTx.FileContainer)
		initFlag(TestTx.Flags)
	}()

	TestTx.Flags.Repository = TestDataDir
	scope := GenerateReferenceScope([]map[string]map[string]interface{}{
		{
			scopeNameCursors: {
				"CUR": &Cursor{
					query: selectQueryForCursorTest,
					mtx:   &sync.Mutex{},
				},
			},
		},
	}, nil, time.Time{}, nil)

	ctx := context.Background()
	_ = scope.OpenCursor(ctx, parser.Identifier{Literal: "cur"}, nil)
	_, _ = scope.FetchCursor(parser.Identifier{Literal: "cur"}, parser.NEXT, 0)

	for _, v := range evaluateTests {
		_ = TestTx.cachedViews.Clean(TestTx.FileContainer)

		if v.Scope == nil {
			v.Scope = scope
		}

		evalCtx := ctx
		if v.ReplaceValues != nil {
			evalCtx = ContextForPreparedStatement(ctx, v.ReplaceValues)
		}
		result, err := Evaluate(evalCtx, v.Scope, v.Expr)
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

var evaluateEmbeddedStringTests = []struct {
	Input  string
	Expect string
	Error  string
}{
	{
		Input:  "str",
		Expect: "str",
	},
	{
		Input:  "str\\tstr",
		Expect: "str\tstr",
	},
	{
		Input:  "str\\\\tstr",
		Expect: "str\\tstr",
	},
	{
		Input:  "str''str",
		Expect: "str''str",
	},
	{
		Input:  "'str''str'",
		Expect: "str'str",
	},
	{
		Input:  "\"str\"\"str\"",
		Expect: "str\"str",
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
		Input:  "abc${'a\\tb'}def",
		Expect: "abca\tbdef",
	},
	{
		Input:  "abc${'a\\\\tb'}def",
		Expect: "abca\\tbdef",
	},
	{
		Input: "@notexist",
		Error: "variable @notexist is undeclared",
	},
	{
		Input: "@#notexist",
		Error: "@#NOTEXIST is an unknown runtime information",
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
	{
		Input: "abc${@var1; @var2;}def",
		Error: "'@var1; @var2;': cannot evaluate as a value",
	},
}

func TestEvaluateEmbeddedString(t *testing.T) {
	scope := NewReferenceScope(TestTx)
	_ = scope.DeclareVariableDirectly(parser.Variable{Name: "var"}, value.NewInteger(1))
	_ = os.Setenv("CSVQ_TEST_FILTER", "FILTER_TEST")

	for _, v := range evaluateEmbeddedStringTests {
		result, err := EvaluateEmbeddedString(context.Background(), scope, v.Input)

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

func generateBenchGroupedViewScope() *ReferenceScope {
	primaries := make([]value.Primary, 10000)
	for i := 0; i < 10000; i++ {
		primaries[i] = value.NewInteger(int64(i))
	}

	view := &View{
		Header: NewHeader("table1", []string{"c1"}),
		RecordSet: []Record{
			{
				NewGroupCell(primaries),
			},
		},
		isGrouped: true,
	}

	return NewReferenceScope(TestTx).CreateScopeForRecordEvaluation(
		view,
		0,
	)
}

func BenchmarkEvaluateCountAllColumns(b *testing.B) {
	ctx := context.Background()
	scope := generateBenchGroupedViewScope()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = Evaluate(ctx, scope, parser.AggregateFunction{
			Name:     "count",
			Distinct: parser.Token{},
			Args: []parser.QueryExpression{
				parser.AllColumns{},
			},
		})
	}
}

func BenchmarkEvaluateCount(b *testing.B) {
	ctx := context.Background()
	scope := generateBenchGroupedViewScope()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = Evaluate(ctx, scope, parser.AggregateFunction{
			Name:     "count",
			Distinct: parser.Token{},
			Args: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "c1"}},
			},
		})
	}
}

func BenchmarkEvaluateSingleThread(b *testing.B) {
	ctx := context.Background()

	for i := 0; i < b.N; i++ {
		scope := NewReferenceScope(TestTx)

		for j := 0; j < 150; j++ {
			_, _ = Evaluate(ctx, scope, parser.Comparison{
				LHS:      parser.NewIntegerValue(1),
				RHS:      parser.NewStringValue("1"),
				Operator: parser.Token{Token: '=', Literal: "="},
			})
		}
	}
}

func BenchmarkEvaluateMultiThread(b *testing.B) {
	ctx := context.Background()

	for i := 0; i < b.N; i++ {
		wg := sync.WaitGroup{}
		for i := 0; i < 3; i++ {
			wg.Add(1)
			go func(thIdx int) {
				scope := NewReferenceScope(TestTx)

				for j := 0; j < 50; j++ {
					_, _ = Evaluate(ctx, scope, parser.Comparison{
						LHS:      parser.NewIntegerValue(1),
						RHS:      parser.NewStringValue("1"),
						Operator: parser.Token{Token: '=', Literal: "="},
					})
				}
				wg.Done()
			}(i)
		}
		wg.Wait()
	}
}

var evaluateFieldReferenceBenchScope = GenerateReferenceScope(nil, nil, time.Time{}, []ReferenceRecord{
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
		cache:       NewFieldIndexCache(10, LimitToUseFieldIndexSliceChache),
	},
})

var evaluateFieldReferenceBenchExpr = parser.FieldReference{
	Column: parser.Identifier{Literal: "column3"},
}

func BenchmarkEvaluateFieldReference(b *testing.B) {
	ctx := context.Background()
	scope := evaluateFieldReferenceBenchScope
	expr := evaluateFieldReferenceBenchExpr

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Evaluate(ctx, scope, expr)
	}
}
