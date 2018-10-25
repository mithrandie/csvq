package text

import (
	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/value"
	"reflect"
	"strings"
	"testing"
)

var fixedLengthReaderReadHeaderTests = []struct {
	Name               string
	Input              string
	DelimiterPositions []int
	WithoutNull        bool
	Output             []string
	Error              string
}{
	{
		Name:               "ReadHeader",
		Input:              "abcdefghi\nklmnopqurst",
		DelimiterPositions: []int{2, 5, 9},
		WithoutNull:        false,
		Output:             []string{"ab", "cde", "fghi"},
	},
	{
		Name:               "ReadHeader with reversed position",
		Input:              "abcdefghi\nklmnopqurst",
		DelimiterPositions: []int{6, 2, 9},
		WithoutNull:        false,
		Error:              "invalid delimiter position: [6, 2, 9]",
	},
}

func TestFixedLengthReader_ReadHeader(t *testing.T) {
	for _, v := range fixedLengthReaderReadHeaderTests {
		r := NewFixedLengthReader(strings.NewReader(v.Input), v.DelimiterPositions)
		r.WithoutNull = v.WithoutNull

		header, err := r.ReadHeader()

		if err != nil {
			if v.Error == "" {
				t.Errorf("%s: unexpected error %q", v.Name, err.Error())
			} else if v.Error != err.Error() {
				t.Errorf("%s: error %q, want error %q", v.Name, err.Error(), v.Error)
			}
			continue
		}
		if 0 < len(v.Error) {
			t.Errorf("%s: no error, want error %q", v.Name, v.Error)
		}
		if !reflect.DeepEqual(header, v.Output) {
			t.Errorf("%s: records = %q, want %q", v.Name, header, v.Output)
		}
	}
}

