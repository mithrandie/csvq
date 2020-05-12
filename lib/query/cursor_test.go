package query

import (
	"context"
	"reflect"
	"sync"
	"testing"

	"github.com/mithrandie/go-text"

	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/value"

	"github.com/mithrandie/ternary"
)

var selectQueryForCursorTest = parser.SelectQuery{
	SelectEntity: parser.SelectEntity{
		SelectClause: parser.SelectClause{
			Fields: []parser.QueryExpression{
				parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
				parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
			},
		},
		FromClause: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{Object: parser.Identifier{Literal: "table1"}},
			},
		},
	},
}

var selectQueryForCursorQueryErrorTest = parser.SelectQuery{
	SelectEntity: parser.SelectEntity{
		SelectClause: parser.SelectClause{
			Fields: []parser.QueryExpression{
				parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}}},
				parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
			},
		},
		FromClause: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{Object: parser.Identifier{Literal: "table1"}},
			},
		},
	},
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
		Result: GenerateCursorMap([]*Cursor{
			{
				name:  "cur",
				query: selectQueryForCursorTest,
				mtx:   &sync.Mutex{},
			},
		}),
	},
	{
		Name: "CursorMap Declare for Statement",
		Expr: parser.CursorDeclaration{
			Cursor:    parser.Identifier{Literal: "stmtcur"},
			Statement: parser.Identifier{Literal: "stmt"},
		},
		Result: GenerateCursorMap([]*Cursor{
			{
				name:  "cur",
				query: selectQueryForCursorTest,
				mtx:   &sync.Mutex{},
			},
			{
				name:      "stmtcur",
				statement: parser.Identifier{Literal: "stmt"},
				mtx:       &sync.Mutex{},
			},
		}),
	},
	{
		Name: "CursorMap Declare Redeclaration Error",
		Expr: parser.CursorDeclaration{
			Cursor: parser.Identifier{Literal: "cur"},
			Query:  parser.SelectQuery{},
		},
		Error: "cursor cur is redeclared",
	},
}

func TestCursorMap_Declare(t *testing.T) {
	cursors := NewCursorMap()

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
		if !SyncMapEqual(cursors, v.Result) {
			t.Errorf("%s: result = %v, want %v", v.Name, cursors, v.Result)
		}
	}
}

var cursorMapAddPseudoCursorTests = []struct {
	Name   string
	Cursor parser.Identifier
	Values []value.Primary
	Result CursorMap
	Error  string
}{
	{
		Name:   "CursorMap AddPseudoCursor",
		Cursor: parser.Identifier{Literal: "pcur"},
		Values: []value.Primary{
			value.NewInteger(1),
			value.NewInteger(2),
		},
		Result: GenerateCursorMap([]*Cursor{
			{
				name: "pcur",
				view: &View{
					Header: NewHeader("", []string{"c1"}),
					RecordSet: RecordSet{
						NewRecord([]value.Primary{value.NewInteger(1)}),
						NewRecord([]value.Primary{value.NewInteger(2)}),
					},
				},
				index:    -1,
				isPseudo: true,
				mtx:      &sync.Mutex{},
			},
		}),
	},
	{
		Name:   "CursorMap AddPseudoCursor Redeclaration Error",
		Cursor: parser.Identifier{Literal: "pcur"},
		Values: []value.Primary{},
		Error:  "cursor pcur is redeclared",
	},
}

