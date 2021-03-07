package value

import (
	"testing"
	"time"

	"github.com/mithrandie/csvq/lib/cmd"

	"github.com/mithrandie/ternary"
)

func TestStrToTime(t *testing.T) {
	formats := []string{"01/02/2006"}
	location, _ := time.LoadLocation("UTC")

	s := "01/02/2006"
	if _, ok := StrToTime(s, formats, location); !ok {
		t.Errorf("failed, want to success for %q", s)
	}

	s = "2006-01-02 15:04:05"
	if _, ok := StrToTime(s, formats, location); !ok {
		t.Errorf("failed, want to success for %q", s)
	}

	s = "2006-01-02"
	if _, ok := StrToTime(s, formats, location); !ok {
		t.Errorf("failed, want to success for %q", s)
	}

	s = "2006-01-02 15:04:05 -08:00"
	if _, ok := StrToTime(s, formats, location); !ok {
		t.Errorf("failed, want to success for %q", s)
	}

	s = "2006-01-02 15:04:05 -0800"
	if _, ok := StrToTime(s, formats, location); !ok {
		t.Errorf("failed, want to success for %q", s)
	}

	s = "2006-01-02 15:04:05 PST"
	if _, ok := StrToTime(s, formats, location); !ok {
		t.Errorf("failed, want to success for %q", s)
	}

	s = "2006/01/02 15:04:05"
	if _, ok := StrToTime(s, formats, location); !ok {
		t.Errorf("failed, want to success for %q", s)
	}

	s = "2006/01/02"
	if _, ok := StrToTime(s, formats, location); !ok {
		t.Errorf("failed, want to success for %q", s)
	}

	s = "2006/01/02 15:04:05 -08:00"
	if _, ok := StrToTime(s, formats, location); !ok {
		t.Errorf("failed, want to success for %q", s)
	}

	s = "2006/01/02 15:04:05 -0800"
	if _, ok := StrToTime(s, formats, location); !ok {
		t.Errorf("failed, want to success for %q", s)
	}

	s = "2006/01/02 15:04:05 -0800"
	if _, ok := StrToTime(s, formats, location); !ok {
		t.Errorf("failed, want to success for %q", s)
	}

	s = "2006/11/2 15:04:05 -0800"
	if _, ok := StrToTime(s, formats, location); !ok {
		t.Errorf("failed, want to success for %q", s)
	}

	s = "2006/01/02 15:04:05 PST"
	if _, ok := StrToTime(s, formats, location); !ok {
		t.Errorf("failed, want to success for %q", s)
	}

	s = "2006-1-2 15:04:05"
	if _, ok := StrToTime(s, formats, location); !ok {
		t.Errorf("failed, want to success for %q", s)
	}

	s = "2006-1-2"
	if _, ok := StrToTime(s, formats, location); !ok {
		t.Errorf("failed, want to success for %q", s)
	}

	s = "2006-1-2 15:04:05 -08:00"
	if _, ok := StrToTime(s, formats, location); !ok {
		t.Errorf("failed, want to success for %q", s)
	}

	s = "2006-1-2 15:04:05 -0800"
	if _, ok := StrToTime(s, formats, location); !ok {
		t.Errorf("failed, want to success for %q", s)
	}

	s = "2006-1-2 15:04:05 PST"
	if _, ok := StrToTime(s, formats, location); !ok {
		t.Errorf("failed, want to success for %q", s)
	}

	s = "2006/1/2 15:04:05"
	if _, ok := StrToTime(s, formats, location); !ok {
		t.Errorf("failed, want to success for %q", s)
	}

	s = "2006/1/2"
	if _, ok := StrToTime(s, formats, location); !ok {
		t.Errorf("failed, want to success for %q", s)
	}

	s = "2006/1/2 15:04:05 -08:00"
	if _, ok := StrToTime(s, formats, location); !ok {
		t.Errorf("failed, want to success for %q", s)
	}

	s = "2006/1/2 15:04:05 -0800"
	if _, ok := StrToTime(s, formats, location); !ok {
		t.Errorf("failed, want to success for %q", s)
	}

	s = "2006/1/2 15:04:05 PST"
	if _, ok := StrToTime(s, formats, location); !ok {
		t.Errorf("failed, want to success for %q", s)
	}

	s = "02 Jan 06 15:04 PDT"
	if _, ok := StrToTime(s, formats, location); !ok {
		t.Errorf("failed, want to success for %q", s)
	}

	s = "02 Jan 06 15:04 -0700"
	if _, ok := StrToTime(s, formats, location); !ok {
		t.Errorf("failed, want to success for %q", s)
	}

	s = "2006-01-02T15:04:05-08:00"
	if _, ok := StrToTime(s, formats, location); !ok {
		t.Errorf("failed, want to success for %q", s)
	}

	s = "2006-01-02T15:04:05+08:00"
	if _, ok := StrToTime(s, formats, location); !ok {
		t.Errorf("failed, want to success for %q", s)
	}

	s = "2006-01-02T15:04:05"
	if _, ok := StrToTime(s, formats, location); !ok {
		t.Errorf("failed, want to success for %q", s)
	}

	s = "e"
	if _, ok := StrToTime(s, formats, location); ok {
		t.Errorf("successeded, want to fail for %q", s)
	}

	s = "2006-01-02"
	tm, _ := StrToTime(s, formats, location)
	if tm.Location() != location {
		t.Errorf("location should be %q", location.String())
	}

	s = "2006-01-02T15:04:05+08:00"
	tm, _ = StrToTime(s, formats, location)
	z, i := tm.Zone()
	if z != "" {
		t.Errorf("zone name shoud be empty")
	}
	if i != 28800 {
		t.Errorf("zone offset should be %d", 28800)
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

func TestFloat64ToTime(t *testing.T) {
	location, _ := time.LoadLocation("UTC")

	f := float64(1136181845)
	expect := time.Date(2006, 1, 2, 6, 4, 5, 0, time.UTC).In(location)
	result := Float64ToTime(f, location)
	if !result.Equal(expect) {
		t.Errorf("result = %q, want %q for %f", result, expect, f)
	}

	f = 1136181845.123
	expect = time.Date(2006, 1, 2, 6, 4, 5, 123000000, time.UTC).In(location)
	result = Float64ToTime(f, location)
	if !result.Equal(expect) {
		t.Errorf("result = %q, want %q for %f", result, expect, f)
	}

	f = 1.123456789012
	expect = time.Date(1970, 1, 1, 0, 0, 1, 123456789, time.UTC).In(location)
	result = Float64ToTime(f, location)
	if !result.Equal(expect) {
		t.Errorf("result = %q, want %q for %f", result, expect, f)
	}
}

func TestParseFloat64(t *testing.T) {
	var p Primary
	var f float64

	f = 1.000
	p = ParseFloat64(f)
	if _, ok := p.(*Integer); !ok {
		t.Errorf("primary type = %T, want Integer for %f", p, f)
	}

	f = 1.234
	p = ParseFloat64(f)
	if _, ok := p.(*Float); !ok {
		t.Errorf("primary type = %T, want Float for %f", p, f)
	}
}

func TestToInteger(t *testing.T) {
	var p Primary
	var i Primary

	p = NewInteger(1)
	i = ToInteger(p)
	if _, ok := i.(*Integer); !ok {
		t.Errorf("primary type = %T, want Integer for %#v", i, p)
	}

	p = NewFloat(1)
	i = ToInteger(p)
	if _, ok := i.(*Integer); !ok {
		t.Errorf("primary type = %T, want Integer for %#v", i, p)
	}

	p = NewFloat(1.6)
	i = ToInteger(p)
	if _, ok := i.(*Null); !ok {
		t.Errorf("primary type = %T, want Null for %#v", i, p)
	}

	p = NewString(" 1")
	i = ToInteger(p)
	if _, ok := i.(*Integer); !ok {
		t.Errorf("primary type = %T, want Integer for %#v", i, p)
	}

	p = NewString("-1")
	i = ToInteger(p)
	if _, ok := i.(*Integer); !ok {
		t.Errorf("primary type = %T, want Integer for %#v", i, p)
	}

	p = NewString("1e+02")
	i = ToInteger(p)
	if _, ok := i.(*Integer); !ok {
		t.Errorf("primary type = %T, want Integer for %#v", i, p)
	}

	p = NewString("1.5")
	i = ToInteger(p)
	if _, ok := i.(*Null); !ok {
		t.Errorf("primary type = %T, want Null for %#v", i, p)
	}

	p = NewString("error")
	i = ToInteger(p)
	if _, ok := i.(*Null); !ok {
		t.Errorf("primary type = %T, want Null for %#v", i, p)
	}

	p = NewString("error")
	i = ToInteger(p)
	if _, ok := i.(*Null); !ok {
		t.Errorf("primary type = %T, want Null for %#v", i, p)
	}

	p = NewString("2002-02-02")
	i = ToInteger(p)
	if _, ok := i.(*Null); !ok {
		t.Errorf("primary type = %T, want Null for %#v", i, p)
	}

	p = NewString("2002/02/02")
	i = ToInteger(p)
	if _, ok := i.(*Null); !ok {
		t.Errorf("primary type = %T, want Null for %#v", i, p)
	}

	p = NewString("03 Mar 12 12:03 PST")
	i = ToInteger(p)
	if _, ok := i.(*Null); !ok {
		t.Errorf("primary type = %T, want Null for %#v", i, p)
	}

	p = NewString("")
	i = ToInteger(p)
	if _, ok := i.(*Null); !ok {
		t.Errorf("primary type = %T, want Null for %#v", i, p)
	}
}

func TestToFloat(t *testing.T) {
	var p Primary
	var f Primary

	p = NewInteger(1)
	f = ToFloat(p)
	if _, ok := f.(*Float); !ok {
		t.Errorf("primary type = %T, want Float for %#v", f, p)
	}

	p = NewFloat(1.234)
	f = ToFloat(p)
	if _, ok := f.(*Float); !ok {
		t.Errorf("primary type = %T, want Float for %#v", f, p)
	}

	p = NewString("1")
	f = ToFloat(p)
	if _, ok := f.(*Float); !ok {
		t.Errorf("primary type = %T, want Float for %#v", f, p)
	}

	p = NewString("error")
	f = ToFloat(p)
	if _, ok := f.(*Null); !ok {
		t.Errorf("primary type = %T, want Null for %#v", f, p)
	}
}

func TestToDatetime(t *testing.T) {
	var p Primary
	var dt Primary

	location, _ := time.LoadLocation("UTC")

	formats := []string{"01022006"}
	p = NewString("02012012")
	dt = ToDatetime(p, formats, location)
	if _, ok := dt.(*Datetime); !ok {
		t.Errorf("primary type = %T, want Datetime for %#v", dt, p)
	} else {
		expect := time.Date(2012, 2, 1, 0, 0, 0, 0, location)
		if !dt.(*Datetime).Raw().Equal(expect) {
			t.Errorf("datetime = %s, want %s for %#v", dt, expect, p)
		}
	}

	p = NewDatetimeFromString("2006-01-02 15:04:05", nil, location)
	dt = ToDatetime(p, nil, location)
	if _, ok := dt.(*Datetime); !ok {
		t.Errorf("primary type = %T, want Datetime for %#v", dt, p)
	}

	p = NewString("2006-01-02 15:04:05")
	dt = ToDatetime(p, nil, location)
	if _, ok := dt.(*Datetime); !ok {
		t.Errorf("primary type = %T, want Datetime for %#v", dt, p)
	}

	p = NewString("error")
	dt = ToDatetime(p, nil, location)
	if _, ok := dt.(*Null); !ok {
		t.Errorf("primary type = %T, want Null for %#v", dt, p)
	}
}

func TestToBoolean(t *testing.T) {
	var p Primary
	var b Primary

	p = NewBoolean(true)
	b = ToBoolean(p)
	if _, ok := b.(*Boolean); !ok {
		t.Errorf("primary type = %T, want Boolean for %#v", b, p)
	}

	p = NewTernary(ternary.TRUE)
	b = ToBoolean(p)
	if _, ok := b.(*Boolean); !ok {
		t.Errorf("primary type = %T, want Boolean for %#v", b, p)
	}

	p = NewInteger(1)
	b = ToBoolean(p)
	if _, ok := b.(*Boolean); !ok {
		t.Errorf("primary type = %T, want Boolean for %#v", b, p)
	}

	p = NewFloat(0)
	b = ToBoolean(p)
	if _, ok := b.(*Boolean); !ok {
		t.Errorf("primary type = %T, want Boolean for %#v", b, p)
	}

	p = NewString("true")
	b = ToBoolean(p)
	if _, ok := b.(*Boolean); !ok {
		t.Errorf("primary type = %T, want Boolean for %#v", b, p)
	}

	p = NewTernary(ternary.UNKNOWN)
	b = ToBoolean(p)
	if _, ok := b.(*Null); !ok {
		t.Errorf("primary type = %T, want Null for %#v", b, p)
	}

	p = NewString("error")
	b = ToBoolean(p)
	if _, ok := b.(*Null); !ok {
		t.Errorf("primary type = %T, want Null for %#v", b, p)
	}
}

func TestToString(t *testing.T) {
	var p Primary
	var s Primary

	location, _ := cmd.GetLocation("UTC")

	p = NewString("str")
	s = ToString(p)
	if _, ok := s.(*String); !ok {
		t.Errorf("primary type = %T, want String for %#v", s, p)
	}

	p = NewInteger(1)
	s = ToString(p)
	if _, ok := s.(*String); !ok {
		t.Errorf("primary type = %T, want String for %#v", s, p)
	}

	p = NewFloat(1)
	s = ToString(p)
	if _, ok := s.(*String); !ok {
		t.Errorf("primary type = %T, want String for %#v", s, p)
	}

	p = NewDatetimeFromString("2006-01-02 15:04:05", nil, location)
	s = ToString(p)
	if _, ok := s.(*Null); !ok {
		t.Errorf("primary type = %T, want Null for %#v", s, p)
	}
}

var maybeNumberTests = []struct {
	Input  string
	Expect bool
}{
	{
		Input:  "",
		Expect: false,
	},
	{
		Input:  "12345",
		Expect: true,
	},
	{
		Input:  "+1234567",
		Expect: true,
	},
	{
		Input:  "12.345",
		Expect: true,
	},
	{
		Input:  "1e+02",
		Expect: true,
	},
	{
		Input:  "1.123e-3",
		Expect: true,
	},
	{
		Input:  "12345.",
		Expect: false,
	},
	{
		Input:  "1e",
		Expect: false,
	},
	{
		Input:  "1ea",
		Expect: false,
	},
	{
		Input:  "abc",
		Expect: false,
	},
	{
		Input:  "+",
		Expect: false,
	},
}

func TestMaybeNumber(t *testing.T) {
	for _, v := range maybeNumberTests {
		result := MaybeNumber(v.Input)
		if result != v.Expect {
			t.Errorf("result = %t, want %t for %q", result, v.Expect, v.Input)
		}
	}
}

func BenchmarkStrToTime1(b *testing.B) {
	formats := []string{"01/02/2006"}
	location, _ := time.LoadLocation("UTC")

	for i := 0; i < b.N; i++ {
		s := "01/02/2006"
		_, _ = StrToTime(s, formats, location)
	}
}

func BenchmarkStrToTime2(b *testing.B) {
	formats := []string{"01/02/2006"}
	location, _ := time.LoadLocation("UTC")

	for i := 0; i < b.N; i++ {
		s := "2006-01-02T15:04:05-07:00"
		_, _ = StrToTime(s, formats, location)
	}
}

func BenchmarkStrToTime3(b *testing.B) {
	formats := []string{"01/02/2006"}
	location, _ := time.LoadLocation("UTC")

	for i := 0; i < b.N; i++ {
		s := "2006-01-02"
		_, _ = StrToTime(s, formats, location)
	}
}

func BenchmarkStrToTime4(b *testing.B) {
	formats := []string{"01/02/2006"}
	location, _ := time.LoadLocation("UTC")

	for i := 0; i < b.N; i++ {
		s := "2006-01-02 15:04:05"
		_, _ = StrToTime(s, formats, location)
	}
}

func BenchmarkStrToTime5(b *testing.B) {
	formats := []string{"01/02/2006"}
	location, _ := time.LoadLocation("UTC")

	for i := 0; i < b.N; i++ {
		s := "2006-01-02 15:04:05 -0700"
		_, _ = StrToTime(s, formats, location)
	}
}

func BenchmarkStrToTime6(b *testing.B) {
	formats := []string{"01/02/2006"}
	location, _ := cmd.GetLocation("UTC")

	for i := 0; i < b.N; i++ {
		s := "02 Jan 06 15:04 PDT"
		_, _ = StrToTime(s, formats, location)
	}
}

func BenchmarkStrToTime7(b *testing.B) {
	formats := []string{"01/02/2006"}
	location, _ := cmd.GetLocation("UTC")

	for i := 0; i < b.N; i++ {
		s := "abcdefghijklmnopq"
		_, _ = StrToTime(s, formats, location)
	}
}

func BenchmarkToInteger(b *testing.B) {
	p := NewString("a")
	for i := 0; i < b.N; i++ {
		_ = ToInteger(p)
	}
}

func BenchmarkToInteger2(b *testing.B) {
	p := NewString("2012-02-02")
	for i := 0; i < b.N; i++ {
		_ = ToInteger(p)
	}
}

func BenchmarkToInteger3(b *testing.B) {
	p := NewString(" 12345")
	for i := 0; i < b.N; i++ {
		_ = ToInteger(p)
	}
}

func BenchmarkToInteger4(b *testing.B) {
	p := NewString(" 123.456")
	for i := 0; i < b.N; i++ {
		_ = ToInteger(p)
	}
}

func BenchmarkToInteger5(b *testing.B) {
	p := NewFloat(123.456)
	for i := 0; i < b.N; i++ {
		_ = ToInteger(p)
	}
}

func BenchmarkToFloat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		p := NewString("a")
		ToFloat(p)
	}
}

func BenchmarkToFloat2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		p := NewString("2012-02-02")
		ToFloat(p)
	}
}

var convertDatetimeFormatBenchString = "%Y-%m-%d %H:%i:%s"

func BenchmarkConvertDatetimeFormat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = ConvertDatetimeFormat(convertDatetimeFormatBenchString)
	}
}
