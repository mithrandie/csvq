package cmd

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/mithrandie/go-text/color"

	"github.com/mithrandie/go-text"

	"github.com/mithrandie/csvq/lib/file"
)

const FlagSign = "@@"
const DelimiteAutomatically = "SPACES"

const (
	RepositoryFlag           = "REPOSITORY"
	TimezoneFlag             = "TIMEZONE"
	DatetimeFormatFlag       = "DATETIME_FORMAT"
	WaitTimeoutFlag          = "WAIT_TIMEOUT"
	DelimiterFlag            = "DELIMITER"
	JsonQueryFlag            = "JSON_QUERY"
	EncodingFlag             = "ENCODING"
	NoHeaderFlag             = "NO_HEADER"
	WithoutNullFlag          = "WITHOUT_NULL"
	FormatFlag               = "FORMAT"
	WriteEncodingFlag        = "WRITE_ENCODING"
	WriteDelimiterFlag       = "WRITE_DELIMITER"
	WithoutHeaderFlag        = "WITHOUT_HEADER"
	LineBreakFlag            = "LINE_BREAK"
	EncloseAll               = "ENCLOSE_ALL"
	PrettyPrintFlag          = "PRETTY_PRINT"
	EastAsianEncodingFlag    = "EAST_ASIAN_ENCODING"
	CountDiacriticalSignFlag = "COUNT_DIACRITICAL_SIGN"
	CountFormatCodeFlag      = "COUNT_FORMAT_CODE"
	ColorFlag                = "COLOR"
	QuietFlag                = "QUIET"
	CPUFlag                  = "CPU"
	StatsFlag                = "STATS"
)

var FlagList = []string{
	RepositoryFlag,
	TimezoneFlag,
	DatetimeFormatFlag,
	WaitTimeoutFlag,
	DelimiterFlag,
	JsonQueryFlag,
	EncodingFlag,
	NoHeaderFlag,
	WithoutNullFlag,
	FormatFlag,
	WriteEncodingFlag,
	WriteDelimiterFlag,
	WithoutHeaderFlag,
	LineBreakFlag,
	EncloseAll,
	PrettyPrintFlag,
	EastAsianEncodingFlag,
	CountDiacriticalSignFlag,
	CountFormatCodeFlag,
	ColorFlag,
	QuietFlag,
	CPUFlag,
	StatsFlag,
}

func FlagSymbol(s string) string {
	return FlagSign + s
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

	// For Import
	Delimiter   rune
	JsonQuery   string
	Encoding    text.Encoding
	NoHeader    bool
	WithoutNull bool

	// For Export
	Format         Format
	WriteEncoding  text.Encoding
	WriteDelimiter rune
	WithoutHeader  bool
	LineBreak      text.LineBreak
	EncloseAll     bool
	PrettyPrint    bool

	// For Calculation of String Width
	EastAsianEncoding    bool
	CountDiacriticalSign bool
	CountFormatCode      bool

	// ANSI Color Sequence
	Color bool

	// System Use
	Quiet bool
	CPU   int
	Stats bool

	// For CSV
	// For Fixed-Length Format
	DelimitAutomatically    bool
	DelimiterPositions      []int
	WriteDelimiterPositions []int

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
			Delimiter:               ',',
			JsonQuery:               "",
			Encoding:                text.UTF8,
			NoHeader:                false,
			WithoutNull:             false,
			Format:                  TEXT,
			WriteEncoding:           text.UTF8,
			WriteDelimiter:          ',',
			WithoutHeader:           false,
			LineBreak:               text.LF,
			EncloseAll:              false,
			PrettyPrint:             false,
			EastAsianEncoding:       false,
			CountDiacriticalSign:    false,
			CountFormatCode:         false,
			Color:                   false,
			Quiet:                   false,
			CPU:                     cpu,
			Stats:                   false,
			DelimitAutomatically:    false,
			DelimiterPositions:      nil,
			WriteDelimiterPositions: nil,
			RetryInterval:           10 * time.Millisecond,
			Now:                     "",
		}
	})
	return flags
}

