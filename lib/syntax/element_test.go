package syntax

import (
	"context"
	"testing"

	"github.com/mithrandie/csvq/lib/file"
	"github.com/mithrandie/csvq/lib/option"

	"github.com/mithrandie/go-text/color"
)

var syntaxTestEnv, _ = option.NewEnvironment(context.Background(), file.DefaultWaitTimeout, file.DefaultRetryDelay)
var syntaxTestPalette, _ = option.NewPalette(syntaxTestEnv)

func TestName_Format(t *testing.T) {
	var e Name = "str"

	expect := "str"
	result := e.Format(nil)
	if result != expect {
		t.Errorf("result = %q, want %q", result, expect)
	}

	expect = syntaxTestPalette.Render(NameEffect, "str")
	result = e.Format(syntaxTestPalette)
	if result != expect {
		t.Errorf("result = %q, want %q", result, expect)
	}
}

var grammarFormatTests = []struct {
	Grammar    Grammar
	UsePalette bool
	Expect     string
}{
	{
		Grammar:    []Element{Keyword("KEY"), Link("link")},
		UsePalette: false,
		Expect:     "KEY <link>",
	},
	{
		Grammar:    []Element{Keyword("KEY"), Link("link")},
		UsePalette: true,
		Expect:     syntaxTestPalette.Render(KeywordEffect, "KEY") + " " + syntaxTestPalette.Render(LinkEffect, "<link>"),
	},
	{
		Grammar:    []Element{Option{String("str1"), String("str2")}},
		UsePalette: false,
		Expect:     "[str1 str2]",
	},
	{
		Grammar:    []Element{Option{String("str1"), String("str2")}},
		UsePalette: true,
		Expect:     "[" + syntaxTestPalette.Render(ItalicEffect, syntaxTestPalette.Render(option.StringEffect, "str1")) + " " + syntaxTestPalette.Render(ItalicEffect, syntaxTestPalette.Render(option.StringEffect, "str2")) + "]",
	},
	{
		Grammar:    []Element{FollowingContinuousOption{String("str1"), String("str2")}},
		UsePalette: false,
		Expect:     " [, str1 str2 ...]",
	},
	{
		Grammar:    []Element{ContinuousOption{String("str1"), String("str2")}},
		UsePalette: false,
		Expect:     "str1 str2 [, str1 str2 ...]",
	},
	{
		Grammar:    []Element{ContinuousOption{String("str1"), String("str2")}},
		UsePalette: true,
		Expect:     syntaxTestPalette.Render(ItalicEffect, syntaxTestPalette.Render(option.StringEffect, "str1")) + " " + syntaxTestPalette.Render(ItalicEffect, syntaxTestPalette.Render(option.StringEffect, "str2")) + " [, " + syntaxTestPalette.Render(ItalicEffect, syntaxTestPalette.Render(option.StringEffect, "str1")) + " " + syntaxTestPalette.Render(ItalicEffect, syntaxTestPalette.Render(option.StringEffect, "str2")) + " ...]",
	},
	{
		Grammar:    []Element{AnyOne{String("str1"), String("str2")}},
		UsePalette: false,
		Expect:     "{str1|str2}",
	},
	{
		Grammar:    []Element{AnyOne{String("str1"), String("str2")}},
		UsePalette: true,
		Expect:     "{" + syntaxTestPalette.Render(ItalicEffect, syntaxTestPalette.Render(option.StringEffect, "str1")) + "|" + syntaxTestPalette.Render(ItalicEffect, syntaxTestPalette.Render(option.StringEffect, "str2")) + "}",
	},
	{
		Grammar:    []Element{Parentheses{String("str1"), String("str2")}},
		UsePalette: false,
		Expect:     "(str1 str2)",
	},
	{
		Grammar:    []Element{Parentheses{String("str1"), String("str2")}},
		UsePalette: true,
		Expect:     "(" + syntaxTestPalette.Render(ItalicEffect, syntaxTestPalette.Render(option.StringEffect, "str1")) + " " + syntaxTestPalette.Render(ItalicEffect, syntaxTestPalette.Render(option.StringEffect, "str2")) + ")",
	},
	{
		Grammar:    []Element{PlainGroup{String("str1"), String("str2")}},
		UsePalette: false,
		Expect:     "str1 str2",
	},
	{
		Grammar:    []Element{PlainGroup{String("str1"), String("str2")}},
		UsePalette: true,
		Expect:     syntaxTestPalette.Render(ItalicEffect, syntaxTestPalette.Render(option.StringEffect, "str1")) + " " + syntaxTestPalette.Render(ItalicEffect, syntaxTestPalette.Render(option.StringEffect, "str2")),
	},
	{
		Grammar:    []Element{ConnectedGroup{String("str1"), String("str2")}},
		UsePalette: false,
		Expect:     "str1str2",
	},
	{
		Grammar:    []Element{ConnectedGroup{String("str1"), String("str2")}},
		UsePalette: true,
		Expect:     syntaxTestPalette.Render(ItalicEffect, syntaxTestPalette.Render(option.StringEffect, "str1")) + syntaxTestPalette.Render(ItalicEffect, syntaxTestPalette.Render(option.StringEffect, "str2")),
	},
	{
		Grammar:    []Element{Function{Name: "fn", Args: []Element{Option{Keyword("DISTINCT")}, String("str"), Option{String("arg1"), String("arg2")}}, AfterArgs: []Element{Keyword("OVER"), Parentheses{Link("partition_clause")}}, Return: Return("string")}},
		UsePalette: false,
		Expect:     "fn([DISTINCT] str::string [, arg1::string [, arg2::string]]) OVER (<partition_clause>)  return::string",
	},
	{
		Grammar:    []Element{Function{Name: "fn", Args: []Element{Option{Keyword("DISTINCT")}, String("str"), Option{String("arg1"), String("arg2")}}, AfterArgs: []Element{Keyword("OVER"), Parentheses{Link("partition_clause")}}, Return: Return("string")}},
		UsePalette: true,
		Expect: syntaxTestPalette.Render(KeywordEffect, "fn") +
			"([" + syntaxTestPalette.Render(KeywordEffect, "DISTINCT") + "] " +
			syntaxTestPalette.Render(ItalicEffect, syntaxTestPalette.Render(option.StringEffect, "str")) + syntaxTestPalette.Render(TypeEffect, "::string") + " [, " +
			syntaxTestPalette.Render(ItalicEffect, syntaxTestPalette.Render(option.StringEffect, "arg1")) + syntaxTestPalette.Render(TypeEffect, "::string") + " [, " +
			syntaxTestPalette.Render(ItalicEffect, syntaxTestPalette.Render(option.StringEffect, "arg2")) + syntaxTestPalette.Render(TypeEffect, "::string") + "]]) " +
			syntaxTestPalette.Render(KeywordEffect, "OVER") + " (" +
			syntaxTestPalette.Render(LinkEffect, "<partition_clause>") + ")  " +
			syntaxTestPalette.Render(TypeEffect, "return::string"),
	},
	{
		Grammar:    []Element{Function{Name: "fn", Args: []Element{String("str"), Integer("int"), Float("float"), Boolean("bool"), Ternary("ternary"), Datetime("dt"), Link("link")}}},
		UsePalette: false,
		Expect:     "fn(str::string, int::integer, float::float, bool::boolean, ternary::ternary, dt::datetime, <link>)",
	},
	{
		Grammar:    []Element{Function{Name: "fn", CustomArgs: Grammar{String("str"), Keyword("FROM"), Integer("num")}}},
		UsePalette: false,
		Expect:     "fn(str::string FROM num::integer)",
	},
	{
		Grammar:    []Element{ArgWithDefValue{Arg: String("str"), Default: String("1")}},
		UsePalette: false,
		Expect:     "str::string DEFAULT 1",
	},
	{
		Grammar:    []Element{ArgWithDefValue{Arg: String("str"), Default: String("1")}},
		UsePalette: true,
		Expect:     syntaxTestPalette.Render(ItalicEffect, syntaxTestPalette.Render(option.StringEffect, "str")) + syntaxTestPalette.Render(TypeEffect, "::string") + " " + syntaxTestPalette.Render(KeywordEffect, "DEFAULT") + " " + syntaxTestPalette.Render(ItalicEffect, syntaxTestPalette.Render(option.StringEffect, "1")),
	},
	{
		Grammar:    []Element{String("str")},
		UsePalette: false,
		Expect:     "str",
	},
	{
		Grammar:    []Element{String("str")},
		UsePalette: true,
		Expect:     syntaxTestPalette.Render(ItalicEffect, syntaxTestPalette.Render(option.StringEffect, "str")),
	},
	{
		Grammar:    []Element{Integer("int")},
		UsePalette: false,
		Expect:     "int",
	},
	{
		Grammar:    []Element{Integer("int")},
		UsePalette: true,
		Expect:     syntaxTestPalette.Render(ItalicEffect, syntaxTestPalette.Render(option.NumberEffect, "int")),
	},
	{
		Grammar:    []Element{Float("float")},
		UsePalette: false,
		Expect:     "float",
	},
	{
		Grammar:    []Element{Float("float")},
		UsePalette: true,
		Expect:     syntaxTestPalette.Render(ItalicEffect, syntaxTestPalette.Render(option.NumberEffect, "float")),
	},
	{
		Grammar:    []Element{Identifier("ident")},
		UsePalette: false,
		Expect:     "ident",
	},
	{
		Grammar:    []Element{Identifier("ident")},
		UsePalette: true,
		Expect:     syntaxTestPalette.Render(option.IdentifierEffect, "ident"),
	},
	{
		Grammar:    []Element{Datetime("dt")},
		UsePalette: false,
		Expect:     "dt",
	},
	{
		Grammar:    []Element{Datetime("dt")},
		UsePalette: true,
		Expect:     syntaxTestPalette.Render(ItalicEffect, syntaxTestPalette.Render(option.DatetimeEffect, "dt")),
	},
	{
		Grammar:    []Element{Boolean("bool")},
		UsePalette: false,
		Expect:     "bool",
	},
	{
		Grammar:    []Element{Boolean("bool")},
		UsePalette: true,
		Expect:     syntaxTestPalette.Render(ItalicEffect, syntaxTestPalette.Render(option.BooleanEffect, "bool")),
	},
	{
		Grammar:    []Element{Ternary("ternary")},
		UsePalette: false,
		Expect:     "ternary",
	},
	{
		Grammar:    []Element{Ternary("ternary")},
		UsePalette: true,
		Expect:     syntaxTestPalette.Render(ItalicEffect, syntaxTestPalette.Render(option.TernaryEffect, "ternary")),
	},
	{
		Grammar:    []Element{Null("null")},
		UsePalette: false,
		Expect:     "null",
	},
	{
		Grammar:    []Element{Null("null")},
		UsePalette: true,
		Expect:     syntaxTestPalette.Render(ItalicEffect, syntaxTestPalette.Render(option.NullEffect, "null")),
	},
	{
		Grammar:    []Element{Variable("@var")},
		UsePalette: false,
		Expect:     "@var",
	},
	{
		Grammar:    []Element{Variable("@var")},
		UsePalette: true,
		Expect:     syntaxTestPalette.Render(VariableEffect, "@var"),
	},
	{
		Grammar:    []Element{Flag("@@flag")},
		UsePalette: false,
		Expect:     "@@flag",
	},
	{
		Grammar:    []Element{Flag("@@flag")},
		UsePalette: true,
		Expect:     syntaxTestPalette.Render(FlagEffect, "@@flag"),
	},
	{
		Grammar:    []Element{Token(".")},
		UsePalette: false,
		Expect:     ".",
	},
	{
		Grammar:    []Element{Token(".")},
		UsePalette: true,
		Expect:     ".",
	},
	{
		Grammar:    []Element{Italic("str")},
		UsePalette: false,
		Expect:     "str",
	},
	{
		Grammar:    []Element{Italic("str")},
		UsePalette: true,
		Expect:     syntaxTestPalette.Render(ItalicEffect, "str"),
	},
	{
		Grammar:    []Element{Return("string")},
		UsePalette: false,
		Expect:     "  return::string",
	},
	{
		Grammar:    []Element{Return("string")},
		UsePalette: true,
		Expect:     "  " + syntaxTestPalette.Render(TypeEffect, "return::string"),
	},
	{
		Grammar: []Element{Description{
			Template: "abc %s %s %s %s %s %s",
			Values:   []Element{String("str"), Ternary("TRUE"), Ternary("ternary"), Null("NULL"), Null("null"), Link("link")},
		}},
		UsePalette: false,
		Expect:     "abc _str_ TRUE _ternary_ NULL _null_ <link>",
	},
	{
		Grammar: []Element{Description{
			Template: "abc %s %s %s %s %s %s",
			Values:   []Element{String("str"), Ternary("TRUE"), Ternary("ternary"), Null("NULL"), Null("null"), Link("link")},
		}},
		UsePalette: true,
		Expect: "abc " +
			syntaxTestPalette.Render(ItalicEffect, syntaxTestPalette.Render(option.StringEffect, "str")) + " " +
			syntaxTestPalette.Render(ItalicEffect, syntaxTestPalette.Render(option.TernaryEffect, "TRUE")) + " " +
			syntaxTestPalette.Render(ItalicEffect, syntaxTestPalette.Render(option.TernaryEffect, "ternary")) + " " +
			syntaxTestPalette.Render(ItalicEffect, syntaxTestPalette.Render(option.NullEffect, "NULL")) + " " +
			syntaxTestPalette.Render(ItalicEffect, syntaxTestPalette.Render(option.NullEffect, "null")) + " " +
			syntaxTestPalette.Render(LinkEffect, "<link>"),
	},
}

func TestGrammar_Format(t *testing.T) {
	var palette *color.Palette

	for _, v := range grammarFormatTests {
		if v.UsePalette {
			palette = syntaxTestPalette
		} else {
			palette = nil
		}
		result := v.Grammar.Format(palette)
		if result != v.Expect {
			t.Errorf("result = %q, want %q for %v", result, v.Expect, v.Grammar)
		}
	}
}
