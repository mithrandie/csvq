package query

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/value"
)

var printTests = []struct {
	Name   string
	Expr   parser.Print
	Result string
	Error  string
}{
	{
		Name: "Print",
		Expr: parser.Print{
			Value: parser.NewStringValue("foo"),
		},
		Result: "'foo'",
	},
	{
		Name: "Print Error",
		Expr: parser.Print{
			Value: parser.Variable{
				Name: "var",
			},
		},
		Error: "[L:- C:-] variable var is undeclared",
	},
}

func TestPrint(t *testing.T) {
	filter := NewEmptyFilter()

	for _, v := range printTests {
		result, err := Print(v.Expr, filter)
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
		if result != v.Result {
			t.Errorf("%s: result = %q, want %q", v.Name, result, v.Result)
		}
	}
}

var printfTests = []struct {
	Name   string
	Expr   parser.Printf
	Result string
	Error  string
}{
	{
		Name: "Printf",
		Expr: parser.Printf{
			Format: parser.NewStringValue("printf test: value1 %q, value2 %q, %a %% %"),
			Values: []parser.QueryExpression{
				parser.NewStringValue("str"),
				parser.NewIntegerValue(1),
			},
		},
		Result: "printf test: value1 'str', value2 1, %a % %",
	},
	{
		Name: "Printf Format Error",
		Expr: parser.Printf{
			Format: parser.Variable{Name: "var"},
			Values: []parser.QueryExpression{
				parser.NewStringValue("str"),
				parser.NewIntegerValue(1),
			},
		},
		Error: "[L:- C:-] variable var is undeclared",
	},
	{
		Name: "Printf Evaluate Error",
		Expr: parser.Printf{
			Format: parser.NewStringValue("printf test: value1 %s"),
			Values: []parser.QueryExpression{
				parser.Variable{
					Name: "var",
				},
			},
		},
		Error: "[L:- C:-] variable var is undeclared",
	},
	{
		Name: "Printf Less Values Error",
		Expr: parser.Printf{
			Format: parser.NewStringValue("printf test: value1 %s, value2 %s, %a %% %"),
			Values: []parser.QueryExpression{
				parser.NewStringValue("str"),
			},
		},
		Error: "[L:- C:-] PRINTF: number of replace values does not match",
	},
	{
		Name: "Printf Greater Values Error",
		Expr: parser.Printf{
			Format: parser.NewStringValue("printf test: value1 %s, value2 %s, %a %% %"),
			Values: []parser.QueryExpression{
				parser.NewStringValue("str"),
				parser.NewIntegerValue(1),
				parser.NewIntegerValue(2),
			},
		},
		Error: "[L:- C:-] PRINTF: number of replace values does not match",
	},
}

func TestPrintf(t *testing.T) {
	filter := NewEmptyFilter()

	for _, v := range printfTests {
		result, err := Printf(v.Expr, filter)
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
		if result != v.Result {
			t.Errorf("%s: result = %q, want %q", v.Name, result, v.Result)
		}
	}
}

var sourceTests = []struct {
	Name   string
	Expr   parser.Source
	Result []parser.Statement
	Error  string
}{
	{
		Name: "Source",
		Expr: parser.Source{
			FilePath: parser.NewStringValue(GetTestFilePath("source.sql")),
		},
		Result: []parser.Statement{
			parser.Print{
				Value: parser.NewStringValue("external executable file"),
			},
		},
	},
	{
		Name: "Source File Argument Evaluation Error",
		Expr: parser.Source{
			FilePath: parser.FieldReference{Column: parser.Identifier{Literal: "ident"}},
		},
		Error: "[L:- C:-] field ident does not exist",
	},
	{
		Name: "Source File Argument Not String Error",
		Expr: parser.Source{
			FilePath: parser.NewNullValueFromString("NULL"),
		},
		Error: "[L:- C:-] SOURCE: argument NULL is not a string",
	},
	{
		Name: "Source File Not Exist Error",
		Expr: parser.Source{
			FilePath: parser.NewStringValue(GetTestFilePath("notexist.sql")),
		},
		Error: fmt.Sprintf("[L:- C:-] SOURCE: file %s does not exist", GetTestFilePath("notexist.sql")),
	},
	{
		Name: "Source File Not Readable Error",
		Expr: parser.Source{
			FilePath: parser.NewStringValue(TestDir),
		},
		Error: fmt.Sprintf("[L:- C:-] SOURCE: file %s is unable to read", TestDir),
	},
	{
		Name: "Source Syntax Error",
		Expr: parser.Source{
			FilePath: parser.NewStringValue(GetTestFilePath("source_syntaxerror.sql")),
		},
		Error: fmt.Sprintf("%s [L:1 C:34] syntax error: unexpected STRING", GetTestFilePath("source_syntaxerror.sql")),
	},
}

