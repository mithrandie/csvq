package query

import (
	"reflect"
	"testing"
	"time"

	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/ternary"
)

type functionTest struct {
	Name   string
	Args   []parser.Primary
	Result parser.Primary
	Error  string
}

func testFunction(t *testing.T, f func([]parser.Primary) (parser.Primary, error), tests []functionTest) {
	for _, v := range tests {
		result, err := f(v.Args)
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

var coalesceTests = []functionTest{
	{
		Name: "Coalesce",
		Args: []parser.Primary{
			parser.NewNull(),
			parser.NewString("str"),
		},
		Result: parser.NewString("str"),
	},
	{
		Name:  "Coalesce Argments Error",
		Args:  []parser.Primary{},
		Error: "function COALESCE takes at least 1 argument",
	},
	{
		Name: "Coalesce No Match",
		Args: []parser.Primary{
			parser.NewNull(),
			parser.NewNull(),
		},
		Result: parser.NewNull(),
	},
}

func TestCoalesce(t *testing.T) {
	testFunction(t, Coalesce, coalesceTests)
}

var ifTests = []functionTest{
	{
		Name: "If True",
		Args: []parser.Primary{
			parser.NewTernary(ternary.TRUE),
			parser.NewInteger(1),
			parser.NewInteger(2),
		},
		Result: parser.NewInteger(1),
	},
	{
		Name: "If False",
		Args: []parser.Primary{
			parser.NewTernary(ternary.FALSE),
			parser.NewInteger(1),
			parser.NewInteger(2),
		},
		Result: parser.NewInteger(2),
	},
	{
		Name: "If Argumants Error",
		Args: []parser.Primary{
			parser.NewTernary(ternary.FALSE),
		},
		Error: "function IF takes 3 arguments",
	},
}

func TestIf(t *testing.T) {
	testFunction(t, If, ifTests)
}

var ifnullTests = []functionTest{
	{
		Name: "Ifnull True",
		Args: []parser.Primary{
			parser.NewNull(),
			parser.NewInteger(2),
		},
		Result: parser.NewInteger(2),
	},
	{
		Name: "Ifnull False",
		Args: []parser.Primary{
			parser.NewInteger(1),
			parser.NewInteger(2),
		},
		Result: parser.NewInteger(1),
	},
	{
		Name: "Ifnull Arguments Error",
		Args: []parser.Primary{
			parser.NewInteger(1),
		},
		Error: "function IFNULL takes 2 arguments",
	},
}

func TestIfnull(t *testing.T) {
	testFunction(t, Ifnull, ifnullTests)
}

var nullifTests = []functionTest{
	{
		Name: "Nullif True",
		Args: []parser.Primary{
			parser.NewInteger(2),
			parser.NewInteger(2),
		},
		Result: parser.NewNull(),
	},
	{
		Name: "Nullif False",
		Args: []parser.Primary{
			parser.NewInteger(1),
			parser.NewInteger(2),
		},
		Result: parser.NewInteger(1),
	},
	{
		Name: "Nullif Arguments Error",
		Args: []parser.Primary{
			parser.NewInteger(1),
		},
		Error: "function NULLIF takes 2 arguments",
	},
}

func TestNullif(t *testing.T) {
	testFunction(t, Nullif, nullifTests)
}

var ceilTests = []functionTest{
	{
		Name: "Ceil",
		Args: []parser.Primary{
			parser.NewFloat(2.345),
		},
		Result: parser.NewInteger(3),
	},
	{
		Name: "Ceil Null",
		Args: []parser.Primary{
			parser.NewNull(),
		},
		Result: parser.NewNull(),
	},
	{
		Name: "Ceil Place is Null",
		Args: []parser.Primary{
			parser.NewFloat(2.345),
			parser.NewNull(),
		},
		Result: parser.NewNull(),
	},
	{
		Name:  "Ceil Arguments Error",
		Args:  []parser.Primary{},
		Error: "function CEIL takes 1 or 2 arguments",
	},
}

func TestCeil(t *testing.T) {
	testFunction(t, Ceil, ceilTests)
}

var floorTests = []functionTest{
	{
		Name: "Floor",
		Args: []parser.Primary{
			parser.NewFloat(2.345),
			parser.NewInteger(1),
		},
		Result: parser.NewFloat(2.3),
	},
	{
		Name: "Floor Null",
		Args: []parser.Primary{
			parser.NewNull(),
		},
		Result: parser.NewNull(),
	},
	{
		Name:  "Floor Arguments Error",
		Args:  []parser.Primary{},
		Error: "function FLOOR takes 1 or 2 arguments",
	},
}

func TestFloor(t *testing.T) {
	testFunction(t, Floor, floorTests)
}

var roundTests = []functionTest{
	{
		Name: "Round",
		Args: []parser.Primary{
			parser.NewFloat(2.4),
		},
		Result: parser.NewInteger(2),
	},
	{
		Name: "Round Negative Number",
		Args: []parser.Primary{
			parser.NewFloat(-2.4),
		},
		Result: parser.NewInteger(-2),
	},
	{
		Name: "Round Null",
		Args: []parser.Primary{
			parser.NewNull(),
		},
		Result: parser.NewNull(),
	},
	{
		Name:  "Round Arguments Error",
		Args:  []parser.Primary{},
		Error: "function ROUND takes 1 or 2 arguments",
	},
}

func TestRound(t *testing.T) {
	testFunction(t, Round, roundTests)
}

var absTests = []functionTest{
	{
		Name: "Abs",
		Args: []parser.Primary{
			parser.NewInteger(-2),
		},
		Result: parser.NewInteger(2),
	},
	{
		Name: "Abs Null",
		Args: []parser.Primary{
			parser.NewNull(),
		},
		Result: parser.NewNull(),
	},
	{
		Name:  "Abs Arguments Error",
		Args:  []parser.Primary{},
		Error: "function ABS takes 1 argument",
	},
}

func TestAbs(t *testing.T) {
	testFunction(t, Abs, absTests)
}

var acosTests = []functionTest{
	{
		Name: "Acos",
		Args: []parser.Primary{
			parser.NewInteger(0),
		},
		Result: parser.NewFloat(1.5707963267948966),
	},
}

func TestAcos(t *testing.T) {
	testFunction(t, Acos, acosTests)
}

var asinTests = []functionTest{
	{
		Name: "Asin",
		Args: []parser.Primary{
			parser.NewFloat(0.1),
		},
		Result: parser.NewFloat(0.1001674211615598),
	},
}

func TestAsin(t *testing.T) {
	testFunction(t, Asin, asinTests)
}

var atanTests = []functionTest{
	{
		Name: "ATAN",
		Args: []parser.Primary{
			parser.NewInteger(2),
		},
		Result: parser.NewFloat(1.1071487177940904),
	},
}

func TestAtan(t *testing.T) {
	testFunction(t, Atan, atanTests)
}

var atan2Tests = []functionTest{
	{
		Name: "Atan2",
		Args: []parser.Primary{
			parser.NewInteger(2),
			parser.NewInteger(2),
		},
		Result: parser.NewFloat(0.7853981633974483),
	},
}

func TestAtan2(t *testing.T) {
	testFunction(t, Atan2, atan2Tests)
}

var cosTests = []functionTest{
	{
		Name: "Cos",
		Args: []parser.Primary{
			parser.NewInteger(2),
		},
		Result: parser.NewFloat(-0.4161468365471424),
	},
}

func TestCos(t *testing.T) {
	testFunction(t, Cos, cosTests)
}

var sinTests = []functionTest{
	{
		Name: "Sin",
		Args: []parser.Primary{
			parser.NewInteger(1),
		},
		Result: parser.NewFloat(0.8414709848078965),
	},
}

func TestSin(t *testing.T) {
	testFunction(t, Sin, sinTests)
}

var tanTests = []functionTest{
	{
		Name: "Tan",
		Args: []parser.Primary{
			parser.NewInteger(2),
		},
		Result: parser.NewFloat(-2.185039863261519),
	},
}

func TestTan(t *testing.T) {
	testFunction(t, Tan, tanTests)
}

var expTests = []functionTest{
	{
		Name: "Exp",
		Args: []parser.Primary{
			parser.NewInteger(2),
		},
		Result: parser.NewFloat(7.38905609893065),
	},
}

func TestExp(t *testing.T) {
	testFunction(t, Exp, expTests)
}

var exp2Tests = []functionTest{
	{
		Name: "Exp2",
		Args: []parser.Primary{
			parser.NewInteger(2),
		},
		Result: parser.NewInteger(4),
	},
}

func TestExp2(t *testing.T) {
	testFunction(t, Exp2, exp2Tests)
}

var expm1Tests = []functionTest{
	{
		Name: "Expm1",
		Args: []parser.Primary{
			parser.NewInteger(2),
		},
		Result: parser.NewFloat(6.38905609893065),
	},
}

func TestExpm1(t *testing.T) {
	testFunction(t, Expm1, expm1Tests)
}

var logTests = []functionTest{
	{
		Name: "Log",
		Args: []parser.Primary{
			parser.NewFloat(2),
		},
		Result: parser.NewFloat(0.6931471805599453),
	},
}

func TestLog(t *testing.T) {
	testFunction(t, Log, logTests)
}

var log10Tests = []functionTest{
	{
		Name: "Log10",
		Args: []parser.Primary{
			parser.NewFloat(100),
		},
		Result: parser.NewInteger(2),
	},
}

func TestLog10(t *testing.T) {
	testFunction(t, Log10, log10Tests)
}

var log2Tests = []functionTest{
	{
		Name: "Log2",
		Args: []parser.Primary{
			parser.NewFloat(16),
		},
		Result: parser.NewInteger(4),
	},
}

func TestLog2(t *testing.T) {
	testFunction(t, Log2, log2Tests)
}

var log1pTests = []functionTest{
	{
		Name: "Log1p",
		Args: []parser.Primary{
			parser.NewFloat(1),
		},
		Result: parser.NewFloat(0.6931471805599453),
	},
}

func TestLog1p(t *testing.T) {
	testFunction(t, Log1p, log1pTests)
}

var sqrtTests = []functionTest{
	{
		Name: "Sqrt",
		Args: []parser.Primary{
			parser.NewFloat(4),
		},
		Result: parser.NewInteger(2),
	},
	{
		Name: "Sqrt Cannot Calculate",
		Args: []parser.Primary{
			parser.NewFloat(-4),
		},
		Result: parser.NewNull(),
	},
}

func TestSqrt(t *testing.T) {
	testFunction(t, Sqrt, sqrtTests)
}

var powTests = []functionTest{
	{
		Name: "Pow",
		Args: []parser.Primary{
			parser.NewFloat(2),
			parser.NewFloat(2),
		},
		Result: parser.NewInteger(4),
	},
	{
		Name: "Pow First Argument is Null",
		Args: []parser.Primary{
			parser.NewNull(),
			parser.NewFloat(2),
		},
		Result: parser.NewNull(),
	},
	{
		Name: "Pow Second Argument is Null",
		Args: []parser.Primary{
			parser.NewFloat(2),
			parser.NewNull(),
		},
		Result: parser.NewNull(),
	},
	{
		Name: "Pow Cannot Calculate",
		Args: []parser.Primary{
			parser.NewFloat(-2),
			parser.NewFloat(2.4),
		},
		Result: parser.NewNull(),
	},
	{
		Name:  "Pow Arguments Error",
		Args:  []parser.Primary{},
		Error: "function POW takes 2 arguments",
	},
}

func TestPow(t *testing.T) {
	testFunction(t, Pow, powTests)
}

var binToDecTests = []functionTest{
	{
		Name: "BinToDec",
		Args: []parser.Primary{
			parser.NewString("1111011"),
		},
		Result: parser.NewInteger(123),
	},
	{
		Name: "BinToDec Null",
		Args: []parser.Primary{
			parser.NewNull(),
		},
		Result: parser.NewNull(),
	},
	{
		Name: "BinToDec Parse Error",
		Args: []parser.Primary{
			parser.NewString("string"),
		},
		Result: parser.NewNull(),
	},
	{
		Name:  "BinToDec Arguments Error",
		Args:  []parser.Primary{},
		Error: "function BIN_TO_DEC takes 1 argument",
	},
}

func TestBinToDec(t *testing.T) {
	testFunction(t, BinToDec, binToDecTests)
}

var octToDecTests = []functionTest{
	{
		Name: "OctToDec",
		Args: []parser.Primary{
			parser.NewString("0173"),
		},
		Result: parser.NewInteger(123),
	},
}

func TestOctToDec(t *testing.T) {
	testFunction(t, OctToDec, octToDecTests)
}

var hexToDecTests = []functionTest{
	{
		Name: "HexToDec",
		Args: []parser.Primary{
			parser.NewString("0x7b"),
		},
		Result: parser.NewInteger(123),
	},
}

func TestHexToDec(t *testing.T) {
	testFunction(t, HexToDec, hexToDecTests)
}

var binTests = []functionTest{
	{
		Name: "Bin",
		Args: []parser.Primary{
			parser.NewInteger(123),
		},
		Result: parser.NewString("1111011"),
	},
	{
		Name: "Bin Null",
		Args: []parser.Primary{
			parser.NewNull(),
		},
		Result: parser.NewNull(),
	},
	{
		Name:  "Bin Arguments Error",
		Args:  []parser.Primary{},
		Error: "function BIN takes 1 argument",
	},
}

func TestBin(t *testing.T) {
	testFunction(t, Bin, binTests)
}

var octTests = []functionTest{
	{
		Name: "Oct",
		Args: []parser.Primary{
			parser.NewInteger(123),
		},
		Result: parser.NewString("173"),
	},
}

func TestOct(t *testing.T) {
	testFunction(t, Oct, octTests)
}

var hexTests = []functionTest{
	{
		Name: "Hex",
		Args: []parser.Primary{
			parser.NewInteger(123),
		},
		Result: parser.NewString("7b"),
	},
}

func TestHex(t *testing.T) {
	testFunction(t, Hex, hexTests)
}

var randTests = []struct {
	Name      string
	Args      []parser.Primary
	RangeLow  float64
	RangeHigh float64
	Error     string
}{
	{
		Name:      "Rand",
		RangeLow:  0.0,
		RangeHigh: 1.0,
	},
	{
		Name: "Rand Range Specified",
		Args: []parser.Primary{
			parser.NewInteger(7),
			parser.NewInteger(12),
		},
		RangeLow:  7.0,
		RangeHigh: 12.0,
	},
	{
		Name: "Range Arguments Error",
		Args: []parser.Primary{
			parser.NewInteger(1),
		},
		Error: "function RAND takes 0 or 2 arguments",
	},
	{
		Name: "Range First Arguments Error",
		Args: []parser.Primary{
			parser.NewString("a"),
			parser.NewInteger(2),
		},
		Error: "function RAND first argument must be parsable as integer",
	},
	{
		Name: "Range Second Arguments Error",
		Args: []parser.Primary{
			parser.NewInteger(1),
			parser.NewString("a"),
		},
		Error: "function RAND second argument must be parsable as integer",
	},
	{
		Name: "Range Arguments Value Error",
		Args: []parser.Primary{
			parser.NewInteger(1),
			parser.NewInteger(1),
		},
		Error: "function RAND second argument must be greater than first argument",
	},
}

func TestRand(t *testing.T) {
	for _, v := range randTests {
		result, err := Rand(v.Args)
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
			f = result.(parser.Float).Value()
		} else {
			f = float64(result.(parser.Integer).Value())
		}

		if f < v.RangeLow || v.RangeHigh < f {
			t.Errorf("%s: result = %f, want in range from %f to %f", v.Name, f, v.RangeLow, v.RangeHigh)
		}
	}
}

var trimTests = []functionTest{
	{
		Name: "Trim",
		Args: []parser.Primary{
			parser.NewString("aabbfoo, baraabb"),
			parser.NewString("ab"),
		},
		Result: parser.NewString("foo, bar"),
	},
	{
		Name: "Trim Spaces",
		Args: []parser.Primary{
			parser.NewString("  foo, bar \n"),
		},
		Result: parser.NewString("foo, bar"),
	},
	{
		Name: "Trim Null",
		Args: []parser.Primary{
			parser.NewNull(),
			parser.NewString("ab"),
		},
		Result: parser.NewNull(),
	},
	{
		Name: "Trim Cutset is Null",
		Args: []parser.Primary{
			parser.NewString("aabbfoo, baraabb"),
			parser.NewNull(),
		},
		Result: parser.NewNull(),
	},
	{
		Name:  "Trim Arguments Error",
		Args:  []parser.Primary{},
		Error: "function TRIM takes 1 or 2 arguments",
	},
}

func TestTrim(t *testing.T) {
	testFunction(t, Trim, trimTests)
}

var ltrimTests = []functionTest{
	{
		Name: "Trim",
		Args: []parser.Primary{
			parser.NewString("aabbfoo, baraabb"),
			parser.NewString("ab"),
		},
		Result: parser.NewString("foo, baraabb"),
	},
	{
		Name: "Ltrim Spaces",
		Args: []parser.Primary{
			parser.NewString("  foo, bar \n"),
		},
		Result: parser.NewString("foo, bar \n"),
	},
}

func TestLtrim(t *testing.T) {
	testFunction(t, Ltrim, ltrimTests)
}

var rtrimTests = []functionTest{
	{
		Name: "Trim",
		Args: []parser.Primary{
			parser.NewString("aabbfoo, baraabb"),
			parser.NewString("ab"),
		},
		Result: parser.NewString("aabbfoo, bar"),
	},
	{
		Name: "Rtrim Spaces",
		Args: []parser.Primary{
			parser.NewString("  foo, bar \n"),
		},
		Result: parser.NewString("  foo, bar"),
	},
}

func TestRtrim(t *testing.T) {
	testFunction(t, Rtrim, rtrimTests)
}

var upperTests = []functionTest{
	{
		Name: "Upper",
		Args: []parser.Primary{
			parser.NewString("Foo"),
		},
		Result: parser.NewString("FOO"),
	},
	{
		Name: "Upper Null",
		Args: []parser.Primary{
			parser.NewNull(),
		},
		Result: parser.NewNull(),
	},
	{
		Name:  "Upper Arguments Error",
		Args:  []parser.Primary{},
		Error: "function UPPER takes 1 argument",
	},
}

func TestUpper(t *testing.T) {
	testFunction(t, Upper, upperTests)
}

var lowerTests = []functionTest{
	{
		Name: "Lower",
		Args: []parser.Primary{
			parser.NewString("Foo"),
		},
		Result: parser.NewString("foo"),
	},
}

func TestLower(t *testing.T) {
	testFunction(t, Lower, lowerTests)
}

var base64EncodeTests = []functionTest{
	{
		Name: "Base64Encode",
		Args: []parser.Primary{
			parser.NewString("Foo"),
		},
		Result: parser.NewString("Rm9v"),
	},
}

func TestBase64Encode(t *testing.T) {
	testFunction(t, Base64Encode, base64EncodeTests)
}

var base64DecodeTests = []functionTest{
	{
		Name: "Base64Decode",
		Args: []parser.Primary{
			parser.NewString("Rm9v"),
		},
		Result: parser.NewString("Foo"),
	},
}

func TestBase64Decode(t *testing.T) {
	testFunction(t, Base64Decode, base64DecodeTests)
}

var hexEncodeTests = []functionTest{
	{
		Name: "HexEncode",
		Args: []parser.Primary{
			parser.NewString("Foo"),
		},
		Result: parser.NewString("466f6f"),
	},
}

func TestHexEncode(t *testing.T) {
	testFunction(t, HexEncode, hexEncodeTests)
}

var hexDecodeTests = []functionTest{
	{
		Name: "HexDecode",
		Args: []parser.Primary{
			parser.NewString("466f6f"),
		},
		Result: parser.NewString("Foo"),
	},
}

func TestHexDecode(t *testing.T) {
	testFunction(t, HexDecode, hexDecodeTests)
}

var lenTests = []functionTest{
	{
		Name: "Len",
		Args: []parser.Primary{
			parser.NewString("日本語"),
		},
		Result: parser.NewInteger(3),
	},
	{
		Name: "Len Null",
		Args: []parser.Primary{
			parser.NewNull(),
		},
		Result: parser.NewNull(),
	},
	{
		Name:  "Len Arguments Error",
		Args:  []parser.Primary{},
		Error: "function LEN takes 1 argument",
	},
}

func TestLen(t *testing.T) {
	testFunction(t, Len, lenTests)
}

var byteLenTests = []functionTest{
	{
		Name: "ByteLen",
		Args: []parser.Primary{
			parser.NewString("日本語"),
		},
		Result: parser.NewInteger(9),
	},
	{
		Name: "ByteLen Null",
		Args: []parser.Primary{
			parser.NewNull(),
		},
		Result: parser.NewNull(),
	},
	{
		Name:  "ByteLen Arguments Error",
		Args:  []parser.Primary{},
		Error: "function BYTE_LEN takes 1 argument",
	},
}

func TestByteLen(t *testing.T) {
	testFunction(t, ByteLen, byteLenTests)
}

var lpadTests = []functionTest{
	{
		Name: "Lpad",
		Args: []parser.Primary{
			parser.NewString("aaaaa"),
			parser.NewInteger(10),
			parser.NewString("01"),
		},
		Result: parser.NewString("01010aaaaa"),
	},
	{
		Name: "Lpad No Padding",
		Args: []parser.Primary{
			parser.NewString("aaaaa"),
			parser.NewInteger(5),
			parser.NewString("01"),
		},
		Result: parser.NewString("aaaaa"),
	},
	{
		Name: "Lpad String is Null",
		Args: []parser.Primary{
			parser.NewNull(),
			parser.NewInteger(10),
			parser.NewString("01"),
		},
		Result: parser.NewNull(),
	},
	{
		Name: "Lpad Length is Null",
		Args: []parser.Primary{
			parser.NewString("aaaaa"),
			parser.NewNull(),
			parser.NewString("01"),
		},
		Result: parser.NewNull(),
	},
	{
		Name: "Lpad Pad String is Null",
		Args: []parser.Primary{
			parser.NewString("aaaaa"),
			parser.NewInteger(10),
			parser.NewNull(),
		},
		Result: parser.NewNull(),
	},
	{
		Name:  "Lpad Arguments Error",
		Args:  []parser.Primary{},
		Error: "function LPAD takes 3 arguments",
	},
}

func TestLpad(t *testing.T) {
	testFunction(t, Lpad, lpadTests)
}

var rpadTests = []functionTest{
	{
		Name: "Rpad",
		Args: []parser.Primary{
			parser.NewString("aaaaa"),
			parser.NewInteger(10),
			parser.NewString("01"),
		},
		Result: parser.NewString("aaaaa01010"),
	},
}

func TestRpad(t *testing.T) {
	testFunction(t, Rpad, rpadTests)
}

var substrTests = []functionTest{
	{
		Name: "Substr",
		Args: []parser.Primary{
			parser.NewString("abcdefghijklmn"),
			parser.NewInteger(-5),
			parser.NewInteger(8),
		},
		Result: parser.NewString("jklmn"),
	},
	{
		Name: "Substr String is Null",
		Args: []parser.Primary{
			parser.NewNull(),
			parser.NewInteger(-5),
			parser.NewInteger(8),
		},
		Result: parser.NewNull(),
	},
	{
		Name: "Substr StartIndex is Null",
		Args: []parser.Primary{
			parser.NewString("abcdefghijklmn"),
			parser.NewNull(),
		},
		Result: parser.NewNull(),
	},
	{
		Name: "Substr Length is Null",
		Args: []parser.Primary{
			parser.NewString("abcdefghijklmn"),
			parser.NewInteger(-5),
			parser.NewNull(),
		},
		Result: parser.NewNull(),
	},
	{
		Name: "Substr Length is Negative",
		Args: []parser.Primary{
			parser.NewString("abcdefghijklmn"),
			parser.NewInteger(-5),
			parser.NewInteger(-1),
		},
		Result: parser.NewNull(),
	},
	{
		Name: "Substr StartIndex is Out Of Index",
		Args: []parser.Primary{
			parser.NewString("abcdefghijklmn"),
			parser.NewInteger(100),
			parser.NewInteger(8),
		},
		Result: parser.NewNull(),
	},
	{
		Name:  "Substr Arguments Error",
		Args:  []parser.Primary{},
		Error: "function SUBSTR takes 2 or 3 arguments",
	},
}

func TestSubstr(t *testing.T) {
	testFunction(t, Substr, substrTests)
}

var replaceTests = []functionTest{
	{
		Name: "Replace",
		Args: []parser.Primary{
			parser.NewString("abcdefg abcdefg"),
			parser.NewString("cd"),
			parser.NewString("CD"),
		},
		Result: parser.NewString("abCDefg abCDefg"),
	},
	{
		Name: "Replace String is Null",
		Args: []parser.Primary{
			parser.NewNull(),
			parser.NewString("cd"),
			parser.NewString("CD"),
		},
		Result: parser.NewNull(),
	},
	{
		Name: "Replace Old String is Null",
		Args: []parser.Primary{
			parser.NewString("abcdefg abcdefg"),
			parser.NewNull(),
			parser.NewString("CD"),
		},
		Result: parser.NewNull(),
	},
	{
		Name: "Replace New String is Null",
		Args: []parser.Primary{
			parser.NewString("abcdefg abcdefg"),
			parser.NewString("cd"),
			parser.NewNull(),
		},
		Result: parser.NewNull(),
	},
	{
		Name:  "Replace Arguments Error",
		Args:  []parser.Primary{},
		Error: "function REPLACE takes 3 arguments",
	},
}

func TestReplace(t *testing.T) {
	testFunction(t, Replace, replaceTests)
}

var md5Tests = []functionTest{
	{
		Name: "Md5",
		Args: []parser.Primary{
			parser.NewString("foo"),
		},
		Result: parser.NewString("acbd18db4cc2f85cedef654fccc4a4d8"),
	},
	{
		Name: "Md5 Null",
		Args: []parser.Primary{
			parser.NewNull(),
		},
		Result: parser.NewNull(),
	},
	{
		Name:  "Md5 Arguments Error",
		Args:  []parser.Primary{},
		Error: "function MD5 takes 1 argument",
	},
}

func TestMd5(t *testing.T) {
	testFunction(t, Md5, md5Tests)
}

var sha1Tests = []functionTest{
	{
		Name: "Sha1",
		Args: []parser.Primary{
			parser.NewString("foo"),
		},
		Result: parser.NewString("0beec7b5ea3f0fdbc95d0dd47f3c5bc275da8a33"),
	},
}

func TestSha1(t *testing.T) {
	testFunction(t, Sha1, sha1Tests)
}

var sha256Tests = []functionTest{
	{
		Name: "Sha256",
		Args: []parser.Primary{
			parser.NewString("foo"),
		},
		Result: parser.NewString("2c26b46b68ffc68ff99b453c1d30413413422d706483bfa0f98a5e886266e7ae"),
	},
}

func TestSha256(t *testing.T) {
	testFunction(t, Sha256, sha256Tests)
}

var sha512Tests = []functionTest{
	{
		Name: "Sha512",
		Args: []parser.Primary{
			parser.NewString("foo"),
		},
		Result: parser.NewString("f7fbba6e0636f890e56fbbf3283e524c6fa3204ae298382d624741d0dc6638326e282c41be5e4254d8820772c5518a2c5a8c0c7f7eda19594a7eb539453e1ed7"),
	},
}

func TestSha512(t *testing.T) {
	testFunction(t, Sha512, sha512Tests)
}

var md5HmacTests = []functionTest{
	{
		Name: "Md5Hmac",
		Args: []parser.Primary{
			parser.NewString("foo"),
			parser.NewString("bar"),
		},
		Result: parser.NewString("31b6db9e5eb4addb42f1a6ca07367adc"),
	},
	{
		Name: "Md5Hmac String is Null",
		Args: []parser.Primary{
			parser.NewNull(),
			parser.NewString("bar"),
		},
		Result: parser.NewNull(),
	},
	{
		Name: "Md5Hmac Key is Null",
		Args: []parser.Primary{
			parser.NewString("foo"),
			parser.NewNull(),
		},
		Result: parser.NewNull(),
	},
	{
		Name:  "Md5Hmac Arguments Error",
		Args:  []parser.Primary{},
		Error: "function MD5_HMAC takes 2 arguments",
	},
}

func TestMd5Hmac(t *testing.T) {
	testFunction(t, Md5Hmac, md5HmacTests)
}

var sha1HmacTests = []functionTest{
	{
		Name: "Sha1Hmac",
		Args: []parser.Primary{
			parser.NewString("foo"),
			parser.NewString("bar"),
		},
		Result: parser.NewString("85d155c55ed286a300bd1cf124de08d87e914f3a"),
	},
}

func TestSha1Hmac(t *testing.T) {
	testFunction(t, Sha1Hmac, sha1HmacTests)
}

var sha256HmacTests = []functionTest{
	{
		Name: "Sha256Hmac",
		Args: []parser.Primary{
			parser.NewString("foo"),
			parser.NewString("bar"),
		},
		Result: parser.NewString("147933218aaabc0b8b10a2b3a5c34684c8d94341bcf10a4736dc7270f7741851"),
	},
}

func TestSha256Hmac(t *testing.T) {
	testFunction(t, Sha256Hmac, sha256HmacTests)
}

var sha512HmacTests = []functionTest{
	{
		Name: "Sha512Hmac",
		Args: []parser.Primary{
			parser.NewString("foo"),
			parser.NewString("bar"),
		},
		Result: parser.NewString("24257d7210582a65c731ec55159c8184cc24c02489453e58587f71f44c23a2d61b4b72154a89d17b2d49448a8452ea066f4fc56a2bcead45c088572ffccdb3d8"),
	},
}

func TestSha512Hmac(t *testing.T) {
	testFunction(t, Sha512Hmac, sha512HmacTests)
}

var nowTests = []functionTest{
	{
		Name:   "Now",
		Result: parser.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 0, GetTestLocation())),
	},
	{
		Name: "Now Arguments Error",
		Args: []parser.Primary{
			parser.NewInteger(1),
		},
		Error: "function NOW takes no arguments",
	},
}

