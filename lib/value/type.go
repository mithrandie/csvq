package value

import (
	"strconv"
	"strings"
	"time"

	"github.com/mithrandie/ternary"
)

func IsNull(v Primary) bool {
	_, ok := v.(Null)
	return ok
}

type RowValue []Primary

type Primary interface {
	String() string
	Ternary() ternary.Value
}

type String struct {
	literal string
}

func (s String) String() string {
	return quoteString(s.literal)
}

func NewString(s string) String {
	return String{
		literal: s,
	}
}

func (s String) Raw() string {
	return s.literal
}

func (s String) Ternary() ternary.Value {
	lit := strings.TrimSpace(s.Raw())
	if b, err := strconv.ParseBool(lit); err == nil {
		return ternary.ConvertFromBool(b)
	}
	return ternary.UNKNOWN
}

type Integer struct {
	value int64
}

func NewIntegerFromString(s string) Integer {
	i, _ := strconv.ParseInt(s, 10, 64)
	return Integer{
		value: i,
	}
}

func NewInteger(i int64) Integer {
	return Integer{
		value: i,
	}
}

func (i Integer) String() string {
	return Int64ToStr(i.value)
}

func (i Integer) Raw() int64 {
	return i.value
}

func (i Integer) Ternary() ternary.Value {
	switch i.Raw() {
	case 0:
		return ternary.FALSE
	case 1:
		return ternary.TRUE
	default:
		return ternary.UNKNOWN
	}
}

type Float struct {
	value float64
}

func NewFloatFromString(s string) Float {
	f, _ := strconv.ParseFloat(s, 64)
	return Float{
		value: f,
	}
}

func NewFloat(f float64) Float {
	return Float{
		value: f,
	}
}

func (f Float) String() string {
	return Float64ToStr(f.value)
}

func (f Float) Raw() float64 {
	return f.value
}

func (f Float) Ternary() ternary.Value {
	switch f.Raw() {
	case 0:
		return ternary.FALSE
	case 1:
		return ternary.TRUE
	default:
		return ternary.UNKNOWN
	}
}

type Boolean struct {
	value bool
}

func NewBoolean(b bool) Boolean {
	return Boolean{
		value: b,
	}
}

func (b Boolean) String() string {
	return strconv.FormatBool(b.value)
}

func (b Boolean) Raw() bool {
	return b.value
}

func (b Boolean) Ternary() ternary.Value {
	return ternary.ConvertFromBool(b.Raw())
}

type Ternary struct {
	value ternary.Value
}

func NewTernaryFromString(s string) Ternary {
	t, _ := ternary.ConvertFromString(s)
	return Ternary{
		value: t,
	}
}

func NewTernary(t ternary.Value) Ternary {
	return Ternary{
		value: t,
	}
}

func (t Ternary) String() string {
	return t.value.String()
}

func (t Ternary) Ternary() ternary.Value {
	return t.value
}

type Datetime struct {
	value time.Time
}

func NewDatetimeFromString(s string) Datetime {
	t, _ := StrToTime(s)
	return Datetime{
		value: t,
	}
}

func NewDatetime(t time.Time) Datetime {
	return Datetime{
		value: t,
	}
}

func (dt Datetime) String() string {
	return quoteString(dt.value.Format(time.RFC3339Nano))
}

func (dt Datetime) Raw() time.Time {
	return dt.value
}

func (dt Datetime) Ternary() ternary.Value {
	return ternary.UNKNOWN
}

func (dt Datetime) Format(s string) string {
	return dt.value.Format(s)
}

type Null struct{}

func NewNull() Null {
	return Null{}
}

func (n Null) String() string {
	return "NULL"
}

func (n Null) Ternary() ternary.Value {
	return ternary.UNKNOWN
}

func quoteString(s string) string {
	return "'" + s + "'"
}
