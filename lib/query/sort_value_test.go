package query

import (
	"bytes"
	"testing"
	"time"

	"github.com/mithrandie/csvq/lib/value"

	"github.com/mithrandie/ternary"
)

func TestSortValues_Serialize(t *testing.T) {
	defer func() {
		TestTx.Flags.StrictEqual = false
	}()

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
	expect := "[N]:[I]1:[F]1.234:[D]1328289495000000000:[D]1328289495123000000:[D]1328289495123456789:[I]0:[S]STR"

	buf := &bytes.Buffer{}
	values.Serialize(buf)
	result := buf.String()
	if result != expect {
		t.Errorf("result = %q, want %q", result, expect)
	}

	TestTx.Flags.StrictEqual = true

	values = SortValues{
		NewSortValue(value.NewNull(), TestTx.Flags),
		NewSortValue(value.NewInteger(1), TestTx.Flags),
		NewSortValue(value.NewFloat(1.234), TestTx.Flags),
		NewSortValue(value.NewDatetimeFromString("2012-02-03T09:18:15-08:00", TestTx.Flags.DatetimeFormat), TestTx.Flags),
		NewSortValue(value.NewDatetimeFromString("2012-02-03T09:18:15.123-08:00", TestTx.Flags.DatetimeFormat), TestTx.Flags),
		NewSortValue(value.NewDatetimeFromString("2012-02-03T09:18:15.123456789-08:00", TestTx.Flags.DatetimeFormat), TestTx.Flags),
		NewSortValue(value.NewBoolean(false), TestTx.Flags),
		NewSortValue(value.NewString("str"), TestTx.Flags),
	}
	expect = "[N]:[I]1:[F]1.234:[D]1328289495000000000:[D]1328289495123000000:[D]1328289495123456789:[B]F:[S]str"

	buf.Reset()
	values.Serialize(buf)
	result = buf.String()
	if result != expect {
		t.Errorf("result = %q, want %q", result, expect)
	}
}

var sortValueLessTests = []struct {
	Name         string
	SortValue    value.Primary
	CompareValue value.Primary
	StrictEqual  bool
	Result       ternary.Value
}{
	{
		Name:         "Integer is less than Integer",
		SortValue:    value.NewInteger(3),
		CompareValue: value.NewInteger(5),
		StrictEqual:  false,
		Result:       ternary.TRUE,
	},
	{
		Name:         "Same Integer",
		SortValue:    value.NewInteger(3),
		CompareValue: value.NewInteger(3),
		StrictEqual:  false,
		Result:       ternary.UNKNOWN,
	},
	{
		Name:         "Integer is less than Float",
		SortValue:    value.NewInteger(3),
		CompareValue: value.NewFloat(5.4),
		StrictEqual:  false,
		Result:       ternary.TRUE,
	},
	{
		Name:         "Integer and Datetime cannot be compared",
		SortValue:    value.NewInteger(3),
		CompareValue: value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456789, GetTestLocation())),
		StrictEqual:  false,
		Result:       ternary.UNKNOWN,
	},
	{
		Name:         "Integer is less than String",
		SortValue:    value.NewInteger(3),
		CompareValue: value.NewString("4a"),
		StrictEqual:  false,
		Result:       ternary.TRUE,
	},
	{
		Name:         "Float is less than Float",
		SortValue:    value.NewFloat(3.4),
		CompareValue: value.NewFloat(5.1),
		StrictEqual:  false,
		Result:       ternary.TRUE,
	},
	{
		Name:         "Same Float",
		SortValue:    value.NewFloat(3.4),
		CompareValue: value.NewFloat(3.4),
		StrictEqual:  false,
		Result:       ternary.UNKNOWN,
	},
	{
		Name:         "Float and Datetime cannot be compared",
		SortValue:    value.NewFloat(3.4),
		CompareValue: value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456789, GetTestLocation())),
		StrictEqual:  false,
		Result:       ternary.UNKNOWN,
	},
	{
		Name:         "Float is less than String",
		SortValue:    value.NewFloat(3.4),
		CompareValue: value.NewString("4a"),
		StrictEqual:  false,
		Result:       ternary.TRUE,
	},
	{
		Name:         "Datetime is less than Datetime",
		SortValue:    value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456789, GetTestLocation())),
		CompareValue: value.NewDatetime(time.Date(2012, 2, 4, 9, 18, 15, 123456789, GetTestLocation())),
		StrictEqual:  false,
		Result:       ternary.TRUE,
	},
	{
		Name:         "Same Datetime",
		SortValue:    value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456789, GetTestLocation())),
		CompareValue: value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456789, GetTestLocation())),
		StrictEqual:  false,
		Result:       ternary.UNKNOWN,
	},
	{
		Name:         "String is less than String",
		SortValue:    value.NewString("aaa"),
		CompareValue: value.NewString("abc"),
		StrictEqual:  false,
		Result:       ternary.TRUE,
	},
	{
		Name:         "String representing integer is less than String  representing integer",
		SortValue:    value.NewString("1"),
		CompareValue: value.NewString("003"),
		StrictEqual:  false,
		Result:       ternary.TRUE,
	},
	{
		Name:         "Character Cases are Ignored",
		SortValue:    value.NewString("aaa"),
		CompareValue: value.NewString("AAA"),
		StrictEqual:  false,
		Result:       ternary.UNKNOWN,
	},
	{
		Name:         "Boolean and Ternary cannot be compared",
		SortValue:    value.NewBoolean(true),
		CompareValue: value.NewTernary(ternary.FALSE),
		StrictEqual:  false,
		Result:       ternary.UNKNOWN,
	},
	{
		Name:         "Integer and Ternary cannot be compared",
		SortValue:    value.NewInteger(3),
		CompareValue: value.NewTernary(ternary.UNKNOWN),
		StrictEqual:  false,
		Result:       ternary.UNKNOWN,
	},
	{
		Name:         "Same Strings when StrictEqual is true",
		SortValue:    value.NewString("abc"),
		CompareValue: value.NewString("abc"),
		StrictEqual:  true,
		Result:       ternary.UNKNOWN,
	},
	{
		Name:         "When StrictEqual is true, Character Cases are not ignored",
		SortValue:    value.NewString("aaa"),
		CompareValue: value.NewString("AAA"),
		StrictEqual:  true,
		Result:       ternary.FALSE,
	},
	{
		Name:         "When StrictEqual is true, Strings representing integers are compared as Strings",
		SortValue:    value.NewString("1"),
		CompareValue: value.NewString("003"),
		StrictEqual:  true,
		Result:       ternary.FALSE,
	},
	{
		Name:         "Float and String are compared even when StrictEqual is true",
		SortValue:    value.NewFloat(3.4),
		CompareValue: value.NewString("4a"),
		StrictEqual:  true,
		Result:       ternary.TRUE,
	},
}

