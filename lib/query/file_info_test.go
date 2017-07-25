package query

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/parser"
)

var fileInfoTests = []struct {
	Name       string
	FilePath   parser.Identifier
	Repository string
	Delimiter  rune
	Result     *FileInfo
	Error      string
}{
	{
		Name:       "CSV",
		FilePath:   parser.Identifier{Literal: "table1"},
		Repository: filepath.Join("..", "..", "testdata", "csv"),
		Delimiter:  cmd.UNDEF,
		Result: &FileInfo{
			Path:      "table1.csv",
			Delimiter: ',',
		},
	},
	{
		Name:       "TSV",
		FilePath:   parser.Identifier{Literal: "table3"},
		Repository: filepath.Join("..", "..", "testdata", "csv"),
		Delimiter:  cmd.UNDEF,
		Result: &FileInfo{
			Path:      "table3.tsv",
			Delimiter: '\t',
		},
	},
	{
		Name:      "Not Exist Error",
		FilePath:  parser.Identifier{Literal: "notexist"},
		Delimiter: cmd.UNDEF,
		Error:     "[L:- C:-] file notexist does not exist",
	},
	{
		Name:      "Directory Error",
		FilePath:  parser.Identifier{Literal: TestDir},
		Delimiter: cmd.UNDEF,
		Error:     fmt.Sprintf("[L:- C:-] file %s is unable to be read", TestDir),
	},
}

func TestNewFileInfo(t *testing.T) {
	for _, v := range fileInfoTests {
		repo := v.Repository
		if 0 < len(repo) {
			dir, _ := os.Getwd()
			repo = filepath.Join(dir, repo)
		}

		fileInfo, err := NewFileInfo(v.FilePath, repo, v.Delimiter)
		if err != nil {
			if len(v.Error) < 1 {
				t.Errorf("%s: unexpected error %q", v.Name, err)
			} else if err.Error() != v.Error {
				t.Errorf("%s: error %q, want error %q", v.Name, err.Error(), v.Error)
			}
			continue
		}
		if 0 < len(v.Error) {
			t.Errorf("%s: no error, want error %q", v.Name, v.Error)
			continue
		}

		abs, _ := filepath.Abs(filepath.Join(repo, v.Result.Path))
		if fileInfo.Path != abs {
			t.Errorf("%s: filepath = %s, want %s", v.Name, filepath.Base(fileInfo.Path), abs)
		}
		if fileInfo.Delimiter != v.Result.Delimiter {
			t.Errorf("%s: delimiter = %q, want %q", v.Name, fileInfo.Delimiter, v.Result.Delimiter)
		}
	}
}
