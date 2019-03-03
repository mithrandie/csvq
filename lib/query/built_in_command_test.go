package query

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/syntax"
	"github.com/mithrandie/csvq/lib/value"

	"github.com/mithrandie/go-text"
	"github.com/mithrandie/go-text/fixedlen"
	"github.com/mithrandie/go-text/json"
)

var echoTests = []struct {
	Name   string
	Expr   parser.Echo
	Result string
	Error  string
}{
	{
		Name: "Echo",
		Expr: parser.Echo{
			Value: parser.NewStringValue("var"),
		},
		Result: "var",
	},
	{
		Name: "Echo Evaluate Error",
		Expr: parser.Echo{
			Value: parser.Variable{
				Name: "var",
			},
		},
		Error: "[L:- C:-] variable @var is undeclared",
	},
}

func TestEcho(t *testing.T) {
	filter := NewEmptyFilter()

	for _, v := range echoTests {
		result, err := Echo(v.Expr, filter)
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
		Result: "\"foo\"",
	},
	{
		Name: "Print Error",
		Expr: parser.Print{
			Value: parser.Variable{
				Name: "var",
			},
		},
		Error: "[L:- C:-] variable @var is undeclared",
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
			Format: parser.NewStringValue("printf test: value1 %q, value2 %q"),
			Values: []parser.QueryExpression{
				parser.NewStringValue("str"),
				parser.NewIntegerValue(1),
			},
		},
		Result: "printf test: value1 \"str\", value2 \"1\"",
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
		Error: "[L:- C:-] variable @var is undeclared",
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
		Error: "[L:- C:-] variable @var is undeclared",
	},
	{
		Name: "Printf Less Values Error",
		Expr: parser.Printf{
			Format: parser.NewStringValue("printf test: value1 %s, value2 %s %%"),
			Values: []parser.QueryExpression{
				parser.NewStringValue("str"),
			},
		},
		Error: "[L:- C:-] number of replace values does not match",
	},
	{
		Name: "Printf Greater Values Error",
		Expr: parser.Printf{
			Format: parser.NewStringValue("printf test: value1 %s, value2 %s %%"),
			Values: []parser.QueryExpression{
				parser.NewStringValue("str"),
				parser.NewIntegerValue(1),
				parser.NewIntegerValue(2),
			},
		},
		Error: "[L:- C:-] number of replace values does not match",
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
		Name: "Source from an identifier",
		Expr: parser.Source{
			FilePath: parser.Identifier{Literal: GetTestFilePath("source.sql")},
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
		Name: "Source File Invalid File Path Error",
		Expr: parser.Source{
			FilePath: parser.NewNullValueFromString("NULL"),
		},
		Error: "[L:- C:-] NULL is a invalid file path",
	},
	{
		Name: "Source File Empty File Path Error",
		Expr: parser.Source{
			FilePath: parser.Identifier{Literal: "", Quoted: true},
		},
		Error: "[L:- C:-] `` is a invalid file path",
	},
	{
		Name: "Source File Not Exist Error",
		Expr: parser.Source{
			FilePath: parser.NewStringValue(GetTestFilePath("notexist.sql")),
		},
		Error: fmt.Sprintf("[L:- C:-] file %s does not exist", GetTestFilePath("notexist.sql")),
	},
	{
		Name: "Source Syntax Error",
		Expr: parser.Source{
			FilePath: parser.NewStringValue(GetTestFilePath("source_syntaxerror.sql")),
		},
		Error: fmt.Sprintf("%s [L:1 C:34] syntax error: unexpected token \"wrong argument\"", GetTestFilePath("source_syntaxerror.sql")),
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

var parseExecuteStatementsTests = []struct {
	Name   string
	Expr   parser.Execute
	Result []parser.Statement
	Error  string
}{
	{
		Name: "ParseExecuteStatements",
		Expr: parser.Execute{
			BaseExpr:   parser.NewBaseExpr(parser.Token{}),
			Statements: parser.NewStringValue("print %q;"),
			Values: []parser.QueryExpression{
				parser.NewStringValue("executable string"),
			},
		},
		Result: []parser.Statement{
			parser.Print{
				Value: parser.NewStringValue("executable string"),
			},
		},
	},
	{
		Name: "ParseExecuteStatements String Evaluation Error",
		Expr: parser.Execute{
			BaseExpr:   parser.NewBaseExpr(parser.Token{}),
			Statements: parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			Values: []parser.QueryExpression{
				parser.NewStringValue("executable string"),
			},
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
	{
		Name: "ParseExecuteStatements Format Error",
		Expr: parser.Execute{
			BaseExpr:   parser.NewBaseExpr(parser.Token{}),
			Statements: parser.NewStringValue("print %q;"),
		},
		Error: "[L:- C:-] number of replace values does not match",
	},
	{
		Name: "ParseExecuteStatements Replace Value Error",
		Expr: parser.Execute{
			BaseExpr:   parser.NewBaseExpr(parser.Token{}),
			Statements: parser.NewStringValue("print %q;"),
			Values: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
			},
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
	{
		Name: "ParseExecuteStatements Parsing Error",
		Expr: parser.Execute{
			BaseExpr:   parser.NewBaseExpr(parser.Token{}),
			Statements: parser.NewStringValue("print;"),
		},
		Error: "(L:0 C:0) EXECUTE [L:1 C:6] syntax error: unexpected token \";\"",
	},
}

func TestParseExecuteStatements(t *testing.T) {
	filter := NewEmptyFilter()

	for _, v := range parseExecuteStatementsTests {
		result, err := ParseExecuteStatements(v.Expr, filter)
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
	Name  string
	Expr  parser.SetFlag
	Error string
}{
	{
		Name: "Set Repository",
		Expr: parser.SetFlag{
			Name:  "repository",
			Value: parser.NewStringValue(TestDir),
		},
	},
	{
		Name: "Set Timezone",
		Expr: parser.SetFlag{
			Name:  "timezone",
			Value: parser.NewStringValue("utc"),
		},
	},
	{
		Name: "Set DatetimeFormat",
		Expr: parser.SetFlag{
			Name:  "datetime_format",
			Value: parser.NewStringValue("%Y%m%d"),
		},
	},
	{
		Name: "Set WaitTimeout",
		Expr: parser.SetFlag{
			Name:  "wait_timeout",
			Value: parser.NewFloatValue(15),
		},
	},
	{
		Name: "Set Delimiter",
		Expr: parser.SetFlag{
			Name:  "delimiter",
			Value: parser.NewStringValue("\\t"),
		},
	},
	{
		Name: "Set JsonQuery",
		Expr: parser.SetFlag{
			Name:  "json_query",
			Value: parser.NewStringValue("{}"),
		},
	},
	{
		Name: "Set Encoding",
		Expr: parser.SetFlag{
			Name:  "encoding",
			Value: parser.NewStringValue("SJIS"),
		},
	},
	{
		Name: "Set NoHeader",
		Expr: parser.SetFlag{
			Name:  "no_header",
			Value: parser.NewTernaryValueFromString("true"),
		},
	},
	{
		Name: "Set WithoutNull",
		Expr: parser.SetFlag{
			Name:  "without_null",
			Value: parser.NewTernaryValueFromString("true"),
		},
	},
	{
		Name: "Set Format",
		Expr: parser.SetFlag{
			Name:  "format",
			Value: parser.NewStringValue("json"),
		},
	},
	{
		Name: "Set WriteEncoding",
		Expr: parser.SetFlag{
			Name:  "write_encoding",
			Value: parser.NewStringValue("SJIS"),
		},
	},
	{
		Name: "Set WriteDelimiter",
		Expr: parser.SetFlag{
			Name:  "write_delimiter",
			Value: parser.NewStringValue("\\t"),
		},
	},
	{
		Name: "Set WithoutHeader",
		Expr: parser.SetFlag{
			Name:  "without_header",
			Value: parser.NewTernaryValueFromString("true"),
		},
	},
	{
		Name: "Set lineBreak",
		Expr: parser.SetFlag{
			Name:  "line_break",
			Value: parser.NewStringValue("CRLF"),
		},
	},
	{
		Name: "Set EncloseAll",
		Expr: parser.SetFlag{
			Name:  "enclose_all",
			Value: parser.NewTernaryValueFromString("true"),
		},
	},
	{
		Name: "Set JsonEscape",
		Expr: parser.SetFlag{
			Name:  "json_escape",
			Value: parser.NewStringValue("hex"),
		},
	},
	{
		Name: "Set PrettyPrint",
		Expr: parser.SetFlag{
			Name:  "pretty_print",
			Value: parser.NewTernaryValueFromString("true"),
		},
	},
	{
		Name: "Set EastAsianEncoding",
		Expr: parser.SetFlag{
			Name:  "east_asian_encoding",
			Value: parser.NewTernaryValueFromString("true"),
		},
	},
	{
		Name: "Set CountDiacriticalSign",
		Expr: parser.SetFlag{
			Name:  "count_diacritical_sign",
			Value: parser.NewTernaryValueFromString("true"),
		},
	},
	{
		Name: "Set CountFormatCode",
		Expr: parser.SetFlag{
			Name:  "count_format_code",
			Value: parser.NewTernaryValueFromString("true"),
		},
	},
	{
		Name: "Set Color",
		Expr: parser.SetFlag{
			Name:  "color",
			Value: parser.NewTernaryValueFromString("true"),
		},
	},
	{
		Name: "Set Quiet",
		Expr: parser.SetFlag{
			Name:  "quiet",
			Value: parser.NewTernaryValueFromString("true"),
		},
	},
	{
		Name: "Set CPU",
		Expr: parser.SetFlag{
			Name:  "cpu",
			Value: parser.NewIntegerValue(int64(runtime.NumCPU())),
		},
	},
	{
		Name: "Set Stats",
		Expr: parser.SetFlag{
			Name:  "stats",
			Value: parser.NewTernaryValueFromString("true"),
		},
	},
	{
		Name: "Set Encoding with Identifier",
		Expr: parser.SetFlag{
			Name:  "encoding",
			Value: parser.Identifier{Literal: "sjis"},
		},
	},
	{
		Name: "Set Delimiter Evaluation Error",
		Expr: parser.SetFlag{
			Name:  "delimiter",
			Value: parser.FieldReference{Column: parser.Identifier{Literal: "err"}},
		},
		Error: "[L:- C:-] field err does not exist",
	},
	{
		Name: "Set Delimiter Value Error",
		Expr: parser.SetFlag{
			Name:  "delimiter",
			Value: parser.NewTernaryValueFromString("true"),
		},
		Error: "[L:- C:-] true for @@delimiter is not allowed",
	},
	{
		Name: "Set WaitTimeout Value Error",
		Expr: parser.SetFlag{
			Name:  "wait_timeout",
			Value: parser.NewTernaryValueFromString("true"),
		},
		Error: "[L:- C:-] true for @@wait_timeout is not allowed",
	},
	{
		Name: "Set WithoutNull Value Error",
		Expr: parser.SetFlag{
			Name:  "without_null",
			Value: parser.NewStringValue("string"),
		},
		Error: "[L:- C:-] 'string' for @@without_null is not allowed",
	},
	{
		Name: "Set CPU Value Error",
		Expr: parser.SetFlag{
			Name:  "cpu",
			Value: parser.NewStringValue("invalid"),
		},
		Error: "[L:- C:-] 'invalid' for @@cpu is not allowed",
	},
	{
		Name: "Invalid Flag Name Error",
		Expr: parser.SetFlag{
			Name:  "invalid",
			Value: parser.NewStringValue("string"),
		},
		Error: "[L:- C:-] @@invalid is an unknown flag",
	},
	{
		Name: "Invalid Flag Value Error",
		Expr: parser.SetFlag{
			Name:  "line_break",
			Value: parser.NewStringValue("invalid"),
		},
		Error: "[L:- C:-] line-break must be one of CRLF|LF|CR",
	},
}

func TestSetFlag(t *testing.T) {
	filter := NewEmptyFilter()

	for _, v := range setFlagTests {
		initCmdFlag()
		err := SetFlag(v.Expr, filter)
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
	}
	initCmdFlag()
}

var addFlagElementTests = []struct {
	Name   string
	Expr   parser.AddFlagElement
	Init   func(*cmd.Flags)
	Expect func() *cmd.Flags
	Error  string
}{
	{
		Name: "Add Element To DatetimeFormat",
		Expr: parser.AddFlagElement{
			Name:  "datetime_format",
			Value: parser.NewStringValue("%Y%m%d"),
		},
		Init: func(flags *cmd.Flags) {
			flags.DatetimeFormat = []string{"%Y:%m:%d"}
		},
		Expect: func() *cmd.Flags {
			expect := new(cmd.Flags)
			initFlag(expect)
			expect.DatetimeFormat = []string{"%Y:%m:%d", "%Y%m%d"}
			return expect
		},
	},
	{
		Name: "Add Element Unsupported Flag Name",
		Expr: parser.AddFlagElement{
			Name:  "format",
			Value: parser.NewStringValue("%Y%m%d"),
		},
		Init: func(flags *cmd.Flags) {
			flags.DatetimeFormat = []string{"%Y:%m:%d"}
		},
		Error: "[L:- C:-] add flag element syntax does not support @@format",
	},
	{
		Name: "Add Element Invalid Flag Name",
		Expr: parser.AddFlagElement{
			Name:  "invalid",
			Value: parser.NewStringValue("%Y%m%d"),
		},
		Init: func(flags *cmd.Flags) {
			flags.DatetimeFormat = []string{"%Y:%m:%d"}
		},
		Error: "[L:- C:-] @@invalid is an unknown flag",
	},
}

func TestAddFlagElement(t *testing.T) {
	filter := NewEmptyFilter()

	for _, v := range addFlagElementTests {
		initCmdFlag()
		flags := cmd.GetFlags()
		v.Init(flags)

		err := AddFlagElement(v.Expr, filter)
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

		expect := v.Expect()
		if !reflect.DeepEqual(flags, expect) {
			t.Errorf("%s: result = %v, want %v", v.Name, flags, expect)
		}
	}
	initCmdFlag()
}

var removeFlagElementTests = []struct {
	Name   string
	Expr   parser.RemoveFlagElement
	Init   func(*cmd.Flags)
	Expect func() *cmd.Flags
	Error  string
}{
	{
		Name: "Remove Element from DatetimeFormat",
		Expr: parser.RemoveFlagElement{
			Name:  "datetime_format",
			Value: parser.NewStringValue("%Y%m%d"),
		},
		Init: func(flags *cmd.Flags) {
			flags.DatetimeFormat = []string{"%Y%m%d", "%Y:%m:%d"}
		},
		Expect: func() *cmd.Flags {
			expect := new(cmd.Flags)
			initFlag(expect)
			expect.DatetimeFormat = []string{"%Y:%m:%d"}
			return expect
		},
	},
	{
		Name: "Remove Element from DatetimeFormat with List Index",
		Expr: parser.RemoveFlagElement{
			Name:  "datetime_format",
			Value: parser.NewIntegerValue(1),
		},
		Init: func(flags *cmd.Flags) {
			flags.DatetimeFormat = []string{"%Y%m%d", "%Y:%m:%d"}
		},
		Expect: func() *cmd.Flags {
			expect := new(cmd.Flags)
			initFlag(expect)
			expect.DatetimeFormat = []string{"%Y%m%d"}
			return expect
		},
	},
	{
		Name: "Remove Element Invalid Flag Value",
		Expr: parser.RemoveFlagElement{
			Name:  "datetime_format",
			Value: parser.NewNullValueFromString("null"),
		},
		Init:  func(flags *cmd.Flags) {},
		Error: "[L:- C:-] null is an invalid value for @@datetime_format to specify the element",
	},
	{
		Name: "Remove Element Evaluation Error",
		Expr: parser.RemoveFlagElement{
			Name:  "format",
			Value: parser.FieldReference{Column: parser.Identifier{Literal: "err"}},
		},
		Init:  func(flags *cmd.Flags) {},
		Error: "[L:- C:-] field err does not exist",
	},
	{
		Name: "Remove Element Unsupported Flag Name",
		Expr: parser.RemoveFlagElement{
			Name:  "format",
			Value: parser.NewIntegerValue(1),
		},
		Init:  func(flags *cmd.Flags) {},
		Error: "[L:- C:-] remove flag element syntax does not support @@format",
	},
	{
		Name: "Remove Element Invalid Flag Name",
		Expr: parser.RemoveFlagElement{
			Name:  "invalid",
			Value: parser.NewIntegerValue(1),
		},
		Init:  func(flags *cmd.Flags) {},
		Error: "[L:- C:-] @@invalid is an unknown flag",
	},
}

func TestRemoveFlagElement(t *testing.T) {
	filter := NewEmptyFilter()

	for _, v := range removeFlagElementTests {
		initCmdFlag()
		flags := cmd.GetFlags()
		v.Init(flags)

		err := RemoveFlagElement(v.Expr, filter)
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

		expect := v.Expect()
		if !reflect.DeepEqual(flags, expect) {
			t.Errorf("%s: result = %v, want %v", v.Name, flags, expect)
		}
	}
	initCmdFlag()
}

var showFlagTests = []struct {
	Name     string
	Expr     parser.ShowFlag
	SetExprs []parser.SetFlag
	Result   string
	Error    string
}{
	{
		Name: "Show Repository",
		Expr: parser.ShowFlag{
			Name: "repository",
		},
		SetExprs: []parser.SetFlag{
			{
				Name:  "repository",
				Value: parser.NewStringValue(TestDir),
			},
		},
		Result: "\033[34;1m@@REPOSITORY:\033[0m \033[32m" + TestDir + "\033[0m",
	},
	{
		Name: "Show Repository Not Set",
		Expr: parser.ShowFlag{
			Name: "repository",
		},
		SetExprs: []parser.SetFlag{
			{
				Name:  "repository",
				Value: parser.NewStringValue(""),
			},
		},
		Result: "\033[34;1m@@REPOSITORY:\033[0m \033[90m(current dir: " + GetWD() + ")\033[0m",
	},
	{
		Name: "Show Timezone",
		Expr: parser.ShowFlag{
			Name: "timezone",
		},
		SetExprs: []parser.SetFlag{
			{
				Name:  "timezone",
				Value: parser.NewStringValue("UTC"),
			},
		},
		Result: "\033[34;1m@@TIMEZONE:\033[0m \033[32mUTC\033[0m",
	},
	{
		Name: "Show DatetimeFormat Not Set",
		Expr: parser.ShowFlag{
			Name: "datetime_format",
		},
		Result: "\033[34;1m@@DATETIME_FORMAT:\033[0m \033[90m(not set)\033[0m",
	},
	{
		Name: "Show DatetimeFormat",
		Expr: parser.ShowFlag{
			Name: "datetime_format",
		},
		SetExprs: []parser.SetFlag{
			{
				Name:  "datetime_format",
				Value: parser.NewStringValue("[\"%Y%m%d\", \"%Y%m%d %H%i%s\"]"),
			},
		},
		Result: "\033[34;1m@@DATETIME_FORMAT:\033[0m \033[32m[\"%Y%m%d\", \"%Y%m%d %H%i%s\"]\033[0m",
	},
	{
		Name: "Show WaitTimeout",
		Expr: parser.ShowFlag{
			Name: "wait_timeout",
		},
		SetExprs: []parser.SetFlag{
			{
				Name:  "wait_timeout",
				Value: parser.NewFloatValue(15),
			},
		},
		Result: "\033[34;1m@@WAIT_TIMEOUT:\033[0m \033[35m15\033[0m",
	},
	{
		Name: "Show Import Format",
		Expr: parser.ShowFlag{
			Name: "import_format",
		},
		SetExprs: []parser.SetFlag{
			{
				Name:  "import_format",
				Value: parser.NewStringValue("tsv"),
			},
		},
		Result: "\033[34;1m@@IMPORT_FORMAT:\033[0m \033[32mTSV\033[0m",
	},
	{
		Name: "Show Delimiter",
		Expr: parser.ShowFlag{
			Name: "delimiter",
		},
		SetExprs: []parser.SetFlag{
			{
				Name:  "delimiter",
				Value: parser.NewStringValue("\t"),
			},
		},
		Result: "\033[34;1m@@DELIMITER:\033[0m \033[32m'\\t'\033[0m",
	},
	{
		Name: "Show Delimiter Positions",
		Expr: parser.ShowFlag{
			Name: "delimiter_positions",
		},
		SetExprs: []parser.SetFlag{
			{
				Name:  "delimiter_positions",
				Value: parser.NewStringValue("s[2, 5, 10]"),
			},
		},
		Result: "\033[34;1m@@DELIMITER_POSITIONS:\033[0m \033[32mS[2, 5, 10]\033[0m",
	},
	{
		Name: "Show Delimiter Positions as spaces",
		Expr: parser.ShowFlag{
			Name: "delimiter_positions",
		},
		SetExprs: []parser.SetFlag{
			{
				Name:  "delimiter_positions",
				Value: parser.NewStringValue("SPACES"),
			},
		},
		Result: "\033[34;1m@@DELIMITER_POSITIONS:\033[0m \033[32mSPACES\033[0m",
	},
	{
		Name: "Show JsonQuery",
		Expr: parser.ShowFlag{
			Name: "json_query",
		},
		SetExprs: []parser.SetFlag{
			{
				Name:  "json_query",
				Value: parser.NewStringValue("{}"),
			},
		},
		Result: "\033[34;1m@@JSON_QUERY:\033[0m \033[32m{}\033[0m",
	},
	{
		Name: "Show JsonQuery Empty",
		Expr: parser.ShowFlag{
			Name: "json_query",
		},
		SetExprs: []parser.SetFlag{},
		Result:   "\033[34;1m@@JSON_QUERY:\033[0m \033[90m(empty)\033[0m",
	},
	{
		Name: "Show Encoding",
		Expr: parser.ShowFlag{
			Name: "encoding",
		},
		SetExprs: []parser.SetFlag{
			{
				Name:  "encoding",
				Value: parser.NewStringValue("SJIS"),
			},
		},
		Result: "\033[34;1m@@ENCODING:\033[0m \033[32mSJIS\033[0m",
	},
	{
		Name: "Show NoHeader",
		Expr: parser.ShowFlag{
			Name: "no_header",
		},
		SetExprs: []parser.SetFlag{
			{
				Name:  "no_header",
				Value: parser.NewTernaryValueFromString("true"),
			},
		},
		Result: "\033[34;1m@@NO_HEADER:\033[0m \033[33;1mtrue\033[0m",
	},
	{
		Name: "Show WithoutNull",
		Expr: parser.ShowFlag{
			Name: "without_null",
		},
		SetExprs: []parser.SetFlag{
			{
				Name:  "without_null",
				Value: parser.NewTernaryValueFromString("true"),
			},
		},
		Result: "\033[34;1m@@WITHOUT_NULL:\033[0m \033[33;1mtrue\033[0m",
	},
	{
		Name: "Show Format",
		Expr: parser.ShowFlag{
			Name: "format",
		},
		SetExprs: []parser.SetFlag{
			{
				Name:  "format",
				Value: parser.NewStringValue("json"),
			},
		},
		Result: "\033[34;1m@@FORMAT:\033[0m \033[32mJSON\033[0m",
	},
	{
		Name: "Show WriteEncoding",
		Expr: parser.ShowFlag{
			Name: "write_encoding",
		},
		SetExprs: []parser.SetFlag{
			{
				Name:  "write_encoding",
				Value: parser.NewStringValue("SJIS"),
			},
		},
		Result: "\033[34;1m@@WRITE_ENCODING:\033[0m \033[32mSJIS\033[0m",
	},
	{
		Name: "Show WriteEncoding Ignored",
		Expr: parser.ShowFlag{
			Name: "write_encoding",
		},
		SetExprs: []parser.SetFlag{
			{
				Name:  "write_encoding",
				Value: parser.NewStringValue("SJIS"),
			},
			{
				Name:  "format",
				Value: parser.NewStringValue("JSON"),
			},
		},
		Result: "\033[34;1m@@WRITE_ENCODING:\033[0m \033[90m(ignored) SJIS\033[0m",
	},
	{
		Name: "Show WriteDelimiter",
		Expr: parser.ShowFlag{
			Name: "write_delimiter",
		},
		SetExprs: []parser.SetFlag{
			{
				Name:  "write_delimiter",
				Value: parser.NewStringValue("\t"),
			},
			{
				Name:  "format",
				Value: parser.NewStringValue("CSV"),
			},
		},
		Result: "\033[34;1m@@WRITE_DELIMITER:\033[0m \033[32m'\\t'\033[0m",
	},
	{
		Name: "Show WriteDelimiter Ignored",
		Expr: parser.ShowFlag{
			Name: "write_delimiter",
		},
		SetExprs: []parser.SetFlag{
			{
				Name:  "write_delimiter",
				Value: parser.NewStringValue("\t"),
			},
			{
				Name:  "format",
				Value: parser.NewStringValue("JSON"),
			},
		},
		Result: "\033[34;1m@@WRITE_DELIMITER:\033[0m \033[90m(ignored) '\\t'\033[0m",
	},
	{
		Name: "Show WriteDelimiterPositions for Single-Line FIXED",
		Expr: parser.ShowFlag{
			Name: "write_delimiter_positions",
		},
		SetExprs: []parser.SetFlag{
			{
				Name:  "write_delimiter_positions",
				Value: parser.NewStringValue("s[2, 5, 10]"),
			},
			{
				Name:  "format",
				Value: parser.NewStringValue("FIXED"),
			},
		},
		Result: "\033[34;1m@@WRITE_DELIMITER_POSITIONS:\033[0m \033[32mS[2, 5, 10]\033[0m",
	},
	{
		Name: "Show WriteDelimiterPositions for FIXED",
		Expr: parser.ShowFlag{
			Name: "write_delimiter_positions",
		},
		SetExprs: []parser.SetFlag{
			{
				Name:  "write_delimiter_positions",
				Value: parser.NewStringValue("spaces"),
			},
			{
				Name:  "format",
				Value: parser.NewStringValue("FIXED"),
			},
		},
		Result: "\033[34;1m@@WRITE_DELIMITER_POSITIONS:\033[0m \033[32mSPACES\033[0m",
	},
	{
		Name: "Show WriteDelimiterPositions Ignored",
		Expr: parser.ShowFlag{
			Name: "write_delimiter_positions",
		},
		SetExprs: []parser.SetFlag{
			{
				Name:  "write_delimiter_positions",
				Value: parser.NewStringValue("spaces"),
			},
			{
				Name:  "format",
				Value: parser.NewStringValue("JSON"),
			},
		},
		Result: "\033[34;1m@@WRITE_DELIMITER_POSITIONS:\033[0m \033[90m(ignored) SPACES\033[0m",
	},
	{
		Name: "Show WithoutHeader",
		Expr: parser.ShowFlag{
			Name: "without_header",
		},
		SetExprs: []parser.SetFlag{
			{
				Name:  "without_header",
				Value: parser.NewTernaryValueFromString("true"),
			},
			{
				Name:  "format",
				Value: parser.NewStringValue("CSV"),
			},
		},
		Result: "\033[34;1m@@WITHOUT_HEADER:\033[0m \033[33;1mtrue\033[0m",
	},
	{
		Name: "Show WithoutHeader Ignored",
		Expr: parser.ShowFlag{
			Name: "without_header",
		},
		SetExprs: []parser.SetFlag{
			{
				Name:  "without_header",
				Value: parser.NewTernaryValueFromString("true"),
			},
			{
				Name:  "format",
				Value: parser.NewStringValue("JSON"),
			},
		},
		Result: "\033[34;1m@@WITHOUT_HEADER:\033[0m \033[90m(ignored) true\033[0m",
	},
	{
		Name: "Show WithoutHeader with Single-Line Fixed-Length",
		Expr: parser.ShowFlag{
			Name: "without_header",
		},
		SetExprs: []parser.SetFlag{
			{
				Name:  "format",
				Value: parser.NewStringValue("fixed"),
			},
			{
				Name:  "write_delimiter_positions",
				Value: parser.NewStringValue("s[2, 5, 10]"),
			},
			{
				Name:  "without_header",
				Value: parser.NewTernaryValueFromString("true"),
			},
		},
		Result: "\033[34;1m@@WITHOUT_HEADER:\033[0m \033[90m(ignored) true\033[0m",
	},
	{
		Name: "Show lineBreak",
		Expr: parser.ShowFlag{
			Name: "line_break",
		},
		SetExprs: []parser.SetFlag{
			{
				Name:  "line_break",
				Value: parser.NewStringValue("CRLF"),
			},
		},
		Result: "\033[34;1m@@LINE_BREAK:\033[0m \033[32mCRLF\033[0m",
	},
	{
		Name: "Show lineBreak with Single-Line Fixed-Length",
		Expr: parser.ShowFlag{
			Name: "line_break",
		},
		SetExprs: []parser.SetFlag{
			{
				Name:  "format",
				Value: parser.NewStringValue("fixed"),
			},
			{
				Name:  "write_delimiter_positions",
				Value: parser.NewStringValue("s[2, 5, 10]"),
			},
			{
				Name:  "line_break",
				Value: parser.NewStringValue("CRLF"),
			},
		},
		Result: "\033[34;1m@@LINE_BREAK:\033[0m \033[90m(ignored) CRLF\033[0m",
	},
	{
		Name: "Show EncloseAll",
		Expr: parser.ShowFlag{
			Name: "enclose_all",
		},
		SetExprs: []parser.SetFlag{
			{
				Name:  "enclose_all",
				Value: parser.NewTernaryValueFromString("true"),
			},
			{
				Name:  "format",
				Value: parser.NewStringValue("CSV"),
			},
		},
		Result: "\033[34;1m@@ENCLOSE_ALL:\033[0m \033[33;1mtrue\033[0m",
	},
	{
		Name: "Show EncloseAll Ignored",
		Expr: parser.ShowFlag{
			Name: "enclose_all",
		},
		SetExprs: []parser.SetFlag{
			{
				Name:  "enclose_all",
				Value: parser.NewTernaryValueFromString("true"),
			},
			{
				Name:  "format",
				Value: parser.NewStringValue("JSON"),
			},
		},
		Result: "\033[34;1m@@ENCLOSE_ALL:\033[0m \033[90m(ignored) true\033[0m",
	},
	{
		Name: "Show JsonEscape",
		Expr: parser.ShowFlag{
			Name: "json_escape",
		},
		SetExprs: []parser.SetFlag{
			{
				Name:  "json_escape",
				Value: parser.NewStringValue("HEXALL"),
			},
			{
				Name:  "format",
				Value: parser.NewStringValue("JSON"),
			},
		},
		Result: "\033[34;1m@@JSON_ESCAPE:\033[0m \033[32mHEXALL\033[0m",
	},
	{
		Name: "Show JsonEscape Ignored",
		Expr: parser.ShowFlag{
			Name: "json_escape",
		},
		SetExprs: []parser.SetFlag{
			{
				Name:  "json_escape",
				Value: parser.NewStringValue("BACKSLASH"),
			},
			{
				Name:  "format",
				Value: parser.NewStringValue("CSV"),
			},
		},
		Result: "\033[34;1m@@JSON_ESCAPE:\033[0m \033[90m(ignored) BACKSLASH\033[0m",
	},
	{
		Name: "Show PrettyPrint",
		Expr: parser.ShowFlag{
			Name: "pretty_print",
		},
		SetExprs: []parser.SetFlag{
			{
				Name:  "pretty_print",
				Value: parser.NewTernaryValueFromString("true"),
			},
			{
				Name:  "format",
				Value: parser.NewStringValue("JSON"),
			},
		},
		Result: "\033[34;1m@@PRETTY_PRINT:\033[0m \033[33;1mtrue\033[0m",
	},
	{
		Name: "Show PrettyPrint Ignored",
		Expr: parser.ShowFlag{
			Name: "pretty_print",
		},
		SetExprs: []parser.SetFlag{
			{
				Name:  "pretty_print",
				Value: parser.NewTernaryValueFromString("true"),
			},
		},
		Result: "\033[34;1m@@PRETTY_PRINT:\033[0m \033[90m(ignored) true\033[0m",
	},
	{
		Name: "Show EastAsianEncoding",
		Expr: parser.ShowFlag{
			Name: "east_asian_encoding",
		},
		SetExprs: []parser.SetFlag{
			{
				Name:  "east_asian_encoding",
				Value: parser.NewTernaryValueFromString("true"),
			},
			{
				Name:  "format",
				Value: parser.NewStringValue("TEXT"),
			},
		},
		Result: "\033[34;1m@@EAST_ASIAN_ENCODING:\033[0m \033[33;1mtrue\033[0m",
	},
	{
		Name: "Show EastAsianEncoding Ignored",
		Expr: parser.ShowFlag{
			Name: "east_asian_encoding",
		},
		SetExprs: []parser.SetFlag{
			{
				Name:  "east_asian_encoding",
				Value: parser.NewTernaryValueFromString("true"),
			},
			{
				Name:  "format",
				Value: parser.NewStringValue("JSON"),
			},
		},
		Result: "\033[34;1m@@EAST_ASIAN_ENCODING:\033[0m \033[90m(ignored) true\033[0m",
	},
	{
		Name: "Show CountDiacriticalSign",
		Expr: parser.ShowFlag{
			Name: "count_diacritical_sign",
		},
		SetExprs: []parser.SetFlag{
			{
				Name:  "count_diacritical_sign",
				Value: parser.NewTernaryValueFromString("true"),
			},
			{
				Name:  "format",
				Value: parser.NewStringValue("TEXT"),
			},
		},
		Result: "\033[34;1m@@COUNT_DIACRITICAL_SIGN:\033[0m \033[33;1mtrue\033[0m",
	},
	{
		Name: "Show CountDiacriticalSign Ignored",
		Expr: parser.ShowFlag{
			Name: "count_diacritical_sign",
		},
		SetExprs: []parser.SetFlag{
			{
				Name:  "count_diacritical_sign",
				Value: parser.NewTernaryValueFromString("true"),
			},
			{
				Name:  "format",
				Value: parser.NewStringValue("JSON"),
			},
		},
		Result: "\033[34;1m@@COUNT_DIACRITICAL_SIGN:\033[0m \033[90m(ignored) true\033[0m",
	},
	{
		Name: "Show CountFormatCode",
		Expr: parser.ShowFlag{
			Name: "count_format_code",
		},
		SetExprs: []parser.SetFlag{
			{
				Name:  "count_format_code",
				Value: parser.NewTernaryValueFromString("true"),
			},
			{
				Name:  "format",
				Value: parser.NewStringValue("TEXT"),
			},
		},
		Result: "\033[34;1m@@COUNT_FORMAT_CODE:\033[0m \033[33;1mtrue\033[0m",
	},
	{
		Name: "Show CountFormatCode Ignored",
		Expr: parser.ShowFlag{
			Name: "count_format_code",
		},
		SetExprs: []parser.SetFlag{
			{
				Name:  "count_format_code",
				Value: parser.NewTernaryValueFromString("true"),
			},
			{
				Name:  "format",
				Value: parser.NewStringValue("JSON"),
			},
		},
		Result: "\033[34;1m@@COUNT_FORMAT_CODE:\033[0m \033[90m(ignored) true\033[0m",
	},
	{
		Name: "Show Color",
		Expr: parser.ShowFlag{
			Name: "color",
		},
		SetExprs: []parser.SetFlag{
			{
				Name:  "color",
				Value: parser.NewTernaryValueFromString("true"),
			},
		},
		Result: "\033[34;1m@@COLOR:\033[0m \033[33;1mtrue\033[0m",
	},
	{
		Name: "Show Quiet",
		Expr: parser.ShowFlag{
			Name: "quiet",
		},
		SetExprs: []parser.SetFlag{
			{
				Name:  "quiet",
				Value: parser.NewTernaryValueFromString("true"),
			},
		},
		Result: "\033[34;1m@@QUIET:\033[0m \033[33;1mtrue\033[0m",
	},
	{
		Name: "Show CPU",
		Expr: parser.ShowFlag{
			Name: "cpu",
		},
		SetExprs: []parser.SetFlag{
			{
				Name:  "cpu",
				Value: parser.NewIntegerValue(1),
			},
		},
		Result: "\033[34;1m@@CPU:\033[0m \033[35m1\033[0m",
	},
	{
		Name: "Show Stats",
		Expr: parser.ShowFlag{
			Name: "stats",
		},
		SetExprs: []parser.SetFlag{
			{
				Name:  "stats",
				Value: parser.NewTernaryValueFromString("true"),
			},
		},
		Result: "\033[34;1m@@STATS:\033[0m \033[33;1mtrue\033[0m",
	},
	{
		Name: "Invalid Flag Name Error",
		Expr: parser.ShowFlag{
			Name: "invalid",
		},
		Error: "[L:- C:-] @@invalid is an unknown flag",
	},
}

func TestShowFlag(t *testing.T) {
	filter := NewEmptyFilter()

	for _, v := range showFlagTests {
		initCmdFlag()
		cmd.GetFlags().SetColor(true)
		for _, expr := range v.SetExprs {
			SetFlag(expr, filter)
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
	initCmdFlag()
}

var showObjectsTests = []struct {
	Name                    string
	Expr                    parser.ShowObjects
	Filter                  *Filter
	ImportFormat            cmd.Format
	Delimiter               rune
	DelimiterPositions      fixedlen.DelimiterPositions
	SingleLine              bool
	JsonQuery               string
	Repository              string
	Format                  cmd.Format
	WriteDelimiter          rune
	WriteDelimiterPositions fixedlen.DelimiterPositions
	WriteAsSingleLine       bool
	ViewCache               ViewMap
	UncommittedViews        *UncommittedViewMap
	Expect                  string
	Error                   string
}{
	{
		Name: "ShowObjects Tables",
		Expr: parser.ShowObjects{Type: parser.Identifier{Literal: "tables"}},
		ViewCache: ViewMap{
			"TABLE1.CSV": &View{
				Header: NewHeader("table1", []string{"col1", "col2"}),
				FileInfo: &FileInfo{
					Path:      "table1.csv",
					Delimiter: '\t',
					Format:    cmd.CSV,
					Encoding:  text.SJIS,
					LineBreak: text.CRLF,
					NoHeader:  true,
				},
			},
			"TABLE1.TSV": &View{
				Header: NewHeader("table1", []string{"col1", "col2"}),
				FileInfo: &FileInfo{
					Path:      "table1.tsv",
					Delimiter: '\t',
					Format:    cmd.TSV,
					Encoding:  text.UTF8,
					LineBreak: text.LF,
					NoHeader:  false,
				},
			},
			"TABLE1.JSON": &View{
				Header: NewHeader("table1", []string{"col1", "col2"}),
				FileInfo: &FileInfo{
					Path:        "table1.json",
					JsonQuery:   "{}",
					Format:      cmd.JSON,
					Encoding:    text.UTF8,
					LineBreak:   text.LF,
					PrettyPrint: false,
				},
			},
			"TABLE2.JSON": &View{
				Header: NewHeader("table2", []string{"col1", "col2"}),
				FileInfo: &FileInfo{
					Path:        "table2.json",
					JsonQuery:   "",
					Format:      cmd.JSON,
					Encoding:    text.UTF8,
					LineBreak:   text.LF,
					JsonEscape:  json.HexDigits,
					PrettyPrint: false,
				},
			},
			"TABLE1.TXT": &View{
				Header: NewHeader("table1", []string{"col1", "col2"}),
				FileInfo: &FileInfo{
					Path:               "table1.txt",
					DelimiterPositions: []int{3, 12},
					Format:             cmd.FIXED,
					Encoding:           text.UTF8,
					LineBreak:          text.LF,
					NoHeader:           false,
				},
			},
			"TABLE2.TXT": &View{
				Header: NewHeader("table2", []string{"col1", "col2"}),
				FileInfo: &FileInfo{
					Path:               "table2.txt",
					DelimiterPositions: []int{3, 12},
					SingleLine:         true,
					Format:             cmd.FIXED,
					Encoding:           text.UTF8,
					LineBreak:          text.LF,
					NoHeader:           false,
				},
			},
		},
		Expect: "\n" +
			"                       Loaded Tables\n" +
			"-----------------------------------------------------------\n" +
			" table1.csv\n" +
			"     Fields: col1, col2\n" +
			"     Format: CSV      Delimiter: '\\t'  Enclose All: false\n" +
			"     Encoding: SJIS   LineBreak: CRLF  Header: false\n" +
			" table1.json\n" +
			"     Fields: col1, col2\n" +
			"     Format: JSON     Escape: BACKSLASH  Query: {}\n" +
			"     Encoding: UTF8   LineBreak: LF    Pretty Print: false\n" +
			" table1.tsv\n" +
			"     Fields: col1, col2\n" +
			"     Format: TSV      Delimiter: '\\t'  Enclose All: false\n" +
			"     Encoding: UTF8   LineBreak: LF    Header: true\n" +
			" table1.txt\n" +
			"     Fields: col1, col2\n" +
			"     Format: FIXED    Delimiter Positions: [3, 12]\n" +
			"     Encoding: UTF8   LineBreak: LF    Header: true\n" +
			" table2.json\n" +
			"     Fields: col1, col2\n" +
			"     Format: JSON     Escape: HEX      Query: (empty)\n" +
			"     Encoding: UTF8   LineBreak: LF    Pretty Print: false\n" +
			" table2.txt\n" +
			"     Fields: col1, col2\n" +
			"     Format: FIXED    Delimiter Positions: S[3, 12]\n" +
			"     Encoding: UTF8\n" +
			"\n",
	},
	{
		Name: "ShowObjects Tables Uncommitted",
		Expr: parser.ShowObjects{Type: parser.Identifier{Literal: "tables"}},
		ViewCache: ViewMap{
			"TABLE1.CSV": &View{
				Header: NewHeader("table1", []string{"col1", "col2"}),
				FileInfo: &FileInfo{
					Path:      "table1.csv",
					Delimiter: '\t',
					Format:    cmd.CSV,
					Encoding:  text.SJIS,
					LineBreak: text.CRLF,
					NoHeader:  true,
				},
			},
			"TABLE1.TSV": &View{
				Header: NewHeader("table1", []string{"col1", "col2"}),
				FileInfo: &FileInfo{
					Path:      "table1.tsv",
					Delimiter: '\t',
					Format:    cmd.TSV,
					Encoding:  text.UTF8,
					LineBreak: text.LF,
					NoHeader:  false,
				},
			},
			"TABLE1.JSON": &View{
				Header: NewHeader("table1", []string{"col1", "col2"}),
				FileInfo: &FileInfo{
					Path:        "table1.json",
					JsonQuery:   "{}",
					Format:      cmd.JSON,
					Encoding:    text.UTF8,
					LineBreak:   text.LF,
					PrettyPrint: false,
				},
			},
			"TABLE2.JSON": &View{
				Header: NewHeader("table2", []string{"col1", "col2"}),
				FileInfo: &FileInfo{
					Path:        "table2.json",
					JsonQuery:   "",
					Format:      cmd.JSON,
					Encoding:    text.UTF8,
					LineBreak:   text.LF,
					PrettyPrint: false,
				},
			},
			"TABLE1.TXT": &View{
				Header: NewHeader("table1", []string{"col1", "col2"}),
				FileInfo: &FileInfo{
					Path:               "table1.txt",
					DelimiterPositions: []int{3, 12},
					Format:             cmd.FIXED,
					Encoding:           text.UTF8,
					LineBreak:          text.LF,
					NoHeader:           false,
				},
			},
			"TABLE2.TXT": &View{
				Header: NewHeader("table2", []string{"col1", "col2"}),
				FileInfo: &FileInfo{
					Path:               "table2.txt",
					DelimiterPositions: []int{3, 12},
					Format:             cmd.FIXED,
					Encoding:           text.UTF8,
					LineBreak:          text.LF,
					NoHeader:           false,
					SingleLine:         true,
				},
			},
		},
		UncommittedViews: &UncommittedViewMap{
			Created: map[string]*FileInfo{
				"TABLE1.TSV": {Path: "table1.tsv"},
			},
			Updated: map[string]*FileInfo{
				"TABLE2.JSON": {Path: "table2.json"},
			},
		},
		Expect: "\n" +
			"           Loaded Tables (Uncommitted: 2 Tables)\n" +
			"-----------------------------------------------------------\n" +
			" table1.csv\n" +
			"     Fields: col1, col2\n" +
			"     Format: CSV      Delimiter: '\\t'  Enclose All: false\n" +
			"     Encoding: SJIS   LineBreak: CRLF  Header: false\n" +
			" table1.json\n" +
			"     Fields: col1, col2\n" +
			"     Format: JSON     Escape: BACKSLASH  Query: {}\n" +
			"     Encoding: UTF8   LineBreak: LF    Pretty Print: false\n" +
			" *Created* table1.tsv\n" +
			"     Fields: col1, col2\n" +
			"     Format: TSV      Delimiter: '\\t'  Enclose All: false\n" +
			"     Encoding: UTF8   LineBreak: LF    Header: true\n" +
			" table1.txt\n" +
			"     Fields: col1, col2\n" +
			"     Format: FIXED    Delimiter Positions: [3, 12]\n" +
			"     Encoding: UTF8   LineBreak: LF    Header: true\n" +
			" *Updated* table2.json\n" +
			"     Fields: col1, col2\n" +
			"     Format: JSON     Escape: BACKSLASH  Query: (empty)\n" +
			"     Encoding: UTF8   LineBreak: LF    Pretty Print: false\n" +
			" table2.txt\n" +
			"     Fields: col1, col2\n" +
			"     Format: FIXED    Delimiter Positions: S[3, 12]\n" +
			"     Encoding: UTF8\n" +
			"\n",
	},
	{
		Name: "ShowObjects Tables Long Fields",
		Expr: parser.ShowObjects{Type: parser.Identifier{Literal: "tables"}},
		ViewCache: ViewMap{
			"TABLE1.CSV": &View{
				Header: NewHeader("table1", []string{"colabcdef1", "colabcdef2", "colabcdef3", "colabcdef4", "colabcdef5", "colabcdef6", "colabcdef7"}),
				FileInfo: &FileInfo{
					Path:      "table1.csv",
					Delimiter: '\t',
					Format:    cmd.CSV,
					Encoding:  text.SJIS,
					LineBreak: text.CRLF,
					NoHeader:  true,
				},
			},
		},
		Expect: "\n" +
			"                              Loaded Tables\n" +
			"--------------------------------------------------------------------------\n" +
			" table1.csv\n" +
			"     Fields: colabcdef1, colabcdef2, colabcdef3, colabcdef4, colabcdef5, \n" +
			"             colabcdef6, colabcdef7\n" +
			"     Format: CSV      Delimiter: '\\t'  Enclose All: false\n" +
			"     Encoding: SJIS   LineBreak: CRLF  Header: false\n" +
			"\n",
	},
	{
		Name:       "ShowObjects No Table is Loaded",
		Expr:       parser.ShowObjects{Type: parser.Identifier{Literal: "tables"}},
		Repository: filepath.Join(TestDir, "test_show_objects_empty"),
		Expect:     "No table is loaded",
	},
	{
		Name: "ShowObjects Views",
		Expr: parser.ShowObjects{Type: parser.Identifier{Literal: "views"}},
		Filter: &Filter{
			TempViews: TemporaryViewScopes{
				ViewMap{
					"VIEW1": &View{
						FileInfo: &FileInfo{
							Path:        "view1",
							IsTemporary: true,
						},
						Header: NewHeader("view1", []string{"column1", "column2"}),
					},
				},
				ViewMap{
					"VIEW1": &View{
						FileInfo: &FileInfo{
							Path:        "view1",
							IsTemporary: true,
						},
						Header: NewHeader("view1", []string{"column1", "column2", "column3"}),
					},
					"VIEW2": &View{
						FileInfo: &FileInfo{
							Path:        "view2",
							IsTemporary: true,
						},
						Header: NewHeader("view2", []string{"column1", "column2"}),
					},
				},
			},
		},
		UncommittedViews: &UncommittedViewMap{
			Created: map[string]*FileInfo{},
			Updated: map[string]*FileInfo{
				"VIEW2": {Path: "view2", IsTemporary: true},
			},
		},
		Expect: "\n" +
			" Views (Uncommitted: 1 View)\n" +
			"------------------------------\n" +
			" view1\n" +
			"     Fields: column1, column2\n" +
			" *Updated* view2\n" +
			"     Fields: column1, column2\n" +
			"\n",
	},
	{
		Name:   "ShowObjects Views Empty",
		Expr:   parser.ShowObjects{Type: parser.Identifier{Literal: "views"}},
		Expect: "No view is declared",
	},
	{
		Name: "ShowObjects Cursors",
		Expr: parser.ShowObjects{Type: parser.Identifier{Literal: "cursors"}},
		Filter: &Filter{
			Cursors: CursorScopes{
				{
					"CUR": &Cursor{
						name:  "cur",
						query: selectQueryForCursorTest,
					},
					"CUR2": &Cursor{
						name:  "cur2",
						query: selectQueryForCursorTest,
						view: &View{
							RecordSet: RecordSet{
								NewRecord([]value.Primary{
									value.NewInteger(1),
									value.NewString("a"),
								}),
								NewRecord([]value.Primary{
									value.NewInteger(2),
									value.NewString("b"),
								}),
							},
						},
						fetched: false,
						index:   -1,
					},
					"CUR3": &Cursor{
						name:  "cur3",
						query: selectQueryForCursorTest,
						view: &View{
							RecordSet: RecordSet{
								NewRecord([]value.Primary{
									value.NewInteger(1),
									value.NewString("a"),
								}),
								NewRecord([]value.Primary{
									value.NewInteger(2),
									value.NewString("b"),
								}),
							},
						},
						fetched: true,
						index:   1,
					},
					"CUR4": &Cursor{
						name:  "cur4",
						query: selectQueryForCursorTest,
						view: &View{
							RecordSet: RecordSet{
								NewRecord([]value.Primary{
									value.NewInteger(1),
									value.NewString("a"),
								}),
								NewRecord([]value.Primary{
									value.NewInteger(2),
									value.NewString("b"),
								}),
							},
						},
						fetched: true,
						index:   2,
					},
				},
			},
		},
		Expect: "\n" +
			"                               Cursors\n" +
			"---------------------------------------------------------------------\n" +
			" cur\n" +
			"     Status: Closed\n" +
			"     Query: select column1, column2 from table1\n" +
			" cur2\n" +
			"     Status: Open    Number of Rows: 2         Pointer: UNKNOWN\n" +
			"     Query: select column1, column2 from table1\n" +
			" cur3\n" +
			"     Status: Open    Number of Rows: 2         Pointer: 1\n" +
			"     Query: select column1, column2 from table1\n" +
			" cur4\n" +
			"     Status: Open    Number of Rows: 2         Pointer: Out of Range\n" +
			"     Query: select column1, column2 from table1\n" +
			"\n",
	},
	{
		Name:   "ShowObjects Cursors Empty",
		Expr:   parser.ShowObjects{Type: parser.Identifier{Literal: "cursors"}},
		Expect: "No cursor is declared",
	},
	{
		Name: "ShowObjects Functions",
		Expr: parser.ShowObjects{Type: parser.Identifier{Literal: "functions"}},
		Filter: &Filter{
			Functions: UserDefinedFunctionScopes{
				UserDefinedFunctionMap{
					"USERFUNC1": &UserDefinedFunction{
						Name: parser.Identifier{Literal: "userfunc1"},
						Parameters: []parser.Variable{
							{Name: "arg1"},
						},
						Statements: []parser.Statement{
							parser.Print{Value: parser.Variable{Name: "arg1"}},
						},
					},
				},
				UserDefinedFunctionMap{
					"USERAGGFUNC": &UserDefinedFunction{
						Name: parser.Identifier{Literal: "useraggfunc"},
						Parameters: []parser.Variable{
							{Name: "arg1"},
							{Name: "arg2"},
						},
						Defaults: map[string]parser.QueryExpression{
							"arg2": parser.NewIntegerValue(1),
						},
						IsAggregate:  true,
						RequiredArgs: 1,
						Cursor:       parser.Identifier{Literal: "column1"},
						Statements: []parser.Statement{
							parser.Print{Value: parser.Variable{Name: "var1"}},
						},
					},
				},
			},
		},
		Expect: "\n" +
			"  Scala Functions\n" +
			"-------------------\n" +
			" userfunc1 (@arg1)\n" +
			"\n" +
			"           Aggregate Functions\n" +
			"-----------------------------------------\n" +
			" useraggfunc (column1, @arg1, @arg2 = 1)\n" +
			"\n",
	},
	{
		Name:   "ShowObjects Functions Empty",
		Expr:   parser.ShowObjects{Type: parser.Identifier{Literal: "functions"}},
		Expect: "No function is declared",
	},
	{
		Name:       "ShowObjects Flags",
		Expr:       parser.ShowObjects{Type: parser.Identifier{Literal: "flags"}},
		Repository: ".",
		Expect: "\n" +
			"                      Flags\n" +
			"--------------------------------------------------\n" +
			"                @@REPOSITORY: .\n" +
			"                  @@TIMEZONE: UTC\n" +
			"           @@DATETIME_FORMAT: (not set)\n" +
			"              @@WAIT_TIMEOUT: 15\n" +
			"             @@IMPORT_FORMAT: CSV\n" +
			"                 @@DELIMITER: ','\n" +
			"       @@DELIMITER_POSITIONS: SPACES\n" +
			"                @@JSON_QUERY: (empty)\n" +
			"                  @@ENCODING: UTF8\n" +
			"                 @@NO_HEADER: false\n" +
			"              @@WITHOUT_NULL: false\n" +
			"                    @@FORMAT: CSV\n" +
			"            @@WRITE_ENCODING: UTF8\n" +
			"           @@WRITE_DELIMITER: ','\n" +
			" @@WRITE_DELIMITER_POSITIONS: (ignored) SPACES\n" +
			"            @@WITHOUT_HEADER: false\n" +
			"                @@LINE_BREAK: LF\n" +
			"               @@ENCLOSE_ALL: false\n" +
			"               @@JSON_ESCAPE: (ignored) BACKSLASH\n" +
			"              @@PRETTY_PRINT: (ignored) false\n" +
			"       @@EAST_ASIAN_ENCODING: (ignored) false\n" +
			"    @@COUNT_DIACRITICAL_SIGN: (ignored) false\n" +
			"         @@COUNT_FORMAT_CODE: (ignored) false\n" +
			"                     @@COLOR: false\n" +
			"                     @@QUIET: false\n" +
			"                       @@CPU: " + strconv.Itoa(cmd.GetFlags().CPU) + "\n" +
			"                     @@STATS: false\n" +
			"\n",
	},
	{
		Name:       "ShowObjects Runtime Information",
		Expr:       parser.ShowObjects{Type: parser.Identifier{Literal: "runinfo"}},
		Repository: ".",
		Expect: "\n" +
			strings.Repeat(" ", (calcShowRuninfoWidth(GetWD())-19)/2) + "Runtime Information\n" +
			strings.Repeat("-", calcShowRuninfoWidth(GetWD())) + "\n" +
			"       @#UNCOMMITTED: false\n" +
			"           @#CREATED: 0\n" +
			"           @#UPDATED: 0\n" +
			"     @#UPDATED_VIEWS: 0\n" +
			"     @#LOADED_TABLES: 0\n" +
			" @#WORKING_DIRECTORY: " + GetWD() + "\n" +
			"           @#VERSION: v1.0.0\n" +
			"\n",
	},
	{
		Name:  "ShowObjects Invalid Object Type",
		Expr:  parser.ShowObjects{Type: parser.Identifier{Literal: "invalid"}},
		Error: "[L:- C:-] object type invalid is invalid",
	},
}

func TestShowObjects(t *testing.T) {
	initCmdFlag()
	flags := cmd.GetFlags()

	for _, v := range showObjectsTests {
		flags.Repository = v.Repository
		flags.ImportFormat = v.ImportFormat
		flags.Delimiter = ','
		if v.Delimiter != 0 {
			flags.Delimiter = v.Delimiter
		}
		flags.DelimiterPositions = v.DelimiterPositions
		flags.SingleLine = v.SingleLine
		flags.JsonQuery = v.JsonQuery
		flags.WriteDelimiter = ','
		if v.WriteDelimiter != 0 {
			flags.WriteDelimiter = v.WriteDelimiter
		}
		flags.WriteDelimiterPositions = v.WriteDelimiterPositions
		flags.WriteAsSingleLine = v.WriteAsSingleLine
		flags.Format = v.Format
		ViewCache.Clean()
		if 0 < len(v.ViewCache) {
			ViewCache = v.ViewCache
		}
		if v.UncommittedViews == nil {
			UncommittedViews = NewUncommittedViewMap()
		} else {
			UncommittedViews = v.UncommittedViews
		}

		var filter *Filter
		if v.Filter != nil {
			filter = v.Filter
		} else {
			filter = NewEmptyFilter()
		}

		result, err := ShowObjects(v.Expr, filter)
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
		if result != v.Expect {
			t.Errorf("%s: result = %s, want %s", v.Name, result, v.Expect)
		}
	}
	ReleaseResources()
	UncommittedViews.Clean()
}

var showFieldsTests = []struct {
	Name             string
	Expr             parser.ShowFields
	Filter           *Filter
	ViewCache        ViewMap
	UncommittedViews *UncommittedViewMap
	Expect           string
	Error            string
}{
	{
		Name: "ShowFields Temporary Table",
		Expr: parser.ShowFields{
			Type:  parser.Identifier{Literal: "fields"},
			Table: parser.Identifier{Literal: "view1"},
		},
		Filter: &Filter{
			TempViews: TemporaryViewScopes{
				ViewMap{
					"VIEW1": &View{
						Header: NewHeader("view1", []string{"column1", "column2"}),
						FileInfo: &FileInfo{
							Path:        "view1",
							IsTemporary: true,
						},
					},
				},
			},
		},
		Expect: "\n" +
			" Fields in view1\n" +
			"-----------------\n" +
			" Type: View\n" +
			" Status: Fixed\n" +
			" Fields:\n" +
			"   1. column1\n" +
			"   2. column2\n" +
			"\n",
	},
	{
		Name: "ShowFields Updated Temporary Table",
		Expr: parser.ShowFields{
			Type:  parser.Identifier{Literal: "fields"},
			Table: parser.Identifier{Literal: "view1"},
		},
		Filter: &Filter{
			TempViews: TemporaryViewScopes{
				ViewMap{
					"VIEW1": &View{
						Header: NewHeader("view1", []string{"column1", "column2"}),
						FileInfo: &FileInfo{
							Path:        "view1",
							IsTemporary: true,
						},
					},
				},
			},
		},
		UncommittedViews: &UncommittedViewMap{
			Created: map[string]*FileInfo{},
			Updated: map[string]*FileInfo{
				"VIEW1": {Path: "view1", IsTemporary: true},
			},
		},
		Expect: "\n" +
			" Fields in view1\n" +
			"-----------------\n" +
			" Type: View\n" +
			" Status: Updated\n" +
			" Fields:\n" +
			"   1. column1\n" +
			"   2. column2\n" +
			"\n",
	},
	{
		Name: "ShowFields",
		Expr: parser.ShowFields{
			Type: parser.Identifier{Literal: "fields"},
			Table: parser.TableObject{
				Type:          parser.Identifier{Literal: "csv"},
				FormatElement: parser.NewStringValue(","),
				Path:          parser.Identifier{Literal: "show_fields_create.csv"},
			},
		},
		ViewCache: ViewMap{
			strings.ToUpper(GetTestFilePath("show_fields_create.csv")): &View{
				Header: NewHeader("show_fields_create", []string{"column1", "column2"}),
				FileInfo: &FileInfo{
					Path:      GetTestFilePath("show_fields_create.csv"),
					Delimiter: ',',
					Format:    cmd.CSV,
					Encoding:  text.UTF8,
					LineBreak: text.LF,
					NoHeader:  false,
				},
			},
		},
		UncommittedViews: &UncommittedViewMap{
			Created: map[string]*FileInfo{
				strings.ToUpper(GetTestFilePath("show_fields_create.csv")): {Path: "show_fields_create.csv"},
			},
			Updated: map[string]*FileInfo{},
		},
		Expect: "\n" +
			strings.Repeat(" ", (calcShowFieldsWidth("show_fields_create.csv", "show_fields_create.csv", 10)-(10+len("show_fields_create.csv")))/2) + "Fields in show_fields_create.csv\n" +
			strings.Repeat("-", calcShowFieldsWidth("show_fields_create.csv", "show_fields_create.csv", 10)) + "\n" +
			" Type: Table\n" +
			" Path: " + GetTestFilePath("show_fields_create.csv") + "\n" +
			" Format: CSV      Delimiter: ','   Enclose All: false\n" +
			" Encoding: UTF8   LineBreak: LF    Header: true\n" +
			" Status: Created\n" +
			" Fields:\n" +
			"   1. column1\n" +
			"   2. column2\n" +
			"\n",
	},
	{
		Name: "ShowFields Created Table",
		Expr: parser.ShowFields{
			Type:  parser.Identifier{Literal: "fields"},
			Table: parser.Identifier{Literal: "show_fields_create.csv"},
		},
		ViewCache: ViewMap{
			strings.ToUpper(GetTestFilePath("show_fields_create.csv")): &View{
				Header: NewHeader("show_fields_create", []string{"column1", "column2"}),
				FileInfo: &FileInfo{
					Path:      GetTestFilePath("show_fields_create.csv"),
					Delimiter: ',',
					Format:    cmd.CSV,
					Encoding:  text.UTF8,
					LineBreak: text.LF,
					NoHeader:  false,
				},
			},
		},
		UncommittedViews: &UncommittedViewMap{
			Created: map[string]*FileInfo{
				strings.ToUpper(GetTestFilePath("show_fields_create.csv")): {Path: "show_fields_create.csv"},
			},
			Updated: map[string]*FileInfo{},
		},
		Expect: "\n" +
			strings.Repeat(" ", (calcShowFieldsWidth("show_fields_create.csv", "show_fields_create.csv", 10)-(10+len("show_fields_create.csv")))/2) + "Fields in show_fields_create.csv\n" +
			strings.Repeat("-", calcShowFieldsWidth("show_fields_create.csv", "show_fields_create.csv", 10)) + "\n" +
			" Type: Table\n" +
			" Path: " + GetTestFilePath("show_fields_create.csv") + "\n" +
			" Format: CSV      Delimiter: ','   Enclose All: false\n" +
			" Encoding: UTF8   LineBreak: LF    Header: true\n" +
			" Status: Created\n" +
			" Fields:\n" +
			"   1. column1\n" +
			"   2. column2\n" +
			"\n",
	},
	{
		Name: "ShowFields Updated Table",
		Expr: parser.ShowFields{
			Type:  parser.Identifier{Literal: "fields"},
			Table: parser.Identifier{Literal: "show_fields_update.csv"},
		},
		ViewCache: ViewMap{
			strings.ToUpper(GetTestFilePath("show_fields_update.csv")): &View{
				Header: NewHeader("show_fields_update", []string{"column1", "column2"}),
				FileInfo: &FileInfo{
					Path:      GetTestFilePath("show_fields_update.csv"),
					Delimiter: ',',
					Format:    cmd.CSV,
					Encoding:  text.UTF8,
					LineBreak: text.LF,
					NoHeader:  false,
				},
			},
		},
		UncommittedViews: &UncommittedViewMap{
			Created: map[string]*FileInfo{},
			Updated: map[string]*FileInfo{
				strings.ToUpper(GetTestFilePath("show_fields_update.csv")): {Path: "show_fields_updated.csv"},
			},
		},
		Expect: "\n" +
			strings.Repeat(" ", (calcShowFieldsWidth("show_fields_update.csv", "show_fields_update.csv", 10)-(10+len("show_fields_update.csv")))/2) + "Fields in show_fields_update.csv\n" +
			strings.Repeat("-", calcShowFieldsWidth("show_fields_create.csv", "show_fields_update.csv", 10)) + "\n" +
			" Type: Table\n" +
			" Path: " + GetTestFilePath("show_fields_update.csv") + "\n" +
			" Format: CSV      Delimiter: ','   Enclose All: false\n" +
			" Encoding: UTF8   LineBreak: LF    Header: true\n" +
			" Status: Updated\n" +
			" Fields:\n" +
			"   1. column1\n" +
			"   2. column2\n" +
			"\n",
	},
	{
		Name: "ShowFields Load Error",
		Expr: parser.ShowFields{
			Type:  parser.Identifier{Literal: "fields"},
			Table: parser.Identifier{Literal: "notexist"},
		},
		Error: "[L:- C:-] file notexist does not exist",
	},
	{
		Name: "ShowFields Invalid Object Type",
		Expr: parser.ShowFields{
			Type:  parser.Identifier{Literal: "invalid"},
			Table: parser.Identifier{Literal: "table2"},
		},
		Error: "[L:- C:-] object type invalid is invalid",
	},
}

func calcShowFieldsWidth(fileName string, fileNameInTitle string, prefixLen int) int {
	w := 54
	pathLen := 8 + len(GetTestFilePath(fileName))
	titleLen := prefixLen + len(fileNameInTitle)

	if w < titleLen {
		w = titleLen
	}
	if w < pathLen {
		w = pathLen
	}
	if 75 < w {
		w = 75
	}
	return w
}

func calcShowRuninfoWidth(wd string) int {
	w := 28
	pathLen := 22 + len(wd)
	if w < pathLen {
		w = pathLen
	}
	w++
	if 75 < w {
		w = 75
	}
	return w
}

func TestShowFields(t *testing.T) {
	initCmdFlag()
	flags := cmd.GetFlags()
	flags.Repository = TestDir

	for _, v := range showFieldsTests {
		ViewCache.Clean()
		if 0 < len(v.ViewCache) {
			ViewCache = v.ViewCache
		}
		if v.UncommittedViews == nil {
			UncommittedViews = NewUncommittedViewMap()
		} else {
			UncommittedViews = v.UncommittedViews
		}

		var filter *Filter
		if v.Filter != nil {
			filter = v.Filter
		} else {
			filter = NewEmptyFilter()
		}

		result, err := ShowFields(v.Expr, filter)
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
		if result != v.Expect {
			t.Errorf("%s: result = %s, want %s", v.Name, result, v.Expect)
		}
	}
	ReleaseResources()
	UncommittedViews.Clean()
}

var setEnvVarTests = []struct {
	Name   string
	Expr   parser.SetEnvVar
	Expect string
	Error  string
}{
	{
		Name: "Set Environment Variable",
		Expr: parser.SetEnvVar{
			EnvVar: parser.EnvironmentVariable{
				Name: "CSVQ_SET_ENV_TEST",
			},
			Value: parser.NewStringValue("foo"),
		},
		Expect: "foo",
	},
	{
		Name: "Set Environment Variable with Identifier",
		Expr: parser.SetEnvVar{
			EnvVar: parser.EnvironmentVariable{
				Name: "CSVQ_SET_ENV_TEST",
			},
			Value: parser.Identifier{Literal: "bar"},
		},
		Expect: "bar",
	},
	{
		Name: "Set Environment Variable with Null",
		Expr: parser.SetEnvVar{
			EnvVar: parser.EnvironmentVariable{
				Name: "CSVQ_SET_ENV_TEST",
			},
			Value: parser.NewNullValue(),
		},
		Expect: "",
	},
	{
		Name: "Set Environment Variable Evaluation Error",
		Expr: parser.SetEnvVar{
			EnvVar: parser.EnvironmentVariable{
				Name: "CSVQ_SET_ENV_TEST",
			},
			Value: parser.FieldReference{Column: parser.Identifier{Literal: "err"}},
		},
		Error: "[L:- C:-] field err does not exist",
	},
}

func TestSetEnvVar(t *testing.T) {
	filter := NewEmptyFilter()

	for _, v := range setEnvVarTests {
		err := SetEnvVar(v.Expr, filter)

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

		val := os.Getenv(v.Expr.EnvVar.Name)
		if val != v.Expect {
			t.Errorf("%s: value = %s, want %s", v.Name, val, v.Expect)
		}
	}
}

var syntaxTests = []struct {
	Expr   parser.Syntax
	Expect string
}{
	{
		Expr: parser.Syntax{},
		Expect: "\n" +
			"        Contents\n" +
			"-------------------------\n" +
			" SELECT Statement\n" +
			"     WITH Clause\n" +
			"     SELECT Clause\n" +
			" INSERT Statement\n" +
			" UPDATE Statement\n" +
			" DELETE Statement\n" +
			" Operators\n" +
			"     Operator Precedence\n" +
			"     String Operators\n" +
			"\n",
	},
	{
		Expr: parser.Syntax{Keywords: []parser.QueryExpression{parser.NewStringValue("select clause")}},
		Expect: "\n" +
			"                Search: select clause\n" +
			"------------------------------------------------------\n" +
			" SELECT Clause\n" +
			"     select_clause\n" +
			"         : SELECT [DISTINCT] <field> [, <field> ...] \n" +
			"\n" +
			"     field\n" +
			"         : <value> \n" +
			"         | <value> AS alias \n" +
			"\n" +
			"\n",
	},
	{
		Expr: parser.Syntax{Keywords: []parser.QueryExpression{parser.NewStringValue(" select  "), parser.NewStringValue("clause")}},
		Expect: "\n" +
			"                Search: select clause\n" +
			"------------------------------------------------------\n" +
			" SELECT Clause\n" +
			"     select_clause\n" +
			"         : SELECT [DISTINCT] <field> [, <field> ...] \n" +
			"\n" +
			"     field\n" +
			"         : <value> \n" +
			"         | <value> AS alias \n" +
			"\n" +
			"\n",
	},
	{
		Expr: parser.Syntax{Keywords: []parser.QueryExpression{parser.FieldReference{Column: parser.Identifier{Literal: "select clause"}}}},
		Expect: "\n" +
			"                Search: select clause\n" +
			"------------------------------------------------------\n" +
			" SELECT Clause\n" +
			"     select_clause\n" +
			"         : SELECT [DISTINCT] <field> [, <field> ...] \n" +
			"\n" +
			"     field\n" +
			"         : <value> \n" +
			"         | <value> AS alias \n" +
			"\n" +
			"\n",
	},
	{
		Expr: parser.Syntax{Keywords: []parser.QueryExpression{parser.NewStringValue("operator prec")}},
		Expect: "\n" +
			"         Search: operator prec\n" +
			"---------------------------------------\n" +
			" Operator Precedence\n" +
			"     Operator Precedence Description. \n" +
			"\n" +
			"\n",
	},
	{
		Expr: parser.Syntax{Keywords: []parser.QueryExpression{parser.NewStringValue("string  op")}},
		Expect: "\n" +
			"       Search: string op\n" +
			"-------------------------------\n" +
			" String Operators\n" +
			"     concatenation\n" +
			"         : <value> || <value> \n" +
			"\n" +
			"         description \n" +
			"\n" +
			"\n",
	},
}

func TestSyntax(t *testing.T) {
	origSyntax := syntax.CsvqSyntax

	syntax.CsvqSyntax = []syntax.Expression{
		{
			Label: "SELECT Statement",
			Grammar: []syntax.Definition{
				{
					Name: "select_statement",
					Group: []syntax.Grammar{
						{syntax.Option{syntax.Link("with_clause")}, syntax.Link("select_query")},
					},
				},
				{
					Name: "select_query",
					Group: []syntax.Grammar{
						{syntax.Link("select_entity"), syntax.Option{syntax.Link("order_by_clause")}, syntax.Option{syntax.Link("limit_clause")}, syntax.Option{syntax.Link("offset_clause")}},
					},
				},
			},
			Children: []syntax.Expression{
				{
					Label: "WITH Clause",
					Grammar: []syntax.Definition{
						{
							Name: "with_clause",
							Group: []syntax.Grammar{
								{syntax.Keyword("WITH"), syntax.ContinuousOption{syntax.Link("common_table_expression")}},
							},
						},
						{
							Name: "common_table_expression",
							Group: []syntax.Grammar{
								{syntax.Option{syntax.Keyword("RECURSIVE")}, syntax.Identifier("table_name"), syntax.Option{syntax.Parentheses{syntax.ContinuousOption{syntax.Identifier("column_name")}}}, syntax.Keyword("AS"), syntax.Parentheses{syntax.Link("select_query")}},
							},
						},
					},
				},
				{
					Label: "SELECT Clause",
					Grammar: []syntax.Definition{
						{
							Name: "select_clause",
							Group: []syntax.Grammar{
								{syntax.Keyword("SELECT"), syntax.Option{syntax.Keyword("DISTINCT")}, syntax.ContinuousOption{syntax.Link("field")}},
							},
						},
						{
							Name: "field",
							Group: []syntax.Grammar{
								{syntax.Link("value")},
								{syntax.Link("value"), syntax.Keyword("AS"), syntax.Identifier("alias")},
							},
						},
					},
				},
			},
		},
		{
			Label: "INSERT Statement",
			Grammar: []syntax.Definition{
				{
					Name: "insert_statement",
					Group: []syntax.Grammar{
						{syntax.Option{syntax.Link("with_clause")}, syntax.Link("insert_query")},
					},
				},
				{
					Name: "insert_query",
					Group: []syntax.Grammar{
						{syntax.Keyword("INSERT"), syntax.Keyword("INTO"), syntax.Identifier("table_name"), syntax.Option{syntax.Parentheses{syntax.ContinuousOption{syntax.Identifier("column_name")}}}, syntax.Keyword("VALUES"), syntax.ContinuousOption{syntax.Link("row_value")}},
						{syntax.Keyword("INSERT"), syntax.Keyword("INTO"), syntax.Identifier("table_name"), syntax.Option{syntax.Parentheses{syntax.ContinuousOption{syntax.Identifier("column_name")}}}, syntax.Link("select_query")},
					},
				},
			},
		},
		{
			Label:   "UPDATE Statement",
			Grammar: []syntax.Definition{},
		},
		{
			Label:   "DELETE Statement",
			Grammar: []syntax.Definition{},
		},
		{
			Label: "Operators",
			Children: []syntax.Expression{
				{
					Label: "Operator Precedence",
					Description: syntax.Description{
						Template: "Operator Precedence Description.",
					},
				},
				{
					Label: "String Operators",
					Grammar: []syntax.Definition{
						{
							Name: "concatenation",
							Group: []syntax.Grammar{
								{syntax.Link("value"), syntax.Keyword("||"), syntax.Link("value")},
							},
							Description: syntax.Description{Template: "description"},
						},
					},
				},
			},
		},
	}

	filter := NewEmptyFilter()

	for _, v := range syntaxTests {
		result := Syntax(v.Expr, filter)
		if result != v.Expect {
			t.Errorf("result = %s, want %s for %v", result, v.Expect, v.Expr)
		}
	}

	syntax.CsvqSyntax = origSyntax
}
