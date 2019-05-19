package excmd

import (
	"reflect"
	"testing"
)

type argumentScannerResult struct {
	Text        string
	ElementType ElementType
}

var argumentScannerScanTests = []struct {
	Input  string
	Expect []argumentScannerResult
	Error  string
}{
	{
		Input:  "",
		Expect: []argumentScannerResult(nil),
	},
	{
		Input: "arg",
		Expect: []argumentScannerResult{
			{Text: "arg", ElementType: FixedString},
		},
	},
	{
		Input: "\\$",
		Expect: []argumentScannerResult{
			{Text: "$", ElementType: FixedString},
		},
	},
	{
		Input: "\\@",
		Expect: []argumentScannerResult{
			{Text: "@", ElementType: FixedString},
		},
	},
	{
		Input: "arg\\@arg\\\\\\$arg\\arg",
		Expect: []argumentScannerResult{
			{Text: "arg@arg\\\\$arg\\arg", ElementType: FixedString},
		},
	},
	{
		Input: "@var",
		Expect: []argumentScannerResult{
			{Text: "var", ElementType: Variable},
		},
	},
	{
		Input: "@%var",
		Expect: []argumentScannerResult{
			{Text: "var", ElementType: EnvironmentVariable},
		},
	},
	{
		Input: "@%`var\\\\var\\`var`",
		Expect: []argumentScannerResult{
			{Text: "var\\var`var", ElementType: EnvironmentVariable},
		},
	},
	{
		Input: "@#var",
		Expect: []argumentScannerResult{
			{Text: "var", ElementType: RuntimeInformation},
		},
	},
	{
		Input: "${print @a}",
		Expect: []argumentScannerResult{
			{Text: "print @a", ElementType: CsvqExpression},
		},
	},
	{
		Input: "${print 'a\\{bc\\\\\\}de'}",
		Expect: []argumentScannerResult{
			{Text: "print 'a{bc\\\\}de'", ElementType: CsvqExpression},
		},
	},
	{
		Input: "cmd --option arg1 'arg 2' arg3",
		Expect: []argumentScannerResult{
			{Text: "cmd --option arg1 'arg 2' arg3", ElementType: FixedString},
		},
	},
	{
		Input: "arg${print @a}arg",
		Expect: []argumentScannerResult{
			{Text: "arg", ElementType: FixedString},
			{Text: "print @a", ElementType: CsvqExpression},
			{Text: "arg", ElementType: FixedString},
		},
	},
	{
		Input: "@%`var",
		Error: "environment variable not terminated",
	},
	{
		Input: "arg@%",
		Error: "invalid variable symbol",
	},
	{
		Input: "arg$arg",
		Error: "invalid command symbol",
	},
	{
		Input: "arg${print @a",
		Error: "command not terminated",
	},
}

func TestArgumentScanner_Scan(t *testing.T) {
	for _, v := range argumentScannerScanTests {
		scanner := new(ArgumentScanner).Init(v.Input)
		var args []argumentScannerResult
		for scanner.Scan() {
			args = append(args, argumentScannerResult{Text: scanner.Text(), ElementType: scanner.ElementType()})
		}
		err := scanner.Err()

		if err != nil {
			if v.Error == "" {
				t.Errorf("unexpected error %q for %q", err.Error(), v.Input)
			} else if v.Error != err.Error() {
				t.Errorf("error %q, want error %q for %q", err.Error(), v.Error, v.Input)
			}
			continue
		}
		if v.Error != "" {
			t.Errorf("no error, want error %q for %q", v.Error, v.Input)
			continue
		}

		if !reflect.DeepEqual(args, v.Expect) {
			t.Errorf("result = %#v, want %#v for %q", args, v.Expect, v.Input)
		}
	}
}
