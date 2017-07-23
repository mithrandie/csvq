package parser

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/ternary"
)

const (
	EOF = -(iota + 1)
	UNCATEGORIZED
)

const (
	TOKEN_FROM   = IDENTIFIER
	TOKEN_TO     = SUBSTITUTION_OP
	KEYWORD_FROM = SELECT
	KEYWORD_TO   = VAR
)

const (
	SUBSTITUTION_OPERATOR = ":="

	VARIABLE_SIGN = '@'
)

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

func TokenLiteral(token int) string {
	if TOKEN_FROM <= token && token <= TOKEN_TO {
		return yyToknames[token-TOKEN_FROM+3]
	}
	return string(token)
}

type Scanner struct {
	src    []rune
	srcPos int
	offset int
	err    error

	line       int
	char       int
	sourceFile string
}

func (s *Scanner) Init(src string, sourceFile string) *Scanner {
	s.src = []rune(src)
	s.srcPos = 0
	s.offset = 0
	s.err = nil
	s.line = 1
	s.char = 0

	if 0 < len(sourceFile) {
		if abs, err := filepath.Abs(sourceFile); err == nil {
			s.sourceFile = abs
		}
	}
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
	if ch == EOF {
		return ch
	}

	s.srcPos++
	s.offset++
	s.char++

	ch = s.checkNewLine(ch)

	return ch
}

func (s *Scanner) checkNewLine(ch rune) rune {
	if ch != '\r' && ch != '\n' {
		return ch
	}

	if ch == '\r' && s.peek() == '\n' {
		s.srcPos++
		s.offset++
	}

	s.line++
	s.char = 0
	return s.src[s.srcPos-1]
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

func (s *Scanner) Scan() (Token, error) {
	ch := s.peek()

	for unicode.IsSpace(ch) {
		s.next()
		ch = s.peek()
	}

	s.offset = 0
	ch = s.next()
	token := ch
	literal := string(ch)
	quoted := false
	line := s.line
	char := s.char

	switch {
	case s.isDecimal(ch):
		token = s.scanNumber(ch)
		literal = s.literal()
	case s.isIdentRune(ch):
		s.scanIdentifier()

		literal = s.literal()
		if _, e := ternary.Parse(literal); e == nil {
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
		} else if literal == SUBSTITUTION_OPERATOR {
			token = SUBSTITUTION_OP
		} else if 1 < len(literal) {
			token = UNCATEGORIZED
		}
	case s.isVariableSign(ch):
		if s.isVariableSign(s.peek()) {
			s.next()
			token = FLAG
		} else {
			token = VARIABLE
		}
		s.scanIdentifier()
		literal = s.literal()
	case s.isCommentRune(ch):
		s.scanComment()
		return s.Scan()
	case s.isLineCommentRune(ch):
		s.scanLineComment()
		return s.Scan()
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

	return Token{Token: int(token), Literal: literal, Quoted: quoted, Line: line, Char: char, SourceFile: s.sourceFile}, s.err
}

func (s *Scanner) scanString(quote rune) {
	for {
		ch := s.next()
		if ch == EOF {
			s.err = errors.New("literal not terminated")
			break
		} else if ch == quote {
			break
		}

		if ch == '\\' {
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
	return ch == '_' || ch == '$' || unicode.IsLetter(ch) || unicode.IsDigit(ch)
}

func (s *Scanner) isDecimal(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

func (s *Scanner) scanNumber(ch rune) rune {
	for s.isDecimal(s.peek()) {
		s.next()
	}

	if s.peek() == '.' {
		s.next()
		for s.isDecimal(s.peek()) {
			s.next()
		}
		return FLOAT
	}

	return INTEGER
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

func (s *Scanner) isVariableSign(ch rune) bool {
	if ch == VARIABLE_SIGN {
		return true
	}
	return false
}

func (s *Scanner) searchKeyword(str string) (int, error) {
	for i := KEYWORD_FROM; i <= KEYWORD_TO; i++ {
		if strings.EqualFold(TokenLiteral(i), str) {
			return i, nil
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

func (s *Scanner) isCommentRune(ch rune) bool {
	if ch == '/' && s.peek() == '*' {
		s.next()
		return true
	}
	return false
}

func (s *Scanner) scanComment() {
	for {
		ch := s.next()
		if ch == EOF {
			break
		} else if ch == '*' {
			if s.peek() == '/' {
				s.next()
				break
			}
		}
	}
	return
}

func (s *Scanner) isLineCommentRune(ch rune) bool {
	if ch == '-' && s.peek() == '-' {
		s.next()
		return true
	}
	return false
}

func (s *Scanner) scanLineComment() {
	for {
		ch := s.peek()
		if ch == '\r' || ch == '\n' || ch == EOF {
			break
		}
		s.next()
	}
	return
}
