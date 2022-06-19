//go:build darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris || windows

package query

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/value"

	"github.com/mithrandie/readline-csvq"
)

var readlineListenerOnChangeTests = []struct {
	Line    string
	Pos     int
	Key     rune
	NewLine string
	NewPos  int
	OK      bool
}{
	{
		Line:    "abcdafghi",
		Pos:     5,
		Key:     'a',
		NewLine: "abcdafghi",
		NewPos:  5,
		OK:      false,
	},
	{
		Line:    "abcd'fghi",
		Pos:     5,
		Key:     '\'',
		NewLine: "abcd''fghi",
		NewPos:  5,
		OK:      true,
	},
	{
		Line:    "abcd'fg'hi",
		Pos:     5,
		Key:     '\'',
		NewLine: "abcd'fg'hi",
		NewPos:  5,
		OK:      false,
	},
	{
		Line:    "abcd'fg''hi",
		Pos:     8,
		Key:     '\'',
		NewLine: "abcd'fg'hi",
		NewPos:  8,
		OK:      true,
	},
	{
		Line:    "abcd(fghi",
		Pos:     5,
		Key:     '(',
		NewLine: "abcd()fghi",
		NewPos:  5,
		OK:      true,
	},
	{
		Line:    "abcd()fghi",
		Pos:     5,
		Key:     '(',
		NewLine: "abcd()fghi",
		NewPos:  5,
		OK:      false,
	},
	{
		Line:    "abcd(fghi)",
		Pos:     5,
		Key:     '(',
		NewLine: "abcd(fghi)",
		NewPos:  5,
		OK:      false,
	},
	{
		Line:    "abcd(fghi))",
		Pos:     10,
		Key:     ')',
		NewLine: "abcd(fghi)",
		NewPos:  10,
		OK:      true,
	},
}

func TestReadlineListener_OnChange(t *testing.T) {
	listener := new(ReadlineListener)

	for _, v := range readlineListenerOnChangeTests {
		newLine, newPos, ok := listener.OnChange([]rune(v.Line), v.Pos, v.Key)
		if string(newLine) != v.NewLine || newPos != v.NewPos || ok != v.OK {
			t.Errorf("result = %q %d %t, want %q %d %t for %q %d %q", string(newLine), newPos, ok, v.NewLine, v.NewPos, v.OK, v.Line, v.Pos, string(v.Key))
		}
	}
}

func TestCompleter_Update(t *testing.T) {
	defer func() {
		TestTx.PreparedStatements = NewPreparedStatementMap()
	}()

	scope := NewReferenceScope(TestTx)
	ctx := context.Background()

	scope.SetTemporaryTable(&View{
		FileInfo: &FileInfo{Path: "view1", ViewType: ViewTypeTemporaryTable},
		Header:   NewHeader("view1", []string{"col1", "col2"}),
	})
	scope.SetTemporaryTable(&View{
		FileInfo: &FileInfo{Path: "view2", ViewType: ViewTypeTemporaryTable},
		Header:   NewHeader("view1", []string{"col3", "col4"}),
	})
	_ = scope.DeclareCursor(parser.CursorDeclaration{Cursor: parser.Identifier{Literal: "cur1"}})
	_ = scope.DeclareFunction(parser.FunctionDeclaration{Name: parser.Identifier{Literal: "scalarfunc"}})
	_ = scope.DeclareAggregateFunction(parser.AggregateDeclaration{Name: parser.Identifier{Literal: "aggfunc"}})
	_ = scope.DeclareVariable(ctx, parser.VariableDeclaration{Assignments: []parser.VariableAssignment{{Variable: parser.Variable{Name: "var"}}}})

	_ = TestTx.PreparedStatements.Prepare(scope.Tx.Flags, parser.StatementPreparation{
		Name:      parser.Identifier{Literal: "stmt"},
		Statement: value.NewString("select 1"),
	})

	c := NewCompleter(scope)
	if len(c.flagList) != len(cmd.FlagList) || !strings.HasPrefix(c.flagList[0], cmd.FlagSign) {
		t.Error("flags are not set correctly")
	}
	if len(c.runinfoList) != len(RuntimeInformatinList) || !strings.HasPrefix(c.runinfoList[0], cmd.RuntimeInformationSign) {
		t.Error("runtime information are not set correctly")
	}
	if len(c.funcs) != len(Functions)+3 {
		t.Error("functions are not set correctly")
	}
	if len(c.aggFuncs) != len(AggregateFunctions)+2 {
		t.Error("aggregate functions are not set correctly")
	}
	if len(c.analyticFuncs) != len(AnalyticFunctions)+len(AggregateFunctions) {
		t.Error("analytic functions are not set correctly")
	}

	c.Update()
	if !reflect.DeepEqual(c.viewList, []string{"view1", "view2"}) {
		t.Error("views are not set correctly")
	}
	if !reflect.DeepEqual(c.cursorList, []string{"cur1"}) {
		t.Error("cursors are not set correctly")
	}
	if !reflect.DeepEqual(c.userFuncList, []string{"aggfunc", "scalarfunc"}) {
		t.Error("user defined functions are not set correctly")
	}
	if len(c.statementList) != 1 {
		t.Error("statement list is not set correctly")
	}
	if len(c.funcList) != len(Functions)+3+1 || !strings.HasSuffix(c.funcList[0], "()") {
		t.Error("function list is not set correctly")
	}
	if len(c.aggFuncList) != len(AggregateFunctions)+2+1 || !strings.HasSuffix(c.aggFuncList[0], "()") {
		t.Error("aggregate function list is not set correctly")
	}
	if len(c.analyticFuncList) != len(AnalyticFunctions)+len(AggregateFunctions)+1 || !strings.HasSuffix(c.analyticFuncList[0], "() OVER ()") {
		t.Error("analytic function list is not set correctly")
	}
	if !reflect.DeepEqual(c.varList, []string{"@var"}) {
		t.Error("variables are not set correctly")
	}
	if !strings.HasPrefix(c.envList[0], cmd.EnvironmentVariableSign) || !strings.HasPrefix(c.enclosedEnvList[0], cmd.EnvironmentVariableSign+"`") {
		t.Error("environment variables are not set correctly")
	}

	expectAllColumns := []string{"col1", "col2", "col3", "col4"}
	if !reflect.DeepEqual(c.allColumns, expectAllColumns) {
		t.Error("columns are not set correctly")
	}
}

var completer = NewCompleter(NewReferenceScope(TestTx))

type completerTest struct {
	Name     string
	Line     string
	OrigLine string
	Index    int
	Expect   readline.CandidateList
}

func testCompleter(t *testing.T, f func(line string, origLine string, index int) readline.CandidateList, tests []completerTest) {
	wd, _ := os.Getwd()

	defer func() {
		_ = os.Chdir(wd)
		_ = TestTx.cachedViews.Clean(TestTx.FileContainer)
		TestTx.PreparedStatements = NewPreparedStatementMap()
		initFlag(TestTx.Flags)
	}()

	completer.runinfoList = []string{"@#INFO1", "@#INFO2"}
	completer.statementList = []string{"stmt"}
	completer.funcs = []string{"NOW"}
	completer.aggFuncs = []string{"COUNT"}
	completer.analyticFuncs = []string{"RANK", "FIRST_VALUE"}
	completer.viewList = []string{"tempview"}
	completer.cursorList = []string{"cur1", "cur2"}
	completer.userFuncs = []string{"userfunc"}
	completer.userAggFuncs = []string{"aggfunc"}
	completer.userFuncList = []string{"userfunc", "aggfunc"}
	completer.funcList = []string{"NOW()", "userfunc()"}
	completer.aggFuncList = []string{"COUNT()", "aggfunc()"}
	completer.analyticFuncList = []string{"RANK() OVER ()", "FIRST_VALUE() OVER ()", "aggfunc"}
	completer.varList = []string{"@var1", "@var2"}
	completer.envList = []string{"@%ENV1", "@%ENV2"}
	completer.enclosedEnvList = []string{"@%`ENV1`", "@%`ENV2`"}
	TestTx.cachedViews.Set(&View{
		FileInfo: &FileInfo{
			Path: filepath.Join(CompletionTestDir, "newtable.csv"),
		},
		Header: NewHeader("newtable", []string{"ncol1", "ncol2", "ncol3"}),
	})
	TestTx.cachedViews.Set(&View{
		FileInfo: &FileInfo{
			Path: filepath.Join(CompletionTestDir, "sub", "table2.csv"),
		},
	})

	TestTx.Flags.Repository = CompletionTestDir

	_ = os.Chdir(CompletionTestDir)
	for _, v := range tests {
		completer.UpdateTokens(v.Line, string([]rune(v.OrigLine)[:v.Index]))
		result := f(v.Line, v.OrigLine, v.Index)
		if !reflect.DeepEqual(result, v.Expect) {
			t.Errorf("result = %v, want %v for %s", result, v.Expect, v.Name)
			for _, c := range result {
				t.Errorf("candidate Result: %s, %t\n", string(c.Name), c.AppendSpace)
			}
			for _, c := range v.Expect {
				t.Errorf("candidate Expect: %s, %t\n", string(c.Name), c.AppendSpace)
			}
		}
	}
}

var completerStatementsTests = []completerTest{
	{
		Name:     "Statements",
		Line:     "",
		OrigLine: "",
		Index:    0,
		Expect: readline.CandidateList{
			{Name: []rune("ADD"), AppendSpace: true},
			{Name: []rune("ALTER"), AppendSpace: true},
			{Name: []rune("CHDIR"), AppendSpace: true},
			{Name: []rune("CLOSE"), AppendSpace: true},
			{Name: []rune("COMMIT")},
			{Name: []rune("CREATE"), AppendSpace: true},
			{Name: []rune("DECLARE"), AppendSpace: true},
			{Name: []rune("DELETE"), AppendSpace: true},
			{Name: []rune("DISPOSE"), AppendSpace: true},
			{Name: []rune("ECHO"), AppendSpace: true},
			{Name: []rune("EXECUTE"), AppendSpace: true},
			{Name: []rune("EXIT")},
			{Name: []rune("FETCH"), AppendSpace: true},
			{Name: []rune("INSERT"), AppendSpace: true},
			{Name: []rune("OPEN"), AppendSpace: true},
			{Name: []rune("PREPARE"), AppendSpace: true},
			{Name: []rune("PRINT"), AppendSpace: true},
			{Name: []rune("PRINTF"), AppendSpace: true},
			{Name: []rune("PWD")},
			{Name: []rune("RELOAD"), AppendSpace: true},
			{Name: []rune("REMOVE"), AppendSpace: true},
			{Name: []rune("REPLACE"), AppendSpace: true},
			{Name: []rune("ROLLBACK")},
			{Name: []rune("SELECT"), AppendSpace: true},
			{Name: []rune("SET"), AppendSpace: true},
			{Name: []rune("SHOW"), AppendSpace: true},
			{Name: []rune("SOURCE"), AppendSpace: true},
			{Name: []rune("SYNTAX"), AppendSpace: true},
			{Name: []rune("UNSET"), AppendSpace: true},
			{Name: []rune("UPDATE"), AppendSpace: true},
			{Name: []rune("VAR"), AppendSpace: true},
			{Name: []rune("WITH"), AppendSpace: true},
		},
	},
	{
		Name:     "Statements WITH",
		Line:     "",
		OrigLine: "with ",
		Index:    5,
		Expect: readline.CandidateList{
			{Name: []rune("RECURSIVE"), AppendSpace: true},
		},
	},
	{
		Name:     "Statements SELECT",
		Line:     "fro",
		OrigLine: "select 1 fro",
		Index:    12,
		Expect: readline.CandidateList{
			{Name: []rune("AS"), AppendSpace: true},
			{Name: []rune("EXCEPT"), AppendSpace: true},
			{Name: []rune("FETCH"), AppendSpace: true},
			{Name: []rune("FOR UPDATE")},
			{Name: []rune("FROM"), AppendSpace: true},
			{Name: []rune("GROUP BY"), AppendSpace: true},
			{Name: []rune("HAVING"), AppendSpace: true},
			{Name: []rune("INTERSECT"), AppendSpace: true},
			{Name: []rune("INTO"), AppendSpace: true},
			{Name: []rune("LIMIT"), AppendSpace: true},
			{Name: []rune("OFFSET"), AppendSpace: true},
			{Name: []rune("ORDER BY"), AppendSpace: true},
			{Name: []rune("UNION"), AppendSpace: true},
			{Name: []rune("WHERE"), AppendSpace: true},
		},
	},
	{
		Name:     "Statements INSERT",
		Line:     "",
		OrigLine: "insert ",
		Index:    7,
		Expect: readline.CandidateList{
			{Name: []rune("INTO"), AppendSpace: true},
		},
	},
	{
		Name:     "Statements UPDATE",
		Line:     "",
		OrigLine: "update tbl ",
		Index:    11,
		Expect: readline.CandidateList{
			{Name: []rune("SET"), AppendSpace: true},
		},
	},
	{
		Name:     "Statements REPLACE",
		Line:     "",
		OrigLine: "replace ",
		Index:    8,
		Expect: readline.CandidateList{
			{Name: []rune("INTO"), AppendSpace: true},
		},
	},
	{
		Name:     "Statements DELETE",
		Line:     "",
		OrigLine: "delete ",
		Index:    7,
		Expect: readline.CandidateList{
			{Name: []rune("FROM"), AppendSpace: true},
		},
	},
	{
		Name:     "Statements CREATE",
		Line:     "",
		OrigLine: "create ",
		Index:    7,
		Expect: readline.CandidateList{
			{Name: []rune("TABLE"), AppendSpace: true},
		},
	},
	{
		Name:     "Statements ALTER",
		Line:     "",
		OrigLine: "alter ",
		Index:    6,
		Expect: readline.CandidateList{
			{Name: []rune("TABLE"), AppendSpace: true},
		},
	},
	{
		Name:     "Statements DECLARE",
		Line:     "",
		OrigLine: "declare cur ",
		Index:    12,
		Expect: readline.CandidateList{
			{Name: []rune("AGGREGATE"), AppendSpace: true},
			{Name: []rune("CURSOR"), AppendSpace: true},
			{Name: []rune("FUNCTION"), AppendSpace: true},
			{Name: []rune("VIEW"), AppendSpace: true},
		},
	},
	{
		Name:     "Statements PREPARE",
		Line:     "",
		OrigLine: "prepare ",
		Index:    8,
		Expect:   readline.CandidateList{},
	},
	{
		Name:     "Statements SET",
		Line:     "",
		OrigLine: "set @@flag ",
		Index:    11,
		Expect: readline.CandidateList{
			{Name: []rune("TO"), AppendSpace: true},
		},
	},
	{
		Name:     "Statements UNSET",
		Line:     "",
		OrigLine: "unset ",
		Index:    6,
		Expect: readline.CandidateList{
			{Name: []rune("@%ENV1")},
			{Name: []rune("@%ENV2")},
		},
	},
	{
		Name:     "Statements ADD",
		Line:     "",
		OrigLine: "add 1 ",
		Index:    6,
		Expect: readline.CandidateList{
			{Name: []rune("TO"), AppendSpace: true},
		},
	},
	{
		Name:     "Statements REMOVE",
		Line:     "",
		OrigLine: "remove 1 ",
		Index:    9,
		Expect: readline.CandidateList{
			{Name: []rune("FROM"), AppendSpace: true},
		},
	},
	{
		Name:     "Statements ECHO",
		Line:     "@",
		OrigLine: "echo @",
		Index:    6,
		Expect: readline.CandidateList{
			{Name: []rune("@var1")},
			{Name: []rune("@var2")},
		},
	},
	{
		Name:     "Statements PRINT",
		Line:     "@",
		OrigLine: "print @",
		Index:    7,
		Expect: readline.CandidateList{
			{Name: []rune("@var1")},
			{Name: []rune("@var2")},
		},
	},
	{
		Name:     "Statements PRINTF",
		Line:     "",
		OrigLine: "printf '%s' ",
		Index:    12,
		Expect: readline.CandidateList{
			{Name: []rune("USING"), AppendSpace: true},
		},
	},
	{
		Name:     "Statements CHDIR",
		Line:     "",
		OrigLine: "chdir ",
		Index:    6,
		Expect: readline.CandidateList{
			{Name: []rune("."), FormatAsIdentifier: true},
			{Name: []rune(".."), FormatAsIdentifier: true},
			{Name: []rune("sub/"), FormatAsIdentifier: true},
		},
	},
	{
		Name:     "Statements EXECUTE",
		Line:     "",
		OrigLine: "execute '%s' ",
		Index:    13,
		Expect: readline.CandidateList{
			{Name: []rune("USING"), AppendSpace: true},
		},
	},
	{
		Name:     "Statements SHOW",
		Line:     "c",
		OrigLine: "show c",
		Index:    6,
		Expect: append(readline.CandidateList{
			{Name: []rune("CURSORS")},
			{Name: []rune("ENV")},
			{Name: []rune("FIELDS"), AppendSpace: true},
			{Name: []rune("FLAGS")},
			{Name: []rune("FUNCTIONS")},
			{Name: []rune("RUNINFO")},
			{Name: []rune("STATEMENTS")},
			{Name: []rune("TABLES")},
			{Name: []rune("VIEWS")},
		}, completer.candidateList(completer.flagList, false)...),
	},
	{
		Name:     "Statements SOURCE",
		Line:     "",
		OrigLine: "source ",
		Index:    7,
		Expect: readline.CandidateList{
			{Name: []rune("."), FormatAsIdentifier: true},
			{Name: []rune(".."), FormatAsIdentifier: true},
			{Name: []rune("source.sql"), FormatAsIdentifier: true},
			{Name: []rune("sub/"), FormatAsIdentifier: true},
		},
	},
	{
		Name:     "Statements RELOAD",
		Line:     "",
		OrigLine: "reload ",
		Index:    7,
		Expect: readline.CandidateList{
			{Name: []rune("CONFIG")},
		},
	},
	{
		Name:     "Statements RELOAD After CONFIG",
		Line:     "",
		OrigLine: "reload config ",
		Index:    14,
		Expect:   readline.CandidateList(nil),
	},
	{
		Name:     "Statements DISPOSE",
		Line:     "@",
		OrigLine: "dispose @",
		Index:    9,
		Expect: readline.CandidateList{
			{Name: []rune("CURSOR"), AppendSpace: true},
			{Name: []rune("FUNCTION"), AppendSpace: true},
			{Name: []rune("PREPARE"), AppendSpace: true},
			{Name: []rune("VIEW"), AppendSpace: true},
			{Name: []rune("@var1")},
			{Name: []rune("@var2")},
		},
	},
	{
		Name:     "Statements OPEN",
		Line:     "",
		OrigLine: "open ",
		Index:    5,
		Expect: readline.CandidateList{
			{Name: []rune("cur1"), AppendSpace: true},
			{Name: []rune("cur2"), AppendSpace: true},
		},
	},
	{
		Name:     "Statements CLOSE",
		Line:     "",
		OrigLine: "close ",
		Index:    6,
		Expect: readline.CandidateList{
			{Name: []rune("cur1")},
			{Name: []rune("cur2")},
		},
	},
	{
		Name:     "Statements FETCH",
		Line:     "",
		OrigLine: "fetch ",
		Index:    6,
		Expect: readline.CandidateList{
			{Name: []rune("cur1"), AppendSpace: true},
			{Name: []rune("cur2"), AppendSpace: true},
			{Name: []rune("NEXT"), AppendSpace: true},
			{Name: []rune("PRIOR"), AppendSpace: true},
			{Name: []rune("FIRST"), AppendSpace: true},
			{Name: []rune("LAST"), AppendSpace: true},
			{Name: []rune("ABSOLUTE"), AppendSpace: true},
			{Name: []rune("RELATIVE"), AppendSpace: true},
		},
	},
	{
		Name:     "Statements Commit",
		Line:     "",
		OrigLine: "commit ",
		Index:    7,
		Expect:   readline.CandidateList(nil),
	},
	{
		Name:     "Statements Rollback",
		Line:     "",
		OrigLine: "rollback ",
		Index:    9,
		Expect:   readline.CandidateList(nil),
	},
	{
		Name:     "Statements Exit",
		Line:     "",
		OrigLine: "exit ",
		Index:    5,
		Expect:   readline.CandidateList(nil),
	},
	{
		Name:     "Statements PWD",
		Line:     "",
		OrigLine: "pwd ",
		Index:    4,
		Expect:   readline.CandidateList(nil),
	},
	{
		Name:     "Statements Syntax",
		Line:     "",
		OrigLine: "syntax ",
		Index:    7,
		Expect:   readline.CandidateList(nil),
	},
	{
		Name:     "Statements Variable",
		Line:     "@",
		OrigLine: "@var := @",
		Index:    9,
		Expect: readline.CandidateList{
			{Name: []rune("@var1")},
			{Name: []rune("@var2")},
		},
	},
	{
		Name:     "Statements TableObject",
		Line:     "",
		OrigLine: "select 1 from csv(",
		Index:    18,
		Expect: readline.CandidateList{
			{Name: []rune("','")},
			{Name: []rune("'\\t'")},
		},
	},
	{
		Name:     "Statements Function",
		Line:     "d",
		OrigLine: "select count(d",
		Index:    14,
		Expect: readline.CandidateList{
			{Name: []rune("DISTINCT"), AppendSpace: true},
		},
	},
}

