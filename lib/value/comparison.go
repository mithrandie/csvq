package value

import (
	"errors"
	"strings"

	"github.com/mithrandie/ternary"
)

type ComparisonResult int

const (
	EQUAL ComparisonResult = iota
	BOOL_EQUAL
	NOT_EQUAL
	LESS
	GREATER
	INCOMMENSURABLE
)

var comparisonResultLiterals = map[ComparisonResult]string{
	EQUAL:           "EQUAL",
	BOOL_EQUAL:      "BOOL_EQUAL",
	NOT_EQUAL:       "NOT_EQUAL",
	LESS:            "LESS",
	GREATER:         "GREATER",
	INCOMMENSURABLE: "INCOMMENSURABLE",
}

func (cr ComparisonResult) String() string {
	return comparisonResultLiterals[cr]
}

func CompareCombinedly(p1 Primary, p2 Primary) ComparisonResult {
	if IsNull(p1) || IsNull(p2) {
		return INCOMMENSURABLE
	}

	if i1 := ToInteger(p1); !IsNull(i1) {
		if i2 := ToInteger(p2); !IsNull(i2) {
			v1 := i1.(Integer).Raw()
			v2 := i2.(Integer).Raw()
			if v1 == v2 {
				return EQUAL
			} else if v1 < v2 {
				return LESS
			} else {
				return GREATER
			}
		}
	}

	if f1 := ToFloat(p1); !IsNull(f1) {
		if f2 := ToFloat(p2); !IsNull(f2) {
			v1 := f1.(Float).Raw()
			v2 := f2.(Float).Raw()
			if v1 == v2 {
				return EQUAL
			} else if v1 < v2 {
				return LESS
			} else {
				return GREATER
			}
		}
	}

	if d1 := ToDatetime(p1); !IsNull(d1) {
		if d2 := ToDatetime(p2); !IsNull(d2) {
			v1 := d1.(Datetime).Raw()
			v2 := d2.(Datetime).Raw()
			if v1.Equal(v2) {
				return EQUAL
			} else if v1.Before(v2) {
				return LESS
			} else {
				return GREATER
			}
		}
	}

	if b1 := ToBoolean(p1); !IsNull(b1) {
		if b2 := ToBoolean(p2); !IsNull(b2) {
			v1 := b1.(Boolean).Raw()
			v2 := b2.(Boolean).Raw()
			if v1 == v2 {
				return BOOL_EQUAL
			} else {
				return NOT_EQUAL
			}
		}
	}

	if s1, ok := p1.(String); ok {
		if s2, ok := p2.(String); ok {
			v1 := strings.ToUpper(strings.TrimSpace(s1.Raw()))
			v2 := strings.ToUpper(strings.TrimSpace(s2.Raw()))

			if v1 == v2 {
				return EQUAL
			} else if v1 < v2 {
				return LESS
			} else {
				return GREATER
			}
		}
	}

	return INCOMMENSURABLE
}

func Equal(p1 Primary, p2 Primary) ternary.Value {
	if r := CompareCombinedly(p1, p2); r != INCOMMENSURABLE {
		return ternary.ConvertFromBool(r == EQUAL || r == BOOL_EQUAL)
	}
	return ternary.UNKNOWN
}

func NotEqual(p1 Primary, p2 Primary) ternary.Value {
	if r := CompareCombinedly(p1, p2); r != INCOMMENSURABLE {
		return ternary.ConvertFromBool(r != EQUAL && r != BOOL_EQUAL)
	}
	return ternary.UNKNOWN
}

func Less(p1 Primary, p2 Primary) ternary.Value {
	if r := CompareCombinedly(p1, p2); r != INCOMMENSURABLE && r != NOT_EQUAL && r != BOOL_EQUAL {
		return ternary.ConvertFromBool(r == LESS)
	}
	return ternary.UNKNOWN
}

func Greater(p1 Primary, p2 Primary) ternary.Value {
	if r := CompareCombinedly(p1, p2); r != INCOMMENSURABLE && r != NOT_EQUAL && r != BOOL_EQUAL {
		return ternary.ConvertFromBool(r == GREATER)
	}
	return ternary.UNKNOWN
}

func LessOrEqual(p1 Primary, p2 Primary) ternary.Value {
	if r := CompareCombinedly(p1, p2); r != INCOMMENSURABLE && r != NOT_EQUAL && r != BOOL_EQUAL {
		return ternary.ConvertFromBool(r != GREATER)
	}
	return ternary.UNKNOWN
}

func GreaterOrEqual(p1 Primary, p2 Primary) ternary.Value {
	if r := CompareCombinedly(p1, p2); r != INCOMMENSURABLE && r != NOT_EQUAL && r != BOOL_EQUAL {
		return ternary.ConvertFromBool(r != LESS)
	}
	return ternary.UNKNOWN
}

func Compare(p1 Primary, p2 Primary, operator string) ternary.Value {
	switch operator {
	case "=":
		return Equal(p1, p2)
	case ">":
		return Greater(p1, p2)
	case "<":
		return Less(p1, p2)
	case ">=":
		return GreaterOrEqual(p1, p2)
	case "<=":
		return LessOrEqual(p1, p2)
	default: //case "<>", "!=":
		return NotEqual(p1, p2)
	}
}

func CompareRowValues(rowValue1 RowValue, rowValue2 RowValue, operator string) (ternary.Value, error) {
	if rowValue1 == nil || rowValue2 == nil {
		return ternary.UNKNOWN, nil
	}

	if len(rowValue1) != len(rowValue2) {
		return ternary.FALSE, errors.New("row value length does not match")
	}

	unknown := false
	for i := 0; i < len(rowValue1); i++ {
		r := CompareCombinedly(rowValue1[i], rowValue2[i])

		if r == INCOMMENSURABLE {
			switch operator {
			case "=", "<>", "!=":
				if i < len(rowValue1)-1 {
					unknown = true
					continue
				}
			}

			return ternary.UNKNOWN, nil
		}

		switch operator {
		case ">", "<", ">=", "<=":
			if r == NOT_EQUAL || r == BOOL_EQUAL {
				return ternary.UNKNOWN, nil
			}
		}

		switch operator {
		case "=":
			if r != EQUAL && r != BOOL_EQUAL {
				return ternary.FALSE, nil
			}
		case ">", ">=":
			switch r {
			case GREATER:
				return ternary.TRUE, nil
			case LESS:
				return ternary.FALSE, nil
			}
		case "<", "<=":
			switch r {
			case LESS:
				return ternary.TRUE, nil
			case GREATER:
				return ternary.FALSE, nil
			}
		case "<>", "!=":
			if r != EQUAL && r != BOOL_EQUAL {
				return ternary.TRUE, nil
			}
		}
	}

	if unknown {
		return ternary.UNKNOWN, nil
	}

	switch operator {
	case ">", "<", "<>", "!=":
		return ternary.FALSE, nil
	}
	return ternary.TRUE, nil
}

func Equivalent(p1 Primary, p2 Primary) ternary.Value {
	if IsNull(p1) && IsNull(p2) {
		return ternary.TRUE
	}
	return Equal(p1, p2)
}
