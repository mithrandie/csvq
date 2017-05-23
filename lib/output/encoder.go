package output

import (
	"fmt"
	"strconv"
	"strings"

	"unicode/utf8"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/query"
)

type field struct {
	value string
	sign  int
}

func Encode(format cmd.Format, result query.Result) string {
	var s string

	switch result.Statement {
	case query.SELECT:
		switch format {
		case cmd.TEXT:
			s = encodeText(result)
		}
	}
	return s
}

func encodeText(result query.Result) string {
	if result.Count < 1 {
		return "Empty\n"
	}

	view := result.View

	header := make([]field, view.FieldLen())
	for i := range view.Header {
		header[i] = field{value: view.Header[i].Label(), sign: -1}
	}

	records := make([][]field, view.RecordLen())
	for i, record := range view.Records {
		records[i] = make([]field, view.FieldLen())
		for j, cell := range record {
			records[i][j] = formatCell(cell)
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

func formatRecord(record []field, fieldLens []int) string {
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

func countRunes(f field) int {
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

func formatCell(c query.Cell) field {
	primary := c.Primary()

	var value string
	var sign int

	sign = 1
	switch primary.(type) {
	case parser.String:
		value = strings.TrimSpace(primary.(parser.String).Value())
		sign = -1
	case parser.Integer:
		value = parser.Int64ToStr(primary.(parser.Integer).Value())
	case parser.Float:
		value = parser.Float64ToStr(primary.(parser.Float).Value())
	case parser.Boolean:
		value = strconv.FormatBool(primary.(parser.Boolean).Bool())
	case parser.Ternary:
		value = primary.(parser.Ternary).Ternary().String()
	case parser.Datetime:
		value = primary.(parser.Datetime).Value().Format("2006-01-02 15:04:05.999999999")
		sign = -1
	case parser.Null:
		value = primary.String()
	}

	return field{value: value, sign: sign}
}
