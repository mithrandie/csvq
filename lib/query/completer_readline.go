//go:build darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris || windows
// +build darwin dragonfly freebsd linux netbsd openbsd solaris windows

package query

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/parser"

	"github.com/mithrandie/go-text"
	"github.com/mithrandie/readline-csvq"
	"github.com/mithrandie/ternary"
)

const (
	dummySubquery    = "____subquery____"
	dummyTableObject = "____table_object____"
	dummyTable       = "____table____"
)

var statementPrefix = []string{
	"WITH",
	"SELECT",
	"INSERT",
	"UPDATE",
	"REPLACE",
	"DELETE",
	"CREATE",
	"ALTER",
	"DECLARE",
	"PREPARE",
	"VAR",
	"SET",
	"UNSET",
	"ADD",
	"REMOVE",
	"ECHO",
	"PRINT",
	"PRINTF",
	"CHDIR",
	"EXECUTE",
	"SHOW",
	"SOURCE",
	"SYNTAX",
	"RELOAD",
}

var singleCommandStatement = []string{
	"COMMIT",
	"ROLLBACK",
	"EXIT",
	"PWD",
}

var delimiterCandidates = []string{
	"','",
	"'\\t'",
}

var delimiterPositionsCandidates = []string{
	"'SPACES'",
	"'S[]'",
	"'[]'",
}

var joinCandidates = []string{
	"JOIN",
	"CROSS",
	"INNER",
	"FULL",
	"LEFT",
	"RIGHT",
	"NATURAL",
}

var tableObjectCandidates = []string{
	"CSV()",
	"FIXED()",
	"JSON()",
	"JSONL()",
	"LTSV()",
}

var exportEncodingsCandidates = []string{
	"SJIS",
	"UTF16",
	"UTF16BE",
	"UTF16BEM",
	"UTF16LE",
	"UTF16LEM",
	"UTF8",
	"UTF8M",
}

type ReadlineListener struct {
	scanner parser.Scanner
}

func skipInputtingEnclosure(line []rune, pos int) []rune {
	tail := line[pos:]
	line = append(line[:pos-1], tail...)
	return line
}

func completeEnclosure(line []rune, pos int, rightEnclosure rune) []rune {
	tail := append([]rune{rightEnclosure}, line[pos:]...)
	line = append(line[:pos], tail...)
	return line
}

func (l ReadlineListener) OnChange(line []rune, pos int, key rune) ([]rune, int, bool) {
	switch {
	case readline.IsQuotationMark(key):
		if !readline.LiteralIsEnclosed(key, line) {
			if pos < len(line) && key == line[pos] {
				return skipInputtingEnclosure(line, pos), pos, true
			} else {
				return completeEnclosure(line, pos, key), pos, true
			}
		}
	case readline.IsBracket(key):
		if !readline.BracketIsEnclosed(key, line) {
			return completeEnclosure(line, pos, readline.RightBracket[key]), pos, true
		}
	case readline.IsRightBracket(key):
		if pos < len(line) && readline.IsRightBracket(line[pos]) && readline.BracketIsEnclosedByRightBracket(key, line) {
			return skipInputtingEnclosure(line, pos), pos, true
		}
	}

	return line, pos, false
}

type Completer struct {
	completer *readline.PrefixCompleter
	scope     *ReferenceScope

	flagList      []string
	runinfoList   []string
	funcs         []string
	aggFuncs      []string
	analyticFuncs []string

	statementList    []string
	userFuncs        []string
	userAggFuncs     []string
	userFuncList     []string
	viewList         []string
	cursorList       []string
	funcList         []string
	aggFuncList      []string
	analyticFuncList []string
	varList          []string
	envList          []string
	enclosedEnvList  []string
	allColumns       []string
	tableColumns     map[string][]string

	tokens            []parser.Token
	lastIdx           int
	selectIntoEnabled bool
}

func NewCompleter(scope *ReferenceScope) *Completer {
	completer := &Completer{
		completer:    readline.NewPrefixCompleter(),
		scope:        scope,
		tableColumns: make(map[string][]string),
	}

	completer.flagList = make([]string, 0, len(cmd.FlagList))
	for _, v := range cmd.FlagList {
		completer.flagList = append(completer.flagList, cmd.FlagSymbol(v))
	}
	completer.runinfoList = make([]string, 0, len(RuntimeInformatinList))
	for _, v := range RuntimeInformatinList {
		completer.runinfoList = append(completer.runinfoList, cmd.RuntimeInformationSymbol(v))
	}

	sort.Strings(completer.flagList)
	sort.Strings(completer.runinfoList)

	completer.funcs = make([]string, 0, len(Functions)+3)
	for k := range Functions {
		completer.funcs = append(completer.funcs, k)
	}
	completer.funcs = append(completer.funcs, "CALL")
	completer.funcs = append(completer.funcs, "NOW")
	completer.funcs = append(completer.funcs, "JSON_OBJECT")

	completer.aggFuncs = make([]string, 0, len(AggregateFunctions)+2)
	completer.analyticFuncs = make([]string, 0, len(AnalyticFunctions)+len(AggregateFunctions))
	for k := range AggregateFunctions {
		completer.aggFuncs = append(completer.aggFuncs, k)
		completer.analyticFuncs = append(completer.analyticFuncs, k)
	}
	completer.aggFuncs = append(completer.aggFuncs, "LISTAGG")
	completer.aggFuncs = append(completer.aggFuncs, "JSON_AGG")
	for k := range AnalyticFunctions {
		completer.analyticFuncs = append(completer.analyticFuncs, k)
	}

	completer.tokens = make([]parser.Token, 0, 20)

	return completer
}

func (c *Completer) Do(line []rune, pos int, index int) (readline.CandidateList, int) {
	return c.completer.Do(line, pos, index)
}

func (c *Completer) Update() {
	c.updateStatements()
	c.updateViews()
	c.updateCursors()
	c.updateFunctions()
	c.updateVariables()
	c.updateEnvironmentVariables()
	c.updateAllColumns()

	completer := readline.NewPrefixCompleter()
	statements := readline.PcItemDynamic(c.Statements)
	statements.AppendOnly = true
	completer.SetChildren([]readline.PrefixCompleterInterface{statements})
	c.completer = completer
}

func (c *Completer) updateStatements() {
	c.statementList = make([]string, 0, 10)
	c.scope.Tx.PreparedStatements.Range(func(key, value interface{}) bool {
		c.statementList = append(c.statementList, value.(*PreparedStatement).Name)
		return true
	})
	sort.Strings(c.statementList)
}

func (c *Completer) updateViews() {
	views := c.scope.AllTemporaryTables()
	viewKeys := views.SortedKeys()
	c.viewList = make([]string, 0, len(viewKeys))
	for _, key := range viewKeys {
		if view, ok := views.Load(key); ok {
			c.viewList = append(c.viewList, view.FileInfo.Path)
		}
	}
}

func (c *Completer) updateCursors() {
	cursors := c.scope.AllCursors()
	cursorKeys := cursors.SortedKeys()
	c.cursorList = make([]string, 0, len(cursorKeys))
	for _, key := range cursorKeys {
		if cur, ok := cursors.Load(key); ok {
			c.cursorList = append(c.cursorList, cur.name)
		}
	}
}

func (c *Completer) updateFunctions() {
	userfuncs, userAggFuncs := c.scope.AllFunctions()
	c.userFuncs = make([]string, 0, userfuncs.Len())
	c.userAggFuncs = make([]string, 0, userAggFuncs.Len())

	funcKeys := make([]string, 0, len(c.funcs)+userfuncs.Len())
	for _, v := range c.funcs {
		funcKeys = append(funcKeys, v+"()")
	}
	userfuncs.Range(func(key, value interface{}) bool {
		f := value.(*UserDefinedFunction)
		funcKeys = append(funcKeys, f.Name.String()+"()")
		c.userFuncs = append(c.userFuncs, f.Name.String())
		return true
	})
	sort.Strings(funcKeys)
	c.funcList = funcKeys

	aggFuncKeys := make([]string, 0, len(c.aggFuncs)+userAggFuncs.Len())
	for _, v := range c.aggFuncs {
		aggFuncKeys = append(aggFuncKeys, v+"()")
	}
	userAggFuncs.Range(func(key, value interface{}) bool {
		f := value.(*UserDefinedFunction)
		aggFuncKeys = append(aggFuncKeys, f.Name.String()+"()")
		c.userAggFuncs = append(c.userAggFuncs, f.Name.String())
		return true
	})
	sort.Strings(aggFuncKeys)
	c.aggFuncList = aggFuncKeys

	c.userFuncList = append(c.userFuncs, c.userAggFuncs...)
	sort.Strings(c.userFuncList)

	analyticFuncKeys := make([]string, 0, len(c.analyticFuncs)+userAggFuncs.Len())
	for _, v := range c.analyticFuncs {
		analyticFuncKeys = append(analyticFuncKeys, v+"() OVER ()")
	}
	userAggFuncs.Range(func(key, value interface{}) bool {
		f := value.(*UserDefinedFunction)
		analyticFuncKeys = append(analyticFuncKeys, f.Name.String()+"() OVER ()")
		return true
	})
	sort.Strings(analyticFuncKeys)
	c.analyticFuncList = analyticFuncKeys
}

func (c *Completer) updateVariables() {
	vars := c.scope.AllVariables()
	varKeys := vars.SortedKeys()

	c.varList = make([]string, 0, len(varKeys))
	for _, k := range varKeys {
		c.varList = append(c.varList, cmd.VariableSymbol(k))
	}
}

func (c *Completer) updateEnvironmentVariables() {
	env := os.Environ()
	c.envList = make([]string, 0, len(env))
	c.enclosedEnvList = make([]string, 0, len(env))
	for _, e := range env {
		words := strings.Split(e, "=")
		c.envList = append(c.envList, cmd.EnvironmentVariableSymbol(words[0]))
		c.enclosedEnvList = append(c.enclosedEnvList, cmd.EnclosedEnvironmentVariableSymbol(words[0]))
	}
	sort.Strings(c.envList)
	sort.Strings(c.enclosedEnvList)
}

func (c *Completer) updateAllColumns() {
	c.allColumns = c.AllColumnList()
}

func (c *Completer) GetStatementPrefix(line string, origLine string, index int) readline.CandidateList {
	prefix := statementPrefix
	if 0 < len(c.cursorList) || 0 < len(c.userFuncList) || 0 < len(c.viewList) || 0 < len(c.varList) || 0 < len(c.statementList) {
		prefix = append(prefix, "DISPOSE")
	}
	if 0 < len(c.cursorList) {
		prefix = append(prefix,
			"OPEN",
			"CLOSE",
			"FETCH",
		)
	}

	var cands readline.CandidateList
	for _, p := range prefix {
		cands = append(cands, c.candidate(p, true))
	}
	for _, p := range singleCommandStatement {
		cands = append(cands, c.candidate(p, false))
	}
	cands = append(cands, c.SearchValues(line, origLine, index)...)

	cands.Sort()
	return cands
}

