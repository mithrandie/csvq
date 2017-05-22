package query

import (
	"reflect"
	"testing"

	"github.com/mithrandie/csvq/lib/parser"
)

var coalesceTests = []struct {
	Name   string
	Args   []parser.Primary
	Result parser.Primary
	Error  string
}{
	{
		Name: "Coalesce",
		Args: []parser.Primary{
			parser.NewNull(),
			parser.NewString("str"),
		},
		Result: parser.NewString("str"),
	},
	{
		Name:  "Coalesce Argment Error",
		Args:  []parser.Primary{},
		Error: "function COALESCE is required at least 1 argument",
	},
	{
		Name: "Coalesce No Match",
		Args: []parser.Primary{
			parser.NewNull(),
			parser.NewNull(),
		},
		Result: parser.NewNull(),
	},
}

func TestCoalesce(t *testing.T) {
	for _, v := range coalesceTests {
		result, err := Coalesce(v.Args)
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
		if !reflect.DeepEqual(result, v.Result) {
			t.Errorf("%s: result = %s, want %s", v.Name, result, v.Result)
		}
	}
}
