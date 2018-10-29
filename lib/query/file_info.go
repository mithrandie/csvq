package query

import (
	"github.com/mithrandie/csvq/lib/text"
	"os"
	"path/filepath"
	"strings"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/parser"
)

type FileInfo struct {
	Path      string
	Delimiter rune

	Format             cmd.Format
	DelimiterPositions text.DelimiterPositions
	JsonQuery          string
	Encoding           cmd.Encoding
	LineBreak          cmd.LineBreak
	NoHeader           bool
	PrettyPrint        bool

	File *os.File

	IsTemporary      bool
	InitialHeader    Header
	InitialRecordSet RecordSet
}

func NewFileInfo(filename parser.Identifier, repository string, delimiter rune, format cmd.Format) (*FileInfo, error) {
	fpath, delimiter, format, err := SearchFilePath(filename, repository, delimiter, format)
	if err != nil {
		return nil, err
	}

	return &FileInfo{
		Path:      fpath,
		Delimiter: delimiter,
		Format:    format,
	}, nil
}

func (f *FileInfo) Equivalent(f2 *FileInfo) bool {
	if f.Path != f2.Path ||
		f.Delimiter != f2.Delimiter ||
		f.Format != f2.Format ||
		(f2.DelimiterPositions != nil && !f.DelimiterPositions.Equal(f2.DelimiterPositions)) ||
		f.JsonQuery != f2.JsonQuery ||
		f.Encoding != f2.Encoding ||
		f.LineBreak != f2.LineBreak ||
		f.NoHeader != f2.NoHeader {
		return false
	}
	return true
}

func SearchFilePath(filename parser.Identifier, repository string, delimiter rune, format cmd.Format) (string, rune, cmd.Format, error) {
	var fpath string
	var err error

	switch format {
	case cmd.CSV:
		fpath, err = SearchCSVFilePath(filename, repository)
	case cmd.TSV:
		fpath, err = SearchTSVFilePath(filename, repository)
		delimiter = '\t'
	case cmd.JSON:
		fpath, err = SearchJsonFilePath(filename, repository)
	default:
		if fpath, err = SearchFilePathFromAllTypes(filename, repository); err == nil {
			switch strings.ToLower(filepath.Ext(fpath)) {
			case cmd.CsvExt:
				format = cmd.CSV
			case cmd.TsvExt:
				delimiter = '\t'
				format = cmd.TSV
			case cmd.JsonExt:
				format = cmd.JSON
			default:
				if format == cmd.AutoSelect {
					format = cmd.GetFlags().ImportFormat
				}
			}
		}
	}

	return fpath, delimiter, format, err
}

func SearchCSVFilePath(filename parser.Identifier, repository string) (string, error) {
	return SearchFilePathWithExtType(filename, repository, []string{cmd.CsvExt})
}

func SearchTSVFilePath(filename parser.Identifier, repository string) (string, error) {
	return SearchFilePathWithExtType(filename, repository, []string{cmd.TsvExt})
}

func SearchJsonFilePath(filename parser.Identifier, repository string) (string, error) {
	return SearchFilePathWithExtType(filename, repository, []string{cmd.JsonExt})
}

func SearchFilePathFromAllTypes(filename parser.Identifier, repository string) (string, error) {
	return SearchFilePathWithExtType(filename, repository, []string{cmd.CsvExt, cmd.TsvExt, cmd.JsonExt})
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

func NewFileInfoForCreate(finename parser.Identifier, repository string, delimiter rune) (*FileInfo, error) {
	fpath := finename.Literal
	if !filepath.IsAbs(fpath) {
		fpath = filepath.Join(repository, fpath)
	}
	fpath, err := filepath.Abs(fpath)
	if err != nil {
		return nil, NewWriteFileError(finename, err.Error())
	}

	var format cmd.Format
	switch strings.ToLower(filepath.Ext(fpath)) {
	case cmd.TsvExt:
		delimiter = '\t'
		format = cmd.TSV
	case cmd.JsonExt:
		format = cmd.JSON
	default:
		format = cmd.CSV
	}

	return &FileInfo{
		Path:      fpath,
		Delimiter: delimiter,
		Format:    format,
	}, nil
}
