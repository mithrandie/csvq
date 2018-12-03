package cmd

import (
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/mithrandie/go-text"
	"github.com/mithrandie/go-text/json"

	"github.com/mithrandie/csvq/lib/file"
)

func TestFlags_SelectImportFormat(t *testing.T) {
	flags := GetFlags()

	flags.SetJsonQuery("{}")
	format := flags.SelectImportFormat()
	expect := JSON
	if format != expect {
		t.Errorf("import-format = %q, want %q", format.String(), expect.String())
	}

	flags.SetJsonQuery("")
	flags.SetDelimiter("SPACES")
	format = flags.SelectImportFormat()
	expect = FIXED
	if format != expect {
		t.Errorf("import-format = %q, want %q", format.String(), expect.String())
	}

	flags.SetDelimiter("\\t")
	format = flags.SelectImportFormat()
	expect = TSV
	if format != expect {
		t.Errorf("import-format = %q, want %q", format.String(), expect.String())
	}

	flags.SetDelimiter(",")
	format = flags.SelectImportFormat()
	expect = CSV
	if format != expect {
		t.Errorf("import-format = %q, want %q", format.String(), expect.String())
	}
}

func TestFlags_SetRepository(t *testing.T) {
	flags := GetFlags()

	flags.SetRepository("")
	if flags.Repository != "" {
		t.Errorf("repository = %s, expect to set %q for %q", flags.Repository, "", "")
	}

	dir := filepath.Join("..", "..", "lib", "cmd")
	absdir, _ := filepath.Abs(dir)
	flags.SetRepository(dir)
	if flags.Repository != absdir {
		t.Errorf("repository = %s, expect to set %s for %s", flags.Repository, absdir, dir)
	}

	expectErr := "repository does not exist"
	err := flags.SetRepository("notexists")
	if err == nil {
		t.Errorf("no error, want error %q for %s", expectErr, "error")
	} else if err.Error() != expectErr {
		t.Errorf("error = %q, want error %q for %s", err.Error(), expectErr, "notexists")
	}

	expectErr = "repository must be a directory path"
	err = flags.SetRepository("flags_test.go")
	if err == nil {
		t.Errorf("no error, want error %q for %s", expectErr, "error")
	} else if err.Error() != expectErr {
		t.Errorf("error = %q, want error %q for %s", err.Error(), expectErr, "flags_test.go")
	}
}

func TestFlags_SetLocation(t *testing.T) {
	flags := GetFlags()

	s := ""
	flags.SetLocation(s)
	if flags.Location != "Local" {
		t.Errorf("location = %s, expect to set %s for %q", flags.Location, "Local", "")
	}

	s = "local"
	flags.SetLocation(s)
	if flags.Location != "Local" {
		t.Errorf("location = %s, expect to set %s for %q", flags.Location, "Local", s)
	}

	s = "utc"
	flags.SetLocation(s)
	if flags.Location != "UTC" {
		t.Errorf("location = %s, expect to set %s for %q", flags.Location, "UTC", s)
	}

	s = "America/NotExist"
	expectErr := "timezone \"America/NotExist\" does not exist"
	err := flags.SetLocation(s)
	if err == nil {
		t.Errorf("no error, want error %q for %s", expectErr, s)
	} else if err.Error() != expectErr {
		t.Errorf("error = %q, want error %q for %s", err.Error(), expectErr, s)
	}
}

func TestFlags_SetDatetimeFormat(t *testing.T) {
	flags := GetFlags()

	format := "%Y-%m-%d"
	flags.SetDatetimeFormat(format)
	expect := []string{
		"%Y-%m-%d",
	}
	if !reflect.DeepEqual(flags.DatetimeFormat, expect) {
		t.Errorf("datetime format = %s, expect to set %s", flags.DatetimeFormat, expect)
	}

	format = ""
	flags.SetDatetimeFormat(format)
	expect = []string{
		"%Y-%m-%d",
	}
	if !reflect.DeepEqual(flags.DatetimeFormat, expect) {
		t.Errorf("datetime format = %s, expect to set %s", flags.DatetimeFormat, expect)
	}

	format = "[\"%Y-%m-%d %H:%i:%s\"]"
	flags.SetDatetimeFormat(format)
	expect = []string{
		"%Y-%m-%d",
		"%Y-%m-%d %H:%i:%s",
	}
	if !reflect.DeepEqual(flags.DatetimeFormat, expect) {
		t.Errorf("datetime format = %s, expect to set %s", flags.DatetimeFormat, expect)
	}
}

