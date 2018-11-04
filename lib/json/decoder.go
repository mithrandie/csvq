package json

import (
	"errors"
	"fmt"
)

type Decoder struct{}

func NewDecoder() *Decoder {
	return &Decoder{}
}

func (d Decoder) Decode(src string) (Structure, EscapeType, error) {
	st, et, err := ParseJson(src)
	if err != nil {
		se := err.(*SyntaxError)
		return st, et, errors.New(fmt.Sprintf("line %d, column %d: %s", se.Line, se.Column, se.Error()))
	}
	return st, et, nil
}
