package query

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

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
		Error: "variable @var is undeclared",
	},
}

func TestEcho(t *testing.T) {
	scope := NewReferenceScope(TestTx)

	for _, v := range echoTests {
		result, err := Echo(context.Background(), scope, v.Expr)
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
		Result: "'foo'",
	},
	{
		Name: "Print Error",
		Expr: parser.Print{
			Value: parser.Variable{
				Name: "var",
			},
		},
		Error: "variable @var is undeclared",
	},
}

func TestPrint(t *testing.T) {
	scope := NewReferenceScope(TestTx)

	for _, v := range printTests {
		result, err := Print(context.Background(), scope, v.Expr)
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
		Result: "printf test: value1 'str', value2 '1'",
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
		Error: "variable @var is undeclared",
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
		Error: "variable @var is undeclared",
	},
	{
		Name: "Printf Less Values Error",
		Expr: parser.Printf{
			Format: parser.NewStringValue("printf test: value1 %s, value2 %s %%"),
			Values: []parser.QueryExpression{
				parser.NewStringValue("str"),
			},
		},
		Error: "number of replace values does not match",
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
		Error: "number of replace values does not match",
	},
}

func TestPrintf(t *testing.T) {
	scope := NewReferenceScope(TestTx)

	for _, v := range printfTests {
		result, err := Printf(context.Background(), scope, v.Expr)
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
		Error: "field ident does not exist",
	},
	{
		Name: "Source File Invalid File Path Error",
		Expr: parser.Source{
			FilePath: parser.NewNullValue(),
		},
		Error: "NULL is a invalid file path",
	},
	{
		Name: "Source File Empty File Path Error",
		Expr: parser.Source{
			FilePath: parser.Identifier{Literal: "", Quoted: true},
		},
		Error: "`` is a invalid file path",
	},
	{
		Name: "Source File Not Exist Error",
		Expr: parser.Source{
			FilePath: parser.NewStringValue(GetTestFilePath("notexist.sql")),
		},
		Error: fmt.Sprintf("file %s does not exist", GetTestFilePath("notexist.sql")),
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
	scope := NewReferenceScope(TestTx)

	for _, v := range sourceTests {
		result, err := Source(context.Background(), scope, v.Expr)
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
		Error: "field notexist does not exist",
	},
	{
		Name: "ParseExecuteStatements Format Error",
		Expr: parser.Execute{
			BaseExpr:   parser.NewBaseExpr(parser.Token{}),
			Statements: parser.NewStringValue("print %q;"),
		},
		Error: "number of replace values does not match",
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
		Error: "field notexist does not exist",
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
	scope := NewReferenceScope(TestTx)

	for _, v := range parseExecuteStatementsTests {
		result, err := ParseExecuteStatements(context.Background(), scope, v.Expr)
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
			Flag:  parser.Flag{Name: "repository"},
			Value: parser.NewStringValue(TestDir),
		},
	},
	{
		Name: "Set Timezone",
		Expr: parser.SetFlag{
			Flag:  parser.Flag{Name: "timezone"},
			Value: parser.NewStringValue("utc"),
		},
	},
	{
		Name: "Set DatetimeFormat",
		Expr: parser.SetFlag{
			Flag:  parser.Flag{Name: "datetime_format"},
			Value: parser.NewStringValue("%Y%m%d"),
		},
	},
	{
		Name: "Set AnsiQuotes",
		Expr: parser.SetFlag{
			Flag:  parser.Flag{Name: "ansi_quotes"},
			Value: parser.NewTernaryValueFromString("true"),
		},
	},
	{
		Name: "Set StrictEqual",
		Expr: parser.SetFlag{
			Flag:  parser.Flag{Name: "strict_equal"},
			Value: parser.NewTernaryValueFromString("true"),
		},
	},
	{
		Name: "Set WaitTimeout",
		Expr: parser.SetFlag{
			Flag:  parser.Flag{Name: "wait_timeout"},
			Value: parser.NewFloatValue(15),
		},
	},
	{
		Name: "Set Delimiter",
		Expr: parser.SetFlag{
			Flag:  parser.Flag{Name: "delimiter"},
			Value: parser.NewStringValue("\\t"),
		},
	},
	{
		Name: "Set JsonQuery",
		Expr: parser.SetFlag{
			Flag:  parser.Flag{Name: "json_query"},
			Value: parser.NewStringValue("{}"),
		},
	},
	{
		Name: "Set Encoding",
		Expr: parser.SetFlag{
			Flag:  parser.Flag{Name: "encoding"},
			Value: parser.NewStringValue("SJIS"),
		},
	},
	{
		Name: "Set NoHeader",
		Expr: parser.SetFlag{
			Flag:  parser.Flag{Name: "no_header"},
			Value: parser.NewTernaryValueFromString("true"),
		},
	},
	{
		Name: "Set WithoutNull",
		Expr: parser.SetFlag{
			Flag:  parser.Flag{Name: "without_null"},
			Value: parser.NewTernaryValueFromString("true"),
		},
	},
	{
		Name: "Set Format",
		Expr: parser.SetFlag{
			Flag:  parser.Flag{Name: "format"},
			Value: parser.NewStringValue("json"),
		},
	},
	{
		Name: "Set WriteEncoding",
		Expr: parser.SetFlag{
			Flag:  parser.Flag{Name: "write_encoding"},
			Value: parser.NewStringValue("SJIS"),
		},
	},
	{
		Name: "Set WriteDelimiter",
		Expr: parser.SetFlag{
			Flag:  parser.Flag{Name: "write_delimiter"},
			Value: parser.NewStringValue("\\t"),
		},
	},
	{
		Name: "Set WithoutHeader",
		Expr: parser.SetFlag{
			Flag:  parser.Flag{Name: "without_header"},
			Value: parser.NewTernaryValueFromString("true"),
		},
	},
	{
		Name: "Set lineBreak",
		Expr: parser.SetFlag{
			Flag:  parser.Flag{Name: "line_break"},
			Value: parser.NewStringValue("CRLF"),
		},
	},
	{
		Name: "Set EncloseAll",
		Expr: parser.SetFlag{
			Flag:  parser.Flag{Name: "enclose_all"},
			Value: parser.NewTernaryValueFromString("true"),
		},
	},
	{
		Name: "Set JsonEscape",
		Expr: parser.SetFlag{
			Flag:  parser.Flag{Name: "json_escape"},
			Value: parser.NewStringValue("hex"),
		},
	},
	{
		Name: "Set PrettyPrint",
		Expr: parser.SetFlag{
			Flag:  parser.Flag{Name: "pretty_print"},
			Value: parser.NewTernaryValueFromString("true"),
		},
	},
	{
		Name: "Set Strip Ending Line Break",
		Expr: parser.SetFlag{
			Flag:  parser.Flag{Name: "strip_ending_line_break"},
			Value: parser.NewTernaryValueFromString("true"),
		},
	},
	{
		Name: "Set EastAsianEncoding",
		Expr: parser.SetFlag{
			Flag:  parser.Flag{Name: "east_asian_encoding"},
			Value: parser.NewTernaryValueFromString("true"),
		},
	},
	{
		Name: "Set CountDiacriticalSign",
		Expr: parser.SetFlag{
			Flag:  parser.Flag{Name: "count_diacritical_sign"},
			Value: parser.NewTernaryValueFromString("true"),
		},
	},
	{
		Name: "Set CountFormatCode",
		Expr: parser.SetFlag{
			Flag:  parser.Flag{Name: "count_format_code"},
			Value: parser.NewTernaryValueFromString("true"),
		},
	},
	{
		Name: "Set Color",
		Expr: parser.SetFlag{
			Flag:  parser.Flag{Name: "color"},
			Value: parser.NewTernaryValueFromString("true"),
		},
	},
	{
		Name: "Set Quiet",
		Expr: parser.SetFlag{
			Flag:  parser.Flag{Name: "quiet"},
			Value: parser.NewTernaryValueFromString("true"),
		},
	},
	{
		Name: "Set LimitRecursion",
		Expr: parser.SetFlag{
			Flag:  parser.Flag{Name: "limit_recursion"},
			Value: parser.NewIntegerValue(int64(10)),
		},
	},
	{
		Name: "Set CPU",
		Expr: parser.SetFlag{
			Flag:  parser.Flag{Name: "cpu"},
			Value: parser.NewIntegerValue(int64(runtime.NumCPU())),
		},
	},
	{
		Name: "Set Stats",
		Expr: parser.SetFlag{
			Flag:  parser.Flag{Name: "stats"},
			Value: parser.NewTernaryValueFromString("true"),
		},
	},
	{
		Name: "Set Encoding with Identifier",
		Expr: parser.SetFlag{
			Flag:  parser.Flag{Name: "encoding"},
			Value: parser.Identifier{Literal: "sjis"},
		},
	},
	{
		Name: "Set Delimiter Evaluation Error",
		Expr: parser.SetFlag{
			Flag:  parser.Flag{Name: "delimiter"},
			Value: parser.FieldReference{Column: parser.Identifier{Literal: "err"}},
		},
		Error: "field err does not exist",
	},
	{
		Name: "Set Delimiter Value Error",
		Expr: parser.SetFlag{
			Flag:  parser.Flag{Name: "delimiter"},
			Value: parser.NewTernaryValueFromString("true"),
		},
		Error: "TRUE for @@DELIMITER is not allowed",
	},
	{
		Name: "Set WaitTimeout Value Error",
		Expr: parser.SetFlag{
			Flag:  parser.Flag{Name: "wait_timeout"},
			Value: parser.NewTernaryValueFromString("true"),
		},
		Error: "TRUE for @@WAIT_TIMEOUT is not allowed",
	},
	{
		Name: "Set WithoutNull Value Error",
		Expr: parser.SetFlag{
			Flag:  parser.Flag{Name: "without_null"},
			Value: parser.NewStringValue("string"),
		},
		Error: "'string' for @@WITHOUT_NULL is not allowed",
	},
	{
		Name: "Set CPU Value Error",
		Expr: parser.SetFlag{
			Flag:  parser.Flag{Name: "cpu"},
			Value: parser.NewStringValue("invalid"),
		},
		Error: "'invalid' for @@CPU is not allowed",
	},
	{
		Name: "Invalid Flag Name Error",
		Expr: parser.SetFlag{
			Flag:  parser.Flag{Name: "invalid"},
			Value: parser.NewStringValue("string"),
		},
		Error: "@@INVALID is an unknown flag",
	},
	{
		Name: "Invalid Flag Value Error",
		Expr: parser.SetFlag{
			Flag:  parser.Flag{Name: "line_break"},
			Value: parser.NewStringValue("invalid"),
		},
		Error: "line-break must be one of CRLF|CR|LF",
	},
}

func TestSetFlag(t *testing.T) {
	defer initFlag(TestTx.Flags)

	ctx := context.Background()
	scope := NewReferenceScope(TestTx)

	for _, v := range setFlagTests {
		initFlag(TestTx.Flags)
		err := SetFlag(ctx, scope, v.Expr)
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
			Flag:  parser.Flag{Name: "datetime_format"},
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
			Flag:  parser.Flag{Name: "format"},
			Value: parser.NewStringValue("%Y%m%d"),
		},
		Init: func(flags *cmd.Flags) {
			flags.DatetimeFormat = []string{"%Y:%m:%d"}
		},
		Error: "add flag element syntax does not support @@FORMAT",
	},
	{
		Name: "Add Element Invalid Flag Name",
		Expr: parser.AddFlagElement{
			Flag:  parser.Flag{Name: "invalid"},
			Value: parser.NewStringValue("%Y%m%d"),
		},
		Init: func(flags *cmd.Flags) {
			flags.DatetimeFormat = []string{"%Y:%m:%d"}
		},
		Error: "@@INVALID is an unknown flag",
	},
}

