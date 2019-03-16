package query

import (
	"context"
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
		Expr:   parser.Variable{Name: "var1"},
		Result: value.NewInteger(1),
	},
	{
		Name:  "VariableScopes Get Undeclared Error",
		Expr:  parser.Variable{Name: "undef"},
		Error: "[L:- C:-] variable @undef is undeclared",
	},
}

func TestVariableScopes_Get(t *testing.T) {
	list := VariableScopes{
		GenerateVariableMap(map[string]value.Primary{
			"var1": value.NewInteger(1),
		}),
		GenerateVariableMap(map[string]value.Primary{
			"var1": value.NewInteger(2),
		}),
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
			Variable: parser.Variable{Name: "var1"},
			Value:    parser.NewIntegerValue(3),
		},
		List: VariableScopes{
			GenerateVariableMap(map[string]value.Primary{
				"var1": value.NewInteger(3),
			}),
			GenerateVariableMap(map[string]value.Primary{
				"var1": value.NewInteger(2),
			}),
		},
		Result: value.NewInteger(3),
	},
	{
		Name: "VariableScopes Substitute Variable Undeclared Error",
		Expr: parser.VariableSubstitution{
			Variable: parser.Variable{Name: "var2"},
			Value:    parser.NewIntegerValue(3),
		},
		Error: "[L:- C:-] variable @var2 is undeclared",
	},
}

