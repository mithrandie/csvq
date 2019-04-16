package query

import (
	"context"
	"errors"
	"sort"
	"strings"
	"sync"

	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/value"

	"github.com/mithrandie/ternary"
)

type CursorScopes []CursorMap

func (list CursorScopes) Declare(expr parser.CursorDeclaration) error {
	return list[0].Declare(expr)
}

func (list CursorScopes) AddPseudoCursor(name parser.Identifier, values []value.Primary) error {
	return list[0].AddPseudoCursor(name, values)
}

func (list CursorScopes) Dispose(name parser.Identifier) error {
	for _, m := range list {
		err := m.Dispose(name)
		if err == nil {
			return nil
		}
		if _, ok := err.(*UndeclaredCursorError); !ok {
			return err
		}
	}
	return NewUndeclaredCursorError(name)
}

func (list CursorScopes) Open(ctx context.Context, name parser.Identifier, values []parser.ReplaceValue) error {
	var err error

	for _, m := range list {
		err = m.Open(ctx, name, values)
		if err == nil {
			return nil
		}
		if _, ok := err.(*UndeclaredCursorError); !ok {
			return err
		}
	}
	return NewUndeclaredCursorError(name)
}

func (list CursorScopes) Close(name parser.Identifier) error {
	for _, m := range list {
		err := m.Close(name)
		if err == nil {
			return nil
		}
		if _, ok := err.(*UndeclaredCursorError); !ok {
			return err
		}
	}
	return NewUndeclaredCursorError(name)
}

func (list CursorScopes) Fetch(name parser.Identifier, position int, number int) ([]value.Primary, error) {
	var values []value.Primary
	var err error

	for _, m := range list {
		values, err = m.Fetch(name, position, number)
		if err == nil {
			return values, nil
		}
		if _, ok := err.(*UndeclaredCursorError); !ok {
			return nil, err
		}
	}
	return nil, NewUndeclaredCursorError(name)
}

func (list CursorScopes) IsOpen(name parser.Identifier) (ternary.Value, error) {
	for _, m := range list {
		if ok, err := m.IsOpen(name); err == nil {
			return ok, nil
		}
	}
	return ternary.FALSE, NewUndeclaredCursorError(name)
}

func (list CursorScopes) IsInRange(name parser.Identifier) (ternary.Value, error) {
	var result ternary.Value
	var err error

	for _, m := range list {
		result, err = m.IsInRange(name)
		if err == nil {
			return result, nil
		}
		if _, ok := err.(*UndeclaredCursorError); !ok {
			return ternary.FALSE, err
		}
	}
	return ternary.FALSE, NewUndeclaredCursorError(name)
}

func (list CursorScopes) Count(name parser.Identifier) (int, error) {
	var count int
	var err error

	for _, m := range list {
		count, err = m.Count(name)
		if err == nil {
			return count, nil
		}
		if _, ok := err.(*UndeclaredCursorError); !ok {
			return 0, err
		}
	}
	return 0, NewUndeclaredCursorError(name)
}

func (list CursorScopes) All() CursorMap {
	all := make(CursorMap, 10)

	for _, m := range list {
		for key, cursor := range m {
			if cursor.isPseudo {
				continue
			}
			if _, ok := all[key]; !ok {
				all[key] = cursor
			}
		}
	}
	return all
}

type CursorMap map[string]*Cursor

func (m CursorMap) Declare(expr parser.CursorDeclaration) error {
	uname := strings.ToUpper(expr.Cursor.Literal)
	if _, ok := m[uname]; ok {
		return NewCursorRedeclaredError(expr.Cursor)
	}
	m[uname] = NewCursor(expr)
	return nil
}

func (m CursorMap) AddPseudoCursor(name parser.Identifier, values []value.Primary) error {
	uname := strings.ToUpper(name.Literal)
	if _, ok := m[uname]; ok {
		return NewCursorRedeclaredError(name)
	}
	m[uname] = NewPseudoCursor(values)
	return nil
}

func (m CursorMap) Dispose(name parser.Identifier) error {
	uname := strings.ToUpper(name.Literal)
	if cur, ok := m[uname]; ok {
		if cur.isPseudo {
			return NewPseudoCursorError(name)
		}
		delete(m, uname)
		return nil
	}
	return NewUndeclaredCursorError(name)
}

func (m CursorMap) Open(ctx context.Context, name parser.Identifier, values []parser.ReplaceValue) error {
	if cur, ok := m[strings.ToUpper(name.Literal)]; ok {
		return cur.Open(ctx, name, values)
	}
	return NewUndeclaredCursorError(name)
}

func (m CursorMap) Close(name parser.Identifier) error {
	if cur, ok := m[strings.ToUpper(name.Literal)]; ok {
		return cur.Close(name)
	}
	return NewUndeclaredCursorError(name)
}

func (m CursorMap) Fetch(name parser.Identifier, position int, number int) ([]value.Primary, error) {
	if cur, ok := m[strings.ToUpper(name.Literal)]; ok {
		return cur.Fetch(name, position, number)
	}
	return nil, NewUndeclaredCursorError(name)
}

