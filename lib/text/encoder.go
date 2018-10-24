package text

import (
	"bytes"
	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/color"
	"github.com/mithrandie/csvq/lib/value"
	"github.com/mithrandie/ternary"
	"strconv"
	"strings"
	"time"
	"unicode"
)

const (
	MarkdownLineBreak = "<br />"
	VLine             = '|'
	HLine             = '-'
	CrossLine         = '+'
	AlignSign         = ':'
	PadChar           = ' '
	LineBreak         = '\n'
	EscapeChar        = '\\'
)

type Align int

const (
	Left   Align = -1
	Center Align = 0
	Right  Align = 1
)

type textField struct {
	lines  []string
	widths []int
	align  Align
	style  color.Style
}

func (tf textField) Width() int {
	w := 0
	for _, v := range tf.widths {
		if w < v {
			w = v
		}
	}
	return w
}

func NewTextField(text string, align Align, style color.Style) *textField {
	lines := strings.Split(text, "\n")
	widths := make([]int, len(lines))

	for i, v := range lines {
		widths[i] = StringWidth(v)
	}

	return &textField{
		lines:  lines,
		widths: widths,
		align:  align,
		style:  style,
	}
}

type Encoder struct {
	Format        cmd.Format
	WithoutHeader bool

	palette *color.Palette

	buf bytes.Buffer
}

func NewEncoder() *Encoder {
	return &Encoder{
		Format:        cmd.TEXT,
		WithoutHeader: false,

		palette: color.NewPalette(),
	}
}

func (e *Encoder) Encode(fieldList []string, recordSet [][]value.Primary) string {
	if len(fieldList) < 1 {
		return color.Warn("Empty Fields")
	}
	if len(recordSet) < 1 {
		return color.Warn("Empty RecordSet")
	}

	withoutHeader := e.WithoutHeader

	switch e.Format {
	case cmd.GFM, cmd.ORG:
		e.palette.Disable()
	default: // cmd.TEXT
		e.palette.Enable()
		withoutHeader = false
	}

	fieldLen := len(fieldList)

	records := make([][]*textField, 0, len(recordSet))
	for _, record := range recordSet {
		r := make([]*textField, 0, fieldLen)
		for _, cell := range record {
			r = append(r, e.ConvertTextCell(cell))
		}
		records = append(records, r)
	}

	fieldWidths := make([]int, fieldLen)

	for _, record := range records {
		for i, f := range record {
			fw := f.Width()
			if fieldWidths[i] < fw {
				fieldWidths[i] = fw
			}
		}
	}

	header := make([]*textField, 0, fieldLen)
	if !withoutHeader {
		for _, field := range fieldList {
			header = append(header, NewTextField(e.Escape(field), Center, color.PlainStyle))
		}
		for i, f := range header {
			fw := f.Width()
			if fieldWidths[i] < fw {
				fieldWidths[i] = f.Width()
			}
			if e.Format == cmd.GFM {
				if fieldWidths[i] < 3 {
					fieldWidths[i] = 3
				}
			}
			if ((fieldWidths[i] - f.Width()) % 2) == 1 {
				fieldWidths[i] = fieldWidths[i] + 1
			}
		}
	}

	lines := make([]string, 0, len(records)+4)

	if e.Format == cmd.TEXT {
		lines = append(lines, e.FormatTextHR(fieldWidths))
	}

	if !withoutHeader {
		lines = append(lines, e.FormatRecord(header, fieldWidths))

		switch e.Format {
		case cmd.GFM:
			lines = append(lines, e.FormatGFMHR(fieldWidths, records[0]))
		case cmd.ORG:
			lines = append(lines, e.FormatOrgHR(fieldWidths))
		default:
			lines = append(lines, e.FormatTextHR(fieldWidths))
		}
	}

	for _, record := range records {
		lines = append(lines, e.FormatRecord(record, fieldWidths))
	}

	if e.Format == cmd.TEXT {
		lines = append(lines, e.FormatTextHR(fieldWidths))
	}

	return strings.Join(lines, string(LineBreak))
}

func (e *Encoder) ConvertTextCell(val value.Primary) *textField {
	var s string

	align := Left
	style := color.PlainStyle

	switch val.(type) {
	case value.Integer, value.Float:
		align = Right
	case value.Boolean, value.Ternary, value.Null:
		align = Center
	}

	switch val.(type) {
	case value.String:
		s = val.(value.String).Raw()
		style = color.StringStyle
	case value.Integer:
		s = val.(value.Integer).String()
		style = color.NumberStyle
	case value.Float:
		s = val.(value.Float).String()
		style = color.NumberStyle
	case value.Boolean:
		s = val.(value.Boolean).String()
		style = color.BooleanStyle
	case value.Ternary:
		t := val.(value.Ternary)
		switch e.Format {
		case cmd.GFM, cmd.ORG:
			if t.Ternary() == ternary.UNKNOWN {
				s = ""
			} else {
				s = strconv.FormatBool(t.Ternary().ParseBool())
			}
		default:
			s = t.Ternary().String()
		}
		style = color.TernaryStyle
	case value.Datetime:
		s = val.(value.Datetime).Format(time.RFC3339Nano)
		style = color.DatetimeStyle
	case value.Null:
		switch e.Format {
		case cmd.GFM, cmd.ORG:
			s = ""
		default:
			s = "NULL"
		}
		style = color.NullStyle
	}

	return NewTextField(e.Escape(s), align, style)
}

