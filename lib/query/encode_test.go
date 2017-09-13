package query

import (
	"io/ioutil"
	"strings"
	"testing"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/value"

	"github.com/mithrandie/ternary"
)

func encodeToSJIS(str string) string {
	r := transform.NewReader(strings.NewReader(str), japanese.ShiftJIS.NewDecoder())
	bytes, _ := ioutil.ReadAll(r)
	return string(bytes)
}

var encodeViewTests = []struct {
	Name           string
	View           *View
	Format         cmd.Format
	LineBreak      cmd.LineBreak
	Encoding       cmd.Encoding
	WriteDelimiter rune
	WithoutHeader  bool
	Result         string
	Error          string
}{
	{
		Name: "Empty RecordSet",
		View: &View{
			Header:    NewHeader("test", []string{"c1", "c2\nsecond line", "c3"}),
			RecordSet: []Record{},
		},
		Format: cmd.TEXT,
		Result: "Empty RecordSet",
	},
	{
		Name: "Empty Fields",
		View: &View{
			Header: NewHeader("", []string{}),
			RecordSet: []Record{
				NewRecord([]value.Primary{value.NewNull()}),
			},
		},
		Format: cmd.TEXT,
		Result: "Empty Fields",
	},
	{
		Name: "Text",
		View: &View{
			Header: NewHeader("test", []string{"c1", "c2\nsecond line", "c3"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{value.NewInteger(-1), value.NewTernary(ternary.UNKNOWN), value.NewBoolean(true)}),
				NewRecord([]value.Primary{value.NewFloat(2.0123), value.NewDatetimeFromString("2016-02-01T16:00:00.123456-07:00"), value.NewString("abcdef")}),
				NewRecord([]value.Primary{value.NewInteger(34567890), value.NewString(" abcdefghijklmnopqrstuvwxyzabcdefg\nhi\"jk日本語あアｱＡ（\n"), value.NewNull()}),
			},
		},
		Format: cmd.TEXT,
		Result: "+----------+------------------------------------+--------+\n" +
			"| c1       | c2                                 | c3     |\n" +
			"|          | second line                        |        |\n" +
			"+----------+------------------------------------+--------+\n" +
			"|       -1 |                            UNKNOWN |   true |\n" +
			"|   2.0123 | 2016-02-01T16:00:00.123456-07:00   | abcdef |\n" +
			"| 34567890 |  abcdefghijklmnopqrstuvwxyzabcdefg |   NULL |\n" +
			"|          | hi\"jk日本語あアｱＡ（               |        |\n" +
			"|          |                                    |        |\n" +
			"+----------+------------------------------------+--------+",
	},
	{
		Name: "CSV",
		View: &View{
			Header: NewHeader("test", []string{"c1", "c2\nsecond line", "c3"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{value.NewInteger(-1), value.NewTernary(ternary.UNKNOWN), value.NewBoolean(true)}),
				NewRecord([]value.Primary{value.NewInteger(-1), value.NewTernary(ternary.FALSE), value.NewBoolean(true)}),
				NewRecord([]value.Primary{value.NewFloat(2.0123), value.NewDatetimeFromString("2016-02-01T16:00:00.123456-07:00"), value.NewString("abcdef")}),
				NewRecord([]value.Primary{value.NewInteger(34567890), value.NewString(" abcdefghijklmnopqrstuvwxyzabcdefg\nhi\"jk\n"), value.NewNull()}),
			},
		},
		Format: cmd.CSV,
		Result: "\"c1\",\"c2\nsecond line\",\"c3\"\n" +
			"-1,,true\n" +
			"-1,false,true\n" +
			"2.0123,\"2016-02-01T16:00:00.123456-07:00\",\"abcdef\"\n" +
			"34567890,\" abcdefghijklmnopqrstuvwxyzabcdefg\nhi\"\"jk\n\",",
	},
	{
		Name: "TSV",
		View: &View{
			Header: NewHeader("test", []string{"c1", "c2\nsecond line", "c3"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{value.NewInteger(-1), value.NewTernary(ternary.UNKNOWN), value.NewBoolean(true)}),
				NewRecord([]value.Primary{value.NewFloat(2.0123), value.NewDatetimeFromString("2016-02-01T16:00:00.123456-07:00"), value.NewString("abcdef")}),
				NewRecord([]value.Primary{value.NewInteger(34567890), value.NewString(" abcdefghijklmnopqrstuvwxyzabcdefg\nhi\"jk\n"), value.NewNull()}),
			},
		},
		Format:         cmd.TSV,
		WriteDelimiter: '\t',
		Result: "\"c1\"\t\"c2\nsecond line\"\t\"c3\"\n" +
			"-1\t\ttrue\n" +
			"2.0123\t\"2016-02-01T16:00:00.123456-07:00\"\t\"abcdef\"\n" +
			"34567890\t\" abcdefghijklmnopqrstuvwxyzabcdefg\nhi\"\"jk\n\"\t",
	},
	{
		Name: "CSV WithoutHeader",
		View: &View{
			Header: NewHeader("test", []string{"c1", "c2\nsecond line", "c3"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{value.NewInteger(-1), value.NewTernary(ternary.UNKNOWN), value.NewBoolean(true)}),
				NewRecord([]value.Primary{value.NewFloat(2.0123), value.NewDatetimeFromString("2016-02-01T16:00:00.123456-07:00"), value.NewString("abcdef")}),
				NewRecord([]value.Primary{value.NewInteger(34567890), value.NewString(" abcdefghijklmnopqrstuvwxyzabcdefg\nhi\"jk\n"), value.NewNull()}),
			},
		},
		Format:        cmd.CSV,
		WithoutHeader: true,
		Result: "-1,,true\n" +
			"2.0123,\"2016-02-01T16:00:00.123456-07:00\",\"abcdef\"\n" +
			"34567890,\" abcdefghijklmnopqrstuvwxyzabcdefg\nhi\"\"jk\n\",",
	},
	{
		Name: "CSV Line Break CRLF",
		View: &View{
			Header: NewHeader("test", []string{"c1", "c2\nsecond line", "c3"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{value.NewInteger(-1), value.NewTernary(ternary.UNKNOWN), value.NewBoolean(true)}),
				NewRecord([]value.Primary{value.NewFloat(2.0123), value.NewDatetimeFromString("2016-02-01T16:00:00.123456-07:00"), value.NewString("abcdef")}),
				NewRecord([]value.Primary{value.NewInteger(34567890), value.NewString(" abcdefghijklmnopqrstuvwxyzabcdefg\nhi\"jk\n"), value.NewNull()}),
			},
		},
		Format:    cmd.CSV,
		LineBreak: cmd.CRLF,
		Result: "\"c1\",\"c2\r\nsecond line\",\"c3\"\r\n" +
			"-1,,true\r\n" +
			"2.0123,\"2016-02-01T16:00:00.123456-07:00\",\"abcdef\"\r\n" +
			"34567890,\" abcdefghijklmnopqrstuvwxyzabcdefg\r\nhi\"\"jk\r\n\",",
	},
	{
		Name: "JSON",
		View: &View{
			Header: NewHeader("test", []string{"c1", "c2\nsecond line", "c3"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{value.NewInteger(-1), value.NewTernary(ternary.UNKNOWN), value.NewBoolean(true)}),
				NewRecord([]value.Primary{value.NewInteger(-1), value.NewTernary(ternary.FALSE), value.NewBoolean(true)}),
				NewRecord([]value.Primary{value.NewFloat(2.0123), value.NewDatetimeFromString("2016-02-01T16:00:00.123456-07:00"), value.NewString("abcdef")}),
				NewRecord([]value.Primary{value.NewInteger(34567890), value.NewString(" abc\\defghi/jk\rlmn\topqrstuvwxyzabcdefg\nhi\"jk\n"), value.NewNull()}),
			},
		},
		Format: cmd.JSON,
		Result: "[" +
			"{" +
			"\"c1\":-1," +
			"\"c2\\nsecond line\":null," +
			"\"c3\":true" +
			"}," +
			"{" +
			"\"c1\":-1," +
			"\"c2\\nsecond line\":false," +
			"\"c3\":true" +
			"}," +
			"{" +
			"\"c1\":2.0123," +
			"\"c2\\nsecond line\":\"2016-02-01T16:00:00.123456-07:00\"," +
			"\"c3\":\"abcdef\"" +
			"}," +
			"{" +
			"\"c1\":34567890," +
			"\"c2\\nsecond line\":\" abc\\\\defghi\\/jk\\rlmn\\topqrstuvwxyzabcdefg\\nhi\\\"jk\\n\"," +
			"\"c3\":null" +
			"}" +
			"]",
	},
	{
		Name: "CSV Encode Character Code",
		View: &View{
			Header: NewHeader("test", []string{"c1", "c2\nsecond line", "c3"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{value.NewInteger(-1), value.NewTernary(ternary.UNKNOWN), value.NewBoolean(true)}),
				NewRecord([]value.Primary{value.NewInteger(-1), value.NewTernary(ternary.FALSE), value.NewBoolean(true)}),
				NewRecord([]value.Primary{value.NewFloat(2.0123), value.NewDatetimeFromString("2016-02-01T16:00:00.123456-07:00"), value.NewString("abcdef")}),
				NewRecord([]value.Primary{value.NewInteger(34567890), value.NewString(" 日本語ghijklmnopqrstuvwxyzabcdefg\nhi\"jk\n"), value.NewNull()}),
			},
		},
		Format:   cmd.CSV,
		Encoding: cmd.SJIS,
		Result: encodeToSJIS("\"c1\",\"c2\nsecond line\",\"c3\"\n" +
			"-1,,true\n" +
			"-1,false,true\n" +
			"2.0123,\"2016-02-01T16:00:00.123456-07:00\",\"abcdef\"\n" +
			"34567890,\" 日本語ghijklmnopqrstuvwxyzabcdefg\nhi\"\"jk\n\","),
	},
}

func TestEncodeView(t *testing.T) {
	flags := cmd.GetFlags()

	for _, v := range encodeViewTests {
		flags.Format = v.Format

		flags.LineBreak = cmd.LF
		if v.LineBreak != "" && v.LineBreak != cmd.LF {
			flags.LineBreak = v.LineBreak
		}
		flags.Encoding = cmd.UTF8
		if v.Encoding != cmd.UTF8 {
			flags.Encoding = v.Encoding
		}
		flags.WithoutHeader = false
		if v.WithoutHeader {
			flags.WithoutHeader = true
		}
		flags.WriteDelimiter = ','
		if v.WriteDelimiter != 0 {
			flags.WriteDelimiter = v.WriteDelimiter
		}

		s, err := EncodeView(v.View, flags.Format, flags.WriteDelimiter, flags.WithoutHeader, flags.Encoding, flags.LineBreak)
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
		if s != v.Result {
			t.Errorf("%s: result = %q, want %q", v.Name, s, v.Result)
		}
	}

	flags.LineBreak = cmd.LF
	flags.Encoding = cmd.UTF8
	flags.WithoutHeader = false
	flags.WriteDelimiter = ','
}
