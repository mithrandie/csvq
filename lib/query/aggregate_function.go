package query

import (
	"sort"
	"strings"

	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/ternary"
)

type AggregateFunction func([]parser.Primary) parser.Primary

var AggregateFunctions = map[string]AggregateFunction{
	"COUNT":  Count,
	"MAX":    Max,
	"MIN":    Min,
	"SUM":    Sum,
	"AVG":    Avg,
	"MEDIAN": Median,
}

func Count(list []parser.Primary) parser.Primary {
	var count int64
	for _, v := range list {
		if !parser.IsNull(v) {
			count++
		}
	}

	return parser.NewInteger(count)
}

func Max(list []parser.Primary) parser.Primary {
	var result parser.Primary
	result = parser.NewNull()

	for _, v := range list {
		if parser.IsNull(v) {
			continue
		}

		if parser.IsNull(result) {
			result = v
			continue
		}

		if GreaterThan(v, result) == ternary.TRUE {
			result = v
		}
	}

	return result
}

func Min(list []parser.Primary) parser.Primary {
	var result parser.Primary
	result = parser.NewNull()

	for _, v := range list {
		if parser.IsNull(v) {
			continue
		}

		if parser.IsNull(result) {
			result = v
			continue
		}

		if LessThan(v, result) == ternary.TRUE {
			result = v
		}
	}

	return result
}

func Sum(list []parser.Primary) parser.Primary {
	var sum float64
	var count int

	for _, v := range list {
		f := parser.PrimaryToFloat(v)
		if parser.IsNull(f) {
			continue
		}

		sum += f.(parser.Float).Value()
		count++
	}

	if count < 1 {
		return parser.NewNull()
	}
	return parser.Float64ToPrimary(sum)
}

func Avg(list []parser.Primary) parser.Primary {
	var sum float64
	var count int

	for _, v := range list {
		f := parser.PrimaryToFloat(v)
		if parser.IsNull(f) {
			continue
		}

		sum += f.(parser.Float).Value()
		count++
	}

	if count < 1 {
		return parser.NewNull()
	}

	avg := sum / float64(count)
	return parser.Float64ToPrimary(avg)
}

func Median(list []parser.Primary) parser.Primary {
	var values []float64

	for _, v := range list {
		if f := parser.PrimaryToFloat(v); !parser.IsNull(f) {
			values = append(values, f.(parser.Float).Value())
			continue
		}
		if d := parser.PrimaryToDatetime(v); !parser.IsNull(d) {
			values = append(values, float64(d.(parser.Datetime).Value().UnixNano())/float64(1000000000))
			continue
		}
	}

	if len(values) < 1 {
		return parser.NewNull()
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
	return parser.Float64ToPrimary(median)
}

func ListAgg(list []parser.Primary, separator string) parser.Primary {
	strlist := []string{}
	for _, v := range list {
		s := parser.PrimaryToString(v)
		if parser.IsNull(s) {
			continue
		}
		strlist = append(strlist, s.(parser.String).Value())
	}

	if len(strlist) < 1 {
		return parser.NewNull()
	}

	return parser.NewString(strings.Join(strlist, separator))
}
