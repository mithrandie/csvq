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
	Format     cmd.Format
	Result     *FileInfo
	Error      string
}{
	{
		Name:       "CSV",
		FilePath:   parser.Identifier{Literal: "table1"},
		Repository: TestDir,
		Delimiter:  cmd.UNDEF,
		Format:     cmd.CSV,
		Result: &FileInfo{
			Path:      "table1.csv",
			Delimiter: ',',
			Format:    cmd.CSV,
		},
	},
	{
		Name:       "CSV with AutoSelect",
		FilePath:   parser.Identifier{Literal: "table1"},
		Repository: TestDir,
		Delimiter:  cmd.UNDEF,
		Format:     cmd.AutoSelect,
		Result: &FileInfo{
			Path:      "table1.csv",
			Delimiter: ',',
			Format:    cmd.CSV,
		},
	},
	{
		Name:       "TSV",
		FilePath:   parser.Identifier{Literal: "table3"},
		Repository: TestDir,
		Delimiter:  cmd.UNDEF,
		Format:     cmd.TSV,
		Result: &FileInfo{
			Path:      "table3.tsv",
			Delimiter: '\t',
			Format:    cmd.TSV,
		},
	},
	{
		Name:       "TSV with AutoSelect",
		FilePath:   parser.Identifier{Literal: "table3"},
		Repository: TestDir,
		Delimiter:  cmd.UNDEF,
		Format:     cmd.AutoSelect,
		Result: &FileInfo{
			Path:      "table3.tsv",
			Delimiter: '\t',
			Format:    cmd.TSV,
		},
	},
	{
		Name:       "JSON",
		FilePath:   parser.Identifier{Literal: "table"},
		Repository: TestDir,
		Delimiter:  cmd.UNDEF,
		Format:     cmd.JSON,
		Result: &FileInfo{
			Path:      "table.json",
			Delimiter: ',',
			Format:    cmd.JSON,
		},
	},
	{
		Name:       "JSON with AutoSelect",
		FilePath:   parser.Identifier{Literal: "table"},
		Repository: TestDir,
		Delimiter:  cmd.UNDEF,
		Format:     cmd.AutoSelect,
		Result: &FileInfo{
			Path:      "table.json",
			Delimiter: ',',
			Format:    cmd.JSON,
		},
	},
	{
		Name:       "Fixed-Length",
		FilePath:   parser.Identifier{Literal: "fixed_length.txt"},
		Repository: TestDir,
		Delimiter:  cmd.UNDEF,
		Format:     cmd.FIXED,
		Result: &FileInfo{
			Path:      "fixed_length.txt",
			Delimiter: ',',
			Format:    cmd.FIXED,
		},
	},
	{
		Name:       "Not Exist Error",
		FilePath:   parser.Identifier{Literal: "notexist"},
		Repository: TestDir,
		Delimiter:  cmd.UNDEF,
		Format:     cmd.CSV,
		Error:      "[L:- C:-] file notexist does not exist",
	},
	{
		Name:       "File Read Error",
		FilePath:   parser.Identifier{Literal: TestDir},
		Repository: TestDir,
		Delimiter:  cmd.UNDEF,
		Format:     cmd.CSV,
		Error:      fmt.Sprintf("[L:- C:-] file %s is unable to be read", TestDir),
	},
	{
		Name:       "Filenames Ambiguous",
		FilePath:   parser.Identifier{Literal: "dup_name"},
		Repository: TestDir,
		Delimiter:  cmd.UNDEF,
		Format:     cmd.AutoSelect,
		Error:      fmt.Sprintf("[L:- C:-] filename dup_name is ambiguous"),
	},
}

func TestNewFileInfo(t *testing.T) {
	for _, v := range fileInfoTests {
		fileInfo, err := NewFileInfo(v.FilePath, v.Repository, v.Delimiter, v.Format)
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
	Result     *FileInfo
	Error      string
}{
	{
		Name:      "CSV",
		FilePath:  parser.Identifier{Literal: "table1.csv"},
		Delimiter: cmd.UNDEF,
		Result: &FileInfo{
			Path:      "table1.csv",
			Delimiter: ',',
			Format:    cmd.CSV,
		},
	},
	{
		Name:      "TSV",
		FilePath:  parser.Identifier{Literal: "table1.tsv"},
		Delimiter: cmd.UNDEF,
		Result: &FileInfo{
			Path:      "table1.tsv",
			Delimiter: '\t',
			Format:    cmd.TSV,
		},
	},
	{
		Name:      "JSON",
		FilePath:  parser.Identifier{Literal: "table1.json"},
		Delimiter: cmd.UNDEF,
		Result: &FileInfo{
			Path:      "table1.json",
			Delimiter: ',',
			Format:    cmd.JSON,
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

		fileInfo, err := NewFileInfoForCreate(v.FilePath, repo, v.Delimiter)
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

var fileInfoEquivalentTests = []struct {
	FileInfo1 *FileInfo
	FileInfo2 *FileInfo
	Expect    bool
}{
	{
		FileInfo1: &FileInfo{
			Path:      "table1.csv",
			Delimiter: ',',
			Format:    cmd.CSV,
		},
		FileInfo2: &FileInfo{
			Path:      "table1.csv",
			Delimiter: ',',
			Format:    cmd.CSV,
		},
		Expect: true,
	},
	{
		FileInfo1: &FileInfo{
			Path:      "table1.csv",
			Delimiter: ',',
			Format:    cmd.CSV,
		},
		FileInfo2: &FileInfo{
			Path:      "table1.csv",
			Delimiter: ',',
			Format:    cmd.TSV,
		},
		Expect: false,
	},
}

func TestFileInfo_Equivalent(t *testing.T) {
	for _, v := range fileInfoEquivalentTests {
		result := v.FileInfo1.Equivalent(v.FileInfo2)
		if result != v.Expect {
			t.Errorf("result = %t, want %t for %v, %v", result, v.Expect, v.FileInfo1, v.FileInfo2)
		}
	}
}
