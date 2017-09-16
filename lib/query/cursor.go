package query

import (
	"strings"

	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/value"

	"fmt"
	"github.com/mithrandie/ternary"
	"sort"
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

func (list CursorScopes) Open(name parser.Identifier, filter *Filter) error {
	var err error

	for _, m := range list {
		err = m.Open(name, filter)
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

func (list CursorScopes) List() []string {
	cursors := make(map[string]string)
	keys := make([]string, 0)

	for _, m := range list {
		for _, cursor := range m {
			if cursor.isPseudo {
				continue
			}
			if !InStrSliceWithCaseInsensitive(cursor.name, keys) {
				cursors[cursor.name] = cursor.query.String()
				keys = append(keys, cursor.name)
			}
		}
	}
	sort.Strings(keys)

	dcls := make([]string, len(keys))
	for i, key := range keys {
		dcls[i] = fmt.Sprintf("%s for %s", key, cursors[key])
	}

	return dcls
}

type CursorMap map[string]*Cursor

func (m CursorMap) Declare(expr parser.CursorDeclaration) error {
	uname := strings.ToUpper(expr.Cursor.Literal)
	if _, ok := m[uname]; ok {
		return NewCursorRedeclaredError(expr.Cursor)
	}
	m[uname] = NewCursor(expr.Cursor.Literal, expr.Query)
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

func (m CursorMap) Open(name parser.Identifier, filter *Filter) error {
	if cur, ok := m[strings.ToUpper(name.Literal)]; ok {
		return cur.Open(name, filter)
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
		return ternary.ConvertFromBool(cur.view != nil), nil
	}
	return ternary.FALSE, NewUndeclaredCursorError(name)
}

func (m CursorMap) IsInRange(name parser.Identifier) (ternary.Value, error) {
	if cur, ok := m[strings.ToUpper(name.Literal)]; ok {
		if cur.view == nil {
			return ternary.FALSE, NewCursorClosedError(name)
		}
		if !cur.fetched {
			return ternary.UNKNOWN, nil
		}
		return ternary.ConvertFromBool(-1 < cur.index && cur.index < cur.view.RecordLen()), nil
	}
	return ternary.FALSE, NewUndeclaredCursorError(name)
}

func (m CursorMap) Count(name parser.Identifier) (int, error) {
	if cur, ok := m[strings.ToUpper(name.Literal)]; ok {
		if cur.view == nil {
			return 0, NewCursorClosedError(name)
		}
		return cur.view.RecordLen(), nil
	}
	return 0, NewUndeclaredCursorError(name)
}

type Cursor struct {
	name    string
	query   parser.SelectQuery
	view    *View
	index   int
	fetched bool

	isPseudo bool
}

func NewCursor(name string, query parser.SelectQuery) *Cursor {
	return &Cursor{
		name:  name,
		query: query,
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
	}
}

func (c *Cursor) Open(name parser.Identifier, filter *Filter) error {
	if c.isPseudo {
		return NewPseudoCursorError(name)
	}

	if c.view != nil {
		return NewCursorOpenError(name)
	}

	view, err := Select(c.query, filter)
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

	c.view = nil
	c.index = 0
	c.fetched = false

	return nil
}

func (c *Cursor) Fetch(name parser.Identifier, position int, number int) ([]value.Primary, error) {
	if c.view == nil {
		return nil, NewCursorClosedError(name)
	}

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