func TestNow(t *testing.T) {
	testFunction(t, Now, nowTests)
}

var datetimeFormatTests = []functionTest{
	{
		Name: "DatetimeFormat",
		Args: []parser.Primary{
			parser.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 0, GetTestLocation())),
			parser.NewString("%Y-%m-%d"),
		},
		Result: parser.NewString("2012-02-03"),
	},
	{
		Name: "DatetimeFormat Datetime is Null",
		Args: []parser.Primary{
			parser.NewNull(),
			parser.NewString("%Y-%m-%d"),
		},
		Result: parser.NewNull(),
	},
	{
		Name: "DatetimeFormat Format is Null",
		Args: []parser.Primary{
			parser.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 0, GetTestLocation())),
			parser.NewNull(),
		},
		Result: parser.NewNull(),
	},
	{
		Name: "DatetimeFormat Arguments Error",
		Args: []parser.Primary{
			parser.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 0, GetTestLocation())),
		},
		Error: "function DATETIME_FORMAT takes 2 arguments",
	},
}

func TestDatetimeFormat(t *testing.T) {
	testFunction(t, DatetimeFormat, datetimeFormatTests)
}

var yearTests = []functionTest{
	{
		Name: "Year",
		Args: []parser.Primary{
			parser.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 0, GetTestLocation())),
		},
		Result: parser.NewInteger(2012),
	},
	{
		Name: "Year Datetime is Null",
		Args: []parser.Primary{
			parser.NewNull(),
		},
		Result: parser.NewNull(),
	},
	{
		Name:  "Year Arguments Error",
		Args:  []parser.Primary{},
		Error: "function YEAR takes 1 argument",
	},
}

