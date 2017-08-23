package parser

import (
	"testing"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/ternary"
)

func TestStrToTime(t *testing.T) {
	flags := cmd.GetFlags()
	flags.DatetimeFormat = "01/02/2006"

	s := "01/02/2006"
	if _, err := StrToTime(s); err != nil {
		t.Errorf("unexpected error %q for %q", err, s)
	}

	s = "2006-01-02 15:04:05"
	if _, err := StrToTime(s); err != nil {
		t.Errorf("unexpected error %q for %q", err, s)
	}

	s = "2006-01-02"
	if _, err := StrToTime(s); err != nil {
		t.Errorf("unexpected error %q for %q", err, s)
	}

	s = "2006-01-02 15:04:05 -08:00"
	if _, err := StrToTime(s); err != nil {
		t.Errorf("unexpected error %q for %q", err, s)
	}

	s = "2006-01-02 15:04:05 -0800"
	if _, err := StrToTime(s); err != nil {
		t.Errorf("unexpected error %q for %q", err, s)
	}

	s = "2006-01-02 15:04:05 PST"
	if _, err := StrToTime(s); err != nil {
		t.Errorf("unexpected error %q for %q", err, s)
	}

	s = "2006/01/02 15:04:05"
	if _, err := StrToTime(s); err != nil {
		t.Errorf("unexpected error %q for %q", err, s)
	}

	s = "2006/01/02"
	if _, err := StrToTime(s); err != nil {
		t.Errorf("unexpected error %q for %q", err, s)
	}

	s = "2006/01/02 15:04:05 -08:00"
	if _, err := StrToTime(s); err != nil {
		t.Errorf("unexpected error %q for %q", err, s)
	}

	s = "2006/01/02 15:04:05 -0800"
	if _, err := StrToTime(s); err != nil {
		t.Errorf("unexpected error %q for %q", err, s)
	}

	s = "2006/01/02 15:04:05 -0800"
	if _, err := StrToTime(s); err != nil {
		t.Errorf("unexpected error %q for %q", err, s)
	}

	s = "2006/11/2 15:04:05 -0800"
	if _, err := StrToTime(s); err != nil {
		t.Errorf("unexpected error %q for %q", err, s)
	}

	s = "2006/01/02 15:04:05 PST"
	if _, err := StrToTime(s); err != nil {
		t.Errorf("unexpected error %q for %q", err, s)
	}

	s = "2006-1-2 15:04:05"
	if _, err := StrToTime(s); err != nil {
		t.Errorf("unexpected error %q for %q", err, s)
	}

	s = "2006-1-2"
	if _, err := StrToTime(s); err != nil {
		t.Errorf("unexpected error %q for %q", err, s)
	}

	s = "2006-1-2 15:04:05 -08:00"
	if _, err := StrToTime(s); err != nil {
		t.Errorf("unexpected error %q for %q", err, s)
	}

	s = "2006-1-2 15:04:05 -0800"
	if _, err := StrToTime(s); err != nil {
		t.Errorf("unexpected error %q for %q", err, s)
	}

	s = "2006-1-2 15:04:05 PST"
	if _, err := StrToTime(s); err != nil {
		t.Errorf("unexpected error %q for %q", err, s)
	}

	s = "2006/1/2 15:04:05"
	if _, err := StrToTime(s); err != nil {
		t.Errorf("unexpected error %q for %q", err, s)
	}

	s = "2006/1/2"
	if _, err := StrToTime(s); err != nil {
		t.Errorf("unexpected error %q for %q", err, s)
	}

	s = "2006/1/2 15:04:05 -08:00"
	if _, err := StrToTime(s); err != nil {
		t.Errorf("unexpected error %q for %q", err, s)
	}

	s = "2006/1/2 15:04:05 -0800"
	if _, err := StrToTime(s); err != nil {
		t.Errorf("unexpected error %q for %q", err, s)
	}

	s = "2006/1/2 15:04:05 PST"
	if _, err := StrToTime(s); err != nil {
		t.Errorf("unexpected error %q for %q", err, s)
	}

	s = "02 Jan 06 15:04 PDT"
	if _, err := StrToTime(s); err != nil {
		t.Errorf("unexpected error %q for %q", err, s)
	}

	s = "02 Jan 06 15:04 -0700"
	if _, err := StrToTime(s); err != nil {
		t.Errorf("unexpected error %q for %q", err, s)
	}

	s = "2006-01-02T15:04:05-08:00"
	if _, err := StrToTime(s); err != nil {
		t.Errorf("unexpected error %q for %q", err, s)
	}

	s = "2006-01-02T15:04:05"
	if _, err := StrToTime(s); err != nil {
		t.Errorf("unexpected error %q for %q", err, s)
	}

	s = "e"
	if _, err := StrToTime(s); err == nil {
		t.Errorf("no errors, want error for %q", s)
	}
}

