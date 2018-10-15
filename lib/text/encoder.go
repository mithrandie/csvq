package text

import (
	"bytes"
	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/color"
	"github.com/mithrandie/csvq/lib/value"
	"strings"
	"time"
	"unicode"
)

const (
	MarkdownLineBreak = "<br />"
	VLine             = '|'
	HLine             = '-'
	CrossLine         = '+'
	AlighSign         = ':'
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

var fullWidthTable = &unicode.RangeTable{
	R16: []unicode.Range16{
		{0x1100, 0x11ff, 1}, //Hangul Jamo
		{0x2460, 0x24ff, 1}, //Enclosed Alphanumerics
		{0x2500, 0x257f, 1}, //Box Drawing
		{0x2e80, 0x2eff, 1}, //CJK Radicals Supplement
		{0x2f00, 0x2fdf, 1}, //CJK Radicals
		{0x2ff0, 0x2fff, 1}, //Ideographic Description Characters
		{0x3000, 0x303e, 1}, //CJK Symbols and Punctuation
		{0x3040, 0x309f, 1}, //Hiragana
		{0x30a0, 0x30ff, 1}, //Katakana
		{0x3100, 0x312f, 1}, //Bopomofo
		{0x3130, 0x318f, 1}, //Hangul Compatibility Jamo
		{0x3190, 0x319f, 1}, //Ideographic Annotations
		{0x31a0, 0x31bf, 1}, //Bopomofo Extended
		{0x31c0, 0x31ef, 1}, //CJK Strokes
		{0x31f0, 0x31ff, 1}, //Phonetic extensions for Ainu
		{0x3200, 0x32ff, 1}, //Enclosed CJK Letters and Months
		{0x3300, 0x33ff, 1}, //CJK Compatibility
		{0x3400, 0x4dbf, 1}, //CJK Unified Ideographs Extension A
		{0x4e00, 0x9fff, 1}, //CJK Unified Ideographs
		{0xa960, 0xa97f, 1}, //Hangul Jamo Extended A
		{0xac00, 0xd7af, 1}, //Hangul Syllables
		{0xd7b0, 0xd7ff, 1}, //Hangul Jamo Extended B
		{0xf900, 0xfaff, 1}, //CJK Compatibility Ideographs
		{0xff01, 0xff60, 1}, //FullWidth ASCII variants
		{0xffe0, 0xffe6, 1}, //FullWidth Symbol variants
	},
	R32: []unicode.Range32{
		{0x1b000, 0x1b0ff, 1}, //Historic Kana
		{0x1f100, 0x1f1ff, 1}, //Enclosed Alphanumeric Supplement
		{0x1f200, 0x1f2ff, 1}, //Enclosed Ideographic Supplement
		{0x20000, 0x2a6df, 1}, //CJK Unified Ideographs Extension B
		{0x2a700, 0x2b73f, 1}, //CJK Unified Ideographs Extension C
		{0x2b740, 0x2b81f, 1}, //CJK Unified Ideographs Extension D
		{0x2b820, 0x2ceaf, 1}, //CJK Unified Ideographs Extension E
		{0x2f800, 0x2fa1f, 1}, //CJK Compatibility Ideographs Supplement
		{0xe0100, 0xe01ef, 1}, //Variation Selectors Supplement
	},
}

var zeroWidthTable = &unicode.RangeTable{
	R16: []unicode.Range16{
		{0x0000, 0x001f, 1}, //Controls
		{0x0300, 0x036f, 1}, //Combining Diacritical Marks
		{0x0591, 0x05af, 1}, //Hebrew Cantillation Marks
		{0x05b0, 0x05bd, 1}, //Hebrew Points
		{0x05bf, 0x05bf, 1}, //Hebrew Points
		{0x05c1, 0x05c2, 1}, //Hebrew Points
		{0x05c4, 0x05c5, 1}, //Hebrew Points
		{0x05c7, 0x05c7, 1}, //Hebrew Points
		{0x064b, 0x0652, 1}, //Arabic Tashkil from ISO 8859-6
		{0x0653, 0x065f, 1}, //Arabic Combining Marks
		{0x0670, 0x0670, 1}, //Arabic Tashkil
		{0x08a0, 0x08ff, 1}, //Arabic Extended-A
		{0x2028, 0x202f, 1}, //Format Characters
		{0xfbb2, 0xfbc1, 1}, //Arabic pedagogical symbols
		{0xfeff, 0xfeff, 1}, //Arabic Zero Width No-Break Space
	},
}

var rightToLeftTable = &unicode.RangeTable{
	R16: []unicode.Range16{
		{0x0590, 0x05ff, 1}, //Hebrew
		{0x0600, 0x06ff, 1}, //Arabic
		{0x0700, 0x074f, 1}, //Syriac
		{0x0750, 0x077f, 1}, //Arabic Supplement
		{0x0860, 0x086f, 1}, //Syriac Supplement
		{0x08a0, 0x08ff, 1}, //Arabic Extended-A
		{0x200f, 0x200f, 1}, //Right-To-Left Mark
		{0x202b, 0x202b, 1}, //Right-To-Left Embedding
		{0x202e, 0x202e, 1}, //Right-To-Left Override
		{0xfb50, 0xfdff, 1}, //Arabic Presentation Forms-A
		{0xfe70, 0xfeff, 1}, //Arabic Presentation Forms-B
	},
	R32: []unicode.Range32{
		{0x1ee00, 0x1eeff, 1}, //Arabic Mathematical Alphabetic Symbols
	},
}

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
		s = val.(value.Ternary).Ternary().String()
		style = color.TernaryStyle
	case value.Datetime:
		s = val.(value.Datetime).Format(time.RFC3339Nano)
		style = color.DatetimeStyle
	case value.Null:
		s = "NULL"
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
			e.buf.WriteRune(AlighSign)
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

func StringWidth(s string) int {
	l := 0
	for _, r := range s {
		if unicode.In(r, fullWidthTable) {
			l = l + 2
		} else if unicode.In(r, zeroWidthTable) {
			// Do Nothing
		} else {
			l = l + 1
		}
	}
	return l
}
