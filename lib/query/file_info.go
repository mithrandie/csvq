package query

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/file"
	"github.com/mithrandie/csvq/lib/parser"

	"github.com/mithrandie/go-text"
	"github.com/mithrandie/go-text/fixedlen"
	"github.com/mithrandie/go-text/json"
)

const (
	TableDelimiter          = "DELIMITER"
	TableDelimiterPositions = "DELIMITER_POSITIONS"
	TableFormat             = "FORMAT"
	TableEncoding           = "ENCODING"
	TableLineBreak          = "LINE_BREAK"
	TableHeader             = "HEADER"
	TableEncloseAll         = "ENCLOSE_ALL"
	TableJsonEscape         = "JSON_ESCAPE"
	TablePrettyPrint        = "PRETTY_PRINT"
)

type ViewType int

const (
	ViewTypeFile ViewType = iota
	ViewTypeTemporaryTable
	ViewTypeStdin
)

var FileAttributeList = []string{
	TableDelimiter,
	TableDelimiterPositions,
	TableFormat,
	TableEncoding,
	TableLineBreak,
	TableHeader,
	TableEncloseAll,
	TableJsonEscape,
	TablePrettyPrint,
}

type TableAttributeUnchangedError struct {
	Path    string
	Message string
}

func NewTableAttributeUnchangedError(fpath string) error {
	return &TableAttributeUnchangedError{
		Path:    fpath,
		Message: "table attributes of %s remain unchanged",
	}
}

func (e TableAttributeUnchangedError) Error() string {
	return fmt.Sprintf(e.Message, e.Path)
}

type FileInfo struct {
	Path string

	Format             cmd.Format
	Delimiter          rune
	DelimiterPositions fixedlen.DelimiterPositions
	JsonQuery          string
	Encoding           text.Encoding
	LineBreak          text.LineBreak
	NoHeader           bool
	EncloseAll         bool
	JsonEscape         json.EscapeType
	PrettyPrint        bool

	SingleLine bool

	Handler *file.Handler

	ForUpdate bool
	ViewType  ViewType

	restorePointHeader    Header
	restorePointRecordSet RecordSet
}

func NewFileInfo(
	filename parser.Identifier,
	repository string,
	options cmd.ImportOptions,
	defaultFormat cmd.Format,
) (*FileInfo, error) {
	fpath, format, err := SearchFilePath(filename, repository, options, defaultFormat)
	if err != nil {
		return nil, err
	}

	delimiter := options.Delimiter
	encoding := options.Encoding
	switch format {
	case cmd.TSV:
		delimiter = '\t'
	case cmd.JSON, cmd.JSONL:
		encoding = text.UTF8
	}

	return &FileInfo{
		Path:      fpath,
		Format:    format,
		Delimiter: delimiter,
		Encoding:  encoding,
	}, nil
}

func (f *FileInfo) SetDelimiter(s string) error {
	delimiter, err := cmd.ParseDelimiter(s)
	if err != nil {
		return err
	}

	var format cmd.Format
	if delimiter == '\t' {
		format = cmd.TSV
	} else {
		format = cmd.CSV
	}

	if f.Delimiter == delimiter && f.Format == format {
		return NewTableAttributeUnchangedError(f.Path)
	}

	f.Delimiter = delimiter
	f.Format = format
	return nil
}

func (f *FileInfo) SetDelimiterPositions(s string) error {
	pos, singleLine, err := cmd.ParseDelimiterPositions(s)
	if err != nil {
		return err
	}
	delimiterPositions := fixedlen.DelimiterPositions(pos)
	format := cmd.FIXED

	if reflect.DeepEqual(f.DelimiterPositions, delimiterPositions) &&
		f.SingleLine == singleLine &&
		f.Format == format {
		return NewTableAttributeUnchangedError(f.Path)
	}

	f.Format = format
	f.DelimiterPositions = delimiterPositions
	f.SingleLine = singleLine

	return nil
}

