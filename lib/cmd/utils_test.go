package cmd

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/mithrandie/go-text"
)

func TestEscapeString(t *testing.T) {
	str := "fo\\o\a\b\f\n\r\t\v\\\\'\"bar\\"
	expect := "fo\\\\o\\a\\b\\f\\n\\r\\t\\v\\\\\\\\\\'\\\"bar\\\\"
	unescaped := EscapeString(str)
	if unescaped != expect {
		t.Errorf("escaped string = %q, want %q", unescaped, expect)
	}
}

func TestUnescapeString(t *testing.T) {
	str := "fo\\o\\a\\b\\f\\n\\r\\t\\v\\\\\\\\'\\\"bar\\"
	expect := "fo\\o\a\b\f\n\r\t\v\\\\'\"bar\\"
	unescaped := UnescapeString(str)
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
	str := "fo\\o\\a\\b\\f\\n\\r\\t\\v\\\\\\\\`bar\\"
	expect := "fo\\o\a\b\f\n\r\t\v\\\\`bar\\"
	unescaped := UnescapeIdentifier(str)
	if unescaped != expect {
		t.Errorf("unescaped identifier = %q, want %q", unescaped, expect)
	}
}

func TestQuoteString(t *testing.T) {
	s := "abc'def"
	expect := "\"abc\\'def\""
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

func TestIsReadableFromPipeOrRedirection(t *testing.T) {
	oldStdin := os.Stdin
	r, _ := os.Open(filepath.Join(TestDataDir, "empty.txt"))
	os.Stdin = r

	result := IsReadableFromPipeOrRedirection()

	r.Close()

	if result != false {
		t.Errorf("readable from pipe or redirection = %t, want %t", result, false)
	}

	oldStdin = os.Stdin
	r, w, _ := os.Pipe()
	os.Stdin = r

	w.Write([]byte("abcde"))
	w.Close()

	result = IsReadableFromPipeOrRedirection()

	r.Close()
	os.Stdin = oldStdin

	if result != true {
		t.Errorf("readable from pipe or redirection = %t, want %t", result, true)
	}
}

var unescapeStringBenchString = "fo\\o\\a\\b\\f\\n\\r\\t\\v\\\\\\\\'\\\"bar\\"
var unescapeStringBenchString2 = "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz"

func BenchmarkUnescapeString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = UnescapeString(unescapeStringBenchString)
	}
}

func BenchmarkUnescapeString2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = UnescapeString(unescapeStringBenchString2)
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

	expectErr := "encoding must be one of UTF8|SJIS"
	_, err = ParseEncoding("error")
	if err == nil {
		t.Errorf("no error, want error %q for %s", expectErr, "error")
	} else if err.Error() != expectErr {
		t.Errorf("error = %q, want error %q for %s", err.Error(), expectErr, "error")
	}
}

func TestParseDelimiter(t *testing.T) {
	var s string
	var delimiter rune
	var delimiterPositions []int
	var delimitAutomatically bool

	var expectD rune
	var expectP []int
	var expectA bool

	s = "\t"
	delimiter = ','
	delimiterPositions = []int{1, 3, 5}
	delimitAutomatically = true

	expectD = '\t'
	expectP = []int(nil)
	expectA = false
	d, p, a, err := ParseDelimiter(s, delimiter, delimiterPositions, delimitAutomatically)
	if err != nil {
		t.Errorf("unexpected error: %q", err.Error())
	} else if expectD != d || !reflect.DeepEqual(expectP, p) || expectA != a {
		t.Errorf("result = %q, %v, %t, expect to set  %q, %v, %t", d, p, a, expectD, expectP, expectA)
	}

	s = "spaces"
	delimiter = ','
	delimiterPositions = []int{1, 3, 5}
	delimitAutomatically = true

	expectD = ','
	expectP = []int(nil)
	expectA = true
	d, p, a, err = ParseDelimiter(s, delimiter, delimiterPositions, delimitAutomatically)
	if err != nil {
		t.Errorf("unexpected error: %q", err.Error())
	} else if expectD != d || !reflect.DeepEqual(expectP, p) || expectA != a {
		t.Errorf("result = %q, %v, %t, expect to set  %q, %v, %t", d, p, a, expectD, expectP, expectA)
	}

	s = "[1, 4, 6]"
	delimiter = ','
	delimiterPositions = nil
	delimitAutomatically = false

	expectD = ','
	expectP = []int{1, 4, 6}
	expectA = false
	d, p, a, err = ParseDelimiter(s, delimiter, delimiterPositions, delimitAutomatically)
	if err != nil {
		t.Errorf("unexpected error: %q", err.Error())
	} else if expectD != d || !reflect.DeepEqual(expectP, p) || expectA != a {
		t.Errorf("result = %q, %v, %t, expect to set  %q, %v, %t", d, p, a, expectD, expectP, expectA)
	}

	s = ""
	delimiter = ','
	delimiterPositions = []int(nil)
	delimitAutomatically = false

	expectErr := "delimiter must be one character, \"SPACES\" or JSON array of integers"
	d, p, a, err = ParseDelimiter(s, delimiter, delimiterPositions, delimitAutomatically)
	if err == nil {
		if err == nil {
			t.Errorf("no error, want error %q for %s", expectErr, "error")
		} else if err.Error() != expectErr {
			t.Errorf("error = %q, want error %q for %s", err.Error(), expectErr, "error")
		}
	}

	s = "invalid"
	delimiter = ','
	delimiterPositions = []int(nil)
	delimitAutomatically = false

	expectErr = "delimiter must be one character, \"SPACES\" or JSON array of integers"
	d, p, a, err = ParseDelimiter(s, delimiter, delimiterPositions, delimitAutomatically)
	if err == nil {
		if err == nil {
			t.Errorf("no error, want error %q for %s", expectErr, "error")
		} else if err.Error() != expectErr {
			t.Errorf("error = %q, want error %q for %s", err.Error(), expectErr, "error")
		}
	}
}
