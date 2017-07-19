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
		parser.FieldReference{View: parser.Identifier{Literal: "t1"}, Column: parser.Identifier{Literal: "c1"}},
		parser.FieldReference{Column: parser.Identifier{Literal: "c3"}},
	}

	result := h.TableColumns()
	if !reflect.DeepEqual(result, expect) {
		t.Errorf("columns = %s, want %s for %#v", result, expect, h)
	}
}

var headerContainsTests = []struct {
	Ref    parser.FieldReference
	Result int
	Error  string
}{
	{
		Ref: parser.FieldReference{
			View:   parser.Identifier{Literal: "t2"},
			Column: parser.Identifier{Literal: "c1"},
		},
		Result: 3,
	},
	{
		Ref: parser.FieldReference{
			Column: parser.Identifier{Literal: "a2"},
		},
		Result: 1,
	},
	{
		Ref: parser.FieldReference{
			Column: parser.Identifier{Literal: "c2"},
		},
		Result: 1,
	},
	{
		Ref: parser.FieldReference{
			Column: parser.Identifier{Literal: "c1"},
		},
		Error: "[L:- C:-] field c1 is ambiguous",
	},
	{
		Ref: parser.FieldReference{
			Column: parser.Identifier{Literal: "d1"},
		},
		Error: "[L:- C:-] field d1 does not exist",
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
		result, err := h.Contains(v.Ref)
		if err != nil {
			if len(v.Error) < 1 {
				t.Errorf("%s: unexpected error %q", v.Ref.String(), err)
			} else if err.Error() != v.Error {
				t.Errorf("%s: error %q, want error %q", v.Ref.String(), err, v.Error)
			}
			continue
		}
		if 0 < len(v.Error) {
			t.Errorf("%s: no error, want error %q", v.Ref.String(), v.Error)
			continue
		}
		if result != v.Result {
			t.Errorf("%s: index = %d, want %d", v.Ref.String(), result, v.Result)
		}
	}
}

func TestNewHeader(t *testing.T) {
	ref := "table1"
	words := []string{"column1", "column2"}
	var expect Header = []HeaderField{
		{
			Reference: "table1",
			Column:    INTERNAL_ID_COLUMN,
		},
		{
			Reference: "table1",
			Column:    "column1",
			FromTable: true,
		},
		{
			Reference: "table1",
			Column:    "column2",
			FromTable: true,
		},
	}
	if !reflect.DeepEqual(NewHeader(ref, words), expect) {
		t.Errorf("header = %s, want %s", NewHeader(ref, words), expect)
	}
}

func TestNewHeaderWithoutId(t *testing.T) {
	ref := "table1"
	words := []string{"column1", "column2"}
	var expect Header = []HeaderField{
		{
			Reference: "table1",
			Column:    "column1",
			FromTable: true,
		},
		{
			Reference: "table1",
			Column:    "column2",
			FromTable: true,
		},
	}
	if !reflect.DeepEqual(NewHeaderWithoutId(ref, words), expect) {
		t.Errorf("header = %s, want %s", NewHeaderWithoutId(ref, words), expect)
	}
}
