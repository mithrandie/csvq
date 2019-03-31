package query

import (
	"bytes"
	"context"
	"sort"
	"strings"

	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/value"
)

var AnalyticFunctions = map[string]AnalyticFunction{
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
	"JSON_AGG":     AnalyticJsonAgg{},
}

type AnalyticFunction interface {
	CheckArgsLen(expr parser.AnalyticFunction) error
	Execute(context.Context, *Filter, Partition, parser.AnalyticFunction) (map[int]value.Primary, error)
}

type Partition []int

func (p Partition) Reverse() {
	sort.Sort(sort.Reverse(sort.IntSlice(p)))
}

type Partitions map[string]Partition

func Analyze(ctx context.Context, view *View, fn parser.AnalyticFunction, partitionIndices []int) error {
	var anfn AnalyticFunction
	var aggfn AggregateFunction
	var udfn *UserDefinedFunction

	var err error

	uname := strings.ToUpper(fn.Name)
	if f, ok := AnalyticFunctions[uname]; ok {
		anfn = f
	} else if f, ok := AggregateFunctions[uname]; ok {
		aggfn = f
	} else {
		if udfn, err = view.Filter.functions.Get(fn, uname); err != nil || !udfn.IsAggregate {
			return NewFunctionNotExistError(fn, fn.Name)
		}
	}

	if anfn != nil {
		if err := anfn.CheckArgsLen(fn); err != nil {
			return err
		}
	} else if aggfn != nil {
		if len(fn.Args) != 1 {
			return NewFunctionArgumentLengthError(fn, fn.Name, []int{1})
		}
	} else {
		if err := udfn.CheckArgsLen(fn, fn.Name, len(fn.Args)-1); err != nil {
			return err
		}
	}

	if view.sortValuesInEachCell == nil {
		view.sortValuesInEachCell = make([][]*SortValue, view.RecordLen())
	}

	partitionKeys := make([]string, view.RecordLen())
	if err = NewGoroutineTaskManager(view.RecordLen(), -1, view.Tx.Flags.CPU).Run(ctx, func(index int) error {
		keyBuf := new(bytes.Buffer)

		if view.sortValuesInEachCell[index] == nil {
			view.sortValuesInEachCell[index] = make([]*SortValue, cap(view.RecordSet[index]))
		}
		if partitionIndices != nil {
			sortValues := make(SortValues, len(partitionIndices))
			for j, idx := range partitionIndices {
				if idx < len(view.sortValuesInEachCell[index]) && view.sortValuesInEachCell[index][idx] != nil {
					sortValues[j] = view.sortValuesInEachCell[index][idx]
				} else {
					sortValues[j] = NewSortValue(view.RecordSet[index][idx].Value(), view.Tx.Flags)
					if idx < len(view.sortValuesInEachCell[index]) {
						view.sortValuesInEachCell[index][idx] = sortValues[j]
					}
				}
			}
			sortValues.Serialize(keyBuf)
		}

		partitionKeys[index] = keyBuf.String()
		return nil
	}); err != nil {
		return err
	}

	partitions := Partitions{}
	partitionMapKeys := make([]string, 0)
	for i, key := range partitionKeys {
		if _, ok := partitions[key]; ok {
			partitions[key] = append(partitions[key], i)
		} else {
			partitions[key] = Partition{i}
			partitionMapKeys = append(partitionMapKeys, key)
		}
	}

	gm := NewGoroutineTaskManager(len(partitionMapKeys), -1, view.Tx.Flags.CPU)
	for i := 0; i < gm.Number; i++ {
		gm.Add()
		go func(thIdx int) {
			start, end := gm.RecordRange(thIdx)
			filter := NewFilterForSequentialEvaluation(view.Filter, view)

		AnalyzeLoop:
			for i := start; i < end; i++ {
				if gm.HasError() || ctx.Err() != nil {
					break AnalyzeLoop
				}

				if anfn != nil {
					list, e := anfn.Execute(ctx, filter, partitions[partitionMapKeys[i]], fn)
					if e != nil {
						gm.SetError(e)
						break AnalyzeLoop
					}
					for idx, val := range list {
						view.RecordSet[idx] = append(view.RecordSet[idx], NewCell(val))
					}
				} else {
					if 0 < len(fn.Args) {
						if _, ok := fn.Args[0].(parser.AllColumns); ok {
							fn.Args[0] = parser.NewIntegerValue(1)
						}
					}

					if aggfn != nil {
						partition := partitions[partitionMapKeys[i]]
						frameSet := WindowFrameSet(partition, fn.AnalyticClause)

						valueCache := make(map[int]value.Primary, len(partition))

						for _, frame := range frameSet {
							values, e := windowValues(ctx, filter, frame, partition, fn, valueCache)
							if e != nil {
								gm.SetError(e)
								break AnalyzeLoop
							}
							val := aggfn(values, view.Tx.Flags)

							for _, idx := range frame.Records {
								view.RecordSet[idx] = append(view.RecordSet[idx], NewCell(val))
							}
						}
					} else { //User Defined Function
						partition := partitions[partitionMapKeys[i]]
						frameSet := WindowFrameSet(partition, fn.AnalyticClause)

						valueCache := make(map[int]value.Primary, len(partition))

						for _, frame := range frameSet {
							values, e := windowValues(ctx, filter, frame, partition, fn, valueCache)
							if e != nil {
								gm.SetError(e)
								break AnalyzeLoop
							}

							for _, idx := range frame.Records {
								filter.records[0].recordIndex = idx

								var args []value.Primary
								argsExprs := fn.Args[1:]
								args = make([]value.Primary, len(argsExprs))
								for i, v := range argsExprs {
									arg, e := filter.Evaluate(ctx, v)
									if e != nil {
										gm.SetError(e)
										break AnalyzeLoop
									}
									args[i] = arg
								}

								val, e := udfn.ExecuteAggregate(ctx, view.Filter, values, args)
								if e != nil {
									gm.SetError(e)
									break AnalyzeLoop
								}

								view.RecordSet[idx] = append(view.RecordSet[idx], NewCell(val))
							}
						}
					}
				}
			}

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
	return nil
}

type WindowFrame struct {
	Low     int
	High    int
	Records []int
}

func WindowFrameSet(partition Partition, expr parser.AnalyticClause) []WindowFrame {
	var singleFrameSet = func(partition Partition) []WindowFrame {
		indices := make([]int, len(partition))
		for i, idx := range partition {
			indices[i] = idx
		}
		return []WindowFrame{{Low: 0, High: len(partition) - 1, Records: indices}}
	}

	var frameIndex = func(current int, length int, framePosition parser.WindowFramePosition) int {
		var idx int

		switch framePosition.Direction {
		case parser.CURRENT:
			idx = current
		case parser.PRECEDING:
			if framePosition.Unbounded {
				idx = 0
			} else {
				idx = current - framePosition.Offset
			}
		case parser.FOLLOWING:
			if framePosition.Unbounded {
				idx = length - 1
			} else {
				idx = current + framePosition.Offset
			}
		}

		return idx
	}

	length := len(partition)

	if expr.OrderByClause == nil {
		return singleFrameSet(partition)
	}

	frameSet := make([]WindowFrame, 0, length)

	var windowClause parser.WindowingClause
	if expr.WindowingClause == nil {
		windowClause = parser.WindowingClause{
			FrameLow: parser.WindowFramePosition{
				Direction: parser.PRECEDING,
				Unbounded: true,
			},
		}
	} else {
		windowClause = expr.WindowingClause.(parser.WindowingClause)
	}
	frameLow := windowClause.FrameLow.(parser.WindowFramePosition)

	if windowClause.FrameHigh == nil {
		for current := 0; current < length; current++ {
			frameSet = append(frameSet, WindowFrame{
				Low:     frameIndex(current, length, frameLow),
				High:    current,
				Records: []int{partition[current]},
			})
		}
	} else {
		frameHigh := windowClause.FrameHigh.(parser.WindowFramePosition)
		if frameLow.Direction == parser.PRECEDING && frameLow.Unbounded && frameHigh.Direction == parser.FOLLOWING && frameHigh.Unbounded {
			return singleFrameSet(partition)
		}

		for current := 0; current < length; current++ {
			frameSet = append(frameSet, WindowFrame{
				Low:     frameIndex(current, length, frameLow),
				High:    frameIndex(current, length, frameHigh),
				Records: []int{partition[current]},
			})
		}
	}

	return frameSet
}

func windowValues(ctx context.Context, filter *Filter, frame WindowFrame, partition Partition, expr parser.AnalyticFunction, valueCache map[int]value.Primary) ([]value.Primary, error) {
	values := make([]value.Primary, 0, frame.High-frame.Low+1)

	for i := frame.Low; i <= frame.High; i++ {
		if i < 0 || len(partition) <= i {
			continue
		}

		recordIdx := partition[i]
		if v, ok := valueCache[recordIdx]; ok {
			values = append(values, v)
		} else {
			filter.records[0].recordIndex = recordIdx
			p, e := filter.Evaluate(ctx, expr.Args[0])
			if e != nil {
				return nil, e
			}
			valueCache[recordIdx] = p
			values = append(values, p)
		}
	}

	if expr.IsDistinct() {
		values = Distinguish(values, filter.tx.Flags)
	}
	return values, nil
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

func (fn RowNumber) Execute(ctx context.Context, filter *Filter, partition Partition, expr parser.AnalyticFunction) (map[int]value.Primary, error) {
	list := make(map[int]value.Primary, len(partition))
	var number int64 = 0
	for _, idx := range partition {
		number++
		list[idx] = value.NewInteger(number)
	}

	return list, nil
}

type Rank struct{}

func (fn Rank) CheckArgsLen(expr parser.AnalyticFunction) error {
	return CheckArgsLen(expr, []int{0})
}

func (fn Rank) Execute(ctx context.Context, filter *Filter, partition Partition, expr parser.AnalyticFunction) (map[int]value.Primary, error) {
	list := make(map[int]value.Primary, len(partition))
	var number int64 = 0
	var rank int64 = 0
	var currentRank SortValues
	for _, idx := range partition {
		number++
		if filter.records[0].view.sortValuesInEachRecord == nil || !filter.records[0].view.sortValuesInEachRecord[idx].EquivalentTo(currentRank) {
			rank = number
			if filter.records[0].view.sortValuesInEachRecord != nil {
				currentRank = filter.records[0].view.sortValuesInEachRecord[idx]
			}
		}
		list[idx] = value.NewInteger(rank)
	}

	return list, nil
}

type DenseRank struct{}

func (fn DenseRank) CheckArgsLen(expr parser.AnalyticFunction) error {
	return CheckArgsLen(expr, []int{0})
}

func (fn DenseRank) Execute(ctx context.Context, filter *Filter, partition Partition, expr parser.AnalyticFunction) (map[int]value.Primary, error) {
	list := make(map[int]value.Primary, len(partition))
	var rank int64 = 0
	var currentRank SortValues
	for _, idx := range partition {
		if filter.records[0].view.sortValuesInEachRecord == nil || !filter.records[0].view.sortValuesInEachRecord[idx].EquivalentTo(currentRank) {
			rank++
			if filter.records[0].view.sortValuesInEachRecord != nil {
				currentRank = filter.records[0].view.sortValuesInEachRecord[idx]
			}
		}
		list[idx] = value.NewInteger(rank)
	}

	return list, nil
}

type CumeDist struct{}

func (fn CumeDist) CheckArgsLen(expr parser.AnalyticFunction) error {
	return CheckArgsLen(expr, []int{0})
}

func (fn CumeDist) Execute(ctx context.Context, filter *Filter, partition Partition, expr parser.AnalyticFunction) (map[int]value.Primary, error) {
	list := make(map[int]value.Primary, len(partition))

	groups := perseCumulativeGroups(partition, filter.records[0].view)
	total := float64(len(partition))
	cumulative := float64(0)
	for _, group := range groups {
		cumulative += float64(len(group))
		dist := cumulative / total

		for _, idx := range group {
			list[idx] = value.NewFloat(dist)
		}
	}

	return list, nil
}

type PercentRank struct{}

func (fn PercentRank) CheckArgsLen(expr parser.AnalyticFunction) error {
	return CheckArgsLen(expr, []int{0})
}

func (fn PercentRank) Execute(ctx context.Context, filter *Filter, partition Partition, expr parser.AnalyticFunction) (map[int]value.Primary, error) {
	list := make(map[int]value.Primary, len(partition))

	groups := perseCumulativeGroups(partition, filter.records[0].view)
	denom := float64(len(partition) - 1)
	cumulative := float64(0)
	for _, group := range groups {
		var dist float64 = 1
		if 0 < denom {
			dist = cumulative / denom
		}

		for _, idx := range group {
			list[idx] = value.NewFloat(dist)
		}

		cumulative += float64(len(group))
	}

	return list, nil
}

func perseCumulativeGroups(partition Partition, view *View) [][]int {
	groups := make([][]int, 0)
	var currentRank SortValues
	for _, idx := range partition {
		if view.sortValuesInEachRecord == nil || !view.sortValuesInEachRecord[idx].EquivalentTo(currentRank) {
			groups = append(groups, []int{idx})
			if view.sortValuesInEachRecord != nil {
				currentRank = view.sortValuesInEachRecord[idx]
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

func (fn NTile) Execute(ctx context.Context, filter *Filter, partition Partition, expr parser.AnalyticFunction) (map[int]value.Primary, error) {
	argsFilter := filter.CreateNode()
	argsFilter.records = nil

	tileNumber := 0
	p, err := argsFilter.Evaluate(ctx, expr.Args[0])
	if err != nil {
		return nil, NewFunctionInvalidArgumentError(expr, expr.Name, "the first argument must be an integer")
	}
	i := value.ToInteger(p)
	if value.IsNull(i) {
		return nil, NewFunctionInvalidArgumentError(expr, expr.Name, "the first argument must be an integer")
	}
	tileNumber = int(i.(value.Integer).Raw())
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

	list := make(map[int]value.Primary, len(partition))
	var tile int64 = 1
	var count = 0
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
		list[idx] = value.NewInteger(tile)
	}

	return list, nil
}

type FirstValue struct{}

func (fn FirstValue) CheckArgsLen(expr parser.AnalyticFunction) error {
	return CheckArgsLen(expr, []int{1})
}

func (fn FirstValue) Execute(ctx context.Context, filter *Filter, partition Partition, expr parser.AnalyticFunction) (map[int]value.Primary, error) {
	return setNthValue(ctx, filter, partition, expr, 1)
}

type LastValue struct{}

func (fn LastValue) CheckArgsLen(expr parser.AnalyticFunction) error {
	return CheckArgsLen(expr, []int{1})
}

func (fn LastValue) Execute(ctx context.Context, filter *Filter, partition Partition, expr parser.AnalyticFunction) (map[int]value.Primary, error) {
	partition.Reverse()
	return setNthValue(ctx, filter, partition, expr, 1)
}

type NthValue struct{}

func (fn NthValue) CheckArgsLen(expr parser.AnalyticFunction) error {
	return CheckArgsLen(expr, []int{2})
}

func (fn NthValue) Execute(ctx context.Context, filter *Filter, partition Partition, expr parser.AnalyticFunction) (map[int]value.Primary, error) {
	argsFilter := filter.CreateNode()
	argsFilter.records = nil

	n := 0
	p, err := argsFilter.Evaluate(ctx, expr.Args[1])
	if err != nil {
		return nil, NewFunctionInvalidArgumentError(expr, expr.Name, "the second argument must be an integer")
	}
	pi := value.ToInteger(p)
	if value.IsNull(pi) {
		return nil, NewFunctionInvalidArgumentError(expr, expr.Name, "the second argument must be an integer")
	}
	n = int(pi.(value.Integer).Raw())
	if n < 1 {
		return nil, NewFunctionInvalidArgumentError(expr, expr.Name, "the second argument must be greater than 0")
	}

	return setNthValue(ctx, filter, partition, expr, n)
}

func setNthValue(ctx context.Context, filter *Filter, partition Partition, expr parser.AnalyticFunction, n int) (map[int]value.Primary, error) {
	frameSet := WindowFrameSet(partition, expr.AnalyticClause)
	list := make(map[int]value.Primary, len(partition))

	valueCache := make(map[int]value.Primary, len(partition))

	for _, frame := range frameSet {
		var val value.Primary = value.NewNull()
		count := 0

		for i := frame.Low; i <= frame.High; i++ {
			if i < 0 || len(partition) <= i {
				continue
			}

			recordIdx := partition[i]
			if v, ok := valueCache[recordIdx]; ok {
				val = v
			} else {
				filter.records[0].recordIndex = recordIdx
				p, err := filter.Evaluate(ctx, expr.Args[0])
				if err != nil {
					return nil, err
				}
				valueCache[recordIdx] = p
				val = p
			}
			if expr.IgnoreNulls && value.IsNull(val) {
				continue
			}

			count++
			if count == n {
				break
			}
		}

		for _, idx := range frame.Records {
			list[idx] = val
		}
	}

	return list, nil
}

type Lag struct{}

func (fn Lag) CheckArgsLen(expr parser.AnalyticFunction) error {
	return CheckArgsLen(expr, []int{1, 3})
}

func (fn Lag) Execute(ctx context.Context, filter *Filter, partition Partition, expr parser.AnalyticFunction) (map[int]value.Primary, error) {
	return setLag(ctx, filter, partition, expr)
}

type Lead struct{}

func (fn Lead) CheckArgsLen(expr parser.AnalyticFunction) error {
	return CheckArgsLen(expr, []int{1, 3})
}

func (fn Lead) Execute(ctx context.Context, filter *Filter, partition Partition, expr parser.AnalyticFunction) (map[int]value.Primary, error) {
	partition.Reverse()
	return setLag(ctx, filter, partition, expr)
}

func setLag(ctx context.Context, filter *Filter, partition Partition, expr parser.AnalyticFunction) (map[int]value.Primary, error) {
	argsFilter := filter.CreateNode()
	argsFilter.records = nil

	offset := 1
	if 1 < len(expr.Args) {
		p, err := argsFilter.Evaluate(ctx, expr.Args[1])
		if err != nil {
			return nil, NewFunctionInvalidArgumentError(expr, expr.Name, "the second argument must be an integer")
		}
		i := value.ToInteger(p)
		if value.IsNull(i) {
			return nil, NewFunctionInvalidArgumentError(expr, expr.Name, "the second argument must be an integer")
		}
		offset = int(i.(value.Integer).Raw())
	}

	var defaultValue value.Primary = value.NewNull()
	if 2 < len(expr.Args) {
		p, err := argsFilter.Evaluate(ctx, expr.Args[2])
		if err != nil {
			return nil, NewFunctionInvalidArgumentError(expr, expr.Name, "the third argument must be a primitive type")
		}
		defaultValue = p
	}

	list := make(map[int]value.Primary, len(partition))
	values := make([]value.Primary, 0)
	for _, idx := range partition {
		filter.records[0].recordIndex = idx
		p, err := filter.Evaluate(ctx, expr.Args[0])
		if err != nil {
			return nil, err
		}

		values = append(values, p)

		lagIdx := len(values) - 1 - offset
		val := defaultValue
		if 0 <= lagIdx && lagIdx < len(values) {
			for i := lagIdx; i >= 0; i-- {
				if expr.IgnoreNulls && value.IsNull(values[i]) {
					continue
				}
				val = values[i]
				break
			}
		}
		list[idx] = val
	}

	return list, nil
}

type AnalyticListAgg struct{}

func (fn AnalyticListAgg) CheckArgsLen(expr parser.AnalyticFunction) error {
	return CheckArgsLen(expr, []int{1, 2})
}

func (fn AnalyticListAgg) Execute(ctx context.Context, filter *Filter, partition Partition, expr parser.AnalyticFunction) (map[int]value.Primary, error) {
	argsFilter := filter.CreateNode()
	argsFilter.records = nil

	separator := ""
	if len(expr.Args) == 2 {
		p, err := argsFilter.Evaluate(ctx, expr.Args[1])
		if err != nil {
			return nil, NewFunctionInvalidArgumentError(expr, expr.Name, "the second argument must be a string")
		}
		s := value.ToString(p)
		if value.IsNull(s) {
			return nil, NewFunctionInvalidArgumentError(expr, expr.Name, "the second argument must be a string")
		}
		separator = s.(value.String).Raw()
	}

	values := make([]value.Primary, len(partition))
	for i, idx := range partition {
		filter.records[0].recordIndex = idx
		val, e := filter.Evaluate(ctx, expr.Args[0])
		if e != nil {
			return nil, e
		}
		values[i] = val
	}
	if expr.IsDistinct() {
		values = Distinguish(values, filter.tx.Flags)
	}

	val := ListAgg(values, separator)

	list := make(map[int]value.Primary, len(partition))
	for _, idx := range partition {
		list[idx] = val
	}

	return list, nil
}

type AnalyticJsonAgg struct{}

func (fn AnalyticJsonAgg) CheckArgsLen(expr parser.AnalyticFunction) error {
	return CheckArgsLen(expr, []int{1})
}

func (fn AnalyticJsonAgg) Execute(ctx context.Context, filter *Filter, partition Partition, expr parser.AnalyticFunction) (map[int]value.Primary, error) {
	argsFilter := filter.CreateNode()
	argsFilter.records = nil

	values := make([]value.Primary, len(partition))
	for i, idx := range partition {
		filter.records[0].recordIndex = idx
		val, e := filter.Evaluate(ctx, expr.Args[0])
		if e != nil {
			return nil, e
		}
		values[i] = val
	}
	if expr.IsDistinct() {
		values = Distinguish(values, filter.tx.Flags)
	}

	val := JsonAgg(values)

	list := make(map[int]value.Primary, len(partition))
	for _, idx := range partition {
		list[idx] = val
	}

	return list, nil
}
