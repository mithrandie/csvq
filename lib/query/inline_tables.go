package query

import (
	"context"
	"strings"

	"github.com/mithrandie/csvq/lib/parser"
)

type InlineTableNodes []InlineTableMap

func (list InlineTableNodes) Set(ctx context.Context, inlineTable parser.InlineTable) error {
	return list[0].Set(ctx, inlineTable)
}

func (list InlineTableNodes) Get(name parser.Identifier) (*View, error) {
	for _, m := range list {
		if view, err := m.Get(name); err == nil {
			return view, nil
		}
	}
	return nil, NewUndefinedInLineTableError(name)
}

func (list InlineTableNodes) Load(ctx context.Context, clause parser.WithClause) error {
	for _, v := range clause.InlineTables {
		inlineTable := v.(parser.InlineTable)
		err := list.Set(ctx, inlineTable)
		if err != nil {
			return err
		}
	}

	return nil
}

type InlineTableMap map[string]*View

func (it InlineTableMap) Set(ctx context.Context, inlineTable parser.InlineTable) error {
	uname := strings.ToUpper(inlineTable.Name.Literal)
	if _, err := it.Get(inlineTable.Name); err == nil {
		return NewInLineTableRedefinedError(inlineTable.Name)
	}

	parentFilter, err := GetFilter(ctx)
	if err != nil {
		return err
	}

	filter := parentFilter.CreateNode()
	ctx = ContextForExecusion(ctx, filter)
	defer filter.CloseNode()

	if inlineTable.IsRecursive() {
		filter.recursiveTable = &inlineTable
	}
	view, err := Select(ctx, inlineTable.Query)
	if err != nil {
		return err
	}

	err = view.Header.Update(inlineTable.Name.Literal, inlineTable.Fields)
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

func (it InlineTableMap) Get(name parser.Identifier) (*View, error) {
	uname := strings.ToUpper(name.Literal)
	if view, ok := it[uname]; ok {
		return view.Copy(), nil
	}
	return nil, NewUndefinedInLineTableError(name)
}

func (it InlineTableMap) Clear() {
	for k := range it {
		delete(it, k)
	}
}
