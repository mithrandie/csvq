package query

import (
	"testing"

	"github.com/mithrandie/csvq/lib/option"
)

func TestObjectWriter_String(t *testing.T) {
	defer initFlag(TestTx.Flags)

	w := NewObjectWriter(TestTx)
	w.MaxWidth = 20

	w.Write("aaa")
	w.BeginBlock()
	w.NewLine()
	w.Write("bbb")
	w.WriteSpaces(4)
	w.Write("bbb")
	w.BeginBlock()
	w.NewLine()
	w.Write("ccc")
	w.BeginBlock()
	w.NewLine()
	w.Write("ddd")
	w.EndBlock()
	w.NewLine()
	w.Write("ccc")
	w.ClearBlock()
	w.NewLine()
	w.Write("aaa")
	w.BeginBlock()
	w.NewLine()
	w.Write("bbbbbbbbbb")
	w.Write(", ")
	w.Write("bbbbbbbbbb")
	w.Write(", ")
	w.Write("bbbbbbbbbbbbbbbbbbbbbbbb")
	w.WriteWithoutLineBreak(", ")
	w.ClearBlock()
	w.NewLine()
	w.Write("aaa")
	w.BeginBlock()
	w.NewLine()
	w.Write("key: ")
	w.BeginSubBlock()
	w.Write("bbbbbbb")
	w.WriteWithoutLineBreak(", ")
	w.Write("bbbbbbb")
	w.EndSubBlock()
	w.NewLine()
	w.Write("bbbbbbb")

	expect := "" +
		" aaa\n" +
		"     bbb    bbb\n" +
		"         ccc\n" +
		"             ddd\n" +
		"         ccc\n" +
		" aaa\n" +
		"     bbbbbbbbbb, \n" +
		"     bbbbbbbbbb, \n" +
		"     bbbbbbbbbbbbbbbbbbbbbbbb, \n" +
		" aaa\n" +
		"     key: bbbbbbb, \n" +
		"          bbbbbbb\n" +
		"     bbbbbbb"
	result := w.String()

	if result != expect {
		t.Errorf("result = %s, want %s", result, expect)
	}

	w = NewObjectWriter(TestTx)
	w.MaxWidth = 20

	w.Title1 = "title"

	w.Write("aaa")
	w.BeginBlock()
	w.NewLine()
	w.Write("bbbbbbbbbb")
	w.Write(", ")
	w.Write("bbbbbbbbbb")
	w.Write(", ")
	w.Write("bbbbbbbbbbbbbbbbbbbbbbbbb")
	w.WriteWithoutLineBreak(", ")
	w.NewLine()
	w.WriteWithAutoLineBreak("aaaaa bbbbb ccccc\n > ddddd \n eeeee")
	w.NewLine()
	w.WriteWithAutoLineBreak("```\naaaaa     bbbbb\n```\nccccc")

	expect = "" +
		"       title\n" +
		"--------------------\n" +
		" aaa\n" +
		"     bbbbbbbbbb, \n" +
		"     bbbbbbbbbb, \n" +
		"     bbbbbbbbbbbbbbbbbbbbbbbbb, \n" +
		"     aaaaa bbbbb\n" +
		"     ccccc\n" +
		"         ddddd\n" +
		"     eeeee\n" +
		"     aaaaa     bbbbb\n" +
		"     ccccc" +
		""
	result = w.String()

	if result != expect {
		t.Errorf("result = %s, want %s", result, expect)
	}

	w = NewObjectWriter(TestTx)
	w.MaxWidth = 20

	w.Title1 = "title"

	w.Write("aaa")

	expect = "" +
		" title\n" +
		"-------\n" +
		" aaa"
	result = w.String()

	if result != expect {
		t.Errorf("result = %s, want %s", result, expect)
	}

	TestTx.UseColor(true)
	defer TestTx.UseColor(false)

	w = NewObjectWriter(TestTx)
	w.MaxWidth = 20

	w.Title1 = "title1"
	w.Title2 = "title2"
	w.Title2Effect = option.IdentifierEffect

	w.Write("aaa")
	w.BeginBlock()
	w.NewLine()
	w.WriteColor("bbbbbbbbbb", option.StringEffect)
	w.Write(", ")
	w.Write("bbbbbbbbbb")

	expect = "" +
		"  title1 \x1b[36;1mtitle2\x1b[0m\n" +
		"------------------\n" +
		" aaa\n" +
		"     \x1b[32mbbbbbbbbbb\x1b[0m, \n" +
		"     bbbbbbbbbb"
	result = w.String()

	if result != expect {
		t.Errorf("result = %s, want %s", result, expect)
	}
}