func TestYear(t *testing.T) {
	testFunction(t, Year, yearTests)
}

var monthTests = []functionTest{
	{
		Name: "Month",
		Args: []parser.Primary{
			parser.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 0, GetTestLocation())),
		},
		Result: parser.NewInteger(2),
	},
}

func TestMonth(t *testing.T) {
	testFunction(t, Month, monthTests)
}

var dayTests = []functionTest{
	{
		Name: "Day",
		Args: []parser.Primary{
			parser.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 0, GetTestLocation())),
		},
		Result: parser.NewInteger(3),
	},
}

func TestDay(t *testing.T) {
	testFunction(t, Day, dayTests)
}

var hourTests = []functionTest{
	{
		Name: "Hour",
		Args: []parser.Primary{
			parser.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 0, GetTestLocation())),
		},
		Result: parser.NewInteger(9),
	},
}

func TestHour(t *testing.T) {
	testFunction(t, Hour, hourTests)
}

var minuteTests = []functionTest{
	{
		Name: "Minute",
		Args: []parser.Primary{
			parser.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 0, GetTestLocation())),
		},
		Result: parser.NewInteger(18),
	},
}

func TestMinute(t *testing.T) {
	testFunction(t, Minute, minuteTests)
}

var secondTests = []functionTest{
	{
		Name: "Second",
		Args: []parser.Primary{
			parser.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 0, GetTestLocation())),
		},
		Result: parser.NewInteger(15),
	},
}

