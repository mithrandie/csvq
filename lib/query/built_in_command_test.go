package query

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"testing"

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
		Result: "'foo'",
	},
	{
		Name: "Print Error",
		Expr: parser.Print{
			Value: parser.Variable{
				Name: "var",
			},
		},
		Error: "[L:- C:-] variable var is undeclared",
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
		Error: "[L:- C:-] variable var is undeclared",
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
		Error: "[L:- C:-] variable var is undeclared",
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

var setFlagTests = []struct {
	Name               string
	Expr               parser.SetFlag
	ResultFlag         string
	ResultStrValue     string
	ResultFloatValue   float64
	ResultIntegerValue int
	ResultBoolValue    bool
	Error              string
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
		Name: "Set lineBreak",
		Expr: parser.SetFlag{
			Name:  "@@line_break",
			Value: value.NewString("CRLF"),
		},
		ResultFlag:     "line_break",
		ResultStrValue: "\r\n",
	},
	{
		Name: "Set Timezone",
		Expr: parser.SetFlag{
			Name:  "@@timezone",
			Value: value.NewString("utc"),
		},
		ResultFlag:     "timezone",
		ResultStrValue: "UTC",
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
		Name: "Set WriteEncoding",
		Expr: parser.SetFlag{
			Name:  "@@write_encoding",
			Value: value.NewString("SJIS"),
		},
		ResultFlag:     "write_encoding",
		ResultStrValue: "SJIS",
	},
	{
		Name: "Set Format",
		Expr: parser.SetFlag{
			Name:  "@@format",
			Value: value.NewString("json"),
		},
		ResultFlag:     "format",
		ResultStrValue: "JSON",
	},
	{
		Name: "Set WriteDelimiter",
		Expr: parser.SetFlag{
			Name:  "@@write_delimiter",
			Value: value.NewString("\t"),
		},
		ResultFlag:     "write_delimiter",
		ResultStrValue: "\t",
	},
	{
		Name: "Set WithoutHeader",
		Expr: parser.SetFlag{
			Name:  "@@without_header",
			Value: value.NewBoolean(true),
		},
		ResultFlag:      "without_header",
		ResultBoolValue: true,
	},
	{
		Name: "Set PrettyPrint",
		Expr: parser.SetFlag{
			Name:  "@@pretty_print",
			Value: value.NewBoolean(true),
		},
		ResultFlag:      "pretty_print",
		ResultBoolValue: true,
	},
	{
		Name: "Set Color",
		Expr: parser.SetFlag{
			Name:  "@@color",
			Value: value.NewBoolean(true),
		},
		ResultFlag:      "color",
		ResultBoolValue: true,
	},
	{
		Name: "Set Quiet",
		Expr: parser.SetFlag{
			Name:  "@@quiet",
			Value: value.NewBoolean(true),
		},
		ResultFlag:      "quiet",
		ResultBoolValue: true,
	},
	{
		Name: "Set CPU",
		Expr: parser.SetFlag{
			Name:  "@@cpu",
			Value: value.NewInteger(int64(runtime.NumCPU())),
		},
		ResultFlag:         "cpu",
		ResultIntegerValue: runtime.NumCPU(),
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
		Name: "Set CPU Value Error",
		Expr: parser.SetFlag{
			Name:  "@@cpu",
			Value: value.NewString("invalid"),
		},
		Error: "[L:- C:-] SET: flag value 'invalid' for @@cpu is invalid",
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
	flags := cmd.GetFlags()

	for _, v := range setFlagTests {
		initFlag()
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
		case "TIMEZONE":
			if flags.Location != v.ResultStrValue {
				t.Errorf("%s: timezone = %q, want %q", v.Name, flags.Location, v.ResultStrValue)
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
		case "NO_HEADER":
			if flags.NoHeader != v.ResultBoolValue {
				t.Errorf("%s: no-header = %t, want %t", v.Name, flags.NoHeader, v.ResultBoolValue)
			}
		case "WITHOUT_NULL":
			if flags.WithoutNull != v.ResultBoolValue {
				t.Errorf("%s: without-null = %t, want %t", v.Name, flags.WithoutNull, v.ResultBoolValue)
			}
		case "WRITE_ENCODING":
			if flags.WriteEncoding.String() != v.ResultStrValue {
				t.Errorf("%s: write-encoding = %q, want %q", v.Name, flags.WriteEncoding.String(), v.ResultStrValue)
			}
		case "FORMAT":
			if flags.Format.String() != v.ResultStrValue {
				t.Errorf("%s: format = %q, want %q", v.Name, flags.Format.String(), v.ResultStrValue)
			}
		case "WRITE_DELIMITER":
			if string(flags.WriteDelimiter) != v.ResultStrValue {
				t.Errorf("%s: write-delimiter = %q, want %q", v.Name, string(flags.WriteDelimiter), v.ResultStrValue)
			}
		case "WITHOUT_HEADER":
			if flags.WithoutHeader != v.ResultBoolValue {
				t.Errorf("%s: without-header = %t, want %t", v.Name, flags.WithoutHeader, v.ResultBoolValue)
			}
		case "PRETTY_PRINT":
			if flags.PrettyPrint != v.ResultBoolValue {
				t.Errorf("%s: pretty-print = %t, want %t", v.Name, flags.PrettyPrint, v.ResultBoolValue)
			}
		case "COLOR":
			if flags.Color != v.ResultBoolValue {
				t.Errorf("%s: color = %t, want %t", v.Name, flags.Color, v.ResultBoolValue)
			}
		case "QUIET":
			if flags.Quiet != v.ResultBoolValue {
				t.Errorf("%s: quiet = %t, want %t", v.Name, flags.Quiet, v.ResultBoolValue)
			}
		case "CPU":
			if flags.CPU != v.ResultIntegerValue {
				t.Errorf("%s: cpu = %d, want %d", v.Name, flags.CPU, v.ResultIntegerValue)
			}
		case "STATS":
			if flags.Stats != v.ResultBoolValue {
				t.Errorf("%s: stats = %t, want %t", v.Name, flags.Stats, v.ResultBoolValue)
			}
		}
	}
	initFlag()
}

var showFlagTests = []struct {
	Name    string
	Expr    parser.ShowFlag
	SetExpr parser.SetFlag
	Result  string
	Error   string
}{
	{
		Name: "Show Delimiter",
		Expr: parser.ShowFlag{
			Name: "@@delimiter",
		},
		SetExpr: parser.SetFlag{
			Name:  "@@delimiter",
			Value: value.NewString("\t"),
		},
		Result: " @@DELIMITER: '\\t'",
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
		Result: " @@ENCODING: SJIS",
	},
	{
		Name: "Show lineBreak",
		Expr: parser.ShowFlag{
			Name: "@@line_break",
		},
		SetExpr: parser.SetFlag{
			Name:  "@@line_break",
			Value: value.NewString("CRLF"),
		},
		Result: " @@LINE_BREAK: CRLF",
	},
	{
		Name: "Show Timezone",
		Expr: parser.ShowFlag{
			Name: "@@timezone",
		},
		SetExpr: parser.SetFlag{
			Name:  "@@timezone",
			Value: value.NewString("UTC"),
		},
		Result: " @@TIMEZONE: UTC",
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
		Result: " @@REPOSITORY: " + TestDir,
	},
	{
		Name: "Show DatetimeFormat Not Set",
		Expr: parser.ShowFlag{
			Name: "@@datetime_format",
		},
		Result: " @@DATETIME_FORMAT: (not set)",
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
		Result: " @@DATETIME_FORMAT: %Y%m%d",
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
		Result: " @@WAIT_TIMEOUT: 15",
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
		Result: " @@NO_HEADER: true",
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
		Result: " @@WITHOUT_NULL: true",
	},
	{
		Name: "Show WriteEncoding",
		Expr: parser.ShowFlag{
			Name: "@@write_encoding",
		},
		SetExpr: parser.SetFlag{
			Name:  "@@write_encoding",
			Value: value.NewString("SJIS"),
		},
		Result: " @@WRITE_ENCODING: SJIS",
	},
	{
		Name: "Show WriteDelimiter",
		Expr: parser.ShowFlag{
			Name: "@@write_delimiter",
		},
		SetExpr: parser.SetFlag{
			Name:  "@@write_delimiter",
			Value: value.NewString("\t"),
		},
		Result: " @@WRITE_DELIMITER: '\\t'",
	},
	{
		Name: "Show Format",
		Expr: parser.ShowFlag{
			Name: "@@format",
		},
		SetExpr: parser.SetFlag{
			Name:  "@@format",
			Value: value.NewString("json"),
		},
		Result: " @@FORMAT: JSON",
	},
	{
		Name: "Show WithoutHeader",
		Expr: parser.ShowFlag{
			Name: "@@without_header",
		},
		SetExpr: parser.SetFlag{
			Name:  "@@without_header",
			Value: value.NewBoolean(true),
		},
		Result: " @@WITHOUT_HEADER: true",
	},
	{
		Name: "Show PrettyPrint",
		Expr: parser.ShowFlag{
			Name: "@@pretty_print",
		},
		SetExpr: parser.SetFlag{
			Name:  "@@pretty_print",
			Value: value.NewBoolean(true),
		},
		Result: " @@PRETTY_PRINT: true",
	},
	{
		Name: "Show Color",
		Expr: parser.ShowFlag{
			Name: "@@color",
		},
		SetExpr: parser.SetFlag{
			Name:  "@@color",
			Value: value.NewBoolean(false),
		},
		Result: " @@COLOR: false",
	},
	{
		Name: "Show Quiet",
		Expr: parser.ShowFlag{
			Name: "@@quiet",
		},
		SetExpr: parser.SetFlag{
			Name:  "@@quiet",
			Value: value.NewBoolean(true),
		},
		Result: " @@QUIET: true",
	},
	{
		Name: "Show CPU",
		Expr: parser.ShowFlag{
			Name: "@@cpu",
		},
		SetExpr: parser.SetFlag{
			Name:  "@@cpu",
			Value: value.NewInteger(1),
		},
		Result: " @@CPU: 1",
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
		Result: " @@STATS: true",
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
	for _, v := range showFlagTests {
		initFlag()
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
	initFlag()
}

var showObjectsTests = []struct {
	Name        string
	Expr        parser.ShowObjects
	Filter      *Filter
	Repository  string
	ViewCache   ViewMap
	ExecResults []ExecResult
	Expect      string
	Error       string
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
					Encoding:  cmd.SJIS,
					LineBreak: cmd.CRLF,
					NoHeader:  true,
				},
			},
			"TABLE1.TSV": &View{
				Header: NewHeader("table1", []string{"col1", "col2"}),
				FileInfo: &FileInfo{
					Path:      "table1.tsv",
					Delimiter: '\t',
					Format:    cmd.TSV,
					Encoding:  cmd.UTF8,
					LineBreak: cmd.LF,
					NoHeader:  false,
				},
			},
			"TABLE1.JSON": &View{
				Header: NewHeader("table1", []string{"col1", "col2"}),
				FileInfo: &FileInfo{
					Path:        "table1.json",
					JsonQuery:   "{}",
					Format:      cmd.JSON,
					Encoding:    cmd.UTF8,
					LineBreak:   cmd.LF,
					PrettyPrint: false,
				},
			},
			"TABLE2.JSON": &View{
				Header: NewHeader("table2", []string{"col1", "col2"}),
				FileInfo: &FileInfo{
					Path:        "table2.json",
					JsonQuery:   "",
					Format:      cmd.JSON,
					Encoding:    cmd.UTF8,
					LineBreak:   cmd.LF,
					PrettyPrint: false,
				},
			},
			"TABLE1.TXT": &View{
				Header: NewHeader("table1", []string{"col1", "col2"}),
				FileInfo: &FileInfo{
					Path:               "table1.txt",
					DelimiterPositions: []int{3, 12},
					Format:             cmd.FIXED,
					Encoding:           cmd.UTF8,
					LineBreak:          cmd.LF,
					NoHeader:           false,
				},
			},
		},
		Expect: "\n" +
			"                      Loaded Tables\n" +
			"----------------------------------------------------------\n" +
			" table1.csv\n" +
			"     Fields: col1, col2\n" +
			"     Format: CSV     Delimiter: '\\t'\n" +
			"     Encoding: SJIS  LineBreak: CRLF  Header: false\n" +
			" table1.json\n" +
			"     Fields: col1, col2\n" +
			"     Format: JSON    Query: {}\n" +
			"     Encoding: UTF8  LineBreak: LF    Pretty Print: false\n" +
			" table1.tsv\n" +
			"     Fields: col1, col2\n" +
			"     Format: TSV     \n" +
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
					Encoding:  cmd.SJIS,
					LineBreak: cmd.CRLF,
					NoHeader:  true,
				},
			},
			"TABLE1.TSV": &View{
				Header: NewHeader("table1", []string{"col1", "col2"}),
				FileInfo: &FileInfo{
					Path:      "table1.tsv",
					Delimiter: '\t',
					Format:    cmd.TSV,
					Encoding:  cmd.UTF8,
					LineBreak: cmd.LF,
					NoHeader:  false,
				},
			},
			"TABLE1.JSON": &View{
				Header: NewHeader("table1", []string{"col1", "col2"}),
				FileInfo: &FileInfo{
					Path:        "table1.json",
					JsonQuery:   "{}",
					Format:      cmd.JSON,
					Encoding:    cmd.UTF8,
					LineBreak:   cmd.LF,
					PrettyPrint: false,
				},
			},
			"TABLE2.JSON": &View{
				Header: NewHeader("table2", []string{"col1", "col2"}),
				FileInfo: &FileInfo{
					Path:        "table2.json",
					JsonQuery:   "",
					Format:      cmd.JSON,
					Encoding:    cmd.UTF8,
					LineBreak:   cmd.LF,
					PrettyPrint: false,
				},
			},
			"TABLE1.TXT": &View{
				Header: NewHeader("table1", []string{"col1", "col2"}),
				FileInfo: &FileInfo{
					Path:               "table1.txt",
					DelimiterPositions: []int{3, 12},
					Format:             cmd.FIXED,
					Encoding:           cmd.UTF8,
					LineBreak:          cmd.LF,
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
			"     Format: CSV     Delimiter: '\\t'\n" +
			"     Encoding: SJIS  LineBreak: CRLF  Header: false\n" +
			" table1.json\n" +
			"     Fields: col1, col2\n" +
			"     Format: JSON    Query: {}\n" +
			"     Encoding: UTF8  LineBreak: LF    Pretty Print: false\n" +
			" *Created* table1.tsv\n" +
			"     Fields: col1, col2\n" +
			"     Format: TSV     \n" +
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
					Encoding:  cmd.SJIS,
					LineBreak: cmd.CRLF,
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
			"     Format: CSV     Delimiter: '\\t'\n" +
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
							{Name: "@arg1"},
						},
						Statements: []parser.Statement{
							parser.Print{Value: parser.Variable{Name: "@arg1"}},
						},
					},
				},
				UserDefinedFunctionMap{
					"USERAGGFUNC": &UserDefinedFunction{
						Name: parser.Identifier{Literal: "useraggfunc"},
						Parameters: []parser.Variable{
							{Name: "@arg1"},
							{Name: "@arg2"},
						},
						Defaults: map[string]parser.QueryExpression{
							"@arg2": parser.NewIntegerValue(1),
						},
						IsAggregate:  true,
						RequiredArgs: 1,
						Cursor:       parser.Identifier{Literal: "column1"},
						Statements: []parser.Statement{
							parser.Print{Value: parser.Variable{Name: "@var1"}},
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
			"             Flags\n" +
			"-------------------------------\n" +
			"        @@DELIMITER: ','\n" +
			"         @@ENCODING: UTF8\n" +
			"       @@LINE_BREAK: LF\n" +
			"         @@TIMEZONE: UTC\n" +
			"       @@REPOSITORY: .\n" +
			"  @@DATETIME_FORMAT: (not set)\n" +
			"        @@NO_HEADER: false\n" +
			"     @@WITHOUT_NULL: false\n" +
			"     @@WAIT_TIMEOUT: 15\n" +
			"   @@WRITE_ENCODING: UTF8\n" +
			"           @@FORMAT: TEXT\n" +
			"  @@WRITE_DELIMITER: ','\n" +
			"   @@WITHOUT_HEADER: false\n" +
			"     @@PRETTY_PRINT: false\n" +
			"            @@COLOR: false\n" +
			"            @@QUIET: false\n" +
			"              @@CPU: " + strconv.Itoa(cmd.GetFlags().CPU) + "\n" +
			"            @@STATS: false\n" +
			"",
	},
	{
		Name:  "ShowObjects Invalid Object Type",
		Expr:  parser.ShowObjects{Type: parser.Identifier{Literal: "invalid"}},
		Error: "[L:- C:-] SHOW: object type invalid is invalid",
	},
}

func TestShowObjects(t *testing.T) {
	initFlag()
	flags := cmd.GetFlags()

	for _, v := range showObjectsTests {
		flags.Repository = v.Repository
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
					Encoding:  cmd.UTF8,
					LineBreak: cmd.LF,
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
			strings.Repeat(" ", (calcShowFieldsWidth("show_fields_create.csv")-(10+len("show_fields_create.csv")))/2) + "Fields in show_fields_create.csv\n" +
			strings.Repeat("-", calcShowFieldsWidth("show_fields_create.csv")) + "\n" +
			" Type: Table\n" +
			" Path: " + GetTestFilePath("show_fields_create.csv") + "\n" +
			" Format: CSV     Delimiter: ','\n" +
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
					Encoding:  cmd.UTF8,
					LineBreak: cmd.LF,
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
			strings.Repeat(" ", (calcShowFieldsWidth("show_fields_update.csv")-(10+len("show_fields_update.csv")))/2) + "Fields in show_fields_update.csv\n" +
			strings.Repeat("-", calcShowFieldsWidth("show_fields_create.csv")) + "\n" +
			" Type: Table\n" +
			" Path: " + GetTestFilePath("show_fields_update.csv") + "\n" +
			" Format: CSV     Delimiter: ','\n" +
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
		Error: "[L:- C:-] SHOW: object type invalid is invalid",
	},
}

func calcShowFieldsWidth(fileName string) int {
	w := 47
	pathLen := 7 + len(GetTestFilePath(fileName))
	titleLen := 10 + len(fileName)

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
