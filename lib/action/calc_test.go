package action

import (
	"io/ioutil"
	"os"
	"testing"

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
		Stdin:  "1,\"a\",1.234,true,unknown,\"2012-01-01 01:00:00\",null",
		Input:  "c1,c2,c3,boolean(c4),null = true,datetime_format(c6, '%Y-%m-%d'),c7",
		Output: "1,a,1.234,true,UNKNOWN,2012-01-01,null",
	},
	{
		Stdin: "foo",
		Input: "from",
		Error: "syntax error: unexpected FROM",
	},
	{
		Stdin: "",
		Input: "md5(c1)",
		Error: "stdin is empty",
	},
	{
		Stdin: "foo",
		Input: "error",
		Error: "field error does not exist",
	},
}

func TestCalc(t *testing.T) {
	for _, v := range calcTests {
		query.ViewCache.Clear()

		var oldStdin *os.File
		if 0 < len(v.Stdin) {
			oldStdin = os.Stdin
			r, w, _ := os.Pipe()
			w.WriteString(v.Stdin)
			w.Close()
			os.Stdin = r
		}
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		err := Calc(v.Input)

		if 0 < len(v.Stdin) {
			os.Stdin = oldStdin
		}
		w.Close()
		os.Stdout = oldStdout
		stdout, _ := ioutil.ReadAll(r)

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

		if string(stdout) != v.Output {
			t.Errorf("%s: output = %q, want %q", v.Input, string(stdout), v.Output)
		}
	}
}
