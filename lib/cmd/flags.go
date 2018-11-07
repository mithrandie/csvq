package cmd

import (
	"errors"
	"fmt"
	"github.com/mithrandie/csvq/lib/color"
	"github.com/mithrandie/csvq/lib/file"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	RepositoryFlag     = "@@REPOSITORY"
	TimezoneFlag       = "@@TIMEZONE"
	DatetimeFormatFlag = "@@DATETIME_FORMAT"
	WaitTimeoutFlag    = "@@WAIT_TIMEOUT"
	DelimiterFlag      = "@@DELIMITER"
	JsonQuery          = "@@JSON_QUERY"
	EncodingFlag       = "@@ENCODING"
	NoHeaderFlag       = "@@NO_HEADER"
	WithoutNullFlag    = "@@WITHOUT_NULL"
	FormatFlag         = "@@FORMAT"
	WriteEncodingFlag  = "@@WRITE_ENCODING"
	WriteDelimiterFlag = "@@WRITE_DELIMITER"
	WithoutHeaderFlag  = "@@WITHOUT_HEADER"
	LineBreakFlag      = "@@LINE_BREAK"
	PrettyPrintFlag    = "@@PRETTY_PRINT"
	ColorFlag          = "@@COLOR"
	QuietFlag          = "@@QUIET"
	CPUFlag            = "@@CPU"
	StatsFlag          = "@@STATS"
)

var FlagList = []string{
	RepositoryFlag,
	TimezoneFlag,
	DatetimeFormatFlag,
	WaitTimeoutFlag,
	DelimiterFlag,
	JsonQuery,
	EncodingFlag,
	NoHeaderFlag,
	WithoutNullFlag,
	FormatFlag,
	WriteEncodingFlag,
	WriteDelimiterFlag,
	WithoutHeaderFlag,
	LineBreakFlag,
	PrettyPrintFlag,
	ColorFlag,
	QuietFlag,
	CPUFlag,
	StatsFlag,
}

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
	AutoSelect Format = -1 + iota
	CSV
	TSV
	FIXED
	JSON
	JSONH
	JSONA
	GFM
	ORG
	TEXT
)

var formatLiterals = map[Format]string{
	CSV:   "CSV",
	TSV:   "TSV",
	FIXED: "FIXED",
	JSON:  "JSON",
	JSONH: "JSONH",
	JSONA: "JSONA",
	GFM:   "GFM",
	ORG:   "ORG",
	TEXT:  "TEXT",
}

func (f Format) String() string {
	return formatLiterals[f]
}

const (
	CsvExt   = ".csv"
	TsvExt   = ".tsv"
	FixedExt = ".txt"
	JsonExt  = ".json"
	GfmExt   = ".md"
	OrgExt   = ".org"
)

type Flags struct {
	// Common Settings
	Repository     string
	Location       string
	DatetimeFormat string
	WaitTimeout    float64

	// For Procedure
	Source string

	// For Import
	Delimiter   rune
	JsonQuery   string
	Encoding    Encoding
	NoHeader    bool
	WithoutNull bool

	// For Export
	OutFile        string
	Format         Format
	WriteEncoding  Encoding
	WriteDelimiter rune
	WithoutHeader  bool
	LineBreak      LineBreak
	PrettyPrint    bool

	// ANSI Color Sequence
	Color bool

	// System Use
	Quiet bool
	CPU   int
	Stats bool

	// Fixed Value
	RetryInterval time.Duration

	// For Fixed-Length Format
	DelimitAutomatically    bool
	DelimiterPositions      []int
	WriteDelimiterPositions []int

	// Use in tests
	Now string
}

var (
	flags    *Flags
	getFlags sync.Once
)

func GetFlags() *Flags {
	getFlags.Do(func() {
		pwd, err := filepath.Abs(".")
		if err != nil {
			pwd = "."
		}

		cpu := runtime.NumCPU() / 2
		if cpu < 1 {
			cpu = 1
		}

		flags = &Flags{
			Repository:              pwd,
			Location:                "Local",
			DatetimeFormat:          "",
			WaitTimeout:             10,
			Source:                  "",
			Delimiter:               ',',
			JsonQuery:               "",
			Encoding:                UTF8,
			NoHeader:                false,
			WithoutNull:             false,
			OutFile:                 "",
			Format:                  TEXT,
			WriteEncoding:           UTF8,
			WriteDelimiter:          ',',
			WithoutHeader:           false,
			LineBreak:               LF,
			PrettyPrint:             false,
			Color:                   false,
			Quiet:                   false,
			CPU:                     cpu,
			Stats:                   false,
			RetryInterval:           10 * time.Millisecond,
			DelimitAutomatically:    false,
			DelimiterPositions:      nil,
			WriteDelimiterPositions: nil,
			Now:                     "",
		}
	})
	return flags
}