func TestCursorMap_AddPseudoCursor(t *testing.T) {
	cursors := NewCursorMap()
	for _, v := range cursorMapAddPseudoCursorTests {
		err := cursors.AddPseudoCursor(v.Cursor, v.Values)
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
		if !SyncMapEqual(cursors, v.Result) {
			t.Errorf("%s: result = %v, want %v", v.Name, cursors, v.Result)
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
		Result: GenerateCursorMap([]*Cursor{
			{
				name: "pcur",
				view: &View{
					Header: NewHeader("", []string{"c1"}),
					RecordSet: RecordSet{
						NewRecord([]value.Primary{value.NewInteger(1)}),
						NewRecord([]value.Primary{value.NewInteger(2)}),
					},
				},
				index:    -1,
				isPseudo: true,
				mtx:      &sync.Mutex{},
			},
		}),
	},
	{
		Name:    "CursorMap Dispose Undeclared Error",
		CurName: parser.Identifier{Literal: "notexist"},
		Error:   "undeclared cursor",
	},
	{
		Name:    "CursorMap Dispose Rseudo Cursor Error",
		CurName: parser.Identifier{Literal: "pcur"},
		Error:   "unpermmitted pseudo cursor usage",
	},
}

func TestCursorMap_Dispose(t *testing.T) {
	cursors := GenerateCursorMap([]*Cursor{
		{
			name:  "cur",
			query: selectQueryForCursorTest,
			mtx:   &sync.Mutex{},
		},
	})
	_ = cursors.AddPseudoCursor(
		parser.Identifier{Literal: "pcur"},
		[]value.Primary{
			value.NewInteger(1),
			value.NewInteger(2),
		},
	)

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
		if !SyncMapEqual(cursors, v.Result) {
			t.Errorf("%s: result = %v, want %v", v.Name, cursors, v.Result)
		}
	}
}

