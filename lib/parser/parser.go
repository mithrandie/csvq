//line parser.y:2
package parser

import __yyfmt__ "fmt"

//line parser.y:2
//line parser.y:5
type yySymType struct {
	yys         int
	program     []Statement
	statement   Statement
	expression  Expression
	expressions []Expression
	identifier  Identifier
	text        String
	integer     Integer
	float       Float
	ternary     Ternary
	null        Null
	token       Token
}

const IDENTIFIER = 57346
const STRING = 57347
const INTEGER = 57348
const FLOAT = 57349
const BOOLEAN = 57350
const TERNARY = 57351
const DATETIME = 57352
const SELECT = 57353
const FROM = 57354
const UPDATE = 57355
const SET = 57356
const DELETE = 57357
const WHERE = 57358
const INSERT = 57359
const INTO = 57360
const VALUES = 57361
const AS = 57362
const DUAL = 57363
const CREATE = 57364
const DROP = 57365
const ALTER = 57366
const TABLE = 57367
const COLUMN = 57368
const ORDER = 57369
const GROUP = 57370
const HAVING = 57371
const BY = 57372
const ASC = 57373
const DESC = 57374
const LIMIT = 57375
const JOIN = 57376
const INNER = 57377
const OUTER = 57378
const LEFT = 57379
const RIGHT = 57380
const FULL = 57381
const CROSS = 57382
const ON = 57383
const USING = 57384
const NATURAL = 57385
const UNION = 57386
const ALL = 57387
const ANY = 57388
const EXISTS = 57389
const IN = 57390
const AND = 57391
const OR = 57392
const NOT = 57393
const BETWEEN = 57394
const LIKE = 57395
const IS = 57396
const NULL = 57397
const DISTINCT = 57398
const WITH = 57399
const TRUE = 57400
const FALSE = 57401
const UNKNOWN = 57402
const CASE = 57403
const WHEN = 57404
const THEN = 57405
const ELSE = 57406
const END = 57407
const COMPARISON_OP = 57408
const STRING_OP = 57409

var yyToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"IDENTIFIER",
	"STRING",
	"INTEGER",
	"FLOAT",
	"BOOLEAN",
	"TERNARY",
	"DATETIME",
	"SELECT",
	"FROM",
	"UPDATE",
	"SET",
	"DELETE",
	"WHERE",
	"INSERT",
	"INTO",
	"VALUES",
	"AS",
	"DUAL",
	"CREATE",
	"DROP",
	"ALTER",
	"TABLE",
	"COLUMN",
	"ORDER",
	"GROUP",
	"HAVING",
	"BY",
	"ASC",
	"DESC",
	"LIMIT",
	"JOIN",
	"INNER",
	"OUTER",
	"LEFT",
	"RIGHT",
	"FULL",
	"CROSS",
	"ON",
	"USING",
	"NATURAL",
	"UNION",
	"ALL",
	"ANY",
	"EXISTS",
	"IN",
	"AND",
	"OR",
	"NOT",
	"BETWEEN",
	"LIKE",
	"IS",
	"NULL",
	"DISTINCT",
	"WITH",
	"TRUE",
	"FALSE",
	"UNKNOWN",
	"CASE",
	"WHEN",
	"THEN",
	"ELSE",
	"END",
	"COMPARISON_OP",
	"STRING_OP",
	"'+'",
	"'-'",
	"'*'",
	"'/'",
	"'('",
	"')'",
	"'%'",
	"','",
	"';'",
}
var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

//line parser.y:734

func SetDebugLevel(level int, verbose bool) {
	yyDebug = level
	yyErrorVerbose = verbose
}

func Parse(s string) ([]Statement, error) {
	l := new(Lexer)
	l.Init(s)
	yyParse(l)
	return l.program, l.err
}

//line yacctab:1
var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
	-1, 18,
	34, 119,
	36, 123,
	-2, 99,
	-1, 29,
	48, 117,
	52, 117,
	53, 117,
	-2, 61,
	-1, 59,
	36, 123,
	-2, 119,
	-1, 86,
	49, 61,
	50, 61,
	-2, 117,
	-1, 90,
	73, 66,
	-2, 115,
}

const yyPrivate = 57344

const yyLast = 384

