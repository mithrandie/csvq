package parser

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/ternary"
)

const (
	EOF = -(iota + 1)
	UNCATEGORIZED
)

var tokenLiterals = map[int]string{
	IDENTIFIER:      "IDENTIFIER",
	STRING:          "STRING",
	INTEGER:         "INTEGER",
	FLOAT:           "FLOAT",
	BOOLEAN:         "BOOLEAN",
	TERNARY:         "TERNARY",
	DATETIME:        "DATETIME",
	VARIABLE:        "VARIABLE",
	SELECT:          "SELECT",
	FROM:            "FROM",
	UPDATE:          "UPDATE",
	SET:             "SET",
	DELETE:          "DELETE",
	WHERE:           "WHERE",
	INSERT:          "INSERT",
	INTO:            "INTO",
	VALUES:          "VLLUES",
	AS:              "AS",
	DUAL:            "DUAL",
	STDIN:           "STDIN",
	CREATE:          "CREATE",
	DROP:            "DROP",
	ALTER:           "ALTER",
	TABLE:           "TABLE",
	COLUMN:          "COLUMN",
	ORDER:           "ORDER",
	GROUP:           "GROUP",
	HAVING:          "HAVING",
	BY:              "BY",
	ASC:             "ASC",
	DESC:            "DESC",
	LIMIT:           "LIMIT",
	JOIN:            "JOIN",
	INNER:           "INNER",
	OUTER:           "OUTER",
	LEFT:            "LEFT",
	RIGHT:           "RIGHT",
	FULL:            "FULL",
	CROSS:           "CROSS",
	ON:              "ON",
	USING:           "USING",
	NATURAL:         "NATURAL",
	UNION:           "UNION",
	ALL:             "ALL",
	ANY:             "ANY",
	EXISTS:          "EXISTS",
	IN:              "IN",
	AND:             "AND",
	OR:              "OR",
	NOT:             "NOT",
	BETWEEN:         "BETWEEN",
	LIKE:            "LIKE",
	IS:              "IS",
	NULL:            "NULL",
	DISTINCT:        "DISTINCT",
	WITH:            "WITH",
	TRUE:            "TRUE",
	FALSE:           "FALSE",
	UNKNOWN:         "UNKNOWN",
	CASE:            "CASE",
	WHEN:            "WHEN",
	THEN:            "THEN",
	ELSE:            "ELSE",
	END:             "END",
	GROUP_CONCAT:    "GROUP_CONCAT",
	SEPARATOR:       "SEPARATOR",
	VAR:             "VAR",
	COMPARISON_OP:   "COMPARISON_OP",
	STRING_OP:       "STRING_OP",
	SUBSTITUTION_OP: "SUBSTITUTION_OP",
}

var keywords = []int{
	SELECT,
	FROM,
	UPDATE,
	SET,
	DELETE,
	WHERE,
	INSERT,
	INTO,
	VALUES,
	AS,
	DUAL,
	STDIN,
	CREATE,
	DROP,
	ALTER,
	TABLE,
	COLUMN,
	ORDER,
	GROUP,
	HAVING,
	BY,
	ASC,
	DESC,
	LIMIT,
	JOIN,
	INNER,
	OUTER,
	LEFT,
	RIGHT,
	FULL,
	CROSS,
	ON,
	USING,
	NATURAL,
	UNION,
	ALL,
	ANY,
	EXISTS,
	IN,
	AND,
	OR,
	NOT,
	BETWEEN,
	LIKE,
	IS,
	NULL,
	DISTINCT,
	WITH,
	CASE,
	WHEN,
	THEN,
	ELSE,
	END,
	GROUP_CONCAT,
	SEPARATOR,
	VAR,
}

var comparisonOperators = []string{
	">",
	"<",
	">=",
	"<=",
	"<>",
	"!=",
}

var stringOperators = []string{
	"||",
}

var substitutionOperators = []string{
	":=",
}

func TokenLiteral(token int) string {
	if val, ok := tokenLiterals[token]; ok {
		return val
	}
	return string(token)
}

type Scanner struct {
	src    []rune
	srcPos int
	offset int
	err    error
}

func (s *Scanner) Init(src string) *Scanner {
	s.src = []rune(src)
	s.srcPos = 0
	s.offset = 0
	s.err = nil
	return s
}

func (s *Scanner) peek() rune {
	if len(s.src) <= s.srcPos {
		return EOF
	}
	return s.src[s.srcPos]
}

func (s *Scanner) next() rune {
	ch := s.peek()
	if ch != EOF {
		s.srcPos++
		s.offset++
	}
	return ch
}

func (s *Scanner) runes() []rune {
	return s.src[(s.srcPos - s.offset):s.srcPos]
}

