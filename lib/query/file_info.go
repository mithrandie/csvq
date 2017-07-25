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

	Temporary      bool
	InitialRecords Records
}

func NewFileInfo(filename parser.Identifier, repository string, delimiter rune) (*FileInfo, error) {
	fpath := filename.Literal
	if !filepath.IsAbs(fpath) {
		fpath = filepath.Join(repository, fpath)
	}

	var info os.FileInfo
	var err error

	if info, err = os.Stat(fpath); err != nil {
		if info, err = os.Stat(fpath + cmd.CSV_EXT); err == nil {
			fpath = fpath + cmd.CSV_EXT
		} else if info, err = os.Stat(fpath + cmd.TSV_EXT); err == nil {
			fpath = fpath + cmd.TSV_EXT
		} else {
			return nil, NewFileNotExistError(filename)
		}
	}

	fpath, err = filepath.Abs(fpath)
	if err != nil {
		return nil, NewFileNotExistError(filename)
	}

	if info.IsDir() {
		return nil, NewFileUnableToReadError(filename)
	}

	if delimiter == cmd.UNDEF {
		if strings.EqualFold(filepath.Ext(fpath), cmd.TSV_EXT) {
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
