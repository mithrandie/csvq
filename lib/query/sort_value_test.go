package query

import (
	"testing"
	"time"

	"github.com/mithrandie/csvq/lib/ternary"
	"github.com/mithrandie/csvq/lib/value"
)

func TestSortValues_Serialize(t *testing.T) {
	values := SortValues{
		NewSortValue(value.NewNull()),
		NewSortValue(value.NewInteger(1)),
		NewSortValue(value.NewFloat(1.234)),
		NewSortValue(value.NewDatetimeFromString("2012-02-03T09:18:15-08:00")),
		NewSortValue(value.NewDatetimeFromString("2012-02-03T09:18:15.123-08:00")),
		NewSortValue(value.NewDatetimeFromString("2012-02-03T09:18:15.123456789-08:00")),
		NewSortValue(value.NewBoolean(false)),
		NewSortValue(value.NewString("str")),
	}
	expect := "[N]:[I]1[B]true:[F]1.234:[I]1328289495:[F]1328289495.123:[D]1328289495123456789:[I]0[B]false:[S]STR"

	result := values.Serialize()
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
		SortValue:    NewSortValue(value.NewInteger(3)),
		CompareValue: NewSortValue(value.NewInteger(5)),
		Result:       ternary.TRUE,
	},
	{
		Name:         "SortValue Less Integer Equal",
		SortValue:    NewSortValue(value.NewInteger(3)),
		CompareValue: NewSortValue(value.NewInteger(3)),
		Result:       ternary.UNKNOWN,
	},
	{
		Name:         "SortValue Less Integer and Float",
		SortValue:    NewSortValue(value.NewInteger(3)),
		CompareValue: NewSortValue(value.NewFloat(5.4)),
		Result:       ternary.TRUE,
	},
	{
		Name:         "SortValue Less Integer and Datetime",
		SortValue:    NewSortValue(value.NewInteger(3)),
		CompareValue: NewSortValue(value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456789, GetTestLocation()))),
		Result:       ternary.TRUE,
	},
	{
		Name:         "SortValue Less Integer and String",
		SortValue:    NewSortValue(value.NewInteger(3)),
		CompareValue: NewSortValue(value.NewString("4a")),
		Result:       ternary.TRUE,
	},
	{
		Name:         "SortValue Less Float",
		SortValue:    NewSortValue(value.NewFloat(3.4)),
		CompareValue: NewSortValue(value.NewFloat(5.1)),
		Result:       ternary.TRUE,
	},
	{
		Name:         "SortValue Less Float Equal",
		SortValue:    NewSortValue(value.NewFloat(3.4)),
		CompareValue: NewSortValue(value.NewFloat(3.4)),
		Result:       ternary.UNKNOWN,
	},
	{
		Name:         "SortValue Less Float and Datetime",
		SortValue:    NewSortValue(value.NewFloat(3.4)),
		CompareValue: NewSortValue(value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456789, GetTestLocation()))),
		Result:       ternary.TRUE,
	},
	{
		Name:         "SortValue Less Float and String",
		SortValue:    NewSortValue(value.NewFloat(3.4)),
		CompareValue: NewSortValue(value.NewString("4a")),
		Result:       ternary.TRUE,
	},
	{
		Name:         "SortValue Less Datetime",
		SortValue:    NewSortValue(value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456789, GetTestLocation()))),
		CompareValue: NewSortValue(value.NewDatetime(time.Date(2012, 2, 4, 9, 18, 15, 123456789, GetTestLocation()))),
		Result:       ternary.TRUE,
	},
	{
		Name:         "SortValue Less Datetime Equal",
		SortValue:    NewSortValue(value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456789, GetTestLocation()))),
		CompareValue: NewSortValue(value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456789, GetTestLocation()))),
		Result:       ternary.UNKNOWN,
	},
	{
		Name:         "SortValue Less String",
		SortValue:    NewSortValue(value.NewString("aaa")),
		CompareValue: NewSortValue(value.NewString("abc")),
		Result:       ternary.TRUE,
	},
	{
		Name:         "SortValue Less String Equal",
		SortValue:    NewSortValue(value.NewString(" aaa ")),
		CompareValue: NewSortValue(value.NewString("AAA")),
		Result:       ternary.UNKNOWN,
	},
	{
		Name:         "SortValue Less Boolean",
		SortValue:    NewSortValue(value.NewBoolean(true)),
		CompareValue: NewSortValue(value.NewTernary(ternary.FALSE)),
		Result:       ternary.UNKNOWN,
	},
	{
		Name:         "SortValue Less Incommensurable Types",
		SortValue:    NewSortValue(value.NewInteger(3)),
		CompareValue: NewSortValue(value.NewTernary(ternary.UNKNOWN)),
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
		SortValue:    NewSortValue(value.NewInteger(3)),
		CompareValue: NewSortValue(value.NewInteger(3)),
		Result:       true,
	},
	{
		Name:         "SortValue EquivalentTo Integer and Boolean",
		SortValue:    NewSortValue(value.NewInteger(1)),
		CompareValue: NewSortValue(value.NewBoolean(true)),
		Result:       true,
	},
	{
		Name:         "SortValue EquivalentTo Integer and DateTime",
		SortValue:    NewSortValue(value.NewInteger(1328289495)),
		CompareValue: NewSortValue(value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 0, GetTestLocation()))),
		Result:       true,
	},
	{
		Name:         "SortValue EquivalentTo Float",
		SortValue:    NewSortValue(value.NewFloat(3.21)),
		CompareValue: NewSortValue(value.NewFloat(3.21)),
		Result:       true,
	},
	{
		Name:         "SortValue EquivalentTo Float and DateTime",
		SortValue:    NewSortValue(value.NewFloat(1328289495.0001)),
		CompareValue: NewSortValue(value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 100000, GetTestLocation()))),
		Result:       true,
	},
	{
		Name:         "SortValue EquivalentTo Datetime",
		SortValue:    NewSortValue(value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456789, GetTestLocation()))),
		CompareValue: NewSortValue(value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456789, GetTestLocation()))),
		Result:       true,
	},
	{
		Name:         "SortValue EquivalentTo Boolean",
		SortValue:    NewSortValue(value.NewBoolean(true)),
		CompareValue: NewSortValue(value.NewBoolean(true)),
		Result:       true,
	},
	{
		Name:         "SortValue EquivalentTo Boolean and Integer",
		SortValue:    NewSortValue(value.NewBoolean(true)),
		CompareValue: NewSortValue(value.NewInteger(1)),
		Result:       true,
	},
	{
		Name:         "SortValue EquivalentTo String",
		SortValue:    NewSortValue(value.NewString("str")),
		CompareValue: NewSortValue(value.NewString("str")),
		Result:       true,
	},
	{
		Name:         "SortValue EquivalentTo Null",
		SortValue:    NewSortValue(value.NewNull()),
		CompareValue: NewSortValue(value.NewNull()),
		Result:       true,
	},
	{
		Name:         "SortValue EquivalentTo String and Null",
		SortValue:    NewSortValue(value.NewString("str")),
		CompareValue: NewSortValue(value.NewNull()),
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

var sortValueLessBench1 = NewSortValue(value.NewInteger(12345))
var sortValueLessBench2 = NewSortValue(value.NewInteger(67890))

func BenchmarkSortValue_Less(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = sortValueLessBench1.Less(sortValueLessBench2)
	}
}

var sortValuesEquivalentBench1 = SortValues{
	NewSortValue(value.NewInteger(12345)),
	NewSortValue(value.NewString("abcdefghijklmnopqrstuvwxymabcdefghijklmnopqrstuvwxyz")),
	NewSortValue(value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 0, GetTestLocation()))),
}

var sortValuesEquivalentBench2 = SortValues{
	NewSortValue(value.NewInteger(12345)),
	NewSortValue(value.NewString("abcdefghijklmnopqrstuvwxymabcdefghijklmnopqrstuvwxyz")),
	NewSortValue(value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 0, GetTestLocation()))),
}

func BenchmarkSortValues_EquivalentTo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = sortValuesEquivalentBench1.EquivalentTo(sortValuesEquivalentBench2)
	}
}
