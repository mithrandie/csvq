package query

import (
	"context"

	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/value"
)

type VariableMap struct {
	*SyncMap
}

func NewVariableMap() VariableMap {
	return VariableMap{
		NewSyncMap(),
	}
}

func (m VariableMap) IsEmpty() bool {
	return m.SyncMap == nil
}

func (m VariableMap) Store(name string, val value.Primary) {
	m.store(name, val)
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

func (m VariableMap) Add(variable parser.Variable, val value.Primary) error {
	if m.Exists(variable.Name) {
		return NewVariableRedeclaredError(variable)
	}
	m.Store(variable.Name, val)
	return nil
}

func (m VariableMap) Set(variable parser.Variable, val value.Primary) error {
	if !m.Exists(variable.Name) {
		return NewUndeclaredVariableError(variable)
	}
	m.Store(variable.Name, val)
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

func (m VariableMap) Declare(ctx context.Context, scope *ReferenceScope, declaration parser.VariableDeclaration) error {
	for _, assignment := range declaration.Assignments {
		var val value.Primary
		var err error
		if assignment.Value == nil {
			val = value.NewNull()
		} else {
			val, err = Evaluate(ctx, scope, assignment.Value)
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

func (m VariableMap) Substitute(ctx context.Context, scope *ReferenceScope, substitution parser.VariableSubstitution) (value.Primary, error) {
	val, err := Evaluate(ctx, scope, substitution.Value)
	if err != nil {
		return nil, err
	}
	err = m.Set(substitution.Variable, val)
	return val, err
}
