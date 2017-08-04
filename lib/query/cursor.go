package query

import (
	"strings"

	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/ternary"
)

type CursorMapList []CursorMap

func (list CursorMapList) Declare(expr parser.CursorDeclaration) error {
	return list[0].Declare(expr)
}

func (list CursorMapList) AddPseudoCursor(name parser.Identifier, values []parser.Primary) error {
	return list[0].AddPseudoCursor(name, values)
}

func (list CursorMapList) Dispose(name parser.Identifier) error {
	for _, m := range list {
		err := m.Dispose(name)
		if err == nil {
			return nil
		}
		if _, ok := err.(*UndefinedCursorError); !ok {
			return err
		}
	}
	return NewUndefinedCursorError(name)
}

func (list CursorMapList) Open(name parser.Identifier, filter Filter) error {
	var err error

	for _, m := range list {
		err = m.Open(name, filter)
		if err == nil {
			return nil
		}
		if _, ok := err.(*UndefinedCursorError); !ok {
			return err
		}
	}
	return NewUndefinedCursorError(name)
}

func (list CursorMapList) Close(name parser.Identifier) error {
	for _, m := range list {
		err := m.Close(name)
		if err == nil {
			return nil
		}
		if _, ok := err.(*UndefinedCursorError); !ok {
			return err
		}
	}
	return NewUndefinedCursorError(name)
}

func (list CursorMapList) Fetch(name parser.Identifier, position int, number int) ([]parser.Primary, error) {
	var values []parser.Primary
	var err error

	for _, m := range list {
		values, err = m.Fetch(name, position, number)
		if err == nil {
			return values, nil
		}
		if _, ok := err.(*UndefinedCursorError); !ok {
			return nil, err
		}
	}
	return nil, NewUndefinedCursorError(name)
}

func (list CursorMapList) IsOpen(name parser.Identifier) (ternary.Value, error) {
	for _, m := range list {
		if ok, err := m.IsOpen(name); err == nil {
			return ok, nil
		}
	}
	return ternary.FALSE, NewUndefinedCursorError(name)
}

func (list CursorMapList) IsInRange(name parser.Identifier) (ternary.Value, error) {
	var result ternary.Value
	var err error

	for _, m := range list {
		result, err = m.IsInRange(name)
		if err == nil {
			return result, nil
		}
		if _, ok := err.(*UndefinedCursorError); !ok {
			return ternary.FALSE, err
		}
	}
	return ternary.FALSE, NewUndefinedCursorError(name)
}

func (list CursorMapList) Count(name parser.Identifier) (int, error) {
	var count int
	var err error

	for _, m := range list {
		count, err = m.Count(name)
		if err == nil {
			return count, nil
		}
		if _, ok := err.(*UndefinedCursorError); !ok {
			return 0, err
		}
	}
	return 0, NewUndefinedCursorError(name)
}

type CursorMap map[string]*Cursor

func (m CursorMap) Declare(expr parser.CursorDeclaration) error {
	uname := strings.ToUpper(expr.Cursor.Literal)
	if _, ok := m[uname]; ok {
		return NewCursorRedeclaredError(expr.Cursor)
	}
	m[uname] = NewCursor(expr.Query)
	return nil
}

func (m CursorMap) AddPseudoCursor(name parser.Identifier, values []parser.Primary) error {
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
		return cur.Close(name)
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

func (m CursorMap) Count(name parser.Identifier) (int, error) {
	if cur, ok := m[strings.ToUpper(name.Literal)]; ok {
		if cur.view == nil {
			return 0, NewCursorClosedError(name)
		}
		return cur.view.RecordLen(), nil
	}
	return 0, NewUndefinedCursorError(name)
}

type Cursor struct {
	query   parser.SelectQuery
	view    *View
	index   int
	fetched bool

	isPseudo bool
}

func NewCursor(query parser.SelectQuery) *Cursor {
	return &Cursor{
		query: query,
	}
}

func NewPseudoCursor(values []parser.Primary) *Cursor {
	header := NewHeader("", []string{"c1"})

	records := make(Records, len(values))
	for i, v := range values {
		records[i] = NewRecord([]parser.Primary{v})
	}
	view := NewView()
	view.Header = header
	view.Records = records

	return &Cursor{
		view:     view,
		index:    -1,
		fetched:  false,
		isPseudo: true,
	}
}

func (c *Cursor) Open(name parser.Identifier, filter Filter) error {
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