func TestFlags_SetWaitTimeout(t *testing.T) {
	flags := GetFlags()

	var f float64 = -1
	flags.SetWaitTimeout(f)
	if flags.WaitTimeout != 0 {
		t.Errorf("wait timeout = %f, expect to set %f for %f", flags.WaitTimeout, 0.0, f)
	}

	f = 15
	flags.SetWaitTimeout(f)
	if flags.WaitTimeout != 15 {
		t.Errorf("wait timeout = %f, expect to set %f for %f", flags.WaitTimeout, 15.0, f)
	}

	if file.WaitTimeout != 15 {
		t.Errorf("wait timeout in the file package = %f, expect to set %f for %f", file.WaitTimeout, 15.0, f)
	}
}

func TestFlags_SetDelimiter(t *testing.T) {
	flags := GetFlags()

	flags.SetDelimiter("")
	if flags.Delimiter != ',' {
		t.Errorf("delimiter = %q, expect to set %q for %q", flags.Delimiter, ',', "")
	}

	flags.SetDelimiter("\\t")
	if flags.Delimiter != '\t' {
		t.Errorf("delimiter = %q, expect to set %q for %q", flags.Delimiter, "\t", "\t")
	}

	flags.SetDelimiter("[1, 2, 3]")
	if flags.DelimitAutomatically != false {
		t.Errorf("delimitAutomatically = %t, expect to set %t for %q", flags.DelimitAutomatically, false, "[1, 2, 3]")
	}
	if !reflect.DeepEqual(flags.DelimiterPositions, []int{1, 2, 3}) {
		t.Errorf("delimitPositions = %v, expect to set %v for %q", flags.DelimiterPositions, []int{1, 2, 3}, "[1, 2, 3]")
	}

	flags.SetDelimiter("spaces")
	if flags.DelimitAutomatically != true {
		t.Errorf("delimitAutomatically = %t, expect to set %t for %q", flags.DelimitAutomatically, true, "spaces")
	}
	if flags.DelimiterPositions != nil {
		t.Errorf("delimitPositions = %v, expect to set %v for %q", flags.DelimiterPositions, nil, "spaces")
	}

	expectErr := "delimiter must be one character, \"SPACES\" or JSON array of integers"
	err := flags.SetDelimiter("[a]")
	if err == nil {
		t.Errorf("no error, want error %q for %s", expectErr, "//")
	} else if err.Error() != expectErr {
		t.Errorf("error = %q, want error %q for %s", err.Error(), expectErr, "//")
	}

	expectErr = "delimiter must be one character, \"SPACES\" or JSON array of integers"
	err = flags.SetDelimiter("//")
	if err == nil {
		t.Errorf("no error, want error %q for %s", expectErr, "//")
	} else if err.Error() != expectErr {
		t.Errorf("error = %q, want error %q for %s", err.Error(), expectErr, "//")
	}
}

func TestFlags_SetJsonQuery(t *testing.T) {
	flags := GetFlags()

	flags.SetJsonQuery("{}")
	if flags.JsonQuery != "{}" {
		t.Errorf("json-query = %q, expect to set %q", flags.JsonQuery, "{}")
	}
}

func TestFlags_SetEncoding(t *testing.T) {
	flags := GetFlags()

	flags.SetEncoding("sjis")
	if flags.Encoding != text.SJIS {
		t.Errorf("encoding = %s, expect to set %s for %s", flags.Encoding, text.SJIS, "sjis")
	}

	expectErr := "encoding must be one of UTF8|SJIS"
	err := flags.SetEncoding("error")
	if err == nil {
		t.Errorf("no error, want error %q for %s", expectErr, "error")
	} else if err.Error() != expectErr {
		t.Errorf("error = %q, want error %q for %s", err.Error(), expectErr, "error")
	}
}

