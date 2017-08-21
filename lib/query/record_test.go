package query

import (
	"reflect"
	"testing"
	"time"

	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/ternary"
)

func TestMergeRecordsList(t *testing.T) {
	list := []Records{
		{
			NewRecord([]parser.Primary{parser.NewInteger(1), parser.NewString("str1")}),
			NewRecord([]parser.Primary{parser.NewInteger(2), parser.NewString("str2")}),
		},
		{
			NewRecord([]parser.Primary{parser.NewInteger(3), parser.NewString("str3")}),
			NewRecord([]parser.Primary{parser.NewInteger(4), parser.NewString("str4")}),
		},
	}
	expect := Records{
		NewRecord([]parser.Primary{parser.NewInteger(1), parser.NewString("str1")}),
		NewRecord([]parser.Primary{parser.NewInteger(2), parser.NewString("str2")}),
		NewRecord([]parser.Primary{parser.NewInteger(3), parser.NewString("str3")}),
		NewRecord([]parser.Primary{parser.NewInteger(4), parser.NewString("str4")}),
	}

	result := MergeRecordsList(list)
	if !reflect.DeepEqual(result, expect) {
		t.Errorf("result = %s, want %s", result, expect)
	}
}

var serializeComparisonKeysBenchmarkRecord = NewRecord([]parser.Primary{
	parser.NewString("str"),
	parser.NewInteger(1),
	parser.NewInteger(0),
	parser.NewInteger(3),
	parser.NewFloat(1.234),
	parser.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 0, GetTestLocation())),
	parser.NewBoolean(true),
	parser.NewBoolean(false),
	parser.NewTernary(ternary.UNKNOWN),
	parser.NewNull(),
})

func BenchmarkRecord_SerializeComparisonKeys(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j < 1000; j++ {
			serializeComparisonKeysBenchmarkRecord.SerializeComparisonKeys()
		}
	}
}
