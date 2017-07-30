package query

import (
	"reflect"
	"testing"

	"github.com/mithrandie/csvq/lib/parser"
)

var userDefinedFunctionListDeclareTests = []struct {
	Name   string
	Expr   parser.FunctionDeclaration
	Result UserDefinedFunctionsList
	Error  string
}{
	{
		Name: "UserDefineFunctionsList Declare",
		Expr: parser.FunctionDeclaration{
			Name: parser.Identifier{Literal: "userfunc1"},
			Parameters: []parser.Variable{
				{Name: "@arg1"},
			},
			Statements: []parser.Statement{
				parser.Print{Value: parser.Variable{Name: "@arg1"}},
			},
		},
		Result: UserDefinedFunctionsList{
			UserDefinedFunctionMap{
				"USERFUNC1": &UserDefinedFunction{
					Name: parser.Identifier{Literal: "userfunc1"},
					Parameters: []parser.Variable{
						{Name: "@arg1"},
					},
					Statements: []parser.Statement{
						parser.Print{Value: parser.Variable{Name: "@arg1"}},
					},
				},
			},
			UserDefinedFunctionMap{
				"USERFUNC2": &UserDefinedFunction{
					Name: parser.Identifier{Literal: "userfunc2"},
					Parameters: []parser.Variable{
						{Name: "@arg1"},
						{Name: "@arg2"},
					},
					Statements: []parser.Statement{
						parser.Print{Value: parser.Variable{Name: "@arg2"}},
					},
				},
			},
		},
	},
}

func TestUserDefinedFunctionsList_Declare(t *testing.T) {
	list := UserDefinedFunctionsList{
		UserDefinedFunctionMap{},
		UserDefinedFunctionMap{
			"USERFUNC2": &UserDefinedFunction{
				Name: parser.Identifier{Literal: "userfunc2"},
				Parameters: []parser.Variable{
					{Name: "@arg1"},
					{Name: "@arg2"},
				},
				Statements: []parser.Statement{
					parser.Print{Value: parser.Variable{Name: "@arg2"}},
				},
			},
		},
	}

	for _, v := range userDefinedFunctionListDeclareTests {
		err := list.Declare(v.Expr)
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
		if !reflect.DeepEqual(list, v.Result) {
			t.Errorf("%s: result = %s, want %s", v.Name, list, v.Result)
		}
	}
}

var userDefinedFunctionListDeclareAggregateTests = []struct {
	Name   string
	Expr   parser.AggregateDeclaration
	Result UserDefinedFunctionsList
	Error  string
}{
	{
		Name: "UserDefineFunctionsList Declare",
		Expr: parser.AggregateDeclaration{
			Name:      parser.Identifier{Literal: "useraggfunc"},
			Parameter: parser.Identifier{Literal: "column1"},
			Statements: []parser.Statement{
				parser.Print{Value: parser.Variable{Name: "@var1"}},
			},
		},
		Result: UserDefinedFunctionsList{
			UserDefinedFunctionMap{
				"USERAGGFUNC": &UserDefinedFunction{
					Name:        parser.Identifier{Literal: "useraggfunc"},
					IsAggregate: true,
					Parameter:   parser.Identifier{Literal: "column1"},
					Statements: []parser.Statement{
						parser.Print{Value: parser.Variable{Name: "@var1"}},
					},
				},
			},
			UserDefinedFunctionMap{
				"USERFUNC2": &UserDefinedFunction{
					Name: parser.Identifier{Literal: "userfunc2"},
					Parameters: []parser.Variable{
						{Name: "@arg1"},
						{Name: "@arg2"},
					},
					Statements: []parser.Statement{
						parser.Print{Value: parser.Variable{Name: "@arg2"}},
					},
				},
			},
		},
	},
}

func TestUserDefinedFunctionsList_DeclareAggregate(t *testing.T) {
	list := UserDefinedFunctionsList{
		UserDefinedFunctionMap{},
		UserDefinedFunctionMap{
			"USERFUNC2": &UserDefinedFunction{
				Name: parser.Identifier{Literal: "userfunc2"},
				Parameters: []parser.Variable{
					{Name: "@arg1"},
					{Name: "@arg2"},
				},
				Statements: []parser.Statement{
					parser.Print{Value: parser.Variable{Name: "@arg2"}},
				},
			},
		},
	}

	for _, v := range userDefinedFunctionListDeclareAggregateTests {
		err := list.DeclareAggregate(v.Expr)
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
		if !reflect.DeepEqual(list, v.Result) {
			t.Errorf("%s: result = %s, want %s", v.Name, list, v.Result)
		}
	}
}

