package action

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/mithrandie/csvq/lib/query"

	"github.com/mithrandie/csvq/lib/cmd"
)

var executeTests = []struct {
	Name    string
	Input   string
	OutFile string
	Output  string
	Stats   bool
	Content string
	Error   string
}{
	{
		Name:    "Select Query Output To File",
		Input:   "select 1 from dual",
		OutFile: GetTestFilePath("select_query_output_file.csv"),
		Content: "" +
			"+---+\n" +
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
		Error: "[L:1 C:8] syntax error: unexpected token \"from\"",
	},
	{
		Name:  "Show Statistics",
		Input: "select 1",
		Stats: true,
	},
}

func initFlags() {
	tf := cmd.GetFlags()
	tf.Repository = TestDir
	tf.Format = cmd.TEXT
	tf.Stats = false
}

func TestRun(t *testing.T) {

	for _, v := range executeTests {
		initFlags()
		tf := cmd.GetFlags()
		if v.Stats {
			tf.Stats = v.Stats
		}

		query.OutFile = nil

		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		err := Run(v.Input, "", v.OutFile)

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

		if v.Stats {
			if !strings.Contains(string(stdout), "Time:") {
				t.Errorf("%s: output = %q, want statistics", v.Name, string(stdout))
			}
		} else {
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
}
