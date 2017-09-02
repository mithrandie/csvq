package cmd

import (
	"os"
	"reflect"
	"testing"
)

func TestGetReader(t *testing.T) {
	fp := os.Stdout

	r := GetReader(fp, UTF8)
	if reflect.TypeOf(r).String() != "*bufio.Reader" {
		t.Errorf("reader = %q, want %q", reflect.TypeOf(r).String(), "*bufio.Reader")
	}

	r = GetReader(fp, SJIS)
	if reflect.TypeOf(r).String() != "*transform.Reader" {
		t.Errorf("reader = %q, want %q", reflect.TypeOf(r).String(), "*transform.Reader")
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
