package json

import (
	"testing"
)

func TestEncodeRune(t *testing.T) {
	r := 'a'
	expect := "\\u0061"
	result := string(EncodeRune(r))
	if expect != result {
		t.Errorf("result = %q, want %q", result, expect)
	}

	r = '𝄞'
	expect = "\\ud834\\udd1e"
	result = string(EncodeRune(r))
	if expect != result {
		t.Errorf("result = %q, want %q", result, expect)
	}
}

func TestEscape(t *testing.T) {
	s := "abc\u0022\u005c\u002f\u0008\u000c\u000a\u000d\u0009\u001f𝄞あ"
	expect := "abc\\\"\\\\\\/\\b\\f\\n\\r\\t\\u001f𝄞あ"
	result := Escape(s)
	if expect != result {
		t.Errorf("result = %q, want %q", result, expect)
	}
}

func TestEscapeWithHexDigits(t *testing.T) {
	s := "abc\u0022\u005c\u002f\u0008\u000c\u000a\u000d\u0009\u001f𝄞あ"
	expect := "abc\\u0022\\u005c\\u002f\\u0008\\u000c\\u000a\\u000d\\u0009\\u001f𝄞あ"
	result := EscapeWithHexDigits(s)
	if expect != result {
		t.Errorf("result = %q, want %q", result, expect)
	}
}

func TestEscapeAll(t *testing.T) {
	s := "abc\u0022\u005c\u002f\u0008\u000c\u000a\u000d\u0009\u001f𝄞あ"
	expect := "\\u0061\\u0062\\u0063\\u0022\\u005c\\u002f\\u0008\\u000c\\u000a\\u000d\\u0009\\u001f\\ud834\\udd1e\\u3042"
	result := EscapeAll(s)
	if expect != result {
		t.Errorf("result = %q, want %q", result, expect)
	}
}

func TestUnescape(t *testing.T) {
	s := "abc\\a\\\"\\\\\\/\\b\\f\\n\\r\\t\\u001F\\u0022\\ud834\\udd1e\\u000\\u3042\\u000"
	expect := "abca\u0022\u005c\u002f\u0008\u000c\u000a\u000d\u0009\u001f\u0022𝄞u000あu000"
	expectEscapeType := AllWithHexDigits
	result, escapeType := Unescape(s)
	if expect != result {
		t.Errorf("result = %q, want %q", result, expect)
	}
	if escapeType != expectEscapeType {
		t.Errorf("escape type = %d, want %d", escapeType, expectEscapeType)
	}
}
