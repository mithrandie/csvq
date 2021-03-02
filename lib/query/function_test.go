package query

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/mithrandie/csvq/lib/cmd"

	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/value"

	"github.com/mithrandie/ternary"
)

type functionTest struct {
	Name     string
	Function parser.Function
	Args     []value.Primary
	Result   value.Primary
	Error    string
}

func testFunction(t *testing.T, f func(parser.Function, []value.Primary, *cmd.Flags) (value.Primary, error), tests []functionTest) {
	for _, v := range tests {
		result, err := f(v.Function, v.Args, TestTx.Flags)
		if err != nil {
			if len(v.Error) < 1 {
				t.Errorf("%s: unexpected error %q", v.Name, err)
			} else if v.Error != "environment-dependent" && err.Error() != v.Error {
				t.Errorf("%s: error %q, want error %q", v.Name, err.Error(), v.Error)
			}
			continue
		}
		if 0 < len(v.Error) {
			t.Errorf("%s: no error, want error %q", v.Name, v.Error)
			continue
		}
		if !reflect.DeepEqual(result, v.Result) {
			t.Errorf("%s: result = %s, want %s", v.Name, result, v.Result)
		}
	}
}

var coalesceTests = []functionTest{
	{
		Name: "Coalesce",
		Function: parser.Function{
			Name: "coalesce",
		},
		Args: []value.Primary{
			value.NewNull(),
			value.NewString("str"),
		},
		Result: value.NewString("str"),
	},
	{
		Name: "Coalesce Argments Error",
		Function: parser.Function{
			Name: "coalesce",
		},
		Args:  []value.Primary{},
		Error: "function coalesce takes at least 1 argument",
	},
	{
		Name: "Coalesce No Match",
		Function: parser.Function{
			Name: "coalesce",
		},
		Args: []value.Primary{
			value.NewNull(),
			value.NewNull(),
		},
		Result: value.NewNull(),
	},
}

func TestCoalesce(t *testing.T) {
	testFunction(t, Coalesce, coalesceTests)
}

var ifTests = []functionTest{
	{
		Name: "If True",
		Function: parser.Function{
			Name: "if",
		},
		Args: []value.Primary{
			value.NewTernary(ternary.TRUE),
			value.NewInteger(1),
			value.NewInteger(2),
		},
		Result: value.NewInteger(1),
	},
	{
		Name: "If False",
		Function: parser.Function{
			Name: "if",
		},
		Args: []value.Primary{
			value.NewTernary(ternary.FALSE),
			value.NewInteger(1),
			value.NewInteger(2),
		},
		Result: value.NewInteger(2),
	},
	{
		Name: "If Argumants Error",
		Function: parser.Function{
			Name: "if",
		},
		Args: []value.Primary{
			value.NewTernary(ternary.FALSE),
		},
		Error: "function if takes exactly 3 arguments",
	},
}

func TestIf(t *testing.T) {
	testFunction(t, If, ifTests)
}

var ifnullTests = []functionTest{
	{
		Name: "Ifnull True",
		Function: parser.Function{
			Name: "ifnull",
		},
		Args: []value.Primary{
			value.NewNull(),
			value.NewInteger(2),
		},
		Result: value.NewInteger(2),
	},
	{
		Name: "Ifnull False",
		Function: parser.Function{
			Name: "ifnull",
		},
		Args: []value.Primary{
			value.NewInteger(1),
			value.NewInteger(2),
		},
		Result: value.NewInteger(1),
	},
	{
		Name: "Ifnull Arguments Error",
		Function: parser.Function{
			Name: "ifnull",
		},
		Args: []value.Primary{
			value.NewInteger(1),
		},
		Error: "function ifnull takes exactly 2 arguments",
	},
}

func TestIfnull(t *testing.T) {
	testFunction(t, Ifnull, ifnullTests)
}

var nullifTests = []functionTest{
	{
		Name: "Nullif True",
		Function: parser.Function{
			Name: "nullif",
		},
		Args: []value.Primary{
			value.NewInteger(2),
			value.NewInteger(2),
		},
		Result: value.NewNull(),
	},
	{
		Name: "Nullif False",
		Function: parser.Function{
			Name: "nullif",
		},
		Args: []value.Primary{
			value.NewInteger(1),
			value.NewInteger(2),
		},
		Result: value.NewInteger(1),
	},
	{
		Name: "Nullif Arguments Error",
		Function: parser.Function{
			Name: "nullif",
		},
		Args: []value.Primary{
			value.NewInteger(1),
		},
		Error: "function nullif takes exactly 2 arguments",
	},
}

func TestNullif(t *testing.T) {
	testFunction(t, Nullif, nullifTests)
}

var ceilTests = []functionTest{
	{
		Name: "Ceil",
		Function: parser.Function{
			Name: "ceil",
		},
		Args: []value.Primary{
			value.NewFloat(2.345),
		},
		Result: value.NewInteger(3),
	},
	{
		Name: "Ceil Null",
		Function: parser.Function{
			Name: "ceil",
		},
		Args: []value.Primary{
			value.NewNull(),
		},
		Result: value.NewNull(),
	},
	{
		Name: "Ceil Place is Null",
		Function: parser.Function{
			Name: "ceil",
		},
		Args: []value.Primary{
			value.NewFloat(2.345),
			value.NewNull(),
		},
		Result: value.NewNull(),
	},
	{
		Name: "Ceil Arguments Error",
		Function: parser.Function{
			Name: "ceil",
		},
		Args:  []value.Primary{},
		Error: "function ceil takes 1 or 2 arguments",
	},
}

func TestCeil(t *testing.T) {
	testFunction(t, Ceil, ceilTests)
}

var floorTests = []functionTest{
	{
		Name: "Floor",
		Function: parser.Function{
			Name: "floor",
		},
		Args: []value.Primary{
			value.NewFloat(2.345),
			value.NewInteger(1),
		},
		Result: value.NewFloat(2.3),
	},
	{
		Name: "Floor Null",
		Function: parser.Function{
			Name: "floor",
		},
		Args: []value.Primary{
			value.NewNull(),
		},
		Result: value.NewNull(),
	},
	{
		Name: "Floor Arguments Error",
		Function: parser.Function{
			Name: "floor",
		},
		Args:  []value.Primary{},
		Error: "function floor takes 1 or 2 arguments",
	},
}

func TestFloor(t *testing.T) {
	testFunction(t, Floor, floorTests)
}

var roundTests = []functionTest{
	{
		Name: "Round",
		Function: parser.Function{
			Name: "round",
		},
		Args: []value.Primary{
			value.NewFloat(2.456),
			value.NewInteger(2),
		},
		Result: value.NewFloat(2.46),
	},
	{
		Name: "Round Negative Number",
		Function: parser.Function{
			Name: "round",
		},
		Args: []value.Primary{
			value.NewFloat(-2.456),
			value.NewInteger(2),
		},
		Result: value.NewFloat(-2.46),
	},
	{
		Name: "Round Null",
		Function: parser.Function{
			Name: "round",
		},
		Args: []value.Primary{
			value.NewNull(),
		},
		Result: value.NewNull(),
	},
	{
		Name: "Round Arguments Error",
		Function: parser.Function{
			Name: "round",
		},
		Args:  []value.Primary{},
		Error: "function round takes 1 or 2 arguments",
	},
}

func TestRound(t *testing.T) {
	testFunction(t, Round, roundTests)
}

var absTests = []functionTest{
	{
		Name: "Abs",
		Function: parser.Function{
			Name: "abs",
		},
		Args: []value.Primary{
			value.NewInteger(-2),
		},
		Result: value.NewInteger(2),
	},
	{
		Name: "Abs Null",
		Function: parser.Function{
			Name: "abs",
		},
		Args: []value.Primary{
			value.NewNull(),
		},
		Result: value.NewNull(),
	},
	{
		Name: "Abs Arguments Error",
		Function: parser.Function{
			Name: "abs",
		},
		Args:  []value.Primary{},
		Error: "function abs takes exactly 1 argument",
	},
}

func TestAbs(t *testing.T) {
	testFunction(t, Abs, absTests)
}

var acosTests = []functionTest{
	{
		Name: "Acos",
		Function: parser.Function{
			Name: "acos",
		},
		Args: []value.Primary{
			value.NewInteger(0),
		},
		Result: value.NewFloat(1.5707963267948966),
	},
}

func TestAcos(t *testing.T) {
	testFunction(t, Acos, acosTests)
}

var asinTests = []functionTest{
	{
		Name: "Asin",
		Function: parser.Function{
			Name: "asin",
		},
		Args: []value.Primary{
			value.NewFloat(0.1),
		},
		Result: value.NewFloat(0.1001674211615598),
	},
}

func TestAsin(t *testing.T) {
	testFunction(t, Asin, asinTests)
}

var atanTests = []functionTest{
	{
		Name: "Atan",
		Function: parser.Function{
			Name: "atan",
		},
		Args: []value.Primary{
			value.NewInteger(2),
		},
		Result: value.NewFloat(1.1071487177940904),
	},
}

func TestAtan(t *testing.T) {
	testFunction(t, Atan, atanTests)
}

var atan2Tests = []functionTest{
	{
		Name: "Atan2",
		Function: parser.Function{
			Name: "atan2",
		},
		Args: []value.Primary{
			value.NewInteger(2),
			value.NewInteger(2),
		},
		Result: value.NewFloat(0.7853981633974483),
	},
}

func TestAtan2(t *testing.T) {
	testFunction(t, Atan2, atan2Tests)
}

var cosTests = []functionTest{
	{
		Name: "Cos",
		Function: parser.Function{
			Name: "cos",
		},
		Args: []value.Primary{
			value.NewInteger(2),
		},
		Result: value.NewFloat(-0.4161468365471424),
	},
}

func TestCos(t *testing.T) {
	testFunction(t, Cos, cosTests)
}

var sinTests = []functionTest{
	{
		Name: "Sin",
		Function: parser.Function{
			Name: "sin",
		},
		Args: []value.Primary{
			value.NewInteger(1),
		},
		Result: value.NewFloat(0.8414709848078965),
	},
}

func TestSin(t *testing.T) {
	testFunction(t, Sin, sinTests)
}

var tanTests = []functionTest{
	{
		Name: "Tan",
		Function: parser.Function{
			Name: "tan",
		},
		Args: []value.Primary{
			value.NewInteger(2),
		},
		Result: value.NewFloat(-2.185039863261519),
	},
}

func TestTan(t *testing.T) {
	testFunction(t, Tan, tanTests)
}

var expTests = []functionTest{
	{
		Name: "Exp",
		Function: parser.Function{
			Name: "exp",
		},
		Args: []value.Primary{
			value.NewInteger(2),
		},
		Result: value.NewFloat(7.38905609893065),
	},
}

func TestExp(t *testing.T) {
	testFunction(t, Exp, expTests)
}

