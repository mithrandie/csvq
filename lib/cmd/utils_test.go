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
	str := "fo\\o\\a\\b\\f\\n\\r\\t\\v\\\\\\\"bar\\"
	expect := "fo\\o\a\b\f\n\r\t\v\\\"bar\\"
	unescaped := UnescapeString(str)
	if unescaped != expect {
		t.Errorf("unescaped string = %q, want %q", unescaped, expect)
	}
}
