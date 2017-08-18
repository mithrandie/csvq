package query

import (
	"reflect"
	"testing"

	"github.com/mithrandie/csvq/lib/parser"
)

var calculateTests = []struct {
	LHS      parser.Primary
	RHS      parser.Primary
	Operator int
	Result   parser.Primary
}{
	{
		LHS:      parser.NewString("9"),
		RHS:      parser.NewNull(),
		Operator: '+',
		Result:   parser.NewNull(),
	},
	{
		LHS:      parser.NewString("9"),
		RHS:      parser.NewString("2"),
		Operator: '+',
		Result:   parser.NewInteger(11),
	},
	{
		LHS:      parser.NewString("9"),
		RHS:      parser.NewString("2"),
		Operator: '-',
		Result:   parser.NewInteger(7),
	},
	{
		LHS:      parser.NewString("9"),
		RHS:      parser.NewString("2"),
		Operator: '*',
		Result:   parser.NewInteger(18),
	},
	{
		LHS:      parser.NewString("9"),
		RHS:      parser.NewString("2"),
		Operator: '%',
		Result:   parser.NewInteger(1),
	},
	{
		LHS:      parser.NewString("9.5"),
		RHS:      parser.NewString("2"),
		Operator: '+',
		Result:   parser.NewFloat(11.5),
	},
	{
		LHS:      parser.NewString("9.5"),
		RHS:      parser.NewString("2"),
		Operator: '-',
		Result:   parser.NewFloat(7.5),
	},
	{
		LHS:      parser.NewString("9.5"),
		RHS:      parser.NewString("2"),
		Operator: '*',
		Result:   parser.NewInteger(19),
	},
	{
		LHS:      parser.NewString("9"),
		RHS:      parser.NewString("2"),
		Operator: '/',
		Result:   parser.NewFloat(4.5),
	},
	{
		LHS:      parser.NewString("8.5"),
		RHS:      parser.NewString("2"),
		Operator: '%',
		Result:   parser.NewFloat(0.5),
	},
}

func TestCalculate(t *testing.T) {
	for _, v := range calculateTests {
		r := Calculate(v.LHS, v.RHS, v.Operator)
		if !reflect.DeepEqual(r, v.Result) {
			t.Errorf("result = %s, want %s for (%s %s %s)", r, v.Result, v.LHS, string(v.Operator), v.RHS)
		}
	}
}
