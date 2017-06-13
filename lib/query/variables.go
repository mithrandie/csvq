package query

import (
	"errors"
	"fmt"

	"github.com/mithrandie/csvq/lib/parser"
)

const (
	AUTO_INCREMENT_KEY = "__auto_increment"
)

type Variables map[string]parser.Primary

func (v Variables) Add(key string, value parser.Primary) error {
	if _, ok := v[key]; ok {
		return errors.New(fmt.Sprintf("variable %s is redeclared", key))
	}
	v[key] = value
	return nil
}

func (v Variables) Set(key string, value parser.Primary) error {
	if _, ok := v[key]; !ok {
		return errors.New(fmt.Sprintf("variable %s is undefined", key))
	}
	v[key] = value
	return nil
}

func (v Variables) Get(key string) (parser.Primary, error) {
	if v, ok := v[key]; ok {
		return v, nil
	}
	return nil, errors.New(fmt.Sprintf("variable %s is undefined", key))
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

func (v Variables) Increment(key string, initialVal int64) parser.Primary {
	if val, err := v.Get(key); err == nil {
		val = Calculate(val, parser.NewInteger(1), '+')
		v.Set(key, val)
		return val
	}

	val := parser.NewInteger(initialVal)
	v.Add(key, val)
	return val
}

func (v Variables) ClearAutoIncrement() {
	v.Delete(AUTO_INCREMENT_KEY)
}
