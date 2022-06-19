package query

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/value"
)

var analyzeTests = []struct {
	Name             string
	CPU              int
	View             *View
	Scope            *ReferenceScope
	Function         parser.AnalyticFunction
	PartitionIndices []int
	Result           *View
	Error            string
}{
	{
		Name: "Analyze AnalyticFunction",
		CPU:  3,
		View: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(1),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(1),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(1),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(2),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(2),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(3),
				}),
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(2),
				}),
			},
			sortValuesInEachRecord: []SortValues{
				{NewSortValue(value.NewInteger(1), TestTx.Flags)},
				{NewSortValue(value.NewInteger(1), TestTx.Flags)},
				{NewSortValue(value.NewInteger(1), TestTx.Flags)},
				{NewSortValue(value.NewInteger(2), TestTx.Flags)},
				{NewSortValue(value.NewInteger(2), TestTx.Flags)},
				{NewSortValue(value.NewInteger(3), TestTx.Flags)},
				{NewSortValue(value.NewInteger(2), TestTx.Flags)},
			},
		},
		Function: parser.AnalyticFunction{
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
		PartitionIndices: []int{0},
		Result: &View{
			Header: append(
				NewHeader("table1", []string{"column1", "column2"}),
				HeaderField{Identifier: "RANK() OVER (PARTITION BY column1 ORDER BY column2)", Column: "RANK() OVER (PARTITION BY column1 ORDER BY column2)"},
			),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(1),
					value.NewInteger(1),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(1),
					value.NewInteger(1),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(1),
					value.NewInteger(1),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(2),
					value.NewInteger(3),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(2),
					value.NewInteger(3),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(3),
					value.NewInteger(5),
				}),
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(2),
					value.NewInteger(2),
				}),
			},
			sortValuesInEachCell: [][]*SortValue{
				{NewSortValue(value.NewString("a"), TestTx.Flags), nil},
				{NewSortValue(value.NewString("b"), TestTx.Flags), nil},
				{NewSortValue(value.NewString("b"), TestTx.Flags), nil},
				{NewSortValue(value.NewString("b"), TestTx.Flags), nil},
				{NewSortValue(value.NewString("b"), TestTx.Flags), nil},
				{NewSortValue(value.NewString("b"), TestTx.Flags), nil},
				{NewSortValue(value.NewString("a"), TestTx.Flags), nil},
			},
			sortValuesInEachRecord: []SortValues{
				{NewSortValue(value.NewInteger(1), TestTx.Flags)},
				{NewSortValue(value.NewInteger(1), TestTx.Flags)},
				{NewSortValue(value.NewInteger(1), TestTx.Flags)},
				{NewSortValue(value.NewInteger(2), TestTx.Flags)},
				{NewSortValue(value.NewInteger(2), TestTx.Flags)},
				{NewSortValue(value.NewInteger(3), TestTx.Flags)},
				{NewSortValue(value.NewInteger(2), TestTx.Flags)},
			},
		},
	},
	{
		Name: "Analyze AnalyticFunction Empty Record",
		CPU:  3,
		View: &View{
			Header:                 NewHeader("table1", []string{"column1", "column2"}),
			RecordSet:              []Record{},
			sortValuesInEachRecord: []SortValues{},
		},
		Function: parser.AnalyticFunction{
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
		PartitionIndices: []int{0},
		Result: &View{
			Header: append(
				NewHeader("table1", []string{"column1", "column2"}),
				HeaderField{Identifier: "RANK() OVER (PARTITION BY column1 ORDER BY column2)", Column: "RANK() OVER (PARTITION BY column1 ORDER BY column2)"},
			),
			RecordSet:              []Record{},
			sortValuesInEachCell:   [][]*SortValue{},
			sortValuesInEachRecord: []SortValues{},
		},
	},
	{
		Name: "Analyze AnalyticFunction Argument Length Error",
		View: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(1),
				}),
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(2),
				}),
			},
		},
		Function: parser.AnalyticFunction{
			Name: "rank",
			Args: []parser.QueryExpression{
				parser.NewIntegerValue(1),
			},
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
		Error: "function rank takes no argument",
	},
	{
		Name: "Analyze AnalyticFunction Execution Error",
		View: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(1),
				}),
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(2),
				}),
			},
		},
		Function: parser.AnalyticFunction{
			Name: "first_value",
			Args: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			},
		},
		Error: "field notexist does not exist",
	},
	{
		Name: "Analyze AggregateFunction",
		View: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(1),
				}),
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(2),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(1),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewNull(),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(1),
				}),
			},
		},
		Function: parser.AnalyticFunction{
			Name: "count",
			Args: []parser.QueryExpression{
				parser.AllColumns{},
			},
			AnalyticClause: parser.AnalyticClause{
				PartitionClause: parser.PartitionClause{
					Values: []parser.QueryExpression{
						parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
					},
				},
			},
		},
		PartitionIndices: []int{0},
		Result: &View{
			Header: append(
				NewHeader("table1", []string{"column1", "column2"}),
				HeaderField{Identifier: "COUNT(*) OVER (PARTITION BY column1)", Column: "COUNT(*) OVER (PARTITION BY column1)"},
			),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(1),
					value.NewInteger(2),
				}),
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(2),
					value.NewInteger(2),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(1),
					value.NewInteger(3),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewNull(),
					value.NewInteger(3),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(1),
					value.NewInteger(3),
				}),
			},
			sortValuesInEachCell: [][]*SortValue{
				{NewSortValue(value.NewString("a"), TestTx.Flags), nil},
				{NewSortValue(value.NewString("a"), TestTx.Flags), nil},
				{NewSortValue(value.NewString("b"), TestTx.Flags), nil},
				{NewSortValue(value.NewString("b"), TestTx.Flags), nil},
				{NewSortValue(value.NewString("b"), TestTx.Flags), nil},
			},
		},
	},
	{
		Name: "Analyze AggregateFunction with Windowing Clause",
		View: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(1),
				}),
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(2),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(1),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewNull(),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(1),
				}),
			},
		},
		Function: parser.AnalyticFunction{
			Name: "sum",
			Args: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
			},
			AnalyticClause: parser.AnalyticClause{
				PartitionClause: parser.PartitionClause{
					Values: []parser.QueryExpression{
						parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
					},
				},
				OrderByClause: parser.OrderByClause{
					Items: []parser.QueryExpression{
						parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
					},
				},
				WindowingClause: parser.WindowingClause{
					FrameLow: parser.WindowFramePosition{
						Direction: parser.Token{Token: parser.PRECEDING},
						Unbounded: parser.Token{Token: parser.UNBOUNDED},
					},
				},
			},
		},
		PartitionIndices: []int{0},
		Result: &View{
			Header: append(
				NewHeader("table1", []string{"column1", "column2"}),
				HeaderField{Identifier: "SUM(column2) OVER (PARTITION BY column1 ORDER BY column1 ROWS UNBOUNDED PRECEDING)", Column: "SUM(column2) OVER (PARTITION BY column1 ORDER BY column1 ROWS UNBOUNDED PRECEDING)"},
			),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(1),
					value.NewFloat(1),
				}),
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(2),
					value.NewFloat(3),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(1),
					value.NewFloat(1),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewNull(),
					value.NewFloat(1),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(1),
					value.NewFloat(2),
				}),
			},
			sortValuesInEachCell: [][]*SortValue{
				{NewSortValue(value.NewString("a"), TestTx.Flags), nil},
				{NewSortValue(value.NewString("a"), TestTx.Flags), nil},
				{NewSortValue(value.NewString("b"), TestTx.Flags), nil},
				{NewSortValue(value.NewString("b"), TestTx.Flags), nil},
				{NewSortValue(value.NewString("b"), TestTx.Flags), nil},
			},
		},
	},
	{
		Name: "Analyze AggregateFunction With Distinct",
		View: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(1),
				}),
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(2),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(1),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewNull(),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(1),
				}),
			},
			sortValuesInEachCell: [][]*SortValue{
				{NewSortValue(value.NewString("a"), TestTx.Flags), nil},
				{NewSortValue(value.NewString("a"), TestTx.Flags), nil},
				{NewSortValue(value.NewString("b"), TestTx.Flags), nil},
				{NewSortValue(value.NewString("b"), TestTx.Flags), nil},
				{NewSortValue(value.NewString("b"), TestTx.Flags), nil},
			},
		},
		Function: parser.AnalyticFunction{
			Name:     "count",
			Distinct: parser.Token{Token: parser.DISTINCT, Literal: "distinct"},
			Args: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
			},
			AnalyticClause: parser.AnalyticClause{
				PartitionClause: parser.PartitionClause{
					Values: []parser.QueryExpression{
						parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
					},
				},
			},
		},
		PartitionIndices: []int{0},
		Result: &View{
			Header: append(
				NewHeader("table1", []string{"column1", "column2"}),
				HeaderField{Identifier: "COUNT(DISTINCT column2) OVER (PARTITION BY column1)", Column: "COUNT(DISTINCT column2) OVER (PARTITION BY column1)"},
			),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(1),
					value.NewInteger(2),
				}),
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(2),
					value.NewInteger(2),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(1),
					value.NewInteger(1),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewNull(),
					value.NewInteger(1),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(1),
					value.NewInteger(1),
				}),
			},
			sortValuesInEachCell: [][]*SortValue{
				{NewSortValue(value.NewString("a"), TestTx.Flags), nil},
				{NewSortValue(value.NewString("a"), TestTx.Flags), nil},
				{NewSortValue(value.NewString("b"), TestTx.Flags), nil},
				{NewSortValue(value.NewString("b"), TestTx.Flags), nil},
				{NewSortValue(value.NewString("b"), TestTx.Flags), nil},
			},
		},
	},
	{
		Name: "Analyze AggregateFunction Argument Length Error",
		View: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(1),
				}),
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(2),
				}),
			},
		},
		Function: parser.AnalyticFunction{
			Name: "count",
			AnalyticClause: parser.AnalyticClause{
				PartitionClause: parser.PartitionClause{
					Values: []parser.QueryExpression{
						parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
					},
				},
			},
		},
		Error: "function count takes exactly 1 argument",
	},
	{
		Name: "Analyze AggregateFunction Aggregate Value Error",
		View: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(1),
				}),
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(2),
				}),
			},
		},
		Function: parser.AnalyticFunction{
			Name: "count",
			Args: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			},
			AnalyticClause: parser.AnalyticClause{
				PartitionClause: parser.PartitionClause{
					Values: []parser.QueryExpression{
						parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
					},
				},
			},
		},
		Error: "field notexist does not exist",
	},
	{
		Name: "Analyze UserDefinedFunction",
		View: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(1),
				}),
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(2),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(1),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewNull(),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(1),
				}),
			},
		},
		Scope: GenerateReferenceScope([]map[string]map[string]interface{}{
			{
				scopeNameFunctions: {
					"USERAGGFUNC": &UserDefinedFunction{
						Name:        parser.Identifier{Literal: "useraggfunc"},
						IsAggregate: true,
						Cursor:      parser.Identifier{Literal: "list"},
						Parameters: []parser.Variable{
							{Name: "default"},
						},
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
								Cursor: parser.Identifier{Literal: "list"},
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
											Operator: parser.Token{Token: '*'},
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
		}, nil, time.Time{}, nil),
		Function: parser.AnalyticFunction{
			Name: "useraggfunc",
			Args: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
				parser.NewIntegerValue(0),
			},
			AnalyticClause: parser.AnalyticClause{
				PartitionClause: parser.PartitionClause{
					Values: []parser.QueryExpression{
						parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
					},
				},
			},
		},
		PartitionIndices: []int{0},
		Result: &View{
			Header: append(
				NewHeader("table1", []string{"column1", "column2"}),
				HeaderField{Identifier: "USERAGGFUNC(column2, 0) OVER (PARTITION BY column1)", Column: "USERAGGFUNC(column2, 0) OVER (PARTITION BY column1)"},
			),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(1),
					value.NewInteger(2),
				}),
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(2),
					value.NewInteger(2),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(1),
					value.NewInteger(1),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewNull(),
					value.NewInteger(1),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(1),
					value.NewInteger(1),
				}),
			},
			sortValuesInEachCell: [][]*SortValue{
				{NewSortValue(value.NewString("a"), TestTx.Flags), nil},
				{NewSortValue(value.NewString("a"), TestTx.Flags), nil},
				{NewSortValue(value.NewString("b"), TestTx.Flags), nil},
				{NewSortValue(value.NewString("b"), TestTx.Flags), nil},
				{NewSortValue(value.NewString("b"), TestTx.Flags), nil},
			},
		},
	},
	{
		Name: "Analyze UserDefinedFunction with Windowing Clause",
		View: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(1),
				}),
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(2),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(1),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewNull(),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(1),
				}),
			},
		},
		Scope: GenerateReferenceScope([]map[string]map[string]interface{}{
			{
				scopeNameFunctions: {
					"USERAGGFUNC": &UserDefinedFunction{
						Name:        parser.Identifier{Literal: "useraggfunc"},
						IsAggregate: true,
						Cursor:      parser.Identifier{Literal: "list"},
						Parameters: []parser.Variable{
							{Name: "default"},
						},
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
								Cursor: parser.Identifier{Literal: "list"},
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
		}, nil, time.Time{}, nil),
		Function: parser.AnalyticFunction{
			Name: "useraggfunc",
			Args: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
				parser.NewIntegerValue(0),
			},
			AnalyticClause: parser.AnalyticClause{
				PartitionClause: parser.PartitionClause{
					Values: []parser.QueryExpression{
						parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
					},
				},
				OrderByClause: parser.OrderByClause{
					Items: []parser.QueryExpression{
						parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
					},
				},
				WindowingClause: parser.WindowingClause{
					FrameLow: parser.WindowFramePosition{
						Direction: parser.Token{Token: parser.PRECEDING},
						Unbounded: parser.Token{Token: parser.UNBOUNDED},
					},
				},
			},
		},
		PartitionIndices: []int{0},
		Result: &View{
			Header: append(
				NewHeader("table1", []string{"column1", "column2"}),
				HeaderField{Identifier: "USERAGGFUNC(column2, 0) OVER (PARTITION BY column1 ORDER BY column1 ROWS UNBOUNDED PRECEDING)", Column: "USERAGGFUNC(column2, 0) OVER (PARTITION BY column1 ORDER BY column1 ROWS UNBOUNDED PRECEDING)"},
			),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(1),
					value.NewInteger(1),
				}),
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(2),
					value.NewInteger(2),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(1),
					value.NewInteger(1),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewNull(),
					value.NewInteger(1),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(1),
					value.NewInteger(1),
				}),
			},
			sortValuesInEachCell: [][]*SortValue{
				{NewSortValue(value.NewString("a"), TestTx.Flags), nil},
				{NewSortValue(value.NewString("a"), TestTx.Flags), nil},
				{NewSortValue(value.NewString("b"), TestTx.Flags), nil},
				{NewSortValue(value.NewString("b"), TestTx.Flags), nil},
				{NewSortValue(value.NewString("b"), TestTx.Flags), nil},
			},
		},
	},
	{
		Name: "Analyze UserDefinedFunction Argument Length Error",
		View: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(1),
				}),
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(2),
				}),
			},
		},
		Scope: GenerateReferenceScope([]map[string]map[string]interface{}{
			{
				scopeNameFunctions: {
					"USERAGGFUNC": &UserDefinedFunction{
						Name:        parser.Identifier{Literal: "useraggfunc"},
						IsAggregate: true,
						Cursor:      parser.Identifier{Literal: "list"},
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
		}, nil, time.Time{}, nil),
		Function: parser.AnalyticFunction{
			Name: "useraggfunc",
			AnalyticClause: parser.AnalyticClause{
				PartitionClause: parser.PartitionClause{
					Values: []parser.QueryExpression{
						parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
					},
				},
			},
		},
		Error: "function useraggfunc takes exactly 2 arguments",
	},
	{
		Name: "Analyze UserDefinedFunction Cursor Value Error",
		View: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(1),
				}),
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(2),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(1),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewNull(),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(1),
				}),
			},
		},
		Scope: GenerateReferenceScope([]map[string]map[string]interface{}{
			{
				scopeNameFunctions: {
					"USERAGGFUNC": &UserDefinedFunction{
						Name:        parser.Identifier{Literal: "useraggfunc"},
						IsAggregate: true,
						Cursor:      parser.Identifier{Literal: "list"},
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
		}, nil, time.Time{}, nil),
		Function: parser.AnalyticFunction{
			Name: "useraggfunc",
			Args: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
				parser.NewIntegerValue(0),
			},
			AnalyticClause: parser.AnalyticClause{
				PartitionClause: parser.PartitionClause{
					Values: []parser.QueryExpression{
						parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
					},
				},
			},
		},
		PartitionIndices: []int{0},
		Error:            "field notexist does not exist",
	},
	{
		Name: "Analyze UserDefinedFunction Argument Value Error",
		View: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(1),
				}),
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(2),
				}),
			},
		},
		Scope: GenerateReferenceScope([]map[string]map[string]interface{}{
			{
				scopeNameFunctions: {
					"USERAGGFUNC": &UserDefinedFunction{
						Name:        parser.Identifier{Literal: "useraggfunc"},
						IsAggregate: true,
						Cursor:      parser.Identifier{Literal: "list"},
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
		}, nil, time.Time{}, nil),
		Function: parser.AnalyticFunction{
			Name: "useraggfunc",
			Args: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
				parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			},
			AnalyticClause: parser.AnalyticClause{
				PartitionClause: parser.PartitionClause{
					Values: []parser.QueryExpression{
						parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
					},
				},
			},
		},
		Error: "field notexist does not exist",
	},
	{
		Name: "Analyze UserDefinedFunction Execution Error",
		View: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(1),
				}),
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(2),
				}),
			},
		},
		Scope: GenerateReferenceScope([]map[string]map[string]interface{}{
			{
				scopeNameFunctions: {
					"USERAGGFUNC": &UserDefinedFunction{
						Name:        parser.Identifier{Literal: "useraggfunc"},
						IsAggregate: true,
						Cursor:      parser.Identifier{Literal: "list"},
						Parameters: []parser.Variable{
							{Name: "default"},
						},
						Statements: []parser.Statement{
							parser.Return{
								Value: parser.Variable{Name: "undefined"},
							},
						},
					},
				},
			},
		}, nil, time.Time{}, nil),
		Function: parser.AnalyticFunction{
			Name: "useraggfunc",
			Args: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
				parser.NewIntegerValue(0),
			},
			AnalyticClause: parser.AnalyticClause{
				PartitionClause: parser.PartitionClause{
					Values: []parser.QueryExpression{
						parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
					},
				},
			},
		},
		Error: "variable @undefined is undeclared",
	},
	{
		Name: "Analyze Not Exist Function Error",
		View: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(1),
				}),
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(2),
				}),
			},
		},
		Function: parser.AnalyticFunction{
			Name: "notexist",
		},
		Error: "function notexist does not exist",
	},
}

