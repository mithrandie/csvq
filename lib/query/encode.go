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

func EncodeView(view *View, fileInfo *FileInfo) (string, error) {
	var s string
	var err error

	switch fileInfo.Format {
	case cmd.FIXED:
		s, err = encodeFixedLengthFormat(view, fileInfo.DelimiterPositions, fileInfo.LineBreak, fileInfo.NoHeader, fileInfo.Encoding)
		if err != nil {
			return "", err
		}
	case cmd.JSON, cmd.JSONH, cmd.JSONA:
		s, err = encodeJson(view, fileInfo.Format, fileInfo.LineBreak, fileInfo.PrettyPrint)
		if err != nil {
			return "", err
		}
	default:
		switch fileInfo.Format {
		case cmd.GFM, cmd.ORG, cmd.TEXT:
			s = encodeText(view, fileInfo.Format, fileInfo.NoHeader)
			if fileInfo.LineBreak != cmd.LF {
				s = convertLineBreak(s, fileInfo.LineBreak)
			}
		case cmd.TSV:
			fileInfo.Delimiter = '\t'
			fallthrough
		default: // cmd.CSV
			s = encodeCSV(view, fileInfo.Delimiter, fileInfo.LineBreak, fileInfo.NoHeader)
		}
	}

	switch fileInfo.Format {
	case cmd.JSON, cmd.JSONH, cmd.JSONA:
		//Do Nothing
	default:
		if fileInfo.Encoding != cmd.UTF8 {
			s, err = encodeCharacterCode(s, fileInfo.Encoding)
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

func encodeFixedLengthFormat(view *View, positions []int, lineBreak cmd.LineBreak, withoutHeader bool, encoding cmd.Encoding) (string, error) {
	header, records := bareValues(view)

	e := text.NewFixedLengthEncoder()
	e.DelimiterPositions = positions
	e.LineBreak = lineBreak
	e.WithoutHeader = withoutHeader
	e.Encoding = encoding
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
