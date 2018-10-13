package json

import (
	"errors"
	"fmt"
	"strconv"
	"unicode"
)

type Scanner struct {
	src    []rune
	srcPos int
	offset int

	line       int
	column     int
	sourceFile string

	err error
}

func (s *Scanner) Init(src string) *Scanner {
	s.src = []rune(src)
	s.srcPos = 0
	s.offset = 0
	s.line = 1
	s.column = 0

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
	s.column++

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
	s.column = 0
	return s.src[s.srcPos-1]
}

func (s *Scanner) runes() []rune {
	return s.src[(s.srcPos - s.offset):s.srcPos]
}

func (s *Scanner) literal() string {
	return string(s.runes())
}

func (s *Scanner) trimQuotes() string {
	runes := s.runes()
	quote := runes[0]
	switch quote {
	case '"':
		if 1 < len(runes) && runes[0] == quote && runes[len(runes)-1] == quote {
			runes = runes[1:(len(runes) - 1)]
		}
	}
	return string(runes)
}

func (s *Scanner) Scan() (Token, error) {
	ch := s.skipSpaces()

	s.offset = 0
	s.next()

	token := ch
	literal := string(ch)
	line := s.line
	column := s.column

	switch {
	case s.isPositiveDecimal(ch) || ch == '-':
		s.scanNumber(ch)
		literal = s.literal()
		if s.err == nil {
			_, err := strconv.ParseFloat(literal, 64)
			if err != nil {
				s.err = errors.New(fmt.Sprintf("could not convert %q into float64", literal))
			}
		}
		token = NUMBER
	case s.isLiteral(ch):
		s.scanLiteral()
		literal = s.literal()

		switch literal {
		case TrueValue, FalseValue:
			token = BOOLEAN
		case NullValue:
			token = NULL
		}
	case ch == '"':
		s.scanString(ch)
		literal = Unescape(s.trimQuotes())
		token = STRING
	}

	return Token{Token: int(token), Literal: literal, Line: line, Column: column}, s.err
}

func (s *Scanner) skipSpaces() rune {
	for s.isWhiteSpace(s.peek()) {
		s.next()
	}
	return s.peek()
}

func (s *Scanner) isWhiteSpace(ch rune) bool {
	for _, r := range WhiteSpaces {
		if r == ch {
			return true
		}
	}
	return false
}

func (s *Scanner) isDecimal(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

func (s *Scanner) isPositiveDecimal(ch rune) bool {
	return '1' <= ch && ch <= '9'
}

func (s *Scanner) isLiteral(ch rune) bool {
	return unicode.IsLower(ch)
}

func (s *Scanner) scanLiteral() {
	for s.isLiteral(s.peek()) {
		s.next()
	}
}

func (s *Scanner) scanString(quote rune) {
	for {
		ch := s.next()
		if ch == EOF {
			s.err = errors.New("string not terminated")
			break
		}

		if ch == quote {
			break
		}

		if ch == EscapeMark {
			switch s.peek() {
			case quote:
				s.next()
			}
		}
	}
}

func (s *Scanner) scanDecimal() {
	for s.isDecimal(s.peek()) {
		s.next()
	}
}

func (s *Scanner) scanNumber(ch rune) {
	if ch == '-' {
		ch = s.next()
	}

	if !s.isDecimal(ch) {
		s.err = errors.New("invalid number")
		return
	}

	if s.isPositiveDecimal(ch) {
		s.scanDecimal()
	}

	if s.peek() == '.' {
		s.next()
		if !s.isDecimal(s.peek()) {
			s.next()
			s.err = errors.New("invalid number")
			return
		}
		s.scanDecimal()
	}

	if s.peek() == 'e' || s.peek() == 'E' {
		s.next()
		if s.peek() == '+' || s.peek() == '-' {
			s.next()
		}
		if !s.isDecimal(s.peek()) {
			s.next()
			s.err = errors.New("invalid number")
			return
		}
		s.scanDecimal()
	}

	return
}
