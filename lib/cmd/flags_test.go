package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/mithrandie/go-text"

	"github.com/mithrandie/csvq/lib/file"
)

func TestFlags_ImportFormat(t *testing.T) {
	flags := GetFlags()

	SetJsonQuery("{}")
	format := flags.ImportFormat()
	expect := JSON
	if format != expect {
		t.Errorf("import-format = %q, want %q", format.String(), expect.String())
	}

	SetJsonQuery("")
	SetDelimiter("SPACES")
	format = flags.ImportFormat()
	expect = FIXED
	if format != expect {
		t.Errorf("import-format = %q, want %q", format.String(), expect.String())
	}

	SetDelimiter("\\t")
	format = flags.ImportFormat()
	expect = TSV
	if format != expect {
		t.Errorf("import-format = %q, want %q", format.String(), expect.String())
	}

	SetDelimiter(",")
	format = flags.ImportFormat()
	expect = CSV
	if format != expect {
		t.Errorf("import-format = %q, want %q", format.String(), expect.String())
	}
}

func TestSetRepository(t *testing.T) {
	flags := GetFlags()

	pwd, _ := os.Getwd()

	SetRepository("")
	if flags.Repository != pwd {
		t.Errorf("repository = %s, expect to set %s for %q", flags.Repository, pwd, "")
	}

	dir := filepath.Join("..", "..", "lib", "cmd")
	absdir, _ := filepath.Abs(dir)
	SetRepository(dir)
	if flags.Repository != absdir {
		t.Errorf("repository = %s, expect to set %s for %s", flags.Repository, absdir, dir)
	}

	expectErr := "repository does not exist"
	err := SetRepository("notexists")
	if err == nil {
		t.Errorf("no error, want error %q for %s", expectErr, "error")
	} else if err.Error() != expectErr {
		t.Errorf("error = %q, want error %q for %s", err.Error(), expectErr, "notexists")
	}

	expectErr = "repository must be a directory path"
	err = SetRepository("flags_test.go")
	if err == nil {
		t.Errorf("no error, want error %q for %s", expectErr, "error")
	} else if err.Error() != expectErr {
		t.Errorf("error = %q, want error %q for %s", err.Error(), expectErr, "flags_test.go")
	}
}

func TestSetLocation(t *testing.T) {
	flags := GetFlags()

	s := ""
	SetLocation(s)
	if flags.Location != "Local" {
		t.Errorf("location = %s, expect to set %s for %q", flags.Location, "Local", "")
	}

	s = "local"
	SetLocation(s)
	if flags.Location != "Local" {
		t.Errorf("location = %s, expect to set %s for %q", flags.Location, "Local", s)
	}

	s = "utc"
	SetLocation(s)
	if flags.Location != "UTC" {
		t.Errorf("location = %s, expect to set %s for %q", flags.Location, "UTC", s)
	}

	s = "America/NotExist"
	expectErr := "timezone does not exist"
	err := SetLocation(s)
	if err == nil {
		t.Errorf("no error, want error %q for %s", expectErr, s)
	} else if err.Error() != expectErr {
		t.Errorf("error = %q, want error %q for %s", err.Error(), expectErr, s)
	}
}

func TestSetDatetimeFormat(t *testing.T) {
	flags := GetFlags()

	format := "%Y-%m-%d"
	SetDatetimeFormat(format)
	if flags.DatetimeFormat != format {
		t.Errorf("datetime format = %s, expect to set %s", flags.DatetimeFormat, format)
	}
}

func TestSetWaitTimeout(t *testing.T) {
	flags := GetFlags()

	var f float64 = -1
	SetWaitTimeout(f)
	if flags.WaitTimeout != 0 {
		t.Errorf("wait timeout = %f, expect to set %f for %f", flags.WaitTimeout, 0.0, f)
	}

	f = 15
	SetWaitTimeout(f)
	if flags.WaitTimeout != 15 {
		t.Errorf("wait timeout = %f, expect to set %f for %f", flags.WaitTimeout, 15.0, f)
	}

	if file.WaitTimeout != 15 {
		t.Errorf("wait timeout in the file package = %f, expect to set %f for %f", file.WaitTimeout, 15.0, f)
	}
}