func TestSource(t *testing.T) {
	filter := NewEmptyFilter()

	for _, v := range sourceTests {
		result, err := Source(v.Expr, filter)
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
		if !reflect.DeepEqual(result, v.Result) {
			t.Errorf("%s: result = %q, want %q", v.Name, result, v.Result)
		}
	}
}

var setFlagTests = []struct {
	Name             string
	Expr             parser.SetFlag
	ResultFlag       string
	ResultStrValue   string
	ResultFloatValue float64
	ResultBoolValue  bool
	Error            string
}{
	{
		Name: "Set Delimiter",
		Expr: parser.SetFlag{
			Name:  "@@delimiter",
			Value: value.NewString("\t"),
		},
		ResultFlag:     "delimiter",
		ResultStrValue: "\t",
	},
	{
		Name: "Set Encoding",
		Expr: parser.SetFlag{
			Name:  "@@encoding",
			Value: value.NewString("SJIS"),
		},
		ResultFlag:     "encoding",
		ResultStrValue: "SJIS",
	},
	{
		Name: "Set LineBreak",
		Expr: parser.SetFlag{
			Name:  "@@line_break",
			Value: value.NewString("CRLF"),
		},
		ResultFlag:     "line_break",
		ResultStrValue: "\r\n",
	},
	{
		Name: "Set Timezone",
		Expr: parser.SetFlag{
			Name:  "@@timezone",
			Value: value.NewString("utc"),
		},
		ResultFlag:     "timezone",
		ResultStrValue: "UTC",
	},
	{
		Name: "Set Repository",
		Expr: parser.SetFlag{
			Name:  "@@repository",
			Value: value.NewString(TestDir),
		},
		ResultFlag:     "repository",
		ResultStrValue: TestDir,
	},
	{
		Name: "Set DatetimeFormat",
		Expr: parser.SetFlag{
			Name:  "@@datetime_format",
			Value: value.NewString("%Y%m%d"),
		},
		ResultFlag:     "datetime_format",
		ResultStrValue: "%Y%m%d",
	},
	{
		Name: "Set WaitTimeout",
		Expr: parser.SetFlag{
			Name:  "@@wait_timeout",
			Value: value.NewFloat(15),
		},
		ResultFlag:       "wait_timeout",
		ResultFloatValue: 15,
	},
	{
		Name: "Set NoHeader",
		Expr: parser.SetFlag{
			Name:  "@@no_header",
			Value: value.NewBoolean(true),
		},
		ResultFlag:      "no_header",
		ResultBoolValue: true,
	},
	{
		Name: "Set WithoutNull",
		Expr: parser.SetFlag{
			Name:  "@@without_null",
			Value: value.NewBoolean(true),
		},
		ResultFlag:      "without_null",
		ResultBoolValue: true,
	},
	{
		Name: "Set Color",
		Expr: parser.SetFlag{
			Name:  "@@color",
			Value: value.NewBoolean(true),
		},
		ResultFlag:      "color",
		ResultBoolValue: true,
	},
	{
		Name: "Set Stats",
		Expr: parser.SetFlag{
			Name:  "@@stats",
			Value: value.NewBoolean(true),
		},
		ResultFlag:      "stats",
		ResultBoolValue: true,
	},
	{
		Name: "Set Delimiter Value Error",
		Expr: parser.SetFlag{
			Name:  "@@delimiter",
			Value: value.NewBoolean(true),
		},
		Error: "[L:- C:-] SET: flag value true for @@delimiter is invalid",
	},
	{
		Name: "Set WaitTimeout Value Error",
		Expr: parser.SetFlag{
			Name:  "@@wait_timeout",
			Value: value.NewBoolean(true),
		},
		Error: "[L:- C:-] SET: flag value true for @@wait_timeout is invalid",
	},
	{
		Name: "Set WithoutNull Value Error",
		Expr: parser.SetFlag{
			Name:  "@@without_null",
			Value: value.NewString("string"),
		},
		Error: "[L:- C:-] SET: flag value 'string' for @@without_null is invalid",
	},
	{
		Name: "Invalid Flag Name Error",
		Expr: parser.SetFlag{
			Name:  "@@invalid",
			Value: value.NewString("string"),
		},
		Error: "[L:- C:-] flag name @@invalid is invalid",
	},
	{
		Name: "Invalid Flag Value Error",
		Expr: parser.SetFlag{
			Name:  "@@line_break",
			Value: value.NewString("invalid"),
		},
		Error: "[L:- C:-] SET: flag value 'invalid' for @@line_break is invalid",
	},
}

