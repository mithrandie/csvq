package query

import (
	"github.com/mithrandie/csvq/lib/parser"
	"reflect"
	"testing"
	"time"
)

var hasAggregateFunctionTests = []struct {
	Name   string
	Expr   parser.QueryExpression
	Result bool
	Error  string
}{
	{
		Name: "Aggregate Function",
		Expr: parser.AggregateFunction{
			Name:     "avg",
			Distinct: parser.Token{},
			Args: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
			},
		},
		Result: true,
	},
	{
		Name: "List Function",
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
		Result: true,
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
		Result: false,
	},
	{
		Name: "User Defined Aggregate Function",
		Expr: parser.Function{
			Name: "useraggfunc",
			Args: []parser.QueryExpression{
				parser.NewIntegerValue(1),
			},
		},
		Result: true,
	},
	{
		Name: "Aggregate Function in Function Arguments",
		Expr: parser.Function{
			Name: "coalesce",
			Args: []parser.QueryExpression{
				parser.AggregateFunction{
					Name:     "avg",
					Distinct: parser.Token{},
					Args: []parser.QueryExpression{
						parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
					},
				},
				parser.NewStringValue("str"),
			},
		},
		Result: true,
	},
	{
		Name: "Analytic Function",
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
							Value: parser.AggregateFunction{
								Name:     "avg",
								Distinct: parser.Token{},
								Args: []parser.QueryExpression{
									parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
								},
							},
						},
					},
				},
			},
		},
		Result: true,
	},
	{
		Name: "Case Expression",
		Expr: parser.CaseExpr{
			Value: parser.NewIntegerValue(2),
			When: []parser.QueryExpression{
				parser.CaseExprWhen{
					Condition: parser.NewIntegerValue(1),
					Result: parser.AggregateFunction{
						Name:     "avg",
						Distinct: parser.Token{},
						Args: []parser.QueryExpression{
							parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
						},
					},
				},
				parser.CaseExprWhen{
					Condition: parser.NewIntegerValue(2),
					Result:    parser.NewStringValue("B"),
				},
			},
		},
		Result: true,
	},
	{
		Name: "Case Expression without Comparison",
		Expr: parser.CaseExpr{
			When: []parser.QueryExpression{
				parser.CaseExprWhen{
					Condition: parser.NewIntegerValue(1),
					Result: parser.AggregateFunction{
						Name:     "avg",
						Distinct: parser.Token{},
						Args: []parser.QueryExpression{
							parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
						},
					},
				},
				parser.CaseExprWhen{
					Condition: parser.NewIntegerValue(2),
					Result:    parser.NewStringValue("B"),
				},
			},
		},
		Result: true,
	},
	{
		Name: "Case Expression with Else",
		Expr: parser.CaseExpr{
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
				Result: parser.AggregateFunction{
					Name:     "avg",
					Distinct: parser.Token{},
					Args: []parser.QueryExpression{
						parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
					},
				},
			},
		},
		Result: true,
	},
}

func TestHasAggregateFunction(t *testing.T) {
	scope := GenerateReferenceScope([]map[string]map[string]interface{}{
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
				"USERAGGFUNC": &UserDefinedFunction{
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
	}, nil, time.Time{}, nil)

	for _, v := range hasAggregateFunctionTests {
		result, err := HasAggregateFunction(v.Expr, scope)
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
		if result != v.Result {
			t.Errorf("%s: result = %t, want %t", v.Name, result, v.Result)
		}
	}
}

var searchAnalyticFunctionsTests = []struct {
	Name   string
	Expr   parser.QueryExpression
	Result []parser.AnalyticFunction
	Error  string
}{
	{
		Name: "Analytic Function",
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
		Result: []parser.AnalyticFunction{
			{
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
		},
	},
	{
		Name: "Nested Analytic Function",
		Expr: parser.AnalyticFunction{
			Name: "sum",
			Args: []parser.QueryExpression{
				parser.AnalyticFunction{
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
			},
			AnalyticClause: parser.AnalyticClause{},
		},
		Result: []parser.AnalyticFunction{
			{
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
			{
				Name: "sum",
				Args: []parser.QueryExpression{
					parser.AnalyticFunction{
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
				},
				AnalyticClause: parser.AnalyticClause{},
			},
		},
	},
	{
		Name: "Arithmetic",
		Expr: parser.Arithmetic{
			LHS: parser.AnalyticFunction{
				Name:           "rank",
				AnalyticClause: parser.AnalyticClause{},
			},
			RHS: parser.AnalyticFunction{
				Name:           "row_number",
				AnalyticClause: parser.AnalyticClause{},
			},
			Operator: parser.Token{Token: '+', Literal: "+"},
		},
		Result: []parser.AnalyticFunction{
			{
				Name:           "row_number",
				AnalyticClause: parser.AnalyticClause{},
			},
			{
				Name:           "rank",
				AnalyticClause: parser.AnalyticClause{},
			},
		},
	},
	{
		Name: "Case Expression",
		Expr: parser.CaseExpr{
			Value: parser.NewIntegerValue(2),
			When: []parser.QueryExpression{
				parser.CaseExprWhen{
					Condition: parser.NewIntegerValue(1),
					Result: parser.AnalyticFunction{
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
				},
				parser.CaseExprWhen{
					Condition: parser.NewIntegerValue(2),
					Result:    parser.NewStringValue("B"),
				},
			},
		},
		Result: []parser.AnalyticFunction{
			{
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
		},
	},
	{
		Name: "Case Expression without Comparison",
		Expr: parser.CaseExpr{
			When: []parser.QueryExpression{
				parser.CaseExprWhen{
					Condition: parser.NewIntegerValue(1),
					Result: parser.AnalyticFunction{
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
				},
				parser.CaseExprWhen{
					Condition: parser.NewIntegerValue(2),
					Result:    parser.NewStringValue("B"),
				},
			},
		},
		Result: []parser.AnalyticFunction{
			{
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
		},
	},
	{
		Name: "Case Expression with Else",
		Expr: parser.CaseExpr{
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
				Result: parser.AnalyticFunction{
					Name:           "rank",
					AnalyticClause: parser.AnalyticClause{},
				},
			},
		},
		Result: []parser.AnalyticFunction{
			{
				Name:           "rank",
				AnalyticClause: parser.AnalyticClause{},
			},
		},
	},
}

func TestSearchAnalyticFunctions(t *testing.T) {
	for _, v := range searchAnalyticFunctionsTests {
		result, err := SearchAnalyticFunctions(v.Expr)
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
			t.Errorf("%s: result = %v, want %v", v.Name, result, v.Result)
		}
	}
}
