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

func TestHeader_TableColumnNames(t *testing.T) {
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
	expect := []string{
		"c1",
		"c3",
	}

	result := h.TableColumnNames()
	if !reflect.DeepEqual(result, expect) {
		t.Errorf("column names = %s, want %s for %#v", result, expect, h)
	}
}

var headerContainsNumberTests = []struct {
	Number parser.ColumnNumber
	Result int
	Error  string
}{
	{
		Number: parser.ColumnNumber{
			View:   parser.Identifier{Literal: "t1"},
			Number: parser.NewInteger(2),
		},
		Result: 1,
	},
	{
		Number: parser.ColumnNumber{
			View:   parser.Identifier{Literal: "t1"},
			Number: parser.NewInteger(0),
		},
		Error: "[L:- C:-] field number t1.0 does not exist",
	},
	{
		Number: parser.ColumnNumber{
			View:   parser.Identifier{Literal: "t1"},
			Number: parser.NewInteger(9),
		},
		Error: "[L:- C:-] field number t1.9 does not exist",
	},
}

func TestHeader_ContainsNumber(t *testing.T) {
	h := Header{
		{
			Reference: "t1",
			Column:    "c1",
			Alias:     "a1",
			Number:    1,
			FromTable: true,
		},
		{
			Reference: "t1",
			Column:    "c2",
			Alias:     "a2",
			Number:    2,
			FromTable: true,
		},
		{
			Column:    "c3",
			FromTable: false,
		},
		{
			Reference: "t2",
			Column:    "c1",
			Alias:     "a3",
			Number:    1,
			FromTable: true,
		},
	}

	for _, v := range headerContainsNumberTests {
		result, err := h.ContainsNumber(v.Number)
		if err != nil {
			if len(v.Error) < 1 {
				t.Errorf("%s: unexpected error %q", v.Number.String(), err)
			} else if err.Error() != v.Error {
				t.Errorf("%s: error %q, want error %q", v.Number.String(), err, v.Error)
			}
			continue
		}
		if 0 < len(v.Error) {
			t.Errorf("%s: no error, want error %q", v.Number.String(), v.Error)
			continue
		}
		if result != v.Result {
			t.Errorf("%s: index = %d, want %d", v.Number.String(), result, v.Result)
		}
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
			Number:    1,
			FromTable: true,
		},
		{
			Reference: "table1",
			Column:    "column2",
			Number:    2,
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
			Number:    1,
			FromTable: true,
		},
		{
			Reference: "table1",
			Column:    "column2",
			Number:    2,
			FromTable: true,
		},
	}
	if !reflect.DeepEqual(NewHeaderWithoutId(ref, words), expect) {
		t.Errorf("header = %s, want %s", NewHeaderWithoutId(ref, words), expect)
	}
}

var headerUpdateTests = []struct {
	Name      string
	Header    Header
	Reference string
	Fields    []parser.Expression
	Result    Header
	Error     string
}{
	{
		Name: "Header Update",
		Header: []HeaderField{
			{
				Reference: "table1",
				Column:    "column1",
				Alias:     "alias1",
			},
			{
				Reference: "table1",
				Column:    "column2",
				Alias:     "alias2",
			},
			{
				Reference: "table2",
				Column:    "column3",
			},
		},
		Reference: "ref1",
		Fields: []parser.Expression{
			parser.Identifier{Literal: "c1"},
			parser.Identifier{Literal: "c2"},
			parser.Identifier{Literal: "c3"},
		},
		Result: []HeaderField{
			{
				Reference: "ref1",
				Column:    "c1",
			},
			{
				Reference: "ref1",
				Column:    "c2",
			},
			{
				Reference: "ref1",
				Column:    "c3",
			},
		},
	},
	{
		Name: "Header Update Without Fields",
		Header: []HeaderField{
			{
				Reference: "table1",
				Column:    "column1",
				Alias:     "alias1",
			},
			{
				Reference: "table1",
				Column:    "column2",
				Alias:     "alias2",
			},
			{
				Reference: "table2",
				Column:    "column3",
			},
		},
		Reference: "ref1",
		Result: []HeaderField{
			{
				Reference: "ref1",
				Column:    "alias1",
			},
			{
				Reference: "ref1",
				Column:    "alias2",
			},
			{
				Reference: "ref1",
				Column:    "column3",
			},
		},
	},
	{
		Name: "Header Update Field Length Error",
		Header: []HeaderField{
			{
				Reference: "table1",
				Column:    "column1",
				Alias:     "alias1",
			},
			{
				Reference: "table1",
				Column:    "column2",
				Alias:     "alias2",
			},
			{
				Reference: "table2",
				Column:    "column3",
			},
		},
		Reference: "ref1",
		Fields: []parser.Expression{
			parser.Identifier{Literal: "c1"},
			parser.Identifier{Literal: "c2"},
		},
		Error: "[L:- C:-] field length does not match",
	},
	{
		Name: "Header Update Field Name Duplicate Error",
		Header: []HeaderField{
			{
				Reference: "table1",
				Column:    "column1",
				Alias:     "alias1",
			},
			{
				Reference: "table1",
				Column:    "column2",
				Alias:     "alias2",
			},
			{
				Reference: "table2",
				Column:    "column3",
			},
		},
		Reference: "ref1",
		Fields: []parser.Expression{
			parser.Identifier{Literal: "c1"},
			parser.Identifier{Literal: "c2"},
			parser.Identifier{Literal: "c2"},
		},
		Error: "[L:- C:-] field name c2 is a duplicate",
	},
}

func TestHeader_Update(t *testing.T) {
	for _, v := range headerUpdateTests {
		err := v.Header.Update(v.Reference, v.Fields)
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
		if !reflect.DeepEqual(v.Header, v.Result) {
			t.Errorf("%s: header = %s, want %s", v.Name, v.Header, v.Result)
		}
	}
}
