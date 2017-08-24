package query

import (
	"reflect"
	"testing"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/parser"
)

var analyzeTests = []struct {
	Name             string
	CPU              int
	View             *View
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
			Records: []Record{
				NewRecord([]parser.Primary{
					parser.NewString("a"),
					parser.NewInteger(1),
				}),
				NewRecord([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(1),
				}),
				NewRecord([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(1),
				}),
				NewRecord([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(2),
				}),
				NewRecord([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(2),
				}),
				NewRecord([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(3),
				}),
				NewRecord([]parser.Primary{
					parser.NewString("a"),
					parser.NewInteger(2),
				}),
			},
			Filter: NewEmptyFilter(),
			recordSortValues: []SortValues{
				{NewSortValue(parser.NewInteger(1))},
				{NewSortValue(parser.NewInteger(1))},
				{NewSortValue(parser.NewInteger(1))},
				{NewSortValue(parser.NewInteger(2))},
				{NewSortValue(parser.NewInteger(2))},
				{NewSortValue(parser.NewInteger(3))},
				{NewSortValue(parser.NewInteger(2))},
			},
		},
		Function: parser.AnalyticFunction{
			Name: "rank",
			AnalyticClause: parser.AnalyticClause{
				Partition: parser.Partition{
					Values: []parser.Expression{
						parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
					},
				},
				OrderByClause: parser.OrderByClause{
					Items: []parser.Expression{
						parser.OrderItem{
							Value: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
						},
					},
				},
			},
		},
		PartitionIndices: []int{0},
		Result: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			Records: []Record{
				NewRecord([]parser.Primary{
					parser.NewString("a"),
					parser.NewInteger(1),
					parser.NewInteger(1),
				}),
				NewRecord([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(1),
					parser.NewInteger(1),
				}),
				NewRecord([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(1),
					parser.NewInteger(1),
				}),
				NewRecord([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(2),
					parser.NewInteger(3),
				}),
				NewRecord([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(2),
					parser.NewInteger(3),
				}),
				NewRecord([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(3),
					parser.NewInteger(5),
				}),
				NewRecord([]parser.Primary{
					parser.NewString("a"),
					parser.NewInteger(2),
					parser.NewInteger(2),
				}),
			},
			Filter: NewEmptyFilter(),
			sortValues: [][]*SortValue{
				{NewSortValue(parser.NewString("a")), nil},
				{NewSortValue(parser.NewString("b")), nil},
				{NewSortValue(parser.NewString("b")), nil},
				{NewSortValue(parser.NewString("b")), nil},
				{NewSortValue(parser.NewString("b")), nil},
				{NewSortValue(parser.NewString("b")), nil},
				{NewSortValue(parser.NewString("a")), nil},
			},
			recordSortValues: []SortValues{
				{NewSortValue(parser.NewInteger(1))},
				{NewSortValue(parser.NewInteger(1))},
				{NewSortValue(parser.NewInteger(1))},
				{NewSortValue(parser.NewInteger(2))},
				{NewSortValue(parser.NewInteger(2))},
				{NewSortValue(parser.NewInteger(3))},
				{NewSortValue(parser.NewInteger(2))},
			},
		},
	},
	{
		Name: "Analyze AnalyticFunction Empty Record",
		CPU:  3,
		View: &View{
			Header:           NewHeader("table1", []string{"column1", "column2"}),
			Records:          []Record{},
			Filter:           NewEmptyFilter(),
			recordSortValues: []SortValues{},
		},
		Function: parser.AnalyticFunction{
			Name: "rank",
			AnalyticClause: parser.AnalyticClause{
				Partition: parser.Partition{
					Values: []parser.Expression{
						parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
					},
				},
				OrderByClause: parser.OrderByClause{
					Items: []parser.Expression{
						parser.OrderItem{
							Value: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
						},
					},
				},
			},
		},
		PartitionIndices: []int{0},
		Result: &View{
			Header:           NewHeader("table1", []string{"column1", "column2"}),
			Records:          []Record{},
			Filter:           NewEmptyFilter(),
			sortValues:       [][]*SortValue{},
			recordSortValues: []SortValues{},
		},
	},
	{
		Name: "Analyze AnalyticFunction Argument Length Error",
		View: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			Records: []Record{
				NewRecord([]parser.Primary{
					parser.NewString("a"),
					parser.NewInteger(1),
				}),
				NewRecord([]parser.Primary{
					parser.NewString("a"),
					parser.NewInteger(2),
				}),
			},
			Filter: NewEmptyFilter(),
		},
		Function: parser.AnalyticFunction{
			Name: "rank",
			Args: []parser.Expression{
				parser.NewIntegerValue(1),
			},
			AnalyticClause: parser.AnalyticClause{
				Partition: parser.Partition{
					Values: []parser.Expression{
						parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
					},
				},
				OrderByClause: parser.OrderByClause{
					Items: []parser.Expression{
						parser.OrderItem{
							Value: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
						},
					},
				},
			},
		},
		Error: "[L:- C:-] function rank takes no argument",
	},
	{
		Name: "Analyze AnalyticFunction Execution Error",
		View: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			Records: []Record{
				NewRecord([]parser.Primary{
					parser.NewString("a"),
					parser.NewInteger(1),
				}),
				NewRecord([]parser.Primary{
					parser.NewString("a"),
					parser.NewInteger(2),
				}),
			},
			Filter: NewEmptyFilter(),
		},
		Function: parser.AnalyticFunction{
			Name: "first_value",
			Args: []parser.Expression{
				parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			},
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
	{
		Name: "Analyze AggregateFunction",
		View: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			Records: []Record{
				NewRecord([]parser.Primary{
					parser.NewString("a"),
					parser.NewInteger(1),
				}),
				NewRecord([]parser.Primary{
					parser.NewString("a"),
					parser.NewInteger(2),
				}),
				NewRecord([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(1),
				}),
				NewRecord([]parser.Primary{
					parser.NewString("b"),
					parser.NewNull(),
				}),
				NewRecord([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(1),
				}),
			},
			Filter: NewEmptyFilter(),
		},
		Function: parser.AnalyticFunction{
			Name: "count",
			Args: []parser.Expression{
				parser.AllColumns{},
			},
			AnalyticClause: parser.AnalyticClause{
				Partition: parser.Partition{
					Values: []parser.Expression{
						parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
					},
				},
			},
		},
		PartitionIndices: []int{0},
		Result: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			Records: []Record{
				NewRecord([]parser.Primary{
					parser.NewString("a"),
					parser.NewInteger(1),
					parser.NewInteger(2),
				}),
				NewRecord([]parser.Primary{
					parser.NewString("a"),
					parser.NewInteger(2),
					parser.NewInteger(2),
				}),
				NewRecord([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(1),
					parser.NewInteger(3),
				}),
				NewRecord([]parser.Primary{
					parser.NewString("b"),
					parser.NewNull(),
					parser.NewInteger(3),
				}),
				NewRecord([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(1),
					parser.NewInteger(3),
				}),
			},
			Filter: NewEmptyFilter(),
			sortValues: [][]*SortValue{
				{NewSortValue(parser.NewString("a")), nil},
				{NewSortValue(parser.NewString("a")), nil},
				{NewSortValue(parser.NewString("b")), nil},
				{NewSortValue(parser.NewString("b")), nil},
				{NewSortValue(parser.NewString("b")), nil},
			},
		},
	},
	{
		Name: "Analyze AggregateFunction With Distinct",
		View: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			Records: []Record{
				NewRecord([]parser.Primary{
					parser.NewString("a"),
					parser.NewInteger(1),
				}),
				NewRecord([]parser.Primary{
					parser.NewString("a"),
					parser.NewInteger(2),
				}),
				NewRecord([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(1),
				}),
				NewRecord([]parser.Primary{
					parser.NewString("b"),
					parser.NewNull(),
				}),
				NewRecord([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(1),
				}),
			},
			Filter: NewEmptyFilter(),
			sortValues: [][]*SortValue{
				{NewSortValue(parser.NewString("a")), nil},
				{NewSortValue(parser.NewString("a")), nil},
				{NewSortValue(parser.NewString("b")), nil},
				{NewSortValue(parser.NewString("b")), nil},
				{NewSortValue(parser.NewString("b")), nil},
			},
		},
		Function: parser.AnalyticFunction{
			Name:     "count",
			Distinct: parser.Token{Token: parser.DISTINCT, Literal: "distinct"},
			Args: []parser.Expression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
			},
			AnalyticClause: parser.AnalyticClause{
				Partition: parser.Partition{
					Values: []parser.Expression{
						parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
					},
				},
			},
		},
		PartitionIndices: []int{0},
		Result: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			Records: []Record{
				NewRecord([]parser.Primary{
					parser.NewString("a"),
					parser.NewInteger(1),
					parser.NewInteger(2),
				}),
				NewRecord([]parser.Primary{
					parser.NewString("a"),
					parser.NewInteger(2),
					parser.NewInteger(2),
				}),
				NewRecord([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(1),
					parser.NewInteger(1),
				}),
				NewRecord([]parser.Primary{
					parser.NewString("b"),
					parser.NewNull(),
					parser.NewInteger(1),
				}),
				NewRecord([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(1),
					parser.NewInteger(1),
				}),
			},
			Filter: NewEmptyFilter(),
			sortValues: [][]*SortValue{
				{NewSortValue(parser.NewString("a")), nil},
				{NewSortValue(parser.NewString("a")), nil},
				{NewSortValue(parser.NewString("b")), nil},
				{NewSortValue(parser.NewString("b")), nil},
				{NewSortValue(parser.NewString("b")), nil},
			},
		},
	},
	{
		Name: "Analyze AggregateFunction Argument Length Error",
		View: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			Records: []Record{
				NewRecord([]parser.Primary{
					parser.NewString("a"),
					parser.NewInteger(1),
				}),
				NewRecord([]parser.Primary{
					parser.NewString("a"),
					parser.NewInteger(2),
				}),
			},
			Filter: NewEmptyFilter(),
		},
		Function: parser.AnalyticFunction{
			Name: "count",
			AnalyticClause: parser.AnalyticClause{
				Partition: parser.Partition{
					Values: []parser.Expression{
						parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
					},
				},
			},
		},
		Error: "[L:- C:-] function count takes exactly 1 argument",
	},
	{
		Name: "Analyze AggregateFunction Aggregate Value Error",
		View: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			Records: []Record{
				NewRecord([]parser.Primary{
					parser.NewString("a"),
					parser.NewInteger(1),
				}),
				NewRecord([]parser.Primary{
					parser.NewString("a"),
					parser.NewInteger(2),
				}),
			},
			Filter: NewEmptyFilter(),
		},
		Function: parser.AnalyticFunction{
			Name: "count",
			Args: []parser.Expression{
				parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			},
			AnalyticClause: parser.AnalyticClause{
				Partition: parser.Partition{
					Values: []parser.Expression{
						parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
					},
				},
			},
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
	{
		Name: "Analyze UserDefinedFunction",
		View: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			Records: []Record{
				NewRecord([]parser.Primary{
					parser.NewString("a"),
					parser.NewInteger(1),
				}),
				NewRecord([]parser.Primary{
					parser.NewString("a"),
					parser.NewInteger(2),
				}),
				NewRecord([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(1),
				}),
				NewRecord([]parser.Primary{
					parser.NewString("b"),
					parser.NewNull(),
				}),
				NewRecord([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(1),
				}),
			},
			Filter: &Filter{
				FunctionsList: UserDefinedFunctionsList{
					{
						"USERAGGFUNC": &UserDefinedFunction{
							Name:        parser.Identifier{Literal: "useraggfunc"},
							IsAggregate: true,
							Cursor:      parser.Identifier{Literal: "list"},
							Parameters: []parser.Variable{
								{Name: "@default"},
							},
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
									Cursor: parser.Identifier{Literal: "list"},
									Statements: []parser.Statement{
										parser.If{
											Condition: parser.Is{
												LHS: parser.Variable{Name: "@fetch"},
												RHS: parser.NewNullValue(),
											},
											Statements: []parser.Statement{
												parser.FlowControl{Token: parser.CONTINUE},
											},
										},
										parser.If{
											Condition: parser.Is{
												LHS: parser.Variable{Name: "@value"},
												RHS: parser.NewNullValue(),
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
		},
		Function: parser.AnalyticFunction{
			Name: "useraggfunc",
			Args: []parser.Expression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
				parser.NewIntegerValue(0),
			},
			AnalyticClause: parser.AnalyticClause{
				Partition: parser.Partition{
					Values: []parser.Expression{
						parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
					},
				},
			},
		},
		PartitionIndices: []int{0},
		Result: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			Records: []Record{
				NewRecord([]parser.Primary{
					parser.NewString("a"),
					parser.NewInteger(1),
					parser.NewInteger(2),
				}),
				NewRecord([]parser.Primary{
					parser.NewString("a"),
					parser.NewInteger(2),
					parser.NewInteger(2),
				}),
				NewRecord([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(1),
					parser.NewInteger(1),
				}),
				NewRecord([]parser.Primary{
					parser.NewString("b"),
					parser.NewNull(),
					parser.NewInteger(1),
				}),
				NewRecord([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(1),
					parser.NewInteger(1),
				}),
			},
			sortValues: [][]*SortValue{
				{NewSortValue(parser.NewString("a")), nil},
				{NewSortValue(parser.NewString("a")), nil},
				{NewSortValue(parser.NewString("b")), nil},
				{NewSortValue(parser.NewString("b")), nil},
				{NewSortValue(parser.NewString("b")), nil},
			},
			Filter: &Filter{
				FunctionsList: UserDefinedFunctionsList{
					{
						"USERAGGFUNC": &UserDefinedFunction{
							Name:        parser.Identifier{Literal: "useraggfunc"},
							IsAggregate: true,
							Cursor:      parser.Identifier{Literal: "list"},
							Parameters: []parser.Variable{
								{Name: "@default"},
							},
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
									Cursor: parser.Identifier{Literal: "list"},
									Statements: []parser.Statement{
										parser.If{
											Condition: parser.Is{
												LHS: parser.Variable{Name: "@fetch"},
												RHS: parser.NewNullValue(),
											},
											Statements: []parser.Statement{
												parser.FlowControl{Token: parser.CONTINUE},
											},
										},
										parser.If{
											Condition: parser.Is{
												LHS: parser.Variable{Name: "@value"},
												RHS: parser.NewNullValue(),
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
		},
	},
	{
		Name: "Analyze UserDefinedFunction Argument Length Error",
		View: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			Records: []Record{
				NewRecord([]parser.Primary{
					parser.NewString("a"),
					parser.NewInteger(1),
				}),
				NewRecord([]parser.Primary{
					parser.NewString("a"),
					parser.NewInteger(2),
				}),
			},
			Filter: &Filter{
				FunctionsList: UserDefinedFunctionsList{
					{
						"USERAGGFUNC": &UserDefinedFunction{
							Name:        parser.Identifier{Literal: "useraggfunc"},
							IsAggregate: true,
							Cursor:      parser.Identifier{Literal: "list"},
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
		},
		Function: parser.AnalyticFunction{
			Name: "useraggfunc",
			AnalyticClause: parser.AnalyticClause{
				Partition: parser.Partition{
					Values: []parser.Expression{
						parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
					},
				},
			},
		},
		Error: "[L:- C:-] function useraggfunc takes exactly 2 arguments",
	},
	{
		Name: "Analyze UserDefinedFunction Argument Value Error",
		View: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			Records: []Record{
				NewRecord([]parser.Primary{
					parser.NewString("a"),
					parser.NewInteger(1),
				}),
				NewRecord([]parser.Primary{
					parser.NewString("a"),
					parser.NewInteger(2),
				}),
			},
			Filter: &Filter{
				FunctionsList: UserDefinedFunctionsList{
					{
						"USERAGGFUNC": &UserDefinedFunction{
							Name:        parser.Identifier{Literal: "useraggfunc"},
							IsAggregate: true,
							Cursor:      parser.Identifier{Literal: "list"},
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
		},
		Function: parser.AnalyticFunction{
			Name: "useraggfunc",
			Args: []parser.Expression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
				parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			},
			AnalyticClause: parser.AnalyticClause{
				Partition: parser.Partition{
					Values: []parser.Expression{
						parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
					},
				},
			},
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
	{
		Name: "Analyze UserDefinedFunction Execution Error",
		View: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			Records: []Record{
				NewRecord([]parser.Primary{
					parser.NewString("a"),
					parser.NewInteger(1),
				}),
				NewRecord([]parser.Primary{
					parser.NewString("a"),
					parser.NewInteger(2),
				}),
			},
			Filter: &Filter{
				FunctionsList: UserDefinedFunctionsList{
					{
						"USERAGGFUNC": &UserDefinedFunction{
							Name:        parser.Identifier{Literal: "useraggfunc"},
							IsAggregate: true,
							Cursor:      parser.Identifier{Literal: "list"},
							Parameters: []parser.Variable{
								{Name: "@default"},
							},
							Statements: []parser.Statement{
								parser.Return{
									Value: parser.Variable{Name: "@undefined"},
								},
							},
						},
					},
				},
			},
		},
		Function: parser.AnalyticFunction{
			Name: "useraggfunc",
			Args: []parser.Expression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
				parser.NewIntegerValue(0),
			},
			AnalyticClause: parser.AnalyticClause{
				Partition: parser.Partition{
					Values: []parser.Expression{
						parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
					},
				},
			},
		},
		Error: "[L:- C:-] variable @undefined is undefined",
	},
	{
		Name: "Analyze Not Exist Function Error",
		View: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			Records: []Record{
				NewRecord([]parser.Primary{
					parser.NewString("a"),
					parser.NewInteger(1),
				}),
				NewRecord([]parser.Primary{
					parser.NewString("a"),
					parser.NewInteger(2),
				}),
			},
			Filter: NewEmptyFilter(),
		},
		Function: parser.AnalyticFunction{
			Name: "notexist",
		},
		Error: "[L:- C:-] function notexist does not exist",
	},
}

func TestAnalyze(t *testing.T) {
	flag := cmd.GetFlags()

	for _, v := range analyzeTests {
		if 0 < v.CPU {
			flag.CPU = v.CPU
		} else {
			v.CPU = 1
		}
		err := Analyze(v.View, v.Function, v.PartitionIndices)
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
			t.Errorf("%s: result = %q, want %q", v.Name, v.View, v.Result)
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
	Result     map[int]parser.Primary
	Error      string
}

var analyticFunctionTestFilter = &Filter{
	Records: []FilterRecord{
		{
			View: &View{
				Header: NewHeader("table1", []string{"column1", "column2"}),
				Records: []Record{
					NewRecord([]parser.Primary{
						parser.NewString("a"),
						parser.NewInteger(100),
					}),
					NewRecord([]parser.Primary{
						parser.NewString("a"),
						parser.NewInteger(200),
					}),
					NewRecord([]parser.Primary{
						parser.NewString("b"),
						parser.NewNull(),
					}),
					NewRecord([]parser.Primary{
						parser.NewString("b"),
						parser.NewInteger(200),
					}),
					NewRecord([]parser.Primary{
						parser.NewString("b"),
						parser.NewInteger(300),
					}),
					NewRecord([]parser.Primary{
						parser.NewString("b"),
						parser.NewInteger(500),
					}),
					NewRecord([]parser.Primary{
						parser.NewString("b"),
						parser.NewInteger(800),
					}),
					NewRecord([]parser.Primary{
						parser.NewString("b"),
						parser.NewNull(),
					}),
				},
				Filter: &Filter{
					VariablesList: VariablesList{{}},
					TempViewsList: TemporaryViewMapList{{}},
					CursorsList:   CursorMapList{{}},
					FunctionsList: UserDefinedFunctionsList{{}},
				},
			},
			RecordIndex: 0,
		},
	},
	VariablesList: VariablesList{{}},
	TempViewsList: TemporaryViewMapList{{}},
	CursorsList:   CursorMapList{{}},
	FunctionsList: UserDefinedFunctionsList{{}},
}

func testAnalyticFunctionExecute(t *testing.T, fn AnalyticFunction, tests []analyticFunctionExecuteTests) {
	for _, v := range tests {
		if v.SortValues != nil {
			list := make([]SortValues, analyticFunctionTestFilter.Records[0].View.RecordLen())
			for i, v := range v.SortValues {
				list[i] = v
			}
			analyticFunctionTestFilter.Records[0].View.recordSortValues = list
		} else {
			analyticFunctionTestFilter.Records[0].View.recordSortValues = nil
		}

		result, err := fn.Execute(v.Items, v.Function, analyticFunctionTestFilter)
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
			Args: []parser.Expression{
				parser.NewIntegerValue(1),
			},
		},
		Error: "[L:- C:-] function row_number takes no argument",
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
		Result: map[int]parser.Primary{
			2: parser.NewInteger(1),
			4: parser.NewInteger(2),
			1: parser.NewInteger(3),
			3: parser.NewInteger(4),
			5: parser.NewInteger(5),
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
			Args: []parser.Expression{
				parser.NewIntegerValue(1),
			},
		},
		Error: "[L:- C:-] function rank takes no argument",
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
			2: {NewSortValue(parser.NewString("1"))},
			4: {NewSortValue(parser.NewString("1"))},
			1: {NewSortValue(parser.NewString("2"))},
			3: {NewSortValue(parser.NewString("2"))},
			5: {NewSortValue(parser.NewString("3"))},
		},
		Function: parser.AnalyticFunction{
			Name: "rank",
		},
		Result: map[int]parser.Primary{
			2: parser.NewInteger(1),
			4: parser.NewInteger(1),
			1: parser.NewInteger(3),
			3: parser.NewInteger(3),
			5: parser.NewInteger(5),
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
			Args: []parser.Expression{
				parser.NewIntegerValue(1),
			},
		},
		Error: "[L:- C:-] function dense_rank takes no argument",
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
			2: {NewSortValue(parser.NewString("1"))},
			4: {NewSortValue(parser.NewString("1"))},
			1: {NewSortValue(parser.NewString("2"))},
			3: {NewSortValue(parser.NewString("2"))},
			5: {NewSortValue(parser.NewString("3"))},
		},
		Function: parser.AnalyticFunction{
			Name: "dense_rank",
		},
		Result: map[int]parser.Primary{
			2: parser.NewInteger(1),
			4: parser.NewInteger(1),
			1: parser.NewInteger(2),
			3: parser.NewInteger(2),
			5: parser.NewInteger(3),
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
			Args: []parser.Expression{
				parser.NewIntegerValue(1),
			},
		},
		Error: "[L:- C:-] function cume_dist takes no argument",
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
			2: {NewSortValue(parser.NewString("1"))},
			4: {NewSortValue(parser.NewString("2"))},
			1: {NewSortValue(parser.NewString("2"))},
			3: {NewSortValue(parser.NewString("3"))},
		},
		Function: parser.AnalyticFunction{
			Name: "cume_dist",
		},
		Result: map[int]parser.Primary{
			2: parser.NewFloat(0.25),
			4: parser.NewFloat(0.75),
			1: parser.NewFloat(0.75),
			3: parser.NewFloat(1),
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
			Args: []parser.Expression{
				parser.NewIntegerValue(1),
			},
		},
		Error: "[L:- C:-] function percent_rank takes no argument",
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
			2: {NewSortValue(parser.NewString("1"))},
			4: {NewSortValue(parser.NewString("2"))},
			1: {NewSortValue(parser.NewString("2"))},
			3: {NewSortValue(parser.NewString("3"))},
			5: {NewSortValue(parser.NewString("4"))},
		},
		Function: parser.AnalyticFunction{
			Name: "percent_rank",
		},
		Result: map[int]parser.Primary{
			2: parser.NewFloat(0),
			4: parser.NewFloat(0.25),
			1: parser.NewFloat(0.25),
			3: parser.NewFloat(0.75),
			5: parser.NewFloat(1),
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
		Error: "[L:- C:-] function ntile takes exactly 1 argument",
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
			Args: []parser.Expression{
				parser.NewIntegerValue(3),
			},
		},
		Result: map[int]parser.Primary{
			1: parser.NewInteger(1),
			2: parser.NewInteger(1),
			3: parser.NewInteger(1),
			4: parser.NewInteger(2),
			5: parser.NewInteger(2),
			6: parser.NewInteger(3),
			7: parser.NewInteger(3),
		},
	},
	{
		Name:  "NTile Execute Greater Tile Number",
		Items: Partition{1, 2},
		Function: parser.AnalyticFunction{
			Name: "ntile",
			Args: []parser.Expression{
				parser.NewIntegerValue(3),
			},
		},
		Result: map[int]parser.Primary{
			1: parser.NewInteger(1),
			2: parser.NewInteger(2),
		},
	},
	{
		Name:  "NTile Execute Argument Evaluation Error",
		Items: Partition{1, 2, 3, 4, 5, 6, 7},
		Function: parser.AnalyticFunction{
			Name: "ntile",
			Args: []parser.Expression{
				parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			},
		},
		Error: "[L:- C:-] the first argument must be an integer for function ntile",
	},
	{
		Name:  "NTile Execute Argument Type Error",
		Items: Partition{1, 2, 3, 4, 5, 6, 7},
		Function: parser.AnalyticFunction{
			Name: "ntile",
			Args: []parser.Expression{
				parser.NewNullValue(),
			},
		},
		Error: "[L:- C:-] the first argument must be an integer for function ntile",
	},
	{
		Name:  "NTile Execute Argument Value Error",
		Items: Partition{1, 2, 3, 4, 5, 6, 7},
		Function: parser.AnalyticFunction{
			Name: "ntile",
			Args: []parser.Expression{
				parser.NewIntegerValue(0),
			},
		},
		Error: "[L:- C:-] the first argument must be greater than 0 for function ntile",
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
		Error: "[L:- C:-] function first_value takes exactly 1 argument",
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
			Args: []parser.Expression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
			},
		},
		Result: map[int]parser.Primary{
			2: parser.NewNull(),
			3: parser.NewNull(),
			4: parser.NewNull(),
			5: parser.NewNull(),
			6: parser.NewNull(),
			7: parser.NewNull(),
		},
	},
	{
		Name:  "FirstValue Execute IgnoreNulls",
		Items: Partition{2, 3, 4, 5, 6, 7},
		Function: parser.AnalyticFunction{
			Name: "first_value",
			Args: []parser.Expression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
			},
			IgnoreNulls: true,
		},
		Result: map[int]parser.Primary{
			2: parser.NewInteger(200),
			3: parser.NewInteger(200),
			4: parser.NewInteger(200),
			5: parser.NewInteger(200),
			6: parser.NewInteger(200),
			7: parser.NewInteger(200),
		},
	},
	{
		Name:  "FirstValue Execute Argument Value Error",
		Items: Partition{2, 3},
		Function: parser.AnalyticFunction{
			Name: "first_value",
			Args: []parser.Expression{
				parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			},
		},
		Error: "[L:- C:-] field notexist does not exist",
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
		Error: "[L:- C:-] function last_value takes exactly 1 argument",
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
			Args: []parser.Expression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
			},
		},
		Result: map[int]parser.Primary{
			2: parser.NewNull(),
			3: parser.NewNull(),
			4: parser.NewNull(),
			5: parser.NewNull(),
			6: parser.NewNull(),
			7: parser.NewNull(),
		},
	},
	{
		Name:  "LastValue Execute IgnoreNulls",
		Items: Partition{2, 3, 4, 5, 6, 7},
		Function: parser.AnalyticFunction{
			Name: "last_value",
			Args: []parser.Expression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
			},
			IgnoreNulls: true,
		},
		Result: map[int]parser.Primary{
			2: parser.NewInteger(800),
			3: parser.NewInteger(800),
			4: parser.NewInteger(800),
			5: parser.NewInteger(800),
			6: parser.NewInteger(800),
			7: parser.NewInteger(800),
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
		Error: "[L:- C:-] function nth_value takes exactly 2 arguments",
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
			Args: []parser.Expression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
				parser.NewIntegerValue(2),
			},
		},
		Result: map[int]parser.Primary{
			2: parser.NewInteger(200),
			3: parser.NewInteger(200),
			4: parser.NewInteger(200),
			5: parser.NewInteger(200),
			6: parser.NewInteger(200),
			7: parser.NewInteger(200),
		},
	},
	{
		Name:  "NthValue Execute Second Argument Evaluation Error",
		Items: Partition{2, 3, 4, 5, 6, 7},
		Function: parser.AnalyticFunction{
			Name: "nth_value",
			Args: []parser.Expression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
				parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			},
		},
		Error: "[L:- C:-] the second argument must be an integer for function nth_value",
	},
	{
		Name:  "NthValue Execute Second Argument Type Error",
		Items: Partition{2, 3, 4, 5, 6, 7},
		Function: parser.AnalyticFunction{
			Name: "nth_value",
			Args: []parser.Expression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
				parser.NewNullValue(),
			},
		},
		Error: "[L:- C:-] the second argument must be an integer for function nth_value",
	},
	{
		Name:  "NthValue Execute Second Argument Value Error",
		Items: Partition{2, 3, 4, 5, 6, 7},
		Function: parser.AnalyticFunction{
			Name: "nth_value",
			Args: []parser.Expression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
				parser.NewIntegerValue(0),
			},
		},
		Error: "[L:- C:-] the second argument must be greater than 0 for function nth_value",
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
		Error: "[L:- C:-] function lag takes at least 1 argument",
	},
	{
		Name: "Lag CheckArgsLen Too Many Error",
		Function: parser.AnalyticFunction{
			Name: "lag",
			Args: []parser.Expression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
				parser.NewIntegerValue(2),
				parser.NewIntegerValue(0),
				parser.NewIntegerValue(0),
			},
		},
		Error: "[L:- C:-] function lag takes at most 3 arguments",
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
			Args: []parser.Expression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
				parser.NewIntegerValue(2),
				parser.NewIntegerValue(0),
			},
		},
		Result: map[int]parser.Primary{
			2: parser.NewInteger(0),
			3: parser.NewInteger(0),
			4: parser.NewNull(),
			5: parser.NewInteger(200),
			6: parser.NewInteger(300),
			7: parser.NewInteger(500),
		},
	},
	{
		Name:  "Lag Execute With Default Value",
		Items: Partition{2, 3, 4, 5, 6, 7},
		Function: parser.AnalyticFunction{
			Name: "lag",
			Args: []parser.Expression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
			},
		},
		Result: map[int]parser.Primary{
			2: parser.NewNull(),
			3: parser.NewNull(),
			4: parser.NewInteger(200),
			5: parser.NewInteger(300),
			6: parser.NewInteger(500),
			7: parser.NewInteger(800),
		},
	},
	{
		Name:  "Lag Execute With IgnoreNulls",
		Items: Partition{2, 3, 4, 5, 6, 7},
		Function: parser.AnalyticFunction{
			Name: "lag",
			Args: []parser.Expression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
				parser.NewIntegerValue(2),
				parser.NewIntegerValue(0),
			},
			IgnoreNulls: true,
		},
		Result: map[int]parser.Primary{
			2: parser.NewInteger(0),
			3: parser.NewInteger(0),
			4: parser.NewInteger(0),
			5: parser.NewInteger(200),
			6: parser.NewInteger(300),
			7: parser.NewInteger(500),
		},
	},
	{
		Name:  "Lag Execute First Argument Evaluation Error",
		Items: Partition{2, 3, 4, 5, 6, 7},
		Function: parser.AnalyticFunction{
			Name: "lag",
			Args: []parser.Expression{
				parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
				parser.NewIntegerValue(2),
				parser.NewIntegerValue(0),
			},
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
	{
		Name:  "Lag Execute Second Argument Evaluation Error",
		Items: Partition{2, 3, 4, 5, 6, 7},
		Function: parser.AnalyticFunction{
			Name: "lag",
			Args: []parser.Expression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
				parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
				parser.NewIntegerValue(0),
			},
		},
		Error: "[L:- C:-] the second argument must be an integer for function lag",
	},
	{
		Name:  "Lag Execute Second Argument Type Error",
		Items: Partition{2, 3, 4, 5, 6, 7},
		Function: parser.AnalyticFunction{
			Name: "lag",
			Args: []parser.Expression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
				parser.NewNullValue(),
				parser.NewIntegerValue(0),
			},
		},
		Error: "[L:- C:-] the second argument must be an integer for function lag",
	},
	{
		Name:  "Lag Execute Third Argument Type Error",
		Items: Partition{2, 3, 4, 5, 6, 7},
		Function: parser.AnalyticFunction{
			Name: "lag",
			Args: []parser.Expression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
				parser.NewIntegerValue(2),
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
			},
		},
		Error: "[L:- C:-] the third argument must be a primitive type for function lag",
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
		Error: "[L:- C:-] function lead takes at least 1 argument",
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
			Args: []parser.Expression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
			},
		},
		Result: map[int]parser.Primary{
			2: parser.NewInteger(200),
			3: parser.NewInteger(300),
			4: parser.NewInteger(500),
			5: parser.NewInteger(800),
			6: parser.NewNull(),
			7: parser.NewNull(),
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
		Error: "[L:- C:-] function listagg takes at least 1 argument",
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
			Args: []parser.Expression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
				parser.NewStringValue(","),
			},
		},
		Result: map[int]parser.Primary{
			0: parser.NewString("100,200,200,300"),
			1: parser.NewString("100,200,200,300"),
			2: parser.NewString("100,200,200,300"),
			3: parser.NewString("100,200,200,300"),
			4: parser.NewString("100,200,200,300"),
		},
	},
	{
		Name:  "AnalyticListAgg Execute With Default Value",
		Items: Partition{0, 1, 2, 3, 4},
		Function: parser.AnalyticFunction{
			Name: "listagg",
			Args: []parser.Expression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
			},
		},
		Result: map[int]parser.Primary{
			0: parser.NewString("100200200300"),
			1: parser.NewString("100200200300"),
			2: parser.NewString("100200200300"),
			3: parser.NewString("100200200300"),
			4: parser.NewString("100200200300"),
		},
	},
	{
		Name:  "AnalyticListAgg Execute With Distinct",
		Items: Partition{0, 1, 2, 3, 4},
		Function: parser.AnalyticFunction{
			Name:     "listagg",
			Distinct: parser.Token{Token: parser.DISTINCT, Literal: "distinct"},
			Args: []parser.Expression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
				parser.NewStringValue(","),
			},
		},
		Result: map[int]parser.Primary{
			0: parser.NewString("100,200,300"),
			1: parser.NewString("100,200,300"),
			2: parser.NewString("100,200,300"),
			3: parser.NewString("100,200,300"),
			4: parser.NewString("100,200,300"),
		},
	},
	{
		Name:  "AnalyticListAgg Execute First Argument Evaluation Error",
		Items: Partition{0, 1, 2, 3, 4},
		Function: parser.AnalyticFunction{
			Name: "listagg",
			Args: []parser.Expression{
				parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
				parser.NewStringValue(","),
			},
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
	{
		Name:  "AnalyticListAgg Execute Second Argument Evaluation Error",
		Items: Partition{0, 1, 2, 3, 4},
		Function: parser.AnalyticFunction{
			Name: "listagg",
			Args: []parser.Expression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
				parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			},
		},
		Error: "[L:- C:-] the second argument must be a string for function listagg",
	},
	{
		Name:  "AnalyticListAgg Execute Second Argument Type Error",
		Items: Partition{0, 1, 2, 3, 4},
		Function: parser.AnalyticFunction{
			Name: "listagg",
			Args: []parser.Expression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
				parser.NewNullValue(),
			},
		},
		Error: "[L:- C:-] the second argument must be a string for function listagg",
	},
}

func TestAnalyticListAgg_Execute(t *testing.T) {
	testAnalyticFunctionExecute(t, AnalyticListAgg{}, analyticListAggExecuteTests)
}
