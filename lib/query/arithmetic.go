package query

import (
	"math"

	"github.com/mithrandie/csvq/lib/value"
)

func Calculate(p1 value.Primary, p2 value.Primary, operator int) value.Primary {
	if operator != '/' {
		if pi1 := value.PrimaryToInteger(p1); !value.IsNull(pi1) {
			if pi2 := value.PrimaryToInteger(p2); !value.IsNull(pi2) {
				return calculateInteger(pi1.(value.Integer).Raw(), pi2.(value.Integer).Raw(), operator)
			}
		}
	}

	pf1 := value.PrimaryToFloat(p1)
	pf2 := value.PrimaryToFloat(p2)

	if value.IsNull(pf1) || value.IsNull(pf2) {
		return value.NewNull()
	}

	f1 := pf1.(value.Float).Raw()
	f2 := pf2.(value.Float).Raw()

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

	return value.Float64ToPrimary(result)
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