var exp2Tests = []functionTest{
	{
		Name: "Exp2",
		Function: parser.Function{
			Name: "exp2",
		},
		Args: []value.Primary{
			value.NewInteger(2),
		},
		Result: value.NewInteger(4),
	},
}

func TestExp2(t *testing.T) {
	testFunction(t, Exp2, exp2Tests)
}

var expm1Tests = []functionTest{
	{
		Name: "Expm1",
		Function: parser.Function{
			Name: "expm1",
		},
		Args: []value.Primary{
			value.NewInteger(2),
		},
		Result: value.NewFloat(6.38905609893065),
	},
}

func TestExpm1(t *testing.T) {
	testFunction(t, Expm1, expm1Tests)
}

var mathLogTests = []functionTest{
	{
		Name: "MathLog",
		Function: parser.Function{
			Name: "log",
		},
		Args: []value.Primary{
			value.NewFloat(2),
		},
		Result: value.NewFloat(0.6931471805599453),
	},
}

func TestMathLog(t *testing.T) {
	testFunction(t, MathLog, mathLogTests)
}

var log10Tests = []functionTest{
	{
		Name: "Log10",
		Function: parser.Function{
			Name: "log10",
		},
		Args: []value.Primary{
			value.NewFloat(100),
		},
		Result: value.NewInteger(2),
	},
}

func TestLog10(t *testing.T) {
	testFunction(t, Log10, log10Tests)
}

var log2Tests = []functionTest{
	{
		Name: "Log2",
		Function: parser.Function{
			Name: "log2",
		},
		Args: []value.Primary{
			value.NewFloat(16),
		},
		Result: value.NewInteger(4),
	},
}

func TestLog2(t *testing.T) {
	testFunction(t, Log2, log2Tests)
}

var log1pTests = []functionTest{
	{
		Name: "Log1p",
		Function: parser.Function{
			Name: "log1p",
		},
		Args: []value.Primary{
			value.NewFloat(1),
		},
		Result: value.NewFloat(0.6931471805599453),
	},
}

func TestLog1p(t *testing.T) {
	testFunction(t, Log1p, log1pTests)
}

var sqrtTests = []functionTest{
	{
		Name: "Sqrt",
		Function: parser.Function{
			Name: "sqrt",
		},
		Args: []value.Primary{
			value.NewFloat(4),
		},
		Result: value.NewInteger(2),
	},
	{
		Name: "Sqrt Cannot Calculate",
		Function: parser.Function{
			Name: "sqrt",
		},
		Args: []value.Primary{
			value.NewFloat(-4),
		},
		Result: value.NewNull(),
	},
}

func TestSqrt(t *testing.T) {
	testFunction(t, Sqrt, sqrtTests)
}

var powTests = []functionTest{
	{
		Name: "Pow",
		Function: parser.Function{
			Name: "pow",
		},
		Args: []value.Primary{
			value.NewFloat(2),
			value.NewFloat(2),
		},
		Result: value.NewInteger(4),
	},
	{
		Name: "Pow First Argument is Null",
		Function: parser.Function{
			Name: "pow",
		},
		Args: []value.Primary{
			value.NewNull(),
			value.NewFloat(2),
		},
		Result: value.NewNull(),
	},
	{
		Name: "Pow Second Argument is Null",
		Function: parser.Function{
			Name: "pow",
		},
		Args: []value.Primary{
			value.NewFloat(2),
			value.NewNull(),
		},
		Result: value.NewNull(),
	},
	{
		Name: "Pow Cannot Calculate",
		Function: parser.Function{
			Name: "pow",
		},
		Args: []value.Primary{
			value.NewFloat(-2),
			value.NewFloat(2.4),
		},
		Result: value.NewNull(),
	},
	{
		Name: "Pow Arguments Error",
		Function: parser.Function{
			Name: "pow",
		},
		Args:  []value.Primary{},
		Error: "function pow takes exactly 2 arguments",
	},
}

func TestPow(t *testing.T) {
	testFunction(t, Pow, powTests)
}

var binToDecTests = []functionTest{
	{
		Name: "BinToDec",
		Function: parser.Function{
			Name: "bin_to_dec",
		},
		Args: []value.Primary{
			value.NewString("1111011"),
		},
		Result: value.NewInteger(123),
	},
	{
		Name: "BinToDec Null",
		Function: parser.Function{
			Name: "bin_to_dec",
		},
		Args: []value.Primary{
			value.NewNull(),
		},
		Result: value.NewNull(),
	},
	{
		Name: "BinToDec Parse Error",
		Function: parser.Function{
			Name: "bin_to_dec",
		},
		Args: []value.Primary{
			value.NewString("string"),
		},
		Result: value.NewNull(),
	},
	{
		Name: "BinToDec Arguments Error",
		Function: parser.Function{
			Name: "bin_to_dec",
		},
		Args:  []value.Primary{},
		Error: "function bin_to_dec takes exactly 1 argument",
	},
}

func TestBinToDec(t *testing.T) {
	testFunction(t, BinToDec, binToDecTests)
}

var octToDecTests = []functionTest{
	{
		Name: "OctToDec",
		Function: parser.Function{
			Name: "oct_to_dec",
		},
		Args: []value.Primary{
			value.NewString("0173"),
		},
		Result: value.NewInteger(123),
	},
}

func TestOctToDec(t *testing.T) {
	testFunction(t, OctToDec, octToDecTests)
}

var hexToDecTests = []functionTest{
	{
		Name: "HexToDec",
		Function: parser.Function{
			Name: "hex_to_dec",
		},
		Args: []value.Primary{
			value.NewString("0x7b"),
		},
		Result: value.NewInteger(123),
	},
}

func TestHexToDec(t *testing.T) {
	testFunction(t, HexToDec, hexToDecTests)
}

var enotationToDecTests = []functionTest{
	{
		Name: "EnotationToDec",
		Function: parser.Function{
			Name: "enotation_to_dec",
		},
		Args: []value.Primary{
			value.NewString("1.23e-11"),
		},
		Result: value.NewFloat(0.0000000000123),
	},
	{
		Name: "EnotationToDec To Integer",
		Function: parser.Function{
			Name: "enotation_to_dec",
		},
		Args: []value.Primary{
			value.NewString("1.23e+12"),
		},
		Result: value.NewInteger(1230000000000),
	},
	{
		Name: "EnotationToDec Null",
		Function: parser.Function{
			Name: "enotation_to_dec",
		},
		Args: []value.Primary{
			value.NewNull(),
		},
		Result: value.NewNull(),
	},
	{
		Name: "EnotationToDec Parse Error",
		Function: parser.Function{
			Name: "enotation_to_dec",
		},
		Args: []value.Primary{
			value.NewString("string"),
		},
		Result: value.NewNull(),
	},
	{
		Name: "EnotationToDec Arguments Error",
		Function: parser.Function{
			Name: "enotation_to_dec",
		},
		Args:  []value.Primary{},
		Error: "function enotation_to_dec takes exactly 1 argument",
	},
}

func TestEnotationToDec(t *testing.T) {
	testFunction(t, EnotationToDec, enotationToDecTests)
}

var binTests = []functionTest{
	{
		Name: "Bin",
		Function: parser.Function{
			Name: "bin",
		},
		Args: []value.Primary{
			value.NewInteger(123),
		},
		Result: value.NewString("1111011"),
	},
	{
		Name: "Bin Null",
		Function: parser.Function{
			Name: "bin",
		},
		Args: []value.Primary{
			value.NewNull(),
		},
		Result: value.NewNull(),
	},
	{
		Name: "Bin Arguments Error",
		Function: parser.Function{
			Name: "bin",
		},
		Args:  []value.Primary{},
		Error: "function bin takes exactly 1 argument",
	},
}

func TestBin(t *testing.T) {
	testFunction(t, Bin, binTests)
}

var octTests = []functionTest{
	{
		Name: "Oct",
		Function: parser.Function{
			Name: "oct",
		},
		Args: []value.Primary{
			value.NewInteger(123),
		},
		Result: value.NewString("173"),
	},
}

func TestOct(t *testing.T) {
	testFunction(t, Oct, octTests)
}

var hexTests = []functionTest{
	{
		Name: "Hex",
		Function: parser.Function{
			Name: "hex",
		},
		Args: []value.Primary{
			value.NewInteger(123),
		},
		Result: value.NewString("7b"),
	},
}

func TestHex(t *testing.T) {
	testFunction(t, Hex, hexTests)
}

var enotationTests = []functionTest{
	{
		Name: "Enotation",
		Function: parser.Function{
			Name: "enotation",
		},
		Args: []value.Primary{
			value.NewFloat(0.0000000000123),
		},
		Result: value.NewString("1.23e-11"),
	},
	{
		Name: "Enotation From Integer",
		Function: parser.Function{
			Name: "enotation",
		},
		Args: []value.Primary{
			value.NewInteger(1230000000000),
		},
		Result: value.NewString("1.23e+12"),
	},
	{
		Name: "Enotation Null",
		Function: parser.Function{
			Name: "enotation",
		},
		Args: []value.Primary{
			value.NewNull(),
		},
		Result: value.NewNull(),
	},
	{
		Name: "Enotation Arguments Error",
		Function: parser.Function{
			Name: "enotation",
		},
		Args:  []value.Primary{},
		Error: "function enotation takes exactly 1 argument",
	},
}

func TestEnotation(t *testing.T) {
	testFunction(t, Enotation, enotationTests)
}

var numberFormatTests = []functionTest{
	{
		Name: "NumberFormat",
		Function: parser.Function{
			Name: "number_format",
		},
		Args: []value.Primary{
			value.NewFloat(123456.789123),
			value.NewInteger(4),
			value.NewString(","),
			value.NewString(" "),
			value.NewString(" "),
		},
		Result: value.NewString("123 456,789 1"),
	},
	{
		Name: "NumberFormat Null",
		Function: parser.Function{
			Name: "number_format",
		},
		Args: []value.Primary{
			value.NewNull(),
		},
		Result: value.NewNull(),
	},
	{
		Name: "NumberFormat Arguments Errlr",
		Function: parser.Function{
			Name: "number_format",
		},
		Args:  []value.Primary{},
		Error: "function number_format takes 1 to 5 arguments",
	},
}

func TestNumberFormat(t *testing.T) {
	testFunction(t, NumberFormat, numberFormatTests)
}

