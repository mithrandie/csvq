package action

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/mithrandie/csvq/lib/file"

	"github.com/mithrandie/csvq/lib/query"
)

var calcTests = []struct {
	Stdin  string
	Input  string
	Output string
	Error  string
}{
	{
		Stdin:  "foo",
		Input:  "md5(c1)",
		Output: "acbd18db4cc2f85cedef654fccc4a4d8",
	},
	{
		Stdin:  "1,\"a\",1.234,true,unknown,\"2012-01-01 01:00:00 +00:00\",",
		Input:  "integer(c1),c2,float(c3),boolean(c4),null = true,datetime(c6),c7",
		Output: "1,a,1.234,true,UNKNOWN,2012-01-01T01:00:00Z,NULL",
	},
	{
		Stdin: "foo",
		Input: "from",
		Error: "syntax error",
	},
	{
		Stdin: "",
		Input: "md5(c1)",
		Error: "STDIN is empty",
	},
	{
		Stdin: "foo",
		Input: "error",
		Error: "field error does not exist",
	},
}

func TestCalc(t *testing.T) {
	tx, _ := query.NewTransaction(context.Background(), file.DefaultWaitTimeout, file.DefaultRetryDelay, query.NewSession())
	scope := query.NewReferenceScope(tx)
	ctx := context.Background()

	for _, v := range calcTests {
		_ = tx.Rollback(scope, nil)

		if 0 < len(v.Stdin) {
			_ = tx.Session.SetStdin(query.NewInput(strings.NewReader(v.Stdin)))
		}
		out := query.NewOutput()
		tx.Session.SetStdout(out)

		err := Calc(ctx, query.NewProcessor(tx), v.Input)

		if 0 < len(v.Stdin) {
			_ = tx.Session.SetStdin(os.Stdin)
		}
		stdout := out.String()

		if err != nil {
			if len(v.Error) < 1 {
				t.Errorf("%s: unexpected error %q", v.Input, err)
			} else if err.Error() != v.Error {
				t.Errorf("%s: error %q, want error %q", v.Input, err.Error(), v.Error)
			}
			continue
		}
		if 0 < len(v.Error) {
			t.Errorf("%s: no error, want error %q", v.Input, v.Error)
			continue
		}

		if stdout != v.Output {
			t.Errorf("%s: output = %q, want %q", v.Input, stdout, v.Output)
		}
	}
}