func (m CursorMap) IsOpen(name parser.Identifier) (ternary.Value, error) {
	if cur, ok := m[strings.ToUpper(name.Literal)]; ok {
		return cur.IsOpen(), nil
	}
	return ternary.FALSE, NewUndeclaredCursorError(name)
}

func (m CursorMap) IsInRange(name parser.Identifier) (ternary.Value, error) {
	if cur, ok := m[strings.ToUpper(name.Literal)]; ok {
		t, err := cur.IsInRange()
		if err != nil {
			return ternary.FALSE, NewCursorClosedError(name)
		}
		return t, nil
	}
	return ternary.FALSE, NewUndeclaredCursorError(name)
}

func (m CursorMap) Count(name parser.Identifier) (int, error) {
	if cur, ok := m[strings.ToUpper(name.Literal)]; ok {
		i, err := cur.Count()
		if err != nil {
			return 0, NewCursorClosedError(name)
		}
		return i, nil
	}
	return 0, NewUndeclaredCursorError(name)
}

func (m CursorMap) Keys() []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func (m CursorMap) SortedKeys() []string {
	keys := m.Keys()
	sort.Strings(keys)
	return keys
}

func (m CursorMap) Clear() {
	for k := range m {
		delete(m, k)
	}
}

type Cursor struct {
	name      string
	query     parser.SelectQuery
	statement parser.Identifier
	view      *View
	index     int
	fetched   bool

	isPseudo bool

	mtx *sync.Mutex
}

func NewCursor(e parser.CursorDeclaration) *Cursor {
	return &Cursor{
		name:      e.Cursor.Literal,
		query:     e.Query,
		statement: e.Statement,
		mtx:       &sync.Mutex{},
	}
}

func NewPseudoCursor(values []value.Primary) *Cursor {
	header := NewHeader("", []string{"c1"})

	records := make(RecordSet, len(values))
	for i, v := range values {
		records[i] = NewRecord([]value.Primary{v})
	}
	view := NewView()
	view.Header = header
	view.RecordSet = records

	return &Cursor{
		view:     view,
		index:    -1,
		fetched:  false,
		isPseudo: true,
		mtx:      &sync.Mutex{},
	}
}

func (c *Cursor) Open(ctx context.Context, name parser.Identifier, values []parser.ReplaceValue) error {
	if c.isPseudo {
		return NewPseudoCursorError(name)
	}

	filter, err := GetFilter(ctx)
	if err != nil {
		return err
	}

	c.mtx.Lock()
	defer c.mtx.Unlock()

	if c.view != nil {
		return NewCursorOpenError(name)
	}

	var view *View
	if c.query.SelectEntity != nil {
		view, err = Select(ctx, c.query)
	} else {
		prepared, e := filter.tx.PreparedStatements.Get(c.statement)
		if e != nil {
			return e
		}
		if len(prepared.Statements) != 1 {
			return NewInvalidCursorStatementError(c.statement)
		}
		stmt, ok := prepared.Statements[0].(parser.SelectQuery)
		if !ok {
			return NewInvalidCursorStatementError(c.statement)
		}
		view, err = Select(ContextForPreparedStatement(ctx, NewReplaceValues(values)), stmt)
	}
	if err != nil {
		return err
	}

	c.view = view
	c.index = -1
	c.fetched = false
	return nil
}

func (c *Cursor) Close(name parser.Identifier) error {
	if c.isPseudo {
		return NewPseudoCursorError(name)
	}

	c.mtx.Lock()

	c.view = nil
	c.index = 0
	c.fetched = false

	c.mtx.Unlock()
	return nil
}

func (c *Cursor) Fetch(name parser.Identifier, position int, number int) ([]value.Primary, error) {
	if c.view == nil {
		return nil, NewCursorClosedError(name)
	}

	c.mtx.Lock()
	defer c.mtx.Unlock()

	if !c.fetched {
		c.fetched = true
	}

	switch position {
	case parser.ABSOLUTE:
		c.index = number
	case parser.RELATIVE:
		c.index = c.index + number
	case parser.FIRST:
		c.index = 0
	case parser.LAST:
		c.index = c.view.RecordLen() - 1
	case parser.PRIOR:
		c.index = c.index - 1
	default: // NEXT
		c.index = c.index + 1
	}

	if c.index < 0 {
		c.index = -1
		return nil, nil
	}

	if c.view.RecordLen() <= c.index {
		c.index = c.view.RecordLen()
		return nil, nil
	}

	list := make([]value.Primary, len(c.view.RecordSet[c.index]))
	for i, cell := range c.view.RecordSet[c.index] {
		list[i] = cell.Value()
	}

	return list, nil
}

func (c *Cursor) IsOpen() ternary.Value {
	return ternary.ConvertFromBool(c.view != nil)
}

func (c *Cursor) IsInRange() (ternary.Value, error) {
	if c.view == nil {
		return ternary.FALSE, errors.New("cursor is closed")
	}
	if !c.fetched {
		return ternary.UNKNOWN, nil
	}
	return ternary.ConvertFromBool(-1 < c.index && c.index < c.view.RecordLen()), nil
}

func (c *Cursor) Count() (int, error) {
	if c.view == nil {
		return 0, errors.New("cursor is closed")
	}
	return c.view.RecordLen(), nil
}

func (c *Cursor) Pointer() (int, error) {
	return c.index, nil
}