func TestAnalyze(t *testing.T) {
	defer initFlag(TestTx.Flags)

	for _, v := range analyzeTests {
		if 0 < v.CPU {
			TestTx.Flags.CPU = v.CPU
		} else {
			TestTx.Flags.CPU = 1
		}

		ctx := context.Background()
		if v.Scope == nil {
			v.Scope = NewReferenceScope(TestTx)
		}

		err := Analyze(ctx, v.Scope, v.View, v.Function, v.PartitionIndices)
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

type analyticFunctionCheckArgsLenTests struct {
	Name     string
	Function parser.AnalyticFunction
	Error    string
}

func testAnalyticFunctionCheckArgsLenTests(t *testing.T, fn AnalyticFunction, tests []analyticFunctionCheckArgsLenTests) {
	for _, v := range tests {
		err := fn.CheckArgsLen(v.Function)
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
	}
}

type analyticFunctionExecuteTests struct {
	Name       string
	Items      Partition
	SortValues map[int]SortValues
	Function   parser.AnalyticFunction
	Result     map[int]value.Primary
	Error      string
}

var analyticFunctionTestScope = GenerateReferenceScope(nil, nil, time.Time{}, []ReferenceRecord{
	{
		view: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(100),
				}),
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(200),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewNull(),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(200),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(300),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(500),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(800),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewNull(),
				}),
			},
		},
		recordIndex: 0,
		cache:       NewFieldIndexCache(10, LimitToUseFieldIndexSliceChache),
	},
})

