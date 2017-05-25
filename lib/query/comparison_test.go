package query

import (
	"testing"

	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/ternary"
)

func TestComparisonResult_String(t *testing.T) {
	if EQUAL.String() != "EQUAL" {
		t.Errorf("string = %s, want %s for %s.String()", EQUAL.String(), "EQUAL", EQUAL)
	}
}

var compareCombinedlyTests = []struct {
	LHS    parser.Primary
	RHS    parser.Primary
	Result ComparisonResult
}{
	{
		LHS:    parser.NewInteger(1),
		RHS:    parser.NewNull(),
		Result: INCOMMENSURABLE,
	},
	{
		LHS:    parser.NewInteger(1),
		RHS:    parser.NewInteger(1),
		Result: EQUAL,
	},
	{
		LHS:    parser.NewInteger(1),
		RHS:    parser.NewInteger(2),
		Result: LESS,
	},
	{
		LHS:    parser.NewInteger(2),
		RHS:    parser.NewInteger(1),
		Result: GREATER,
	},
	{
		LHS:    parser.NewDatetimeFromString("2006-01-02T15:04:05-07:00"),
		RHS:    parser.NewDatetimeFromString("2006-01-02T15:04:05-07:00"),
		Result: EQUAL,
	},
	{
		LHS:    parser.NewDatetimeFromString("2006-01-02T15:04:05-07:00"),
		RHS:    parser.NewDatetimeFromString("2006-02-02T15:04:05-07:00"),
		Result: LESS,
	},
	{
		LHS:    parser.NewDatetimeFromString("2006-02-02T15:04:05-07:00"),
		RHS:    parser.NewDatetimeFromString("2006-01-02T15:04:05-07:00"),
		Result: GREATER,
	},
	{
		LHS:    parser.NewBoolean(true),
		RHS:    parser.NewBoolean(true),
		Result: EQUAL,
	},
	{
		LHS:    parser.NewBoolean(true),
		RHS:    parser.NewBoolean(false),
		Result: NOT_EQUAL,
	},
	{
		LHS:    parser.NewString("A"),
		RHS:    parser.NewString("a"),
		Result: EQUAL,
	},
	{
		LHS:    parser.NewString("A"),
		RHS:    parser.NewString("B"),
		Result: LESS,
	},
	{
		LHS:    parser.NewString("B"),
		RHS:    parser.NewString("A"),
		Result: GREATER,
	},
	{
		LHS:    parser.NewString("B"),
		RHS:    parser.NewTernary(ternary.TRUE),
		Result: INCOMMENSURABLE,
	},
}

func TestCompareCombinedly(t *testing.T) {
	for _, v := range compareCombinedlyTests {
		r := CompareCombinedly(v.LHS, v.RHS)
		if r != v.Result {
			t.Errorf("result = %s, want %s for comparison with %s and %s", r, v.Result, v.LHS, v.RHS)
		}
	}
}

var compareTests = []struct {
	LHS    parser.Primary
	RHS    parser.Primary
	Op     string
	Result ternary.Value
}{
	{
		LHS:    parser.NewInteger(1),
		RHS:    parser.NewInteger(2),
		Op:     "=",
		Result: ternary.FALSE,
	},
	{
		LHS:    parser.NewInteger(1),
		RHS:    parser.NewInteger(1),
		Op:     "=",
		Result: ternary.TRUE,
	},
	{
		LHS:    parser.NewInteger(1),
		RHS:    parser.NewNull(),
		Op:     "=",
		Result: ternary.UNKNOWN,
	},
	{
		LHS:    parser.NewInteger(1),
		RHS:    parser.NewInteger(2),
		Op:     ">",
		Result: ternary.FALSE,
	},
	{
		LHS:    parser.NewInteger(2),
		RHS:    parser.NewInteger(1),
		Op:     ">",
		Result: ternary.TRUE,
	},
	{
		LHS:    parser.NewInteger(1),
		RHS:    parser.NewNull(),
		Op:     ">",
		Result: ternary.UNKNOWN,
	},
	{
		LHS:    parser.NewInteger(2),
		RHS:    parser.NewInteger(1),
		Op:     "<",
		Result: ternary.FALSE,
	},
	{
		LHS:    parser.NewInteger(1),
		RHS:    parser.NewInteger(2),
		Op:     "<",
		Result: ternary.TRUE,
	},
	{
		LHS:    parser.NewInteger(1),
		RHS:    parser.NewNull(),
		Op:     "<",
		Result: ternary.UNKNOWN,
	},
	{
		LHS:    parser.NewInteger(1),
		RHS:    parser.NewInteger(2),
		Op:     ">=",
		Result: ternary.FALSE,
	},
	{
		LHS:    parser.NewInteger(2),
		RHS:    parser.NewInteger(2),
		Op:     ">=",
		Result: ternary.TRUE,
	},
	{
		LHS:    parser.NewInteger(1),
		RHS:    parser.NewNull(),
		Op:     ">=",
		Result: ternary.UNKNOWN,
	},
	{
		LHS:    parser.NewInteger(2),
		RHS:    parser.NewInteger(1),
		Op:     "<=",
		Result: ternary.FALSE,
	},
	{
		LHS:    parser.NewInteger(2),
		RHS:    parser.NewInteger(2),
		Op:     "<=",
		Result: ternary.TRUE,
	},
	{
		LHS:    parser.NewInteger(1),
		RHS:    parser.NewNull(),
		Op:     "<=",
		Result: ternary.UNKNOWN,
	},
	{
		LHS:    parser.NewInteger(2),
		RHS:    parser.NewInteger(2),
		Op:     "<>",
		Result: ternary.FALSE,
	},
	{
		LHS:    parser.NewInteger(2),
		RHS:    parser.NewInteger(1),
		Op:     "<>",
		Result: ternary.TRUE,
	},
	{
		LHS:    parser.NewInteger(1),
		RHS:    parser.NewNull(),
		Op:     "<>",
		Result: ternary.UNKNOWN,
	},
}

