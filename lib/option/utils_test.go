package option

import (
	"reflect"
	"testing"

	"github.com/mithrandie/go-text"
)

func TestEscapeString(t *testing.T) {
	str := "fo\\o\a\b\f\n\r\t\v\\\\'\"bar\\"
	expect := "fo\\\\o\\a\\b\\f\\n\\r\\t\\v\\\\\\\\\\'\"bar\\\\"
	unescaped := EscapeString(str)
	if unescaped != expect {
		t.Errorf("escaped string = %q, want %q", unescaped, expect)
	}
}

func TestUnescapeString(t *testing.T) {
	str := "fo\\o\\a\\b\\f\\n\\r\\t\\v\\\\\\\\'\\\"bar''\"\"\\"
	expect := "fo\\o\a\b\f\n\r\t\v\\\\'\"bar'\"\"\\"
	unescaped := UnescapeString(str, '\'')
	if unescaped != expect {
		t.Errorf("unescaped string = %q, want %q", unescaped, expect)
	}
}

func TestEscapeIdentifier(t *testing.T) {
	str := "fo\\o\a\b\f\n\r\t\v\\\\`bar\\"
	expect := "fo\\\\o\\a\\b\\f\\n\\r\\t\\v\\\\\\\\\\`bar\\\\"
	unescaped := EscapeIdentifier(str)
	if unescaped != expect {
		t.Errorf("escaped identifier = %q, want %q", unescaped, expect)
	}
}

func TestUnescapeIdentifier(t *testing.T) {
	str := "fo\\o\\a\\b\\f\\n\\r\\t\\v\\\\\\\\`bar``\\"
	expect := "fo\\o\a\b\f\n\r\t\v\\\\`bar`\\"
	unescaped := UnescapeIdentifier(str, '`')
	if unescaped != expect {
		t.Errorf("unescaped identifier = %q, want %q", unescaped, expect)
	}
}

func TestQuoteString(t *testing.T) {
	s := "abc'def"
	expect := "'abc\\'def'"
	result := QuoteString(s)
	if result != expect {
		t.Errorf("quoted string = %q, want %q for %q", result, expect, s)
	}
}

func TestQuoteIdentifier(t *testing.T) {
	s := "abc`def"
	expect := "`abc\\`def`"
	result := QuoteIdentifier(s)
	if result != expect {
		t.Errorf("quoted identifier = %q, want %q for %q", result, expect, s)
	}
}

func TestVariableSymbol(t *testing.T) {
	s := "var"
	expect := "@var"
	result := VariableSymbol(s)
	if result != expect {
		t.Errorf("variable symbol = %q, want %q for %q", result, expect, s)
	}
}

func TestFlagSymbol(t *testing.T) {
	s := "flag"
	expect := "@@flag"
	result := FlagSymbol(s)
	if result != expect {
		t.Errorf("flag symbol = %q, want %q for %q", result, expect, s)
	}
}

func TestEnvironmentVariableSymbol(t *testing.T) {
	s := "env"
	expect := "@%env"
	result := EnvironmentVariableSymbol(s)
	if result != expect {
		t.Errorf("environment variable symbol = %q, want %q for %q", result, expect, s)
	}

	s = "1env"
	expect = "@%`1env`"
	result = EnvironmentVariableSymbol(s)
	if result != expect {
		t.Errorf("environment variable symbol = %q, want %q for %q", result, expect, s)
	}
}

func TestEnclosedEnvironmentVariableSymbol(t *testing.T) {
	s := "env"
	expect := "@%`env`"
	result := EnclosedEnvironmentVariableSymbol(s)
	if result != expect {
		t.Errorf("environment variable symbol = %q, want %q for %q", result, expect, s)
	}
}

func TestRuntimeInformationSymbol(t *testing.T) {
	s := "info"
	expect := "@#info"
	result := RuntimeInformationSymbol(s)
	if result != expect {
		t.Errorf("runtime information symbol = %q, want %q for %q", result, expect, s)
	}
}

func TestMustBeEnclosed(t *testing.T) {
	var expect bool

	s := ""
	expect = false
	result := MustBeEnclosed(s)
	if result != expect {
		t.Errorf("must be enclosed = %t, want %t for %q", result, expect, s)
	}

	s = "1abc123"
	expect = true
	result = MustBeEnclosed(s)
	if result != expect {
		t.Errorf("must be enclosed = %t, want %t for %q", result, expect, s)
	}

	s = "abc123"
	expect = false
	result = MustBeEnclosed(s)
	if result != expect {
		t.Errorf("must be enclosed = %t, want %t for %q", result, expect, s)
	}

	s = "abc12#3"
	expect = true
	result = MustBeEnclosed(s)
	if result != expect {
		t.Errorf("must be enclosed = %t, want %t for %q", result, expect, s)
	}
}

func TestFormatInt(t *testing.T) {
	i := 1234567
	sep := ","
	expect := "1,234,567"
	result := FormatInt(i, sep)
	if result != expect {
		t.Errorf("format int = %q, want %q for %d", result, expect, i)
	}
}

