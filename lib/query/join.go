package query

import (
	"sync"

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
			LHS:      parser.FieldReference{BaseExpr: v.GetBaseExpr(), View: viewName, Column: v.(parser.Identifier)},
			RHS:      parser.FieldReference{BaseExpr: v.GetBaseExpr(), View: joinViewName, Column: v.(parser.Identifier)},
			Operator: "=",
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
			records[idx] = MergeRecord(viewRecord, joinViewRecord)
			idx++
		}
	}

	view.Header = mergedHeader
	view.Records = records
	view.FileInfo = nil
}

func InnerJoin(view *View, joinView *View, condition parser.Expression, parentFilter *Filter) error {
	if condition == nil {
		CrossJoin(view, joinView)
		return nil
	}

	mergedHeader := MergeHeader(view.Header, joinView.Header)

	cpu := NumberOfCPU(view.RecordLen())

	var err error
	recordsList := make([]Records, cpu)

	wg := sync.WaitGroup{}
	for i := 0; i < cpu; i++ {
		wg.Add(1)
		go func(thIdx int) {
			records := Records{}
			start, end := RecordRange(thIdx, view.RecordLen(), cpu)

			filter := NewFilterForRecord(
				&View{
					Header:  mergedHeader,
					Records: make(Records, 1),
				},
				0,
				parentFilter,
			)

		InnerJoinLoop:
			for _, viewRecord := range view.Records[start:end] {
				for _, joinViewRecord := range joinView.Records {
					if err != nil {
						break InnerJoinLoop
					}

					mergedRecord := MergeRecord(viewRecord, joinViewRecord)
					filter.Records[0].View.Records[0] = mergedRecord

					primary, e := filter.Evaluate(condition)
					if e != nil {
						err = e
						break InnerJoinLoop
					}
					if primary.Ternary() == ternary.TRUE {
						records = append(records, mergedRecord)
					}
				}
			}

			recordsList[thIdx] = records
			wg.Done()
		}(i)
	}
	wg.Wait()

	if err != nil {
		return err
	}

	view.Header = mergedHeader
	view.Records = MergeRecordsList(recordsList)
	view.FileInfo = nil
	return nil
}

func OuterJoin(view *View, joinView *View, condition parser.Expression, direction int, parentFilter *Filter) error {
	if direction == parser.TOKEN_UNDEFINED {
		direction = parser.LEFT
	}

	mergedHeader := MergeHeader(view.Header, joinView.Header)

	if direction == parser.RIGHT {
		view, joinView = joinView, view
	}

	viewEmptyRecord := NewEmptyRecord(view.FieldLen())
	joinViewEmptyRecord := NewEmptyRecord(joinView.FieldLen())

	cpu := NumberOfCPU(view.RecordLen())

	var err error
	recordsList := make([]Records, cpu)
	joinViewMatchesList := make([][]bool, cpu)

	wg := sync.WaitGroup{}
	for i := 0; i < cpu; i++ {
		wg.Add(1)
		go func(thIdx int) {
			records := Records{}
			start, end := RecordRange(thIdx, view.RecordLen(), cpu)

			filter := NewFilterForRecord(
				&View{
					Header:  mergedHeader,
					Records: make(Records, 1),
				},
				0,
				parentFilter,
			)

			joinViewMatches := make([]bool, joinView.RecordLen())

		OuterJoinLoop:
			for _, viewRecord := range view.Records[start:end] {
				match := false
				for j, joinViewRecord := range joinView.Records {
					if err != nil {
						break OuterJoinLoop
					}

					var mergedRecord Record
					switch direction {
					case parser.RIGHT:
						mergedRecord = MergeRecord(joinViewRecord, viewRecord)
					default:
						mergedRecord = MergeRecord(viewRecord, joinViewRecord)
					}
					filter.Records[0].View.Records[0] = mergedRecord

					primary, e := filter.Evaluate(condition)
					if e != nil {
						err = e
						break OuterJoinLoop
					}
					if primary.Ternary() == ternary.TRUE {
						if direction == parser.FULL && !joinViewMatches[j] {
							joinViewMatches[j] = true
						}
						records = append(records, mergedRecord)
						match = true
					}
				}

				if !match {
					var record Record
					switch direction {
					case parser.RIGHT:
						record = MergeRecord(joinViewEmptyRecord, viewRecord)
					default:
						record = MergeRecord(viewRecord, joinViewEmptyRecord)
					}
					records = append(records, record)

				}
			}

			recordsList[thIdx] = records
			joinViewMatchesList[thIdx] = joinViewMatches
			wg.Done()
		}(i)
	}
	wg.Wait()

	if err != nil {
		return err
	}

	if direction == parser.FULL {
		for i := 0; i < joinView.RecordLen(); i++ {
			match := false
			for _, joinViewMatches := range joinViewMatchesList {
				if joinViewMatches[i] {
					match = true
					break
				}
			}
			if !match {
				record := MergeRecord(viewEmptyRecord, joinView.Records[i])
				recordsList[len(recordsList)-1] = append(recordsList[len(recordsList)-1], record)
			}
		}
	}

	if direction == parser.RIGHT {
		view, joinView = joinView, view
	}

	view.Header = mergedHeader
	view.Records = MergeRecordsList(recordsList)
	view.FileInfo = nil
	return nil
}
