package cmd

import (
	"errors"
	"os"
	"path"
	"reflect"
	"strings"
	"sync"
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

func (lb LineBreak) Value() string {
	return reflect.ValueOf(lb).String()
}

type Format int

const (
	TEXT Format = iota
	CSV
	TSV
	JSON
)

const (
	CSV_EXT = ".csv"
	TSV_EXT = ".tsv"
)

type Flags struct {
	// Global Options
	Delimiter   rune
	Encoding    Encoding
	Repository  string
	Source      string
	NoHeader    bool
	WithoutNull bool

	// Write Subcommand Options
	WriteEncoding  Encoding
	LineBreak      LineBreak
	OutFile        string
	Format         Format
	WriteDelimiter rune
	WithoutHeader  bool

	// Use in tests
	Location string
	Now      string
}

var (
	flags    *Flags
	getFlags sync.Once
)

func GetFlags() *Flags {
	getFlags.Do(func() {
		flags = &Flags{
			Delimiter:      UNDEF,
			Encoding:       UTF8,
			Repository:     ".",
			Source:         "",
			NoHeader:       false,
			WithoutNull:    false,
			WriteEncoding:  UTF8,
			LineBreak:      LF,
			OutFile:        "",
			Format:         TEXT,
			WriteDelimiter: ',',
			WithoutHeader:  false,
			Location:       "Local",
			Now:            "",
		}
	})
	return flags
}

func SetDelimiter(s string) error {
	if len(s) < 1 {
		return nil
	}

	s = UnescapeString(s)

	runes := []rune(s)
	if 1 < len(runes) {
		return errors.New("delimiter must be 1 character")
	}

	f := GetFlags()
	f.Delimiter = runes[0]
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

func SetNoHeader(b bool) error {
	f := GetFlags()
	f.NoHeader = b
	return nil
}

func SetWithoutNull(b bool) error {
	f := GetFlags()
	f.WithoutNull = b
	return nil
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

func SetOut(s string) error {
	if len(s) < 1 {
		return nil
	}

	_, err := os.Stat(s)
	if err == nil {
		return errors.New("file passed in out option already exists")
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
		switch strings.ToUpper(path.Ext(f.OutFile)) {
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

	if len(s) < 1 {
		if f.Format == TSV {
			f.WriteDelimiter = '\t'
		}
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

func SetWithoutHeader(b bool) error {
	f := GetFlags()
	f.WithoutHeader = b
	return nil
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
