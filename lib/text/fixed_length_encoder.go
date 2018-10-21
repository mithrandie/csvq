package text

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/value"
	"github.com/mithrandie/ternary"
	"strconv"
	"time"
)

type FixedLengthEncoder struct {
	DelimiterPositions []int
	LineBreak          cmd.LineBreak
	WithoutHeader      bool
	Encoding           cmd.Encoding

	buf bytes.Buffer
}

func NewFixedLengthEncoder() *FixedLengthEncoder {
	return &FixedLengthEncoder{
		DelimiterPositions: nil,
		LineBreak:          cmd.LF,
		WithoutHeader:      false,
	}
}

func (e *FixedLengthEncoder) Encode(fieldList []string, recordSet [][]value.Primary) (string, error) {
	prevPos := 0
	for _, endPos := range e.DelimiterPositions {
		if endPos < 0 || endPos <= prevPos {
			return "", errors.New(fmt.Sprintf("invalid delimiter position: %s", FormatIntSlice(e.DelimiterPositions)))
		}
		prevPos = endPos
	}

	var err error

	e.buf.Reset()

	if !e.WithoutHeader {
		start := 0
		for i, end := range e.DelimiterPositions {
			size := end - start
			if i < len(fieldList) {
				if err = e.addHeader(fieldList[i], size); err != nil {
					return e.buf.String(), err
				}
			} else {
				e.buf.Write(bytes.Repeat([]byte(string(PadChar)), size))
			}
			start = end
		}
	}

	for _, record := range recordSet {
		if 0 < e.buf.Len() {
			e.addLineBreak()
		}

		start := 0
		for i, end := range e.DelimiterPositions {
			size := end - start
			if i < len(record) {
				if err = e.addField(record[i], size); err != nil {
					return e.buf.String(), err
				}
			} else {
				e.buf.Write(bytes.Repeat([]byte(string(PadChar)), size))
			}
			start = end
		}
	}

	return e.buf.String(), nil
}

func (e *FixedLengthEncoder) addHeader(s string, byteSize int) error {
	size := ByteSize(s, e.Encoding)
	if byteSize < size {
		return errors.New(fmt.Sprintf("value is too long: %q for %d byte(s) length field", s, byteSize))
	}

	e.buf.WriteString(s)
	e.buf.Write(bytes.Repeat([]byte(string(PadChar)), byteSize-size))

	return nil
}

func (e *FixedLengthEncoder) addField(val value.Primary, byteSize int) error {
	var s string
	switch val.(type) {
	case value.String:
		s = val.(value.String).Raw()
	case value.Integer:
		s = val.(value.Integer).String()
	case value.Float:
		s = val.(value.Float).String()
	case value.Boolean:
		s = val.(value.Boolean).String()
	case value.Ternary:
		t := val.(value.Ternary)
		if t.Ternary() == ternary.UNKNOWN {
			s = ""
		} else {
			s = strconv.FormatBool(t.Ternary().ParseBool())
		}
	case value.Datetime:
		s = val.(value.Datetime).Format(time.RFC3339Nano)
	case value.Null:
		s = ""
	}

	size := ByteSize(s, e.Encoding)
	if byteSize < size {
		return errors.New(fmt.Sprintf("value is too long: %q for %d byte(s) length field", s, byteSize))
	}

	switch val.(type) {
	case value.Integer, value.Float:
		e.buf.Write(bytes.Repeat([]byte(string(PadChar)), byteSize-size))
		e.buf.WriteString(s)
	default:
		e.buf.WriteString(s)
		e.buf.Write(bytes.Repeat([]byte(string(PadChar)), byteSize-size))
	}

	return nil
}

func (e *FixedLengthEncoder) addLineBreak() {
	e.buf.WriteString(e.LineBreak.Value())
}
