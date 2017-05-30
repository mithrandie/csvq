package query

import (
	"reflect"
	"testing"

	"github.com/mithrandie/csvq/lib/parser"
)

type variableTests struct {
	Name   string
	Expr   parser.Expression
	Filter Filter
	Result Variables
	Error  string
}

var variablesDeclareTests = []variableTests{
	{
		Name: "Decrare Variable",
		Expr: parser.VariableDeclaration{
			Assignments: []parser.Expression{
				parser.VariableAssignment{
					Name: "var1",
				},
			},
		},
		Result: Variables{
			"var1": parser.NewNull(),
		},
	},
	{
		Name: "Decrare Variable With Initial Value",
		Expr: parser.VariableDeclaration{
			Assignments: []parser.Expression{
				parser.VariableAssignment{
					Name:  "var2",
					Value: parser.NewInteger(1),
				},
			},
		},
		Result: Variables{
			"var1": parser.NewNull(),
			"var2": parser.NewInteger(1),
		},
	},
	{
		Name: "Decrare Variable Redeclaration Error",
		Expr: parser.VariableDeclaration{
			Assignments: []parser.Expression{
				parser.VariableAssignment{
					Name:  "var2",
					Value: parser.NewInteger(1),
				},
			},
		},
		Error: "variable var2 is redeclared",
	},
	{
		Name: "Decrare Variable Filter Error",
		Expr: parser.VariableDeclaration{
			Assignments: []parser.Expression{
				parser.VariableAssignment{
					Name:  "var3",
					Value: parser.Identifier{Literal: "notexist"},
				},
			},
		},
		Error: "identifier = notexist: field does not exist",
	},
}

func TestVariables_Decrare(t *testing.T) {
	vars := Variables{}

	for _, v := range variablesDeclareTests {
		err := vars.Decrare(v.Expr.(parser.VariableDeclaration), v.Filter)
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
		Error: "variable var2 is undefined",
	},
	{
		Name: "Substitute Variable Filter Error",
		Expr: parser.VariableSubstitution{
			Variable: parser.Variable{Name: "var1"},
			Value:    parser.Identifier{Literal: "notexist"},
		},
		Error: "identifier = notexist: field does not exist",
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

func TestVariables_ClearAutoIncrement(t *testing.T) {
	vars := Variables{
		AUTO_INCREMENT_KEY: parser.NewInteger(1),
	}
	vars.ClearAutoIncrement()
	if _, err := vars.Get(AUTO_INCREMENT_KEY); err == nil {
		t.Error("auto increment key in variables is not deleted")
	}
}