func (f *FileInfo) SetFormat(s string) error {
	format, escapeType, err := cmd.ParseFormat(s, f.JsonEscape)
	if err != nil {
		return err
	}

	if f.Format == format &&
		f.JsonEscape == escapeType {
		return NewTableAttributeUnchangedError(f.Path)
	}

	delimiter := f.Delimiter
	encoding := f.Encoding

	switch format {
	case cmd.TSV:
		delimiter = '\t'
	case cmd.JSON, cmd.JSONL:
		encoding = text.UTF8
	}

	f.Format = format
	f.JsonEscape = escapeType
	f.Delimiter = delimiter
	f.Encoding = encoding
	return nil
}

func (f *FileInfo) SetEncoding(s string) error {
	encoding, err := cmd.ParseEncoding(s)
	if err != nil || encoding == text.AUTO {
		return errors.New("encoding must be one of UTF8|UTF8M|UTF16|UTF16BE|UTF16LE|UTF16BEM|UTF16LEM|SJIS")
	}

	switch f.Format {
	case cmd.JSON, cmd.JSONL:
		if encoding != text.UTF8 {
			return errors.New("json format is supported only UTF8")
		}
	}

	if f.Encoding == encoding {
		return NewTableAttributeUnchangedError(f.Path)
	}

	f.Encoding = encoding
	return nil
}

func (f *FileInfo) SetLineBreak(s string) error {
	lb, err := cmd.ParseLineBreak(s)
	if err != nil {
		return err
	}

	if f.LineBreak == lb {
		return NewTableAttributeUnchangedError(f.Path)
	}

	f.LineBreak = lb
	return nil
}

func (f *FileInfo) SetNoHeader(b bool) error {
	if b == f.NoHeader {
		return NewTableAttributeUnchangedError(f.Path)
	}
	f.NoHeader = b
	return nil
}

func (f *FileInfo) SetEncloseAll(b bool) error {
	if b == f.EncloseAll {
		return NewTableAttributeUnchangedError(f.Path)
	}
	f.EncloseAll = b
	return nil
}

func (f *FileInfo) SetJsonEscape(s string) error {
	escape, err := cmd.ParseJsonEscapeType(s)
	if err != nil {
		return err
	}

	if escape == f.JsonEscape {
		return NewTableAttributeUnchangedError(f.Path)
	}

	f.JsonEscape = escape
	return nil
}

func (f *FileInfo) SetPrettyPrint(b bool) error {
	if b == f.PrettyPrint {
		return NewTableAttributeUnchangedError(f.Path)
	}
	f.PrettyPrint = b
	return nil
}

func (f *FileInfo) IsFile() bool {
	return f.ViewType == ViewTypeFile
}

func (f *FileInfo) IsTemporaryTable() bool {
	return f.ViewType == ViewTypeTemporaryTable
}

func (f *FileInfo) IsStdin() bool {
	return f.ViewType == ViewTypeStdin
}

func (f *FileInfo) ExportOptions(tx *Transaction) cmd.ExportOptions {
	ops := tx.Flags.ExportOptions.Copy()
	ops.Format = f.Format
	ops.Delimiter = f.Delimiter
	ops.DelimiterPositions = f.DelimiterPositions
	ops.SingleLine = f.SingleLine
	ops.Encoding = f.Encoding
	ops.LineBreak = f.LineBreak
	ops.WithoutHeader = f.NoHeader
	ops.EncloseAll = f.EncloseAll
	ops.JsonEscape = f.JsonEscape
	ops.PrettyPrint = f.PrettyPrint
	return ops
}

func SearchFilePath(filename parser.Identifier, repository string, options cmd.ImportOptions, defaultFormat cmd.Format) (string, cmd.Format, error) {
	var fpath string
	var err error

	format := options.Format

	switch format {
	case cmd.CSV, cmd.TSV:
		fpath, err = SearchCSVFilePath(filename, repository)
	case cmd.JSON:
		fpath, err = SearchJsonFilePath(filename, repository)
	case cmd.JSONL:
		fpath, err = SearchJsonlFilePath(filename, repository)
	case cmd.FIXED:
		fpath, err = SearchFixedLengthFilePath(filename, repository)
	case cmd.LTSV:
		fpath, err = SearchLTSVFilePath(filename, repository)
	default: // AutoSelect
		if fpath, err = SearchFilePathFromAllTypes(filename, repository); err == nil {
			switch strings.ToLower(filepath.Ext(fpath)) {
			case cmd.CsvExt:
				format = cmd.CSV
			case cmd.TsvExt:
				format = cmd.TSV
			case cmd.JsonExt:
				format = cmd.JSON
			case cmd.JsonlExt:
				format = cmd.JSONL
			case cmd.LtsvExt:
				format = cmd.LTSV
			default:
				format = defaultFormat
			}
		}
	}

	return fpath, format, err
}

