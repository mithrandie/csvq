package query

import (
	"math"
	"reflect"
	"testing"

	"github.com/mithrandie/csvq/lib/value"
)

var calculateTests = []struct {
	LHS      value.Primary
	RHS      value.Primary
	Operator int
	Result   value.Primary
	Error    string
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
		Result:   value.NewFloat(19),
	},
	{
		LHS:      value.NewString("9"),
		RHS:      value.NewString("2"),
		Operator: '/',
		Result:   value.NewInteger(4),
	},
	{
		LHS:      value.NewString("8.5"),
		RHS:      value.NewString("2"),
		Operator: '%',
		Result:   value.NewFloat(0.5),
	},
	{
		LHS:      value.NewString("8"),
		RHS:      value.NewString("0"),
		Operator: '/',
		Error:    "integer devided by zero",
	},
	{
		LHS:      value.NewString("8"),
		RHS:      value.NewString("0"),
		Operator: '%',
		Error:    "integer devided by zero",
	},
	{
		LHS:      value.NewFloat(math.Inf(1)),
		RHS:      value.NewFloat(100),
		Operator: '+',
		Result:   value.NewFloat(math.Inf(1)),
	},
	{
		LHS:      value.NewFloat(math.Inf(1)),
		RHS:      value.NewFloat(math.Inf(-1)),
		Operator: '-',
		Result:   value.NewFloat(math.Inf(1)),
	},
	{
		LHS:      value.NewFloat(math.Inf(1)),
		RHS:      value.NewFloat(math.NaN()),
		Operator: '+',
		Result:   value.NewFloat(math.NaN()),
	},
}

func TestCalculate(t *testing.T) {
	for _, v := range calculateTests {
		r, err := Calculate(v.LHS, v.RHS, v.Operator)

		if err != nil {
			if len(v.Error) < 1 {
				t.Errorf("unexpected error %q for (%s %s %s)", err, v.LHS, string(rune(v.Operator)), v.RHS)
			} else if err.Error() != v.Error {
				t.Errorf("error %q, want error %q for (%s %s %s)", err.Error(), v.Error, v.LHS, string(rune(v.Operator)), v.RHS)
			}
			continue
		}
		if 0 < len(v.Error) {
			t.Errorf("no error, want error %q for (%s %s %s)", v.Error, v.LHS, string(rune(v.Operator)), v.RHS)
			continue
		}

		if !reflect.DeepEqual(r, v.Result) {
			if expectF, ok := r.(*value.Float); ok && math.IsNaN(expectF.Raw()) {
				if retF, ok := r.(*value.Float); ok && math.IsNaN(retF.Raw()) {
					continue
				}
			}

			t.Errorf("result = %s, want %s for (%s %s %s)", r, v.Result, v.LHS, string(rune(v.Operator)), v.RHS)
		}
	}
}