var randTests = []struct {
	Name      string
	Function  parser.Function
	Args      []value.Primary
	RangeLow  float64
	RangeHigh float64
	Error     string
}{
	{
		Name: "Rand",
		Function: parser.Function{
			Name: "rand",
		},
		RangeLow:  0.0,
		RangeHigh: 1.0,
	},
	{
		Name: "Rand Range Specified",
		Function: parser.Function{
			Name: "rand",
		},
		Args: []value.Primary{
			value.NewInteger(7),
			value.NewInteger(12),
		},
		RangeLow:  7.0,
		RangeHigh: 12.0,
	},
	{
		Name: "Range Arguments Error",
		Function: parser.Function{
			Name: "rand",
		},
		Args: []value.Primary{
			value.NewInteger(1),
		},
		Error: "function rand takes 0 or 2 arguments",
	},
	{
		Name: "Range First Arguments Error",
		Function: parser.Function{
			Name: "rand",
		},
		Args: []value.Primary{
			value.NewString("a"),
			value.NewInteger(2),
		},
		Error: "the first argument must be an integer for function rand",
	},
	{
		Name: "Range Second Arguments Error",
		Function: parser.Function{
			Name: "rand",
		},
		Args: []value.Primary{
			value.NewInteger(1),
			value.NewString("a"),
		},
		Error: "the second argument must be an integer for function rand",
	},
	{
		Name: "Range Arguments Value Error",
		Function: parser.Function{
			Name: "rand",
		},
		Args: []value.Primary{
			value.NewInteger(1),
			value.NewInteger(1),
		},
		Error: "the second argument must be greater than the first argument for function rand",
	},
}

func TestRand(t *testing.T) {
	for _, v := range randTests {
		result, err := Rand(v.Function, v.Args, TestTx.Flags)
		if err != nil {
			if len(v.Error) < 1 {
				t.Errorf("%s: unexpected error %q", v.Name, err)
			} else if err.Error() != v.Error {
				t.Errorf("%s: error %q, want error %q", v.Name, err.Error(), v.Error)
			}
			continue
		}
		if 0 < len(v.Error) {
			t.Errorf("%s: no error, want error %q", v.Name, v.Error)
			continue
		}

		var f float64
		if len(v.Args) < 1 {
			f = result.(*value.Float).Raw()
		} else {
			f = float64(result.(*value.Integer).Raw())
		}

		if f < v.RangeLow || v.RangeHigh < f {
			t.Errorf("%s: result = %f, want in range from %f to %f", v.Name, f, v.RangeLow, v.RangeHigh)
		}
	}
}

var trimTests = []functionTest{
	{
		Name: "Trim",
		Function: parser.Function{
			Name: "trim",
		},
		Args: []value.Primary{
			value.NewString("aabbfoo, baraabb"),
			value.NewString("ab"),
		},
		Result: value.NewString("foo, bar"),
	},
	{
		Name: "Trim Spaces",
		Function: parser.Function{
			Name: "trim",
		},
		Args: []value.Primary{
			value.NewString("  foo, bar \n"),
		},
		Result: value.NewString("foo, bar"),
	},
	{
		Name: "Trim Null",
		Function: parser.Function{
			Name: "trim",
		},
		Args: []value.Primary{
			value.NewNull(),
			value.NewString("ab"),
		},
		Result: value.NewNull(),
	},
	{
		Name: "Trim Cutset is Null",
		Function: parser.Function{
			Name: "trim",
		},
		Args: []value.Primary{
			value.NewString("aabbfoo, baraabb"),
			value.NewNull(),
		},
		Result: value.NewNull(),
	},
	{
		Name: "Trim Arguments Error",
		Function: parser.Function{
			Name: "trim",
		},
		Args:  []value.Primary{},
		Error: "function trim takes 1 or 2 arguments",
	},
}

func TestTrim(t *testing.T) {
	testFunction(t, Trim, trimTests)
}

var ltrimTests = []functionTest{
	{
		Name: "Ltrim",
		Function: parser.Function{
			Name: "ltrim",
		},
		Args: []value.Primary{
			value.NewString("aabbfoo, baraabb"),
			value.NewString("ab"),
		},
		Result: value.NewString("foo, baraabb"),
	},
	{
		Name: "Ltrim Spaces",
		Function: parser.Function{
			Name: "ltrim",
		},
		Args: []value.Primary{
			value.NewString("  foo, bar \n"),
		},
		Result: value.NewString("foo, bar \n"),
	},
}

func TestLtrim(t *testing.T) {
	testFunction(t, Ltrim, ltrimTests)
}

var rtrimTests = []functionTest{
	{
		Name: "Rtrim",
		Function: parser.Function{
			Name: "rtrim",
		},
		Args: []value.Primary{
			value.NewString("aabbfoo, baraabb"),
			value.NewString("ab"),
		},
		Result: value.NewString("aabbfoo, bar"),
	},
	{
		Name: "Rtrim Spaces",
		Function: parser.Function{
			Name: "rtrim",
		},
		Args: []value.Primary{
			value.NewString("  foo, bar \n"),
		},
		Result: value.NewString("  foo, bar"),
	},
}

func TestRtrim(t *testing.T) {
	testFunction(t, Rtrim, rtrimTests)
}

var upperTests = []functionTest{
	{
		Name: "Upper",
		Function: parser.Function{
			Name: "upper",
		},
		Args: []value.Primary{
			value.NewString("Foo"),
		},
		Result: value.NewString("FOO"),
	},
	{
		Name: "Upper Null",
		Function: parser.Function{
			Name: "upper",
		},
		Args: []value.Primary{
			value.NewNull(),
		},
		Result: value.NewNull(),
	},
	{
		Name: "Upper Arguments Error",
		Function: parser.Function{
			Name: "upper",
		},
		Args:  []value.Primary{},
		Error: "function upper takes exactly 1 argument",
	},
}

func TestUpper(t *testing.T) {
	testFunction(t, Upper, upperTests)
}

var lowerTests = []functionTest{
	{
		Name: "Lower",
		Function: parser.Function{
			Name: "lower",
		},
		Args: []value.Primary{
			value.NewString("Foo"),
		},
		Result: value.NewString("foo"),
	},
}

func TestLower(t *testing.T) {
	testFunction(t, Lower, lowerTests)
}

var base64EncodeTests = []functionTest{
	{
		Name: "Base64Encode",
		Function: parser.Function{
			Name: "base64_encode",
		},
		Args: []value.Primary{
			value.NewString("Foo"),
		},
		Result: value.NewString("Rm9v"),
	},
}

func TestBase64Encode(t *testing.T) {
	testFunction(t, Base64Encode, base64EncodeTests)
}

var base64DecodeTests = []functionTest{
	{
		Name: "Base64Decode",
		Function: parser.Function{
			Name: "base64_decode",
		},
		Args: []value.Primary{
			value.NewString("Rm9v"),
		},
		Result: value.NewString("Foo"),
	},
}

func TestBase64Decode(t *testing.T) {
	testFunction(t, Base64Decode, base64DecodeTests)
}

var hexEncodeTests = []functionTest{
	{
		Name: "HexEncode",
		Function: parser.Function{
			Name: "hex_encode",
		},
		Args: []value.Primary{
			value.NewString("Foo"),
		},
		Result: value.NewString("466f6f"),
	},
}

func TestHexEncode(t *testing.T) {
	testFunction(t, HexEncode, hexEncodeTests)
}

var hexDecodeTests = []functionTest{
	{
		Name: "HexDecode",
		Function: parser.Function{
			Name: "hex_decode",
		},
		Args: []value.Primary{
			value.NewString("466f6f"),
		},
		Result: value.NewString("Foo"),
	},
}

func TestHexDecode(t *testing.T) {
	testFunction(t, HexDecode, hexDecodeTests)
}

var lenTests = []functionTest{
	{
		Name: "Len",
		Function: parser.Function{
			Name: "len",
		},
		Args: []value.Primary{
			value.NewString("日本語"),
		},
		Result: value.NewInteger(3),
	},
	{
		Name: "Len Null",
		Function: parser.Function{
			Name: "len",
		},
		Args: []value.Primary{
			value.NewNull(),
		},
		Result: value.NewNull(),
	},
	{
		Name: "Len Arguments Error",
		Function: parser.Function{
			Name: "len",
		},
		Args:  []value.Primary{},
		Error: "function len takes exactly 1 argument",
	},
}

func TestLen(t *testing.T) {
	testFunction(t, Len, lenTests)
}

var byteLenTests = []functionTest{
	{
		Name: "ByteLen",
		Function: parser.Function{
			Name: "byte_len",
		},
		Args: []value.Primary{
			value.NewString("abc日本語"),
		},
		Result: value.NewInteger(12),
	},
	{
		Name: "ByteLen SJIS",
		Function: parser.Function{
			Name: "byte_len",
		},
		Args: []value.Primary{
			value.NewString("abc日本語"),
			value.NewString("sjis"),
		},
		Result: value.NewInteger(9),
	},
	{
		Name: "ByteLen Null",
		Function: parser.Function{
			Name: "byte_len",
		},
		Args: []value.Primary{
			value.NewNull(),
		},
		Result: value.NewNull(),
	},
	{
		Name: "ByteLen Arguments Error",
		Function: parser.Function{
			Name: "byte_len",
		},
		Args:  []value.Primary{},
		Error: "function byte_len takes 1 or 2 arguments",
	},
	{
		Name: "ByteLen Invalid Encoding Error",
		Function: parser.Function{
			Name: "byte_len",
		},
		Args: []value.Primary{
			value.NewString("abc日本語"),
			value.NewString("invalid"),
		},
		Error: "encoding must be one of UTF8|UTF16|SJIS for function byte_len",
	},
}

func TestByteLen(t *testing.T) {
	testFunction(t, ByteLen, byteLenTests)
}

var widthTests = []functionTest{
	{
		Name: "Width",
		Function: parser.Function{
			Name: "width",
		},
		Args: []value.Primary{
			value.NewString("abc日本語"),
		},
		Result: value.NewInteger(9),
	},
	{
		Name: "Width Arguments Is Null",
		Function: parser.Function{
			Name: "width",
		},
		Args: []value.Primary{
			value.NewNull(),
		},
		Result: value.NewNull(),
	},
	{
		Name: "Width Arguments Error",
		Function: parser.Function{
			Name: "width",
		},
		Args:  []value.Primary{},
		Error: "function width takes exactly 1 argument",
	},
}

func TestWidth(t *testing.T) {
	testFunction(t, Width, widthTests)
}