func SearchCSVFilePath(filename parser.Identifier, repository string) (string, error) {
	return SearchFilePathWithExtType(filename, repository, []string{cmd.CsvExt, cmd.TsvExt, cmd.TextExt})
}

func SearchJsonFilePath(filename parser.Identifier, repository string) (string, error) {
	return SearchFilePathWithExtType(filename, repository, []string{cmd.JsonExt})
}

func SearchJsonlFilePath(filename parser.Identifier, repository string) (string, error) {
	return SearchFilePathWithExtType(filename, repository, []string{cmd.JsonlExt})
}

func SearchFixedLengthFilePath(filename parser.Identifier, repository string) (string, error) {
	return SearchFilePathWithExtType(filename, repository, []string{cmd.TextExt})
}

func SearchLTSVFilePath(filename parser.Identifier, repository string) (string, error) {
	return SearchFilePathWithExtType(filename, repository, []string{cmd.LtsvExt, cmd.TextExt})
}

func SearchFilePathFromAllTypes(filename parser.Identifier, repository string) (string, error) {
	return SearchFilePathWithExtType(filename, repository, []string{cmd.CsvExt, cmd.TsvExt, cmd.JsonExt, cmd.JsonlExt, cmd.LtsvExt, cmd.TextExt})
}

func SearchFilePathWithExtType(filename parser.Identifier, repository string, extTypes []string) (string, error) {
	fpath := filename.Literal
	if !filepath.IsAbs(fpath) {
		if len(repository) < 1 {
			repository, _ = os.Getwd()
		}
		fpath = filepath.Join(repository, fpath)
	}

	var info os.FileInfo
	var err error

	if info, err = os.Stat(fpath); err != nil {
		pathes := make([]string, 0, len(extTypes))
		infoList := make([]os.FileInfo, 0, len(extTypes))
		for _, ext := range extTypes {
			if i, err := os.Stat(fpath + ext); err == nil {
				pathes = append(pathes, fpath+ext)
				infoList = append(infoList, i)
			}
		}
		switch {
		case len(pathes) < 1:
			return fpath, NewFileNotExistError(filename)
		case 1 < len(pathes):
			return fpath, NewFileNameAmbiguousError(filename)
		}
		fpath = pathes[0]
		info = infoList[0]
	}

	fpath, err = filepath.Abs(fpath)
	if err != nil {
		return fpath, NewFileNotExistError(filename)
	}

	if info.IsDir() {
		return fpath, NewFileUnableToReadError(filename)
	}

	return fpath, nil
}

func NewFileInfoForCreate(filename parser.Identifier, repository string, delimiter rune, encoding text.Encoding) (*FileInfo, error) {
	fpath, err := CreateFilePath(filename, repository)
	if err != nil {
		return nil, NewIOError(filename, err.Error())
	}

	var format cmd.Format
	switch strings.ToLower(filepath.Ext(fpath)) {
	case cmd.TsvExt:
		delimiter = '\t'
		format = cmd.TSV
	case cmd.JsonExt:
		encoding = text.UTF8
		format = cmd.JSON
	case cmd.JsonlExt:
		encoding = text.UTF8
		format = cmd.JSONL
	case cmd.LtsvExt:
		format = cmd.LTSV
	case cmd.GfmExt:
		format = cmd.GFM
	case cmd.OrgExt:
		format = cmd.ORG
	default:
		format = cmd.CSV
	}

	return &FileInfo{
		Path:      fpath,
		Delimiter: delimiter,
		Format:    format,
		Encoding:  encoding,
	}, nil
}

func CreateFilePath(filename parser.Identifier, repository string) (string, error) {
	fpath := filename.Literal
	if !filepath.IsAbs(fpath) {
		if len(repository) < 1 {
			repository, _ = os.Getwd()
		}
		fpath = filepath.Join(repository, fpath)
	}
	return filepath.Abs(fpath)
}