func TestCompleter_Statements(t *testing.T) {
	testCompleter(t, completer.Statements, completerStatementsTests)
}

var completerTableObjectArgsTests = []completerTest{
	{
		Name:     "TableObjectArgs CSV Delimiter",
		Line:     "",
		OrigLine: "csv(",
		Index:    4,
		Expect: readline.CandidateList{
			{Name: []rune("','")},
			{Name: []rune("'\\t'")},
		},
	},
	{
		Name:     "TableObjectArgs FIXED Delimiter",
		Line:     "",
		OrigLine: "fixed(",
		Index:    6,
		Expect: readline.CandidateList{
			{Name: []rune("'SPACES'")},
			{Name: []rune("'S[]'")},
			{Name: []rune("'[]'")},
		},
	},
	{
		Name:     "TableObjectArgs CSV Files",
		Line:     "",
		OrigLine: "csv(',',",
		Index:    8,
		Expect: readline.CandidateList{
			{Name: []rune(filepath.Join(CompletionTestDir, "sub", "table2.csv")), FormatAsIdentifier: true},
			{Name: []rune("newtable.csv"), FormatAsIdentifier: true},
			{Name: []rune("tempview"), FormatAsIdentifier: true},
			{Name: []rune("."), FormatAsIdentifier: true},
			{Name: []rune(".."), FormatAsIdentifier: true},
			{Name: []rune("sub/"), FormatAsIdentifier: true},
			{Name: []rune("table1.csv"), FormatAsIdentifier: true},
		},
	},
	{
		Name:     "TableObjectArgs CSV Encoding",
		Line:     "",
		OrigLine: "csv(',', filepath, ",
		Index:    19,
		Expect: readline.CandidateList{
			{Name: []rune("AUTO")},
			{Name: []rune("SJIS")},
			{Name: []rune("UTF16")},
			{Name: []rune("UTF16BE")},
			{Name: []rune("UTF16BEM")},
			{Name: []rune("UTF16LE")},
			{Name: []rune("UTF16LEM")},
			{Name: []rune("UTF8")},
			{Name: []rune("UTF8M")},
		},
	},
	{
		Name:     "TableObjectArgs CSV NoHeader",
		Line:     "",
		OrigLine: "csv(',', filepath, utf8,",
		Index:    24,
		Expect: readline.CandidateList{
			{Name: []rune("TRUE")},
			{Name: []rune("FALSE")},
		},
	},
	{
		Name:     "TableObjectArgs CSV WithoutNull",
		Line:     "",
		OrigLine: "csv(',', filepath, utf8, false,",
		Index:    31,
		Expect: readline.CandidateList{
			{Name: []rune("TRUE")},
			{Name: []rune("FALSE")},
		},
	},
	{
		Name:     "TableObjectArgs CSV Value",
		Line:     "@",
		OrigLine: "csv(@",
		Index:    5,
		Expect: readline.CandidateList{
			{Name: []rune("','")},
			{Name: []rune("'\\t'")},
			{Name: []rune("@var1")},
			{Name: []rune("@var2")},
		},
	},
	{
		Name:     "TableObjectArgs LTSV Files",
		Line:     "",
		OrigLine: "ltsv(",
		Index:    5,
		Expect: readline.CandidateList{
			{Name: []rune(filepath.Join(CompletionTestDir, "sub", "table2.csv")), FormatAsIdentifier: true},
			{Name: []rune("newtable.csv"), FormatAsIdentifier: true},
			{Name: []rune("tempview"), FormatAsIdentifier: true},
			{Name: []rune("."), FormatAsIdentifier: true},
			{Name: []rune(".."), FormatAsIdentifier: true},
			{Name: []rune("sub/"), FormatAsIdentifier: true},
			{Name: []rune("table1.csv"), FormatAsIdentifier: true},
		},
	},
	{
		Name:     "TableObjectArgs LTSV Encoding",
		Line:     "",
		OrigLine: "ltsv(filepath, ",
		Index:    15,
		Expect: readline.CandidateList{
			{Name: []rune("AUTO")},
			{Name: []rune("SJIS")},
			{Name: []rune("UTF16")},
			{Name: []rune("UTF16BE")},
			{Name: []rune("UTF16BEM")},
			{Name: []rune("UTF16LE")},
			{Name: []rune("UTF16LEM")},
			{Name: []rune("UTF8")},
			{Name: []rune("UTF8M")},
		},
	},
	{
		Name:     "TableObjectArgs LTSV WithoutNull",
		Line:     "",
		OrigLine: "ltsv(filepath, utf8, ",
		Index:    21,
		Expect: readline.CandidateList{
			{Name: []rune("TRUE")},
			{Name: []rune("FALSE")},
		},
	},
}

func TestCompleter_TableObjectArgs(t *testing.T) {
	testCompleter(t, completer.TableObjectArgs, completerTableObjectArgsTests)
}

var completerFunctionArgs = []completerTest{
	{
		Name:     "FunctionArgs",
		Line:     "@v",
		OrigLine: "trim(@v",
		Index:    7,
		Expect: readline.CandidateList{
			{Name: []rune("@var1")},
			{Name: []rune("@var2")},
		},
	},
	{
		Name:     "FunctionArgs DISTINCT",
		Line:     "d",
		OrigLine: "count(d",
		Index:    7,
		Expect: readline.CandidateList{
			{Name: []rune("DISTINCT"), AppendSpace: true},
		},
	},
	{
		Name:     "FunctionArgs After OVER",
		Line:     "",
		OrigLine: "count(1) over (",
		Index:    15,
		Expect: readline.CandidateList{
			{Name: []rune("ORDER BY"), AppendSpace: true},
			{Name: []rune("PARTITION BY"), AppendSpace: true},
		},
	},
	{
		Name:     "FunctionArgs After PARTITION",
		Line:     "",
		OrigLine: "count(1) over (partition ",
		Index:    25,
		Expect: readline.CandidateList{
			{Name: []rune("BY"), AppendSpace: true},
		},
	},
	{
		Name:     "FunctionArgs After PARTITION BY",
		Line:     "@v",
		OrigLine: "count(1) over (partition by @v",
		Index:    30,
		Expect: readline.CandidateList{
			{Name: []rune("@var1")},
			{Name: []rune("@var2")},
		},
	},
	{
		Name:     "FunctionArgs After Partition Clause",
		Line:     "",
		OrigLine: "count(1) over (partition by f1 ",
		Index:    31,
		Expect: readline.CandidateList{
			{Name: []rune("ORDER BY"), AppendSpace: true},
		},
	},
	{
		Name:     "FunctionArgs After ORDER",
		Line:     "",
		OrigLine: "count(1) over (order ",
		Index:    21,
		Expect: readline.CandidateList{
			{Name: []rune("BY"), AppendSpace: true},
		},
	},
	{
		Name:     "FunctionArgs After ORDER BY",
		Line:     "@v",
		OrigLine: "count(1) over (order by @v",
		Index:    26,
		Expect: readline.CandidateList{
			{Name: []rune("@var1")},
			{Name: []rune("@var2")},
		},
	},
	{
		Name:     "FunctionArgs After Order Value in Order By Clause",
		Line:     "",
		OrigLine: "count(1) over (order by f1 ",
		Index:    27,
		Expect: readline.CandidateList{
			{Name: []rune("ASC")},
			{Name: []rune("DESC")},
			{Name: []rune("NULLS FIRST")},
			{Name: []rune("NULLS LAST")},
			{Name: []rune("ROWS"), AppendSpace: true},
		},
	},
	{
		Name:     "FunctionArgs After ASC in Order By Clause",
		Line:     "",
		OrigLine: "count(1) over (order by f1 asc ",
		Index:    31,
		Expect: readline.CandidateList{
			{Name: []rune("NULLS FIRST")},
			{Name: []rune("NULLS LAST")},
			{Name: []rune("ROWS"), AppendSpace: true},
		},
	},
	{
		Name:     "FunctionArgs After NULLS in Order By Clause",
		Line:     "",
		OrigLine: "count(1) over (order by f1 asc nulls ",
		Index:    37,
		Expect: readline.CandidateList{
			{Name: []rune("FIRST")},
			{Name: []rune("LAST")},
		},
	},
	{
		Name:     "FunctionArgs After NULLS FIRST in Order By Clause",
		Line:     "",
		OrigLine: "count(1) over (order by f1 asc nulls first ",
		Index:    43,
		Expect: readline.CandidateList{
			{Name: []rune("ROWS"), AppendSpace: true},
		},
	},
	{
		Name:     "FunctionArgs After ROWS in Windowing Clause",
		Line:     "",
		OrigLine: "count(1) over (order by f1 rows ",
		Index:    32,
		Expect: readline.CandidateList{
			{Name: []rune("BETWEEN"), AppendSpace: true},
			{Name: []rune("CURRENT ROW")},
			{Name: []rune("UNBOUNDED PRECEDING")},
		},
	},
	{
		Name:     "FunctionArgs After CURRENT in Windowing Clause",
		Line:     "",
		OrigLine: "count(1) over (order by f1 rows current ",
		Index:    40,
		Expect: readline.CandidateList{
			{Name: []rune("ROW")},
		},
	},
	{
		Name:     "FunctionArgs After UNBOUNDED in Windowing Clause",
		Line:     "",
		OrigLine: "count(1) over (order by f1 rows unbounded ",
		Index:    42,
		Expect: readline.CandidateList{
			{Name: []rune("PRECEDING")},
		},
	},
	{
		Name:     "FunctionArgs After BETWEEN in Windowing Clause",
		Line:     "",
		OrigLine: "count(1) over (order by f1 rows between ",
		Index:    40,
		Expect: readline.CandidateList{
			{Name: []rune("CURRENT ROW"), AppendSpace: true},
			{Name: []rune("UNBOUNDED PRECEDING"), AppendSpace: true},
		},
	},
	{
		Name:     "FunctionArgs After CURRENT in Windowing Clause",
		Line:     "",
		OrigLine: "count(1) over (order by f1 rows between current ",
		Index:    48,
		Expect: readline.CandidateList{
			{Name: []rune("ROW"), AppendSpace: true},
		},
	},
	{
		Name:     "FunctionArgs After UNBOUNDED in Windowing Clause",
		Line:     "",
		OrigLine: "count(1) over (order by f1 rows between unbounded ",
		Index:    50,
		Expect: readline.CandidateList{
			{Name: []rune("PRECEDING"), AppendSpace: true},
		},
	},
	{
		Name:     "FunctionArgs After Offset in Windowing Clause",
		Line:     "",
		OrigLine: "count(1) over (order by f1 rows between 1 ",
		Index:    42,
		Expect: readline.CandidateList{
			{Name: []rune("FOLLOWING"), AppendSpace: true},
			{Name: []rune("PRECEDING"), AppendSpace: true},
		},
	},
	{
		Name:     "FunctionArgs After Low Frame in Windowing Clause",
		Line:     "",
		OrigLine: "count(1) over (order by f1 rows between 1 preceding ",
		Index:    52,
		Expect: readline.CandidateList{
			{Name: []rune("AND"), AppendSpace: true},
		},
	},
	{
		Name:     "FunctionArgs After AND in Windowing Clause",
		Line:     "",
		OrigLine: "count(1) over (order by f1 rows between 1 preceding and ",
		Index:    56,
		Expect: readline.CandidateList{
			{Name: []rune("CURRENT ROW")},
			{Name: []rune("UNBOUNDED FOLLOWING")},
		},
	},
	{
		Name:     "FunctionArgs After UNBOUNDED in Windowing Clause High Frame",
		Line:     "",
		OrigLine: "count(1) over (order by f1 rows between 1 preceding and unbounded ",
		Index:    66,
		Expect: readline.CandidateList{
			{Name: []rune("FOLLOWING")},
		},
	},
	{
		Name:     "FunctionArgs After UNBOUNDED in Windowing Clause High Frame",
		Line:     "",
		OrigLine: "count(1) over (order by f1 rows between 1 preceding and current ",
		Index:    64,
		Expect: readline.CandidateList{
			{Name: []rune("ROW")},
		},
	},
	{
		Name:     "FunctionArgs After Offset in Windowing Clause High Frame",
		Line:     "",
		OrigLine: "count(1) over (order by f1 rows between 1 preceding and 1 ",
		Index:    58,
		Expect: readline.CandidateList{
			{Name: []rune("FOLLOWING")},
			{Name: []rune("PRECEDING")},
		},
	},
	{
		Name:     "FunctionArgs Substring After Extraction String",
		Line:     "",
		OrigLine: "substring('abc' ",
		Index:    16,
		Expect: readline.CandidateList{
			{Name: []rune("FROM"), AppendSpace: true},
		},
	},
	{
		Name:     "FunctionArgs Substring After Position",
		Line:     "",
		OrigLine: "substring('abc' from 2 ",
		Index:    23,
		Expect: readline.CandidateList{
			{Name: []rune("FOR"), AppendSpace: true},
		},
	},
	{
		Name:     "FunctionArgs Substring After FOR",
		Line:     "",
		OrigLine: "substring('abc' from 2 for ",
		Index:    27,
		Expect:   readline.CandidateList{},
	},
}

func TestCompleter_FunctionArgs(t *testing.T) {
	testCompleter(t, completer.FunctionArgs, completerFunctionArgs)
}

var completerWithArgsTests = []completerTest{
	{
		Name:     "WithArgs",
		Line:     "",
		OrigLine: "with ",
		Index:    5,
		Expect: readline.CandidateList{
			{Name: []rune("RECURSIVE"), AppendSpace: true},
		},
	},
	{
		Name:     "WithArgs After Table Name",
		Line:     "",
		OrigLine: "with tbl ",
		Index:    9,
		Expect: readline.CandidateList{
			{Name: []rune("AS"), AppendSpace: true},
		},
	},
	{
		Name:     "WithArgs After Table Name with RECURSIVE",
		Line:     "",
		OrigLine: "with recursive tbl ",
		Index:    19,
		Expect: readline.CandidateList{
			{Name: []rune("AS"), AppendSpace: true},
		},
	},
	{
		Name:     "WithArgs in Column Enumulation",
		Line:     "",
		OrigLine: "with tbl (col1, ",
		Index:    16,
		Expect:   readline.CandidateList{},
	},
	{
		Name:     "WithArgs After Column Enumulation",
		Line:     "",
		OrigLine: "with tbl (col1, col2) ",
		Index:    22,
		Expect: readline.CandidateList{
			{Name: []rune("AS"), AppendSpace: true},
		},
	},
	{
		Name:     "WithArgs After AS",
		Line:     "",
		OrigLine: "with tbl (col1, col2) as (",
		Index:    26,
		Expect: readline.CandidateList{
			{Name: []rune("SELECT"), AppendSpace: true},
		},
	},
	{
		Name:     "WithArgs Second Table",
		Line:     "",
		OrigLine: "with tbl1 as (select 1 as 'a'), ",
		Index:    32,
		Expect: readline.CandidateList{
			{Name: []rune("RECURSIVE"), AppendSpace: true},
		},
	},
	{
		Name:     "WithArgs After Table Definition",
		Line:     "",
		OrigLine: "with tbl1 as (select 1 as 'a') ",
		Index:    31,
		Expect: readline.CandidateList{
			{Name: []rune("DELETE"), AppendSpace: true},
			{Name: []rune("INSERT"), AppendSpace: true},
			{Name: []rune("REPLACE"), AppendSpace: true},
			{Name: []rune("SELECT"), AppendSpace: true},
			{Name: []rune("UPDATE"), AppendSpace: true},
		},
	},
	{
		Name:     "WithArgs Select Query",
		Line:     "fro",
		OrigLine: "with tbl1 as (select 1 as 'a') select 1 fro",
		Index:    43,
		Expect: readline.CandidateList{
			{Name: []rune("AS"), AppendSpace: true},
			{Name: []rune("EXCEPT"), AppendSpace: true},
			{Name: []rune("FETCH"), AppendSpace: true},
			{Name: []rune("FOR UPDATE")},
			{Name: []rune("FROM"), AppendSpace: true},
			{Name: []rune("GROUP BY"), AppendSpace: true},
			{Name: []rune("HAVING"), AppendSpace: true},
			{Name: []rune("INTERSECT"), AppendSpace: true},
			{Name: []rune("INTO"), AppendSpace: true},
			{Name: []rune("LIMIT"), AppendSpace: true},
			{Name: []rune("OFFSET"), AppendSpace: true},
			{Name: []rune("ORDER BY"), AppendSpace: true},
			{Name: []rune("UNION"), AppendSpace: true},
			{Name: []rune("WHERE"), AppendSpace: true},
		},
	},
	{
		Name:     "WithArgs Insert Query",
		Line:     "",
		OrigLine: "with tbl1 as (select 1 as 'a') insert ",
		Index:    37,
		Expect: readline.CandidateList{
			{Name: []rune("INTO"), AppendSpace: true},
		},
	},
	{
		Name:     "WithArgs Update Query",
		Line:     "",
		OrigLine: "with tbl1 as (select 1 as 'a') update tbl ",
		Index:    41,
		Expect: readline.CandidateList{
			{Name: []rune("SET"), AppendSpace: true},
		},
	},
	{
		Name:     "WithArgs Replace Query",
		Line:     "",
		OrigLine: "with tbl1 as (select 1 as 'a') replace ",
		Index:    38,
		Expect: readline.CandidateList{
			{Name: []rune("INTO"), AppendSpace: true},
		},
	},
	{
		Name:     "WithArgs Delete Query",
		Line:     "",
		OrigLine: "with tbl1 as (select 1 as 'a') delete ",
		Index:    37,
		Expect: readline.CandidateList{
			{Name: []rune("FROM"), AppendSpace: true},
		},
	},
}