func TestAddFlagElement(t *testing.T) {
	defer initFlag(TestTx.Flags)

	ctx := context.Background()
	scope := NewReferenceScope(TestTx)

	for _, v := range addFlagElementTests {
		initFlag(TestTx.Flags)
		v.Init(TestTx.Flags)

		err := AddFlagElement(ctx, scope, v.Expr)
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
		if !reflect.DeepEqual(TestTx.Flags, expect) {
			t.Errorf("%s: result = %v, want %v", v.Name, TestTx.Flags, expect)
		}
	}
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
			Flag:  parser.Flag{Name: "datetime_format"},
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
			Flag:  parser.Flag{Name: "datetime_format"},
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
			Flag:  parser.Flag{Name: "datetime_format"},
			Value: parser.NewNullValue(),
		},
		Init:  func(flags *cmd.Flags) {},
		Error: "NULL is an invalid value for @@DATETIME_FORMAT to specify the element",
	},
	{
		Name: "Remove Element Evaluation Error",
		Expr: parser.RemoveFlagElement{
			Flag:  parser.Flag{Name: "format"},
			Value: parser.FieldReference{Column: parser.Identifier{Literal: "err"}},
		},
		Init:  func(flags *cmd.Flags) {},
		Error: "field err does not exist",
	},
	{
		Name: "Remove Element Unsupported Flag Name",
		Expr: parser.RemoveFlagElement{
			Flag:  parser.Flag{Name: "format"},
			Value: parser.NewIntegerValue(1),
		},
		Init:  func(flags *cmd.Flags) {},
		Error: "remove flag element syntax does not support @@FORMAT",
	},
	{
		Name: "Remove Element Invalid Flag Name",
		Expr: parser.RemoveFlagElement{
			Flag:  parser.Flag{Name: "invalid"},
			Value: parser.NewIntegerValue(1),
		},
		Init:  func(flags *cmd.Flags) {},
		Error: "@@INVALID is an unknown flag",
	},
}

