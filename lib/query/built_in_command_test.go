package query

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/mithrandie/go-text"

	"github.com/mithrandie/go-text/fixedlen"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/value"
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
		Name: "Source File Not Readable Error",
		Expr: parser.Source{
			FilePath: parser.NewStringValue(TestDir),
		},
		Error: fmt.Sprintf("[L:- C:-] file %s is unable to read", TestDir),
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
		Error: "[L:- C:-] flag @@invalid does not exist",
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
		initFlag()
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
	initFlag()
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
				Value: parser.NewStringValue("%Y%m%d"),
			},
		},
		Result: "\033[34;1m@@DATETIME_FORMAT:\033[0m \033[32m%Y%m%d\033[0m",
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
		Name: "Show Delimiter for CSV",
		Expr: parser.ShowFlag{
			Name: "delimiter",
		},
		SetExprs: []parser.SetFlag{
			{
				Name:  "delimiter",
				Value: parser.NewStringValue("\t"),
			},
		},
		Result: "\033[34;1m@@DELIMITER:\033[0m \033[32m'\\t'\033[0m\033[34;1m | \033[0m\033[90mSPACES\033[0m",
	},
	{
		Name: "Show Delimiter for FIXED",
		Expr: parser.ShowFlag{
			Name: "delimiter",
		},
		SetExprs: []parser.SetFlag{
			{
				Name:  "delimiter",
				Value: parser.NewStringValue("SPACES"),
			},
		},
		Result: "\033[34;1m@@DELIMITER:\033[0m \033[90m','\033[0m\033[34;1m | \033[0m\033[32mSPACES\033[0m",
	},
	{
		Name: "Show Delimiter Ignored",
		Expr: parser.ShowFlag{
			Name: "delimiter",
		},
		SetExprs: []parser.SetFlag{
			{
				Name:  "json_query",
				Value: parser.NewStringValue("{}"),
			},
		},
		Result: "\033[34;1m@@DELIMITER:\033[0m \033[90m(ignored) ',' | SPACES\033[0m",
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
		Name: "Show JsonQuery Ignored",
		Expr: parser.ShowFlag{
			Name: "json_query",
		},
		SetExprs: []parser.SetFlag{},
		Result:   "\033[34;1m@@JSON_QUERY:\033[0m \033[90m(ignored) (empty)\033[0m",
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
		Name: "Show WriteDelimiter for CSV",
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
		Result: "\033[34;1m@@WRITE_DELIMITER:\033[0m \033[32m'\\t'\033[0m\033[34;1m | \033[0m\033[90mSPACES\033[0m",
	},
	{
		Name: "Show WriteDelimiter for FIXED",
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
				Value: parser.NewStringValue("FIXED"),
			},
		},
		Result: "\033[34;1m@@WRITE_DELIMITER:\033[0m \033[90m'\\t'\033[0m\033[34;1m | \033[0m\033[32mSPACES\033[0m",
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
		Result: "\033[34;1m@@WRITE_DELIMITER:\033[0m \033[90m(ignored) '\\t' | SPACES\033[0m",
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
		Error: "[L:- C:-] flag @@invalid does not exist",
	},
}

