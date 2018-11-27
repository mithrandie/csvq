package query

import (
	"reflect"
	"testing"

	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/value"
)

var userDefinedFunctionScopesDeclareTests = []struct {
	Name   string
	Expr   parser.FunctionDeclaration
	Result UserDefinedFunctionScopes
	Error  string
}{
	{
		Name: "UserDefinedFunctionScopes Declare",
		Expr: parser.FunctionDeclaration{
			Name: parser.Identifier{Literal: "userfunc1"},
			Parameters: []parser.VariableAssignment{
				{Variable: parser.Variable{Name: "arg1"}},
			},
			Statements: []parser.Statement{
				parser.Print{Value: parser.Variable{Name: "arg1"}},
			},
		},
		Result: UserDefinedFunctionScopes{
			UserDefinedFunctionMap{
				"USERFUNC1": &UserDefinedFunction{
					Name: parser.Identifier{Literal: "userfunc1"},
					Parameters: []parser.Variable{
						{Name: "arg1"},
					},
					Defaults:     map[string]parser.QueryExpression{},
					RequiredArgs: 1,
					Statements: []parser.Statement{
						parser.Print{Value: parser.Variable{Name: "arg1"}},
					},
				},
			},
			UserDefinedFunctionMap{
				"USERFUNC2": &UserDefinedFunction{
					Name: parser.Identifier{Literal: "userfunc2"},
					Parameters: []parser.Variable{
						{Name: "arg1"},
						{Name: "arg2"},
					},
					Defaults:     map[string]parser.QueryExpression{},
					RequiredArgs: 2,
					Statements: []parser.Statement{
						parser.Print{Value: parser.Variable{Name: "arg2"}},
					},
				},
			},
		},
	},
}

func TestUserDefinedFunctionScopes_Declare(t *testing.T) {
	list := UserDefinedFunctionScopes{
		UserDefinedFunctionMap{},
		UserDefinedFunctionMap{
			"USERFUNC2": &UserDefinedFunction{
				Name: parser.Identifier{Literal: "userfunc2"},
				Parameters: []parser.Variable{
					{Name: "arg1"},
					{Name: "arg2"},
				},
				Defaults:     map[string]parser.QueryExpression{},
				RequiredArgs: 2,
				Statements: []parser.Statement{
					parser.Print{Value: parser.Variable{Name: "arg2"}},
				},
			},
		},
	}

	for _, v := range userDefinedFunctionScopesDeclareTests {
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
			t.Errorf("%s: result = %v, want %v", v.Name, list, v.Result)
		}
	}
}

var userDefinedFunctionScopesDeclareAggregateTests = []struct {
	Name   string
	Expr   parser.AggregateDeclaration
	Result UserDefinedFunctionScopes
	Error  string
}{
	{
		Name: "UserDefinedFunctionScopes DeclareAggregate",
		Expr: parser.AggregateDeclaration{
			Name:   parser.Identifier{Literal: "useraggfunc"},
			Cursor: parser.Identifier{Literal: "column1"},
			Parameters: []parser.VariableAssignment{
				{
					Variable: parser.Variable{Name: "arg1"},
				},
				{
					Variable: parser.Variable{Name: "arg2"},
					Value:    parser.NewIntegerValue(1),
				},
			},
			Statements: []parser.Statement{
				parser.Print{Value: parser.Variable{Name: "var1"}},
			},
		},
		Result: UserDefinedFunctionScopes{
			UserDefinedFunctionMap{
				"USERAGGFUNC": &UserDefinedFunction{
					Name: parser.Identifier{Literal: "useraggfunc"},
					Parameters: []parser.Variable{
						{Name: "arg1"},
						{Name: "arg2"},
					},
					Defaults: map[string]parser.QueryExpression{
						"arg2": parser.NewIntegerValue(1),
					},
					IsAggregate:  true,
					RequiredArgs: 1,
					Cursor:       parser.Identifier{Literal: "column1"},
					Statements: []parser.Statement{
						parser.Print{Value: parser.Variable{Name: "var1"}},
					},
				},
			},
			UserDefinedFunctionMap{
				"USERFUNC2": &UserDefinedFunction{
					Name: parser.Identifier{Literal: "userfunc2"},
					Parameters: []parser.Variable{
						{Name: "arg1"},
						{Name: "arg2"},
					},
					Defaults:     map[string]parser.QueryExpression{},
					RequiredArgs: 2,
					Statements: []parser.Statement{
						parser.Print{Value: parser.Variable{Name: "arg2"}},
					},
				},
			},
		},
	},
}