func TestSecond(t *testing.T) {
	testFunction(t, Second, secondTests)
}

var millisecondTests = []functionTest{
	{
		Name: "Millisecond",
		Args: []parser.Primary{
			parser.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456789, GetTestLocation())),
		},
		Result: parser.NewInteger(123),
	},
}

func TestMillisecond(t *testing.T) {
	testFunction(t, Millisecond, millisecondTests)
}

var microsecondTests = []functionTest{
	{
		Name: "Microsecond",
		Args: []parser.Primary{
			parser.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456789, GetTestLocation())),
		},
		Result: parser.NewInteger(123457),
	},
}

func TestMicrosecond(t *testing.T) {
	testFunction(t, Microsecond, microsecondTests)
}

var nanosecondTests = []functionTest{
	{
		Name: "Nanosecond",
		Args: []parser.Primary{
			parser.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456789, GetTestLocation())),
		},
		Result: parser.NewInteger(123456789),
	},
}

func TestNanosecond(t *testing.T) {
	testFunction(t, Nanosecond, nanosecondTests)
}

var weekdayTests = []functionTest{
	{
		Name: "Weekday",
		Args: []parser.Primary{
			parser.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 0, GetTestLocation())),
		},
		Result: parser.NewInteger(5),
	},
}

func TestWeekday(t *testing.T) {
	testFunction(t, Weekday, weekdayTests)
}