func (c *Completer) Statements(line string, origLine string, index int) readline.CandidateList {
	origRunes := []rune(origLine)
	c.UpdateTokens(line, string(origRunes[:index]))

	token := parser.EOF
	if (1 == len(c.tokens) && unicode.IsSpace(origRunes[index-1])) || 1 < len(c.tokens) {
		token = c.tokens[0].Token
	}

	switch token {
	case parser.WITH:
		return c.WithArgs(line, origLine, index)
	case parser.SELECT:
		return c.SelectArgs(line, origLine, index)
	case parser.INSERT:
		return c.InsertArgs(line, origLine, index)
	case parser.UPDATE:
		return c.UpdateArgs(line, origLine, index)
	case parser.REPLACE:
		return c.ReplaceArgs(line, origLine, index)
	case parser.DELETE:
		return c.DeleteArgs(line, origLine, index)
	case parser.CREATE:
		return c.CreateArgs(line, origLine, index)
	case parser.ALTER:
		return c.AlterArgs(line, origLine, index)
	case parser.DECLARE, parser.VAR:
		return c.DeclareArgs(line, origLine, index)
	case parser.PREPARE:
		return c.PrepareArgs(line, origLine, index)
	case parser.SET:
		return c.SetArgs(line, origLine, index)
	case parser.UNSET:
		return c.candidateList(c.environmentVariableList(line), false)
	case parser.ADD:
		return c.AddFlagArgs(line, origLine, index)
	case parser.REMOVE:
		return c.RemoveFlagArgs(line, origLine, index)
	case parser.ECHO:
		return c.SearchValues(line, origLine, index)
	case parser.PRINT:
		return c.SearchValues(line, origLine, index)
	case parser.PRINTF:
		return c.UsingArgs(line, origLine, index)
	case parser.CHDIR:
		return c.SearchDirs(line, origLine, index)
	case parser.EXECUTE:
		return c.UsingArgs(line, origLine, index)
	case parser.SHOW:
		return c.ShowArgs(line, origLine, index)
	case parser.SOURCE:
		return c.SearchExecutableFiles(line, origLine, index)
	case parser.RELOAD:
		if 0 < len(line) && len(c.tokens) == 2 || len(line) < 1 && len(c.tokens) == 1 {
			return readline.CandidateList{c.candidate("CONFIG", false)}
		} else {
			return nil
		}
	case parser.DISPOSE:
		return c.DisposeArgs(line, origLine, index)
	case parser.OPEN:
		return c.UsingArgs(line, origLine, index)
	case parser.CLOSE:
		return c.candidateList(c.cursorList, false)
	case parser.FETCH:
		return c.FetchArgs(line, origLine, index)
	case parser.COMMIT:
		return nil
	case parser.ROLLBACK:
		return nil
	case parser.EXIT:
		return nil
	case parser.PWD:
		return nil
	case parser.SYNTAX:
		return nil
	case parser.EOF:
		return c.GetStatementPrefix(line, origLine, index)
	default:
		if 0 < len(c.tokens) {
			switch {
			case c.isTableObject(c.tokens[0]):
				return c.TableObjectArgs(line, origLine, index)
			case c.isFunction(c.tokens[0]):
				return c.FunctionArgs(line, origLine, index)
			}
		}
	}

	cands := c.SearchValues(line, origLine, index)
	cands.Sort()
	return cands
}

func (c *Completer) TableObjectArgs(line string, origLine string, index int) readline.CandidateList {
	commaCnt := 0
	for _, t := range c.tokens {
		if t.Token == ',' {
			commaCnt++
		}
	}

	var cands readline.CandidateList

	switch strings.ToUpper(c.tokens[0].Literal) {
	case "LTSV":
		switch commaCnt {
		case 0:
			if c.tokens[c.lastIdx].Token == '(' {
				cands = c.SearchAllTables(line, origLine, index)
			}
		case 1:
			if c.tokens[c.lastIdx].Token == ',' {
				cands = c.candidateList(c.encodingList(), false)
			}
		case 2:
			if c.tokens[c.lastIdx].Token == ',' {
				cands = c.candidateList([]string{ternary.TRUE.String(), ternary.FALSE.String()}, false)
			}
		}
	default:
		switch commaCnt {
		case 0:
			if c.tokens[c.lastIdx].Token == '(' {
				switch strings.ToUpper(c.tokens[0].Literal) {
				case cmd.CSV.String():
					cands = c.candidateList(delimiterCandidates, false)
				case cmd.FIXED.String():
					cands = c.candidateList(delimiterPositionsCandidates, false)
				}
			}
		case 1:
			if c.tokens[c.lastIdx].Token == ',' {
				cands = c.SearchAllTables(line, origLine, index)
			}
		case 2:
			if c.tokens[c.lastIdx].Token == ',' {
				switch strings.ToUpper(c.tokens[0].Literal) {
				case cmd.CSV.String(), cmd.FIXED.String():
					cands = c.candidateList(c.encodingList(), false)
				}
			}
		case 3, 4:
			if c.tokens[c.lastIdx].Token == ',' {
				switch strings.ToUpper(c.tokens[0].Literal) {
				case cmd.CSV.String(), cmd.FIXED.String():
					cands = c.candidateList([]string{ternary.TRUE.String(), ternary.FALSE.String()}, false)
				}
			}
		}
	}

	cands = append(cands, c.SearchValues(line, origLine, index)...)
	return cands
}

func (c *Completer) FunctionArgs(line string, origLine string, index int) readline.CandidateList {
	if c.tokens[0].Token == parser.SUBSTRING {
		return c.substringArgs(line, origLine, index)
	} else {
		return c.functionArgs(line, origLine, index)
	}
}

func (c *Completer) substringArgs(line string, origLine string, index int) readline.CandidateList {
	return c.completeArgs(
		line,
		origLine,
		index,
		func(i int) (keywords []string, customList readline.CandidateList, breakLoop bool) {
			customList = append(customList, c.SearchValues(line, origLine, index)...)
			customList.Sort()

			switch c.tokens[i].Token {
			case parser.FOR:
				//Do nothing
			case parser.FROM:
				if i < c.lastIdx {
					keywords = append(keywords, "FOR")
				}
			case parser.SUBSTRING:
				if i < c.lastIdx-1 {
					keywords = append(keywords, "FROM")
				}
			default:
				return keywords, customList, false
			}

			return keywords, customList, true
		},
	)
}

func (c *Completer) functionArgs(line string, origLine string, index int) readline.CandidateList {
	return c.completeArgs(
		line,
		origLine,
		index,
		func(i int) (keywords []string, customList readline.CandidateList, breakLoop bool) {
			switch c.tokens[i].Token {
			case parser.BETWEEN:
				if i == c.lastIdx {
					customList = append(customList, c.candidateList([]string{
						"UNBOUNDED PRECEDING",
						"CURRENT ROW",
					}, true)...)
				} else if i == c.lastIdx-1 {
					switch c.tokens[c.lastIdx].Token {
					case parser.UNBOUNDED:
						customList = c.candidateList([]string{"PRECEDING"}, true)
					case parser.CURRENT:
						customList = c.candidateList([]string{"ROW"}, true)
					default:
						customList = c.candidateList([]string{
							"PRECEDING",
							"FOLLOWING",
						}, true)
					}
				} else if c.tokens[c.lastIdx].Token == parser.PRECEDING ||
					c.tokens[c.lastIdx].Token == parser.FOLLOWING ||
					c.tokens[c.lastIdx].Token == parser.ROW {

					keywords = append(keywords, "AND")
				} else if c.tokens[c.lastIdx].Token == parser.AND {
					customList = append(customList, c.candidateList([]string{
						"UNBOUNDED FOLLOWING",
						"CURRENT ROW",
					}, false)...)
				} else {
					switch c.tokens[c.lastIdx].Token {
					case parser.UNBOUNDED:
						customList = c.candidateList([]string{"FOLLOWING"}, false)
					case parser.CURRENT:
						customList = c.candidateList([]string{"ROW"}, false)
					default:
						customList = c.candidateList([]string{
							"PRECEDING",
							"FOLLOWING",
						}, false)
					}
				}
			case parser.ROWS:
				if i == c.lastIdx {
					customList = append(customList, c.candidateList([]string{
						"UNBOUNDED PRECEDING",
						"CURRENT ROW",
					}, false)...)

					customList = append(customList, c.candidate("BETWEEN", true))
				} else if i == c.lastIdx-1 {
					switch c.tokens[c.lastIdx].Token {
					case parser.CURRENT:
						customList = c.candidateList([]string{"ROW"}, false)
					default:
						customList = c.candidateList([]string{"PRECEDING"}, false)
					}
				}
			case parser.ORDER:
				if i == c.lastIdx {
					keywords = append(keywords, "BY")
				} else {
					if i < c.lastIdx-1 && c.tokens[c.lastIdx].Token != ',' {
						canUseWindowingClause := false

						switch c.tokens[c.lastIdx].Token {
						case parser.ASC, parser.DESC:
							customList = append(customList, c.candidateList([]string{
								"NULLS FIRST",
								"NULLS LAST",
							}, false)...)
							canUseWindowingClause = true
						case parser.NULLS:
							customList = append(customList, c.candidateList([]string{
								"FIRST",
								"LAST",
							}, false)...)
						case parser.FIRST, parser.LAST:
							canUseWindowingClause = true
						default:
							customList = append(customList, c.candidateList([]string{
								"ASC",
								"DESC",
								"NULLS FIRST",
								"NULLS LAST",
							}, false)...)
							customList = append(customList, c.SearchValues(line, origLine, index)...)

							canUseWindowingClause = true
						}

						if canUseWindowingClause {
							funcName := strings.ToUpper(c.tokens[0].Literal)
							if funcName == "FIRST_VALUE" ||
								funcName == "LAST_VALUE" ||
								funcName == "NTH_VALUE" ||
								(funcName != "LISTAGG" && funcName != "JSON_AGG" && InStrSliceWithCaseInsensitive(funcName, c.aggFuncs)) ||
								InStrSliceWithCaseInsensitive(funcName, c.userAggFuncs) {

								customList = append(customList, c.candidate("ROWS", true))
							}
						}
					} else {
						customList = append(customList, c.SearchValues(line, origLine, index)...)
					}
				}
			case parser.PARTITION:
				if i == c.lastIdx {
					keywords = append(keywords,
						"BY",
					)
				} else {
					if i < c.lastIdx-1 && c.tokens[c.lastIdx].Token != ',' {
						keywords = append(keywords,
							"ORDER BY",
						)
					}
					customList = append(customList, c.SearchValues(line, origLine, index)...)
				}
			case parser.OVER:
				if i == c.lastIdx-1 && c.tokens[c.lastIdx].Token == '(' {
					keywords = append(keywords,
						"PARTITION BY",
						"ORDER BY",
					)
				}
			default:
				if i == 0 {
					if 0 < len(line) && i == c.lastIdx-1 {
						if InStrSliceWithCaseInsensitive(c.tokens[i].Literal, c.aggFuncs) ||
							InStrSliceWithCaseInsensitive(c.tokens[i].Literal, c.userAggFuncs) {
							keywords = append(keywords, "DISTINCT")
						}
					}
					customList = append(customList, c.SearchValues(line, origLine, index)...)
				} else {
					return nil, nil, false
				}
			}

			if customList != nil {
				customList.Sort()
			}
			return keywords, customList, true
		},
	)
}

