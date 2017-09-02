package query

import (
	"reflect"
	"testing"

	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/value"
)

var variableScopesGet = []struct {
	Name   string
	Expr   parser.Variable
	Result value.Primary
	Error  string
}{
	{
		Name:   "VariableScopes Get",
		Expr:   parser.Variable{Name: "@var1"},
		Result: value.NewInteger(1),
	},
	{
		Name:  "VariableScopes Get Undefined Error",
		Expr:  parser.Variable{Name: "@undef"},
		Error: "[L:- C:-] variable @undef is undefined",
	},
}

func TestVariableScopes_Get(t *testing.T) {
	list := VariableScopes{
		VariableMap{
			"@var1": value.NewInteger(1),
		},
		VariableMap{
			"@var1": value.NewInteger(2),
		},
	}

	for _, v := range variableScopesGet {
		result, err := list.Get(v.Expr)
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

var variableScopesSubstituteTests = []struct {
	Name   string
	Expr   parser.VariableSubstitution
	Filter *Filter
	List   VariableScopes
	Result value.Primary
	Error  string
}{
	{
		Name: "VariableScopes Substitute",
		Expr: parser.VariableSubstitution{
			Variable: parser.Variable{Name: "@var1"},
			Value:    parser.NewIntegerValue(3),
		},
		List: VariableScopes{
			VariableMap{
				"@var1": value.NewInteger(3),
			},
			VariableMap{
				"@var1": value.NewInteger(2),
			},
		},
		Result: value.NewInteger(3),
	},
	{
		Name: "VariableScopes Substitute Variable Undefined Error",
		Expr: parser.VariableSubstitution{
			Variable: parser.Variable{Name: "var2"},
			Value:    parser.NewIntegerValue(3),
		},
		Error: "[L:- C:-] variable var2 is undefined",
	},
}

func TestVariableScopes_Substitute(t *testing.T) {
	list := VariableScopes{
		VariableMap{
			"@var1": value.NewInteger(1),
		},
		VariableMap{
			"@var1": value.NewInteger(2),
		},
	}

	for _, v := range variableScopesSubstituteTests {
		result, err := list.Substitute(v.Expr, v.Filter)
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
		if !reflect.DeepEqual(list, v.List) {
			t.Errorf("%s: list = %s, want %s", v.Name, list, v.List)
		}
		if !reflect.DeepEqual(result, v.Result) {
			t.Errorf("%s: result = %s, want %s", v.Name, result, v.Result)
		}
	}
}

var variableScopesDisposeTests = []struct {
	Name  string
	Expr  parser.Variable
	List  VariableScopes
	Error string
}{
	{
		Name: "VariableScopes Dispose",
		Expr: parser.Variable{Name: "@var1"},
		List: VariableScopes{
			VariableMap{},
			VariableMap{
				"@var1": value.NewInteger(2),
			},
		},
	},
	{
		Name:  "VariableScopes Dispose Undefined Error",
		Expr:  parser.Variable{Name: "@undef"},
		Error: "[L:- C:-] variable @undef is undefined",
	},
}

func TestVariableScopes_Dispose(t *testing.T) {
	list := VariableScopes{
		VariableMap{
			"@var1": value.NewInteger(1),
		},
		VariableMap{
			"@var1": value.NewInteger(2),
		},
	}

	for _, v := range variableScopesDisposeTests {
		err := list.Dispose(v.Expr)
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
		if !reflect.DeepEqual(list, v.List) {
			t.Errorf("%s: list = %s, want %s", v.Name, list, v.List)
		}
	}
}

type variableMapTests struct {
	Name   string
	Expr   parser.Expression
	Filter *Filter
	Result VariableMap
	Error  string
}

var variableMapDeclareTests = []variableMapTests{
	{
		Name: "Declare Variable",
		Expr: parser.VariableDeclaration{
			Assignments: []parser.VariableAssignment{
				{
					Variable: parser.Variable{Name: "@var1"},
				},
			},
		},
		Result: VariableMap{
			"@var1": value.NewNull(),
		},
	},
	{
		Name: "Declare Variable With Initial Value",
		Expr: parser.VariableDeclaration{
			Assignments: []parser.VariableAssignment{
				{
					Variable: parser.Variable{Name: "@var2"},
					Value:    parser.NewIntegerValue(1),
				},
			},
		},
		Result: VariableMap{
			"@var1": value.NewNull(),
			"@var2": value.NewInteger(1),
		},
	},
	{
		Name: "Declare Variable Redeclaration Error",
		Expr: parser.VariableDeclaration{
			Assignments: []parser.VariableAssignment{
				{
					Variable: parser.Variable{Name: "@var2"},
					Value:    parser.NewIntegerValue(1),
				},
			},
		},
		Error: "[L:- C:-] variable @var2 is redeclared",
	},
	{
		Name: "Declare Variable Filter Error",
		Expr: parser.VariableDeclaration{
			Assignments: []parser.VariableAssignment{
				{
					Variable: parser.Variable{Name: "@var3"},
					Value:    parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
				},
			},
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
}

func TestVariableMap_Declare(t *testing.T) {
	vars := VariableMap{}

	for _, v := range variableMapDeclareTests {
		if v.Filter == nil {
			v.Filter = NewEmptyFilter()
		}

		err := vars.Declare(v.Expr.(parser.VariableDeclaration), v.Filter)
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
		if !reflect.DeepEqual(vars, v.Result) {
			t.Errorf("%s: result = %s, want %s", v.Name, vars, v.Result)
		}
	}
}

var variableMapSubstituteTests = []variableMapTests{
	{
		Name: "Substitute Variable",
		Expr: parser.VariableSubstitution{
			Variable: parser.Variable{Name: "var1"},
			Value:    parser.NewIntegerValue(2),
		},
		Result: VariableMap{
			"var1": value.NewInteger(2),
		},
	},
	{
		Name: "Substitute Variable Undefined Error",
		Expr: parser.VariableSubstitution{
			Variable: parser.Variable{Name: "var2"},
			Value:    parser.NewIntegerValue(2),
		},
		Error: "[L:- C:-] variable var2 is undefined",
	},
	{
		Name: "Substitute Variable Filter Error",
		Expr: parser.VariableSubstitution{
			Variable: parser.Variable{Name: "var1"},
			Value:    parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
}

func TestVariableMap_Substitute(t *testing.T) {
	vars := VariableMap{
		"var1": value.NewInteger(1),
	}

	for _, v := range variableMapSubstituteTests {
		if v.Filter == nil {
			v.Filter = NewEmptyFilter()
		}

		_, err := vars.Substitute(v.Expr.(parser.VariableSubstitution), v.Filter)
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
		if !reflect.DeepEqual(vars, v.Result) {
			t.Errorf("%s: result = %s, want %s", v.Name, vars, v.Result)
		}
	}
}

var variableMapDisposeTests = []variableMapTests{
	{
		Name:   "Dispose Variable",
		Expr:   parser.Variable{Name: "var1"},
		Result: VariableMap{},
	},
	{
		Name:  "Dispose Variable Undefined Error",
		Expr:  parser.Variable{Name: "var2"},
		Error: "[L:- C:-] variable var2 is undefined",
	},
}

func TestVariableMap_Dispose(t *testing.T) {
	vars := VariableMap{
		"var1": value.NewInteger(1),
	}

	for _, v := range variableMapDisposeTests {
		err := vars.Dispose(v.Expr.(parser.Variable))
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
		if !reflect.DeepEqual(vars, v.Result) {
			t.Errorf("%s: result = %s, want %s", v.Name, vars, v.Result)
		}
	}
}
