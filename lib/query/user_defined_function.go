package query

import (
	"strings"

	"github.com/mithrandie/csvq/lib/parser"
)

type UserDefinedFunctionMap map[string]*UserDefinedFunction

func (m UserDefinedFunctionMap) Declare(expr parser.FunctionDeclaration) error {
	uname := strings.ToUpper(expr.Name.Literal)

	if _, ok := Functions[uname]; ok {
		return NewBuiltInFunctionDeclaredError(expr.Name)
	}
	if _, ok := AggregateFunctions[uname]; ok {
		return NewBuiltInFunctionDeclaredError(expr.Name)
	}
	if _, ok := AnalyticFunctions[uname]; ok {
		return NewBuiltInFunctionDeclaredError(expr.Name)
	}
	if uname == "GROUP_CONCAT" {
		return NewBuiltInFunctionDeclaredError(expr.Name)
	}
	if _, ok := m[uname]; ok {
		return NewFunctionRedeclaredError(expr.Name)
	}

	m[uname] = &UserDefinedFunction{
		Name:       expr.Name,
		Parameters: expr.Parameters,
		Statements: expr.Statements,
	}
	return nil
}

func (m UserDefinedFunctionMap) Get(fn parser.Function) (*UserDefinedFunction, error) {
	uname := strings.ToUpper(fn.Name)
	if fn, ok := m[uname]; ok {
		return fn, nil
	}
	return nil, NewFunctionNotExistError(fn.Name, fn)
}

type UserDefinedFunction struct {
	Name       parser.Identifier
	Parameters []parser.Variable
	Statements []parser.Statement
}

func (fn *UserDefinedFunction) Execute(args []parser.Primary, filter Filter) (parser.Primary, error) {
	if len(args) != len(fn.Parameters) {
		return nil, NewFunctionArgumentLengthError(fn.Name.Literal, []int{len(fn.Parameters)}, fn.Name)
	}

	proc := NewProcedure()
	vars := Variables{}
	for i, v := range fn.Parameters {
		if err := vars.Add(v, args[i]); err != nil {
			return nil, err
		}
	}
	proc.Filter = NewFilter(
		append([]Variables{vars}, filter.VariablesList...),
		append([]ViewMap{{}}, filter.TempViewsList...),
		append([]CursorMap{{}}, filter.CursorsList...),
	)

	if _, err := proc.Execute(fn.Statements); err != nil {
		return nil, err
	}

	ret := proc.ReturnVal
	if ret == nil {
		ret = parser.NewNull()
	}
	return ret, nil
}
