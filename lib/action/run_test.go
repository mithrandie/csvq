package action

import (
	"context"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/mithrandie/csvq/lib/file"

	"github.com/mithrandie/csvq/lib/query"
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

func TestRun(t *testing.T) {
	tx, _ := query.NewTransaction(context.Background(), file.DefaultWaitTimeout, file.DefaultRetryDelay, query.NewSession())
	tx.UseColor(false)
	ctx := context.Background()

	for _, v := range executeTests {
		if v.Stats {
			tx.Flags.Stats = v.Stats
		}

		tx.Session.SetOutFile(nil)

		out := query.NewOutput()
		tx.Session.SetStdout(out)

		proc := query.NewProcessor(tx)
		err := Run(ctx, proc, v.Input, "", v.OutFile)

		stdout := out.String()

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
			if !strings.Contains(stdout, "Time:") {
				t.Errorf("%s: output = %q, want statistics", v.Name, stdout)
			}
		} else {
			if stdout != v.Output {
				t.Errorf("%s: output = %q, want %q", v.Name, stdout, v.Output)
			}

			if 0 < len(v.OutFile) {
				fp, _ := os.Open(v.OutFile)
				buf, _ := io.ReadAll(fp)
				if string(buf) != v.Content {
					t.Errorf("%s: content = %q, want %q", v.Name, string(buf), v.Content)
				}
			}
		}
	}
}
