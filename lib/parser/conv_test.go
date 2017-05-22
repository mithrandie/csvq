package parser

import (
	"testing"

	"github.com/mithrandie/csvq/lib/ternary"
)

func TestStrToTime(t *testing.T) {
	s := "2006-01-02 15:04:05"
	if _, err := StrToTime(s); err != nil {
		t.Errorf("unexpected error %q for %q", err, s)
	}

	s = "2006-01-02T15:04:05-08:00"
	if _, err := StrToTime(s); err != nil {
		t.Errorf("unexpected error %q for %q", err, s)
	}

	s = "e"
	if _, err := StrToTime(s); err == nil {
		t.Errorf("no errors, want error for %q", s)
	}
}

func TestFloat64ToPrimary(t *testing.T) {
	var p Primary
	var f float64

	f = 1.000
	p = Float64ToPrimary(f)
	if _, ok := p.(Integer); !ok {
		t.Errorf("primary type = %T, want Integer for %d", p, f)
	}

	f = 1.234
	p = Float64ToPrimary(f)
	if _, ok := p.(Float); !ok {
		t.Errorf("primary type = %T, want Float for %d", p, f)
	}
}

func TestPrimaryToFloat(t *testing.T) {
	var p Primary
	var f Primary

	p = NewInteger(1)
	f = PrimaryToFloat(p)
	if _, ok := f.(Float); !ok {
		t.Errorf("primary type = %T, want Float for %#v", f, p)
	}

	p = NewFloat(1.234)
	f = PrimaryToFloat(p)
	if _, ok := f.(Float); !ok {
		t.Errorf("primary type = %T, want Float for %#v", f, p)
	}

	p = NewString("1")
	f = PrimaryToFloat(p)
	if _, ok := f.(Float); !ok {
		t.Errorf("primary type = %T, want Float for %#v", f, p)
	}

	p = NewString("error")
	f = PrimaryToFloat(p)
	if _, ok := f.(Null); !ok {
		t.Errorf("primary type = %T, want Null for %#v", f, p)
	}
}

func TestPrimaryToDatetime(t *testing.T) {
	var p Primary
	var dt Primary

	p = NewInteger(1136181845)
	dt = PrimaryToDatetime(p)
	if _, ok := dt.(Datetime); !ok {
		t.Errorf("primary type = %T, want Datetime for %#v", dt, p)
	}

	p = NewFloat(1136181845)
	dt = PrimaryToDatetime(p)
	if _, ok := dt.(Datetime); !ok {
		t.Errorf("primary type = %T, want Datetime for %#v", dt, p)
	}

	p = NewDatetimeFromString("2006-01-02 15:04:05")
	dt = PrimaryToDatetime(p)
	if _, ok := dt.(Datetime); !ok {
		t.Errorf("primary type = %T, want Datetime for %#v", dt, p)
	}

	p = NewString("1136181845.12345")
	dt = PrimaryToDatetime(p)
	if _, ok := dt.(Datetime); !ok {
		t.Errorf("primary type = %T, want Datetime for %#v", dt, p)
	}

	p = NewString("2006-01-02 15:04:05")
	dt = PrimaryToDatetime(p)
	if _, ok := dt.(Datetime); !ok {
		t.Errorf("primary type = %T, want Datetime for %#v", dt, p)
	}

	p = NewString("error")
	dt = PrimaryToDatetime(p)
	if _, ok := dt.(Null); !ok {
		t.Errorf("primary type = %T, want Null for %#v", dt, p)
	}
}

func TestPrimaryToBoolean(t *testing.T) {
	var p Primary
	var b Primary

	p = NewBoolean(true)
	b = PrimaryToBoolean(p)
	if _, ok := b.(Boolean); !ok {
		t.Errorf("primary type = %T, want Boolean for %#v", b, p)
	}

	p = NewTernary(ternary.TRUE)
	b = PrimaryToBoolean(p)
	if _, ok := b.(Boolean); !ok {
		t.Errorf("primary type = %T, want Boolean for %#v", b, p)
	}

	p = NewString("true")
	b = PrimaryToBoolean(p)
	if _, ok := b.(Boolean); !ok {
		t.Errorf("primary type = %T, want Boolean for %#v", b, p)
	}

	p = NewTernary(ternary.UNKNOWN)
	b = PrimaryToBoolean(p)
	if _, ok := b.(Null); !ok {
		t.Errorf("primary type = %T, want Null for %#v", b, p)
	}

	p = NewString("error")
	b = PrimaryToBoolean(p)
	if _, ok := b.(Null); !ok {
		t.Errorf("primary type = %T, want Null for %#v", b, p)
	}
}

func TestPrimaryToString(t *testing.T) {
	var p Primary
	var s Primary

	p = NewString("str")
	s = PrimaryToString(p)
	if _, ok := s.(String); !ok {
		t.Errorf("primary type = %T, want String for %#v", s, p)
	}

	p = NewInteger(1)
	s = PrimaryToString(p)
	if _, ok := s.(String); !ok {
		t.Errorf("primary type = %T, want String for %#v", s, p)
	}

	p = NewFloat(1)
	s = PrimaryToString(p)
	if _, ok := s.(String); !ok {
		t.Errorf("primary type = %T, want String for %#v", s, p)
	}

	p = NewDatetimeFromString("2006-01-02 15:04:05")
	s = PrimaryToString(p)
	if _, ok := s.(Null); !ok {
		t.Errorf("primary type = %T, want Null for %#v", s, p)
	}
}
