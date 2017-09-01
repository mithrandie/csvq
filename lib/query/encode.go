package query

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/ternary"
	"github.com/mithrandie/csvq/lib/value"
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
	values []string
	widths []int
	sign   int
}

func (tf textField) width() int {
	w := 0
	for _, v := range tf.widths {
		if w < v {
			w = v
		}
	}
	return w
}

func NewTextField(s string, sign int) textField {
	values := strings.Split(s, "\n")
	widths := make([]int, len(values))

	for i, v := range values {
		widths[i] = stringWidth(v)
	}

	return textField{
		values: values,
		widths: widths,
		sign:   sign,
	}
}

func EncodeView(view *View, format cmd.Format, delimiter rune, withoutHeader bool, encoding cmd.Encoding, lineBreak cmd.LineBreak) (string, error) {
	var s string
	var err error

	switch format {
	case cmd.CSV, cmd.TSV:
		s = encodeCSV(view, string(delimiter), withoutHeader)
	case cmd.JSON:
		s = encodeJson(view)
	default:
		s = encodeText(view)
	}

	if encoding != cmd.UTF8 {
		s, err = encodeCharacterCode(s, encoding)
		if err != nil {
			return "", err
		}
	}
	if lineBreak != cmd.LF {
		s = convertLineBreak(s, lineBreak)
	}

	return s, nil
}

