package value

import (
	"testing"

	"github.com/mithrandie/ternary"
)

func TestComparisonResult_String(t *testing.T) {
	if IsEqual.String() != "IsEqual" {
		t.Errorf("string = %s, want %s for %s.String()", IsEqual.String(), "IsEqual", IsEqual)
	}
}

var compareCombinedlyTests = []struct {
	LHS    Primary
	RHS    Primary
	Result ComparisonResult
}{
	{
		LHS:    NewInteger(1),
		RHS:    NewNull(),
		Result: IsIncommensurable,
	},
	{
		LHS:    NewIntegerFromString("1"),
		RHS:    NewInteger(1),
		Result: IsEqual,
	},
	{
		LHS:    NewInteger(1),
		RHS:    NewInteger(2),
		Result: IsLess,
	},
	{
		LHS:    NewInteger(2),
		RHS:    NewInteger(1),
		Result: IsGreater,
	},
	{
		LHS:    NewFloatFromString("1.5"),
		RHS:    NewFloat(1.5),
		Result: IsEqual,
	},
	{
		LHS:    NewFloat(1.5),
		RHS:    NewFloat(2.0),
		Result: IsLess,
	},
	{
		LHS:    NewFloat(1.5),
		RHS:    NewFloat(1.0),
		Result: IsGreater,
	},
	{
		LHS:    NewDatetimeFromString("2006-01-02T15:04:05-07:00", nil),
		RHS:    NewDatetimeFromString("2006-01-02T15:04:05-07:00", nil),
		Result: IsEqual,
	},
	{
		LHS:    NewDatetimeFromString("2006-01-02T15:04:05-07:00", nil),
		RHS:    NewDatetimeFromString("2006-02-02T15:04:05-07:00", nil),
		Result: IsLess,
	},
	{
		LHS:    NewDatetimeFromString("2006-02-02T15:04:05-07:00", nil),
		RHS:    NewDatetimeFromString("2006-01-02T15:04:05-07:00", nil),
		Result: IsGreater,
	},
	{
		LHS:    NewDatetimeFromString("2006-02-02T15:04:05-07:00", nil),
		RHS:    NewString("abc"),
		Result: IsIncommensurable,
	},
	{
		LHS:    NewBoolean(true),
		RHS:    NewBoolean(true),
		Result: IsBoolEqual,
	},
	{
		LHS:    NewBoolean(true),
		RHS:    NewBoolean(false),
		Result: IsNotEqual,
	},
	{
		LHS:    NewString(" A "),
		RHS:    NewString("a"),
		Result: IsEqual,
	},
	{
		LHS:    NewString("A"),
		RHS:    NewString("B"),
		Result: IsLess,
	},
	{
		LHS:    NewString("B"),
		RHS:    NewString("A"),
		Result: IsGreater,
	},
	{
		LHS:    NewString("B"),
		RHS:    NewTernaryFromString("true"),
		Result: IsIncommensurable,
	},
	{
		LHS:    NewString("1"),
		RHS:    NewString("A"),
		Result: IsLess,
	},
}

func TestCompareCombinedly(t *testing.T) {
	for _, v := range compareCombinedlyTests {
		r := CompareCombinedly(v.LHS, v.RHS, nil)
		if r != v.Result {
			t.Errorf("result = %s, want %s for comparison with %s and %s", r, v.Result, v.LHS, v.RHS)
		}
	}
}

var identicalTests = []struct {
	LHS    Primary
	RHS    Primary
	Result ternary.Value
}{
	{
		LHS:    NewNull(),
		RHS:    NewString("R"),
		Result: ternary.UNKNOWN,
	},
	{
		LHS:    NewTernary(ternary.UNKNOWN),
		RHS:    NewString("R"),
		Result: ternary.UNKNOWN,
	},
	{
		LHS:    NewString("L"),
		RHS:    NewNull(),
		Result: ternary.UNKNOWN,
	},
	{
		LHS:    NewString("L"),
		RHS:    NewTernary(ternary.UNKNOWN),
		Result: ternary.UNKNOWN,
	},
	{
		LHS:    NewInteger(1),
		RHS:    NewInteger(1),
		Result: ternary.TRUE,
	},
	{
		LHS:    NewFloat(1),
		RHS:    NewFloat(1),
		Result: ternary.TRUE,
	},
	{
		LHS:    NewDatetimeFromString("2006-02-02T15:04:05-07:00", nil),
		RHS:    NewDatetimeFromString("2006-02-02T15:04:05-07:00", nil),
		Result: ternary.TRUE,
	},
	{
		LHS:    NewBoolean(true),
		RHS:    NewBoolean(true),
		Result: ternary.TRUE,
	},
	{
		LHS:    NewTernary(ternary.TRUE),
		RHS:    NewTernary(ternary.TRUE),
		Result: ternary.TRUE,
	},
	{
		LHS:    NewString("L"),
		RHS:    NewString("R"),
		Result: ternary.FALSE,
	},
	{
		LHS:    NewString("A"),
		RHS:    NewString("A"),
		Result: ternary.TRUE,
	},
	{
		LHS:    NewInteger(1),
		RHS:    NewFloat(1),
		Result: ternary.FALSE,
	},
}

