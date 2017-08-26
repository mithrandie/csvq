package query

import (
	"testing"
	"time"

	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/ternary"
)

func TestSortValues_Serialize(t *testing.T) {
	values := SortValues{
		NewSortValue(parser.NewNull()),
		NewSortValue(parser.NewInteger(1)),
		NewSortValue(parser.NewFloat(1.234)),
		NewSortValue(parser.NewDatetimeFromString("2012-02-03T09:18:15-08:00")),
		NewSortValue(parser.NewDatetimeFromString("2012-02-03T09:18:15.123-08:00")),
		NewSortValue(parser.NewDatetimeFromString("2012-02-03T09:18:15.123456789-08:00")),
		NewSortValue(parser.NewBoolean(false)),
		NewSortValue(parser.NewString("str")),
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
		SortValue:    NewSortValue(parser.NewInteger(3)),
		CompareValue: NewSortValue(parser.NewInteger(5)),
		Result:       ternary.TRUE,
	},
	{
		Name:         "SortValue Less Integer Equal",
		SortValue:    NewSortValue(parser.NewInteger(3)),
		CompareValue: NewSortValue(parser.NewInteger(3)),
		Result:       ternary.UNKNOWN,
	},
	{
		Name:         "SortValue Less Integer and Float",
		SortValue:    NewSortValue(parser.NewInteger(3)),
		CompareValue: NewSortValue(parser.NewFloat(5.4)),
		Result:       ternary.TRUE,
	},
	{
		Name:         "SortValue Less Integer and Datetime",
		SortValue:    NewSortValue(parser.NewInteger(3)),
		CompareValue: NewSortValue(parser.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456789, GetTestLocation()))),
		Result:       ternary.TRUE,
	},
	{
		Name:         "SortValue Less Integer and String",
		SortValue:    NewSortValue(parser.NewInteger(3)),
		CompareValue: NewSortValue(parser.NewString("4a")),
		Result:       ternary.TRUE,
	},
	{
		Name:         "SortValue Less Float",
		SortValue:    NewSortValue(parser.NewFloat(3.4)),
		CompareValue: NewSortValue(parser.NewFloat(5.1)),
		Result:       ternary.TRUE,
	},
	{
		Name:         "SortValue Less Float Equal",
		SortValue:    NewSortValue(parser.NewFloat(3.4)),
		CompareValue: NewSortValue(parser.NewFloat(3.4)),
		Result:       ternary.UNKNOWN,
	},
	{
		Name:         "SortValue Less Float and Datetime",
		SortValue:    NewSortValue(parser.NewFloat(3.4)),
		CompareValue: NewSortValue(parser.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456789, GetTestLocation()))),
		Result:       ternary.TRUE,
	},
	{
		Name:         "SortValue Less Float and String",
		SortValue:    NewSortValue(parser.NewFloat(3.4)),
		CompareValue: NewSortValue(parser.NewString("4a")),
		Result:       ternary.TRUE,
	},
	{
		Name:         "SortValue Less Datetime",
		SortValue:    NewSortValue(parser.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456789, GetTestLocation()))),
		CompareValue: NewSortValue(parser.NewDatetime(time.Date(2012, 2, 4, 9, 18, 15, 123456789, GetTestLocation()))),
		Result:       ternary.TRUE,
	},
	{
		Name:         "SortValue Less Datetime Equal",
		SortValue:    NewSortValue(parser.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456789, GetTestLocation()))),
		CompareValue: NewSortValue(parser.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456789, GetTestLocation()))),
		Result:       ternary.UNKNOWN,
	},
	{
		Name:         "SortValue Less String",
		SortValue:    NewSortValue(parser.NewString("aaa")),
		CompareValue: NewSortValue(parser.NewString("abc")),
		Result:       ternary.TRUE,
	},
	{
		Name:         "SortValue Less String Equal",
		SortValue:    NewSortValue(parser.NewString(" aaa ")),
		CompareValue: NewSortValue(parser.NewString("AAA")),
		Result:       ternary.UNKNOWN,
	},
	{
		Name:         "SortValue Less Boolean",
		SortValue:    NewSortValue(parser.NewBoolean(true)),
		CompareValue: NewSortValue(parser.NewTernary(ternary.FALSE)),
		Result:       ternary.UNKNOWN,
	},
	{
		Name:         "SortValue Less Incommensurable Types",
		SortValue:    NewSortValue(parser.NewInteger(3)),
		CompareValue: NewSortValue(parser.NewTernary(ternary.UNKNOWN)),
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
		SortValue:    NewSortValue(parser.NewInteger(3)),
		CompareValue: NewSortValue(parser.NewInteger(3)),
		Result:       true,
	},
	{
		Name:         "SortValue EquivalentTo Integer and Boolean",
		SortValue:    NewSortValue(parser.NewInteger(1)),
		CompareValue: NewSortValue(parser.NewBoolean(true)),
		Result:       true,
	},
	{
		Name:         "SortValue EquivalentTo Integer and DateTime",
		SortValue:    NewSortValue(parser.NewInteger(1328289495)),
		CompareValue: NewSortValue(parser.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 0, GetTestLocation()))),
		Result:       true,
	},
	{
		Name:         "SortValue EquivalentTo Float",
		SortValue:    NewSortValue(parser.NewFloat(3.21)),
		CompareValue: NewSortValue(parser.NewFloat(3.21)),
		Result:       true,
	},
	{
		Name:         "SortValue EquivalentTo Float and DateTime",
		SortValue:    NewSortValue(parser.NewFloat(1328289495.0001)),
		CompareValue: NewSortValue(parser.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 100000, GetTestLocation()))),
		Result:       true,
	},
	{
		Name:         "SortValue EquivalentTo Datetime",
		SortValue:    NewSortValue(parser.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456789, GetTestLocation()))),
		CompareValue: NewSortValue(parser.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456789, GetTestLocation()))),
		Result:       true,
	},
	{
		Name:         "SortValue EquivalentTo Boolean",
		SortValue:    NewSortValue(parser.NewBoolean(true)),
		CompareValue: NewSortValue(parser.NewBoolean(true)),
		Result:       true,
	},
	{
		Name:         "SortValue EquivalentTo Boolean and Integer",
		SortValue:    NewSortValue(parser.NewBoolean(true)),
		CompareValue: NewSortValue(parser.NewInteger(1)),
		Result:       true,
	},
	{
		Name:         "SortValue EquivalentTo String",
		SortValue:    NewSortValue(parser.NewString("str")),
		CompareValue: NewSortValue(parser.NewString("str")),
		Result:       true,
	},
	{
		Name:         "SortValue EquivalentTo Null",
		SortValue:    NewSortValue(parser.NewNull()),
		CompareValue: NewSortValue(parser.NewNull()),
		Result:       true,
	},
	{
		Name:         "SortValue EquivalentTo String and Null",
		SortValue:    NewSortValue(parser.NewString("str")),
		CompareValue: NewSortValue(parser.NewNull()),
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

var sortValueLessBench1 = NewSortValue(parser.NewInteger(12345))
var sortValueLessBench2 = NewSortValue(parser.NewInteger(67890))

func BenchmarkSortValue_Less(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = sortValueLessBench1.Less(sortValueLessBench2)
	}
}

var sortValuesEquivalentBench1 = SortValues{
	NewSortValue(parser.NewInteger(12345)),
	NewSortValue(parser.NewString("abcdefghijklmnopqrstuvwxymabcdefghijklmnopqrstuvwxyz")),
	NewSortValue(parser.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 0, GetTestLocation()))),
}

var sortValuesEquivalentBench2 = SortValues{
	NewSortValue(parser.NewInteger(12345)),
	NewSortValue(parser.NewString("abcdefghijklmnopqrstuvwxymabcdefghijklmnopqrstuvwxyz")),
	NewSortValue(parser.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 0, GetTestLocation()))),
}

func BenchmarkSortValues_EquivalentTo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = sortValuesEquivalentBench1.EquivalentTo(sortValuesEquivalentBench2)
	}
}
