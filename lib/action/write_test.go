package action

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/mithrandie/csvq/lib/cmd"
)

var executeTests = []struct {
	Name       string
	Input      string
	OutFile    string
	UpdateFile string
	Format     cmd.Format
	Output     string
	Content    string
	Error      string
}{
	{
		Name:  "Select Query",
		Input: "select 1 from dual",
		Output: "+---+\n" +
			"| 1 |\n" +
			"+---+\n" +
			"| 1 |\n" +
			"+---+\n",
	},
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
		Name:       "Insert Query",
		Input:      "insert into insert_query values (4, 'str4'), (5, 'str5')",
		Output:     fmt.Sprintf("%d record(s) inserted on %q\n", 2, GetTestFilePath("insert_query.csv")),
		UpdateFile: GetTestFilePath("insert_query.csv"),
		Content: "\"column1\",\"column2\"\n" +
			"\"1\",\"str1\"\n" +
			"\"2\",\"str2\"\n" +
			"\"3\",\"str3\"\n" +
			"4,\"str4\"\n" +
			"5,\"str5\"",
	},
	{
		Name:       "Update Query",
		Input:      "update update_query set column2 = 'update' where column1 = 2",
		Output:     fmt.Sprintf("%d record(s) updated on %q\n", 1, GetTestFilePath("update_query.csv")),
		UpdateFile: GetTestFilePath("update_query.csv"),
		Content: "\"column1\",\"column2\"\n" +
			"\"1\",\"str1\"\n" +
			"\"2\",\"update\"\n" +
			"\"3\",\"str3\"",
	},
	{
		Name:   "Update Query No Record Updated",
		Input:  "update update_query set column2 = 'update' where false",
		Output: fmt.Sprintf("no record updated on %q\n", GetTestFilePath("update_query.csv")),
	},
	{
		Name:       "Delete Query",
		Input:      "delete from delete_query where column1 = 2",
		Output:     fmt.Sprintf("%d record(s) deleted on %q\n", 1, GetTestFilePath("delete_query.csv")),
		UpdateFile: GetTestFilePath("delete_query.csv"),
		Content: "\"column1\",\"column2\"\n" +
			"\"1\",\"str1\"\n" +
			"\"3\",\"str3\"",
	},
	{
		Name:   "Delete Query No Record Deleted",
		Input:  "delete from delete_query where false",
		Output: fmt.Sprintf("no record deleted on %q\n", GetTestFilePath("delete_query.csv")),
	},
	{
		Name:       "Create Table",
		Input:      "create table create_table.csv (column1, column2)",
		Output:     fmt.Sprintf("file %q is created\n", GetTestFilePath("create_table.csv")),
		UpdateFile: GetTestFilePath("create_table.csv"),
		Content:    "\"column1\",\"column2\"\n",
	},
	{
		Name:       "Add Columns",
		Input:      "alter table add_columns add column3",
		Output:     fmt.Sprintf("%d field(s) added on %q\n", 1, GetTestFilePath("add_columns.csv")),
		UpdateFile: GetTestFilePath("add_columns.csv"),
		Content: "\"column1\",\"column2\",\"column3\"\n" +
			"\"1\",\"str1\",\n" +
			"\"2\",\"str2\",\n" +
			"\"3\",\"str3\",",
	},
	{
		Name:       "Drop Columns",
		Input:      "alter table drop_columns drop column1",
		Output:     fmt.Sprintf("%d field(s) dropped on %q\n", 1, GetTestFilePath("drop_columns.csv")),
		UpdateFile: GetTestFilePath("drop_columns.csv"),
		Content: "\"column2\"\n" +
			"\"str1\"\n" +
			"\"str2\"\n" +
			"\"str3\"",
	},
	{
		Name:       "Rename Column",
		Input:      "alter table rename_column rename column1 to newcolumn",
		Output:     fmt.Sprintf("%d field(s) renamed on %q\n", 1, GetTestFilePath("rename_column.csv")),
		UpdateFile: GetTestFilePath("rename_column.csv"),
		Content: "\"newcolumn\",\"column2\"\n" +
			"\"1\",\"str1\"\n" +
			"\"2\",\"str2\"\n" +
			"\"3\",\"str3\"",
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
		if v.Format != tf.Format {
			tf.Format = v.Format
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

		if 0 < len(v.UpdateFile) {
			fp, _ := os.Open(v.UpdateFile)
			buf, _ := ioutil.ReadAll(fp)
			if string(buf) != v.Content {
				t.Errorf("%s: content = %q, want %q", v.Name, string(buf), v.Content)
			}
		}
	}
}
