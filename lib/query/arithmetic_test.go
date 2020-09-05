package query

import (
	"reflect"
	"testing"

	"github.com/mithrandie/csvq/lib/value"
)

var calculateTests = []struct {
	LHS      value.Primary
	RHS      value.Primary
	Operator int
	Result   value.Primary
}{
	{
		LHS:      value.NewString("9"),
		RHS:      value.NewNull(),
		Operator: '+',
		Result:   value.NewNull(),
	},
	{
		LHS:      value.NewNull(),
		RHS:      value.NewString("9"),
		Operator: '+',
		Result:   value.NewNull(),
	},
	{
		LHS:      value.NewString("9"),
		RHS:      value.NewString("2"),
		Operator: '+',
		Result:   value.NewInteger(11),
	},
	{
		LHS:      value.NewString("9"),
		RHS:      value.NewString("2"),
		Operator: '-',
		Result:   value.NewInteger(7),
	},
	{
		LHS:      value.NewString("9"),
		RHS:      value.NewString("2"),
		Operator: '*',
		Result:   value.NewInteger(18),
	},
	{
		LHS:      value.NewString("9"),
		RHS:      value.NewString("2"),
		Operator: '%',
		Result:   value.NewInteger(1),
	},
	{
		LHS:      value.NewString("9.5"),
		RHS:      value.NewString("2"),
		Operator: '+',
		Result:   value.NewFloat(11.5),
	},
	{
		LHS:      value.NewString("9.5"),
		RHS:      value.NewString("2"),
		Operator: '-',
		Result:   value.NewFloat(7.5),
	},
	{
		LHS:      value.NewString("9.5"),
		RHS:      value.NewString("2"),
		Operator: '*',
		Result:   value.NewInteger(19),
	},
	{
		LHS:      value.NewString("9"),
		RHS:      value.NewString("2"),
		Operator: '/',
		Result:   value.NewFloat(4.5),
	},
	{
		LHS:      value.NewString("8.5"),
		RHS:      value.NewString("2"),
		Operator: '%',
		Result:   value.NewFloat(0.5),
	},
}

func TestCalculate(t *testing.T) {
	for _, v := range calculateTests {
		r := Calculate(v.LHS, v.RHS, v.Operator)
		if !reflect.DeepEqual(r, v.Result) {
			t.Errorf("result = %s, want %s for (%s %s %s)", r, v.Result, v.LHS, string(rune(v.Operator)), v.RHS)
		}
	}
}
