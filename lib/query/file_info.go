package query

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/mithrandie/go-text"
	"github.com/mithrandie/go-text/fixedlen"
	"github.com/mithrandie/go-text/json"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/file"
	"github.com/mithrandie/csvq/lib/parser"
)

const (
	TableDelimiter   = "DELIMITER"
	TableFormat      = "FORMAT"
	TableEncoding    = "ENCODING"
	TableLineBreak   = "LINE_BREAK"
	TableHeader      = "HEADER"
	TableEncloseAll  = "ENCLOSE_ALL"
	TableJsonEscape  = "JSON_ESCAPE"
	TablePrettyPring = "PRETTY_PRINT"
)

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
	Path      string
	Delimiter rune

	Format             cmd.Format
	DelimiterPositions fixedlen.DelimiterPositions
	JsonQuery          string
	Encoding           text.Encoding
	LineBreak          text.LineBreak
	NoHeader           bool
	EncloseAll         bool
	JsonEscape         json.EscapeType
	PrettyPrint        bool

	Handler *file.Handler

	IsTemporary      bool
	InitialHeader    Header
	InitialRecordSet RecordSet
}

func NewFileInfo(
	filename parser.Identifier,
	repository string,
	format cmd.Format,
	delimiter rune,
	encoding text.Encoding,
) (*FileInfo, error) {
	fpath, format, err := SearchFilePath(filename, repository, format)
	if err != nil {
		return nil, err
	}

	switch format {
	case cmd.TSV:
		delimiter = '\t'
	case cmd.JSON:
		encoding = text.UTF8
	}

	return &FileInfo{
		Path:      fpath,
		Format:    format,
		Delimiter: delimiter,
		Encoding:  encoding,
	}, nil
}

func (f *FileInfo) Equivalent(f2 *FileInfo) bool {
	if f.Path != f2.Path ||
		f.Format != f2.Format ||
		f.Delimiter != f2.Delimiter ||
		(f2.DelimiterPositions != nil && !f.DelimiterPositions.Equal(f2.DelimiterPositions)) ||
		f.JsonQuery != f2.JsonQuery ||
		f.Encoding != f2.Encoding ||
		f.NoHeader != f2.NoHeader {
		return false
	}
	return true
}

func (f *FileInfo) SetDelimiter(s string) error {
	delimiter, dp, auto, err := cmd.ParseDelimiter(
		s,
		f.Delimiter,
		f.DelimiterPositions,
		f.Format == cmd.FIXED && f.DelimiterPositions == nil,
	)
	if err != nil {
		return err
	}
	delimiterPositions := fixedlen.DelimiterPositions(dp)

	var format cmd.Format
	if auto || delimiterPositions != nil {
		format = cmd.FIXED
	} else if delimiter == '\t' {
		format = cmd.TSV
	} else {
		format = cmd.CSV
	}

	if f.Delimiter == delimiter &&
		reflect.DeepEqual(f.DelimiterPositions, delimiterPositions) &&
		f.Format == format {
		return NewTableAttributeUnchangedError(f.Path)
	}

	f.Delimiter = delimiter
	f.DelimiterPositions = delimiterPositions
	f.Format = format

	return nil
}

func (f *FileInfo) SetFormat(s string) error {
	format, escapeType, err := cmd.ParseFormat(s, f.JsonEscape)
	if err != nil {
		return err
	}

	delimiter := f.Delimiter
	encoding := f.Encoding

	switch format {
	case cmd.TSV:
		delimiter = '\t'
	case cmd.JSON:
		encoding = text.UTF8
	}

	if f.Delimiter == delimiter &&
		f.Encoding == encoding &&
		f.Format == format &&
		f.JsonEscape == escapeType {
		return NewTableAttributeUnchangedError(f.Path)
	}

	f.Format = format
	f.JsonEscape = escapeType
	f.Delimiter = delimiter
	f.Encoding = encoding
	return nil
}

func (f *FileInfo) SetEncoding(s string) error {
	encoding, err := cmd.ParseEncoding(s)
	if err != nil {
		return err
	}

	switch f.Format {
	case cmd.JSON:
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

func (f *FileInfo) Close() error {
	if f.Handler == nil {
		return nil
	}
	return f.Handler.Close()
}

func (f *FileInfo) CloseWithErrors() error {
	if f.Handler == nil {
		return nil
	}
	return f.Handler.CloseWithErrors()
}

func (f *FileInfo) Commit() error {
	if f.Handler == nil {
		return nil
	}
	return f.Handler.Commit()
}

func SearchFilePath(filename parser.Identifier, repository string, format cmd.Format) (string, cmd.Format, error) {
	var fpath string
	var err error

	switch format {
	case cmd.CSV, cmd.TSV:
		fpath, err = SearchCSVFilePath(filename, repository)
	case cmd.JSON:
		fpath, err = SearchJsonFilePath(filename, repository)
	case cmd.FIXED:
		fpath, err = SearchFixedLengthFilePath(filename, repository)
	default: // AutoSelect
		if fpath, err = SearchFilePathFromAllTypes(filename, repository); err == nil {
			switch strings.ToLower(filepath.Ext(fpath)) {
			case cmd.CsvExt:
				format = cmd.CSV
			case cmd.TsvExt:
				format = cmd.TSV
			case cmd.FixedExt:
				format = cmd.FIXED
			case cmd.JsonExt:
				format = cmd.JSON
			default:
				format = cmd.GetFlags().SelectImportFormat()
			}
		}
	}

	return fpath, format, err
}

func SearchCSVFilePath(filename parser.Identifier, repository string) (string, error) {
	return SearchFilePathWithExtType(filename, repository, []string{cmd.CsvExt, cmd.TsvExt})
}

func SearchJsonFilePath(filename parser.Identifier, repository string) (string, error) {
	return SearchFilePathWithExtType(filename, repository, []string{cmd.JsonExt})
}

func SearchFixedLengthFilePath(filename parser.Identifier, repository string) (string, error) {
	return SearchFilePathWithExtType(filename, repository, []string{cmd.FixedExt})
}

func SearchFilePathFromAllTypes(filename parser.Identifier, repository string) (string, error) {
	return SearchFilePathWithExtType(filename, repository, []string{cmd.CsvExt, cmd.TsvExt, cmd.JsonExt, cmd.FixedExt})
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
		return nil, NewWriteFileError(filename, err.Error())
	}

	var format cmd.Format
	switch strings.ToLower(filepath.Ext(fpath)) {
	case cmd.TsvExt:
		delimiter = '\t'
		format = cmd.TSV
	case cmd.FixedExt:
		format = cmd.FIXED
	case cmd.JsonExt:
		encoding = text.UTF8
		format = cmd.JSON
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
