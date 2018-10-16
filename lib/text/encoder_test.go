package text

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
	Format        cmd.Format
	WithoutHeader bool
	Expect        string
}{
	{
		Name:          "Empty Fields",
		FieldList:     []string{},
		RecordSet:     [][]value.Primary{},
		Format:        cmd.TEXT,
		WithoutHeader: false,
		Expect:        "Empty Fields",
	},
	{
		Name:          "Empty RecordSet",
		FieldList:     []string{"c1", "c2"},
		RecordSet:     [][]value.Primary{},
		Format:        cmd.TEXT,
		WithoutHeader: false,
		Expect:        "Empty RecordSet",
	},
	{
		Name:      "Text",
		FieldList: []string{"c1", "c2\nsecond line", "c3"},
		RecordSet: [][]value.Primary{
			{value.NewInteger(-1), value.NewTernary(ternary.UNKNOWN), value.NewBoolean(false)},
			{value.NewFloat(2.0123), value.NewDatetimeFromString("2016-02-01T16:00:00.123456-07:00"), value.NewString("abcdef")},
			{value.NewInteger(34567890), value.NewString(" ab|cdefghijklmnopqrstuvwxyzabcdefg\nhi\"jk日本語あアｱＡ（\n"), value.NewNull()},
		},
		Format:        cmd.TEXT,
		WithoutHeader: false,
		Expect: "+----------+-------------------------------------+--------+\n" +
			"|    c1    |                 c2                  |   c3   |\n" +
			"|          |             second line             |        |\n" +
			"+----------+-------------------------------------+--------+\n" +
			"|       -1 |               UNKNOWN               | false  |\n" +
			"|   2.0123 | 2016-02-01T16:00:00.123456-07:00    | abcdef |\n" +
			"| 34567890 |  ab|cdefghijklmnopqrstuvwxyzabcdefg |  NULL  |\n" +
			"|          | hi\"jk日本語あアｱＡ（                |        |\n" +
			"|          |                                     |        |\n" +
			"+----------+-------------------------------------+--------+",
	},
	{
		Name:      "GFM",
		FieldList: []string{"c1", "c2\nsecond line", "c3"},
		RecordSet: [][]value.Primary{
			{value.NewFloat(2.0123), value.NewDatetimeFromString("2016-02-01T16:00:00.123456-07:00"), value.NewString("abcdef")},
			{value.NewInteger(34567890), value.NewString(" ab|cdefghijklmnopqrstuvwxyzabcdefg\nhi\"jk日本語あアｱＡ（\n"), value.NewNull()},
		},
		Format:        cmd.GFM,
		WithoutHeader: false,
		Expect: "" +
			"|    c1    |                          c2<br />second line                          |   c3   |\n" +
			"| -------: | --------------------------------------------------------------------- | ------ |\n" +
			"|   2.0123 | 2016-02-01T16:00:00.123456-07:00                                      | abcdef |\n" +
			"| 34567890 |  ab\\|cdefghijklmnopqrstuvwxyzabcdefg<br />hi\"jk日本語あアｱＡ（<br />  |  NULL  |",
	},
	{
		Name:      "Org",
		FieldList: []string{"c1", "c2\nsecond line", "c3"},
		RecordSet: [][]value.Primary{
			{value.NewFloat(2.0123), value.NewDatetimeFromString("2016-02-01T16:00:00.123456-07:00"), value.NewString("abcdef")},
			{value.NewInteger(34567890), value.NewString(" ab|cdefghijklmnopqrstuvwxyzabcdefg\nhi\"jk日本語あアｱＡ（\n"), value.NewNull()},
		},
		Format:        cmd.ORG,
		WithoutHeader: false,
		Expect: "" +
			"|    c1    |                          c2<br />second line                          |   c3   |\n" +
			"|----------+-----------------------------------------------------------------------+--------|\n" +
			"|   2.0123 | 2016-02-01T16:00:00.123456-07:00                                      | abcdef |\n" +
			"| 34567890 |  ab\\|cdefghijklmnopqrstuvwxyzabcdefg<br />hi\"jk日本語あアｱＡ（<br />  |  NULL  |",
	},
	{
		Name:      "Text Special Characters",
		FieldList: []string{"c1", "c2"},
		RecordSet: [][]value.Primary{
			{value.NewString("abc\r\ndef"), value.NewString("العَرَبِيَّة")},
		},
		Format:        cmd.TEXT,
		WithoutHeader: false,
		Expect: "" +
			"+------+----------+\n" +
			"|  c1  |    c2    |\n" +
			"+------+----------+\n" +
			"| abc  |  العَرَبِيَّة |\n" +
			"| def  |          |\n" +
			"+------+----------+",
	},
	{
		Name:      "Text Narrow Fields",
		FieldList: []string{"c1", "c2"},
		RecordSet: [][]value.Primary{
			{value.NewInteger(1), value.NewInteger(2)},
		},
		Format:        cmd.TEXT,
		WithoutHeader: false,
		Expect: "" +
			"+----+----+\n" +
			"| c1 | c2 |\n" +
			"+----+----+\n" +
			"|  1 |  2 |\n" +
			"+----+----+",
	},
	{
		Name:      "GFM Narrow Fields",
		FieldList: []string{"c1", "c2"},
		RecordSet: [][]value.Primary{
			{value.NewInteger(1), value.NewInteger(2)},
		},
		Format:        cmd.GFM,
		WithoutHeader: false,
		Expect: "" +
			"|  c1  |  c2  |\n" +
			"| ---: | ---: |\n" +
			"|    1 |    2 |",
	},
	{
		Name:      "Org WithoutHeader",
		FieldList: []string{"c1", "c2\nsecond line", "c3"},
		RecordSet: [][]value.Primary{
			{value.NewFloat(2.0123), value.NewDatetimeFromString("2016-02-01T16:00:00.123456-07:00"), value.NewString("abcdef")},
			{value.NewInteger(34567890), value.NewString(" ab|cdefghijklmnopqrstuvwxyzabcdefg\nhi\"jk日本語あアｱＡ（\n"), value.NewNull()},
		},
		Format:        cmd.ORG,
		WithoutHeader: true,
		Expect: "" +
			"|   2.0123 | 2016-02-01T16:00:00.123456-07:00                                     | abcdef |\n" +
			"| 34567890 |  ab\\|cdefghijklmnopqrstuvwxyzabcdefg<br />hi\"jk日本語あアｱＡ（<br /> |  NULL  |",
	},
}

func TestEncoder_Encode(t *testing.T) {
	e := NewEncoder()

	for _, v := range encoderEncodeTests {
		e.Format = v.Format
		e.WithoutHeader = v.WithoutHeader

		result := e.Encode(v.FieldList, v.RecordSet)

		if result != v.Expect {
			t.Errorf("%s: result = %q, want %q", v.Name, result, v.Expect)
		}
	}
}
