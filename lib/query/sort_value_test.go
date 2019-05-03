package query

import (
	"bytes"
	"testing"
	"time"

	"github.com/mithrandie/csvq/lib/value"

	"github.com/mithrandie/ternary"
)

func TestSortValues_Serialize(t *testing.T) {
	values := SortValues{
		NewSortValue(value.NewNull(), TestTx.Flags),
		NewSortValue(value.NewInteger(1), TestTx.Flags),
		NewSortValue(value.NewFloat(1.234), TestTx.Flags),
		NewSortValue(value.NewDatetimeFromString("2012-02-03T09:18:15-08:00", TestTx.Flags.DatetimeFormat), TestTx.Flags),
		NewSortValue(value.NewDatetimeFromString("2012-02-03T09:18:15.123-08:00", TestTx.Flags.DatetimeFormat), TestTx.Flags),
		NewSortValue(value.NewDatetimeFromString("2012-02-03T09:18:15.123456789-08:00", TestTx.Flags.DatetimeFormat), TestTx.Flags),
		NewSortValue(value.NewBoolean(false), TestTx.Flags),
		NewSortValue(value.NewString("str"), TestTx.Flags),
	}
	expect := "[N]:[I]1:[F]\x58\x39\xb4\xc8\x76\xbe\xf3\x3f:[D]\x00\xa6\x5b\x14\x42\x08\x6f\x12:[D]\xc0\x7a\xb0\x1b\x42\x08\x6f\x12:[D]\x15\x73\xb7\x1b\x42\x08\x6f\x12:[I]0:[S]STR"

	buf := &bytes.Buffer{}
	values.Serialize(buf)
	result := buf.String()
	if result != expect {
		t.Errorf("result = %q, want %q", result, expect)
	}
}

var sortValueLessTests = []struct {
	Name         string
	SortValue    *SortValue
	CompareValue *SortValue
	Result       ternary.Value
}{
	{
		Name:         "SortValue Less Integer",
		SortValue:    NewSortValue(value.NewInteger(3), TestTx.Flags),
		CompareValue: NewSortValue(value.NewInteger(5), TestTx.Flags),
		Result:       ternary.TRUE,
	},
	{
		Name:         "SortValue Less Integer Equal",
		SortValue:    NewSortValue(value.NewInteger(3), TestTx.Flags),
		CompareValue: NewSortValue(value.NewInteger(3), TestTx.Flags),
		Result:       ternary.UNKNOWN,
	},
	{
		Name:         "SortValue Less Integer and Float",
		SortValue:    NewSortValue(value.NewInteger(3), TestTx.Flags),
		CompareValue: NewSortValue(value.NewFloat(5.4), TestTx.Flags),
		Result:       ternary.TRUE,
	},
	{
		Name:         "SortValue Less Integer and Datetime",
		SortValue:    NewSortValue(value.NewInteger(3), TestTx.Flags),
		CompareValue: NewSortValue(value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456789, GetTestLocation())), TestTx.Flags),
		Result:       ternary.UNKNOWN,
	},
	{
		Name:         "SortValue Less Integer and String",
		SortValue:    NewSortValue(value.NewInteger(3), TestTx.Flags),
		CompareValue: NewSortValue(value.NewString("4a"), TestTx.Flags),
		Result:       ternary.TRUE,
	},
	{
		Name:         "SortValue Less Float",
		SortValue:    NewSortValue(value.NewFloat(3.4), TestTx.Flags),
		CompareValue: NewSortValue(value.NewFloat(5.1), TestTx.Flags),
		Result:       ternary.TRUE,
	},
	{
		Name:         "SortValue Less Float Equal",
		SortValue:    NewSortValue(value.NewFloat(3.4), TestTx.Flags),
		CompareValue: NewSortValue(value.NewFloat(3.4), TestTx.Flags),
		Result:       ternary.UNKNOWN,
	},
	{
		Name:         "SortValue Less Float and Datetime",
		SortValue:    NewSortValue(value.NewFloat(3.4), TestTx.Flags),
		CompareValue: NewSortValue(value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456789, GetTestLocation())), TestTx.Flags),
		Result:       ternary.UNKNOWN,
	},
	{
		Name:         "SortValue Less Float and String",
		SortValue:    NewSortValue(value.NewFloat(3.4), TestTx.Flags),
		CompareValue: NewSortValue(value.NewString("4a"), TestTx.Flags),
		Result:       ternary.TRUE,
	},
	{
		Name:         "SortValue Less Datetime",
		SortValue:    NewSortValue(value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456789, GetTestLocation())), TestTx.Flags),
		CompareValue: NewSortValue(value.NewDatetime(time.Date(2012, 2, 4, 9, 18, 15, 123456789, GetTestLocation())), TestTx.Flags),
		Result:       ternary.TRUE,
	},
	{
		Name:         "SortValue Less Datetime Equal",
		SortValue:    NewSortValue(value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456789, GetTestLocation())), TestTx.Flags),
		CompareValue: NewSortValue(value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456789, GetTestLocation())), TestTx.Flags),
		Result:       ternary.UNKNOWN,
	},
	{
		Name:         "SortValue Less String",
		SortValue:    NewSortValue(value.NewString("aaa"), TestTx.Flags),
		CompareValue: NewSortValue(value.NewString("abc"), TestTx.Flags),
		Result:       ternary.TRUE,
	},
	{
		Name:         "SortValue Less String Equal",
		SortValue:    NewSortValue(value.NewString(" aaa "), TestTx.Flags),
		CompareValue: NewSortValue(value.NewString("AAA"), TestTx.Flags),
		Result:       ternary.UNKNOWN,
	},
	{
		Name:         "SortValue Less Boolean",
		SortValue:    NewSortValue(value.NewBoolean(true), TestTx.Flags),
		CompareValue: NewSortValue(value.NewTernary(ternary.FALSE), TestTx.Flags),
		Result:       ternary.UNKNOWN,
	},
	{
		Name:         "SortValue Less Incommensurable Types",
		SortValue:    NewSortValue(value.NewInteger(3), TestTx.Flags),
		CompareValue: NewSortValue(value.NewTernary(ternary.UNKNOWN), TestTx.Flags),
		Result:       ternary.UNKNOWN,
	},
}

