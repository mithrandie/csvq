package query

import (
	"strings"
	"time"

	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/ternary"
)

type SortValueType int

const (
	SORT_VALUE_NULL SortValueType = iota
	SORT_VALUE_INTEGER
	SORT_VALUE_FLOAT
	SORT_VALUE_DATETIME
	SORT_VALUE_STRING
)

type SortValues []*SortValue

func (values SortValues) Less(compareValues SortValues, directions []int, nullPositions []int) bool {
	for i, value := range values {
		t := value.Less(compareValues[i])
		if t != ternary.UNKNOWN {
			if directions[i] == parser.ASC {
				return t == ternary.TRUE
			} else {
				return t == ternary.FALSE
			}
		}

		if value.Type == SORT_VALUE_NULL && compareValues[i].Type != SORT_VALUE_NULL {
			if nullPositions[i] == parser.FIRST {
				return true
			} else {
				return false
			}
		}
		if value.Type != SORT_VALUE_NULL && compareValues[i].Type == SORT_VALUE_NULL {
			if nullPositions[i] == parser.FIRST {
				return false
			} else {
				return true
			}
		}
	}
	return false
}

func (values SortValues) EquivalentTo(compareValues SortValues) bool {
	for i, value := range values {
		if !value.EquivalentTo(compareValues[i]) {
			return false
		}
	}
	return true
}

type SortValue struct {
	Type SortValueType

	Integer  int64
	Float    float64
	Datetime time.Time
	String   string
	Boolean  bool
}

func NewSortValue(value parser.Primary) *SortValue {
	sortValue := &SortValue{}

	if parser.IsNull(value) {
		sortValue.Type = SORT_VALUE_NULL
	} else if in := parser.PrimaryToInteger(value); !parser.IsNull(in) {
		sortValue.Type = SORT_VALUE_INTEGER
		sortValue.Integer = in.(parser.Integer).Value()
	} else if f := parser.PrimaryToFloat(value); !parser.IsNull(f) {
		sortValue.Type = SORT_VALUE_FLOAT
		sortValue.Float = f.(parser.Float).Value()
	} else if dt := parser.PrimaryToDatetime(value); !parser.IsNull(dt) {
		sortValue.Type = SORT_VALUE_DATETIME
		sortValue.Datetime = dt.(parser.Datetime).Value()
	} else if b := parser.PrimaryToBoolean(value); !parser.IsNull(b) {
		sortValue.Type = SORT_VALUE_NULL
	} else if s, ok := value.(parser.String); ok {
		sortValue.Type = SORT_VALUE_STRING
		sortValue.String = strings.ToUpper(strings.TrimSpace(s.Value()))
	} else {
		sortValue.Type = SORT_VALUE_NULL
	}

	return sortValue
}

func (v SortValue) Less(compareValue *SortValue) ternary.Value {
	if v.Type == SORT_VALUE_INTEGER && compareValue.Type == SORT_VALUE_FLOAT {
		f := float64(v.Integer)
		return ternary.ParseBool(f < compareValue.Float)
	}
	if v.Type == SORT_VALUE_FLOAT && compareValue.Type == SORT_VALUE_INTEGER {
		f := float64(compareValue.Integer)
		return ternary.ParseBool(v.Float < f)
	}

	if v.Type != compareValue.Type {
		return ternary.UNKNOWN
	}

	switch v.Type {
	case SORT_VALUE_INTEGER:
		if v.Integer == compareValue.Integer {
			return ternary.UNKNOWN
		}
		return ternary.ParseBool(v.Integer < compareValue.Integer)
	case SORT_VALUE_FLOAT:
		if v.Float == compareValue.Float {
			return ternary.UNKNOWN
		}
		return ternary.ParseBool(v.Float < compareValue.Float)
	case SORT_VALUE_DATETIME:
		if v.Datetime.Equal(compareValue.Datetime) {
			return ternary.UNKNOWN
		}
		return ternary.ParseBool(v.Datetime.Before(compareValue.Datetime))
	case SORT_VALUE_STRING:
		if v.String == compareValue.String {
			return ternary.UNKNOWN
		}
		return ternary.ParseBool(v.String < compareValue.String)
	}

	return ternary.UNKNOWN
}

func (v SortValue) EquivalentTo(compareValue *SortValue) bool {
	if v.Type != compareValue.Type {
		return false
	}

	switch v.Type {
	case SORT_VALUE_INTEGER:
		return v.Integer == compareValue.Integer
	case SORT_VALUE_FLOAT:
		return v.Float == compareValue.Float
	case SORT_VALUE_DATETIME:
		return v.Datetime.Equal(compareValue.Datetime)
	case SORT_VALUE_STRING:
		return v.String == compareValue.String
	default: //SORT_VALUE_NULL
		return true
	}
}