func testAnalyticFunctionExecute(t *testing.T, fn AnalyticFunction, tests []analyticFunctionExecuteTests) {
	for _, v := range tests {
		if v.SortValues != nil {
			list := make([]SortValues, analyticFunctionTestScope.Records[0].view.RecordLen())
			for i, v := range v.SortValues {
				list[i] = v
			}
			analyticFunctionTestScope.Records[0].view.sortValuesInEachRecord = list
		} else {
			analyticFunctionTestScope.Records[0].view.sortValuesInEachRecord = nil
		}

		result, err := fn.Execute(context.Background(), analyticFunctionTestScope, v.Items, v.Function)
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

var rowNumberCheckArgsLenTests = []analyticFunctionCheckArgsLenTests{
	{
		Name: "RowNumber CheckArgsLen",
		Function: parser.AnalyticFunction{
			Name: "row_number",
			Args: nil,
		},
	},
	{
		Name: "RowNumber CheckArgsLen Error",
		Function: parser.AnalyticFunction{
			Name: "row_number",
			Args: []parser.QueryExpression{
				parser.NewIntegerValue(1),
			},
		},
		Error: "function row_number takes no argument",
	},
}

func TestRowNumber_CheckArgsLen(t *testing.T) {
	testAnalyticFunctionCheckArgsLenTests(t, RowNumber{}, rowNumberCheckArgsLenTests)
}

var rowNumberExecuteTests = []analyticFunctionExecuteTests{
	{
		Name:  "RowNumber Execute",
		Items: Partition{2, 4, 1, 3, 5},
		Function: parser.AnalyticFunction{
			Name: "row_number",
		},
		Result: map[int]value.Primary{
			2: value.NewInteger(1),
			4: value.NewInteger(2),
			1: value.NewInteger(3),
			3: value.NewInteger(4),
			5: value.NewInteger(5),
		},
	},
}

func TestRowNumber_Execute(t *testing.T) {
	testAnalyticFunctionExecute(t, RowNumber{}, rowNumberExecuteTests)
}

var rankCheckArgsLenTests = []analyticFunctionCheckArgsLenTests{
	{
		Name: "Rank CheckArgsLen Error",
		Function: parser.AnalyticFunction{
			Name: "rank",
			Args: []parser.QueryExpression{
				parser.NewIntegerValue(1),
			},
		},
		Error: "function rank takes no argument",
	},
}

func TestRank_CheckArgsLen(t *testing.T) {
	testAnalyticFunctionCheckArgsLenTests(t, Rank{}, rankCheckArgsLenTests)
}

var rankExecuteTests = []analyticFunctionExecuteTests{
	{
		Name:  "Rank Execute",
		Items: Partition{2, 4, 1, 3, 5},
		SortValues: map[int]SortValues{
			2: {NewSortValue(value.NewString("1"), TestTx.Flags)},
			4: {NewSortValue(value.NewString("1"), TestTx.Flags)},
			1: {NewSortValue(value.NewString("2"), TestTx.Flags)},
			3: {NewSortValue(value.NewString("2"), TestTx.Flags)},
			5: {NewSortValue(value.NewString("3"), TestTx.Flags)},
		},
		Function: parser.AnalyticFunction{
			Name: "rank",
		},
		Result: map[int]value.Primary{
			2: value.NewInteger(1),
			4: value.NewInteger(1),
			1: value.NewInteger(3),
			3: value.NewInteger(3),
			5: value.NewInteger(5),
		},
	},
}

func TestRank_Execute(t *testing.T) {
	testAnalyticFunctionExecute(t, Rank{}, rankExecuteTests)
}

var denseRankCheckArgsLenTests = []analyticFunctionCheckArgsLenTests{
	{
		Name: "DenseRank CheckArgsLen Error",
		Function: parser.AnalyticFunction{
			Name: "dense_rank",
			Args: []parser.QueryExpression{
				parser.NewIntegerValue(1),
			},
		},
		Error: "function dense_rank takes no argument",
	},
}

func TestDenseRank_CheckArgsLen(t *testing.T) {
	testAnalyticFunctionCheckArgsLenTests(t, DenseRank{}, denseRankCheckArgsLenTests)
}

var denseRankExecuteTests = []analyticFunctionExecuteTests{
	{
		Name:  "DenseRank Execute",
		Items: Partition{2, 4, 1, 3, 5},
		SortValues: map[int]SortValues{
			2: {NewSortValue(value.NewString("1"), TestTx.Flags)},
			4: {NewSortValue(value.NewString("1"), TestTx.Flags)},
			1: {NewSortValue(value.NewString("2"), TestTx.Flags)},
			3: {NewSortValue(value.NewString("2"), TestTx.Flags)},
			5: {NewSortValue(value.NewString("3"), TestTx.Flags)},
		},
		Function: parser.AnalyticFunction{
			Name: "dense_rank",
		},
		Result: map[int]value.Primary{
			2: value.NewInteger(1),
			4: value.NewInteger(1),
			1: value.NewInteger(2),
			3: value.NewInteger(2),
			5: value.NewInteger(3),
		},
	},
}

func TestDenseRank_Execute(t *testing.T) {
	testAnalyticFunctionExecute(t, DenseRank{}, denseRankExecuteTests)
}

var cumeDistCheckArgsLenTests = []analyticFunctionCheckArgsLenTests{
	{
		Name: "CumeDist CheckArgsLen Error",
		Function: parser.AnalyticFunction{
			Name: "cume_dist",
			Args: []parser.QueryExpression{
				parser.NewIntegerValue(1),
			},
		},
		Error: "function cume_dist takes no argument",
	},
}

func TestCumeDist_CheckArgsLen(t *testing.T) {
	testAnalyticFunctionCheckArgsLenTests(t, CumeDist{}, cumeDistCheckArgsLenTests)
}

var cumeDistExecuteTests = []analyticFunctionExecuteTests{
	{
		Name:  "CumeDist Execute",
		Items: Partition{2, 4, 1, 3},
		SortValues: map[int]SortValues{
			2: {NewSortValue(value.NewString("1"), TestTx.Flags)},
			4: {NewSortValue(value.NewString("2"), TestTx.Flags)},
			1: {NewSortValue(value.NewString("2"), TestTx.Flags)},
			3: {NewSortValue(value.NewString("3"), TestTx.Flags)},
		},
		Function: parser.AnalyticFunction{
			Name: "cume_dist",
		},
		Result: map[int]value.Primary{
			2: value.NewFloat(0.25),
			4: value.NewFloat(0.75),
			1: value.NewFloat(0.75),
			3: value.NewFloat(1),
		},
	},
}

func TestCumeDist_Execute(t *testing.T) {
	testAnalyticFunctionExecute(t, CumeDist{}, cumeDistExecuteTests)
}

var percentRankCheckArgsLenTests = []analyticFunctionCheckArgsLenTests{
	{
		Name: "PercentRank CheckArgsLen Error",
		Function: parser.AnalyticFunction{
			Name: "percent_rank",
			Args: []parser.QueryExpression{
				parser.NewIntegerValue(1),
			},
		},
		Error: "function percent_rank takes no argument",
	},
}

func TestPercentRank_CheckArgsLen(t *testing.T) {
	testAnalyticFunctionCheckArgsLenTests(t, PercentRank{}, percentRankCheckArgsLenTests)
}

var percentRankExecuteTests = []analyticFunctionExecuteTests{
	{
		Name:  "PercentRank Execute",
		Items: Partition{2, 4, 1, 3, 5},
		SortValues: map[int]SortValues{
			2: {NewSortValue(value.NewString("1"), TestTx.Flags)},
			4: {NewSortValue(value.NewString("2"), TestTx.Flags)},
			1: {NewSortValue(value.NewString("2"), TestTx.Flags)},
			3: {NewSortValue(value.NewString("3"), TestTx.Flags)},
			5: {NewSortValue(value.NewString("4"), TestTx.Flags)},
		},
		Function: parser.AnalyticFunction{
			Name: "percent_rank",
		},
		Result: map[int]value.Primary{
			2: value.NewFloat(0),
			4: value.NewFloat(0.25),
			1: value.NewFloat(0.25),
			3: value.NewFloat(0.75),
			5: value.NewFloat(1),
		},
	},
}

func TestPercentRank_Execute(t *testing.T) {
	testAnalyticFunctionExecute(t, PercentRank{}, percentRankExecuteTests)
}

var nTileCheckArgsLenTests = []analyticFunctionCheckArgsLenTests{
	{
		Name: "NTile CheckArgsLen Error",
		Function: parser.AnalyticFunction{
			Name: "ntile",
		},
		Error: "function ntile takes exactly 1 argument",
	},
}

func TestNTile_CheckArgsLen(t *testing.T) {
	testAnalyticFunctionCheckArgsLenTests(t, NTile{}, nTileCheckArgsLenTests)
}

var ntileValueExecuteTests = []analyticFunctionExecuteTests{
	{
		Name:  "NTile Execute",
		Items: Partition{1, 2, 3, 4, 5, 6, 7},
		Function: parser.AnalyticFunction{
			Name: "ntile",
			Args: []parser.QueryExpression{
				parser.NewIntegerValue(3),
			},
		},
		Result: map[int]value.Primary{
			1: value.NewInteger(1),
			2: value.NewInteger(1),
			3: value.NewInteger(1),
			4: value.NewInteger(2),
			5: value.NewInteger(2),
			6: value.NewInteger(3),
			7: value.NewInteger(3),
		},
	},
	{
		Name:  "NTile Execute Greater Tile Number",
		Items: Partition{1, 2},
		Function: parser.AnalyticFunction{
			Name: "ntile",
			Args: []parser.QueryExpression{
				parser.NewIntegerValue(3),
			},
		},
		Result: map[int]value.Primary{
			1: value.NewInteger(1),
			2: value.NewInteger(2),
		},
	},
	{
		Name:  "NTile Execute Argument Evaluation Error",
		Items: Partition{1, 2, 3, 4, 5, 6, 7},
		Function: parser.AnalyticFunction{
			Name: "ntile",
			Args: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			},
		},
		Error: "the first argument must be an integer for function ntile",
	},
	{
		Name:  "NTile Execute Argument Type Error",
		Items: Partition{1, 2, 3, 4, 5, 6, 7},
		Function: parser.AnalyticFunction{
			Name: "ntile",
			Args: []parser.QueryExpression{
				parser.NewNullValue(),
			},
		},
		Error: "the first argument must be an integer for function ntile",
	},
	{
		Name:  "NTile Execute Argument Value Error",
		Items: Partition{1, 2, 3, 4, 5, 6, 7},
		Function: parser.AnalyticFunction{
			Name: "ntile",
			Args: []parser.QueryExpression{
				parser.NewIntegerValue(0),
			},
		},
		Error: "the first argument must be greater than 0 for function ntile",
	},
}

func TestNTile_Execute(t *testing.T) {
	testAnalyticFunctionExecute(t, NTile{}, ntileValueExecuteTests)
}

var firstValueCheckArgsLenTests = []analyticFunctionCheckArgsLenTests{
	{
		Name: "FirstValue CheckArgsLen Error",
		Function: parser.AnalyticFunction{
			Name: "first_value",
		},
		Error: "function first_value takes exactly 1 argument",
	},
}

func TestFirstValue_CheckArgsLen(t *testing.T) {
	testAnalyticFunctionCheckArgsLenTests(t, FirstValue{}, firstValueCheckArgsLenTests)
}

var firstValueExecuteTests = []analyticFunctionExecuteTests{
	{
		Name:  "FirstValue Execute",
		Items: Partition{2, 3, 4, 5, 6, 7},
		Function: parser.AnalyticFunction{
			Name: "first_value",
			Args: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
			},
		},
		Result: map[int]value.Primary{
			2: value.NewNull(),
			3: value.NewNull(),
			4: value.NewNull(),
			5: value.NewNull(),
			6: value.NewNull(),
			7: value.NewNull(),
		},
	},
	{
		Name:  "FirstValue Execute IgnoreNulls",
		Items: Partition{2, 3, 4, 5, 6, 7},
		Function: parser.AnalyticFunction{
			Name: "first_value",
			Args: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
			},
			IgnoreType: parser.Token{Token: parser.NULLS, Literal: "nulls"},
		},
		Result: map[int]value.Primary{
			2: value.NewInteger(200),
			3: value.NewInteger(200),
			4: value.NewInteger(200),
			5: value.NewInteger(200),
			6: value.NewInteger(200),
			7: value.NewInteger(200),
		},
	},
	{
		Name:  "FirstValue Execute Argument Value Error",
		Items: Partition{2, 3},
		Function: parser.AnalyticFunction{
			Name: "first_value",
			Args: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			},
		},
		Error: "field notexist does not exist",
	},
}

