package query

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/mithrandie/csvq/lib/option"
	"github.com/mithrandie/csvq/lib/parser"

	"github.com/mithrandie/go-text"
)

var fileInfoTests = []struct {
	Name       string
	FilePath   parser.Identifier
	Repository string
	Format     option.Format
	Delimiter  rune
	Encoding   text.Encoding
	Result     *FileInfo
	Error      string
}{
	{
		Name:       "CSV",
		FilePath:   parser.Identifier{Literal: "table1"},
		Repository: TestDir,
		Format:     option.CSV,
		Delimiter:  ',',
		Encoding:   text.UTF8,
		Result: &FileInfo{
			Path:      "table1.csv",
			Delimiter: ',',
			Format:    option.CSV,
			Encoding:  text.UTF8,
		},
	},
	{
		Name:       "CSV with AutoSelect",
		FilePath:   parser.Identifier{Literal: "table1"},
		Repository: TestDir,
		Format:     option.AutoSelect,
		Delimiter:  ',',
		Encoding:   text.UTF8,
		Result: &FileInfo{
			Path:      "table1.csv",
			Delimiter: ',',
			Format:    option.CSV,
			Encoding:  text.UTF8,
		},
	},
	{
		Name:       "TSV",
		FilePath:   parser.Identifier{Literal: "table3"},
		Repository: TestDir,
		Format:     option.TSV,
		Delimiter:  ',',
		Encoding:   text.UTF8,
		Result: &FileInfo{
			Path:      "table3.tsv",
			Delimiter: '\t',
			Format:    option.TSV,
			Encoding:  text.UTF8,
		},
	},
	{
		Name:       "TSV with AutoSelect",
		FilePath:   parser.Identifier{Literal: "table3"},
		Repository: TestDir,
		Format:     option.AutoSelect,
		Delimiter:  ',',
		Encoding:   text.UTF8,
		Result: &FileInfo{
			Path:      "table3.tsv",
			Delimiter: '\t',
			Format:    option.TSV,
			Encoding:  text.UTF8,
		},
	},
	{
		Name:       "JSON",
		FilePath:   parser.Identifier{Literal: "table"},
		Repository: TestDir,
		Format:     option.JSON,
		Delimiter:  ',',
		Encoding:   text.UTF16,
		Result: &FileInfo{
			Path:      "table.json",
			Delimiter: ',',
			Format:    option.JSON,
			Encoding:  text.UTF8,
		},
	},
	{
		Name:       "JSON with AutoSelect",
		FilePath:   parser.Identifier{Literal: "table"},
		Repository: TestDir,
		Format:     option.AutoSelect,
		Delimiter:  ',',
		Encoding:   text.UTF16,
		Result: &FileInfo{
			Path:      "table.json",
			Delimiter: ',',
			Format:    option.JSON,
			Encoding:  text.UTF8,
		},
	},
	{
		Name:       "JSONL",
		FilePath:   parser.Identifier{Literal: "table7"},
		Repository: TestDir,
		Format:     option.JSONL,
		Delimiter:  ',',
		Encoding:   text.UTF16,
		Result: &FileInfo{
			Path:      "table7.jsonl",
			Delimiter: ',',
			Format:    option.JSONL,
			Encoding:  text.UTF8,
		},
	},
	{
		Name:       "JSONL with AutoSelect",
		FilePath:   parser.Identifier{Literal: "table7"},
		Repository: TestDir,
		Format:     option.AutoSelect,
		Delimiter:  ',',
		Encoding:   text.UTF16,
		Result: &FileInfo{
			Path:      "table7.jsonl",
			Delimiter: ',',
			Format:    option.JSONL,
			Encoding:  text.UTF8,
		},
	},
	{
		Name:       "LTSV",
		FilePath:   parser.Identifier{Literal: "table6"},
		Repository: TestDir,
		Format:     option.LTSV,
		Delimiter:  ',',
		Encoding:   text.UTF8,
		Result: &FileInfo{
			Path:      "table6.ltsv",
			Delimiter: ',',
			Format:    option.LTSV,
			Encoding:  text.UTF8,
		},
	},
	{
		Name:       "LTSV with AutoSelect",
		FilePath:   parser.Identifier{Literal: "table6"},
		Repository: TestDir,
		Format:     option.AutoSelect,
		Delimiter:  ',',
		Encoding:   text.UTF8,
		Result: &FileInfo{
			Path:      "table6.ltsv",
			Delimiter: ',',
			Format:    option.LTSV,
			Encoding:  text.UTF8,
		},
	},
	{
		Name:       "Fixed-Length",
		FilePath:   parser.Identifier{Literal: "fixed_length.txt"},
		Repository: TestDir,
		Format:     option.FIXED,
		Delimiter:  ',',
		Encoding:   text.UTF8,
		Result: &FileInfo{
			Path:      "fixed_length.txt",
			Delimiter: ',',
			Format:    option.FIXED,
			Encoding:  text.UTF8,
		},
	},
	{
		Name:       "Import Format",
		FilePath:   parser.Identifier{Literal: "autoselect"},
		Repository: TestDir,
		Format:     option.AutoSelect,
		Delimiter:  ',',
		Encoding:   text.UTF8,
		Result: &FileInfo{
			Path:      "autoselect",
			Delimiter: ',',
			Format:    option.CSV,
			Encoding:  text.UTF8,
		},
	},
	{
		Name:       "Not Exist Error",
		FilePath:   parser.Identifier{Literal: "notexist"},
		Repository: TestDir,
		Format:     option.CSV,
		Delimiter:  ',',
		Encoding:   text.UTF8,
		Error:      "file notexist does not exist",
	},
	{
		Name:       "File Read Error",
		FilePath:   parser.Identifier{Literal: TestDir},
		Repository: TestDir,
		Format:     option.CSV,
		Delimiter:  ',',
		Encoding:   text.UTF8,
		Error:      fmt.Sprintf("file %s is unable to be read", TestDir),
	},
	{
		Name:       "Filenames Ambiguous",
		FilePath:   parser.Identifier{Literal: "dup_name"},
		Repository: TestDir,
		Format:     option.AutoSelect,
		Delimiter:  ',',
		Encoding:   text.UTF8,
		Error:      fmt.Sprintf("filename dup_name is ambiguous"),
	},
}

