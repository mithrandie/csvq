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

var cursorMapListDeclareTests = []struct {
	Name   string
	Expr   parser.CursorDeclaration
	Result CursorMapList
	Error  string
}{
	{
		Name: "CursorMapList Declare",
		Expr: parser.CursorDeclaration{
			Cursor: parser.Identifier{Literal: "cur"},
			Query:  selectQueryForCursorTest,
		},
		Result: CursorMapList{
			{
				"CUR": &Cursor{
					query: selectQueryForCursorTest,
				},
			},
			{},
		},
	},
}

func TestCursorMapList_Declare(t *testing.T) {
	list := CursorMapList{
		{},
		{},
	}

	for _, v := range cursorMapListDeclareTests {
		err := list.Declare(v.Expr)
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
		if !reflect.DeepEqual(list, v.Result) {
			t.Errorf("%s: result = %s, want %s", v.Name, list, v.Result)
		}
	}
}

var cursorMapListDisposeTests = []struct {
	Name    string
	CurName parser.Identifier
	Result  CursorMapList
	Error   string
}{
	{
		Name:    "CursorMapList Dispose",
		CurName: parser.Identifier{Literal: "cur"},
		Result: CursorMapList{
			{},
			{},
		},
	},
	{
		Name:    "CursorMapList Dispose Undefined Error",
		CurName: parser.Identifier{Literal: "notexist"},
		Error:   "[L:- C:-] cursor notexist is undefined",
	},
}

func TestCursorMapList_Dispose(t *testing.T) {
	list := CursorMapList{
		{},
		{
			"CUR": &Cursor{
				query: selectQueryForCursorTest,
			},
		},
	}

	for _, v := range cursorMapListDisposeTests {
		err := list.Dispose(v.CurName)
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
		if !reflect.DeepEqual(list, v.Result) {
			t.Errorf("%s: result = %s, want %s", v.Name, list, v.Result)
		}
	}
}

