package query

import (
	"context"
	"sort"
	"sync"

	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/value"
)

type VariableScopes []VariableMap

func (list VariableScopes) Declare(ctx context.Context, filter *Filter, expr parser.VariableDeclaration) error {
	return list[0].Declare(ctx, filter, expr)
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

func (list VariableScopes) Substitute(ctx context.Context, filter *Filter, expr parser.VariableSubstitution) (value value.Primary, err error) {
	for i := range list {
		if value, err = list[i].Substitute(ctx, filter, expr); err == nil {
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
		for k, v := range list[i].variables {
			if _, ok := all.variables[k]; !ok {
				all.variables[k] = v
			}
		}
	}
	return all
}

type VariableMap struct {
	mtx       *sync.RWMutex
	variables map[string]value.Primary
}

func NewVariableMap() VariableMap {
	return VariableMap{
		mtx:       &sync.RWMutex{},
		variables: make(map[string]value.Primary),
	}
}

func (m *VariableMap) Add(variable parser.Variable, value value.Primary) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	if _, ok := m.variables[variable.Name]; ok {
		return NewVariableRedeclaredError(variable)
	}
	m.variables[variable.Name] = value
	return nil
}

func (m *VariableMap) Set(variable parser.Variable, value value.Primary) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	if _, ok := m.variables[variable.Name]; !ok {
		return NewUndeclaredVariableError(variable)
	}
	m.variables[variable.Name] = value
	return nil
}

func (m *VariableMap) Get(variable parser.Variable) (value.Primary, error) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	if v, ok := m.variables[variable.Name]; ok {
		return v, nil
	}
	return nil, NewUndeclaredVariableError(variable)
}

func (m *VariableMap) SortedKeys() []string {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	keys := make([]string, 0, len(m.variables))
	for k := range m.variables {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (m *VariableMap) Dispose(variable parser.Variable) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	if _, ok := m.variables[variable.Name]; !ok {
		return NewUndeclaredVariableError(variable)
	}
	delete(m.variables, variable.Name)
	return nil
}

func (m *VariableMap) Declare(ctx context.Context, filter *Filter, declaration parser.VariableDeclaration) error {
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

func (m *VariableMap) Substitute(ctx context.Context, filter *Filter, substitution parser.VariableSubstitution) (value.Primary, error) {
	val, err := filter.Evaluate(ctx, substitution.Value)
	if err != nil {
		return nil, err
	}
	return m.SubstituteDirectly(substitution.Variable, val)
}

func (m *VariableMap) SubstituteDirectly(variable parser.Variable, value value.Primary) (value.Primary, error) {
	err := m.Set(variable, value)
	return value, err
}