func TestSortValue_Less(t *testing.T) {
	for _, v := range sortValueLessTests {
		result := v.SortValue.Less(v.CompareValue)
		if result != v.Result {
			t.Errorf("%s: result = %s, want %s", v.Name, result, v.Result)
		}
	}
}

var sortValueEquivalentToTests = []struct {
	Name         string
	SortValue    *SortValue
	CompareValue *SortValue
	Result       bool
}{
	{
		Name:         "SortValue EquivalentTo Integer",
		SortValue:    NewSortValue(value.NewInteger(3), TestTx.Flags),
		CompareValue: NewSortValue(value.NewInteger(3), TestTx.Flags),
		Result:       true,
	},
	{
		Name:         "SortValue EquivalentTo Integer and Boolean",
		SortValue:    NewSortValue(value.NewInteger(1), TestTx.Flags),
		CompareValue: NewSortValue(value.NewBoolean(true), TestTx.Flags),
		Result:       true,
	},
	{
		Name:         "SortValue EquivalentTo Integer and DateTime",
		SortValue:    NewSortValue(value.NewInteger(1328260695), TestTx.Flags),
		CompareValue: NewSortValue(value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 0, GetTestLocation())), TestTx.Flags),
		Result:       false,
	},
	{
		Name:         "SortValue EquivalentTo Float",
		SortValue:    NewSortValue(value.NewFloat(3.21), TestTx.Flags),
		CompareValue: NewSortValue(value.NewFloat(3.21), TestTx.Flags),
		Result:       true,
	},
	{
		Name:         "SortValue EquivalentTo Float and DateTime",
		SortValue:    NewSortValue(value.NewFloat(1328260695.0001), TestTx.Flags),
		CompareValue: NewSortValue(value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 100000, GetTestLocation())), TestTx.Flags),
		Result:       false,
	},
	{
		Name:         "SortValue EquivalentTo Datetime",
		SortValue:    NewSortValue(value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456789, GetTestLocation())), TestTx.Flags),
		CompareValue: NewSortValue(value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456789, GetTestLocation())), TestTx.Flags),
		Result:       true,
	},
	{
		Name:         "SortValue EquivalentTo Boolean",
		SortValue:    NewSortValue(value.NewBoolean(true), TestTx.Flags),
		CompareValue: NewSortValue(value.NewBoolean(true), TestTx.Flags),
		Result:       true,
	},
	{
		Name:         "SortValue EquivalentTo Boolean and Integer",
		SortValue:    NewSortValue(value.NewBoolean(true), TestTx.Flags),
		CompareValue: NewSortValue(value.NewInteger(1), TestTx.Flags),
		Result:       true,
	},
	{
		Name:         "SortValue EquivalentTo String",
		SortValue:    NewSortValue(value.NewString("str"), TestTx.Flags),
		CompareValue: NewSortValue(value.NewString("str"), TestTx.Flags),
		Result:       true,
	},
	{
		Name:         "SortValue EquivalentTo Null",
		SortValue:    NewSortValue(value.NewNull(), TestTx.Flags),
		CompareValue: NewSortValue(value.NewNull(), TestTx.Flags),
		Result:       true,
	},
	{
		Name:         "SortValue EquivalentTo String and Null",
		SortValue:    NewSortValue(value.NewString("str"), TestTx.Flags),
		CompareValue: NewSortValue(value.NewNull(), TestTx.Flags),
		Result:       false,
	},
}

func TestSortValue_EquivalentTo(t *testing.T) {
	for _, v := range sortValueEquivalentToTests {
		result := v.SortValue.EquivalentTo(v.CompareValue)
		if result != v.Result {
			t.Errorf("%s: result = %t, want %t", v.Name, result, v.Result)
		}
	}
}

var sortValueLessBench1 = NewSortValue(value.NewInteger(12345), TestTx.Flags)
var sortValueLessBench2 = NewSortValue(value.NewInteger(67890), TestTx.Flags)

func BenchmarkSortValue_Less(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = sortValueLessBench1.Less(sortValueLessBench2)
	}
}

var sortValuesEquivalentBench1 = SortValues{
	NewSortValue(value.NewInteger(12345), TestTx.Flags),
	NewSortValue(value.NewString("abcdefghijklmnopqrstuvwxymabcdefghijklmnopqrstuvwxyz"), TestTx.Flags),
	NewSortValue(value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 0, GetTestLocation())), TestTx.Flags),
}

var sortValuesEquivalentBench2 = SortValues{
	NewSortValue(value.NewInteger(12345), TestTx.Flags),
	NewSortValue(value.NewString("abcdefghijklmnopqrstuvwxymabcdefghijklmnopqrstuvwxyz"), TestTx.Flags),
	NewSortValue(value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 0, GetTestLocation())), TestTx.Flags),
}

func BenchmarkSortValues_EquivalentTo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = sortValuesEquivalentBench1.EquivalentTo(sortValuesEquivalentBench2)
	}
}