func TestCompleter_WithArgs(t *testing.T) {
	testCompleter(t, completer.WithArgs, completerWithArgsTests)
}

var completerSelectArgsTests = []completerTest{
	{
		Name:     "SelectArgs",
		Line:     "d",
		OrigLine: "select d",
		Index:    8,
		Expect: readline.CandidateList{
			{Name: []rune("DISTINCT"), AppendSpace: true},
		},
	},
	{
		Name:     "SelectArgs After Column Name",
		Line:     "",
		OrigLine: "select col1 ",
		Index:    12,
		Expect: readline.CandidateList{
			{Name: []rune("AS"), AppendSpace: true},
			{Name: []rune("EXCEPT"), AppendSpace: true},
			{Name: []rune("FETCH"), AppendSpace: true},
			{Name: []rune("FOR UPDATE")},
			{Name: []rune("FROM"), AppendSpace: true},
			{Name: []rune("GROUP BY"), AppendSpace: true},
			{Name: []rune("HAVING"), AppendSpace: true},
			{Name: []rune("INTERSECT"), AppendSpace: true},
			{Name: []rune("INTO"), AppendSpace: true},
			{Name: []rune("LIMIT"), AppendSpace: true},
			{Name: []rune("OFFSET"), AppendSpace: true},
			{Name: []rune("ORDER BY"), AppendSpace: true},
			{Name: []rune("UNION"), AppendSpace: true},
			{Name: []rune("WHERE"), AppendSpace: true},
		},
	},
	{
		Name:     "SelectArgs After Comma",
		Line:     "",
		OrigLine: "select col1, ",
		Index:    13,
		Expect:   readline.CandidateList{},
	},
	{
		Name:     "SelectArgs After Column Name with DISTINCT",
		Line:     "",
		OrigLine: "select distinct col1 ",
		Index:    21,
		Expect: readline.CandidateList{
			{Name: []rune("AS"), AppendSpace: true},
			{Name: []rune("EXCEPT"), AppendSpace: true},
			{Name: []rune("FETCH"), AppendSpace: true},
			{Name: []rune("FOR UPDATE")},
			{Name: []rune("FROM"), AppendSpace: true},
			{Name: []rune("GROUP BY"), AppendSpace: true},
			{Name: []rune("HAVING"), AppendSpace: true},
			{Name: []rune("INTERSECT"), AppendSpace: true},
			{Name: []rune("INTO"), AppendSpace: true},
			{Name: []rune("LIMIT"), AppendSpace: true},
			{Name: []rune("OFFSET"), AppendSpace: true},
			{Name: []rune("ORDER BY"), AppendSpace: true},
			{Name: []rune("UNION"), AppendSpace: true},
			{Name: []rune("WHERE"), AppendSpace: true},
		},
	},
	{
		Name:     "SelectArgs Analytic Function in Select Clause",
		Line:     "r",
		OrigLine: "select col1, r",
		Index:    14,
		Expect: readline.CandidateList{
			{Name: []rune("RANK() OVER ()")},
		},
	},
	{
		Name:     "SelectArgs Before OVER of Analytic Function in Select Clause",
		Line:     "",
		OrigLine: "select col1, first_value(col1) ",
		Index:    31,
		Expect: readline.CandidateList{
			{Name: []rune("AS"), AppendSpace: true},
			{Name: []rune("EXCEPT"), AppendSpace: true},
			{Name: []rune("FETCH"), AppendSpace: true},
			{Name: []rune("FOR UPDATE")},
			{Name: []rune("FROM"), AppendSpace: true},
			{Name: []rune("GROUP BY"), AppendSpace: true},
			{Name: []rune("HAVING"), AppendSpace: true},
			{Name: []rune("IGNORE NULLS"), AppendSpace: true},
			{Name: []rune("INTERSECT"), AppendSpace: true},
			{Name: []rune("INTO"), AppendSpace: true},
			{Name: []rune("LIMIT"), AppendSpace: true},
			{Name: []rune("OFFSET"), AppendSpace: true},
			{Name: []rune("ORDER BY"), AppendSpace: true},
			{Name: []rune("OVER"), AppendSpace: true},
			{Name: []rune("UNION"), AppendSpace: true},
			{Name: []rune("WHERE"), AppendSpace: true},
		},
	},
	{
		Name:     "SelectArgs After INTO",
		Line:     "",
		OrigLine: "select 1 into ",
		Index:    14,
		Expect: readline.CandidateList{
			{Name: []rune("@var1")},
			{Name: []rune("@var2")},
		},
	},
	{
		Name:     "SelectArgs After Variable in INTO Clause",
		Line:     "",
		OrigLine: "select 1 into @var1 ",
		Index:    20,
		Expect: readline.CandidateList{
			{Name: []rune("FETCH"), AppendSpace: true},
			{Name: []rune("FOR UPDATE")},
			{Name: []rune("FROM"), AppendSpace: true},
			{Name: []rune("GROUP BY"), AppendSpace: true},
			{Name: []rune("HAVING"), AppendSpace: true},
			{Name: []rune("LIMIT"), AppendSpace: true},
			{Name: []rune("OFFSET"), AppendSpace: true},
			{Name: []rune("ORDER BY"), AppendSpace: true},
			{Name: []rune("WHERE"), AppendSpace: true},
		},
	},
	{
		Name:     "SelectArgs After FROM",
		Line:     "",
		OrigLine: "select 1 from ",
		Index:    14,
		Expect: readline.CandidateList{
			{Name: []rune("CSV()")},
			{Name: []rune("FIXED()")},
			{Name: []rune("JSON()")},
			{Name: []rune("JSONL()")},
			{Name: []rune("JSON_TABLE()")},
			{Name: []rune("LTSV()")},
			{Name: []rune(filepath.Join(CompletionTestDir, "sub", "table2.csv")), FormatAsIdentifier: true},
			{Name: []rune("newtable.csv"), FormatAsIdentifier: true},
			{Name: []rune("tempview"), FormatAsIdentifier: true},
			{Name: []rune("."), FormatAsIdentifier: true},
			{Name: []rune(".."), FormatAsIdentifier: true},
			{Name: []rune("sub/"), FormatAsIdentifier: true},
			{Name: []rune("table1.csv"), FormatAsIdentifier: true},
		},
	},
	{
		Name:     "SelectArgs After Table Name in From Clause",
		Line:     "",
		OrigLine: "select 1 from tbl ",
		Index:    18,
		Expect: readline.CandidateList{
			{Name: []rune("AS"), AppendSpace: true},
			{Name: []rune("CROSS"), AppendSpace: true},
			{Name: []rune("EXCEPT"), AppendSpace: true},
			{Name: []rune("FETCH"), AppendSpace: true},
			{Name: []rune("FOR UPDATE")},
			{Name: []rune("FULL"), AppendSpace: true},
			{Name: []rune("GROUP BY"), AppendSpace: true},
			{Name: []rune("HAVING"), AppendSpace: true},
			{Name: []rune("INNER"), AppendSpace: true},
			{Name: []rune("INTERSECT"), AppendSpace: true},
			{Name: []rune("JOIN"), AppendSpace: true},
			{Name: []rune("LEFT"), AppendSpace: true},
			{Name: []rune("LIMIT"), AppendSpace: true},
			{Name: []rune("NATURAL"), AppendSpace: true},
			{Name: []rune("OFFSET"), AppendSpace: true},
			{Name: []rune("ORDER BY"), AppendSpace: true},
			{Name: []rune("RIGHT"), AppendSpace: true},
			{Name: []rune("UNION"), AppendSpace: true},
			{Name: []rune("WHERE"), AppendSpace: true},
		},
	},
	{
		Name:     "SelectArgs After Subquery in From Clause",
		Line:     "a",
		OrigLine: "select 1 from (select 1) a",
		Index:    26,
		Expect: readline.CandidateList{
			{Name: []rune("AS"), AppendSpace: true},
			{Name: []rune("CROSS"), AppendSpace: true},
			{Name: []rune("EXCEPT"), AppendSpace: true},
			{Name: []rune("FETCH"), AppendSpace: true},
			{Name: []rune("FOR UPDATE")},
			{Name: []rune("FULL"), AppendSpace: true},
			{Name: []rune("GROUP BY"), AppendSpace: true},
			{Name: []rune("HAVING"), AppendSpace: true},
			{Name: []rune("INNER"), AppendSpace: true},
			{Name: []rune("INTERSECT"), AppendSpace: true},
			{Name: []rune("JOIN"), AppendSpace: true},
			{Name: []rune("LEFT"), AppendSpace: true},
			{Name: []rune("LIMIT"), AppendSpace: true},
			{Name: []rune("NATURAL"), AppendSpace: true},
			{Name: []rune("OFFSET"), AppendSpace: true},
			{Name: []rune("ORDER BY"), AppendSpace: true},
			{Name: []rune("RIGHT"), AppendSpace: true},
			{Name: []rune("UNION"), AppendSpace: true},
			{Name: []rune("WHERE"), AppendSpace: true},
		},
	},
	{
		Name:     "SelectArgs Subquery in From Clause",
		Line:     "",
		OrigLine: "select 1 from (",
		Index:    15,
		Expect: readline.CandidateList{
			{Name: []rune("CSV()")},
			{Name: []rune("FIXED()")},
			{Name: []rune("JSON()")},
			{Name: []rune("JSONL()")},
			{Name: []rune("JSON_TABLE()")},
			{Name: []rune("LTSV()")},
			{Name: []rune(filepath.Join(CompletionTestDir, "sub", "table2.csv")), FormatAsIdentifier: true},
			{Name: []rune("newtable.csv"), FormatAsIdentifier: true},
			{Name: []rune("tempview"), FormatAsIdentifier: true},
			{Name: []rune("."), FormatAsIdentifier: true},
			{Name: []rune(".."), FormatAsIdentifier: true},
			{Name: []rune("sub/"), FormatAsIdentifier: true},
			{Name: []rune("table1.csv"), FormatAsIdentifier: true},
			{Name: []rune("SELECT"), AppendSpace: true},
		},
	},
	{
		Name:     "SelectArgs After Table Name Alias in From Clause",
		Line:     "",
		OrigLine: "select 1 from tb1 t ",
		Index:    20,
		Expect: readline.CandidateList{
			{Name: []rune("CROSS"), AppendSpace: true},
			{Name: []rune("EXCEPT"), AppendSpace: true},
			{Name: []rune("FETCH"), AppendSpace: true},
			{Name: []rune("FOR UPDATE")},
			{Name: []rune("FULL"), AppendSpace: true},
			{Name: []rune("GROUP BY"), AppendSpace: true},
			{Name: []rune("HAVING"), AppendSpace: true},
			{Name: []rune("INNER"), AppendSpace: true},
			{Name: []rune("INTERSECT"), AppendSpace: true},
			{Name: []rune("JOIN"), AppendSpace: true},
			{Name: []rune("LEFT"), AppendSpace: true},
			{Name: []rune("LIMIT"), AppendSpace: true},
			{Name: []rune("NATURAL"), AppendSpace: true},
			{Name: []rune("OFFSET"), AppendSpace: true},
			{Name: []rune("ORDER BY"), AppendSpace: true},
			{Name: []rune("RIGHT"), AppendSpace: true},
			{Name: []rune("UNION"), AppendSpace: true},
			{Name: []rune("WHERE"), AppendSpace: true},
		},
	},
	{
		Name:     "SelectArgs After Parentheses in From Clause",
		Line:     "",
		OrigLine: "select 1 from (tbl t join (tbl2 t2)) ",
		Index:    37,
		Expect: readline.CandidateList{
			{Name: []rune("CROSS"), AppendSpace: true},
			{Name: []rune("EXCEPT"), AppendSpace: true},
			{Name: []rune("FETCH"), AppendSpace: true},
			{Name: []rune("FOR UPDATE")},
			{Name: []rune("FULL"), AppendSpace: true},
			{Name: []rune("GROUP BY"), AppendSpace: true},
			{Name: []rune("HAVING"), AppendSpace: true},
			{Name: []rune("INNER"), AppendSpace: true},
			{Name: []rune("INTERSECT"), AppendSpace: true},
			{Name: []rune("JOIN"), AppendSpace: true},
			{Name: []rune("LEFT"), AppendSpace: true},
			{Name: []rune("LIMIT"), AppendSpace: true},
			{Name: []rune("NATURAL"), AppendSpace: true},
			{Name: []rune("OFFSET"), AppendSpace: true},
			{Name: []rune("ORDER BY"), AppendSpace: true},
			{Name: []rune("RIGHT"), AppendSpace: true},
			{Name: []rune("UNION"), AppendSpace: true},
			{Name: []rune("WHERE"), AppendSpace: true},
		},
	},
	{
		Name:     "SelectArgs After AS in From Clause",
		Line:     "",
		OrigLine: "select 1 from tb1 as ",
		Index:    21,
		Expect:   readline.CandidateList{},
	},
	{
		Name:     "SelectArgs After Table Name Alias with AS in From Clause",
		Line:     "",
		OrigLine: "select 1 from tb1 as t ",
		Index:    23,
		Expect: readline.CandidateList{
			{Name: []rune("CROSS"), AppendSpace: true},
			{Name: []rune("EXCEPT"), AppendSpace: true},
			{Name: []rune("FETCH"), AppendSpace: true},
			{Name: []rune("FOR UPDATE")},
			{Name: []rune("FULL"), AppendSpace: true},
			{Name: []rune("GROUP BY"), AppendSpace: true},
			{Name: []rune("HAVING"), AppendSpace: true},
			{Name: []rune("INNER"), AppendSpace: true},
			{Name: []rune("INTERSECT"), AppendSpace: true},
			{Name: []rune("JOIN"), AppendSpace: true},
			{Name: []rune("LEFT"), AppendSpace: true},
			{Name: []rune("LIMIT"), AppendSpace: true},
			{Name: []rune("NATURAL"), AppendSpace: true},
			{Name: []rune("OFFSET"), AppendSpace: true},
			{Name: []rune("ORDER BY"), AppendSpace: true},
			{Name: []rune("RIGHT"), AppendSpace: true},
			{Name: []rune("UNION"), AppendSpace: true},
			{Name: []rune("WHERE"), AppendSpace: true},
		},
	},
	{
		Name:     "SelectArgs After INNER in From Clause",
		Line:     "",
		OrigLine: "select 1 from tb1 as t inner ",
		Index:    29,
		Expect: readline.CandidateList{
			{Name: []rune("JOIN"), AppendSpace: true},
		},
	},
	{
		Name:     "SelectArgs After FULL in From Clause",
		Line:     "",
		OrigLine: "select 1 from tb1 as t full ",
		Index:    28,
		Expect: readline.CandidateList{
			{Name: []rune("JOIN"), AppendSpace: true},
			{Name: []rune("OUTER"), AppendSpace: true},
		},
	},
	{
		Name:     "SelectArgs After NATURAL in From Clause",
		Line:     "",
		OrigLine: "select 1 from tb1 as t natural ",
		Index:    31,
		Expect: readline.CandidateList{
			{Name: []rune("INNER"), AppendSpace: true},
			{Name: []rune("LEFT"), AppendSpace: true},
			{Name: []rune("RIGHT"), AppendSpace: true},
		},
	},
	{
		Name:     "SelectArgs After CROSS JOIN in From Clause",
		Line:     "",
		OrigLine: "select 1 from t1 cross join t2 ",
		Index:    31,
		Expect: readline.CandidateList{
			{Name: []rune("AS"), AppendSpace: true},
			{Name: []rune("CROSS"), AppendSpace: true},
			{Name: []rune("EXCEPT"), AppendSpace: true},
			{Name: []rune("FETCH"), AppendSpace: true},
			{Name: []rune("FOR UPDATE")},
			{Name: []rune("FULL"), AppendSpace: true},
			{Name: []rune("GROUP BY"), AppendSpace: true},
			{Name: []rune("HAVING"), AppendSpace: true},
			{Name: []rune("INNER"), AppendSpace: true},
			{Name: []rune("INTERSECT"), AppendSpace: true},
			{Name: []rune("JOIN"), AppendSpace: true},
			{Name: []rune("LEFT"), AppendSpace: true},
			{Name: []rune("LIMIT"), AppendSpace: true},
			{Name: []rune("NATURAL"), AppendSpace: true},
			{Name: []rune("OFFSET"), AppendSpace: true},
			{Name: []rune("ORDER BY"), AppendSpace: true},
			{Name: []rune("RIGHT"), AppendSpace: true},
			{Name: []rune("UNION"), AppendSpace: true},
			{Name: []rune("WHERE"), AppendSpace: true},
		},
	},
	{
		Name:     "SelectArgs After INNER JOIN in From Clause",
		Line:     "",
		OrigLine: "select 1 from t1 inner join t2 ",
		Index:    31,
		Expect: readline.CandidateList{
			{Name: []rune("AS"), AppendSpace: true},
			{Name: []rune("ON"), AppendSpace: true},
			{Name: []rune("USING ()"), AppendSpace: true},
		},
	},
	{
		Name:     "SelectArgs After FULL JOIN in From Clause",
		Line:     "",
		OrigLine: "select 1 from t1 full outer join t2 ",
		Index:    36,
		Expect: readline.CandidateList{
			{Name: []rune("AS"), AppendSpace: true},
			{Name: []rune("ON"), AppendSpace: true},
		},
	},
	{
		Name:     "SelectArgs After JOIN Keyword",
		Line:     "",
		OrigLine: "select 1 from t1, ",
		Index:    18,
		Expect: readline.CandidateList{
			{Name: []rune("CSV()")},
			{Name: []rune("FIXED()")},
			{Name: []rune("JSON()")},
			{Name: []rune("JSONL()")},
			{Name: []rune("JSON_TABLE()")},
			{Name: []rune("LTSV()")},
			{Name: []rune(filepath.Join(CompletionTestDir, "sub", "table2.csv")), FormatAsIdentifier: true},
			{Name: []rune("newtable.csv"), FormatAsIdentifier: true},
			{Name: []rune("tempview"), FormatAsIdentifier: true},
			{Name: []rune("."), FormatAsIdentifier: true},
			{Name: []rune(".."), FormatAsIdentifier: true},
			{Name: []rune("sub/"), FormatAsIdentifier: true},
			{Name: []rune("table1.csv"), FormatAsIdentifier: true},
			{Name: []rune("LATERAL"), AppendSpace: true},
		},
	},
	{
		Name:     "SelectArgs After JOIN Keyword",
		Line:     "",
		OrigLine: "select 1 from t1 inner join ",
		Index:    28,
		Expect: readline.CandidateList{
			{Name: []rune("CSV()")},
			{Name: []rune("FIXED()")},
			{Name: []rune("JSON()")},
			{Name: []rune("JSONL()")},
			{Name: []rune("JSON_TABLE()")},
			{Name: []rune("LTSV()")},
			{Name: []rune(filepath.Join(CompletionTestDir, "sub", "table2.csv")), FormatAsIdentifier: true},
			{Name: []rune("newtable.csv"), FormatAsIdentifier: true},
			{Name: []rune("tempview"), FormatAsIdentifier: true},
			{Name: []rune("."), FormatAsIdentifier: true},
			{Name: []rune(".."), FormatAsIdentifier: true},
			{Name: []rune("sub/"), FormatAsIdentifier: true},
			{Name: []rune("table1.csv"), FormatAsIdentifier: true},
			{Name: []rune("LATERAL"), AppendSpace: true},
		},
	},
	{
		Name:     "SelectArgs After JOIN Keyword",
		Line:     "",
		OrigLine: "select 1 from t1 right join ",
		Index:    28,
		Expect: readline.CandidateList{
			{Name: []rune("CSV()")},
			{Name: []rune("FIXED()")},
			{Name: []rune("JSON()")},
			{Name: []rune("JSONL()")},
			{Name: []rune("JSON_TABLE()")},
			{Name: []rune("LTSV()")},
			{Name: []rune(filepath.Join(CompletionTestDir, "sub", "table2.csv")), FormatAsIdentifier: true},
			{Name: []rune("newtable.csv"), FormatAsIdentifier: true},
			{Name: []rune("tempview"), FormatAsIdentifier: true},
			{Name: []rune("."), FormatAsIdentifier: true},
			{Name: []rune(".."), FormatAsIdentifier: true},
			{Name: []rune("sub/"), FormatAsIdentifier: true},
			{Name: []rune("table1.csv"), FormatAsIdentifier: true},
		},
	},
	{
		Name:     "SelectArgs After JOIN Keyword",
		Line:     "",
		OrigLine: "select 1 from t1 full outer join ",
		Index:    33,
		Expect: readline.CandidateList{
			{Name: []rune("CSV()")},
			{Name: []rune("FIXED()")},
			{Name: []rune("JSON()")},
			{Name: []rune("JSONL()")},
			{Name: []rune("JSON_TABLE()")},
			{Name: []rune("LTSV()")},
			{Name: []rune(filepath.Join(CompletionTestDir, "sub", "table2.csv")), FormatAsIdentifier: true},
			{Name: []rune("newtable.csv"), FormatAsIdentifier: true},
			{Name: []rune("tempview"), FormatAsIdentifier: true},
			{Name: []rune("."), FormatAsIdentifier: true},
			{Name: []rune(".."), FormatAsIdentifier: true},
			{Name: []rune("sub/"), FormatAsIdentifier: true},
			{Name: []rune("table1.csv"), FormatAsIdentifier: true},
		},
	},
	{
		Name:     "SelectArgs After USING in From Clause",
		Line:     "",
		OrigLine: "select 1 from t1 inner join t2 using (",
		Index:    38,
		Expect:   readline.CandidateList{},
	},
	{
		Name:     "SelectArgs After ON in From Clause",
		Line:     "",
		OrigLine: "select 1 from t1 inner join t2 on ",
		Index:    34,
		Expect:   readline.CandidateList{},
	},
	{
		Name:     "SelectArgs First Item After ON in From Clause",
		Line:     "j",
		OrigLine: "select 1 from t1 inner join t2 on j",
		Index:    35,
		Expect: readline.CandidateList{
			{Name: []rune("JSON_ROW()")},
		},
	},
	{
		Name:     "SelectArgs After Values in From Clause",
		Line:     "",
		OrigLine: "select 1 from t1 inner join t2 on 1 + (1 + 2) ",
		Index:    46,
		Expect: readline.CandidateList{
			{Name: []rune("CROSS"), AppendSpace: true},
			{Name: []rune("EXCEPT"), AppendSpace: true},
			{Name: []rune("FETCH"), AppendSpace: true},
			{Name: []rune("FOR UPDATE")},
			{Name: []rune("FULL"), AppendSpace: true},
			{Name: []rune("GROUP BY"), AppendSpace: true},
			{Name: []rune("HAVING"), AppendSpace: true},
			{Name: []rune("INNER"), AppendSpace: true},
			{Name: []rune("INTERSECT"), AppendSpace: true},
			{Name: []rune("JOIN"), AppendSpace: true},
			{Name: []rune("LEFT"), AppendSpace: true},
			{Name: []rune("LIMIT"), AppendSpace: true},
			{Name: []rune("NATURAL"), AppendSpace: true},
			{Name: []rune("OFFSET"), AppendSpace: true},
			{Name: []rune("ORDER BY"), AppendSpace: true},
			{Name: []rune("RIGHT"), AppendSpace: true},
			{Name: []rune("UNION"), AppendSpace: true},
			{Name: []rune("WHERE"), AppendSpace: true},
		},
	},
	{
		Name:     "SelectArgs After WHERE",
		Line:     "",
		OrigLine: "select 1 from t1 where ",
		Index:    23,
		Expect:   readline.CandidateList{},
	},
	{
		Name:     "SelectArgs First Item in Where Clause",
		Line:     "j",
		OrigLine: "select 1 from t1 where j",
		Index:    24,
		Expect: readline.CandidateList{
			{Name: []rune("JSON_ROW()")},
		},
	},
	{
		Name:     "SelectArgs After First Item in Where Clause",
		Line:     "",
		OrigLine: "select 1 from t1 where 1 ",
		Index:    25,
		Expect: readline.CandidateList{
			{Name: []rune("EXCEPT"), AppendSpace: true},
			{Name: []rune("FETCH"), AppendSpace: true},
			{Name: []rune("FOR UPDATE")},
			{Name: []rune("GROUP BY"), AppendSpace: true},
			{Name: []rune("HAVING"), AppendSpace: true},
			{Name: []rune("INTERSECT"), AppendSpace: true},
			{Name: []rune("LIMIT"), AppendSpace: true},
			{Name: []rune("OFFSET"), AppendSpace: true},
			{Name: []rune("ORDER BY"), AppendSpace: true},
			{Name: []rune("UNION"), AppendSpace: true},
		},
	},
	{
		Name:     "SelectArgs After GROUP",
		Line:     "",
		OrigLine: "select 1 from t1 where 1 group ",
		Index:    31,
		Expect: readline.CandidateList{
			{Name: []rune("BY"), AppendSpace: true},
		},
	},
	{
		Name:     "SelectArgs First Item After GROUP BY",
		Line:     "@",
		OrigLine: "select 1 from t1 where 1 group by @",
		Index:    35,
		Expect: readline.CandidateList{
			{Name: []rune("@var1")},
			{Name: []rune("@var2")},
			{Name: []rune("JSON_ROW()")},
		},
	},
	{
		Name:     "SelectArgs After First Item in Group By Clause",
		Line:     "",
		OrigLine: "select 1 from t1 where 1 group by 1 ",
		Index:    36,
		Expect: readline.CandidateList{
			{Name: []rune("EXCEPT"), AppendSpace: true},
			{Name: []rune("FETCH"), AppendSpace: true},
			{Name: []rune("FOR UPDATE")},
			{Name: []rune("HAVING"), AppendSpace: true},
			{Name: []rune("INTERSECT"), AppendSpace: true},
			{Name: []rune("LIMIT"), AppendSpace: true},
			{Name: []rune("OFFSET"), AppendSpace: true},
			{Name: []rune("ORDER BY"), AppendSpace: true},
			{Name: []rune("UNION"), AppendSpace: true},
		},
	},
	{
		Name:     "SelectArgs After HAVING",
		Line:     "",
		OrigLine: "select 1 from t1 having ",
		Index:    24,
		Expect:   readline.CandidateList{},
	},
	{
		Name:     "SelectArgs After First Item in Having Clause",
		Line:     "",
		OrigLine: "select 1 from t1 having 1 ",
		Index:    26,
		Expect: readline.CandidateList{
			{Name: []rune("EXCEPT"), AppendSpace: true},
			{Name: []rune("FETCH"), AppendSpace: true},
			{Name: []rune("FOR UPDATE")},
			{Name: []rune("INTERSECT"), AppendSpace: true},
			{Name: []rune("LIMIT"), AppendSpace: true},
			{Name: []rune("OFFSET"), AppendSpace: true},
			{Name: []rune("ORDER BY"), AppendSpace: true},
			{Name: []rune("UNION"), AppendSpace: true},
		},
	},
	{
		Name:     "SelectArgs After ORDER",
		Line:     "",
		OrigLine: "select 1 from t1 where 1 order ",
		Index:    31,
		Expect: readline.CandidateList{
			{Name: []rune("BY"), AppendSpace: true},
		},
	},
	{
		Name:     "SelectArgs After ORDER BY",
		Line:     "@v",
		OrigLine: "select 1 from t1 where 1 order by @v",
		Index:    36,
		Expect: readline.CandidateList{
			{Name: []rune("@var1")},
			{Name: []rune("@var2")},
		},
	},
	{
		Name:     "SelectArgs After First Item in Order By Clause",
		Line:     "",
		OrigLine: "select 1 from t1 where 1 order by 1 ",
		Index:    36,
		Expect: readline.CandidateList{
			{Name: []rune("ASC")},
			{Name: []rune("DESC")},
			{Name: []rune("FETCH"), AppendSpace: true},
			{Name: []rune("FOR UPDATE")},
			{Name: []rune("LIMIT"), AppendSpace: true},
			{Name: []rune("NULLS FIRST")},
			{Name: []rune("NULLS LAST")},
			{Name: []rune("OFFSET"), AppendSpace: true},
		},
	},
	{
		Name:     "SelectArgs After DESC in Order By Clause",
		Line:     "",
		OrigLine: "select 1 from t1 where 1 order by 1 desc ",
		Index:    41,
		Expect: readline.CandidateList{
			{Name: []rune("FETCH"), AppendSpace: true},
			{Name: []rune("FOR UPDATE")},
			{Name: []rune("LIMIT"), AppendSpace: true},
			{Name: []rune("NULLS FIRST")},
			{Name: []rune("NULLS LAST")},
			{Name: []rune("OFFSET"), AppendSpace: true},
		},
	},
	{
		Name:     "SelectArgs After NULLS in Order By Clause",
		Line:     "",
		OrigLine: "select 1 from t1 where 1 order by 1 desc nulls ",
		Index:    47,
		Expect: readline.CandidateList{
			{Name: []rune("FIRST")},
			{Name: []rune("LAST")},
		},
	},
	{
		Name:     "SelectArgs After NULLS FIRST in Order By Clause",
		Line:     "",
		OrigLine: "select 1 from t1 where 1 order by 1 desc nulls first ",
		Index:    53,
		Expect: readline.CandidateList{
			{Name: []rune("FETCH"), AppendSpace: true},
			{Name: []rune("FOR UPDATE")},
			{Name: []rune("LIMIT"), AppendSpace: true},
			{Name: []rune("OFFSET"), AppendSpace: true},
		},
	},
	{
		Name:     "SelectArgs After LIMIT",
		Line:     "@v",
		OrigLine: "select 1 from t1 where 1 limit @v",
		Index:    33,
		Expect: readline.CandidateList{
			{Name: []rune("@var1")},
			{Name: []rune("@var2")},
		},
	},
	{
		Name:     "SelectArgs After Frist Item in Limit Clause",
		Line:     "",
		OrigLine: "select 1 from t1 where 1 limit 1 ",
		Index:    33,
		Expect: readline.CandidateList{
			{Name: []rune("FOR UPDATE")},
			{Name: []rune("OFFSET"), AppendSpace: true},
			{Name: []rune("PERCENT")},
			{Name: []rune("WITH TIES")},
		},
	},
	{
		Name:     "SelectArgs After PERCENT in Limit Clause",
		Line:     "",
		OrigLine: "select 1 from t1 where 1 limit 1 percent ",
		Index:    41,
		Expect: readline.CandidateList{
			{Name: []rune("FOR UPDATE")},
			{Name: []rune("OFFSET"), AppendSpace: true},
			{Name: []rune("WITH TIES")},
		},
	},
	{
		Name:     "SelectArgs After WITH in Limit Clause",
		Line:     "",
		OrigLine: "select 1 from t1 where 1 limit 1 percent with ",
		Index:    46,
		Expect: readline.CandidateList{
			{Name: []rune("TIES")},
		},
	},
	{
		Name:     "SelectArgs After WITH TIES in Limit Clause",
		Line:     "",
		OrigLine: "select 1 from t1 where 1 limit 1 percent with ties ",
		Index:    51,
		Expect: readline.CandidateList{
			{Name: []rune("FOR UPDATE")},
			{Name: []rune("OFFSET"), AppendSpace: true},
		},
	},
	{
		Name:     "SelectArgs After OFFSET",
		Line:     "@",
		OrigLine: "select 1 from t1 where 1 offset @",
		Index:    33,
		Expect: readline.CandidateList{
			{Name: []rune("@var1")},
			{Name: []rune("@var2")},
		},
	},
	{
		Name:     "SelectArgs After Value in OFFSET Clause",
		Line:     "",
		OrigLine: "select 1 from t1 where 1 offset 2 ",
		Index:    34,
		Expect: readline.CandidateList{
			{Name: []rune("FETCH"), AppendSpace: true},
			{Name: []rune("FOR UPDATE")},
			{Name: []rune("ROWS")},
		},
	},
	{
		Name:     "SelectArgs After Value in OFFSET Clause with LIMIT Clause",
		Line:     "",
		OrigLine: "select 1 from t1 where 1 limit 1 offset 2 ",
		Index:    42,
		Expect: readline.CandidateList{
			{Name: []rune("FOR UPDATE")},
			{Name: []rune("ROWS")},
		},
	},
	{
		Name:     "SelectArgs After 1 in OFFSET Clause",
		Line:     "",
		OrigLine: "select 1 from t1 where 1 offset 1 ",
		Index:    34,
		Expect: readline.CandidateList{
			{Name: []rune("FETCH"), AppendSpace: true},
			{Name: []rune("FOR UPDATE")},
			{Name: []rune("ROW")},
		},
	},
	{
		Name:     "SelectArgs After FETCH without OFFSET Clause",
		Line:     "",
		OrigLine: "select 1 from t1 where 1 fetch ",
		Index:    31,
		Expect: readline.CandidateList{
			{Name: []rune("FIRST"), AppendSpace: true},
		},
	},
	{
		Name:     "SelectArgs After FETCH",
		Line:     "",
		OrigLine: "select 1 from t1 where 1 offset 10 rows fetch ",
		Index:    46,
		Expect: readline.CandidateList{
			{Name: []rune("NEXT"), AppendSpace: true},
		},
	},
	{
		Name:     "SelectArgs After Value in FETCH Clause",
		Line:     "",
		OrigLine: "select 1 from t1 where 1 offset 10 rows fetch next 2 ",
		Index:    53,
		Expect: readline.CandidateList{
			{Name: []rune("PERCENT"), AppendSpace: true},
			{Name: []rune("ROWS"), AppendSpace: true},
		},
	},
	{
		Name:     "SelectArgs After 1 in FETCH Clause",
		Line:     "",
		OrigLine: "select 1 from t1 where 1 offset 10 rows fetch next 1 ",
		Index:    53,
		Expect: readline.CandidateList{
			{Name: []rune("PERCENT"), AppendSpace: true},
			{Name: []rune("ROW"), AppendSpace: true},
		},
	},
	{
		Name:     "SelectArgs After ROW in FETCH Clause",
		Line:     "",
		OrigLine: "select 1 from t1 where 1 offset 10 rows fetch next 1 row ",
		Index:    57,
		Expect: readline.CandidateList{
			{Name: []rune("FOR UPDATE")},
			{Name: []rune("ONLY")},
			{Name: []rune("WITH TIES")},
		},
	},
	{
		Name:     "SelectArgs After WITH in FETCH Clause",
		Line:     "",
		OrigLine: "select 1 from t1 where 1 offset 10 rows fetch next 1 row with ",
		Index:    62,
		Expect: readline.CandidateList{
			{Name: []rune("TIES")},
		},
	},
	{
		Name:     "SelectArgs After ONLY in FETCH Clause",
		Line:     "",
		OrigLine: "select 1 from t1 where 1 offset 10 rows fetch next 1 row only ",
		Index:    62,
		Expect: readline.CandidateList{
			{Name: []rune("FOR UPDATE")},
		},
	},
	{
		Name:     "SelectArgs After UNION",
		Line:     "",
		OrigLine: "select 1 from t1 where 1 union ",
		Index:    31,
		Expect: readline.CandidateList{
			{Name: []rune("ALL"), AppendSpace: true},
			{Name: []rune("SELECT"), AppendSpace: true},
		},
	},
	{
		Name:     "SelectArgs After UNION ALL",
		Line:     "",
		OrigLine: "select 1 from t1 where 1 union all ",
		Index:    35,
		Expect: readline.CandidateList{
			{Name: []rune("SELECT"), AppendSpace: true},
		},
	},
	{
		Name:     "SelectArgs After FOR",
		Line:     "",
		OrigLine: "select col1 for ",
		Index:    16,
		Expect: readline.CandidateList{
			{Name: []rune("UPDATE")},
		},
	},
}

