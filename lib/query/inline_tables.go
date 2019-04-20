package query

import (
	"context"
	"strings"

	"github.com/mithrandie/csvq/lib/parser"
)

type InlineTableMap map[string]*View

func (it InlineTableMap) Set(ctx context.Context, scope *ReferenceScope, inlineTable parser.InlineTable) error {
	uname := strings.ToUpper(inlineTable.Name.Literal)
	if _, err := it.Get(inlineTable.Name); err == nil {
		return NewInLineTableRedefinedError(inlineTable.Name)
	}

	if inlineTable.IsRecursive() {
		scope.RecursiveTable = &inlineTable
	}

	view, err := Select(ctx, scope, inlineTable.Query)
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