func (c *Completer) WithArgs(line string, origLine string, index int) readline.CandidateList {
	var searchQuery = func() int {
		blockLevel := 0
		for i := 0; i < len(c.tokens); i++ {
			switch c.tokens[i].Token {
			case ')':
				blockLevel++
			case '(':
				blockLevel--
			case parser.SELECT, parser.INSERT, parser.UPDATE, parser.REPLACE, parser.DELETE:
				if blockLevel == 0 {
					return i
				}
			}
		}
		return -1
	}

	if queryStart := searchQuery(); -1 < queryStart {
		c.tokens = c.tokens[queryStart:]
		c.SetLastIndex(line)

		switch c.tokens[0].Token {
		case parser.SELECT:
			return c.SelectArgs(line, origLine, index)
		case parser.INSERT:
			return c.InsertArgs(line, origLine, index)
		case parser.UPDATE:
			return c.UpdateArgs(line, origLine, index)
		case parser.REPLACE:
			return c.ReplaceArgs(line, origLine, index)
		case parser.DELETE:
			return c.DeleteArgs(line, origLine, index)
		}
	}

	blockLevel := 0

	return c.completeArgs(
		line,
		origLine,
		index,
		func(i int) (keywords []string, customList readline.CandidateList, breakLoop bool) {
			switch c.tokens[i].Token {
			case ')':
				blockLevel++
				if blockLevel == 1 && i == c.lastIdx {
					keywords = append(keywords, "AS")
				}
			case '(':
				blockLevel--
				if blockLevel == -1 && i == c.lastIdx && c.tokens[c.lastIdx-1].Token == parser.AS {
					keywords = append(keywords, "SELECT")
				}
			case parser.WITH:
				if blockLevel == 0 {
					switch {
					case i == c.lastIdx || c.tokens[c.lastIdx].Token == ',':
						keywords = append(keywords, "RECURSIVE")
					case c.tokens[c.lastIdx].Token == parser.IDENTIFIER:
						if c.tokens[c.lastIdx].Literal == dummySubquery && c.tokens[c.lastIdx-1].Token == parser.AS {
							keywords = append(keywords,
								"SELECT",
								"INSERT",
								"UPDATE",
								"REPLACE",
								"DELETE",
							)
						} else {
							keywords = append(keywords, "AS")
						}
					}
				}
			default:
				return nil, nil, false
			}
			return keywords, customList, true
		},
	)
}

func (c *Completer) combineTableAlias(fromIdx int) {
	combined := make([]parser.Token, 0, cap(c.tokens))
	temp := make([]parser.Token, 0, cap(c.tokens))
	blockLevel := 0
	for i := 0; i < len(c.tokens); i++ {
		if i <= fromIdx {
			combined = append(combined, c.tokens[i])
			continue
		}

		if 0 != blockLevel {
			temp = append(temp, c.tokens[i])

			switch c.tokens[i].Token {
			case '(':
				blockLevel++
			case ')':
				blockLevel--
				if 0 == blockLevel {
					combined = append(combined, parser.Token{Token: parser.IDENTIFIER, Literal: dummyTable, Quoted: true})
					temp = temp[:0]
				}
			}
			continue
		}

		if c.tokens[i].Token == '(' {
			blockLevel++
			temp = append(temp, c.tokens[i])
			continue
		}

		combined = append(combined, c.tokens[i])

		if c.tokens[i].Token == parser.IDENTIFIER {
			combined[len(combined)-1].Quoted = false
			if i+2 <= c.lastIdx && c.tokens[i+1].Token == parser.AS && c.tokens[i+2].Token == parser.IDENTIFIER {
				i = i + 2
				combined[len(combined)-1].Quoted = true
			} else if i+1 <= c.lastIdx && c.tokens[i+1].Token == parser.IDENTIFIER {
				i = i + 1
				combined[len(combined)-1].Quoted = true
			}
		}
	}
	if 0 < len(temp) {
		combined = append(combined, temp...)
	}

	c.tokens = combined
}

func (c *Completer) allTableCandidates(line string, origLine string, index int) readline.CandidateList {
	list := c.candidateList(append(tableObjectCandidates, "JSON_TABLE()"), false)
	list.Sort()
	list = append(list, c.SearchAllTables(line, origLine, index)...)
	return list
}

func (c *Completer) allTableCandidatesForUpdate(line string, origLine string, index int) readline.CandidateList {
	list := c.candidateList(tableObjectCandidates, false)
	list.Sort()
	list = append(list, c.SearchAllTables(line, origLine, index)...)
	return list
}

func (c *Completer) allTableCandidatesWithSpaceForUpdate(line string, origLine string, index int) readline.CandidateList {
	list := c.candidateList(tableObjectCandidates, true)
	list.Sort()
	list = append(list, c.SearchAllTablesWithSpace(line, origLine, index)...)
	return list
}

func (c *Completer) fromClause(i int, line string, origLine string, index int) (tables readline.CandidateList, customList readline.CandidateList, restrict bool) {
	c.combineTableAlias(i)
	c.SetLastIndex(line)

	var isInUsing = func() bool {
		for j := len(c.tokens) - 1; j > i; j-- {
			switch c.tokens[j].Token {
			case parser.IDENTIFIER, ',':
				//OK
			case '(':
				return c.tokens[j-1].Token == parser.USING
			default:
				return false
			}
		}
		return false
	}

	var isInOn = func() bool {
		blockLevel := 0
		for j := len(c.tokens) - 1; j > i; j-- {
			switch c.tokens[j].Token {
			case '(':
				blockLevel--
			case ')':
				blockLevel++
			case parser.JOIN:
				if blockLevel <= 0 {
					return false
				}
			case parser.ON:
				if blockLevel <= 0 {
					return true
				}
			}
		}
		return false
	}

	var joinConditionRequired = func() bool {
		if c.tokens[c.lastIdx].Token == '(' {
			return false
		}

		baseIdx := c.lastIdx - 1
		if c.tokens[baseIdx].Token == parser.LATERAL {
			baseIdx--
		}

		if c.tokens[baseIdx].Token != parser.JOIN {
			return false
		}
		if c.tokens[baseIdx-1].Token == parser.CROSS ||
			c.tokens[baseIdx-1].Token == parser.NATURAL ||
			c.tokens[baseIdx-2].Token == parser.NATURAL ||
			c.tokens[baseIdx-3].Token == parser.NATURAL {
			return false
		}
		return true
	}

	if isInUsing() {
		restrict = true
		return
	}

	if isInOn() {
		customList = append(customList, c.whereClause(line, origLine, index)...)
		if c.tokens[c.lastIdx].Token == parser.ON {
			restrict = true
		} else {
			customList = append(customList, c.candidateList(joinCandidates, true)...)
		}
		return
	}

	switch c.tokens[c.lastIdx].Token {
	case parser.LATERAL:
		restrict = true
	case parser.CROSS, parser.INNER, parser.OUTER:
		customList = append(customList, c.candidate("JOIN", true))
		restrict = true
	case parser.LEFT, parser.RIGHT, parser.FULL:
		customList = append(customList, c.candidateList([]string{"JOIN", "OUTER"}, true)...)
		restrict = true
	case parser.NATURAL:
		customList = append(customList, c.candidateList([]string{"INNER", "LEFT", "RIGHT"}, true)...)
		restrict = true
	case parser.AS, parser.USING:
		restrict = true
	}

	if joinConditionRequired() {
		customList = append(customList, c.candidate("ON", true))

		if c.tokens[c.lastIdx-2].Token != parser.FULL &&
			c.tokens[c.lastIdx-3].Token != parser.FULL {

			customList = append(customList, c.candidate("USING ()", true))
		}

		if c.tokens[c.lastIdx].Quoted == false {
			customList = append(customList, c.candidate("AS", true))
		}
		restrict = true
	} else {
		canUseLateral := false
		switch c.tokens[c.lastIdx].Token {
		case ',':
			canUseLateral = true
		case parser.JOIN:
			canUseLateral = true
			switch c.tokens[c.lastIdx-1].Token {
			case parser.RIGHT, parser.FULL:
				canUseLateral = false
			case parser.OUTER:
				switch c.tokens[c.lastIdx-2].Token {
				case parser.RIGHT, parser.FULL:
					canUseLateral = false
				}
			}

		}

		if canUseLateral {
			customList = append(customList, c.candidate("LATERAL", true))
		}

		switch c.tokens[c.lastIdx].Token {
		case '(':
			customList = append(customList, c.candidate("SELECT", true))
			fallthrough
		case parser.FROM, parser.JOIN, ',':
			tables = c.allTableCandidates(line, origLine, index)
			restrict = true
		}
	}

	if restrict {
		return
	}

	if c.tokens[c.lastIdx].Token == parser.IDENTIFIER {
		customList = append(customList, c.candidateList(joinCandidates, true)...)

		if c.tokens[c.lastIdx].Quoted == false {
			customList = append(customList, c.candidate("AS", true))
		}
	}
	return
}

func (c *Completer) whereClause(line string, origLine string, index int) readline.CandidateList {
	cands := c.SearchValues(line, origLine, index)
	if 0 < len(line) {
		cands = append(cands, c.candidate("JSON_ROW()", false))
	}
	return cands
}

