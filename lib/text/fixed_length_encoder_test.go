package text

import (
	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/value"
	"github.com/mithrandie/ternary"
	"testing"
)

var fixedLengthEncoderEncodeTests = []struct {
	Name               string
	FieldList          []string
	RecordSet          [][]value.Primary
	DelimiterPositions []int
	LineBreak          cmd.LineBreak
	WithoutHeader      bool
	Encoding           cmd.Encoding
	Expect             string
	Error              string
}{
	{
		Name:      "Fixed-Length Encode",
		FieldList: []string{"c1", "c2", "c3"},
		RecordSet: [][]value.Primary{
			{value.NewInteger(-1), value.NewTernary(ternary.UNKNOWN), value.NewBoolean(false)},
			{value.NewFloat(2.0123), value.NewDatetimeFromString("2016-02-01T16:00:00.123456-07:00"), value.NewString("abcdef")},
			{value.NewTernary(ternary.TRUE), value.NewString(" abc"), value.NewNull()},
		},
		DelimiterPositions: []int{10, 42, 50},
		LineBreak:          cmd.LF,
		WithoutHeader:      false,
		Encoding:           cmd.UTF8,
		Expect: "" +
			"c1        c2                              c3      \n" +
			"        -1                                false   \n" +
			"    2.01232016-02-01T16:00:00.123456-07:00abcdef  \n" +
			"true       abc                                    ",
	},
	{
		Name:      "Fixed-Length Encode SJIS Format",
		FieldList: []string{"c1", "c2", "c3"},
		RecordSet: [][]value.Primary{
			{value.NewString("abc"), value.NewString("日本語"), value.NewString("def")},
			{value.NewString("ghi"), value.NewString("jkl"), value.NewString("mno")},
		},
		DelimiterPositions: []int{5, 15, 20},
		LineBreak:          cmd.LF,
		WithoutHeader:      false,
		Encoding:           cmd.SJIS,
		Expect: "" +
			"c1   c2        c3   \n" +
			"abc  日本語    def  \n" +
			"ghi  jkl       mno  ",
	},
	{
		Name:      "Fixed-Length Encode Without Header",
		FieldList: []string{"c1", "c2", "c3"},
		RecordSet: [][]value.Primary{
			{value.NewString("abc"), value.NewString("def"), value.NewString("def")},
			{value.NewString("ghi"), value.NewString("jkl"), value.NewString("mno")},
		},
		DelimiterPositions: []int{5, 15, 20},
		LineBreak:          cmd.LF,
		WithoutHeader:      true,
		Encoding:           cmd.UTF8,
		Expect: "" +
			"abc  def       def  \n" +
			"ghi  jkl       mno  ",
	},
	{
		Name:      "Fixed-Length Encode with Empty Field",
		FieldList: []string{"c1", "c2", "c3"},
		RecordSet: [][]value.Primary{
			{value.NewString("abc"), value.NewString("def"), value.NewString("def")},
			{value.NewString("ghi"), value.NewString("jkl"), value.NewString("mno")},
		},
		DelimiterPositions: []int{5, 15, 20, 25},
		LineBreak:          cmd.LF,
		WithoutHeader:      false,
		Encoding:           cmd.UTF8,
		Expect: "" +
			"c1   c2        c3        \n" +
			"abc  def       def       \n" +
			"ghi  jkl       mno       ",
	},
	{
		Name:      "Fixed-Length Encode Invalid Positions",
		FieldList: []string{"c1", "c2", "c3"},
		RecordSet: [][]value.Primary{
			{value.NewString("abc"), value.NewString("def"), value.NewString("def")},
			{value.NewString("ghi"), value.NewString("jkl"), value.NewString("mno")},
		},
		DelimiterPositions: []int{5, 15, 10},
		LineBreak:          cmd.LF,
		WithoutHeader:      false,
		Encoding:           cmd.UTF8,
		Error:              "invalid delimiter position: [5, 15, 10]",
	},
	{
		Name:      "Fixed-Length Encode Field Length too long in Header",
		FieldList: []string{"cccccc1", "c2", "c3"},
		RecordSet: [][]value.Primary{
			{value.NewString("abc"), value.NewString("def"), value.NewString("def")},
			{value.NewString("ghi"), value.NewString("jkl"), value.NewString("mno")},
		},
		DelimiterPositions: []int{5, 15, 20},
		LineBreak:          cmd.LF,
		WithoutHeader:      false,
		Encoding:           cmd.UTF8,
		Error:              "value is too long: \"cccccc1\" for 5 byte(s) length field",
	},
	{
		Name:      "Fixed-Length Encode Field Length too long in Record",
		FieldList: []string{"c1", "c2", "c3"},
		RecordSet: [][]value.Primary{
			{value.NewString("abcabc"), value.NewString("def"), value.NewString("def")},
			{value.NewString("ghi"), value.NewString("jkl"), value.NewString("mno")},
		},
		DelimiterPositions: []int{5, 15, 20},
		LineBreak:          cmd.LF,
		WithoutHeader:      false,
		Encoding:           cmd.UTF8,
		Error:              "value is too long: \"abcabc\" for 5 byte(s) length field",
	},
}

func TestFixedLengthEncoder_Encode(t *testing.T) {
	e := NewFixedLengthEncoder()

	for _, v := range fixedLengthEncoderEncodeTests {
		e.DelimiterPositions = v.DelimiterPositions
		e.LineBreak = v.LineBreak
		e.WithoutHeader = v.WithoutHeader
		e.Encoding = v.Encoding

		result, err := e.Encode(v.FieldList, v.RecordSet)

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
		if result != v.Expect {
			t.Errorf("%s: result = %q, want %q", v.Name, result, v.Expect)
		}
	}
}
