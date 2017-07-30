package query

import (
	"strings"

	"github.com/mithrandie/csvq/lib/parser"
)

type UserDefinedFunctionsList []UserDefinedFunctionMap

func (list UserDefinedFunctionsList) Declare(expr parser.FunctionDeclaration) error {
	return list[0].Declare(expr)
}

func (list UserDefinedFunctionsList) DeclareAggregate(expr parser.AggregateDeclaration) error {
	return list[0].DeclareAggregate(expr)
}

func (list UserDefinedFunctionsList) Get(expr parser.Expression, name string) (*UserDefinedFunction, error) {
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

	uname := strings.ToUpper(expr.Name.Literal)

	m[uname] = &UserDefinedFunction{
		Name:       expr.Name,
		Statements: expr.Statements,
		Parameters: expr.Parameters,
	}
	return nil
}

func (m UserDefinedFunctionMap) DeclareAggregate(expr parser.AggregateDeclaration) error {
	if err := m.CheckDuplicate(expr.Name); err != nil {
		return err
	}

	uname := strings.ToUpper(expr.Name.Literal)

	m[uname] = &UserDefinedFunction{
		Name:        expr.Name,
		Statements:  expr.Statements,
		IsAggregate: true,
		Parameter:   expr.Parameter,
	}
	return nil
}

func (m UserDefinedFunctionMap) CheckDuplicate(name parser.Identifier) error {
	uname := strings.ToUpper(name.Literal)

	if _, ok := Functions[uname]; ok {
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

func (m UserDefinedFunctionMap) Get(fn parser.Expression, name string) (*UserDefinedFunction, error) {
	uname := strings.ToUpper(name)
	if fn, ok := m[uname]; ok {
		return fn, nil
	}
	return nil, NewFunctionNotExistError(fn, name)
}

type UserDefinedFunction struct {
	Name        parser.Identifier
	Statements  []parser.Statement
	IsAggregate bool

	// For Scala Functions
	Parameters []parser.Variable

	// For Aggregate Functions
	Parameter parser.Identifier
}

func (fn *UserDefinedFunction) Execute(args []parser.Primary, filter Filter) (parser.Primary, error) {
	var err error

	proc := NewProcedure()

	if fn.IsAggregate {
		proc.Filter = fn.generateAggregateFilter(args, filter)
	} else {
		proc.Filter, err = fn.generateScalaFilter(args, filter)
		if err != nil {
			return nil, err
		}
	}

	if _, err = proc.Execute(fn.Statements); err != nil {
		return nil, err
	}

	ret := proc.ReturnVal
	if ret == nil {
		ret = parser.NewNull()
	}

	return ret, nil
}

func (fn *UserDefinedFunction) generateScalaFilter(args []parser.Primary, filter Filter) (Filter, error) {
	if len(args) != len(fn.Parameters) {
		return filter, NewFunctionArgumentLengthError(fn.Name, fn.Name.Literal, []int{len(fn.Parameters)})
	}

	f := filter.CreateChildScope()

	for i, v := range fn.Parameters {
		if err := f.VariablesList[0].Add(v, args[i]); err != nil {
			return filter, err
		}
	}
	return f, nil
}

func (fn *UserDefinedFunction) generateAggregateFilter(values []parser.Primary, filter Filter) Filter {
	f := filter.CreateChildScope()
	f.CursorsList.AddPseudoCursor(fn.Parameter, values)
	return f
}