func TestIdentical(t *testing.T) {
	for _, v := range identicalTests {
		r := Identical(v.LHS, v.RHS)
		if r != v.Result {
			t.Errorf("result = %s, want %s for comparison with %s and %s", r, v.Result, v.LHS, v.RHS)
		}
	}
}

var compareTests = []struct {
	LHS    Primary
	RHS    Primary
	Op     string
	Result ternary.Value
}{
	{
		LHS:    NewInteger(1),
		RHS:    NewInteger(2),
		Op:     "=",
		Result: ternary.FALSE,
	},
	{
		LHS:    NewInteger(1),
		RHS:    NewInteger(1),
		Op:     "=",
		Result: ternary.TRUE,
	},
	{
		LHS:    NewInteger(1),
		RHS:    NewNull(),
		Op:     "=",
		Result: ternary.UNKNOWN,
	},
	{
		LHS:    NewString("0001"),
		RHS:    NewInteger(1),
		Op:     "=",
		Result: ternary.TRUE,
	},
	{
		LHS:    NewString("0001"),
		RHS:    NewInteger(1),
		Op:     "==",
		Result: ternary.FALSE,
	},
	{
		LHS:    NewInteger(1),
		RHS:    NewInteger(2),
		Op:     ">",
		Result: ternary.FALSE,
	},
	{
		LHS:    NewInteger(2),
		RHS:    NewInteger(1),
		Op:     ">",
		Result: ternary.TRUE,
	},
	{
		LHS:    NewInteger(1),
		RHS:    NewNull(),
		Op:     ">",
		Result: ternary.UNKNOWN,
	},
	{
		LHS:    NewInteger(2),
		RHS:    NewInteger(1),
		Op:     "<",
		Result: ternary.FALSE,
	},
	{
		LHS:    NewInteger(1),
		RHS:    NewInteger(2),
		Op:     "<",
		Result: ternary.TRUE,
	},
	{
		LHS:    NewInteger(1),
		RHS:    NewNull(),
		Op:     "<",
		Result: ternary.UNKNOWN,
	},
	{
		LHS:    NewInteger(1),
		RHS:    NewInteger(2),
		Op:     ">=",
		Result: ternary.FALSE,
	},
	{
		LHS:    NewInteger(2),
		RHS:    NewInteger(2),
		Op:     ">=",
		Result: ternary.TRUE,
	},
	{
		LHS:    NewInteger(1),
		RHS:    NewNull(),
		Op:     ">=",
		Result: ternary.UNKNOWN,
	},
	{
		LHS:    NewInteger(2),
		RHS:    NewInteger(1),
		Op:     "<=",
		Result: ternary.FALSE,
	},
	{
		LHS:    NewInteger(2),
		RHS:    NewInteger(2),
		Op:     "<=",
		Result: ternary.TRUE,
	},
	{
		LHS:    NewInteger(1),
		RHS:    NewNull(),
		Op:     "<=",
		Result: ternary.UNKNOWN,
	},
	{
		LHS:    NewInteger(2),
		RHS:    NewInteger(2),
		Op:     "<>",
		Result: ternary.FALSE,
	},
	{
		LHS:    NewInteger(2),
		RHS:    NewInteger(1),
		Op:     "<>",
		Result: ternary.TRUE,
	},
	{
		LHS:    NewInteger(1),
		RHS:    NewNull(),
		Op:     "<>",
		Result: ternary.UNKNOWN,
	},
}

