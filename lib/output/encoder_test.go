package output

import (
	"testing"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/query"
	"github.com/mithrandie/csvq/lib/ternary"
)

var encodeTests = []struct {
	Name   string
	Stmt   query.Statement
	View   *query.View
	Count  int
	Format cmd.Format
	Result string
	Error  string
}{
	{
		Name: "Text Empty",
		Stmt: query.SELECT,
		View: &query.View{
			Header:  query.NewHeader("test", []string{"c1", "c2\nsecond line", "c3"}),
			Records: []query.Record{},
		},
		Count:  0,
		Format: cmd.TEXT,
		Result: "Empty\n",
	},
	{
		Name: "Text",
		Stmt: query.SELECT,
		View: &query.View{
			Header: query.NewHeader("test", []string{"c1", "c2\nsecond line", "c3"}),
			Records: []query.Record{
				query.NewRecord([]parser.Primary{parser.NewInteger(-1), parser.NewTernary(ternary.UNKNOWN), parser.NewBoolean(true)}),
				query.NewRecord([]parser.Primary{parser.NewFloat(2.0123), parser.NewDatetimeFromString("2016-02-01 16:00:00.123456"), parser.NewString("abcdef")}),
				query.NewRecord([]parser.Primary{parser.NewInteger(34567890), parser.NewString(" abcdefghijklmnopqrstuvwxyzabcdefg\nhi\"jk\n"), parser.NewNull()}),
			},
		},
		Count:  3,
		Format: cmd.TEXT,
		Result: "+----------+-----------------------------------+--------+\n" +
			"| c1       | c2                                | c3     |\n" +
			"|          | second line                       |        |\n" +
			"+----------+-----------------------------------+--------+\n" +
			"|       -1 |                           UNKNOWN |   true |\n" +
			"|   2.0123 | 2016-02-01 16:00:00.123456        | abcdef |\n" +
			"| 34567890 | abcdefghijklmnopqrstuvwxyzabcdefg |   NULL |\n" +
			"|          | hi\"jk                             |        |\n" +
			"+----------+-----------------------------------+--------+\n",
	},
	{
		Name: "CSV",
		Stmt: query.SELECT,
		View: &query.View{
			Header: query.NewHeader("test", []string{"c1", "c2\nsecond line", "c3"}),
			Records: []query.Record{
				query.NewRecord([]parser.Primary{parser.NewInteger(-1), parser.NewTernary(ternary.UNKNOWN), parser.NewBoolean(true)}),
				query.NewRecord([]parser.Primary{parser.NewFloat(2.0123), parser.NewDatetimeFromString("2016-02-01 16:00:00.123456"), parser.NewString("abcdef")}),
				query.NewRecord([]parser.Primary{parser.NewInteger(34567890), parser.NewString(" abcdefghijklmnopqrstuvwxyzabcdefg\nhi\"jk\n"), parser.NewNull()}),
			},
		},
		Count:  3,
		Format: cmd.CSV,
		Result: "\"c1\",\"c2\nsecond line\",\"c3\"\n" +
			"-1,false,true\n" +
			"2.0123,2016-02-01 16:00:00.123456,\"abcdef\"\n" +
			"34567890,\" abcdefghijklmnopqrstuvwxyzabcdefg\nhi\"\"jk\n\",",
	},
	{
		Name: "TSV",
		Stmt: query.SELECT,
		View: &query.View{
			Header: query.NewHeader("test", []string{"c1", "c2\nsecond line", "c3"}),
			Records: []query.Record{
				query.NewRecord([]parser.Primary{parser.NewInteger(-1), parser.NewTernary(ternary.UNKNOWN), parser.NewBoolean(true)}),
				query.NewRecord([]parser.Primary{parser.NewFloat(2.0123), parser.NewDatetimeFromString("2016-02-01 16:00:00.123456"), parser.NewString("abcdef")}),
				query.NewRecord([]parser.Primary{parser.NewInteger(34567890), parser.NewString(" abcdefghijklmnopqrstuvwxyzabcdefg\nhi\"jk\n"), parser.NewNull()}),
			},
		},
		Count:  3,
		Format: cmd.TSV,
		Result: "\"c1\"\t\"c2\nsecond line\"\t\"c3\"\n" +
			"-1\tfalse\ttrue\n" +
			"2.0123\t2016-02-01 16:00:00.123456\t\"abcdef\"\n" +
			"34567890\t\" abcdefghijklmnopqrstuvwxyzabcdefg\nhi\"\"jk\n\"\t",
	},
	{
		Name: "JSON",
		Stmt: query.SELECT,
		View: &query.View{
			Header: query.NewHeader("test", []string{"c1", "c2\nsecond line", "c3"}),
			Records: []query.Record{
				query.NewRecord([]parser.Primary{parser.NewInteger(-1), parser.NewTernary(ternary.UNKNOWN), parser.NewBoolean(true)}),
				query.NewRecord([]parser.Primary{parser.NewFloat(2.0123), parser.NewDatetimeFromString("2016-02-01 16:00:00.123456"), parser.NewString("abcdef")}),
				query.NewRecord([]parser.Primary{parser.NewInteger(34567890), parser.NewString(" abc\\defghi/jklmn\topqrstuvwxyzabcdefg\nhi\"jk\n"), parser.NewNull()}),
			},
		},
		Count:  3,
		Format: cmd.JSON,
		Result: "[" +
			"{" +
			"\"c1\":-1," +
			"\"c2\\nsecond line\":false," +
			"\"c3\":true" +
			"}," +
			"{" +
			"\"c1\":2.0123," +
			"\"c2\\nsecond line\":\"2016-02-01 16:00:00.123456\"," +
			"\"c3\":\"abcdef\"" +
			"}," +
			"{" +
			"\"c1\":34567890," +
			"\"c2\\nsecond line\":\" abc\\\\defghi\\/jklmn\\topqrstuvwxyzabcdefg\\nhi\\\"jk\\n\"," +
			"\"c3\":null" +
			"}" +
			"]",
	},
}

func TestEncode(t *testing.T) {
	for _, v := range encodeTests {
		result := query.Result{
			Statement: v.Stmt,
			View:      v.View,
			Count:     v.Count,
		}

		s := Encode(v.Format, result)
		if s != v.Result {
			t.Errorf("%s, result = %q, want %q", v.Name, s, v.Result)
		}
	}
}
