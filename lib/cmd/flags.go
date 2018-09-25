package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/mithrandie/csvq/lib/file"
)

const UNDEF = -1

type Encoding string

const (
	UTF8 Encoding = "UTF8"
	SJIS Encoding = "SJIS"
)

func (e Encoding) String() string {
	return reflect.ValueOf(e).String()
}

type LineBreak string

const (
	CR   LineBreak = "\r"
	LF   LineBreak = "\n"
	CRLF LineBreak = "\r\n"
)

var lineBreakLiterals = map[LineBreak]string{
	CR:   "CR",
	LF:   "LF",
	CRLF: "CRLF",
}

func (lb LineBreak) Value() string {
	return reflect.ValueOf(lb).String()
}

func (lb LineBreak) String() string {
	return lineBreakLiterals[lb]
}

type Format int

const (
	TEXT Format = iota
	CSV
	TSV
	JSON
)

var formatLiterals = map[Format]string{
	TEXT: "TEXT",
	CSV:  "CSV",
	TSV:  "TSV",
	JSON: "JSON",
}

func (f Format) String() string {
	return formatLiterals[f]
}

const (
	CsvExt = ".csv"
	TsvExt = ".tsv"
)

type Flags struct {
	// Global Options
	Delimiter      rune
	Encoding       Encoding
	LineBreak      LineBreak
	Location       string
	Repository     string
	Source         string
	DatetimeFormat string
	WaitTimeout    float64
	NoHeader       bool
	WithoutNull    bool

	// For Output
	WriteEncoding  Encoding
	OutFile        string
	Format         Format
	WriteDelimiter rune
	WithoutHeader  bool

	// System Use
	Quiet bool
	CPU   int
	Stats bool

	// Fixed Value
	RetryInterval time.Duration

	// Use in tests
	Now string
}

var (
	flags    *Flags
	getFlags sync.Once
)

func GetFlags() *Flags {
	pwd, err := filepath.Abs(".")
	if err != nil {
		pwd = "."
	}

	cpu := runtime.NumCPU() / 2
	if cpu < 1 {
		cpu = 1
	}

	getFlags.Do(func() {
		flags = &Flags{
			Delimiter:      UNDEF,
			Encoding:       UTF8,
			LineBreak:      LF,
			Location:       "Local",
			Repository:     pwd,
			Source:         "",
			DatetimeFormat: "",
			WaitTimeout:    10,
			NoHeader:       false,
			WithoutNull:    false,
			WriteEncoding:  UTF8,
			OutFile:        "",
			Format:         TEXT,
			WriteDelimiter: ',',
			WithoutHeader:  false,
			Quiet:          false,
			CPU:            cpu,
			Stats:          false,
			RetryInterval:  10 * time.Millisecond,
			Now:            "",
		}
	})
	return flags
}

func SetDelimiter(s string) error {
	var delimiter rune

	if len(s) < 1 {
		delimiter = UNDEF
	} else {
		s = UnescapeString(s)

		runes := []rune(s)
		if 1 < len(runes) {
			return errors.New("delimiter must be 1 character")
		}
		delimiter = runes[0]
	}

	f := GetFlags()
	f.Delimiter = delimiter
	return nil
}

func SetEncoding(s string) error {
	encoding, err := ParseEncoding(s)
	if err != nil {
		return err
	}

	f := GetFlags()
	f.Encoding = encoding
	return nil
}

func SetLineBreak(s string) error {
	if len(s) < 1 {
		return nil
	}

	var lb LineBreak
	switch strings.ToUpper(s) {
	case "CRLF":
		lb = CRLF
	case "CR":
		lb = CR
	case "LF":
		lb = LF
	default:
		return errors.New("line-break must be one of crlf|lf|cr")
	}

	f := GetFlags()
	f.LineBreak = lb
	return nil
}