var unixTimeTests = []functionTest{
	{
		Name: "UnixTime",
		Args: []parser.Primary{
			parser.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 0, GetTestLocation())),
		},
		Result: parser.NewInteger(1328289495),
	},
}

func TestUnixTime(t *testing.T) {
	testFunction(t, UnixTime, unixTimeTests)
}

var unixNanoTimeTests = []functionTest{
	{
		Name: "UnixNanoTime",
		Args: []parser.Primary{
			parser.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456789, GetTestLocation())),
		},
		Result: parser.NewInteger(1328289495123456789),
	},
}

func TestUnixNanoTime(t *testing.T) {
	testFunction(t, UnixNanoTime, unixNanoTimeTests)
}

var dayOfYearTests = []functionTest{
	{
		Name: "DayOfYear",
		Args: []parser.Primary{
			parser.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 0, GetTestLocation())),
		},
		Result: parser.NewInteger(34),
	},
}

func TestDayOfYear(t *testing.T) {
	testFunction(t, DayOfYear, dayOfYearTests)
}

var weekOfYearTests = []functionTest{
	{
		Name: "WeekOfYear",
		Args: []parser.Primary{
			parser.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 0, GetTestLocation())),
		},
		Result: parser.NewInteger(5),
	},
}

func TestWeekOfYear(t *testing.T) {
	testFunction(t, WeekOfYear, weekOfYearTests)
}

