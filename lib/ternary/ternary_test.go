package ternary

import (
	"testing"
)

func TestValue_String(t *testing.T) {
	s := FALSE.String()
	if s != "FALSE" {
		t.Errorf("string = %q, want %q for %s.String()", s, "FALSE", FALSE)
	}

	s = UNKNOWN.String()
	if s != "UNKNOWN" {
		t.Errorf("string = %q, want %q for %s.String()", s, "UNKNOWN", UNKNOWN)
	}

	s = TRUE.String()
	if s != "TRUE" {
		t.Errorf("string = %q, want %q for %s.String()", s, "TRUE", TRUE)
	}
}

func TestValue_EqualTo(t *testing.T) {
	r := FALSE.EqualTo(TRUE)
	if r != FALSE {
		t.Errorf("ternary = %s, want %s for %s.EqualTo(%s)", r, FALSE, FALSE, TRUE)
	}

	r = UNKNOWN.EqualTo(UNKNOWN)
	if r != TRUE {
		t.Errorf("ternary = %s, want %s for %s.EqualTo(%s)", r, TRUE, UNKNOWN, UNKNOWN)
	}
}

func TestValue_BoolValue(t *testing.T) {
	b := FALSE.BoolValue()
	if b != false {
		t.Errorf("bool value = %t, want %t for %s", b, false, FALSE)
	}

	b = UNKNOWN.BoolValue()
	if b != false {
		t.Errorf("bool value = %t, want %t for %s", b, false, UNKNOWN)
	}

	b = TRUE.BoolValue()
	if b != true {
		t.Errorf("bool value = %t, want %t for %s", b, true, TRUE)
	}
}

var parseTests = []struct {
	Str    string
	Result Value
	Err    string
}{
	{
		Str:    "false",
		Result: FALSE,
	},
	{
		Str:    "unknown",
		Result: UNKNOWN,
	},
	{
		Str:    "true",
		Result: TRUE,
	},
	{
		Str:    "-1",
		Result: FALSE,
	},
	{
		Str:    "0",
		Result: UNKNOWN,
	},
	{
		Str:    "1",
		Result: TRUE,
	},
	{
		Str: "ParseError",
		Err: "parsing \"ParseError\": invalid syntax",
	},
}

func TestParse(t *testing.T) {
	for _, test := range parseTests {
		v, err := Parse(test.Str)
		if err != nil {
			if len(test.Err) < 1 {
				t.Errorf("unexpected error: %q", err.Error())
			} else if err.Error() != test.Err {
				t.Errorf("error = %q, want error %q for %s", err.Error(), test.Err, test.Str)
			}
			continue
		}
		if 0 < len(test.Err) {
			t.Errorf("no error, want error %q for %s", test.Err, test.Str)
			continue
		}
		if v != test.Result {
			t.Errorf("ternary = %s, want %s for %q", v, test.Result, test.Str)
		}
	}
}

func TestParseBool(t *testing.T) {
	r := ParseBool(false)
	if r != FALSE {
		t.Errorf("ternary = %s, want %s for %t", r, FALSE, false)
	}

	r = ParseBool(true)
	if r != TRUE {
		t.Errorf("ternary = %s, want %s for %t", r, TRUE, true)
	}
}

var notTests = []struct {
	Value  Value
	Result Value
}{
	{
		Value:  FALSE,
		Result: TRUE,
	},
	{
		Value:  TRUE,
		Result: FALSE,
	},
	{
		Value:  UNKNOWN,
		Result: UNKNOWN,
	},
}

func TestNot(t *testing.T) {
	for _, test := range notTests {
		v := Not(test.Value)
		if v != test.Result {
			t.Errorf("ternary = %s, want %s for \"not %s\"", v, test.Result, test.Value)
		}
	}
}

var andTests = []struct {
	Value1 Value
	Value2 Value
	Result Value
}{
	{
		Value1: FALSE,
		Value2: FALSE,
		Result: FALSE,
	},
	{
		Value1: FALSE,
		Value2: UNKNOWN,
		Result: FALSE,
	},
	{
		Value1: FALSE,
		Value2: TRUE,
		Result: FALSE,
	},
	{
		Value1: UNKNOWN,
		Value2: FALSE,
		Result: FALSE,
	},
	{
		Value1: UNKNOWN,
		Value2: UNKNOWN,
		Result: UNKNOWN,
	},
	{
		Value1: UNKNOWN,
		Value2: TRUE,
		Result: UNKNOWN,
	},
	{
		Value1: TRUE,
		Value2: FALSE,
		Result: FALSE,
	},
	{
		Value1: TRUE,
		Value2: UNKNOWN,
		Result: UNKNOWN,
	},
	{
		Value1: TRUE,
		Value2: TRUE,
		Result: TRUE,
	},
}

func TestAnd(t *testing.T) {
	for _, test := range andTests {
		v := And(test.Value1, test.Value2)
		if v != test.Result {
			t.Errorf("ternary = %s, want %s for \"%s and %s\"", v, test.Result, test.Value1, test.Value2)
		}
	}
}

var orTests = []struct {
	Value1 Value
	Value2 Value
	Result Value
}{
	{
		Value1: FALSE,
		Value2: FALSE,
		Result: FALSE,
	},
	{
		Value1: FALSE,
		Value2: UNKNOWN,
		Result: UNKNOWN,
	},
	{
		Value1: FALSE,
		Value2: TRUE,
		Result: TRUE,
	},
	{
		Value1: UNKNOWN,
		Value2: FALSE,
		Result: UNKNOWN,
	},
	{
		Value1: UNKNOWN,
		Value2: UNKNOWN,
		Result: UNKNOWN,
	},
	{
		Value1: UNKNOWN,
		Value2: TRUE,
		Result: TRUE,
	},
	{
		Value1: TRUE,
		Value2: FALSE,
		Result: TRUE,
	},
	{
		Value1: TRUE,
		Value2: UNKNOWN,
		Result: TRUE,
	},
	{
		Value1: TRUE,
		Value2: TRUE,
		Result: TRUE,
	},
}

func TestOr(t *testing.T) {
	for _, test := range orTests {
		v := Or(test.Value1, test.Value2)
		if v != test.Result {
			t.Errorf("ternary = %s, want %s for \"%s or %s\"", v, test.Result, test.Value1, test.Value2)
		}
	}
}
