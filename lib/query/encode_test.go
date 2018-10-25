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
	Name                    string
	View                    *View
	Format                  cmd.Format
	LineBreak               cmd.LineBreak
	WriteEncoding           cmd.Encoding
	WriteDelimiter          rune
	WriteDelimiterPositions []int
	WithoutHeader           bool
	PrettyPrint             bool
	Result                  string
	Error                   string
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
		Result: "+----------+-------------------------------------+--------+\n" +
			"|    c1    |                 c2                  |   c3   |\n" +
			"|          |             second line             |        |\n" +
			"+----------+-------------------------------------+--------+\n" +
			"|       -1 |               UNKNOWN               |  true  |\n" +
			"|   2.0123 | 2016-02-01T16:00:00.123456-07:00    | abcdef |\n" +
			"| 34567890 |  abcdefghijklmnopqrstuvwxyzabcdefg  |  NULL  |\n" +
			"|          | hi\"jk日本語あアｱＡ（                |        |\n" +
			"|          |                                     |        |\n" +
			"+----------+-------------------------------------+--------+",
	},
	{
		Name: "Fixed-Length Format",
		View: &View{
			Header: NewHeader("test", []string{"c1", "c2", "c3"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{value.NewInteger(-1), value.NewTernary(ternary.UNKNOWN), value.NewBoolean(false)}),
				NewRecord([]value.Primary{value.NewFloat(2.0123), value.NewDatetimeFromString("2016-02-01T16:00:00.123456-07:00"), value.NewString("abcdef")}),
			},
		},
		Format:                  cmd.FIXED,
		WriteDelimiterPositions: []int{10, 42, 50},
		Result: "" +
			"c1        c2                              c3      \n" +
			"        -1                                false   \n" +
			"    2.01232016-02-01T16:00:00.123456-07:00abcdef  ",
	},
	{
		Name: "GFM LineBreak CRLF",
		View: &View{
			Header: NewHeader("test", []string{"c1", "c2\nsecond line", "c3"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{value.NewFloat(2.0123), value.NewDatetimeFromString("2016-02-01T16:00:00.123456-07:00"), value.NewString("abcdef")}),
				NewRecord([]value.Primary{value.NewInteger(34567890), value.NewString(" ab|cdefghijklmnopqrstuvwxyzabcdefg\nhi\"jk日本語あアｱＡ（\n"), value.NewNull()}),
			},
		},
		Format:    cmd.GFM,
		LineBreak: cmd.CRLF,
		Result: "" +
			"|    c1    |                          c2<br />second line                          |   c3   |\r\n" +
			"| -------: | --------------------------------------------------------------------- | ------ |\r\n" +
			"|   2.0123 | 2016-02-01T16:00:00.123456-07:00                                      | abcdef |\r\n" +
			"| 34567890 |  ab\\|cdefghijklmnopqrstuvwxyzabcdefg<br />hi\"jk日本語あアｱＡ（<br />  |        |",
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
		Result: "\"c1\",\"c2\nsecond line\",\"c3\"\r\n" +
			"-1,,true\r\n" +
			"2.0123,\"2016-02-01T16:00:00.123456-07:00\",\"abcdef\"\r\n" +
			"34567890,\" abcdefghijklmnopqrstuvwxyzabcdefg\nhi\"\"jk\n\",",
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
		Name: "JSONH",
		View: &View{
			Header: NewHeader("test", []string{"c1"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{value.NewString("a")}),
				NewRecord([]value.Primary{value.NewString("b")}),
				NewRecord([]value.Primary{value.NewString("abc\\def")}),
			},
		},
		Format: cmd.JSONH,
		Result: "[" +
			"{" +
			"\"c1\":\"a\"" +
			"}," +
			"{" +
			"\"c1\":\"b\"" +
			"}," +
			"{" +
			"\"c1\":\"abc\\u005cdef\"" +
			"}" +
			"]",
	},
	{
		Name: "JSONA",
		View: &View{
			Header: NewHeader("test", []string{"c1"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{value.NewString("a")}),
				NewRecord([]value.Primary{value.NewString("b")}),
				NewRecord([]value.Primary{value.NewString("abc\\def")}),
			},
		},
		Format: cmd.JSONA,
		Result: "[" +
			"{" +
			"\"\\u0063\\u0031\":\"\\u0061\"" +
			"}," +
			"{" +
			"\"\\u0063\\u0031\":\"\\u0062\"" +
			"}," +
			"{" +
			"\"\\u0063\\u0031\":\"\\u0061\\u0062\\u0063\\u005c\\u0064\\u0065\\u0066\"" +
			"}" +
			"]",
	},
	{
		Name: "JSONH Pretty Print",
		View: &View{
			Header: NewHeader("test", []string{"c1"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{value.NewString("a")}),
				NewRecord([]value.Primary{value.NewString("b")}),
				NewRecord([]value.Primary{value.NewString("abc\\def")}),
			},
		},
		Format:      cmd.JSONH,
		PrettyPrint: true,
		Result: "[\n" +
			"  {\n" +
			"    \"c1\": \"a\"\n" +
			"  },\n" +
			"  {\n" +
			"    \"c1\": \"b\"\n" +
			"  },\n" +
			"  {\n" +
			"    \"c1\": \"abc\\u005cdef\"\n" +
			"  }\n" +
			"]",
	},
	{
		Name: "Fixed-Length Format Invalid Positions",
		View: &View{
			Header: NewHeader("test", []string{"c1", "c2", "c3"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{value.NewInteger(-1), value.NewTernary(ternary.UNKNOWN), value.NewBoolean(false)}),
				NewRecord([]value.Primary{value.NewFloat(2.0123), value.NewDatetimeFromString("2016-02-01T16:00:00.123456-07:00"), value.NewString("abcdef")}),
			},
		},
		Format:                  cmd.FIXED,
		WriteDelimiterPositions: []int{10, 42, -1},
		Error:                   "invalid delimiter position: [10, 42, -1]",
	},
	{
		Name: "JSONH Column Name Convert Error",
		View: &View{
			Header: NewHeader("test", []string{"c1.."}),
			RecordSet: []Record{
				NewRecord([]value.Primary{value.NewString("a")}),
				NewRecord([]value.Primary{value.NewString("b")}),
				NewRecord([]value.Primary{value.NewString("abc\\def")}),
			},
		},
		Format: cmd.JSONH,
		Error:  "encoding to json failed: unexpected token \".\" at column 4 in \"c1..\"",
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
		Format:        cmd.CSV,
		WriteEncoding: cmd.SJIS,
		Result: encodeToSJIS("\"c1\",\"c2\nsecond line\",\"c3\"\n" +
			"-1,,true\n" +
			"-1,false,true\n" +
			"2.0123,\"2016-02-01T16:00:00.123456-07:00\",\"abcdef\"\n" +
			"34567890,\" 日本語ghijklmnopqrstuvwxyzabcdefg\nhi\"\"jk\n\","),
	},
}

func TestEncodeView(t *testing.T) {
	for _, v := range encodeViewTests {
		if v.WriteEncoding == "" {
			v.WriteEncoding = cmd.UTF8
		}
		if v.LineBreak == "" {
			v.LineBreak = cmd.LF
		}
		if v.WriteDelimiter == cmd.UNDEF {
			v.WriteDelimiter = ','
		}

		fileInfo := &FileInfo{
			Format:             v.Format,
			Delimiter:          v.WriteDelimiter,
			DelimiterPositions: v.WriteDelimiterPositions,
			Encoding:           v.WriteEncoding,
			LineBreak:          v.LineBreak,
			NoHeader:           v.WithoutHeader,
			PrettyPrint:        v.PrettyPrint,
		}

		s, err := EncodeView(v.View, fileInfo)
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
}
