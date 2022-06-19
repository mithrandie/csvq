package query

import (
	"bytes"
	"context"
	"testing"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/value"

	"github.com/mithrandie/go-text"
	"github.com/mithrandie/go-text/json"

	"github.com/mithrandie/ternary"
)

var encodeViewTests = []struct {
	Name                    string
	View                    *View
	Format                  cmd.Format
	LineBreak               text.LineBreak
	WriteEncoding           text.Encoding
	WriteDelimiter          rune
	WriteDelimiterPositions []int
	WriteAsSingleLine       bool
	WithoutHeader           bool
	EncloseAll              bool
	JsonEscape              json.EscapeType
	PrettyPrint             bool
	UseColor                bool
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
		Error:  "empty result set",
	},
	{
		Name: "Empty Fields",
		View: &View{
			Header: NewHeader("", []string{}),
			RecordSet: []Record{
				NewRecord([]value.Primary{}),
			},
		},
		Format: cmd.TEXT,
		Error:  "empty result set",
	},
	{
		Name: "Text",
		View: &View{
			Header: NewHeader("test", []string{"c1", "c2\nsecond line", "c3"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{value.NewInteger(-1), value.NewTernary(ternary.UNKNOWN), value.NewBoolean(true)}),
				NewRecord([]value.Primary{value.NewFloat(2.0123), value.NewDatetimeFromString("2016-02-01T16:00:00.123456-07:00", nil, GetTestLocation()), value.NewString("abcdef")}),
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
		Name: "Text with colors",
		View: &View{
			Header: NewHeader("test", []string{"c1", "c2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{value.NewInteger(-1), value.NewString("abcde")}),
				NewRecord([]value.Primary{value.NewFloat(2.0123), value.NewString("abcdef\r\nghijkl")}),
			},
		},
		Format:   cmd.TEXT,
		UseColor: true,
		Result: "" +
			"+--------+--------+\n" +
			"|   c1   |   c2   |\n" +
			"+--------+--------+\n" +
			"|     \033[35m-1\033[0m | \033[32mabcde\033[0m  |\n" +
			"| \033[35m2.0123\033[0m | \033[32mabcdef\033[0m |\n" +
			"|        | \033[32mghijkl\033[0m |\n" +
			"+--------+--------+",
	},
	{
		Name: "Box",
		View: &View{
			Header: NewHeader("test", []string{"c1", "c2\nsecond line", "c3"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{value.NewInteger(-1), value.NewTernary(ternary.UNKNOWN), value.NewBoolean(true)}),
				NewRecord([]value.Primary{value.NewFloat(2.0123), value.NewDatetimeFromString("2016-02-01T16:00:00.123456-07:00", nil, GetTestLocation()), value.NewString("abcdef")}),
				NewRecord([]value.Primary{value.NewInteger(34567890), value.NewString(" abcdefghijklmnopqrstuvwxyzabcdefg\nhi\"jk日本語あアｱＡ（\n"), value.NewNull()}),
			},
		},
		Format: cmd.BOX,
		Result: "┌──────────┬─────────────────────────────────────┬────────┐\n" +
			"│    c1    │                 c2                  │   c3   │\n" +
			"│          │             second line             │        │\n" +
			"├──────────┼─────────────────────────────────────┼────────┤\n" +
			"│       -1 │               UNKNOWN               │  true  │\n" +
			"│   2.0123 │ 2016-02-01T16:00:00.123456-07:00    │ abcdef │\n" +
			"│ 34567890 │  abcdefghijklmnopqrstuvwxyzabcdefg  │  NULL  │\n" +
			"│          │ hi\"jk日本語あアｱＡ（                │        │\n" +
			"│          │                                     │        │\n" +
			"└──────────┴─────────────────────────────────────┴────────┘",
	},
	{
		Name: "Box with colors",
		View: &View{
			Header: NewHeader("test", []string{"c1", "c2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{value.NewInteger(-1), value.NewString("abcde")}),
				NewRecord([]value.Primary{value.NewFloat(2.0123), value.NewString("abcdef\r\nghijkl")}),
			},
		},
		Format:   cmd.BOX,
		UseColor: true,
		Result: "" +
			"┌────────┬────────┐\n" +
			"│   c1   │   c2   │\n" +
			"├────────┼────────┤\n" +
			"│     \033[35m-1\033[0m │ \033[32mabcde\033[0m  │\n" +
			"│ \033[35m2.0123\033[0m │ \033[32mabcdef\033[0m │\n" +
			"│        │ \033[32mghijkl\033[0m │\n" +
			"└────────┴────────┘",
	},
	{
		Name: "Fixed-Length Format",
		View: &View{
			Header: NewHeader("test", []string{"c1", "c2", "c3"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{value.NewInteger(-1), value.NewTernary(ternary.UNKNOWN), value.NewBoolean(false)}),
				NewRecord([]value.Primary{value.NewFloat(2.0123), value.NewDatetimeFromString("2016-02-01T16:00:00.123456-07:00", nil, GetTestLocation()), value.NewString("abcdef")}),
			},
		},
		Format:                  cmd.FIXED,
		WriteDelimiterPositions: []int{10, 42, 50},
		Result: "" +
			"c1        c2                              c3      \n" +
			"        -1                                 false  \n" +
			"    2.01232016-02-01T16:00:00.123456-07:00abcdef  ",
	},
	{
		Name: "Fixed-Length Format Single Line",
		View: &View{
			Header: NewHeader("test", []string{"c1", "c2", "c3"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{value.NewInteger(-1), value.NewTernary(ternary.UNKNOWN), value.NewBoolean(false)}),
				NewRecord([]value.Primary{value.NewFloat(2.0123), value.NewDatetimeFromString("2016-02-01T16:00:00.123456-07:00", nil, GetTestLocation()), value.NewString("abcdef")}),
			},
		},
		Format:                  cmd.FIXED,
		WriteDelimiterPositions: []int{10, 42, 50},
		WriteAsSingleLine:       true,
		Result: "" +
			"        -1                                 false  " +
			"    2.01232016-02-01T16:00:00.123456-07:00abcdef  ",
	},
	{
		Name: "Fixed-Length Format Auto Filled",
		View: &View{
			Header: NewHeader("test", []string{"c1", "c2", "c3"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{value.NewInteger(-1), value.NewTernary(ternary.UNKNOWN), value.NewBoolean(false)}),
				NewRecord([]value.Primary{value.NewFloat(2.0123), value.NewDatetimeFromString("2016-02-01T16:00:00.123456-07:00", nil, GetTestLocation()), value.NewString("abcdef")}),
			},
		},
		Format:                  cmd.FIXED,
		WriteDelimiterPositions: nil,
		Result: "" +
			"c1     c2                               c3    \n" +
			"    -1                                  false \n" +
			"2.0123 2016-02-01T16:00:00.123456-07:00 abcdef",
	},
	{
		Name: "Fixed-Length Format Data Empty",
		View: &View{
			Header:    NewHeader("test", []string{"c1", "c2", "c3"}),
			RecordSet: []Record{},
		},
		Format:                  cmd.FIXED,
		WriteDelimiterPositions: []int{10, 42, 50},
		WithoutHeader:           true,
		Error:                   "data empty",
	},
	{
		Name: "Fixed-Length Format Auto Filled Data Empty",
		View: &View{
			Header:    NewHeader("test", []string{"c1", "c2", "c3"}),
			RecordSet: []Record{},
		},
		Format:                  cmd.FIXED,
		WriteDelimiterPositions: nil,
		WithoutHeader:           true,
		Error:                   "data empty",
	},
	{
		Name: "GFM LineBreak CRLF",
		View: &View{
			Header: NewHeader("test", []string{"c1", "c2\nsecond line", "c3"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{value.NewFloat(2.0123), value.NewDatetimeFromString("2016-02-01T16:00:00.123456-07:00", nil, GetTestLocation()), value.NewString("abcdef")}),
				NewRecord([]value.Primary{value.NewInteger(34567890), value.NewString(" ab|cdefghijklmnopqrstuvwxyzabcdefg\nhi\"jk日本語あアｱＡ（\n"), value.NewNull()}),
			},
		},
		Format:    cmd.GFM,
		LineBreak: text.CRLF,
		Result: "" +
			"|    c1    |                          c2<br />second line                          |   c3   |\r\n" +
			"| -------: | --------------------------------------------------------------------- | ------ |\r\n" +
			"|   2.0123 | 2016-02-01T16:00:00.123456-07:00                                      | abcdef |\r\n" +
			"| 34567890 |  ab\\|cdefghijklmnopqrstuvwxyzabcdefg<br />hi\"jk日本語あアｱＡ（<br />  |        |",
	},
	{
		Name: "GFM Data Empty",
		View: &View{
			Header:    NewHeader("test", []string{"c1", "c2\nsecond line", "c3"}),
			RecordSet: []Record{},
		},
		Format:        cmd.GFM,
		WithoutHeader: true,
		Error:         "data empty",
	},
	{
		Name: "Org-mode Table",
		View: &View{
			Header: NewHeader("test", []string{"c1", "c2\nsecond line", "c3"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{value.NewFloat(2.0123), value.NewDatetimeFromString("2016-02-01T16:00:00.123456-07:00", nil, GetTestLocation()), value.NewString("abcdef")}),
				NewRecord([]value.Primary{value.NewInteger(34567890), value.NewString(" ab|cdefghijklmnopqrstuvwxyzabcdefg\nhi\"jk日本語あアｱＡ（\n"), value.NewNull()}),
			},
		},
		Format:    cmd.ORG,
		LineBreak: text.LF,
		Result: "" +
			"|    c1    |                          c2<br />second line                          |   c3   |\n" +
			"|----------+-----------------------------------------------------------------------+--------|\n" +
			"|   2.0123 | 2016-02-01T16:00:00.123456-07:00                                      | abcdef |\n" +
			"| 34567890 |  ab\\|cdefghijklmnopqrstuvwxyzabcdefg<br />hi\"jk日本語あアｱＡ（<br />  |        |",
	},
	{
		Name: "Org-mode Table Data Empty",
		View: &View{
			Header:    NewHeader("test", []string{"c1", "c2\nsecond line", "c3"}),
			RecordSet: []Record{},
		},
		Format:        cmd.ORG,
		WithoutHeader: true,
		Error:         "data empty",
	},
	{
		Name: "TSV",
		View: &View{
			Header: NewHeader("test", []string{"c1", "c2\nsecond line", "c3"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{value.NewInteger(-1), value.NewTernary(ternary.UNKNOWN), value.NewBoolean(true)}),
				NewRecord([]value.Primary{value.NewFloat(2.0123), value.NewDatetimeFromString("2016-02-01T16:00:00.123456-07:00", nil, GetTestLocation()), value.NewString("abcdef")}),
				NewRecord([]value.Primary{value.NewInteger(34567890), value.NewString(" abcdefghijklmnopqrstuvwxyzabcdefg\nhi\"jk\n"), value.NewNull()}),
			},
		},
		Format:         cmd.TSV,
		WriteDelimiter: '\t',
		EncloseAll:     true,
		Result: "\"c1\"\t\"c2\nsecond line\"\t\"c3\"\n" +
			"-1\t\ttrue\n" +
			"2.0123\t\"2016-02-01T16:00:00.123456-07:00\"\t\"abcdef\"\n" +
			"34567890\t\" abcdefghijklmnopqrstuvwxyzabcdefg\nhi\"\"jk\n\"\t",
	},
	{
		Name: "TSV",
		View: &View{
			Header:    NewHeader("test", []string{"c1", "c2\nsecond line", "c3"}),
			RecordSet: []Record{},
		},
		Format:         cmd.TSV,
		WriteDelimiter: '\t',
		WithoutHeader:  true,
		Error:          "data empty",
	},
	{
		Name: "CSV WithoutHeader",
		View: &View{
			Header: NewHeader("test", []string{"c1", "c2\nsecond line", "c3"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{value.NewInteger(-1), value.NewTernary(ternary.UNKNOWN), value.NewBoolean(true)}),
				NewRecord([]value.Primary{value.NewFloat(2.0123), value.NewDatetimeFromString("2016-02-01T16:00:00.123456-07:00", nil, GetTestLocation()), value.NewString("abcdef")}),
				NewRecord([]value.Primary{value.NewInteger(34567890), value.NewString(" abcdefghijklmnopqrstuvwxyzabcdefg\nhi\"jk\n"), value.NewNull()}),
			},
		},
		Format:        cmd.CSV,
		WithoutHeader: true,
		EncloseAll:    true,
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
				NewRecord([]value.Primary{value.NewFloat(2.0123), value.NewDatetimeFromString("2016-02-01T16:00:00.123456-07:00", nil, GetTestLocation()), value.NewString("abcdef")}),
				NewRecord([]value.Primary{value.NewInteger(34567890), value.NewString(" abcdefghijklmnopqrstuvwxyzabcdefg\nhi\"jk\n"), value.NewNull()}),
			},
		},
		Format:     cmd.CSV,
		LineBreak:  text.CRLF,
		EncloseAll: true,
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
				NewRecord([]value.Primary{value.NewFloat(2.0123), value.NewDatetimeFromString("2016-02-01T16:00:00.123456-07:00", nil, GetTestLocation()), value.NewString("abcdef")}),
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
		Format:     cmd.JSON,
		JsonEscape: json.HexDigits,
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
		Format:     cmd.JSON,
		JsonEscape: json.AllWithHexDigits,
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
		Format:      cmd.JSON,
		JsonEscape:  json.HexDigits,
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
		Name: "JSONL",
		View: &View{
			Header: NewHeader("test", []string{"c1", "c2\nsecond line", "c3"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{value.NewInteger(-1), value.NewTernary(ternary.UNKNOWN), value.NewBoolean(true)}),
				NewRecord([]value.Primary{value.NewInteger(-1), value.NewTernary(ternary.FALSE), value.NewBoolean(true)}),
				NewRecord([]value.Primary{value.NewFloat(2.0123), value.NewDatetimeFromString("2016-02-01T16:00:00.123456-07:00", nil, GetTestLocation()), value.NewString("abcdef")}),
				NewRecord([]value.Primary{value.NewInteger(34567890), value.NewString(" abc\\defghi/jk\rlmn\topqrstuvwxyzabcdefg\nhi\"jk\n"), value.NewNull()}),
			},
		},
		Format: cmd.JSONL,
		Result: "{" +
			"\"c1\":-1," +
			"\"c2\\nsecond line\":null," +
			"\"c3\":true" +
			"}\n" +
			"{" +
			"\"c1\":-1," +
			"\"c2\\nsecond line\":false," +
			"\"c3\":true" +
			"}\n" +
			"{" +
			"\"c1\":2.0123," +
			"\"c2\\nsecond line\":\"2016-02-01T16:00:00.123456-07:00\"," +
			"\"c3\":\"abcdef\"" +
			"}\n" +
			"{" +
			"\"c1\":34567890," +
			"\"c2\\nsecond line\":\" abc\\\\defghi\\/jk\\rlmn\\topqrstuvwxyzabcdefg\\nhi\\\"jk\\n\"," +
			"\"c3\":null" +
			"}\n",
	},
	{
		Name: "LTSV",
		View: &View{
			Header: NewHeader("test", []string{"c1", "c2", "c3"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{value.NewInteger(-1), value.NewTernary(ternary.FALSE), value.NewBoolean(true)}),
				NewRecord([]value.Primary{value.NewFloat(2.0123), value.NewDatetimeFromString("2016-02-01T16:00:00.123456-07:00", nil, GetTestLocation()), value.NewString("abcdef")}),
			},
		},
		Format: cmd.LTSV,
		Result: "c1:-1\tc2:false\tc3:true\n" +
			"c1:2.0123\tc2:2016-02-01T16:00:00.123456-07:00\tc3:abcdef",
	},
	{
		Name: "LTSV Data Empty",
		View: &View{
			Header:    NewHeader("test", []string{"c1", "c2", "c3"}),
			RecordSet: []Record{},
		},
		Format: cmd.LTSV,
		Error:  "data empty",
	},
	{
		Name: "Fixed-Length Format Invalid Positions",
		View: &View{
			Header: NewHeader("test", []string{"c1", "c2", "c3"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{value.NewInteger(-1), value.NewTernary(ternary.UNKNOWN), value.NewBoolean(false)}),
				NewRecord([]value.Primary{value.NewFloat(2.0123), value.NewDatetimeFromString("2016-02-01T16:00:00.123456-07:00", nil, GetTestLocation()), value.NewString("abcdef")}),
			},
		},
		Format:                  cmd.FIXED,
		WriteDelimiterPositions: []int{10, 42, -1},
		Error:                   "data encode error: invalid delimiter position: [10, 42, -1]",
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
		Format:     cmd.JSON,
		JsonEscape: json.HexDigits,
		Error:      "data encode error: unexpected token \".\" at column 4 in \"c1..\"",
	},
	{
		Name: "LTSV Invalid Character in Label",
		View: &View{
			Header: NewHeader("test", []string{"c1:", "c2", "c3"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{value.NewInteger(-1), value.NewTernary(ternary.FALSE), value.NewBoolean(true)}),
				NewRecord([]value.Primary{value.NewFloat(2.0123), value.NewDatetimeFromString("2016-02-01T16:00:00.123456-07:00", nil, GetTestLocation()), value.NewString("abcdef")}),
			},
		},
		Format: cmd.LTSV,
		Error:  "data encode error: unpermitted character in label: U+003A",
	},
	{
		Name: "LTSV Invalid Character in Field-Value",
		View: &View{
			Header: NewHeader("test", []string{"c1", "c2", "c3"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{value.NewInteger(-1), value.NewTernary(ternary.FALSE), value.NewBoolean(true)}),
				NewRecord([]value.Primary{value.NewFloat(2.0123), value.NewDatetimeFromString("2016-02-01T16:00:00.123456-07:00", nil, GetTestLocation()), value.NewString("abc\tdef")}),
			},
		},
		Format: cmd.LTSV,
		Error:  "data encode error: unpermitted character in field-value: U+0009",
	},
	{
		Name: "CSV Encode Character Code",
		View: &View{
			Header: NewHeader("test", []string{"c1", "c2\nsecond line", "c3"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{value.NewInteger(-1), value.NewTernary(ternary.UNKNOWN), value.NewBoolean(true)}),
				NewRecord([]value.Primary{value.NewInteger(-1), value.NewTernary(ternary.FALSE), value.NewBoolean(true)}),
				NewRecord([]value.Primary{value.NewFloat(2.0123), value.NewDatetimeFromString("2016-02-01T16:00:00.123456-07:00", nil, GetTestLocation()), value.NewString("abcdef")}),
				NewRecord([]value.Primary{value.NewInteger(34567890), value.NewString(" 日本語ghijklmnopqrstuvwxyzabcdefg\nhi\"jk\n"), value.NewNull()}),
			},
		},
		Format:        cmd.CSV,
		WriteEncoding: text.SJIS,
		EncloseAll:    true,
		Result: "\"c1\",\"c2\nsecond line\",\"c3\"\n" +
			"-1,,true\n" +
			"-1,false,true\n" +
			"2.0123,\"2016-02-01T16:00:00.123456-07:00\",\"abcdef\"\n" +
			"34567890,\" " + string([]byte{0x93, 0xfa, 0x96, 0x7b, 0x8c, 0xea}) + "ghijklmnopqrstuvwxyzabcdefg\nhi\"\"jk\n\",",
	},
}

func TestEncodeView(t *testing.T) {
	defer TestTx.UseColor(false)

	buf := &bytes.Buffer{}
	ctx := context.Background()

	for _, v := range encodeViewTests {
		if v.WriteEncoding == text.AUTO {
			v.WriteEncoding = text.UTF8
		}
		if v.LineBreak == "" {
			v.LineBreak = text.LF
		}
		if v.WriteDelimiter == 0 {
			v.WriteDelimiter = ','
		}
		TestTx.UseColor(v.UseColor)

		options := TestTx.Flags.ExportOptions.Copy()
		options.Format = v.Format
		options.Delimiter = v.WriteDelimiter
		options.DelimiterPositions = v.WriteDelimiterPositions
		options.Encoding = v.WriteEncoding
		options.LineBreak = v.LineBreak
		options.WithoutHeader = v.WithoutHeader
		options.EncloseAll = v.EncloseAll
		options.JsonEscape = v.JsonEscape
		options.PrettyPrint = v.PrettyPrint
		options.SingleLine = v.WriteAsSingleLine

		buf.Reset()
		_, err := EncodeView(ctx, buf, v.View, options, TestTx.Palette)
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

		result := buf.String()
		if result != v.Result {
			t.Errorf("%s: result = %s, want %s", v.Name, result, v.Result)
		}
	}
}