var userDefinedFunctionListGetTests = []struct {
	Name     string
	Function parser.Expression
	FuncName string
	Result   *UserDefinedFunction
	Error    string
}{
	{
		Name: "UserDefineFunctionsList Get",
		Function: parser.Function{
			Name: "userfunc2",
		},
		FuncName: "userfunc2",
		Result: &UserDefinedFunction{
			Name: parser.Identifier{Literal: "userfunc2"},
			Parameters: []parser.Variable{
				{Name: "@arg1"},
				{Name: "@arg2"},
			},
			Statements: []parser.Statement{
				parser.Print{Value: parser.Variable{Name: "@arg2"}},
			},
		},
	},
	{
		Name: "UserDefineFunctionsList Get Not Exist Error",
		Function: parser.Function{
			Name: "notexist",
		},
		FuncName: "notexist",
		Error:    "[L:- C:-] function notexist does not exist",
	},
}

func TestUserDefinedFunctionsList_Get(t *testing.T) {
	list := UserDefinedFunctionsList{
		UserDefinedFunctionMap{
			"USERFUNC1": &UserDefinedFunction{
				Name: parser.Identifier{Literal: "userfunc1"},
				Parameters: []parser.Variable{
					{Name: "@arg1"},
				},
				Statements: []parser.Statement{
					parser.Print{Value: parser.Variable{Name: "@arg1"}},
				},
			},
		},
		UserDefinedFunctionMap{
			"USERFUNC2": &UserDefinedFunction{
				Name: parser.Identifier{Literal: "userfunc2"},
				Parameters: []parser.Variable{
					{Name: "@arg1"},
					{Name: "@arg2"},
				},
				Statements: []parser.Statement{
					parser.Print{Value: parser.Variable{Name: "@arg2"}},
				},
			},
		},
	}

	for _, v := range userDefinedFunctionListGetTests {
		fn, err := list.Get(v.Function, v.FuncName)
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
		if !reflect.DeepEqual(fn, v.Result) {
			t.Errorf("%s: result = %s, want %s", v.Name, fn, v.Result)
		}
	}
}

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

var userDefinedFunctionMapDeclareAggregateTests = []struct {
	Name   string
	Expr   parser.AggregateDeclaration
	Result UserDefinedFunctionMap
	Error  string
}{
	{
		Name: "UserDefinedFunctionMap DeclareAggregate",
		Expr: parser.AggregateDeclaration{
			Name:      parser.Identifier{Literal: "useraggfunc"},
			Parameter: parser.Identifier{Literal: "column1"},
			Statements: []parser.Statement{
				parser.Print{Value: parser.Variable{Name: "@var1"}},
			},
		},
		Result: UserDefinedFunctionMap{
			"USERAGGFUNC": &UserDefinedFunction{
				Name:        parser.Identifier{Literal: "useraggfunc"},
				IsAggregate: true,
				Parameter:   parser.Identifier{Literal: "column1"},
				Statements: []parser.Statement{
					parser.Print{Value: parser.Variable{Name: "@var1"}},
				},
			},
		},
	},
	{
		Name: "UserDefinedFunctionMap DeclareAggregate Redeclaration Error",
		Expr: parser.AggregateDeclaration{
			Name:      parser.Identifier{Literal: "useraggfunc"},
			Parameter: parser.Identifier{Literal: "column1"},
			Statements: []parser.Statement{
				parser.Print{Value: parser.Variable{Name: "@var1"}},
			},
		},
		Error: "[L:- C:-] function useraggfunc is redeclared",
	},
}