func TestCompare(t *testing.T) {
	for _, v := range compareTests {
		r := Compare(v.LHS, v.RHS, v.Op)
		if r != v.Result {
			t.Errorf("result = %s, want %s for (%s %s %s)", r, v.Result, v.LHS, v.Op, v.RHS)
		}
	}
}

var equivalentToTests = []struct {
	LHS    parser.Primary
	RHS    parser.Primary
	Result ternary.Value
}{
	{
		LHS:    parser.NewNull(),
		RHS:    parser.NewNull(),
		Result: ternary.TRUE,
	},
	{
		LHS:    parser.NewInteger(1),
		RHS:    parser.NewNull(),
		Result: ternary.UNKNOWN,
	},
}

func TestEquivalentTo(t *testing.T) {
	for _, v := range equivalentToTests {
		r := EquivalentTo(v.LHS, v.RHS)
		if r != v.Result {
			t.Errorf("result = %s, want %s for (%s is equivalent to %s)", r, v.Result, v.LHS, v.RHS)
		}
	}
}

var isTests = []struct {
	LHS    parser.Primary
	RHS    parser.Primary
	Result ternary.Value
}{
	{
		LHS:    parser.NewBoolean(true),
		RHS:    parser.NewTernary(ternary.TRUE),
		Result: ternary.TRUE,
	},
	{
		LHS:    parser.NewNull(),
		RHS:    parser.NewNull(),
		Result: ternary.TRUE,
	},
	{
		LHS:    parser.NewString("foo"),
		RHS:    parser.NewNull(),
		Result: ternary.FALSE,
	},
}

func TestIs(t *testing.T) {
	for _, v := range isTests {
		r := Is(v.LHS, v.RHS)
		if r != v.Result {
			t.Errorf("result = %s, want %s for (%s is %s)", r, v.Result, v.LHS, v.RHS)
		}
	}
}

var betweenTests = []struct {
	LHS    parser.Primary
	Low    parser.Primary
	High   parser.Primary
	Result ternary.Value
}{
	{
		LHS:    parser.NewInteger(2),
		Low:    parser.NewInteger(1),
		High:   parser.NewInteger(3),
		Result: ternary.TRUE,
	},
	{
		LHS:    parser.NewInteger(2),
		Low:    parser.NewInteger(3),
		High:   parser.NewInteger(1),
		Result: ternary.FALSE,
	},
	{
		LHS:    parser.NewInteger(2),
		Low:    parser.NewInteger(2),
		High:   parser.NewInteger(2),
		Result: ternary.TRUE,
	},
	{
		LHS:    parser.NewInteger(2),
		Low:    parser.NewInteger(1),
		High:   parser.NewNull(),
		Result: ternary.UNKNOWN,
	},
}

func TestBetween(t *testing.T) {
	for _, v := range betweenTests {
		r := Between(v.LHS, v.Low, v.High)
		if r != v.Result {
			t.Errorf("result = %s, want %s for (%s between %s and %s)", r, v.Result, v.LHS, v.Low, v.High)
		}
	}
}

