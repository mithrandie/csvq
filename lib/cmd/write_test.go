package cmd

import (
	"github.com/mithrandie/go-file"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

type writeTest struct {
	Name     string
	Filename string
	Content  string
	Result   string
	Error    string
}

var createFileTests = []writeTest{
	{
		Name:     "Create",
		Filename: "create.txt",
		Content:  "write",
		Result:   "write",
	},
	{
		Name:     "Output to Stdout",
		Filename: "",
		Content:  "write",
		Result:   "write",
	},
	{
		Name:     "File Exists Error",
		Filename: "create.txt",
		Error:    "environment-dependent",
	},
	{
		Name:     "File Open Error",
		Filename: filepath.Join("notexistdir", "create.txt"),
		Error:    "environment-dependent",
	},
}

func TestCreateFile(t *testing.T) {
	for _, v := range createFileTests {
		if len(v.Filename) < 1 {
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			ToStdout(v.Content)

			w.Close()
			os.Stdout = oldStdout

			buf, _ := ioutil.ReadAll(r)
			if string(buf) != v.Result {
				t.Errorf("%s: content = %q, want %q", v.Name, string(buf), v.Result)
			}
		} else {
			filename := GetTestFilePath(v.Filename)
			err := CreateFile(filename, v.Content)
			if err != nil {
				if len(v.Error) < 1 {
					t.Errorf("%s: unexpected error %q", v.Name, err)
				} else if v.Error != "environment-dependent" && err.Error() != v.Error {
					t.Errorf("%s: error %q, want error %q", v.Name, err.Error(), v.Error)
				}
				continue
			}
			if 0 < len(v.Error) {
				t.Errorf("%s: no error, want error %q", v.Name, v.Error)
				continue
			}

			fp, _ := os.Open(filename)
			buf, _ := ioutil.ReadAll(fp)
			if string(buf) != v.Result {
				t.Errorf("%s: content = %q, want %q", v.Name, string(buf), v.Result)
			}
		}
	}
}

var updateFileTests = []writeTest{
	{
		Name:     "Update",
		Filename: "create.txt",
		Content:  "truncate and write",
		Result:   "truncate and write",
	},
}

func TestUpdateFile(t *testing.T) {
	for _, v := range updateFileTests {
		filename := GetTestFilePath(v.Filename)
		fp, _ := file.OpenToUpdate(filename)
		err := UpdateFile(fp, v.Content)
		if err != nil {
			if len(v.Error) < 1 {
				t.Errorf("%s: unexpected error %q", v.Name, err)
			} else if v.Error != "environment-dependent" && err.Error() != v.Error {
				t.Errorf("%s: error %q, want error %q", v.Name, err.Error(), v.Error)
			}
			continue
		}
		if 0 < len(v.Error) {
			t.Errorf("%s: no error, want error %q", v.Name, v.Error)
			continue
		}

		fp, _ = os.Open(filename)
		buf, _ := ioutil.ReadAll(fp)
		if string(buf) != v.Result {
			t.Errorf("%s: content = %q, want %q", v.Name, string(buf), v.Result)
		}
	}
}

func TestTryCreateFile(t *testing.T) {
	err := TryCreateFile(GetTestFilePath("table1.csv"))
	if err == nil {
		t.Error("Create table1.csv: no error, want error")
	}

	err = TryCreateFile(GetTestFilePath("notexist.csv"))
	if err != nil {
		t.Errorf("Create notexist.csv: unexpected error %q", err)
	} else {
		if _, err := os.Stat(GetTestFilePath("notexist.csv")); err == nil {
			t.Errorf("Create notexist.csv: temporary file does not removed")
		}
	}
}