func TestUserDefinedFunctionMap_DeclareAggregate(t *testing.T) {
	funcs := UserDefinedFunctionMap{}

	for _, v := range userDefinedFunctionMapDeclareAggregateTests {
		err := funcs.DeclareAggregate(v.Expr)
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

var userDefinedFunctionMapCheckDuplicateTests = []struct {
	Name     string
	FuncName parser.Identifier
	Result   bool
	Error    string
}{
	{
		Name:     "UserDefinedFunctionMap CheckDuplicate Redeclaration Error",
		FuncName: parser.Identifier{Literal: "userfunc"},
		Error:    "[L:- C:-] function userfunc is redeclared",
	},
	{
		Name:     "UserDefinedFunctionMap CheckDuplicate Duplicate with Built-in Function Error",
		FuncName: parser.Identifier{Literal: "now"},
		Error:    "[L:- C:-] function now is a built-in function",
	},
	{
		Name:     "UserDefinedFunctionMap CheckDuplicate Duplicate with Aggregate Function Error",
		FuncName: parser.Identifier{Literal: "count"},
		Error:    "[L:- C:-] function count is a built-in function",
	},
	{
		Name:     "UserDefinedFunctionMap CheckDuplicate Duplicate with Analytic Function Error",
		FuncName: parser.Identifier{Literal: "row_number"},
		Error:    "[L:- C:-] function row_number is a built-in function",
	},
	{
		Name:     "UserDefinedFunctionMap CheckDuplicate OK",
		FuncName: parser.Identifier{Literal: "undefined"},
		Result:   true,
	},
}

func TestUserDefinedFunctionMap_CheckDuplicate(t *testing.T) {
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

	for _, v := range userDefinedFunctionMapCheckDuplicateTests {
		err := funcs.CheckDuplicate(v.FuncName)
		if err != nil {
			if v.Result {
				continue
			}

			if len(v.Error) < 1 {
				t.Errorf("%s: unexpected error %q", v.Name, err)
			} else if err.Error() != v.Error {
				t.Errorf("%s: error %q, want error %q", v.Name, err.Error(), v.Error)
			}
			continue
		}
		if 0 < len(v.Error) {
			t.Errorf("%s: no error, want error %q", v.Name, v.Error)
		}
	}
}

var userDefinedFunctionMapGetTests = []struct {
	Name     string
	Function parser.Expression
	FuncName string
	Result   *UserDefinedFunction
	Error    string
}{
	{
		Name: "UserDefinedFunctionMap Get",
		Function: parser.Function{
			Name: "userfunc",
		},
		FuncName: "userfunc",
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
		FuncName: "notexist",
		Error:    "[L:- C:-] function notexist does not exist",
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
		result, err := funcs.Get(v.Function, v.FuncName)
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
	{
		Name: "UserDefinedFunction Execute Aggregate",
		Func: &UserDefinedFunction{
			Name:        parser.Identifier{Literal: "useraggfunc"},
			IsAggregate: true,
			Parameter:   parser.Identifier{Literal: "column1"},
			Statements: []parser.Statement{
				parser.VariableDeclaration{
					Assignments: []parser.Expression{
						parser.VariableAssignment{
							Variable: parser.Variable{Name: "@value"},
							Value:    parser.NewInteger(0),
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
						parser.VariableSubstitution{
							Variable: parser.Variable{Name: "@value"},
							Value: parser.Arithmetic{
								LHS:      parser.Variable{Name: "@value"},
								RHS:      parser.Variable{Name: "@fetch"},
								Operator: '+',
							},
						},
					},
				},
				parser.Return{
					Value: parser.Variable{Name: "@value"},
				},
			},
		},
		Args: []parser.Primary{
			parser.NewInteger(1),
			parser.NewInteger(2),
			parser.NewInteger(3),
		},
		Result: parser.NewInteger(6),
	},
	{
		Name: "UserDefinedFunction Execute Aggregate",
		Func: &UserDefinedFunction{
			Name:        parser.Identifier{Literal: "useraggfunc"},
			IsAggregate: true,
			Parameter:   parser.Identifier{Literal: "column1"},
			Statements: []parser.Statement{
				parser.VariableDeclaration{
					Assignments: []parser.Expression{
						parser.VariableAssignment{
							Variable: parser.Variable{Name: "@value"},
							Value:    parser.NewInteger(0),
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
						parser.VariableSubstitution{
							Variable: parser.Variable{Name: "@value"},
							Value: parser.Arithmetic{
								LHS:      parser.Variable{Name: "@value"},
								RHS:      parser.Variable{Name: "@fetch"},
								Operator: '+',
							},
						},
					},
				},
				parser.Return{
					Value: parser.Variable{Name: "@value"},
				},
			},
		},
		Args: []parser.Primary{
			parser.NewInteger(1),
			parser.NewInteger(2),
			parser.NewInteger(3),
		},
		Result: parser.NewInteger(6),
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
		[]UserDefinedFunctionMap{{}},
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