func TestCompleter_SelectArgs(t *testing.T) {
	testCompleter(t, completer.SelectArgs, completerSelectArgsTests)
}

var completerInsertArgsTests = []completerTest{
	{
		Name:     "InsertArgs",
		Line:     "",
		OrigLine: "insert ",
		Index:    7,
		Expect: readline.CandidateList{
			{Name: []rune("INTO"), AppendSpace: true},
		},
	},
	{
		Name:     "InsertArgs After INTO",
		Line:     "",
		OrigLine: "insert into ",
		Index:    12,
		Expect: readline.CandidateList{
			{Name: []rune("CSV()"), AppendSpace: true},
			{Name: []rune("FIXED()"), AppendSpace: true},
			{Name: []rune("JSON()"), AppendSpace: true},
			{Name: []rune("JSONL()"), AppendSpace: true},
			{Name: []rune("LTSV()"), AppendSpace: true},
			{Name: []rune(filepath.Join(CompletionTestDir, "sub", "table2.csv")), FormatAsIdentifier: true, AppendSpace: true},
			{Name: []rune("newtable.csv"), FormatAsIdentifier: true, AppendSpace: true},
			{Name: []rune("tempview"), FormatAsIdentifier: true, AppendSpace: true},
			{Name: []rune("."), FormatAsIdentifier: true, AppendSpace: true},
			{Name: []rune(".."), FormatAsIdentifier: true, AppendSpace: true},
			{Name: []rune("sub/"), FormatAsIdentifier: true, AppendSpace: true},
			{Name: []rune("table1.csv"), FormatAsIdentifier: true, AppendSpace: true},
		},
	},
	{
		Name:     "InsertArgs After Table Name",
		Line:     "",
		OrigLine: "insert into tbl ",
		Index:    16,
		Expect: readline.CandidateList{
			{Name: []rune("SELECT"), AppendSpace: true},
			{Name: []rune("VALUES"), AppendSpace: true},
		},
	},
	{
		Name:     "InsertArgs In Column Enumulation",
		Line:     "",
		OrigLine: "insert into tbl (col, ",
		Index:    22,
		Expect:   readline.CandidateList{},
	},
	{
		Name:     "InsertArgs After Column Enumulation",
		Line:     "",
		OrigLine: "insert into tbl (col1, col2) ",
		Index:    29,
		Expect: readline.CandidateList{
			{Name: []rune("SELECT"), AppendSpace: true},
			{Name: []rune("VALUES"), AppendSpace: true},
		},
	},
	{
		Name:     "InsertArgs After VALUES",
		Line:     "@",
		OrigLine: "insert into tbl (col1, col2) values (@",
		Index:    38,
		Expect: readline.CandidateList{
			{Name: []rune("@var1")},
			{Name: []rune("@var2")},
			{Name: []rune("JSON_ROW()")},
		},
	},
	{
		Name:     "InsertArgs After SELECT",
		Line:     "fro",
		OrigLine: "insert into tbl select 1 fro",
		Index:    28,
		Expect: readline.CandidateList{
			{Name: []rune("AS"), AppendSpace: true},
			{Name: []rune("EXCEPT"), AppendSpace: true},
			{Name: []rune("FETCH"), AppendSpace: true},
			{Name: []rune("FOR UPDATE")},
			{Name: []rune("FROM"), AppendSpace: true},
			{Name: []rune("GROUP BY"), AppendSpace: true},
			{Name: []rune("HAVING"), AppendSpace: true},
			{Name: []rune("INTERSECT"), AppendSpace: true},
			{Name: []rune("LIMIT"), AppendSpace: true},
			{Name: []rune("OFFSET"), AppendSpace: true},
			{Name: []rune("ORDER BY"), AppendSpace: true},
			{Name: []rune("UNION"), AppendSpace: true},
			{Name: []rune("WHERE"), AppendSpace: true},
		},
	},
}

