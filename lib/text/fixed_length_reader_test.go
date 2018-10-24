package text

import (
	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/value"
	"reflect"
	"strings"
	"testing"
)

var fieldToPrimaryTests = []struct {
	Field  Field
	Expect value.Primary
}{
	{
		Field:  nil,
		Expect: value.NewNull(),
	},
	{
		Field:  []byte{},
		Expect: value.NewString(""),
	},
}

func TestField_ToPrimary(t *testing.T) {
	for _, v := range fieldToPrimaryTests {
		result := v.Field.ToPrimary()
		if !reflect.DeepEqual(result, v.Expect) {
			t.Errorf("result = %#v, want %#v for %#v", result, v.Expect, v.Field)
		}
	}
}

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
	Output             [][]Field
	Error              string
}{
	{
		Name:               "ReadAll",
		Input:              "abcdefghi\nklmnopqurst",
		DelimiterPositions: []int{2, 5, 11},
		WithoutNull:        false,
		Encoding:           cmd.UTF8,
		Output: [][]Field{
			{NewField("ab"), NewField("cde"), NewField("fghi")},
			{NewField("kl"), NewField("mno"), NewField("pqurst")},
		},
	},
	{
		Name:               "ReadAll with empty fields",
		Input:              "ab\nklmnopqurst",
		DelimiterPositions: []int{2, 5, 11},
		WithoutNull:        false,
		Encoding:           cmd.UTF8,
		Output: [][]Field{
			{NewField("ab"), nil, nil},
			{NewField("kl"), NewField("mno"), NewField("pqurst")},
		},
	},
	{
		Name:               "ReadAll with empty fields without nulls",
		Input:              "ab\nklmnopqurst",
		DelimiterPositions: []int{2, 5, 11},
		WithoutNull:        true,
		Encoding:           cmd.UTF8,
		Output: [][]Field{
			{NewField("ab"), NewField(""), NewField("")},
			{NewField("kl"), NewField("mno"), NewField("pqurst")},
		},
	},
	{
		Name:               "ReadAll with nil",
		Input:              "abcdefghi\nklmnopqurst",
		DelimiterPositions: nil,
		WithoutNull:        false,
		Encoding:           cmd.UTF8,
		Output: [][]Field{
			{},
			{},
		},
	},
	{
		Name:               "ReadAll with empty positions",
		Input:              "abcdefghi\nklmnopqurst",
		DelimiterPositions: []int{},
		WithoutNull:        false,
		Encoding:           cmd.UTF8,
		Output: [][]Field{
			{},
			{},
		},
	},
	{
		Name:               "ReadAll with trimming value",
		Input:              "abcdefghi\nk   lm  no    pqurst",
		DelimiterPositions: []int{3, 11, 20},
		WithoutNull:        false,
		Encoding:           cmd.UTF8,
		Output: [][]Field{
			{NewField("abc"), NewField("defghi"), nil},
			{NewField("k"), NewField("lm  no"), NewField("pqurst")},
		},
	},
	{
		Name:               "ReadAll with cutting length",
		Input:              "abcdefghi\nklmnopqurst",
		DelimiterPositions: []int{2, 5, 9},
		WithoutNull:        false,
		Encoding:           cmd.UTF8,
		Output: [][]Field{
			{NewField("ab"), NewField("cde"), NewField("fghi")},
			{NewField("kl"), NewField("mno"), NewField("pqur")},
		},
	},
	{
		Name:               "ReadAll from SJIS Text",
		Input:              "abcde日本語fghi\nklmnopqurst",
		DelimiterPositions: []int{5, 11, 15},
		WithoutNull:        false,
		Encoding:           cmd.SJIS,
		Output: [][]Field{
			{NewField("abcde"), NewField("日本語"), NewField("fghi")},
			{NewField("klmno"), NewField("pqurst"), nil},
		},
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
	}
}
