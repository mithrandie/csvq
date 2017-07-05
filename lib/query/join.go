package query

import (
	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/ternary"
)

func ParseJoinCondition(join parser.Join, view *View, joinView *View) parser.Expression {
	if join.Natural.IsEmpty() && join.Condition == nil {
		return nil
	}

	var using []parser.Expression

	if !join.Natural.IsEmpty() {
		for _, f1 := range view.Header {
			if f1.Column == INTERNAL_ID_COLUMN {
				continue
			}

			for _, f2 := range joinView.Header {
				if f2.Column == INTERNAL_ID_COLUMN {
					continue
				}

				if f1.Column == f2.Column {
					using = append(using, parser.Identifier{Literal: f1.Column})
				}
			}
		}
	} else {
		cond := join.Condition.(parser.JoinCondition)
		if cond.On != nil {
			return cond.On
		}

		using = cond.Using
	}

	if len(using) < 1 {
		return nil
	}

	viewName := join.Table.Name()
	joinViewName := join.JoinTable.Name()

	comps := make([]parser.Comparison, len(using))
	for i, v := range using {
		comps[i] = parser.Comparison{
			LHS:      parser.Identifier{Literal: viewName + "." + v.(parser.Identifier).Literal},
			RHS:      parser.Identifier{Literal: joinViewName + "." + v.(parser.Identifier).Literal},
			Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: "="},
		}
	}

	if len(comps) == 1 {
		return comps[0]
	}

	logic := parser.Logic{
		LHS:      comps[0],
		RHS:      comps[1],
		Operator: parser.Token{Token: parser.AND, Literal: parser.TokenLiteral(parser.AND)},
	}
	for i := 2; i < len(comps); i++ {
		logic = parser.Logic{
			LHS:      logic,
			RHS:      comps[i],
			Operator: parser.Token{Token: parser.AND, Literal: parser.TokenLiteral(parser.AND)},
		}
	}
	return logic
}

func CrossJoin(view *View, joinView *View) {
	mergedHeader := MergeHeader(view.Header, joinView.Header)
	records := make([]Record, view.RecordLen()*joinView.RecordLen())

	idx := 0
	for _, viewRecord := range view.Records {
		for _, joinViewRecord := range joinView.Records {
			records[idx] = append(viewRecord, joinViewRecord...)
			idx++
		}
	}

	view.Header = mergedHeader
	view.Records = records
	view.FileInfo = nil
}

func InnerJoin(view *View, joinView *View, condition parser.Expression, parentFilter Filter) error {
	mergedHeader := MergeHeader(view.Header, joinView.Header)

	filter := Filter{
		Records: []FilterRecord{
			{View: view},
			{View: joinView},
		},
		CommonTables: CommonTables{},
	}
	filter = filter.Merge(parentFilter)

	records := []Record{}
	for i, viewRecord := range view.Records {
		for j, joinViewRecord := range joinView.Records {
			filter.Records[0].RecordIndex = i
			filter.Records[1].RecordIndex = j
			primary, err := filter.Evaluate(condition)
			if err != nil {
				return err
			}
			if primary.Ternary() == ternary.TRUE {
				record := append(viewRecord, joinViewRecord...)
				records = append(records, record)
			}
		}
	}

	view.Header = mergedHeader
	view.Records = records
	view.FileInfo = nil
	return nil
}

func OuterJoin(view *View, joinView *View, condition parser.Expression, direction int, parentFilter Filter) error {
	if direction == parser.TOKEN_UNDEFINED {
		direction = parser.LEFT
	}

	mergedHeader := MergeHeader(view.Header, joinView.Header)

	if direction == parser.RIGHT {
		view, joinView = joinView, view
	}

	filter := Filter{
		Records: []FilterRecord{
			{View: view},
			{View: joinView},
		},
		CommonTables: CommonTables{},
	}
	filter = filter.Merge(parentFilter)

	records := []Record{}
	joinViewMatches := make([]bool, len(joinView.Records))
	for i, viewRecord := range view.Records {
		match := false
		for j, joinViewRecord := range joinView.Records {
			filter.Records[0].RecordIndex = i
			filter.Records[1].RecordIndex = j
			primary, err := filter.Evaluate(condition)
			if err != nil {
				return err
			}
			if primary.Ternary() == ternary.TRUE {
				var record Record
				switch direction {
				case parser.RIGHT:
					record = append(joinViewRecord, viewRecord...)
				default:
					record = append(viewRecord, joinViewRecord...)
					if !joinViewMatches[j] {
						joinViewMatches[j] = true
					}
				}
				records = append(records, record)
				match = true
			}
		}

		if !match {
			empty := NewEmptyRecord(joinView.FieldLen())
			var record Record
			switch direction {
			case parser.RIGHT:
				record = append(empty, viewRecord...)
			default:
				record = append(viewRecord, empty...)
			}
			records = append(records, record)

		}
	}

	if direction == parser.FULL {
		for i, match := range joinViewMatches {
			if !match {
				empty := NewEmptyRecord(view.FieldLen())
				record := append(empty, joinView.Records[i]...)
				records = append(records, record)
			}
		}
	}

	if direction == parser.RIGHT {
		view, joinView = joinView, view
	}

	view.Header = mergedHeader
	view.Records = records
	view.FileInfo = nil
	return nil
}