var lpadTests = []functionTest{
	{
		Name: "Lpad",
		Function: parser.Function{
			Name: "lpad",
		},
		Args: []value.Primary{
			value.NewString("aaaaa"),
			value.NewInteger(10),
			value.NewString("01"),
		},
		Result: value.NewString("01010aaaaa"),
	},
	{
		Name: "Lpad by Byte Length",
		Function: parser.Function{
			Name: "lpad",
		},
		Args: []value.Primary{
			value.NewString("日本語"),
			value.NewInteger(12),
			value.NewString("空白"),
			value.NewString("byte"),
			value.NewString("sjis"),
		},
		Result: value.NewString("空白空日本語"),
	},
	{
		Name: "Lpad by Byte Length Error",
		Function: parser.Function{
			Name: "lpad",
		},
		Args: []value.Primary{
			value.NewString("日本語"),
			value.NewInteger(11),
			value.NewString("空白"),
			value.NewString("byte"),
			value.NewString("sjis"),
		},
		Error: "cannot split pad string in a byte array of a character for function lpad",
	},
	{
		Name: "Lpad PadType Error",
		Function: parser.Function{
			Name: "lpad",
		},
		Args: []value.Primary{
			value.NewString("日本語"),
			value.NewInteger(11),
			value.NewString("空白"),
			value.NewString("invalid"),
		},
		Error: "padding type must be one of LEN|BYTE|WIDTH for function lpad",
	},
	{
		Name: "Lpad Encoding Error",
		Function: parser.Function{
			Name: "lpad",
		},
		Args: []value.Primary{
			value.NewString("日本語"),
			value.NewInteger(11),
			value.NewString("空白"),
			value.NewString("byte"),
			value.NewString("invalid"),
		},
		Error: "encoding must be one of UTF8|UTF16|SJIS for function lpad",
	},
	{
		Name: "Lpad by Width",
		Function: parser.Function{
			Name: "lpad",
		},
		Args: []value.Primary{
			value.NewString("日本語"),
			value.NewInteger(12),
			value.NewString("空白"),
			value.NewString("width"),
		},
		Result: value.NewString("空白空日本語"),
	},
	{
		Name: "Lpad by Width Length Error",
		Function: parser.Function{
			Name: "lpad",
		},
		Args: []value.Primary{
			value.NewString("日本語"),
			value.NewInteger(11),
			value.NewString("空白"),
			value.NewString("width"),
		},
		Error: "cannot split pad string in a byte array of a character for function lpad",
	},
	{
		Name: "Lpad No Padding",
		Function: parser.Function{
			Name: "lpad",
		},
		Args: []value.Primary{
			value.NewString("aaaaa"),
			value.NewInteger(5),
			value.NewString("01"),
		},
		Result: value.NewString("aaaaa"),
	},
	{
		Name: "Lpad String is Null",
		Function: parser.Function{
			Name: "lpad",
		},
		Args: []value.Primary{
			value.NewNull(),
			value.NewInteger(10),
			value.NewString("01"),
		},
		Result: value.NewNull(),
	},
	{
		Name: "Lpad Length is Null",
		Function: parser.Function{
			Name: "lpad",
		},
		Args: []value.Primary{
			value.NewString("aaaaa"),
			value.NewNull(),
			value.NewString("01"),
		},
		Result: value.NewNull(),
	},
	{
		Name: "Lpad Pad String is Null",
		Function: parser.Function{
			Name: "lpad",
		},
		Args: []value.Primary{
			value.NewString("aaaaa"),
			value.NewInteger(10),
			value.NewNull(),
		},
		Result: value.NewNull(),
	},
	{
		Name: "Lpad Arguments Error",
		Function: parser.Function{
			Name: "lpad",
		},
		Args:  []value.Primary{},
		Error: "function lpad takes 3 to 5 arguments",
	},
}

func TestLpad(t *testing.T) {
	testFunction(t, Lpad, lpadTests)
}

var rpadTests = []functionTest{
	{
		Name: "Rpad",
		Function: parser.Function{
			Name: "rpad",
		},
		Args: []value.Primary{
			value.NewString("aaaaa"),
			value.NewInteger(10),
			value.NewString("01"),
		},
		Result: value.NewString("aaaaa01010"),
	},
}

func TestRpad(t *testing.T) {
	testFunction(t, Rpad, rpadTests)
}

var substringTests = []functionTest{
	{
		Name: "Substring with a positive argument",
		Function: parser.Function{
			Name: "substring",
		},
		Args: []value.Primary{
			value.NewString("abcdefghijklmn"),
			value.NewInteger(5),
		},
		Result: value.NewString("efghijklmn"),
	},
	{
		Name: "Substring with a negative argument",
		Function: parser.Function{
			Name: "substring",
		},
		Args: []value.Primary{
			value.NewString("abcdefghijklmn"),
			value.NewInteger(-5),
		},
		Result: value.NewString("jklmn"),
	},
	{
		Name: "Substring with two positive argument",
		Function: parser.Function{
			Name: "substring",
		},
		Args: []value.Primary{
			value.NewString("abcdefghijklmn"),
			value.NewInteger(5),
			value.NewInteger(3),
		},
		Result: value.NewString("efg"),
	},
	{
		Name: "Substring starting with zero",
		Function: parser.Function{
			Name: "substring",
		},
		Args: []value.Primary{
			value.NewString("abcdefghijklmn"),
			value.NewInteger(0),
			value.NewInteger(3),
		},
		Result: value.NewString("abc"),
	},
	{
		Name: "Substring",
		Function: parser.Function{
			Name: "substring",
		},
		Args: []value.Primary{
			value.NewString("abcdefghijklmn"),
			value.NewInteger(-5),
			value.NewInteger(8),
		},
		Result: value.NewString("jklmn"),
	},
	{
		Name: "Substring String is Null",
		Function: parser.Function{
			Name: "substring",
		},
		Args: []value.Primary{
			value.NewNull(),
			value.NewInteger(-5),
			value.NewInteger(8),
		},
		Result: value.NewNull(),
	},
	{
		Name: "Substring StartIndex is Null",
		Function: parser.Function{
			Name: "substring",
		},
		Args: []value.Primary{
			value.NewString("abcdefghijklmn"),
			value.NewNull(),
		},
		Result: value.NewNull(),
	},
	{
		Name: "Substring Length is Null",
		Function: parser.Function{
			Name: "substring",
		},
		Args: []value.Primary{
			value.NewString("abcdefghijklmn"),
			value.NewInteger(-5),
			value.NewNull(),
		},
		Result: value.NewNull(),
	},
	{
		Name: "Substring Length is Negative",
		Function: parser.Function{
			Name: "substring",
		},
		Args: []value.Primary{
			value.NewString("abcdefghijklmn"),
			value.NewInteger(-5),
			value.NewInteger(-1),
		},
		Result: value.NewNull(),
	},
	{
		Name: "Substring StartIndex is Out Of Index",
		Function: parser.Function{
			Name: "substring",
		},
		Args: []value.Primary{
			value.NewString("abcdefghijklmn"),
			value.NewInteger(100),
			value.NewInteger(8),
		},
		Result: value.NewNull(),
	},
	{
		Name: "Substring Arguments Error",
		Function: parser.Function{
			Name: "substring",
		},
		Args:  []value.Primary{},
		Error: "function substring takes 2 or 3 arguments",
	},
}

func TestSubstring(t *testing.T) {
	testFunction(t, Substring, substringTests)
}

var substrTests = []functionTest{
	{
		Name: "Substr with a positive argument",
		Function: parser.Function{
			Name: "substr",
		},
		Args: []value.Primary{
			value.NewString("abcdefghijklmn"),
			value.NewInteger(5),
		},
		Result: value.NewString("fghijklmn"),
	},
	{
		Name: "Substr with a negative argument",
		Function: parser.Function{
			Name: "substr",
		},
		Args: []value.Primary{
			value.NewString("abcdefghijklmn"),
			value.NewInteger(-5),
		},
		Result: value.NewString("jklmn"),
	},
	{
		Name: "Substr with two positive argument",
		Function: parser.Function{
			Name: "substr",
		},
		Args: []value.Primary{
			value.NewString("abcdefghijklmn"),
			value.NewInteger(5),
			value.NewInteger(3),
		},
		Result: value.NewString("fgh"),
	},
}

func TestSubstr(t *testing.T) {
	testFunction(t, Substr, substrTests)
}

var instrTests = []functionTest{
	{
		Name: "Instr",
		Function: parser.Function{
			Name: "instr",
		},
		Args: []value.Primary{
			value.NewString("abcdefghijklmn"),
			value.NewString("def"),
		},
		Result: value.NewInteger(3),
	},
	{
		Name: "Instr String is Null",
		Function: parser.Function{
			Name: "instr",
		},
		Args: []value.Primary{
			value.NewNull(),
			value.NewString("def"),
		},
		Result: value.NewNull(),
	},
	{
		Name: "Instr Substring is Null",
		Function: parser.Function{
			Name: "instr",
		},
		Args: []value.Primary{
			value.NewString("abcdefghijklmn"),
			value.NewNull(),
		},
		Result: value.NewNull(),
	},
	{
		Name: "Instr Substring does not exist",
		Function: parser.Function{
			Name: "instr",
		},
		Args: []value.Primary{
			value.NewString("abcdefghijklmn"),
			value.NewString("zzz"),
		},
		Result: value.NewNull(),
	},
	{
		Name: "Instr Arguments Error",
		Function: parser.Function{
			Name: "instr",
		},
		Args:  []value.Primary{},
		Error: "function instr takes exactly 2 arguments",
	},
}

func TestInstr(t *testing.T) {
	testFunction(t, Instr, instrTests)
}

var listElemTests = []functionTest{
	{
		Name: "ListElem",
		Function: parser.Function{
			Name: "list_elem",
		},
		Args: []value.Primary{
			value.NewString("abc def ghi"),
			value.NewString(" "),
			value.NewInteger(1),
		},
		Result: value.NewString("def"),
	},
	{
		Name: "ListElem String is Null",
		Function: parser.Function{
			Name: "list_elem",
		},
		Args: []value.Primary{
			value.NewNull(),
			value.NewString(" "),
			value.NewInteger(1),
		},
		Result: value.NewNull(),
	},
	{
		Name: "ListElem Separator is Null",
		Function: parser.Function{
			Name: "list_elem",
		},
		Args: []value.Primary{
			value.NewString("abc def ghi"),
			value.NewNull(),
			value.NewInteger(1),
		},
		Result: value.NewNull(),
	},
	{
		Name: "ListElem Index is Null",
		Function: parser.Function{
			Name: "list_elem",
		},
		Args: []value.Primary{
			value.NewString("abc def ghi"),
			value.NewString(" "),
			value.NewNull(),
		},
		Result: value.NewNull(),
	},
	{
		Name: "ListElem Index is negative value",
		Function: parser.Function{
			Name: "list_elem",
		},
		Args: []value.Primary{
			value.NewString("abc def ghi"),
			value.NewString(" "),
			value.NewInteger(-1),
		},
		Result: value.NewNull(),
	},
	{
		Name: "ListElem Index does not exist",
		Function: parser.Function{
			Name: "list_elem",
		},
		Args: []value.Primary{
			value.NewString("abc def ghi"),
			value.NewString(" "),
			value.NewInteger(100),
		},
		Result: value.NewNull(),
	},
	{
		Name: "ListElem Arguments Error",
		Function: parser.Function{
			Name: "list_elem",
		},
		Args:  []value.Primary{},
		Error: "function list_elem takes exactly 3 arguments",
	},
}

func TestListElem(t *testing.T) {
	testFunction(t, ListElem, listElemTests)
}