var addYearTests = []functionTest{
	{
		Name: "AddYear",
		Args: []parser.Primary{
			parser.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456789, GetTestLocation())),
			parser.NewInteger(2),
		},
		Result: parser.NewDatetime(time.Date(2014, 2, 3, 9, 18, 15, 123456789, GetTestLocation())),
	},
	{
		Name: "AddYear Datetime is Null",
		Args: []parser.Primary{
			parser.NewNull(),
			parser.NewInteger(2),
		},
		Result: parser.NewNull(),
	},
	{
		Name: "AddYear Duration is Null",
		Args: []parser.Primary{
			parser.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456789, GetTestLocation())),
			parser.NewNull(),
		},
		Result: parser.NewNull(),
	},
	{
		Name:  "AddYear Arguments Error",
		Args:  []parser.Primary{},
		Error: "function ADD_YEAR takes 2 arguments",
	},
}

func TestAddYear(t *testing.T) {
	testFunction(t, AddYear, addYearTests)
}

var addMonthTests = []functionTest{
	{
		Name: "AddMonth",
		Args: []parser.Primary{
			parser.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456789, GetTestLocation())),
			parser.NewInteger(2),
		},
		Result: parser.NewDatetime(time.Date(2012, 4, 3, 9, 18, 15, 123456789, GetTestLocation())),
	},
}

