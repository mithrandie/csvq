package query

import (
	"testing"
	"time"

	"github.com/mithrandie/csvq/lib/value"

	"github.com/mithrandie/ternary"
)

func TestSerializeComparisonKeys(t *testing.T) {
	values := []value.Primary{
		value.NewString("str"),
		value.NewInteger(1),
		value.NewInteger(0),
		value.NewInteger(3),
		value.NewFloat(1.234),
		value.NewDatetimeFromString("2012-02-03T09:18:15-08:00"),
		value.NewDatetimeFromString("2012-02-03T09:18:15.123-08:00"),
		value.NewDatetimeFromString("2012-02-03T09:18:15.123456789-08:00"),
		value.NewBoolean(true),
		value.NewBoolean(false),
		value.NewTernary(ternary.UNKNOWN),
		value.NewNull(),
	}
	expect := "[S]STR:[I]1[B]true:[I]0[B]false:[I]3:[F]1.234:[I]1328289495:[F]1328289495.123:[D]1328289495123456789:[I]1[B]true:[I]0[B]false:[N]:[N]"

	result := SerializeComparisonKeys(values)
	if result != expect {
		t.Errorf("result = %q, want %q", result, expect)
	}
}

var formatStringTests = []struct {
	Name   string
	Format string
	Args   []value.Primary
	Result string
	Error  string
}{
	{
		Name:   "FormatString Integer",
		Format: "integer: %b %o %d %x %X %b %+d % d",
		Args: []value.Primary{
			value.NewInteger(123),
			value.NewInteger(123),
			value.NewInteger(123),
			value.NewInteger(123),
			value.NewInteger(123),
			value.NewNull(),
			value.NewInteger(123),
			value.NewInteger(123),
		},
		Result: "integer: 1111011 173 123 7b 7B  +123  123",
	},
	{
		Name:   "FormatString Float",
		Format: "float: %e %E %f %e %+f %.2f %.6f %.6f %.6e %.6e % f",
		Args: []value.Primary{
			value.NewFloat(0.0000000000123),
			value.NewFloat(0.0000000000123),
			value.NewFloat(0.0000000000123),
			value.NewNull(),
			value.NewFloat(123.456),
			value.NewFloat(123.456),
			value.NewFloat(123.456),
			value.NewFloat(0),
			value.NewFloat(0),
			value.NewFloat(1.23e-2),
			value.NewFloat(123.456),
		},
		Result: "float: 1.23e-11 1.23E-11 0.0000000000123  +123.456 123.46 123.456000 0.000000 0.000000e+00 1.230000e-02  123.456",
	},
	{
		Name:   "FormatString String",
		Format: "string: %s %s %s %s %s %s %s",
		Args: []value.Primary{
			value.NewString("str"),
			value.NewInteger(1),
			value.NewFloat(1.234),
			value.NewBoolean(true),
			value.NewTernary(ternary.TRUE),
			value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 0, GetTestLocation())),
			value.NewNull(),
		},
		Result: "string: str 1 1.234 true TRUE 2012-02-03T09:18:15Z NULL",
	},
	{
		Name:   "FormatString Quoted String",
		Format: "string: %q %q %q %q %q %q %q",
		Args: []value.Primary{
			value.NewString("st'r"),
			value.NewInteger(1),
			value.NewFloat(1.234),
			value.NewBoolean(true),
			value.NewTernary(ternary.TRUE),
			value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 0, GetTestLocation())),
			value.NewNull(),
		},
		Result: "string: \"st\\'r\" \"1\" \"1.234\" \"true\" \"TRUE\" \"2012-02-03T09:18:15Z\" \"NULL\"",
	},
	{
		Name:   "FormatString Quoted Identifier",
		Format: "string: %i %i %i %i %i %i %i",
		Args: []value.Primary{
			value.NewString("st`r"),
			value.NewInteger(1),
			value.NewFloat(1.234),
			value.NewBoolean(true),
			value.NewTernary(ternary.TRUE),
			value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 0, GetTestLocation())),
			value.NewNull(),
		},
		Result: "string: `st\\`r` `1` `1.234` `true` `TRUE` `2012-02-03T09:18:15Z` `NULL`",
	},
	{
		Name:   "FormatString Type",
		Format: "type: %T",
		Args: []value.Primary{
			value.NewString("str"),
		},
		Result: "type: String",
	},
	{
		Name:   "FormatString Padding",
		Format: "padding: %6d %+06d %+06d %2d %010d % 4d %6f %6s %-6s",
		Args: []value.Primary{
			value.NewInteger(123),
			value.NewInteger(123),
			value.NewInteger(-123),
			value.NewInteger(123),
			value.NewInteger(123),
			value.NewInteger(-1),
			value.NewFloat(123.4),
			value.NewString("str"),
			value.NewString("str"),
		},
		Result: "padding:    123 +00123 -00123 123 0000000123   -1  123.4    str str   ",
	},
	{
		Name:   "FormatString Etc.",
		Format: "string: %s %% %a %",
		Args: []value.Primary{
			value.NewString("str"),
		},
		Result: "string: str % %a %",
	},
	{
		Name:   "FormatString PlaceHolder Too Little Error",
		Format: "string: %s %s",
		Args: []value.Primary{
			value.NewString("str"),
		},
		Error: "[L:- C:-] number of replace values does not match",
	},
	{
		Name:   "FormatString PlaceHolder Too Many Error",
		Format: "string: %s %s",
		Args: []value.Primary{
			value.NewString("str"),
			value.NewString("str"),
			value.NewString("str"),
		},
		Error: "[L:- C:-] number of replace values does not match",
	},
}

func TestFormatString(t *testing.T) {
	for _, v := range formatStringTests {
		result, err := FormatString(v.Format, v.Args)
		if err != nil {
			if len(v.Error) < 1 {
				t.Errorf("%s: unexpected error %q", v.Name, err)
			} else if err.Error() != v.Error {
				t.Errorf("%s: error %q, want error %q", v.Name, err.Error(), v.Error)
			}
			continue
		}
		if 0 < len(v.Error) {
			t.Errorf("%s: no error, want error %q", v.Name, v.Error)
			continue
		}
		if result != v.Result {
			t.Errorf("%s: result = %q, want %q", v.Name, result, v.Result)
		}
	}
}

func BenchmarkDistinguish(b *testing.B) {
	values := make([]value.Primary, 10000)
	for i := 0; i < 100; i++ {
		for j := 0; j < 100; j++ {
			values[i*100+j] = value.NewInteger(int64(j))
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Distinguish(values)
	}
}
