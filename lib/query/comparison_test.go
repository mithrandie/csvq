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

var compareRowValuesTests = []struct {
	LHS    []parser.Primary
	RHS    []parser.Primary
	Op     string
	Result ternary.Value
	Error  string
}{
	{
		LHS: []parser.Primary{
			parser.NewInteger(1),
			parser.NewInteger(2),
			parser.NewInteger(3),
		},
		RHS: []parser.Primary{
			parser.NewInteger(1),
			parser.NewInteger(2),
			parser.NewInteger(3),
		},
		Op:     "=",
		Result: ternary.TRUE,
	},
	{
		LHS: []parser.Primary{
			parser.NewInteger(1),
			parser.NewInteger(2),
			parser.NewInteger(3),
		},
		RHS: []parser.Primary{
			parser.NewInteger(1),
			parser.NewNull(),
			parser.NewInteger(3),
		},
		Op:     "=",
		Result: ternary.UNKNOWN,
	},
	{
		LHS: []parser.Primary{
			parser.NewInteger(1),
			parser.NewInteger(2),
			parser.NewInteger(3),
		},
		RHS: []parser.Primary{
			parser.NewInteger(1),
			parser.NewInteger(9),
			parser.NewInteger(3),
		},
		Op:     "=",
		Result: ternary.FALSE,
	},
	{
		LHS: []parser.Primary{
			parser.NewInteger(1),
			parser.NewInteger(2),
			parser.NewInteger(3),
		},
		RHS:    []parser.Primary(nil),
		Op:     "=",
		Result: ternary.UNKNOWN,
	},
	{
		LHS: []parser.Primary{
			parser.NewInteger(1),
			parser.NewInteger(2),
			parser.NewInteger(3),
		},
		RHS: []parser.Primary{
			parser.NewInteger(1),
			parser.NewInteger(3),
		},
		Op:    "=",
		Error: "row value length does not match",
	},
}

func TestCompareRowValues(t *testing.T) {
	for _, v := range compareRowValuesTests {
		r, err := CompareRowValues(v.LHS, v.RHS, v.Op)
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
	{
		LHS:     parser.NewString("abcde"),
		Pattern: parser.NewString("abc\\_e"),
		Result:  ternary.FALSE,
	},
	{
		LHS:     parser.NewString("abc_e"),
		Pattern: parser.NewString("abc\\_e"),
		Result:  ternary.TRUE,
	},
	{
		LHS:     parser.NewString("abcde"),
		Pattern: parser.NewString("abc\\%e"),
		Result:  ternary.FALSE,
	},
	{
		LHS:     parser.NewString("a\\bc%e"),
		Pattern: parser.NewString("a\\bc\\%e"),
		Result:  ternary.TRUE,
	},
	{
		LHS:     parser.NewString("abcdef"),
		Pattern: parser.NewString("abcde\\"),
		Result:  ternary.FALSE,
	},
	{
		LHS:     parser.NewString("abcde"),
		Pattern: parser.NewString("abc"),
		Result:  ternary.FALSE,
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

var inRowValueListTests = []struct {
	LHS      []parser.Primary
	List     [][]parser.Primary
	Type     int
	Operator string
	Result   ternary.Value
	Error    string
}{
	{
		LHS: []parser.Primary{
			parser.NewInteger(1),
			parser.NewInteger(2),
			parser.NewInteger(3),
		},
		List: [][]parser.Primary{
			{
				parser.NewInteger(1),
				parser.NewInteger(2),
				parser.NewInteger(3),
			},
			{
				parser.NewInteger(4),
				parser.NewInteger(5),
				parser.NewInteger(6),
			},
		},
		Type:     parser.ANY,
		Operator: "=",
		Result:   ternary.TRUE,
	},
	{
		LHS: []parser.Primary{
			parser.NewInteger(1),
			parser.NewInteger(2),
			parser.NewInteger(3),
		},
		List: [][]parser.Primary{
			{
				parser.NewInteger(1),
				parser.NewInteger(2),
				parser.NewInteger(3),
			},
			{
				parser.NewInteger(1),
				parser.NewInteger(2),
				parser.NewInteger(3),
			},
		},
		Type:     parser.ALL,
		Operator: "=",
		Result:   ternary.TRUE,
	},
	{
		LHS: []parser.Primary{
			parser.NewInteger(1),
			parser.NewInteger(2),
			parser.NewInteger(3),
		},
		List: [][]parser.Primary{
			{
				parser.NewInteger(1),
				parser.NewInteger(2),
				parser.NewInteger(3),
			},
			{
				parser.NewInteger(1),
				parser.NewInteger(3),
			},
		},
		Type:     parser.ALL,
		Operator: "=",
		Error:    "row value length does not match",
	},
}

func TestInRowValueList(t *testing.T) {
	for _, v := range inRowValueListTests {
		r, err := InRowValueList(v.LHS, v.List, v.Type, v.Operator)
		if err != nil {
			if len(v.Error) < 1 {
				t.Errorf("unexpected error %q for (%s %s %s %s)", err, v.LHS, v.Operator, parser.TokenLiteral(v.Type), v.List)
			} else if err.Error() != v.Error {
				t.Errorf("error %q, want error %q for (%s %s %s %s)", err.Error(), v.Error, v.LHS, v.Operator, parser.TokenLiteral(v.Type), v.List)
			}
			continue
		}
		if 0 < len(v.Error) {
			t.Errorf("no error, want error %q for (%s %s %s %s)", v.Error, v.LHS, v.Operator, parser.TokenLiteral(v.Type), v.List)
			continue
		}
		if r != v.Result {
			t.Errorf("result = %s, want %s for (%s %s %s %s)", r, v.Result, v.LHS, v.Operator, parser.TokenLiteral(v.Type), v.List)
		}
	}
}