func TestAddMonth(t *testing.T) {
	testFunction(t, AddMonth, addMonthTests)
}

var addDayTests = []functionTest{
	{
		Name: "AddDay",
		Args: []parser.Primary{
			parser.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456789, GetTestLocation())),
			parser.NewInteger(2),
		},
		Result: parser.NewDatetime(time.Date(2012, 2, 5, 9, 18, 15, 123456789, GetTestLocation())),
	},
}

func TestAddDay(t *testing.T) {
	testFunction(t, AddDay, addDayTests)
}

var addHourTests = []functionTest{
	{
		Name: "AddHour",
		Args: []parser.Primary{
			parser.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456789, GetTestLocation())),
			parser.NewInteger(2),
		},
		Result: parser.NewDatetime(time.Date(2012, 2, 3, 11, 18, 15, 123456789, GetTestLocation())),
	},
}

func TestAddHour(t *testing.T) {
	testFunction(t, AddHour, addHourTests)
}

var addMinuteTests = []functionTest{
	{
		Name: "AddMinute",
		Args: []parser.Primary{
			parser.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456789, GetTestLocation())),
			parser.NewInteger(2),
		},
		Result: parser.NewDatetime(time.Date(2012, 2, 3, 9, 20, 15, 123456789, GetTestLocation())),
	},
}

func TestAddMinute(t *testing.T) {
	testFunction(t, AddMinute, addMinuteTests)
}

var addSecondTests = []functionTest{
	{
		Name: "AddSecond",
		Args: []parser.Primary{
			parser.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456789, GetTestLocation())),
			parser.NewInteger(2),
		},
		Result: parser.NewDatetime(time.Date(2012, 2, 3, 9, 18, 17, 123456789, GetTestLocation())),
	},
}

func TestAddSecond(t *testing.T) {
	testFunction(t, AddSecond, addSecondTests)
}

var addMilliTests = []functionTest{
	{
		Name: "AddMilli",
		Args: []parser.Primary{
			parser.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456789, GetTestLocation())),
			parser.NewInteger(2),
		},
		Result: parser.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 125456789, GetTestLocation())),
	},
}

func TestAddMilli(t *testing.T) {
	testFunction(t, AddMilli, addMilliTests)
}

var addMicroTests = []functionTest{
	{
		Name: "AddMicro",
		Args: []parser.Primary{
			parser.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456789, GetTestLocation())),
			parser.NewInteger(2),
		},
		Result: parser.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123458789, GetTestLocation())),
	},
}

func TestAddMicro(t *testing.T) {
	testFunction(t, AddMicro, addMicroTests)
}

var addNanoTests = []functionTest{
	{
		Name: "AddNano",
		Args: []parser.Primary{
			parser.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456789, GetTestLocation())),
			parser.NewInteger(2),
		},
		Result: parser.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456791, GetTestLocation())),
	},
}

func TestAddNano(t *testing.T) {
	testFunction(t, AddNano, addNanoTests)
}

var dateDiffTests = []functionTest{
	{
		Name: "DateDiff",
		Args: []parser.Primary{
			parser.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456789, GetTestLocation())),
			parser.NewDatetime(time.Date(2012, 2, 5, 1, 18, 55, 123456789, GetTestLocation())),
		},
		Result: parser.NewInteger(-2),
	},
	{
		Name: "DateDiff Datetime1 is Null",
		Args: []parser.Primary{
			parser.NewNull(),
			parser.NewDatetime(time.Date(2012, 2, 5, 1, 18, 55, 123456789, GetTestLocation())),
		},
		Result: parser.NewNull(),
	},
	{
		Name: "DateDiff Datetime2 is Null",
		Args: []parser.Primary{
			parser.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456789, GetTestLocation())),
			parser.NewNull(),
		},
		Result: parser.NewNull(),
	},
	{
		Name:  "DateDiff Arguments Error",
		Args:  []parser.Primary{},
		Error: "function DATE_DIFF takes 2 arguments",
	},
}

func TestDateDiff(t *testing.T) {
	testFunction(t, DateDiff, dateDiffTests)
}

