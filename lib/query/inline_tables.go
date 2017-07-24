package query

import (
	"strings"

	"github.com/mithrandie/csvq/lib/parser"
)

type InlineTables map[string]*View

func (it InlineTables) Set(inlineTable parser.InlineTable, parentFilter Filter) error {
	uname := strings.ToUpper(inlineTable.Name.Literal)
	if _, err := it.Get(inlineTable.Name); err == nil {
		return NewInLineTableRedeclaredError(inlineTable.Name)
	}

	filter := parentFilter.CreateNode()
	filter.InlineTables = it
	if inlineTable.IsRecursive() {
		filter.RecursiveTable = &inlineTable
	}
	view, err := Select(inlineTable.Query, filter)
	if err != nil {
		return err
	}

	err = view.UpdateHeader(inlineTable.Name.Literal, inlineTable.Fields)
	if err != nil {
		if _, ok := err.(*FieldLengthNotMatchError); ok {
			return NewInlineTableFieldLengthError(inlineTable.Query, inlineTable.Name, len(inlineTable.Fields))
		}
		return err
	}

	view.FileInfo = nil
	it[uname] = view

	return nil
}

func (it InlineTables) Get(name parser.Identifier) (*View, error) {
	uname := strings.ToUpper(name.Literal)
	if view, ok := it[uname]; ok {
		return view.Copy(), nil
	}
	return nil, NewUndefinedInLineTableError(name)
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

func (it InlineTables) Load(clause parser.WithClause, parentFilter Filter) error {
	for _, v := range clause.InlineTables {
		inlineTable := v.(parser.InlineTable)
		err := it.Set(inlineTable, parentFilter)
		if err != nil {
			return err
		}
	}

	return nil
}
