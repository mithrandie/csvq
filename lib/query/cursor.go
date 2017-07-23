package query

import (
	"strings"

	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/ternary"
)

type CursorMap map[string]*Cursor

func (m CursorMap) Declare(expr parser.CursorDeclaration) error {
	uname := strings.ToUpper(expr.Cursor.Literal)
	if _, ok := m[uname]; ok {
		return NewCursorRedeclaredError(expr.Cursor)
	}
	m[uname] = NewCursor(expr.Query)
	return nil
}

func (m CursorMap) Dispose(name parser.Identifier) error {
	uname := strings.ToUpper(name.Literal)
	if _, ok := m[uname]; ok {
		delete(m, uname)
		return nil
	}
	return NewUndefinedCursorError(name)
}

func (m CursorMap) Open(name parser.Identifier, filter Filter) error {
	if cur, ok := m[strings.ToUpper(name.Literal)]; ok {
		return cur.Open(name, filter)
	}
	return NewUndefinedCursorError(name)
}

func (m CursorMap) Close(name parser.Identifier) error {
	if cur, ok := m[strings.ToUpper(name.Literal)]; ok {
		cur.Close()
		return nil
	}
	return NewUndefinedCursorError(name)
}

func (m CursorMap) Fetch(name parser.Identifier, position int, number int) ([]parser.Primary, error) {
	if cur, ok := m[strings.ToUpper(name.Literal)]; ok {
		return cur.Fetch(name, position, number)
	}
	return nil, NewUndefinedCursorError(name)
}

func (m CursorMap) IsOpen(name parser.Identifier) (ternary.Value, error) {
	if cur, ok := m[strings.ToUpper(name.Literal)]; ok {
		return ternary.ParseBool(cur.view != nil), nil
	}
	return ternary.FALSE, NewUndefinedCursorError(name)
}

func (m CursorMap) IsInRange(name parser.Identifier) (ternary.Value, error) {
	if cur, ok := m[strings.ToUpper(name.Literal)]; ok {
		if cur.view == nil {
			return ternary.FALSE, NewCursorClosedError(name)
		}
		if !cur.fetched {
			return ternary.UNKNOWN, nil
		}
		return ternary.ParseBool(-1 < cur.index && cur.index < cur.view.RecordLen()), nil
	}
	return ternary.FALSE, NewUndefinedCursorError(name)
}

type Cursor struct {
	query   parser.SelectQuery
	view    *View
	index   int
	fetched bool
}

func NewCursor(query parser.SelectQuery) *Cursor {
	return &Cursor{
		query: query,
	}
}

func (c *Cursor) Open(name parser.Identifier, filter Filter) error {
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

func (c *Cursor) Close() {
	c.view = nil
	c.index = 0
	c.fetched = false
}

func (c *Cursor) Fetch(name parser.Identifier, position int, number int) ([]parser.Primary, error) {
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

	list := make([]parser.Primary, len(c.view.Records[c.index]))
	for i, cell := range c.view.Records[c.index] {
		list[i] = cell.Primary()
	}

	return list, nil
}
