package query

import (
	"sort"
	"strings"
	"sync"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/parser"
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
	Execute(Partition, parser.AnalyticFunction, *Filter) (map[int]parser.Primary, error)
}

type Partition []int

func (p Partition) Reverse() {
	sort.Sort(sort.Reverse(sort.IntSlice(p)))
}

type Partitions map[string]Partition

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

	cpu := NumberOfCPU(view.RecordLen())
	partitionKeys := make([]string, view.RecordLen())

	wg := sync.WaitGroup{}
	for i := 0; i < cpu; i++ {
		wg.Add(1)
		go func(thIdx int) {
			start, end := RecordRange(thIdx, view.RecordLen(), cpu)
			filter := NewFilterForSequentialEvaluation(view, view.Filter)

		AnalyzePrepareLoop:
			for i := start; i < end; i++ {
				if err != nil {
					break AnalyzePrepareLoop
				}

				filter.Records[0].RecordIndex = i

				var partitionKey string
				if fn.AnalyticClause.PartitionValues() != nil {
					partitionValues, e := filter.evalValues(fn.AnalyticClause.PartitionValues())
					if e != nil {
						err = e
						break AnalyzePrepareLoop
					}
					partitionKey = SerializeComparisonKeys(partitionValues)
				}

				partitionKeys[i] = partitionKey
			}

			wg.Done()
		}(i)
	}

	wg.Wait()

	if err != nil {
		return err
	}

	partitions := Partitions{}
	partitionMapKeys := []string{}
	for i, key := range partitionKeys {
		if _, ok := partitions[key]; ok {
			partitions[key] = append(partitions[key], i)
		} else {
			partitions[key] = Partition{i}
			partitionMapKeys = append(partitionMapKeys, key)
		}
	}

	cpu = cmd.GetFlags().CPU
	if cpu < len(partitionMapKeys) {
		cpu = len(partitionMapKeys)
	}

	for i := 0; i < cpu; i++ {
		wg.Add(1)
		go func(thIdx int) {
			start, end := RecordRange(thIdx, len(partitionMapKeys), cpu)
			filter := NewFilterForSequentialEvaluation(view, view.Filter)

		AnalyzeLoop:
			for i := start; i < end; i++ {
				if fnType == ANALYTIC {
					list, e := anfn.Execute(partitions[partitionMapKeys[i]], fn, filter)
					if e != nil {
						err = e
						break AnalyzeLoop
					}
					for idx, value := range list {
						view.Records[idx] = append(view.Records[idx], NewCell(value))
					}
				} else {
					if 0 < len(fn.Args) {
						if _, ok := fn.Args[0].(parser.AllColumns); ok {
							fn.Args[0] = parser.NewIntegerValue(1)
						}
					}

					values, e := view.ListValuesForAnalyticFunctions(fn, partitions[partitionMapKeys[i]])
					if e != nil {
						err = e
						break AnalyzeLoop
					}

					if fnType == AGGREGATE {
						value := aggfn(values)
						for _, idx := range partitions[partitionMapKeys[i]] {
							view.Records[idx] = append(view.Records[idx], NewCell(value))
						}
					} else { //User Defined Function
						for _, idx := range partitions[partitionMapKeys[i]] {
							filter.Records[0].RecordIndex = idx

							var args []parser.Primary
							argsExprs := fn.Args[1:]
							args = make([]parser.Primary, len(argsExprs))
							for i, v := range argsExprs {
								arg, e := filter.Evaluate(v)
								if e != nil {
									err = e
									break AnalyzeLoop
								}
								args[i] = arg
							}

							value, e := udfn.ExecuteAggregate(values, args, view.Filter)
							if e != nil {
								err = e
								break AnalyzeLoop
							}

							view.Records[idx] = append(view.Records[idx], NewCell(value))
						}
					}
				}
			}

			wg.Done()
		}(i)
	}

	wg.Wait()

	return err
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

func (fn RowNumber) Execute(partition Partition, expr parser.AnalyticFunction, filter *Filter) (map[int]parser.Primary, error) {
	list := make(map[int]parser.Primary, len(partition))
	var number int64 = 0
	for _, idx := range partition {
		number++
		list[idx] = parser.NewInteger(number)
	}

	return list, nil
}

type Rank struct{}

func (fn Rank) CheckArgsLen(expr parser.AnalyticFunction) error {
	return CheckArgsLen(expr, []int{0})
}

func (fn Rank) Execute(partition Partition, expr parser.AnalyticFunction, filter *Filter) (map[int]parser.Primary, error) {
	list := make(map[int]parser.Primary, len(partition))
	var number int64 = 0
	var rank int64 = 0
	var currentRank SortValues
	for _, idx := range partition {
		number++
		if filter.Records[0].View.recordSortValues == nil || !filter.Records[0].View.recordSortValues[idx].EquivalentTo(currentRank) {
			rank = number
			if filter.Records[0].View.recordSortValues != nil {
				currentRank = filter.Records[0].View.recordSortValues[idx]
			}
		}
		list[idx] = parser.NewInteger(rank)
	}

	return list, nil
}

type DenseRank struct{}

func (fn DenseRank) CheckArgsLen(expr parser.AnalyticFunction) error {
	return CheckArgsLen(expr, []int{0})
}

func (fn DenseRank) Execute(partition Partition, expr parser.AnalyticFunction, filter *Filter) (map[int]parser.Primary, error) {
	list := make(map[int]parser.Primary, len(partition))
	var rank int64 = 0
	var currentRank SortValues
	for _, idx := range partition {
		if filter.Records[0].View.recordSortValues == nil || !filter.Records[0].View.recordSortValues[idx].EquivalentTo(currentRank) {
			rank++
			if filter.Records[0].View.recordSortValues != nil {
				currentRank = filter.Records[0].View.recordSortValues[idx]
			}
		}
		list[idx] = parser.NewInteger(rank)
	}

	return list, nil
}

