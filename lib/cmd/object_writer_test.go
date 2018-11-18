package cmd

import (
	"testing"
)

func TestObjectWriter_String(t *testing.T) {
	w := NewObjectWriter()
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
	w.Write("bbbbbbbbbbbbbbbbbbbbbbbbb")
	w.WriteWithoutLineBreak(", ")
	w.ClearBlock()
	w.NewLine()
	w.Write("aaa")
	w.BeginBlock()
	w.NewLine()
	w.Write("key: ")
	w.BeginSubBlock()
	w.Write("bbbbbbbb")
	w.WriteWithoutLineBreak(", ")
	w.Write("bbbbbbbb")
	w.EndSubBlock()
	w.NewLine()
	w.Write("bbbbbbbb")

	expect := "" +
		" aaa\n" +
		"     bbb    bbb\n" +
		"         ccc\n" +
		"             ddd\n" +
		"         ccc\n" +
		" aaa\n" +
		"     bbbbbbbbbb, \n" +
		"     bbbbbbbbbb, \n" +
		"     bbbbbbbbbbbbbbbbbbbbbbbbb, \n" +
		" aaa\n" +
		"     key: bbbbbbbb, \n" +
		"          bbbbbbbb\n" +
		"     bbbbbbbb"
	result := w.String()

	if result != expect {
		t.Errorf("result = %q, want %q", result, expect)
	}

	w = NewObjectWriter()
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

	expect = "" +
		"       title\n" +
		"--------------------\n" +
		" aaa\n" +
		"     bbbbbbbbbb, \n" +
		"     bbbbbbbbbb, \n" +
		"     bbbbbbbbbbbbbbbbbbbbbbbbb, "
	result = w.String()

	if result != expect {
		t.Errorf("result = %q, want %q", result, expect)
	}

	w = NewObjectWriter()
	w.MaxWidth = 20

	w.Title1 = "title"

	w.Write("aaa")

	expect = "" +
		" title\n" +
		"-------\n" +
		" aaa"
	result = w.String()

	if result != expect {
		t.Errorf("result = %q, want %q", result, expect)
	}

	GetFlags().SetColor(true)
	w = NewObjectWriter()
	w.MaxWidth = 20

	w.Title1 = "title1"
	w.Title2 = "title2"
	w.Title2Effect = IdentifierEffect

	w.Write("aaa")
	w.BeginBlock()
	w.NewLine()
	w.WriteColor("bbbbbbbbbb", StringEffect)
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
		t.Errorf("result = %q, want %q", result, expect)
	}

	GetFlags().SetColor(false)
}
