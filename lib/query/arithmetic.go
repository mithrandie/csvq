package query

import (
	"math"

	"github.com/mithrandie/csvq/lib/parser"
)

func Calculate(p1 parser.Primary, p2 parser.Primary, operator int) parser.Primary {
	if operator != '/' {
		if pi1 := parser.PrimaryToInteger(p1); !parser.IsNull(pi1) {
			if pi2 := parser.PrimaryToInteger(p2); !parser.IsNull(pi2) {
				return calculateInteger(pi1.(parser.Integer).Value(), pi2.(parser.Integer).Value(), operator)
			}
		}
	}

	pf1 := parser.PrimaryToFloat(p1)
	pf2 := parser.PrimaryToFloat(p2)

	if parser.IsNull(pf1) || parser.IsNull(pf2) {
		return parser.NewNull()
	}

	f1 := pf1.(parser.Float).Value()
	f2 := pf2.(parser.Float).Value()

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

	return parser.Float64ToPrimary(result)
}

func calculateInteger(i1 int64, i2 int64, operator int) parser.Primary {
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

	return parser.NewInteger(result)
}
