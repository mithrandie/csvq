package query

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/mithrandie/csvq/lib/file"

	"github.com/mithrandie/go-text"

	"github.com/mithrandie/go-text/fixedlen"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/parser"
)

const (
	TableDelimiter   = "DELIMITER"
	TableFormat      = "FORMAT"
	TableEncoding    = "ENCODING"
	TableLineBreak   = "LINE_BREAK"
	TableHeader      = "HEADER"
	TableEncloseAll  = "ENCLOSE_ALL"
	TablePrettyPring = "PRETTY_PRINT"
)

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
	delimiter, delimiterPositions, auto, err := cmd.ParseDelimiter(
		s,
		f.Delimiter,
		f.DelimiterPositions,
		f.Format == cmd.FIXED && f.DelimiterPositions == nil,
	)
	if err != nil {
		return err
	}

	f.Delimiter = delimiter
	f.DelimiterPositions = delimiterPositions
	if auto || f.DelimiterPositions != nil {
		f.Format = cmd.FIXED
	} else if delimiter == '\t' {
		f.Format = cmd.TSV
	} else {
		f.Format = cmd.CSV
	}

	return nil
}

func (f *FileInfo) SetFormat(s string) error {
	format, err := cmd.ParseFormat(s)
	if err != nil {
		return err
	}

	switch format {
	case cmd.TSV:
		f.Delimiter = '\t'
	case cmd.JSON, cmd.JSONH, cmd.JSONA:
		f.Encoding = text.UTF8
	}

	f.Format = format
	return nil
}

func (f *FileInfo) SetEncoding(s string) error {
	encoding, err := cmd.ParseEncoding(s)
	if err != nil {
		return err
	}

	switch f.Format {
	case cmd.JSON, cmd.JSONH, cmd.JSONA:
		if encoding != text.UTF8 {
			return errors.New("json format is supported only UTF8")
		}
	}

	f.Encoding = encoding
	return nil
}

func (f *FileInfo) SetLineBreak(s string) error {
	lb, err := cmd.ParseLineBreak(s)
	if err != nil {
		return err
	}

	f.LineBreak = lb
	return nil
}

func (f *FileInfo) SetNoHeader(b bool) {
	f.NoHeader = b
}

func (f *FileInfo) SetEncloseAll(b bool) {
	f.EncloseAll = b
}

func (f *FileInfo) SetPrettyPrint(b bool) {
	f.PrettyPrint = b
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
		fpath = filepath.Join(repository, fpath)
	}
	return filepath.Abs(fpath)
}
