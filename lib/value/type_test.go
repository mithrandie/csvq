package value

import (
	"testing"
	"time"

	"github.com/mithrandie/ternary"
)

func TestIsNull(t *testing.T) {
	var p Primary

	p = NewInteger(1)
	if IsNull(p) {
		t.Errorf("value %#p is evaluated as is a null, but it is not so", p)
	}

	p = NewNull()
	if !IsNull(p) {
		t.Errorf("value %#p is evaluated as is not a null, but it is so", p)
	}
}

func TestString_String(t *testing.T) {
	s := "abcde"
	p := NewString(s)
	expect := "'" + s + "'"
	if p.String() != expect {
		t.Errorf("string = %q, want %q for %#v", p.String(), expect, p)
	}
}

func TestString_Value(t *testing.T) {
	s := "abcde"
	p := NewString(s)
	if p.Raw() != s {
		t.Errorf("value = %q, want %q for %#v", p.Raw(), s, p)
	}
}

func TestString_Ternary(t *testing.T) {
	s := " 1"
	p := NewString(s)
	if p.Ternary() != ternary.TRUE {
		t.Errorf("ternary = %s, want %s for %#v", p.Ternary(), ternary.TRUE, p)
	}

	s = "0"
	p = NewString(s)
	if p.Ternary() != ternary.FALSE {
		t.Errorf("ternary = %s, want %s for %#v", p.Ternary(), ternary.FALSE, p)
	}
	s = "unknown"
	p = NewString(s)
	if p.Ternary() != ternary.UNKNOWN {
		t.Errorf("ternary = %s, want %s for %#v", p.Ternary(), ternary.UNKNOWN, p)
	}
}

func TestInteger_String(t *testing.T) {
	s := "1"
	p := NewInteger(1)
	if p.String() != s {
		t.Errorf("string = %q, want %q for %#v", p.String(), s, p)
	}
}

func TestInteger_Value(t *testing.T) {
	i := NewInteger(1)
	expect := int64(1)

	if i.Raw() != expect {
		t.Errorf("value = %d, want %d for %#v", i.Raw(), expect, i)
	}
}

func TestInteger_Ternary(t *testing.T) {
	p := NewInteger(1)
	if p.Ternary() != ternary.TRUE {
		t.Errorf("ternary = %s, want %s for %#v", p.Ternary(), ternary.TRUE, p)
	}
	p = NewInteger(0)
	if p.Ternary() != ternary.FALSE {
		t.Errorf("ternary = %s, want %s for %#v", p.Ternary(), ternary.FALSE, p)
	}
	p = NewInteger(2)
	if p.Ternary() != ternary.UNKNOWN {
		t.Errorf("ternary = %s, want %s for %#v", p.Ternary(), ternary.UNKNOWN, p)
	}
}

func TestFloat_String(t *testing.T) {
	s := "1.234"
	p := NewFloat(1.234)
	if p.String() != s {
		t.Errorf("string = %q, want %q for %#v", p.String(), s, p)
	}
}

func TestFloat_Value(t *testing.T) {
	f := NewFloat(1.234)
	expect := float64(1.234)

	if f.Raw() != expect {
		t.Errorf("value = %f, want %f for %#v", f.Raw(), expect, f)
	}
}

func TestFloat_Ternary(t *testing.T) {
	p := NewFloat(1)
	if p.Ternary() != ternary.TRUE {
		t.Errorf("ternary = %s, want %s for %#v", p.Ternary(), ternary.TRUE, p)
	}
	p = NewFloat(0)
	if p.Ternary() != ternary.FALSE {
		t.Errorf("ternary = %s, want %s for %#v", p.Ternary(), ternary.FALSE, p)
	}
	p = NewFloat(2)
	if p.Ternary() != ternary.UNKNOWN {
		t.Errorf("ternary = %s, want %s for %#v", p.Ternary(), ternary.UNKNOWN, p)
	}
}

func TestBoolean_String(t *testing.T) {
	s := "true"
	p := NewBoolean(true)
	if p.String() != s {
		t.Errorf("string = %q, want %q for %#v", p.String(), s, p)
	}
}

func TestBoolean_Value(t *testing.T) {
	p := NewBoolean(true)
	if p.Raw() != true {
		t.Errorf("bool = %t, want %t for %#v", p.Raw(), true, p)
	}
}

func TestBoolean_Ternary(t *testing.T) {
	p := NewBoolean(true)
	if p.Ternary() != ternary.TRUE {
		t.Errorf("ternary = %s, want %s for %#v", p.Ternary(), ternary.TRUE, p)
	}
}

func TestTernary_String(t *testing.T) {
	s := "TRUE"
	p := NewTernary(ternary.TRUE)
	if p.String() != s {
		t.Errorf("string = %q, want %q for %#v", p.String(), s, p)
	}
}

func TestTernary_Ternary(t *testing.T) {
	p := NewTernary(ternary.TRUE)
	if p.Ternary() != ternary.TRUE {
		t.Errorf("ternary = %s, want %s for %#v", p.Ternary(), ternary.TRUE, p)
	}
}

func TestDatetime_String(t *testing.T) {
	s := "2012-01-01T12:34:56Z"
	p := NewDatetimeFromString(s)

	expect := "'" + s + "'"
	if p.String() != expect {
		t.Errorf("string = %q, want %q for %#v", p.String(), expect, p)
	}
}

func TestDatetime_Value(t *testing.T) {
	d := NewDatetimeFromString("2012-01-01 12:34:56")
	expect := time.Date(2012, time.January, 1, 12, 34, 56, 0, time.Local)

	if d.Raw() != expect {
		t.Errorf("value = %v, want %v for %#v", d.Raw(), expect, d)
	}

	d = NewDatetimeFromString("2012-01-01T12:34:56-08:00")
	l, _ := time.LoadLocation("America/Los_Angeles")
	expect = time.Date(2012, time.January, 1, 12, 34, 56, 0, l)

	if d.Raw().Sub(expect).Seconds() != 0 {
		t.Errorf("value = %v, want %v for %#v", d.Raw(), expect, d)
	}
}

func TestDatetime_Ternary(t *testing.T) {
	p := NewDatetimeFromString("2012-01-01T12:34:56-08:00")
	if p.Ternary() != ternary.UNKNOWN {
		t.Errorf("ternary = %s, want %s for %#v", p.Ternary(), ternary.UNKNOWN, p)
	}
}

func TestDatetime_Format(t *testing.T) {
	dtstring := "2012-08-01T04:03:05.123-08:00"
	dt := NewDatetimeFromString(dtstring)
	expect := "2012-08-01T04:03:05-08:00"
	if dt.Format(time.RFC3339) != expect {
		t.Errorf("result = %q, want %q for %q ", dt.Format(time.RFC3339), expect, dtstring)
	}
}

func TestNull_String(t *testing.T) {
	p := NewNull()
	if p.String() != "NULL" {
		t.Errorf("string = %q, want %q for %#v", p.String(), "NULL", p)
	}
}

func TestNull_Ternary(t *testing.T) {
	p := NewNull()
	if p.Ternary() != ternary.UNKNOWN {
		t.Errorf("ternary = %s, want %s for %#v", p.Ternary(), ternary.UNKNOWN, p)
	}
}
