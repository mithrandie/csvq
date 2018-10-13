package json

import (
	"strconv"
	"strings"
)

const (
	BeginArray     = '['
	BeginObject    = '{'
	EndArray       = ']'
	EndObject      = '}'
	NameSeparator  = ':'
	ValueSeparator = ','
	QuotationMark  = '"'
	EscapeMark     = '\\'
)

var WhiteSpaces = []rune{
	32, //Space
	9,  //Horizontal Tab
	10, //Line Feed
	13, //Carriage Return
}

const (
	FalseValue = "false"
	TrueValue  = "true"
	NullValue  = "null"
)

const EOF = -1

type Structure interface {
	String() string
}

type Object struct {
	Members []ObjectMember
}

func NewObject(capacity int) Object {
	return Object{
		Members: make([]ObjectMember, 0, capacity),
	}
}

func (obj *Object) Len() int {
	return len(obj.Members)
}

func (obj *Object) Add(key string, val Structure) {
	obj.Members = append(obj.Members, ObjectMember{Key: key, Value: val})
}

func (obj *Object) Exists(key string) bool {
	for _, m := range obj.Members {
		if m.Key == key {
			return true
		}
	}
	return false
}

func (obj *Object) Value(key string) Structure {
	for _, m := range obj.Members {
		if m.Key == key {
			return m.Value
		}
	}
	return nil
}

func (obj *Object) Update(key string, val Structure) {
	for i, m := range obj.Members {
		if m.Key == key {
			obj.Members[i].Value = val
			break
		}
	}
}

func (obj Object) String() string {
	strs := make([]string, 0, obj.Len())
	for _, m := range obj.Members {
		strs = append(strs, Quote(Escape(m.Key))+string(NameSeparator)+m.Value.String())
	}
	return string(BeginObject) + strings.Join(strs[:], string(ValueSeparator)) + string(EndObject)
}

type ObjectMember struct {
	Key   string
	Value Structure
}

type Array []Structure

func (ar Array) String() string {
	strs := make([]string, 0, len(ar))
	for _, v := range ar {
		strs = append(strs, v.String())
	}
	return string(BeginArray) + strings.Join(strs[:], string(ValueSeparator)) + string(EndArray)
}

type Number float64

func (n Number) String() string {
	return strconv.FormatFloat(float64(n), 'f', -1, 64)
}

type String string

func (s String) String() string {
	return Quote(Escape(string(s)))
}

type Boolean bool

func (b Boolean) String() string {
	if b {
		return TrueValue
	}
	return FalseValue
}

type Null struct{}

func (n Null) String() string {
	return NullValue
}