var yyAct = [...]int{

	41, 197, 37, 29, 179, 193, 43, 165, 155, 136,
	47, 18, 19, 102, 36, 158, 9, 12, 45, 201,
	67, 69, 85, 192, 176, 78, 20, 27, 76, 80,
	81, 82, 83, 79, 164, 84, 86, 77, 76, 80,
	81, 82, 83, 57, 126, 84, 76, 80, 81, 82,
	83, 72, 92, 84, 88, 82, 83, 84, 94, 84,
	87, 200, 189, 19, 80, 81, 82, 83, 89, 106,
	84, 107, 75, 74, 110, 85, 153, 20, 78, 108,
	113, 114, 194, 195, 121, 122, 123, 124, 125, 22,
	77, 76, 80, 81, 82, 83, 127, 132, 84, 138,
	19, 139, 111, 112, 22, 19, 16, 130, 129, 185,
	90, 140, 117, 150, 20, 142, 144, 24, 23, 20,
	174, 131, 149, 90, 152, 135, 50, 52, 148, 50,
	145, 146, 13, 156, 151, 160, 147, 131, 17, 159,
	60, 53, 19, 63, 19, 64, 65, 66, 61, 85,
	75, 59, 103, 168, 156, 170, 20, 23, 20, 172,
	75, 74, 169, 175, 181, 138, 182, 139, 58, 177,
	19, 143, 23, 141, 183, 188, 104, 156, 187, 44,
	191, 186, 190, 62, 20, 94, 198, 54, 196, 51,
	109, 99, 178, 181, 184, 182, 162, 199, 194, 195,
	101, 105, 198, 202, 22, 49, 50, 52, 163, 53,
	46, 22, 49, 50, 52, 98, 53, 46, 6, 97,
	22, 49, 50, 52, 56, 53, 46, 22, 100, 134,
	73, 95, 22, 49, 50, 52, 15, 53, 46, 22,
	49, 50, 52, 70, 53, 46, 6, 39, 22, 11,
	6, 40, 119, 22, 39, 54, 118, 120, 40, 8,
	42, 48, 54, 39, 68, 157, 91, 40, 48, 51,
	28, 54, 32, 116, 115, 38, 51, 48, 63, 32,
	64, 65, 66, 54, 1, 51, 25, 7, 32, 48,
	54, 76, 80, 81, 82, 83, 48, 51, 84, 173,
	93, 26, 21, 128, 51, 171, 31, 93, 22, 49,
	50, 52, 34, 53, 46, 22, 49, 50, 52, 30,
	53, 46, 35, 76, 80, 81, 82, 83, 180, 137,
	84, 76, 80, 81, 82, 83, 33, 126, 84, 161,
	63, 133, 64, 65, 66, 61, 166, 167, 59, 63,
	96, 64, 65, 66, 61, 71, 4, 59, 4, 54,
	55, 14, 10, 5, 3, 48, 54, 2, 0, 0,
	0, 0, 48, 51, 154, 0, 93, 0, 0, 0,
	51, 0, 0, 93,
}
var yyPact = [...]int{

	239, -1000, 239, -60, -1000, 237, 76, -1000, -1000, -1000,
	220, 85, 200, -1000, 196, 216, -1000, -1000, 108, 244,
	223, -1000, -1000, 239, -1000, -24, 210, 111, -1000, 24,
	-1000, -1000, 207, -1000, -1000, -1000, -1000, -1000, -1000, 46,
	216, 38, -1000, -1000, -1000, -1000, -1000, -1000, 311, -1000,
	-1000, 120, -1000, -1000, -1000, 190, 185, 111, 157, 243,
	116, 142, 100, -1000, -1000, -1000, -1000, -1000, 249, -1000,
	249, 6, 200, 249, 216, 216, 311, 228, 98, 204,
	311, 311, 311, 311, 311, -1000, -29, 23, -1000, -1000,
	76, 59, -21, 235, -1000, -1000, 202, 216, 249, 100,
	139, 116, 137, -1000, 100, -1000, -1000, -1000, -1000, -1000,
	-1000, 101, -1000, -4, -21, 46, 46, 132, 311, 41,
	311, -15, -15, -17, -17, -21, -1000, -1000, 3, 304,
	75, 311, 264, 163, 178, 111, -1000, -41, 38, -1000,
	305, 100, 128, 100, 314, -1000, -1000, -1000, -1000, 256,
	235, -1000, -21, -1000, -1000, -1000, 224, 55, 59, 311,
	-39, -1000, 123, 249, 249, -1000, 216, 37, 314, 100,
	305, 311, -11, 311, -1000, -21, 311, -1000, 123, -1000,
	-52, 51, 167, -1000, 111, 249, 314, -1000, -21, -1000,
	-1000, -21, 249, -1000, -1000, -1000, -1000, -12, -56, -1000,
	-1000, 249, -1000,
}
var yyPgo = [...]int{

	0, 284, 367, 364, 355, 363, 362, 361, 360, 350,
	341, 339, 336, 3, 329, 328, 14, 322, 319, 312,
	306, 27, 2, 303, 11, 302, 7, 301, 286, 275,
	266, 265, 8, 9, 4, 138, 1, 117, 15, 0,
	260, 6, 179, 18, 10, 17, 33, 5, 168, 13,
	140, 259,
}
var yyR1 = [...]int{

	0, 1, 1, 2, 3, 4, 5, 6, 6, 6,
	7, 7, 8, 8, 9, 9, 10, 10, 11, 11,
	12, 12, 12, 12, 12, 12, 12, 13, 13, 13,
	13, 13, 13, 13, 14, 14, 15, 15, 47, 47,
	47, 16, 17, 18, 18, 18, 18, 18, 18, 18,
	18, 18, 18, 19, 19, 19, 19, 19, 20, 20,
	20, 21, 21, 21, 21, 22, 23, 23, 23, 24,
	24, 24, 24, 24, 24, 24, 25, 25, 25, 25,
	25, 26, 26, 26, 27, 27, 28, 28, 29, 30,
	30, 31, 31, 32, 32, 33, 33, 34, 34, 35,
	35, 36, 36, 37, 37, 38, 38, 39, 40, 41,
	41, 42, 42, 43, 44, 45, 45, 46, 46, 48,
	48, 49, 49, 50, 50, 50, 50, 51, 51,
}
var yyR2 = [...]int{

	0, 0, 2, 2, 1, 7, 3, 0, 2, 2,
	0, 2, 0, 3, 0, 2, 0, 3, 0, 2,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 3, 1, 1, 2, 2, 0, 1,
	1, 3, 3, 3, 4, 4, 6, 6, 4, 4,
	4, 4, 2, 3, 3, 3, 3, 3, 3, 3,
	2, 1, 1, 1, 3, 4, 0, 2, 2, 1,
	2, 3, 1, 2, 3, 1, 5, 5, 6, 6,
	4, 0, 2, 4, 1, 1, 1, 3, 5, 0,
	1, 0, 2, 1, 3, 1, 3, 1, 3, 1,
	3, 1, 3, 1, 3, 4, 2, 1, 1, 1,
	2, 1, 2, 1, 1, 0, 1, 0, 1, 0,
	1, 0, 1, 0, 1, 1, 1, 0, 1,
}
var yyChk = [...]int{

	-1000, -1, -2, -3, -4, -5, 11, -1, -51, 76,
	-6, 12, -45, 56, -7, 16, 21, -35, -24, -39,
	-16, -25, 4, 72, -37, -28, -27, -21, 70, -13,
	-18, -20, 72, -12, -19, -17, -16, -22, -29, 47,
	51, -39, -40, -41, -42, -43, 10, -44, 61, 5,
	6, 69, 7, 9, 55, -8, 28, -21, -48, 43,
	-50, 40, 75, 35, 37, 38, 39, -39, 20, -39,
	20, -4, 75, 20, 50, 49, 67, 66, 54, -46,
	68, 69, 70, 71, 74, 51, -13, -21, -16, -21,
	72, -30, -13, 72, -41, -42, -9, 29, 30, 34,
	-48, -50, -49, 36, 34, -35, -39, -39, 73, -37,
	-39, -21, -21, -13, -13, 46, 45, -46, 52, 48,
	53, -13, -13, -13, -13, -13, 73, 73, -23, -45,
	-38, 62, -13, -10, 27, -21, -33, -14, -39, -22,
	-24, 34, -49, 34, -24, -16, -16, -43, -44, -13,
	72, -16, -13, 73, 70, -32, -13, -31, -38, 64,
	-13, -11, 33, 30, 75, -26, 41, 42, -24, 34,
	-24, 49, -32, 75, 65, -13, 63, -41, 69, -34,
	-15, -39, -22, -33, -21, 72, -24, -26, -13, 73,
	-32, -13, 75, -47, 31, 32, -47, -36, -39, -34,
	73, 75, -36,
}
var yyDef = [...]int{

	1, -2, 1, 127, 4, 7, 115, 2, 3, 128,
	10, 0, 0, 116, 12, 0, 8, 9, -2, 69,
	72, 75, 107, 0, 6, 103, 86, 84, 85, -2,
	62, 63, 0, 27, 28, 29, 30, 31, 32, 0,
	0, 20, 21, 22, 23, 24, 25, 26, 89, 108,
	109, 0, 111, 113, 114, 14, 0, 11, 0, -2,
	121, 0, 0, 120, 124, 125, 126, 70, 0, 73,
	0, 0, 0, 0, 0, 0, 0, 0, 117, 0,
	0, 0, 0, 0, 0, 118, -2, 0, 52, 60,
	-2, 0, 90, 0, 110, 112, 16, 0, 0, 0,
	0, 121, 0, 122, 0, 100, 71, 74, 41, 104,
	87, 58, 59, 42, 43, 0, 0, 0, 0, 0,
	0, 53, 54, 55, 56, 57, 33, 64, 0, 0,
	91, 0, 0, 18, 0, 15, 13, 95, 34, 35,
	81, 0, 0, 0, 80, 50, 51, 44, 45, 0,
	0, 48, 49, 65, 67, 68, 93, 0, 106, 0,
	0, 5, 0, 0, 0, 76, 0, 0, 77, 0,
	81, 0, 0, 0, 88, 92, 0, 19, 0, 17,
	97, 38, 38, 96, 82, 0, 79, 78, 46, 47,
	94, 105, 0, 36, 39, 40, 37, 0, 101, 98,
	83, 0, 102,
}
var yyTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 74, 3, 3,
	72, 73, 70, 68, 75, 69, 3, 71, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 76,
}
var yyTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, 34, 35, 36, 37, 38, 39, 40, 41,
	42, 43, 44, 45, 46, 47, 48, 49, 50, 51,
	52, 53, 54, 55, 56, 57, 58, 59, 60, 61,
	62, 63, 64, 65, 66, 67,
}
var yyTok3 = [...]int{
	0,
}

var yyErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

//line yaccpar:1

/*	parser for yacc output	*/

var (
	yyDebug        = 0
	yyErrorVerbose = false
)

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

type yyParser interface {
	Parse(yyLexer) int
	Lookahead() int
}

type yyParserImpl struct {
	lval  yySymType
	stack [yyInitialStackSize]yySymType
	char  int
}

func (p *yyParserImpl) Lookahead() int {
	return p.char
}

func yyNewParser() yyParser {
	return &yyParserImpl{}
}

const yyFlag = -1000

func yyTokname(c int) string {
	if c >= 1 && c-1 < len(yyToknames) {
		if yyToknames[c-1] != "" {
			return yyToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func yyStatname(s int) string {
	if s >= 0 && s < len(yyStatenames) {
		if yyStatenames[s] != "" {
			return yyStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func yyErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !yyErrorVerbose {
		return "syntax error"
	}

	for _, e := range yyErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + yyTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := yyPact[state]
	for tok := TOKSTART; tok-1 < len(yyToknames); tok++ {
		if n := base + tok; n >= 0 && n < yyLast && yyChk[yyAct[n]] == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if yyDef[state] == -2 {
		i := 0
		for yyExca[i] != -1 || yyExca[i+1] != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; yyExca[i] >= 0; i += 2 {
			tok := yyExca[i]
			if tok < TOKSTART || yyExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if yyExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += yyTokname(tok)
	}
	return res
}

func yylex1(lex yyLexer, lval *yySymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = yyTok1[0]
		goto out
	}
	if char < len(yyTok1) {
		token = yyTok1[char]
		goto out
	}
	if char >= yyPrivate {
		if char < yyPrivate+len(yyTok2) {
			token = yyTok2[char-yyPrivate]
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		token = yyTok3[i+0]
		if token == char {
			token = yyTok3[i+1]
			goto out
		}
	}

out:
	if token == 0 {
		token = yyTok2[1] /* unknown char */
	}
	if yyDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", yyTokname(token), uint(char))
	}
	return char, token
}

func yyParse(yylex yyLexer) int {
	return yyNewParser().Parse(yylex)
}

func (yyrcvr *yyParserImpl) Parse(yylex yyLexer) int {
	var yyn int
	var yyVAL yySymType
	var yyDollar []yySymType
	_ = yyDollar // silence set and not used
	yyS := yyrcvr.stack[:]

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yyrcvr.char = -1
	yytoken := -1 // yyrcvr.char translated into internal numbering
	defer func() {
		// Make sure we report no lookahead when not parsing.
		yystate = -1
		yyrcvr.char = -1
		yytoken = -1
	}()
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	if yyDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", yyTokname(yytoken), yyStatname(yystate))
	}

	yyp++
	if yyp >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyS[yyp] = yyVAL
	yyS[yyp].yys = yystate

yynewstate:
	yyn = yyPact[yystate]
	if yyn <= yyFlag {
		goto yydefault /* simple state */
	}
	if yyrcvr.char < 0 {
		yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
	}
	yyn += yytoken
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = yyAct[yyn]
	if yyChk[yyn] == yytoken { /* valid shift */
		yyrcvr.char = -1
		yytoken = -1
		yyVAL = yyrcvr.lval
		yystate = yyn
		if Errflag > 0 {
			Errflag--
		}
		goto yystack
	}

yydefault:
	/* default state action */
	yyn = yyDef[yystate]
	if yyn == -2 {
		if yyrcvr.char < 0 {
			yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
		}

		/* look through exception table */
		xi := 0
		for {
			if yyExca[xi+0] == -1 && yyExca[xi+1] == yystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			yyn = yyExca[xi+0]
			if yyn < 0 || yyn == yytoken {
				break
			}
		}
		yyn = yyExca[xi+1]
		if yyn < 0 {
			goto ret0
		}
	}
	if yyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			yylex.Error(yyErrorMessage(yystate, yytoken))
			Nerrs++
			if yyDebug >= 1 {
				__yyfmt__.Printf("%s", yyStatname(yystate))
				__yyfmt__.Printf(" saw %s\n", yyTokname(yytoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for yyp >= 0 {
				yyn = yyPact[yyS[yyp].yys] + yyErrCode
				if yyn >= 0 && yyn < yyLast {
					yystate = yyAct[yyn] /* simulate a shift of "error" */
					if yyChk[yystate] == yyErrCode {
						goto yystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if yyDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", yyS[yyp].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", yyTokname(yytoken))
			}
			if yytoken == yyEofCode {
				goto ret1
			}
			yyrcvr.char = -1
			yytoken = -1
			goto yynewstate /* try again in the same state */
		}
	}

	/* reduction by production yyn */
	if yyDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", yyn, yyStatname(yystate))
	}

	yynt := yyn
	yypt := yyp
	_ = yypt // guard against "declared and not used"

	yyp -= yyR2[yyn]
	// yyp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if yyp+1 >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyVAL = yyS[yyp+1]

	/* consult goto table to find next state */
	yyn = yyR1[yyn]
	yyg := yyPgo[yyn]
	yyj := yyg + yyS[yyp].yys + 1

	if yyj >= yyLast {
		yystate = yyAct[yyg]
	} else {
		yystate = yyAct[yyj]
		if yyChk[yystate] != -yyn {
			yystate = yyAct[yyg]
		}
	}
	// dummy call; replaced with literal code
	switch yynt {

	case 1:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:93
		{
			yyVAL.program = nil
			yylex.(*Lexer).program = yyVAL.program
		}
	case 2:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:98
		{
			yyVAL.program = append([]Statement{yyDollar[1].statement}, yyDollar[2].program...)
			yylex.(*Lexer).program = yyVAL.program
		}
	case 3:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:105
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 4:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:111
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 5:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:117
		{
			yyVAL.expression = SelectQuery{
				SelectClause:  yyDollar[1].expression,
				FromClause:    yyDollar[2].expression,
				WhereClause:   yyDollar[3].expression,
				GroupByClause: yyDollar[4].expression,
				HavingClause:  yyDollar[5].expression,
				OrderByClause: yyDollar[6].expression,
				LimitClause:   yyDollar[7].expression,
			}
		}
	case 6:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:131
		{
			yyVAL.expression = SelectClause{Select: yyDollar[1].token.Literal, Distinct: yyDollar[2].token, Fields: yyDollar[3].expressions}
		}
	case 7:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:137
		{
			yyVAL.expression = nil
		}
	case 8:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:141
		{
			yyVAL.expression = FromClause{From: yyDollar[1].token.Literal, Tables: []Expression{Dual{Dual: yyDollar[2].token.Literal}}}
		}
	case 9:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:145
		{
			yyVAL.expression = FromClause{From: yyDollar[1].token.Literal, Tables: yyDollar[2].expressions}
		}
	case 10:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:151
		{
			yyVAL.expression = nil
		}
	case 11:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:155
		{
			yyVAL.expression = WhereClause{Where: yyDollar[1].token.Literal, Filter: yyDollar[2].expression}
		}
	case 12:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:161
		{
			yyVAL.expression = nil
		}
	case 13:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:165
		{
			yyVAL.expression = GroupByClause{GroupBy: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Items: yyDollar[3].expressions}
		}
	case 14:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:171
		{
			yyVAL.expression = nil
		}
	case 15:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:175
		{
			yyVAL.expression = HavingClause{Having: yyDollar[1].token.Literal, Filter: yyDollar[2].expression}
		}
	case 16:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:181
		{
			yyVAL.expression = nil
		}
	case 17:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:185
		{
			yyVAL.expression = OrderByClause{OrderBy: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Items: yyDollar[3].expressions}
		}
	case 18:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:191
		{
			yyVAL.expression = nil
		}
	case 19:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:195
		{
			yyVAL.expression = LimitClause{Limit: yyDollar[1].token.Literal, Number: yyDollar[2].integer.Value()}
		}
	case 20:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:201
		{
			yyVAL.expression = yyDollar[1].identifier
		}
	case 21:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:205
		{
			yyVAL.expression = yyDollar[1].text
		}
	case 22:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:209
		{
			yyVAL.expression = yyDollar[1].integer
		}
	case 23:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:213
		{
			yyVAL.expression = yyDollar[1].float
		}
	case 24:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:217
		{
			yyVAL.expression = yyDollar[1].ternary
		}
	case 25:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:221
		{
			yyVAL.expression = NewDatetimeFromString(yyDollar[1].token.Literal)
		}
	case 26:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:225
		{
			yyVAL.expression = yyDollar[1].null
		}
	case 27:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:231
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 28:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:235
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 29:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:239
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 30:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:243
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 31:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:247
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 32:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:251
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 33:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:255
		{
			yyVAL.expression = Parentheses{Expr: yyDollar[2].expression}
		}
	case 34:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:261
		{
			yyVAL.expression = yyDollar[1].identifier
		}
	case 35:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:265
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 36:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:271
		{
			yyVAL.expression = OrderItem{Item: yyDollar[1].identifier, Direction: yyDollar[2].token}
		}
	case 37:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:275
		{
			yyVAL.expression = OrderItem{Item: yyDollar[1].expression, Direction: yyDollar[2].token}
		}
	case 38:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:281
		{
			yyVAL.token = Token{}
		}
	case 39:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:285
		{
			yyVAL.token = yyDollar[1].token
		}
	case 40:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:289
		{
			yyVAL.token = yyDollar[1].token
		}
	case 41:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:295
		{
			yyVAL.expression = Subquery{Query: yyDollar[2].expression.(SelectQuery)}
		}
	case 42:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:301
		{
			var item1 []Expression
			var item2 []Expression

			c1, ok := yyDollar[1].expression.(Concat)
			if ok {
				item1 = c1.Items
			} else {
				item1 = []Expression{yyDollar[1].expression}
			}

			c2, ok := yyDollar[3].expression.(Concat)
			if ok {
				item2 = c2.Items
			} else {
				item2 = []Expression{yyDollar[3].expression}
			}

			yyVAL.expression = Concat{Items: append(item1, item2...)}
		}
	case 43:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:324
		{
			yyVAL.expression = Comparison{LHS: yyDollar[1].expression, Operator: yyDollar[2].token, RHS: yyDollar[3].expression}
		}
	case 44:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:328
		{
			yyVAL.expression = Is{Is: yyDollar[2].token.Literal, LHS: yyDollar[1].expression, RHS: yyDollar[4].ternary, Negation: yyDollar[3].token}
		}
	case 45:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:332
		{
			yyVAL.expression = Is{Is: yyDollar[2].token.Literal, LHS: yyDollar[1].expression, RHS: yyDollar[4].null, Negation: yyDollar[3].token}
		}
	case 46:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:336
		{
			yyVAL.expression = Between{Between: yyDollar[3].token.Literal, And: yyDollar[5].token.Literal, LHS: yyDollar[1].expression, Low: yyDollar[4].expression, High: yyDollar[6].expression, Negation: yyDollar[2].token}
		}
	case 47:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:340
		{
			yyVAL.expression = In{In: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, List: yyDollar[5].expressions, Negation: yyDollar[2].token}
		}
	case 48:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:344
		{
			yyVAL.expression = In{In: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Query: yyDollar[4].expression.(Subquery), Negation: yyDollar[2].token}
		}
	case 49:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:348
		{
			yyVAL.expression = Like{Like: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Pattern: yyDollar[4].expression, Negation: yyDollar[2].token}
		}
	case 50:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:352
		{
			yyVAL.expression = Any{Any: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token, Query: yyDollar[4].expression.(Subquery)}
		}
	case 51:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:356
		{
			yyVAL.expression = All{All: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token, Query: yyDollar[4].expression.(Subquery)}
		}
	case 52:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:360
		{
			yyVAL.expression = Exists{Exists: yyDollar[1].token.Literal, Query: yyDollar[2].expression.(Subquery)}
		}
	case 53:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:366
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('+'), RHS: yyDollar[3].expression}
		}
	case 54:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:370
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('-'), RHS: yyDollar[3].expression}
		}
	case 55:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:374
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('*'), RHS: yyDollar[3].expression}
		}
	case 56:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:378
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('/'), RHS: yyDollar[3].expression}
		}
	case 57:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:382
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('%'), RHS: yyDollar[3].expression}
		}
	case 58:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:388
		{
			yyVAL.expression = Logic{LHS: yyDollar[1].expression, Operator: yyDollar[2].token, RHS: yyDollar[3].expression}
		}
	case 59:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:392
		{
			yyVAL.expression = Logic{LHS: yyDollar[1].expression, Operator: yyDollar[2].token, RHS: yyDollar[3].expression}
		}
	case 60:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:396
		{
			yyVAL.expression = Logic{LHS: nil, Operator: yyDollar[1].token, RHS: yyDollar[2].expression}
		}
	case 61:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:402
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 62:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:406
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 63:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:410
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 64:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:414
		{
			yyVAL.expression = Parentheses{Expr: yyDollar[2].expression}
		}
	case 65:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:420
		{
			yyVAL.expression = Function{Name: yyDollar[1].identifier.Literal, Option: yyDollar[3].expression.(Option)}
		}
	case 66:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:426
		{
			yyVAL.expression = Option{}
		}
	case 67:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:430
		{
			yyVAL.expression = Option{Distinct: yyDollar[1].token, Args: []Expression{AllColumns{}}}
		}
	case 68:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:434
		{
			yyVAL.expression = Option{Distinct: yyDollar[1].token, Args: yyDollar[2].expressions}
		}
	case 69:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:440
		{
			yyVAL.expression = Table{Object: yyDollar[1].identifier}
		}
	case 70:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:444
		{
			yyVAL.expression = Table{Object: yyDollar[1].identifier, Alias: yyDollar[2].identifier}
		}
	case 71:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:448
		{
			yyVAL.expression = Table{Object: yyDollar[1].identifier, As: yyDollar[2].token, Alias: yyDollar[3].identifier}
		}
	case 72:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:452
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression}
		}
	case 73:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:456
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression, Alias: yyDollar[2].identifier}
		}
	case 74:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:460
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression, As: yyDollar[2].token, Alias: yyDollar[3].identifier}
		}
	case 75:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:464
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression}
		}
	case 76:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:470
		{
			yyVAL.expression = Join{Join: yyDollar[3].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[4].expression.(Table), Natural: Token{}, JoinType: yyDollar[2].token, Condition: yyDollar[5].expression}
		}
	case 77:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:474
		{
			yyVAL.expression = Join{Join: yyDollar[4].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[5].expression.(Table), Natural: yyDollar[2].token, JoinType: yyDollar[3].token, Condition: nil}
		}
	case 78:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:478
		{
			yyVAL.expression = Join{Join: yyDollar[4].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[5].expression.(Table), Natural: Token{}, JoinType: yyDollar[3].token, Direction: yyDollar[2].token, Condition: yyDollar[6].expression}
		}
	case 79:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:482
		{
			yyVAL.expression = Join{Join: yyDollar[5].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[6].expression.(Table), Natural: yyDollar[2].token, JoinType: yyDollar[4].token, Direction: yyDollar[3].token, Condition: nil}
		}
	case 80:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:486
		{
			yyVAL.expression = Join{Join: yyDollar[3].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[4].expression.(Table), Natural: Token{}, JoinType: yyDollar[2].token, Condition: nil}
		}
	case 81:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:492
		{
			yyVAL.expression = nil
		}
	case 82:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:496
		{
			yyVAL.expression = JoinCondition{Literal: yyDollar[1].token.Literal, On: yyDollar[2].expression}
		}
	case 83:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:500
		{
			yyVAL.expression = JoinCondition{Literal: yyDollar[1].token.Literal, Using: yyDollar[3].expressions}
		}
	case 84:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:506
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 85:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:510
		{
			yyVAL.expression = AllColumns{}
		}
	case 86:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:516
		{
			yyVAL.expression = Field{Object: yyDollar[1].expression}
		}
	case 87:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:520
		{
			yyVAL.expression = Field{Object: yyDollar[1].expression, As: yyDollar[2].token, Alias: yyDollar[3].identifier}
		}
	case 88:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:526
		{
			yyVAL.expression = Case{Case: yyDollar[1].token.Literal, End: yyDollar[5].token.Literal, Value: yyDollar[2].expression, When: yyDollar[3].expressions, Else: yyDollar[4].expression}
		}
	case 89:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:532
		{
			yyVAL.expression = nil
		}
	case 90:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:536
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 91:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:542
		{
			yyVAL.expression = nil
		}
	case 92:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:546
		{
			yyVAL.expression = CaseElse{Else: yyDollar[1].token.Literal, Result: yyDollar[2].expression}
		}
	case 93:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:552
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 94:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:556
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 95:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:562
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 96:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:566
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 97:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:572
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 98:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:576
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 99:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:582
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 100:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:586
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 101:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:592
		{
			yyVAL.expressions = []Expression{yyDollar[1].identifier}
		}
	case 102:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:596
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].identifier}, yyDollar[3].expressions...)
		}
	case 103:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:602
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 104:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:606
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 105:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:612
		{
			yyVAL.expressions = []Expression{CaseWhen{When: yyDollar[1].token.Literal, Then: yyDollar[3].token.Literal, Condition: yyDollar[2].expression, Result: yyDollar[4].expression}}
		}
	case 106:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:616
		{
			yyVAL.expressions = append(yyDollar[1].expressions, yyDollar[2].expressions...)
		}
	case 107:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:622
		{
			yyVAL.identifier = Identifier{Literal: yyDollar[1].token.Literal, Quoted: yyDollar[1].token.Quoted}
		}
	case 108:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:628
		{
			yyVAL.text = NewString(yyDollar[1].token.Literal)
		}
	case 109:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:634
		{
			yyVAL.integer = NewIntegerFromString(yyDollar[1].token.Literal)
		}
	case 110:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:638
		{
			i := yyDollar[2].integer.Value() * -1
			yyVAL.integer = NewInteger(i)
		}
	case 111:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:645
		{
			yyVAL.float = NewFloatFromString(yyDollar[1].token.Literal)
		}
	case 112:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:649
		{
			f := yyDollar[2].float.Value() * -1
			yyVAL.float = NewFloat(f)
		}
	case 113:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:656
		{
			yyVAL.ternary = NewTernaryFromString(yyDollar[1].token.Literal)
		}
	case 114:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:662
		{
			yyVAL.null = NewNullFromString(yyDollar[1].token.Literal)
		}
	case 115:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:668
		{
			yyVAL.token = Token{}
		}
	case 116:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:672
		{
			yyVAL.token = yyDollar[1].token
		}
	case 117:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:678
		{
			yyVAL.token = Token{}
		}
	case 118:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:682
		{
			yyVAL.token = yyDollar[1].token
		}
	case 119:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:688
		{
			yyVAL.token = Token{}
		}
	case 120:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:692
		{
			yyVAL.token = yyDollar[1].token
		}
	case 121:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:698
		{
			yyVAL.token = Token{}
		}
	case 122:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:702
		{
			yyVAL.token = yyDollar[1].token
		}
	case 123:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:708
		{
			yyVAL.token = Token{}
		}
	case 124:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:712
		{
			yyVAL.token = yyDollar[1].token
		}
	case 125:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:716
		{
			yyVAL.token = yyDollar[1].token
		}
	case 126:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:720
		{
			yyVAL.token = yyDollar[1].token
		}
	case 127:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:726
		{
			yyVAL.token = Token{}
		}
	case 128:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:730
		{
			yyVAL.token = Token{Token: ';', Literal: string(';')}
		}
	}
	goto yystack /* stack new state and value */
}