var replaceFnTests = []functionTest{
	{
		Name: "Replace",
		Function: parser.Function{
			Name: "replace",
		},
		Args: []value.Primary{
			value.NewString("abcdefg abcdefg"),
			value.NewString("cd"),
			value.NewString("CD"),
		},
		Result: value.NewString("abCDefg abCDefg"),
	},
	{
		Name: "Replace String is Null",
		Function: parser.Function{
			Name: "replace",
		},
		Args: []value.Primary{
			value.NewNull(),
			value.NewString("cd"),
			value.NewString("CD"),
		},
		Result: value.NewNull(),
	},
	{
		Name: "Replace Old String is Null",
		Function: parser.Function{
			Name: "replace",
		},
		Args: []value.Primary{
			value.NewString("abcdefg abcdefg"),
			value.NewNull(),
			value.NewString("CD"),
		},
		Result: value.NewNull(),
	},
	{
		Name: "Replace New String is Null",
		Function: parser.Function{
			Name: "replace",
		},
		Args: []value.Primary{
			value.NewString("abcdefg abcdefg"),
			value.NewString("cd"),
			value.NewNull(),
		},
		Result: value.NewNull(),
	},
	{
		Name: "Replace Arguments Error",
		Function: parser.Function{
			Name: "replace",
		},
		Args:  []value.Primary{},
		Error: "function replace takes exactly 3 arguments",
	},
}

func TestReplaceFn(t *testing.T) {
	testFunction(t, ReplaceFn, replaceFnTests)
}

var formatTests = []functionTest{
	{
		Name: "Format",
		Function: parser.Function{
			Name: "format",
		},
		Args: []value.Primary{
			value.NewString("string = %q, integer = %q"),
			value.NewString("str"),
			value.NewInteger(1),
		},
		Result: value.NewString("string = 'str', integer = '1'"),
	},
	{
		Name: "Format Argument Length Error",
		Function: parser.Function{
			Name: "format",
		},
		Args:  []value.Primary{},
		Error: "function format takes at least 1 argument",
	},
	{
		Name: "Format Argument Value Error",
		Function: parser.Function{
			Name: "format",
		},
		Args: []value.Primary{
			value.NewNull(),
			value.NewString("str"),
			value.NewInteger(1),
		},
		Error: "the first argument must be a string for function format",
	},
	{
		Name: "Format Replace Holder Length Not Match Error",
		Function: parser.Function{
			Name: "format",
		},
		Args: []value.Primary{
			value.NewString("string = %s, integer = %s"),
			value.NewString("str"),
		},
		Error: "number of replace values does not match for function format",
	},
}

func TestFormat(t *testing.T) {
	testFunction(t, Format, formatTests)
}

var jsonValueTests = []functionTest{
	{
		Name: "JsonValue",
		Function: parser.Function{
			Name: "json_value",
		},
		Args: []value.Primary{
			value.NewString("key1.key2"),
			value.NewString("{\"key1\":{\"key2\":\"value\"}}"),
		},
		Result: value.NewString("value"),
	},
	{
		Name: "JsonValue Query is Null",
		Function: parser.Function{
			Name: "json_value",
		},
		Args: []value.Primary{
			value.NewNull(),
			value.NewString("{\"key1\":{\"key2\":\"value\"}}"),
		},
		Result: value.NewNull(),
	},
	{
		Name: "JsonValue Json-Text is Null",
		Function: parser.Function{
			Name: "json_value",
		},
		Args: []value.Primary{
			value.NewString("key1.key2"),
			value.NewNull(),
		},
		Result: value.NewNull(),
	},
	{
		Name: "JsonValue Arguments Error",
		Function: parser.Function{
			Name: "json_value",
		},
		Args: []value.Primary{
			value.NewString("key1.key2"),
		},
		Error: "function json_value takes exactly 2 arguments",
	},
	{
		Name: "JsonValue Json Loading Error",
		Function: parser.Function{
			Name: "json_value",
		},
		Args: []value.Primary{
			value.NewString("key1.key2"),
			value.NewString("{key1:{\"key2\":\"value\"}}"),
		},
		Error: "line 1, column 2: unexpected token \"key\" for function json_value",
	},
}

func TestJsonValue(t *testing.T) {
	testFunction(t, JsonValue, jsonValueTests)
}

var md5Tests = []functionTest{
	{
		Name: "Md5",
		Function: parser.Function{
			Name: "md5",
		},
		Args: []value.Primary{
			value.NewString("foo"),
		},
		Result: value.NewString("acbd18db4cc2f85cedef654fccc4a4d8"),
	},
	{
		Name: "Md5 Null",
		Function: parser.Function{
			Name: "md5",
		},
		Args: []value.Primary{
			value.NewNull(),
		},
		Result: value.NewNull(),
	},
	{
		Name: "Md5 Arguments Error",
		Function: parser.Function{
			Name: "md5",
		},
		Args:  []value.Primary{},
		Error: "function md5 takes exactly 1 argument",
	},
}

func TestMd5(t *testing.T) {
	testFunction(t, Md5, md5Tests)
}

var sha1Tests = []functionTest{
	{
		Name: "Sha1",
		Function: parser.Function{
			Name: "sha1",
		},
		Args: []value.Primary{
			value.NewString("foo"),
		},
		Result: value.NewString("0beec7b5ea3f0fdbc95d0dd47f3c5bc275da8a33"),
	},
}

func TestSha1(t *testing.T) {
	testFunction(t, Sha1, sha1Tests)
}

var sha256Tests = []functionTest{
	{
		Name: "Sha256",
		Function: parser.Function{
			Name: "sha256",
		},
		Args: []value.Primary{
			value.NewString("foo"),
		},
		Result: value.NewString("2c26b46b68ffc68ff99b453c1d30413413422d706483bfa0f98a5e886266e7ae"),
	},
}

func TestSha256(t *testing.T) {
	testFunction(t, Sha256, sha256Tests)
}

var sha512Tests = []functionTest{
	{
		Name: "Sha512",
		Function: parser.Function{
			Name: "sha512",
		},
		Args: []value.Primary{
			value.NewString("foo"),
		},
		Result: value.NewString("f7fbba6e0636f890e56fbbf3283e524c6fa3204ae298382d624741d0dc6638326e282c41be5e4254d8820772c5518a2c5a8c0c7f7eda19594a7eb539453e1ed7"),
	},
}

func TestSha512(t *testing.T) {
	testFunction(t, Sha512, sha512Tests)
}

var md5HmacTests = []functionTest{
	{
		Name: "Md5Hmac",
		Function: parser.Function{
			Name: "md5_hmac",
		},
		Args: []value.Primary{
			value.NewString("foo"),
			value.NewString("bar"),
		},
		Result: value.NewString("31b6db9e5eb4addb42f1a6ca07367adc"),
	},
	{
		Name: "Md5Hmac String is Null",
		Function: parser.Function{
			Name: "md5_hmac",
		},
		Args: []value.Primary{
			value.NewNull(),
			value.NewString("bar"),
		},
		Result: value.NewNull(),
	},
	{
		Name: "Md5Hmac Key is Null",
		Function: parser.Function{
			Name: "md5_hmac",
		},
		Args: []value.Primary{
			value.NewString("foo"),
			value.NewNull(),
		},
		Result: value.NewNull(),
	},
	{
		Name: "Md5Hmac Arguments Error",
		Function: parser.Function{
			Name: "md5_hmac",
		},
		Args:  []value.Primary{},
		Error: "function md5_hmac takes exactly 2 arguments",
	},
}

func TestMd5Hmac(t *testing.T) {
	testFunction(t, Md5Hmac, md5HmacTests)
}

var sha1HmacTests = []functionTest{
	{
		Name: "Sha1Hmac",
		Function: parser.Function{
			Name: "sha1_hmac",
		},
		Args: []value.Primary{
			value.NewString("foo"),
			value.NewString("bar"),
		},
		Result: value.NewString("85d155c55ed286a300bd1cf124de08d87e914f3a"),
	},
}

func TestSha1Hmac(t *testing.T) {
	testFunction(t, Sha1Hmac, sha1HmacTests)
}

var sha256HmacTests = []functionTest{
	{
		Name: "Sha256Hmac",
		Function: parser.Function{
			Name: "sha256_hmac",
		},
		Args: []value.Primary{
			value.NewString("foo"),
			value.NewString("bar"),
		},
		Result: value.NewString("147933218aaabc0b8b10a2b3a5c34684c8d94341bcf10a4736dc7270f7741851"),
	},
}

func TestSha256Hmac(t *testing.T) {
	testFunction(t, Sha256Hmac, sha256HmacTests)
}

var sha512HmacTests = []functionTest{
	{
		Name: "Sha512Hmac",
		Function: parser.Function{
			Name: "sha512_hmac",
		},
		Args: []value.Primary{
			value.NewString("foo"),
			value.NewString("bar"),
		},
		Result: value.NewString("24257d7210582a65c731ec55159c8184cc24c02489453e58587f71f44c23a2d61b4b72154a89d17b2d49448a8452ea066f4fc56a2bcead45c088572ffccdb3d8"),
	},
}

func TestSha512Hmac(t *testing.T) {
	testFunction(t, Sha512Hmac, sha512HmacTests)
}

var datetimeFormatTests = []functionTest{
	{
		Name: "DatetimeFormat",
		Function: parser.Function{
			Name: "datetime_format",
		},
		Args: []value.Primary{
			value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 0, GetTestLocation())),
			value.NewString("%Y-%m-%d"),
		},
		Result: value.NewString("2012-02-03"),
	},
	{
		Name: "DatetimeFormat Datetime is Null",
		Function: parser.Function{
			Name: "datetime_format",
		},
		Args: []value.Primary{
			value.NewNull(),
			value.NewString("%Y-%m-%d"),
		},
		Result: value.NewNull(),
	},
	{
		Name: "DatetimeFormat Format is Null",
		Function: parser.Function{
			Name: "datetime_format",
		},
		Args: []value.Primary{
			value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 0, GetTestLocation())),
			value.NewNull(),
		},
		Result: value.NewNull(),
	},
	{
		Name: "DatetimeFormat Arguments Error",
		Function: parser.Function{
			Name: "datetime_format",
		},
		Args: []value.Primary{
			value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 0, GetTestLocation())),
		},
		Error: "function datetime_format takes exactly 2 arguments",
	},
}

func TestDatetimeFormat(t *testing.T) {
	testFunction(t, DatetimeFormat, datetimeFormatTests)
}

var yearTests = []functionTest{
	{
		Name: "Year",
		Function: parser.Function{
			Name: "year",
		},
		Args: []value.Primary{
			value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 0, GetTestLocation())),
		},
		Result: value.NewInteger(2012),
	},
	{
		Name: "Year Datetime is Null",
		Function: parser.Function{
			Name: "year",
		},
		Args: []value.Primary{
			value.NewNull(),
		},
		Result: value.NewNull(),
	},
	{
		Name: "Year Arguments Error",
		Function: parser.Function{
			Name: "year",
		},
		Args:  []value.Primary{},
		Error: "function year takes exactly 1 argument",
	},
}

func TestYear(t *testing.T) {
	testFunction(t, Year, yearTests)
}

var monthTests = []functionTest{
	{
		Name: "Month",
		Function: parser.Function{
			Name: "month",
		},
		Args: []value.Primary{
			value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 0, GetTestLocation())),
		},
		Result: value.NewInteger(2),
	},
}