func TestRemoveFlagElement(t *testing.T) {
	defer initFlag(TestTx.Flags)

	ctx := context.Background()
	scope := NewReferenceScope(TestTx)

	for _, v := range removeFlagElementTests {
		initFlag(TestTx.Flags)
		v.Init(TestTx.Flags)

		err := RemoveFlagElement(ctx, scope, v.Expr)
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
		if !reflect.DeepEqual(TestTx.Flags, expect) {
			t.Errorf("%s: result = %v, want %v", v.Name, TestTx.Flags, expect)
		}
	}
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
			Flag: parser.Flag{Name: "repository"},
		},
		SetExprs: []parser.SetFlag{
			{
				Flag:  parser.Flag{Name: "repository"},
				Value: parser.NewStringValue(TestDir),
			},
		},
		Result: "\033[34;1m@@REPOSITORY:\033[0m \033[32m" + TestDir + "\033[0m",
	},
	{
		Name: "Show Repository Not Set",
		Expr: parser.ShowFlag{
			Flag: parser.Flag{Name: "repository"},
		},
		SetExprs: []parser.SetFlag{
			{
				Flag:  parser.Flag{Name: "repository"},
				Value: parser.NewStringValue(""),
			},
		},
		Result: "\033[34;1m@@REPOSITORY:\033[0m \033[90m(current dir: " + GetWD() + ")\033[0m",
	},
	{
		Name: "Show Timezone",
		Expr: parser.ShowFlag{
			Flag: parser.Flag{Name: "timezone"},
		},
		SetExprs: []parser.SetFlag{
			{
				Flag:  parser.Flag{Name: "timezone"},
				Value: parser.NewStringValue("UTC"),
			},
		},
		Result: "\033[34;1m@@TIMEZONE:\033[0m \033[32mUTC\033[0m",
	},
	{
		Name: "Show DatetimeFormat Not Set",
		Expr: parser.ShowFlag{
			Flag: parser.Flag{Name: "datetime_format"},
		},
		Result: "\033[34;1m@@DATETIME_FORMAT:\033[0m \033[90m(not set)\033[0m",
	},
	{
		Name: "Show DatetimeFormat",
		Expr: parser.ShowFlag{
			Flag: parser.Flag{Name: "datetime_format"},
		},
		SetExprs: []parser.SetFlag{
			{
				Flag:  parser.Flag{Name: "datetime_format"},
				Value: parser.NewStringValue("[\"%Y%m%d\", \"%Y%m%d %H%i%s\"]"),
			},
		},
		Result: "\033[34;1m@@DATETIME_FORMAT:\033[0m \033[32m[\"%Y%m%d\", \"%Y%m%d %H%i%s\"]\033[0m",
	},
	{
		Name: "Show AnsiQuotes",
		Expr: parser.ShowFlag{
			Flag: parser.Flag{Name: "ansi_quotes"},
		},
		SetExprs: []parser.SetFlag{
			{
				Flag:  parser.Flag{Name: "ansi_quotes"},
				Value: parser.NewTernaryValueFromString("true"),
			},
		},
		Result: "\033[34;1m@@ANSI_QUOTES:\033[0m \033[33;1mtrue\033[0m",
	},
	{
		Name: "Show StrictEqual",
		Expr: parser.ShowFlag{
			Flag: parser.Flag{Name: "strict_equal"},
		},
		SetExprs: []parser.SetFlag{
			{
				Flag:  parser.Flag{Name: "strict_equal"},
				Value: parser.NewTernaryValueFromString("true"),
			},
		},
		Result: "\033[34;1m@@STRICT_EQUAL:\033[0m \033[33;1mtrue\033[0m",
	},
	{
		Name: "Show WaitTimeout",
		Expr: parser.ShowFlag{
			Flag: parser.Flag{Name: "wait_timeout"},
		},
		SetExprs: []parser.SetFlag{
			{
				Flag:  parser.Flag{Name: "wait_timeout"},
				Value: parser.NewFloatValue(15),
			},
		},
		Result: "\033[34;1m@@WAIT_TIMEOUT:\033[0m \033[35m15\033[0m",
	},
	{
		Name: "Show Import Format",
		Expr: parser.ShowFlag{
			Flag: parser.Flag{Name: "import_format"},
		},
		SetExprs: []parser.SetFlag{
			{
				Flag:  parser.Flag{Name: "import_format"},
				Value: parser.NewStringValue("tsv"),
			},
		},
		Result: "\033[34;1m@@IMPORT_FORMAT:\033[0m \033[32mTSV\033[0m",
	},
	{
		Name: "Show Delimiter",
		Expr: parser.ShowFlag{
			Flag: parser.Flag{Name: "delimiter"},
		},
		SetExprs: []parser.SetFlag{
			{
				Flag:  parser.Flag{Name: "delimiter"},
				Value: parser.NewStringValue("\t"),
			},
		},
		Result: "\033[34;1m@@DELIMITER:\033[0m \033[32m'\\t'\033[0m",
	},
	{
		Name: "Show Delimiter Positions",
		Expr: parser.ShowFlag{
			Flag: parser.Flag{Name: "delimiter_positions"},
		},
		SetExprs: []parser.SetFlag{
			{
				Flag:  parser.Flag{Name: "delimiter_positions"},
				Value: parser.NewStringValue("s[2, 5, 10]"),
			},
		},
		Result: "\033[34;1m@@DELIMITER_POSITIONS:\033[0m \033[32mS[2, 5, 10]\033[0m",
	},
	{
		Name: "Show Delimiter Positions as spaces",
		Expr: parser.ShowFlag{
			Flag: parser.Flag{Name: "delimiter_positions"},
		},
		SetExprs: []parser.SetFlag{
			{
				Flag:  parser.Flag{Name: "delimiter_positions"},
				Value: parser.NewStringValue("SPACES"),
			},
		},
		Result: "\033[34;1m@@DELIMITER_POSITIONS:\033[0m \033[32mSPACES\033[0m",
	},
	{
		Name: "Show JsonQuery",
		Expr: parser.ShowFlag{
			Flag: parser.Flag{Name: "json_query"},
		},
		SetExprs: []parser.SetFlag{
			{
				Flag:  parser.Flag{Name: "json_query"},
				Value: parser.NewStringValue("{}"),
			},
		},
		Result: "\033[34;1m@@JSON_QUERY:\033[0m \033[32m{}\033[0m",
	},
	{
		Name: "Show JsonQuery Empty",
		Expr: parser.ShowFlag{
			Flag: parser.Flag{Name: "json_query"},
		},
		SetExprs: []parser.SetFlag{},
		Result:   "\033[34;1m@@JSON_QUERY:\033[0m \033[90m(empty)\033[0m",
	},
	{
		Name: "Show Encoding",
		Expr: parser.ShowFlag{
			Flag: parser.Flag{Name: "encoding"},
		},
		SetExprs: []parser.SetFlag{
			{
				Flag:  parser.Flag{Name: "encoding"},
				Value: parser.NewStringValue("SJIS"),
			},
		},
		Result: "\033[34;1m@@ENCODING:\033[0m \033[32mSJIS\033[0m",
	},
	{
		Name: "Show NoHeader",
		Expr: parser.ShowFlag{
			Flag: parser.Flag{Name: "no_header"},
		},
		SetExprs: []parser.SetFlag{
			{
				Flag:  parser.Flag{Name: "no_header"},
				Value: parser.NewTernaryValueFromString("true"),
			},
		},
		Result: "\033[34;1m@@NO_HEADER:\033[0m \033[33;1mtrue\033[0m",
	},
	{
		Name: "Show WithoutNull",
		Expr: parser.ShowFlag{
			Flag: parser.Flag{Name: "without_null"},
		},
		SetExprs: []parser.SetFlag{
			{
				Flag:  parser.Flag{Name: "without_null"},
				Value: parser.NewTernaryValueFromString("true"),
			},
		},
		Result: "\033[34;1m@@WITHOUT_NULL:\033[0m \033[33;1mtrue\033[0m",
	},
	{
		Name: "Show Format",
		Expr: parser.ShowFlag{
			Flag: parser.Flag{Name: "format"},
		},
		SetExprs: []parser.SetFlag{
			{
				Flag:  parser.Flag{Name: "format"},
				Value: parser.NewStringValue("json"),
			},
		},
		Result: "\033[34;1m@@FORMAT:\033[0m \033[32mJSON\033[0m",
	},
	{
		Name: "Show WriteEncoding",
		Expr: parser.ShowFlag{
			Flag: parser.Flag{Name: "write_encoding"},
		},
		SetExprs: []parser.SetFlag{
			{
				Flag:  parser.Flag{Name: "write_encoding"},
				Value: parser.NewStringValue("SJIS"),
			},
		},
		Result: "\033[34;1m@@WRITE_ENCODING:\033[0m \033[32mSJIS\033[0m",
	},
	{
		Name: "Show WriteEncoding Ignored",
		Expr: parser.ShowFlag{
			Flag: parser.Flag{Name: "write_encoding"},
		},
		SetExprs: []parser.SetFlag{
			{
				Flag:  parser.Flag{Name: "write_encoding"},
				Value: parser.NewStringValue("SJIS"),
			},
			{
				Flag:  parser.Flag{Name: "format"},
				Value: parser.NewStringValue("JSON"),
			},
		},
		Result: "\033[34;1m@@WRITE_ENCODING:\033[0m \033[90m(ignored) SJIS\033[0m",
	},
	{
		Name: "Show WriteDelimiter",
		Expr: parser.ShowFlag{
			Flag: parser.Flag{Name: "write_delimiter"},
		},
		SetExprs: []parser.SetFlag{
			{
				Flag:  parser.Flag{Name: "write_delimiter"},
				Value: parser.NewStringValue("\t"),
			},
			{
				Flag:  parser.Flag{Name: "format"},
				Value: parser.NewStringValue("CSV"),
			},
		},
		Result: "\033[34;1m@@WRITE_DELIMITER:\033[0m \033[32m'\\t'\033[0m",
	},
	{
		Name: "Show WriteDelimiter Ignored",
		Expr: parser.ShowFlag{
			Flag: parser.Flag{Name: "write_delimiter"},
		},
		SetExprs: []parser.SetFlag{
			{
				Flag:  parser.Flag{Name: "write_delimiter"},
				Value: parser.NewStringValue("\t"),
			},
			{
				Flag:  parser.Flag{Name: "format"},
				Value: parser.NewStringValue("JSON"),
			},
		},
		Result: "\033[34;1m@@WRITE_DELIMITER:\033[0m \033[90m(ignored) '\\t'\033[0m",
	},
	{
		Name: "Show WriteDelimiterPositions for Single-Line FIXED",
		Expr: parser.ShowFlag{
			Flag: parser.Flag{Name: "write_delimiter_positions"},
		},
		SetExprs: []parser.SetFlag{
			{
				Flag:  parser.Flag{Name: "write_delimiter_positions"},
				Value: parser.NewStringValue("s[2, 5, 10]"),
			},
			{
				Flag:  parser.Flag{Name: "format"},
				Value: parser.NewStringValue("FIXED"),
			},
		},
		Result: "\033[34;1m@@WRITE_DELIMITER_POSITIONS:\033[0m \033[32mS[2, 5, 10]\033[0m",
	},
	{
		Name: "Show WriteDelimiterPositions for FIXED",
		Expr: parser.ShowFlag{
			Flag: parser.Flag{Name: "write_delimiter_positions"},
		},
		SetExprs: []parser.SetFlag{
			{
				Flag:  parser.Flag{Name: "write_delimiter_positions"},
				Value: parser.NewStringValue("spaces"),
			},
			{
				Flag:  parser.Flag{Name: "format"},
				Value: parser.NewStringValue("FIXED"),
			},
		},
		Result: "\033[34;1m@@WRITE_DELIMITER_POSITIONS:\033[0m \033[32mSPACES\033[0m",
	},
	{
		Name: "Show WriteDelimiterPositions Ignored",
		Expr: parser.ShowFlag{
			Flag: parser.Flag{Name: "write_delimiter_positions"},
		},
		SetExprs: []parser.SetFlag{
			{
				Flag:  parser.Flag{Name: "write_delimiter_positions"},
				Value: parser.NewStringValue("spaces"),
			},
			{
				Flag:  parser.Flag{Name: "format"},
				Value: parser.NewStringValue("JSON"),
			},
		},
		Result: "\033[34;1m@@WRITE_DELIMITER_POSITIONS:\033[0m \033[90m(ignored) SPACES\033[0m",
	},
	{
		Name: "Show WithoutHeader",
		Expr: parser.ShowFlag{
			Flag: parser.Flag{Name: "without_header"},
		},
		SetExprs: []parser.SetFlag{
			{
				Flag:  parser.Flag{Name: "without_header"},
				Value: parser.NewTernaryValueFromString("true"),
			},
			{
				Flag:  parser.Flag{Name: "format"},
				Value: parser.NewStringValue("CSV"),
			},
		},
		Result: "\033[34;1m@@WITHOUT_HEADER:\033[0m \033[33;1mtrue\033[0m",
	},
	{
		Name: "Show WithoutHeader Ignored",
		Expr: parser.ShowFlag{
			Flag: parser.Flag{Name: "without_header"},
		},
		SetExprs: []parser.SetFlag{
			{
				Flag:  parser.Flag{Name: "without_header"},
				Value: parser.NewTernaryValueFromString("true"),
			},
			{
				Flag:  parser.Flag{Name: "format"},
				Value: parser.NewStringValue("JSON"),
			},
		},
		Result: "\033[34;1m@@WITHOUT_HEADER:\033[0m \033[90m(ignored) true\033[0m",
	},
	{
		Name: "Show WithoutHeader with Single-Line Fixed-Length",
		Expr: parser.ShowFlag{
			Flag: parser.Flag{Name: "without_header"},
		},
		SetExprs: []parser.SetFlag{
			{
				Flag:  parser.Flag{Name: "format"},
				Value: parser.NewStringValue("fixed"),
			},
			{
				Flag:  parser.Flag{Name: "write_delimiter_positions"},
				Value: parser.NewStringValue("s[2, 5, 10]"),
			},
			{
				Flag:  parser.Flag{Name: "without_header"},
				Value: parser.NewTernaryValueFromString("true"),
			},
		},
		Result: "\033[34;1m@@WITHOUT_HEADER:\033[0m \033[90m(ignored) true\033[0m",
	},
	{
		Name: "Show lineBreak",
		Expr: parser.ShowFlag{
			Flag: parser.Flag{Name: "line_break"},
		},
		SetExprs: []parser.SetFlag{
			{
				Flag:  parser.Flag{Name: "line_break"},
				Value: parser.NewStringValue("CRLF"),
			},
		},
		Result: "\033[34;1m@@LINE_BREAK:\033[0m \033[32mCRLF\033[0m",
	},
	{
		Name: "Show lineBreak with Single-Line Fixed-Length",
		Expr: parser.ShowFlag{
			Flag: parser.Flag{Name: "line_break"},
		},
		SetExprs: []parser.SetFlag{
			{
				Flag:  parser.Flag{Name: "format"},
				Value: parser.NewStringValue("fixed"),
			},
			{
				Flag:  parser.Flag{Name: "write_delimiter_positions"},
				Value: parser.NewStringValue("s[2, 5, 10]"),
			},
			{
				Flag:  parser.Flag{Name: "line_break"},
				Value: parser.NewStringValue("CRLF"),
			},
		},
		Result: "\033[34;1m@@LINE_BREAK:\033[0m \033[90m(ignored) CRLF\033[0m",
	},
	{
		Name: "Show EncloseAll",
		Expr: parser.ShowFlag{
			Flag: parser.Flag{Name: "enclose_all"},
		},
		SetExprs: []parser.SetFlag{
			{
				Flag:  parser.Flag{Name: "enclose_all"},
				Value: parser.NewTernaryValueFromString("true"),
			},
			{
				Flag:  parser.Flag{Name: "format"},
				Value: parser.NewStringValue("CSV"),
			},
		},
		Result: "\033[34;1m@@ENCLOSE_ALL:\033[0m \033[33;1mtrue\033[0m",
	},
	{
		Name: "Show EncloseAll Ignored",
		Expr: parser.ShowFlag{
			Flag: parser.Flag{Name: "enclose_all"},
		},
		SetExprs: []parser.SetFlag{
			{
				Flag:  parser.Flag{Name: "enclose_all"},
				Value: parser.NewTernaryValueFromString("true"),
			},
			{
				Flag:  parser.Flag{Name: "format"},
				Value: parser.NewStringValue("JSON"),
			},
		},
		Result: "\033[34;1m@@ENCLOSE_ALL:\033[0m \033[90m(ignored) true\033[0m",
	},
	{
		Name: "Show JsonEscape",
		Expr: parser.ShowFlag{
			Flag: parser.Flag{Name: "json_escape"},
		},
		SetExprs: []parser.SetFlag{
			{
				Flag:  parser.Flag{Name: "json_escape"},
				Value: parser.NewStringValue("HEXALL"),
			},
			{
				Flag:  parser.Flag{Name: "format"},
				Value: parser.NewStringValue("JSON"),
			},
		},
		Result: "\033[34;1m@@JSON_ESCAPE:\033[0m \033[32mHEXALL\033[0m",
	},
	{
		Name: "Show JsonEscape Ignored",
		Expr: parser.ShowFlag{
			Flag: parser.Flag{Name: "json_escape"},
		},
		SetExprs: []parser.SetFlag{
			{
				Flag:  parser.Flag{Name: "json_escape"},
				Value: parser.NewStringValue("BACKSLASH"),
			},
			{
				Flag:  parser.Flag{Name: "format"},
				Value: parser.NewStringValue("CSV"),
			},
		},
		Result: "\033[34;1m@@JSON_ESCAPE:\033[0m \033[90m(ignored) BACKSLASH\033[0m",
	},
	{
		Name: "Show PrettyPrint",
		Expr: parser.ShowFlag{
			Flag: parser.Flag{Name: "pretty_print"},
		},
		SetExprs: []parser.SetFlag{
			{
				Flag:  parser.Flag{Name: "pretty_print"},
				Value: parser.NewTernaryValueFromString("true"),
			},
			{
				Flag:  parser.Flag{Name: "format"},
				Value: parser.NewStringValue("JSON"),
			},
		},
		Result: "\033[34;1m@@PRETTY_PRINT:\033[0m \033[33;1mtrue\033[0m",
	},
	{
		Name: "Show StripEndingLineBreak",
		Expr: parser.ShowFlag{
			Flag: parser.Flag{Name: "strip_ending_line_break"},
		},
		SetExprs: []parser.SetFlag{
			{
				Flag:  parser.Flag{Name: "strip_ending_line_break"},
				Value: parser.NewTernaryValueFromString("true"),
			},
		},
		Result: "\033[34;1m@@STRIP_ENDING_LINE_BREAK:\033[0m \033[33;1mtrue\033[0m",
	},
	{
		Name: "Show PrettyPrint Ignored",
		Expr: parser.ShowFlag{
			Flag: parser.Flag{Name: "pretty_print"},
		},
		SetExprs: []parser.SetFlag{
			{
				Flag:  parser.Flag{Name: "pretty_print"},
				Value: parser.NewTernaryValueFromString("true"),
			},
		},
		Result: "\033[34;1m@@PRETTY_PRINT:\033[0m \033[90m(ignored) true\033[0m",
	},
	{
		Name: "Show EastAsianEncoding",
		Expr: parser.ShowFlag{
			Flag: parser.Flag{Name: "east_asian_encoding"},
		},
		SetExprs: []parser.SetFlag{
			{
				Flag:  parser.Flag{Name: "east_asian_encoding"},
				Value: parser.NewTernaryValueFromString("true"),
			},
			{
				Flag:  parser.Flag{Name: "format"},
				Value: parser.NewStringValue("TEXT"),
			},
		},
		Result: "\033[34;1m@@EAST_ASIAN_ENCODING:\033[0m \033[33;1mtrue\033[0m",
	},
	{
		Name: "Show EastAsianEncoding Ignored",
		Expr: parser.ShowFlag{
			Flag: parser.Flag{Name: "east_asian_encoding"},
		},
		SetExprs: []parser.SetFlag{
			{
				Flag:  parser.Flag{Name: "east_asian_encoding"},
				Value: parser.NewTernaryValueFromString("true"),
			},
			{
				Flag:  parser.Flag{Name: "format"},
				Value: parser.NewStringValue("JSON"),
			},
		},
		Result: "\033[34;1m@@EAST_ASIAN_ENCODING:\033[0m \033[90m(ignored) true\033[0m",
	},
	{
		Name: "Show CountDiacriticalSign",
		Expr: parser.ShowFlag{
			Flag: parser.Flag{Name: "count_diacritical_sign"},
		},
		SetExprs: []parser.SetFlag{
			{
				Flag:  parser.Flag{Name: "count_diacritical_sign"},
				Value: parser.NewTernaryValueFromString("true"),
			},
			{
				Flag:  parser.Flag{Name: "format"},
				Value: parser.NewStringValue("TEXT"),
			},
		},
		Result: "\033[34;1m@@COUNT_DIACRITICAL_SIGN:\033[0m \033[33;1mtrue\033[0m",
	},
	{
		Name: "Show CountDiacriticalSign Ignored",
		Expr: parser.ShowFlag{
			Flag: parser.Flag{Name: "count_diacritical_sign"},
		},
		SetExprs: []parser.SetFlag{
			{
				Flag:  parser.Flag{Name: "count_diacritical_sign"},
				Value: parser.NewTernaryValueFromString("true"),
			},
			{
				Flag:  parser.Flag{Name: "format"},
				Value: parser.NewStringValue("JSON"),
			},
		},
		Result: "\033[34;1m@@COUNT_DIACRITICAL_SIGN:\033[0m \033[90m(ignored) true\033[0m",
	},
	{
		Name: "Show CountFormatCode",
		Expr: parser.ShowFlag{
			Flag: parser.Flag{Name: "count_format_code"},
		},
		SetExprs: []parser.SetFlag{
			{
				Flag:  parser.Flag{Name: "count_format_code"},
				Value: parser.NewTernaryValueFromString("true"),
			},
			{
				Flag:  parser.Flag{Name: "format"},
				Value: parser.NewStringValue("TEXT"),
			},
		},
		Result: "\033[34;1m@@COUNT_FORMAT_CODE:\033[0m \033[33;1mtrue\033[0m",
	},
	{
		Name: "Show CountFormatCode Ignored",
		Expr: parser.ShowFlag{
			Flag: parser.Flag{Name: "count_format_code"},
		},
		SetExprs: []parser.SetFlag{
			{
				Flag:  parser.Flag{Name: "count_format_code"},
				Value: parser.NewTernaryValueFromString("true"),
			},
			{
				Flag:  parser.Flag{Name: "format"},
				Value: parser.NewStringValue("JSON"),
			},
		},
		Result: "\033[34;1m@@COUNT_FORMAT_CODE:\033[0m \033[90m(ignored) true\033[0m",
	},
	{
		Name: "Show Color",
		Expr: parser.ShowFlag{
			Flag: parser.Flag{Name: "color"},
		},
		SetExprs: []parser.SetFlag{
			{
				Flag:  parser.Flag{Name: "color"},
				Value: parser.NewTernaryValueFromString("true"),
			},
		},
		Result: "\033[34;1m@@COLOR:\033[0m \033[33;1mtrue\033[0m",
	},
	{
		Name: "Show Quiet",
		Expr: parser.ShowFlag{
			Flag: parser.Flag{Name: "quiet"},
		},
		SetExprs: []parser.SetFlag{
			{
				Flag:  parser.Flag{Name: "quiet"},
				Value: parser.NewTernaryValueFromString("true"),
			},
		},
		Result: "\033[34;1m@@QUIET:\033[0m \033[33;1mtrue\033[0m",
	},
	{
		Name: "Show LimitRecursion",
		Expr: parser.ShowFlag{
			Flag: parser.Flag{Name: "limit_recursion"},
		},
		SetExprs: []parser.SetFlag{
			{
				Flag:  parser.Flag{Name: "limit_recursion"},
				Value: parser.NewIntegerValue(3),
			},
		},
		Result: "\033[34;1m@@LIMIT_RECURSION:\033[0m \033[35m3\033[0m",
	},
	{
		Name: "Show LimitRecursion No Limit",
		Expr: parser.ShowFlag{
			Flag: parser.Flag{Name: "limit_recursion"},
		},
		SetExprs: []parser.SetFlag{
			{
				Flag:  parser.Flag{Name: "limit_recursion"},
				Value: parser.NewIntegerValue(-100),
			},
		},
		Result: "\033[34;1m@@LIMIT_RECURSION:\033[0m \033[90m(no limit)\033[0m",
	},
	{
		Name: "Show CPU",
		Expr: parser.ShowFlag{
			Flag: parser.Flag{Name: "cpu"},
		},
		SetExprs: []parser.SetFlag{
			{
				Flag:  parser.Flag{Name: "cpu"},
				Value: parser.NewIntegerValue(1),
			},
		},
		Result: "\033[34;1m@@CPU:\033[0m \033[35m1\033[0m",
	},
	{
		Name: "Show Stats",
		Expr: parser.ShowFlag{
			Flag: parser.Flag{Name: "stats"},
		},
		SetExprs: []parser.SetFlag{
			{
				Flag:  parser.Flag{Name: "stats"},
				Value: parser.NewTernaryValueFromString("true"),
			},
		},
		Result: "\033[34;1m@@STATS:\033[0m \033[33;1mtrue\033[0m",
	},
	{
		Name: "Invalid Flag Name Error",
		Expr: parser.ShowFlag{
			Flag: parser.Flag{Name: "invalid"},
		},
		Error: "@@INVALID is an unknown flag",
	},
}