func TestFirstValue_Execute(t *testing.T) {
	testAnalyticFunctionExecute(t, FirstValue{}, firstValueExecuteTests)
}

var lastValueCheckArgsLenTests = []analyticFunctionCheckArgsLenTests{
	{
		Name: "LastValue CheckArgsLen Error",
		Function: parser.AnalyticFunction{
			Name: "last_value",
		},
		Error: "function last_value takes exactly 1 argument",
	},
}

func TestLastValue_CheckArgsLen(t *testing.T) {
	testAnalyticFunctionCheckArgsLenTests(t, LastValue{}, lastValueCheckArgsLenTests)
}

var lastValueExecuteTests = []analyticFunctionExecuteTests{
	{
		Name:  "LastValue Execute",
		Items: Partition{2, 3, 4, 5, 6, 7},
		Function: parser.AnalyticFunction{
			Name: "last_value",
			Args: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
			},
		},
		Result: map[int]value.Primary{
			2: value.NewNull(),
			3: value.NewNull(),
			4: value.NewNull(),
			5: value.NewNull(),
			6: value.NewNull(),
			7: value.NewNull(),
		},
	},
	{
		Name:  "LastValue Execute IgnoreNulls",
		Items: Partition{2, 3, 4, 5, 6, 7},
		Function: parser.AnalyticFunction{
			Name: "last_value",
			Args: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
			},
			IgnoreType: parser.Token{Token: parser.NULLS, Literal: "nulls"},
		},
		Result: map[int]value.Primary{
			2: value.NewInteger(800),
			3: value.NewInteger(800),
			4: value.NewInteger(800),
			5: value.NewInteger(800),
			6: value.NewInteger(800),
			7: value.NewInteger(800),
		},
	},
}

