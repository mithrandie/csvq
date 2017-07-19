package query

import (
	"reflect"
	"testing"

	"github.com/mithrandie/csvq/lib/parser"
)

var userFunctionMapDeclareTests = []struct {
	Name   string
	Expr   parser.FunctionDeclaration
	Result UserFunctionMap
	Error  string
}{
	{
		Name: "UserFunctionMap Declare",
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
		Result: UserFunctionMap{
			"USERFUNC": &UserFunction{
				Name: "userfunc",
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
		Name: "UserFunctionMap Declare Redeclaration Error",
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
		Error: "function userfunc is redeclared",
	},
	{
		Name: "UserFunctionMap Declare Duplicate with Built-in Function Error",
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
		Error: "function now is redeclared",
	},
	{
		Name: "UserFunctionMap Declare Duplicate with Aggregate Function Error",
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
		Error: "function count is redeclared",
	},
	{
		Name: "UserFunctionMap Declare Duplicate with Analytic Function Error",
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
		Error: "function row_number is redeclared",
	},
	{
		Name: "UserFunctionMap Declare Duplicate with GroupConcat Function Error",
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
		Error: "function group_concat is redeclared",
	},
}

func TestUserFunctionMap_Declare(t *testing.T) {
	funcs := UserFunctionMap{}

	for _, v := range userFunctionMapDeclareTests {
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

var userFunctionMapGetTests = []struct {
	Name   string
	Key    string
	Result *UserFunction
	Error  string
}{
	{
		Name: "UserFunctionMap Get",
		Key:  "userfunc",
		Result: &UserFunction{
			Name: "userfunc",
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
		Name:  "UserFunctionMap Get Not Exist Error",
		Key:   "notexist",
		Error: "function notexist does not exist",
	},
}

func TestUserFunctionMap_Get(t *testing.T) {
	funcs := UserFunctionMap{
		"USERFUNC": &UserFunction{
			Name: "userfunc",
			Parameters: []parser.Variable{
				{Name: "@arg1"},
				{Name: "@arg2"},
			},
			Statements: []parser.Statement{
				parser.Print{Value: parser.Variable{Name: "@var1"}},
			},
		},
	}

	for _, v := range userFunctionMapGetTests {
		result, err := funcs.Get(v.Key)
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

var userFunctionExecuteTests = []struct {
	Name   string
	Func   *UserFunction
	Args   []parser.Primary
	Result parser.Primary
	Error  string
}{
	{
		Name: "UserFunction Execute",
		Func: &UserFunction{
			Name: "userfunc",
			Parameters: []parser.Variable{
				{Name: "@arg1"},
				{Name: "@arg2"},
			},
			Statements: []parser.Statement{
				parser.VariableDeclaration{
					Assignments: []parser.Expression{
						parser.VariableAssignment{
							Name: "@var2",
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
		Name: "UserFunction Execute No Return Statement",
		Func: &UserFunction{
			Name: "userfunc",
			Parameters: []parser.Variable{
				{Name: "@arg1"},
				{Name: "@arg2"},
			},
			Statements: []parser.Statement{
				parser.VariableDeclaration{
					Assignments: []parser.Expression{
						parser.VariableAssignment{
							Name: "@var2",
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
		Name: "UserFunction Execute Arguments Error",
		Func: &UserFunction{
			Name: "userfunc",
			Parameters: []parser.Variable{
				{Name: "@arg1"},
				{Name: "@arg2"},
			},
			Statements: []parser.Statement{
				parser.VariableDeclaration{
					Assignments: []parser.Expression{
						parser.VariableAssignment{
							Name: "@var2",
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
		Error: "declared function userfunc takes 2 argument(s)",
	},
	{
		Name: "UserFunction Execute Execution Error",
		Func: &UserFunction{
			Name: "userfunc",
			Parameters: []parser.Variable{
				{Name: "@arg1"},
				{Name: "@arg2"},
			},
			Statements: []parser.Statement{
				parser.VariableDeclaration{
					Assignments: []parser.Expression{
						parser.VariableAssignment{
							Name: "@var2",
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
		Error: "field notexist does not exist",
	},
}

func TestUserFunction_Execute(t *testing.T) {
	vars := Variables{
		"@var1": parser.NewInteger(1),
	}
	filter := NewFilter([]Variables{vars})

	for _, v := range userFunctionExecuteTests {
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
