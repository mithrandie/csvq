package query

import (
	"reflect"
	"testing"

	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/value"
)

func TestHeader_TableColumns(t *testing.T) {
	h := Header{
		{
			View:        "t1",
			Column:      "c1",
			Aliases:     []string{"a1"},
			IsFromTable: true,
		},
		{
			View:        "t1",
			Column:      "c2",
			Aliases:     []string{"a3"},
			IsFromTable: false,
		},
		{
			Column:      "c3",
			IsFromTable: true,
		},
	}
	expect := []parser.QueryExpression{
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
			View:        "t1",
			Column:      "c1",
			Aliases:     []string{"a1"},
			IsFromTable: true,
		},
		{
			View:        "t1",
			Column:      "c2",
			Aliases:     []string{"a3"},
			IsFromTable: false,
		},
		{
			Column:      "c3",
			IsFromTable: true,
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
	Expr   parser.QueryExpression
	Result int
	Ok     bool
}{
	{
		Expr: parser.AggregateFunction{
			Name: "count",
			Args: []parser.QueryExpression{
				parser.AllColumns{},
			},
		},
		Result: 5,
		Ok:     true,
	},
	{
		Expr: parser.FieldReference{
			View:   parser.Identifier{Literal: "t2"},
			Column: parser.Identifier{Literal: "c1"},
		},
		Result: 4,
		Ok:     true,
	},
	{
		Expr: parser.ColumnNumber{
			View:   parser.Identifier{Literal: "t1"},
			Number: value.NewInteger(2),
		},
		Result: 1,
		Ok:     true,
	},
	{
		Expr:   parser.NewIntegerValueFromString("1"),
		Result: 6,
		Ok:     true,
	},
	{
		Expr:   parser.NewIntegerValueFromString("2"),
		Result: -1,
		Ok:     false,
	},
	{
		Expr: parser.ColumnNumber{
			View:   parser.Identifier{Literal: "t1"},
			Number: value.NewInteger(999),
		},
		Result: -1,
		Ok:     false,
	},
}

func TestHeader_ContainsObject(t *testing.T) {
	h := Header{
		{
			View:        "t1",
			Column:      "c1",
			Aliases:     []string{"a1"},
			Number:      1,
			IsFromTable: true,
		},
		{
			View:        "t1",
			Column:      "c2",
			Aliases:     []string{"a2"},
			Number:      2,
			IsFromTable: true,
		},
		{
			View:        "t1",
			Column:      "count(*)",
			Number:      3,
			IsFromTable: true,
		},
		{
			Column:      "c3",
			IsFromTable: false,
		},
		{
			View:        "t2",
			Column:      "c1",
			Aliases:     []string{"a3"},
			Number:      1,
			IsFromTable: true,
		},
		{
			Column:      "count(*)",
			IsFromTable: false,
		},
		{
			Column:      "1",
			IsFromTable: false,
		},
		{
			Column:      "1",
			IsFromTable: false,
		},
	}

	for _, v := range headerContainsObjectTests {
		result, ok := h.ContainsObject(v.Expr)
		if ok != v.Ok {
			t.Errorf("%s: contains flag = %t, want %t", v.Expr.String(), ok, v.Ok)
			continue
		}
		if result != v.Result {
			t.Errorf("%s: index = %d, want %d", v.Expr.String(), result, v.Result)
		}
	}
}

var headerFieldNumberIndexTests = []struct {
	Number parser.ColumnNumber
	Result int
	Error  string
}{
	{
		Number: parser.ColumnNumber{
			View:   parser.Identifier{Literal: "t1"},
			Number: value.NewInteger(2),
		},
		Result: 1,
	},
	{
		Number: parser.ColumnNumber{
			View:   parser.Identifier{Literal: "t1"},
			Number: value.NewInteger(0),
		},
		Error: "field not exists",
	},
	{
		Number: parser.ColumnNumber{
			View:   parser.Identifier{Literal: "t1"},
			Number: value.NewInteger(9),
		},
		Error: "field not exists",
	},
}

func TestHeader_FieldNumberIndex(t *testing.T) {
	h := Header{
		{
			View:        "t1",
			Column:      "c1",
			Aliases:     []string{"a1"},
			Number:      1,
			IsFromTable: true,
		},
		{
			View:        "t1",
			Column:      "c2",
			Aliases:     []string{"a2"},
			Number:      2,
			IsFromTable: true,
		},
		{
			Column:      "c3",
			IsFromTable: false,
		},
		{
			View:        "t2",
			Column:      "c1",
			Aliases:     []string{"a3"},
			Number:      1,
			IsFromTable: true,
		},
	}

	for _, v := range headerFieldNumberIndexTests {
		result, err := h.FieldNumberIndex(v.Number)
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

var headerFieldIndexTests = []struct {
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
			Column: parser.Identifier{Literal: "c4"},
		},
		Result: 5,
	},
	{
		Ref: parser.FieldReference{
			Column: parser.Identifier{Literal: " c4 "},
		},
		Result: 5,
	},
	{
		Ref: parser.FieldReference{
			Column: parser.Identifier{Literal: "c5"},
		},
		Result: 6,
	},
	{
		Ref: parser.FieldReference{
			Column: parser.Identifier{Literal: "c1"},
		},
		Error: "field ambiguous",
	},
	{
		Ref: parser.FieldReference{
			Column: parser.Identifier{Literal: "d1"},
		},
		Error: "field not exists",
	},
}

func TestHeader_FieldIndex(t *testing.T) {
	h := Header{
		{
			View:        "t1",
			Column:      "c1",
			Aliases:     []string{"a1"},
			IsFromTable: true,
		},
		{
			View:        "t1",
			Column:      "c2",
			Aliases:     []string{"a2"},
			IsFromTable: false,
		},
		{
			Column:      "c3",
			IsFromTable: true,
		},
		{
			View:        "t2",
			Column:      "c1",
			Aliases:     []string{"a3"},
			IsFromTable: true,
		},
		{
			View:        "t3",
			Column:      "c4",
			IsFromTable: true,
		},
		{
			Column:       "c4",
			IsFromTable:  true,
			IsJoinColumn: true,
		},
		{
			View:        "t4",
			Column:      "  c5  ",
			IsFromTable: true,
		},
	}

	for _, v := range headerFieldIndexTests {
		result, err := h.FieldIndex(v.Ref)
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
			Column: InternalIdColumn,
		},
		{
			View:        "table1",
			Column:      "column1",
			Number:      1,
			IsFromTable: true,
		},
		{
			View:        "table1",
			Column:      "column2",
			Number:      2,
			IsFromTable: true,
		},
	}
	if !reflect.DeepEqual(NewHeaderWithId(ref, words), expect) {
		t.Errorf("header = %v, want %v", NewHeaderWithId(ref, words), expect)
	}
}

func TestNewHeaderWithoutId(t *testing.T) {
	ref := "table1"
	words := []string{"column1", "column2"}
	var expect Header = []HeaderField{
		{
			View:        "table1",
			Column:      "column1",
			Number:      1,
			IsFromTable: true,
		},
		{
			View:        "table1",
			Column:      "column2",
			Number:      2,
			IsFromTable: true,
		},
	}
	if !reflect.DeepEqual(NewHeader(ref, words), expect) {
		t.Errorf("header = %v, want %v", NewHeader(ref, words), expect)
	}
}

var headerUpdateTests = []struct {
	Name      string
	Header    Header
	Reference string
	Fields    []parser.QueryExpression
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
		Fields: []parser.QueryExpression{
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
		Fields: []parser.QueryExpression{
			parser.Identifier{Literal: "c1"},
			parser.Identifier{Literal: "c2"},
		},
		Error: "field length does not match",
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
		Fields: []parser.QueryExpression{
			parser.Identifier{Literal: "c1"},
			parser.Identifier{Literal: "c2"},
			parser.Identifier{Literal: "c2"},
		},
		Error: "field name c2 is a duplicate",
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
			t.Errorf("%s: header = %v, want %v", v.Name, v.Header, v.Result)
		}
	}
}
