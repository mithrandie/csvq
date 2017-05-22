package main

import (
	"testing"

	"github.com/mithrandie/csvq/lib/cmd"
)

var executeTests = []struct {
	Name       string
	Input      string
	Repository string
	OutFile    string
	Format     cmd.Format
	Error      string
}{
	{
		Name:  "Query Execution Error",
		Input: "select from",
		Error: "syntax error: unexpected FROM",
	},
}

func initFlags() {
	tf := cmd.GetFlags()
	tf.Repository = "."
	tf.OutFile = ""
	tf.Format = cmd.STDOUT
}

func TestWrite(t *testing.T) {
	for _, v := range executeTests {
		initFlags()
		tf := cmd.GetFlags()
		if v.Repository != tf.Repository {
			tf.Repository = v.Repository
		}
		if v.OutFile != tf.OutFile {
			tf.OutFile = v.OutFile
		}
		if v.Format != tf.Format {
			tf.Format = v.Format
		}

		err := Write(v.Input)
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
	}
}

func ExampleWrite() {
	initFlags()
	tf := cmd.GetFlags()
	tf.OutFile = ""
	tf.Format = cmd.STDOUT

	Write("select 1 from dual where false")
	Write("select 1 from dual")
	//OUTPUT:
	//Empty
	//+---+
	//| 1 |
	//+---+
	//| 1 |
	//+---+
}