func TestLastValue_Execute(t *testing.T) {
	testAnalyticFunctionExecute(t, LastValue{}, lastValueExecuteTests)
}

var nthValueCheckArgsLenTests = []analyticFunctionCheckArgsLenTests{
	{
		Name: "NthValue CheckArgsLen Error",
		Function: parser.AnalyticFunction{
			Name: "nth_value",
		},
		Error: "function nth_value takes exactly 2 arguments",
	},
}

func TestNthValue_CheckArgsLen(t *testing.T) {
	testAnalyticFunctionCheckArgsLenTests(t, NthValue{}, nthValueCheckArgsLenTests)
}

var nthValueExecuteTests = []analyticFunctionExecuteTests{
	{
		Name:  "NthValue Execute",
		Items: Partition{2, 3, 4, 5, 6, 7},
		Function: parser.AnalyticFunction{
			Name: "nth_value",
			Args: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
				parser.NewIntegerValue(2),
			},
		},
		Result: map[int]value.Primary{
			2: value.NewInteger(200),
			3: value.NewInteger(200),
			4: value.NewInteger(200),
			5: value.NewInteger(200),
			6: value.NewInteger(200),
			7: value.NewInteger(200),
		},
	},
	{
		Name:  "NthValue with Start Specified Windowing Clause Execute",
		Items: Partition{2, 3, 4, 5, 6, 7},
		Function: parser.AnalyticFunction{
			Name: "nth_value",
			Args: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
				parser.NewIntegerValue(2),
			},
			AnalyticClause: parser.AnalyticClause{
				OrderByClause: parser.OrderByClause{
					Items: []parser.QueryExpression{
						parser.OrderItem{Value: parser.Identifier{Literal: "column2"}},
					},
				},
				WindowingClause: parser.WindowingClause{
					FrameLow: parser.WindowFramePosition{
						Direction: parser.Token{Token: parser.PRECEDING},
						Offset:    2,
					},
				},
			},
		},
		Result: map[int]value.Primary{
			2: value.NewNull(),
			3: value.NewInteger(200),
			4: value.NewInteger(200),
			5: value.NewInteger(300),
			6: value.NewInteger(500),
			7: value.NewInteger(800),
		},
	},
	{
		Name:  "NthValue with Rows Specified Windowing Clause Execute",
		Items: Partition{2, 3, 4, 5, 6, 7},
		Function: parser.AnalyticFunction{
			Name: "nth_value",
			Args: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
				parser.NewIntegerValue(2),
			},
			AnalyticClause: parser.AnalyticClause{
				OrderByClause: parser.OrderByClause{
					Items: []parser.QueryExpression{
						parser.OrderItem{Value: parser.Identifier{Literal: "column2"}},
					},
				},
				WindowingClause: parser.WindowingClause{
					FrameLow: parser.WindowFramePosition{
						Direction: parser.Token{Token: parser.CURRENT},
					},
					FrameHigh: parser.WindowFramePosition{
						Direction: parser.Token{Token: parser.FOLLOWING},
						Unbounded: parser.Token{Token: parser.UNBOUNDED},
					},
				},
			},
		},
		Result: map[int]value.Primary{
			2: value.NewInteger(200),
			3: value.NewInteger(300),
			4: value.NewInteger(500),
			5: value.NewInteger(800),
			6: value.NewNull(),
			7: value.NewNull(),
		},
	},
	{
		Name:  "NthValue with Rows Specified Windowing Clause Execute 2",
		Items: Partition{2, 3, 4, 5, 6, 7},
		Function: parser.AnalyticFunction{
			Name: "nth_value",
			Args: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
				parser.NewIntegerValue(2),
			},
			AnalyticClause: parser.AnalyticClause{
				OrderByClause: parser.OrderByClause{
					Items: []parser.QueryExpression{
						parser.OrderItem{Value: parser.Identifier{Literal: "column2"}},
					},
				},
				WindowingClause: parser.WindowingClause{
					FrameLow: parser.WindowFramePosition{
						Direction: parser.Token{Token: parser.PRECEDING},
						Unbounded: parser.Token{Token: parser.UNBOUNDED},
					},
					FrameHigh: parser.WindowFramePosition{
						Direction: parser.Token{Token: parser.FOLLOWING},
						Offset:    2,
					},
				},
			},
		},
		Result: map[int]value.Primary{
			2: value.NewInteger(200),
			3: value.NewInteger(200),
			4: value.NewInteger(200),
			5: value.NewInteger(200),
			6: value.NewInteger(200),
			7: value.NewInteger(200),
		},
	},
	{
		Name:  "NthValue with Rows Specified Windowing Clause Execute 3",
		Items: Partition{2, 3, 4, 5, 6, 7},
		Function: parser.AnalyticFunction{
			Name: "nth_value",
			Args: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
				parser.NewIntegerValue(2),
			},
			AnalyticClause: parser.AnalyticClause{
				OrderByClause: parser.OrderByClause{
					Items: []parser.QueryExpression{
						parser.OrderItem{Value: parser.Identifier{Literal: "column2"}},
					},
				},
				WindowingClause: parser.WindowingClause{
					FrameLow: parser.WindowFramePosition{
						Direction: parser.Token{Token: parser.PRECEDING},
						Unbounded: parser.Token{Token: parser.UNBOUNDED},
					},
					FrameHigh: parser.WindowFramePosition{
						Direction: parser.Token{Token: parser.FOLLOWING},
						Unbounded: parser.Token{Token: parser.UNBOUNDED},
					},
				},
			},
		},
		Result: map[int]value.Primary{
			2: value.NewInteger(200),
			3: value.NewInteger(200),
			4: value.NewInteger(200),
			5: value.NewInteger(200),
			6: value.NewInteger(200),
			7: value.NewInteger(200),
		},
	},
	{
		Name:  "NthValue with Default Windowing Clause Execute",
		Items: Partition{2, 3, 4, 5, 6, 7},
		Function: parser.AnalyticFunction{
			Name: "nth_value",
			Args: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
				parser.NewIntegerValue(2),
			},
			AnalyticClause: parser.AnalyticClause{
				OrderByClause: parser.OrderByClause{
					Items: []parser.QueryExpression{
						parser.OrderItem{Value: parser.Identifier{Literal: "column2"}},
					},
				},
			},
		},
		Result: map[int]value.Primary{
			2: value.NewNull(),
			3: value.NewInteger(200),
			4: value.NewInteger(200),
			5: value.NewInteger(200),
			6: value.NewInteger(200),
			7: value.NewInteger(200),
		},
	},
	{
		Name:  "NthValue Execute Second Argument Evaluation Error",
		Items: Partition{2, 3, 4, 5, 6, 7},
		Function: parser.AnalyticFunction{
			Name: "nth_value",
			Args: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
				parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			},
		},
		Error: "the second argument must be an integer for function nth_value",
	},
	{
		Name:  "NthValue Execute Second Argument Type Error",
		Items: Partition{2, 3, 4, 5, 6, 7},
		Function: parser.AnalyticFunction{
			Name: "nth_value",
			Args: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
				parser.NewNullValue(),
			},
		},
		Error: "the second argument must be an integer for function nth_value",
	},
	{
		Name:  "NthValue Execute Second Argument Value Error",
		Items: Partition{2, 3, 4, 5, 6, 7},
		Function: parser.AnalyticFunction{
			Name: "nth_value",
			Args: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
				parser.NewIntegerValue(0),
			},
		},
		Error: "the second argument must be greater than 0 for function nth_value",
	},
}