var likeTests = []struct {
	LHS     parser.Primary
	Pattern parser.Primary
	Result  ternary.Value
}{
	{
		LHS:     parser.NewString("str"),
		Pattern: parser.NewNull(),
		Result:  ternary.UNKNOWN,
	},
	{
		LHS:     parser.NewString("str"),
		Pattern: parser.NewBoolean(true),
		Result:  ternary.UNKNOWN,
	},
	{
		LHS:     parser.NewBoolean(true),
		Pattern: parser.NewString("str"),
		Result:  ternary.UNKNOWN,
	},
	{
		LHS:     parser.NewString("str"),
		Pattern: parser.NewString("str"),
		Result:  ternary.TRUE,
	},
	{
		LHS:     parser.NewString("str"),
		Pattern: parser.NewString(""),
		Result:  ternary.FALSE,
	},
	{
		LHS:     parser.NewString("abcdefghijk"),
		Pattern: parser.NewString("lmn"),
		Result:  ternary.FALSE,
	},
	{
		LHS:     parser.NewString("abc"),
		Pattern: parser.NewString("_____"),
		Result:  ternary.FALSE,
	},
	{
		LHS:     parser.NewString("abcde"),
		Pattern: parser.NewString("___"),
		Result:  ternary.FALSE,
	},
	{
		LHS:     parser.NewString("abcde"),
		Pattern: parser.NewString("___%__"),
		Result:  ternary.TRUE,
	},
	{
		LHS:     parser.NewString("abcde"),
		Pattern: parser.NewString("_%__"),
		Result:  ternary.TRUE,
	},
	{
		LHS:     parser.NewString("abcdefghijkabcdefghijk"),
		Pattern: parser.NewString("%def%_abc%"),
		Result:  ternary.TRUE,
	},
}

func TestLike(t *testing.T) {
	for _, v := range likeTests {
		r := Like(v.LHS, v.Pattern)
		if r != v.Result {
			t.Errorf("result = %s, want %s for (%s like %s)", r, v.Result, v.LHS, v.Pattern)
		}
	}
}

var anyTests = []struct {
	LHS      parser.Primary
	List     []parser.Primary
	Operator string
	Result   ternary.Value
}{
	{
		LHS: parser.NewInteger(3),
		List: []parser.Primary{
			parser.NewInteger(1),
			parser.NewInteger(2),
			parser.NewInteger(3),
		},
		Operator: "=",
		Result:   ternary.TRUE,
	},
	{
		LHS: parser.NewInteger(5),
		List: []parser.Primary{
			parser.NewInteger(1),
			parser.NewNull(),
			parser.NewInteger(3),
		},
		Operator: "=",
		Result:   ternary.UNKNOWN,
	},
	{
		LHS: parser.NewInteger(5),
		List: []parser.Primary{
			parser.NewInteger(1),
			parser.NewInteger(2),
			parser.NewInteger(3),
		},
		Operator: "=",
		Result:   ternary.FALSE,
	},
	{
		LHS:      parser.NewInteger(5),
		List:     nil,
		Operator: "=",
		Result:   ternary.FALSE,
	},
}

func TestAny(t *testing.T) {
	for _, v := range anyTests {
		r := Any(v.LHS, v.List, v.Operator)
		if r != v.Result {
			t.Errorf("result = %s, want %s for (%s %s any (%s))", r, v.Result, v.LHS, v.Operator, v.List)
		}
	}
}

var allTests = []struct {
	LHS      parser.Primary
	List     []parser.Primary
	Operator string
	Result   ternary.Value
}{
	{
		LHS: parser.NewInteger(5),
		List: []parser.Primary{
			parser.NewInteger(1),
			parser.NewInteger(2),
			parser.NewInteger(3),
		},
		Operator: ">",
		Result:   ternary.TRUE,
	},
	{
		LHS: parser.NewInteger(5),
		List: []parser.Primary{
			parser.NewInteger(1),
			parser.NewNull(),
			parser.NewInteger(3),
		},
		Operator: ">",
		Result:   ternary.UNKNOWN,
	},
	{
		LHS: parser.NewInteger(3),
		List: []parser.Primary{
			parser.NewInteger(1),
			parser.NewInteger(2),
			parser.NewInteger(3),
		},
		Operator: ">",
		Result:   ternary.FALSE,
	},
	{
		LHS:      parser.NewInteger(5),
		List:     nil,
		Operator: "=",
		Result:   ternary.TRUE,
	},
}

func TestAll(t *testing.T) {
	for _, v := range allTests {
		r := All(v.LHS, v.List, v.Operator)
		if r != v.Result {
			t.Errorf("result = %s, want %s for (%s %s all (%s))", r, v.Result, v.LHS, v.Operator, v.List)
		}
	}
}
