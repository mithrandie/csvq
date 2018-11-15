package query

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/json"
	"github.com/mithrandie/csvq/lib/value"

	"github.com/mithrandie/go-text"
	"github.com/mithrandie/go-text/color"
	"github.com/mithrandie/go-text/csv"
	"github.com/mithrandie/go-text/fixedlen"
	txjson "github.com/mithrandie/go-text/json"
	"github.com/mithrandie/go-text/table"
	"github.com/mithrandie/ternary"
)

func EncodeView(fp io.Writer, view *View, fileInfo *FileInfo) error {
	switch fileInfo.Format {
	case cmd.FIXED:
		return encodeFixedLengthFormat(fp, view, fileInfo.DelimiterPositions, fileInfo.LineBreak, fileInfo.NoHeader, fileInfo.Encoding)
	case cmd.JSON, cmd.JSONH, cmd.JSONA:
		return encodeJson(fp, view, fileInfo.Format, fileInfo.LineBreak, fileInfo.PrettyPrint)
	case cmd.GFM, cmd.ORG, cmd.TEXT:
		return encodeText(fp, view, fileInfo.Format, fileInfo.LineBreak, fileInfo.NoHeader, fileInfo.Encoding)
	case cmd.TSV:
		fileInfo.Delimiter = '\t'
		fallthrough
	default: // cmd.CSV
		return encodeCSV(fp, view, fileInfo.Delimiter, fileInfo.LineBreak, fileInfo.NoHeader, fileInfo.Encoding, fileInfo.EncloseAll)
	}
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

func encodeCSV(fp io.Writer, view *View, delimiter rune, lineBreak text.LineBreak, withoutHeader bool, encoding text.Encoding, encloseAll bool) error {
	header, records := bareValues(view)

	w := csv.NewWriter(fp, lineBreak, encoding)
	w.Delimiter = delimiter

	fields := make([]csv.Field, len(header))

	if !withoutHeader {
		for i, v := range header {
			fields[i] = csv.NewField(v, encloseAll)
		}
		if err := w.Write(fields); err != nil {
			return err
		}
	}

	for _, record := range records {
		for i, v := range record {
			str, e, _ := ConvertFieldContents(v, false)
			quote := false
			if encloseAll && (e == cmd.StringEffect || e == cmd.DatetimeEffect) {
				quote = true
			}
			fields[i] = csv.NewField(str, quote)
		}
		if err := w.Write(fields); err != nil {
			return err
		}
	}
	w.Flush()
	return nil
}

func encodeFixedLengthFormat(fp io.Writer, view *View, positions []int, lineBreak text.LineBreak, withoutHeader bool, encoding text.Encoding) error {
	header, records := bareValues(view)

	if positions == nil {
		m := fixedlen.NewMeasure()
		m.Encoding = encoding

		fieldList := make([][]fixedlen.Field, 0, len(records)+1)
		if !withoutHeader {
			fields := make([]fixedlen.Field, 0, len(header))
			for _, v := range header {
				fields = append(fields, fixedlen.NewField(v, text.NotAligned))
			}
			fieldList = append(fieldList, fields)
			m.Measure(fields)
		}

		for _, record := range records {
			fields := make([]fixedlen.Field, 0, len(record))
			for _, v := range record {
				str, _, a := ConvertFieldContents(v, false)
				fields = append(fields, fixedlen.NewField(str, a))
			}
			fieldList = append(fieldList, fields)
			m.Measure(fields)
		}

		positions = m.GeneratePositions()
		w := fixedlen.NewWriter(fp, positions, lineBreak, encoding)
		w.InsertSpace = true
		for _, fields := range fieldList {
			if err := w.Write(fields); err != nil {
				return err
			}
		}
		w.Flush()

	} else {
		w := fixedlen.NewWriter(fp, positions, lineBreak, encoding)

		fields := make([]fixedlen.Field, len(header))

		if !withoutHeader {
			for i, v := range header {
				fields[i] = fixedlen.NewField(v, text.NotAligned)
			}
			if err := w.Write(fields); err != nil {
				return err
			}
		}

		for _, record := range records {
			for i, v := range record {
				str, _, a := ConvertFieldContents(v, false)
				fields[i] = fixedlen.NewField(str, a)
			}
			if err := w.Write(fields); err != nil {
				return err
			}
		}
		w.Flush()
	}
	return nil
}

func encodeJson(fp io.Writer, view *View, format cmd.Format, lineBreak text.LineBreak, prettyPrint bool) error {
	header, records := bareValues(view)

	data, err := json.ConvertTableValueToJsonStructure(header, records)
	if err != nil {
		return errors.New(fmt.Sprintf("encoding to json failed: %s", err.Error()))
	}

	e := txjson.NewEncoder()
	switch format {
	case cmd.JSONH:
		e.EscapeType = txjson.HexDigits
	case cmd.JSONA:
		e.EscapeType = txjson.AllWithHexDigits
	}
	e.LineBreak = lineBreak
	e.PrettyPrint = prettyPrint
	if prettyPrint && cmd.GetFlags().Color {
		e.Palette = cmd.GetPalette()
	}

	s := e.Encode(data)
	cmd.GetPalette().Enable()

	w := bufio.NewWriter(fp)
	if _, err := w.WriteString(s); err != nil {
		return err
	}
	return w.Flush()
}

func encodeText(fp io.Writer, view *View, format cmd.Format, lineBreak text.LineBreak, withoutHeader bool, encoding text.Encoding) error {
	header, records := bareValues(view)

	isPlainTable := false

	var tableFormat = table.PlainTable
	switch format {
	case cmd.GFM:
		tableFormat = table.GFMTable
	case cmd.ORG:
		tableFormat = table.OrgTable
	default:
		if len(header) < 1 {
			w := bufio.NewWriter(fp)
			if _, err := w.WriteString(color.Warn("Empty Fields")); err != nil {
				return err
			}
			return w.Flush()
		}
		if len(records) < 1 {
			w := bufio.NewWriter(fp)
			if _, err := w.WriteString(color.Warn("Empty RecordSet")); err != nil {
				return err
			}
			return w.Flush()
		}
		isPlainTable = true
	}

	e := table.NewEncoder(tableFormat, len(records))
	e.LineBreak = lineBreak
	e.EastAsianEncoding = cmd.GetFlags().EastAsianEncoding
	e.CountDiacriticalSign = cmd.GetFlags().CountDiacriticalSign
	e.CountFormatCode = cmd.GetFlags().CountFormatCode
	e.WithoutHeader = withoutHeader
	e.Encoding = encoding

	palette := cmd.GetPalette()

	if !withoutHeader {
		hfields := make([]table.Field, 0, len(header))
		for _, v := range header {
			hfields = append(hfields, table.NewField(v, text.Centering))
		}
		e.SetHeader(hfields)
	}

	aligns := make([]text.FieldAlignment, 0, len(header))
	for i, record := range records {
		rfields := make([]table.Field, 0, len(header))
		for _, v := range record {
			str, effect, align := ConvertFieldContents(v, isPlainTable)
			if format == cmd.TEXT {
				str = palette.Render(effect, str)
			}
			rfields = append(rfields, table.NewField(str, align))

			if i == 0 {
				aligns = append(aligns, align)
			}
		}
		e.AppendRecord(rfields)
	}

	if format == cmd.GFM {
		e.SetFieldAlignments(aligns)
	}

	s, err := e.Encode()
	if err != nil {
		return err
	}
	w := bufio.NewWriter(fp)
	if _, err := w.WriteString(s); err != nil {
		return err
	}
	return w.Flush()
}

func ConvertFieldContents(val value.Primary, forTextTable bool) (string, string, text.FieldAlignment) {
	var s string
	var effect = cmd.NoEffect
	var align = text.NotAligned

	switch val.(type) {
	case value.String:
		s = val.(value.String).Raw()
		effect = cmd.StringEffect
	case value.Integer:
		s = val.(value.Integer).String()
		effect = cmd.NumberEffect
		align = text.RightAligned
	case value.Float:
		s = val.(value.Float).String()
		effect = cmd.NumberEffect
		align = text.RightAligned
	case value.Boolean:
		s = val.(value.Boolean).String()
		effect = cmd.BooleanEffect
		align = text.Centering
	case value.Ternary:
		t := val.(value.Ternary)
		if forTextTable {
			s = t.Ternary().String()
			effect = cmd.TernaryEffect
			align = text.Centering
		} else if t.Ternary() != ternary.UNKNOWN {
			s = strconv.FormatBool(t.Ternary().ParseBool())
			effect = cmd.BooleanEffect
			align = text.Centering
		}
	case value.Datetime:
		s = val.(value.Datetime).Format(time.RFC3339Nano)
		effect = cmd.DatetimeEffect
	case value.Null:
		if forTextTable {
			s = "NULL"
			effect = cmd.NullEffect
			align = text.Centering
		}
	}

	return s, effect, align
}
