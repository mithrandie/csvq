package query

import (
	"strings"

	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/ternary"
	"time"
)

type SortValueType int

const (
	SORT_VALUE_NULL SortValueType = iota
	SORT_VALUE_INTEGER
	SORT_VALUE_FLOAT
	SORT_VALUE_DATETIME
	SORT_VALUE_BOOLEAN
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
	if compareValues == nil {
		return false
	}

	for i, value := range values {
		if !value.EquivalentTo(compareValues[i]) {
			return false
		}
	}
	return true
}

func (values SortValues) Serialize() string {
	list := make([]string, len(values))

	for i, value := range values {
		switch value.Type {
		case SORT_VALUE_NULL:
			list[i] = serializeNull()
		case SORT_VALUE_INTEGER:
			list[i] = serializeInteger(value.Integer)
		case SORT_VALUE_FLOAT:
			list[i] = serializeFlaot(value.Float)
		case SORT_VALUE_DATETIME:
			list[i] = serializeDatetime(value.Datetime)
		case SORT_VALUE_BOOLEAN:
			list[i] = serializeBoolean(value.Boolean)
		case SORT_VALUE_STRING:
			list[i] = serializeString(value.String)
		}
	}

	return strings.Join(list, ":")
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
	} else if i := parser.PrimaryToInteger(value); !parser.IsNull(i) {
		sortValue.Type = SORT_VALUE_INTEGER
		s := parser.PrimaryToString(value)
		sortValue.Integer = i.(parser.Integer).Value()
		sortValue.Float = float64(sortValue.Integer)
		sortValue.Datetime = time.Unix(sortValue.Integer, 0)
		sortValue.String = s.(parser.String).Value()
	} else if f := parser.PrimaryToFloat(value); !parser.IsNull(f) {
		sortValue.Type = SORT_VALUE_FLOAT
		s := parser.PrimaryToString(value)
		sortValue.Float = f.(parser.Float).Value()
		sortValue.Datetime = parser.Float64ToTime(sortValue.Float)
		sortValue.String = s.(parser.String).Value()
	} else if dt := parser.PrimaryToDatetime(value); !parser.IsNull(dt) {
		t := dt.(parser.Datetime).Value()
		if t.Nanosecond() > 0 {
			f := float64(t.Unix()) + float64(t.Nanosecond())/float64(1000000000)
			t2 := parser.Float64ToTime(f)
			if t.Equal(t2) {
				sortValue.Type = SORT_VALUE_FLOAT
				sortValue.Float = f
				sortValue.Datetime = t
				sortValue.String = parser.Float64ToStr(f)
			} else {
				sortValue.Type = SORT_VALUE_DATETIME
				sortValue.Datetime = t
			}
		} else {
			sortValue.Type = SORT_VALUE_INTEGER
			i := t.Unix()
			sortValue.Integer = i
			sortValue.Float = float64(i)
			sortValue.Datetime = t
			sortValue.String = parser.Int64ToStr(i)
		}
	} else if b := parser.PrimaryToBoolean(value); !parser.IsNull(b) {
		sortValue.Type = SORT_VALUE_BOOLEAN
		sortValue.Boolean = b.(parser.Boolean).Value()
		if sortValue.Boolean {
			sortValue.Integer = 1
		} else {
			sortValue.Integer = 0
		}
	} else if s, ok := value.(parser.String); ok {
		sortValue.Type = SORT_VALUE_STRING
		sortValue.String = strings.ToUpper(strings.TrimSpace(s.Value()))
	} else {
		sortValue.Type = SORT_VALUE_NULL
	}

	return sortValue
}

func (v *SortValue) Less(compareValue *SortValue) ternary.Value {
	switch v.Type {
	case SORT_VALUE_INTEGER:
		switch compareValue.Type {
		case SORT_VALUE_INTEGER:
			if v.Integer == compareValue.Integer {
				return ternary.UNKNOWN
			}
			return ternary.ParseBool(v.Integer < compareValue.Integer)
		case SORT_VALUE_FLOAT:
			return ternary.ParseBool(v.Float < compareValue.Float)
		case SORT_VALUE_DATETIME:
			return ternary.ParseBool(v.Datetime.Before(compareValue.Datetime))
		case SORT_VALUE_STRING:
			return ternary.ParseBool(v.String < compareValue.String)
		}
	case SORT_VALUE_FLOAT:
		switch compareValue.Type {
		case SORT_VALUE_INTEGER, SORT_VALUE_FLOAT:
			if v.Float == compareValue.Float {
				return ternary.UNKNOWN
			}
			return ternary.ParseBool(v.Float < compareValue.Float)
		case SORT_VALUE_DATETIME:
			return ternary.ParseBool(v.Datetime.Before(compareValue.Datetime))
		case SORT_VALUE_STRING:
			return ternary.ParseBool(v.String < compareValue.String)
		}
	case SORT_VALUE_DATETIME:
		switch compareValue.Type {
		case SORT_VALUE_INTEGER, SORT_VALUE_FLOAT, SORT_VALUE_DATETIME:
			if v.Datetime.Equal(compareValue.Datetime) {
				return ternary.UNKNOWN
			}
			return ternary.ParseBool(v.Datetime.Before(compareValue.Datetime))
		}
	case SORT_VALUE_STRING:
		switch compareValue.Type {
		case SORT_VALUE_INTEGER, SORT_VALUE_FLOAT, SORT_VALUE_STRING:
			if v.String == compareValue.String {
				return ternary.UNKNOWN
			}
			return ternary.ParseBool(v.String < compareValue.String)
		}
	}

	return ternary.UNKNOWN
}

func (v *SortValue) EquivalentTo(compareValue *SortValue) bool {
	switch v.Type {
	case SORT_VALUE_INTEGER:
		switch compareValue.Type {
		case SORT_VALUE_INTEGER, SORT_VALUE_BOOLEAN:
			return v.Integer == compareValue.Integer
		}
	case SORT_VALUE_FLOAT:
		switch compareValue.Type {
		case SORT_VALUE_FLOAT:
			return v.Float == compareValue.Float
		}
	case SORT_VALUE_DATETIME:
		switch compareValue.Type {
		case SORT_VALUE_DATETIME:
			return v.Datetime.Equal(compareValue.Datetime)
		}
	case SORT_VALUE_BOOLEAN:
		switch compareValue.Type {
		case SORT_VALUE_INTEGER:
			return v.Integer == compareValue.Integer
		case SORT_VALUE_BOOLEAN:
			return v.Boolean == compareValue.Boolean
		}
	case SORT_VALUE_STRING:
		switch compareValue.Type {
		case SORT_VALUE_STRING:
			return v.String == compareValue.String
		}
	case SORT_VALUE_NULL:
		return compareValue.Type == SORT_VALUE_NULL
	}

	return false
}
