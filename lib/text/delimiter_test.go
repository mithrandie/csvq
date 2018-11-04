package text

import (
	"github.com/mithrandie/csvq/lib/cmd"
	"reflect"
	"strings"
	"testing"
)

var delimiterPositionsEqualTests = []struct {
	Position1 DelimiterPositions
	Position2 DelimiterPositions
	Expect    bool
}{
	{
		Position1: nil,
		Position2: []int{},
		Expect:    false,
	},
	{
		Position1: []int{4, 7},
		Position2: []int{4, 7, 10},
		Expect:    false,
	},
	{
		Position1: []int{4, 7},
		Position2: []int{4, 8},
		Expect:    false,
	},
	{
		Position1: []int{4, 7, 10},
		Position2: []int{4, 7, 10},
		Expect:    true,
	},
}

func TestDelimiterPositions_Equal(t *testing.T) {
	for _, v := range delimiterPositionsEqualTests {
		result := v.Position1.Equal(v.Position2)
		if result != v.Expect {
			t.Errorf("result = %t, want %t for %v, %v", result, v.Expect, v.Position1, v.Position2)
		}
	}
}

var delimiterPositionsStringTests = []struct {
	Position DelimiterPositions
	Expect   string
}{
	{
		Position: nil,
		Expect:   "SPACES",
	},
	{
		Position: []int{},
		Expect:   "[]",
	},
	{
		Position: []int{4, 7},
		Expect:   "[4, 7]",
	},
}

func TestDelimiterPositions_String(t *testing.T) {
	for _, v := range delimiterPositionsStringTests {
		result := v.Position.String()
		if result != v.Expect {
			t.Errorf("result = %s, want %s for %v", result, v.Expect, v.Position)
		}
	}
}