func TestCompleter_InsertArgs(t *testing.T) {
	testCompleter(t, completer.InsertArgs, completerInsertArgsTests)
}

var completerUpdateArgsTests = []completerTest{
	{
		Name:     "UpdateArgs",
		Line:     "",
		OrigLine: "update ",
		Index:    7,
		Expect: readline.CandidateList{
			{Name: []rune("CSV()")},
			{Name: []rune("FIXED()")},
			{Name: []rune("JSON()")},
			{Name: []rune("JSONL()")},
			{Name: []rune("LTSV()")},
			{Name: []rune(filepath.Join(CompletionTestDir, "sub", "table2.csv")), FormatAsIdentifier: true},
			{Name: []rune("newtable.csv"), FormatAsIdentifier: true},
			{Name: []rune("tempview"), FormatAsIdentifier: true},
			{Name: []rune("."), FormatAsIdentifier: true},
			{Name: []rune(".."), FormatAsIdentifier: true},
			{Name: []rune("sub/"), FormatAsIdentifier: true},
			{Name: []rune("table1.csv"), FormatAsIdentifier: true},
		},
	},
	{
		Name:     "UpdateArgs After Table Name",
		Line:     "",
		OrigLine: "update tbl ",
		Index:    11,
		Expect: readline.CandidateList{
			{Name: []rune("SET"), AppendSpace: true},
		},
	},
	{
		Name:     "UpdateArgs After Substitution Operator in Set Clause",
		Line:     "@",
		OrigLine: "update tbl set col = @",
		Index:    22,
		Expect: readline.CandidateList{
			{Name: []rune("@var1")},
			{Name: []rune("@var2")},
		},
	},
	{
		Name:     "UpdateArgs After Substitution in Set Clause",
		Line:     "",
		OrigLine: "update tbl set col = 1 ",
		Index:    23,
		Expect: readline.CandidateList{
			{Name: []rune("FROM"), AppendSpace: true},
			{Name: []rune("WHERE"), AppendSpace: true},
		},
	},
	{
		Name:     "UpdateArgs After First Item in From Clause",
		Line:     "",
		OrigLine: "update tbl set col = 1 from 1 ",
		Index:    30,
		Expect: readline.CandidateList{
			{Name: []rune("WHERE"), AppendSpace: true},
		},
	},
	{
		Name:     "UpdateArgs After WHERE",
		Line:     "@",
		OrigLine: "update tbl set col = 1 where @",
		Index:    30,
		Expect: readline.CandidateList{
			{Name: []rune("@var1")},
			{Name: []rune("@var2")},
			{Name: []rune("JSON_ROW()")},
		},
	},
}

func TestCompleter_UpdateArgs(t *testing.T) {
	testCompleter(t, completer.UpdateArgs, completerUpdateArgsTests)
}

var completerReplaceArgsTests = []completerTest{
	{
		Name:     "ReplaceArgs",
		Line:     "",
		OrigLine: "replace ",
		Index:    8,
		Expect: readline.CandidateList{
			{Name: []rune("INTO"), AppendSpace: true},
		},
	},
	{
		Name:     "ReplaceArgs After INTO",
		Line:     "",
		OrigLine: "replace into ",
		Index:    13,
		Expect: readline.CandidateList{
			{Name: []rune("CSV()"), AppendSpace: true},
			{Name: []rune("FIXED()"), AppendSpace: true},
			{Name: []rune("JSON()"), AppendSpace: true},
			{Name: []rune("JSONL()"), AppendSpace: true},
			{Name: []rune("LTSV()"), AppendSpace: true},
			{Name: []rune(filepath.Join(CompletionTestDir, "sub", "table2.csv")), FormatAsIdentifier: true, AppendSpace: true},
			{Name: []rune("newtable.csv"), FormatAsIdentifier: true, AppendSpace: true},
			{Name: []rune("tempview"), FormatAsIdentifier: true, AppendSpace: true},
			{Name: []rune("."), FormatAsIdentifier: true, AppendSpace: true},
			{Name: []rune(".."), FormatAsIdentifier: true, AppendSpace: true},
			{Name: []rune("sub/"), FormatAsIdentifier: true, AppendSpace: true},
			{Name: []rune("table1.csv"), FormatAsIdentifier: true, AppendSpace: true},
		},
	},
	{
		Name:     "ReplaceArgs After Table Name",
		Line:     "",
		OrigLine: "replace into tbl ",
		Index:    17,
		Expect: readline.CandidateList{
			{Name: []rune("USING ()"), AppendSpace: true},
		},
	},
	{
		Name:     "ReplaceArgs In Column Enumulation",
		Line:     "",
		OrigLine: "replace into tbl (col, ",
		Index:    23,
		Expect:   readline.CandidateList{},
	},
	{
		Name:     "ReplaceArgs After Column Enumulation",
		Line:     "",
		OrigLine: "replace into tbl (col1, col2) ",
		Index:    30,
		Expect: readline.CandidateList{
			{Name: []rune("USING ()"), AppendSpace: true},
		},
	},
	{
		Name:     "ReplaceArgs After USING",
		Line:     "",
		OrigLine: "replace into tbl (col1, col2) using (",
		Index:    37,
		Expect:   readline.CandidateList{},
	},
	{
		Name:     "ReplaceArgs After Key Enumulation",
		Line:     "",
		OrigLine: "replace into tbl (col1, col2) using (col1) ",
		Index:    43,
		Expect: readline.CandidateList{
			{Name: []rune("SELECT"), AppendSpace: true},
			{Name: []rune("VALUES"), AppendSpace: true},
		},
	},
	{
		Name:     "ReplaceArgs After VALUES",
		Line:     "@",
		OrigLine: "insert into tbl (col1, col2) using (col1) values (@",
		Index:    51,
		Expect: readline.CandidateList{
			{Name: []rune("@var1")},
			{Name: []rune("@var2")},
			{Name: []rune("JSON_ROW()")},
		},
	},
	{
		Name:     "ReplaceArgs After SELECT",
		Line:     "fro",
		OrigLine: "insert into tbl using (col1) select 1 fro",
		Index:    41,
		Expect: readline.CandidateList{
			{Name: []rune("AS"), AppendSpace: true},
			{Name: []rune("EXCEPT"), AppendSpace: true},
			{Name: []rune("FETCH"), AppendSpace: true},
			{Name: []rune("FOR UPDATE")},
			{Name: []rune("FROM"), AppendSpace: true},
			{Name: []rune("GROUP BY"), AppendSpace: true},
			{Name: []rune("HAVING"), AppendSpace: true},
			{Name: []rune("INTERSECT"), AppendSpace: true},
			{Name: []rune("LIMIT"), AppendSpace: true},
			{Name: []rune("OFFSET"), AppendSpace: true},
			{Name: []rune("ORDER BY"), AppendSpace: true},
			{Name: []rune("UNION"), AppendSpace: true},
			{Name: []rune("WHERE"), AppendSpace: true},
		},
	},
}

func TestCompleter_ReplaceArgs(t *testing.T) {
	testCompleter(t, completer.ReplaceArgs, completerReplaceArgsTests)
}

var completerDeleteArgsTests = []completerTest{
	{
		Name:     "DeleteArgs",
		Line:     "",
		OrigLine: "delete ",
		Index:    7,
		Expect: readline.CandidateList{
			{Name: []rune("FROM"), AppendSpace: true},
		},
	},
	{
		Name:     "DeleteArgs After FROM",
		Line:     "",
		OrigLine: "delete from ",
		Index:    12,
		Expect: readline.CandidateList{
			{Name: []rune("CSV()")},
			{Name: []rune("FIXED()")},
			{Name: []rune("JSON()")},
			{Name: []rune("JSONL()")},
			{Name: []rune("JSON_TABLE()")},
			{Name: []rune("LTSV()")},
			{Name: []rune(filepath.Join(CompletionTestDir, "sub", "table2.csv")), FormatAsIdentifier: true},
			{Name: []rune("newtable.csv"), FormatAsIdentifier: true},
			{Name: []rune("tempview"), FormatAsIdentifier: true},
			{Name: []rune("."), FormatAsIdentifier: true},
			{Name: []rune(".."), FormatAsIdentifier: true},
			{Name: []rune("sub/"), FormatAsIdentifier: true},
			{Name: []rune("table1.csv"), FormatAsIdentifier: true},
		},
	},
	{
		Name:     "DeleteArgs After First Item in From Clause",
		Line:     "",
		OrigLine: "delete from tbl ",
		Index:    16,
		Expect: readline.CandidateList{
			{Name: []rune("WHERE"), AppendSpace: true},
		},
	},
	{
		Name:     "DeleteArgs After FROM for Multiple Files",
		Line:     "",
		OrigLine: "delete t1 from ",
		Index:    15,
		Expect: readline.CandidateList{
			{Name: []rune("CSV()")},
			{Name: []rune("FIXED()")},
			{Name: []rune("JSON()")},
			{Name: []rune("JSONL()")},
			{Name: []rune("JSON_TABLE()")},
			{Name: []rune("LTSV()")},
			{Name: []rune(filepath.Join(CompletionTestDir, "sub", "table2.csv")), FormatAsIdentifier: true},
			{Name: []rune("newtable.csv"), FormatAsIdentifier: true},
			{Name: []rune("tempview"), FormatAsIdentifier: true},
			{Name: []rune("."), FormatAsIdentifier: true},
			{Name: []rune(".."), FormatAsIdentifier: true},
			{Name: []rune("sub/"), FormatAsIdentifier: true},
			{Name: []rune("table1.csv"), FormatAsIdentifier: true},
		},
	},
	{
		Name:     "DeleteArgs After First Item in From Clause for Multiple Files",
		Line:     "",
		OrigLine: "delete tbl from tbl ",
		Index:    20,
		Expect: readline.CandidateList{
			{Name: []rune("AS"), AppendSpace: true},
			{Name: []rune("CROSS"), AppendSpace: true},
			{Name: []rune("FULL"), AppendSpace: true},
			{Name: []rune("INNER"), AppendSpace: true},
			{Name: []rune("JOIN"), AppendSpace: true},
			{Name: []rune("LEFT"), AppendSpace: true},
			{Name: []rune("NATURAL"), AppendSpace: true},
			{Name: []rune("RIGHT"), AppendSpace: true},
			{Name: []rune("WHERE"), AppendSpace: true},
		},
	},
	{
		Name:     "DeleteArgs WHERE",
		Line:     "@",
		OrigLine: "delete from tbl where @",
		Index:    23,
		Expect: readline.CandidateList{
			{Name: []rune("@var1")},
			{Name: []rune("@var2")},
			{Name: []rune("JSON_ROW()")},
		},
	},
}

func TestCompleter_DeleteArgs(t *testing.T) {
	testCompleter(t, completer.DeleteArgs, completerDeleteArgsTests)
}

var completerCreateArgsTests = []completerTest{
	{
		Name:     "CreateArgs",
		Line:     "",
		OrigLine: "create ",
		Index:    7,
		Expect: readline.CandidateList{
			{Name: []rune("TABLE"), AppendSpace: true},
		},
	},
	{
		Name:     "CreateArgs After Table Name",
		Line:     "",
		OrigLine: "create table newtable ",
		Index:    22,
		Expect: readline.CandidateList{
			{Name: []rune("AS"), AppendSpace: true},
			{Name: []rune("SELECT"), AppendSpace: true},
		},
	},
	{
		Name:     "CreateArgs After Column Enumulation",
		Line:     "",
		OrigLine: "create table newtable (col1, col2) ",
		Index:    35,
		Expect: readline.CandidateList{
			{Name: []rune("AS"), AppendSpace: true},
			{Name: []rune("SELECT"), AppendSpace: true},
		},
	},
	{
		Name:     "CreateArgs After AS",
		Line:     "",
		OrigLine: "create table newtable (col1, col2) as ",
		Index:    38,
		Expect: readline.CandidateList{
			{Name: []rune("SELECT"), AppendSpace: true},
		},
	},
	{
		Name:     "CreateArgs After SELECT",
		Line:     "fro",
		OrigLine: "create table select 1 fro",
		Index:    25,
		Expect: readline.CandidateList{
			{Name: []rune("AS"), AppendSpace: true},
			{Name: []rune("EXCEPT"), AppendSpace: true},
			{Name: []rune("FETCH"), AppendSpace: true},
			{Name: []rune("FOR UPDATE")},
			{Name: []rune("FROM"), AppendSpace: true},
			{Name: []rune("GROUP BY"), AppendSpace: true},
			{Name: []rune("HAVING"), AppendSpace: true},
			{Name: []rune("INTERSECT"), AppendSpace: true},
			{Name: []rune("LIMIT"), AppendSpace: true},
			{Name: []rune("OFFSET"), AppendSpace: true},
			{Name: []rune("ORDER BY"), AppendSpace: true},
			{Name: []rune("UNION"), AppendSpace: true},
			{Name: []rune("WHERE"), AppendSpace: true},
		},
	},
	{
		Name:     "CreateArgs Ignore AS in select query",
		Line:     "fro",
		OrigLine: "create table select 1 as fro",
		Index:    28,
		Expect: readline.CandidateList{
			{Name: []rune("AS"), AppendSpace: true},
			{Name: []rune("EXCEPT"), AppendSpace: true},
			{Name: []rune("FETCH"), AppendSpace: true},
			{Name: []rune("FOR UPDATE")},
			{Name: []rune("FROM"), AppendSpace: true},
			{Name: []rune("GROUP BY"), AppendSpace: true},
			{Name: []rune("HAVING"), AppendSpace: true},
			{Name: []rune("INTERSECT"), AppendSpace: true},
			{Name: []rune("LIMIT"), AppendSpace: true},
			{Name: []rune("OFFSET"), AppendSpace: true},
			{Name: []rune("ORDER BY"), AppendSpace: true},
			{Name: []rune("UNION"), AppendSpace: true},
			{Name: []rune("WHERE"), AppendSpace: true},
		},
	},
}

func TestCompleter_CreateArgs(t *testing.T) {
	testCompleter(t, completer.CreateArgs, completerCreateArgsTests)
}

