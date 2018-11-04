package cmd

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestGetReader(t *testing.T) {
	fp := os.Stdout

	r := GetReader(fp, UTF8)
	if reflect.TypeOf(r).String() != "*os.File" {
		t.Errorf("reader = %q, want %q", reflect.TypeOf(r).String(), "*os.File")
	}

	r = GetReader(fp, SJIS)
	if reflect.TypeOf(r).String() != "*transform.Reader" {
		t.Errorf("reader = %q, want %q", reflect.TypeOf(r).String(), "*transform.Reader")
	}
}

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

func TestHumarizeNumber(t *testing.T) {
	number := "1234567"
	expect := "1,234,567"
	result := HumarizeNumber(number)
	if result != expect {
		t.Errorf("humarized = %q, want %q", result, expect)
	}

	number = "123456"
	expect = "123,456"
	result = HumarizeNumber(number)
	if result != expect {
		t.Errorf("humarized = %q, want %q", result, expect)
	}

	number = "123"
	expect = "123"
	result = HumarizeNumber(number)
	if result != expect {
		t.Errorf("humarized = %q, want %q", result, expect)
	}

	number = "1234.5678"
	expect = "1,234.5678"
	result = HumarizeNumber(number)
	if result != expect {
		t.Errorf("humarized = %q, want %q", result, expect)
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
	if e != UTF8 {
		t.Errorf("encoding = %s, expect to set %s for %s", e, UTF8, "utf8")
	}

	e, err = ParseEncoding("sjis")
	if err != nil {
		t.Errorf("unexpected error: %q", err.Error())
	}
	if e != SJIS {
		t.Errorf("encoding = %s, expect to set %s for %s", e, SJIS, "sjis")
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