func TestFlags_SetNoHeader(t *testing.T) {
	flags := GetFlags()

	flags.SetNoHeader(true)
	if !flags.NoHeader {
		t.Errorf("no-header = %t, expect to set %t", flags.NoHeader, true)
	}
}

func TestFlags_SetWithoutNull(t *testing.T) {
	flags := GetFlags()

	flags.SetWithoutNull(true)
	if !flags.WithoutNull {
		t.Errorf("without-null = %t, expect to set %t", flags.WithoutNull, true)
	}
}

func TestFlags_SetFormat(t *testing.T) {
	flags := GetFlags()

	flags.SetFormat("", "")
	if flags.Format != TEXT {
		t.Errorf("format = %s, expect to set %s for empty string", flags.Format, TEXT)
	}

	flags.SetFormat("", "foo.csv")
	if flags.Format != CSV {
		t.Errorf("format = %s, expect to set %s for empty string with file %q", flags.Format, CSV, "foo.csv")
	}

	flags.SetFormat("", "foo.tsv")
	if flags.Format != TSV {
		t.Errorf("format = %s, expect to set %s for empty string with file %q", flags.Format, TSV, "foo.tsv")
	}

	flags.SetFormat("", "foo.txt")
	if flags.Format != FIXED {
		t.Errorf("format = %s, expect to set %s for empty string with file %q", flags.Format, FIXED, "foo.txt")
	}

	flags.SetFormat("", "foo.json")
	if flags.Format != JSON {
		t.Errorf("format = %s, expect to set %s for empty string with file %q", flags.Format, JSON, "foo.json")
	}

	flags.SetFormat("", "foo.md")
	if flags.Format != GFM {
		t.Errorf("format = %s, expect to set %s for empty string with file %q", flags.Format, GFM, "foo.md")
	}

	flags.SetFormat("", "foo.org")
	if flags.Format != ORG {
		t.Errorf("format = %s, expect to set %s for empty string with file %q", flags.Format, ORG, "foo.org")
	}

	flags.SetFormat("csv", "")
	if flags.Format != CSV {
		t.Errorf("format = %s, expect to set %s for %s", flags.Format, CSV, "csv")
	}

	flags.SetFormat("tsv", "")
	if flags.Format != TSV {
		t.Errorf("format = %s, expect to set %s for %s", flags.Format, TSV, "tsv")
	}

	flags.SetFormat("fixed", "")
	if flags.Format != FIXED {
		t.Errorf("format = %s, expect to set %s for %s", flags.Format, FIXED, "fixed")
	}

	flags.SetFormat("json", "")
	if flags.Format != JSON {
		t.Errorf("format = %s, expect to set %s for %s", flags.Format, JSON, "json")
	}

	flags.SetFormat("jsonh", "")
	if flags.Format != JSON {
		t.Errorf("format = %s, expect to set %s for %s", flags.Format, JSON, "jsonh")
	}
	if flags.JsonEscape != json.HexDigits {
		t.Errorf("json escape type = %v, expect to set %v for %s", flags.JsonEscape, json.HexDigits, "jsonh")
	}

	flags.SetFormat("jsona", "")
	if flags.Format != JSON {
		t.Errorf("format = %s, expect to set %s for %s", flags.Format, JSON, "jsona")
	}
	if flags.JsonEscape != json.AllWithHexDigits {
		t.Errorf("json escape type = %v, expect to set %v for %s", flags.JsonEscape, json.AllWithHexDigits, "jsonh")
	}

	flags.SetFormat("gfm", "")
	if flags.Format != GFM {
		t.Errorf("format = %s, expect to set %s for %s", flags.Format, GFM, "gfm")
	}

	flags.SetFormat("org", "")
	if flags.Format != ORG {
		t.Errorf("format = %s, expect to set %s for %s", flags.Format, ORG, "org")
	}

	flags.SetFormat("text", "")
	if flags.Format != TEXT {
		t.Errorf("format = %s, expect to set %s for %s", flags.Format, TEXT, "text")
	}

	expectErr := "format must be one of CSV|TSV|FIXED|JSON|GFM|ORG|TEXT|JSONH|JSONA"
	err := flags.SetFormat("error", "")
	if err == nil {
		t.Errorf("no error, want error %q for %s", expectErr, "error")
	} else if err.Error() != expectErr {
		t.Errorf("error = %q, want error %q for %s", err.Error(), expectErr, "error")
	}
}