func TestNthValue_Execute(t *testing.T) {
	testAnalyticFunctionExecute(t, NthValue{}, nthValueExecuteTests)
}

var lagCheckArgsLenTests = []analyticFunctionCheckArgsLenTests{
	{
		Name: "Lag CheckArgsLen Too Little Error",
		Function: parser.AnalyticFunction{
			Name: "lag",
		},
		Error: "function lag takes at least 1 argument",
	},
	{
		Name: "Lag CheckArgsLen Too Many Error",
		Function: parser.AnalyticFunction{
			Name: "lag",
			Args: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
				parser.NewIntegerValue(2),
				parser.NewIntegerValue(0),
				parser.NewIntegerValue(0),
			},
		},
		Error: "function lag takes at most 3 arguments",
	},
}

func TestLag_CheckArgsLen(t *testing.T) {
	testAnalyticFunctionCheckArgsLenTests(t, Lag{}, lagCheckArgsLenTests)
}

var lagExecuteTests = []analyticFunctionExecuteTests{
	{
		Name:  "Lag Execute",
		Items: Partition{2, 3, 4, 5, 6, 7},
		Function: parser.AnalyticFunction{
			Name: "lag",
			Args: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
				parser.NewIntegerValue(2),
				parser.NewIntegerValue(0),
			},
		},
		Result: map[int]value.Primary{
			2: value.NewInteger(0),
			3: value.NewInteger(0),
			4: value.NewNull(),
			5: value.NewInteger(200),
			6: value.NewInteger(300),
			7: value.NewInteger(500),
		},
	},
	{
		Name:  "Lag Execute With Default Value",
		Items: Partition{2, 3, 4, 5, 6, 7},
		Function: parser.AnalyticFunction{
			Name: "lag",
			Args: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
			},
		},
		Result: map[int]value.Primary{
			2: value.NewNull(),
			3: value.NewNull(),
			4: value.NewInteger(200),
			5: value.NewInteger(300),
			6: value.NewInteger(500),
			7: value.NewInteger(800),
		},
	},
	{
		Name:  "Lag Execute With IgnoreNulls",
		Items: Partition{2, 3, 4, 5, 6, 7},
		Function: parser.AnalyticFunction{
			Name: "lag",
			Args: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
				parser.NewIntegerValue(2),
				parser.NewIntegerValue(0),
			},
			IgnoreType: parser.Token{Token: parser.NULLS, Literal: "nulls"},
		},
		Result: map[int]value.Primary{
			2: value.NewInteger(0),
			3: value.NewInteger(0),
			4: value.NewInteger(0),
			5: value.NewInteger(200),
			6: value.NewInteger(300),
			7: value.NewInteger(500),
		},
	},
	{
		Name:  "Lag Execute First Argument Evaluation Error",
		Items: Partition{2, 3, 4, 5, 6, 7},
		Function: parser.AnalyticFunction{
			Name: "lag",
			Args: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
				parser.NewIntegerValue(2),
				parser.NewIntegerValue(0),
			},
		},
		Error: "field notexist does not exist",
	},
	{
		Name:  "Lag Execute Second Argument Evaluation Error",
		Items: Partition{2, 3, 4, 5, 6, 7},
		Function: parser.AnalyticFunction{
			Name: "lag",
			Args: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
				parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
				parser.NewIntegerValue(0),
			},
		},
		Error: "the second argument must be an integer for function lag",
	},
	{
		Name:  "Lag Execute Second Argument Type Error",
		Items: Partition{2, 3, 4, 5, 6, 7},
		Function: parser.AnalyticFunction{
			Name: "lag",
			Args: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
				parser.NewNullValue(),
				parser.NewIntegerValue(0),
			},
		},
		Error: "the second argument must be an integer for function lag",
	},
	{
		Name:  "Lag Execute Third Argument Evaluation Error",
		Items: Partition{2, 3, 4, 5, 6, 7},
		Function: parser.AnalyticFunction{
			Name: "lag",
			Args: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
				parser.NewIntegerValue(2),
				parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			},
		},
		Error: "field notexist does not exist",
	},
}

