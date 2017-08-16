package query

import (
	"strings"

	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/ternary"
)

var AnalyticFunctions map[string]AnalyticFunction = map[string]AnalyticFunction{
	"ROW_NUMBER":   RowNumber{},
	"RANK":         Rank{},
	"DENSE_RANK":   DenseRank{},
	"CUME_DIST":    CumeDist{},
	"PERCENT_RANK": PercentRank{},
	"NTILE":        NTile{},
	"FIRST_VALUE":  FirstValue{},
	"LAST_VALUE":   LastValue{},
	"NTH_VALUE":    NthValue{},
	"LAG":          Lag{},
	"LEAD":         Lead{},
	"LISTAGG":      AnalyticListAgg{},
}

type AnalyticFunction interface {
	CheckArgsLen(expr parser.AnalyticFunction) error
	Execute(PartitionItemList, parser.AnalyticFunction, *Filter) (map[int]parser.Primary, error)
}

type Partition struct {
	PartitionValues []parser.Primary
	Items           PartitionItemList
}

func (p Partition) Match(values []parser.Primary) bool {
	for i, v := range p.PartitionValues {
		if EquivalentTo(v, values[i]) != ternary.TRUE {
			return false
		}
	}
	return true
}

type PartitionItem struct {
	OrderValues []parser.Primary
	RecordIndex int
}

func (item PartitionItem) IsSameRank(partitionItem PartitionItem) bool {
	if partitionItem.OrderValues == nil && item.OrderValues != nil {
		return false
	}
	for i, v := range partitionItem.OrderValues {
		if EquivalentTo(v, item.OrderValues[i]) != ternary.TRUE {
			return false
		}
	}
	return true
}

type PartitionItemList []PartitionItem

func (list PartitionItemList) Reverse() PartitionItemList {
	reverse := make([]PartitionItem, len(list))
	lastIdx := len(list) - 1
	for i, item := range list {
		reverse[lastIdx-i] = item
	}
	return reverse
}

type PartitionList []Partition

func (list PartitionList) SearchIndex(partitionValues []parser.Primary) int {
	for idx, v := range list {
		if v.Match(partitionValues) {
			return idx
		}
	}
	return -1
}

func Analyze(view *View, fn parser.AnalyticFunction) error {
	const (
		ANALYTIC = iota
		AGGREGATE
		USER_DEFINED
	)

	var anfn AnalyticFunction
	var aggfn AggregateFunction
	var udfn *UserDefinedFunction

	fnType := -1
	var err error

	uname := strings.ToUpper(fn.Name)
	if f, ok := AnalyticFunctions[uname]; ok {
		anfn = f
		fnType = ANALYTIC
	} else if f, ok := AggregateFunctions[uname]; ok {
		aggfn = f
		fnType = AGGREGATE
	} else {
		if udfn, err = view.Filter.FunctionsList.Get(fn, uname); err != nil || !udfn.IsAggregate {
			return NewFunctionNotExistError(fn, fn.Name)
		}
		fnType = USER_DEFINED
	}

	switch fnType {
	case ANALYTIC:
		if err := anfn.CheckArgsLen(fn); err != nil {
			return err
		}
	case AGGREGATE:
		if len(fn.Args) != 1 {
			return NewFunctionArgumentLengthError(fn, fn.Name, []int{1})
		}
	case USER_DEFINED:
		if err := udfn.CheckArgsLen(fn, fn.Name, len(fn.Args)-1); err != nil {
			return err
		}
	}

	partitions := PartitionList{}

	filter := NewFilterForSequentialEvaluation(view, view.Filter)
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

		pitem := PartitionItem{
			OrderValues: orderValues,
			RecordIndex: i,
		}

		if idx := partitions.SearchIndex(partitionValues); -1 < idx {
			partitions[idx].Items = append(partitions[idx].Items, pitem)
		} else {
			partitions = append(
				partitions,
				Partition{
					PartitionValues: partitionValues,
					Items:           PartitionItemList{pitem},
				},
			)
		}
	}

	for _, partition := range partitions {
		if fnType == ANALYTIC {
			list, err := anfn.Execute(partition.Items, fn, filter)
			if err != nil {
				return err
			}
			for idx, value := range list {
				view.Records[idx] = append(view.Records[idx], NewCell(value))
			}
		} else {
			if 0 < len(fn.Args) {
				if _, ok := fn.Args[0].(parser.AllColumns); ok {
					fn.Args[0] = parser.NewInteger(1)
				}
			}

			values, err := view.ListValuesForAnalyticFunctions(fn, partition.Items)
			if err != nil {
				return err
			}

			if fnType == AGGREGATE {
				value := aggfn(values)
				for _, item := range partition.Items {
					view.Records[item.RecordIndex] = append(view.Records[item.RecordIndex], NewCell(value))
				}
			} else { //User Defined Function
				for _, item := range partition.Items {
					filter.Records[0].RecordIndex = item.RecordIndex

					var args []parser.Primary
					argsExprs := fn.Args[1:]
					args = make([]parser.Primary, len(argsExprs))
					for i, v := range argsExprs {
						arg, err := filter.Evaluate(v)
						if err != nil {
							return err
						}
						args[i] = arg
					}

					value, err := udfn.ExecuteAggregate(values, args, view.Filter)
					if err != nil {
						return err
					}

					view.Records[item.RecordIndex] = append(view.Records[item.RecordIndex], NewCell(value))
				}
			}
		}
	}

	return nil
}

