package query

import (
	"reflect"
	"testing"

	"github.com/mithrandie/csvq/lib/parser"
)

type variableTests struct {
	Name   string
	Expr   parser.ProcExpr
	Filter Filter
	Result Variables
	Error  string
}

var variablesDeclareTests = []variableTests{
	{
		Name: "Declare Variable",
		Expr: parser.VariableDeclaration{
			Assignments: []parser.Expression{
				parser.VariableAssignment{
					Variable: parser.Variable{Name: "@var1"},
				},
			},
		},
		Result: Variables{
			"@var1": parser.NewNull(),
		},
	},
	{
		Name: "Declare Variable With Initial Value",
		Expr: parser.VariableDeclaration{
			Assignments: []parser.Expression{
				parser.VariableAssignment{
					Variable: parser.Variable{Name: "@var2"},
					Value:    parser.NewInteger(1),
				},
			},
		},
		Result: Variables{
			"@var1": parser.NewNull(),
			"@var2": parser.NewInteger(1),
		},
	},
	{
		Name: "Declare Variable Redeclaration Error",
		Expr: parser.VariableDeclaration{
			Assignments: []parser.Expression{
				parser.VariableAssignment{
					Variable: parser.Variable{Name: "@var2"},
					Value:    parser.NewInteger(1),
				},
			},
		},
		Error: "[L:- C:-] variable @var2 is redeclared",
	},
	{
		Name: "Declare Variable Filter Error",
		Expr: parser.VariableDeclaration{
			Assignments: []parser.Expression{
				parser.VariableAssignment{
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
			Value:    parser.NewInteger(2),
		},
		Result: Variables{
			"var1": parser.NewInteger(2),
		},
	},
	{
		Name: "Substitute Variable Undefined Error",
		Expr: parser.VariableSubstitution{
			Variable: parser.Variable{Name: "var2"},
			Value:    parser.NewInteger(2),
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
		"var1": parser.NewInteger(1),
	}

	for _, v := range variablesSubstituteTests {
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
		"var1": parser.NewInteger(1),
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
