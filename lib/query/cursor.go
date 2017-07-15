package query

import (
	"errors"
	"fmt"

	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/ternary"
	"strings"
)

type CursorMap map[string]*Cursor

func (m CursorMap) Add(key string, query parser.SelectQuery) error {
	uname := strings.ToUpper(key)
	if _, ok := m[uname]; ok {
		return errors.New(fmt.Sprintf("cursor %s already exists", key))
	}
	m[uname] = NewCursor(key, query)
	return nil
}

func (m CursorMap) Dispose(key string) {
	uname := strings.ToUpper(key)
	if cur, ok := m[uname]; ok {
		cur.Close()
		delete(m, uname)
	}
}

func (m CursorMap) Open(key string) error {
	if cur, ok := m[strings.ToUpper(key)]; ok {
		return cur.Open()
	}
	return errors.New(fmt.Sprintf("cursor %s does not exist", key))
}

func (m CursorMap) Close(key string) error {
	if cur, ok := m[strings.ToUpper(key)]; ok {
		cur.Close()
		return nil
	}
	return errors.New(fmt.Sprintf("cursor %s does not exist", key))
}

func (m CursorMap) Fetch(key string, position int, number int) ([]parser.Primary, error) {
	if cur, ok := m[strings.ToUpper(key)]; ok {
		return cur.Fetch(position, number)
	}
	return nil, errors.New(fmt.Sprintf("cursor %s does not exist", key))
}

func (m CursorMap) IsOpen(key string) (ternary.Value, error) {
	if cur, ok := m[strings.ToUpper(key)]; ok {
		return ternary.ParseBool(cur.view != nil), nil
	}
	return ternary.FALSE, errors.New(fmt.Sprintf("cursor %s does not exist", key))
}

func (m CursorMap) IsInRange(key string) (ternary.Value, error) {
	if cur, ok := m[strings.ToUpper(key)]; ok {
		if cur.view == nil {
			return ternary.FALSE, errors.New(fmt.Sprintf("cursor %s is closed", key))
		}
		if !cur.fetched {
			return ternary.UNKNOWN, nil
		}
		return ternary.ParseBool(-1 < cur.index && cur.index < cur.view.RecordLen()), nil
	}
	return ternary.FALSE, errors.New(fmt.Sprintf("cursor %s does not exist", key))
}

type Cursor struct {
	name    string
	query   parser.SelectQuery
	view    *View
	index   int
	fetched bool
}

func NewCursor(name string, query parser.SelectQuery) *Cursor {
	return &Cursor{
		name:  name,
		query: query,
	}
}

func (c *Cursor) Open() error {
	if c.view != nil {
		return errors.New(fmt.Sprintf("cursor %s is already open", c.name))
	}

	view, err := Select(c.query)
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

func (c *Cursor) Fetch(position int, number int) ([]parser.Primary, error) {
	if c.view == nil {
		return nil, errors.New(fmt.Sprintf("cursor %s is closed", c.name))
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
