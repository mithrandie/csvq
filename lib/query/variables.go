package query

import (
	"github.com/mithrandie/csvq/lib/parser"
)

type Variables map[string]parser.Primary

func (v Variables) Add(key string, value parser.Primary) error {
	if _, ok := v[key]; ok {
		return NewRedeclaredVariableError(key)
	}
	v[key] = value
	return nil
}

func (v Variables) Set(key string, value parser.Primary) error {
	if _, ok := v[key]; !ok {
		return NewUndefinedVariableError(key)
	}
	v[key] = value
	return nil
}

func (v Variables) Get(key string) (parser.Primary, error) {
	if v, ok := v[key]; ok {
		return v, nil
	}
	return nil, NewUndefinedVariableError(key)
}

func (v Variables) Delete(key string) {
	if _, ok := v[key]; ok {
		delete(v, key)
	}
}

func (v Variables) Decrare(declaration parser.VariableDeclaration, filter Filter) error {
	for _, a := range declaration.Assignments {
		assignment := a.(parser.VariableAssignment)
		var val parser.Primary
		var err error
		if assignment.Value == nil {
			val = parser.NewNull()
		} else {
			val, err = filter.Evaluate(assignment.Value)
			if err != nil {
				return err
			}
		}
		err = v.Add(assignment.Name, val)
		if err != nil {
			return err
		}
	}
	return nil
}

func (v Variables) Substitute(substitution parser.VariableSubstitution, filter Filter) (parser.Primary, error) {
	val, err := filter.Evaluate(substitution.Value)
	if err != nil {
		return nil, err
	}
	err = v.Set(substitution.Variable.Name, val)
	if err != nil {
		return nil, err
	}
	return val, nil
}