func (c *Completer) SelectArgs(line string, origLine string, index int) readline.CandidateList {
	isSelectInto := false
	if c.selectIntoEnabled {
		for i := len(c.tokens) - 1; i >= 0; i-- {
			if c.tokens[i].Token == parser.INTO {
				isSelectInto = true
				break
			}
		}
	}

	return c.completeArgs(
		line,
		origLine,
		index,
		func(i int) (keywords []string, customList readline.CandidateList, breakLoop bool) {
			switch c.tokens[i].Token {
			case parser.FOR:
				if i == c.lastIdx {
					customList = append(customList, c.candidateList([]string{
						"UPDATE",
					}, false)...)
					customList.Sort()
				}
				return nil, customList, true
			case parser.LIMIT:
				if i < c.lastIdx {
					switch c.tokens[c.lastIdx].Token {
					case parser.PERCENT:
						customList = append(customList, c.candidateList([]string{
							"OFFSET",
						}, true)...)
						customList = append(customList, c.candidateList([]string{
							"WITH TIES",
							"FOR UPDATE",
						}, false)...)
					case parser.WITH:
						customList = append(customList, c.candidate("TIES", false))
					case parser.TIES:
						customList = append(customList, c.candidateList([]string{
							"OFFSET",
						}, true)...)
						customList = append(customList, c.candidate("FOR UPDATE", false))
					default:
						customList = append(customList, c.candidateList([]string{
							"OFFSET",
						}, true)...)
						customList = append(customList, c.candidateList([]string{
							"PERCENT",
							"WITH TIES",
							"FOR UPDATE",
						}, false)...)
						customList = append(customList, c.SearchValues(line, origLine, index)...)
					}
				} else {
					customList = append(customList, c.SearchValues(line, origLine, index)...)
				}
				customList.Sort()
				return keywords, customList, true
			case parser.FETCH:
				if i == c.lastIdx {
					afterOffset := false

				CompleterSelectArgsSearchOffsetLoop:
					for j := c.lastIdx - 1; j >= 0; j-- {
						switch c.tokens[j].Token {
						case parser.OFFSET:
							afterOffset = true
							break CompleterSelectArgsSearchOffsetLoop
						case parser.ORDER, parser.HAVING, parser.GROUP, parser.FROM, parser.SELECT:
							break CompleterSelectArgsSearchOffsetLoop
						}
					}

					if afterOffset {
						keywords = append(keywords, "NEXT")
					} else {
						keywords = append(keywords, "FIRST")
					}
				} else {
					switch c.tokens[c.lastIdx].Token {
					case parser.ROW, parser.ROWS, parser.PERCENT:
						customList = append(customList, c.candidateList([]string{
							"ONLY",
							"WITH TIES",
							"FOR UPDATE",
						}, false)...)
					case parser.WITH:
						customList = append(customList, c.candidateList([]string{
							"TIES",
						}, false)...)
					case parser.ONLY, parser.TIES:
						customList = append(customList, c.candidateList([]string{
							"FOR UPDATE",
						}, false)...)
					case parser.FIRST, parser.NEXT:
						//Do nothing
					default:
						if c.tokens[c.lastIdx].Literal == "1" {
							keywords = append(keywords, "ROW")
						} else {
							keywords = append(keywords, "ROWS")
						}
						keywords = append(keywords, "PERCENT")
					}
				}

				if customList != nil {
					customList.Sort()
				}
				return keywords, customList, true
			case parser.OFFSET:
				if i < c.lastIdx {
					afterLimit := false

				CompleterSelectArgsSearchLimitLoop:
					for j := c.lastIdx - 1; j >= 0; j-- {
						switch c.tokens[j].Token {
						case parser.LIMIT:
							afterLimit = true
							break CompleterSelectArgsSearchLimitLoop
						case parser.ORDER, parser.HAVING, parser.GROUP, parser.FROM, parser.SELECT:
							break CompleterSelectArgsSearchLimitLoop
						}
					}

					if !afterLimit {
						customList = append(customList, c.candidateList([]string{
							"FETCH",
						}, true)...)
					}

					customList = append(customList, c.candidateList([]string{
						"FOR UPDATE",
					}, false)...)
					if c.tokens[c.lastIdx].Token != parser.ROW && c.tokens[c.lastIdx].Token != parser.ROWS {
						if c.tokens[c.lastIdx].Literal == "1" {
							customList = append(customList, c.candidateList([]string{
								"ROW",
							}, false)...)
						} else {
							customList = append(customList, c.candidateList([]string{
								"ROWS",
							}, false)...)
						}
					}
				}
				customList = append(customList, c.SearchValues(line, origLine, index)...)
				customList.Sort()
				return nil, customList, true
			case parser.ORDER:
				if i == c.lastIdx {
					keywords = append(keywords, "BY")
				} else if i < c.lastIdx-1 && c.tokens[c.lastIdx].Token != ',' {
					switch c.tokens[c.lastIdx].Token {
					case parser.ASC, parser.DESC:
						customList = append(customList, c.candidateList([]string{
							"NULLS FIRST",
							"NULLS LAST",
							"FOR UPDATE",
						}, false)...)
						customList = append(customList, c.candidateList([]string{
							"OFFSET",
							"FETCH",
							"LIMIT",
						}, true)...)
					case parser.NULLS:
						customList = append(customList, c.candidateList([]string{
							"FIRST",
							"LAST",
						}, false)...)
					case parser.FIRST, parser.LAST:
						customList = append(customList, c.candidateList([]string{
							"OFFSET",
							"FETCH",
							"LIMIT",
						}, true)...)
						customList = append(customList, c.candidateList([]string{
							"FOR UPDATE",
						}, false)...)
					default:
						customList = append(customList, c.candidateList([]string{
							"ASC",
							"DESC",
							"NULLS FIRST",
							"NULLS LAST",
							"FOR UPDATE",
						}, false)...)
						customList = append(customList, c.candidateList([]string{
							"OFFSET",
							"FETCH",
							"LIMIT",
						}, true)...)
						customList = append(customList, c.SearchValues(line, origLine, index)...)
						customList = append(customList, c.aggregateFunctionCandidateList(line)...)
						customList = append(customList, c.analyticFunctionCandidateList(line)...)
					}
				} else {
					customList = append(customList, c.SearchValues(line, origLine, index)...)
					customList = append(customList, c.aggregateFunctionCandidateList(line)...)
					customList = append(customList, c.analyticFunctionCandidateList(line)...)
				}

				if customList != nil {
					customList.Sort()
				}
				return keywords, customList, true
			case parser.ALL:
				if i == c.lastIdx && 0 <= c.lastIdx-1 {
					switch c.tokens[i-1].Token {
					case parser.UNION, parser.EXCEPT, parser.INTERSECT:
						keywords = append(keywords, "SELECT")
						return keywords, nil, true
					}
				}
			case parser.UNION, parser.EXCEPT, parser.INTERSECT:
				if i == c.lastIdx {
					keywords = append(keywords,
						"ALL",
						"SELECT",
					)
				}
				return keywords, nil, true
			case parser.HAVING:
				customList = append(customList, c.whereClause(line, origLine, index)...)
				customList = append(customList, c.aggregateFunctionCandidateList(line)...)
				if i < c.lastIdx {
					customList = append(customList, c.candidateList([]string{
						"ORDER BY",
						"OFFSET",
						"FETCH",
						"LIMIT",
					}, true)...)
					if !isSelectInto {
						customList = append(customList, c.candidateList([]string{
							"UNION",
							"EXCEPT",
							"INTERSECT",
						}, true)...)
					}
					customList = append(customList, c.candidateList([]string{
						"FOR UPDATE",
					}, false)...)
				}
				customList.Sort()
				return keywords, customList, true
			case parser.GROUP:
				if i == c.lastIdx {
					keywords = append(keywords, "BY")
				} else {
					customList = append(customList, c.whereClause(line, origLine, index)...)
					if i < c.lastIdx-1 && c.tokens[c.lastIdx].Token != ',' {
						customList = append(customList, c.candidateList([]string{
							"HAVING",
							"ORDER BY",
							"OFFSET",
							"FETCH",
							"LIMIT",
						}, true)...)
						if !isSelectInto {
							customList = append(customList, c.candidateList([]string{
								"UNION",
								"EXCEPT",
								"INTERSECT",
							}, true)...)
						}
						customList = append(customList, c.candidateList([]string{
							"FOR UPDATE",
						}, false)...)
					}
					customList.Sort()
				}
				return keywords, customList, true
			case parser.WHERE:
				customList = append(customList, c.whereClause(line, origLine, index)...)
				if i < c.lastIdx {
					customList = append(customList, c.candidateList([]string{
						"GROUP BY",
						"HAVING",
						"ORDER BY",
						"OFFSET",
						"FETCH",
						"LIMIT",
					}, true)...)
					if !isSelectInto {
						customList = append(customList, c.candidateList([]string{
							"UNION",
							"EXCEPT",
							"INTERSECT",
						}, true)...)
					}
					customList = append(customList, c.candidateList([]string{
						"FOR UPDATE",
					}, false)...)
				}
				customList.Sort()
				return keywords, customList, true
			case parser.FROM:
				tables, clist, restrict := c.fromClause(i, line, origLine, index)
				if !restrict {
					if i < c.lastIdx && c.tokens[c.lastIdx].Token != ',' {
						clist = append(clist, c.candidateList([]string{
							"WHERE",
							"GROUP BY",
							"HAVING",
							"ORDER BY",
							"OFFSET",
							"FETCH",
							"LIMIT",
						}, true)...)
						if !isSelectInto {
							clist = append(clist, c.candidateList([]string{
								"UNION",
								"EXCEPT",
								"INTERSECT",
							}, true)...)
						}
						clist = append(clist, c.candidateList([]string{
							"FOR UPDATE",
						}, false)...)
					}
				}
				clist.Sort()
				return nil, append(tables, clist...), true
			case parser.INTO:
				switch c.tokens[c.lastIdx].Token {
				case parser.INTO, ',':
					customList = append(customList, c.candidateList(c.varList, false)...)
				case parser.VARIABLE:
					customList = append(customList, c.candidateList([]string{
						"FROM",
						"WHERE",
						"GROUP BY",
						"HAVING",
						"ORDER BY",
						"OFFSET",
						"FETCH",
						"LIMIT",
					}, true)...)
					customList = append(customList, c.candidateList([]string{
						"FOR UPDATE",
					}, false)...)
				}
				if customList != nil {
					customList.Sort()
				}
				return keywords, customList, true
			case parser.SELECT:
				if i == c.lastIdx {
					if 0 < len(line) {
						keywords = append(keywords, "DISTINCT")
					}
				} else {
					lastIdx := c.lastIdx
					if i+1 < len(c.tokens) && c.tokens[i+1].Token == parser.DISTINCT {
						lastIdx--
					}
					if i < lastIdx && c.tokens[c.lastIdx].Token != ',' {
						customList = append(customList, c.candidateList([]string{
							"AS",
							"FROM",
							"WHERE",
							"GROUP BY",
							"HAVING",
							"ORDER BY",
							"OFFSET",
							"FETCH",
							"LIMIT",
							"UNION",
							"EXCEPT",
							"INTERSECT",
						}, true)...)
						customList = append(customList, c.candidateList([]string{
							"FOR UPDATE",
						}, false)...)
						if c.selectIntoEnabled {
							customList = append(customList, c.candidateList([]string{
								"INTO",
							}, true)...)
						}
					}
				}
				customList = append(customList, c.SearchValues(line, origLine, index)...)
				customList = append(customList, c.aggregateFunctionCandidateList(line)...)
				customList = append(customList, c.analyticFunctionCandidateList(line)...)
				customList.Sort()
				return keywords, customList, true
			}
			return nil, nil, false
		},
	)
}

func (c *Completer) InsertArgs(line string, origLine string, index int) readline.CandidateList {
	return c.completeArgs(
		line,
		origLine,
		index,
		func(i int) (keywords []string, customList readline.CandidateList, breakLoop bool) {
			switch c.tokens[i].Token {
			case parser.SELECT:
				return nil, c.SelectArgs(line, origLine, index), true
			case parser.VALUES:
				customList = c.SearchValues(line, origLine, index)
				if 0 < len(line) {
					customList = append(customList, c.candidate("JSON_ROW()", false))
				}
				customList.Sort()
				return nil, customList, true
			case parser.INTO:
				if i == c.lastIdx {
					customList = c.allTableCandidatesWithSpaceForUpdate(line, origLine, index)
				} else {
					if c.tokens[c.lastIdx-1].Token == parser.INTO || c.tokens[c.lastIdx].Token == ')' {
						keywords = append(keywords, "VALUES", "SELECT")
					}
				}
				return keywords, customList, true
			case parser.INSERT:
				if i == c.lastIdx {
					keywords = append(keywords, "INTO")
				}
				return keywords, nil, true
			}
			return nil, nil, false
		},
	)
}