func TestSetFlag(t *testing.T) {
	flags := cmd.GetFlags()

	for _, v := range setFlagTests {
		initFlag()
		err := SetFlag(v.Expr)
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

		switch strings.ToUpper(v.ResultFlag) {
		case "DELIMITER":
			if string(flags.Delimiter) != v.ResultStrValue {
				t.Errorf("%s: delimiter = %q, want %q", v.Name, string(flags.Delimiter), v.ResultStrValue)
			}
		case "ENCODING":
			if flags.Encoding.String() != v.ResultStrValue {
				t.Errorf("%s: encoding = %q, want %q", v.Name, flags.Encoding.String(), v.ResultStrValue)
			}
		case "LINE_BREAK":
			if flags.LineBreak.Value() != v.ResultStrValue {
				t.Errorf("%s: line-break = %q, want %q", v.Name, flags.LineBreak.Value(), v.ResultStrValue)
			}
		case "TIMEZONE":
			if flags.Location != v.ResultStrValue {
				t.Errorf("%s: timezone = %q, want %q", v.Name, flags.Location, v.ResultStrValue)
			}
		case "REPOSITORY":
			if flags.Repository != v.ResultStrValue {
				t.Errorf("%s: repository = %q, want %q", v.Name, flags.Repository, v.ResultStrValue)
			}
		case "DATETIME_FORMAT":
			if flags.DatetimeFormat != v.ResultStrValue {
				t.Errorf("%s: datetime-format = %q, want %q", v.Name, flags.DatetimeFormat, v.ResultStrValue)
			}
		case "WAIT_TIMEOUT":
			if flags.WaitTimeout != v.ResultFloatValue {
				t.Errorf("%s: wait-timeout = %f, want %f", v.Name, flags.WaitTimeout, v.ResultFloatValue)
			}
		case "NO-HEADER":
			if flags.NoHeader != v.ResultBoolValue {
				t.Errorf("%s: no-header = %t, want %t", v.Name, flags.NoHeader, v.ResultBoolValue)
			}
		case "WITHOUT-NULL":
			if flags.WithoutNull != v.ResultBoolValue {
				t.Errorf("%s: without-null = %t, want %t", v.Name, flags.WithoutNull, v.ResultBoolValue)
			}
		case "COLOR":
			if flags.Color != v.ResultBoolValue {
				t.Errorf("%s: color = %t, want %t", v.Name, flags.Stats, v.ResultBoolValue)
			}
		case "STATS":
			if flags.Stats != v.ResultBoolValue {
				t.Errorf("%s: stats = %t, want %t", v.Name, flags.Stats, v.ResultBoolValue)
			}
		}
	}
	initFlag()
}