func TestMonth(t *testing.T) {
	testFunction(t, Month, monthTests)
}

var dayTests = []functionTest{
	{
		Name: "Day",
		Function: parser.Function{
			Name: "day",
		},
		Args: []value.Primary{
			value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 0, GetTestLocation())),
		},
		Result: value.NewInteger(3),
	},
}

func TestDay(t *testing.T) {
	testFunction(t, Day, dayTests)
}

var hourTests = []functionTest{
	{
		Name: "Hour",
		Function: parser.Function{
			Name: "hour",
		},
		Args: []value.Primary{
			value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 0, GetTestLocation())),
		},
		Result: value.NewInteger(9),
	},
}

func TestHour(t *testing.T) {
	testFunction(t, Hour, hourTests)
}

var minuteTests = []functionTest{
	{
		Name: "Minute",
		Function: parser.Function{
			Name: "minute",
		},
		Args: []value.Primary{
			value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 0, GetTestLocation())),
		},
		Result: value.NewInteger(18),
	},
}

func TestMinute(t *testing.T) {
	testFunction(t, Minute, minuteTests)
}

var secondTests = []functionTest{
	{
		Name: "Second",
		Function: parser.Function{
			Name: "second",
		},
		Args: []value.Primary{
			value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 0, GetTestLocation())),
		},
		Result: value.NewInteger(15),
	},
}

func TestSecond(t *testing.T) {
	testFunction(t, Second, secondTests)
}

var millisecondTests = []functionTest{
	{
		Name: "Millisecond",
		Function: parser.Function{
			Name: "millisecond",
		},
		Args: []value.Primary{
			value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456789, GetTestLocation())),
		},
		Result: value.NewInteger(123),
	},
}

func TestMillisecond(t *testing.T) {
	testFunction(t, Millisecond, millisecondTests)
}

var microsecondTests = []functionTest{
	{
		Name: "Microsecond",
		Function: parser.Function{
			Name: "microsecond",
		},
		Args: []value.Primary{
			value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456789, GetTestLocation())),
		},
		Result: value.NewInteger(123457),
	},
}

func TestMicrosecond(t *testing.T) {
	testFunction(t, Microsecond, microsecondTests)
}

var nanosecondTests = []functionTest{
	{
		Name: "Nanosecond",
		Function: parser.Function{
			Name: "nanosecond",
		},
		Args: []value.Primary{
			value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456789, GetTestLocation())),
		},
		Result: value.NewInteger(123456789),
	},
}

func TestNanosecond(t *testing.T) {
	testFunction(t, Nanosecond, nanosecondTests)
}

var weekdayTests = []functionTest{
	{
		Name: "Weekday",
		Function: parser.Function{
			Name: "weekday",
		},
		Args: []value.Primary{
			value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 0, GetTestLocation())),
		},
		Result: value.NewInteger(5),
	},
}

func TestWeekday(t *testing.T) {
	testFunction(t, Weekday, weekdayTests)
}

var unixTimeTests = []functionTest{
	{
		Name: "UnixTime",
		Function: parser.Function{
			Name: "unix_time",
		},
		Args: []value.Primary{
			value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 0, GetTestLocation())),
		},
		Result: value.NewInteger(1328260695),
	},
}

func TestUnixTime(t *testing.T) {
	testFunction(t, UnixTime, unixTimeTests)
}

var unixNanoTimeTests = []functionTest{
	{
		Name: "UnixNanoTime",
		Function: parser.Function{
			Name: "unix_nano_time",
		},
		Args: []value.Primary{
			value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456789, GetTestLocation())),
		},
		Result: value.NewInteger(1328260695123456789),
	},
}

func TestUnixNanoTime(t *testing.T) {
	testFunction(t, UnixNanoTime, unixNanoTimeTests)
}

var dayOfYearTests = []functionTest{
	{
		Name: "DayOfYear",
		Function: parser.Function{
			Name: "day_of_year",
		},
		Args: []value.Primary{
			value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 0, GetTestLocation())),
		},
		Result: value.NewInteger(34),
	},
}

func TestDayOfYear(t *testing.T) {
	testFunction(t, DayOfYear, dayOfYearTests)
}

var weekOfYearTests = []functionTest{
	{
		Name: "WeekOfYear",
		Function: parser.Function{
			Name: "week_of_year",
		},
		Args: []value.Primary{
			value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 0, GetTestLocation())),
		},
		Result: value.NewInteger(5),
	},
}

func TestWeekOfYear(t *testing.T) {
	testFunction(t, WeekOfYear, weekOfYearTests)
}

var addYearTests = []functionTest{
	{
		Name: "AddYear",
		Function: parser.Function{
			Name: "add_year",
		},
		Args: []value.Primary{
			value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456789, GetTestLocation())),
			value.NewInteger(2),
		},
		Result: value.NewDatetime(time.Date(2014, 2, 3, 9, 18, 15, 123456789, GetTestLocation())),
	},
	{
		Name: "AddYear Datetime is Null",
		Function: parser.Function{
			Name: "add_year",
		},
		Args: []value.Primary{
			value.NewNull(),
			value.NewInteger(2),
		},
		Result: value.NewNull(),
	},
	{
		Name: "AddYear Duration is Null",
		Function: parser.Function{
			Name: "add_year",
		},
		Args: []value.Primary{
			value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456789, GetTestLocation())),
			value.NewNull(),
		},
		Result: value.NewNull(),
	},
	{
		Name: "AddYear Arguments Error",
		Function: parser.Function{
			Name: "add_year",
		},
		Args:  []value.Primary{},
		Error: "function add_year takes exactly 2 arguments",
	},
}

func TestAddYear(t *testing.T) {
	testFunction(t, AddYear, addYearTests)
}

var addMonthTests = []functionTest{
	{
		Name: "AddMonth",
		Function: parser.Function{
			Name: "add_month",
		},
		Args: []value.Primary{
			value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456789, GetTestLocation())),
			value.NewInteger(2),
		},
		Result: value.NewDatetime(time.Date(2012, 4, 3, 9, 18, 15, 123456789, GetTestLocation())),
	},
}

func TestAddMonth(t *testing.T) {
	testFunction(t, AddMonth, addMonthTests)
}

var addDayTests = []functionTest{
	{
		Name: "AddDay",
		Function: parser.Function{
			Name: "add_day",
		},
		Args: []value.Primary{
			value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456789, GetTestLocation())),
			value.NewInteger(2),
		},
		Result: value.NewDatetime(time.Date(2012, 2, 5, 9, 18, 15, 123456789, GetTestLocation())),
	},
}

func TestAddDay(t *testing.T) {
	testFunction(t, AddDay, addDayTests)
}

var addHourTests = []functionTest{
	{
		Name: "AddHour",
		Function: parser.Function{
			Name: "add_hour",
		},
		Args: []value.Primary{
			value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456789, GetTestLocation())),
			value.NewInteger(2),
		},
		Result: value.NewDatetime(time.Date(2012, 2, 3, 11, 18, 15, 123456789, GetTestLocation())),
	},
}

func TestAddHour(t *testing.T) {
	testFunction(t, AddHour, addHourTests)
}

var addMinuteTests = []functionTest{
	{
		Name: "AddMinute",
		Function: parser.Function{
			Name: "add_minute",
		},
		Args: []value.Primary{
			value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456789, GetTestLocation())),
			value.NewInteger(2),
		},
		Result: value.NewDatetime(time.Date(2012, 2, 3, 9, 20, 15, 123456789, GetTestLocation())),
	},
}

func TestAddMinute(t *testing.T) {
	testFunction(t, AddMinute, addMinuteTests)
}

var addSecondTests = []functionTest{
	{
		Name: "AddSecond",
		Function: parser.Function{
			Name: "add_second",
		},
		Args: []value.Primary{
			value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456789, GetTestLocation())),
			value.NewInteger(2),
		},
		Result: value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 17, 123456789, GetTestLocation())),
	},
}

func TestAddSecond(t *testing.T) {
	testFunction(t, AddSecond, addSecondTests)
}

var addMilliTests = []functionTest{
	{
		Name: "AddMilli",
		Function: parser.Function{
			Name: "add_milli",
		},
		Args: []value.Primary{
			value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456789, GetTestLocation())),
			value.NewInteger(2),
		},
		Result: value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 125456789, GetTestLocation())),
	},
}

func TestAddMilli(t *testing.T) {
	testFunction(t, AddMilli, addMilliTests)
}

var addMicroTests = []functionTest{
	{
		Name: "AddMicro",
		Function: parser.Function{
			Name: "add_micro",
		},
		Args: []value.Primary{
			value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456789, GetTestLocation())),
			value.NewInteger(2),
		},
		Result: value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123458789, GetTestLocation())),
	},
}

func TestAddMicro(t *testing.T) {
	testFunction(t, AddMicro, addMicroTests)
}

var addNanoTests = []functionTest{
	{
		Name: "AddNano",
		Function: parser.Function{
			Name: "add_nano",
		},
		Args: []value.Primary{
			value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456789, GetTestLocation())),
			value.NewInteger(2),
		},
		Result: value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456791, GetTestLocation())),
	},
}

func TestAddNano(t *testing.T) {
	testFunction(t, AddNano, addNanoTests)
}

var truncMonthTests = []functionTest{
	{
		Name: "TruncMonth",
		Function: parser.Function{
			Name: "trunc_month",
		},
		Args: []value.Primary{
			value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456789, GetTestLocation())),
		},
		Result: value.NewDatetime(time.Date(2012, 1, 1, 0, 0, 0, 0, GetTestLocation())),
	},
	{
		Name: "TruncMonth Argument Error",
		Function: parser.Function{
			Name: "trunc_month",
		},
		Args:  []value.Primary{},
		Error: "function trunc_month takes exactly 1 argument",
	},
	{
		Name: "TruncMonth Argument Is Null",
		Function: parser.Function{
			Name: "trunc_month",
		},
		Args: []value.Primary{
			value.NewNull(),
		},
		Result: value.NewNull(),
	},
}

func TestTruncMonth(t *testing.T) {
	testFunction(t, TruncMonth, truncMonthTests)
}

var truncDayTests = []functionTest{
	{
		Name: "TruncDay",
		Function: parser.Function{
			Name: "trunc_day",
		},
		Args: []value.Primary{
			value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456789, GetTestLocation())),
		},
		Result: value.NewDatetime(time.Date(2012, 2, 1, 0, 0, 0, 0, GetTestLocation())),
	},
}

func TestTruncDay(t *testing.T) {
	testFunction(t, TruncDay, truncDayTests)
}

var truncTimeTests = []functionTest{
	{
		Name: "TruncTime",
		Function: parser.Function{
			Name: "trunc_time",
		},
		Args: []value.Primary{
			value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456789, GetTestLocation())),
		},
		Result: value.NewDatetime(time.Date(2012, 2, 3, 0, 0, 0, 0, GetTestLocation())),
	},
}

