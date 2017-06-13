package query

import (
	"errors"
	"fmt"

	"github.com/mithrandie/csvq/lib/parser"
)

type CursorMap map[string]*Cursor

func (m CursorMap) Add(key string, query parser.SelectQuery) error {
	if _, ok := m[key]; ok {
		return errors.New(fmt.Sprintf("cursor %s already exists", key))
	}
	m[key] = NewCursor(key, query)
	return nil
}

func (m CursorMap) Dispose(key string) {
	if cur, ok := m[key]; ok {
		cur.Close()
		delete(m, key)
	}
}

func (m CursorMap) Open(key string) error {
	if cur, ok := m[key]; ok {
		return cur.Open()
	}
	return errors.New(fmt.Sprintf("cursor %s does not exist", key))
}

func (m CursorMap) Close(key string) error {
	if cur, ok := m[key]; ok {
		cur.Close()
		return nil
	}
	return errors.New(fmt.Sprintf("cursor %s does not exist", key))
}

func (m CursorMap) Fetch(key string) ([]parser.Primary, error) {
	if cur, ok := m[key]; ok {
		return cur.Fetch()
	}
	return nil, errors.New(fmt.Sprintf("cursor %s does not exist", key))
}

type Cursor struct {
	name  string
	query parser.SelectQuery
	view  *View
	index int
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

	view, err := Select(c.query, nil)
	if err != nil {
		return err
	}

	c.view = view
	c.index = 0
	return nil
}

func (c *Cursor) Close() {
	c.view = nil
	c.index = 0
}

func (c *Cursor) Fetch() ([]parser.Primary, error) {
	if c.view == nil {
		return nil, errors.New(fmt.Sprintf("cursor %s is closed", c.name))
	}

	if c.view.RecordLen() <= c.index {
		return nil, nil
	}

	list := make([]parser.Primary, len(c.view.Records[c.index]))
	for i, cell := range c.view.Records[c.index] {
		list[i] = cell.Primary()
	}

	c.index++
	return list, nil
}