var showFlagTests = []struct {
	Name    string
	Expr    parser.ShowFlag
	SetExpr parser.SetFlag
	Result  string
	Error   string
}{
	{
		Name: "Show Delimiter Not Set",
		Expr: parser.ShowFlag{
			Name: "@@delimiter",
		},
		Result: "(not set)",
	},
	{
		Name: "Show Delimiter",
		Expr: parser.ShowFlag{
			Name: "@@delimiter",
		},
		SetExpr: parser.SetFlag{
			Name:  "@@delimiter",
			Value: value.NewString("\t"),
		},
		Result: "'\\t'",
	},
	{
		Name: "Show Encoding",
		Expr: parser.ShowFlag{
			Name: "@@encoding",
		},
		SetExpr: parser.SetFlag{
			Name:  "@@encoding",
			Value: value.NewString("SJIS"),
		},
		Result: "SJIS",
	},
	{
		Name: "Show LineBreak",
		Expr: parser.ShowFlag{
			Name: "@@line_break",
		},
		SetExpr: parser.SetFlag{
			Name:  "@@line_break",
			Value: value.NewString("CRLF"),
		},
		Result: "CRLF",
	},
	{
		Name: "Show Timezone",
		Expr: parser.ShowFlag{
			Name: "@@timezone",
		},
		SetExpr: parser.SetFlag{
			Name:  "@@timezone",
			Value: value.NewString("UTC"),
		},
		Result: "UTC",
	},
	{
		Name: "Show Repository",
		Expr: parser.ShowFlag{
			Name: "@@repository",
		},
		SetExpr: parser.SetFlag{
			Name:  "@@repository",
			Value: value.NewString(TestDir),
		},
		Result: TestDir,
	},
	{
		Name: "Show DatetimeFormat Not Set",
		Expr: parser.ShowFlag{
			Name: "@@datetime_format",
		},
		Result: "(not set)",
	},
	{
		Name: "Show DatetimeFormat",
		Expr: parser.ShowFlag{
			Name: "@@datetime_format",
		},
		SetExpr: parser.SetFlag{
			Name:  "@@datetime_format",
			Value: value.NewString("%Y%m%d"),
		},
		Result: "%Y%m%d",
	},
	{
		Name: "Show WaitTimeout",
		Expr: parser.ShowFlag{
			Name: "@@wait_timeout",
		},
		SetExpr: parser.SetFlag{
			Name:  "@@wait_timeout",
			Value: value.NewFloat(15),
		},
		Result: "15",
	},
	{
		Name: "Show NoHeader",
		Expr: parser.ShowFlag{
			Name: "@@no_header",
		},
		SetExpr: parser.SetFlag{
			Name:  "@@no_header",
			Value: value.NewBoolean(true),
		},
		Result: "true",
	},
	{
		Name: "Show WithoutNull",
		Expr: parser.ShowFlag{
			Name: "@@without_null",
		},
		SetExpr: parser.SetFlag{
			Name:  "@@without_null",
			Value: value.NewBoolean(true),
		},
		Result: "true",
	},
	{
		Name: "Show Color",
		Expr: parser.ShowFlag{
			Name: "@@color",
		},
		SetExpr: parser.SetFlag{
			Name:  "@@color",
			Value: value.NewBoolean(true),
		},
		Result: "true",
	},
	{
		Name: "Show Stats",
		Expr: parser.ShowFlag{
			Name: "@@stats",
		},
		SetExpr: parser.SetFlag{
			Name:  "@@stats",
			Value: value.NewBoolean(true),
		},
		Result: "true",
	},
	{
		Name: "Invalid Flag Name Error",
		Expr: parser.ShowFlag{
			Name: "@@invalid",
		},
		Error: "[L:- C:-] flag name @@invalid is invalid",
	},
}

func TestShowFlag(t *testing.T) {
	for _, v := range showFlagTests {
		initFlag()
		if v.SetExpr.Value != nil {
			SetFlag(v.SetExpr)
		}
		result, err := ShowFlag(v.Expr)
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
		if result != v.Result {
			t.Errorf("%s: result = %s, want %s", v.Name, result, v.Result)
		}
	}
	initFlag()
}

