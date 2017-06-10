package action

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/mithrandie/csvq/lib/cmd"
)

var executeTests = []struct {
	Name    string
	Input   string
	OutFile string
	Output  string
	Content string
	Error   string
}{
	{
		Name:    "Select Query Output To A File",
		Input:   "select 1 from dual",
		OutFile: GetTestFilePath("select_query_output_a_file.csv"),
		Content: "+---+\n" +
			"| 1 |\n" +
			"+---+\n" +
			"| 1 |\n" +
			"+---+\n",
	},
	{
		Name:   "Print",
		Input:  "var @a := 1; print @a;",
		Output: "1\n",
	},
	{
		Name:  "Query Execution Error",
		Input: "select from",
		Error: "syntax error: unexpected FROM",
	},
}

func initFlags() {
	tf := cmd.GetFlags()
	tf.Repository = TestDir
	tf.OutFile = ""
	tf.Format = cmd.TEXT
}

func TestWrite(t *testing.T) {

	for _, v := range executeTests {
		initFlags()
		tf := cmd.GetFlags()
		if v.OutFile != tf.OutFile {
			tf.OutFile = v.OutFile
		}

		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		err := Write(v.Input)

		w.Close()
		os.Stdout = oldStdout
		stdout, _ := ioutil.ReadAll(r)

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

		if string(stdout) != v.Output {
			t.Errorf("%s: output = %q, want %q", v.Name, string(stdout), v.Output)
		}

		if 0 < len(v.OutFile) {
			fp, _ := os.Open(v.OutFile)
			buf, _ := ioutil.ReadAll(fp)
			if string(buf) != v.Content {
				t.Errorf("%s: content = %q, want %q", v.Name, string(buf), v.Content)
			}
		}
	}
}
