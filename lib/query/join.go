package query

import (
	"context"
	"math"

	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/ternary"
)

func ParseJoinCondition(join parser.Join, view *View, joinView *View) (parser.QueryExpression, []parser.FieldReference, []parser.FieldReference, error) {
	if join.Natural.IsEmpty() && join.Condition == nil {
		return nil, nil, nil, nil
	}

	var using []parser.QueryExpression

	if !join.Natural.IsEmpty() {
		for _, field := range view.Header {
			if field.Column == InternalIdColumn {
				continue
			}
			ref := parser.FieldReference{BaseExpr: parser.NewBaseExpr(join.Natural), Column: parser.Identifier{Literal: field.Column}}
			if _, err := joinView.FieldIndex(ref); err != nil {
				if _, ok := err.(*FieldAmbiguousError); ok {
					return nil, nil, nil, err
				}
				continue
			}
			using = append(using, parser.Identifier{BaseExpr: parser.NewBaseExpr(join.Natural), Literal: field.Column})
		}
	} else {
		cond := join.Condition.(parser.JoinCondition)
		if cond.On != nil {
			return cond.On, nil, nil, nil
		}

		using = cond.Using
	}

	if len(using) < 1 {
		return nil, nil, nil, nil
	}

	usingFields := make([]string, len(using))
	for i, v := range using {
		usingFields[i] = v.(parser.Identifier).Literal
	}

	includeFields := make([]parser.FieldReference, len(using))
	excludeFields := make([]parser.FieldReference, len(using))

	comps := make([]parser.Comparison, len(using))
	for i, v := range using {
		var lhs parser.FieldReference
		var rhs parser.FieldReference
		fieldref := parser.FieldReference{BaseExpr: v.GetBaseExpr(), Column: v.(parser.Identifier)}

		lhsidx, err := view.FieldIndex(fieldref)
		if err != nil {
			return nil, nil, nil, err
		}
		lhs = parser.FieldReference{BaseExpr: v.GetBaseExpr(), View: parser.Identifier{Literal: view.Header[lhsidx].View}, Column: v.(parser.Identifier)}

		rhsidx, err := joinView.FieldIndex(fieldref)
		if err != nil {
			return nil, nil, nil, err
		}
		rhs = parser.FieldReference{BaseExpr: v.GetBaseExpr(), View: parser.Identifier{Literal: joinView.Header[rhsidx].View}, Column: v.(parser.Identifier)}

		comps[i] = parser.Comparison{
			LHS:      lhs,
			RHS:      rhs,
			Operator: "=",
		}

		if join.Direction.Token == parser.RIGHT {
			includeFields[i] = rhs
			excludeFields[i] = lhs
		} else {
			includeFields[i] = lhs
			excludeFields[i] = rhs
		}
	}

	if len(comps) == 1 {
		return comps[0], includeFields, excludeFields, nil
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
	return logic, includeFields, excludeFields, nil
}

func CrossJoin(ctx context.Context, view *View, joinView *View) error {
	mergedHeader := MergeHeader(view.Header, joinView.Header)
	records := make(RecordSet, view.RecordLen()*joinView.RecordLen())

	if err := NewGoroutineTaskManager(view.RecordLen(), CalcMinimumRequired(view.RecordLen(), joinView.RecordLen(), MinimumRequiredPerCPUCore)).Run(ctx, func(index int) error {
		start := index * joinView.RecordLen()
		for i := 0; i < joinView.RecordLen(); i++ {
			records[start+i] = append(view.RecordSet[index], joinView.RecordSet[i]...)
		}
		return nil
	}); err != nil {
		return err
	}

	view.Header = mergedHeader
	view.RecordSet = records
	view.FileInfo = nil
	return nil
}

func InnerJoin(ctx context.Context, view *View, joinView *View, condition parser.QueryExpression, parentFilter *Filter) error {
	if condition == nil {
		return CrossJoin(ctx, view, joinView)
	}

	mergedHeader := MergeHeader(view.Header, joinView.Header)

	gm := NewGoroutineTaskManager(view.RecordLen(), CalcMinimumRequired(view.RecordLen(), joinView.RecordLen(), MinimumRequiredPerCPUCore))
	recordsList := make([]RecordSet, gm.Number)
	for i := 0; i < gm.Number; i++ {
		gm.Add()
		go func(thIdx int) {
			start, end := gm.RecordRange(thIdx)
			records := make(RecordSet, 0, end-start)
			filter := NewFilterForRecord(
				&View{
					Header:    mergedHeader,
					RecordSet: make(RecordSet, 1),
				},
				0,
				parentFilter,
			)

		InnerJoinLoop:
			for i := start; i < end; i++ {
				for j := 0; j < joinView.RecordLen(); j++ {
					if gm.HasError() || ctx.Err() != nil {
						break InnerJoinLoop
					}

					mergedRecord := append(view.RecordSet[i], joinView.RecordSet[j]...)
					filter.Records[0].View.RecordSet[0] = mergedRecord

					primary, e := filter.Evaluate(ctx, condition)
					if e != nil {
						gm.SetError(e)
						break InnerJoinLoop
					}
					if primary.Ternary() == ternary.TRUE {
						records = append(records, mergedRecord)
					}
				}
			}

			recordsList[thIdx] = records
			gm.Done()
		}(i)
	}
	gm.Wait()

	if gm.HasError() {
		return gm.Err()
	}
	if ctx.Err() != nil {
		return NewContextIsDone(ctx.Err().Error())
	}

	view.Header = mergedHeader
	view.RecordSet = MergeRecordSetList(recordsList)
	view.FileInfo = nil
	return nil
}

func OuterJoin(ctx context.Context, view *View, joinView *View, condition parser.QueryExpression, direction int, parentFilter *Filter) error {
	if direction == parser.TokenUndefined {
		direction = parser.LEFT
	}

	mergedHeader := MergeHeader(view.Header, joinView.Header)

	if direction == parser.RIGHT {
		view, joinView = joinView, view
	}

	viewEmptyRecord := NewEmptyRecord(view.FieldLen())
	joinViewEmptyRecord := NewEmptyRecord(joinView.FieldLen())

	gm := NewGoroutineTaskManager(view.RecordLen(), CalcMinimumRequired(view.RecordLen(), joinView.RecordLen(), MinimumRequiredPerCPUCore))

	recordsList := make([]RecordSet, gm.Number)
	joinViewMatchesList := make([][]bool, gm.Number)

	for i := 0; i < gm.Number; i++ {
		gm.Add()
		go func(thIdx int) {
			start, end := gm.RecordRange(thIdx)
			records := make(RecordSet, 0, end-start)
			filter := NewFilterForRecord(
				&View{
					Header:    mergedHeader,
					RecordSet: make(RecordSet, 1),
				},
				0,
				parentFilter,
			)

			joinViewMatches := make([]bool, joinView.RecordLen())

		OuterJoinLoop:
			for i := start; i < end; i++ {
				match := false
				for j := 0; j < joinView.RecordLen(); j++ {
					if gm.HasError() || ctx.Err() != nil {
						break OuterJoinLoop
					}

					var mergedRecord Record
					switch direction {
					case parser.RIGHT:
						mergedRecord = append(joinView.RecordSet[j], view.RecordSet[i]...)
					default:
						mergedRecord = append(view.RecordSet[i], joinView.RecordSet[j]...)
					}
					filter.Records[0].View.RecordSet[0] = mergedRecord

					primary, e := filter.Evaluate(ctx, condition)
					if e != nil {
						gm.SetError(e)
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
						record = append(joinViewEmptyRecord, view.RecordSet[i]...)
					default:
						record = append(view.RecordSet[i], joinViewEmptyRecord...)
					}
					records = append(records, record)

				}
			}

			recordsList[thIdx] = records
			joinViewMatchesList[thIdx] = joinViewMatches
			gm.Done()
		}(i)
	}
	gm.Wait()

	if gm.HasError() {
		return gm.Err()
	}
	if ctx.Err() != nil {
		return NewContextIsDone(ctx.Err().Error())
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
				record := append(viewEmptyRecord, joinView.RecordSet[i]...)
				recordsList[len(recordsList)-1] = append(recordsList[len(recordsList)-1], record)
			}
		}
	}

	if direction == parser.RIGHT {
		view, joinView = joinView, view
	}

	view.Header = mergedHeader
	view.RecordSet = MergeRecordSetList(recordsList)
	view.FileInfo = nil
	return nil
}

func CalcMinimumRequired(i1 int, i2 int, defaultMinimumRequired int) int {
	if i1 < 1 || i2 < 1 {
		return defaultMinimumRequired
	}

	p := i1 * i2
	if p <= defaultMinimumRequired {
		return defaultMinimumRequired
	}
	return int(math.Ceil(float64(i1) / math.Floor(float64(p)/float64(defaultMinimumRequired))))
}
