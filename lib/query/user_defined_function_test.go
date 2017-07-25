package query

import (
	"reflect"
	"testing"

	"github.com/mithrandie/csvq/lib/parser"
)

var userDefinedFunctionMapDeclareTests = []struct {
	Name   string
	Expr   parser.FunctionDeclaration
	Result UserDefinedFunctionMap
	Error  string
}{
	{
		Name: "UserDefinedFunctionMap Declare",
		Expr: parser.FunctionDeclaration{
			Name: parser.Identifier{Literal: "userfunc"},
			Parameters: []parser.Variable{
				{Name: "@arg1"},
				{Name: "@arg2"},
			},
			Statements: []parser.Statement{
				parser.Print{Value: parser.Variable{Name: "@var1"}},
			},
		},
		Result: UserDefinedFunctionMap{
			"USERFUNC": &UserDefinedFunction{
				Name: parser.Identifier{Literal: "userfunc"},
				Parameters: []parser.Variable{
					{Name: "@arg1"},
					{Name: "@arg2"},
				},
				Statements: []parser.Statement{
					parser.Print{Value: parser.Variable{Name: "@var1"}},
				},
			},
		},
	},
	{
		Name: "UserDefinedFunctionMap Declare Redeclaration Error",
		Expr: parser.FunctionDeclaration{
			Name: parser.Identifier{Literal: "userfunc"},
			Parameters: []parser.Variable{
				{Name: "@arg1"},
				{Name: "@arg2"},
			},
			Statements: []parser.Statement{
				parser.Print{Value: parser.Variable{Name: "@var1"}},
			},
		},
		Error: "[L:- C:-] function userfunc is redeclared",
	},
	{
		Name: "UserDefinedFunctionMap Declare Duplicate with Built-in Function Error",
		Expr: parser.FunctionDeclaration{
			Name: parser.Identifier{Literal: "now"},
			Parameters: []parser.Variable{
				{Name: "@arg1"},
				{Name: "@arg2"},
			},
			Statements: []parser.Statement{
				parser.Print{Value: parser.Variable{Name: "@var1"}},
			},
		},
		Error: "[L:- C:-] function now is a built-in function",
	},
	{
		Name: "UserDefinedFunctionMap Declare Duplicate with Aggregate Function Error",
		Expr: parser.FunctionDeclaration{
			Name: parser.Identifier{Literal: "count"},
			Parameters: []parser.Variable{
				{Name: "@arg1"},
				{Name: "@arg2"},
			},
			Statements: []parser.Statement{
				parser.Print{Value: parser.Variable{Name: "@var1"}},
			},
		},
		Error: "[L:- C:-] function count is a built-in function",
	},
	{
		Name: "UserDefinedFunctionMap Declare Duplicate with Analytic Function Error",
		Expr: parser.FunctionDeclaration{
			Name: parser.Identifier{Literal: "row_number"},
			Parameters: []parser.Variable{
				{Name: "@arg1"},
				{Name: "@arg2"},
			},
			Statements: []parser.Statement{
				parser.Print{Value: parser.Variable{Name: "@var1"}},
			},
		},
		Error: "[L:- C:-] function row_number is a built-in function",
	},
	{
		Name: "UserDefinedFunctionMap Declare Duplicate with GroupConcat Function Error",
		Expr: parser.FunctionDeclaration{
			Name: parser.Identifier{Literal: "group_concat"},
			Parameters: []parser.Variable{
				{Name: "@arg1"},
				{Name: "@arg2"},
			},
			Statements: []parser.Statement{
				parser.Print{Value: parser.Variable{Name: "@var1"}},
			},
		},
		Error: "[L:- C:-] function group_concat is a built-in function",
	},
}

func TestUserDefinedFunctionMap_Declare(t *testing.T) {
	funcs := UserDefinedFunctionMap{}

	for _, v := range userDefinedFunctionMapDeclareTests {
		err := funcs.Declare(v.Expr)
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
		if !reflect.DeepEqual(funcs, v.Result) {
			t.Errorf("%s: result = %s, want %s", v.Name, funcs, v.Result)
		}
	}
}

var userDefinedFunctionMapGetTests = []struct {
	Name     string
	Function parser.Function
	Result   *UserDefinedFunction
	Error    string
}{
	{
		Name: "UserDefinedFunctionMap Get",
		Function: parser.Function{
			Name: "userfunc",
		},
		Result: &UserDefinedFunction{
			Name: parser.Identifier{Literal: "userfunc"},
			Parameters: []parser.Variable{
				{Name: "@arg1"},
				{Name: "@arg2"},
			},
			Statements: []parser.Statement{
				parser.Print{Value: parser.Variable{Name: "@var1"}},
			},
		},
	},
	{
		Name: "UserDefinedFunctionMap Get Not Exist Error",
		Function: parser.Function{
			Name: "notexist",
		},
		Error: "[L:- C:-] function notexist does not exist",
	},
}

func TestUserDefinedFunctionMap_Get(t *testing.T) {
	funcs := UserDefinedFunctionMap{
		"USERFUNC": &UserDefinedFunction{
			Name: parser.Identifier{Literal: "userfunc"},
			Parameters: []parser.Variable{
				{Name: "@arg1"},
				{Name: "@arg2"},
			},
			Statements: []parser.Statement{
				parser.Print{Value: parser.Variable{Name: "@var1"}},
			},
		},
	}

	for _, v := range userDefinedFunctionMapGetTests {
		result, err := funcs.Get(v.Function)
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
			t.Errorf("%s: result = %s, want %s", v.Name, result, v.Result)
		}
	}
}