func (f *Flags) SelectImportFormat() Format {
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

func (f *Flags) SetRepository(s string) error {
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

	f.Repository = path
	return nil
}

func (f *Flags) SetLocation(s string) error {
	if len(s) < 1 || strings.EqualFold(s, "Local") {
		s = "Local"
	} else if strings.EqualFold(s, "UTC") {
		s = "UTC"
	}

	location, err := time.LoadLocation(s)
	if err != nil {
		return errors.New("timezone does not exist")
	}

	f.Location = location.String()
	time.Local = location
	return nil
}

func (f *Flags) SetDatetimeFormat(s string) {
	f.DatetimeFormat = s
}

func (f *Flags) SetWaitTimeout(t float64) {
	if t < 0 {
		t = 0
	}

	flags.WaitTimeout = t
	file.UpdateWaitTimeout(flags.WaitTimeout, flags.RetryInterval)
	return
}

func (f *Flags) SetDelimiter(s string) error {
	if len(s) < 1 {
		return nil
	}

	delimiter, delimiterPositions, delimitAutomatically, err := ParseDelimiter(s, f.Delimiter, f.DelimiterPositions, f.DelimitAutomatically)
	if err != nil {
		return err
	}

	f.Delimiter = delimiter
	f.DelimiterPositions = delimiterPositions
	f.DelimitAutomatically = delimitAutomatically
	return nil
}

func (f *Flags) SetJsonQuery(s string) {
	f.JsonQuery = strings.TrimSpace(s)
}

func (f *Flags) SetEncoding(s string) error {
	if len(s) < 1 {
		return nil
	}

	encoding, err := ParseEncoding(s)
	if err != nil {
		return err
	}

	f.Encoding = encoding
	return nil
}

func (f *Flags) SetNoHeader(b bool) {
	f.NoHeader = b
}

func (f *Flags) SetWithoutNull(b bool) {
	f.WithoutNull = b
}

func (f *Flags) SetFormat(s string, outfile string) error {
	var fm Format

	switch s {
	case "":
		switch strings.ToLower(filepath.Ext(outfile)) {
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

func (f *Flags) SetWriteEncoding(s string) error {
	if len(s) < 1 {
		return nil
	}

	encoding, err := ParseEncoding(s)
	if err != nil {
		return err
	}

	f.WriteEncoding = encoding
	return nil
}

func (f *Flags) SetWriteDelimiter(s string) error {
	if len(s) < 1 {
		return nil
	}

	delimiter, delimiterPositions, _, err := ParseDelimiter(s, f.WriteDelimiter, f.WriteDelimiterPositions, false)
	if err != nil {
		return errors.New("write-delimiter must be one character, \"SPACES\" or JSON array of integers")
	}

	f.WriteDelimiter = delimiter
	f.WriteDelimiterPositions = delimiterPositions
	return nil
}

func (f *Flags) SetWithoutHeader(b bool) {
	f.WithoutHeader = b
}

func (f *Flags) SetLineBreak(s string) error {
	if len(s) < 1 {
		return nil
	}

	lb, err := ParseLineBreak(s)
	if err != nil {
		return err
	}

	f.LineBreak = lb
	return nil
}

func (f *Flags) SetPrettyPrint(b bool) {
	f.PrettyPrint = b
}

func (f *Flags) SetEncloseAll(b bool) {
	f.EncloseAll = b
}

func (f *Flags) SetColor(b bool) {
	f.Color = b
	color.UseEffect = b
}

func (f *Flags) SetEastAsianEncoding(b bool) {
	f.EastAsianEncoding = b
}

func (f *Flags) SetCountDiacriticalSign(b bool) {
	f.CountDiacriticalSign = b
}

func (f *Flags) SetCountFormatCode(b bool) {
	f.CountFormatCode = b
}

func (f *Flags) SetQuiet(b bool) {
	f.Quiet = b
}

func (f *Flags) SetCPU(i int) {
	if i < 1 {
		i = 1
	}

	if runtime.NumCPU() < i {
		i = runtime.NumCPU()
	}

	f.CPU = i
}

func (f *Flags) SetStats(b bool) {
	f.Stats = b
}
