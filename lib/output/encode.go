package output

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/query"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

type textField struct {
	value string
	sign  int
}

func Encode(result query.Result) (string, error) {
	var s string
	var err error
	flags := cmd.GetFlags()

	switch result.Statement {
	case query.SELECT:
		switch flags.Format {
		case cmd.TEXT:
			s = encodeText(result)
		case cmd.CSV, cmd.TSV:
			s = encodeCSV(result, string(flags.WriteDelimiter), flags.WithoutHeader)
		case cmd.JSON:
			s = encodeJson(result)
		}
	}

	if flags.WriteEncoding != cmd.UTF8 {
		s, err = encodeCharacterCode(s, flags.WriteEncoding)
		if err != nil {
			return "", err
		}
	}
	if flags.LineBreak != cmd.LF {
		s = convertLineBreak(s, flags.LineBreak)
	}

	return s, nil
}

func encodeCharacterCode(str string, enc cmd.Encoding) (string, error) {
	var e *encoding.Encoder

	switch enc {
	case cmd.SJIS:
		e = japanese.ShiftJIS.NewEncoder()
	}

	r := transform.NewReader(strings.NewReader(str), e)
	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func convertLineBreak(str string, lb cmd.LineBreak) string {
	return strings.Replace(str, "\n", lb.Value(), -1)
}

func encodeText(result query.Result) string {
	if result.Count < 1 {
		return "Empty\n"
	}

	view := result.View

	header := make([]textField, view.FieldLen())
	for i := range view.Header {
		header[i] = textField{value: view.Header[i].Label(), sign: -1}
	}

	records := make([][]textField, view.RecordLen())
	for i, record := range view.Records {
		records[i] = make([]textField, view.FieldLen())
		for j, cell := range record {
			records[i][j] = formatTextCell(cell)
		}
	}

	fieldLens := make([]int, len(header))

	for i, f := range header {
		fieldLens[i] = countRunes(f)
	}
	for _, record := range records {
		for i, f := range record {
			flen := countRunes(f)
			if fieldLens[i] < flen {
				fieldLens[i] = flen
			}
		}
	}

	s := make([]string, len(records)+4)
	s[0] = formatHR(fieldLens)

	s[1] = formatRecord(header, fieldLens)

	s[2] = formatHR(fieldLens)

	for i, record := range records {
		s[i+3] = formatRecord(record, fieldLens)
	}

	s[len(s)-1] = formatHR(fieldLens)
	return strings.Join(s, "")
}

func formatHR(lens []int) string {
	s := make([]string, len(lens)+1)
	for i, l := range lens {
		s[i] = "+" + strings.Repeat("-", l+2)
	}
	s[len(s)-1] = "+\n"
	return strings.Join(s, "")
}

func formatRecord(record []textField, fieldLens []int) string {
	row := make([][]string, len(record))
	for i, f := range record {
		row[i] = strings.Split(f.value, "\n")
	}

	lineCount := 0
	for _, lines := range row {
		n := len(lines)
		if lineCount < n {
			lineCount = n
		}
	}

	s := make([]string, lineCount)

	for lineIdx := 0; lineIdx < lineCount; lineIdx++ {
		sl := make([]string, len(row)+1)
		for fieldIdx, lines := range row {
			if lineIdx < len(lines) {
				sl[fieldIdx] = fmt.Sprintf("| %"+strconv.Itoa(record[fieldIdx].sign*fieldLens[fieldIdx])+"s ", lines[lineIdx])
			} else {
				sl[fieldIdx] = fmt.Sprintf("| %"+strconv.Itoa(fieldLens[fieldIdx])+"s ", "")
			}
		}
		sl[len(sl)-1] = "|\n"
		s[lineIdx] = strings.Join(sl, "")
	}

	return strings.Join(s, "")
}

func countRunes(f textField) int {
	i := 0
	lines := strings.Split(f.value, "\n")
	for _, line := range lines {
		count := utf8.RuneCountInString(line)
		if i < count {
			i = count
		}
	}
	return i
}

func formatTextCell(c query.Cell) textField {
	primary := c.Primary()

	var s string
	var sign int

	sign = 1
	switch primary.(type) {
	case parser.String:
		s = strings.TrimSpace(primary.(parser.String).Value())
		sign = -1
	case parser.Integer:
		s = parser.Int64ToStr(primary.(parser.Integer).Value())
	case parser.Float:
		s = parser.Float64ToStr(primary.(parser.Float).Value())
	case parser.Boolean:
		s = strconv.FormatBool(primary.(parser.Boolean).Bool())
	case parser.Ternary:
		s = primary.(parser.Ternary).Ternary().String()
	case parser.Datetime:
		s = primary.(parser.Datetime).Format()
		sign = -1
	case parser.Null:
		s = primary.String()
	}

	return textField{value: s, sign: sign}
}

func encodeCSV(result query.Result, delimiter string, withoutHeader bool) string {
	view := result.View

	var header string
	if !withoutHeader {
		h := make([]string, view.FieldLen())
		for i := range view.Header {
			h[i] = quote(escapeCSVString(view.Header[i].Label()))
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

func formatCSVCell(c query.Cell) string {
	primary := c.Primary()

	var s string

	switch primary.(type) {
	case parser.String:
		s = quote(escapeCSVString(primary.(parser.String).Value()))
	case parser.Integer:
		s = parser.Int64ToStr(primary.(parser.Integer).Value())
	case parser.Float:
		s = parser.Float64ToStr(primary.(parser.Float).Value())
	case parser.Boolean:
		s = strconv.FormatBool(primary.(parser.Boolean).Bool())
	case parser.Ternary:
		s = strconv.FormatBool(primary.(parser.Ternary).Bool())
	case parser.Datetime:
		s = quote(escapeCSVString(primary.(parser.Datetime).Format()))
	case parser.Null:
		s = ""
	}

	return s
}

func escapeCSVString(s string) string {
	return strings.Replace(s, "\"", "\"\"", -1)
}

func encodeJson(result query.Result) string {
	view := result.View
	records := make([]string, view.RecordLen())

	for i, record := range view.Records {
		cells := make([]string, view.FieldLen())
		for j, cell := range record {
			cells[j] = quote(escapeJsonString(view.Header[j].Label())) + ":" + formatJsonCell(cell)
		}
		records[i] = "{" + strings.Join(cells, ",") + "}"
	}

	return "[" + strings.Join(records, ",") + "]"
}

func formatJsonCell(c query.Cell) string {
	primary := c.Primary()

	var s string

	switch primary.(type) {
	case parser.String:
		s = quote(escapeJsonString(primary.(parser.String).Value()))
	case parser.Integer:
		s = parser.Int64ToStr(primary.(parser.Integer).Value())
	case parser.Float:
		s = parser.Float64ToStr(primary.(parser.Float).Value())
	case parser.Boolean:
		s = strconv.FormatBool(primary.(parser.Boolean).Bool())
	case parser.Ternary:
		s = strconv.FormatBool(primary.(parser.Ternary).Bool())
	case parser.Datetime:
		s = quote(escapeJsonString(primary.(parser.Datetime).Format()))
	case parser.Null:
		s = "null"
	}

	return s
}

func escapeJsonString(s string) string {
	s = strings.Replace(s, "\\", "\\\\", -1)
	s = strings.Replace(s, "\"", "\\\"", -1)
	s = strings.Replace(s, "/", "\\/", -1)
	s = strings.Replace(s, "\n", "\\n", -1)
	s = strings.Replace(s, "\r", "\\r", -1)
	s = strings.Replace(s, "\t", "\\t", -1)
	return s
}

func quote(s string) string {
	return "\"" + s + "\""
}
