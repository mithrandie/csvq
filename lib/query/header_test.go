package query

import (
	"reflect"
	"testing"

	"github.com/mithrandie/csvq/lib/parser"
)

func TestHeader_TableColumns(t *testing.T) {
	h := Header{
		{
			View:      "t1",
			Column:    "c1",
			Aliases:   []string{"a1"},
			FromTable: true,
		},
		{
			View:      "t1",
			Column:    "c2",
			Aliases:   []string{"a3"},
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
			View:      "t1",
			Column:    "c1",
			Aliases:   []string{"a1"},
			FromTable: true,
		},
		{
			View:      "t1",
			Column:    "c2",
			Aliases:   []string{"a3"},
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

var headerContainsObjectTests = []struct {
	Expr   parser.Expression
	Result int
	Error  string
}{
	{
		Expr: parser.AggregateFunction{
			Name: "count",
			Args: []parser.Expression{
				parser.AllColumns{},
			},
		},
		Result: 5,
	},
	{
		Expr: parser.FieldReference{
			View:   parser.Identifier{Literal: "t2"},
			Column: parser.Identifier{Literal: "c1"},
		},
		Result: 4,
	},
	{
		Expr: parser.ColumnNumber{
			View:   parser.Identifier{Literal: "t1"},
			Number: parser.NewInteger(2),
		},
		Result: 1,
	},
	{
		Expr:  parser.NewIntegerValue(1),
		Error: "[L:- C:-] field 1 is ambiguous",
	},
	{
		Expr:  parser.NewIntegerValue(2),
		Error: "[L:- C:-] field 2 does not exist",
	},
}

func TestHeader_ContainsObject(t *testing.T) {
	h := Header{
		{
			View:      "t1",
			Column:    "c1",
			Aliases:   []string{"a1"},
			Number:    1,
			FromTable: true,
		},
		{
			View:      "t1",
			Column:    "c2",
			Aliases:   []string{"a2"},
			Number:    2,
			FromTable: true,
		},
		{
			View:      "t1",
			Column:    "count(*)",
			Number:    3,
			FromTable: true,
		},
		{
			Column:    "c3",
			FromTable: false,
		},
		{
			View:      "t2",
			Column:    "c1",
			Aliases:   []string{"a3"},
			Number:    1,
			FromTable: true,
		},
		{
			Column:    "count(*)",
			FromTable: false,
		},
		{
			Column:    "1",
			FromTable: false,
		},
		{
			Column:    "1",
			FromTable: false,
		},
	}

	for _, v := range headerContainsObjectTests {
		result, err := h.ContainsObject(v.Expr)
		if err != nil {
			if len(v.Error) < 1 {
				t.Errorf("%s: unexpected error %q", v.Expr.String(), err)
			} else if err.Error() != v.Error {
				t.Errorf("%s: error %q, want error %q", v.Expr.String(), err, v.Error)
			}
			continue
		}
		if 0 < len(v.Error) {
			t.Errorf("%s: no error, want error %q", v.Expr.String(), v.Error)
			continue
		}
		if result != v.Result {
			t.Errorf("%s: index = %d, want %d", v.Expr.String(), result, v.Result)
		}
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
		Error: "[L:- C:-] field t1.0 does not exist",
	},
	{
		Number: parser.ColumnNumber{
			View:   parser.Identifier{Literal: "t1"},
			Number: parser.NewInteger(9),
		},
		Error: "[L:- C:-] field t1.9 does not exist",
	},
}

func TestHeader_ContainsNumber(t *testing.T) {
	h := Header{
		{
			View:      "t1",
			Column:    "c1",
			Aliases:   []string{"a1"},
			Number:    1,
			FromTable: true,
		},
		{
			View:      "t1",
			Column:    "c2",
			Aliases:   []string{"a2"},
			Number:    2,
			FromTable: true,
		},
		{
			Column:    "c3",
			FromTable: false,
		},
		{
			View:      "t2",
			Column:    "c1",
			Aliases:   []string{"a3"},
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
			View:      "t1",
			Column:    "c1",
			Aliases:   []string{"a1"},
			FromTable: true,
		},
		{
			View:      "t1",
			Column:    "c2",
			Aliases:   []string{"a2"},
			FromTable: false,
		},
		{
			Column:    "c3",
			FromTable: true,
		},
		{
			View:      "t2",
			Column:    "c1",
			Aliases:   []string{"a3"},
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
			View:   "table1",
			Column: INTERNAL_ID_COLUMN,
		},
		{
			View:      "table1",
			Column:    "column1",
			Number:    1,
			FromTable: true,
		},
		{
			View:      "table1",
			Column:    "column2",
			Number:    2,
			FromTable: true,
		},
	}
	if !reflect.DeepEqual(NewHeaderWithId(ref, words), expect) {
		t.Errorf("header = %s, want %s", NewHeaderWithId(ref, words), expect)
	}
}

func TestNewHeaderWithoutId(t *testing.T) {
	ref := "table1"
	words := []string{"column1", "column2"}
	var expect Header = []HeaderField{
		{
			View:      "table1",
			Column:    "column1",
			Number:    1,
			FromTable: true,
		},
		{
			View:      "table1",
			Column:    "column2",
			Number:    2,
			FromTable: true,
		},
	}
	if !reflect.DeepEqual(NewHeader(ref, words), expect) {
		t.Errorf("header = %s, want %s", NewHeader(ref, words), expect)
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
				View:    "table1",
				Column:  "column1",
				Aliases: []string{"alias1"},
			},
			{
				View:    "table1",
				Column:  "column2",
				Aliases: []string{"alias2"},
			},
			{
				View:   "table2",
				Column: "column3",
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
				View:   "ref1",
				Column: "c1",
			},
			{
				View:   "ref1",
				Column: "c2",
			},
			{
				View:   "ref1",
				Column: "c3",
			},
		},
	},
	{
		Name: "Header Update Without Fields",
		Header: []HeaderField{
			{
				View:    "table1",
				Column:  "column1",
				Aliases: []string{"alias1"},
			},
			{
				View:    "table1",
				Column:  "column2",
				Aliases: []string{"alias2"},
			},
			{
				View:   "table2",
				Column: "column3",
			},
		},
		Reference: "ref1",
		Result: []HeaderField{
			{
				View:   "ref1",
				Column: "column1",
			},
			{
				View:   "ref1",
				Column: "column2",
			},
			{
				View:   "ref1",
				Column: "column3",
			},
		},
	},
	{
		Name: "Header Update Field Length Error",
		Header: []HeaderField{
			{
				View:    "table1",
				Column:  "column1",
				Aliases: []string{"alias1"},
			},
			{
				View:    "table1",
				Column:  "column2",
				Aliases: []string{"alias2"},
			},
			{
				View:   "table2",
				Column: "column3",
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
				View:    "table1",
				Column:  "column1",
				Aliases: []string{"alias1"},
			},
			{
				View:    "table1",
				Column:  "column2",
				Aliases: []string{"alias2"},
			},
			{
				View:   "table2",
				Column: "column3",
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