func CheckArgsLen(expr parser.AnalyticFunction, length []int) error {
	if len(length) == 1 {
		if len(expr.Args) != length[0] {
			return NewFunctionArgumentLengthError(expr, expr.Name, length)
		}
	} else {
		if len(expr.Args) < length[0] {
			return NewFunctionArgumentLengthErrorWithCustomArgs(expr, expr.Name, "at least "+FormatCount(length[0], "argument"))
		}
		if length[1] < len(expr.Args) {
			return NewFunctionArgumentLengthErrorWithCustomArgs(expr, expr.Name, "at most "+FormatCount(length[1], "argument"))
		}
	}
	return nil
}

type RowNumber struct{}

func (fn RowNumber) CheckArgsLen(expr parser.AnalyticFunction) error {
	return CheckArgsLen(expr, []int{0})
}

func (fn RowNumber) Execute(items PartitionItemList, expr parser.AnalyticFunction, filter *Filter) (map[int]parser.Primary, error) {
	list := make(map[int]parser.Primary, len(items))
	var number int64 = 0
	for _, item := range items {
		number++
		list[item.RecordIndex] = parser.NewInteger(number)
	}

	return list, nil
}

type Rank struct{}

func (fn Rank) CheckArgsLen(expr parser.AnalyticFunction) error {
	return CheckArgsLen(expr, []int{0})
}

func (fn Rank) Execute(items PartitionItemList, expr parser.AnalyticFunction, filter *Filter) (map[int]parser.Primary, error) {
	list := make(map[int]parser.Primary, len(items))
	var number int64 = 0
	var rank int64 = 0
	var currentRank PartitionItem
	for _, item := range items {
		number++
		if !item.IsSameRank(currentRank) {
			rank = number
			currentRank = item
		}
		list[item.RecordIndex] = parser.NewInteger(rank)
	}

	return list, nil
}

type DenseRank struct{}

func (fn DenseRank) CheckArgsLen(expr parser.AnalyticFunction) error {
	return CheckArgsLen(expr, []int{0})
}

func (fn DenseRank) Execute(items PartitionItemList, expr parser.AnalyticFunction, filter *Filter) (map[int]parser.Primary, error) {
	list := make(map[int]parser.Primary, len(items))
	var rank int64 = 0
	var currentRank PartitionItem
	for _, item := range items {
		if !item.IsSameRank(currentRank) {
			rank++
			currentRank = item
		}
		list[item.RecordIndex] = parser.NewInteger(rank)
	}

	return list, nil
}

type CumeDist struct{}

func (fn CumeDist) CheckArgsLen(expr parser.AnalyticFunction) error {
	return CheckArgsLen(expr, []int{0})
}

func (fn CumeDist) Execute(items PartitionItemList, expr parser.AnalyticFunction, filter *Filter) (map[int]parser.Primary, error) {
	list := make(map[int]parser.Primary, len(items))

	groups := perseCumulativeGroups(items)
	total := float64(len(items))
	cumulative := float64(0)
	for _, group := range groups {
		cumulative += float64(len(group))
		dist := cumulative / total

		for _, idx := range group {
			list[idx] = parser.NewFloat(dist)
		}
	}

	return list, nil
}

type PercentRank struct{}

func (fn PercentRank) CheckArgsLen(expr parser.AnalyticFunction) error {
	return CheckArgsLen(expr, []int{0})
}

