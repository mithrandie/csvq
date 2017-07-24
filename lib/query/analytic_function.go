package query

import (
	"sync"

	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/ternary"
)

var AnalyticFunctions map[string]func(*View, parser.AnalyticFunction) error
var defineAnalyticFunctions sync.Once

func DefineAnalyticFunctions() {
	defineAnalyticFunctions.Do(func() {
		AnalyticFunctions = map[string]func(*View, parser.AnalyticFunction) error{
			"ROW_NUMBER": RowNumber,
			"RANK":       Rank,
			"DENSE_RANK": DenseRank,
		}
	})
}

type partitionValue struct {
	partitionValues []parser.Primary
	orderValues     []parser.Primary
	values          map[string]float64
}

func (pv partitionValue) match(values []parser.Primary) bool {
	for i, v := range pv.partitionValues {
		if EquivalentTo(v, values[i]) != ternary.TRUE {
			return false
		}
	}
	return true
}

func (pv partitionValue) isSameRank(orderValues []parser.Primary) bool {
	for i, v := range pv.orderValues {
		if EquivalentTo(v, orderValues[i]) != ternary.TRUE {
			return false
		}
	}
	return true
}

type partitionValues []partitionValue

func (pv partitionValues) searchIndex(values []parser.Primary) int {
	for idx, v := range pv {
		if v.match(values) {
			return idx
		}
	}
	return -1
}

func RowNumber(view *View, fn parser.AnalyticFunction) error {
	if fn.Args != nil {
		return NewFunctionArgumentLengthError(fn.Name, []int{0}, fn)
	}

	partitions := partitionValues{}

	filter := NewFilterForSequentialEvaluation(view, view.ParentFilter)
	for i := range view.Records {
		filter.Records[0].RecordIndex = i
		partitionValues, err := filter.evalValues(fn.AnalyticClause.PartitionValues())
		if err != nil {
			return err
		}

		var idx int
		if idx = partitions.searchIndex(partitionValues); -1 < idx {
			partitions[idx].values["number"]++
		} else {
			partitions = append(partitions, partitionValue{
				partitionValues: partitionValues,
				values: map[string]float64{
					"number": 1,
				},
			})
			idx = len(partitions) - 1
		}

		view.Records[i] = append(view.Records[i], NewCell(parser.NewInteger(int64(partitions[idx].values["number"]))))
	}

	return nil
}

func Rank(view *View, fn parser.AnalyticFunction) error {
	if fn.Args != nil {
		return NewFunctionArgumentLengthError(fn.Name, []int{0}, fn)
	}

	partitions := partitionValues{}

	filter := NewFilterForSequentialEvaluation(view, view.ParentFilter)
	for i := range view.Records {
		filter.Records[0].RecordIndex = i
		partitionValues, err := filter.evalValues(fn.AnalyticClause.PartitionValues())
		if err != nil {
			return err
		}

		orderValues, err := filter.evalValues(fn.AnalyticClause.OrderValues())
		if err != nil {
			return err
		}

		var idx int
		if idx = partitions.searchIndex(partitionValues); -1 < idx {
			partitions[idx].values["number"]++
			if !partitions[idx].isSameRank(orderValues) {
				partitions[idx].values["rank"] = partitions[idx].values["number"]
			}
		} else {
			partitions = append(partitions, partitionValue{
				partitionValues: partitionValues,
				orderValues:     orderValues,
				values: map[string]float64{
					"number": 1,
					"rank":   1,
				},
			})
			idx = len(partitions) - 1
		}

		view.Records[i] = append(view.Records[i], NewCell(parser.NewInteger(int64(partitions[idx].values["rank"]))))
	}

	return nil
}

func DenseRank(view *View, fn parser.AnalyticFunction) error {
	if fn.Args != nil {
		return NewFunctionArgumentLengthError(fn.Name, []int{0}, fn)
	}

	partitions := partitionValues{}

	filter := NewFilterForSequentialEvaluation(view, view.ParentFilter)
	for i := range view.Records {
		filter.Records[0].RecordIndex = i
		partitionValues, err := filter.evalValues(fn.AnalyticClause.PartitionValues())
		if err != nil {
			return err
		}

		orderValues, err := filter.evalValues(fn.AnalyticClause.OrderValues())
		if err != nil {
			return err
		}

		var idx int
		if idx = partitions.searchIndex(partitionValues); -1 < idx {
			if !partitions[idx].isSameRank(orderValues) {
				partitions[idx].values["rank"]++
			}
		} else {
			partitions = append(partitions, partitionValue{
				partitionValues: partitionValues,
				orderValues:     orderValues,
				values: map[string]float64{
					"rank": 1,
				},
			})
			idx = len(partitions) - 1
		}

		view.Records[i] = append(view.Records[i], NewCell(parser.NewInteger(int64(partitions[idx].values["rank"]))))
	}

	return nil
}
