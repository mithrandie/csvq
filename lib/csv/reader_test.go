package csv

import (
	"reflect"
	"strings"
	"testing"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/parser"
)

var readAllTests = []struct {
	Name      string
	Delimiter rune
	Input     string
	Output    [][]parser.Primary
	LineBreak cmd.LineBreak
	Error     string
}{
	{
		Name:  "NewLineLF",
		Input: "a,b,c\nd,e,f",
		Output: [][]parser.Primary{
			{parser.NewString("a"), parser.NewString("b"), parser.NewString("c")},
			{parser.NewString("d"), parser.NewString("e"), parser.NewString("f")},
		},
		LineBreak: cmd.LF,
	},
	{
		Name:  "NewLineCR",
		Input: "a,b,c\rd,e,f",
		Output: [][]parser.Primary{
			{parser.NewString("a"), parser.NewString("b"), parser.NewString("c")},
			{parser.NewString("d"), parser.NewString("e"), parser.NewString("f")},
		},
		LineBreak: cmd.CR,
	},
	{
		Name:  "NewLineCRLF",
		Input: "a,b,c\r\nd,e,f",
		Output: [][]parser.Primary{
			{parser.NewString("a"), parser.NewString("b"), parser.NewString("c")},
			{parser.NewString("d"), parser.NewString("e"), parser.NewString("f")},
		},
		LineBreak: cmd.CRLF,
	},
	{
		Name:      "TabDelimiter",
		Delimiter: '\t',
		Input:     "a\tb\tc\nd\te\tf",
		Output: [][]parser.Primary{
			{parser.NewString("a"), parser.NewString("b"), parser.NewString("c")},
			{parser.NewString("d"), parser.NewString("e"), parser.NewString("f")},
		},
		LineBreak: cmd.LF,
	},
	{
		Name:  "QuotedString",
		Input: "a,\"b\",\"ccc\ncc\"\nd,e,",
		Output: [][]parser.Primary{
			{parser.NewString("a"), parser.NewString("b"), parser.NewString("ccc\ncc")},
			{parser.NewString("d"), parser.NewString("e"), parser.NewNull()},
		},
		LineBreak: cmd.LF,
	},
	{
		Name:  "EscapeDoubleQuote",
		Input: "a,\"b\",\"ccc\"\"cc\"\nd,e,\"\"",
		Output: [][]parser.Primary{
			{parser.NewString("a"), parser.NewString("b"), parser.NewString("ccc\"cc")},
			{parser.NewString("d"), parser.NewString("e"), parser.NewString("")},
		},
		LineBreak: cmd.LF,
	},
	{
		Name:  "DoubleQuoteInNoQuoteField",
		Input: "a,b,ccc\"cc\nd,e,",
		Output: [][]parser.Primary{
			{parser.NewString("a"), parser.NewString("b"), parser.NewString("ccc\"cc")},
			{parser.NewString("d"), parser.NewString("e"), parser.NewNull()},
		},
		LineBreak: cmd.LF,
	},
	{
		Name:  "SingleValue",
		Input: "a",
		Output: [][]parser.Primary{
			{parser.NewString("a")},
		},
		LineBreak: "",
	},
	{
		Name:  "Trailing empty lines",
		Input: "a,b,c\nd,e,f\n\n",
		Output: [][]parser.Primary{
			{parser.NewString("a"), parser.NewString("b"), parser.NewString("c")},
			{parser.NewString("d"), parser.NewString("e"), parser.NewString("f")},
		},
		LineBreak: cmd.LF,
	},
	{
		Name:  "ExtraneousQuote",
		Input: "a,\"b\",\"ccc\ncc\nd,e,",
		Error: "line 3, column 5: extraneous \" in field",
	},
	{
		Name:  "UnexpectedQuote",
		Input: "a,\"b\",\"ccc\"cc\nd,e,",
		Error: "line 1, column 11: unexpected \" in field",
	},
	{
		Name:  "NumberOfFieldsIsLess",
		Input: "a,b,c\nd,e",
		Error: "line 2, column 4: wrong number of fields in line",
	},
	{
		Name:  "NumberOfFieldsIsGreater",
		Input: "a,b,c\nd,e,f,g\nh,i,j",
		Error: "line 2, column 6: wrong number of fields in line",
	},
}

func TestReader_ReadAll(t *testing.T) {
	for _, v := range readAllTests {
		r := NewReader(strings.NewReader(v.Input))

		if v.Delimiter != 0 {
			r.Delimiter = v.Delimiter
		}

		records, err := r.ReadAll()

		if err != nil {
			if v.Error == "" {
				t.Errorf("%s: unexpected error %q", v.Name, err.Error())
			} else if v.Error != err.Error() {
				t.Errorf("%s: error %q, want error %q", v.Name, err.Error(), v.Error)
			}
			continue
		}

		if !reflect.DeepEqual(records, v.Output) {
			t.Errorf("%s: records = %q, want %q", v.Name, records, v.Output)
		}

		if r.LineBreak != v.LineBreak {
			t.Errorf("%s: line break = %q, want %q", v.Name, r.LineBreak, v.LineBreak)
		}
	}
}

func TestReader_ReadHeader(t *testing.T) {
	input := "h1,h2,h3\na,b,c\nd,e,f"
	outHeader := []string{"h1", "h2", "h3"}
	output := [][]parser.Primary{
		{parser.NewString("a"), parser.NewString("b"), parser.NewString("c")},
		{parser.NewString("d"), parser.NewString("e"), parser.NewString("f")},
	}

	r := NewReader(strings.NewReader(input))
	header, err := r.ReadHeader()
	if err != nil {
		t.Errorf("unexpected error %q", err.Error())
	}
	if !reflect.DeepEqual(header, outHeader) {
		t.Errorf("header = %q, want %q", header, outHeader)
	}

	records, err := r.ReadAll()
	if err != nil {
		t.Errorf("unexpected error %q", err.Error())
	}
	if !reflect.DeepEqual(records, output) {
		t.Errorf("records = %q, want %q", records, output)
	}
}
