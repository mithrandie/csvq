package query

import (
	"errors"
	"sync"

	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/ternary"
)

var AnalyticFunctions map[string]func(*View, []parser.Expression, parser.AnalyticClause) error
var defineAnalyticFunctions sync.Once

func DefineAnalyticFunctions() {
	defineAnalyticFunctions.Do(func() {
		AnalyticFunctions = map[string]func(*View, []parser.Expression, parser.AnalyticClause) error{
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

func RowNumber(view *View, args []parser.Expression, clause parser.AnalyticClause) error {
	if args != nil {
		return errors.New("analytic function ROW_NUMBER takes no argument")
	}

	partitions := partitionValues{}

	var filter Filter = append([]FilterRecord{{View: view, RecordIndex: 0}}, view.parentFilter...)
	for i := range view.Records {
		filter[0].RecordIndex = i
		partitionValues, err := filter.evalValues(clause.PartitionValues())
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

func Rank(view *View, args []parser.Expression, clause parser.AnalyticClause) error {
	if args != nil {
		return errors.New("analytic function RANK takes no argument")
	}

	partitions := partitionValues{}

	var filter Filter = append([]FilterRecord{{View: view, RecordIndex: 0}}, view.parentFilter...)
	for i := range view.Records {
		filter[0].RecordIndex = i
		partitionValues, err := filter.evalValues(clause.PartitionValues())
		if err != nil {
			return err
		}

		orderValues, err := filter.evalValues(clause.OrderValues())
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

func DenseRank(view *View, args []parser.Expression, clause parser.AnalyticClause) error {
	if args != nil {
		return errors.New("analytic function DENSE_RANK takes no argument")
	}

	partitions := partitionValues{}

	var filter Filter = append([]FilterRecord{{View: view, RecordIndex: 0}}, view.parentFilter...)
	for i := range view.Records {
		filter[0].RecordIndex = i
		partitionValues, err := filter.evalValues(clause.PartitionValues())
		if err != nil {
			return err
		}

		orderValues, err := filter.evalValues(clause.OrderValues())
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