func TestUserDefinedFunctionScopes_DeclareAggregate(t *testing.T) {
	list := UserDefinedFunctionScopes{
		UserDefinedFunctionMap{},
		UserDefinedFunctionMap{
			"USERFUNC2": &UserDefinedFunction{
				Name: parser.Identifier{Literal: "userfunc2"},
				Parameters: []parser.Variable{
					{Name: "arg1"},
					{Name: "arg2"},
				},
				Defaults:     map[string]parser.QueryExpression{},
				RequiredArgs: 2,
				Statements: []parser.Statement{
					parser.Print{Value: parser.Variable{Name: "arg2"}},
				},
			},
		},
	}

	for _, v := range userDefinedFunctionScopesDeclareAggregateTests {
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
			t.Errorf("%s: result = %v, want %v", v.Name, list, v.Result)
		}
	}
}

var userDefinedFunctionScopesGetTests = []struct {
	Name     string
	Function parser.QueryExpression
	FuncName string
	Result   *UserDefinedFunction
	Error    string
}{
	{
		Name: "UserDefinedFunctionScopes Get",
		Function: parser.Function{
			Name: "userfunc2",
		},
		FuncName: "userfunc2",
		Result: &UserDefinedFunction{
			Name: parser.Identifier{Literal: "userfunc2"},
			Parameters: []parser.Variable{
				{Name: "arg1"},
				{Name: "arg2"},
			},
			Statements: []parser.Statement{
				parser.Print{Value: parser.Variable{Name: "arg2"}},
			},
		},
	},
	{
		Name: "UserDefinedFunctionScopes Get Not Exist Error",
		Function: parser.Function{
			Name: "notexist",
		},
		FuncName: "notexist",
		Error:    "[L:- C:-] function notexist does not exist",
	},
}

func TestUserDefinedFunctionScopes_Get(t *testing.T) {
	list := UserDefinedFunctionScopes{
		UserDefinedFunctionMap{
			"USERFUNC1": &UserDefinedFunction{
				Name: parser.Identifier{Literal: "userfunc1"},
				Parameters: []parser.Variable{
					{Name: "arg1"},
				},
				Statements: []parser.Statement{
					parser.Print{Value: parser.Variable{Name: "arg1"}},
				},
			},
		},
		UserDefinedFunctionMap{
			"USERFUNC2": &UserDefinedFunction{
				Name: parser.Identifier{Literal: "userfunc2"},
				Parameters: []parser.Variable{
					{Name: "arg1"},
					{Name: "arg2"},
				},
				Statements: []parser.Statement{
					parser.Print{Value: parser.Variable{Name: "arg2"}},
				},
			},
		},
	}

	for _, v := range userDefinedFunctionScopesGetTests {
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
			t.Errorf("%s: result = %v, want %v", v.Name, fn, v.Result)
		}
	}
}

var userDefinedFunctionScopesDisposeTests = []struct {
	Name     string
	FuncName parser.Identifier
	Result   UserDefinedFunctionScopes
	Error    string
}{
	{
		Name:     "UserDefinedFunctionScopes Despose",
		FuncName: parser.Identifier{Literal: "userfunc2"},
		Result: UserDefinedFunctionScopes{
			UserDefinedFunctionMap{
				"USERFUNC1": &UserDefinedFunction{
					Name: parser.Identifier{Literal: "userfunc1"},
					Parameters: []parser.Variable{
						{Name: "arg1"},
					},
					Statements: []parser.Statement{
						parser.Print{Value: parser.Variable{Name: "arg1"}},
					},
				},
			},
			UserDefinedFunctionMap{},
		},
	},
	{
		Name:     "UserDefinedFunctionScopes Despose Not Exist Error",
		FuncName: parser.Identifier{Literal: "notexist"},
		Error:    "[L:- C:-] function notexist does not exist",
	},
}