func TestTruncTime(t *testing.T) {
	testFunction(t, TruncTime, truncTimeTests)
}

var truncMinuteTests = []functionTest{
	{
		Name: "TruncMinute",
		Function: parser.Function{
			Name: "trunc_minute",
		},
		Args: []value.Primary{
			value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456789, GetTestLocation())),
		},
		Result: value.NewDatetime(time.Date(2012, 2, 3, 9, 0, 0, 0, GetTestLocation())),
	},
	{
		Name: "TruncMinute Argument Error",
		Function: parser.Function{
			Name: "trunc_minute",
		},
		Args:  []value.Primary{},
		Error: "function trunc_minute takes exactly 1 argument",
	},
	{
		Name: "TruncMinute Argument Is Null",
		Function: parser.Function{
			Name: "trunc_minute",
		},
		Args: []value.Primary{
			value.NewNull(),
		},
		Result: value.NewNull(),
	},
}

func TestTruncMinute(t *testing.T) {
	testFunction(t, TruncMinute, truncMinuteTests)
}

var truncSecondTests = []functionTest{
	{
		Name: "TruncSecond",
		Function: parser.Function{
			Name: "trunc_second",
		},
		Args: []value.Primary{
			value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456789, GetTestLocation())),
		},
		Result: value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 0, 0, GetTestLocation())),
	},
}

func TestTruncSecond(t *testing.T) {
	testFunction(t, TruncSecond, truncSecondTests)
}

var truncMilliTests = []functionTest{
	{
		Name: "TruncMilli",
		Function: parser.Function{
			Name: "trunc_milli",
		},
		Args: []value.Primary{
			value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456789, GetTestLocation())),
		},
		Result: value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 0, GetTestLocation())),
	},
}

func TestTruncMilli(t *testing.T) {
	testFunction(t, TruncMilli, truncMilliTests)
}

var truncMicroTests = []functionTest{
	{
		Name: "TruncMicro",
		Function: parser.Function{
			Name: "trunc_micro",
		},
		Args: []value.Primary{
			value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456789, GetTestLocation())),
		},
		Result: value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123000000, GetTestLocation())),
	},
}

func TestTruncateMicro(t *testing.T) {
	testFunction(t, TruncMicro, truncMicroTests)
}

var truncNanoTests = []functionTest{
	{
		Name: "TruncNano",
		Function: parser.Function{
			Name: "trunc_nano",
		},
		Args: []value.Primary{
			value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456789, GetTestLocation())),
		},
		Result: value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456000, GetTestLocation())),
	},
}

func TestTruncateNano(t *testing.T) {
	testFunction(t, TruncNano, truncNanoTests)
}

var dateDiffTests = []functionTest{
	{
		Name: "DateDiff",
		Function: parser.Function{
			Name: "date_diff",
		},
		Args: []value.Primary{
			value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456789, GetTestLocation())),
			value.NewDatetime(time.Date(2012, 2, 5, 1, 18, 55, 123456789, GetTestLocation())),
		},
		Result: value.NewInteger(-2),
	},
	{
		Name: "DateDiff Datetime1 is Null",
		Function: parser.Function{
			Name: "date_diff",
		},
		Args: []value.Primary{
			value.NewNull(),
			value.NewDatetime(time.Date(2012, 2, 5, 1, 18, 55, 123456789, GetTestLocation())),
		},
		Result: value.NewNull(),
	},
	{
		Name: "DateDiff Datetime2 is Null",
		Function: parser.Function{
			Name: "date_diff",
		},
		Args: []value.Primary{
			value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456789, GetTestLocation())),
			value.NewNull(),
		},
		Result: value.NewNull(),
	},
	{
		Name: "DateDiff Arguments Error",
		Function: parser.Function{
			Name: "date_diff",
		},
		Args:  []value.Primary{},
		Error: "function date_diff takes exactly 2 arguments",
	},
}

func TestDateDiff(t *testing.T) {
	testFunction(t, DateDiff, dateDiffTests)
}

var timeDiffTests = []functionTest{
	{
		Name: "TimeDiff",
		Function: parser.Function{
			Name: "time_diff",
		},
		Args: []value.Primary{
			value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456789, GetTestLocation())),
			value.NewDatetime(time.Date(2012, 2, 3, 1, 18, 55, 123000000, GetTestLocation())),
		},
		Result: value.NewFloat(28760.000456789),
	},
	{
		Name: "TimeDiff Datetime1 is Null",
		Function: parser.Function{
			Name: "time_diff",
		},
		Args: []value.Primary{
			value.NewNull(),
			value.NewDatetime(time.Date(2012, 2, 5, 1, 18, 55, 123456789, GetTestLocation())),
		},
		Result: value.NewNull(),
	},
	{
		Name: "TimeDiff Datetime2 is Null",
		Function: parser.Function{
			Name: "time_diff",
		},
		Args: []value.Primary{
			value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456789, GetTestLocation())),
			value.NewNull(),
		},
		Result: value.NewNull(),
	},
	{
		Name: "TimeDiff Arguments Error",
		Function: parser.Function{
			Name: "time_diff",
		},
		Args:  []value.Primary{},
		Error: "function time_diff takes exactly 2 arguments",
	},
}

func TestTimeDiff(t *testing.T) {
	testFunction(t, TimeDiff, timeDiffTests)
}

var timeNanoDiffTests = []functionTest{
	{
		Name: "TimeNanoDiff",
		Function: parser.Function{
			Name: "time_nano_diff",
		},
		Args: []value.Primary{
			value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456789, GetTestLocation())),
			value.NewDatetime(time.Date(2012, 2, 3, 1, 18, 55, 123000000, GetTestLocation())),
		},
		Result: value.NewInteger(28760000456789),
	},
}

func TestTimeNanoDiff(t *testing.T) {
	testFunction(t, TimeNanoDiff, timeNanoDiffTests)
}

var utcTests = []functionTest{
	{
		Name: "UTC",
		Function: parser.Function{
			Name: "utc",
		},
		Args: []value.Primary{
			value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456789, GetTestLocation())),
		},
		Result: value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456789, time.UTC)),
	},
	{
		Name: "UTC Argument Error",
		Function: parser.Function{
			Name: "utc",
		},
		Args:  []value.Primary{},
		Error: "function utc takes exactly 1 argument",
	},
	{
		Name: "UTC Argument Is Null",
		Function: parser.Function{
			Name: "utc",
		},
		Args: []value.Primary{
			value.NewNull(),
		},
		Result: value.NewNull(),
	},
}

func TestUTC(t *testing.T) {
	testFunction(t, UTC, utcTests)
}

var nanoToDatetimeTests = []functionTest{
	{
		Name: "NanoToDatetime",
		Function: parser.Function{
			Name: "nano_to_datetime",
		},
		Args: []value.Primary{
			value.NewInteger(1328260695000000001),
		},
		Result: value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 1, GetTestLocation())),
	},
	{
		Name: "NanoToDatetime Invalid Argument",
		Function: parser.Function{
			Name: "nano_to_datetime",
		},
		Args: []value.Primary{
			value.NewString("abc"),
		},
		Result: value.NewNull(),
	},
	{
		Name: "NanoToDatetime Arguments Error",
		Function: parser.Function{
			Name: "nano_to_datetime",
		},
		Args:  []value.Primary{},
		Error: "function nano_to_datetime takes exactly 1 argument",
	},
}

func TestNanoToDatetime(t *testing.T) {
	testFunction(t, NanoToDatetime, nanoToDatetimeTests)
}

var stringTests = []functionTest{
	{
		Name: "String from Integer",
		Function: parser.Function{
			Name: "string",
		},
		Args: []value.Primary{
			value.NewInteger(2),
		},
		Result: value.NewString("2"),
	},
	{
		Name: "String from Boolean",
		Function: parser.Function{
			Name: "string",
		},
		Args: []value.Primary{
			value.NewBoolean(true),
		},
		Result: value.NewString("true"),
	},
	{
		Name: "String from Ternary",
		Function: parser.Function{
			Name: "string",
		},
		Args: []value.Primary{
			value.NewTernary(ternary.TRUE),
		},
		Result: value.NewString("TRUE"),
	},
	{
		Name: "String from Datetime",
		Function: parser.Function{
			Name: "string",
		},
		Args: []value.Primary{
			value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 0, GetTestLocation())),
		},
		Result: value.NewString("2012-02-03T09:18:15Z"),
	},
	{
		Name: "String Arguments Error",
		Function: parser.Function{
			Name: "string",
		},
		Args:  []value.Primary{},
		Error: "function string takes exactly 1 argument",
	},
}

func TestString(t *testing.T) {
	testFunction(t, String, stringTests)
}

var integerTests = []functionTest{
	{
		Name: "Integer from Integer",
		Function: parser.Function{
			Name: "integer",
		},
		Args: []value.Primary{
			value.NewInteger(2),
		},
		Result: value.NewInteger(2),
	},
	{
		Name: "Integer from String",
		Function: parser.Function{
			Name: "integer",
		},
		Args: []value.Primary{
			value.NewString("2"),
		},
		Result: value.NewInteger(2),
	},
	{
		Name: "Integer from E-Notation",
		Function: parser.Function{
			Name: "integer",
		},
		Args: []value.Primary{
			value.NewString("2e+02"),
		},
		Result: value.NewInteger(200),
	},
	{
		Name: "Integer from Float",
		Function: parser.Function{
			Name: "integer",
		},
		Args: []value.Primary{
			value.NewFloat(1.7),
		},
		Result: value.NewInteger(2),
	},
	{
		Name: "Float Null",
		Function: parser.Function{
			Name: "float",
		},
		Args: []value.Primary{
			value.NewNull(),
		},
		Result: value.NewNull(),
	},
	{
		Name: "Integer from Datetime",
		Function: parser.Function{
			Name: "integer",
		},
		Args: []value.Primary{
			value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 0, GetTestLocation())),
		},
		Result: value.NewInteger(1328260695),
	},
	{
		Name: "Integer Arguments Error",
		Function: parser.Function{
			Name: "integer",
		},
		Args:  []value.Primary{},
		Error: "function integer takes exactly 1 argument",
	},
}

func TestInteger(t *testing.T) {
	testFunction(t, Integer, integerTests)
}

var floatTests = []functionTest{
	{
		Name: "Float from String",
		Function: parser.Function{
			Name: "float",
		},
		Args: []value.Primary{
			value.NewString("2"),
		},
		Result: value.NewFloat(2),
	},
	{
		Name: "Float from Datetime",
		Function: parser.Function{
			Name: "float",
		},
		Args: []value.Primary{
			value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123450000, GetTestLocation())),
		},
		Result: value.NewFloat(1328260695.12345),
	},
	{
		Name: "Float Arguments Error",
		Function: parser.Function{
			Name: "float",
		},
		Args:  []value.Primary{},
		Error: "function float takes exactly 1 argument",
	},
}

func TestFloat(t *testing.T) {
	testFunction(t, Float, floatTests)
}