func TestCompare(t *testing.T) {
	for _, v := range compareTests {
		r := Compare(v.LHS, v.RHS, v.Op, nil)
		if r != v.Result {
			t.Errorf("result = %s, want %s for (%s %s %s)", r, v.Result, v.LHS, v.Op, v.RHS)
		}
	}
}

var compareRowValuesTests = []struct {
	LHS    RowValue
	RHS    RowValue
	Op     string
	Result ternary.Value
	Error  string
}{
	{
		LHS: RowValue{
			NewInteger(1),
			NewInteger(2),
			NewInteger(3),
		},
		RHS: RowValue{
			NewInteger(1),
			NewInteger(2),
			NewInteger(3),
		},
		Op:     "=",
		Result: ternary.TRUE,
	},
	{
		LHS: RowValue{
			NewInteger(1),
			NewInteger(2),
			NewInteger(3),
		},
		RHS: RowValue{
			NewInteger(1),
			NewNull(),
			NewInteger(3),
		},
		Op:     "=",
		Result: ternary.UNKNOWN,
	},
	{
		LHS: RowValue{
			NewInteger(1),
			NewInteger(2),
			NewInteger(3),
		},
		RHS: RowValue{
			NewInteger(1),
			NewNull(),
			NewInteger(2),
		},
		Op:     "=",
		Result: ternary.FALSE,
	},
	{
		LHS: RowValue{
			NewInteger(1),
			NewInteger(2),
			NewInteger(3),
		},
		RHS: RowValue{
			NewInteger(1),
			NewInteger(2),
			NewNull(),
		},
		Op:     "=",
		Result: ternary.UNKNOWN,
	},
	{
		LHS: RowValue{
			NewInteger(1),
			NewInteger(2),
			NewInteger(3),
		},
		RHS: RowValue{
			NewInteger(1),
			NewInteger(9),
			NewInteger(3),
		},
		Op:     "=",
		Result: ternary.FALSE,
	},
	{
		LHS: RowValue{
			NewInteger(1),
			NewInteger(2),
			NewInteger(3),
		},
		RHS: RowValue{
			NewInteger(1),
			NewInteger(2),
			NewInteger(3),
		},
		Op:     "==",
		Result: ternary.TRUE,
	},
	{
		LHS: RowValue{
			NewInteger(1),
			NewInteger(2),
			NewInteger(3),
		},
		RHS: RowValue{
			NewInteger(1),
			NewNull(),
			NewInteger(2),
		},
		Op:     "==",
		Result: ternary.FALSE,
	},
	{
		LHS: RowValue{
			NewInteger(1),
			NewInteger(2),
			NewInteger(3),
		},
		RHS: RowValue{
			NewInteger(1),
			NewNull(),
			NewInteger(3),
		},
		Op:     "==",
		Result: ternary.UNKNOWN,
	},
	{
		LHS: RowValue{
			NewInteger(1),
			NewInteger(2),
			NewInteger(3),
		},
		RHS: RowValue{
			NewInteger(1),
			NewInteger(0),
			NewInteger(3),
		},
		Op:     "<>",
		Result: ternary.TRUE,
	},
	{
		LHS: RowValue{
			NewInteger(1),
			NewInteger(2),
			NewInteger(3),
		},
		RHS: RowValue{
			NewInteger(1),
			NewNull(),
			NewInteger(2),
		},
		Op:     "<>",
		Result: ternary.TRUE,
	},
	{
		LHS: RowValue{
			NewInteger(1),
			NewInteger(2),
			NewInteger(3),
		},
		RHS: RowValue{
			NewInteger(1),
			NewNull(),
			NewInteger(3),
		},
		Op:     "<>",
		Result: ternary.UNKNOWN,
	},
	{
		LHS: RowValue{
			NewInteger(1),
			NewInteger(2),
			NewInteger(3),
		},
		RHS: RowValue{
			NewInteger(1),
			NewInteger(2),
			NewInteger(3),
		},
		Op:     "!=",
		Result: ternary.FALSE,
	},
	{
		LHS: RowValue{
			NewInteger(1),
			NewInteger(2),
			NewInteger(3),
		},
		RHS: RowValue{
			NewInteger(1),
			NewInteger(2),
			NewInteger(2),
		},
		Op:     ">",
		Result: ternary.TRUE,
	},
	{
		LHS: RowValue{
			NewInteger(1),
			NewInteger(2),
			NewInteger(3),
		},
		RHS: RowValue{
			NewInteger(1),
			NewInteger(2),
			NewInteger(4),
		},
		Op:     ">",
		Result: ternary.FALSE,
	},
	{
		LHS: RowValue{
			NewInteger(1),
			NewInteger(2),
			NewInteger(3),
		},
		RHS: RowValue{
			NewInteger(1),
			NewInteger(2),
			NewInteger(3),
		},
		Op:     ">",
		Result: ternary.FALSE,
	},
	{
		LHS: RowValue{
			NewInteger(1),
			NewInteger(2),
			NewInteger(3),
		},
		RHS: RowValue{
			NewInteger(1),
			NewInteger(2),
			NewInteger(3),
		},
		Op:     ">=",
		Result: ternary.TRUE,
	},
	{
		LHS: RowValue{
			NewInteger(1),
			NewBoolean(true),
			NewInteger(3),
		},
		RHS: RowValue{
			NewInteger(1),
			NewBoolean(false),
			NewInteger(2),
		},
		Op:     ">",
		Result: ternary.UNKNOWN,
	},
	{
		LHS: RowValue{
			NewInteger(1),
			NewInteger(2),
			NewInteger(3),
		},
		RHS: RowValue{
			NewInteger(1),
			NewInteger(2),
			NewInteger(4),
		},
		Op:     "<",
		Result: ternary.TRUE,
	},
	{
		LHS: RowValue{
			NewInteger(1),
			NewInteger(2),
			NewInteger(3),
		},
		RHS: RowValue{
			NewInteger(1),
			NewInteger(2),
			NewInteger(2),
		},
		Op:     "<",
		Result: ternary.FALSE,
	},
	{
		LHS: RowValue{
			NewInteger(1),
			NewInteger(2),
			NewInteger(3),
		},
		RHS: RowValue{
			NewInteger(1),
			NewInteger(2),
			NewInteger(3),
		},
		Op:     "<",
		Result: ternary.FALSE,
	},
	{
		LHS: RowValue{
			NewInteger(1),
			NewInteger(2),
			NewInteger(3),
		},
		RHS: RowValue{
			NewInteger(1),
			NewInteger(2),
			NewInteger(3),
		},
		Op:     "<=",
		Result: ternary.TRUE,
	},
	{
		LHS: RowValue{
			NewInteger(1),
			NewInteger(2),
			NewInteger(3),
		},
		RHS:    RowValue(nil),
		Op:     "=",
		Result: ternary.UNKNOWN,
	},
	{
		LHS: RowValue{
			NewInteger(1),
			NewInteger(2),
			NewInteger(3),
		},
		RHS: RowValue{
			NewInteger(1),
			NewInteger(3),
		},
		Op:    "=",
		Error: "row value length does not match",
	},
}