func (c *Completer) UpdateArgs(line string, origLine string, index int) readline.CandidateList {
	return c.completeArgs(
		line,
		origLine,
		index,
		func(i int) (keywords []string, customList readline.CandidateList, breakLoop bool) {
			switch c.tokens[i].Token {
			case parser.WHERE:
				customList = append(customList, c.whereClause(line, origLine, index)...)
				customList.Sort()
				return nil, customList, true
			case parser.FROM:
				tables, clist, restrict := c.fromClause(i, line, origLine, index)
				if !restrict {
					if i < c.lastIdx && c.tokens[c.lastIdx].Token != ',' {
						clist = append(clist, c.candidate("WHERE", true))
					}
				}
				clist.Sort()
				return nil, append(tables, clist...), true
			case parser.SET:
				if c.tokens[c.lastIdx].Token != parser.SET && c.tokens[c.lastIdx].Token != ',' &&
					c.tokens[c.lastIdx-1].Token != parser.SET && c.tokens[c.lastIdx-1].Token != ',' {

					customList = append(customList, c.SearchValues(line, origLine, index)...)
					customList.Sort()
					if c.tokens[c.lastIdx].Token != '=' {
						keywords = append(keywords, "FROM", "WHERE")
					}
				}
				return keywords, customList, true
			case parser.UPDATE:
				if c.tokens[c.lastIdx].Token == parser.UPDATE || c.tokens[c.lastIdx].Token == ',' {
					customList = c.allTableCandidatesForUpdate(line, origLine, index)
				} else {
					keywords = append(keywords, "SET")
				}
				return keywords, customList, true
			}
			return nil, nil, false
		},
	)
}

func (c *Completer) ReplaceArgs(line string, origLine string, index int) readline.CandidateList {
	usingExists := false

	return c.completeArgs(
		line,
		origLine,
		index,
		func(i int) (keywords []string, customList readline.CandidateList, breakLoop bool) {
			switch c.tokens[i].Token {
			case parser.SELECT:
				return nil, c.SelectArgs(line, origLine, index), true
			case parser.VALUES:
				customList = c.SearchValues(line, origLine, index)
				if 0 < len(line) {
					customList = append(customList, c.candidate("JSON_ROW()", false))
				}
				customList.Sort()
				return nil, customList, true
			case parser.USING:
				usingExists = true
			case parser.INTO:
				if i == c.lastIdx {
					customList = c.allTableCandidatesWithSpaceForUpdate(line, origLine, index)
				} else if c.tokens[c.lastIdx-1].Token == parser.INTO || (!usingExists && c.tokens[c.lastIdx].Token == ')') {
					customList = append(customList, c.candidate("USING ()", true))
				} else if usingExists && c.tokens[c.lastIdx].Token == ')' {
					keywords = append(keywords, "VALUES", "SELECT")
				}
				return keywords, customList, true
			case parser.REPLACE:
				if i == c.lastIdx {
					keywords = append(keywords, "INTO")
				}
				return keywords, nil, true
			}
			return nil, nil, false
		},
	)
}

func (c *Completer) DeleteArgs(line string, origLine string, index int) readline.CandidateList {
	return c.completeArgs(
		line,
		origLine,
		index,
		func(i int) (keywords []string, customList readline.CandidateList, breakLoop bool) {
			switch c.tokens[i].Token {
			case parser.WHERE:
				customList = append(customList, c.whereClause(line, origLine, index)...)
				customList.Sort()
				return nil, customList, true
			case parser.FROM:
				if c.tokens[i-1].Token == parser.DELETE {
					if c.tokens[c.lastIdx].Token == parser.FROM || c.tokens[c.lastIdx].Token == ',' {
						customList = c.allTableCandidates(line, origLine, index)
					} else {
						customList = append(customList, c.candidate("WHERE", true))
					}
				} else {
					tables, clist, restrict := c.fromClause(i, line, origLine, index)
					if !restrict {
						if i < c.lastIdx && c.tokens[c.lastIdx].Token != ',' {
							clist = append(clist, c.candidate("WHERE", true))
						}
					}
					clist.Sort()
					customList = append(tables, clist...)
				}
				return keywords, customList, true
			case parser.DELETE:
				if c.tokens[c.lastIdx].Token != ',' {
					keywords = append(keywords, "FROM")
				}
				return keywords, nil, true
			}
			return nil, nil, false
		},
	)
}

func (c *Completer) CreateArgs(line string, origLine string, index int) readline.CandidateList {
	var isPrepositionOfSelect = func(i int) bool {
		for j := i; j >= 0; j-- {
			if c.tokens[j].Token == ';' || c.tokens[j].Token == parser.TABLE {
				break
			}
			if c.tokens[j].Token == parser.SELECT {
				return false
			}
		}
		return true
	}

	return c.completeArgs(
		line,
		origLine,
		index,
		func(i int) (keywords []string, customList readline.CandidateList, breakLoop bool) {
			switch c.tokens[i].Token {
			case parser.SELECT:
				return nil, c.SelectArgs(line, origLine, index), true
			case parser.AS:
				if isPrepositionOfSelect(i) {
					return []string{"SELECT"}, nil, true
				}
			case parser.TABLE:
				if (c.tokens[c.lastIdx].Token == ')' && c.BracketIsEnclosed()) ||
					i == c.lastIdx-1 {
					return []string{"AS", "SELECT"}, nil, true
				}
			case parser.CREATE:
				if i == c.lastIdx {
					return []string{"TABLE"}, nil, true
				}
			}
			return nil, nil, false
		},
	)
}

func (c *Completer) AlterArgs(line string, origLine string, index int) readline.CandidateList {
	operations := []string{
		"ADD",
		"DROP",
		"RENAME",
		"SET",
	}

	addPositions := []string{
		"FIRST",
		"LAST",
		"AFTER",
		"BEFORE",
	}

	var columnsInTable = func(i int, appendSpace bool) readline.CandidateList {
		var tableName string
		for j := i; j >= 0; j-- {
			if c.tokens[j].Token == parser.TABLE {
				tableName = c.tokens[j+1].Literal
				break
			}
		}

		var clist readline.CandidateList
		if 0 < len(tableName) {
			clist = c.identifierList(c.ColumnList(tableName, c.scope.Tx.Flags.Repository), appendSpace)
		}
		return clist
	}

	return c.completeArgs(
		line,
		origLine,
		index,
		func(i int) (keywords []string, customList readline.CandidateList, breakLoop bool) {
			switch c.tokens[i].Token {
			case parser.AFTER, parser.BEFORE:
				return nil, columnsInTable(i, false), true
			case parser.ADD:
				prevToken := c.tokens[c.lastIdx].Token

				if i+1 < len(c.tokens) && c.tokens[i+1].Token == '(' {
					if c.BracketIsEnclosed() {
						if prevToken == ')' {
							return nil, c.candidateList(addPositions, true), true
						}
					} else {
						switch {
						case prevToken == '(' || prevToken == ',':
							//Column
						case c.tokens[c.lastIdx-1].Token == '(' || c.tokens[c.lastIdx-1].Token == ',':
							return []string{"DEFAULT"}, nil, true
						default:
							return nil, c.SearchValues(line, origLine, index), true
						}
					}
				} else {
					switch i {
					case c.lastIdx:
						//Column
					case c.lastIdx - 1:
						return []string{"DEFAULT"}, c.candidateList(addPositions, true), true
					default:
						return nil, c.SearchValues(line, origLine, index), true
					}
				}
			case parser.DROP:
				if i+1 < len(c.tokens) && c.tokens[i+1].Token == '(' {
					if !c.BracketIsEnclosed() {
						return nil, columnsInTable(i, false), true
					}
				} else {
					switch i {
					case c.lastIdx:
						return nil, columnsInTable(i, false), true
					}
				}
			case parser.RENAME:
				switch i {
				case c.lastIdx:
					return nil, columnsInTable(i, true), true
				case c.lastIdx - 1:
					return []string{"TO"}, nil, true
				}
			case parser.SET:
				switch i {
				case c.lastIdx:
					return FileAttributeList, nil, true
				case c.lastIdx - 1:
					return []string{"TO"}, nil, true
				case c.lastIdx - 2:
					switch strings.ToUpper(c.tokens[c.lastIdx-1].Literal) {
					case TableFormat:
						return nil, c.candidateList(c.tableFormatList(), false), true
					case TableDelimiter:
						return nil, c.candidateList(delimiterCandidates, false), true
					case TableDelimiterPositions:
						return nil, c.candidateList(delimiterPositionsCandidates, false), true
					case TableEncoding:
						return nil, c.candidateList(exportEncodingsCandidates, false), true
					case TableLineBreak:
						return nil, c.candidateList(c.lineBreakList(), false), true
					case TableJsonEscape:
						return nil, c.candidateList(c.jsonEscapeTypeList(), false), true
					case TableHeader, TableEncloseAll, TablePrettyPrint:
						return nil, c.candidateList([]string{ternary.TRUE.String(), ternary.FALSE.String()}, false), true
					}
				}
			case parser.TABLE:
				switch i {
				case c.lastIdx:
					return nil, c.allTableCandidatesWithSpaceForUpdate(line, origLine, index), true
				case c.lastIdx - 1:
					return operations, nil, true
				}
			case parser.ALTER:
				if i == c.lastIdx {
					return []string{"TABLE"}, nil, true
				}
			default:
				return nil, nil, false
			}
			return nil, nil, true
		},
	)
}

func (c *Completer) DeclareArgs(line string, origLine string, index int) readline.CandidateList {
	return c.completeArgs(
		line,
		origLine,
		index,
		func(i int) (keywords []string, customList readline.CandidateList, breakLoop bool) {
			switch c.tokens[i].Token {
			case parser.SELECT:
				return nil, c.SelectArgs(line, origLine, index), true
			case parser.CURSOR:
				if i == c.lastIdx {
					return []string{"FOR"}, nil, true
				} else if i == c.lastIdx-1 && c.tokens[c.lastIdx].Token == parser.FOR {
					return []string{"SELECT"}, c.candidateList(c.statementList, false), true
				}
			case parser.VIEW:
				switch c.tokens[c.lastIdx].Token {
				case parser.AS:
					return []string{"SELECT"}, nil, true
				case parser.VIEW, ')':
					return []string{"AS"}, nil, true
				}
			case parser.AGGREGATE, parser.FUNCTION:
			case parser.VAR:
				switch c.tokens[c.lastIdx].Token {
				case parser.VARIABLE, ',':
				default:
					return nil, c.SearchValues(line, origLine, index), true
				}
			case parser.DECLARE:
				if i == c.lastIdx-1 && c.tokens[c.lastIdx].Token != parser.VARIABLE {
					obj := []string{
						"CURSOR",
						"VIEW",
						"FUNCTION",
						"AGGREGATE",
					}
					return obj, nil, true
				} else if i < c.lastIdx-1 && c.tokens[i+1].Token == parser.VARIABLE {
					switch c.tokens[c.lastIdx].Token {
					case parser.VARIABLE, ',':
					default:
						return nil, c.SearchValues(line, origLine, index), true
					}
				}
			default:
				return nil, nil, false
			}

			if c.tokens[c.lastIdx].Token == parser.SUBSTITUTION_OP {
				return nil, c.SearchValues(line, origLine, index), true
			}

			return nil, nil, true
		},
	)
}