var cursorMapOpenTests = []struct {
	Name      string
	CurName   parser.Identifier
	CurValues []parser.ReplaceValue
	Result    CursorMap
	Error     string
}{
	{
		Name:    "CursorMap Open",
		CurName: parser.Identifier{Literal: "cur"},
		Result: GenerateCursorMap([]*Cursor{
			{
				name:  "cur",
				query: selectQueryForCursorTest,
				view: &View{
					Header: NewHeader("table1", []string{"column1", "column2"}),
					RecordSet: []Record{
						NewRecord([]value.Primary{
							value.NewString("1"),
							value.NewString("str1"),
						}),
						NewRecord([]value.Primary{
							value.NewString("2"),
							value.NewString("str2"),
						}),
						NewRecord([]value.Primary{
							value.NewString("3"),
							value.NewString("str3"),
						}),
					},
					FileInfo: &FileInfo{
						Path:      GetTestFilePath("table1.csv"),
						Delimiter: ',',
						NoHeader:  false,
						Encoding:  text.UTF8,
						LineBreak: text.LF,
					},
				},
				index: -1,
				mtx:   &sync.Mutex{},
			},
			{
				name:  "cur2",
				query: selectQueryForCursorQueryErrorTest,
				mtx:   &sync.Mutex{},
			},
			{
				name: "pcur",
				view: &View{
					Header: NewHeader("", []string{"c1"}),
					RecordSet: RecordSet{
						NewRecord([]value.Primary{value.NewInteger(1)}),
						NewRecord([]value.Primary{value.NewInteger(2)}),
					},
				},
				index:    -1,
				isPseudo: true,
				mtx:      &sync.Mutex{},
			},
			{
				name:      "stmt",
				statement: parser.Identifier{Literal: "stmt"},
				mtx:       &sync.Mutex{},
			},
			{
				name:      "not_exist_stmt",
				statement: parser.Identifier{Literal: "not_exist_stmt"},
				mtx:       &sync.Mutex{},
			},
			{
				name:      "invalid_stmt",
				statement: parser.Identifier{Literal: "invalid_stmt"},
				mtx:       &sync.Mutex{},
			},
			{
				name:      "invalid_stmt2",
				statement: parser.Identifier{Literal: "invalid_stmt2"},
				mtx:       &sync.Mutex{},
			},
		}),
	},
	{
		Name:    "CursorMap Open Statement",
		CurName: parser.Identifier{Literal: "stmt"},
		CurValues: []parser.ReplaceValue{
			{Value: parser.NewIntegerValueFromString("2")},
		},
		Result: GenerateCursorMap([]*Cursor{
			{
				name:      "stmt",
				statement: parser.Identifier{Literal: "stmt"},
				view: &View{
					Header: NewHeader("table1", []string{"column1", "column2"}),
					RecordSet: []Record{
						NewRecord([]value.Primary{
							value.NewString("2"),
							value.NewString("str2"),
						}),
					},
					FileInfo: &FileInfo{
						Path:      GetTestFilePath("table1.csv"),
						Delimiter: ',',
						NoHeader:  false,
						Encoding:  text.UTF8,
						LineBreak: text.LF,
					},
				},
				index: -1,
				mtx:   &sync.Mutex{},
			},
			{
				name:  "cur",
				query: selectQueryForCursorTest,
				view: &View{
					Header: NewHeader("table1", []string{"column1", "column2"}),
					RecordSet: []Record{
						NewRecord([]value.Primary{
							value.NewString("1"),
							value.NewString("str1"),
						}),
						NewRecord([]value.Primary{
							value.NewString("2"),
							value.NewString("str2"),
						}),
						NewRecord([]value.Primary{
							value.NewString("3"),
							value.NewString("str3"),
						}),
					},
					FileInfo: &FileInfo{
						Path:      GetTestFilePath("table1.csv"),
						Delimiter: ',',
						NoHeader:  false,
						Encoding:  text.UTF8,
						LineBreak: text.LF,
					},
				},
				index: -1,
				mtx:   &sync.Mutex{},
			},
			{
				name:  "cur2",
				query: selectQueryForCursorQueryErrorTest,
				mtx:   &sync.Mutex{},
			},
			{
				name: "pcur",
				view: &View{
					Header: NewHeader("", []string{"c1"}),
					RecordSet: RecordSet{
						NewRecord([]value.Primary{value.NewInteger(1)}),
						NewRecord([]value.Primary{value.NewInteger(2)}),
					},
				},
				index:    -1,
				isPseudo: true,
				mtx:      &sync.Mutex{},
			},
			{
				name:      "not_exist_stmt",
				statement: parser.Identifier{Literal: "not_exist_stmt"},
				mtx:       &sync.Mutex{},
			},
			{
				name:      "invalid_stmt",
				statement: parser.Identifier{Literal: "invalid_stmt"},
				mtx:       &sync.Mutex{},
			},
			{
				name:      "invalid_stmt2",
				statement: parser.Identifier{Literal: "invalid_stmt2"},
				mtx:       &sync.Mutex{},
			},
		}),
	},
	{
		Name:    "CursorMap Open Undeclared Error",
		CurName: parser.Identifier{Literal: "notexist"},
		Error:   "undeclared cursor",
	},
	{
		Name:    "CursorMap Open Open Error",
		CurName: parser.Identifier{Literal: "cur"},
		Error:   "cursor cur is already open",
	},
	{
		Name:    "CursorMap Open Query Error",
		CurName: parser.Identifier{Literal: "cur2"},
		Error:   "field notexist does not exist",
	},
	{
		Name:    "CursorMap Open Rseudo Cursor Error",
		CurName: parser.Identifier{Literal: "pcur"},
		Error:   "cursor pcur is a pseudo cursor",
	},
	{
		Name:    "CursorMap Open Unprepared Statement",
		CurName: parser.Identifier{Literal: "not_exist_stmt"},
		Error:   "statement not_exist_stmt does not exist",
	},
	{
		Name:    "CursorMap Open Not Select Query Error",
		CurName: parser.Identifier{Literal: "invalid_stmt"},
		Error:   "invalid cursor statement: invalid_stmt",
	},
	{
		Name:    "CursorMap Open Multiple Statements Error",
		CurName: parser.Identifier{Literal: "invalid_stmt2"},
		Error:   "invalid cursor statement: invalid_stmt2",
	},
}