func TestUserDefinedFunctionScopes_Dispose(t *testing.T) {
	list := UserDefinedFunctionScopes{
		UserDefinedFunctionMap{
			"USERFUNC1": &UserDefinedFunction{
				Name: parser.Identifier{Literal: "userfunc1"},
				Parameters: []parser.Variable{
					{Name: "arg1"},
				},
				Statements: []parser.Statement{
					parser.Print{Value: parser.Variable{Name: "arg1"}},
				},
			},
		},
		UserDefinedFunctionMap{
			"USERFUNC2": &UserDefinedFunction{
				Name: parser.Identifier{Literal: "userfunc2"},
				Parameters: []parser.Variable{
					{Name: "arg1"},
					{Name: "arg2"},
				},
				Statements: []parser.Statement{
					parser.Print{Value: parser.Variable{Name: "arg2"}},
				},
			},
		},
	}

	for _, v := range userDefinedFunctionScopesDisposeTests {
		err := list.Dispose(v.FuncName)
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
			t.Errorf("%s: result = %v, want %v", v.Name, list, v.Result)
		}
	}
}

func TestUserDefinedFunctionScopes_All(t *testing.T) {
	list := UserDefinedFunctionScopes{
		UserDefinedFunctionMap{
			"USERFUNC1": &UserDefinedFunction{
				Name: parser.Identifier{Literal: "userfunc1"},
				Parameters: []parser.Variable{
					{Name: "arg1"},
				},
				Statements: []parser.Statement{
					parser.Print{Value: parser.Variable{Name: "arg1"}},
				},
			},
		},
		UserDefinedFunctionMap{
			"USERAGGFUNC": &UserDefinedFunction{
				Name: parser.Identifier{Literal: "useraggfunc"},
				Parameters: []parser.Variable{
					{Name: "arg1"},
					{Name: "arg2"},
				},
				Defaults: map[string]parser.QueryExpression{
					"@arg2": parser.NewIntegerValue(1),
				},
				IsAggregate:  true,
				RequiredArgs: 1,
				Cursor:       parser.Identifier{Literal: "column1"},
				Statements: []parser.Statement{
					parser.Print{Value: parser.Variable{Name: "var1"}},
				},
			},
			"USERFUNC2": &UserDefinedFunction{
				Name: parser.Identifier{Literal: "userfunc2"},
				Parameters: []parser.Variable{
					{Name: "arg1"},
					{Name: "arg2"},
				},
				Defaults: map[string]parser.QueryExpression{
					"@arg2": parser.NewIntegerValue(3),
				},
				RequiredArgs: 1,
				Statements: []parser.Statement{
					parser.Print{Value: parser.Variable{Name: "arg2"}},
				},
			},
			"USERFUNC1": &UserDefinedFunction{
				Name: parser.Identifier{Literal: "userfunc1"},
				Parameters: []parser.Variable{
					{Name: "arg2"},
				},
				Statements: []parser.Statement{
					parser.Print{Value: parser.Variable{Name: "arg1"}},
				},
			},
		},
	}

	expectScala := UserDefinedFunctionMap{
		"USERFUNC1": &UserDefinedFunction{
			Name: parser.Identifier{Literal: "userfunc1"},
			Parameters: []parser.Variable{
				{Name: "arg1"},
			},
			Statements: []parser.Statement{
				parser.Print{Value: parser.Variable{Name: "arg1"}},
			},
		},
		"USERFUNC2": &UserDefinedFunction{
			Name: parser.Identifier{Literal: "userfunc2"},
			Parameters: []parser.Variable{
				{Name: "arg1"},
				{Name: "arg2"},
			},
			Defaults: map[string]parser.QueryExpression{
				"@arg2": parser.NewIntegerValue(3),
			},
			RequiredArgs: 1,
			Statements: []parser.Statement{
				parser.Print{Value: parser.Variable{Name: "arg2"}},
			},
		},
	}

	expectAgg := UserDefinedFunctionMap{
		"USERAGGFUNC": &UserDefinedFunction{
			Name: parser.Identifier{Literal: "useraggfunc"},
			Parameters: []parser.Variable{
				{Name: "arg1"},
				{Name: "arg2"},
			},
			Defaults: map[string]parser.QueryExpression{
				"@arg2": parser.NewIntegerValue(1),
			},
			IsAggregate:  true,
			RequiredArgs: 1,
			Cursor:       parser.Identifier{Literal: "column1"},
			Statements: []parser.Statement{
				parser.Print{Value: parser.Variable{Name: "var1"}},
			},
		},
	}

	scala, agg := list.All()
	if !reflect.DeepEqual(scala, expectScala) {
		t.Errorf("scala: result = %v, want %v", scala, expectScala)
	}
	if !reflect.DeepEqual(agg, expectAgg) {
		t.Errorf("aggregate: result = %v, want %v", agg, expectAgg)
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
			Parameters: []parser.VariableAssignment{
				{
					Variable: parser.Variable{Name: "arg1"},
				},
				{
					Variable: parser.Variable{Name: "arg2"},
				},
			},
			Statements: []parser.Statement{
				parser.Print{Value: parser.Variable{Name: "var1"}},
			},
		},
		Result: UserDefinedFunctionMap{
			"USERFUNC": &UserDefinedFunction{
				Name: parser.Identifier{Literal: "userfunc"},
				Parameters: []parser.Variable{
					{Name: "arg1"},
					{Name: "arg2"},
				},
				Defaults:     map[string]parser.QueryExpression{},
				RequiredArgs: 2,
				Statements: []parser.Statement{
					parser.Print{Value: parser.Variable{Name: "var1"}},
				},
			},
		},
	},
	{
		Name: "UserDefinedFunctionMap Declare Redeclaration Error",
		Expr: parser.FunctionDeclaration{
			Name: parser.Identifier{Literal: "userfunc"},
			Parameters: []parser.VariableAssignment{
				{
					Variable: parser.Variable{Name: "arg1"},
				},
				{
					Variable: parser.Variable{Name: "arg2"},
				},
			},
			Statements: []parser.Statement{
				parser.Print{Value: parser.Variable{Name: "var1"}},
			},
		},
		Error: "[L:- C:-] function userfunc is redeclared",
	},
	{
		Name: "UserDefinedFunctionMap Declare Duplicate Prameters Error",
		Expr: parser.FunctionDeclaration{
			Name: parser.Identifier{Literal: "userfunc2"},
			Parameters: []parser.VariableAssignment{
				{
					Variable: parser.Variable{Name: "arg1"},
				},
				{
					Variable: parser.Variable{Name: "arg1"},
				},
			},
			Statements: []parser.Statement{
				parser.Print{Value: parser.Variable{Name: "var1"}},
			},
		},
		Error: "[L:- C:-] parameter @arg1 is a duplicate",
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
			t.Errorf("%s: result = %v, want %v", v.Name, funcs, v.Result)
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
			Name:   parser.Identifier{Literal: "useraggfunc"},
			Cursor: parser.Identifier{Literal: "column1"},
			Parameters: []parser.VariableAssignment{
				{
					Variable: parser.Variable{Name: "arg1"},
				},
				{
					Variable: parser.Variable{Name: "arg2"},
				},
			},
			Statements: []parser.Statement{
				parser.Print{Value: parser.Variable{Name: "var1"}},
			},
		},
		Result: UserDefinedFunctionMap{
			"USERAGGFUNC": &UserDefinedFunction{
				Name:        parser.Identifier{Literal: "useraggfunc"},
				IsAggregate: true,
				Cursor:      parser.Identifier{Literal: "column1"},
				Parameters: []parser.Variable{
					{Name: "arg1"},
					{Name: "arg2"},
				},
				RequiredArgs: 2,
				Defaults:     map[string]parser.QueryExpression{},
				Statements: []parser.Statement{
					parser.Print{Value: parser.Variable{Name: "var1"}},
				},
			},
		},
	},
	{
		Name: "UserDefinedFunctionMap DeclareAggregate Redeclaration Error",
		Expr: parser.AggregateDeclaration{
			Name:   parser.Identifier{Literal: "useraggfunc"},
			Cursor: parser.Identifier{Literal: "column1"},
			Statements: []parser.Statement{
				parser.Print{Value: parser.Variable{Name: "var1"}},
			},
		},
		Error: "[L:- C:-] function useraggfunc is redeclared",
	},
	{
		Name: "UserDefinedFunctionMap DeclareAggregate Duplicate Parameters Error",
		Expr: parser.AggregateDeclaration{
			Name:   parser.Identifier{Literal: "useraggfunc2"},
			Cursor: parser.Identifier{Literal: "column1"},
			Parameters: []parser.VariableAssignment{
				{
					Variable: parser.Variable{Name: "arg1"},
				},
				{
					Variable: parser.Variable{Name: "arg1"},
				},
			},
			Statements: []parser.Statement{
				parser.Print{Value: parser.Variable{Name: "var1"}},
			},
		},
		Error: "[L:- C:-] parameter @arg1 is a duplicate",
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
			t.Errorf("%s: result = %v, want %v", v.Name, funcs, v.Result)
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
				{Name: "arg1"},
				{Name: "arg2"},
			},
			Statements: []parser.Statement{
				parser.Print{Value: parser.Variable{Name: "var1"}},
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
	Function parser.QueryExpression
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
				{Name: "arg1"},
				{Name: "arg2"},
			},
			Statements: []parser.Statement{
				parser.Print{Value: parser.Variable{Name: "var1"}},
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
				{Name: "arg1"},
				{Name: "arg2"},
			},
			Statements: []parser.Statement{
				parser.Print{Value: parser.Variable{Name: "var1"}},
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
			t.Errorf("%s: result = %v, want %v", v.Name, result, v.Result)
		}
	}
}