var fixedLengthReaderReadAllTests = []struct {
	Name               string
	Input              string
	DelimiterPositions []int
	WithoutNull        bool
	Encoding           cmd.Encoding
	Output             [][]value.Field
	ExpectLineBreak    cmd.LineBreak
	Error              string
}{
	{
		Name:               "ReadAll",
		Input:              "abcdefghi\nklmnopqurst",
		DelimiterPositions: []int{2, 5, 11},
		WithoutNull:        false,
		Encoding:           cmd.UTF8,
		Output: [][]value.Field{
			{value.NewField("ab"), value.NewField("cde"), value.NewField("fghi")},
			{value.NewField("kl"), value.NewField("mno"), value.NewField("pqurst")},
		},
		ExpectLineBreak: cmd.LF,
	},
	{
		Name:               "ReadAll with empty fields",
		Input:              "ab\nklmnopqurst",
		DelimiterPositions: []int{2, 5, 11},
		WithoutNull:        false,
		Encoding:           cmd.UTF8,
		Output: [][]value.Field{
			{value.NewField("ab"), nil, nil},
			{value.NewField("kl"), value.NewField("mno"), value.NewField("pqurst")},
		},
		ExpectLineBreak: cmd.LF,
	},
	{
		Name:               "ReadAll with empty fields without nulls",
		Input:              "ab\nklmnopqurst",
		DelimiterPositions: []int{2, 5, 11},
		WithoutNull:        true,
		Encoding:           cmd.UTF8,
		Output: [][]value.Field{
			{value.NewField("ab"), value.NewField(""), value.NewField("")},
			{value.NewField("kl"), value.NewField("mno"), value.NewField("pqurst")},
		},
		ExpectLineBreak: cmd.LF,
	},
	{
		Name:               "ReadAll with nil",
		Input:              "abcdefghi\nklmnopqurst",
		DelimiterPositions: nil,
		WithoutNull:        false,
		Encoding:           cmd.UTF8,
		Output: [][]value.Field{
			{},
			{},
		},
		ExpectLineBreak: cmd.LF,
	},
	{
		Name:               "ReadAll with empty positions",
		Input:              "abcdefghi\nklmnopqurst",
		DelimiterPositions: []int{},
		WithoutNull:        false,
		Encoding:           cmd.UTF8,
		Output: [][]value.Field{
			{},
			{},
		},
		ExpectLineBreak: cmd.LF,
	},
	{
		Name:               "ReadAll with trimming value",
		Input:              "abcdefghi\nk   lm  no    pqurst",
		DelimiterPositions: []int{3, 11, 20},
		WithoutNull:        false,
		Encoding:           cmd.UTF8,
		Output: [][]value.Field{
			{value.NewField("abc"), value.NewField("defghi"), nil},
			{value.NewField("k"), value.NewField("lm  no"), value.NewField("pqurst")},
		},
		ExpectLineBreak: cmd.LF,
	},
	{
		Name:               "ReadAll with cutting length",
		Input:              "abcdefghi\nklmnopqurst\nklmnopqurst",
		DelimiterPositions: []int{2, 5, 9},
		WithoutNull:        false,
		Encoding:           cmd.UTF8,
		Output: [][]value.Field{
			{value.NewField("ab"), value.NewField("cde"), value.NewField("fghi")},
			{value.NewField("kl"), value.NewField("mno"), value.NewField("pqur")},
			{value.NewField("kl"), value.NewField("mno"), value.NewField("pqur")},
		},
		ExpectLineBreak: cmd.LF,
	},
	{
		Name:               "ReadAll from SJIS Text",
		Input:              "abcde日本語fghi\nklmnopqurst",
		DelimiterPositions: []int{5, 11, 15},
		WithoutNull:        false,
		Encoding:           cmd.SJIS,
		Output: [][]value.Field{
			{value.NewField("abcde"), value.NewField("日本語"), value.NewField("fghi")},
			{value.NewField("klmno"), value.NewField("pqurst"), nil},
		},
		ExpectLineBreak: cmd.LF,
	},
	{
		Name:               "ReadAll LineBreak CR",
		Input:              "abcdefghi\rklmnopqurst",
		DelimiterPositions: []int{2, 5, 11},
		WithoutNull:        false,
		Encoding:           cmd.UTF8,
		Output: [][]value.Field{
			{value.NewField("ab"), value.NewField("cde"), value.NewField("fghi")},
			{value.NewField("kl"), value.NewField("mno"), value.NewField("pqurst")},
		},
		ExpectLineBreak: cmd.CR,
	},
	{
		Name:               "ReadAll LineBreak CRLF",
		Input:              "abcdefghi\r\nklmnopqurst",
		DelimiterPositions: []int{2, 5, 11},
		WithoutNull:        false,
		Encoding:           cmd.UTF8,
		Output: [][]value.Field{
			{value.NewField("ab"), value.NewField("cde"), value.NewField("fghi")},
			{value.NewField("kl"), value.NewField("mno"), value.NewField("pqurst")},
		},
		ExpectLineBreak: cmd.CRLF,
	},
	{
		Name:               "ReadAll LineBreak CR with cutting length",
		Input:              "abcdefghi\rklmnopqurst\rklmnopqurst",
		DelimiterPositions: []int{2, 5, 9},
		WithoutNull:        false,
		Encoding:           cmd.UTF8,
		Output: [][]value.Field{
			{value.NewField("ab"), value.NewField("cde"), value.NewField("fghi")},
			{value.NewField("kl"), value.NewField("mno"), value.NewField("pqur")},
			{value.NewField("kl"), value.NewField("mno"), value.NewField("pqur")},
		},
		ExpectLineBreak: cmd.CR,
	},
	{
		Name:               "ReadAll LineBreak CRLF with cutting length",
		Input:              "abcdefghi\r\nklmnopqurst\r\nklmnopqurst",
		DelimiterPositions: []int{2, 5, 9},
		WithoutNull:        false,
		Encoding:           cmd.UTF8,
		Output: [][]value.Field{
			{value.NewField("ab"), value.NewField("cde"), value.NewField("fghi")},
			{value.NewField("kl"), value.NewField("mno"), value.NewField("pqur")},
			{value.NewField("kl"), value.NewField("mno"), value.NewField("pqur")},
		},
		ExpectLineBreak: cmd.CRLF,
	},
	{
		Name:               "ReadAll with negative position",
		Input:              "abcdefghi\nklmnopqurst",
		DelimiterPositions: []int{-2, 5},
		WithoutNull:        false,
		Encoding:           cmd.UTF8,
		Error:              "invalid delimiter position: [-2, 5]",
	},
	{
		Name:               "ReadAll with reversed position",
		Input:              "abcdefghi\nklmnopqurst",
		DelimiterPositions: []int{6, 2, 9},
		WithoutNull:        false,
		Encoding:           cmd.UTF8,
		Error:              "invalid delimiter position: [6, 2, 9]",
	},
	{
		Name:               "ReadAll with position error",
		Input:              "abcde日本語fghi\nklmnopqurst",
		DelimiterPositions: []int{5, 10, 15},
		WithoutNull:        false,
		Encoding:           cmd.UTF8,
		Error:              "cannot delimit lines at the position of byte array of a character",
	},
	{
		Name:               "ReadAll from SJIS Text with position error",
		Input:              "abcde日本語fghi\nklmnopqurst",
		DelimiterPositions: []int{5, 10, 15},
		WithoutNull:        false,
		Encoding:           cmd.SJIS,
		Error:              "cannot delimit lines at the position of byte array of a character",
	},
}

func TestFixedLengthReader_ReadAll(t *testing.T) {
	for _, v := range fixedLengthReaderReadAllTests {
		r := NewFixedLengthReader(strings.NewReader(v.Input), v.DelimiterPositions)
		r.WithoutNull = v.WithoutNull
		r.Encoding = v.Encoding

		records, err := r.ReadAll()

		if err != nil {
			if v.Error == "" {
				t.Errorf("%s: unexpected error %q", v.Name, err.Error())
			} else if v.Error != err.Error() {
				t.Errorf("%s: error %q, want error %q", v.Name, err.Error(), v.Error)
			}
			continue
		}
		if 0 < len(v.Error) {
			t.Errorf("%s: no error, want error %q", v.Name, v.Error)
		}
		if !reflect.DeepEqual(records, v.Output) {
			t.Errorf("%s: records = %q, want %q", v.Name, records, v.Output)
		}
		if r.LineBreak != v.ExpectLineBreak {
			t.Errorf("%s: detected line-break = %s, want %s", v.Name, r.LineBreak, v.ExpectLineBreak)
		}
	}
}