func TestShowFlag(t *testing.T) {
	defer func() {
		TestTx.UseColor(false)
		initFlag(TestTx.Flags)
	}()

	TestTx.UseColor(true)
	ctx := context.Background()
	scope := NewReferenceScope(TestTx)

	for _, v := range showFlagTests {
		initFlag(TestTx.Flags)
		TestTx.UseColor(true)
		for _, expr := range v.SetExprs {
			_ = SetFlag(ctx, scope, expr)
		}
		result, err := ShowFlag(TestTx, v.Expr)
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

var showObjectsTests = []struct {
	Name                    string
	Expr                    parser.ShowObjects
	Scope                   *ReferenceScope
	PreparedStatements      PreparedStatementMap
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
	UncommittedViews        UncommittedViews
	Expect                  string
	Error                   string
}{
	{
		Name: "ShowObjects Tables",
		Expr: parser.ShowObjects{Type: parser.Identifier{Literal: "tables"}},
		ViewCache: GenerateViewMap([]*View{
			{
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
			{
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
			{
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
			{
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
			{
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
			{
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
		}),
		Expect: "\n" +
			"                      Loaded Tables\n" +
			"----------------------------------------------------------\n" +
			" table1.csv\n" +
			"     Fields: col1, col2\n" +
			"     Format: CSV     Delimiter: '\\t'  Enclose All: false\n" +
			"     Encoding: SJIS  LineBreak: CRLF  Header: false\n" +
			" table1.json\n" +
			"     Fields: col1, col2\n" +
			"     Format: JSON    Escape: BACKSLASH  Query: {}\n" +
			"     Encoding: UTF8  LineBreak: LF    Pretty Print: false\n" +
			" table1.tsv\n" +
			"     Fields: col1, col2\n" +
			"     Format: TSV     Delimiter: '\\t'  Enclose All: false\n" +
			"     Encoding: UTF8  LineBreak: LF    Header: true\n" +
			" table1.txt\n" +
			"     Fields: col1, col2\n" +
			"     Format: FIXED   Delimiter Positions: [3, 12]\n" +
			"     Encoding: UTF8  LineBreak: LF    Header: true\n" +
			" table2.json\n" +
			"     Fields: col1, col2\n" +
			"     Format: JSON    Escape: HEX      Query: (empty)\n" +
			"     Encoding: UTF8  LineBreak: LF    Pretty Print: false\n" +
			" table2.txt\n" +
			"     Fields: col1, col2\n" +
			"     Format: FIXED   Delimiter Positions: S[3, 12]\n" +
			"     Encoding: UTF8\n" +
			"\n",
	},
	{
		Name: "ShowObjects Tables Uncommitted",
		Expr: parser.ShowObjects{Type: parser.Identifier{Literal: "tables"}},
		ViewCache: GenerateViewMap([]*View{
			{
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
			{
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
			{
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
			{
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
			{
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
			{
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
		}),
		UncommittedViews: UncommittedViews{
			mtx: &sync.RWMutex{},
			Created: map[string]*FileInfo{
				"TABLE1.TSV": {Path: "table1.tsv"},
			},
			Updated: map[string]*FileInfo{
				"TABLE2.JSON": {Path: "table2.json"},
			},
		},
		Expect: "\n" +
			"          Loaded Tables (Uncommitted: 2 Tables)\n" +
			"----------------------------------------------------------\n" +
			" table1.csv\n" +
			"     Fields: col1, col2\n" +
			"     Format: CSV     Delimiter: '\\t'  Enclose All: false\n" +
			"     Encoding: SJIS  LineBreak: CRLF  Header: false\n" +
			" table1.json\n" +
			"     Fields: col1, col2\n" +
			"     Format: JSON    Escape: BACKSLASH  Query: {}\n" +
			"     Encoding: UTF8  LineBreak: LF    Pretty Print: false\n" +
			" *Created* table1.tsv\n" +
			"     Fields: col1, col2\n" +
			"     Format: TSV     Delimiter: '\\t'  Enclose All: false\n" +
			"     Encoding: UTF8  LineBreak: LF    Header: true\n" +
			" table1.txt\n" +
			"     Fields: col1, col2\n" +
			"     Format: FIXED   Delimiter Positions: [3, 12]\n" +
			"     Encoding: UTF8  LineBreak: LF    Header: true\n" +
			" *Updated* table2.json\n" +
			"     Fields: col1, col2\n" +
			"     Format: JSON    Escape: BACKSLASH  Query: (empty)\n" +
			"     Encoding: UTF8  LineBreak: LF    Pretty Print: false\n" +
			" table2.txt\n" +
			"     Fields: col1, col2\n" +
			"     Format: FIXED   Delimiter Positions: S[3, 12]\n" +
			"     Encoding: UTF8\n" +
			"\n",
	},
	{
		Name: "ShowObjects Tables Long Fields",
		Expr: parser.ShowObjects{Type: parser.Identifier{Literal: "tables"}},
		ViewCache: GenerateViewMap([]*View{
			{
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
		}),
		Expect: "\n" +
			"                        Loaded Tables\n" +
			"--------------------------------------------------------------\n" +
			" table1.csv\n" +
			"     Fields: colabcdef1, colabcdef2, colabcdef3, colabcdef4, \n" +
			"             colabcdef5, colabcdef6, colabcdef7\n" +
			"     Format: CSV     Delimiter: '\\t'  Enclose All: false\n" +
			"     Encoding: SJIS  LineBreak: CRLF  Header: false\n" +
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
		Scope: GenerateReferenceScope([]map[string]map[string]interface{}{
			{
				scopeNameTempTables: {
					"VIEW1": &View{
						FileInfo: &FileInfo{
							Path:     "view1",
							ViewType: ViewTypeTemporaryTable,
						},
						Header: NewHeader("view1", []string{"column1", "column2"}),
					},
				},
			},
			{
				scopeNameTempTables: {
					"VIEW1": &View{
						FileInfo: &FileInfo{
							Path:     "view1",
							ViewType: ViewTypeTemporaryTable,
						},
						Header: NewHeader("view1", []string{"column1", "column2", "column3"}),
					},
					"VIEW2": &View{
						FileInfo: &FileInfo{
							Path:     "view2",
							ViewType: ViewTypeTemporaryTable,
						},
						Header: NewHeader("view2", []string{"column1", "column2"}),
					},
				},
			},
		}, nil, time.Time{}, nil),
		UncommittedViews: UncommittedViews{
			mtx:     &sync.RWMutex{},
			Created: map[string]*FileInfo{},
			Updated: map[string]*FileInfo{
				"VIEW2": {
					Path:     "view2",
					ViewType: ViewTypeTemporaryTable,
				},
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
		Scope: GenerateReferenceScope([]map[string]map[string]interface{}{
			{
				scopeNameCursors: {
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
					"CUR5": &Cursor{
						name:      "stmtcur",
						statement: parser.Identifier{Literal: "stmt"},
					},
				},
			},
		}, nil, time.Time{}, nil),
		Expect: "\n" +
			"                               Cursors\n" +
			"---------------------------------------------------------------------\n" +
			" cur\n" +
			"     Status: Closed\n" +
			"     Query:\n" +
			"       SELECT column1, column2 FROM table1\n" +
			" cur2\n" +
			"     Status: Open    Number of Rows: 2         Pointer: UNKNOWN\n" +
			"     Query:\n" +
			"       SELECT column1, column2 FROM table1\n" +
			" cur3\n" +
			"     Status: Open    Number of Rows: 2         Pointer: 1\n" +
			"     Query:\n" +
			"       SELECT column1, column2 FROM table1\n" +
			" cur4\n" +
			"     Status: Open    Number of Rows: 2         Pointer: Out of Range\n" +
			"     Query:\n" +
			"       SELECT column1, column2 FROM table1\n" +
			" stmtcur\n" +
			"     Status: Closed\n" +
			"     Statement: stmt\n" +
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
		Scope: GenerateReferenceScope([]map[string]map[string]interface{}{
			{
				scopeNameFunctions: {
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
			},
			{
				scopeNameFunctions: {
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
		}, nil, time.Time{}, nil),
		Expect: "\n" +
			" Scalar Functions\n" +
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
		Name: "ShowObjects Statements",
		Expr: parser.ShowObjects{Type: parser.Identifier{Literal: "statements"}},
		PreparedStatements: GenerateStatementMap([]*PreparedStatement{
			{
				Name:            "stmt1",
				StatementString: "select 1;\nselect 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 22, 22;",
				Statements: []parser.Statement{
					parser.SelectQuery{
						SelectEntity: parser.SelectEntity{
							SelectClause: parser.SelectClause{
								BaseExpr: parser.NewBaseExpr(parser.Token{Line: 1, Char: 1, SourceFile: "stmt"}),
								Fields: []parser.QueryExpression{
									parser.Field{
										Object: parser.NewIntegerValueFromString("1"),
									},
								},
							},
						},
					},
					parser.SelectQuery{
						SelectEntity: parser.SelectEntity{
							SelectClause: parser.SelectClause{
								BaseExpr: parser.NewBaseExpr(parser.Token{Line: 2, Char: 1, SourceFile: "stmt"}),
								Fields: []parser.QueryExpression{
									parser.Field{
										Object: parser.NewIntegerValueFromString("2"),
									},
								},
							},
						},
					},
				},
				HolderNumber: 0,
			},
			{
				Name:            "stmt2",
				StatementString: "select ?",
				Statements: []parser.Statement{
					parser.SelectQuery{
						SelectEntity: parser.SelectEntity{
							SelectClause: parser.SelectClause{
								BaseExpr: parser.NewBaseExpr(parser.Token{Line: 1, Char: 1, SourceFile: "stmt"}),
								Fields: []parser.QueryExpression{
									parser.Field{
										Object: parser.Placeholder{Literal: "?", Ordinal: 1},
									},
								},
							},
						},
					},
				},
				HolderNumber: 1,
			},
		}),
		Expect: "\n" +
			"                          Prepared Statements\n" +
			"-----------------------------------------------------------------------\n" +
			" stmt1\n" +
			"     Placeholder Number: 0\n" +
			"     Statement:\n" +
			"       select 1;\n" +
			"       select 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17,\\\n" +
			"       18, 19, 20, 22, 22;\n" +
			" stmt2\n" +
			"     Placeholder Number: 1\n" +
			"     Statement:\n" +
			"       select ?\n" +
			"\n",
	},
	{
		Name:   "ShowObjects Statements Empty",
		Expr:   parser.ShowObjects{Type: parser.Identifier{Literal: "statements"}},
		Expect: "No statement is prepared",
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
			"               @@ANSI_QUOTES: false\n" +
			"              @@STRICT_EQUAL: false\n" +
			"              @@WAIT_TIMEOUT: 15\n" +
			"             @@IMPORT_FORMAT: CSV\n" +
			"                 @@DELIMITER: ','\n" +
			"       @@DELIMITER_POSITIONS: SPACES\n" +
			"                @@JSON_QUERY: (empty)\n" +
			"                  @@ENCODING: AUTO\n" +
			"                 @@NO_HEADER: false\n" +
			"              @@WITHOUT_NULL: false\n" +
			"   @@STRIP_ENDING_LINE_BREAK: false\n" +
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
			"           @@LIMIT_RECURSION: 5\n" +
			"                       @@CPU: " + strconv.Itoa(TestTx.Flags.CPU) + "\n" +
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
		Error: "object type invalid is invalid",
	},
}

func TestShowObjects(t *testing.T) {
	defer func() {
		_ = TestTx.ReleaseResources()
		TestTx.uncommittedViews.Clean()
		TestTx.PreparedStatements = NewPreparedStatementMap()
		initFlag(TestTx.Flags)
	}()

	for _, v := range showObjectsTests {
		initFlag(TestTx.Flags)

		TestTx.Flags.Repository = v.Repository
		TestTx.Flags.ImportOptions.Format = v.ImportFormat
		TestTx.Flags.ImportOptions.Delimiter = ','
		if v.Delimiter != 0 {
			TestTx.Flags.ImportOptions.Delimiter = v.Delimiter
		}
		TestTx.Flags.ImportOptions.DelimiterPositions = v.DelimiterPositions
		TestTx.Flags.ImportOptions.SingleLine = v.SingleLine
		TestTx.Flags.ImportOptions.JsonQuery = v.JsonQuery
		TestTx.Flags.ExportOptions.Delimiter = ','
		if v.WriteDelimiter != 0 {
			TestTx.Flags.ExportOptions.Delimiter = v.WriteDelimiter
		}
		TestTx.Flags.ExportOptions.DelimiterPositions = v.WriteDelimiterPositions
		TestTx.Flags.ExportOptions.SingleLine = v.WriteAsSingleLine
		TestTx.Flags.ExportOptions.Format = v.Format
		_ = TestTx.cachedViews.Clean(TestTx.FileContainer)
		if v.ViewCache.SyncMap != nil {
			TestTx.cachedViews = v.ViewCache
		}
		if v.UncommittedViews.mtx == nil {
			TestTx.uncommittedViews = NewUncommittedViews()
		} else {
			TestTx.uncommittedViews = v.UncommittedViews
		}
		TestTx.PreparedStatements = NewPreparedStatementMap()
		if v.PreparedStatements.SyncMap != nil {
			TestTx.PreparedStatements = v.PreparedStatements
		}

		if v.Scope == nil {
			v.Scope = NewReferenceScope(TestTx)
		}
		result, err := ShowObjects(v.Scope, v.Expr)
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
}

var showFieldsTests = []struct {
	Name             string
	Expr             parser.ShowFields
	Scope            *ReferenceScope
	ViewCache        ViewMap
	UncommittedViews UncommittedViews
	Expect           string
	Error            string
}{
	{
		Name: "ShowFields Temporary Table",
		Expr: parser.ShowFields{
			Type:  parser.Identifier{Literal: "fields"},
			Table: parser.Identifier{Literal: "view1"},
		},
		Scope: GenerateReferenceScope([]map[string]map[string]interface{}{
			{
				scopeNameTempTables: {
					"VIEW1": &View{
						Header: NewHeader("view1", []string{"column1", "column2"}),
						FileInfo: &FileInfo{
							Path:     "view1",
							ViewType: ViewTypeTemporaryTable,
						},
					},
				},
			},
		}, nil, time.Time{}, nil),
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
		Name: "ShowFields Stdin Table",
		Expr: parser.ShowFields{
			Type:  parser.Identifier{Literal: "fields"},
			Table: parser.Stdin{},
		},
		Scope: GenerateReferenceScope([]map[string]map[string]interface{}{
			{
				scopeNameTempTables: {
					"STDIN": &View{
						Header: NewHeader("stdin", []string{"column1", "column2"}),
						FileInfo: &FileInfo{
							Path:     "stdin",
							ViewType: ViewTypeStdin,
						},
					},
				},
			},
		}, nil, time.Time{}, nil),
		Expect: "\n" +
			" Fields in STDIN\n" +
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
		Scope: GenerateReferenceScope([]map[string]map[string]interface{}{
			{
				scopeNameTempTables: {
					"VIEW1": &View{
						Header: NewHeader("view1", []string{"column1", "column2"}),
						FileInfo: &FileInfo{
							Path:     "view1",
							ViewType: ViewTypeTemporaryTable,
						},
					},
				},
			},
		}, nil, time.Time{}, nil),
		UncommittedViews: UncommittedViews{
			mtx:     &sync.RWMutex{},
			Created: map[string]*FileInfo{},
			Updated: map[string]*FileInfo{
				"VIEW1": {
					Path:     "view1",
					ViewType: ViewTypeTemporaryTable,
				},
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
				Type:          parser.Token{Token: parser.CSV, Literal: "csv"},
				FormatElement: parser.NewStringValue(","),
				Path:          parser.Identifier{Literal: "show_fields_create.csv"},
			},
		},
		ViewCache: GenerateViewMap([]*View{
			{
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
		}),
		UncommittedViews: UncommittedViews{
			mtx: &sync.RWMutex{},
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
			" Format: CSV     Delimiter: ','   Enclose All: false\n" +
			" Encoding: UTF8  LineBreak: LF    Header: true\n" +
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
		ViewCache: GenerateViewMap([]*View{
			{
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
		}),
		UncommittedViews: UncommittedViews{
			mtx: &sync.RWMutex{},
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
			" Format: CSV     Delimiter: ','   Enclose All: false\n" +
			" Encoding: UTF8  LineBreak: LF    Header: true\n" +
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
		ViewCache: GenerateViewMap([]*View{
			{
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
		}),
		UncommittedViews: UncommittedViews{
			mtx:     &sync.RWMutex{},
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
			" Format: CSV     Delimiter: ','   Enclose All: false\n" +
			" Encoding: UTF8  LineBreak: LF    Header: true\n" +
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
		Error: "file notexist does not exist",
	},
	{
		Name: "ShowFields Invalid Object Type",
		Expr: parser.ShowFields{
			Type:  parser.Identifier{Literal: "invalid"},
			Table: parser.Identifier{Literal: "table2"},
		},
		Error: "object type invalid is invalid",
	},
}

func calcShowFieldsWidth(fileName string, fileNameInTitle string, prefixLen int) int {
	w := 53
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
	defer func() {
		_ = TestTx.ReleaseResources()
		TestTx.uncommittedViews.Clean()
		initFlag(TestTx.Flags)
	}()

	initFlag(TestTx.Flags)
	TestTx.Flags.Repository = TestDir
	ctx := context.Background()

	for _, v := range showFieldsTests {
		_ = TestTx.cachedViews.Clean(TestTx.FileContainer)
		if v.ViewCache.SyncMap != nil {
			TestTx.cachedViews = v.ViewCache
		}
		if v.UncommittedViews.mtx == nil {
			TestTx.uncommittedViews = NewUncommittedViews()
		} else {
			TestTx.uncommittedViews = v.UncommittedViews
		}

		if v.Scope == nil {
			v.Scope = NewReferenceScope(TestTx)
		}

		result, err := ShowFields(ctx, v.Scope, v.Expr)
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
		Error: "field err does not exist",
	},
}

func TestSetEnvVar(t *testing.T) {
	ctx := context.Background()
	scope := NewReferenceScope(TestTx)

	for _, v := range setEnvVarTests {
		err := SetEnvVar(ctx, scope, v.Expr)

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
			"-----------------------------------------------------\n" +
			" SELECT Clause\n" +
			"     select_clause\n" +
			"         : SELECT [DISTINCT] <field> [, <field> ...]\n" +
			"\n" +
			"     field\n" +
			"         : <value>\n" +
			"         | <value> AS alias\n" +
			"\n" +
			"\n",
	},
	{
		Expr: parser.Syntax{Keywords: []parser.QueryExpression{parser.NewStringValue(" select  "), parser.NewStringValue("clause")}},
		Expect: "\n" +
			"                Search: select clause\n" +
			"-----------------------------------------------------\n" +
			" SELECT Clause\n" +
			"     select_clause\n" +
			"         : SELECT [DISTINCT] <field> [, <field> ...]\n" +
			"\n" +
			"     field\n" +
			"         : <value>\n" +
			"         | <value> AS alias\n" +
			"\n" +
			"\n",
	},
	{
		Expr: parser.Syntax{Keywords: []parser.QueryExpression{parser.FieldReference{Column: parser.Identifier{Literal: "select clause"}}}},
		Expect: "\n" +
			"                Search: select clause\n" +
			"-----------------------------------------------------\n" +
			" SELECT Clause\n" +
			"     select_clause\n" +
			"         : SELECT [DISTINCT] <field> [, <field> ...]\n" +
			"\n" +
			"     field\n" +
			"         : <value>\n" +
			"         | <value> AS alias\n" +
			"\n" +
			"\n",
	},
	{
		Expr: parser.Syntax{Keywords: []parser.QueryExpression{parser.NewStringValue("operator prec")}},
		Expect: "\n" +
			"        Search: operator prec\n" +
			"--------------------------------------\n" +
			" Operator Precedence\n" +
			"     Operator Precedence Description.\n" +
			"\n" +
			"\n",
	},
	{
		Expr: parser.Syntax{Keywords: []parser.QueryExpression{parser.NewStringValue("string  op")}},
		Expect: "\n" +
			"      Search: string op\n" +
			"------------------------------\n" +
			" String Operators\n" +
			"     concatenation\n" +
			"         : <value> || <value>\n" +
			"\n" +
			"         description\n" +
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

	ctx := context.Background()
	scope := NewReferenceScope(TestTx)

	for _, v := range syntaxTests {
		result, _ := Syntax(ctx, scope, v.Expr)
		if result != v.Expect {
			t.Errorf("result = %s, want %s for %v", result, v.Expect, v.Expr)
		}
	}

	syntax.CsvqSyntax = origSyntax
}
