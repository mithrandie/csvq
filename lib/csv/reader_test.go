package csv

import (
	"reflect"
	"strings"
	"testing"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/value"
)

var readAllTests = []struct {
	Name      string
	Delimiter rune
	Input     string
	Output    [][]value.Field
	LineBreak cmd.LineBreak
	Error     string
}{
	{
		Name:  "NewLineLF",
		Input: "a,b,c\nd,e,f",
		Output: [][]value.Field{
			{value.NewField("a"), value.NewField("b"), value.NewField("c")},
			{value.NewField("d"), value.NewField("e"), value.NewField("f")},
		},
		LineBreak: cmd.LF,
	},
	{
		Name:  "NewLineCR",
		Input: "a,b,c\rd,e,f",
		Output: [][]value.Field{
			{value.NewField("a"), value.NewField("b"), value.NewField("c")},
			{value.NewField("d"), value.NewField("e"), value.NewField("f")},
		},
		LineBreak: cmd.CR,
	},
	{
		Name:  "NewLineCRLF",
		Input: "a,b,c\r\nd,e,f",
		Output: [][]value.Field{
			{value.NewField("a"), value.NewField("b"), value.NewField("c")},
			{value.NewField("d"), value.NewField("e"), value.NewField("f")},
		},
		LineBreak: cmd.CRLF,
	},
	{
		Name:      "TabDelimiter",
		Delimiter: '\t',
		Input:     "a\tb\tc\nd\te\tf",
		Output: [][]value.Field{
			{value.NewField("a"), value.NewField("b"), value.NewField("c")},
			{value.NewField("d"), value.NewField("e"), value.NewField("f")},
		},
		LineBreak: cmd.LF,
	},
	{
		Name:  "QuotedString",
		Input: "a,\"b\",\"ccc\ncc\"\nd,e,",
		Output: [][]value.Field{
			{value.NewField("a"), value.NewField("b"), value.NewField("ccc\ncc")},
			{value.NewField("d"), value.NewField("e"), nil},
		},
		LineBreak: cmd.LF,
	},
	{
		Name:  "EscapeDoubleQuote",
		Input: "a,\"b\",\"ccc\"\"cc\"\nd,e,\"\"",
		Output: [][]value.Field{
			{value.NewField("a"), value.NewField("b"), value.NewField("ccc\"cc")},
			{value.NewField("d"), value.NewField("e"), value.NewField("")},
		},
		LineBreak: cmd.LF,
	},
	{
		Name:  "DoubleQuoteInNoQuoteField",
		Input: "a,b,ccc\"cc\nd,e,",
		Output: [][]value.Field{
			{value.NewField("a"), value.NewField("b"), value.NewField("ccc\"cc")},
			{value.NewField("d"), value.NewField("e"), nil},
		},
		LineBreak: cmd.LF,
	},
	{
		Name:  "SingleValue",
		Input: "a",
		Output: [][]value.Field{
			{value.NewField("a")},
		},
		LineBreak: "",
	},
	{
		Name:  "Trailing empty lines",
		Input: "a,b,c\nd,e,f\n\n",
		Output: [][]value.Field{
			{value.NewField("a"), value.NewField("b"), value.NewField("c")},
			{value.NewField("d"), value.NewField("e"), value.NewField("f")},
		},
		LineBreak: cmd.LF,
	},
	{
		Name:  "Different Line Breaks",
		Input: "a,b,\"c\r\nd\"\ne,f,g",
		Output: [][]value.Field{
			{value.NewField("a"), value.NewField("b"), value.NewField("c\r\nd")},
			{value.NewField("e"), value.NewField("f"), value.NewField("g")},
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
		Input: "a,b,c\nd,e\nf,g,h",
		Error: "line 2, column 0: wrong number of fields in line",
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
	input := "h1,h2 ,h3\na,b,c\nd,e,f"
	outHeader := []string{"h1", "h2 ", "h3"}
	output := [][]value.Field{
		{value.NewField("a"), value.NewField("b"), value.NewField("c")},
		{value.NewField("d"), value.NewField("e"), value.NewField("f")},
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

	input = "h1,\"h2 ,h3\na,b,c\nd,e,f"
	expectErr := "line 3, column 6: extraneous \" in field"

	r = NewReader(strings.NewReader(input))
	_, err = r.ReadHeader()
	if err == nil {
		t.Errorf("no error, want error %q", expectErr)
	} else if err.Error() != expectErr {
		t.Errorf("error = %q, want error %q", err.Error(), expectErr)
	}
}

var readerReadAllBenchmarkText = strings.Repeat("aaaaaa,\"bbbbbb\",cccccc\n", 10000)

func BenchmarkReader_ReadAll(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r := strings.NewReader(readerReadAllBenchmarkText)
		reader := NewReader(r)
		reader.Delimiter = ','
		reader.WithoutNull = false
		reader.ReadAll()
	}
}

var row = []value.Field{
	[]byte("aaaaaaaaaa"),
	[]byte("bbbbbbbbbb"),
	nil,
	[]byte("cccccccccc"),
}

func BenchmarkField_ToPrimary(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j < 10000; j++ {
			fields := make([]value.Primary, len(row))
			for i, v := range row {
				fields[i] = v.ToPrimary()
			}
		}
	}
}