var completerAlterArgsTests = []completerTest{
	{
		Name:     "AlterArgs",
		Line:     "",
		OrigLine: "alter ",
		Index:    6,
		Expect: readline.CandidateList{
			{Name: []rune("TABLE"), AppendSpace: true},
		},
	},
	{
		Name:     "AlterArgs After TABLE",
		Line:     "",
		OrigLine: "alter table ",
		Index:    12,
		Expect: readline.CandidateList{
			{Name: []rune("CSV()"), AppendSpace: true},
			{Name: []rune("FIXED()"), AppendSpace: true},
			{Name: []rune("JSON()"), AppendSpace: true},
			{Name: []rune("JSONL()"), AppendSpace: true},
			{Name: []rune("LTSV()"), AppendSpace: true},
			{Name: []rune(filepath.Join(CompletionTestDir, "sub", "table2.csv")), FormatAsIdentifier: true, AppendSpace: true},
			{Name: []rune("newtable.csv"), FormatAsIdentifier: true, AppendSpace: true},
			{Name: []rune("tempview"), FormatAsIdentifier: true, AppendSpace: true},
			{Name: []rune("."), FormatAsIdentifier: true, AppendSpace: true},
			{Name: []rune(".."), FormatAsIdentifier: true, AppendSpace: true},
			{Name: []rune("sub/"), FormatAsIdentifier: true, AppendSpace: true},
			{Name: []rune("table1.csv"), FormatAsIdentifier: true, AppendSpace: true},
		},
	},
	{
		Name:     "AlterArgs After Table Name",
		Line:     "",
		OrigLine: "alter table `newtable.csv`",
		Index:    26,
		Expect: readline.CandidateList{
			{Name: []rune("ADD"), AppendSpace: true},
			{Name: []rune("DROP"), AppendSpace: true},
			{Name: []rune("RENAME"), AppendSpace: true},
			{Name: []rune("SET"), AppendSpace: true},
		},
	},
	{
		Name:     "AlterArgs After ADD",
		Line:     "",
		OrigLine: "alter table `newtable.csv` add ",
		Index:    31,
		Expect:   readline.CandidateList{},
	},
	{
		Name:     "AlterArgs After New Column Name",
		Line:     "",
		OrigLine: "alter table `newtable.csv` add newcol ",
		Index:    38,
		Expect: readline.CandidateList{
			{Name: []rune("DEFAULT"), AppendSpace: true},
			{Name: []rune("FIRST"), AppendSpace: true},
			{Name: []rune("LAST"), AppendSpace: true},
			{Name: []rune("AFTER"), AppendSpace: true},
			{Name: []rune("BEFORE"), AppendSpace: true},
		},
	},
	{
		Name:     "AlterArgs After DEFAULT",
		Line:     "@",
		OrigLine: "alter table `newtable.csv` add newcol default @",
		Index:    47,
		Expect: readline.CandidateList{
			{Name: []rune("@var1")},
			{Name: []rune("@var2")},
		},
	},
	{
		Name:     "AlterArgs After BEFORE",
		Line:     "",
		OrigLine: "alter table `newtable.csv` add newcol default 1 before ",
		Index:    55,
		Expect: readline.CandidateList{
			{Name: []rune("ncol1"), FormatAsIdentifier: true},
			{Name: []rune("ncol2"), FormatAsIdentifier: true},
			{Name: []rune("ncol3"), FormatAsIdentifier: true},
		},
	},
	{
		Name:     "AlterArgs After Column Name in Brackets",
		Line:     "",
		OrigLine: "alter table `newtable.csv` add (newcol ",
		Index:    39,
		Expect: readline.CandidateList{
			{Name: []rune("DEFAULT"), AppendSpace: true},
		},
	},
	{
		Name:     "AlterArgs After DEFAULT in Brackets",
		Line:     "@",
		OrigLine: "alter table `newtable.csv` add (newcol default @",
		Index:    48,
		Expect: readline.CandidateList{
			{Name: []rune("@var1")},
			{Name: []rune("@var2")},
		},
	},
	{
		Name:     "AlterArgs After Multiple Columns Specified",
		Line:     "",
		OrigLine: "alter table `newtable.csv` add (newcol default 1, newcol2) ",
		Index:    59,
		Expect: readline.CandidateList{
			{Name: []rune("FIRST"), AppendSpace: true},
			{Name: []rune("LAST"), AppendSpace: true},
			{Name: []rune("AFTER"), AppendSpace: true},
			{Name: []rune("BEFORE"), AppendSpace: true},
		},
	},
	{
		Name:     "AlterArgs After DROP",
		Line:     "",
		OrigLine: "alter table `newtable.csv` drop ",
		Index:    32,
		Expect: readline.CandidateList{
			{Name: []rune("ncol1"), FormatAsIdentifier: true},
			{Name: []rune("ncol2"), FormatAsIdentifier: true},
			{Name: []rune("ncol3"), FormatAsIdentifier: true},
		},
	},
	{
		Name:     "AlterArgs After DROP in Brackets",
		Line:     "",
		OrigLine: "alter table `newtable.csv` drop (",
		Index:    33,
		Expect: readline.CandidateList{
			{Name: []rune("ncol1"), FormatAsIdentifier: true},
			{Name: []rune("ncol2"), FormatAsIdentifier: true},
			{Name: []rune("ncol3"), FormatAsIdentifier: true},
		},
	},
	{
		Name:     "AlterArgs After RENAME",
		Line:     "",
		OrigLine: "alter table `newtable.csv` rename ",
		Index:    34,
		Expect: readline.CandidateList{
			{Name: []rune("ncol1"), FormatAsIdentifier: true, AppendSpace: true},
			{Name: []rune("ncol2"), FormatAsIdentifier: true, AppendSpace: true},
			{Name: []rune("ncol3"), FormatAsIdentifier: true, AppendSpace: true},
		},
	},
	{
		Name:     "AlterArgs After Rename Column Name",
		Line:     "",
		OrigLine: "alter table `newtable.csv` rename ncol1 ",
		Index:    40,
		Expect: readline.CandidateList{
			{Name: []rune("TO"), AppendSpace: true},
		},
	},
	{
		Name:     "AlterArgs After SET",
		Line:     "",
		OrigLine: "alter table `newtable.csv` set ",
		Index:    31,
		Expect: readline.CandidateList{
			{Name: []rune("DELIMITER"), AppendSpace: true},
			{Name: []rune("DELIMITER_POSITIONS"), AppendSpace: true},
			{Name: []rune("ENCLOSE_ALL"), AppendSpace: true},
			{Name: []rune("ENCODING"), AppendSpace: true},
			{Name: []rune("FORMAT"), AppendSpace: true},
			{Name: []rune("HEADER"), AppendSpace: true},
			{Name: []rune("JSON_ESCAPE"), AppendSpace: true},
			{Name: []rune("LINE_BREAK"), AppendSpace: true},
			{Name: []rune("PRETTY_PRINT"), AppendSpace: true},
		},
	},
	{
		Name:     "AlterArgs After Set Item",
		Line:     "",
		OrigLine: "alter table `newtable.csv` set delimiter ",
		Index:    41,
		Expect: readline.CandidateList{
			{Name: []rune("TO"), AppendSpace: true},
		},
	},
	{
		Name:     "AlterArgs Set Delimiter Values",
		Line:     "",
		OrigLine: "alter table `newtable.csv` set delimiter to ",
		Index:    43,
		Expect: readline.CandidateList{
			{Name: []rune("','")},
			{Name: []rune("'\\t'")},
		},
	},
	{
		Name:     "AlterArgs Set DelimiterPositions Values",
		Line:     "",
		OrigLine: "alter table `newtable.csv` set delimiter_positions to ",
		Index:    53,
		Expect: readline.CandidateList{
			{Name: []rune("'SPACES'")},
			{Name: []rune("'S[]'")},
			{Name: []rune("'[]'")},
		},
	},
	{
		Name:     "AlterArgs Set Format Values",
		Line:     "",
		OrigLine: "alter table `newtable.csv` set format to ",
		Index:    40,
		Expect: readline.CandidateList{
			{Name: []rune("BOX")},
			{Name: []rune("CSV")},
			{Name: []rune("FIXED")},
			{Name: []rune("GFM")},
			{Name: []rune("JSON")},
			{Name: []rune("JSONL")},
			{Name: []rune("LTSV")},
			{Name: []rune("ORG")},
			{Name: []rune("TEXT")},
			{Name: []rune("TSV")},
		},
	},
	{
		Name:     "AlterArgs Set Encoding Values",
		Line:     "",
		OrigLine: "alter table `newtable.csv` set encoding to ",
		Index:    42,
		Expect: readline.CandidateList{
			{Name: []rune("SJIS")},
			{Name: []rune("UTF16")},
			{Name: []rune("UTF16BE")},
			{Name: []rune("UTF16BEM")},
			{Name: []rune("UTF16LE")},
			{Name: []rune("UTF16LEM")},
			{Name: []rune("UTF8")},
			{Name: []rune("UTF8M")},
		},
	},
	{
		Name:     "AlterArgs Set LineBreak Values",
		Line:     "",
		OrigLine: "alter table `newtable.csv` set line_break to ",
		Index:    44,
		Expect: readline.CandidateList{
			{Name: []rune("CR")},
			{Name: []rune("CRLF")},
			{Name: []rune("LF")},
		},
	},
	{
		Name:     "AlterArgs Set JsonEscape Values",
		Line:     "",
		OrigLine: "alter table `newtable.csv` set json_escape to ",
		Index:    45,
		Expect: readline.CandidateList{
			{Name: []rune("BACKSLASH")},
			{Name: []rune("HEX")},
			{Name: []rune("HEXALL")},
		},
	},
	{
		Name:     "AlterArgs Set Header Values",
		Line:     "",
		OrigLine: "alter table `newtable.csv` set header to ",
		Index:    40,
		Expect: readline.CandidateList{
			{Name: []rune("TRUE")},
			{Name: []rune("FALSE")},
		},
	},
}

func TestCompleter_AlterArgs(t *testing.T) {
	testCompleter(t, completer.AlterArgs, completerAlterArgsTests)
}

var completerDeclareArgsTests = []completerTest{
	{
		Name:     "DeclareArgs",
		Line:     "",
		OrigLine: "declare ",
		Index:    8,
		Expect:   readline.CandidateList{},
	},
	{
		Name:     "DeclareArgs Declare Variable",
		Line:     "@",
		OrigLine: "declare @var := @",
		Index:    17,
		Expect: readline.CandidateList{
			{Name: []rune("@var1")},
			{Name: []rune("@var2")},
		},
	},
	{
		Name:     "DeclareArgs After Variable",
		Line:     "",
		OrigLine: "declare @var ",
		Index:    13,
		Expect:   readline.CandidateList{},
	},
	{
		Name:     "DeclareArgs Declare Variable with VAR",
		Line:     "@",
		OrigLine: "var @var := @",
		Index:    13,
		Expect: readline.CandidateList{
			{Name: []rune("@var1")},
			{Name: []rune("@var2")},
		},
	},
	{
		Name:     "DeclareArgs Declare Variable with DECLARE",
		Line:     "@",
		OrigLine: "declare @var := @",
		Index:    17,
		Expect: readline.CandidateList{
			{Name: []rune("@var1")},
			{Name: []rune("@var2")},
		},
	},
	{
		Name:     "DeclareArgs After Identifier",
		Line:     "",
		OrigLine: "declare cur ",
		Index:    12,
		Expect: readline.CandidateList{
			{Name: []rune("AGGREGATE"), AppendSpace: true},
			{Name: []rune("CURSOR"), AppendSpace: true},
			{Name: []rune("FUNCTION"), AppendSpace: true},
			{Name: []rune("VIEW"), AppendSpace: true},
		},
	},
	{
		Name:     "DeclareArgs After VIEW",
		Line:     "",
		OrigLine: "declare v view ",
		Index:    15,
		Expect: readline.CandidateList{
			{Name: []rune("AS"), AppendSpace: true},
		},
	},
	{
		Name:     "DeclareArgs After AS in View Declaration",
		Line:     "",
		OrigLine: "declare v view as ",
		Index:    18,
		Expect: readline.CandidateList{
			{Name: []rune("SELECT"), AppendSpace: true},
		},
	},
	{
		Name:     "DeclareArgs After CURSOR",
		Line:     "",
		OrigLine: "declare cur cursor ",
		Index:    19,
		Expect: readline.CandidateList{
			{Name: []rune("FOR"), AppendSpace: true},
		},
	},
	{
		Name:     "DeclareArgs After FOR",
		Line:     "",
		OrigLine: "declare cur cursor for ",
		Index:    23,
		Expect: readline.CandidateList{
			{Name: []rune("SELECT"), AppendSpace: true},
			{Name: []rune("stmt")},
		},
	},
	{
		Name:     "DeclareArgs After SELECT",
		Line:     "fro",
		OrigLine: "declare cur cursor for select 1 fro",
		Index:    35,
		Expect: readline.CandidateList{
			{Name: []rune("AS"), AppendSpace: true},
			{Name: []rune("EXCEPT"), AppendSpace: true},
			{Name: []rune("FETCH"), AppendSpace: true},
			{Name: []rune("FOR UPDATE")},
			{Name: []rune("FROM"), AppendSpace: true},
			{Name: []rune("GROUP BY"), AppendSpace: true},
			{Name: []rune("HAVING"), AppendSpace: true},
			{Name: []rune("INTERSECT"), AppendSpace: true},
			{Name: []rune("LIMIT"), AppendSpace: true},
			{Name: []rune("OFFSET"), AppendSpace: true},
			{Name: []rune("ORDER BY"), AppendSpace: true},
			{Name: []rune("UNION"), AppendSpace: true},
			{Name: []rune("WHERE"), AppendSpace: true},
		},
	},
}

func TestCompleter_DeclareArgs(t *testing.T) {
	testCompleter(t, completer.DeclareArgs, completerDeclareArgsTests)
}

var completerPrepareArgsTests = []completerTest{
	{
		Name:     "PrepareArgs",
		Line:     "",
		OrigLine: "prepare ",
		Index:    8,
		Expect:   readline.CandidateList{},
	},
	{
		Name:     "PrepareArgs After Statement Name",
		Line:     "",
		OrigLine: "prepare stmt ",
		Index:    13,
		Expect: readline.CandidateList{
			{Name: []rune("FROM"), AppendSpace: true},
		},
	},
}

func TestCompleter_PrepareArgs(t *testing.T) {
	testCompleter(t, completer.PrepareArgs, completerPrepareArgsTests)
}

var completerFetchArgsTests = []completerTest{
	{
		Name:     "FetchArgs",
		Line:     "",
		OrigLine: "fetch ",
		Index:    6,
		Expect: readline.CandidateList{
			{Name: []rune("cur1"), AppendSpace: true},
			{Name: []rune("cur2"), AppendSpace: true},
			{Name: []rune("NEXT"), AppendSpace: true},
			{Name: []rune("PRIOR"), AppendSpace: true},
			{Name: []rune("FIRST"), AppendSpace: true},
			{Name: []rune("LAST"), AppendSpace: true},
			{Name: []rune("ABSOLUTE"), AppendSpace: true},
			{Name: []rune("RELATIVE"), AppendSpace: true},
		},
	},
	{
		Name:     "FetchArgs After Cursor Name",
		Line:     "",
		OrigLine: "fetch cur1 ",
		Index:    11,
		Expect: readline.CandidateList{
			{Name: []rune("INTO"), AppendSpace: true},
		},
	},
	{
		Name:     "FetchArgs After NEXT",
		Line:     "",
		OrigLine: "fetch next ",
		Index:    11,
		Expect: readline.CandidateList{
			{Name: []rune("cur1"), AppendSpace: true},
			{Name: []rune("cur2"), AppendSpace: true},
		},
	},
	{
		Name:     "FetchArgs After Cursor Name with NEXT",
		Line:     "",
		OrigLine: "fetch next cur1 ",
		Index:    16,
		Expect: readline.CandidateList{
			{Name: []rune("INTO"), AppendSpace: true},
		},
	},
	{
		Name:     "FetchArgs After ABSOLUTE",
		Line:     "@",
		OrigLine: "fetch absolute @",
		Index:    16,
		Expect: readline.CandidateList{
			{Name: []rune("@var1"), AppendSpace: true},
			{Name: []rune("@var2"), AppendSpace: true},
		},
	},
	{
		Name:     "FetchArgs After Absolute Value",
		Line:     "",
		OrigLine: "fetch absolute 1 ",
		Index:    17,
		Expect: readline.CandidateList{
			{Name: []rune("cur1"), AppendSpace: true},
			{Name: []rune("cur2"), AppendSpace: true},
		},
	},
	{
		Name:     "FetchArgs After Cursor Name with Absolute Value",
		Line:     "",
		OrigLine: "fetch absolute 1 cur1 ",
		Index:    22,
		Expect: readline.CandidateList{
			{Name: []rune("INTO"), AppendSpace: true},
		},
	},
	{
		Name:     "FetchArgs After INTO",
		Line:     "@",
		OrigLine: "fetch absolute 1 cur1 into @",
		Index:    28,
		Expect: readline.CandidateList{
			{Name: []rune("@var1")},
			{Name: []rune("@var2")},
		},
	},
	{
		Name:     "FetchArgs Invalid Syntax",
		Line:     "",
		OrigLine: "fetch cur1 with 1 ",
		Index:    18,
		Expect:   readline.CandidateList{},
	},
	{
		Name:     "FetchArgs Invalid Syntax with NEXT",
		Line:     "",
		OrigLine: "fetch next cur1 with 1 ",
		Index:    22,
		Expect:   readline.CandidateList{},
	},
}

func TestCompleter_FetchArgs(t *testing.T) {
	testCompleter(t, completer.FetchArgs, completerFetchArgsTests)
}

var completerSetArgsTests = []completerTest{
	{
		Name:     "SetArgs",
		Line:     "",
		OrigLine: "set ",
		Index:    4,
		Expect: append(completer.candidateList(completer.flagList, true), readline.CandidateList{
			{Name: []rune("@%ENV1"), AppendSpace: true},
			{Name: []rune("@%ENV2"), AppendSpace: true},
		}...),
	},
	{
		Name:     "SetArgs After Set Item",
		Line:     "",
		OrigLine: "set @flag ",
		Index:    10,
		Expect: readline.CandidateList{
			{Name: []rune("TO"), AppendSpace: true},
		},
	},
	{
		Name:     "SetArgs After TO for Repository Flag",
		Line:     "",
		OrigLine: "set @@repository to ",
		Index:    20,
		Expect: readline.CandidateList{
			{Name: []rune("."), FormatAsIdentifier: true},
			{Name: []rune(".."), FormatAsIdentifier: true},
			{Name: []rune("sub/"), FormatAsIdentifier: true},
		},
	},
	{
		Name:     "SetArgs After TO for Timezone Flag",
		Line:     "",
		OrigLine: "set @@timezone to ",
		Index:    18,
		Expect: readline.CandidateList{
			{Name: []rune("Local")},
			{Name: []rune("UTC")},
		},
	},
	{
		Name:     "SetArgs After TO for Import Format Flag",
		Line:     "",
		OrigLine: "set @@import_format to ",
		Index:    23,
		Expect: readline.CandidateList{
			{Name: []rune("CSV")},
			{Name: []rune("FIXED")},
			{Name: []rune("JSON")},
			{Name: []rune("JSONL")},
			{Name: []rune("LTSV")},
			{Name: []rune("TSV")},
		},
	},
	{
		Name:     "SetArgs After TO for Delimiter Flag",
		Line:     "",
		OrigLine: "set @@delimiter to ",
		Index:    19,
		Expect: readline.CandidateList{
			{Name: []rune("','")},
			{Name: []rune("'\\t'")},
		},
	},
	{
		Name:     "SetArgs After TO for Delimiter Positions Flag",
		Line:     "",
		OrigLine: "set @@delimiter_positions to ",
		Index:    29,
		Expect: readline.CandidateList{
			{Name: []rune("'SPACES'")},
			{Name: []rune("'S[]'")},
			{Name: []rune("'[]'")},
		},
	},
	{
		Name:     "SetArgs After TO for Encoding Flag",
		Line:     "",
		OrigLine: "set @@encoding to ",
		Index:    18,
		Expect: readline.CandidateList{
			{Name: []rune("AUTO")},
			{Name: []rune("SJIS")},
			{Name: []rune("UTF16")},
			{Name: []rune("UTF16BE")},
			{Name: []rune("UTF16BEM")},
			{Name: []rune("UTF16LE")},
			{Name: []rune("UTF16LEM")},
			{Name: []rune("UTF8")},
			{Name: []rune("UTF8M")},
		},
	},
	{
		Name:     "SetArgs After TO for Write-Encoding Flag",
		Line:     "",
		OrigLine: "set @@write_encoding to ",
		Index:    24,
		Expect: readline.CandidateList{
			{Name: []rune("SJIS")},
			{Name: []rune("UTF16")},
			{Name: []rune("UTF16BE")},
			{Name: []rune("UTF16BEM")},
			{Name: []rune("UTF16LE")},
			{Name: []rune("UTF16LEM")},
			{Name: []rune("UTF8")},
			{Name: []rune("UTF8M")},
		},
	},
	{
		Name:     "SetArgs After TO for Color Flag",
		Line:     "",
		OrigLine: "set @@color to ",
		Index:    15,
		Expect: readline.CandidateList{
			{Name: []rune("TRUE")},
			{Name: []rune("FALSE")},
		},
	},
	{
		Name:     "SetArgs After TO for Format Flag",
		Line:     "",
		OrigLine: "set @@format to ",
		Index:    16,
		Expect: readline.CandidateList{
			{Name: []rune("BOX")},
			{Name: []rune("CSV")},
			{Name: []rune("FIXED")},
			{Name: []rune("GFM")},
			{Name: []rune("JSON")},
			{Name: []rune("JSONL")},
			{Name: []rune("LTSV")},
			{Name: []rune("ORG")},
			{Name: []rune("TEXT")},
			{Name: []rune("TSV")},
		},
	},
	{
		Name:     "SetArgs After TO for LineBreak Flag",
		Line:     "",
		OrigLine: "set @@line_break to ",
		Index:    20,
		Expect: readline.CandidateList{
			{Name: []rune("CR")},
			{Name: []rune("CRLF")},
			{Name: []rune("LF")},
		},
	},
	{
		Name:     "SetArgs After TO for Json Escape Flag",
		Line:     "",
		OrigLine: "set @@json_escape to ",
		Index:    21,
		Expect: readline.CandidateList{
			{Name: []rune("BACKSLASH")},
			{Name: []rune("HEX")},
			{Name: []rune("HEXALL")},
		},
	},
	{
		Name:     "SetArgs After TO",
		Line:     "@",
		OrigLine: "set @%ENV1 to @",
		Index:    15,
		Expect: readline.CandidateList{
			{Name: []rune("@var1")},
			{Name: []rune("@var2")},
		},
	},
	{
		Name:     "SetArgs Invalid Syntax",
		Line:     "",
		OrigLine: "set @%ENV1 1 ",
		Index:    13,
		Expect:   readline.CandidateList{},
	},
}