var userDefinedFunctionExecuteTests = []struct {
	Name   string
	Func   *UserDefinedFunction
	Args   []parser.Primary
	Result parser.Primary
	Error  string
}{
	{
		Name: "UserDefinedFunction Execute",
		Func: &UserDefinedFunction{
			Name: parser.Identifier{Literal: "userfunc"},
			Parameters: []parser.Variable{
				{Name: "@arg1"},
				{Name: "@arg2"},
			},
			Statements: []parser.Statement{
				parser.VariableDeclaration{
					Assignments: []parser.Expression{
						parser.VariableAssignment{
							Variable: parser.Variable{Name: "@var2"},
							Value: parser.Arithmetic{
								LHS: parser.Arithmetic{
									LHS:      parser.Variable{Name: "@arg1"},
									RHS:      parser.Variable{Name: "@arg2"},
									Operator: '+',
								},
								RHS:      parser.Variable{Name: "@var1"},
								Operator: '+',
							},
						},
					},
				},
				parser.Return{
					Value: parser.Variable{Name: "@var2"},
				},
			},
		},
		Args: []parser.Primary{
			parser.NewInteger(2),
			parser.NewInteger(3),
		},
		Result: parser.NewInteger(6),
	},
	{
		Name: "UserDefinedFunction Execute No Return Statement",
		Func: &UserDefinedFunction{
			Name: parser.Identifier{Literal: "userfunc"},
			Parameters: []parser.Variable{
				{Name: "@arg1"},
				{Name: "@arg2"},
			},
			Statements: []parser.Statement{
				parser.VariableDeclaration{
					Assignments: []parser.Expression{
						parser.VariableAssignment{
							Variable: parser.Variable{Name: "@var2"},
							Value: parser.Arithmetic{
								LHS: parser.Arithmetic{
									LHS:      parser.Variable{Name: "@arg1"},
									RHS:      parser.Variable{Name: "@arg2"},
									Operator: '+',
								},
								RHS:      parser.Variable{Name: "@var1"},
								Operator: '+',
							},
						},
					},
				},
			},
		},
		Args: []parser.Primary{
			parser.NewInteger(2),
			parser.NewInteger(3),
		},
		Result: parser.NewNull(),
	},
	{
		Name: "UserDefinedFunction Execute Arguments Error",
		Func: &UserDefinedFunction{
			Name: parser.Identifier{Literal: "userfunc"},
			Parameters: []parser.Variable{
				{Name: "@arg1"},
				{Name: "@arg2"},
			},
			Statements: []parser.Statement{
				parser.VariableDeclaration{
					Assignments: []parser.Expression{
						parser.VariableAssignment{
							Variable: parser.Variable{Name: "@var2"},
							Value: parser.Arithmetic{
								LHS: parser.Arithmetic{
									LHS:      parser.Variable{Name: "@arg1"},
									RHS:      parser.Variable{Name: "@arg2"},
									Operator: '+',
								},
								RHS:      parser.Variable{Name: "@var1"},
								Operator: '+',
							},
						},
					},
				},
				parser.Return{
					Value: parser.Variable{Name: "@var2"},
				},
			},
		},
		Args: []parser.Primary{
			parser.NewInteger(2),
		},
		Error: "[L:- C:-] function userfunc takes 2 arguments",
	},
	{
		Name: "UserDefinedFunction Execute Parameter Duplicate Error",
		Func: &UserDefinedFunction{
			Name: parser.Identifier{Literal: "userfunc"},
			Parameters: []parser.Variable{
				{Name: "@arg1"},
				{Name: "@arg1"},
			},
			Statements: []parser.Statement{
				parser.VariableDeclaration{
					Assignments: []parser.Expression{
						parser.VariableAssignment{
							Variable: parser.Variable{Name: "@var2"},
							Value: parser.Subquery{
								Query: parser.SelectQuery{
									SelectEntity: parser.SelectEntity{
										SelectClause: parser.SelectClause{
											Fields: []parser.Expression{
												parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}}},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		Args: []parser.Primary{
			parser.NewInteger(2),
			parser.NewInteger(3),
		},
		Error: "[L:- C:-] variable @arg1 is redeclared",
	},
	{
		Name: "UserDefinedFunction Execute Execution Error",
		Func: &UserDefinedFunction{
			Name: parser.Identifier{Literal: "userfunc"},
			Parameters: []parser.Variable{
				{Name: "@arg1"},
				{Name: "@arg2"},
			},
			Statements: []parser.Statement{
				parser.VariableDeclaration{
					Assignments: []parser.Expression{
						parser.VariableAssignment{
							Variable: parser.Variable{Name: "@var2"},
							Value: parser.Subquery{
								Query: parser.SelectQuery{
									SelectEntity: parser.SelectEntity{
										SelectClause: parser.SelectClause{
											Fields: []parser.Expression{
												parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}}},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		Args: []parser.Primary{
			parser.NewInteger(2),
			parser.NewInteger(3),
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
}

func TestUserDefinedFunction_Execute(t *testing.T) {
	vars := Variables{
		"@var1": parser.NewInteger(1),
	}
	filter := NewFilter(
		[]Variables{vars},
		[]ViewMap{{}},
		[]CursorMap{{}},
	)

	for _, v := range userDefinedFunctionExecuteTests {
		result, err := v.Func.Execute(v.Args, filter)
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
			t.Errorf("%s: result = %s, want %s", v.Name, result, v.Result)
		}
	}
}