func TestNewFileInfo(t *testing.T) {
	options := TestTx.Flags.ImportOptions.Copy()

	for _, v := range fileInfoTests {
		options.Format = v.Format
		options.Delimiter = v.Delimiter
		options.Encoding = v.Encoding
		fileInfo, err := NewFileInfo(v.FilePath, v.Repository, options, TestTx.Flags.ImportOptions.Format)
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

		abs, _ := filepath.Abs(filepath.Join(v.Repository, v.Result.Path))
		if fileInfo.Path != abs {
			t.Errorf("%s: FileInfo.Path = %s, want %s", v.Name, filepath.Base(fileInfo.Path), abs)
		}
		if fileInfo.Delimiter != v.Result.Delimiter {
			t.Errorf("%s: FileInfo.Delimiter = %q, want %q", v.Name, fileInfo.Delimiter, v.Result.Delimiter)
		}
		if fileInfo.Format != v.Result.Format {
			t.Errorf("%s: FileInfo.Format = %s, want %s", v.Name, fileInfo.Format, v.Result.Format)
		}
	}
}

var fileInfoForCreateTests = []struct {
	Name       string
	FilePath   parser.Identifier
	Repository string
	Delimiter  rune
	Encoding   text.Encoding
	Result     *FileInfo
	Error      string
}{
	{
		Name:      "CSV",
		FilePath:  parser.Identifier{Literal: "table1.csv"},
		Delimiter: ',',
		Encoding:  text.UTF8,
		Result: &FileInfo{
			Path:      "table1.csv",
			Delimiter: ',',
			Format:    option.CSV,
			Encoding:  text.UTF8,
		},
	},
	{
		Name:      "TSV",
		FilePath:  parser.Identifier{Literal: "table1.tsv"},
		Delimiter: ',',
		Encoding:  text.UTF8,
		Result: &FileInfo{
			Path:      "table1.tsv",
			Delimiter: '\t',
			Format:    option.TSV,
			Encoding:  text.UTF8,
		},
	},
	{
		Name:      "JSON",
		FilePath:  parser.Identifier{Literal: "table1.json"},
		Delimiter: ',',
		Encoding:  text.SJIS,
		Result: &FileInfo{
			Path:      "table1.json",
			Delimiter: ',',
			Format:    option.JSON,
			Encoding:  text.UTF8,
		},
	},
	{
		Name:      "JSONL",
		FilePath:  parser.Identifier{Literal: "table1.jsonl"},
		Delimiter: ',',
		Encoding:  text.SJIS,
		Result: &FileInfo{
			Path:      "table1.jsonl",
			Delimiter: ',',
			Format:    option.JSONL,
			Encoding:  text.UTF8,
		},
	},
	{
		Name:      "LTSV",
		FilePath:  parser.Identifier{Literal: "table1.ltsv"},
		Delimiter: ',',
		Encoding:  text.UTF8,
		Result: &FileInfo{
			Path:      "table1.ltsv",
			Delimiter: ',',
			Format:    option.LTSV,
			Encoding:  text.UTF8,
		},
	},
	{
		Name:      "GFM",
		FilePath:  parser.Identifier{Literal: "table1.md"},
		Delimiter: ',',
		Encoding:  text.UTF8,
		Result: &FileInfo{
			Path:      "table1.md",
			Delimiter: ',',
			Format:    option.GFM,
			Encoding:  text.UTF8,
		},
	},
	{
		Name:      "ORG",
		FilePath:  parser.Identifier{Literal: "table1.org"},
		Delimiter: ',',
		Encoding:  text.UTF8,
		Result: &FileInfo{
			Path:      "table1.org",
			Delimiter: ',',
			Format:    option.ORG,
			Encoding:  text.UTF8,
		},
	},
}

func TestNewFileInfoForCreate(t *testing.T) {
	for _, v := range fileInfoForCreateTests {
		repo := v.Repository
		if 0 < len(repo) {
			dir, _ := os.Getwd()
			repo = filepath.Join(dir, repo)
		}

		fileInfo, err := NewFileInfoForCreate(v.FilePath, repo, v.Delimiter, v.Encoding)
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
			t.Errorf("%s: FileInfo.Path = %s, want %s", v.Name, fileInfo.Path, abs)
		}
		if fileInfo.Delimiter != v.Result.Delimiter {
			t.Errorf("%s: FileInfo.Delimiter = %q, want %q", v.Name, fileInfo.Delimiter, v.Result.Delimiter)
		}
		if fileInfo.Format != v.Result.Format {
			t.Errorf("%s: FileInfo.Format = %s, want %s", v.Name, fileInfo.Format, v.Result.Format)
		}
	}
}