func TestCompleter_SetArgs(t *testing.T) {
	testCompleter(t, completer.SetArgs, completerSetArgsTests)
}

var completerUsingArgsTests = []completerTest{
	{
		Name:     "UsingArgs After EXECUTE for Prepared Statement",
		Line:     "",
		OrigLine: "execute ",
		Index:    8,
		Expect: readline.CandidateList{
			{Name: []rune("stmt"), AppendSpace: true},
		},
	},
	{
		Name:     "UsingArgs After EXECUTE",
		Line:     "@",
		OrigLine: "execute @",
		Index:    9,
		Expect: readline.CandidateList{
			{Name: []rune("@var1")},
			{Name: []rune("@var2")},
		},
	},
	{
		Name:     "UsingArgs After Value",
		Line:     "",
		OrigLine: "execute '%s' ",
		Index:    13,
		Expect: readline.CandidateList{
			{Name: []rune("USING"), AppendSpace: true},
		},
	},
	{
		Name:     "UsingArgs After Value with PRINTF",
		Line:     "",
		OrigLine: "printf '%s' ",
		Index:    12,
		Expect: readline.CandidateList{
			{Name: []rune("USING"), AppendSpace: true},
		},
	},
	{
		Name:     "UsingArgs After USING",
		Line:     "@",
		OrigLine: "execute '%s' using @",
		Index:    20,
		Expect: readline.CandidateList{
			{Name: []rune("@var1")},
			{Name: []rune("@var2")},
		},
	},
	{
		Name:     "UsingArgs After Value in USING",
		Line:     "",
		OrigLine: "execute '%s' using 'a'",
		Index:    22,
		Expect: readline.CandidateList{
			{Name: []rune("AS"), AppendSpace: true},
		},
	},
}

func TestCompleter_UsingArgs(t *testing.T) {
	testCompleter(t, completer.UsingArgs, completerUsingArgsTests)
}

var completerAddFlagArgsTests = []completerTest{
	{
		Name:     "AddFlagArgs",
		Line:     "",
		OrigLine: "add ",
		Index:    4,
		Expect:   readline.CandidateList{},
	},
	{
		Name:     "AddFlagArgs After ADD",
		Line:     "@",
		OrigLine: "add @",
		Index:    5,
		Expect: readline.CandidateList{
			{Name: []rune("@var1")},
			{Name: []rune("@var2")},
		},
	},
	{
		Name:     "AddFlagArgs After Value",
		Line:     "",
		OrigLine: "add '%Y' ",
		Index:    9,
		Expect: readline.CandidateList{
			{Name: []rune("TO"), AppendSpace: true},
		},
	},
	{
		Name:     "AddFlagArgs After TO",
		Line:     "",
		OrigLine: "add '%Y' to ",
		Index:    12,
		Expect: readline.CandidateList{
			{Name: []rune("@@DATETIME_FORMAT")},
		},
	},
}

func TestCompleter_AddFlagArgs(t *testing.T) {
	testCompleter(t, completer.AddFlagArgs, completerAddFlagArgsTests)
}

var completerRemoveFlagArgsTests = []completerTest{
	{
		Name:     "RemoveFlagArgs",
		Line:     "",
		OrigLine: "remove ",
		Index:    7,
		Expect:   readline.CandidateList{},
	},
	{
		Name:     "AddFlagArgs After REMOVE",
		Line:     "@",
		OrigLine: "remove @",
		Index:    8,
		Expect: readline.CandidateList{
			{Name: []rune("@var1")},
			{Name: []rune("@var2")},
		},
	},
	{
		Name:     "RemoveFlagArgs After Value",
		Line:     "",
		OrigLine: "remove 1 ",
		Index:    9,
		Expect: readline.CandidateList{
			{Name: []rune("FROM"), AppendSpace: true},
		},
	},
	{
		Name:     "RemoveFlagArgs After FROM",
		Line:     "",
		OrigLine: "remove 1 from ",
		Index:    14,
		Expect: readline.CandidateList{
			{Name: []rune("@@DATETIME_FORMAT")},
		},
	},
}

func TestCompleter_RemoveFlagArgs(t *testing.T) {
	testCompleter(t, completer.RemoveFlagArgs, completerRemoveFlagArgsTests)
}

var completerDisposeArgsTests = []completerTest{
	{
		Name:     "DisposeArgs",
		Line:     "",
		OrigLine: "dispose ",
		Index:    8,
		Expect: readline.CandidateList{
			{Name: []rune("CURSOR"), AppendSpace: true},
			{Name: []rune("FUNCTION"), AppendSpace: true},
			{Name: []rune("PREPARE"), AppendSpace: true},
			{Name: []rune("VIEW"), AppendSpace: true},
			{Name: []rune("@var1")},
			{Name: []rune("@var2")},
		},
	},
	{
		Name:     "DisposeArgs After CURSOR",
		Line:     "",
		OrigLine: "dispose cursor ",
		Index:    15,
		Expect: readline.CandidateList{
			{Name: []rune("cur1")},
			{Name: []rune("cur2")},
		},
	},
	{
		Name:     "DisposeArgs After FUNCTION",
		Line:     "",
		OrigLine: "dispose function ",
		Index:    17,
		Expect: readline.CandidateList{
			{Name: []rune("userfunc")},
			{Name: []rune("aggfunc")},
		},
	},
	{
		Name:     "DisposeArgs After VIEW",
		Line:     "",
		OrigLine: "dispose view ",
		Index:    13,
		Expect: readline.CandidateList{
			{Name: []rune("tempview")},
		},
	},
	{
		Name:     "DisposeArgs After PREPARE",
		Line:     "",
		OrigLine: "dispose prepare ",
		Index:    16,
		Expect: readline.CandidateList{
			{Name: []rune("stmt")},
		},
	},
	{
		Name:     "DisposeArgs Invalid Syntax",
		Line:     "",
		OrigLine: "dispose 1 ",
		Index:    10,
		Expect:   readline.CandidateList{},
	},
}

func TestCompleter_DisposeArgs(t *testing.T) {
	testCompleter(t, completer.DisposeArgs, completerDisposeArgsTests)
}

var completerShowArgsTests = []completerTest{
	{
		Name:     "ShowArgs",
		Line:     "",
		OrigLine: "show ",
		Index:    5,
		Expect: append(readline.CandidateList{
			{Name: []rune("CURSORS")},
			{Name: []rune("ENV")},
			{Name: []rune("FIELDS"), AppendSpace: true},
			{Name: []rune("FLAGS")},
			{Name: []rune("FUNCTIONS")},
			{Name: []rune("RUNINFO")},
			{Name: []rune("STATEMENTS")},
			{Name: []rune("TABLES")},
			{Name: []rune("VIEWS")},
		}, completer.candidateList(completer.flagList, false)...),
	},
	{
		Name:     "ShowArgs Object After SHOW",
		Line:     "cu",
		OrigLine: "show cu",
		Index:    7,
		Expect: append(readline.CandidateList{
			{Name: []rune("CURSORS")},
			{Name: []rune("ENV")},
			{Name: []rune("FIELDS"), AppendSpace: true},
			{Name: []rune("FLAGS")},
			{Name: []rune("FUNCTIONS")},
			{Name: []rune("RUNINFO")},
			{Name: []rune("STATEMENTS")},
			{Name: []rune("TABLES")},
			{Name: []rune("VIEWS")},
		}, completer.candidateList(completer.flagList, false)...),
	},
	{
		Name:     "ShowArgs After SHOW with Line",
		Line:     "cu",
		OrigLine: "show cu",
		Index:    7,
		Expect: append(readline.CandidateList{
			{Name: []rune("CURSORS")},
			{Name: []rune("ENV")},
			{Name: []rune("FIELDS"), AppendSpace: true},
			{Name: []rune("FLAGS")},
			{Name: []rune("FUNCTIONS")},
			{Name: []rune("RUNINFO")},
			{Name: []rune("STATEMENTS")},
			{Name: []rune("TABLES")},
			{Name: []rune("VIEWS")},
		}, completer.candidateList(completer.flagList, false)...),
	},
	{
		Name:     "ShowArgs After FIELD",
		Line:     "",
		OrigLine: "show fields ",
		Index:    12,
		Expect: readline.CandidateList{
			{Name: []rune("FROM"), AppendSpace: true},
		},
	},
	{
		Name:     "ShowArgs After FROM",
		Line:     "",
		OrigLine: "show fields from ",
		Index:    17,
		Expect: readline.CandidateList{
			{Name: []rune("CSV()")},
			{Name: []rune("FIXED()")},
			{Name: []rune("JSON()")},
			{Name: []rune("JSONL()")},
			{Name: []rune("LTSV()")},
			{Name: []rune(filepath.Join(CompletionTestDir, "sub", "table2.csv")), FormatAsIdentifier: true},
			{Name: []rune("newtable.csv"), FormatAsIdentifier: true},
			{Name: []rune("tempview"), FormatAsIdentifier: true},
			{Name: []rune("."), FormatAsIdentifier: true},
			{Name: []rune(".."), FormatAsIdentifier: true},
			{Name: []rune("sub/"), FormatAsIdentifier: true},
			{Name: []rune("table1.csv"), FormatAsIdentifier: true},
		},
	},
	{
		Name:     "ShowArgs Invalid Syntax",
		Line:     "",
		OrigLine: "show fields 1 ",
		Index:    14,
		Expect:   readline.CandidateList{},
	},
}

func TestCompleter_ShowArgs(t *testing.T) {
	testCompleter(t, completer.ShowArgs, completerShowArgsTests)
}

var completerSearchAllTablesTests = []completerTest{
	{
		Name:     "SearchAllTables",
		Line:     "",
		OrigLine: "select * from ",
		Index:    14,
		Expect: readline.CandidateList{
			{Name: []rune(filepath.Join(CompletionTestDir, "sub", "table2.csv")), FormatAsIdentifier: true},
			{Name: []rune("newtable.csv"), FormatAsIdentifier: true},
			{Name: []rune("tempview"), FormatAsIdentifier: true},
			{Name: []rune("."), FormatAsIdentifier: true},
			{Name: []rune(".."), FormatAsIdentifier: true},
			{Name: []rune("sub/"), FormatAsIdentifier: true},
			{Name: []rune("table1.csv"), FormatAsIdentifier: true},
		},
	},
	{
		Name:     "SearchAllTables Exclude Duplicate",
		Line:     "sub/",
		OrigLine: "select * from sub/",
		Index:    18,
		Expect: readline.CandidateList{
			{Name: []rune(filepath.Join(CompletionTestDir, "sub", "table2.csv")), FormatAsIdentifier: true},
			{Name: []rune("newtable.csv"), FormatAsIdentifier: true},
			{Name: []rune("tempview"), FormatAsIdentifier: true},
		},
	},
}

func TestCompleter_SearchAllTables(t *testing.T) {
	testCompleter(t, completer.SearchAllTables, completerSearchAllTablesTests)
}

var completerSearchExcutableFilesTests = []completerTest{
	{
		Name:     "SearchExecutableFiles",
		Line:     "",
		OrigLine: "source ",
		Index:    7,
		Expect: readline.CandidateList{
			{Name: []rune("."), FormatAsIdentifier: true},
			{Name: []rune(".."), FormatAsIdentifier: true},
			{Name: []rune("source.sql"), FormatAsIdentifier: true},
			{Name: []rune("sub/"), FormatAsIdentifier: true},
		},
	},
	{
		Name:     "SearchExecutableFiles with Variables",
		Line:     "@",
		OrigLine: "chdir @",
		Index:    8,
		Expect: readline.CandidateList{
			{Name: []rune("@var1")},
			{Name: []rune("@var2")},
		},
	},
}

func TestCompleter_SearchExecutableFiles(t *testing.T) {
	testCompleter(t, completer.SearchExecutableFiles, completerSearchExcutableFilesTests)
}

var completerSearchDirsTests = []completerTest{
	{
		Name:     "SearchDirs",
		Line:     "",
		OrigLine: "chdir ",
		Index:    6,
		Expect: readline.CandidateList{
			{Name: []rune("."), FormatAsIdentifier: true},
			{Name: []rune(".."), FormatAsIdentifier: true},
			{Name: []rune("sub/"), FormatAsIdentifier: true},
		},
	},
	{
		Name:     "SearchDirs with Variables",
		Line:     "@",
		OrigLine: "chdir @",
		Index:    7,
		Expect: readline.CandidateList{
			{Name: []rune("@var1")},
			{Name: []rune("@var2")},
		},
	},
}

func TestCompleter_SearchDirs(t *testing.T) {
	testCompleter(t, completer.SearchDirs, completerSearchDirsTests)
}

var completerSearchValuesWithSpaceTests = []completerTest{
	{
		Name:     "SearchValuesWithSpace Environment Variables",
		Line:     "@%",
		OrigLine: "select @%",
		Index:    9,
		Expect: readline.CandidateList{
			{Name: []rune("@%ENV1"), AppendSpace: true},
			{Name: []rune("@%ENV2"), AppendSpace: true},
		},
	},
}

func TestCompleter_SearchValuesWithSpace(t *testing.T) {
	testCompleter(t, completer.SearchValuesWithSpace, completerSearchValuesWithSpaceTests)
}

var completerSearchValuesTests = []completerTest{
	{
		Name:     "SearchValues Quotatin Not Enclosed",
		Line:     "'str",
		OrigLine: "select 'str",
		Index:    11,
		Expect: readline.CandidateList{
			{Name: []rune("'str'")},
		},
	},
	{
		Name:     "SearchValues Cursor Status",
		Line:     "",
		OrigLine: "select cursor ",
		Index:    14,
		Expect: readline.CandidateList{
			{Name: []rune("cur1"), AppendSpace: true},
			{Name: []rune("cur2"), AppendSpace: true},
		},
	},
	{
		Name:     "SearchValues Case Expression",
		Line:     "w",
		OrigLine: "select case w",
		Index:    13,
		Expect:   readline.CandidateList{{Name: []rune("WHEN"), AppendSpace: true}},
	},
	{
		Name:     "SearchValues Environment Variables",
		Line:     "@%",
		OrigLine: "select @%",
		Index:    9,
		Expect: readline.CandidateList{
			{Name: []rune("@%ENV1")},
			{Name: []rune("@%ENV2")},
		},
	},
	{
		Name:     "SearchValues Enclosed Environment Variables",
		Line:     "@%`",
		OrigLine: "select @%`",
		Index:    10,
		Expect: readline.CandidateList{
			{Name: []rune("@%`ENV1`")},
			{Name: []rune("@%`ENV2`")},
		},
	},
	{
		Name:     "SearchValues Keyword",
		Line:     "c",
		OrigLine: "select c",
		Index:    8,
		Expect: readline.CandidateList{
			{Name: []rune("COUNT()"), AppendSpace: false},
			{Name: []rune("CASE"), AppendSpace: true},
			{Name: []rune("CURSOR"), AppendSpace: true},
		},
	},
	{
		Name:     "SearchValues Subquery",
		Line:     "s",
		OrigLine: "@a := (s",
		Index:    8,
		Expect: readline.CandidateList{
			{Name: []rune("SELECT"), AppendSpace: true},
		},
	},
}

func TestCompleter_SearchValues(t *testing.T) {
	testCompleter(t, completer.SearchValues, completerSearchValuesTests)
}

var completerCursorStatusTests = []completerTest{
	{
		Name:     "CursorStatus Empty Line",
		Line:     "",
		OrigLine: "",
		Index:    0,
		Expect:   readline.CandidateList{},
	},
	{
		Name:     "CursorStatus After CURSOR",
		Line:     "",
		OrigLine: "select cursor ",
		Index:    14,
		Expect: readline.CandidateList{
			{Name: []rune("cur1"), AppendSpace: true},
			{Name: []rune("cur2"), AppendSpace: true},
		},
	},
	{
		Name:     "CursorStatus After cursor name",
		Line:     "",
		OrigLine: "select cursor cur1 ",
		Index:    19,
		Expect: readline.CandidateList{
			{Name: []rune("IS"), AppendSpace: true},
			{Name: []rune("COUNT")},
		},
	},
	{
		Name:     "CursorStatus After IS",
		Line:     "",
		OrigLine: "select cursor cur1 is ",
		Index:    22,
		Expect: readline.CandidateList{
			{Name: []rune("NOT"), AppendSpace: true},
			{Name: []rune("IN RANGE")},
			{Name: []rune("OPEN")},
		},
	},
	{
		Name:     "CursorStatus After NOT",
		Line:     "",
		OrigLine: "select cursor cur1 is not ",
		Index:    26,
		Expect: readline.CandidateList{
			{Name: []rune("IN RANGE")},
			{Name: []rune("OPEN")},
		},
	},
	{
		Name:     "CursorStatus After IN",
		Line:     "",
		OrigLine: "select cursor cur1 is not in ",
		Index:    29,
		Expect: readline.CandidateList{
			{Name: []rune("RANGE")},
		},
	},
	{
		Name:     "CursorStatus Terminated",
		Line:     "",
		OrigLine: "select cursor cur1 is not in range ",
		Index:    35,
		Expect:   readline.CandidateList{},
	},
}