func TestCursorMap_Open(t *testing.T) {
	defer func() {
		TestTx.PreparedStatements = NewPreparedStatementMap()
		_ = TestTx.cachedViews.Clean(TestTx.FileContainer)
		initFlag(TestTx.Flags)
	}()

	scope := NewReferenceScope(TestTx)
	TestTx.Flags.Repository = TestDir
	_ = TestTx.PreparedStatements.Prepare(scope.Tx.Flags, parser.StatementPreparation{
		Name:      parser.Identifier{Literal: "stmt"},
		Statement: value.NewString("select * from table1 where column1 = ?"),
	})
	_ = TestTx.PreparedStatements.Prepare(scope.Tx.Flags, parser.StatementPreparation{
		Name:      parser.Identifier{Literal: "invalid_stmt"},
		Statement: value.NewString("insert into table1 values (?, ?)"),
	})
	_ = TestTx.PreparedStatements.Prepare(scope.Tx.Flags, parser.StatementPreparation{
		Name:      parser.Identifier{Literal: "invalid_stmt2"},
		Statement: value.NewString("select 1; insert into table1 values (?, ?);"),
	})

	scope.blocks[0].cursors = GenerateCursorMap([]*Cursor{
		{
			name:  "cur",
			query: selectQueryForCursorTest,
			mtx:   &sync.Mutex{},
		},
		{
			name:  "cur2",
			query: selectQueryForCursorQueryErrorTest,
			mtx:   &sync.Mutex{},
		},
		{
			name:      "stmt",
			statement: parser.Identifier{Literal: "stmt"},
			mtx:       &sync.Mutex{},
		},
		{
			name:      "not_exist_stmt",
			statement: parser.Identifier{Literal: "not_exist_stmt"},
			mtx:       &sync.Mutex{},
		},
		{
			name:      "invalid_stmt",
			statement: parser.Identifier{Literal: "invalid_stmt"},
			mtx:       &sync.Mutex{},
		},
		{
			name:      "invalid_stmt2",
			statement: parser.Identifier{Literal: "invalid_stmt2"},
			mtx:       &sync.Mutex{},
		},
	})
	_ = scope.blocks[0].cursors.AddPseudoCursor(
		parser.Identifier{Literal: "pcur"},
		[]value.Primary{
			value.NewInteger(1),
			value.NewInteger(2),
		},
	)

	ctx := context.Background()
	for _, v := range cursorMapOpenTests {
		_ = TestTx.cachedViews.Clean(TestTx.FileContainer)
		err := scope.blocks[0].cursors.Open(ctx, scope, v.CurName, v.CurValues)
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
		if !SyncMapEqual(scope.blocks[0].cursors, v.Result) {
			t.Errorf("%s: result = %v, want %v", v.Name, scope.blocks[0].cursors, v.Result)
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
		Result: GenerateCursorMap([]*Cursor{
			{
				name:  "cur",
				query: selectQueryForCursorTest,
				mtx:   &sync.Mutex{},
			},
			{
				name: "pcur",
				view: &View{
					Header: NewHeader("", []string{"c1"}),
					RecordSet: RecordSet{
						NewRecord([]value.Primary{value.NewInteger(1)}),
						NewRecord([]value.Primary{value.NewInteger(2)}),
					},
				},
				index:    -1,
				isPseudo: true,
				mtx:      &sync.Mutex{},
			},
		}),
	},
	{
		Name:    "CursorMap Close Rseudo Cursor Error",
		CurName: parser.Identifier{Literal: "pcur"},
		Error:   "cursor pcur is a pseudo cursor",
	},
	{
		Name:    "CursorMap Close Undeclared Error",
		CurName: parser.Identifier{Literal: "notexist"},
		Error:   "undeclared cursor",
	},
}

func TestCursorMap_Close(t *testing.T) {
	defer func() {
		_ = TestTx.cachedViews.Clean(TestTx.FileContainer)
		initFlag(TestTx.Flags)
	}()

	TestTx.Flags.Repository = TestDir

	cursors := GenerateCursorMap([]*Cursor{
		{
			name:  "cur",
			query: selectQueryForCursorTest,
			mtx:   &sync.Mutex{},
		},
	})
	_ = cursors.AddPseudoCursor(
		parser.Identifier{Literal: "pcur"},
		[]value.Primary{
			value.NewInteger(1),
			value.NewInteger(2),
		},
	)
	scope := NewReferenceScope(TestTx)
	ctx := context.Background()
	_ = cursors.Open(ctx, scope, parser.Identifier{Literal: "cur"}, nil)

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
		if !SyncMapEqual(cursors, v.Result) {
			t.Errorf("%s: result = %v, want %v", v.Name, cursors, v.Result)
		}
	}
}

