package query

import (
	"bytes"
	"reflect"
	"testing"
	"time"

	"github.com/mithrandie/csvq/lib/value"

	"github.com/mithrandie/ternary"
)

func TestMergeRecordSetList(t *testing.T) {
	list := []RecordSet{
		{
			NewRecord([]value.Primary{value.NewInteger(1), value.NewString("str1")}),
			NewRecord([]value.Primary{value.NewInteger(2), value.NewString("str2")}),
		},
		{
			NewRecord([]value.Primary{value.NewInteger(3), value.NewString("str3")}),
			NewRecord([]value.Primary{value.NewInteger(4), value.NewString("str4")}),
		},
	}
	expect := RecordSet{
		NewRecord([]value.Primary{value.NewInteger(1), value.NewString("str1")}),
		NewRecord([]value.Primary{value.NewInteger(2), value.NewString("str2")}),
		NewRecord([]value.Primary{value.NewInteger(3), value.NewString("str3")}),
		NewRecord([]value.Primary{value.NewInteger(4), value.NewString("str4")}),
	}

	result := MergeRecordSetList(list)
	if !reflect.DeepEqual(result, expect) {
		t.Errorf("result = %s, want %s", result, expect)
	}
}

var recordSerializeComparisonKeysBenchmarkRecord = NewRecord([]value.Primary{
	value.NewString("str"),
	value.NewInteger(1),
	value.NewInteger(0),
	value.NewInteger(3),
	value.NewFloat(1.234),
	value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 0, GetTestLocation())),
	value.NewBoolean(true),
	value.NewBoolean(false),
	value.NewTernary(ternary.UNKNOWN),
	value.NewNull(),
})

func BenchmarkRecord_SerializeComparisonKeys(b *testing.B) {
	buf := &bytes.Buffer{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		recordSerializeComparisonKeysBenchmarkRecord.SerializeComparisonKeys(buf, TestTx.Flags)
	}
}