var booleanTests = []functionTest{
	{
		Name: "Boolean from String",
		Function: parser.Function{
			Name: "boolean",
		},
		Args: []value.Primary{
			value.NewString("true"),
		},
		Result: value.NewBoolean(true),
	},
	{
		Name: "Boolean Arguments Error",
		Function: parser.Function{
			Name: "boolean",
		},
		Args:  []value.Primary{},
		Error: "function boolean takes exactly 1 argument",
	},
}

func TestBoolean(t *testing.T) {
	testFunction(t, Boolean, booleanTests)
}

var ternaryTest = []functionTest{
	{
		Name: "Ternary from String",
		Function: parser.Function{
			Name: "ternary",
		},
		Args: []value.Primary{
			value.NewString("true"),
		},
		Result: value.NewTernary(ternary.TRUE),
	},
	{
		Name: "Ternary Arguments Error",
		Function: parser.Function{
			Name: "ternary",
		},
		Args:  []value.Primary{},
		Error: "function ternary takes exactly 1 argument",
	},
}

func TestTernary(t *testing.T) {
	testFunction(t, Ternary, ternaryTest)
}

var datetimeTests = []functionTest{
	{
		Name: "Datetime",
		Function: parser.Function{
			Name: "datetime",
		},
		Args: []value.Primary{
			value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 0, GetTestLocation())),
		},
		Result: value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 0, GetTestLocation())),
	},
	{
		Name: "Datetime from Integer",
		Function: parser.Function{
			Name: "datetime",
		},
		Args: []value.Primary{
			value.NewInteger(1136181845),
		},
		Result: value.NewDatetime(time.Date(2006, 1, 2, 6, 4, 5, 0, GetTestLocation())),
	},
	{
		Name: "Datetime from Float",
		Function: parser.Function{
			Name: "datetime",
		},
		Args: []value.Primary{
			value.NewFloat(1136181845.123),
		},
		Result: value.NewDatetime(time.Date(2006, 1, 2, 6, 4, 5, 123000000, GetTestLocation())),
	},
	{
		Name: "Datetime from String",
		Function: parser.Function{
			Name: "datetime",
		},
		Args: []value.Primary{
			value.NewString("2012-02-03 09:18:15"),
		},
		Result: value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 0, GetTestLocation())),
	},
	{
		Name: "Datetime from String representing Integer",
		Function: parser.Function{
			Name: "datetime",
		},
		Args: []value.Primary{
			value.NewString("1136181845"),
		},
		Result: value.NewDatetime(time.Date(2006, 1, 2, 6, 4, 5, 0, GetTestLocation())),
	},
	{
		Name: "Datetime from String representing Float",
		Function: parser.Function{
			Name: "datetime",
		},
		Args: []value.Primary{
			value.NewString("1136181845.123"),
		},
		Result: value.NewDatetime(time.Date(2006, 1, 2, 6, 4, 5, 123000000, GetTestLocation())),
	},
	{
		Name: "Datetime from String with time zone conversion",
		Function: parser.Function{
			Name: "datetime",
		},
		Args: []value.Primary{
			value.NewString("2012-02-03T09:18:15-07:00"),
			value.NewString(TestLocation),
		},
		Result: value.NewDatetime(time.Date(2012, 2, 3, 16, 18, 15, 0, GetTestLocation())),
	},
	{
		Name: "Datetime Invalid String",
		Function: parser.Function{
			Name: "datetime",
		},
		Args: []value.Primary{
			value.NewString("abcde"),
		},
		Result: value.NewNull(),
	},
	{
		Name: "Datetime Arguments Error",
		Function: parser.Function{
			Name: "datetime",
		},
		Args:  []value.Primary{},
		Error: "function datetime takes 1 or 2 arguments",
	},
	{
		Name: "Datetime Second Argument Not String Error",
		Function: parser.Function{
			Name: "datetime",
		},
		Args: []value.Primary{
			value.NewString("2012-02-03T09:18:15-07:00"),
			value.NewNull(),
		},
		Error: "failed to load time zone NULL for function datetime",
	},
	{
		Name: "Datetime Second Argument Invalid Location Error",
		Function: parser.Function{
			Name: "datetime",
		},
		Args: []value.Primary{
			value.NewString("2012-02-03T09:18:15-07:00"),
			value.NewString("Err"),
		},
		Error: "failed to load time zone 'Err' for function datetime",
	},
}

func TestDatetime(t *testing.T) {
	testFunction(t, Datetime, datetimeTests)
}

var callTests = []functionTest{
	{
		Name: "Call Argument Error",
		Function: parser.Function{
			Name: "call",
		},
		Args:  []value.Primary{},
		Error: "function call takes at least 1 argument",
	},
	{
		Name: "Call Command Error",
		Function: parser.Function{
			Name: "call",
		},
		Args: []value.Primary{
			value.NewString("notexistcommand"),
		},
		Error: "environment-dependent",
	},
}

func TestCall(t *testing.T) {
	ctx := context.Background()
	for _, v := range callTests {
		result, err := Call(ctx, v.Function, v.Args)
		if err != nil {
			if len(v.Error) < 1 {
				t.Errorf("%s: unexpected error %q", v.Name, err)
			} else if v.Error != "environment-dependent" && err.Error() != v.Error {
				t.Errorf("%s: error %q, want error %q", v.Name, err.Error(), v.Error)
			}
			continue
		}
		if 0 < len(v.Error) {
			t.Errorf("%s: no error, want error %q", v.Name, v.Error)
			continue
		}
		if !reflect.DeepEqual(result, v.Result) {
			t.Errorf("%s: result = %s, want %s", v.Name, result, v.Result)
		}
	}
}

var nowTests = []struct {
	Name     string
	Function parser.Function
	Args     []value.Primary
	Scope    *ReferenceScope
	Result   value.Primary
	Error    string
}{
	{
		Name: "Now From Current Time",
		Function: parser.Function{
			Name: "now",
		},
		Scope:  NewReferenceScope(TestTx),
		Result: value.NewDatetime(NowForTest),
	},
	{
		Name: "Now From Filter",
		Function: parser.Function{
			Name: "now",
		},
		Scope:  GenerateReferenceScope(nil, nil, time.Date(2013, 2, 3, 0, 0, 0, 0, GetTestLocation()), nil),
		Result: value.NewDatetime(time.Date(2013, 2, 3, 0, 0, 0, 0, GetTestLocation())),
	},
	{
		Name: "Now Arguments Error",
		Function: parser.Function{
			Name: "now",
		},
		Args: []value.Primary{
			value.NewInteger(1),
		},
		Scope: NewReferenceScope(TestTx),
		Error: "function now takes no argument",
	},
}

func TestNow(t *testing.T) {
	initFlag(TestTx.Flags)
	for _, v := range nowTests {
		result, err := Now(v.Scope, v.Function, v.Args)
		if err != nil {
			if len(v.Error) < 1 {
				t.Errorf("%s: unexpected error %q", v.Name, err)
			} else if err.Error() != v.Error {
				t.Errorf("%s: error %q, want error %q", v.Name, err.Error(), v.Error)
			}
			continue
		}
		if 0 < len(v.Error) {
			t.Errorf("%s: no error, want error %q", v.Name, v.Error)
			continue
		}
		if !reflect.DeepEqual(result, v.Result) {
			t.Errorf("%s: result = %s, want %s", v.Name, result, v.Result)
		}
	}
}

var jsonObjectTests = []struct {
	Name     string
	Function parser.Function
	Scope    *ReferenceScope
	Result   value.Primary
	Error    string
}{
	{
		Name: "Json Object",
		Function: parser.Function{
			Name: "json_object",
			Args: []parser.QueryExpression{
				parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
			},
		},
		Scope: GenerateReferenceScope(nil, nil, time.Time{}, []ReferenceRecord{
			{
				view: &View{
					Header: NewHeaderWithId("table1", []string{"column1", "column2"}),
					RecordSet: []Record{
						NewRecordWithId(0, []value.Primary{value.NewInteger(1), value.NewInteger(2)}),
						NewRecordWithId(1, []value.Primary{value.NewInteger(11), value.NewInteger(12)}),
					},
				},
				recordIndex: 1,
				cache:       NewFieldIndexCache(10, LimitToUseFieldIndexSliceChache),
			},
		}),
		Result: value.NewString("{\"column1\":11}"),
	},
	{
		Name: "Json Object with All Columns",
		Function: parser.Function{
			Name: "json_object",
		},
		Scope: GenerateReferenceScope(nil, nil, time.Time{}, []ReferenceRecord{
			{
				view: &View{
					Header: NewHeaderWithId("table1", []string{"column1", "column2.child1"}),
					RecordSet: []Record{
						NewRecordWithId(0, []value.Primary{value.NewInteger(1), value.NewInteger(2)}),
						NewRecordWithId(1, []value.Primary{value.NewInteger(11), value.NewInteger(12)}),
					},
				},
				recordIndex: 1,
				cache:       NewFieldIndexCache(10, LimitToUseFieldIndexSliceChache),
			},
		}),
		Result: value.NewString("{\"column1\":11,\"column2\":{\"child1\":12}}"),
	},
	{
		Name: "Json Object Unpermitted Statement Error",
		Function: parser.Function{
			Name: "json_object",
			Args: []parser.QueryExpression{
				parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
			},
		},
		Scope:  NewReferenceScope(TestTx),
		Result: value.NewNull(),
	},
	{
		Name: "Json Object Path Error",
		Function: parser.Function{
			Name: "json_object",
		},
		Scope: GenerateReferenceScope(nil, nil, time.Time{}, []ReferenceRecord{
			{
				view: &View{
					Header: NewHeaderWithId("table1", []string{"column1", "column2.."}),
					RecordSet: []Record{
						NewRecordWithId(0, []value.Primary{value.NewInteger(1), value.NewInteger(2)}),
						NewRecordWithId(1, []value.Primary{value.NewInteger(11), value.NewInteger(12)}),
					},
				},
				recordIndex: 1,
				cache:       NewFieldIndexCache(10, LimitToUseFieldIndexSliceChache),
			},
		}),
		Error: "unexpected token \".\" at column 9 in \"column2..\" for function json_object",
	},
}

func TestJsonObject(t *testing.T) {
	for _, v := range jsonObjectTests {
		if v.Scope == nil {
			v.Scope = NewReferenceScope(TestTx)
		}

		result, err := JsonObject(context.Background(), v.Scope, v.Function)
		if err != nil {
			if len(v.Error) < 1 {
				t.Errorf("%s: unexpected error %q", v.Name, err)
			} else if err.Error() != v.Error {
				t.Errorf("%s: error %q, want error %q", v.Name, err.Error(), v.Error)
			}
			continue
		}
		if 0 < len(v.Error) {
			t.Errorf("%s: no error, want error %q", v.Name, v.Error)
			continue
		}
		if !reflect.DeepEqual(result, v.Result) {
			t.Errorf("%s: result = %s, want %s", v.Name, result, v.Result)
		}
	}
}