var delimiterDelimitTests = []struct {
	Input    string
	NoHeader bool
	Encoding cmd.Encoding
	Expect   []int
	Error    string
}{
	{
		Input: "" +
			"aaa bbb ccc ddd\n" +
			"\n" +
			"     \n" +
			"aaa bbb ccc ddd\n",
		NoHeader: true,
		Encoding: cmd.UTF8,
		Expect:   []int{3, 7, 11, 15},
	},
	{
		Input: "" +
			"  aaa bbb ccc ddd\n" +
			" aaa bbbbb ccc ddd\n",
		NoHeader: true,
		Encoding: cmd.UTF8,
		Expect:   []int{5, 10, 14, 18},
	},
	{
		Input: "" +
			"aaa bbbbbbb ccc ddd\n" +
			"aaa  bb bb  ccc ddd\n" +
			"aaa  bb bb  ccc ddd\n" +
			"aaa  bb bb  ccc ddd\n" +
			"aaa  bb bb  ccc ddd\n" +
			"aaa  bb bb  ccc ddd\n" +
			"aaa  bb bb  ccc ddd\n" +
			"aaa  bb bb  ccc ddd\n" +
			"aaa  bb bb  ccc ddd\n",
		NoHeader: false,
		Encoding: cmd.UTF8,
		Expect:   []int{3, 11, 15, 19},
	},
	{
		Input: "" +
			"aaa bbbbccc ddd eee\n" +
			"aaa  bb cc  ddd eee\n" +
			"aaa  bb cc  ddd eee\n" +
			"aaa  bb cc  ddd eee\n" +
			"aaa  bb cc  ddd eee\n" +
			"aaa  bb cc  ddd eee\n" +
			"aaa  bb cc  ddd eee\n" +
			"aaa  bb cc  ddd eee\n" +
			"aaa  bb cc  ddd eee\n" +
			"aaa  bb cc  ddd eee\n",
		NoHeader: true,
		Encoding: cmd.UTF8,
		Expect:   []int{3, 8, 11, 15, 19},
	},
	{
		Input: "" +
			"aaa bbb ccc ddd\n" +
			"aaa bbb ccc ddd\n" +
			"aaa bbb ccc ddd ddd ddd\n" +
			"aaa bbb ccc ddd\n",
		NoHeader: true,
		Encoding: cmd.UTF8,
		Expect:   []int{3, 7, 11, 23},
	},
	{
		Input: "" +
			"aaaa    bb  ccc ddd\n" +
			"aaaa   bbb  ccc ddd\n" +
			"aaaa   bbb  ccc ddd\n" +
			"aaaa bbbbb  ccc ddd\n" +
			"aaaa   bbb  ccc ddd\n",
		NoHeader: true,
		Encoding: cmd.UTF8,
		Expect:   []int{4, 10, 15, 19},
	},
	{
		Input: "" +
			"aaaaa   bb  ccc ddd\n" +
			"aaaaa   bb  ccc ddd\n" +
			"aaaaa   bb  ccc ddd\n" +
			"aaaaaaa bb  ccc ddd\n" +
			"aaaaa   bb  ccc ddd\n",
		NoHeader: true,
		Encoding: cmd.UTF8,
		Expect:   []int{7, 10, 15, 19},
	},
	{
		Input: "" +
			"aaaa   bbbb ccc ddd\n" +
			"aaaa   bbbb ccc ddd\n" +
			"aaaa   bbbb ccc ddd\n" +
			"aaaa   bbbb ccc ddd\n" +
			"aaaa   bbbb ccc ddd\n" +
			"aaaa   bbbb ccc ddd\n" +
			"aaaa   bbbb ccc ddd\n" +
			"aaaa   bbbb ccc ddd\n" +
			"aaaaaabbbb  ccc ddd ddd\n" +
			"aaaaa bbbb  ccc ddd\n",
		NoHeader: true,
		Encoding: cmd.UTF8,
		Expect:   []int{6, 11, 15, 23},
	},
	{
		Input: "" +
			"aaaa   bbbb  ccc ddd\n" +
			"aaaa   bbbb  ccc ddd\n" +
			"aaaa   bbbb  ccc ddd\n" +
			"aaaabbbbbb   ccc ddd\n" +
			"aaaa   bbbb  ccc ddd\n" +
			"aaaa   bbbb  ccc ddd\n" +
			"aaaa   bbbb  ccc ddd\n" +
			"aaaa   bbbb  ccc ddd\n" +
			"aaaa   bbbb  ccc ddd\n" +
			"aaaa  bbbbb  ccc ddd\n",
		NoHeader: true,
		Encoding: cmd.UTF8,
		Expect:   []int{4, 11, 16, 20},
	},
	{
		Input: "" +
			"aaaa    bb         eee fff\n" +
			"aaaa   bbb ccc ddd     fff\n" +
			"aaaa  bbbb ccc ddd      fff\n" +
			"aaaa bbbbb ccc ddd     fff\n" +
			"aaaa  bbbb ccc ddd     fff\n",
		NoHeader: true,
		Encoding: cmd.UTF8,
		Expect:   []int{4, 10, 14, 18, 22, 27},
	},
	{
		Input: "" +
			"aaaa    bb    ddd eee fff\n" +
			"aaaa   bbb ccc        fff\n" +
			"aaaa  bbbb ccc        fff\n" +
			"aaaa bbbbb ccc        fff\n" +
			"aaaa bbbbb ccc        fff\n" +
			"aaaa  bbbb ccc        fff\n",
		NoHeader: true,
		Encoding: cmd.UTF8,
		Expect:   []int{4, 10, 14, 17, 21, 25},
	},
	{
		Input: "" +
			"aaaa bb  cc  ccc ddd\n" +
			"aaaa bbb     ccc ddd\n" +
			"aaaa bbb     ccc ddd\n" +
			"aaaa bbb     ccc ddd\n" +
			"aaaa bbb     ccc ddd\n" +
			"aaaa bbb     ccc ddd\n" +
			"aaaa bbb     ccc ddd\n" +
			"aaaa bbb     ccc ddd\n" +
			"aaaa bbb  cccccc ddd\n" +
			"aaaa bbb    cccc ddd\n",
		NoHeader: true,
		Encoding: cmd.UTF8,
		Expect:   []int{4, 8, 12, 16, 20},
	},
	{
		Input: "" +
			"aaaa  bbbb  ccc ddd\n" +
			"aaaa  bbbb  ccc ddd ddd ddd\n" +
			"aaaa bbbbb  ccc ddddddd ddd\n" +
			"aaaa bbbbb  ccc ddd\n",
		NoHeader: true,
		Encoding: cmd.UTF8,
		Expect:   []int{4, 10, 15, 27},
	},
	{
		Input: "" +
			"aa  aa  bbbbb   ccc  ddddd\n" +
			"aa  aaaa bbbb   ccc  ddd ddd ddd\n" +
			"aa  aa  bbbbb  ccc ddddddd ddd\n" +
			"aa  aa  bbbbb  ccc ddddddd ddd\n" +
			"aa  aa  bbbbb  ccc ddddddd ddd\n" +
			"aa  aa  bbbbb  ccc ddddddd ddd\n" +
			"aa  aa  bbbbb  ccc ddddddd ddd\n" +
			"aa  aa  bbbbb  ccc ddddddd ddd\n" +
			"aa  aa  bbbbb  ccc ddddddd ddd\n" +
			"aa  aa  bbbbb   ccc  ddddd\n",
		NoHeader: true,
		Encoding: cmd.UTF8,
		Expect:   []int{2, 8, 13, 19, 32},
	},
	{
		Input: "" +
			"aaaaa        cccccc\n" +
			"aaaaaaa bbb    cccccc\n" +
			"aaaaaa bbb     cccccc",
		NoHeader: true,
		Encoding: cmd.UTF8,
		Expect:   []int{7, 11, 21},
	},
	{
		Input: "" +
			"aaaaaa ３文字 ccccc\n" +
			"aaaaaa bbbbbbbbb ccccc\n" +
			"aaaaaa bbbbbbbbb ccccc",
		NoHeader: true,
		Encoding: cmd.UTF8,
		Expect:   []int{6, 16, 22},
	},
	{
		Input: "" +
			"aaaaaa 日本語ﾆﾎﾝｺﾞ ccccc\n" +
			"aaaaaa bbbbbbbbbbb ccccc\n" +
			"aaaaaa bbbbbbbbbbb ccccc",
		NoHeader: true,
		Encoding: cmd.SJIS,
		Expect:   []int{6, 18, 24},
	},
	{
		Input: "" +
			"AAAAAAA\n" +
			"aaa     bbb\n" +
			"aaa     bbb\n" +
			"aaa     bbb",
		NoHeader: false,
		Encoding: cmd.UTF8,
		Expect:   []int{7, 11},
	},
}

func TestDelimiter_Delimit(t *testing.T) {
	for _, v := range delimiterDelimitTests {
		d := NewDelimiter(strings.NewReader(v.Input))
		d.NoHeader = v.NoHeader
		d.Encoding = v.Encoding

		result, err := d.Delimit()

		if err != nil {
			if v.Error == "" {
				t.Errorf("unexpected error %q for %q", err.Error(), v.Input)
			} else if v.Error != err.Error() {
				t.Errorf("error %q, want error %q for %q", err.Error(), v.Error, v.Input)
			}
			continue
		}
		if 0 < len(v.Error) {
			t.Errorf("no error, want error %q for %q", v.Error, v.Input)
			continue
		}
		if !reflect.DeepEqual(result, v.Expect) {
			t.Errorf("result = %v, want %v for %q", result, v.Expect, v.Input)
		}
	}
}