func TestCompleter_CursorStatus(t *testing.T) {
	testCompleter(t, completer.CursorStatus, completerCursorStatusTests)
}

var completerCaseExpressionTests = []completerTest{
	{
		Name:     "CaseExpression Empty Line",
		Line:     "",
		OrigLine: "",
		Index:    0,
		Expect:   readline.CandidateList{},
	},
	{
		Name:     "CaseExpression Enclosed",
		Line:     "",
		OrigLine: "select case c1 when true then 1 else 2 end, c2",
		Index:    46,
		Expect:   readline.CandidateList{},
	},
	{
		Name:     "CaseExpression Quotation Not Enclosed",
		Line:     "'abc def",
		OrigLine: "select case c1 when true then 1 else 'abc def",
		Index:    45,
		Expect:   readline.CandidateList{{Name: []rune("'abc def'")}},
	},
	{
		Name:     "Case Statement",
		Line:     "'abc def",
		OrigLine: "case c1 when true then 1 ",
		Index:    25,
		Expect:   readline.CandidateList{},
	},
	{
		Name:     "CaseExpression After CASE",
		Line:     "",
		OrigLine: "select case ",
		Index:    12,
		Expect:   readline.CandidateList{{Name: []rune("WHEN"), AppendSpace: true}},
	},
	{
		Name:     "CaseExpression After WHEN",
		Line:     "",
		OrigLine: "select case when ",
		Index:    17,
		Expect:   readline.CandidateList{},
	},
	{
		Name:     "CaseExpression After WHEN and a Value",
		Line:     "",
		OrigLine: "select case when 1 ",
		Index:    19,
		Expect:   readline.CandidateList{{Name: []rune("THEN"), AppendSpace: true}},
	},
	{
		Name:     "CaseExpression After THEN",
		Line:     "",
		OrigLine: "select case when 1 then ",
		Index:    24,
		Expect:   readline.CandidateList{},
	},
	{
		Name:     "CaseExpression After THEN and a Value",
		Line:     "",
		OrigLine: "select case when 1 then 'a' ",
		Index:    28,
		Expect: readline.CandidateList{
			{Name: []rune("ELSE"), AppendSpace: true},
			{Name: []rune("WHEN"), AppendSpace: true},
			{Name: []rune("END")},
		},
	},
	{
		Name:     "CaseExpression After ELSE",
		Line:     "",
		OrigLine: "select case when 1 then 'a' else ",
		Index:    33,
		Expect:   readline.CandidateList{},
	},
	{
		Name:     "CaseExpression After ELSE and a Value",
		Line:     "",
		OrigLine: "select case when 1 then 'a' else 'b' ",
		Index:    37,
		Expect: readline.CandidateList{
			{Name: []rune("END")},
		},
	},
	{
		Name:     "Case Statement in Multiple Statements",
		Line:     "'abc def",
		OrigLine: "select 1;case c1 when true then 1 ",
		Index:    34,
		Expect:   readline.CandidateList{},
	},
	{
		Name:     "CaseExpression After CASEin Multiple Statements",
		Line:     "",
		OrigLine: "select 1; select case ",
		Index:    22,
		Expect:   readline.CandidateList{{Name: []rune("WHEN"), AppendSpace: true}},
	},
}

func TestCompleter_CaseExpression(t *testing.T) {
	testCompleter(t, completer.CaseExpression, completerCaseExpressionTests)
}

var completerEncloseQuotationTests = []completerTest{
	{
		Name:     "Enclose Quotation",
		Line:     "'abc",
		OrigLine: "select 'abc",
		Index:    11,
		Expect:   readline.CandidateList{{Name: []rune("'abc'")}},
	},
	{
		Name:     "Enclose Quotation Not Complete",
		Line:     "",
		OrigLine: "select 'abc'",
		Index:    12,
		Expect:   readline.CandidateList(nil),
	},
	{
		Name:     "Enclose Quotation Not Enclosed",
		Line:     "'abc def",
		OrigLine: "select case c1 when true then 1 else 'abc def",
		Index:    45,
		Expect:   readline.CandidateList{{Name: []rune("'abc def'")}},
	},
}

func TestCompleter_EncloseQuotation(t *testing.T) {
	testCompleter(t, completer.EncloseQuotation, completerEncloseQuotationTests)
}

var completerListFilesTests = []struct {
	Name       string
	Line       string
	IncludeExt []string
	Repository string
	Expect     []string
}{
	{
		Name:       "Directory",
		Line:       "",
		IncludeExt: nil,
		Repository: "",
		Expect: []string{
			".",
			"..",
			"sub/",
		},
	},
	{
		Name:       "Executable Files",
		Line:       "",
		IncludeExt: []string{".sql"},
		Repository: "",
		Expect: []string{
			".",
			"..",
			"source.sql",
			"sub/",
		},
	},
	{
		Name:       "CSV Files in Sub Directory",
		Line:       "sub/",
		IncludeExt: []string{".csv"},
		Repository: "",
		Expect: []string{
			"sub/table2.csv",
		},
	},
	{
		Name:       "CSV Files with Identifier Quotation Mark",
		Line:       "`sub/",
		IncludeExt: []string{".csv"},
		Repository: "",
		Expect: []string{
			"sub/table2.csv",
		},
	},
	{
		Name:       "CSV Files with Not Exist Directory",
		Line:       filepath.Join(CompletionTestDir, "notexist"),
		IncludeExt: []string{".csv"},
		Repository: "",
		Expect:     []string{},
	},
	{
		Name:       "CSV Files with Completion",
		Line:       "su",
		IncludeExt: []string{".csv"},
		Repository: "",
		Expect: []string{
			"sub/",
		},
	},
	{
		Name:       "CSV Files in Repository",
		Line:       "",
		IncludeExt: []string{".csv"},
		Repository: CompletionTestDir,
		Expect: []string{
			".",
			"..",
			"sub/",
			"table1.csv",
		},
	},
	{
		Name:       "CSV Files with Current Directory Mark",
		Line:       ".",
		IncludeExt: []string{".csv"},
		Repository: CompletionTestDir,
		Expect: []string{
			"./sub/",
			"./table1.csv",
		},
	},
	{
		Name:       "CSV Files with Absolute Path",
		Line:       CompletionTestDir,
		IncludeExt: []string{".csv"},
		Repository: "",
		Expect: []string{
			filepath.Join(CompletionTestDir, "sub") + "/",
			filepath.Join(CompletionTestDir, "table1.csv"),
		},
	},
}

func TestCompleter_ListFiles(t *testing.T) {
	wd, _ := os.Getwd()

	_ = os.Chdir(CompletionTestDir)
	completer := NewCompleter(NewReferenceScope(TestTx))
	for _, v := range completerListFilesTests {
		result := completer.ListFiles(v.Line, v.IncludeExt, v.Repository)
		if !reflect.DeepEqual(result, v.Expect) {
			t.Errorf("result = %v, want %v for %s", result, v.Expect, v.Name)
		}
	}

	_ = os.Chdir(wd)
}

func TestCompleter_AllColumnList(t *testing.T) {
	defer func() {
		_ = TestTx.cachedViews.Clean(TestTx.FileContainer)
	}()

	scope := NewReferenceScope(TestTx)
	scope.SetTemporaryTable(&View{
		FileInfo: &FileInfo{Path: "view1", ViewType: ViewTypeTemporaryTable},
		Header:   NewHeader("view1", []string{"v1col1", "v1col2", "v1col3"}),
	})
	scope.SetTemporaryTable(&View{FileInfo: &FileInfo{Path: "view2", ViewType: ViewTypeTemporaryTable}})
	TestTx.cachedViews.Set(
		&View{
			FileInfo: &FileInfo{
				Path: filepath.Join(CompletionTestDir, "newtable.csv"),
			},
			Header: NewHeader("newtable", []string{"ncol1", "col2", "ncol3"}),
		},
	)
	TestTx.cachedViews.Set(
		&View{
			FileInfo: &FileInfo{
				Path: filepath.Join(CompletionTestDir, "table1.csv"),
			},
			Header: NewHeader("newtable", []string{"tcol1", "col2", "tcol3"}),
		},
	)

	expect := []string{"col2", "ncol1", "ncol3", "tcol1", "tcol3", "v1col1", "v1col2", "v1col3"}

	completer := NewCompleter(scope)
	result := completer.AllColumnList()
	if !reflect.DeepEqual(result, expect) {
		t.Errorf("result = %v, want %v", result, expect)
	}
}

var completerColumnListTests = []struct {
	TableName string
	Expect    []string
}{
	{
		TableName: "view1",
		Expect:    []string{"v1col1", "v1col2", "v1col3"},
	},
	{
		TableName: "newtable.csv",
		Expect:    []string{"ncol3", "ncol2", "ncol1"},
	},
	{
		TableName: "table1",
		Expect:    []string{"tcol1", "tcol2", "tcol3"},
	},
	{
		TableName: "notexist",
		Expect:    []string(nil),
	},
	{ //Use Saved List
		TableName: "view1",
		Expect:    []string{"v1col1", "v1col2", "v1col3"},
	},
}

func TestCompleter_ColumnList(t *testing.T) {
	wd, _ := os.Getwd()

	defer func() {
		_ = os.Chdir(wd)
		_ = TestTx.cachedViews.Clean(TestTx.FileContainer)
		initFlag(TestTx.Flags)
	}()

	scope := NewReferenceScope(TestTx)
	scope.SetTemporaryTable(&View{
		FileInfo: &FileInfo{Path: "view1", ViewType: ViewTypeTemporaryTable},
		Header:   NewHeader("view1", []string{"v1col1", "v1col2", "v1col3"}),
	})
	scope.SetTemporaryTable(&View{FileInfo: &FileInfo{Path: "view2", ViewType: ViewTypeTemporaryTable}})
	TestTx.cachedViews.Set(
		&View{
			FileInfo: &FileInfo{
				Path: filepath.Join(CompletionTestDir, "newtable.csv"),
			},
			Header: NewHeader("newtable", []string{"ncol3", "ncol2", "ncol1"}),
		},
	)
	TestTx.cachedViews.Set(
		&View{
			FileInfo: &FileInfo{
				Path: filepath.Join(CompletionTestDir, "table1.csv"),
			},
			Header: NewHeader("newtable", []string{"tcol1", "tcol2", "tcol3"}),
		},
	)

	TestTx.Flags.ExportOptions.Format = cmd.CSV

	_ = os.Chdir(CompletionTestDir)
	completer := NewCompleter(scope)
	for _, v := range completerColumnListTests {
		result := completer.ColumnList(v.TableName, CompletionTestDir)
		if !reflect.DeepEqual(result, v.Expect) {
			t.Errorf("result = %v, want %v for %s", result, v.Expect, v.TableName)
		}
	}
}

var completerUpdateTokensTests = []struct {
	SearchWord string
	Statements string
	Expect     []parser.Token
	LastIdx    int
}{
	{
		SearchWord: "",
		Statements: "",
		Expect:     []parser.Token{},
		LastIdx:    -1,
	},
	{
		SearchWord: "",
		Statements: "select 1",
		Expect: []parser.Token{
			{Token: parser.SELECT, Literal: "select", Line: 1, Char: 1},
			{Token: parser.INTEGER, Literal: "1", Line: 1, Char: 8},
		},
		LastIdx: 1,
	},
	{
		SearchWord: "",
		Statements: "select 1; select 2",
		Expect: []parser.Token{
			{Token: parser.SELECT, Literal: "select", Line: 1, Char: 11},
			{Token: parser.INTEGER, Literal: "2", Line: 1, Char: 18},
		},
		LastIdx: 1,
	},
	{
		SearchWord: "",
		Statements: "select 1 from (select 2",
		Expect: []parser.Token{
			{Token: parser.SELECT, Literal: "select", Line: 1, Char: 16},
			{Token: parser.INTEGER, Literal: "2", Line: 1, Char: 23},
		},
		LastIdx: 1,
	},
	{
		SearchWord: "",
		Statements: "select 1 from (select 2, fn()) join",
		Expect: []parser.Token{
			{Token: parser.SELECT, Literal: "select", Line: 1, Char: 1},
			{Token: parser.INTEGER, Literal: "1", Line: 1, Char: 8},
			{Token: parser.FROM, Literal: "from", Line: 1, Char: 10},
			{Token: parser.IDENTIFIER, Literal: dummySubquery},
			{Token: parser.JOIN, Literal: "join", Line: 1, Char: 32},
		},
		LastIdx: 4,
	},
	{
		SearchWord: "",
		Statements: "select 1 from csv((','), `table.csv`) join",
		Expect: []parser.Token{
			{Token: parser.SELECT, Literal: "select", Line: 1, Char: 1},
			{Token: parser.INTEGER, Literal: "1", Line: 1, Char: 8},
			{Token: parser.FROM, Literal: "from", Line: 1, Char: 10},
			{Token: parser.IDENTIFIER, Literal: "table.csv"},
			{Token: parser.JOIN, Literal: "join", Line: 1, Char: 39},
		},
		LastIdx: 4,
	},
	{
		SearchWord: "",
		Statements: "select 1 from csv(',') join",
		Expect: []parser.Token{
			{Token: parser.SELECT, Literal: "select", Line: 1, Char: 1},
			{Token: parser.INTEGER, Literal: "1", Line: 1, Char: 8},
			{Token: parser.FROM, Literal: "from", Line: 1, Char: 10},
			{Token: parser.IDENTIFIER, Literal: dummyTableObject},
			{Token: parser.JOIN, Literal: "join", Line: 1, Char: 24},
		},
		LastIdx: 4,
	},
	{
		SearchWord: "",
		Statements: "select 1 from csv((','), ",
		Expect: []parser.Token{
			{Token: parser.CSV, Literal: "csv", Line: 1, Char: 15},
			{Token: '(', Literal: "(", Line: 1, Char: 18},
			{Token: '(', Literal: "(", Line: 1, Char: 19},
			{Token: parser.STRING, Literal: ",", Line: 1, Char: 20},
			{Token: ')', Literal: ")", Line: 1, Char: 23},
			{Token: ',', Literal: ",", Line: 1, Char: 24},
		},
		LastIdx: 5,
	},
	{
		SearchWord: "1",
		Statements: "select 1",
		Expect: []parser.Token{
			{Token: parser.SELECT, Literal: "select", Line: 1, Char: 1},
			{Token: parser.IDENTIFIER, Literal: "1", Line: 1, Char: 8},
		},
		LastIdx: 0,
	},
	{
		SearchWord: "",
		Statements: "select trim(",
		Expect: []parser.Token{
			{Token: parser.IDENTIFIER, Literal: "trim", Line: 1, Char: 8},
			{Token: '(', Literal: "(", Line: 1, Char: 12},
		},
		LastIdx: 1,
	},
	{
		SearchWord: "",
		Statements: "select trim('a')",
		Expect: []parser.Token{
			{Token: parser.SELECT, Literal: "select", Line: 1, Char: 1},
			{Token: parser.FUNCTION, Literal: "trim"},
		},
		LastIdx: 1,
	},
	{
		SearchWord: "",
		Statements: "select substring(",
		Expect: []parser.Token{
			{Token: parser.SUBSTRING, Literal: "substring", Line: 1, Char: 8},
			{Token: '(', Literal: "(", Line: 1, Char: 17},
		},
		LastIdx: 1,
	},
	{
		SearchWord: "",
		Statements: "select substring('abc' from 2 for 1)",
		Expect: []parser.Token{
			{Token: parser.SELECT, Literal: "select", Line: 1, Char: 1},
			{Token: parser.FUNCTION, Literal: "substring"},
		},
		LastIdx: 1,
	},
	{
		SearchWord: "",
		Statements: "select userfunc(",
		Expect: []parser.Token{
			{Token: parser.IDENTIFIER, Literal: "userfunc", Line: 1, Char: 8},
			{Token: '(', Literal: "(", Line: 1, Char: 16},
		},
		LastIdx: 1,
	},
	{
		SearchWord: "",
		Statements: "select count(",
		Expect: []parser.Token{
			{Token: parser.COUNT, Literal: "count", Line: 1, Char: 8},
			{Token: '(', Literal: "(", Line: 1, Char: 13},
		},
		LastIdx: 1,
	},
	{
		SearchWord: "",
		Statements: "select rank() over ",
		Expect: []parser.Token{
			{Token: parser.SELECT, Literal: "select", Line: 1, Char: 1},
			{Token: parser.FUNCTION, Literal: "rank"},
			{Token: parser.OVER, Literal: "over", Line: 1, Char: 15},
		},
		LastIdx: 2,
	},
	{
		SearchWord: "",
		Statements: "select rank() over (",
		Expect: []parser.Token{
			{Token: parser.ANALYTIC_FUNCTION, Literal: "rank", Line: 1, Char: 8},
			{Token: '(', Literal: "(", Line: 1, Char: 12},
			{Token: ')', Literal: ")", Line: 1, Char: 13},
			{Token: parser.OVER, Literal: "over", Line: 1, Char: 15},
			{Token: '(', Literal: "(", Line: 1, Char: 20},
		},
		LastIdx: 4,
	},
	{
		SearchWord: "",
		Statements: "select count((1)) as",
		Expect: []parser.Token{
			{Token: parser.SELECT, Literal: "select", Line: 1, Char: 1},
			{Token: parser.FUNCTION, Literal: "count"},
			{Token: parser.AS, Literal: "as", Line: 1, Char: 19},
		},
		LastIdx: 2,
	},
	{
		SearchWord: "",
		Statements: "select rank() over () as",
		Expect: []parser.Token{
			{Token: parser.SELECT, Literal: "select", Line: 1, Char: 1},
			{Token: parser.FUNCTION, Literal: "rank"},
			{Token: parser.AS, Literal: "as", Line: 1, Char: 23},
		},
		LastIdx: 2,
	},
	{
		SearchWord: "",
		Statements: "select first_value() ignore nulls over () as",
		Expect: []parser.Token{
			{Token: parser.SELECT, Literal: "select", Line: 1, Char: 1},
			{Token: parser.FUNCTION, Literal: "first_value"},
			{Token: parser.AS, Literal: "as", Line: 1, Char: 43},
		},
		LastIdx: 2,
	},
}

func TestCompleter_UpdateTokens(t *testing.T) {
	c := NewCompleter(NewReferenceScope(TestTx))
	c.userFuncs = []string{"userfunc"}
	c.userAggFuncs = []string{"aggfunc"}
	c.userFuncList = []string{"userfunc", "aggfunc"}
	for _, v := range completerUpdateTokensTests {
		c.UpdateTokens(v.SearchWord, v.Statements)
		if !reflect.DeepEqual(c.tokens, v.Expect) {
			t.Errorf("tokens = %v, want %v for %q, %q", c.tokens, v.Expect, v.SearchWord, v.Statements)
		}
		if c.lastIdx != v.LastIdx {
			t.Errorf("last index = %d, want %d for %q, %q", c.lastIdx, v.LastIdx, v.SearchWord, v.Statements)
		}
	}
}
