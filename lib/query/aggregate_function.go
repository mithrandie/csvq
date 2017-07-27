package query

import (
	"strings"

	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/ternary"
)

var AggregateFunctions = map[string]func(bool, []parser.Primary) parser.Primary{
	"COUNT": Count,
	"MAX":   Max,
	"MIN":   Min,
	"SUM":   Sum,
	"AVG":   Avg,
}

func Count(distinct bool, list []parser.Primary) parser.Primary {
	if distinct {
		list = distinguish(list)
	}

	var count int64
	for _, v := range list {
		if !parser.IsNull(v) {
			count++
		}
	}

	return parser.NewInteger(count)
}

func Max(distinct bool, list []parser.Primary) parser.Primary {
	if distinct {
		list = distinguish(list)
	}

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

func Min(distinct bool, list []parser.Primary) parser.Primary {
	if distinct {
		list = distinguish(list)
	}

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

func Sum(distinct bool, list []parser.Primary) parser.Primary {
	if distinct {
		list = distinguish(list)
	}

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

func Avg(distinct bool, list []parser.Primary) parser.Primary {
	if distinct {
		list = distinguish(list)
	}

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

func ListAgg(distinct bool, list []parser.Primary, separator string) parser.Primary {
	if distinct {
		list = distinguish(list)
	}

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

func distinguish(list []parser.Primary) []parser.Primary {
	var in = func(list []parser.Primary, item parser.Primary) bool {
		for _, v := range list {
			if EquivalentTo(item, v) == ternary.TRUE {
				return true
			}
		}
		return false
	}

	distinguished := []parser.Primary{}
	for _, v := range list {
		if !in(distinguished, v) {
			distinguished = append(distinguished, v)
		}
	}
	return distinguished
}
