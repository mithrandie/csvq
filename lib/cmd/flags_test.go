package cmd

import (
	"path/filepath"
	"runtime"
	"testing"
)

func TestLineBreak_Value(t *testing.T) {
	lb := CRLF
	if lb.Value() != "\r\n" {
		t.Errorf("value = %q, want %q for %s", lb.Value(), "\\r\\n", "CRLF")
	}
}

func TestSetDelimiter(t *testing.T) {
	flags := GetFlags()

	SetDelimiter("")
	if flags.Delimiter != UNDEF {
		t.Errorf("delimiter = %q, expect to set %q for %q", flags.Delimiter, UNDEF, "")
	}

	SetDelimiter("\\t")
	if flags.Delimiter != '\t' {
		t.Errorf("delimiter = %q, expect to set %q for %q", flags.Delimiter, "\t", "\t")
	}

	expectErr := "delimiter must be 1 character"
	err := SetDelimiter("//")
	if err == nil {
		t.Errorf("no error, want error %q for %s", expectErr, "//")
	} else if err.Error() != expectErr {
		t.Errorf("error = %q, want error %q for %s", err.Error(), expectErr, "//")
	}
}

func TestSetEncoding(t *testing.T) {
	flags := GetFlags()

	SetEncoding("sjis")
	if flags.Encoding != SJIS {
		t.Errorf("encoding = %s, expect to set %s for %s", flags.Encoding, SJIS, "sjis")
	}

	expectErr := "encoding must be one of utf8|sjis"
	err := SetEncoding("error")
	if err == nil {
		t.Errorf("no error, want error %q for %s", expectErr, "error")
	} else if err.Error() != expectErr {
		t.Errorf("error = %q, want error %q for %s", err.Error(), expectErr, "error")
	}
}

func TestSetLineBreak(t *testing.T) {
	flags := GetFlags()

	SetLineBreak("")
	if flags.LineBreak != LF {
		t.Errorf("line-break = %s, expect to set %s for %q", flags.LineBreak, LF, "")
	}

	SetLineBreak("crlf")
	if flags.LineBreak != CRLF {
		t.Errorf("line-break = %s, expect to set %s for %s", flags.LineBreak, CRLF, "crlf")
	}

	SetLineBreak("cr")
	if flags.LineBreak != CR {
		t.Errorf("line-break = %s, expect to set %s for %s", flags.LineBreak, CR, "cr")
	}

	SetLineBreak("lf")
	if flags.LineBreak != LF {
		t.Errorf("line-break = %s, expect to set %s for %s", flags.LineBreak, LF, "LF")
	}

	expectErr := "line-break must be one of crlf|lf|cr"
	err := SetLineBreak("error")
	if err == nil {
		t.Errorf("no error, want error %q for %s", expectErr, "error")
	} else if err.Error() != expectErr {
		t.Errorf("error = %q, want error %q for %s", err.Error(), expectErr, "error")
	}
}