func (fn PercentRank) Execute(items PartitionItemList, expr parser.AnalyticFunction, filter *Filter) (map[int]parser.Primary, error) {
	list := make(map[int]parser.Primary, len(items))

	groups := perseCumulativeGroups(items)
	denom := float64(len(items) - 1)
	cumulative := float64(0)
	for _, group := range groups {
		var dist float64 = 1
		if 0 < denom {
			dist = cumulative / denom
		}

		for _, idx := range group {
			list[idx] = parser.NewFloat(dist)
		}

		cumulative += float64(len(group))
	}

	return list, nil
}

func perseCumulativeGroups(items PartitionItemList) [][]int {
	groups := [][]int{}
	var currentRank PartitionItem
	for _, item := range items {
		if !item.IsSameRank(currentRank) {
			groups = append(groups, []int{item.RecordIndex})
			currentRank = item
		} else {
			groups[len(groups)-1] = append(groups[len(groups)-1], item.RecordIndex)
		}
	}
	return groups
}

type NTile struct{}

func (fn NTile) CheckArgsLen(expr parser.AnalyticFunction) error {
	return CheckArgsLen(expr, []int{1})
}

func (fn NTile) Execute(items PartitionItemList, expr parser.AnalyticFunction, filter *Filter) (map[int]parser.Primary, error) {
	argsFilter := filter.CreateNode()
	argsFilter.Records = nil

	tileNumber := 0
	p, err := argsFilter.Evaluate(expr.Args[0])
	if err != nil {
		return nil, NewFunctionInvalidArgumentError(expr, expr.Name, "the first argument must be an integer")
	}
	i := parser.PrimaryToInteger(p)
	if parser.IsNull(i) {
		return nil, NewFunctionInvalidArgumentError(expr, expr.Name, "the first argument must be an integer")
	}
	tileNumber = int(i.(parser.Integer).Value())
	if tileNumber < 1 {
		return nil, NewFunctionInvalidArgumentError(expr, expr.Name, "the first argument must be greater than 0")
	}

	total := len(items)
	perTile := total / tileNumber
	mod := total % tileNumber

	if perTile < 1 {
		perTile = 1
		mod = 0
	}

	list := make(map[int]parser.Primary, len(items))
	var tile int64 = 1
	var count int = 0
	for _, item := range items {
		count++

		switch {
		case perTile+1 < count:
			tile++
			count = 1
		case perTile+1 == count:
			if 0 < mod {
				mod--
			} else {
				tile++
				count = 1
			}
		}
		list[item.RecordIndex] = parser.NewInteger(tile)
	}

	return list, nil
}

type FirstValue struct{}

func (fn FirstValue) CheckArgsLen(expr parser.AnalyticFunction) error {
	return CheckArgsLen(expr, []int{1})
}

func (fn FirstValue) Execute(items PartitionItemList, expr parser.AnalyticFunction, filter *Filter) (map[int]parser.Primary, error) {
	return setNthValue(items, expr, filter, 1)
}

type LastValue struct{}

func (fn LastValue) CheckArgsLen(expr parser.AnalyticFunction) error {
	return CheckArgsLen(expr, []int{1})
}

func (fn LastValue) Execute(items PartitionItemList, expr parser.AnalyticFunction, filter *Filter) (map[int]parser.Primary, error) {
	return setNthValue(items.Reverse(), expr, filter, 1)
}

type NthValue struct{}

func (fn NthValue) CheckArgsLen(expr parser.AnalyticFunction) error {
	return CheckArgsLen(expr, []int{2})
}

func (fn NthValue) Execute(items PartitionItemList, expr parser.AnalyticFunction, filter *Filter) (map[int]parser.Primary, error) {
	argsFilter := filter.CreateNode()
	argsFilter.Records = nil

	n := 0
	p, err := argsFilter.Evaluate(expr.Args[1])
	if err != nil {
		return nil, NewFunctionInvalidArgumentError(expr, expr.Name, "the second argument must be an integer")
	}
	pi := parser.PrimaryToInteger(p)
	if parser.IsNull(pi) {
		return nil, NewFunctionInvalidArgumentError(expr, expr.Name, "the second argument must be an integer")
	}
	n = int(pi.(parser.Integer).Value())
	if n < 1 {
		return nil, NewFunctionInvalidArgumentError(expr, expr.Name, "the second argument must be greater than 0")
	}

	return setNthValue(items, expr, filter, n)
}

