package query

import (
	"reflect"
	"testing"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/ternary"
)

var selectQueryForCursorTest parser.SelectQuery = parser.SelectQuery{
	SelectEntity: parser.SelectEntity{
		SelectClause: parser.SelectClause{
			Select: "select",
			Fields: []parser.Expression{
				parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
				parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
			},
		},
		FromClause: parser.FromClause{
			From: "from",
			Tables: []parser.Expression{
				parser.Table{Object: parser.Identifier{Literal: "table1"}},
			},
		},
	},
}

var selectQueryForCursorQueryErrorTest parser.SelectQuery = parser.SelectQuery{
	SelectEntity: parser.SelectEntity{
		SelectClause: parser.SelectClause{
			Select: "select",
			Fields: []parser.Expression{
				parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}}},
				parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
			},
		},
		FromClause: parser.FromClause{
			From: "from",
			Tables: []parser.Expression{
				parser.Table{Object: parser.Identifier{Literal: "table1"}},
			},
		},
	},
}

var cursorMapAddTests = []struct {
	Name   string
	Key    string
	Query  parser.SelectQuery
	Result CursorMap
	Error  string
}{
	{
		Name:  "CursorMap Add",
		Key:   "cur",
		Query: selectQueryForCursorTest,
		Result: CursorMap{
			"CUR": &Cursor{
				name:  "cur",
				query: selectQueryForCursorTest,
			},
		},
	},
	{
		Name:  "CursorMap Add Already Exist",
		Key:   "cur",
		Query: parser.SelectQuery{},
		Error: "cursor cur already exists",
	},
}

func TestCursorMap_Add(t *testing.T) {
	cursors := CursorMap{}

	for _, v := range cursorMapAddTests {
		err := cursors.Add(v.Key, v.Query)
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
		if !reflect.DeepEqual(cursors, v.Result) {
			t.Errorf("%s: result = %s, want %s", v.Name, cursors, v.Result)
		}
	}
}

var cursorMapDisposeTests = []struct {
	Name   string
	Key    string
	Result CursorMap
}{
	{
		Name:   "CursorMap Dispose",
		Key:    "cur",
		Result: CursorMap{},
	},
}

func TestCursorMap_Dispose(t *testing.T) {
	cursors := CursorMap{
		"CUR": &Cursor{
			name:  "cur",
			query: selectQueryForCursorTest,
		},
	}

	for _, v := range cursorMapDisposeTests {
		cursors.Dispose(v.Key)
		if !reflect.DeepEqual(cursors, v.Result) {
			t.Errorf("%s: result = %s, want %s", v.Name, cursors, v.Result)
		}
	}
}

var cursorMapOpenTests = []struct {
	Name   string
	Key    string
	Result CursorMap
	Error  string
}{
	{
		Name: "CursorMap Open",
		Key:  "cur",
		Result: CursorMap{
			"CUR": &Cursor{
				name:  "cur",
				query: selectQueryForCursorTest,
				view: &View{
					Header: NewHeaderWithoutId("table1", []string{"column1", "column2"}),
					Records: []Record{
						NewRecordWithoutId([]parser.Primary{
							parser.NewString("1"),
							parser.NewString("str1"),
						}),
						NewRecordWithoutId([]parser.Primary{
							parser.NewString("2"),
							parser.NewString("str2"),
						}),
						NewRecordWithoutId([]parser.Primary{
							parser.NewString("3"),
							parser.NewString("str3"),
						}),
					},
					FileInfo: &FileInfo{
						Path:      GetTestFilePath("table1.csv"),
						Delimiter: ',',
						NoHeader:  false,
						Encoding:  cmd.UTF8,
						LineBreak: cmd.LF,
					},
				},
				index: -1,
			},
			"CUR2": &Cursor{
				name:  "cur2",
				query: selectQueryForCursorQueryErrorTest,
			},
		},
	},
	{
		Name:  "CursorMap Open Not Exists Error",
		Key:   "notexist",
		Error: "cursor notexist does not exist",
	},
	{
		Name:  "CursorMap Open Already Open Error",
		Key:   "cur",
		Error: "cursor cur is already open",
	},
	{
		Name:  "CursorMap Open Query Error",
		Key:   "cur2",
		Error: "field notexist does not exist",
	},
}

