package query

import (
	"bytes"
	"testing"

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
		value.NewDatetimeFromString("2012-02-03T09:18:15-08:00", TestTx.Flags.DatetimeFormat),
		value.NewDatetimeFromString("2012-02-03T09:18:15.123-08:00", TestTx.Flags.DatetimeFormat),
		value.NewDatetimeFromString("2012-02-03T09:18:15.123456789-08:00", TestTx.Flags.DatetimeFormat),
		value.NewBoolean(true),
		value.NewBoolean(false),
		value.NewTernary(ternary.UNKNOWN),
		value.NewNull(),
	}
	expect := "[S]STR:[I]1[B]true:[I]0[B]false:[I]3:[F]1.234:[I]1328289495:[F]1328289495.123:[D]1328289495123456789:[I]1[B]true:[I]0[B]false:[N]:[N]"

	buf := &bytes.Buffer{}
	SerializeComparisonKeys(buf, values, TestTx.Flags)
	result := buf.String()
	if result != expect {
		t.Errorf("result = %q, want %q", result, expect)
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
		Distinguish(values, TestTx.Flags)
	}
}

func BenchmarkSerializeKey(b *testing.B) {
	plist := []value.Primary{
		value.NewString("str"),
		value.NewInteger(123),
	}

	buf := &bytes.Buffer{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		SerializeComparisonKeys(buf, plist, TestTx.Flags)
	}
}
