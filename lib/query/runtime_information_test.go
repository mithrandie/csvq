package query

import (
	"reflect"
	"testing"

	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/value"
)

var getRuntimeInformationTests = []struct {
	Input  parser.RuntimeInformation
	Expect value.Primary
	Error  string
}{
	{
		Input:  parser.RuntimeInformation{Name: "uncommitted"},
		Expect: value.NewBoolean(true),
	},
	{
		Input:  parser.RuntimeInformation{Name: "created"},
		Expect: value.NewInteger(2),
	},
	{
		Input:  parser.RuntimeInformation{Name: "updated"},
		Expect: value.NewInteger(3),
	},
	{
		Input:  parser.RuntimeInformation{Name: "updated_views"},
		Expect: value.NewInteger(1),
	},
	{
		Input:  parser.RuntimeInformation{Name: "loaded_tables"},
		Expect: value.NewInteger(4),
	},
	{
		Input:  parser.RuntimeInformation{Name: "working_directory"},
		Expect: value.NewString(GetWD()),
	},
	{
		Input:  parser.RuntimeInformation{Name: "version"},
		Expect: value.NewString("v1.0.0"),
	},
	{
		Input: parser.RuntimeInformation{Name: "invalid"},
		Error: "[L:- C:-] runtime information @#invalid does not exist",
	},
}

func TestGetRuntimeInformation(t *testing.T) {
	ViewCache = ViewMap{
		"TABLE1": &View{},
		"TABLE2": &View{},
		"TABLE3": &View{},
		"TABLE4": &View{},
	}
	UncommittedViews = &UncommittedViewMap{
		Created: map[string]*FileInfo{
			"TABLE1": {},
			"TABLE2": {},
		},
		Updated: map[string]*FileInfo{
			"TABLE3": {},
			"TABLE4": {},
			"TABLE5": {},
			"VIEW1":  {IsTemporary: true},
		},
	}

	for _, v := range getRuntimeInformationTests {
		result, err := GetRuntimeInformation(v.Input)

		if err != nil {
			if v.Error == "" {
				t.Errorf("unexpected error %q for %q", err.Error(), v.Input)
			} else if v.Error != err.Error() {
				t.Errorf("error %q, want error %q for %q", err.Error(), v.Error, v.Input)
			}
			continue
		}
		if v.Error != "" {
			t.Errorf("no error, want error %q for %q", v.Error, v.Input)
			continue
		}

		if !reflect.DeepEqual(result, v.Expect) {
			t.Errorf("result = %#v, want %#v for %q", result, v.Expect, v.Input)
		}
	}

	ViewCache = make(ViewMap)
	UncommittedViews = NewUncommittedViewMap()
}