func TestSetLocation(t *testing.T) {
	flags := GetFlags()

	s := ""
	SetLocation(s)
	if flags.Location != "Local" {
		t.Errorf("location = %s, expect to set %s for %q", flags.Location, "Local", "")
	}

	s = "Local"
	SetLocation(s)
	if flags.Location != s {
		t.Errorf("location = %s, expect to set %s for %q", flags.Location, s, s)
	}

	s = "America/Los_Angeles"
	SetLocation(s)
	if flags.Location != s {
		t.Errorf("location = %s, expect to set %s for %q", flags.Location, s, s)
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

func TestSetRepository(t *testing.T) {
	flags := GetFlags()

	SetRepository("")
	if flags.Repository != "." {
		t.Errorf("repository = %s, expect to set %s for %q", flags.Repository, ".", "")
	}

	dir := filepath.Join("..", "..", "lib", "cmd")
	SetRepository(dir)
	if flags.Repository != dir {
		t.Errorf("repository = %s, expect to set %s for %s", flags.Repository, dir, dir)
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

func TestSetSource(t *testing.T) {
	flags := GetFlags()

	s := filepath.Join("..", "..", "lib", "cmd", "flags_test.go")
	SetSource(s)
	if flags.Source != s {
		t.Errorf("source = %s, expect to set %s for %s", flags.Source, s, s)
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

func TestSetDatetimeFormat(t *testing.T) {
	flags := GetFlags()

	format := "%Y-%m-%d"
	SetDatetimeFormat(format)
	if flags.DatetimeFormat != format {
		t.Errorf("datetime format = %s, expect to set %s", flags.DatetimeFormat, format)
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

func TestSetWriteEncoding(t *testing.T) {
	flags := GetFlags()

	SetWriteEncoding("sjis")
	if flags.Encoding != SJIS {
		t.Errorf("encoding = %s, expect to set %s for %s", flags.WriteEncoding, SJIS, "sjis")
	}

	expectErr := "encoding must be one of utf8|sjis"
	err := SetWriteEncoding("error")
	if err == nil {
		t.Errorf("no error, want error %q for %s", expectErr, "error")
	} else if err.Error() != expectErr {
		t.Errorf("error = %q, want error %q for %s", err.Error(), expectErr, "error")
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

	expectErr := "file passed in out option already exists"
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

	SetOut("foo.json")
	SetFormat("")
	if flags.Format != JSON {
		t.Errorf("format = %s, expect to set %s for empty string with file %q", flags.Format, JSON, "foo.json")
	}

	SetFormat("csv")
	if flags.Format != CSV {
		t.Errorf("format = %s, expect to set %s for %s", flags.Format, CSV, "csv")
	}

	SetFormat("tsv")
	if flags.Format != TSV {
		t.Errorf("format = %s, expect to set %s for %s", flags.Format, TSV, "tsv")
	}

	SetFormat("json")
	if flags.Format != JSON {
		t.Errorf("format = %s, expect to set %s for %s", flags.Format, JSON, "json")
	}

	SetFormat("text")
	if flags.Format != TEXT {
		t.Errorf("format = %s, expect to set %s for %s", flags.Format, TEXT, "text")
	}

	expectErr := "format must be one of csv|tsv|json|text"
	err := SetFormat("error")
	if err == nil {
		t.Errorf("no error, want error %q for %s", expectErr, "error")
	} else if err.Error() != expectErr {
		t.Errorf("error = %q, want error %q for %s", err.Error(), expectErr, "error")
	}
}

func TestSetWriteDelimiter(t *testing.T) {
	flags := GetFlags()

	flags.Format = CSV
	SetWriteDelimiter("")
	if flags.WriteDelimiter != ',' {
		t.Errorf("write-delimiter = %q, expect to set %q for %q, format = %s", flags.WriteDelimiter, ',', "", flags.Format)
	}

	flags.Format = TSV
	SetWriteDelimiter("")
	if flags.WriteDelimiter != '\t' {
		t.Errorf("write-delimiter = %q, expect to set %q for %q, format = %s", flags.WriteDelimiter, '\t', "", flags.Format)
	}

	flags.Format = CSV
	flags.WriteDelimiter = ','
	SetWriteDelimiter("\\t")
	if flags.WriteDelimiter != '\t' {
		t.Errorf("write-delimiter = %q, expect to set %q for %q", flags.WriteDelimiter, "\t", "\t")
	}

	expectErr := "write-delimiter must be 1 character"
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

func TestSetQuiet(t *testing.T) {
	flags := GetFlags()

	SetQuiet(true)
	if !flags.Quiet {
		t.Errorf("silent = %t, expect to set %t", flags.Quiet, true)
	}
}

func TestSetCPU(t *testing.T) {
	flags := GetFlags()

	SetCPU(3)
	if 3 != flags.CPU {
		t.Errorf("cpu = %d, expect to set %d", flags.CPU, 3)
	}

	SetCPU(0)
	expect := runtime.NumCPU() / 2
	if expect < 1 {
		expect = 1
	}
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

func TestParseEncoding(t *testing.T) {
	e, err := ParseEncoding("")
	if err != nil {
		t.Errorf("unexpected error: %q", err.Error())
	}
	if e != UTF8 {
		t.Errorf("encoding = %s, expect to set %s for %s", e, UTF8, "")
	}

	e, err = ParseEncoding("utf8")
	if err != nil {
		t.Errorf("unexpected error: %q", err.Error())
	}
	if e != UTF8 {
		t.Errorf("encoding = %s, expect to set %s for %s", e, UTF8, "utf8")
	}

	e, err = ParseEncoding("sjis")
	if err != nil {
		t.Errorf("unexpected error: %q", err.Error())
	}
	if e != SJIS {
		t.Errorf("encoding = %s, expect to set %s for %s", e, SJIS, "sjis")
	}

	expectErr := "encoding must be one of utf8|sjis"
	_, err = ParseEncoding("error")
	if err == nil {
		t.Errorf("no error, want error %q for %s", expectErr, "error")
	} else if err.Error() != expectErr {
		t.Errorf("error = %q, want error %q for %s", err.Error(), expectErr, "error")
	}
}