var userDefinedFunctionExecuteTests = []struct {
	Name   string
	Func   *UserDefinedFunction
	Args   []value.Primary
	Result value.Primary
	Error  string
}{
	{
		Name: "UserDefinedFunction Execute",
		Func: &UserDefinedFunction{
			Name: parser.Identifier{Literal: "userfunc"},
			Parameters: []parser.Variable{
				{Name: "arg1"},
				{Name: "arg2"},
			},
			Defaults: map[string]parser.QueryExpression{
				"arg2": parser.NewIntegerValue(3),
			},
			RequiredArgs: 1,
			Statements: []parser.Statement{
				parser.VariableDeclaration{
					Assignments: []parser.VariableAssignment{
						{
							Variable: parser.Variable{Name: "var2"},
							Value: parser.Arithmetic{
								LHS: parser.Arithmetic{
									LHS:      parser.Variable{Name: "arg1"},
									RHS:      parser.Variable{Name: "arg2"},
									Operator: '+',
								},
								RHS:      parser.Variable{Name: "var1"},
								Operator: '+',
							},
						},
					},
				},
				parser.Return{
					Value: parser.Variable{Name: "var2"},
				},
			},
		},
		Args: []value.Primary{
			value.NewInteger(2),
		},
		Result: value.NewInteger(6),
	},
	{
		Name: "UserDefinedFunction Execute No Return Statement",
		Func: &UserDefinedFunction{
			Name: parser.Identifier{Literal: "userfunc"},
			Parameters: []parser.Variable{
				{Name: "arg1"},
				{Name: "arg2"},
			},
			Statements: []parser.Statement{
				parser.VariableDeclaration{
					Assignments: []parser.VariableAssignment{
						{
							Variable: parser.Variable{Name: "var2"},
							Value: parser.Arithmetic{
								LHS: parser.Arithmetic{
									LHS:      parser.Variable{Name: "arg1"},
									RHS:      parser.Variable{Name: "arg2"},
									Operator: '+',
								},
								RHS:      parser.Variable{Name: "var1"},
								Operator: '+',
							},
						},
					},
				},
			},
		},
		Args: []value.Primary{
			value.NewInteger(2),
			value.NewInteger(3),
		},
		Result: value.NewNull(),
	},
	{
		Name: "UserDefinedFunction Execute Argument Length Error",
		Func: &UserDefinedFunction{
			Name: parser.Identifier{Literal: "userfunc"},
			Parameters: []parser.Variable{
				{Name: "arg1"},
				{Name: "arg2"},
			},
			Statements: []parser.Statement{
				parser.VariableDeclaration{
					Assignments: []parser.VariableAssignment{
						{
							Variable: parser.Variable{Name: "var2"},
							Value: parser.Arithmetic{
								LHS: parser.Arithmetic{
									LHS:      parser.Variable{Name: "arg1"},
									RHS:      parser.Variable{Name: "arg2"},
									Operator: '+',
								},
								RHS:      parser.Variable{Name: "var1"},
								Operator: '+',
							},
						},
					},
				},
				parser.Return{
					Value: parser.Variable{Name: "var2"},
				},
			},
		},
		Args: []value.Primary{
			value.NewInteger(2),
		},
		Error: "[L:- C:-] function userfunc takes exactly 2 arguments",
	},
	{
		Name: "UserDefinedFunction Execute Argument Evaluation Error",
		Func: &UserDefinedFunction{
			Name: parser.Identifier{Literal: "userfunc"},
			Parameters: []parser.Variable{
				{Name: "arg1"},
				{Name: "arg2"},
			},
			Defaults: map[string]parser.QueryExpression{
				"arg2": parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			},
			RequiredArgs: 1,
			Statements: []parser.Statement{
				parser.VariableDeclaration{
					Assignments: []parser.VariableAssignment{
						{
							Variable: parser.Variable{Name: "var2"},
							Value: parser.Arithmetic{
								LHS: parser.Arithmetic{
									LHS:      parser.Variable{Name: "arg1"},
									RHS:      parser.Variable{Name: "arg2"},
									Operator: '+',
								},
								RHS:      parser.Variable{Name: "var1"},
								Operator: '+',
							},
						},
					},
				},
				parser.Return{
					Value: parser.Variable{Name: "var2"},
				},
			},
		},
		Args: []value.Primary{
			value.NewInteger(2),
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
	{
		Name: "UserDefinedFunction Execute Execution Error",
		Func: &UserDefinedFunction{
			Name: parser.Identifier{Literal: "userfunc"},
			Parameters: []parser.Variable{
				{Name: "arg1"},
				{Name: "arg2"},
			},
			Statements: []parser.Statement{
				parser.VariableDeclaration{
					Assignments: []parser.VariableAssignment{
						{
							Variable: parser.Variable{Name: "var2"},
							Value: parser.Subquery{
								Query: parser.SelectQuery{
									SelectEntity: parser.SelectEntity{
										SelectClause: parser.SelectClause{
											Fields: []parser.QueryExpression{
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
		Args: []value.Primary{
			value.NewInteger(2),
			value.NewInteger(3),
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
}

func TestUserDefinedFunction_Execute(t *testing.T) {
	vars := GenerateVariableMap(map[string]value.Primary{
		"var1": value.NewInteger(1),
	})
	filter := NewFilter(
		[]VariableMap{vars},
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

var userDefinedFunctionExecuteAggregateTests = []struct {
	Name   string
	Func   *UserDefinedFunction
	Values []value.Primary
	Args   []value.Primary
	Result value.Primary
	Error  string
}{
	{
		Name: "UserDefinedFunction Execute Aggregate",
		Func: &UserDefinedFunction{
			Name:        parser.Identifier{Literal: "useraggfunc"},
			IsAggregate: true,
			Cursor:      parser.Identifier{Literal: "list"},
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
		Values: []value.Primary{
			value.NewInteger(1),
			value.NewInteger(2),
			value.NewInteger(3),
		},
		Result: value.NewInteger(6),
	},
	{
		Name: "UserDefinedFunction Execute Aggregate With Arguments",
		Func: &UserDefinedFunction{
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
		Values: []value.Primary{
			value.NewNull(),
			value.NewNull(),
			value.NewNull(),
		},
		Args: []value.Primary{
			value.NewInteger(0),
		},
		Result: value.NewInteger(0),
	},
	{
		Name: "UserDefinedFunction Aggregate Argument Length Error",
		Func: &UserDefinedFunction{
			Name:        parser.Identifier{Literal: "useraggfunc"},
			IsAggregate: true,
			Cursor:      parser.Identifier{Literal: "list"},
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
		Values: []value.Primary{
			value.NewInteger(1),
			value.NewInteger(2),
			value.NewInteger(3),
		},
		Args: []value.Primary{
			value.NewInteger(0),
		},
		Error: "[L:- C:-] function useraggfunc takes exactly 1 argument",
	},
}

func TestUserDefinedFunction_ExecuteAggregate(t *testing.T) {
	filter := NewEmptyFilter()

	for _, v := range userDefinedFunctionExecuteAggregateTests {
		result, err := v.Func.ExecuteAggregate(v.Values, v.Args, filter)
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

var userDefinedFunctionCheckArgsLenTests = []struct {
	Name    string
	Func    *UserDefinedFunction
	ArgsLen int
	Error   string
}{
	{
		Name: "UserDefinedFunction CheckArgsLen",
		Func: &UserDefinedFunction{
			Name: parser.Identifier{Literal: "userfunc"},
			Parameters: []parser.Variable{
				{Name: "arg1"},
				{Name: "arg2"},
			},
			Defaults: map[string]parser.QueryExpression{
				"arg2": parser.NewIntegerValue(3),
			},
			RequiredArgs: 1,
			Statements:   []parser.Statement{},
		},
		ArgsLen: 1,
		Error:   "",
	},
	{
		Name: "UserDefinedFunction CheckArgsLen Argument Length Error",
		Func: &UserDefinedFunction{
			Name: parser.Identifier{Literal: "userfunc"},
			Parameters: []parser.Variable{
				{Name: "arg1"},
				{Name: "arg2"},
			},
			RequiredArgs: 2,
			Statements:   []parser.Statement{},
		},
		ArgsLen: 1,
		Error:   "[L:- C:-] function userfunc takes exactly 2 arguments",
	},
	{
		Name: "UserDefinedFunction CheckArgsLen Too Little Argument Error",
		Func: &UserDefinedFunction{
			Name: parser.Identifier{Literal: "userfunc"},
			Parameters: []parser.Variable{
				{Name: "arg1"},
				{Name: "arg2"},
			},
			Defaults: map[string]parser.QueryExpression{
				"arg2": parser.NewIntegerValue(3),
			},
			RequiredArgs: 1,
			Statements:   []parser.Statement{},
		},
		ArgsLen: 0,
		Error:   "[L:- C:-] function userfunc takes at least 1 argument",
	},
	{
		Name: "UserDefinedFunction CheckArgsLen Too Many Argument Length Error",
		Func: &UserDefinedFunction{
			Name: parser.Identifier{Literal: "userfunc"},
			Parameters: []parser.Variable{
				{Name: "arg1"},
				{Name: "arg2"},
			},
			Defaults: map[string]parser.QueryExpression{
				"arg2": parser.NewIntegerValue(3),
			},
			RequiredArgs: 1,
			Statements:   []parser.Statement{},
		},
		ArgsLen: 3,
		Error:   "[L:- C:-] function userfunc takes at most 2 arguments",
	},
	{
		Name: "UserDefinedFunction CheckArgsLen Aggregate Argument Length Error",
		Func: &UserDefinedFunction{
			Name:        parser.Identifier{Literal: "userfunc"},
			IsAggregate: true,
			Parameters: []parser.Variable{
				{Name: "arg1"},
				{Name: "arg2"},
			},
			RequiredArgs: 2,
			Cursor:       parser.Identifier{Literal: "list"},
			Statements:   []parser.Statement{},
		},
		ArgsLen: 1,
		Error:   "[L:- C:-] function userfunc takes exactly 3 arguments",
	},
}

func TestUserDefinedFunction_CheckArgsLen(t *testing.T) {
	for _, v := range userDefinedFunctionCheckArgsLenTests {
		err := v.Func.CheckArgsLen(parser.Identifier{Literal: "userfunc"}, "userfunc", v.ArgsLen)
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
