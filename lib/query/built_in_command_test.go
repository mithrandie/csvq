package query

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/value"
	"github.com/mithrandie/go-file"
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
			Value: parser.NewStringValue("foo"),
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
		Error: "[L:- C:-] variable var is undefined",
	},
}

func TestPrint(t *testing.T) {
	filter := NewEmptyFilter()

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
			Format: parser.NewStringValue("printf test: value1 %q, value2 %q, %a %% %"),
			Values: []parser.QueryExpression{
				parser.NewStringValue("str"),
				parser.NewIntegerValue(1),
			},
		},
		Result: "printf test: value1 'str', value2 1, %a % %",
	},
	{
		Name: "Printf Format Error",
		Expr: parser.Printf{
			Format: parser.Variable{Name: "var"},
			Values: []parser.QueryExpression{
				parser.NewStringValue("str"),
				parser.NewIntegerValue(1),
			},
		},
		Error: "[L:- C:-] variable var is undefined",
	},
	{
		Name: "Printf Evaluate Error",
		Expr: parser.Printf{
			Format: parser.NewStringValue("printf test: value1 %s"),
			Values: []parser.QueryExpression{
				parser.Variable{
					Name: "var",
				},
			},
		},
		Error: "[L:- C:-] variable var is undefined",
	},
	{
		Name: "Printf Less Values Error",
		Expr: parser.Printf{
			Format: parser.NewStringValue("printf test: value1 %s, value2 %s, %a %% %"),
			Values: []parser.QueryExpression{
				parser.NewStringValue("str"),
			},
		},
		Error: "[L:- C:-] PRINTF: number of replace values does not match",
	},
	{
		Name: "Printf Greater Values Error",
		Expr: parser.Printf{
			Format: parser.NewStringValue("printf test: value1 %s, value2 %s, %a %% %"),
			Values: []parser.QueryExpression{
				parser.NewStringValue("str"),
				parser.NewIntegerValue(1),
				parser.NewIntegerValue(2),
			},
		},
		Error: "[L:- C:-] PRINTF: number of replace values does not match",
	},
}

func TestPrintf(t *testing.T) {
	filter := NewEmptyFilter()

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
			FilePath: parser.NewStringValue(GetTestFilePath("source.sql")),
		},
		Result: []parser.Statement{
			parser.Print{
				Value: parser.NewStringValue("external executable file"),
			},
		},
	},
	{
		Name: "Source File Argument Evaluation Error",
		Expr: parser.Source{
			FilePath: parser.FieldReference{Column: parser.Identifier{Literal: "ident"}},
		},
		Error: "[L:- C:-] field ident does not exist",
	},
	{
		Name: "Source File Argument Not String Error",
		Expr: parser.Source{
			FilePath: parser.NewNullValueFromString("NULL"),
		},
		Error: "[L:- C:-] SOURCE: argument NULL is not a string",
	},
	{
		Name: "Source File Not Exist Error",
		Expr: parser.Source{
			FilePath: parser.NewStringValue(GetTestFilePath("notexist.sql")),
		},
		Error: fmt.Sprintf("[L:- C:-] SOURCE: file %s does not exist", GetTestFilePath("notexist.sql")),
	},
	{
		Name: "Source File Not Readable Error",
		Expr: parser.Source{
			FilePath: parser.NewStringValue(TestDir),
		},
		Error: fmt.Sprintf("[L:- C:-] SOURCE: file %s is unable to read", TestDir),
	},
	{
		Name: "Source Syntax Error",
		Expr: parser.Source{
			FilePath: parser.NewStringValue(GetTestFilePath("source_syntaxerror.sql")),
		},
		Error: fmt.Sprintf("%s [L:1 C:34] syntax error: unexpected STRING", GetTestFilePath("source_syntaxerror.sql")),
	},
}

