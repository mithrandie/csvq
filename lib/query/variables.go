package query

import (
	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/value"
)

type VariablesList []Variables

func (list VariablesList) Declare(expr parser.VariableDeclaration, filter *Filter) error {
	return list[0].Declare(expr, filter)
}

func (list VariablesList) Get(expr parser.Variable) (value value.Primary, err error) {
	for _, v := range list {
		if value, err = v.Get(expr); err == nil {
			return
		}
	}
	err = NewUndefinedVariableError(expr)
	return
}

func (list VariablesList) Substitute(expr parser.VariableSubstitution, filter *Filter) (value value.Primary, err error) {
	for _, v := range list {
		if value, err = v.Substitute(expr, filter); err == nil {
			return
		}
		if _, ok := err.(*UndefinedVariableError); !ok {
			return
		}
	}
	err = NewUndefinedVariableError(expr.Variable)
	return
}

func (list VariablesList) SubstitutePrimary(variable parser.Variable, value value.Primary) (value.Primary, error) {
	var err error
	for _, v := range list {
		if value, err = v.SubstitutePrimary(variable, value); err == nil {
			return value, nil
		}
	}
	return nil, NewUndefinedVariableError(variable)
}

func (list VariablesList) Dispose(expr parser.Variable) error {
	for _, v := range list {
		if err := v.Dispose(expr); err == nil {
			return nil
		}
	}
	return NewUndefinedVariableError(expr)
}

type Variables map[string]value.Primary

func (v Variables) Add(variable parser.Variable, value value.Primary) error {
	if _, ok := v[variable.Name]; ok {
		return NewVariableRedeclaredError(variable)
	}
	v[variable.Name] = value
	return nil
}

func (v Variables) Set(variable parser.Variable, value value.Primary) error {
	if _, ok := v[variable.Name]; !ok {
		return NewUndefinedVariableError(variable)
	}
	v[variable.Name] = value
	return nil
}

func (v Variables) Get(variable parser.Variable) (value.Primary, error) {
	if v, ok := v[variable.Name]; ok {
		return v, nil
	}
	return nil, NewUndefinedVariableError(variable)
}

func (v Variables) Dispose(variable parser.Variable) error {
	if _, ok := v[variable.Name]; !ok {
		return NewUndefinedVariableError(variable)
	}
	delete(v, variable.Name)
	return nil
}

func (v Variables) Declare(declaration parser.VariableDeclaration, filter *Filter) error {
	for _, a := range declaration.Assignments {
		assignment := a.(parser.VariableAssignment)
		var val value.Primary
		var err error
		if assignment.Value == nil {
			val = value.NewNull()
		} else {
			val, err = filter.Evaluate(assignment.Value)
			if err != nil {
				return err
			}
		}
		err = v.Add(assignment.Variable, val)
		if err != nil {
			return err
		}
	}
	return nil
}

func (v Variables) Substitute(substitution parser.VariableSubstitution, filter *Filter) (value.Primary, error) {
	val, err := filter.Evaluate(substitution.Value)
	if err != nil {
		return nil, err
	}
	return v.SubstitutePrimary(substitution.Variable, val)
}

func (v Variables) SubstitutePrimary(variable parser.Variable, value value.Primary) (value.Primary, error) {
	err := v.Set(variable, value)
	return value, err
}