var cursorMapFetchTests = []struct {
	Name     string
	CurName  parser.Identifier
	Position int
	Number   int
	Result   []value.Primary
	Error    string
}{
	{
		Name:     "CursorMap Fetch First Time",
		CurName:  parser.Identifier{Literal: "cur"},
		Position: parser.NEXT,
		Result: []value.Primary{
			value.NewString("1"),
			value.NewString("str1"),
		},
	},
	{
		Name:     "CursorMap Fetch Second Time",
		CurName:  parser.Identifier{Literal: "cur"},
		Position: parser.NEXT,
		Result: []value.Primary{
			value.NewString("2"),
			value.NewString("str2"),
		},
	},
	{
		Name:     "CursorMap Fetch Third Time",
		CurName:  parser.Identifier{Literal: "cur"},
		Position: parser.NEXT,
		Result: []value.Primary{
			value.NewString("3"),
			value.NewString("str3"),
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
		Result: []value.Primary{
			value.NewString("1"),
			value.NewString("str1"),
		},
	},
	{
		Name:     "CursorMap Fetch Last",
		CurName:  parser.Identifier{Literal: "cur"},
		Position: parser.LAST,
		Result: []value.Primary{
			value.NewString("3"),
			value.NewString("str3"),
		},
	},
	{
		Name:     "CursorMap Fetch Prior",
		CurName:  parser.Identifier{Literal: "cur"},
		Position: parser.PRIOR,
		Result: []value.Primary{
			value.NewString("2"),
			value.NewString("str2"),
		},
	},
	{
		Name:     "CursorMap Fetch Absolute",
		CurName:  parser.Identifier{Literal: "cur"},
		Position: parser.ABSOLUTE,
		Number:   1,
		Result: []value.Primary{
			value.NewString("2"),
			value.NewString("str2"),
		},
	},
	{
		Name:     "CursorMap Fetch Relative",
		CurName:  parser.Identifier{Literal: "cur"},
		Position: parser.RELATIVE,
		Number:   -1,
		Result: []value.Primary{
			value.NewString("1"),
			value.NewString("str1"),
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
		Name:     "CursorMap Fetch Undeclared Error",
		CurName:  parser.Identifier{Literal: "notexist"},
		Position: parser.NEXT,
		Error:    "undeclared cursor",
	},
	{
		Name:     "CursorMap Fetch Closed Error",
		CurName:  parser.Identifier{Literal: "cur2"},
		Position: parser.NEXT,
		Error:    "cursor cur2 is closed",
	},
}

func TestCursorMap_Fetch(t *testing.T) {
	defer func() {
		_ = TestTx.cachedViews.Clean(TestTx.FileContainer)
		initFlag(TestTx.Flags)
	}()

	TestTx.Flags.Repository = TestDir

	cursors := GenerateCursorMap([]*Cursor{
		{
			name:  "cur",
			query: selectQueryForCursorTest,
			mtx:   &sync.Mutex{},
		},
		{
			name:  "cur2",
			query: selectQueryForCursorTest,
			mtx:   &sync.Mutex{},
		},
	})
	scope := NewReferenceScope(TestTx)
	ctx := context.Background()
	_ = cursors.Open(ctx, scope, parser.Identifier{Literal: "cur"}, nil)

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
		Name:    "CursorMap IsOpen Undeclared Error",
		CurName: parser.Identifier{Literal: "notexist"},
		Error:   "undeclared cursor",
	},
}

func TestCursorMap_IsOpen(t *testing.T) {
	defer func() {
		_ = TestTx.cachedViews.Clean(TestTx.FileContainer)
		initFlag(TestTx.Flags)
	}()

	TestTx.Flags.Repository = TestDir

	cursors := GenerateCursorMap([]*Cursor{
		{
			name:  "cur",
			query: selectQueryForCursorTest,
			mtx:   &sync.Mutex{},
		},
		{
			name:  "cur2",
			query: selectQueryForCursorTest,
			mtx:   &sync.Mutex{},
		},
	})
	scope := NewReferenceScope(TestTx)
	ctx := context.Background()
	_ = cursors.Open(ctx, scope, parser.Identifier{Literal: "cur"}, nil)

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
		Error:   "cursor cur3 is closed",
	},
	{
		Name:    "CursorMap Is In Range Undeclared Error",
		CurName: parser.Identifier{Literal: "notexist"},
		Error:   "undeclared cursor",
	},
}

