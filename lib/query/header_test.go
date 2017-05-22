package query

import (
	"reflect"
	"testing"

	"github.com/mithrandie/csvq/lib/parser"
)

func TestHeaderField_Label(t *testing.T) {
	hf := HeaderField{
		Column: "c1",
		Alias:  "a1",
	}
	expect := "a1"

	if hf.Label() != expect {
		t.Errorf("label = %s, want %s for %#v", hf.Label(), expect, hf)
	}

	hf = HeaderField{
		Column: "c1",
	}
	expect = "c1"

	if hf.Label() != expect {
		t.Errorf("label = %s, want %s for %#v", hf.Label(), expect, hf)
	}
}

func TestHeader_TableColumns(t *testing.T) {
	h := Header{
		{
			Reference: "t1",
			Column:    "c1",
			Alias:     "a1",
			FromTable: true,
		},
		{
			Reference: "t1",
			Column:    "c2",
			Alias:     "a3",
			FromTable: false,
		},
		{
			Column:    "c3",
			FromTable: true,
		},
	}
	expect := []parser.Expression{
		parser.Identifier{Literal: "t1.c1"},
		parser.Identifier{Literal: "c3"},
	}

	result := h.TableColumns()
	if !reflect.DeepEqual(result, expect) {
		t.Errorf("columns = %s, want %s for %#v", result, expect, h)
	}
}

var headerContainsTests = []struct {
	Ref    string
	Column string
	Result int
	Error  string
}{
	{
		Ref:    "t2",
		Column: "c1",
		Result: 3,
	},
	{
		Column: "a2",
		Result: 1,
	},
	{
		Column: "c2",
		Result: 1,
	},
	{
		Column: "c1",
		Error:  "identifier = c1: field is ambiguous",
	},
	{
		Column: "d1",
		Error:  "identifier = d1: field does not exist",
	},
}

func TestHeader_Contains(t *testing.T) {
	h := Header{
		{
			Reference: "t1",
			Column:    "c1",
			Alias:     "a1",
			FromTable: true,
		},
		{
			Reference: "t1",
			Column:    "c2",
			Alias:     "a2",
			FromTable: false,
		},
		{
			Column:    "c3",
			FromTable: true,
		},
		{
			Reference: "t2",
			Column:    "c1",
			Alias:     "a3",
			FromTable: true,
		},
	}

	for _, v := range headerContainsTests {
		identifier := v.Column
		if 0 < len(v.Ref) {
			identifier = v.Ref + "." + identifier
		}

		result, err := h.Contains(v.Ref, v.Column)
		if err != nil {
			if len(v.Error) < 1 {
				t.Errorf("%s: unexpected error %q", identifier, err)
			} else if err.Error() != v.Error {
				t.Errorf("%s: error %q, want error %q", identifier, err, v.Error)
			}
			continue
		}
		if 0 < len(v.Error) {
			t.Errorf("%s: no error, want error %q", identifier, v.Error)
			continue
		}
		if result != v.Result {
			t.Errorf("%s: index = %d, want %d", identifier, result, v.Result)
		}
	}
}