func (c *Completer) PrepareArgs(line string, origLine string, index int) readline.CandidateList {
	return c.completeArgs(
		line,
		origLine,
		index,
		func(i int) (keywords []string, customList readline.CandidateList, breakLoop bool) {
			switch c.tokens[i].Token {
			case parser.PREPARE:
				if i == c.lastIdx-1 {
					return []string{"FROM"}, nil, true
				}
			default:
				return nil, nil, false
			}

			return nil, nil, true
		},
	)
}

func (c *Completer) FetchArgs(line string, origLine string, index int) readline.CandidateList {
	positions := []string{
		"NEXT",
		"PRIOR",
		"FIRST",
		"LAST",
		"ABSOLUTE",
		"RELATIVE",
	}

	return c.completeArgs(
		line,
		origLine,
		index,
		func(i int) (keywords []string, customList readline.CandidateList, breakLoop bool) {
			switch c.tokens[i].Token {
			case parser.INTO:
				return nil, c.SearchValues(line, origLine, index), true
			case parser.NEXT, parser.PRIOR, parser.FIRST, parser.LAST:
				switch i {
				case c.lastIdx:
					return c.cursorList, nil, true
				case c.lastIdx - 1:
					return []string{"INTO"}, nil, true
				}
				return nil, nil, true
			case parser.ABSOLUTE, parser.RELATIVE:
				switch i {
				case c.lastIdx:
					return nil, c.SearchValuesWithSpace(line, origLine, index), true
				default:
					if c.tokens[c.lastIdx].Token == parser.IDENTIFIER {
						return []string{"INTO"}, nil, true
					} else {
						return c.cursorList, nil, true
					}
				}
			case parser.FETCH:
				switch i {
				case c.lastIdx:
					return c.cursorList, c.candidateList(positions, true), true
				case c.lastIdx - 1:
					return []string{"INTO"}, nil, true
				}
				return nil, nil, true
			}
			return nil, nil, false
		},
	)
}

func (c *Completer) SetArgs(line string, origLine string, index int) readline.CandidateList {
	return c.completeArgs(
		line,
		origLine,
		index,
		func(i int) (keywords []string, customList readline.CandidateList, breakLoop bool) {
			switch c.tokens[i].Token {
			case parser.TO:
				if i == c.lastIdx && c.tokens[c.lastIdx-1].Token == parser.FLAG {
					switch strings.ToUpper(c.tokens[c.lastIdx-1].Literal) {
					case cmd.RepositoryFlag:
						return nil, c.SearchDirs(line, origLine, index), true
					case cmd.TimezoneFlag:
						return nil, c.candidateList([]string{"Local", "UTC"}, false), true
					case cmd.ImportFormatFlag:
						return nil, c.candidateList(c.importFormatList(), false), true
					case cmd.DelimiterFlag, cmd.ExportDelimiterFlag:
						return nil, c.candidateList(delimiterCandidates, false), true
					case cmd.DelimiterPositionsFlag, cmd.ExportDelimiterPositionsFlag:
						return nil, c.candidateList(delimiterPositionsCandidates, false), true
					case cmd.EncodingFlag:
						return nil, c.candidateList(c.encodingList(), false), true
					case cmd.ExportEncodingFlag:
						return nil, c.candidateList(exportEncodingsCandidates, false), true
					case cmd.AnsiQuotesFlag, cmd.StrictEqualFlag, cmd.AllowUnevenFieldsFlag,
						cmd.NoHeaderFlag, cmd.WithoutNullFlag,
						cmd.WithoutHeaderFlag, cmd.EncloseAllFlag, cmd.PrettyPrintFlag,
						cmd.StripEndingLineBreakFlag, cmd.EastAsianEncodingFlag,
						cmd.CountDiacriticalSignFlag, cmd.CountFormatCodeFlag,
						cmd.ColorFlag, cmd.QuietFlag, cmd.StatsFlag:
						return nil, c.candidateList([]string{ternary.TRUE.String(), ternary.FALSE.String()}, false), true
					case cmd.FormatFlag:
						return nil, c.candidateList(c.tableFormatList(), false), true
					case cmd.LineBreakFlag:
						return nil, c.candidateList(c.lineBreakList(), false), true
					case cmd.JsonEscapeFlag:
						return nil, c.candidateList(c.jsonEscapeTypeList(), false), true
					}
				}
				return nil, c.SearchValues(line, origLine, index), true
			case parser.SET:
				switch i {
				case c.lastIdx:
					return nil, append(c.candidateList(c.flagList, true), c.candidateList(c.environmentVariableList(line), true)...), true
				case c.lastIdx - 1:
					return []string{"TO"}, nil, true
				}
				return nil, nil, true
			}
			return nil, nil, false
		},
	)
}

func (c *Completer) UsingArgs(line string, origLine string, index int) readline.CandidateList {
	return c.completeArgs(
		line,
		origLine,
		index,
		func(i int) (keywords []string, customList readline.CandidateList, breakLoop bool) {
			switch c.tokens[i].Token {
			case parser.USING:
				if c.tokens[c.lastIdx-1].Token == parser.USING || c.tokens[c.lastIdx-1].Token == ',' {
					keywords = []string{"AS"}
				}
				return keywords, c.SearchValues(line, origLine, index), true
			case parser.EXECUTE, parser.PRINTF, parser.OPEN:
				if i == c.lastIdx {
					switch c.tokens[i].Token {
					case parser.EXECUTE:
						keywords = c.statementList
					case parser.OPEN:
						keywords = c.cursorList
					}
				} else {
					keywords = append(keywords, "USING")
				}
				return keywords, c.SearchValues(line, origLine, index), true
			}
			return nil, nil, false
		},
	)
}

func (c *Completer) AddFlagArgs(line string, origLine string, index int) readline.CandidateList {
	return c.completeArgs(
		line,
		origLine,
		index,
		func(i int) (keywords []string, customList readline.CandidateList, breakLoop bool) {
			switch c.tokens[i].Token {
			case parser.TO:
				return nil, c.candidateList([]string{cmd.FlagSymbol(cmd.DatetimeFormatFlag)}, false), true
			case parser.ADD:
				if i < c.lastIdx {
					keywords = append(keywords, "TO")
				}
				return keywords, c.SearchValues(line, origLine, index), true
			}
			return nil, nil, false
		},
	)
}

func (c *Completer) RemoveFlagArgs(line string, origLine string, index int) readline.CandidateList {
	return c.completeArgs(
		line,
		origLine,
		index,
		func(i int) (keywords []string, customList readline.CandidateList, breakLoop bool) {
			switch c.tokens[i].Token {
			case parser.FROM:
				return nil, c.candidateList([]string{cmd.FlagSymbol(cmd.DatetimeFormatFlag)}, false), true
			case parser.REMOVE:
				if i < c.lastIdx {
					keywords = append(keywords, "FROM")
				}
				return keywords, c.SearchValues(line, origLine, index), true
			}
			return nil, nil, false
		},
	)

}

func (c *Completer) DisposeArgs(line string, origLine string, index int) readline.CandidateList {
	return c.completeArgs(
		line,
		origLine,
		index,
		func(i int) (keywords []string, customList readline.CandidateList, breakLoop bool) {
			switch c.tokens[i].Token {
			case parser.CURSOR:
				switch i {
				case c.lastIdx:
					return nil, c.candidateList(c.cursorList, false), true
				}
			case parser.FUNCTION:
				switch i {
				case c.lastIdx:
					return nil, c.candidateList(c.userFuncList, false), true
				}
			case parser.VIEW:
				switch i {
				case c.lastIdx:
					return nil, c.candidateList(c.viewList, false), true
				}
			case parser.PREPARE:
				switch i {
				case c.lastIdx:
					return nil, c.candidateList(c.statementList, false), true
				}
			case parser.DISPOSE:
				switch i {
				case c.lastIdx:
					var items []string
					if 0 < len(c.cursorList) {
						items = append(items, "CURSOR")
					}
					if 0 < len(c.userFuncList) {
						items = append(items, "FUNCTION")
					}
					if 0 < len(c.viewList) {
						items = append(items, "VIEW")
					}
					if 0 < len(c.statementList) {
						items = append(items, "PREPARE")
					}
					sort.Strings(items)
					list := append(c.candidateList(items, true), c.candidateList(c.varList, false)...)
					return nil, list, true
				}
			default:
				return nil, nil, false
			}
			return nil, nil, true
		},
	)
}

func (c *Completer) ShowArgs(line string, origLine string, index int) readline.CandidateList {
	var showChild = func() readline.CandidateList {
		cands := c.candidateList(ShowObjectList, false)
		cands = append(cands, c.candidate("FIELDS", true))
		cands.Sort()
		cands = append(cands, c.candidateList(c.flagList, false)...)
		return cands
	}

	return c.completeArgs(
		line,
		origLine,
		index,
		func(i int) (keywords []string, customList readline.CandidateList, breakLoop bool) {
			switch c.tokens[i].Token {
			case parser.FROM:
				switch i {
				case c.lastIdx:
					return nil, c.allTableCandidatesForUpdate(line, origLine, index), true
				}
			case parser.SHOW:
				switch i {
				case c.lastIdx:
					return nil, showChild(), true
				case c.lastIdx - 1:
					if c.tokens[c.lastIdx].Token == parser.IDENTIFIER && strings.ToUpper(c.tokens[c.lastIdx].Literal) == "FIELDS" {
						return []string{"FROM"}, nil, true
					}
				}
			default:
				return nil, nil, false
			}
			return nil, nil, true
		},
	)
}

func (c *Completer) SearchAllTablesWithSpace(line string, origLine string, index int) readline.CandidateList {
	cands := c.SearchAllTables(line, origLine, index)
	for i := range cands {
		cands[i].AppendSpace = true
	}
	return cands
}