func TestSetSource(t *testing.T) {
	flags := GetFlags()

	SetSource("")
	if flags.Source != "" {
		t.Errorf("source = %s, expect to set %q for %q", flags.Source, "", "")
	}

	s := filepath.Join("..", "..", "lib", "cmd", "flags_test.go")
	expect, _ := filepath.Abs(s)
	SetSource(s)
	if flags.Source != expect {
		t.Errorf("source = %s, expect to set %s for %s", flags.Source, expect, s)
	}

	s = filepath.Join("..", "..", "lib", "cmd", "notexist")
	expectErr := "source file does not exist"
	err := SetSource(s)
	if err == nil {
		t.Errorf("no error, want error %q for %s", expectErr, "error")
	} else if err.Error() != expectErr {
		t.Errorf("error = %q, want error %q for %s", err.Error(), expectErr, "notexists")
	}

	s = filepath.Join("..", "..", "lib", "cmd")
	expectErr = "source file must be a readable file"
	err = SetSource(s)
	if err == nil {
		t.Errorf("no error, want error %q for %s", expectErr, "error")
	} else if err.Error() != expectErr {
		t.Errorf("error = %q, want error %q for %s", err.Error(), expectErr, "notexists")
	}
}

func TestSetDelimiter(t *testing.T) {
	flags := GetFlags()

	SetDelimiter("")
	if flags.Delimiter != ',' {
		t.Errorf("delimiter = %q, expect to set %q for %q", flags.Delimiter, ',', "")
	}

	SetDelimiter("\\t")
	if flags.Delimiter != '\t' {
		t.Errorf("delimiter = %q, expect to set %q for %q", flags.Delimiter, "\t", "\t")
	}

	SetDelimiter("[1, 2, 3]")
	if flags.DelimitAutomatically != false {
		t.Errorf("delimitAutomatically = %t, expect to set %t for %q", flags.DelimitAutomatically, false, "[1, 2, 3]")
	}
	if !reflect.DeepEqual(flags.DelimiterPositions, []int{1, 2, 3}) {
		t.Errorf("delimitPositions = %v, expect to set %v for %q", flags.DelimiterPositions, []int{1, 2, 3}, "[1, 2, 3]")
	}

	SetDelimiter("spaces")
	if flags.DelimitAutomatically != true {
		t.Errorf("delimitAutomatically = %t, expect to set %t for %q", flags.DelimitAutomatically, true, "spaces")
	}
	if flags.DelimiterPositions != nil {
		t.Errorf("delimitPositions = %v, expect to set %v for %q", flags.DelimiterPositions, nil, "spaces")
	}

	expectErr := "delimiter must be one character, \"SPACES\" or JSON array of integers"
	err := SetDelimiter("[a]")
	if err == nil {
		t.Errorf("no error, want error %q for %s", expectErr, "//")
	} else if err.Error() != expectErr {
		t.Errorf("error = %q, want error %q for %s", err.Error(), expectErr, "//")
	}

	expectErr = "delimiter must be one character, \"SPACES\" or JSON array of integers"
	err = SetDelimiter("//")
	if err == nil {
		t.Errorf("no error, want error %q for %s", expectErr, "//")
	} else if err.Error() != expectErr {
		t.Errorf("error = %q, want error %q for %s", err.Error(), expectErr, "//")
	}
}

func TestSetJsonQuery(t *testing.T) {
	flags := GetFlags()

	SetJsonQuery("{}")
	if flags.JsonQuery != "{}" {
		t.Errorf("json-query = %q, expect to set %q", flags.JsonQuery, "{}")
	}
}

func TestSetEncoding(t *testing.T) {
	flags := GetFlags()

	SetEncoding("sjis")
	if flags.Encoding != text.SJIS {
		t.Errorf("encoding = %s, expect to set %s for %s", flags.Encoding, text.SJIS, "sjis")
	}

	expectErr := "encoding must be one of UTF8|SJIS"
	err := SetEncoding("error")
	if err == nil {
		t.Errorf("no error, want error %q for %s", expectErr, "error")
	} else if err.Error() != expectErr {
		t.Errorf("error = %q, want error %q for %s", err.Error(), expectErr, "error")
	}
}

func TestSetNoHeader(t *testing.T) {
	flags := GetFlags()

	SetNoHeader(true)
	if !flags.NoHeader {
		t.Errorf("no-header = %t, expect to set %t", flags.NoHeader, true)
	}
}

func TestSetWithoutNull(t *testing.T) {
	flags := GetFlags()

	SetWithoutNull(true)
	if !flags.WithoutNull {
		t.Errorf("without-null = %t, expect to set %t", flags.WithoutNull, true)
	}
}

func TestSetOut(t *testing.T) {
	flags := GetFlags()

	SetOut("outfile")
	if flags.OutFile != "outfile" {
		t.Errorf("out-file = %s, expect to set %s for %s", flags.OutFile, "outfile", "outfile")
	}

	err := SetOut("")
	if err != nil {
		t.Errorf("unexpected error %q for %q", err.Error(), "")
	}

	expectErr := fmt.Sprintf("file %q already exists", "flags_test.go")
	err = SetOut("flags_test.go")
	if err == nil {
		t.Errorf("no error, want error %q for %q", expectErr, "flags_test.go")
	} else if err.Error() != expectErr {
		t.Errorf("error = %q, want error %q for %q", err.Error(), expectErr, "flags_test.go")
	}
}