func TestFlags_SetWriteEncoding(t *testing.T) {
	flags := GetFlags()

	flags.SetWriteEncoding("sjis")
	if flags.Encoding != text.SJIS {
		t.Errorf("encoding = %s, expect to set %s for %s", flags.WriteEncoding, text.SJIS, "sjis")
	}

	expectErr := "encoding must be one of UTF8|SJIS"
	err := flags.SetWriteEncoding("error")
	if err == nil {
		t.Errorf("no error, want error %q for %s", expectErr, "error")
	} else if err.Error() != expectErr {
		t.Errorf("error = %q, want error %q for %s", err.Error(), expectErr, "error")
	}
}

func TestFlags_SetWriteDelimiter(t *testing.T) {
	flags := GetFlags()

	flags.SetWriteDelimiter("")
	if flags.WriteDelimiter != ',' {
		t.Errorf("write-delimiter = %q, expect to set %q for %q, format = %s", flags.WriteDelimiter, ',', "", flags.Format)
	}

	flags.SetWriteDelimiter("\\t")
	if flags.WriteDelimiter != '\t' {
		t.Errorf("write-delimiter = %q, expect to set %q for %q", flags.WriteDelimiter, "\t", "\t")
	}

	flags.SetWriteDelimiter("[1, 2, 3]")
	if !reflect.DeepEqual(flags.WriteDelimiterPositions, []int{1, 2, 3}) {
		t.Errorf("writeDelimitPositions = %v, expect to set %v for %q", flags.WriteDelimiterPositions, []int{1, 2, 3}, "[1, 2, 3]")
	}

	expectErr := "write-delimiter must be one character, \"SPACES\" or JSON array of integers"
	err := flags.SetWriteDelimiter("//")
	if err == nil {
		t.Errorf("no error, want error %q for %s", expectErr, "//")
	} else if err.Error() != expectErr {
		t.Errorf("error = %q, want error %q for %s", err.Error(), expectErr, "//")
	}
}

func TestFlags_SetWithoutHeader(t *testing.T) {
	flags := GetFlags()

	flags.SetWithoutHeader(true)
	if !flags.NoHeader {
		t.Errorf("without-header = %t, expect to set %t", flags.WithoutHeader, true)
	}
}

func TestFlags_SetLineBreak(t *testing.T) {
	flags := GetFlags()

	flags.SetLineBreak("")
	if flags.LineBreak != text.LF {
		t.Errorf("line-break = %s, expect to set %s for %q", flags.LineBreak, text.LF, "")
	}

	flags.SetLineBreak("crlf")
	if flags.LineBreak != text.CRLF {
		t.Errorf("line-break = %s, expect to set %s for %s", flags.LineBreak, text.CRLF, "crlf")
	}

	flags.SetLineBreak("cr")
	if flags.LineBreak != text.CR {
		t.Errorf("line-break = %s, expect to set %s for %s", flags.LineBreak, text.CR, "cr")
	}

	flags.SetLineBreak("lf")
	if flags.LineBreak != text.LF {
		t.Errorf("line-break = %s, expect to set %s for %s", flags.LineBreak, text.LF, "LF")
	}

	expectErr := "line-break must be one of CRLF|LF|CR"
	err := flags.SetLineBreak("error")
	if err == nil {
		t.Errorf("no error, want error %q for %s", expectErr, "error")
	} else if err.Error() != expectErr {
		t.Errorf("error = %q, want error %q for %s", err.Error(), expectErr, "error")
	}
}