func TestCursorMap_Open(t *testing.T) {
	tf := cmd.GetFlags()
	tf.Repository = TestDir

	cursors := CursorMap{
		"CUR": &Cursor{
			name:  "cur",
			query: selectQueryForCursorTest,
		},
		"CUR2": &Cursor{
			name:  "cur2",
			query: selectQueryForCursorQueryErrorTest,
		},
	}

	for _, v := range cursorMapOpenTests {
		err := cursors.Open(v.Key)
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
		if !reflect.DeepEqual(cursors, v.Result) {
			t.Errorf("%s: result = %s, want %s", v.Name, cursors, v.Result)
		}
	}
}

var cursorMapCloseTests = []struct {
	Name   string
	Key    string
	Result CursorMap
	Error  string
}{
	{
		Name: "CursorMap Close",
		Key:  "cur",
		Result: CursorMap{
			"CUR": &Cursor{
				name:  "cur",
				query: selectQueryForCursorTest,
			},
		},
	},
	{
		Name:  "CursorMap Close Not Exist Error",
		Key:   "notexist",
		Error: "cursor notexist does not exist",
	},
}

func TestCursorMap_Close(t *testing.T) {
	tf := cmd.GetFlags()
	tf.Repository = TestDir

	cursors := CursorMap{
		"CUR": &Cursor{
			name:  "cur",
			query: selectQueryForCursorTest,
		},
	}
	cursors.Open("cur")

	for _, v := range cursorMapCloseTests {
		err := cursors.Close(v.Key)
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
		if !reflect.DeepEqual(cursors, v.Result) {
			t.Errorf("%s: result = %s, want %s", v.Name, cursors, v.Result)
		}
	}
}

var cursorMapFetchTests = []struct {
	Name     string
	Key      string
	Position int
	Number   int
	Result   []parser.Primary
	Error    string
}{
	{
		Name:     "CursorMap Fetch First Time",
		Key:      "cur",
		Position: parser.NEXT,
		Result: []parser.Primary{
			parser.NewString("1"),
			parser.NewString("str1"),
		},
	},
	{
		Name:     "CursorMap Fetch Second Time",
		Key:      "cur",
		Position: parser.NEXT,
		Result: []parser.Primary{
			parser.NewString("2"),
			parser.NewString("str2"),
		},
	},
	{
		Name:     "CursorMap Fetch Third Time",
		Key:      "cur",
		Position: parser.NEXT,
		Result: []parser.Primary{
			parser.NewString("3"),
			parser.NewString("str3"),
		},
	},
	{
		Name:     "CursorMap Fetch Fourth Time",
		Key:      "cur",
		Position: parser.NEXT,
		Result:   nil,
	},
	{
		Name:     "CursorMap Fetch First",
		Key:      "cur",
		Position: parser.FIRST,
		Result: []parser.Primary{
			parser.NewString("1"),
			parser.NewString("str1"),
		},
	},
	{
		Name:     "CursorMap Fetch Last",
		Key:      "cur",
		Position: parser.LAST,
		Result: []parser.Primary{
			parser.NewString("3"),
			parser.NewString("str3"),
		},
	},
	{
		Name:     "CursorMap Fetch Prior",
		Key:      "cur",
		Position: parser.PRIOR,
		Result: []parser.Primary{
			parser.NewString("2"),
			parser.NewString("str2"),
		},
	},
	{
		Name:     "CursorMap Fetch Absolute",
		Key:      "cur",
		Position: parser.ABSOLUTE,
		Number:   1,
		Result: []parser.Primary{
			parser.NewString("2"),
			parser.NewString("str2"),
		},
	},
	{
		Name:     "CursorMap Fetch Relative",
		Key:      "cur",
		Position: parser.RELATIVE,
		Number:   -1,
		Result: []parser.Primary{
			parser.NewString("1"),
			parser.NewString("str1"),
		},
	},
	{
		Name:     "CursorMap Fetch Prior to Last",
		Key:      "cur",
		Position: parser.ABSOLUTE,
		Number:   -2,
		Result:   nil,
	},
	{
		Name:     "CursorMap Fetch Later than Last",
		Key:      "cur",
		Position: parser.ABSOLUTE,
		Number:   100,
		Result:   nil,
	},
	{
		Name:     "CursorMap Fetch Not Exist Error",
		Key:      "notexist",
		Position: parser.NEXT,
		Error:    "cursor notexist does not exist",
	},
	{
		Name:     "CursorMap Fetch Closed Error",
		Key:      "cur2",
		Position: parser.NEXT,
		Error:    "cursor cur2 is closed",
	},
}

