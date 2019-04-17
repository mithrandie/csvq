package query

import (
	"reflect"
	"testing"

	"github.com/mithrandie/csvq/lib/parser"
)

var aliasMapAddTests = []struct {
	Name   string
	Alias  parser.Identifier
	Path   string
	Result AliasMap
	Error  string
}{
	{
		Name:  "AliasMap Add",
		Alias: parser.Identifier{Literal: "tbl"},
		Path:  "/path/to/tbl1.csv",
		Result: AliasMap{
			"TBL": "/PATH/TO/TBL1.CSV",
		},
	},
	{
		Name:  "AliasMap Add Table Name Duplicate Error",
		Alias: parser.Identifier{Literal: "tbl"},
		Path:  "/path/to/tbl2.csv",
		Error: "table name tbl is a duplicate",
	},
}

func TestAliasMap_Add(t *testing.T) {
	aliases := AliasMap{}

	for _, v := range aliasMapAddTests {
		err := aliases.Add(v.Alias, v.Path)
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
		if !reflect.DeepEqual(aliases, v.Result) {
			t.Errorf("%s: result = %s, want %s", v.Name, aliases, v.Result)
		}
	}
}

var aliasMapGetTests = []struct {
	Name   string
	Alias  parser.Identifier
	Result string
	Error  string
}{
	{
		Name:   "AliasMap Get",
		Alias:  parser.Identifier{Literal: "tbl"},
		Result: "/PATH/TO/TBL1.CSV",
	},
	{
		Name:  "AliasMap Get Empty Path",
		Alias: parser.Identifier{Literal: "tbl2"},
		Error: "table tbl2 is not loaded",
	},
	{
		Name:  "AliasMap Get Not Loaded Error",
		Alias: parser.Identifier{Literal: "notloaded"},
		Error: "table notloaded is not loaded",
	},
}

func TestAliasMap_Get(t *testing.T) {
	aliases := AliasMap{
		"TBL":  "/PATH/TO/TBL1.CSV",
		"TBL2": "",
	}

	for _, v := range aliasMapGetTests {
		result, err := aliases.Get(v.Alias)
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
		if result != v.Result {
			t.Errorf("%s: result = %s, want %s", v.Name, result, v.Result)
		}
	}
}
