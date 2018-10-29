package cmd

import (
	"encoding/json"
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
	GFM:   "MD",
	ORG:   "ORG",
	TEXT:  "TEXT",
}

func (f Format) String() string {
	return formatLiterals[f]
}

const (
	CsvExt  = ".csv"
	TsvExt  = ".tsv"
	JsonExt = ".json"
)

type Flags struct {
	// Global Options
	Delimiter      rune
	JsonQuery      string
	Encoding       Encoding
	LineBreak      LineBreak
	Location       string
	Repository     string
	Source         string
	DatetimeFormat string
	NoHeader       bool
	WithoutNull    bool
	WaitTimeout    float64

	// For Output
	WriteEncoding  Encoding
	OutFile        string
	Format         Format
	PrettyPrint    bool
	Color          bool
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

	// For Fixed-Length Format
	ImportFormat            Format
	DelimiterPositions      []int
	WriteDelimiterPositions []int
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
			Delimiter:          ',',
			JsonQuery:          "",
			Encoding:           UTF8,
			LineBreak:          LF,
			Location:           "Local",
			Repository:         pwd,
			Source:             "",
			DatetimeFormat:     "",
			NoHeader:           false,
			WithoutNull:        false,
			WaitTimeout:        10,
			WriteEncoding:      UTF8,
			OutFile:            "",
			Format:             TEXT,
			WriteDelimiter:     ',',
			WithoutHeader:      false,
			Quiet:              false,
			CPU:                cpu,
			Stats:              false,
			RetryInterval:      10 * time.Millisecond,
			Now:                "",
			ImportFormat:       CSV,
			DelimiterPositions: nil,
		}
	})
	return flags
}

func SetDelimiter(s string) error {
	importFormat, delimiter, delimiterPositions, err := ParseDelimiter(s)
	if err != nil {
		return err
	}

	f := GetFlags()
	f.Delimiter = delimiter
	f.ImportFormat = importFormat
	f.DelimiterPositions = delimiterPositions
	return nil
}

func SetJsonQuery(s string) {
	f := GetFlags()
	f.JsonQuery = s
	return
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

	lb, err := ParseLineBreak(s)
	if err != nil {
		return err
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

func SetWaitTimeout(f float64) {
	if f < 0 {
		f = 0
	}

	flags := GetFlags()
	flags.WaitTimeout = f
	file.UpdateWaitTimeout(flags.WaitTimeout, flags.RetryInterval)
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
	case "FIXED":
		fm = FIXED
	case "JSON":
		fm = JSON
	case "JSONH":
		fm = JSONH
	case "JSONA":
		fm = JSONA
	case "GFM":
		fm = GFM
	case "ORG":
		fm = ORG
	case "TEXT":
		fm = TEXT
	default:
		return errors.New("format must be one of CSV|TSV|FIXED|JSON|JSONH|JSONA|GFM|ORG|TEXT")
	}

	f.Format = fm
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

func SetWriteDelimiter(s string) error {
	f := GetFlags()

	if f.Format == TSV {
		f.WriteDelimiter = '\t'
		return nil
	}

	var delimiter = ','
	var delimiterPositions []int = nil

	if s == "[]" || 2 < len(s) {
		var positions []int
		err := json.Unmarshal([]byte(s), &positions)
		if err != nil {
			return errors.New("write-delimiter must be one character or JSON array of integers")
		}
		delimiterPositions = positions
	} else if 0 < len(s) {
		s = UnescapeString(s)

		runes := []rune(s)
		if 1 < len(runes) {
			return errors.New("write-delimiter must be one character or JSON array of integers")
		}
		delimiter = runes[0]
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