var convertDatetimeFormatTests = []struct {
	Datetime string
	Format   string
	Result   string
}{
	{
		Format: "datetime: %Y-%m-%d %H:%i:%s %% %g",
		Result: "datetime: 2006-01-02 15:04:05 % g",
	},
	{
		Format: "%a",
		Result: "Mon",
	},
	{
		Format: "%b",
		Result: "Jan",
	},
	{
		Format: "%c",
		Result: "1",
	},
	{
		Format: "%E",
		Result: "_2",
	},
	{
		Format: "%e",
		Result: "2",
	},
	{
		Format: "%F",
		Result: ".999999",
	},
	{
		Format: "%f",
		Result: ".000000",
	},
	{
		Format: "%h",
		Result: "03",
	},
	{
		Format: "%l",
		Result: "3",
	},
	{
		Format: "%M",
		Result: "January",
	},
	{
		Format: "%N",
		Result: ".999999999",
	},
	{
		Format: "%n",
		Result: ".000000000",
	},
	{
		Format: "%p",
		Result: "PM",
	},
	{
		Format: "%r",
		Result: "03:04:05 PM",
	},
	{
		Format: "%T",
		Result: "15:04:05",
	},
	{
		Format: "%W",
		Result: "Monday",
	},
	{
		Format: "%y",
		Result: "06",
	},
	{
		Format: "%Z",
		Result: "Z07:00",
	},
	{
		Format: "%z",
		Result: "MST",
	},
}