func (c *Completer) SearchAllTables(line string, _ string, _ int) readline.CandidateList {
	tableKeys := c.scope.Tx.cachedViews.SortedKeys()
	files := c.ListFiles(line, []string{cmd.CsvExt, cmd.TsvExt, cmd.JsonExt, cmd.JsonlExt, cmd.LtsvExt, cmd.TextExt}, c.scope.Tx.Flags.Repository)

	defaultDir := c.scope.Tx.Flags.Repository
	if len(defaultDir) < 1 {
		defaultDir, _ = os.Getwd()
	}

	items := make([]string, 0, len(tableKeys)+len(files)+len(c.viewList))
	tablePath := make(map[string]bool)
	for _, k := range tableKeys {
		if view, ok := c.scope.Tx.cachedViews.Load(k); ok {
			lpath := view.FileInfo.Path
			tablePath[lpath] = true
			if filepath.Dir(lpath) == defaultDir {
				items = append(items, filepath.Base(lpath))
			} else {
				items = append(items, lpath)
			}
		}
	}

	items = append(items, c.viewList...)
	sort.Strings(items)

	for _, f := range files {
		if f != "." && f != ".." {
			abs := f
			if !filepath.IsAbs(abs) {
				abs = filepath.Join(defaultDir, abs)
			}
			if _, ok := tablePath[abs]; ok {
				continue
			}
		}
		items = append(items, f)
	}

	cands := make(readline.CandidateList, 0, len(items))
	for _, t := range items {
		cands = append(cands, readline.Candidate{Name: []rune(t), FormatAsIdentifier: true, AppendSpace: false})
	}
	return cands
}

func (c *Completer) SearchExecutableFiles(line string, origLine string, index int) readline.CandidateList {
	cands := c.SearchValues(line, origLine, index)
	files := c.ListFiles(line, []string{cmd.SqlExt, cmd.CsvqProcExt}, "")
	return append(cands, c.identifierList(files, false)...)
}

func (c *Completer) SearchDirs(line string, origLine string, index int) readline.CandidateList {
	cands := c.SearchValues(line, origLine, index)
	files := c.ListFiles(line, nil, "")
	return append(cands, c.identifierList(files, false)...)
}

func (c *Completer) SearchValuesWithSpace(line string, origLine string, index int) readline.CandidateList {
	cands := c.SearchValues(line, origLine, index)
	for i := range cands {
		cands[i].AppendSpace = true
	}
	return cands
}

func (c *Completer) SearchValues(line string, origLine string, index int) readline.CandidateList {
	if cands := c.EncloseQuotation(line, origLine, index); cands != nil {
		return cands
	}

	searchWord := strings.ToUpper(line)

	if 0 < len(c.cursorList) {
		if cands := c.CursorStatus(line, origLine, index); 0 < len(cands) {
			return cands
		}
	}

	var cands readline.CandidateList
	if len(searchWord) < 1 {
		return cands
	}

	var list []string
	if 1 < len(line) {
		list = append(list, c.runinfoList...)
		list = append(list, c.environmentVariableList(line)...)
	}
	list = append(list, c.varList...)
	list = append(list, c.funcList...)
	list = append(list,
		"TRUE",
		"FALSE",
		"UNKNOWN",
		"NULL",
	)

	for _, s := range list {
		if strings.HasPrefix(strings.ToUpper(s), searchWord) {
			cands = append(cands, readline.Candidate{Name: []rune(s), FormatAsIdentifier: false, AppendSpace: false})
		}
	}

	list = list[:0]
	list = append(list,
		"AND",
		"OR",
		"NOT",
		"IS",
		"BETWEEN",
		"LIKE",
		"IN",
		"ANY",
		"ALL",
		"EXISTS",
		"CASE",
	)

	if 0 < len(c.cursorList) {
		list = append(list, "CURSOR")
	}

	if 0 <= c.lastIdx && c.tokens[c.lastIdx].Token == '(' {
		list = append(list, "SELECT")
	}

	for _, s := range list {
		if strings.HasPrefix(strings.ToUpper(s), searchWord) {
			cands = append(cands, readline.Candidate{Name: []rune(s), FormatAsIdentifier: false, AppendSpace: true})
		}
	}

	if caseCands := c.CaseExpression(line, origLine, index); 0 < len(caseCands) {
		cands = append(cands, caseCands...)
	}

	return cands
}

func (c *Completer) CursorStatus(line string, origLine string, index int) readline.CandidateList {
	return c.completeArgs(
		line,
		origLine,
		index,
		func(i int) (keywords []string, customList readline.CandidateList, breakLoop bool) {
			switch c.tokens[i].Token {
			case parser.IN:
				switch i {
				case c.lastIdx:
					if (0 < c.lastIdx-3 &&
						c.tokens[c.lastIdx-3].Token == parser.CURSOR &&
						c.tokens[c.lastIdx-1].Token == parser.IS) ||
						0 < c.lastIdx-4 &&
							c.tokens[c.lastIdx-4].Token == parser.CURSOR &&
							c.tokens[c.lastIdx-2].Token == parser.IS &&
							c.tokens[c.lastIdx-1].Token == parser.NOT {
						return nil, c.candidateList([]string{"RANGE"}, false), true
					}
				}
			case parser.NOT:
				switch i {
				case c.lastIdx:
					if 0 < c.lastIdx-3 &&
						c.tokens[c.lastIdx-3].Token == parser.CURSOR &&
						c.tokens[c.lastIdx-1].Token == parser.IS {
						return nil, c.candidateList([]string{
							"IN RANGE",
							"OPEN",
						}, false), true
					}
				}
			case parser.IS:
				switch i {
				case c.lastIdx:
					if 0 < c.lastIdx-2 && c.tokens[c.lastIdx-2].Token == parser.CURSOR {
						return []string{"NOT"}, c.candidateList([]string{
							"IN RANGE",
							"OPEN",
						}, false), true
					}
				}
			case parser.CURSOR:
				switch i {
				case c.lastIdx:
					return c.cursorList, nil, true
				case c.lastIdx - 1:
					return []string{"IS"}, c.candidateList([]string{"COUNT"}, false), true
				}
			case parser.IDENTIFIER, parser.COUNT, parser.OPEN, parser.RANGE:
				return nil, nil, false
			}
			return nil, nil, true
		},
	)
}

func (c *Completer) caseExpressionIsNotEnclosed() bool {
	if 0 < len(c.tokens) && c.tokens[0].Token == parser.CASE {
		return false
	}

	var blockLevel = 0
	for i := 0; i < len(c.tokens); i++ {
		switch c.tokens[i].Token {
		case parser.CASE:
			blockLevel++
		case parser.END:
			blockLevel--
		}
	}
	return 0 < blockLevel
}

func (c *Completer) CaseExpression(line string, origLine string, index int) readline.CandidateList {
	caseExperEnclosed := true

	return c.completeArgs(
		line,
		origLine,
		index,
		func(i int) (keywords []string, customList readline.CandidateList, breakLoop bool) {
			if caseExperEnclosed {
				if caseExperEnclosed = !c.caseExpressionIsNotEnclosed(); caseExperEnclosed {
					return nil, nil, true
				}
			}

			switch c.tokens[i].Token {
			case parser.ELSE:
				if i < c.lastIdx {
					return nil, c.filteredCandidateList(line, []string{"END"}, false), true
				}
				return nil, nil, true
			case parser.THEN:
				if i < c.lastIdx {
					return []string{"WHEN", "ELSE"}, c.filteredCandidateList(line, []string{"END"}, false), true
				}
				return nil, nil, true
			case parser.WHEN:
				if i < c.lastIdx {
					return []string{"THEN"}, nil, true
				}
				return nil, nil, true
			case parser.CASE:
				return []string{"WHEN"}, nil, true
			}
			return nil, nil, false
		},
	)
}

func (c *Completer) EncloseQuotation(line string, origLine string, _ int) readline.CandidateList {
	runes := []rune(line)
	if 0 < len(runes) && readline.IsQuotationMark(runes[0]) && !readline.LiteralIsEnclosed(runes[0], []rune(origLine)) {
		return c.candidateList([]string{string(append([]rune(line), runes[0]))}, false)
	}

	return nil
}

func (c *Completer) ListFiles(path string, includeExt []string, repository string) []string {
	list := make([]string, 0, 10)

	if 0 < len(path) && (path[0] == '"' || path[0] == '\'' || path[0] == '`') {
		path = path[1:]
	}
	searchWord := strings.ToUpper(path)

	var defaultDir string
	if len(path) < 1 || (!filepath.IsAbs(path) && path != "." && path != ".." && filepath.Base(path) == path) {
		if 0 < len(repository) {
			defaultDir = repository
		} else {
			defaultDir, _ = os.Getwd()
		}
		path = defaultDir

		for _, v := range []string{".", ".."} {
			if len(searchWord) < 1 || strings.HasPrefix(strings.ToUpper(v), searchWord) {
				list = append(list, v)
			}
		}
	}

	if _, err := os.Stat(path); err != nil {
		path = filepath.Dir(path)
	}

	if files, err := ioutil.ReadDir(path); err == nil {

		for _, f := range files {
			if f.Name()[0] == '.' {
				continue
			}

			if !f.IsDir() && (len(includeExt) < 1 || !InStrSliceWithCaseInsensitive(filepath.Ext(f.Name()), includeExt)) {
				continue
			}

			fpath := f.Name()
			if len(defaultDir) < 1 {
				if path == "." || path == "."+string(os.PathSeparator) {
					fpath = "." + string(os.PathSeparator) + fpath
				} else {
					fpath = filepath.Join(path, fpath)
				}
			}
			if f.IsDir() {
				fpath = fpath + string(os.PathSeparator)
			}
			if len(searchWord) < 1 || strings.HasPrefix(strings.ToUpper(fpath), searchWord) {
				list = append(list, fpath)
			}
		}
	}

	return list
}

func (c *Completer) AllColumnList() []string {
	m := make(map[string]bool)
	c.scope.blocks[0].temporaryTables.Range(func(key, value interface{}) bool {
		col := c.columnList(value.(*View))
		for _, s := range col {
			if _, ok := m[s]; !ok {
				m[s] = true
			}
		}
		return true
	})

	c.scope.Tx.cachedViews.Range(func(key, value interface{}) bool {
		col := c.columnList(value.(*View))
		for _, s := range col {
			if _, ok := m[s]; !ok {
				m[s] = true
			}
		}
		return true
	})

	list := make([]string, 0, len(m))
	for k := range m {
		list = append(list, k)
	}
	sort.Strings(list)
	return list
}

func (c *Completer) ColumnList(tableName string, repository string) []string {
	if list, ok := c.tableColumns[tableName]; ok {
		return list
	}

	if view, ok := c.scope.blocks[0].temporaryTables.Load(tableName); ok {
		list := c.columnList(view)
		c.tableColumns[tableName] = list
		return list
	}

	if fpath, err := CreateFilePath(parser.Identifier{Literal: tableName}, repository); err == nil {
		if view, ok := c.scope.Tx.cachedViews.Load(fpath); ok {
			list := c.columnList(view)
			c.tableColumns[tableName] = list
			return list
		}
	}
	if fpath, err := SearchFilePathFromAllTypes(parser.Identifier{Literal: tableName}, repository); err == nil {
		if view, ok := c.scope.Tx.cachedViews.Load(fpath); ok {
			list := c.columnList(view)
			c.tableColumns[tableName] = list
			return list
		}
	}

	return nil
}

func (*Completer) columnList(view *View) []string {
	var list []string
	for _, h := range view.Header {
		list = append(list, h.Column)
	}
	return list
}

