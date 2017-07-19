package query

import (
	"errors"
	"fmt"
	"strings"

	"github.com/mithrandie/csvq/lib/parser"
)

type UserFunctionMap map[string]*UserFunction

func (m UserFunctionMap) Declare(expr parser.FunctionDeclaration) error {
	uname := strings.ToUpper(expr.Name.Literal)

	if _, ok := Functions[uname]; ok {
		return errors.New(fmt.Sprintf("function %s is redeclared", expr.Name.Literal))
	}
	if _, ok := AggregateFunctions[uname]; ok {
		return errors.New(fmt.Sprintf("function %s is redeclared", expr.Name.Literal))
	}
	if _, ok := AnalyticFunctions[uname]; ok {
		return errors.New(fmt.Sprintf("function %s is redeclared", expr.Name.Literal))
	}
	if uname == "GROUP_CONCAT" {
		return errors.New(fmt.Sprintf("function %s is redeclared", expr.Name.Literal))
	}
	if _, ok := m[uname]; ok {
		return errors.New(fmt.Sprintf("function %s is redeclared", expr.Name.Literal))
	}

	m[uname] = &UserFunction{
		Name:       expr.Name.Literal,
		Parameters: expr.Parameters,
		Statements: expr.Statements,
	}
	return nil
}

func (m UserFunctionMap) Get(key string) (*UserFunction, error) {
	uname := strings.ToUpper(key)
	if fn, ok := m[uname]; ok {
		return fn, nil
	}
	return nil, errors.New(fmt.Sprintf("function %s does not exist", key))
}

type UserFunction struct {
	Name       string
	Parameters []parser.Variable
	Statements []parser.Statement
}

func (fn *UserFunction) Execute(args []parser.Primary, filter Filter) (parser.Primary, error) {
	if len(args) != len(fn.Parameters) {
		return nil, errors.New(fmt.Sprintf("declared function %s takes %d argument(s)", fn.Name, len(fn.Parameters)))
	}

	proc := NewProcedure()
	vars := Variables{}
	for i, v := range fn.Parameters {
		vars.Add(v.Name, args[i])
	}
	proc.VariablesList = append([]Variables{vars}, filter.VariablesList...)

	if _, err := proc.Execute(fn.Statements); err != nil {
		return nil, err
	}

	ret := proc.ReturnVal
	if ret == nil {
		ret = parser.NewNull()
	}
	return ret, nil
}
