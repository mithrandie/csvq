package cmd

import (
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/mithrandie/go-text"
	"github.com/mithrandie/go-text/json"

	"github.com/mithrandie/csvq/lib/file"
)

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

	if file.WaitTimeout != 15*time.Second {
		t.Errorf("wait timeout in the file package = %v, expect to set %v for %f", file.WaitTimeout, 15*time.Second, f)
	}
}

func TestFlags_SetImportFormat(t *testing.T) {
	flags := GetFlags()

	flags.SetImportFormat("")
	if flags.ImportFormat != CSV {
		t.Errorf("importFormat = %s, expect to set %s for empty string", flags.ImportFormat, CSV)
	}

	flags.SetImportFormat("json")
	if flags.ImportFormat != JSON {
		t.Errorf("importFormat = %s, expect to set %s for empty string", flags.ImportFormat, JSON)
	}

	expectErr := "import format must be one of CSV|TSV|FIXED|JSON|LTSV"
	err := flags.SetImportFormat("error")
	if err == nil {
		t.Errorf("no error, want error %q for %s", expectErr, "error")
	} else if err.Error() != expectErr {
		t.Errorf("error = %q, want error %q for %s", err.Error(), expectErr, "error")
	}

	err = flags.SetImportFormat("text")
	if err == nil {
		t.Errorf("no error, want error %q for %s", expectErr, "error")
	} else if err.Error() != expectErr {
		t.Errorf("error = %q, want error %q for %s", err.Error(), expectErr, "error")
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

	expectErr := "delimiter must be one character"
	err := flags.SetDelimiter("[a]")
	if err == nil {
		t.Errorf("no error, want error %q for %s", expectErr, "//")
	} else if err.Error() != expectErr {
		t.Errorf("error = %q, want error %q for %s", err.Error(), expectErr, "//")
	}

	expectErr = "delimiter must be one character"
	err = flags.SetDelimiter("//")
	if err == nil {
		t.Errorf("no error, want error %q for %s", expectErr, "//")
	} else if err.Error() != expectErr {
		t.Errorf("error = %q, want error %q for %s", err.Error(), expectErr, "//")
	}
}

func TestFlags_SetDelimiterPositions(t *testing.T) {
	flags := GetFlags()

	flags.SetDelimiterPositions("")
	if flags.DelimiterPositions != nil {
		t.Errorf("delimiter-positions = %v, expect to set %v for %q", flags.DelimiterPositions, nil, "")
	}

	flags.SetDelimiterPositions("s[1, 2, 3]")
	if flags.SingleLine != true {
		t.Errorf("singleLine = %t, expect to set %t for %q", flags.SingleLine, true, "s[1, 2, 3]")
	}
	if !reflect.DeepEqual(flags.DelimiterPositions, []int{1, 2, 3}) {
		t.Errorf("delimitPositions = %v, expect to set %v for %q", flags.DelimiterPositions, []int{1, 2, 3}, "[1, 2, 3]")
	}

	flags.SetDelimiterPositions("[1, 2, 3]")
	if flags.SingleLine != false {
		t.Errorf("singleLine = %t, expect to set %t for %q", flags.SingleLine, false, "[1, 2, 3]")
	}
	if !reflect.DeepEqual(flags.DelimiterPositions, []int{1, 2, 3}) {
		t.Errorf("delimitPositions = %v, expect to set %v for %q", flags.DelimiterPositions, []int{1, 2, 3}, "[1, 2, 3]")
	}

	flags.SetDelimiterPositions("spaces")
	if flags.SingleLine != false {
		t.Errorf("singleLine = %t, expect to set %t for %q", flags.SingleLine, false, "spaces")
	}
	if flags.DelimiterPositions != nil {
		t.Errorf("delimitPositions = %v, expect to set %v for %q", flags.DelimiterPositions, nil, "spaces")
	}

	expectErr := "delimiter positions must be \"SPACES\" or a JSON array of integers"
	err := flags.SetDelimiterPositions("[a]")
	if err == nil {
		t.Errorf("no error, want error %q for %s", expectErr, "//")
	} else if err.Error() != expectErr {
		t.Errorf("error = %q, want error %q for %s", err.Error(), expectErr, "//")
	}

	err = flags.SetDelimiterPositions("//")
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

	expectErr := "encoding must be one of UTF8|UTF8M|SJIS"
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

	flags.SetFormat("", "foo.json")
	if flags.Format != JSON {
		t.Errorf("format = %s, expect to set %s for empty string with file %q", flags.Format, JSON, "foo.json")
	}

	flags.SetFormat("", "foo.ltsv")
	if flags.Format != LTSV {
		t.Errorf("format = %s, expect to set %s for empty string with file %q", flags.Format, LTSV, "foo.ltsv")
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

	flags.SetFormat("ltsv", "")
	if flags.Format != LTSV {
		t.Errorf("format = %s, expect to set %s for %s", flags.Format, LTSV, "ltsv")
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

	expectErr := "format must be one of CSV|TSV|FIXED|JSON|LTSV|GFM|ORG|TEXT"
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

	expectErr := "encoding must be one of UTF8|UTF8M|SJIS"
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

	expectErr := "write-delimiter must be one character"
	err := flags.SetWriteDelimiter("//")
	if err == nil {
		t.Errorf("no error, want error %q for %s", expectErr, "//")
	} else if err.Error() != expectErr {
		t.Errorf("error = %q, want error %q for %s", err.Error(), expectErr, "//")
	}
}

func TestFlags_SetWriteDelimiterPositions(t *testing.T) {
	flags := GetFlags()

	flags.SetWriteDelimiterPositions("s[1, 2, 3]")
	if flags.WriteAsSingleLine != true {
		t.Errorf("WriteAsSingleLine = %t, expect to set %t for %q", flags.WriteAsSingleLine, true, "s[1, 2, 3]")
	}
	if !reflect.DeepEqual(flags.WriteDelimiterPositions, []int{1, 2, 3}) {
		t.Errorf("WriteDelimiterPositions = %v, expect to set %v for %q", flags.WriteDelimiterPositions, []int{1, 2, 3}, "s[1, 2, 3]")
	}

	flags.SetWriteDelimiterPositions("[1, 2, 3]")
	if flags.WriteAsSingleLine != false {
		t.Errorf("WriteAsSingleLine = %t, expect to set %t for %q", flags.WriteAsSingleLine, false, "[1, 2, 3]")
	}
	if !reflect.DeepEqual(flags.WriteDelimiterPositions, []int{1, 2, 3}) {
		t.Errorf("WriteDelimiterPositions = %v, expect to set %v for %q", flags.WriteDelimiterPositions, []int{1, 2, 3}, "[1, 2, 3]")
	}

	expectErr := "write-delimiter-positions must be \"SPACES\" or a JSON array of integers"
	err := flags.SetWriteDelimiterPositions("//")
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
	expectErr := "json escape type must be one of BACKSLASH|HEX|HEXALL"
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
