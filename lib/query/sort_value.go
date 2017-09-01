package query

import (
	"strings"

	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/ternary"
	"github.com/mithrandie/csvq/lib/value"
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
	for i, val := range values {
		t := val.Less(compareValues[i])
		if t != ternary.UNKNOWN {
			if directions[i] == parser.ASC {
				return t == ternary.TRUE
			} else {
				return t == ternary.FALSE
			}
		}

		if val.Type == SORT_VALUE_NULL && compareValues[i].Type != SORT_VALUE_NULL {
			if nullPositions[i] == parser.FIRST {
				return true
			} else {
				return false
			}
		}
		if val.Type != SORT_VALUE_NULL && compareValues[i].Type == SORT_VALUE_NULL {
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

	for i, val := range values {
		if !val.EquivalentTo(compareValues[i]) {
			return false
		}
	}
	return true
}

func (values SortValues) Serialize() string {
	list := make([]string, len(values))

	for i, val := range values {
		switch val.Type {
		case SORT_VALUE_NULL:
			list[i] = serializeNull()
		case SORT_VALUE_INTEGER:
			list[i] = serializeInteger(val.Integer)
		case SORT_VALUE_FLOAT:
			list[i] = serializeFlaot(val.Float)
		case SORT_VALUE_DATETIME:
			list[i] = serializeDatetimeFromUnixNano(val.Datetime)
		case SORT_VALUE_BOOLEAN:
			list[i] = serializeBoolean(val.Boolean)
		case SORT_VALUE_STRING:
			list[i] = serializeString(val.String)
		}
	}

	return strings.Join(list, ":")
}

type SortValue struct {
	Type SortValueType

	Integer  int64
	Float    float64
	Datetime int64
	String   string
	Boolean  bool
}

func NewSortValue(val value.Primary) *SortValue {
	sortValue := &SortValue{}

	if value.IsNull(val) {
		sortValue.Type = SORT_VALUE_NULL
	} else if i := value.ToInteger(val); !value.IsNull(i) {
		s := value.ToString(val)
		sortValue.Type = SORT_VALUE_INTEGER
		sortValue.Integer = i.(value.Integer).Raw()
		sortValue.Float = float64(sortValue.Integer)
		sortValue.Datetime = sortValue.Integer * 1e9
		sortValue.String = s.(value.String).Raw()
	} else if f := value.ToFloat(val); !value.IsNull(f) {
		s := value.ToString(val)
		sortValue.Type = SORT_VALUE_FLOAT
		sortValue.Float = f.(value.Float).Raw()
		sortValue.Datetime = int64(sortValue.Float * 1e9)
		sortValue.String = s.(value.String).Raw()
	} else if dt := value.ToDatetime(val); !value.IsNull(dt) {
		t := dt.(value.Datetime).Raw()
		if t.Nanosecond() > 0 {
			f := float64(t.Unix()) + float64(t.Nanosecond())/1e9
			t2 := value.Float64ToTime(f)
			if t.Equal(t2) {
				sortValue.Type = SORT_VALUE_FLOAT
				sortValue.Float = f
				sortValue.Datetime = t.UnixNano()
				sortValue.String = value.Float64ToStr(f)
			} else {
				sortValue.Type = SORT_VALUE_DATETIME
				sortValue.Datetime = t.UnixNano()
			}
		} else {
			sortValue.Type = SORT_VALUE_INTEGER
			i := t.Unix()
			sortValue.Integer = i
			sortValue.Float = float64(i)
			sortValue.Datetime = t.UnixNano()
			sortValue.String = value.Int64ToStr(i)
		}
	} else if b := value.ToBoolean(val); !value.IsNull(b) {
		sortValue.Type = SORT_VALUE_BOOLEAN
		sortValue.Boolean = b.(value.Boolean).Raw()
		if sortValue.Boolean {
			sortValue.Integer = 1
		} else {
			sortValue.Integer = 0
		}
	} else if s, ok := val.(value.String); ok {
		sortValue.Type = SORT_VALUE_STRING
		sortValue.String = strings.ToUpper(strings.TrimSpace(s.Raw()))
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
			return ternary.ParseBool(v.Datetime < compareValue.Datetime)
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
			return ternary.ParseBool(v.Datetime < compareValue.Datetime)
		case SORT_VALUE_STRING:
			return ternary.ParseBool(v.String < compareValue.String)
		}
	case SORT_VALUE_DATETIME:
		switch compareValue.Type {
		case SORT_VALUE_INTEGER, SORT_VALUE_FLOAT, SORT_VALUE_DATETIME:
			if v.Datetime == compareValue.Datetime {
				return ternary.UNKNOWN
			}
			return ternary.ParseBool(v.Datetime < compareValue.Datetime)
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
			return v.Datetime == compareValue.Datetime
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
