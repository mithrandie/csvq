package query

import (
	"github.com/mithrandie/csvq/lib/parser"
	"reflect"
	"testing"
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