var showObjectsTests = []struct {
	Name       string
	Expr       parser.ShowObjects
	Filter     *Filter
	Repository string
	ViewCache  ViewMap
	Result     string
	Error      string
}{
	{
		Name: "ShowObjects Tables",
		Expr: parser.ShowObjects{Type: parser.TABLES},
		ViewCache: ViewMap{
			"DUMMY.CSV": &View{
				FileInfo: &FileInfo{
					Path: filepath.Join(TestDir, "dummy.csv"),
				},
			},
			"TABLE1.CSV": &View{
				FileInfo: &FileInfo{
					Path: filepath.Join(filepath.Join(TestDir, "test_show_objects"), "table1.csv"),
				},
			},
		},
		Result: "\n" + fmt.Sprintf("    Tables in %s\n", filepath.Join(TestDir, "test_show_objects")) + strings.Repeat("-", len(filepath.Join(TestDir, "test_show_objects"))+18) + "\n" +
			filepath.Join("table1.csv") + "\n" +
			filepath.Join("table2.csv") + "\n" +
			"\n" + "    Tables in other directories\n" + "-----------------------------------\n" +
			filepath.Join(TestDir, "dummy.csv") + "\n",
	},
	{
		Name:       "ShowObjects Tables Empty",
		Expr:       parser.ShowObjects{Type: parser.TABLES},
		Repository: filepath.Join(TestDir, "test_show_objects_empty"),
		Result:     fmt.Sprintf("Repository %q is empty", filepath.Join(TestDir, "test_show_objects_empty")),
	},
	{
		Name: "ShowObjects Views",
		Expr: parser.ShowObjects{Type: parser.VIEWS},
		Filter: &Filter{
			TempViews: TemporaryViewScopes{
				ViewMap{
					"VIEW1": &View{
						FileInfo: &FileInfo{
							Path:        "view1",
							IsTemporary: true,
						},
						Header: NewHeader("view1", []string{"column1", "column2"}),
					},
				},
				ViewMap{
					"VIEW1": &View{
						FileInfo: &FileInfo{
							Path:        "view1",
							IsTemporary: true,
						},
						Header: NewHeader("view1", []string{"column1", "column2", "column3"}),
					},
					"VIEW2": &View{
						FileInfo: &FileInfo{
							Path:        "view2",
							IsTemporary: true,
						},
						Header: NewHeader("view2", []string{"column1", "column2"}),
					},
				},
			},
		},
		Result: "\n" + "    Views\n" + "-------------\n" +
			"view1 (column1, column2)\nview2 (column1, column2)\n",
	},
	{
		Name:   "ShowObjects Views Empty",
		Expr:   parser.ShowObjects{Type: parser.VIEWS},
		Result: "No view is declared",
	},
	{
		Name: "ShowObjects Cursors",
		Expr: parser.ShowObjects{Type: parser.CURSORS},
		Filter: &Filter{
			Cursors: CursorScopes{
				{
					"CUR": &Cursor{
						name:  "cur",
						query: selectQueryForCursorTest,
					},
					"CUR2": &Cursor{
						name:  "cur2",
						query: selectQueryForCursorTest,
					},
				},
			},
		},
		Result: "\n" + "    Cursors\n" + "---------------\n" +
			"cur for select column1, column2 from table1\ncur2 for select column1, column2 from table1\n",
	},
	{
		Name:   "ShowObjects Cursors Empty",
		Expr:   parser.ShowObjects{Type: parser.CURSORS},
		Result: "No cursor is declared",
	},
	{
		Name: "ShowObjects Functions",
		Expr: parser.ShowObjects{Type: parser.FUNCTIONS},
		Filter: &Filter{
			Functions: UserDefinedFunctionScopes{
				UserDefinedFunctionMap{
					"USERFUNC1": &UserDefinedFunction{
						Name: parser.Identifier{Literal: "userfunc1"},
						Parameters: []parser.Variable{
							{Name: "@arg1"},
						},
						Statements: []parser.Statement{
							parser.Print{Value: parser.Variable{Name: "@arg1"}},
						},
					},
				},
				UserDefinedFunctionMap{
					"USERAGGFUNC": &UserDefinedFunction{
						Name: parser.Identifier{Literal: "useraggfunc"},
						Parameters: []parser.Variable{
							{Name: "@arg1"},
							{Name: "@arg2"},
						},
						Defaults: map[string]parser.QueryExpression{
							"@arg2": parser.NewIntegerValue(1),
						},
						IsAggregate:  true,
						RequiredArgs: 1,
						Cursor:       parser.Identifier{Literal: "column1"},
						Statements: []parser.Statement{
							parser.Print{Value: parser.Variable{Name: "@var1"}},
						},
					},
				},
			},
		},
		Result: "\n" + "    Scala Functions\n" + "-----------------------\n" +
			"userfunc1 (@arg1)\n" +
			"\n" + "    Aggregate Functions\n" + "---------------------------\n" +
			"useraggfunc (column1, @arg1, @arg2 = 1)\n",
	},
	{
		Name:   "ShowObjects Functions Empty",
		Expr:   parser.ShowObjects{Type: parser.FUNCTIONS},
		Result: "No function is declared",
	},
}