var cursorMapListOpenTests = []struct {
	Name    string
	CurName parser.Identifier
	Result  CursorMapList
	Error   string
}{
	{
		Name:    "CursorMapList Open",
		CurName: parser.Identifier{Literal: "cur"},
		Result: CursorMapList{
			{},
			CursorMap{
				"CUR": &Cursor{
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
			},
		},
	},
	{
		Name:    "CursorMapList Open Undefined Error",
		CurName: parser.Identifier{Literal: "notexist"},
		Error:   "[L:- C:-] cursor notexist is undefined",
	},
	{
		Name:    "CursorMapList Open Open Error",
		CurName: parser.Identifier{Literal: "cur"},
		Error:   "[L:- C:-] cursor cur is already open",
	},
}

func TestCursorMapList_Open(t *testing.T) {
	initFlag()
	tf := cmd.GetFlags()
	tf.Repository = TestDir

	list := CursorMapList{
		{},
		{
			"CUR": &Cursor{
				query: selectQueryForCursorTest,
			},
		},
	}

	for _, v := range cursorMapListOpenTests {
		ViewCache.Clear()

		err := list.Open(v.CurName, NewEmptyFilter())
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
		if !reflect.DeepEqual(list, v.Result) {
			t.Errorf("%s: result = %s, want %s", v.Name, list, v.Result)
		}
	}
}

var cursorMapListCloseTests = []struct {
	Name    string
	CurName parser.Identifier
	Result  CursorMapList
	Error   string
}{
	{
		Name:    "CursorMapList Close",
		CurName: parser.Identifier{Literal: "cur"},
		Result: CursorMapList{
			{},
			CursorMap{
				"CUR": &Cursor{
					query: selectQueryForCursorTest,
				},
			},
		},
	},
	{
		Name:    "CursorMapList Close Undefined Error",
		CurName: parser.Identifier{Literal: "notexist"},
		Error:   "[L:- C:-] cursor notexist is undefined",
	},
}

func TestCursorMapList_Close(t *testing.T) {
	initFlag()
	tf := cmd.GetFlags()
	tf.Repository = TestDir

	list := CursorMapList{
		{},
		{
			"CUR": &Cursor{
				query: selectQueryForCursorTest,
			},
		},
	}

	ViewCache.Clear()
	list[1]["CUR"].Open(parser.Identifier{Literal: "cur"}, NewEmptyFilter())

	for _, v := range cursorMapListCloseTests {
		err := list.Close(v.CurName)
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
		if !reflect.DeepEqual(list, v.Result) {
			t.Errorf("%s: result = %s, want %s", v.Name, list, v.Result)
		}
	}
}

var cursorMapListFetchTests = []struct {
	Name     string
	CurName  parser.Identifier
	Position int
	Number   int
	Result   []parser.Primary
	Error    string
}{
	{
		Name:     "CursorMapList Fetch",
		CurName:  parser.Identifier{Literal: "cur"},
		Position: parser.NEXT,
		Result: []parser.Primary{
			parser.NewString("1"),
			parser.NewString("str1"),
		},
	},
	{
		Name:     "CursorMapList Fetch Undefined Error",
		CurName:  parser.Identifier{Literal: "notexist"},
		Position: parser.NEXT,
		Error:    "[L:- C:-] cursor notexist is undefined",
	},
	{
		Name:     "CursorMapList Fetch Closed Error",
		CurName:  parser.Identifier{Literal: "cur2"},
		Position: parser.NEXT,
		Error:    "[L:- C:-] cursor cur2 is closed",
	},
}

func TestCursorMapList_Fetch(t *testing.T) {
	initFlag()
	tf := cmd.GetFlags()
	tf.Repository = TestDir

	list := CursorMapList{
		{},
		{
			"CUR": &Cursor{
				query: selectQueryForCursorTest,
			},
			"CUR2": &Cursor{
				query: selectQueryForCursorTest,
			},
		},
	}

	ViewCache.Clear()
	list[1]["CUR"].Open(parser.Identifier{Literal: "cur"}, NewEmptyFilter())

	for _, v := range cursorMapListFetchTests {
		result, err := list.Fetch(v.CurName, v.Position, v.Number)
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

var cursorMapListIsOpenTests = []struct {
	Name    string
	CurName parser.Identifier
	Result  ternary.Value
	Error   string
}{
	{
		Name:    "CursorMapList IsOpen",
		CurName: parser.Identifier{Literal: "cur"},
		Result:  ternary.TRUE,
	},
	{
		Name:    "CursorMapList IsOpen Undefined Error",
		CurName: parser.Identifier{Literal: "notexist"},
		Error:   "[L:- C:-] cursor notexist is undefined",
	},
}

func TestCursorMapList_IsOpen(t *testing.T) {
	initFlag()
	tf := cmd.GetFlags()
	tf.Repository = TestDir

	list := CursorMapList{
		{},
		{
			"CUR": &Cursor{
				query: selectQueryForCursorTest,
			},
		},
	}

	ViewCache.Clear()
	list[1]["CUR"].Open(parser.Identifier{Literal: "cur"}, NewEmptyFilter())

	for _, v := range cursorMapListIsOpenTests {
		result, err := list.IsOpen(v.CurName)
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

var cursorMapListIsInRangeTests = []struct {
	Name    string
	CurName parser.Identifier
	Index   int
	Result  ternary.Value
	Error   string
}{
	{
		Name:    "CursorMapList Is In Range UNKNOWN",
		CurName: parser.Identifier{Literal: "cur"},
		Result:  ternary.UNKNOWN,
	},
	{
		Name:    "CursorMap Is In Range Not Open Error",
		CurName: parser.Identifier{Literal: "cur2"},
		Error:   "[L:- C:-] cursor cur2 is closed",
	},
	{
		Name:    "CursorMap Is In Range Undefined Error",
		CurName: parser.Identifier{Literal: "notexist"},
		Error:   "[L:- C:-] cursor notexist is undefined",
	},
}

func TestCursorMapList_IsInRange(t *testing.T) {
	initFlag()
	tf := cmd.GetFlags()
	tf.Repository = TestDir

	list := CursorMapList{
		{},
		{
			"CUR": &Cursor{
				query: selectQueryForCursorTest,
			},
			"CUR2": &Cursor{
				query: selectQueryForCursorTest,
			},
		},
	}

	ViewCache.Clear()
	list[1]["CUR"].Open(parser.Identifier{Literal: "cur"}, NewEmptyFilter())

	for _, v := range cursorMapListIsInRangeTests {
		result, err := list.IsInRange(v.CurName)
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

var cursorMapDeclareTests = []struct {
	Name   string
	Expr   parser.CursorDeclaration
	Result CursorMap
	Error  string
}{
	{
		Name: "CursorMap Declare",
		Expr: parser.CursorDeclaration{
			Cursor: parser.Identifier{Literal: "cur"},
			Query:  selectQueryForCursorTest,
		},
		Result: CursorMap{
			"CUR": &Cursor{
				query: selectQueryForCursorTest,
			},
		},
	},
	{
		Name: "CursorMap Declare Redeclaration Error",
		Expr: parser.CursorDeclaration{
			Cursor: parser.Identifier{Literal: "cur"},
			Query:  parser.SelectQuery{},
		},
		Error: "[L:- C:-] cursor cur is redeclared",
	},
}

func TestCursorMap_Declare(t *testing.T) {
	cursors := CursorMap{}

	for _, v := range cursorMapDeclareTests {
		err := cursors.Declare(v.Expr)
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
	Name    string
	CurName parser.Identifier
	Result  CursorMap
	Error   string
}{
	{
		Name:    "CursorMap Dispose",
		CurName: parser.Identifier{Literal: "cur"},
		Result:  CursorMap{},
	},
	{
		Name:    "CursorMap Dispose Undefined Error",
		CurName: parser.Identifier{Literal: "notexist"},
		Error:   "[L:- C:-] cursor notexist is undefined",
	},
}

func TestCursorMap_Dispose(t *testing.T) {
	cursors := CursorMap{
		"CUR": &Cursor{
			query: selectQueryForCursorTest,
		},
	}

	for _, v := range cursorMapDisposeTests {
		err := cursors.Dispose(v.CurName)
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

var cursorMapOpenTests = []struct {
	Name    string
	CurName parser.Identifier
	Result  CursorMap
	Error   string
}{
	{
		Name:    "CursorMap Open",
		CurName: parser.Identifier{Literal: "cur"},
		Result: CursorMap{
			"CUR": &Cursor{
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
				query: selectQueryForCursorQueryErrorTest,
			},
		},
	},
	{
		Name:    "CursorMap Open Undefined Error",
		CurName: parser.Identifier{Literal: "notexist"},
		Error:   "[L:- C:-] cursor notexist is undefined",
	},
	{
		Name:    "CursorMap Open Open Error",
		CurName: parser.Identifier{Literal: "cur"},
		Error:   "[L:- C:-] cursor cur is already open",
	},
	{
		Name:    "CursorMap Open Query Error",
		CurName: parser.Identifier{Literal: "cur2"},
		Error:   "[L:- C:-] field notexist does not exist",
	},
}

func TestCursorMap_Open(t *testing.T) {
	initFlag()
	tf := cmd.GetFlags()
	tf.Repository = TestDir

	cursors := CursorMap{
		"CUR": &Cursor{
			query: selectQueryForCursorTest,
		},
		"CUR2": &Cursor{
			query: selectQueryForCursorQueryErrorTest,
		},
	}

	for _, v := range cursorMapOpenTests {
		ViewCache.Clear()
		err := cursors.Open(v.CurName, NewEmptyFilter())
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
	Name    string
	CurName parser.Identifier
	Result  CursorMap
	Error   string
}{
	{
		Name:    "CursorMap Close",
		CurName: parser.Identifier{Literal: "cur"},
		Result: CursorMap{
			"CUR": &Cursor{
				query: selectQueryForCursorTest,
			},
		},
	},
	{
		Name:    "CursorMap Close Undefined Error",
		CurName: parser.Identifier{Literal: "notexist"},
		Error:   "[L:- C:-] cursor notexist is undefined",
	},
}

func TestCursorMap_Close(t *testing.T) {
	tf := cmd.GetFlags()
	tf.Repository = TestDir

	cursors := CursorMap{
		"CUR": &Cursor{
			query: selectQueryForCursorTest,
		},
	}
	ViewCache.Clear()
	cursors.Open(parser.Identifier{Literal: "cur"}, NewEmptyFilter())

	for _, v := range cursorMapCloseTests {
		err := cursors.Close(v.CurName)
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
	CurName  parser.Identifier
	Position int
	Number   int
	Result   []parser.Primary
	Error    string
}{
	{
		Name:     "CursorMap Fetch First Time",
		CurName:  parser.Identifier{Literal: "cur"},
		Position: parser.NEXT,
		Result: []parser.Primary{
			parser.NewString("1"),
			parser.NewString("str1"),
		},
	},
	{
		Name:     "CursorMap Fetch Second Time",
		CurName:  parser.Identifier{Literal: "cur"},
		Position: parser.NEXT,
		Result: []parser.Primary{
			parser.NewString("2"),
			parser.NewString("str2"),
		},
	},
	{
		Name:     "CursorMap Fetch Third Time",
		CurName:  parser.Identifier{Literal: "cur"},
		Position: parser.NEXT,
		Result: []parser.Primary{
			parser.NewString("3"),
			parser.NewString("str3"),
		},
	},
	{
		Name:     "CursorMap Fetch Fourth Time",
		CurName:  parser.Identifier{Literal: "cur"},
		Position: parser.NEXT,
		Result:   nil,
	},
	{
		Name:     "CursorMap Fetch First",
		CurName:  parser.Identifier{Literal: "cur"},
		Position: parser.FIRST,
		Result: []parser.Primary{
			parser.NewString("1"),
			parser.NewString("str1"),
		},
	},
	{
		Name:     "CursorMap Fetch Last",
		CurName:  parser.Identifier{Literal: "cur"},
		Position: parser.LAST,
		Result: []parser.Primary{
			parser.NewString("3"),
			parser.NewString("str3"),
		},
	},
	{
		Name:     "CursorMap Fetch Prior",
		CurName:  parser.Identifier{Literal: "cur"},
		Position: parser.PRIOR,
		Result: []parser.Primary{
			parser.NewString("2"),
			parser.NewString("str2"),
		},
	},
	{
		Name:     "CursorMap Fetch Absolute",
		CurName:  parser.Identifier{Literal: "cur"},
		Position: parser.ABSOLUTE,
		Number:   1,
		Result: []parser.Primary{
			parser.NewString("2"),
			parser.NewString("str2"),
		},
	},
	{
		Name:     "CursorMap Fetch Relative",
		CurName:  parser.Identifier{Literal: "cur"},
		Position: parser.RELATIVE,
		Number:   -1,
		Result: []parser.Primary{
			parser.NewString("1"),
			parser.NewString("str1"),
		},
	},
	{
		Name:     "CursorMap Fetch Prior to Last",
		CurName:  parser.Identifier{Literal: "cur"},
		Position: parser.ABSOLUTE,
		Number:   -2,
		Result:   nil,
	},
	{
		Name:     "CursorMap Fetch Later than Last",
		CurName:  parser.Identifier{Literal: "cur"},
		Position: parser.ABSOLUTE,
		Number:   100,
		Result:   nil,
	},
	{
		Name:     "CursorMap Fetch Undefined Error",
		CurName:  parser.Identifier{Literal: "notexist"},
		Position: parser.NEXT,
		Error:    "[L:- C:-] cursor notexist is undefined",
	},
	{
		Name:     "CursorMap Fetch Closed Error",
		CurName:  parser.Identifier{Literal: "cur2"},
		Position: parser.NEXT,
		Error:    "[L:- C:-] cursor cur2 is closed",
	},
}

func TestCursorMap_Fetch(t *testing.T) {
	tf := cmd.GetFlags()
	tf.Repository = TestDir

	cursors := CursorMap{
		"CUR": &Cursor{
			query: selectQueryForCursorTest,
		},
		"CUR2": &Cursor{
			query: selectQueryForCursorTest,
		},
	}
	ViewCache.Clear()
	cursors.Open(parser.Identifier{Literal: "cur"}, NewEmptyFilter())

	for _, v := range cursorMapFetchTests {
		result, err := cursors.Fetch(v.CurName, v.Position, v.Number)
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
	Name    string
	CurName parser.Identifier
	Result  ternary.Value
	Error   string
}{
	{
		Name:    "CursorMap IsOpen TRUE",
		CurName: parser.Identifier{Literal: "cur"},
		Result:  ternary.TRUE,
	},
	{
		Name:    "CursorMap IsOpen FALSE",
		CurName: parser.Identifier{Literal: "cur2"},
		Result:  ternary.FALSE,
	},
	{
		Name:    "CursorMap IsOpen Undefined Error",
		CurName: parser.Identifier{Literal: "notexist"},
		Error:   "[L:- C:-] cursor notexist is undefined",
	},
}

func TestCursorMap_IsOpen(t *testing.T) {
	tf := cmd.GetFlags()
	tf.Repository = TestDir

	cursors := CursorMap{
		"CUR": &Cursor{
			query: selectQueryForCursorTest,
		},
		"CUR2": &Cursor{
			query: selectQueryForCursorTest,
		},
	}
	ViewCache.Clear()
	cursors.Open(parser.Identifier{Literal: "cur"}, NewEmptyFilter())

	for _, v := range cursorMapIsOpenTests {
		result, err := cursors.IsOpen(v.CurName)
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
	Name    string
	CurName parser.Identifier
	Index   int
	Result  ternary.Value
	Error   string
}{
	{
		Name:    "CursorMap Is In Range UNKNOWN",
		CurName: parser.Identifier{Literal: "cur"},
		Result:  ternary.UNKNOWN,
	},
	{
		Name:    "CursorMap Is In Range TRUE",
		CurName: parser.Identifier{Literal: "cur2"},
		Index:   1,
		Result:  ternary.TRUE,
	},
	{
		Name:    "CursorMap Is In Range FALSE",
		CurName: parser.Identifier{Literal: "cur2"},
		Index:   -1,
		Result:  ternary.FALSE,
	},
	{
		Name:    "CursorMap Is In Range Not Open Error",
		CurName: parser.Identifier{Literal: "cur3"},
		Error:   "[L:- C:-] cursor cur3 is closed",
	},
	{
		Name:    "CursorMap Is In Range Undefined Error",
		CurName: parser.Identifier{Literal: "notexist"},
		Error:   "[L:- C:-] cursor notexist is undefined",
	},
}

func TestCursorMap_IsInRange(t *testing.T) {
	tf := cmd.GetFlags()
	tf.Repository = TestDir

	cursors := CursorMap{
		"CUR": &Cursor{
			query: selectQueryForCursorTest,
		},
		"CUR2": &Cursor{
			query: selectQueryForCursorTest,
		},
		"CUR3": &Cursor{
			query: selectQueryForCursorTest,
		},
	}
	ViewCache.Clear()
	cursors.Open(parser.Identifier{Literal: "cur"}, NewEmptyFilter())
	ViewCache.Clear()
	cursors.Open(parser.Identifier{Literal: "cur2"}, NewEmptyFilter())
	cursors.Fetch(parser.Identifier{Literal: "cur2"}, parser.NEXT, 0)

	for _, v := range cursorMapIsInRangeTests {
		if 0 != v.Index {
			cursors["CUR2"].index = v.Index
		}
		result, err := cursors.IsInRange(v.CurName)
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