var formatNumberTests = []struct {
	Float              float64
	Precision          int
	DecimalPoint       string
	ThousandsSeparator string
	DecimalSeparator   string
	Expect             string
}{
	{
		Float:              0,
		Precision:          0,
		DecimalPoint:       ".",
		ThousandsSeparator: ",",
		DecimalSeparator:   " ",
		Expect:             "0",
	},
	{
		Float:              123456.789123,
		Precision:          4,
		DecimalPoint:       ".",
		ThousandsSeparator: ",",
		DecimalSeparator:   " ",
		Expect:             "123,456.789 1",
	},
	{
		Float:              123456.7891,
		Precision:          -1,
		DecimalPoint:       ".",
		ThousandsSeparator: ",",
		DecimalSeparator:   "",
		Expect:             "123,456.7891",
	},
}

func TestFormatNumber(t *testing.T) {
	for _, v := range formatNumberTests {
		result := FormatNumber(v.Float, v.Precision, v.DecimalPoint, v.ThousandsSeparator, v.DecimalSeparator)
		if result != v.Expect {
			t.Errorf("result = %s, want %s for %f, %d, %q, %q, %q", result, v.Expect, v.Float, v.Precision, v.DecimalPoint, v.ThousandsSeparator, v.DecimalSeparator)
		}
	}
}

func TestParseEncoding(t *testing.T) {
	e, err := ParseEncoding("utf8")
	if err != nil {
		t.Errorf("unexpected error: %q", err.Error())
	}
	if e != text.UTF8 {
		t.Errorf("encoding = %s, expect to set %s for %s", e, text.UTF8, "utf8")
	}

	e, err = ParseEncoding("sjis")
	if err != nil {
		t.Errorf("unexpected error: %q", err.Error())
	}
	if e != text.SJIS {
		t.Errorf("encoding = %s, expect to set %s for %s", e, text.SJIS, "sjis")
	}

	expectErr := "encoding must be one of AUTO|UTF8|UTF8M|UTF16|UTF16BE|UTF16LE|UTF16BEM|UTF16LEM|SJIS"
	_, err = ParseEncoding("error")
	if err == nil {
		t.Errorf("no error, want error %q for %s", expectErr, "error")
	} else if err.Error() != expectErr {
		t.Errorf("error = %q, want error %q for %s", err.Error(), expectErr, "error")
	}
}

func TestParseDelimiter(t *testing.T) {
	var s string

	var expect rune

	s = "\t"
	expect = '\t'
	result, err := ParseDelimiter(s)
	if err != nil {
		t.Errorf("unexpected error: %q", err.Error())
	} else if expect != result {
		t.Errorf("result = %q, expect to set  %q", result, expect)
	}

	s = ""
	expectErr := "delimiter must be one character"
	result, err = ParseDelimiter(s)
	if err == nil {
		t.Errorf("no error, want error %q for %s", expectErr, "error")
	} else if err.Error() != expectErr {
		t.Errorf("error = %q, want error %q for %s", err.Error(), expectErr, "error")
	}

	s = "invalid"
	result, err = ParseDelimiter(s)
	if err == nil {
		t.Errorf("no error, want error %q for %s", expectErr, "error")
	} else if err.Error() != expectErr {
		t.Errorf("error = %q, want error %q for %s", err.Error(), expectErr, "error")
	}
}

func TestParseDelimiterPositions(t *testing.T) {
	var s string

	var expectP []int
	var expectSL bool

	s = "spaces"
	expectP = []int(nil)
	expectSL = false
	p, sl, err := ParseDelimiterPositions(s)
	if err != nil {
		t.Errorf("unexpected error: %q", err.Error())
	} else if !reflect.DeepEqual(expectP, p) || expectSL != sl {
		t.Errorf("result = %v, %t, expect to set  %v, %t", p, sl, expectP, expectSL)
	}

	s = "[1, 4, 6]"
	expectP = []int{1, 4, 6}
	expectSL = false
	p, sl, err = ParseDelimiterPositions(s)
	if err != nil {
		t.Errorf("unexpected error: %q", err.Error())
	} else if !reflect.DeepEqual(expectP, p) || expectSL != sl {
		t.Errorf("result = %v, %t, expect to set  %v, %t", p, sl, expectP, expectSL)
	}

	s = "S[1, 4, 6]"
	expectP = []int{1, 4, 6}
	expectSL = true
	p, sl, err = ParseDelimiterPositions(s)
	if err != nil {
		t.Errorf("unexpected error: %q", err.Error())
	} else if !reflect.DeepEqual(expectP, p) || expectSL != sl {
		t.Errorf("result = %v, %t, expect to set  %v, %t", p, sl, expectP, expectSL)
	}

	s = ""
	expectP = []int(nil)
	expectErr := "delimiter positions must be \"SPACES\" or a JSON array of integers"
	expectSL = false
	p, sl, err = ParseDelimiterPositions(s)
	if err == nil {
		t.Errorf("no error, want error %q for %s", expectErr, "error")
	} else if err.Error() != expectErr {
		t.Errorf("error = %q, want error %q for %s", err.Error(), expectErr, "error")
	}

	s = "invalid"
	p, sl, err = ParseDelimiterPositions(s)
	if err == nil {
		t.Errorf("no error, want error %q for %s", expectErr, "error")
	} else if err.Error() != expectErr {
		t.Errorf("error = %q, want error %q for %s", err.Error(), expectErr, "error")
	}
}

var unescapeStringBenchString = "fo\\o\\a\\b\\f\\n\\r\\t\\v\\\\\\\\'\\\"bar\\"
var unescapeStringBenchString2 = "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz"

func BenchmarkUnescapeString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = UnescapeString(unescapeStringBenchString, '\'')
	}
}

func BenchmarkUnescapeString2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = UnescapeString(unescapeStringBenchString2, '\'')
	}
}