func TestShowFlag(t *testing.T) {
	filter := NewEmptyFilter()

	for _, v := range showFlagTests {
		initFlag()
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
	initFlag()
}

var showObjectsTests = []struct {
	Name                    string
	Expr                    parser.ShowObjects
	Filter                  *Filter
	Delimiter               rune
	DelimiterPositions      fixedlen.DelimiterPositions
	DelimitAutomatically    bool
	JsonQuery               string
	Repository              string
	Format                  cmd.Format
	WriteDelimiter          rune
	WriteDelimiterPositions fixedlen.DelimiterPositions
	ViewCache               ViewMap
	ExecResults             []ExecResult
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
		},
		Expect: "\n" +
			"                      Loaded Tables\n" +
			"----------------------------------------------------------\n" +
			" table1.csv\n" +
			"     Fields: col1, col2\n" +
			"     Format: CSV     Delimiter: '\\t'  Enclose All: false\n" +
			"     Encoding: SJIS  LineBreak: CRLF  Header: false\n" +
			" table1.json\n" +
			"     Fields: col1, col2\n" +
			"     Format: JSON    Query: {}\n" +
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
			"     Format: JSON    Query: (empty)\n" +
			"     Encoding: UTF8  LineBreak: LF    Pretty Print: false\n" +
			"",
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
		},
		ExecResults: []ExecResult{
			{
				Type: CreateTableQuery,
				FileInfo: &FileInfo{
					Path: "table1.tsv",
				},
			},
			{
				Type: UpdateQuery,
				FileInfo: &FileInfo{
					Path: "table2.json",
				},
				OperatedCount: 2,
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
			"     Format: JSON    Query: {}\n" +
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
			"     Format: JSON    Query: (empty)\n" +
			"     Encoding: UTF8  LineBreak: LF    Pretty Print: false\n" +
			"",
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
			"     Format: CSV     Delimiter: '\\t'  Enclose All: false\n" +
			"     Encoding: SJIS  LineBreak: CRLF  Header: false\n" +
			"",
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
		ExecResults: []ExecResult{
			{
				Type: UpdateQuery,
				FileInfo: &FileInfo{
					Path:        "view2",
					IsTemporary: true,
				},
				OperatedCount: 2,
			},
		},
		Expect: "\n" +
			" Views (Uncommitted: 1 View)\n" +
			"------------------------------\n" +
			" view1\n" +
			"     Fields: column1, column2\n" +
			" *Updated* view2\n" +
			"     Fields: column1, column2\n",
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
			"",
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
			" useraggfunc (column1, @arg1, @arg2 = 1)\n",
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
			"                    Flags\n" +
			"---------------------------------------------\n" +
			"             @@REPOSITORY: .\n" +
			"               @@TIMEZONE: UTC\n" +
			"        @@DATETIME_FORMAT: (not set)\n" +
			"           @@WAIT_TIMEOUT: 15\n" +
			"              @@DELIMITER: ',' | SPACES\n" +
			"             @@JSON_QUERY: (ignored) (empty)\n" +
			"               @@ENCODING: UTF8\n" +
			"              @@NO_HEADER: false\n" +
			"           @@WITHOUT_NULL: false\n" +
			"                 @@FORMAT: CSV\n" +
			"         @@WRITE_ENCODING: UTF8\n" +
			"        @@WRITE_DELIMITER: ',' | SPACES\n" +
			"         @@WITHOUT_HEADER: false\n" +
			"             @@LINE_BREAK: LF\n" +
			"            @@ENCLOSE_ALL: false\n" +
			"           @@PRETTY_PRINT: (ignored) false\n" +
			"    @@EAST_ASIAN_ENCODING: (ignored) false\n" +
			" @@COUNT_DIACRITICAL_SIGN: (ignored) false\n" +
			"      @@COUNT_FORMAT_CODE: (ignored) false\n" +
			"                  @@COLOR: false\n" +
			"                  @@QUIET: false\n" +
			"                    @@CPU: " + strconv.Itoa(cmd.GetFlags().CPU) + "\n" +
			"                  @@STATS: false\n" +
			"",
	},
	{
		Name:  "ShowObjects Invalid Object Type",
		Expr:  parser.ShowObjects{Type: parser.Identifier{Literal: "invalid"}},
		Error: "[L:- C:-] object type invalid is invalid",
	},
}

func TestShowObjects(t *testing.T) {
	initFlag()
	flags := cmd.GetFlags()

	for _, v := range showObjectsTests {
		flags.Repository = v.Repository
		flags.Delimiter = ','
		if v.Delimiter != 0 {
			flags.Delimiter = v.Delimiter
		}
		flags.DelimiterPositions = v.DelimiterPositions
		flags.DelimitAutomatically = v.DelimitAutomatically
		flags.JsonQuery = v.JsonQuery
		flags.WriteDelimiterPositions = v.WriteDelimiterPositions
		flags.Format = v.Format
		ViewCache.Clean()
		ExecResults = make([]ExecResult, 0)
		if 0 < len(v.ViewCache) {
			ViewCache = v.ViewCache
		}
		if 0 < len(v.ExecResults) {
			ExecResults = v.ExecResults
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
}

var showFieldsTests = []struct {
	Name        string
	Expr        parser.ShowFields
	Filter      *Filter
	ViewCache   ViewMap
	ExecResults []ExecResult
	Expect      string
	Error       string
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
			"   2. column2\n",
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
		ExecResults: []ExecResult{
			{
				Type: UpdateQuery,
				FileInfo: &FileInfo{
					Path:        "view1",
					IsTemporary: true,
				},
				OperatedCount: 2,
			},
		},
		Expect: "\n" +
			" Fields in view1\n" +
			"-----------------\n" +
			" Type: View\n" +
			" Status: Updated\n" +
			" Fields:\n" +
			"   1. column1\n" +
			"   2. column2\n",
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
		ExecResults: []ExecResult{
			{
				Type: CreateTableQuery,
				FileInfo: &FileInfo{
					Path: GetTestFilePath("show_fields_create.csv"),
				},
			},
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
			"   2. column2\n",
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
		ExecResults: []ExecResult{
			{
				Type: UpdateQuery,
				FileInfo: &FileInfo{
					Path: GetTestFilePath("show_fields_update.csv"),
				},
				OperatedCount: 2,
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
			"   2. column2\n",
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

func TestShowFields(t *testing.T) {
	initFlag()
	flags := cmd.GetFlags()
	flags.Repository = TestDir

	for _, v := range showFieldsTests {
		ViewCache.Clean()
		ExecResults = make([]ExecResult, 0)
		if 0 < len(v.ViewCache) {
			ViewCache = v.ViewCache
		}
		if 0 < len(v.ExecResults) {
			ExecResults = v.ExecResults
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
}
