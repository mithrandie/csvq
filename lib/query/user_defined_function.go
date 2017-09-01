package query

import (
	"fmt"
	"strings"

	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/value"
)

type UserDefinedFunctionsList []UserDefinedFunctionMap

func (list UserDefinedFunctionsList) Declare(expr parser.FunctionDeclaration) error {
	return list[0].Declare(expr)
}

func (list UserDefinedFunctionsList) DeclareAggregate(expr parser.AggregateDeclaration) error {
	return list[0].DeclareAggregate(expr)
}

func (list UserDefinedFunctionsList) Get(expr parser.QueryExpression, name string) (*UserDefinedFunction, error) {
	for _, v := range list {
		if fn, err := v.Get(expr, name); err == nil {
			return fn, nil
		}
	}
	return nil, NewFunctionNotExistError(expr, name)
}

type UserDefinedFunctionMap map[string]*UserDefinedFunction

func (m UserDefinedFunctionMap) Declare(expr parser.FunctionDeclaration) error {
	if err := m.CheckDuplicate(expr.Name); err != nil {
		return err
	}

	parameters, defaults, required, err := m.parserParameters(expr.Parameters)
	if err != nil {
		return err
	}

	uname := strings.ToUpper(expr.Name.Literal)

	m[uname] = &UserDefinedFunction{
		Name:         expr.Name,
		Statements:   expr.Statements,
		Parameters:   parameters,
		Defaults:     defaults,
		RequiredArgs: required,
	}
	return nil
}

func (m UserDefinedFunctionMap) DeclareAggregate(expr parser.AggregateDeclaration) error {
	if err := m.CheckDuplicate(expr.Name); err != nil {
		return err
	}

	parameters, defaults, required, err := m.parserParameters(expr.Parameters)
	if err != nil {
		return err
	}

	uname := strings.ToUpper(expr.Name.Literal)

	m[uname] = &UserDefinedFunction{
		Name:         expr.Name,
		Statements:   expr.Statements,
		Parameters:   parameters,
		Defaults:     defaults,
		RequiredArgs: required,
		IsAggregate:  true,
		Cursor:       expr.Cursor,
	}
	return nil
}

func (m UserDefinedFunctionMap) parserParameters(parameters []parser.Expression) ([]parser.Variable, map[string]parser.QueryExpression, int, error) {
	var isDuplicate = func(variable parser.Variable, variables []parser.Variable) bool {
		for _, v := range variables {
			if variable.Name == v.Name {
				return true
			}
		}
		return false
	}

	variables := make([]parser.Variable, len(parameters))
	defaults := make(map[string]parser.QueryExpression)

	required := 0
	for i, parameter := range parameters {
		assignment := parameter.(parser.VariableAssignment)

		if isDuplicate(assignment.Variable, variables) {
			return nil, nil, 0, NewDuplicateParameterError(assignment.Variable)
		}

		variables[i] = assignment.Variable
		if assignment.Value == nil {
			required = i + 1
		} else {
			defaults[assignment.Variable.String()] = assignment.Value
		}
	}
	return variables, defaults, required, nil
}

func (m UserDefinedFunctionMap) CheckDuplicate(name parser.Identifier) error {
	uname := strings.ToUpper(name.Literal)

	if _, ok := Functions[uname]; ok || uname == "NOW" {
		return NewBuiltInFunctionDeclaredError(name)
	}
	if _, ok := AggregateFunctions[uname]; ok {
		return NewBuiltInFunctionDeclaredError(name)
	}
	if _, ok := AnalyticFunctions[uname]; ok {
		return NewBuiltInFunctionDeclaredError(name)
	}
	if _, ok := m[uname]; ok {
		return NewFunctionRedeclaredError(name)
	}
	return nil
}

func (m UserDefinedFunctionMap) Get(fn parser.QueryExpression, name string) (*UserDefinedFunction, error) {
	uname := strings.ToUpper(name)
	if fn, ok := m[uname]; ok {
		return fn, nil
	}
	return nil, NewFunctionNotExistError(fn, name)
}

type UserDefinedFunction struct {
	Name         parser.Identifier
	Statements   []parser.Statement
	Parameters   []parser.Variable
	Defaults     map[string]parser.QueryExpression
	RequiredArgs int

	IsAggregate bool
	Cursor      parser.Identifier // For Aggregate Functions
}

func (fn *UserDefinedFunction) Execute(args []value.Primary, filter *Filter) (value.Primary, error) {
	childScope := filter.CreateChildScope()
	return fn.execute(args, childScope)
}

func (fn *UserDefinedFunction) ExecuteAggregate(values []value.Primary, args []value.Primary, filter *Filter) (value.Primary, error) {
	childScope := filter.CreateChildScope()
	childScope.CursorsList.AddPseudoCursor(fn.Cursor, values)
	return fn.execute(args, childScope)
}

func (fn *UserDefinedFunction) CheckArgsLen(expr parser.QueryExpression, name string, argsLen int) error {
	parametersLen := len(fn.Parameters)
	requiredLen := fn.RequiredArgs
	if fn.IsAggregate {
		parametersLen++
		requiredLen++
	}

	if len(fn.Defaults) < 1 {
		if argsLen != len(fn.Parameters) {
			return NewFunctionArgumentLengthError(expr, name, []int{parametersLen})
		}
	} else if argsLen < fn.RequiredArgs {
		return NewFunctionArgumentLengthErrorWithCustomArgs(expr, name, fmt.Sprintf("at least %s", FormatCount(requiredLen, "argument")))
	} else if len(fn.Parameters) < argsLen {
		return NewFunctionArgumentLengthErrorWithCustomArgs(expr, name, fmt.Sprintf("at most %s", FormatCount(parametersLen, "argument")))
	}

	return nil
}

func (fn *UserDefinedFunction) execute(args []value.Primary, filter *Filter) (value.Primary, error) {
	if err := fn.CheckArgsLen(fn.Name, fn.Name.Literal, len(args)); err != nil {
		return nil, err
	}

	for i, v := range fn.Parameters {
		if i < len(args) {
			filter.VariablesList[0].Add(v, args[i])
		} else {
			defaultValue, _ := fn.Defaults[v.String()]
			val, err := filter.Evaluate(defaultValue)
			if err != nil {
				return nil, err
			}
			filter.VariablesList[0].Add(v, val)
		}
	}

	proc := NewProcedure()
	proc.Filter = filter

	if _, err := proc.Execute(fn.Statements); err != nil {
		return nil, err
	}

	ret := proc.ReturnVal
	if ret == nil {
		ret = value.NewNull()
	}

	return ret, nil
}