func (e *Encoder) FormatRecord(record []*textField, fieldWidths []int) string {
	lineLen := 0
	for _, tf := range record {
		n := len(tf.lines)
		if lineLen < n {
			lineLen = n
		}
	}

	lines := make([]string, 0, lineLen)

	for lineIdx := 0; lineIdx < lineLen; lineIdx++ {
		e.buf.Reset()

		for fieldIdx, tf := range record {
			e.buf.WriteRune(VLine)
			e.buf.WriteRune(PadChar)

			if len(tf.lines) <= lineIdx || len(tf.lines[lineIdx]) < 1 {
				e.buf.Write(bytes.Repeat([]byte(string(PadChar)), fieldWidths[fieldIdx]+1))
				continue
			}

			padLen := fieldWidths[fieldIdx] - tf.widths[lineIdx]
			cellAlign := tf.align
			if cellAlign == Left && unicode.In([]rune(tf.lines[lineIdx])[0], rightToLeftTable) {
				cellAlign = Right
			}

			val := e.palette.Color(tf.lines[lineIdx], tf.style)

			switch cellAlign {
			case Center:
				halfPadLen := padLen / 2
				e.buf.Write(bytes.Repeat([]byte(string(PadChar)), halfPadLen))
				e.buf.WriteString(val)
				e.buf.Write(bytes.Repeat([]byte(string(PadChar)), (padLen-halfPadLen)+1))
			case Right:
				e.buf.Write(bytes.Repeat([]byte(string(PadChar)), padLen))
				e.buf.WriteString(val)
				e.buf.WriteRune(PadChar)
			default: // Left
				e.buf.WriteString(val)
				e.buf.Write(bytes.Repeat([]byte(string(PadChar)), padLen+1))
			}
		}
		e.buf.WriteRune(VLine)
		lines = append(lines, e.buf.String())
	}

	return strings.Join(lines, string(LineBreak))
}

func (e *Encoder) FormatGFMHR(widths []int, record []*textField) string {
	e.buf.Reset()

	for i, w := range widths {
		e.buf.WriteRune(VLine)
		e.buf.WriteRune(PadChar)
		if record[i].align == Right {
			e.buf.Write(bytes.Repeat([]byte(string(HLine)), w-1))
			e.buf.WriteRune(AlignSign)
		} else {
			e.buf.Write(bytes.Repeat([]byte(string(HLine)), w))
		}
		e.buf.WriteRune(PadChar)
	}
	e.buf.WriteRune(VLine)
	return e.buf.String()
}

func (e *Encoder) FormatOrgHR(widths []int) string {
	e.buf.Reset()

	e.buf.WriteRune(VLine)
	for i, w := range widths {
		if 0 < i {
			e.buf.WriteRune(CrossLine)
		}
		e.buf.Write(bytes.Repeat([]byte(string(HLine)), w+2))
	}
	e.buf.WriteRune(VLine)
	return e.buf.String()
}

func (e *Encoder) FormatTextHR(widths []int) string {
	e.buf.Reset()

	for _, w := range widths {
		e.buf.WriteRune(CrossLine)
		e.buf.Write(bytes.Repeat([]byte(string(HLine)), w+2))
	}
	e.buf.WriteRune(CrossLine)
	return e.buf.String()
}

func (e *Encoder) Escape(s string) string {
	e.buf.Reset()

	runes := []rune(s)
	pos := 0

	for {
		if len(runes) <= pos {
			break
		}

		r := runes[pos]
		switch r {
		case '\r':
			if (pos+1) < len(runes) && runes[pos+1] == '\n' {
				pos++
			}
			fallthrough
		case '\n':
			switch e.Format {
			case cmd.GFM, cmd.ORG:
				e.buf.WriteString(MarkdownLineBreak)
			default:
				e.buf.WriteRune(LineBreak)
			}
		case VLine:
			switch e.Format {
			case cmd.GFM, cmd.ORG:
				e.buf.WriteRune(EscapeChar)
			}
			e.buf.WriteRune(r)
		default:
			e.buf.WriteRune(r)
		}

		pos++
	}
	return e.buf.String()
}