func TestLag_Execute(t *testing.T) {
	testAnalyticFunctionExecute(t, Lag{}, lagExecuteTests)
}

var leadCheckArgsLenTests = []analyticFunctionCheckArgsLenTests{
	{
		Name: "Lead CheckArgsLen Too Little Error",
		Function: parser.AnalyticFunction{
			Name: "lead",
		},
		Error: "function lead takes at least 1 argument",
	},
}

func TestLead_CheckArgsLen(t *testing.T) {
	testAnalyticFunctionCheckArgsLenTests(t, Lead{}, leadCheckArgsLenTests)
}

var leadExecuteTests = []analyticFunctionExecuteTests{
	{
		Name:  "Lead Execute With Default Value",
		Items: Partition{2, 3, 4, 5, 6, 7},
		Function: parser.AnalyticFunction{
			Name: "lead",
			Args: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
			},
		},
		Result: map[int]value.Primary{
			2: value.NewInteger(200),
			3: value.NewInteger(300),
			4: value.NewInteger(500),
			5: value.NewInteger(800),
			6: value.NewNull(),
			7: value.NewNull(),
		},
	},
}

func TestLead_Execute(t *testing.T) {
	testAnalyticFunctionExecute(t, Lead{}, leadExecuteTests)
}

var analyticListAggCheckArgsLenTests = []analyticFunctionCheckArgsLenTests{
	{
		Name: "ListAgg CheckArgsLen Too Little Error",
		Function: parser.AnalyticFunction{
			Name: "listagg",
		},
		Error: "function listagg takes at least 1 argument",
	},
}

