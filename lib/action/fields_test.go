package action

import (
	"testing"

	"github.com/mithrandie/csvq/lib/query"
)

var showFieldsTests = []struct {
	Name  string
	Input string
	Error string
}{
	{
		Name:  "File Not Exist Error",
		Input: "notexist",
		Error: "file notexist does not exist",
	},
}

func TestShowFields(t *testing.T) {
	for _, v := range showFieldsTests {
		_ = query.ReleaseResources()
		proc := query.NewProcedure()
		err := ShowFields(proc, v.Input)
		if err != nil {
			if len(v.Error) < 1 {
				t.Errorf("%s: unexpected error %q", v.Name, err)
			} else if err.Error() != v.Error {
				t.Errorf("%s: error %q, want error %q", v.Name, err.Error(), v.Error)
			}
			continue
		}
		if 0 < len(v.Error) {
			t.Errorf("%s: no error, want error %q", v.Name, v.Error)
			continue
		}
	}
	_ = query.ReleaseResources()
}