func TestSetFormat(t *testing.T) {
	flags := GetFlags()

	SetFormat("")
	if flags.Format != TEXT {
		t.Errorf("format = %s, expect to set %s for empty string", flags.Format, TEXT)
	}

	SetOut("foo.csv")
	SetFormat("")
	if flags.Format != CSV {
		t.Errorf("format = %s, expect to set %s for empty string with file %q", flags.Format, CSV, "foo.csv")
	}

	SetOut("foo.tsv")
	SetFormat("")
	if flags.Format != TSV {
		t.Errorf("format = %s, expect to set %s for empty string with file %q", flags.Format, TSV, "foo.tsv")
	}

	SetOut("foo.txt")
	SetFormat("")
	if flags.Format != FIXED {
		t.Errorf("format = %s, expect to set %s for empty string with file %q", flags.Format, FIXED, "foo.txt")
	}

	SetOut("foo.json")
	SetFormat("")
	if flags.Format != JSON {
		t.Errorf("format = %s, expect to set %s for empty string with file %q", flags.Format, JSON, "foo.json")
	}

	SetOut("foo.md")
	SetFormat("")
	if flags.Format != GFM {
		t.Errorf("format = %s, expect to set %s for empty string with file %q", flags.Format, GFM, "foo.md")
	}

	SetOut("foo.org")
	SetFormat("")
	if flags.Format != ORG {
		t.Errorf("format = %s, expect to set %s for empty string with file %q", flags.Format, ORG, "foo.org")
	}

	SetFormat("csv")
	if flags.Format != CSV {
		t.Errorf("format = %s, expect to set %s for %s", flags.Format, CSV, "csv")
	}

	SetFormat("tsv")
	if flags.Format != TSV {
		t.Errorf("format = %s, expect to set %s for %s", flags.Format, TSV, "tsv")
	}

	SetFormat("fixed")
	if flags.Format != FIXED {
		t.Errorf("format = %s, expect to set %s for %s", flags.Format, FIXED, "fixed")
	}

	SetFormat("json")
	if flags.Format != JSON {
		t.Errorf("format = %s, expect to set %s for %s", flags.Format, JSON, "json")
	}

	SetFormat("jsonh")
	if flags.Format != JSONH {
		t.Errorf("format = %s, expect to set %s for %s", flags.Format, JSONH, "jsonh")
	}

	SetFormat("jsona")
	if flags.Format != JSONA {
		t.Errorf("format = %s, expect to set %s for %s", flags.Format, JSONA, "jsona")
	}

	SetFormat("gfm")
	if flags.Format != GFM {
		t.Errorf("format = %s, expect to set %s for %s", flags.Format, GFM, "gfm")
	}

	SetFormat("org")
	if flags.Format != ORG {
		t.Errorf("format = %s, expect to set %s for %s", flags.Format, ORG, "org")
	}

	SetFormat("text")
	if flags.Format != TEXT {
		t.Errorf("format = %s, expect to set %s for %s", flags.Format, TEXT, "text")
	}

	expectErr := "format must be one of CSV|TSV|FIXED|JSON|JSONH|JSONA|GFM|ORG|TEXT"
	err := SetFormat("error")
	if err == nil {
		t.Errorf("no error, want error %q for %s", expectErr, "error")
	} else if err.Error() != expectErr {
		t.Errorf("error = %q, want error %q for %s", err.Error(), expectErr, "error")
	}
}

func TestSetWriteEncoding(t *testing.T) {
	flags := GetFlags()

	SetWriteEncoding("sjis")
	if flags.Encoding != text.SJIS {
		t.Errorf("encoding = %s, expect to set %s for %s", flags.WriteEncoding, text.SJIS, "sjis")
	}

	expectErr := "encoding must be one of UTF8|SJIS"
	err := SetWriteEncoding("error")
	if err == nil {
		t.Errorf("no error, want error %q for %s", expectErr, "error")
	} else if err.Error() != expectErr {
		t.Errorf("error = %q, want error %q for %s", err.Error(), expectErr, "error")
	}
}