func TestSource(t *testing.T) {
	filter := NewEmptyFilter()

	for _, v := range sourceTests {
		result, err := Source(v.Expr, filter)
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
	Name             string
	Expr             parser.SetFlag
	ResultFlag       string
	ResultStrValue   string
	ResultFloatValue float64
	ResultBoolValue  bool
	Error            string
}{
	{
		Name: "Set Delimiter",
		Expr: parser.SetFlag{
			Name:  "@@delimiter",
			Value: value.NewString("\t"),
		},
		ResultFlag:     "delimiter",
		ResultStrValue: "\t",
	},
	{
		Name: "Set Encoding",
		Expr: parser.SetFlag{
			Name:  "@@encoding",
			Value: value.NewString("SJIS"),
		},
		ResultFlag:     "encoding",
		ResultStrValue: "SJIS",
	},
	{
		Name: "Set LineBreak",
		Expr: parser.SetFlag{
			Name:  "@@line_break",
			Value: value.NewString("CRLF"),
		},
		ResultFlag:     "line_break",
		ResultStrValue: "\r\n",
	},
	{
		Name: "Set Repository",
		Expr: parser.SetFlag{
			Name:  "@@repository",
			Value: value.NewString(TestDir),
		},
		ResultFlag:     "repository",
		ResultStrValue: TestDir,
	},
	{
		Name: "Set DatetimeFormat",
		Expr: parser.SetFlag{
			Name:  "@@datetime_format",
			Value: value.NewString("%Y%m%d"),
		},
		ResultFlag:     "datetime_format",
		ResultStrValue: "%Y%m%d",
	},
	{
		Name: "Set WaitTimeout",
		Expr: parser.SetFlag{
			Name:  "@@wait_timeout",
			Value: value.NewFloat(15),
		},
		ResultFlag:       "wait_timeout",
		ResultFloatValue: 15,
	},
	{
		Name: "Set NoHeader",
		Expr: parser.SetFlag{
			Name:  "@@no_header",
			Value: value.NewBoolean(true),
		},
		ResultFlag:      "no_header",
		ResultBoolValue: true,
	},
	{
		Name: "Set WithoutNull",
		Expr: parser.SetFlag{
			Name:  "@@without_null",
			Value: value.NewBoolean(true),
		},
		ResultFlag:      "without_null",
		ResultBoolValue: true,
	},
	{
		Name: "Set Stats",
		Expr: parser.SetFlag{
			Name:  "@@stats",
			Value: value.NewBoolean(true),
		},
		ResultFlag:      "stats",
		ResultBoolValue: true,
	},
	{
		Name: "Set Delimiter Value Error",
		Expr: parser.SetFlag{
			Name:  "@@delimiter",
			Value: value.NewBoolean(true),
		},
		Error: "[L:- C:-] SET: flag value true for @@delimiter is invalid",
	},
	{
		Name: "Set WaitTimeout Value Error",
		Expr: parser.SetFlag{
			Name:  "@@wait_timeout",
			Value: value.NewBoolean(true),
		},
		Error: "[L:- C:-] SET: flag value true for @@wait_timeout is invalid",
	},
	{
		Name: "Set WithoutNull Value Error",
		Expr: parser.SetFlag{
			Name:  "@@without_null",
			Value: value.NewString("string"),
		},
		Error: "[L:- C:-] SET: flag value 'string' for @@without_null is invalid",
	},
	{
		Name: "Invalid Flag Name Error",
		Expr: parser.SetFlag{
			Name:  "@@invalid",
			Value: value.NewString("string"),
		},
		Error: "[L:- C:-] flag name @@invalid is invalid",
	},
	{
		Name: "Invalid Flag Value Error",
		Expr: parser.SetFlag{
			Name:  "@@line_break",
			Value: value.NewString("invalid"),
		},
		Error: "[L:- C:-] SET: flag value 'invalid' for @@line_break is invalid",
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
			if string(flags.Delimiter) != v.ResultStrValue {
				t.Errorf("%s: delimiter = %q, want %q", v.Name, string(flags.Delimiter), v.ResultStrValue)
			}
		case "ENCODING":
			if flags.Encoding.String() != v.ResultStrValue {
				t.Errorf("%s: encoding = %q, want %q", v.Name, flags.Encoding.String(), v.ResultStrValue)
			}
		case "LINE_BREAK":
			if flags.LineBreak.Value() != v.ResultStrValue {
				t.Errorf("%s: line-break = %q, want %q", v.Name, flags.LineBreak.Value(), v.ResultStrValue)
			}
		case "REPOSITORY":
			if flags.Repository != v.ResultStrValue {
				t.Errorf("%s: repository = %q, want %q", v.Name, flags.Repository, v.ResultStrValue)
			}
		case "DATETIME_FORMAT":
			if flags.DatetimeFormat != v.ResultStrValue {
				t.Errorf("%s: datetime-format = %q, want %q", v.Name, flags.DatetimeFormat, v.ResultStrValue)
			}
		case "WAIT_TIMEOUT":
			if flags.WaitTimeout != v.ResultFloatValue {
				t.Errorf("%s: wait-timeout = %f, want %f", v.Name, flags.WaitTimeout, v.ResultFloatValue)
			}
		case "NO-HEADER":
			if flags.NoHeader != v.ResultBoolValue {
				t.Errorf("%s: no-header = %t, want %t", v.Name, flags.NoHeader, v.ResultBoolValue)
			}
		case "WITHOUT-NULL":
			if flags.WithoutNull != v.ResultBoolValue {
				t.Errorf("%s: without-null = %t, want %t", v.Name, flags.WithoutNull, v.ResultBoolValue)
			}
		case "STATS":
			if flags.Stats != v.ResultBoolValue {
				t.Errorf("%s: stats = %t, want %t", v.Name, flags.Stats, v.ResultBoolValue)
			}
		}
	}

	expr := parser.SetFlag{
		Name:  "@@wait_timeout",
		Value: value.NewFloat(20),
	}
	SetFlag(expr)

	if FileLocks.WaitTimeout != 20 {
		t.Errorf("filelocks wait-timeout = %f, want %f", FileLocks.WaitTimeout, 20)
	}
	if file.WaitTimeout != 20 {
		t.Errorf("package file wait-timeout = %f, want %f", file.WaitTimeout, 20)
	}
}