type CumeDist struct{}

func (fn CumeDist) CheckArgsLen(expr parser.AnalyticFunction) error {
	return CheckArgsLen(expr, []int{0})
}

func (fn CumeDist) Execute(partition Partition, expr parser.AnalyticFunction, filter *Filter) (map[int]parser.Primary, error) {
	list := make(map[int]parser.Primary, len(partition))

	groups := perseCumulativeGroups(partition, filter.Records[0].View)
	total := float64(len(partition))
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

func (fn PercentRank) Execute(partition Partition, expr parser.AnalyticFunction, filter *Filter) (map[int]parser.Primary, error) {
	list := make(map[int]parser.Primary, len(partition))

	groups := perseCumulativeGroups(partition, filter.Records[0].View)
	denom := float64(len(partition) - 1)
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

func perseCumulativeGroups(partition Partition, view *View) [][]int {
	groups := [][]int{}
	var currentRank SortValues
	for _, idx := range partition {
		if view.recordSortValues == nil || !view.recordSortValues[idx].EquivalentTo(currentRank) {
			groups = append(groups, []int{idx})
			if view.recordSortValues != nil {
				currentRank = view.recordSortValues[idx]
			}
		} else {
			groups[len(groups)-1] = append(groups[len(groups)-1], idx)
		}
	}
	return groups
}

type NTile struct{}

func (fn NTile) CheckArgsLen(expr parser.AnalyticFunction) error {
	return CheckArgsLen(expr, []int{1})
}

func (fn NTile) Execute(partition Partition, expr parser.AnalyticFunction, filter *Filter) (map[int]parser.Primary, error) {
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

	total := len(partition)
	perTile := total / tileNumber
	mod := total % tileNumber

	if perTile < 1 {
		perTile = 1
		mod = 0
	}

	list := make(map[int]parser.Primary, len(partition))
	var tile int64 = 1
	var count int = 0
	for _, idx := range partition {
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
		list[idx] = parser.NewInteger(tile)
	}

	return list, nil
}

type FirstValue struct{}

func (fn FirstValue) CheckArgsLen(expr parser.AnalyticFunction) error {
	return CheckArgsLen(expr, []int{1})
}

func (fn FirstValue) Execute(partition Partition, expr parser.AnalyticFunction, filter *Filter) (map[int]parser.Primary, error) {
	return setNthValue(partition, expr, filter, 1)
}

type LastValue struct{}

func (fn LastValue) CheckArgsLen(expr parser.AnalyticFunction) error {
	return CheckArgsLen(expr, []int{1})
}

func (fn LastValue) Execute(partition Partition, expr parser.AnalyticFunction, filter *Filter) (map[int]parser.Primary, error) {
	partition.Reverse()
	return setNthValue(partition, expr, filter, 1)
}

type NthValue struct{}

func (fn NthValue) CheckArgsLen(expr parser.AnalyticFunction) error {
	return CheckArgsLen(expr, []int{2})
}

func (fn NthValue) Execute(partition Partition, expr parser.AnalyticFunction, filter *Filter) (map[int]parser.Primary, error) {
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

	return setNthValue(partition, expr, filter, n)
}

func setNthValue(partition Partition, expr parser.AnalyticFunction, filter *Filter, n int) (map[int]parser.Primary, error) {
	var value parser.Primary = parser.NewNull()

	count := 0
	if n <= len(partition) {
		for _, idx := range partition {
			filter.Records[0].RecordIndex = idx
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

	list := make(map[int]parser.Primary, len(partition))
	for _, idx := range partition {
		list[idx] = value
	}

	return list, nil
}

type Lag struct{}

func (fn Lag) CheckArgsLen(expr parser.AnalyticFunction) error {
	return CheckArgsLen(expr, []int{1, 3})
}

func (fn Lag) Execute(partition Partition, expr parser.AnalyticFunction, filter *Filter) (map[int]parser.Primary, error) {
	return setLag(partition, expr, filter)
}

type Lead struct{}

func (fn Lead) CheckArgsLen(expr parser.AnalyticFunction) error {
	return CheckArgsLen(expr, []int{1, 3})
}

func (fn Lead) Execute(partition Partition, expr parser.AnalyticFunction, filter *Filter) (map[int]parser.Primary, error) {
	partition.Reverse()
	return setLag(partition, expr, filter)
}

func setLag(partition Partition, expr parser.AnalyticFunction, filter *Filter) (map[int]parser.Primary, error) {
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

	list := make(map[int]parser.Primary, len(partition))
	values := []parser.Primary{}
	for _, idx := range partition {
		filter.Records[0].RecordIndex = idx
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
		list[idx] = value
	}

	return list, nil
}

type AnalyticListAgg struct{}

func (fn AnalyticListAgg) CheckArgsLen(expr parser.AnalyticFunction) error {
	return CheckArgsLen(expr, []int{1, 2})
}

func (fn AnalyticListAgg) Execute(partition Partition, expr parser.AnalyticFunction, filter *Filter) (map[int]parser.Primary, error) {
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

	values, err := filter.Records[0].View.ListValuesForAnalyticFunctions(expr, partition)
	if err != nil {
		return nil, err
	}

	value := ListAgg(values, separator)

	list := make(map[int]parser.Primary, len(partition))
	for _, idx := range partition {
		list[idx] = value
	}

	return list, nil
}