func TestFlags_SetEncloseAll(t *testing.T) {
	flags := GetFlags()

	flags.SetEncloseAll(true)
	if !flags.EncloseAll {
		t.Errorf("enclose-all = %t, expect to set %t", flags.EncloseAll, true)
	}
}

func TestFlags_SetJsonEscape(t *testing.T) {
	flags := GetFlags()

	s := "backslash"
	flags.SetJsonEscape(s)
	if flags.JsonEscape != json.Backslash {
		t.Errorf("json-escape = %v, expect to set %v", flags.JsonEscape, json.Backslash)
	}

	s = "hex"
	flags.SetJsonEscape(s)
	if flags.JsonEscape != json.HexDigits {
		t.Errorf("json-escape = %v, expect to set %v", flags.JsonEscape, json.HexDigits)
	}

	s = "hexall"
	flags.SetJsonEscape(s)
	if flags.JsonEscape != json.AllWithHexDigits {
		t.Errorf("json-escape = %v, expect to set %v", flags.JsonEscape, json.AllWithHexDigits)
	}

	s = "error"
	expectErr := "json-escape must be one of BACKSLASH|HEX|HEXALL"
	err := flags.SetJsonEscape(s)
	if err == nil {
		t.Errorf("no error, want error %q for %s", expectErr, s)
	} else if err.Error() != expectErr {
		t.Errorf("error = %q, want error %q for %s", err.Error(), expectErr, s)
	}
}

func TestFlags_SetPrettyPrint(t *testing.T) {
	flags := GetFlags()

	flags.SetPrettyPrint(true)
	if !flags.PrettyPrint {
		t.Errorf("pretty-print = %t, expect to set %t", flags.PrettyPrint, true)
	}
}

func TestFlags_SetEastAsianEncoding(t *testing.T) {
	flags := GetFlags()

	flags.SetEastAsianEncoding(true)
	if !flags.EastAsianEncoding {
		t.Errorf("east-asian-encoding = %t, expect to set %t", flags.EastAsianEncoding, true)
	}
}

func TestFlags_SetCountDiacriticalSign(t *testing.T) {
	flags := GetFlags()

	flags.SetCountDiacriticalSign(true)
	if !flags.CountDiacriticalSign {
		t.Errorf("count-diacritical-sign = %t, expect to set %t", flags.CountDiacriticalSign, true)
	}
}

func TestFlags_SetCountFormatCode(t *testing.T) {
	flags := GetFlags()

	flags.SetCountFormatCode(true)
	if !flags.CountFormatCode {
		t.Errorf("count-format-code = %t, expect to set %t", flags.CountFormatCode, true)
	}
}

func TestFlags_SetColor(t *testing.T) {
	flags := GetFlags()

	flags.SetColor(true)
	if !flags.Color {
		t.Errorf("color = %t, expect to set %t", flags.Color, true)
	}
	flags.SetColor(false)
}

func TestFlags_SetQuiet(t *testing.T) {
	flags := GetFlags()

	flags.SetQuiet(true)
	if !flags.Quiet {
		t.Errorf("silent = %t, expect to set %t", flags.Quiet, true)
	}
}

func TestFlags_SetCPU(t *testing.T) {
	flags := GetFlags()

	flags.SetCPU(0)
	expect := 1
	if expect != flags.CPU {
		t.Errorf("cpu = %d, expect to set %d", flags.CPU, 1)
	}

	flags.SetCPU(runtime.NumCPU() + 100)
	if runtime.NumCPU() != flags.CPU {
		t.Errorf("cpu = %d, expect to set %d", flags.CPU, runtime.NumCPU())
	}
}

func TestFlags_SetStats(t *testing.T) {
	flags := GetFlags()

	flags.SetStats(true)
	if !flags.Stats {
		t.Errorf("stats = %t, expect to set %t", flags.Stats, true)
	}
}