func TestCursorMap_IsInRange(t *testing.T) {
	defer func() {
		_ = TestTx.cachedViews.Clean(TestTx.FileContainer)
		initFlag(TestTx.Flags)
	}()

	TestTx.Flags.Repository = TestDir

	cursors := GenerateCursorMap([]*Cursor{
		{
			name:  "cur",
			query: selectQueryForCursorTest,
			mtx:   &sync.Mutex{},
		},
		{
			name:  "cur2",
			query: selectQueryForCursorTest,
			mtx:   &sync.Mutex{},
		},
		{
			name:  "cur3",
			query: selectQueryForCursorTest,
			mtx:   &sync.Mutex{},
		},
	})
	scope := NewReferenceScope(TestTx)
	ctx := context.Background()
	_ = TestTx.cachedViews.Clean(TestTx.FileContainer)
	_ = cursors.Open(ctx, scope, parser.Identifier{Literal: "cur"}, nil)
	_ = TestTx.cachedViews.Clean(TestTx.FileContainer)
	_ = cursors.Open(ctx, scope, parser.Identifier{Literal: "cur2"}, nil)
	_, _ = cursors.Fetch(parser.Identifier{Literal: "cur2"}, parser.NEXT, 0)

	for _, v := range cursorMapIsInRangeTests {
		if 0 != v.Index {
			c, _ := cursors.Load("CUR2")
			c.index = v.Index
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

var cursorMapCountTests = []struct {
	Name    string
	CurName parser.Identifier
	Result  int
	Error   string
}{
	{
		Name:    "CursorMap Count",
		CurName: parser.Identifier{Literal: "cur"},
		Result:  3,
	},
	{
		Name:    "CursorMap Count Not Open Error",
		CurName: parser.Identifier{Literal: "cur2"},
		Error:   "cursor cur2 is closed",
	},
	{
		Name:    "CursorMap Count Undeclared Error",
		CurName: parser.Identifier{Literal: "notexist"},
		Error:   "undeclared cursor",
	},
}

func TestCursorMap_Count(t *testing.T) {
	defer func() {
		_ = TestTx.cachedViews.Clean(TestTx.FileContainer)
		initFlag(TestTx.Flags)
	}()

	TestTx.Flags.Repository = TestDir

	cursors := GenerateCursorMap([]*Cursor{
		{
			name:  "cur",
			query: selectQueryForCursorTest,
			mtx:   &sync.Mutex{},
		},
		{
			name:  "cur2",
			query: selectQueryForCursorTest,
			mtx:   &sync.Mutex{},
		},
	})
	_ = TestTx.cachedViews.Clean(TestTx.FileContainer)
	scope := NewReferenceScope(TestTx)
	ctx := context.Background()
	_ = cursors.Open(ctx, scope, parser.Identifier{Literal: "cur"}, nil)

	for _, v := range cursorMapCountTests {
		result, err := cursors.Count(v.CurName)
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
			t.Errorf("%s: result = %v, want %v", v.Name, result, v.Result)
		}
	}
}
