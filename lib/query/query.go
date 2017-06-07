package query

import (
	"github.com/mithrandie/csvq/lib/parser"
)

type StatementType int

const (
	SELECT StatementType = iota
	INSERT
)

type Result struct {
	Type  StatementType
	View  *View
	Count int
}

func Execute(input string) ([]Result, error) {
	results := []Result{}

	parser.SetDebugLevel(0, true)
	program, err := parser.Parse(input)
	if err != nil {
		return nil, err
	}

	for _, stmt := range program {
		Variable.ClearAutoIncrement()

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
			view, err := Select(stmt.(parser.SelectQuery), nil)
			if err != nil {
				return nil, err
			}
			results = append(results, Result{
				Type:  SELECT,
				View:  view,
				Count: view.RecordLen(),
			})
		case parser.InsertQuery:
			view, err := Insert(stmt.(parser.InsertQuery))
			if err != nil {
				return nil, err
			}
			results = append(results, Result{
				Type:  INSERT,
				View:  view,
				Count: view.OperatedRecords,
			})
		}
	}

	return results, nil
}

func Select(query parser.SelectQuery, parentFilter Filter) (*View, error) {
	if query.FromClause == nil {
		query.FromClause = parser.FromClause{}
	}
	view, err := NewView(query.FromClause.(parser.FromClause), parentFilter)
	if err != nil {
		return nil, err
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

func Insert(query parser.InsertQuery) (*View, error) {
	fromClause := parser.FromClause{
		Tables: []parser.Expression{
			parser.Table{Object: query.Table},
		},
	}
	selectClause := parser.SelectClause{
		Fields: []parser.Expression{
			parser.Field{Object: parser.AllColumns{}},
		},
	}

	view, err := NewView(fromClause, nil)
	if err != nil {
		return nil, err
	}

	fields := query.Fields
	if fields == nil {
		fields = view.Header.TableColumns()
	}

	if query.ValuesList != nil {
		if err := view.InsertValues(fields, query.ValuesList); err != nil {
			return nil, err
		}
	} else {
		if err := view.InsertFromQuery(fields, query.Query.(parser.SelectQuery)); err != nil {
			return nil, err
		}
	}

	view.Select(selectClause)
	view.Fix()

	return view, nil
}
