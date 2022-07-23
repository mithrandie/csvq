package query

import (
	"testing"
	"time"

	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/value"

	"github.com/mithrandie/ternary"
)

func TestFormatTableName(t *testing.T) {
	path := "/path/to/file.txt"
	expect := "file"
	result := FormatTableName(path)
	if result != expect {
		t.Errorf("table name = %q, want %q for %q", result, expect, path)
	}

	path = "/path/to/file"
	expect = "file"
	result = FormatTableName(path)
	if result != expect {
		t.Errorf("table name = %q, want %q for %q", result, expect, path)
	}

	path = "file.txt"
	expect = "file"
	result = FormatTableName(path)
	if result != expect {
		t.Errorf("table name = %q, want %q for %q", result, expect, path)
	}

	path = ""
	expect = ""
	result = FormatTableName(path)
	if result != expect {
		t.Errorf("table name = %q, want %q for %q", result, expect, path)
	}
}

func TestFormatFieldIdentifier(t *testing.T) {
	location, _ := time.LoadLocation("UTC")

	var e parser.QueryExpression = parser.NewStringValue("str")
	expect := "@__PT:S:str"
	result := FormatFieldIdentifier(e)
	if result != expect {
		t.Errorf("field identifier = %q, want %q for %#v", result, expect, e)
	}

	e = parser.NewIntegerValue(1)
	expect = "@__PT:I:1"
	result = FormatFieldIdentifier(e)
	if result != expect {
		t.Errorf("field identifier = %q, want %q for %#v", result, expect, e)
	}

	e = parser.NewFloatValue(1.2)
	expect = "@__PT:F:1.2"
	result = FormatFieldIdentifier(e)
	if result != expect {
		t.Errorf("field identifier = %q, want %q for %#v", result, expect, e)
	}

	e = parser.PrimitiveType{
		Value: value.NewBoolean(true),
	}
	expect = "@__PT:B:true"
	result = FormatFieldIdentifier(e)
	if result != expect {
		t.Errorf("field identifier = %q, want %q for %#v", result, expect, e)
	}

	e = parser.NewTernaryValue(ternary.TRUE)
	expect = "@__PT:T:TRUE"
	result = FormatFieldIdentifier(e)
	if result != expect {
		t.Errorf("field identifier = %q, want %q for %#v", result, expect, e)
	}

	e = parser.NewDatetimeValueFromString("2006-01-02 15:04:05 -08:00", nil, location)
	expect = "@__PT:D:2006-01-02T15:04:05-08:00"
	result = FormatFieldIdentifier(e)
	if result != expect {
		t.Errorf("field identifier = %q, want %q for %#v", result, expect, e)
	}

	e = parser.NewNullValue()
	expect = "@__PT:N:NULL"
	result = FormatFieldIdentifier(e)
	if result != expect {
		t.Errorf("field identifier = %q, want %q for %#v", result, expect, e)
	}

	e = parser.FieldReference{Column: parser.Identifier{Literal: "column1", Quoted: true}}
	expect = "@__IDENT:column1"
	result = FormatFieldIdentifier(e)
	if result != expect {
		t.Errorf("field identifier = %q, want %q for %#v", result, expect, e)
	}

	e = parser.ColumnNumber{View: parser.Identifier{Literal: "table1"}, Number: value.NewInteger(1)}
	expect = "table1.1"
	result = FormatFieldIdentifier(e)
	if result != expect {
		t.Errorf("field identifier = %q, want %q for %#v", result, expect, e)
	}
}

func TestFormatFieldLabel(t *testing.T) {
	location, _ := time.LoadLocation("UTC")

	var e parser.QueryExpression = parser.NewStringValue("str")
	expect := "str"
	result := FormatFieldLabel(e)
	if result != expect {
		t.Errorf("field label = %q, want %q for %#v", result, expect, e)
	}

	e = parser.NewDatetimeValueFromString("2006-01-02 15:04:05 -08:00", nil, location)
	expect = "2006-01-02T15:04:05-08:00"
	result = FormatFieldLabel(e)
	if result != expect {
		t.Errorf("field label = %q, want %q for %#v", result, expect, e)
	}

	e = parser.NewIntegerValue(1)
	expect = "1"
	result = FormatFieldLabel(e)
	if result != expect {
		t.Errorf("field label = %q, want %q for %#v", result, expect, e)
	}

	e = parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}
	expect = "column1"
	result = FormatFieldLabel(e)
	if result != expect {
		t.Errorf("field label = %q, want %q for %#v", result, expect, e)
	}

	e = parser.ColumnNumber{View: parser.Identifier{Literal: "table1"}, Number: value.NewInteger(1)}
	expect = "table1.1"
	result = FormatFieldLabel(e)
	if result != expect {
		t.Errorf("field label = %q, want %q for %#v", result, expect, e)
	}
}