func SetLocation(s string) error {
	if len(s) < 1 || strings.EqualFold(s, "Local") {
		s = "Local"
	} else if strings.EqualFold(s, "UTC") {
		s = "UTC"
	}

	location, err := time.LoadLocation(s)
	if err != nil {
		return errors.New("timezone does not exist")
	}

	f := GetFlags()
	f.Location = s
	time.Local = location
	return nil
}

func SetRepository(s string) error {
	if len(s) < 1 {
		return nil
	}

	stat, err := os.Stat(s)
	if err != nil {
		return errors.New("repository does not exist")
	}
	if !stat.IsDir() {
		return errors.New("repository must be a directory path")
	}

	f := GetFlags()
	f.Repository = s
	return nil
}

func SetSource(s string) error {
	if len(s) < 1 {
		return nil
	}

	stat, err := os.Stat(s)
	if err != nil {
		return errors.New("source file does not exist")
	}
	if stat.IsDir() {
		return errors.New("source file must be a readable file")
	}

	f := GetFlags()
	f.Source = s
	return nil
}

func SetDatetimeFormat(s string) {
	f := GetFlags()
	f.DatetimeFormat = s
	return
}

func SetWaitTimeout(f float64) {
	if f < 0 {
		f = 0
	}

	flags := GetFlags()
	flags.WaitTimeout = f
	file.UpdateWaitTimeout(flags.WaitTimeout, flags.RetryInterval)
	return
}

func SetNoHeader(b bool) {
	f := GetFlags()
	f.NoHeader = b
	return
}

func SetWithoutNull(b bool) {
	f := GetFlags()
	f.WithoutNull = b
	return
}

func SetWriteEncoding(s string) error {
	encoding, err := ParseEncoding(s)
	if err != nil {
		return err
	}

	f := GetFlags()
	f.WriteEncoding = encoding
	return nil
}

func SetOut(s string) error {
	if 0 < len(s) {
		_, err := os.Stat(s)
		if err == nil {
			return errors.New(fmt.Sprintf("file %q already exists", s))
		}
	}

	f := GetFlags()
	f.OutFile = s
	return nil
}

func SetFormat(s string) error {
	var fm Format
	f := GetFlags()

	switch strings.ToUpper(s) {
	case "":
		switch strings.ToUpper(filepath.Ext(f.OutFile)) {
		case ".CSV":
			fm = CSV
		case ".TSV":
			fm = TSV
		case ".JSON":
			fm = JSON
		default:
			return nil
		}
	case "CSV":
		fm = CSV
	case "TSV":
		fm = TSV
	case "JSON":
		fm = JSON
	case "TEXT":
		fm = TEXT
	default:
		return errors.New("format must be one of csv|tsv|json|text")
	}

	f.Format = fm
	return nil
}

func SetWriteDelimiter(s string) error {
	f := GetFlags()

	if f.Format == TSV {
		f.WriteDelimiter = '\t'
		return nil
	}

	if len(s) < 1 {
		f.WriteDelimiter = ','
		return nil
	}

	s = UnescapeString(s)

	runes := []rune(s)
	if 1 < len(runes) {
		return errors.New("write-delimiter must be 1 character")
	}

	f.WriteDelimiter = runes[0]
	return nil
}

func SetWithoutHeader(b bool) {
	f := GetFlags()
	f.WithoutHeader = b
	return
}

func ParseEncoding(s string) (Encoding, error) {
	if len(s) < 1 {
		return UTF8, nil
	}

	var encoding Encoding
	switch strings.ToUpper(s) {
	case "UTF8":
		encoding = UTF8
	case "SJIS":
		encoding = SJIS
	default:
		return UTF8, errors.New("encoding must be one of utf8|sjis")
	}
	return encoding, nil
}

func SetQuiet(b bool) {
	f := GetFlags()
	f.Quiet = b
	return
}

func SetCPU(i int) {
	if i <= 0 {
		return
	}

	if runtime.NumCPU() < i {
		i = runtime.NumCPU()
	}

	f := GetFlags()
	f.CPU = i
	return
}

func SetStats(b bool) {
	f := GetFlags()
	f.Stats = b
	return
}