var showFlagTests = []struct {
	Name    string
	Expr    parser.ShowFlag
	SetExpr parser.SetFlag
	Result  string
	Error   string
}{
	{
		Name: "Show Delimiter Not Set",
		Expr: parser.ShowFlag{
			Name: "@@delimiter",
		},
		Result: "(not set)",
	},
	{
		Name: "Show Delimiter",
		Expr: parser.ShowFlag{
			Name: "@@delimiter",
		},
		SetExpr: parser.SetFlag{
			Name:  "@@delimiter",
			Value: value.NewString("\t"),
		},
		Result: "'\\t'",
	},
	{
		Name: "Show Encoding",
		Expr: parser.ShowFlag{
			Name: "@@encoding",
		},
		SetExpr: parser.SetFlag{
			Name:  "@@encoding",
			Value: value.NewString("SJIS"),
		},
		Result: "SJIS",
	},
	{
		Name: "Show LineBreak",
		Expr: parser.ShowFlag{
			Name: "@@line_break",
		},
		SetExpr: parser.SetFlag{
			Name:  "@@line_break",
			Value: value.NewString("CRLF"),
		},
		Result: "CRLF",
	},
	{
		Name: "Show Repository",
		Expr: parser.ShowFlag{
			Name: "@@repository",
		},
		SetExpr: parser.SetFlag{
			Name:  "@@repository",
			Value: value.NewString(TestDir),
		},
		Result: TestDir,
	},
	{
		Name: "Show DatetimeFormat Not Set",
		Expr: parser.ShowFlag{
			Name: "@@datetime_format",
		},
		Result: "(not set)",
	},
	{
		Name: "Show DatetimeFormat",
		Expr: parser.ShowFlag{
			Name: "@@datetime_format",
		},
		SetExpr: parser.SetFlag{
			Name:  "@@datetime_format",
			Value: value.NewString("%Y%m%d"),
		},
		Result: "%Y%m%d",
	},
	{
		Name: "Show WaitTimeout",
		Expr: parser.ShowFlag{
			Name: "@@wait_timeout",
		},
		SetExpr: parser.SetFlag{
			Name:  "@@wait_timeout",
			Value: value.NewFloat(15),
		},
		Result: "15",
	},
	{
		Name: "Show NoHeader",
		Expr: parser.ShowFlag{
			Name: "@@no_header",
		},
		SetExpr: parser.SetFlag{
			Name:  "@@no_header",
			Value: value.NewBoolean(true),
		},
		Result: "true",
	},
	{
		Name: "Show WithoutNull",
		Expr: parser.ShowFlag{
			Name: "@@without_null",
		},
		SetExpr: parser.SetFlag{
			Name:  "@@without_null",
			Value: value.NewBoolean(true),
		},
		Result: "true",
	},
	{
		Name: "Show Stats",
		Expr: parser.ShowFlag{
			Name: "@@stats",
		},
		SetExpr: parser.SetFlag{
			Name:  "@@stats",
			Value: value.NewBoolean(true),
		},
		Result: "true",
	},
	{
		Name: "Invalid Flag Name Error",
		Expr: parser.ShowFlag{
			Name: "@@invalid",
		},
		Error: "[L:- C:-] flag name @@invalid is invalid",
	},
}

func TestShowFlag(t *testing.T) {
	initFlag()

	for _, v := range showFlagTests {
		if v.SetExpr.Value != nil {
			SetFlag(v.SetExpr)
		}
		result, err := ShowFlag(v.Expr)
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
			t.Errorf("%s: result = %s, want %s", v.Name, result, v.Result)
		}
	}
}
