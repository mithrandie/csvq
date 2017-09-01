package query

import (
	"reflect"
	"testing"

	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/value"
)

var variablesListGet = []struct {
	Name   string
	Expr   parser.Variable
	Result value.Primary
	Error  string
}{
	{
		Name:   "VariablesList Get",
		Expr:   parser.Variable{Name: "@var1"},
		Result: value.NewInteger(1),
	},
	{
		Name:  "VariablesList Get Undefined Error",
		Expr:  parser.Variable{Name: "@undef"},
		Error: "[L:- C:-] variable @undef is undefined",
	},
}

func TestVariablesList_Get(t *testing.T) {
	list := VariablesList{
		Variables{
			"@var1": value.NewInteger(1),
		},
		Variables{
			"@var1": value.NewInteger(2),
		},
	}

	for _, v := range variablesListGet {
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

var variablesListSubstituteTests = []struct {
	Name   string
	Expr   parser.VariableSubstitution
	Filter *Filter
	List   VariablesList
	Result value.Primary
	Error  string
}{
	{
		Name: "VariablesList Substitute",
		Expr: parser.VariableSubstitution{
			Variable: parser.Variable{Name: "@var1"},
			Value:    parser.NewIntegerValue(3),
		},
		List: VariablesList{
			Variables{
				"@var1": value.NewInteger(3),
			},
			Variables{
				"@var1": value.NewInteger(2),
			},
		},
		Result: value.NewInteger(3),
	},
	{
		Name: "VariablesList Substitute Variable Undefined Error",
		Expr: parser.VariableSubstitution{
			Variable: parser.Variable{Name: "var2"},
			Value:    parser.NewIntegerValue(3),
		},
		Error: "[L:- C:-] variable var2 is undefined",
	},
}

func TestVariablesList_Substitute(t *testing.T) {
	list := VariablesList{
		Variables{
			"@var1": value.NewInteger(1),
		},
		Variables{
			"@var1": value.NewInteger(2),
		},
	}

	for _, v := range variablesListSubstituteTests {
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

var variablesListDisposeTests = []struct {
	Name  string
	Expr  parser.Variable
	List  VariablesList
	Error string
}{
	{
		Name: "VariablesList Dispose",
		Expr: parser.Variable{Name: "@var1"},
		List: VariablesList{
			Variables{},
			Variables{
				"@var1": value.NewInteger(2),
			},
		},
	},
	{
		Name:  "VariablesList Dispose Undefined Error",
		Expr:  parser.Variable{Name: "@undef"},
		Error: "[L:- C:-] variable @undef is undefined",
	},
}

func TestVariablesList_Dispose(t *testing.T) {
	list := VariablesList{
		Variables{
			"@var1": value.NewInteger(1),
		},
		Variables{
			"@var1": value.NewInteger(2),
		},
	}

	for _, v := range variablesListDisposeTests {
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

type variableTests struct {
	Name   string
	Expr   parser.Expression
	Filter *Filter
	Result Variables
	Error  string
}

var variablesDeclareTests = []variableTests{
	{
		Name: "Declare Variable",
		Expr: parser.VariableDeclaration{
			Assignments: []parser.VariableAssignment{
				{
					Variable: parser.Variable{Name: "@var1"},
				},
			},
		},
		Result: Variables{
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
		Result: Variables{
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

func TestVariables_Declare(t *testing.T) {
	vars := Variables{}

	for _, v := range variablesDeclareTests {
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

var variablesSubstituteTests = []variableTests{
	{
		Name: "Substitute Variable",
		Expr: parser.VariableSubstitution{
			Variable: parser.Variable{Name: "var1"},
			Value:    parser.NewIntegerValue(2),
		},
		Result: Variables{
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

func TestVariables_Substitute(t *testing.T) {
	vars := Variables{
		"var1": value.NewInteger(1),
	}

	for _, v := range variablesSubstituteTests {
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

var variablesDisposeTests = []variableTests{
	{
		Name:   "Dispose Variable",
		Expr:   parser.Variable{Name: "var1"},
		Result: Variables{},
	},
	{
		Name:  "Dispose Variable Undefined Error",
		Expr:  parser.Variable{Name: "var2"},
		Error: "[L:- C:-] variable var2 is undefined",
	},
}

func TestVariables_Dispose(t *testing.T) {
	vars := Variables{
		"var1": value.NewInteger(1),
	}

	for _, v := range variablesDisposeTests {
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