func TestCompareRowValues(t *testing.T) {
	for _, v := range compareRowValuesTests {
		r, err := CompareRowValues(v.LHS, v.RHS, v.Op, nil)
		if err != nil {
			if len(v.Error) < 1 {
				t.Errorf("unexpected error %q for (%s %s %s)", err, v.LHS, v.Op, v.RHS)
			} else if err.Error() != v.Error {
				t.Errorf("error %q, want error %q for (%s %s %s)", err.Error(), v.Error, v.LHS, v.Op, v.RHS)
			}
			continue
		}
		if 0 < len(v.Error) {
			t.Errorf("no error, want error %q for (%s %s %s)", v.Error, v.LHS, v.Op, v.RHS)
			continue
		}
		if r != v.Result {
			t.Errorf("result = %s, want %s for (%s %s %s)", r, v.Result, v.LHS, v.Op, v.RHS)
		}
	}
}

var equivalentToTests = []struct {
	LHS    Primary
	RHS    Primary
	Result ternary.Value
}{
	{
		LHS:    NewNull(),
		RHS:    NewNull(),
		Result: ternary.TRUE,
	},
	{
		LHS:    NewInteger(1),
		RHS:    NewNull(),
		Result: ternary.UNKNOWN,
	},
}

func TestEquivalentTo(t *testing.T) {
	for _, v := range equivalentToTests {
		r := Equivalent(v.LHS, v.RHS, nil)
		if r != v.Result {
			t.Errorf("result = %s, want %s for (%s is equivalent to %s)", r, v.Result, v.LHS, v.RHS)
		}
	}
}
