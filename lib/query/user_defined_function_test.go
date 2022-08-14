package query

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/value"
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
		Result: GenerateUserDefinedFunctionMap([]*UserDefinedFunction{
			{
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
		}),
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
		Error: "function userfunc is redeclared",
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
		Error: "parameter @arg1 is a duplicate",
	},
}

func TestUserDefinedFunctionMap_Declare(t *testing.T) {
	funcs := NewUserDefinedFunctionMap()

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
		if !SyncMapEqual(funcs, v.Result) {
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
		Result: GenerateUserDefinedFunctionMap([]*UserDefinedFunction{
			{
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
		}),
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
		Error: "function useraggfunc is redeclared",
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
		Error: "parameter @arg1 is a duplicate",
	},
}

func TestUserDefinedFunctionMap_DeclareAggregate(t *testing.T) {
	funcs := NewUserDefinedFunctionMap()

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
		if !SyncMapEqual(funcs, v.Result) {
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
		Error:    "function userfunc is redeclared",
	},
	{
		Name:     "UserDefinedFunctionMap CheckDuplicate Duplicate with Built-in Function Error",
		FuncName: parser.Identifier{Literal: "now"},
		Error:    "function now is a built-in function",
	},
	{
		Name:     "UserDefinedFunctionMap CheckDuplicate Duplicate with Aggregate Function Error",
		FuncName: parser.Identifier{Literal: "count"},
		Error:    "function count is a built-in function",
	},
	{
		Name:     "UserDefinedFunctionMap CheckDuplicate Duplicate with Analytic Function Error",
		FuncName: parser.Identifier{Literal: "row_number"},
		Error:    "function row_number is a built-in function",
	},
	{
		Name:     "UserDefinedFunctionMap CheckDuplicate OK",
		FuncName: parser.Identifier{Literal: "undefined"},
		Result:   true,
	},
}

func TestUserDefinedFunctionMap_CheckDuplicate(t *testing.T) {
	funcs := GenerateUserDefinedFunctionMap([]*UserDefinedFunction{
		{
			Name: parser.Identifier{Literal: "userfunc"},
			Parameters: []parser.Variable{
				{Name: "arg1"},
				{Name: "arg2"},
			},
			Statements: []parser.Statement{
				parser.Print{Value: parser.Variable{Name: "var1"}},
			},
		},
	})

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
	OK       bool
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
		OK: true,
	},
	{
		Name: "UserDefinedFunctionMap Get Not Exist Error",
		Function: parser.Function{
			Name: "notexist",
		},
		FuncName: "notexist",
		OK:       false,
	},
}

func TestUserDefinedFunctionMap_Get(t *testing.T) {
	funcs := GenerateUserDefinedFunctionMap([]*UserDefinedFunction{
		{
			Name: parser.Identifier{Literal: "userfunc"},
			Parameters: []parser.Variable{
				{Name: "arg1"},
				{Name: "arg2"},
			},
			Statements: []parser.Statement{
				parser.Print{Value: parser.Variable{Name: "var1"}},
			},
		},
	})

	for _, v := range userDefinedFunctionMapGetTests {
		result, ok := funcs.Get(v.FuncName)
		if ok != v.OK {
			t.Errorf("%s: result = %t, want %t", v.Name, ok, v.OK)
			continue
		}
		if ok && v.OK && !reflect.DeepEqual(result, v.Result) {
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
									Operator: parser.Token{Token: '+', Literal: "+"},
								},
								RHS:      parser.Variable{Name: "var1"},
								Operator: parser.Token{Token: '+', Literal: "+"},
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
									Operator: parser.Token{Token: '+', Literal: "+"},
								},
								RHS:      parser.Variable{Name: "var1"},
								Operator: parser.Token{Token: '+', Literal: "+"},
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
									Operator: parser.Token{Token: '+', Literal: "+"},
								},
								RHS:      parser.Variable{Name: "var1"},
								Operator: parser.Token{Token: '+', Literal: "+"},
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
		Error: "function userfunc takes exactly 2 arguments",
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
									Operator: parser.Token{Token: '+', Literal: "+"},
								},
								RHS:      parser.Variable{Name: "var1"},
								Operator: parser.Token{Token: '+', Literal: "+"},
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
		Error: "field notexist does not exist",
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
		Error: "field notexist does not exist",
	},
}

func TestUserDefinedFunction_Execute(t *testing.T) {
	scope := GenerateReferenceScope([]map[string]map[string]interface{}{
		{
			scopeNameVariables: {
				"var1": value.NewInteger(1),
			},
		},
	}, nil, time.Time{}, nil)

	ctx := context.Background()
	for _, v := range userDefinedFunctionExecuteTests {
		result, err := v.Func.Execute(ctx, scope, v.Args)
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
								Operator: parser.Token{Token: '*', Literal: "*"},
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
								Operator: parser.Token{Token: '*', Literal: "*"},
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
								Operator: parser.Token{Token: '*', Literal: "*"},
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
		Error: "function useraggfunc takes exactly 1 argument",
	},
}

func TestUserDefinedFunction_ExecuteAggregate(t *testing.T) {
	scope := NewReferenceScope(TestTx)
	ctx := context.Background()

	for _, v := range userDefinedFunctionExecuteAggregateTests {
		result, err := v.Func.ExecuteAggregate(ctx, scope, v.Values, v.Args)
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
		Error:   "function userfunc takes exactly 2 arguments",
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
		Error:   "function userfunc takes at least 1 argument",
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
		Error:   "function userfunc takes at most 2 arguments",
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
		Error:   "function userfunc takes exactly 3 arguments",
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
