package csv

import (
	"bytes"
	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/value"
	"github.com/mithrandie/ternary"
	"strconv"
	"strings"
	"time"
)

type Encoder struct {
	Delimiter     rune
	LineBreak     cmd.LineBreak
	WithoutHeader bool

	buf bytes.Buffer
}

func NewEncoder() *Encoder {
	return &Encoder{
		Delimiter:     ',',
		LineBreak:     cmd.LF,
		WithoutHeader: false,
	}
}

func (e *Encoder) Encode(fieldList []string, recordSet [][]value.Primary) string {
	fieldLen := len(fieldList)

	lines := make([]string, 0, len(recordSet)+1)

	if !e.WithoutHeader {
		line := make([]string, 0, fieldLen)
		for _, f := range fieldList {
			line = append(line, e.FormatString(f))
		}
		lines = append(lines, strings.Join(line, string(e.Delimiter)))
	}

	for _, record := range recordSet {
		line := make([]string, 0, fieldLen)
		for _, cell := range record {
			line = append(line, e.ConvertCSVCell(cell))
		}
		lines = append(lines, strings.Join(line, string(e.Delimiter)))
	}

	return strings.Join(lines, e.LineBreak.Value())
}

func (e *Encoder) ConvertCSVCell(val value.Primary) string {
	var s string

	switch val.(type) {
	case value.String:
		s = e.FormatString(val.(value.String).Raw())
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
		s = e.FormatString(val.(value.Datetime).Format(time.RFC3339Nano))
	case value.Null:
		s = ""
	}

	return s
}

func (e *Encoder) FormatString(s string) string {
	e.buf.Reset()
	e.buf.WriteRune('"')

	runes := []rune(s)
	pos := 0

	for {
		if len(runes) <= pos {
			break
		}

		r := runes[pos]
		switch r {
		case '"':
			e.buf.WriteRune(r)
			e.buf.WriteRune(r)
		default:
			e.buf.WriteRune(r)
		}

		pos++
	}

	e.buf.WriteRune('"')
	return e.buf.String()
}
