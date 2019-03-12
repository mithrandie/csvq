package query

import (
	"context"
	"reflect"
	"sort"
	"sync"

	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/value"
)

type VariableScopes []VariableMap

func (list VariableScopes) Declare(ctx context.Context, expr parser.VariableDeclaration, filter *Filter) error {
	return list[0].Declare(ctx, expr, filter)
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

func (list VariableScopes) Substitute(ctx context.Context, expr parser.VariableSubstitution, filter *Filter) (value value.Primary, err error) {
	for _, v := range list {
		if value, err = v.Substitute(ctx, expr, filter); err == nil {
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

func (list VariableScopes) Equal(list2 VariableScopes) bool {
	if len(list) != len(list2) {
		return false
	}
	for i := 0; i < len(list); i++ {
		if !list[i].Equal(&list2[i]) {
			return false
		}
	}
	return true
}

func (list VariableScopes) All() VariableMap {
	all := NewVariableMap()
	for _, m := range list {
		m.variables.Range(func(key, value interface{}) bool {
			if _, ok := all.variables.Load(key); !ok {
				all.variables.Store(key, value)
			}
			return true
		})
	}
	return all
}

type VariableMap struct {
	variables sync.Map
}

func NewVariableMap() VariableMap {
	return VariableMap{
		variables: sync.Map{},
	}
}

func (m *VariableMap) Add(variable parser.Variable, value value.Primary) error {
	if _, ok := m.variables.Load(variable.Name); ok {
		return NewVariableRedeclaredError(variable)
	}
	m.variables.Store(variable.Name, value)
	return nil
}

func (m *VariableMap) Set(variable parser.Variable, value value.Primary) error {
	if _, ok := m.variables.Load(variable.Name); !ok {
		return NewUndeclaredVariableError(variable)
	}
	m.variables.Store(variable.Name, value)
	return nil
}

func (m *VariableMap) Get(variable parser.Variable) (value.Primary, error) {
	if v, ok := m.variables.Load(variable.Name); ok {
		return v.(value.Primary), nil
	}
	return nil, NewUndeclaredVariableError(variable)
}

func (m *VariableMap) SortedKeys() []string {
	keys := make([]string, 0, 10)
	m.variables.Range(func(key, value interface{}) bool {
		keys = append(keys, key.(string))
		return true
	})
	sort.Strings(keys)
	return keys
}

func (m *VariableMap) Dispose(variable parser.Variable) error {
	if _, ok := m.variables.Load(variable.Name); !ok {
		return NewUndeclaredVariableError(variable)
	}
	m.variables.Delete(variable.Name)
	return nil
}

func (m *VariableMap) Declare(ctx context.Context, declaration parser.VariableDeclaration, filter *Filter) error {
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

func (m *VariableMap) Substitute(ctx context.Context, substitution parser.VariableSubstitution, filter *Filter) (value.Primary, error) {
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

func (m *VariableMap) Equal(m2 *VariableMap) bool {
	mvalues := make(map[interface{}]interface{})
	m2values := make(map[interface{}]interface{})

	m.variables.Range(func(key, value interface{}) bool {
		mvalues[key] = value
		return true
	})
	m2.variables.Range(func(key, value interface{}) bool {
		m2values[key] = value
		return true
	})

	return reflect.DeepEqual(mvalues, m2values)
}
