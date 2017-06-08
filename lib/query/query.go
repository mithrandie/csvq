package query

import (
	"errors"

	"github.com/mithrandie/csvq/lib/parser"
)

type StatementType int

const (
	SELECT StatementType = iota
	INSERT
	UPDATE
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
		ViewCache.Clear()

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
		case parser.UpdateQuery:
			views, err := Update(stmt.(parser.UpdateQuery))
			if err != nil {
				return nil, err
			}
			for _, view := range views {
				results = append(results, Result{
					Type:  UPDATE,
					View:  view,
					Count: view.OperatedRecords,
				})
			}
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
	view, err := NewViewFromIdentifier(query.Table, nil)
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

	if err := view.SelectAllColumns(); err != nil {
		return nil, err
	}

	view.Fix()

	return view, nil
}

func Update(query parser.UpdateQuery) ([]*View, error) {
	if query.FromClause == nil {
		query.FromClause = parser.FromClause{Tables: query.Tables}
	}

	view, err := NewView(query.FromClause.(parser.FromClause), nil)
	if err != nil {
		return nil, err
	}

	if query.WhereClause != nil {
		if err := view.Where(query.WhereClause.(parser.WhereClause)); err != nil {
			return nil, err
		}
		view.Extract()
	}

	viewsToUpdate := make(map[string]*View)
	updatedIndices := make(map[string][]int)
	for _, v := range query.Tables {
		table := v.(parser.Table)
		if viewsToUpdate[table.Name()], err = ViewCache.Get(table.Name()); err != nil {
			return nil, err
		}
		updatedIndices[table.Name()] = []int{}
	}

	for i := range view.Records {
		var filter Filter = []FilterRecord{{View: view, RecordIndex: i}}

		for _, v := range query.SetList {
			uset := v.(parser.UpdateSet)

			value, err := filter.Evaluate(uset.Value)
			if err != nil {
				return nil, err
			}

			viewref, err := view.FieldRef(uset.Field)
			if err != nil {
				return nil, err
			}

			internalId, _ := view.InternalRecordId(viewref, i)

			if InIntArray(internalId, updatedIndices[viewref]) {
				return nil, errors.New("record to update is ambiguous")
			}

			fieldIdx, _ := viewsToUpdate[viewref].FieldIndex(uset.Field)

			viewsToUpdate[viewref].Records[internalId][fieldIdx] = NewCell(value)
			updatedIndices[viewref] = append(updatedIndices[viewref], internalId)
		}
	}

	views := []*View{}
	for k, v := range viewsToUpdate {
		if err := v.SelectAllColumns(); err != nil {
			return nil, err
		}

		v.Fix()
		v.OperatedRecords = len(updatedIndices[k])
		views = append(views, v)
	}

	return views, nil
}
