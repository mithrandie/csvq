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
	str := "\\a\\b\\f\\n\\r\\t\\v\\\\\\\""
	expect := "\a\b\f\n\r\t\v\\\""
	unescaped := UnescapeString(str)
	if unescaped != expect {
		t.Errorf("unescaped string = %q, want %q", unescaped, expect)
	}
}
