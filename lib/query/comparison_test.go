package query

import (
	"testing"

	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/value"

	"github.com/mithrandie/ternary"
)

var isTests = []struct {
	LHS    value.Primary
	RHS    value.Primary
	Result ternary.Value
}{
	{
		LHS:    value.NewBoolean(true),
		RHS:    value.NewTernary(ternary.TRUE),
		Result: ternary.TRUE,
	},
	{
		LHS:    value.NewNull(),
		RHS:    value.NewNull(),
		Result: ternary.TRUE,
	},
	{
		LHS:    value.NewString("foo"),
		RHS:    value.NewNull(),
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
	LHS     value.Primary
	Pattern value.Primary
	Result  ternary.Value
}{
	{
		LHS:     value.NewString("str"),
		Pattern: value.NewNull(),
		Result:  ternary.UNKNOWN,
	},
	{
		LHS:     value.NewString("str"),
		Pattern: value.NewBoolean(true),
		Result:  ternary.UNKNOWN,
	},
	{
		LHS:     value.NewBoolean(true),
		Pattern: value.NewString("str"),
		Result:  ternary.UNKNOWN,
	},
	{
		LHS:     value.NewString("str"),
		Pattern: value.NewString("str"),
		Result:  ternary.TRUE,
	},
	{
		LHS:     value.NewString("str"),
		Pattern: value.NewString(""),
		Result:  ternary.FALSE,
	},
	{
		LHS:     value.NewString("abcdefghijk"),
		Pattern: value.NewString("lmn"),
		Result:  ternary.FALSE,
	},
	{
		LHS:     value.NewString("abc"),
		Pattern: value.NewString("_____"),
		Result:  ternary.FALSE,
	},
	{
		LHS:     value.NewString("abcde"),
		Pattern: value.NewString("___"),
		Result:  ternary.FALSE,
	},
	{
		LHS:     value.NewString("abcde"),
		Pattern: value.NewString("___%__"),
		Result:  ternary.TRUE,
	},
	{
		LHS:     value.NewString("abcde"),
		Pattern: value.NewString("_%__"),
		Result:  ternary.TRUE,
	},
	{
		LHS:     value.NewString("abcdefghijkabcdefghijk"),
		Pattern: value.NewString("%def%_abc%"),
		Result:  ternary.TRUE,
	},
	{
		LHS:     value.NewString("abcde"),
		Pattern: value.NewString("abc\\_e"),
		Result:  ternary.FALSE,
	},
	{
		LHS:     value.NewString("abc_e"),
		Pattern: value.NewString("abc\\_e"),
		Result:  ternary.TRUE,
	},
	{
		LHS:     value.NewString("abcde"),
		Pattern: value.NewString("abc\\%e"),
		Result:  ternary.FALSE,
	},
	{
		LHS:     value.NewString("a\\bc%e"),
		Pattern: value.NewString("a\\bc\\%e"),
		Result:  ternary.TRUE,
	},
	{
		LHS:     value.NewString("abcdef"),
		Pattern: value.NewString("abcde\\"),
		Result:  ternary.FALSE,
	},
	{
		LHS:     value.NewString("abcde"),
		Pattern: value.NewString("abc"),
		Result:  ternary.FALSE,
	},
	{
		LHS:     value.NewString("abcdecba"),
		Pattern: value.NewString("%c_a"),
		Result:  ternary.TRUE,
	},
	{
		LHS:     value.NewString("aaaaa"),
		Pattern: value.NewString("%a"),
		Result:  ternary.TRUE,
	},
	{
		LHS:     value.NewString("abababc"),
		Pattern: value.NewString("%aba_c"),
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

var inRowValueListTests = []struct {
	LHS      value.RowValue
	List     []value.RowValue
	Type     int
	Operator string
	Result   ternary.Value
	Error    string
}{
	{
		LHS: value.RowValue{
			value.NewInteger(1),
			value.NewInteger(2),
			value.NewInteger(3),
		},
		List: []value.RowValue{
			{
				value.NewInteger(1),
				value.NewInteger(2),
				value.NewInteger(3),
			},
			{
				value.NewInteger(4),
				value.NewInteger(5),
				value.NewInteger(6),
			},
		},
		Type:     parser.ANY,
		Operator: "=",
		Result:   ternary.TRUE,
	},
	{
		LHS: value.RowValue{
			value.NewInteger(1),
			value.NewInteger(2),
			value.NewInteger(3),
		},
		List: []value.RowValue{
			{
				value.NewInteger(1),
				value.NewInteger(2),
				value.NewInteger(3),
			},
			{
				value.NewInteger(1),
				value.NewInteger(2),
				value.NewInteger(3),
			},
		},
		Type:     parser.ALL,
		Operator: "=",
		Result:   ternary.TRUE,
	},
	{
		LHS: value.RowValue{
			value.NewInteger(1),
			value.NewInteger(2),
			value.NewInteger(3),
		},
		List: []value.RowValue{
			{
				value.NewInteger(1),
				value.NewInteger(2),
				value.NewInteger(3),
			},
			{
				value.NewInteger(1),
				value.NewInteger(3),
			},
		},
		Type:     parser.ALL,
		Operator: "=",
		Error:    "row value length does not match at index 1",
	},
}

func TestInRowValueList(t *testing.T) {
	for _, v := range inRowValueListTests {
		r, err := InRowValueList(v.LHS, v.List, v.Type, v.Operator, TestTx.Flags.DatetimeFormat)
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