func TestConvertDatetimeFormat(t *testing.T) {
	for _, v := range convertDatetimeFormatTests {
		converted := ConvertDatetimeFormat(v.Format)
		if converted != v.Result {
			t.Errorf("result = %q, want %q for %q", converted, v.Result, v.Format)
		}
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

func TestPrimaryToInteger(t *testing.T) {
	var p Primary
	var i Primary

	p = NewInteger(1)
	i = PrimaryToInteger(p)
	if _, ok := i.(Integer); !ok {
		t.Errorf("primary type = %T, want Integer for %#v", i, p)
	}

	p = NewFloat(1)
	i = PrimaryToInteger(p)
	if _, ok := i.(Integer); !ok {
		t.Errorf("primary type = %T, want Integer for %#v", i, p)
	}

	p = NewFloat(1.6)
	i = PrimaryToInteger(p)
	if _, ok := i.(Null); !ok {
		t.Errorf("primary type = %T, want Null for %#v", i, p)
	}

	p = NewString("1")
	i = PrimaryToInteger(p)
	if _, ok := i.(Integer); !ok {
		t.Errorf("primary type = %T, want Integer for %#v", i, p)
	}

	p = NewString("-1")
	i = PrimaryToInteger(p)
	if _, ok := i.(Integer); !ok {
		t.Errorf("primary type = %T, want Integer for %#v", i, p)
	}

	p = NewString("1e+02")
	i = PrimaryToInteger(p)
	if _, ok := i.(Integer); !ok {
		t.Errorf("primary type = %T, want Integer for %#v", i, p)
	}

	p = NewString("1.5")
	i = PrimaryToInteger(p)
	if _, ok := i.(Null); !ok {
		t.Errorf("primary type = %T, want Null for %#v", i, p)
	}

	p = NewString("error")
	i = PrimaryToInteger(p)
	if _, ok := i.(Null); !ok {
		t.Errorf("primary type = %T, want Null for %#v", i, p)
	}

	p = NewString("error")
	i = PrimaryToInteger(p)
	if _, ok := i.(Null); !ok {
		t.Errorf("primary type = %T, want Null for %#v", i, p)
	}

	p = NewString("2002-02-02")
	i = PrimaryToInteger(p)
	if _, ok := i.(Null); !ok {
		t.Errorf("primary type = %T, want Null for %#v", i, p)
	}

	p = NewString("2002/02/02")
	i = PrimaryToInteger(p)
	if _, ok := i.(Null); !ok {
		t.Errorf("primary type = %T, want Null for %#v", i, p)
	}

	p = NewString("03 Mar 12 12:03 PST")
	i = PrimaryToInteger(p)
	if _, ok := i.(Null); !ok {
		t.Errorf("primary type = %T, want Null for %#v", i, p)
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

	p = NewString("1136181845")
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

	p = NewInteger(1)
	b = PrimaryToBoolean(p)
	if _, ok := b.(Boolean); !ok {
		t.Errorf("primary type = %T, want Boolean for %#v", b, p)
	}

	p = NewFloat(0)
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

func BenchmarkStrToTime1(b *testing.B) {
	flags := cmd.GetFlags()
	flags.DatetimeFormat = "01/02/2006"

	for i := 0; i < b.N; i++ {
		s := "01/02/2006"
		StrToTime(s)
	}
}

func BenchmarkStrToTime2(b *testing.B) {
	flags := cmd.GetFlags()
	flags.DatetimeFormat = ""

	for i := 0; i < b.N; i++ {
		s := "2006-01-02T15:04:05-07:00"
		StrToTime(s)
	}
}

func BenchmarkStrToTime3(b *testing.B) {
	flags := cmd.GetFlags()
	flags.DatetimeFormat = ""

	for i := 0; i < b.N; i++ {
		s := "2006-01-02"
		StrToTime(s)
	}
}

func BenchmarkStrToTime4(b *testing.B) {
	flags := cmd.GetFlags()
	flags.DatetimeFormat = ""

	for i := 0; i < b.N; i++ {
		s := "2006-01-02 15:04:05"
		StrToTime(s)
	}
}

func BenchmarkStrToTime5(b *testing.B) {
	flags := cmd.GetFlags()
	flags.DatetimeFormat = ""

	for i := 0; i < b.N; i++ {
		s := "2006-01-02 15:04:05 -0700"
		StrToTime(s)
	}
}

func BenchmarkStrToTime6(b *testing.B) {
	flags := cmd.GetFlags()
	flags.DatetimeFormat = ""

	for i := 0; i < b.N; i++ {
		s := "02 Jan 06 15:04 PDT"
		StrToTime(s)
	}
}

func BenchmarkStrToTime7(b *testing.B) {
	flags := cmd.GetFlags()
	flags.DatetimeFormat = ""

	for i := 0; i < b.N; i++ {
		s := "abcdefghijklmnopq"
		StrToTime(s)
	}
}

func BenchmarkPrimaryToInteger(b *testing.B) {
	for i := 0; i < b.N; i++ {
		p := NewString("a")
		PrimaryToInteger(p)
	}
}

func BenchmarkPrimaryToInteger2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		p := NewString("2012-02-02")
		PrimaryToInteger(p)
	}
}

func BenchmarkPrimaryToFloat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		p := NewString("a")
		PrimaryToFloat(p)
	}
}

func BenchmarkPrimaryToFloat2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		p := NewString("2012-02-02")
		PrimaryToFloat(p)
	}
}
