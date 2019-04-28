package query

import (
	"testing"
	"time"

	"github.com/mithrandie/csvq/lib/value"
	"github.com/mithrandie/ternary"
)

var stringFormatterFormatTests = []struct {
	Format string
	Values []value.Primary
	Expect string
	Error  string
}{
	{
		Format: "--%d--%+6d--% d--%06d--%-10d--%4d--%04d--%d--",
		Values: []value.Primary{
			value.NewInteger(123),
			value.NewInteger(123),
			value.NewInteger(123),
			value.NewInteger(123),
			value.NewInteger(123),
			value.NewNull(),
			value.NewNull(),
			value.NewNull(),
		},
		Expect: "--123--  +123-- 123--000123--123       --    --    ----",
	},
	{
		Format: "--%b--%o--%x--%X--% x--",
		Values: []value.Primary{
			value.NewInteger(123),
			value.NewInteger(123),
			value.NewInteger(123),
			value.NewInteger(123),
			value.NewInteger(123),
		},
		Expect: "--1111011--173--7b--7B-- 7b--",
	},
	{
		Format: "--%e--%E--%f--%.2f--%.f--%06f--%02f--%0f--%-6f--%6f--%2f--% f--%4f--%04f--%f--",
		Values: []value.Primary{
			value.NewFloat(123.456),
			value.NewFloat(123.456),
			value.NewFloat(123.456),
			value.NewFloat(123.456),
			value.NewFloat(123.456),
			value.NewFloat(-1.2),
			value.NewFloat(1.2),
			value.NewFloat(-1.2),
			value.NewFloat(1.2),
			value.NewFloat(1.2),
			value.NewFloat(1.2),
			value.NewFloat(1.2),
			value.NewNull(),
			value.NewNull(),
			value.NewNull(),
		},
		Expect: "--1.23456e+02--1.23456E+02--123.456--123.46--123---001.2--1.2---1.2--1.2   --   1.2--1.2-- 1.2--    --    ----",
	},
	{
		Format: "--%s--%q--%i--%T--%.2i--%.2T--%%",
		Values: []value.Primary{
			value.NewString("str"),
			value.NewString("str"),
			value.NewString("str"),
			value.NewString("str"),
			value.NewString("str"),
			value.NewString("str"),
		},
		Expect: "--str--'str'--`str`--String--`st`--St--%",
	},
	{
		Format: "--%s--%s--%s--",
		Values: []value.Primary{
			value.NewString("str"),
			value.NewTernary(ternary.TRUE),
			value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 0, GetTestLocation())),
		},
		Expect: "--str--TRUE--2012-02-03T09:18:15Z--",
	},
	{
		Format: "--%s--",
		Values: []value.Primary{},
		Error:  "number of replace values does not match",
	},
	{
		Format: "--%s--",
		Values: []value.Primary{
			value.NewString("str"),
			value.NewString("str"),
		},
		Error: "number of replace values does not match",
	},
	{
		Format: "--%w--",
		Values: []value.Primary{
			value.NewString("str"),
		},
		Error: "\"w\" is an unknown placeholder",
	},
	{
		Format: "--%0",
		Values: []value.Primary{
			value.NewString("str"),
		},
		Error: "unexpected termination of format string",
	},
}

func TestStringFormatter_Format(t *testing.T) {
	f := NewStringFormatter()

	for _, v := range stringFormatterFormatTests {
		result, err := f.Format(v.Format, v.Values)
		if err != nil {
			if len(v.Error) < 1 {
				t.Errorf("unexpected error %q for %q, %v", err.Error(), v.Format, v.Values)
			} else if err.Error() != v.Error {
				t.Errorf("error %q, want error %q for %q, %v", err, v.Error, v.Format, v.Values)
			}
			continue
		}
		if 0 < len(v.Error) {
			t.Errorf("no error, want error %q for %q, %v", v.Error, v.Format, v.Values)
			continue
		}
		if result != v.Expect {
			t.Errorf("result = %q, want %q for %q, %v", result, v.Expect, v.Format, v.Values)
		}
	}
}
