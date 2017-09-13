package query

import (
	"sort"
	"strings"

	"github.com/mithrandie/csvq/lib/value"

	"github.com/mithrandie/ternary"
)

type AggregateFunction func([]value.Primary) value.Primary

var AggregateFunctions = map[string]AggregateFunction{
	"COUNT":  Count,
	"MAX":    Max,
	"MIN":    Min,
	"SUM":    Sum,
	"AVG":    Avg,
	"MEDIAN": Median,
}

func Count(list []value.Primary) value.Primary {
	var count int64
	for _, v := range list {
		if !value.IsNull(v) {
			count++
		}
	}

	return value.NewInteger(count)
}

func Max(list []value.Primary) value.Primary {
	var result value.Primary
	result = value.NewNull()

	for _, v := range list {
		if value.IsNull(v) {
			continue
		}

		if value.IsNull(result) {
			result = v
			continue
		}

		if value.Greater(v, result) == ternary.TRUE {
			result = v
		}
	}

	return result
}

func Min(list []value.Primary) value.Primary {
	var result value.Primary
	result = value.NewNull()

	for _, v := range list {
		if value.IsNull(v) {
			continue
		}

		if value.IsNull(result) {
			result = v
			continue
		}

		if value.Less(v, result) == ternary.TRUE {
			result = v
		}
	}

	return result
}

func Sum(list []value.Primary) value.Primary {
	var sum float64
	var count int

	for _, v := range list {
		f := value.ToFloat(v)
		if value.IsNull(f) {
			continue
		}

		sum += f.(value.Float).Raw()
		count++
	}

	if count < 1 {
		return value.NewNull()
	}
	return value.ParseFloat64(sum)
}

func Avg(list []value.Primary) value.Primary {
	var sum float64
	var count int

	for _, v := range list {
		f := value.ToFloat(v)
		if value.IsNull(f) {
			continue
		}

		sum += f.(value.Float).Raw()
		count++
	}

	if count < 1 {
		return value.NewNull()
	}

	avg := sum / float64(count)
	return value.ParseFloat64(avg)
}

func Median(list []value.Primary) value.Primary {
	var values []float64

	for _, v := range list {
		if f := value.ToFloat(v); !value.IsNull(f) {
			values = append(values, f.(value.Float).Raw())
			continue
		}
		if d := value.ToDatetime(v); !value.IsNull(d) {
			values = append(values, float64(d.(value.Datetime).Raw().UnixNano())/1e9)
			continue
		}
	}

	if len(values) < 1 {
		return value.NewNull()
	}

	sort.Float64s(values)

	var median float64
	if len(values)%2 == 1 {
		idx := ((len(values) + 1) / 2) - 1
		median = values[idx]
	} else {
		idx := (len(values) / 2) - 1
		median = (values[idx] + values[idx+1]) / float64(2)
	}
	return value.ParseFloat64(median)
}

func ListAgg(list []value.Primary, separator string) value.Primary {
	strlist := []string{}
	for _, v := range list {
		s := value.ToString(v)
		if value.IsNull(s) {
			continue
		}
		strlist = append(strlist, s.(value.String).Raw())
	}

	if len(strlist) < 1 {
		return value.NewNull()
	}

	return value.NewString(strings.Join(strlist, separator))
}
