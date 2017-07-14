package query

import (
	"errors"
	"fmt"
	"strings"

	"github.com/mithrandie/csvq/lib/parser"
)

type InlineTables map[string]*View

func (it InlineTables) Set(inlineTable parser.InlineTable) error {
	uname := strings.ToUpper(inlineTable.Name.Literal)
	if _, err := it.Get(uname); err == nil {
		return errors.New(fmt.Sprintf("inline table %s already exists", inlineTable.Name.Literal))
	}

	filter := Filter{}
	filter.InlineTables = it
	if inlineTable.IsRecursive() {
		filter.RecursiveTable = inlineTable
	}
	view, err := SelectAsSubquery(inlineTable.Query, filter)
	if err != nil {
		return err
	}

	err = view.UpdateHeader(inlineTable.Name.Literal, inlineTable.Columns)
	if err != nil {
		return err
	}

	it[uname] = view
	return nil
}

func (it InlineTables) Get(name string) (*View, error) {
	uname := strings.ToUpper(name)
	if view, ok := it[uname]; ok {
		return view.Copy(), nil
	}
	return nil, errors.New(fmt.Sprintf("inline table %s does not exist", name))
}

func (it InlineTables) Copy() InlineTables {
	table := InlineTables{}
	for k, v := range it {
		table[k] = v
	}
	return table
}

func (it InlineTables) Merge(tables InlineTables) InlineTables {
	table := it.Copy()
	for k, v := range tables {
		table[k] = v
	}
	return table
}

func (it InlineTables) Load(clause parser.WithClause) error {
	for _, v := range clause.InlineTables {
		inlineTable := v.(parser.InlineTable)
		err := it.Set(inlineTable)
		if err != nil {
			return err
		}

		view, _ := it.Get(inlineTable.Name.Literal)
		err = view.UpdateHeader(inlineTable.Name.Literal, inlineTable.Columns)
		if err != nil {
			return err
		}
	}

	return nil
}
