package query

import (
	"github.com/mithrandie/csvq/lib/parser"
)

type Statement int

const (
	SELECT Statement = iota
)

type Result struct {
	Statement Statement
	View      *View
	Count     int
}

func Execute(input string) ([]Result, error) {
	results := []Result{}

	parser.SetDebugLevel(0, true)
	program, err := parser.Parse(input)
	if err != nil {
		return nil, err
	}

	for _, stmt := range program {
		switch stmt.(type) {
		case parser.VariableDeclaration:
			if err := Variable.Decrare(stmt.(parser.VariableDeclaration), nil); err != nil {
				return nil, err
			}
		case parser.VariableSubstitution:
			if _, err := Variable.Substitute(stmt.(parser.VariableSubstitution), nil); err != nil {
				return nil, err
			}
		case parser.SelectQuery:
			Variable.ClearAutoIncrement()

			view, err := ExecuteSelect(stmt.(parser.SelectQuery), nil)
			if err != nil {
				return nil, err
			}
			results = append(results, Result{
				Statement: SELECT,
				View:      view,
				Count:     view.RecordLen(),
			})
		}
	}

	return results, nil
}

func ExecuteSelect(query parser.SelectQuery, parentFilter Filter) (*View, error) {
	var view *View

	if query.FromClause == nil {
		view = NewDualView()
	} else {
		v, err := NewView(query.FromClause.(parser.FromClause), parentFilter)
		if err != nil {
			return nil, err
		}
		view = v
	}

	if query.WhereClause != nil {
		if err := view.Where(query.WhereClause.(parser.WhereClause)); err != nil {
			return nil, err
		}
		view.Extract()
	}

	if query.GroupByClause != nil {
		if err := view.GroupBy(query.GroupByClause.(parser.GroupByClause)); err != nil {
			return nil, err
		}
	}

	if query.HavingClause != nil {
		if err := view.Having(query.HavingClause.(parser.HavingClause)); err != nil {
			return nil, err
		}
		view.Extract()
	}

	if err := view.Select(query.SelectClause.(parser.SelectClause)); err != nil {
		return nil, err
	}

	if query.OrderByClause != nil {
		if err := view.OrderBy(query.OrderByClause.(parser.OrderByClause)); err != nil {
			return nil, err
		}
	}

	if query.LimitClause != nil {
		view.Limit(query.LimitClause.(parser.LimitClause))
	}

	view.Fix()

	return view, nil
}
