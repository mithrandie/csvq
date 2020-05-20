package query

import (
	"reflect"
	"sync"
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
		Error: "@#INVALID is an unknown runtime information",
	},
}

func TestGetRuntimeInformation(t *testing.T) {
	defer func() {
		_ = TestTx.cachedViews.Clean(TestTx.FileContainer)
		TestTx.uncommittedViews.Clean()
		initFlag(TestTx.Flags)
	}()

	TestTx.cachedViews = GenerateViewMap([]*View{
		{FileInfo: &FileInfo{Path: "table1"}},
		{FileInfo: &FileInfo{Path: "table2"}},
		{FileInfo: &FileInfo{Path: "table3"}},
		{FileInfo: &FileInfo{Path: "table4"}},
	})
	TestTx.uncommittedViews = UncommittedViews{
		mtx: &sync.RWMutex{},
		Created: map[string]*FileInfo{
			"TABLE1": {},
			"TABLE2": {},
		},
		Updated: map[string]*FileInfo{
			"TABLE3": {},
			"TABLE4": {},
			"TABLE5": {},
			"VIEW1":  {ViewType: ViewTypeTemporaryTable},
		},
	}

	for _, v := range getRuntimeInformationTests {
		result, err := GetRuntimeInformation(TestTx, v.Input)

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
}
