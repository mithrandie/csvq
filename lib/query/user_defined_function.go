package query

import (
	"context"
	"fmt"
	"strings"

	"sort"

	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/value"
)

type UserDefinedFunctionScopes []UserDefinedFunctionMap

func (list UserDefinedFunctionScopes) Declare(expr parser.FunctionDeclaration) error {
	return list[0].Declare(expr)
}

func (list UserDefinedFunctionScopes) DeclareAggregate(expr parser.AggregateDeclaration) error {
	return list[0].DeclareAggregate(expr)
}

func (list UserDefinedFunctionScopes) Get(expr parser.QueryExpression, name string) (*UserDefinedFunction, error) {
	for _, v := range list {
		if fn, err := v.Get(expr, name); err == nil {
			return fn, nil
		}
	}
	return nil, NewFunctionNotExistError(expr, name)
}

func (list UserDefinedFunctionScopes) Dispose(name parser.Identifier) error {
	for _, m := range list {
		err := m.Dispose(name)
		if err == nil {
			return nil
		}
	}
	return NewFunctionNotExistError(name, name.Literal)
}

func (list UserDefinedFunctionScopes) All() (UserDefinedFunctionMap, UserDefinedFunctionMap) {
	scalaAll := make(UserDefinedFunctionMap, 10)
	aggregateAll := make(UserDefinedFunctionMap, 10)

	for _, m := range list {
		for key, fn := range m {
			if fn.IsAggregate {
				if _, ok := aggregateAll[key]; !ok {
					aggregateAll[key] = fn
				}
			} else {
				if _, ok := scalaAll[key]; !ok {
					scalaAll[key] = fn
				}
			}
		}
	}

	return scalaAll, aggregateAll
}

type UserDefinedFunctionMap map[string]*UserDefinedFunction

func (m UserDefinedFunctionMap) Declare(expr parser.FunctionDeclaration) error {
	if err := m.CheckDuplicate(expr.Name); err != nil {
		return err
	}

	parameters, defaults, required, err := m.parseParameters(expr.Parameters)
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

	parameters, defaults, required, err := m.parseParameters(expr.Parameters)
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

func (m UserDefinedFunctionMap) parseParameters(parameters []parser.VariableAssignment) ([]parser.Variable, map[string]parser.QueryExpression, int, error) {
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
	for i, assignment := range parameters {
		if isDuplicate(assignment.Variable, variables) {
			return nil, nil, 0, NewDuplicateParameterError(assignment.Variable)
		}

		variables[i] = assignment.Variable
		if assignment.Value == nil {
			required = i + 1
		} else {
			defaults[assignment.Variable.Name] = assignment.Value
		}
	}
	return variables, defaults, required, nil
}

func (m UserDefinedFunctionMap) CheckDuplicate(name parser.Identifier) error {
	uname := strings.ToUpper(name.Literal)

	if _, ok := Functions[uname]; ok || uname == "CALL" || uname == "NOW" || uname == "JSON_OBJECT" {
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

func (m UserDefinedFunctionMap) Keys() []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func (m UserDefinedFunctionMap) SortedKeys() []string {
	keys := m.Keys()
	sort.Strings(keys)
	return keys
}

func (m UserDefinedFunctionMap) Dispose(name parser.Identifier) error {
	uname := strings.ToUpper(name.Literal)
	if _, ok := m[uname]; ok {
		delete(m, uname)
		return nil
	}
	return NewFunctionNotExistError(name, name.Literal)
}

func (m UserDefinedFunctionMap) Clear() {
	for k := range m {
		delete(m, k)
	}
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

func (fn *UserDefinedFunction) Execute(ctx context.Context, filter *Filter, args []value.Primary) (value.Primary, error) {
	childScope := filter.CreateChildScope()
	defer childScope.CloseScope()

	return fn.execute(ctx, childScope, args)
}

func (fn *UserDefinedFunction) ExecuteAggregate(ctx context.Context, filter *Filter, values []value.Primary, args []value.Primary) (value.Primary, error) {
	childScope := filter.CreateChildScope()
	defer childScope.CloseScope()

	if err := childScope.cursors.AddPseudoCursor(filter.tx, fn.Cursor, values); err != nil {
		return nil, err
	}
	return fn.execute(ctx, childScope, args)
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

func (fn *UserDefinedFunction) execute(ctx context.Context, filter *Filter, args []value.Primary) (value.Primary, error) {
	if err := fn.CheckArgsLen(fn.Name, fn.Name.Literal, len(args)); err != nil {
		return nil, err
	}

	for i, v := range fn.Parameters {
		if i < len(args) {
			if err := filter.variables[0].Add(v, args[i]); err != nil {
				return nil, err
			}
		} else {
			defaultValue, _ := fn.Defaults[v.Name]
			val, err := filter.Evaluate(ctx, defaultValue)
			if err != nil {
				return nil, err
			}
			if err = filter.variables[0].Add(v, val); err != nil {
				return nil, err
			}
		}
	}

	proc := NewProcessorWithFilter(filter.tx, filter)
	if _, err := proc.execute(ctx, fn.Statements); err != nil {
		return nil, err
	}

	ret := proc.returnVal
	if ret == nil {
		ret = value.NewNull()
	}

	return ret, nil
}
