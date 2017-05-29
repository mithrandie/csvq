package query

import (
	"os"
	"path"
	"reflect"
	"testing"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/parser"
)

var executeTests = []struct {
	Input  string
	Result []Result
	Error  string
}{
	{
		Input: "select 1",
		Result: []Result{
			{
				Statement: SELECT,
				View: &View{
					Header: []HeaderField{
						{
							Alias: "1",
						},
					},
					Records: []Record{
						{
							NewCell(parser.NewInteger(1)),
						},
					},
				},
				Count: 1,
			},
		},
	},
	{
		Input: "select column1 from table1 where column1 = 1 group by column1 having sum(column1) > 0 order by column1 limit 10",
		Result: []Result{
			{
				Statement: SELECT,
				View: &View{
					Header: []HeaderField{
						{
							Reference:  "table1",
							Column:     "column1",
							FromTable:  true,
							IsGroupKey: true,
						},
					},
					Records: []Record{
						{
							NewCell(parser.NewString("1")),
						},
					},
					FileInfo: &FileInfo{
						Path:      "table1.csv",
						Delimiter: ',',
					},
				},
				Count: 1,
			},
		},
	},
	{
		Input: "select from notexist",
		Error: "syntax error: unexpected FROM",
	},
	{
		Input: "select column1 from notexist",
		Error: "file notexist does not exist",
	},
	{
		Input: "select column1 from table1 where notexist = 1",
		Error: "identifier = notexist: field does not exist",
	},
	{
		Input: "select column1 from table1 group by notexist",
		Error: "identifier = notexist: field does not exist",
	},
	{
		Input: "select column1 from table1 having notexist",
		Error: "identifier = notexist: field does not exist",
	},
	{
		Input: "select column1 from table1 order by notexist",
		Error: "identifier = notexist: field does not exist",
	},
	{
		Input: "select notexist",
		Error: "identifier = notexist: field does not exist",
	},
}

func TestExecute(t *testing.T) {
	tf := cmd.GetFlags()
	dir, _ := os.Getwd()
	tf.Repository = path.Join(dir, "..", "..", "testdata", "csv")

	for _, v := range executeTests {
		results, err := Execute(v.Input)
		if err != nil {
			if len(v.Error) < 1 {
				t.Errorf("unexpected error %q for %q", err, v.Input)
			} else if err.Error() != v.Error {
				t.Errorf("error %q, want error %q for %q", err, v.Error, v.Input)
			}
			continue
		}
		if 0 < len(v.Error) {
			t.Errorf("no error, want error %q for %q", v.Error, v.Input)
			continue
		}

		for i, result := range results {
			if result.View.FileInfo != nil {
				if path.Base(result.View.FileInfo.Path) != v.Result[i].View.FileInfo.Path {
					t.Errorf("filepath = %s, want %s for %q", path.Base(result.View.FileInfo.Path), v.Result[i].View.FileInfo.Path, v.Input)
				}
				if result.View.FileInfo.Delimiter != v.Result[i].View.FileInfo.Delimiter {
					t.Errorf("delimiter = %q, want %q for %q", result.View.FileInfo.Delimiter, v.Result[i].View.FileInfo.Delimiter, v.Input)
				}
			}
			result.View.FileInfo = nil
			v.Result[i].View.FileInfo = nil
		}
		if !reflect.DeepEqual(results, v.Result) {
			t.Errorf("results = %q, want %q for %q", results, v.Result, v.Input)
		}
	}
}