var timeDiffTests = []functionTest{
	{
		Name: "TimeDiff",
		Args: []parser.Primary{
			parser.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456789, GetTestLocation())),
			parser.NewDatetime(time.Date(2012, 2, 3, 1, 18, 55, 123000000, GetTestLocation())),
		},
		Result: parser.NewFloat(28760.000456789),
	},
	{
		Name: "TimeDiff Datetime1 is Null",
		Args: []parser.Primary{
			parser.NewNull(),
			parser.NewDatetime(time.Date(2012, 2, 5, 1, 18, 55, 123456789, GetTestLocation())),
		},
		Result: parser.NewNull(),
	},
	{
		Name: "TimeDiff Datetime2 is Null",
		Args: []parser.Primary{
			parser.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456789, GetTestLocation())),
			parser.NewNull(),
		},
		Result: parser.NewNull(),
	},
	{
		Name:  "TimeDiff Arguments Error",
		Args:  []parser.Primary{},
		Error: "function TIME_DIFF takes 2 arguments",
	},
}

func TestTimeDiff(t *testing.T) {
	testFunction(t, TimeDiff, timeDiffTests)
}

var autoIncrementTests = []functionTest{
	{
		Name: "AutoIncrement Arguments Error",
		Args: []parser.Primary{
			parser.NewInteger(10),
			parser.NewInteger(10),
		},
		Error: "function AUTO_INCREMENT takes at most 1 argument",
	},
	{
		Name: "AutoIncrement First Time",
		Args: []parser.Primary{
			parser.NewInteger(10),
		},
		Result: parser.NewInteger(10),
	},
	{
		Name: "AutoIncrement Second Time",
		Args: []parser.Primary{
			parser.NewInteger(10),
		},
		Result: parser.NewInteger(11),
	},
}

var stringTests = []functionTest{
	{
		Name: "String from Integer",
		Args: []parser.Primary{
			parser.NewInteger(2),
		},
		Result: parser.NewString("2"),
	},
	{
		Name: "String from Boolean",
		Args: []parser.Primary{
			parser.NewBoolean(true),
		},
		Result: parser.NewString("true"),
	},
	{
		Name: "String from Ternary",
		Args: []parser.Primary{
			parser.NewTernary(ternary.TRUE),
		},
		Result: parser.NewString("TRUE"),
	},
	{
		Name: "String from Datetime",
		Args: []parser.Primary{
			parser.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 0, GetTestLocation())),
		},
		Result: parser.NewString("2012-02-03 09:18:15"),
	},
	{
		Name:  "String Arguments Error",
		Args:  []parser.Primary{},
		Error: "function STRING takes 1 argument",
	},
}

func TestString(t *testing.T) {
	testFunction(t, String, stringTests)
}

var integerTests = []functionTest{
	{
		Name: "Integer from String",
		Args: []parser.Primary{
			parser.NewString("2"),
		},
		Result: parser.NewInteger(2),
	},
	{
		Name: "Integer from Float",
		Args: []parser.Primary{
			parser.NewFloat(1.7),
		},
		Result: parser.NewInteger(2),
	},
	{
		Name: "Integer from Datetime",
		Args: []parser.Primary{
			parser.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 0, GetTestLocation())),
		},
		Result: parser.NewInteger(1328289495),
	},
	{
		Name:  "Integer Arguments Error",
		Args:  []parser.Primary{},
		Error: "function INTEGER takes 1 argument",
	},
}

func TestInteger(t *testing.T) {
	testFunction(t, Integer, integerTests)
}

var floatTests = []functionTest{
	{
		Name: "Float from String",
		Args: []parser.Primary{
			parser.NewString("2"),
		},
		Result: parser.NewFloat(2),
	},
	{
		Name: "Float from Datetime",
		Args: []parser.Primary{
			parser.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 0, GetTestLocation())),
		},
		Result: parser.NewFloat(1328289495),
	},
	{
		Name:  "Float Arguments Error",
		Args:  []parser.Primary{},
		Error: "function FLOAT takes 1 argument",
	},
}

func TestFloat(t *testing.T) {
	testFunction(t, Float, floatTests)
}

var booleanTests = []functionTest{
	{
		Name: "Boolean from String",
		Args: []parser.Primary{
			parser.NewString("true"),
		},
		Result: parser.NewBoolean(true),
	},
	{
		Name:  "Boolean Arguments Error",
		Args:  []parser.Primary{},
		Error: "function BOOLEAN takes 1 argument",
	},
}

func TestBoolean(t *testing.T) {
	testFunction(t, Boolean, booleanTests)
}

var ternaryTest = []functionTest{
	{
		Name: "Ternary from String",
		Args: []parser.Primary{
			parser.NewString("true"),
		},
		Result: parser.NewTernary(ternary.TRUE),
	},
	{
		Name:  "Ternary Arguments Error",
		Args:  []parser.Primary{},
		Error: "function TERNARY takes 1 argument",
	},
}

func TestTernary(t *testing.T) {
	testFunction(t, Ternary, ternaryTest)
}

var datetimeTests = []functionTest{
	{
		Name: "Datetime from String",
		Args: []parser.Primary{
			parser.NewString("2012-02-03 09:18:15"),
		},
		Result: parser.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 0, GetTestLocation())),
	},
	{
		Name:  "Datetime Arguments Error",
		Args:  []parser.Primary{},
		Error: "function DATETIME takes 1 argument",
	},
}

func TestDatetime(t *testing.T) {
	testFunction(t, Datetime, datetimeTests)
}

var callTests = []functionTest{
	{
		Name: "Call",
		Args: []parser.Primary{
			parser.NewString("echo"),
			parser.NewString("foo"),
			parser.NewInteger(1),
			parser.NewFloat(1.234),
			parser.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 0, GetTestLocation())),
			parser.NewBoolean(true),
			parser.NewTernary(ternary.TRUE),
			parser.NewNull(),
		},
		Result: parser.NewString("foo 1 1.234 2012-02-03 09:18:15 true true \n"),
	},
	{
		Name:  "Call Argument Error",
		Args:  []parser.Primary{},
		Error: "function CALL takes at least 1 argument",
	},
	{
		Name: "Call Command Error",
		Args: []parser.Primary{
			parser.NewString("notexistcommand"),
		},
		Error: "exec: \"notexistcommand\": executable file not found in $PATH",
	},
}

func TestCall(t *testing.T) {
	testFunction(t, Call, callTests)
}
