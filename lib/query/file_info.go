package query

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/parser"
)

type FileInfo struct {
	Path      string
	Delimiter rune
	NoHeader  bool
	Encoding  cmd.Encoding
	LineBreak cmd.LineBreak
	File      *os.File

	IsTemporary      bool
	InitialHeader    Header
	InitialRecordSet RecordSet
}

func NewFileInfo(filename parser.Identifier, repository string, delimiter rune) (*FileInfo, error) {
	fpath, err := searchFilePath(filename, repository, []string{cmd.CsvExt, cmd.TsvExt})
	if err != nil {
		return nil, err
	}

	if delimiter == cmd.UNDEF {
		if strings.EqualFold(filepath.Ext(fpath), cmd.TsvExt) {
			delimiter = '\t'
		} else {
			delimiter = ','
		}
	}

	return &FileInfo{
		Path:      fpath,
		Delimiter: delimiter,
	}, nil
}

func searchFilePath(filename parser.Identifier, repository string, extTypes []string) (string, error) {
	fpath := filename.Literal
	if !filepath.IsAbs(fpath) {
		fpath = filepath.Join(repository, fpath)
	}

	var info os.FileInfo
	var err error

	if info, err = os.Stat(fpath); err != nil {
		found := false
		for _, ext := range extTypes {
			if info, err = os.Stat(fpath + ext); err == nil {
				fpath = fpath + ext
				found = true
				break
			}
		}
		if !found {
			return fpath, NewFileNotExistError(filename)
		}
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

	if delimiter == cmd.UNDEF {
		if strings.EqualFold(filepath.Ext(fpath), cmd.TsvExt) {
			delimiter = '\t'
		} else {
			delimiter = ','
		}
	}

	return &FileInfo{
		Path:      fpath,
		Delimiter: delimiter,
	}, nil
}