func TestShowObjects(t *testing.T) {
	tableDir := filepath.Join(TestDir, "test_show_objects")
	emptyDir := filepath.Join(TestDir, "test_show_objects_empty")
	if _, err := os.Stat(tableDir); os.IsNotExist(err) {
		os.Mkdir(tableDir, 0755)
	}
	if _, err := os.Stat(emptyDir); os.IsNotExist(err) {
		os.Mkdir(emptyDir, 0755)
	}
	copyfile(filepath.Join(tableDir, "table1.csv"), filepath.Join(TestDataDir, "table1.csv"))
	copyfile(filepath.Join(tableDir, "table2.csv"), filepath.Join(TestDataDir, "table2.csv"))

	flags := cmd.GetFlags()

	for _, v := range showObjectsTests {
		if 0 < len(v.Repository) {
			flags.Repository = v.Repository
		} else {
			flags.Repository = tableDir
		}
		ViewCache.Clean()
		if 0 < len(v.ViewCache) {
			ViewCache = v.ViewCache
		}

		var filter *Filter
		if v.Filter != nil {
			filter = v.Filter
		} else {
			filter = NewEmptyFilter()
		}

		result, err := ShowObjects(v.Expr, filter)
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
		if result != v.Result {
			t.Errorf("%s: result = %s, want %s", v.Name, result, v.Result)
		}
	}
	ReleaseResources()
}

var showFieldsTests = []struct {
	Name      string
	Expr      parser.ShowFields
	Filter    *Filter
	ViewCache ViewMap
	Result    string
	Error     string
}{
	{
		Name: "ShowFields Temporary Table",
		Expr: parser.ShowFields{
			Table: parser.Identifier{Literal: "view1"},
		},
		Filter: &Filter{
			TempViews: TemporaryViewScopes{
				ViewMap{
					"VIEW1": &View{
						Header: NewHeader("view1", []string{"column1", "column2"}),
						FileInfo: &FileInfo{
							Path: "view1",
						},
					},
				},
			},
		},
		Result: "\n" + "    Fields in view1" + "\n" + "-----------------------\n" +
			"1. column1\n2. column2\n",
	},
	{
		Name: "ShowFields Created Table",
		Expr: parser.ShowFields{
			Table: parser.Identifier{Literal: GetTestFilePath("show_fields_create.csv")},
		},
		ViewCache: ViewMap{
			strings.ToUpper(GetTestFilePath("show_fields_create.csv")): &View{
				Header: NewHeader("show_fields_create", []string{"column1", "column2"}),
				FileInfo: &FileInfo{
					Path: GetTestFilePath("show_fields_create.csv"),
				},
			},
		},
		Result: "\n" + fmt.Sprintf("    Fields in %s", GetTestFilePath("show_fields_create.csv")) + "\n" + strings.Repeat("-", len(GetTestFilePath("show_fields_create.csv"))+18) + "\n" +
			"1. column1\n2. column2\n",
	},
	{
		Name: "ShowFields Cached Table",
		Expr: parser.ShowFields{
			Table: parser.Identifier{Literal: "table1"},
		},
		ViewCache: ViewMap{
			strings.ToUpper(GetTestFilePath("table1.csv")): &View{
				Header: NewHeader("table1", []string{"column1", "column2"}),
				FileInfo: &FileInfo{
					Path: GetTestFilePath("table1.csv"),
				},
			},
		},
		Result: "\n" + "    Fields in table1" + "\n" + "------------------------\n" +
			"1. column1\n2. column2\n",
	},
	{
		Name: "ShowFields Load From File",
		Expr: parser.ShowFields{
			Table: parser.Identifier{Literal: "table2"},
		},
		Result: "\n" + "    Fields in table2" + "\n" + "------------------------\n" +
			"1. column3\n2. column4\n",
	},
	{
		Name: "ShowFields Load Error",
		Expr: parser.ShowFields{
			Table: parser.Identifier{Literal: "notexist"},
		},
		Error: "[L:- C:-] file notexist does not exist",
	},
}

func TestShowFields(t *testing.T) {
	initFlag()
	flags := cmd.GetFlags()
	flags.Repository = TestDir

	for _, v := range showFieldsTests {
		ViewCache.Clean()
		if 0 < len(v.ViewCache) {
			ViewCache = v.ViewCache
		}

		var filter *Filter
		if v.Filter != nil {
			filter = v.Filter
		} else {
			filter = NewEmptyFilter()
		}

		result, err := ShowFields(v.Expr, filter)
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
		if result != v.Result {
			t.Errorf("%s: result = %s, want %s", v.Name, result, v.Result)
		}
	}
	ReleaseResources()
}