func setNthValue(items PartitionItemList, expr parser.AnalyticFunction, filter *Filter, n int) (map[int]parser.Primary, error) {
	var value parser.Primary = parser.NewNull()

	count := 0
	if n <= len(items) {
		for _, item := range items {
			filter.Records[0].RecordIndex = item.RecordIndex
			p, err := filter.Evaluate(expr.Args[0])
			if err != nil {
				return nil, err
			}

			if expr.IgnoreNulls && parser.IsNull(p) {
				continue
			}

			count++
			if count == n {
				value = p
				break
			}
		}
	}

	list := make(map[int]parser.Primary, len(items))
	for _, item := range items {
		list[item.RecordIndex] = value
	}

	return list, nil
}

type Lag struct{}

func (fn Lag) CheckArgsLen(expr parser.AnalyticFunction) error {
	return CheckArgsLen(expr, []int{1, 3})
}

func (fn Lag) Execute(items PartitionItemList, expr parser.AnalyticFunction, filter *Filter) (map[int]parser.Primary, error) {
	return setLag(items, expr, filter)
}

type Lead struct{}

func (fn Lead) CheckArgsLen(expr parser.AnalyticFunction) error {
	return CheckArgsLen(expr, []int{1, 3})
}

func (fn Lead) Execute(items PartitionItemList, expr parser.AnalyticFunction, filter *Filter) (map[int]parser.Primary, error) {
	return setLag(items.Reverse(), expr, filter)
}

func setLag(items PartitionItemList, expr parser.AnalyticFunction, filter *Filter) (map[int]parser.Primary, error) {
	argsFilter := filter.CreateNode()
	argsFilter.Records = nil

	offset := 1
	if 1 < len(expr.Args) {
		p, err := argsFilter.Evaluate(expr.Args[1])
		if err != nil {
			return nil, NewFunctionInvalidArgumentError(expr, expr.Name, "the second argument must be an integer")
		}
		i := parser.PrimaryToInteger(p)
		if parser.IsNull(i) {
			return nil, NewFunctionInvalidArgumentError(expr, expr.Name, "the second argument must be an integer")
		}
		offset = int(i.(parser.Integer).Value())
	}

	var defaultValue parser.Primary = parser.NewNull()
	if 2 < len(expr.Args) {
		p, err := argsFilter.Evaluate(expr.Args[2])
		if err != nil {
			return nil, NewFunctionInvalidArgumentError(expr, expr.Name, "the third argument must be a primitive type")
		}
		defaultValue = p
	}

	list := make(map[int]parser.Primary, len(items))
	values := []parser.Primary{}
	for _, item := range items {
		filter.Records[0].RecordIndex = item.RecordIndex
		p, err := filter.Evaluate(expr.Args[0])
		if err != nil {
			return nil, err
		}

		values = append(values, p)

		lagIdx := len(values) - 1 - offset
		value := defaultValue
		if 0 <= lagIdx && lagIdx < len(values) {
			for i := lagIdx; i >= 0; i-- {
				if expr.IgnoreNulls && parser.IsNull(values[i]) {
					continue
				}
				value = values[i]
				break
			}
		}
		list[item.RecordIndex] = value
	}

	return list, nil
}

type AnalyticListAgg struct{}

func (fn AnalyticListAgg) CheckArgsLen(expr parser.AnalyticFunction) error {
	return CheckArgsLen(expr, []int{1, 2})
}

func (fn AnalyticListAgg) Execute(items PartitionItemList, expr parser.AnalyticFunction, filter *Filter) (map[int]parser.Primary, error) {
	argsFilter := filter.CreateNode()
	argsFilter.Records = nil

	separator := ""
	if len(expr.Args) == 2 {
		p, err := argsFilter.Evaluate(expr.Args[1])
		if err != nil {
			return nil, NewFunctionInvalidArgumentError(expr, expr.Name, "the second argument must be a string")
		}
		s := parser.PrimaryToString(p)
		if parser.IsNull(s) {
			return nil, NewFunctionInvalidArgumentError(expr, expr.Name, "the second argument must be a string")
		}
		separator = s.(parser.String).Value()
	}

	values, err := filter.Records[0].View.ListValuesForAnalyticFunctions(expr, items)
	if err != nil {
		return nil, err
	}

	value := ListAgg(values, separator)

	list := make(map[int]parser.Primary, len(items))
	for _, item := range items {
		list[item.RecordIndex] = value
	}

	return list, nil
}
