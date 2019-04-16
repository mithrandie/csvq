package query

import (
	"context"

	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/value"
)

type VariableScopes []VariableMap

func (list VariableScopes) Declare(ctx context.Context, expr parser.VariableDeclaration) error {
	return list[0].Declare(ctx, expr)
}

func (list VariableScopes) Get(expr parser.Variable) (value value.Primary, err error) {
	for i := range list {
		if value, err = list[i].Get(expr); err == nil {
			return
		}
	}
	err = NewUndeclaredVariableError(expr)
	return
}

func (list VariableScopes) Substitute(ctx context.Context, expr parser.VariableSubstitution) (value value.Primary, err error) {
	for i := range list {
		if value, err = list[i].Substitute(ctx, expr); err == nil {
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
	for i := range list {
		if value, err = list[i].SubstituteDirectly(variable, value); err == nil {
			return value, nil
		}
	}
	return nil, NewUndeclaredVariableError(variable)
}

func (list VariableScopes) Dispose(expr parser.Variable) error {
	for i := range list {
		if err := list[i].Dispose(expr); err == nil {
			return nil
		}
	}
	return NewUndeclaredVariableError(expr)
}

func (list VariableScopes) All() VariableMap {
	all := NewVariableMap()
	for i := range list {
		list[i].Range(func(key, val interface{}) bool {
			if !all.Exists(key.(string)) {
				all.Store(key.(string), val.(value.Primary))
			}
			return true
		})
	}
	return all
}

type VariableMap struct {
	*SyncMap
}

func NewVariableMap() VariableMap {
	return VariableMap{
		NewSyncMap(),
	}
}

func (m VariableMap) Store(name string, value value.Primary) {
	m.store(name, value)
}

func (m VariableMap) Load(name string) (value.Primary, bool) {
	if v, ok := m.load(name); ok {
		return v.(value.Primary), true
	}
	return nil, false
}

func (m VariableMap) Delete(name string) {
	m.delete(name)
}

func (m VariableMap) Exists(name string) bool {
	return m.exists(name)
}

func (m VariableMap) Add(variable parser.Variable, value value.Primary) error {
	if m.Exists(variable.Name) {
		return NewVariableRedeclaredError(variable)
	}
	m.Store(variable.Name, value)
	return nil
}

func (m VariableMap) Set(variable parser.Variable, value value.Primary) error {
	if !m.Exists(variable.Name) {
		return NewUndeclaredVariableError(variable)
	}
	m.Store(variable.Name, value)
	return nil
}

func (m VariableMap) Get(variable parser.Variable) (value.Primary, error) {
	if v, ok := m.Load(variable.Name); ok {
		return v, nil
	}
	return nil, NewUndeclaredVariableError(variable)
}

func (m VariableMap) Dispose(variable parser.Variable) error {
	if !m.Exists(variable.Name) {
		return NewUndeclaredVariableError(variable)
	}
	m.Delete(variable.Name)
	return nil
}

func (m VariableMap) Declare(ctx context.Context, declaration parser.VariableDeclaration) error {
	filter, err := GetFilter(ctx)
	if err != nil {
		return err
	}

	for _, assignment := range declaration.Assignments {
		var val value.Primary
		var err error
		if assignment.Value == nil {
			val = value.NewNull()
		} else {
			val, err = filter.Evaluate(ctx, assignment.Value)
			if err != nil {
				return err
			}
		}
		err = m.Add(assignment.Variable, val)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m VariableMap) Substitute(ctx context.Context, substitution parser.VariableSubstitution) (value.Primary, error) {
	filter, err := GetFilter(ctx)
	if err != nil {
		return nil, err
	}

	val, err := filter.Evaluate(ctx, substitution.Value)
	if err != nil {
		return nil, err
	}
	return m.SubstituteDirectly(substitution.Variable, val)
}

func (m VariableMap) SubstituteDirectly(variable parser.Variable, value value.Primary) (value.Primary, error) {
	err := m.Set(variable, value)
	return value, err
}
