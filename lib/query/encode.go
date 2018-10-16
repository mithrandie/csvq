package query

import (
	"errors"
	"fmt"
	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/csv"
	"github.com/mithrandie/csvq/lib/json"
	"github.com/mithrandie/csvq/lib/text"
	"github.com/mithrandie/csvq/lib/value"
	"io/ioutil"
	"strings"
)

func EncodeView(view *View, format cmd.Format, delimiter rune, withoutHeader bool, encoding cmd.Encoding, lineBreak cmd.LineBreak) (string, error) {
	var s string
	var err error

	switch format {
	case cmd.JSON, cmd.JSONH, cmd.JSONA:
		s, err = encodeJson(view, format, lineBreak, cmd.GetFlags().PrettyPrint)
		if err != nil {
			return "", err
		}
	default:
		switch format {
		case cmd.GFM, cmd.ORG, cmd.TEXT:
			s = encodeText(view, format, withoutHeader)
			if lineBreak != cmd.LF {
				s = convertLineBreak(s, lineBreak)
			}
		default: // cmd.CSV, cmd.TSV:
			s = encodeCSV(view, delimiter, lineBreak, withoutHeader)
		}

		if encoding != cmd.UTF8 {
			s, err = encodeCharacterCode(s, encoding)
			if err != nil {
				return "", err
			}
		}
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

func bareValues(view *View) ([]string, [][]value.Primary) {
	header := view.Header.TableColumnNames()
	records := make([][]value.Primary, 0, view.RecordLen())
	for _, record := range view.RecordSet {
		row := make([]value.Primary, 0, view.FieldLen())
		for _, cell := range record {
			row = append(row, cell.Value())
		}
		records = append(records, row)
	}
	return header, records
}

func encodeCSV(view *View, delimiter rune, lineBreak cmd.LineBreak, withoutHeader bool) string {
	header, records := bareValues(view)

	e := csv.NewEncoder()
	e.Delimiter = delimiter
	e.LineBreak = lineBreak
	e.WithoutHeader = withoutHeader
	return e.Encode(header, records)
}

func encodeJson(view *View, format cmd.Format, lineBreak cmd.LineBreak, prettyPrint bool) (string, error) {
	header, records := bareValues(view)

	data, err := json.ConvertTableValueToJsonStructure(header, records)
	if err != nil {
		return "", errors.New(fmt.Sprintf("encoding to json failed: %s", err.Error()))
	}

	e := json.NewEncoder()
	switch format {
	case cmd.JSONH:
		e.EscapeType = json.HexDigits
	case cmd.JSONA:
		e.EscapeType = json.AllWithHexDigits
	}
	e.LineBreak = lineBreak
	e.PrettyPrint = prettyPrint

	s := e.Encode(data, e.PrettyPrint)
	return s, nil
}

func encodeText(view *View, format cmd.Format, withoutHeader bool) string {
	header, records := bareValues(view)

	e := text.NewEncoder()
	e.Format = format
	e.WithoutHeader = withoutHeader
	return e.Encode(header, records)
}