func (f *Flags) ImportFormat() Format {
	if 0 < len(f.JsonQuery) {
		return JSON
	}
	if f.DelimitAutomatically || f.DelimiterPositions != nil {
		return FIXED
	}
	if f.Delimiter == '\t' {
		return TSV
	}
	return CSV
}

func SetRepository(s string) error {
	if len(s) < 1 {
		return nil
	}

	path, err := filepath.Abs(s)
	if err != nil {
		path = s
	}

	stat, err := os.Stat(path)
	if err != nil {
		return errors.New("repository does not exist")
	}
	if !stat.IsDir() {
		return errors.New("repository must be a directory path")
	}

	f := GetFlags()
	f.Repository = path
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
	if abs, err := filepath.Abs(s); err == nil {
		s = abs
	}

	f := GetFlags()
	f.Source = s
	return nil
}

func SetDelimiter(s string) error {
	if len(s) < 1 {
		return nil
	}

	f := GetFlags()

	delimiter, delimiterPositions, delimitAutomatically, err := ParseDelimiter(s, f.Delimiter, f.DelimiterPositions, f.DelimitAutomatically)
	if err != nil {
		return err
	}

	f.Delimiter = delimiter
	f.DelimiterPositions = delimiterPositions
	f.DelimitAutomatically = delimitAutomatically
	return nil
}

func SetJsonQuery(s string) {
	f := GetFlags()
	f.JsonQuery = strings.TrimSpace(s)
	return
}

func SetEncoding(s string) error {
	if len(s) < 1 {
		return nil
	}

	encoding, err := ParseEncoding(s)
	if err != nil {
		return err
	}

	f := GetFlags()
	f.Encoding = encoding
	return nil
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

	switch s {
	case "":
		switch strings.ToLower(filepath.Ext(f.OutFile)) {
		case CsvExt:
			fm = CSV
		case TsvExt:
			fm = TSV
		case FixedExt:
			fm = FIXED
		case JsonExt:
			fm = JSON
		case GfmExt:
			fm = GFM
		case OrgExt:
			fm = ORG
		default:
			return nil
		}
	default:
		var err error
		if fm, err = ParseFormat(s); err != nil {
			return err
		}
	}

	f.Format = fm
	return nil
}

func SetWriteEncoding(s string) error {
	if len(s) < 1 {
		return nil
	}

	encoding, err := ParseEncoding(s)
	if err != nil {
		return err
	}

	f := GetFlags()
	f.WriteEncoding = encoding
	return nil
}

func SetWriteDelimiter(s string) error {
	if len(s) < 1 {
		return nil
	}

	f := GetFlags()

	delimiter, delimiterPositions, _, err := ParseDelimiter(s, f.WriteDelimiter, f.WriteDelimiterPositions, false)
	if err != nil {
		return errors.New("write-delimiter must be one character, \"SPACES\" or JSON array of integers")
	}

	f.WriteDelimiter = delimiter
	f.WriteDelimiterPositions = delimiterPositions
	return nil
}

func SetWithoutHeader(b bool) {
	f := GetFlags()
	f.WithoutHeader = b
	return
}

func SetLineBreak(s string) error {
	if len(s) < 1 {
		return nil
	}

	lb, err := ParseLineBreak(s)
	if err != nil {
		return err
	}

	f := GetFlags()
	f.LineBreak = lb
	return nil
}

func SetPrettyPrint(b bool) {
	f := GetFlags()
	f.PrettyPrint = b
	return
}

func SetColor(b bool) {
	f := GetFlags()
	f.Color = b
	color.UseEscapeSequences = b
	return
}

func SetQuiet(b bool) {
	f := GetFlags()
	f.Quiet = b
	return
}

func SetCPU(i int) {
	if i < 1 {
		i = 1
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
