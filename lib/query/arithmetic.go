package query

import (
	"math"

	"github.com/mithrandie/csvq/lib/value"
)

func Calculate(p1 value.Primary, p2 value.Primary, operator int) value.Primary {
	if operator != '/' {
		if pi1 := value.ToInteger(p1); !value.IsNull(pi1) {
			if pi2 := value.ToInteger(p2); !value.IsNull(pi2) {
				ret := calculateInteger(pi1.(*value.Integer).Raw(), pi2.(*value.Integer).Raw(), operator)
				value.Discard(pi1)
				value.Discard(pi2)
				return ret
			}
			value.Discard(pi1)
		}
	}

	pf1 := value.ToFloat(p1)
	if value.IsNull(pf1) {
		return value.NewNull()
	}
	f1 := pf1.(*value.Float).Raw()
	value.Discard(pf1)

	pf2 := value.ToFloat(p2)
	if value.IsNull(pf2) {
		return value.NewNull()
	}
	f2 := pf2.(*value.Float).Raw()
	value.Discard(pf2)

	result := 0.0
	switch operator {
	case '+':
		result = f1 + f2
	case '-':
		result = f1 - f2
	case '*':
		result = f1 * f2
	case '/':
		result = f1 / f2
	case '%':
		result = math.Remainder(f1, f2)
	}

	return value.ParseFloat64(result)
}

func calculateInteger(i1 int64, i2 int64, operator int) value.Primary {
	var result int64 = 0
	switch operator {
	case '+':
		result = i1 + i2
	case '-':
		result = i1 - i2
	case '*':
		result = i1 * i2
	case '%':
		result = i1 % i2
	}

	return value.NewInteger(result)
}