func encodeCharacterCode(str string, enc cmd.Encoding) (string, error) {
	r := cmd.GetReader(strings.NewReader(str), enc)
	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func convertLineBreak(str string, lb cmd.LineBreak) string {
	return strings.Replace(str, "\n", lb.Value(), -1)
}

func encodeText(view *View) string {
	if view.FieldLen() < 1 {
		return "Empty Fields"
	}
	if view.RecordLen() < 1 {
		return "Empty Records"
	}

	header := make([]textField, view.FieldLen())
	for i := range view.Header {
		header[i] = NewTextField(view.Header[i].Column, -1)
	}

	records := make([][]textField, view.RecordLen())
	for i, record := range view.Records {
		records[i] = make([]textField, view.FieldLen())
		for j, cell := range record {
			records[i][j] = formatTextCell(cell)
		}
	}

	fieldWidths := make([]int, len(header))

	for i, f := range header {
		fieldWidths[i] = f.width()
	}
	for _, record := range records {
		for i, f := range record {
			flen := f.width()
			if fieldWidths[i] < flen {
				fieldWidths[i] = flen
			}
		}
	}

	s := make([]string, len(records)+4)
	s[0] = formatHR(fieldWidths)

	s[1] = formatRecord(header, fieldWidths)

	s[2] = formatHR(fieldWidths)

	for i, record := range records {
		s[i+3] = formatRecord(record, fieldWidths)
	}

	s[len(s)-1] = formatHR(fieldWidths)
	return strings.Join(s, "\n")
}

func formatHR(lens []int) string {
	s := make([]string, len(lens)+1)
	for i, l := range lens {
		s[i] = "+" + strings.Repeat("-", l+2)
	}
	s[len(s)-1] = "+"
	return strings.Join(s, "")
}

func formatRecord(record []textField, fieldWidths []int) string {
	lineCount := 0
	for _, tf := range record {
		n := len(tf.values)
		if lineCount < n {
			lineCount = n
		}
	}

	s := make([]string, lineCount)

	for lineIdx := 0; lineIdx < lineCount; lineIdx++ {
		sl := make([]string, len(record)+1)
		for fieldIdx, tf := range record {
			var value string
			if lineIdx < len(tf.values) {
				pad := strings.Repeat(" ", fieldWidths[fieldIdx]-tf.widths[lineIdx])
				if tf.sign < 0 {
					runes := []rune(tf.values[lineIdx])
					if 0 < len(runes) && unicode.In(runes[0], rightToLeftTable) {
						value = pad + tf.values[lineIdx]
					} else {
						value = tf.values[lineIdx] + pad
					}
				} else {
					value = pad + tf.values[lineIdx]
				}
			} else {
				value = strings.Repeat(" ", fieldWidths[fieldIdx])
			}
			sl[fieldIdx] = fmt.Sprintf("| %s ", value)
		}
		sl[len(sl)-1] = "|"
		s[lineIdx] = strings.Join(sl, "")
	}

	return strings.Join(s, "\n")
}

func stringWidth(s string) int {
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

func formatTextCell(c Cell) textField {
	primary := c.Primary()

	var s string
	var sign int

	sign = 1
	switch primary.(type) {
	case value.String:
		s = primary.(value.String).Raw()
		sign = -1
	case value.Integer:
		s = primary.(value.Integer).String()
	case value.Float:
		s = primary.(value.Float).String()
	case value.Boolean:
		s = primary.(value.Boolean).String()
	case value.Ternary:
		s = primary.(value.Ternary).Ternary().String()
	case value.Datetime:
		s = primary.(value.Datetime).Format(time.RFC3339Nano)
		sign = -1
	case value.Null:
		s = "NULL"
	}

	return NewTextField(s, sign)
}

func encodeCSV(view *View, delimiter string, withoutHeader bool) string {
	var header string
	if !withoutHeader {
		h := make([]string, view.FieldLen())
		for i := range view.Header {
			h[i] = quote(escapeCSVString(view.Header[i].Column))
		}
		header = strings.Join(h, delimiter)
	}

	records := make([]string, view.RecordLen())
	for i, record := range view.Records {
		cells := make([]string, view.FieldLen())
		for j, cell := range record {
			cells[j] = formatCSVCell(cell)
		}
		records[i] = strings.Join(cells, delimiter)
	}

	s := strings.Join(records, "\n")
	if !withoutHeader {
		s = header + "\n" + s
	}
	return s
}

func formatCSVCell(c Cell) string {
	primary := c.Primary()

	var s string

	switch primary.(type) {
	case value.String:
		s = quote(escapeCSVString(primary.(value.String).Raw()))
	case value.Integer:
		s = primary.(value.Integer).String()
	case value.Float:
		s = primary.(value.Float).String()
	case value.Boolean:
		s = primary.(value.Boolean).String()
	case value.Ternary:
		t := primary.(value.Ternary)
		if t.Ternary() == ternary.UNKNOWN {
			s = ""
		} else {
			s = strconv.FormatBool(t.Ternary().BoolValue())
		}
	case value.Datetime:
		s = quote(escapeCSVString(primary.(value.Datetime).Format(time.RFC3339Nano)))
	case value.Null:
		s = ""
	}

	return s
}

func escapeCSVString(s string) string {
	return strings.Replace(s, "\"", "\"\"", -1)
}

func encodeJson(view *View) string {
	records := make([]string, view.RecordLen())

	for i, record := range view.Records {
		cells := make([]string, view.FieldLen())
		for j, cell := range record {
			cells[j] = quote(escapeJsonString(view.Header[j].Column)) + ":" + formatJsonCell(cell)
		}
		records[i] = "{" + strings.Join(cells, ",") + "}"
	}

	return "[" + strings.Join(records, ",") + "]"
}

func formatJsonCell(c Cell) string {
	primary := c.Primary()

	var s string

	switch primary.(type) {
	case value.String:
		s = quote(escapeJsonString(primary.(value.String).Raw()))
	case value.Integer:
		s = primary.(value.Integer).String()
	case value.Float:
		s = primary.(value.Float).String()
	case value.Boolean:
		s = primary.(value.Boolean).String()
	case value.Ternary:
		t := primary.(value.Ternary)
		if t.Ternary() == ternary.UNKNOWN {
			s = "null"
		} else {
			s = strconv.FormatBool(t.Ternary().BoolValue())
		}
	case value.Datetime:
		s = quote(escapeJsonString(primary.(value.Datetime).Format(time.RFC3339Nano)))
	case value.Null:
		s = "null"
	}

	return s
}

func escapeJsonString(s string) string {
	runes := []rune(s)
	encoded := []rune{}

	for _, r := range runes {
		switch r {
		case '\\', '"', '/':
			encoded = append(encoded, '\\', r)
		case '\n':
			encoded = append(encoded, '\\', 'n')
		case '\r':
			encoded = append(encoded, '\\', 'r')
		case '\t':
			encoded = append(encoded, '\\', 't')
		default:
			encoded = append(encoded, r)
		}
	}
	return string(encoded)
}

func quote(s string) string {
	return "\"" + s + "\""
}