func (s *Scanner) literal() string {
	return string(s.runes())
}

func (s *Scanner) unescapeTokenString() string {
	runes := s.runes()
	quote := runes[0]
	switch quote {
	case '"', '\'', '`':
		if runes[len(runes)-1] == quote {
			runes = runes[1:(len(runes) - 1)]
		} else {
			runes = runes[1:]
		}

		escaped := []rune{}
		for i := 0; i < len(runes); i++ {
			if runes[i] == '\\' && (i+1) < len(runes) && runes[i+1] == quote {
				i++
			}
			escaped = append(escaped, runes[i])
		}
		runes = escaped
	}
	return string(runes)
}

func (s *Scanner) Scan() (int, string, bool, error) {
	ch := s.peek()

	for s.isWhiteSpace(ch) {
		s.next()
		ch = s.peek()
	}

	s.offset = 0
	ch = s.next()
	token := ch
	literal := string(ch)
	quoted := false

	switch {
	case s.isIdentRune(ch):
		s.scanIdentifier()

		literal = s.literal()
		if _, e := strconv.ParseInt(literal, 10, 64); e == nil {
			token = INTEGER
		} else if _, e := strconv.ParseFloat(literal, 64); e == nil {
			token = FLOAT
		} else if _, e := ternary.Parse(literal); e == nil {
			token = TERNARY
		} else if t, e := s.searchKeyword(literal); e == nil {
			token = rune(t)
		} else {
			token = IDENTIFIER
		}
	case s.isOperatorRune(ch):
		s.scanOperator()

		literal = s.literal()
		if e := s.searchComparisonOperators(literal); e == nil {
			token = COMPARISON_OP
		} else if e := s.searchStringOperators(literal); e == nil {
			token = STRING_OP
		} else if e := s.searchSubstitutionOperators(literal); e == nil {
			token = SUBSTITUTION_OP
		} else if 1 < len(literal) {
			token = UNCATEGORIZED
		}
	case s.isVariableRune(ch):
		s.scanIdentifier()
		literal = s.literal()
		token = VARIABLE
	default:
		switch ch {
		case EOF:
			break
		case '"', '\'':
			s.scanString(ch)
			literal = s.unescapeTokenString()
			if _, e := StrToTime(literal); e == nil {
				token = DATETIME
			} else {
				token = STRING
				literal = cmd.UnescapeString(literal)
			}
		case '`':
			s.scanString(ch)
			literal = s.unescapeTokenString()
			token = IDENTIFIER
			quoted = true
		}
	}

	return int(token), literal, quoted, s.err
}

func (s *Scanner) isWhiteSpace(ch rune) bool {
	switch ch {
	case '\t', '\n', '\r', ' ':
		return true
	}
	return false
}

func (s *Scanner) scanString(quote rune) {
	for {
		ch := s.next()
		if ch == EOF {
			s.err = errors.New("literal not terminated")
			break
		} else if ch == quote {
			break
		} else if ch == '\\' {
			s.next()
		}
	}

	return
}

func (s *Scanner) scanIdentifier() {
	for s.isIdentRune(s.peek()) {
		s.next()
	}
	return
}

func (s *Scanner) isIdentRune(ch rune) bool {
	return ch == '_' || ch == '$' || ch == '.' || unicode.IsLetter(ch) || unicode.IsDigit(ch)
}

func (s *Scanner) scanOperator() {
	for s.isOperatorRune(s.peek()) {
		s.next()
	}
	return
}

func (s *Scanner) isOperatorRune(ch rune) bool {
	switch ch {
	case '=', '>', '<', '!', '|', ':':
		return true
	}
	return false
}

func (s *Scanner) isVariableRune(ch rune) bool {
	if ch == '@' {
		return true
	}
	return false
}

func (s *Scanner) searchKeyword(str string) (int, error) {
	for _, v := range keywords {
		if strings.EqualFold(TokenLiteral(v), str) {
			return v, nil
		}
	}
	return IDENTIFIER, errors.New(fmt.Sprintf("%q is not a keyword", str))
}

func (s *Scanner) searchComparisonOperators(str string) error {
	for _, v := range comparisonOperators {
		if v == str {
			return nil
		}
	}
	return errors.New(fmt.Sprintf("%q is not a comparison operator", str))
}

func (s *Scanner) searchStringOperators(str string) error {
	for _, v := range stringOperators {
		if v == str {
			return nil
		}
	}
	return errors.New(fmt.Sprintf("%q is not a string operator", str))
}

func (s *Scanner) searchSubstitutionOperators(str string) error {
	for _, v := range substitutionOperators {
		if v == str {
			return nil
		}
	}
	return errors.New(fmt.Sprintf("%q is not a substitution operator", str))
}
