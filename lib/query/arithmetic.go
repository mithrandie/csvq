package query

import (
	"math"

	"github.com/mithrandie/csvq/lib/parser"
)

func Calculate(p1 parser.Primary, p2 parser.Primary, operator int) parser.Primary {
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
