package query

import (
	"bytes"
	"reflect"
	"testing"
	"time"

	"github.com/mithrandie/csvq/lib/value"

	"github.com/mithrandie/ternary"
)

func TestUintPool_Exists(t *testing.T) {
	pool := &UintPool{
		limitToUseSlice: 2,
		m: map[uint]bool{
			1: true,
			2: true,
			3: true,
		},
		values: []uint{1, 2, 3},
	}

	result := pool.Exists(1)
	if !result {
		t.Errorf("result = %t, want %t", false, true)
	}

	result = pool.Exists(9)
	if result {
		t.Errorf("result = %t, want %t", true, false)
	}

	pool = &UintPool{
		limitToUseSlice: 2,
		m: map[uint]bool{
			1: true,
		},
		values: []uint{1},
	}

	result = pool.Exists(1)
	if !result {
		t.Errorf("result = %t, want %t", false, true)
	}

	result = pool.Exists(9)
	if result {
		t.Errorf("result = %t, want %t", true, false)
	}
}

func TestUintPool_Add(t *testing.T) {
	pool := NewUintPool(1, 2)

	pool.Add(1)
	expect := &UintPool{
		limitToUseSlice: 2,
		m: map[uint]bool{
			1: true,
		},
		values: []uint{1},
	}

	if !reflect.DeepEqual(pool, expect) {
		t.Errorf("pool = %v, want %v", pool, expect)
	}
}

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
	expect := "[S]STR:[I]1:[I]0:[I]\x03\x00\x00\x00\x00\x00\x00\x00:[F]\x58\x39\xb4\xc8\x76\xbe\xf3\x3f:[D]\x00\xa6\x5b\x14\x42\x08\x6f\x12:[D]\xc0\x7a\xb0\x1b\x42\x08\x6f\x12:[D]\x15\x73\xb7\x1b\x42\x08\x6f\x12:[I]1:[I]0:[N]:[N]"

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

func benchmarkComparisonKeys(b *testing.B, plist []value.Primary) {
	buf := &bytes.Buffer{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		SerializeComparisonKeys(buf, plist, TestTx.Flags)
	}
}

func BenchmarkComparisonKeysString(b *testing.B) {
	plist := []value.Primary{
		value.NewString("abcdefghi"),
		value.NewString("jklmn"),
	}

	benchmarkComparisonKeys(b, plist)
}

func BenchmarkComparisonKeysInteger(b *testing.B) {
	plist := []value.Primary{
		value.NewInteger(123456789),
		value.NewInteger(1),
	}

	benchmarkComparisonKeys(b, plist)
}

func BenchmarkComparisonKeysFlaot(b *testing.B) {
	plist := []value.Primary{
		value.NewFloat(1.234e-9),
		value.NewFloat(123456.7890123),
	}

	benchmarkComparisonKeys(b, plist)
}

func BenchmarkComparisonKeysDatetime(b *testing.B) {
	plist := []value.Primary{
		value.NewDatetime(time.Date(2012, 2, 4, 9, 18, 15, 0, time.Local)),
		value.NewDatetime(time.Date(2015, 3, 4, 9, 17, 15, 0, time.Local)),
	}

	benchmarkComparisonKeys(b, plist)
}

func BenchmarkComparisonKeysBoolean(b *testing.B) {
	plist := []value.Primary{
		value.NewBoolean(true),
		value.NewBoolean(false),
	}

	benchmarkComparisonKeys(b, plist)
}

func BenchmarkComparisonKeysTernary(b *testing.B) {
	plist := []value.Primary{
		value.NewTernary(ternary.TRUE),
		value.NewTernary(ternary.FALSE),
		value.NewTernary(ternary.UNKNOWN),
	}

	benchmarkComparisonKeys(b, plist)
}

func BenchmarkComparisonKeysNull(b *testing.B) {
	plist := []value.Primary{
		value.NewNull(),
	}

	benchmarkComparisonKeys(b, plist)
}
