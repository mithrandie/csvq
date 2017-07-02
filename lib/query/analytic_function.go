package query

import (
	"errors"

	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/ternary"
)

var AnalyticFunctions = map[string]func(*View, []parser.Expression, []partitionValue) error{
	"ROW_NUMBER": RowNumber,
	"RANK":       Rank,
	"DENSE_RANK": DenseRank,
}

type partitionValue struct {
	values      []parser.Primary
	orderValues []parser.Primary
	number      float64
	rank        float64
}

func (pv partitionValue) match(values []parser.Primary) bool {
	for i, v := range pv.values {
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

func RowNumber(view *View, args []parser.Expression, partitinList []partitionValue) error {
	if args != nil {
		return errors.New("function ROW_NUMBER takes no argument")
	}

	partitions := partitionValues{}

	for i := range view.Records {
		var idx int
		if idx = partitions.searchIndex(partitinList[i].values); -1 < idx {
			partitions[idx].number++
		} else {
			partitions = append(partitions, partitionValue{
				values: partitinList[i].values,
				number: 1,
			})
			idx = len(partitions) - 1
		}

		view.Records[i] = append(view.Records[i], NewCell(parser.NewInteger(int64(partitions[idx].number))))
	}

	return nil
}

func Rank(view *View, args []parser.Expression, partitinList []partitionValue) error {
	if args != nil {
		return errors.New("function RANK takes no argument")
	}

	partitions := partitionValues{}

	for i := range view.Records {
		var idx int
		if idx = partitions.searchIndex(partitinList[i].values); -1 < idx {
			partitions[idx].number++
			if !partitions[idx].isSameRank(partitinList[i].orderValues) {
				partitions[idx].rank = partitions[idx].number
			}
		} else {
			partitions = append(partitions, partitionValue{
				values:      partitinList[i].values,
				orderValues: partitinList[i].orderValues,
				number:      1,
				rank:        1,
			})
			idx = len(partitions) - 1
		}

		view.Records[i] = append(view.Records[i], NewCell(parser.NewInteger(int64(partitions[idx].rank))))
	}

	return nil
}

func DenseRank(view *View, args []parser.Expression, partitinList []partitionValue) error {
	if args != nil {
		return errors.New("function DENSE_RANK takes no argument")
	}

	partitions := partitionValues{}

	for i := range view.Records {
		var idx int
		if idx = partitions.searchIndex(partitinList[i].values); -1 < idx {
			if !partitions[idx].isSameRank(partitinList[i].orderValues) {
				partitions[idx].rank++
			}
		} else {
			partitions = append(partitions, partitionValue{
				values:      partitinList[i].values,
				orderValues: partitinList[i].orderValues,
				rank:        1,
			})
			idx = len(partitions) - 1
		}

		view.Records[i] = append(view.Records[i], NewCell(parser.NewInteger(int64(partitions[idx].rank))))
	}

	return nil
}