func (c *Completer) completeArgs(
	line string,
	origLine string,
	index int,
	fn func(i int) (keywords []string, customList readline.CandidateList, breakLoop bool),
) readline.CandidateList {
	if cands := c.EncloseQuotation(line, origLine, index); cands != nil {
		return cands
	}

	baseInterval := 1
	if 0 < len(line) {
		baseInterval++
	}
	if 0 < len(c.tokens) && 0 < len(line) {
		c.tokens[len(c.tokens)-1].Token = parser.IDENTIFIER
	}

	var keywords []string
	var customList readline.CandidateList
	var breakLoop bool

	for i := len(c.tokens) - 1; i >= 0; i-- {
		if keywords, customList, breakLoop = fn(i); breakLoop {
			break
		}
	}

	cands := c.filteredCandidateList(line, keywords, true)
	cands.Sort()
	if 0 < len(customList) {
		cands = append(cands, customList...)
	}
	return cands
}

func (c *Completer) UpdateTokens(line string, origLine string) {
	c.tokens = c.tokens[:0]
	s := new(parser.Scanner)
	s.Init(origLine, "", false, c.scope.Tx.Flags.AnsiQuotes)
	for {
		t, _ := s.Scan()
		if t.Token == parser.EOF {
			break
		}
		c.tokens = append(c.tokens, t)
	}

	if 0 < len(c.tokens) {
		c.tokens = c.tokens[c.searchStartIndex():]
	}

	c.combineSubqueryTokens()
	c.combineTableObject()
	c.combineFunction()
	c.SetLastIndex(line)
}

func (c *Completer) SetLastIndex(line string) {
	c.lastIdx = len(c.tokens) - 1
	if 0 < len(c.tokens) && 0 < len(line) {
		c.tokens[len(c.tokens)-1].Token = parser.IDENTIFIER
		c.lastIdx--
	}
}

func (c *Completer) searchStartIndex() int {
	idx := 0
	blockLevel := 0
	isStatement := false

StartIndexLoop:
	for i := len(c.tokens) - 1; i >= 0; i-- {
		switch c.tokens[i].Token {
		case ';':
			idx = i + 1
			isStatement = true
			break StartIndexLoop
		case '(':
			blockLevel--
			switch {
			case blockLevel < 0:
				switch {
				case i+1 < len(c.tokens) && c.tokens[i+1].Token == parser.SELECT:
					idx = i + 1
					break StartIndexLoop
				case 0 <= i-1 && (c.isTableObject(c.tokens[i-1]) || c.isFunction(c.tokens[i-1])):
					idx = i - 1
					break StartIndexLoop
				}
			}
		case ')':
			blockLevel++
		}
	}

	if 0 < len(c.varList) &&
		(idx == 0 || isStatement) &&
		(idx < len(c.tokens) && (c.tokens[idx].Token == parser.SELECT || c.tokens[idx].Token == parser.WITH)) {
		c.selectIntoEnabled = true
	} else {
		c.selectIntoEnabled = false
	}
	return idx
}

func (c *Completer) combineSubqueryTokens() {
	combined := make([]parser.Token, 0, cap(c.tokens))
	blockLevel := 0
	for i := 0; i < len(c.tokens); i++ {
		if 0 < blockLevel {
			switch c.tokens[i].Token {
			case '(':
				blockLevel++
			case ')':
				blockLevel--
				if blockLevel == 0 {
					combined = append(combined, parser.Token{Token: parser.IDENTIFIER, Literal: dummySubquery})
				}
			}
			continue
		}

		if c.tokens[i].Token == '(' && i+1 < len(c.tokens) && c.tokens[i+1].Token == parser.SELECT {
			blockLevel++
			i++
		} else {
			combined = append(combined, c.tokens[i])
		}

	}
	c.tokens = combined
}

func (c *Completer) combineTableObject() {
	combined := make([]parser.Token, 0, cap(c.tokens))
	blockLevel := 0
	tableIdx := 0
	for i := 0; i < len(c.tokens); i++ {
		if 0 < blockLevel {
			switch c.tokens[i].Token {
			case '(':
				blockLevel++
			case ')':
				blockLevel--
				if blockLevel == 0 {
					lit := dummyTableObject
					if 0 < tableIdx {
						lit = c.tokens[tableIdx].Literal
					}
					combined = append(combined, parser.Token{Token: parser.IDENTIFIER, Literal: lit})
				}
			case ',':
				if tableIdx == 0 && blockLevel == 1 {
					tableIdx = i + 1
				}
			}
			continue
		}

		if 1 < i && c.isTableObject(c.tokens[i]) && i+1 < len(c.tokens) && c.tokens[i+1].Token == '(' {
			blockLevel++
			i++
		} else {
			combined = append(combined, c.tokens[i])
		}

	}
	c.tokens = combined
}

func (c *Completer) combineFunction() {
	combined := make([]parser.Token, 0, cap(c.tokens))
	blockLevel := 0
	funcName := ""
	for i := 0; i < len(c.tokens); i++ {
		if 0 < blockLevel {
			switch c.tokens[i].Token {
			case '(':
				blockLevel++
			case ')':
				blockLevel--
				if blockLevel == 0 {
					if i+2 < len(c.tokens) && c.tokens[i+1].Token == parser.OVER && c.tokens[i+2].Token == '(' {
						blockLevel++
						i = i + 2
					} else if i+4 < len(c.tokens) && c.tokens[i+3].Token == parser.OVER && c.tokens[i+4].Token == '(' {
						blockLevel++
						i = i + 4
					} else {
						combined = append(combined, parser.Token{Token: parser.FUNCTION, Literal: funcName})
					}
				}
			}
			continue
		}

		if 0 < i && c.isFunction(c.tokens[i]) && i+1 < len(c.tokens) && c.tokens[i+1].Token == '(' {
			funcName = c.tokens[i].Literal
			blockLevel++
			i++
		} else {
			combined = append(combined, c.tokens[i])
		}

	}
	c.tokens = combined
}

func (c *Completer) isTableObject(token parser.Token) bool {
	switch token.Token {
	case parser.CSV, parser.JSON, parser.JSONL, parser.FIXED, parser.LTSV, parser.JSON_TABLE:
		return true
	}
	return false
}

func (c *Completer) isFunction(token parser.Token) bool {
	if token.Token == parser.IDENTIFIER {
		if _, ok := Functions[strings.ToUpper(token.Literal)]; ok {
			return true
		}
		return InStrSliceWithCaseInsensitive(token.Literal, c.userFuncList)
	}

	return token.Token == parser.SUBSTRING ||
		token.Token == parser.JSON_OBJECT ||
		token.Token == parser.IF ||
		token.Token == parser.AGGREGATE_FUNCTION ||
		token.Token == parser.COUNT ||
		token.Token == parser.LIST_FUNCTION ||
		token.Token == parser.ANALYTIC_FUNCTION ||
		token.Token == parser.FUNCTION_NTH ||
		token.Token == parser.FUNCTION_WITH_INS
}

func (c *Completer) BracketIsEnclosed() bool {
	var blockLevel = 0
	for i := 0; i < len(c.tokens); i++ {
		switch c.tokens[i].Token {
		case '(':
			blockLevel++
		case ')':
			blockLevel--
		}
	}
	return blockLevel < 1
}

func (c *Completer) candidateList(list []string, appendSpace bool) readline.CandidateList {
	cands := make(readline.CandidateList, 0, len(list))
	for _, v := range list {
		cands = append(cands, c.candidate(v, appendSpace))
	}
	return cands
}

func (c *Completer) identifierList(list []string, appendSpace bool) readline.CandidateList {
	cands := make(readline.CandidateList, 0, len(list))
	for _, v := range list {
		cands = append(cands, c.identifier(v, appendSpace))
	}
	return cands
}

func (c *Completer) filteredCandidateList(line string, list []string, appendSpace bool) readline.CandidateList {
	searchWord := strings.ToUpper(line)

	cands := make(readline.CandidateList, 0, len(list))
	for _, v := range list {
		if len(searchWord) < 1 || strings.HasPrefix(strings.ToUpper(v), searchWord) {
			cands = append(cands, c.candidate(v, appendSpace))
		}
	}
	return cands
}

func (c *Completer) candidate(candidate string, appendSpace bool) readline.Candidate {
	return readline.Candidate{Name: []rune(candidate), FormatAsIdentifier: false, AppendSpace: appendSpace}
}

func (c *Completer) identifier(candidate string, appendSpace bool) readline.Candidate {
	return readline.Candidate{Name: []rune(candidate), FormatAsIdentifier: true, AppendSpace: appendSpace}
}

func (c *Completer) aggregateFunctionCandidateList(line string) readline.CandidateList {
	if len(line) < 1 {
		return nil
	}
	return c.filteredCandidateList(line, c.aggFuncList, false)
}

func (c *Completer) analyticFunctionCandidateList(line string) readline.CandidateList {
	var cands readline.CandidateList
	if 0 <= c.lastIdx && c.tokens[c.lastIdx].Token == parser.FUNCTION {
		if InStrSliceWithCaseInsensitive(c.tokens[c.lastIdx].Literal, []string{
			"FIRST_VALUE",
			"LAST_VALUE",
			"NTH_VALUE",
			"LAG_VALUE",
			"LEAD_VALUE",
		}) {
			cands = append(cands, c.candidate("IGNORE NULLS", true))
		}

		if InStrSliceWithCaseInsensitive(c.tokens[c.lastIdx].Literal, c.analyticFuncs) {
			cands = append(cands, c.candidate("OVER", true))
		}
	}

	if 0 < len(line) {
		cands = append(cands, c.filteredCandidateList(line, c.analyticFuncList, false)...)
	}
	return cands
}

func (c *Completer) environmentVariableList(line string) []string {
	if 2 < len(line) && strings.HasPrefix(line, cmd.EnvironmentVariableSign+"`") {
		return c.enclosedEnvList
	}
	return c.envList
}

func (c *Completer) tableFormatList() []string {
	list := make([]string, 0, len(cmd.FormatLiteral))
	for _, v := range cmd.FormatLiteral {
		list = append(list, v)
	}
	sort.Strings(list)
	return list
}

func (c *Completer) importFormatList() []string {
	list := make([]string, 0, len(cmd.ImportFormats))
	for _, v := range cmd.ImportFormats {
		list = append(list, cmd.FormatLiteral[v])
	}
	sort.Strings(list)
	return list
}

func (c *Completer) encodingList() []string {
	list := make([]string, 0, len(text.EncodingLiteral))
	for _, v := range text.EncodingLiteral {
		list = append(list, v)
	}
	sort.Strings(list)
	return list
}

func (c *Completer) lineBreakList() []string {
	list := make([]string, 0, len(text.LineBreakLiteral))
	for _, v := range text.LineBreakLiteral {
		list = append(list, v)
	}
	sort.Strings(list)
	return list
}

func (c *Completer) jsonEscapeTypeList() []string {
	list := make([]string, 0, len(cmd.JsonEscapeTypeLiteral))
	for _, v := range cmd.JsonEscapeTypeLiteral {
		list = append(list, v)
	}
	sort.Strings(list)
	return list
}