func TestVariableScopes_Substitute(t *testing.T) {
	list := VariableScopes{
		GenerateVariableMap(map[string]value.Primary{
			"var1": value.NewInteger(1),
		}),
		GenerateVariableMap(map[string]value.Primary{
			"var1": value.NewInteger(2),
		}),
	}

	for _, v := range variableScopesSubstituteTests {
		result, err := list.Substitute(context.Background(), v.Filter, v.Expr)
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
		if !list.Equal(v.List) {
			t.Errorf("%s: list = %v, want %v", v.Name, list, v.List)
		}
		if !reflect.DeepEqual(result, v.Result) {
			t.Errorf("%s: result = %v, want %v", v.Name, result, v.Result)
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
		Expr: parser.Variable{Name: "var1"},
		List: VariableScopes{
			NewVariableMap(),
			GenerateVariableMap(map[string]value.Primary{
				"var1": value.NewInteger(2),
			}),
		},
	},
	{
		Name:  "VariableScopes Dispose Undeclared Error",
		Expr:  parser.Variable{Name: "undef"},
		Error: "[L:- C:-] variable @undef is undeclared",
	},
}

func TestVariableScopes_Dispose(t *testing.T) {
	list := VariableScopes{
		GenerateVariableMap(map[string]value.Primary{
			"var1": value.NewInteger(1),
		}),
		GenerateVariableMap(map[string]value.Primary{
			"var1": value.NewInteger(2),
		}),
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
		if !list.Equal(v.List) {
			t.Errorf("%s: list = %v, want %v", v.Name, list, v.List)
		}
	}
}

var variableScopesEqualTests = []struct {
	Name   string
	List1  VariableScopes
	List2  VariableScopes
	Expect bool
}{
	{
		Name: "Equal",
		List1: VariableScopes{
			GenerateVariableMap(map[string]value.Primary{
				"var1": value.NewInteger(1),
			}),
			GenerateVariableMap(map[string]value.Primary{
				"var1": value.NewInteger(2),
			}),
		},
		List2: VariableScopes{
			GenerateVariableMap(map[string]value.Primary{
				"var1": value.NewInteger(1),
			}),
			GenerateVariableMap(map[string]value.Primary{
				"var1": value.NewInteger(2),
			}),
		},
		Expect: true,
	},
	{
		Name: "Different Length",
		List1: VariableScopes{
			GenerateVariableMap(map[string]value.Primary{
				"var1": value.NewInteger(1),
			}),
			GenerateVariableMap(map[string]value.Primary{
				"var1": value.NewInteger(2),
			}),
			GenerateVariableMap(map[string]value.Primary{
				"var1": value.NewInteger(3),
			}),
		},
		List2: VariableScopes{
			GenerateVariableMap(map[string]value.Primary{
				"var1": value.NewInteger(1),
			}),
			GenerateVariableMap(map[string]value.Primary{
				"var1": value.NewInteger(2),
			}),
		},
		Expect: false,
	},
	{
		Name: "Different Value",
		List1: VariableScopes{
			GenerateVariableMap(map[string]value.Primary{
				"var1": value.NewInteger(1),
			}),
			GenerateVariableMap(map[string]value.Primary{
				"var1": value.NewInteger(2),
			}),
		},
		List2: VariableScopes{
			GenerateVariableMap(map[string]value.Primary{
				"var1": value.NewInteger(1),
			}),
			GenerateVariableMap(map[string]value.Primary{
				"var1": value.NewInteger(3),
			}),
		},
		Expect: false,
	},
}

func TestVariableScopes_Equal(t *testing.T) {
	for _, v := range variableScopesEqualTests {
		result := v.List1.Equal(v.List2)
		if result != v.Expect {
			t.Errorf("%s: result = %t, want %t", v.Name, result, v.Expect)
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
					Variable: parser.Variable{Name: "var1"},
				},
			},
		},
		Result: GenerateVariableMap(map[string]value.Primary{
			"var1": value.NewNull(),
		}),
	},
	{
		Name: "Declare Variable With Initial Value",
		Expr: parser.VariableDeclaration{
			Assignments: []parser.VariableAssignment{
				{
					Variable: parser.Variable{Name: "var2"},
					Value:    parser.NewIntegerValue(1),
				},
			},
		},
		Result: GenerateVariableMap(map[string]value.Primary{
			"var1": value.NewNull(),
			"var2": value.NewInteger(1),
		}),
	},
	{
		Name: "Declare Variable Redeclaration Error",
		Expr: parser.VariableDeclaration{
			Assignments: []parser.VariableAssignment{
				{
					Variable: parser.Variable{Name: "var2"},
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
					Variable: parser.Variable{Name: "var3"},
					Value:    parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
				},
			},
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
}

func TestVariableMap_Declare(t *testing.T) {
	vars := NewVariableMap()

	for _, v := range variableMapDeclareTests {
		if v.Filter == nil {
			v.Filter = NewEmptyFilter(TestTx)
		}

		err := vars.Declare(context.Background(), v.Filter, v.Expr.(parser.VariableDeclaration))
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
		if !vars.Equal(&v.Result) {
			t.Errorf("%s: result = %v, want %v", v.Name, vars, v.Result)
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
		Result: GenerateVariableMap(map[string]value.Primary{
			"var1": value.NewInteger(2),
		}),
	},
	{
		Name: "Substitute Variable Undeclared Error",
		Expr: parser.VariableSubstitution{
			Variable: parser.Variable{Name: "var2"},
			Value:    parser.NewIntegerValue(2),
		},
		Error: "[L:- C:-] variable @var2 is undeclared",
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
	vars := GenerateVariableMap(map[string]value.Primary{
		"var1": value.NewInteger(1),
	})

	for _, v := range variableMapSubstituteTests {
		if v.Filter == nil {
			v.Filter = NewEmptyFilter(TestTx)
		}

		_, err := vars.Substitute(context.Background(), v.Filter, v.Expr.(parser.VariableSubstitution))
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
		if !vars.Equal(&v.Result) {
			t.Errorf("%s: result = %v, want %v", v.Name, vars, v.Result)
		}
	}
}

var variableMapDisposeTests = []variableMapTests{
	{
		Name:   "Dispose Variable",
		Expr:   parser.Variable{Name: "var1"},
		Result: NewVariableMap(),
	},
	{
		Name:  "Dispose Variable Undeclared Error",
		Expr:  parser.Variable{Name: "var2"},
		Error: "[L:- C:-] variable @var2 is undeclared",
	},
}

func TestVariableMap_Dispose(t *testing.T) {
	vars := GenerateVariableMap(map[string]value.Primary{
		"var1": value.NewInteger(1),
	})

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
		if !vars.Equal(&v.Result) {
			t.Errorf("%s: result = %v, want %v", v.Name, vars, v.Result)
		}
	}
}
