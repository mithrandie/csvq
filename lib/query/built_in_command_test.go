package query

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/parser"
)

var printTests = []struct {
	Name   string
	Expr   parser.Print
	Result string
	Error  string
}{
	{
		Name: "Print",
		Expr: parser.Print{
			Value: parser.NewString("foo"),
		},
		Result: "'foo'",
	},
	{
		Name: "Print Error",
		Expr: parser.Print{
			Value: parser.Variable{
				Name: "var",
			},
		},
		Error: "variable var is undefined",
	},
}

func TestPrint(t *testing.T) {
	filter := NewFilter([]Variables{{}})

	for _, v := range printTests {
		result, err := Print(v.Expr, filter)
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
		if result != v.Result {
			t.Errorf("%s: result = %q, want %q", v.Name, result, v.Result)
		}
	}
}

var printfTests = []struct {
	Name   string
	Expr   parser.Printf
	Result string
	Error  string
}{
	{
		Name: "Printf",
		Expr: parser.Printf{
			Values: []parser.Expression{
				parser.NewString("printf test: value1 %s, value2 %s, %a %% %"),
				parser.NewString("str"),
				parser.NewInteger(1),
			},
		},
		Result: "printf test: value1 'str', value2 1, %a % %",
	},
	{
		Name: "Printf Less Values Error",
		Expr: parser.Printf{
			Values: []parser.Expression{
				parser.NewString("printf test: value1 %s, value2 %s, %a %% %"),
				parser.NewString("str"),
			},
		},
		Error: "print format \"printf test: value1 %s, value2 %s, %a %% %\": number of replace values does not match",
	},
	{
		Name: "Printf Greater Values Error",
		Expr: parser.Printf{
			Values: []parser.Expression{
				parser.NewString("printf test: value1 %s, value2 %s, %a %% %"),
				parser.NewString("str"),
				parser.NewInteger(1),
				parser.NewInteger(2),
			},
		},
		Error: "print format \"printf test: value1 %s, value2 %s, %a %% %\": number of replace values does not match",
	},
}

func TestPrintf(t *testing.T) {
	filter := NewFilter([]Variables{{}})

	for _, v := range printfTests {
		result, err := Printf(v.Expr, filter)
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
		if result != v.Result {
			t.Errorf("%s: result = %q, want %q", v.Name, result, v.Result)
		}
	}
}

var sourceTests = []struct {
	Name   string
	Expr   parser.Source
	Result []parser.Statement
	Error  string
}{
	{
		Name: "Source",
		Expr: parser.Source{
			FilePath: GetTestFilePath("source.sql"),
		},
		Result: []parser.Statement{
			parser.Print{
				Value: parser.NewString("external executable file"),
			},
		},
	},
	{
		Name: "Source File Not Exist Error",
		Expr: parser.Source{
			FilePath: GetTestFilePath("notexist.sql"),
		},
		Error: fmt.Sprintf("source file %q does not exist", GetTestFilePath("notexist.sql")),
	},
}

func TestSource(t *testing.T) {
	for _, v := range sourceTests {
		result, err := Source(v.Expr)
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
		if !reflect.DeepEqual(result, v.Result) {
			t.Errorf("%s: result = %q, want %q", v.Name, result, v.Result)
		}
	}
}

var setFlagTests = []struct {
	Name            string
	Expr            parser.SetFlag
	ResultFlag      string
	ResultStlValue  string
	ResultBoolValue bool
	Error           string
}{
	{
		Name: "Set Delimiter",
		Expr: parser.SetFlag{
			Name:  "@@delimiter",
			Value: parser.NewString("\t"),
		},
		ResultFlag:     "delimiter",
		ResultStlValue: "\t",
	},
	{
		Name: "Set Encoding",
		Expr: parser.SetFlag{
			Name:  "@@encoding",
			Value: parser.NewString("SJIS"),
		},
		ResultFlag:     "encoding",
		ResultStlValue: "SJIS",
	},
	{
		Name: "Set LineBreak",
		Expr: parser.SetFlag{
			Name:  "@@line_break",
			Value: parser.NewString("CRLF"),
		},
		ResultFlag:     "line_break",
		ResultStlValue: "\r\n",
	},
	{
		Name: "Set Repository",
		Expr: parser.SetFlag{
			Name:  "@@repository",
			Value: parser.NewString(TestDir),
		},
		ResultFlag:     "repository",
		ResultStlValue: TestDir,
	},
	{
		Name: "Set DatetimeFormat",
		Expr: parser.SetFlag{
			Name:  "@@datetime_format",
			Value: parser.NewString("%Y%m%d"),
		},
		ResultFlag:     "datetime_format",
		ResultStlValue: "%Y%m%d",
	},
	{
		Name: "Set NoHeader",
		Expr: parser.SetFlag{
			Name:  "@@no_header",
			Value: parser.NewBoolean(true),
		},
		ResultFlag:      "no_header",
		ResultBoolValue: true,
	},
	{
		Name: "Set WithoutNull",
		Expr: parser.SetFlag{
			Name:  "@@without_null",
			Value: parser.NewBoolean(true),
		},
		ResultFlag:      "without_null",
		ResultBoolValue: true,
	},
	{
		Name: "Set Delimiter Value Error",
		Expr: parser.SetFlag{
			Name:  "@@delimiter",
			Value: parser.NewBoolean(true),
		},
		Error: "invalid flag value: @@delimiter = true",
	},
	{
		Name: "Set WithoutNull Value Error",
		Expr: parser.SetFlag{
			Name:  "@@without_null",
			Value: parser.NewString("string"),
		},
		Error: "invalid flag value: @@without_null = 'string'",
	},
	{
		Name: "Invalid Flag Error",
		Expr: parser.SetFlag{
			Name:  "@@invalid",
			Value: parser.NewString("string"),
		},
		Error: "invalid flag name: @@invalid",
	},
}

func TestSetFlag(t *testing.T) {
	initFlag()
	flags := cmd.GetFlags()

	for _, v := range setFlagTests {
		err := SetFlag(v.Expr)
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

		switch strings.ToUpper(v.ResultFlag) {
		case "DELIMITER":
			if string(flags.Delimiter) != v.ResultStlValue {
				t.Errorf("%s: delimiter = %q, want %q", v.Name, string(flags.Delimiter), v.ResultStlValue)
			}
		case "ENCODING":
			if flags.Encoding.String() != v.ResultStlValue {
				t.Errorf("%s: encoding = %q, want %q", v.Name, flags.Encoding.String(), v.ResultStlValue)
			}
		case "LINE_BREAK":
			if flags.LineBreak.Value() != v.ResultStlValue {
				t.Errorf("%s: line-break = %q, want %q", v.Name, flags.LineBreak.Value(), v.ResultStlValue)
			}
		case "REPOSITORY":
			if flags.Repository != v.ResultStlValue {
				t.Errorf("%s: repository = %q, want %q", v.Name, flags.Repository, v.ResultStlValue)
			}
		case "DATETIME_FORMAT":
			if flags.DatetimeFormat != v.ResultStlValue {
				t.Errorf("%s: datetime-format = %q, want %q", v.Name, flags.DatetimeFormat, v.ResultStlValue)
			}
		case "NO-HEADER":
			if flags.NoHeader != v.ResultBoolValue {
				t.Errorf("%s: no-header = %t, want %t", v.Name, flags.NoHeader, v.ResultBoolValue)
			}
		case "WITHOUT-NULL":
			if flags.WithoutNull != v.ResultBoolValue {
				t.Errorf("%s: without-null = %t, want %t", v.Name, flags.WithoutNull, v.ResultBoolValue)
			}
		}
	}
}
