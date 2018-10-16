package csv

import (
	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/value"
	"github.com/mithrandie/ternary"
	"testing"
)

var encoderEncodeTests = []struct {
	Name          string
	FieldList     []string
	RecordSet     [][]value.Primary
	Delimiter     rune
	LineBreak     cmd.LineBreak
	WithoutHeader bool
	Expect        string
}{
	{
		Name:      "CSV",
		FieldList: []string{"c1", "c2\nsecond line", "c3"},
		RecordSet: [][]value.Primary{
			{value.NewInteger(-1), value.NewTernary(ternary.UNKNOWN), value.NewBoolean(true)},
			{value.NewInteger(-1), value.NewTernary(ternary.FALSE), value.NewBoolean(true)},
			{value.NewFloat(2.0123), value.NewDatetimeFromString("2016-02-01T16:00:00.123456-07:00"), value.NewString("abcdef")},
			{value.NewInteger(34567890), value.NewString(" abcdefghijklmnopqrstuvwxyzabcdefg\nhi\"jk\n"), value.NewNull()},
		},
		Delimiter:     ',',
		LineBreak:     cmd.LF,
		WithoutHeader: false,
		Expect: "\"c1\",\"c2\nsecond line\",\"c3\"\n" +
			"-1,,true\n" +
			"-1,false,true\n" +
			"2.0123,\"2016-02-01T16:00:00.123456-07:00\",\"abcdef\"\n" +
			"34567890,\" abcdefghijklmnopqrstuvwxyzabcdefg\nhi\"\"jk\n\",",
	},
	{
		Name:      "TSV",
		FieldList: []string{"c1", "c2\nsecond line", "c3"},
		RecordSet: [][]value.Primary{
			{value.NewInteger(-1), value.NewTernary(ternary.UNKNOWN), value.NewBoolean(true)},
			{value.NewFloat(2.0123), value.NewDatetimeFromString("2016-02-01T16:00:00.123456-07:00"), value.NewString("abcdef")},
			{value.NewInteger(34567890), value.NewString(" abcdefghijklmnopqrstuvwxyzabcdefg\nhi\"jk\n"), value.NewNull()},
		},
		Delimiter:     '\t',
		LineBreak:     cmd.LF,
		WithoutHeader: false,
		Expect: "\"c1\"\t\"c2\nsecond line\"\t\"c3\"\n" +
			"-1\t\ttrue\n" +
			"2.0123\t\"2016-02-01T16:00:00.123456-07:00\"\t\"abcdef\"\n" +
			"34567890\t\" abcdefghijklmnopqrstuvwxyzabcdefg\nhi\"\"jk\n\"\t",
	},
}

func TestEncoder_Encode(t *testing.T) {
	e := NewEncoder()

	for _, v := range encoderEncodeTests {
		e.Delimiter = v.Delimiter
		e.LineBreak = v.LineBreak
		e.WithoutHeader = v.WithoutHeader

		result := e.Encode(v.FieldList, v.RecordSet)

		if result != v.Expect {
			t.Errorf("%s: result = %q, want %q", v.Name, result, v.Expect)
		}
	}
}
