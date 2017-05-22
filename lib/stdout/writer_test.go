package stdout

import (
	"testing"

	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/query"
	"github.com/mithrandie/csvq/lib/ternary"
)

func TestFormat(t *testing.T) {
	header := []string{"c1", "c2\nsecond line", "c3"}
	values := [][]parser.Primary{
		{parser.NewInteger(-1), parser.NewTernary(ternary.UNKNOWN), parser.NewBoolean(true)},
		{parser.NewFloat(2.0123), parser.NewDatetimeFromString("2016-02-01 16:00:00.123456"), parser.NewString("abcdef")},
		{parser.NewInteger(34567890), parser.NewString(" abcdefghijklmnopqrstuvwxyzabcdefg\nhijk\n"), parser.NewNull()},
	}

	view := new(query.View)
	view.Header = query.NewHeader("test", header)
	view.Records = make([]query.Record, len(values))
	for i, v := range values {
		view.Records[i] = query.NewRecord(v)
	}

	expect := `+----------+-----------------------------------+--------+
| c1       | c2                                | c3     |
|          | second line                       |        |
+----------+-----------------------------------+--------+
|       -1 |                           UNKNOWN |   true |
|   2.0123 | 2016-02-01 16:00:00.123456        | abcdef |
| 34567890 | abcdefghijklmnopqrstuvwxyzabcdefg |   NULL |
|          | hijk                              |        |
+----------+-----------------------------------+--------+
`

	result := Format(view)

	if result != expect {
		t.Errorf("result = %q, want %q for %s", result, expect, view)
	}
}

func ExampleWrite() {
	Write("write test")
	//Output:
	//write test
}
