package query

import (
	"reflect"
	"testing"

	"github.com/mithrandie/csvq/lib/value"
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
