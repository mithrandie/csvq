package query

import (
	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/value"
)

type VariableScopes []VariableMap

func (list VariableScopes) Declare(expr parser.VariableDeclaration, filter *Filter) error {
	return list[0].Declare(expr, filter)
}

func (list VariableScopes) Get(expr parser.Variable) (value value.Primary, err error) {
	for _, v := range list {
		if value, err = v.Get(expr); err == nil {
			return
		}
	}
	err = NewUndeclaredVariableError(expr)
	return
}

func (list VariableScopes) Substitute(expr parser.VariableSubstitution, filter *Filter) (value value.Primary, err error) {
	for _, v := range list {
		if value, err = v.Substitute(expr, filter); err == nil {
			return
		}
		if _, ok := err.(*UndeclaredVariableError); !ok {
			return
		}
	}
	err = NewUndeclaredVariableError(expr.Variable)
	return
}

func (list VariableScopes) SubstituteDirectly(variable parser.Variable, value value.Primary) (value.Primary, error) {
	var err error
	for _, v := range list {
		if value, err = v.SubstituteDirectly(variable, value); err == nil {
			return value, nil
		}
	}
	return nil, NewUndeclaredVariableError(variable)
}

func (list VariableScopes) Dispose(expr parser.Variable) error {
	for _, v := range list {
		if err := v.Dispose(expr); err == nil {
			return nil
		}
	}
	return NewUndeclaredVariableError(expr)
}

type VariableMap map[string]value.Primary

func (v VariableMap) Add(variable parser.Variable, value value.Primary) error {
	if _, ok := v[variable.Name]; ok {
		return NewVariableRedeclaredError(variable)
	}
	v[variable.Name] = value
	return nil
}

func (v VariableMap) Set(variable parser.Variable, value value.Primary) error {
	if _, ok := v[variable.Name]; !ok {
		return NewUndeclaredVariableError(variable)
	}
	v[variable.Name] = value
	return nil
}

func (v VariableMap) Get(variable parser.Variable) (value.Primary, error) {
	if v, ok := v[variable.Name]; ok {
		return v, nil
	}
	return nil, NewUndeclaredVariableError(variable)
}

func (v VariableMap) Dispose(variable parser.Variable) error {
	if _, ok := v[variable.Name]; !ok {
		return NewUndeclaredVariableError(variable)
	}
	delete(v, variable.Name)
	return nil
}

func (v VariableMap) Declare(declaration parser.VariableDeclaration, filter *Filter) error {
	for _, assignment := range declaration.Assignments {
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

func (v VariableMap) Substitute(substitution parser.VariableSubstitution, filter *Filter) (value.Primary, error) {
	val, err := filter.Evaluate(substitution.Value)
	if err != nil {
		return nil, err
	}
	return v.SubstituteDirectly(substitution.Variable, val)
}

func (v VariableMap) SubstituteDirectly(variable parser.Variable, value value.Primary) (value.Primary, error) {
	err := v.Set(variable, value)
	return value, err
}