func TestSortValue_Less(t *testing.T) {
	defer func() {
		TestTx.Flags.StrictEqual = true
	}()

	for _, v := range sortValueLessTests {
		TestTx.Flags.StrictEqual = v.StrictEqual
		sv1 := NewSortValue(v.SortValue, TestTx.Flags)
		sv2 := NewSortValue(v.CompareValue, TestTx.Flags)
		result := sv1.Less(sv2)
		if result != v.Result {
			t.Errorf("%s: result = %s, want %s", v.Name, result, v.Result)
		}
	}
}

var sortValueEquivalentToTests = []struct {
	Name         string
	SortValue    value.Primary
	CompareValue value.Primary
	StrictEqual  bool
	Result       bool
}{
	{
		Name:         "Same Integer",
		SortValue:    value.NewInteger(3),
		CompareValue: value.NewInteger(3),
		StrictEqual:  false,
		Result:       true,
	},
	{
		Name:         "Integer and Boolean with equivalent values",
		SortValue:    value.NewInteger(1),
		CompareValue: value.NewBoolean(true),
		StrictEqual:  false,
		Result:       true,
	},
	{
		Name:         "Integer and DateTime with equivalent values",
		SortValue:    value.NewInteger(1328260695),
		CompareValue: value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 0, GetTestLocation())),
		StrictEqual:  false,
		Result:       false,
	},
	{
		Name:         "Same Float",
		SortValue:    value.NewFloat(3.21),
		CompareValue: value.NewFloat(3.21),
		StrictEqual:  false,
		Result:       true,
	},
	{
		Name:         "Float and DateTime with equivalent values",
		SortValue:    value.NewFloat(1328260695.0001),
		CompareValue: value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 100000, GetTestLocation())),
		StrictEqual:  false,
		Result:       false,
	},
	{
		Name:         "Same Datetime",
		SortValue:    value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456789, GetTestLocation())),
		CompareValue: value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 123456789, GetTestLocation())),
		StrictEqual:  false,
		Result:       true,
	},
	{
		Name:         "Same Boolean",
		SortValue:    value.NewBoolean(true),
		CompareValue: value.NewBoolean(true),
		StrictEqual:  false,
		Result:       true,
	},
	{
		Name:         "Boolean and Integer with equivalent values",
		SortValue:    value.NewBoolean(true),
		CompareValue: value.NewInteger(1),
		StrictEqual:  false,
		Result:       true,
	},
	{
		Name:         "String and String with equivalent values",
		SortValue:    value.NewString("Str"),
		CompareValue: value.NewString("str"),
		StrictEqual:  false,
		Result:       true,
	},
	{
		Name:         "Null and Null",
		SortValue:    value.NewNull(),
		CompareValue: value.NewNull(),
		StrictEqual:  false,
		Result:       true,
	},
	{
		Name:         "String and Null",
		SortValue:    value.NewString("str"),
		CompareValue: value.NewNull(),
		StrictEqual:  false,
		Result:       false,
	},
	{
		Name:         "When StrictEqual is true, Strings with different case",
		SortValue:    value.NewString("Str"),
		CompareValue: value.NewString("str"),
		StrictEqual:  true,
		Result:       false,
	},
	{
		Name:         "String and Integer with equivalent values",
		SortValue:    value.NewString("001"),
		CompareValue: value.NewInteger(1),
		StrictEqual:  false,
		Result:       true,
	},
	{
		Name:         "When StrictEqual is true, Strings and Integer with equivalent values",
		SortValue:    value.NewString("001"),
		CompareValue: value.NewInteger(1),
		StrictEqual:  true,
		Result:       false,
	},
}

func TestSortValue_EquivalentTo(t *testing.T) {
	defer func() {
		TestTx.Flags.StrictEqual = true
	}()

	for _, v := range sortValueEquivalentToTests {
		TestTx.Flags.StrictEqual = v.StrictEqual
		sv1 := NewSortValue(v.SortValue, TestTx.Flags)
		sv2 := NewSortValue(v.CompareValue, TestTx.Flags)
		result := sv1.EquivalentTo(sv2)
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
