package query

import (
	"strings"
	"unicode/utf8"

	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/value"

	"github.com/mithrandie/ternary"
)

func Is(p1 value.Primary, p2 value.Primary) ternary.Value {
	if value.IsNull(p2) {
		return ternary.ConvertFromBool(value.IsNull(p1))
	}

	return ternary.Equal(p1.Ternary(), p2.Ternary())
}

func Like(p1 value.Primary, p2 value.Primary) ternary.Value {
	if value.IsNull(p1) || value.IsNull(p2) {
		return ternary.UNKNOWN
	}

	s1 := value.ToString(p1)
	if value.IsNull(s1) {
		return ternary.UNKNOWN
	}
	s2 := value.ToString(p2)
	if value.IsNull(s2) {
		return ternary.UNKNOWN
	}

	s := strings.ToUpper(p1.(value.String).Raw())
	pattern := strings.ToUpper(p2.(value.String).Raw())

	if s == pattern {
		return ternary.TRUE
	}
	if len(pattern) < 1 {
		return ternary.FALSE
	}

	patternRunes := []rune(pattern)
	patternPos := 0

	for {
		anyRunesMinLen, anyRunexMaxLen, search, pos := stringPattern(patternRunes, patternPos)
		patternPos = pos

		anyString := s
		if 0 < len(search) {
			idx := strings.Index(s, search)
			if idx < 0 {
				return ternary.FALSE
			}
			anyString = s[:idx]
		}

		if utf8.RuneCountInString(anyString) < anyRunesMinLen {
			return ternary.FALSE
		}
		if -1 < anyRunexMaxLen && anyRunexMaxLen < utf8.RuneCountInString(anyString) {
			return ternary.FALSE
		}

		if len(patternRunes) <= patternPos {
			if len(anyString+search) < len(s) {
				return ternary.FALSE
			}
			break
		}

		s = s[len(anyString+search):]
	}

	return ternary.TRUE
}

func stringPattern(pattern []rune, position int) (int, int, string, int) {
	anyRunesMinLen := 0
	anyRunesMaxLen := 0
	search := []rune{}
	returnPostion := position

	escaped := false
	for i := position; i < len(pattern); i++ {
		r := pattern[i]

		if escaped {
			switch r {
			case '%', '_':
				search = append(search, r)
			default:
				search = append(search, '\\', r)
			}
			returnPostion++
			escaped = false
			continue
		}

		if (r == '%' || r == '_') && 0 < len(search) {
			break
		}
		returnPostion++

		switch r {
		case '%':
			anyRunesMaxLen = -1
		case '_':
			anyRunesMinLen++
			if -1 < anyRunesMaxLen {
				anyRunesMaxLen++
			}
		case '\\':
			escaped = true
		default:
			search = append(search, r)
		}
	}
	if escaped {
		search = append(search, '\\')
	}

	return anyRunesMinLen, anyRunesMaxLen, string(search), returnPostion
}

func InRowValueList(rowValue value.RowValue, list []value.RowValue, matchType int, operator string) (ternary.Value, error) {
	results := make([]ternary.Value, len(list))

	for i, v := range list {
		t, err := value.CompareRowValues(rowValue, v, operator)
		if err != nil {
			return ternary.FALSE, NewRowValueLengthInListError(i)
		}
		switch matchType {
		case parser.ANY:
			if t == ternary.TRUE {
				return ternary.TRUE, nil
			}
		default: // parser.ALL
			if t == ternary.FALSE {
				return ternary.FALSE, nil
			}
		}

		results[i] = t
	}

	switch matchType {
	case parser.ANY:
		return ternary.Any(results), nil
	default: // parser.ALL
		return ternary.All(results), nil
	}
}

func Any(rowValue value.RowValue, list []value.RowValue, operator string) (ternary.Value, error) {
	return InRowValueList(rowValue, list, parser.ANY, operator)
}

func All(rowValue value.RowValue, list []value.RowValue, operator string) (ternary.Value, error) {
	return InRowValueList(rowValue, list, parser.ALL, operator)
}