func TestSetWriteDelimiter(t *testing.T) {
	flags := GetFlags()

	SetWriteDelimiter("")
	if flags.WriteDelimiter != ',' {
		t.Errorf("write-delimiter = %q, expect to set %q for %q, format = %s", flags.WriteDelimiter, ',', "", flags.Format)
	}

	SetWriteDelimiter("\\t")
	if flags.WriteDelimiter != '\t' {
		t.Errorf("write-delimiter = %q, expect to set %q for %q", flags.WriteDelimiter, "\t", "\t")
	}

	SetWriteDelimiter("[1, 2, 3]")
	if !reflect.DeepEqual(flags.WriteDelimiterPositions, []int{1, 2, 3}) {
		t.Errorf("writeDelimitPositions = %v, expect to set %v for %q", flags.WriteDelimiterPositions, []int{1, 2, 3}, "[1, 2, 3]")
	}

	expectErr := "write-delimiter must be one character, \"SPACES\" or JSON array of integers"
	err := SetWriteDelimiter("//")
	if err == nil {
		t.Errorf("no error, want error %q for %s", expectErr, "//")
	} else if err.Error() != expectErr {
		t.Errorf("error = %q, want error %q for %s", err.Error(), expectErr, "//")
	}
}

func TestSetWithoutHeader(t *testing.T) {
	flags := GetFlags()

	SetWithoutHeader(true)
	if !flags.NoHeader {
		t.Errorf("without-header = %t, expect to set %t", flags.WithoutHeader, true)
	}
}

func TestSetLineBreak(t *testing.T) {
	flags := GetFlags()

	SetLineBreak("")
	if flags.LineBreak != text.LF {
		t.Errorf("line-break = %s, expect to set %s for %q", flags.LineBreak, text.LF, "")
	}

	SetLineBreak("crlf")
	if flags.LineBreak != text.CRLF {
		t.Errorf("line-break = %s, expect to set %s for %s", flags.LineBreak, text.CRLF, "crlf")
	}

	SetLineBreak("cr")
	if flags.LineBreak != text.CR {
		t.Errorf("line-break = %s, expect to set %s for %s", flags.LineBreak, text.CR, "cr")
	}

	SetLineBreak("lf")
	if flags.LineBreak != text.LF {
		t.Errorf("line-break = %s, expect to set %s for %s", flags.LineBreak, text.LF, "LF")
	}

	expectErr := "line-break must be one of CRLF|LF|CR"
	err := SetLineBreak("error")
	if err == nil {
		t.Errorf("no error, want error %q for %s", expectErr, "error")
	} else if err.Error() != expectErr {
		t.Errorf("error = %q, want error %q for %s", err.Error(), expectErr, "error")
	}
}

func TestSetEncloseAll(t *testing.T) {
	flags := GetFlags()

	SetEncloseAll(true)
	if !flags.EncloseAll {
		t.Errorf("enclose-all = %t, expect to set %t", flags.EncloseAll, true)
	}
}

func TestSetPrettyPrint(t *testing.T) {
	flags := GetFlags()

	SetPrettyPrint(true)
	if !flags.PrettyPrint {
		t.Errorf("pretty-print = %t, expect to set %t", flags.PrettyPrint, true)
	}
}

func TestSetEastAsianEncoding(t *testing.T) {
	flags := GetFlags()

	SetEastAsianEncoding(true)
	if !flags.EastAsianEncoding {
		t.Errorf("east-asian-encoding = %t, expect to set %t", flags.EastAsianEncoding, true)
	}
}

func TestSetCountDiacriticalSign(t *testing.T) {
	flags := GetFlags()

	SetCountDiacriticalSign(true)
	if !flags.CountDiacriticalSign {
		t.Errorf("count-diacritical-sign = %t, expect to set %t", flags.CountDiacriticalSign, true)
	}
}

func TestSetColor(t *testing.T) {
	flags := GetFlags()

	SetColor(true)
	if !flags.Color {
		t.Errorf("color = %t, expect to set %t", flags.Color, true)
	}
	SetColor(false)
}

func TestSetQuiet(t *testing.T) {
	flags := GetFlags()

	SetQuiet(true)
	if !flags.Quiet {
		t.Errorf("silent = %t, expect to set %t", flags.Quiet, true)
	}
}

func TestSetCPU(t *testing.T) {
	flags := GetFlags()

	SetCPU(0)
	expect := 1
	if expect != flags.CPU {
		t.Errorf("cpu = %d, expect to set %d", flags.CPU, 1)
	}

	SetCPU(runtime.NumCPU() + 100)
	if runtime.NumCPU() != flags.CPU {
		t.Errorf("cpu = %d, expect to set %d", flags.CPU, runtime.NumCPU())
	}
}

func TestSetStats(t *testing.T) {
	flags := GetFlags()

	SetStats(true)
	if !flags.Stats {
		t.Errorf("stats = %t, expect to set %t", flags.Stats, true)
	}
}
