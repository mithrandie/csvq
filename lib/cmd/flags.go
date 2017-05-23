package cmd

import (
	"errors"
	"os"
	"strings"
	"sync"
)

const UNDEF = -1

type Encoding int

const (
	UTF8 Encoding = iota
	SJIS
)

type LineBreak string

const (
	CR   LineBreak = "\r"
	LF   LineBreak = "\n"
	CRLF LineBreak = "\r\n"
)

type Format int

const (
	TEXT Format = iota
	CSV
	TSV
	JSON
	INSERT
)

const (
	CSV_EXT = ".csv"
	TSV_EXT = ".tsv"
)

type Flags struct {
	Delimiter   rune
	Encoding    Encoding
	Repository  string
	NoHeader    bool
	WithoutNull bool

	WriteEncoding Encoding
	LineBreak     LineBreak
	OutFile       string
	Format        Format
}

var (
	flags    *Flags
	getFlags sync.Once
)

func GetFlags() *Flags {
	getFlags.Do(func() {
		flags = &Flags{
			Delimiter:     UNDEF,
			Encoding:      UTF8,
			Repository:    ".",
			NoHeader:      false,
			WithoutNull:   false,
			WriteEncoding: UTF8,
			LineBreak:     LF,
			OutFile:       "",
			Format:        TEXT,
		}
	})
	return flags
}

func SetDelimiter(s string) error {
	if len(s) < 1 {
		return nil
	}

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
	if err == nil || !os.IsNotExist(err) {
		return errors.New("file passed in out option is already exist")
	}

	f := GetFlags()
	f.OutFile = s
	return nil
}

func SetFormat(s string) error {
	if len(s) < 1 {
		return nil
	}

	var fm Format
	switch strings.ToUpper(s) {
	case "CSV":
		fm = CSV
	case "TSV":
		fm = TSV
	case "JSON":
		fm = JSON
	case "INSERT":
		fm = INSERT
	case "TEXT":
		fm = TEXT
	default:
		return errors.New("format must be one of csv|tsv|json|insert|stdout")
	}

	f := GetFlags()
	f.Format = fm
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