func TestCursorMap_Fetch(t *testing.T) {
	tf := cmd.GetFlags()
	tf.Repository = TestDir

	cursors := CursorMap{
		"CUR": &Cursor{
			name:  "cur",
			query: selectQueryForCursorTest,
		},
		"CUR2": &Cursor{
			name:  "cur2",
			query: selectQueryForCursorTest,
		},
	}
	cursors.Open("cur")

	for _, v := range cursorMapFetchTests {
		result, err := cursors.Fetch(v.Key, v.Position, v.Number)
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
			t.Errorf("%s: result = %s, want %s", v.Name, result, v.Result)
		}
	}
}

var cursorMapIsOpenTests = []struct {
	Name   string
	Key    string
	Result ternary.Value
	Error  string
}{
	{
		Name:   "CursorMap IsOpen TRUE",
		Key:    "cur",
		Result: ternary.TRUE,
	},
	{
		Name:   "CursorMap IsOpen FALSE",
		Key:    "cur2",
		Result: ternary.FALSE,
	},
	{
		Name:  "CursorMap IsOpen Not Exist Error",
		Key:   "notexist",
		Error: "cursor notexist does not exist",
	},
}

func TestCursorMap_IsOpen(t *testing.T) {
	tf := cmd.GetFlags()
	tf.Repository = TestDir

	cursors := CursorMap{
		"CUR": &Cursor{
			name:  "cur",
			query: selectQueryForCursorTest,
		},
		"CUR2": &Cursor{
			name:  "cur2",
			query: selectQueryForCursorTest,
		},
	}
	cursors.Open("cur")

	for _, v := range cursorMapIsOpenTests {
		result, err := cursors.IsOpen(v.Key)
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
}

var cursorMapIsInRangeTests = []struct {
	Name   string
	Key    string
	Index  int
	Result ternary.Value
	Error  string
}{
	{
		Name:   "CursorMap Is In Range UNKNOWN",
		Key:    "cur",
		Result: ternary.UNKNOWN,
	},
	{
		Name:   "CursorMap Is In Range TRUE",
		Key:    "cur2",
		Index:  1,
		Result: ternary.TRUE,
	},
	{
		Name:   "CursorMap Is In Range FALSE",
		Key:    "cur2",
		Index:  -1,
		Result: ternary.FALSE,
	},
	{
		Name:  "CursorMap Is In Range Not Open Error",
		Key:   "cur3",
		Error: "cursor cur3 is closed",
	},
	{
		Name:  "CursorMap Is In Range Not Exist Error",
		Key:   "notexist",
		Error: "cursor notexist does not exist",
	},
}

func TestCursorMap_IsInRange(t *testing.T) {
	tf := cmd.GetFlags()
	tf.Repository = TestDir

	cursors := CursorMap{
		"CUR": &Cursor{
			name:  "cur",
			query: selectQueryForCursorTest,
		},
		"CUR2": &Cursor{
			name:  "cur2",
			query: selectQueryForCursorTest,
		},
		"CUR3": &Cursor{
			name:  "cur3",
			query: selectQueryForCursorTest,
		},
	}
	cursors.Open("cur")
	cursors.Open("cur2")
	cursors.Fetch("cur2", parser.NEXT, 0)

	for _, v := range cursorMapIsInRangeTests {
		if 0 != v.Index {
			cursors["CUR2"].index = v.Index
		}
		result, err := cursors.IsInRange(v.Key)
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
}