func TestAnalyticListAgg_CheckArgsLen(t *testing.T) {
	testAnalyticFunctionCheckArgsLenTests(t, AnalyticListAgg{}, analyticListAggCheckArgsLenTests)
}

var analyticListAggExecuteTests = []analyticFunctionExecuteTests{
	{
		Name:  "AnalyticListAgg Execute",
		Items: Partition{0, 1, 2, 3, 4},
		Function: parser.AnalyticFunction{
			Name: "listagg",
			Args: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
				parser.NewStringValue(","),
			},
		},
		Result: map[int]value.Primary{
			0: value.NewString("100,200,200,300"),
			1: value.NewString("100,200,200,300"),
			2: value.NewString("100,200,200,300"),
			3: value.NewString("100,200,200,300"),
			4: value.NewString("100,200,200,300"),
		},
	},
	{
		Name:  "AnalyticListAgg Execute With Default Value",
		Items: Partition{0, 1, 2, 3, 4},
		Function: parser.AnalyticFunction{
			Name: "listagg",
			Args: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
			},
		},
		Result: map[int]value.Primary{
			0: value.NewString("100200200300"),
			1: value.NewString("100200200300"),
			2: value.NewString("100200200300"),
			3: value.NewString("100200200300"),
			4: value.NewString("100200200300"),
		},
	},
	{
		Name:  "AnalyticListAgg Execute With Distinct",
		Items: Partition{0, 1, 2, 3, 4},
		Function: parser.AnalyticFunction{
			Name:     "listagg",
			Distinct: parser.Token{Token: parser.DISTINCT, Literal: "distinct"},
			Args: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
				parser.NewStringValue(","),
			},
		},
		Result: map[int]value.Primary{
			0: value.NewString("100,200,300"),
			1: value.NewString("100,200,300"),
			2: value.NewString("100,200,300"),
			3: value.NewString("100,200,300"),
			4: value.NewString("100,200,300"),
		},
	},
	{
		Name:  "AnalyticListAgg Execute First Argument Evaluation Error",
		Items: Partition{0, 1, 2, 3, 4},
		Function: parser.AnalyticFunction{
			Name: "listagg",
			Args: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
				parser.NewStringValue(","),
			},
		},
		Error: "field notexist does not exist",
	},
	{
		Name:  "AnalyticListAgg Execute Second Argument Evaluation Error",
		Items: Partition{0, 1, 2, 3, 4},
		Function: parser.AnalyticFunction{
			Name: "listagg",
			Args: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
				parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			},
		},
		Error: "the second argument must be a string for function listagg",
	},
	{
		Name:  "AnalyticListAgg Execute Second Argument Type Error",
		Items: Partition{0, 1, 2, 3, 4},
		Function: parser.AnalyticFunction{
			Name: "listagg",
			Args: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
				parser.NewNullValue(),
			},
		},
		Error: "the second argument must be a string for function listagg",
	},
}

func TestAnalyticListAgg_Execute(t *testing.T) {
	testAnalyticFunctionExecute(t, AnalyticListAgg{}, analyticListAggExecuteTests)
}

var analyticJsonAggCheckArgsLenTests = []analyticFunctionCheckArgsLenTests{
	{
		Name: "JsonAgg CheckArgsLen Too Little Error",
		Function: parser.AnalyticFunction{
			Name: "json_agg",
		},
		Error: "function json_agg takes exactly 1 argument",
	},
}

func TestAnalyticJsonAgg_CheckArgsLen(t *testing.T) {
	testAnalyticFunctionCheckArgsLenTests(t, AnalyticJsonAgg{}, analyticJsonAggCheckArgsLenTests)
}

var analyticJsonAggExecuteTests = []analyticFunctionExecuteTests{
	{
		Name:  "AnalyticJsonAgg Execute",
		Items: Partition{0, 1, 2, 3, 4},
		Function: parser.AnalyticFunction{
			Name: "json_agg",
			Args: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
			},
		},
		Result: map[int]value.Primary{
			0: value.NewString("[100,200,null,200,300]"),
			1: value.NewString("[100,200,null,200,300]"),
			2: value.NewString("[100,200,null,200,300]"),
			3: value.NewString("[100,200,null,200,300]"),
			4: value.NewString("[100,200,null,200,300]"),
		},
	},
	{
		Name:  "AnalyticJsonAgg Execute With Distinct",
		Items: Partition{0, 1, 2, 3, 4},
		Function: parser.AnalyticFunction{
			Name:     "json_agg",
			Distinct: parser.Token{Token: parser.DISTINCT, Literal: "distinct"},
			Args: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
			},
		},
		Result: map[int]value.Primary{
			0: value.NewString("[100,200,null,300]"),
			1: value.NewString("[100,200,null,300]"),
			2: value.NewString("[100,200,null,300]"),
			3: value.NewString("[100,200,null,300]"),
			4: value.NewString("[100,200,null,300]"),
		},
	},
	{
		Name:  "AnalyticJsonAgg Execute First Argument Evaluation Error",
		Items: Partition{0, 1, 2, 3, 4},
		Function: parser.AnalyticFunction{
			Name: "json_agg",
			Args: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			},
		},
		Error: "field notexist does not exist",
	},
}

func TestAnalyticJsonAgg_Execute(t *testing.T) {
	testAnalyticFunctionExecute(t, AnalyticJsonAgg{}, analyticJsonAggExecuteTests)
}
