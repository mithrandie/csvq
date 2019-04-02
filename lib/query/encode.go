package query

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/json"
	"github.com/mithrandie/csvq/lib/value"

	"github.com/mithrandie/go-text"
	"github.com/mithrandie/go-text/csv"
	"github.com/mithrandie/go-text/fixedlen"
	txjson "github.com/mithrandie/go-text/json"
	"github.com/mithrandie/go-text/ltsv"
	"github.com/mithrandie/go-text/table"
	"github.com/mithrandie/ternary"
)

type EmptyResultSetError struct{}

func (e EmptyResultSetError) Error() string {
	return "empty result set"
}

func NewEmptyResultSetError() *EmptyResultSetError {
	return &EmptyResultSetError{}
}

func EncodeView(ctx context.Context, fp io.Writer, view *View, fileInfo *FileInfo, flags *cmd.Flags) (string, error) {
	switch fileInfo.Format {
	case cmd.FIXED:
		return "", encodeFixedLengthFormat(ctx, fp, view, fileInfo.DelimiterPositions, fileInfo.LineBreak, fileInfo.NoHeader, fileInfo.Encoding, fileInfo.SingleLine)
	case cmd.JSON:
		return "", encodeJson(ctx, fp, view, fileInfo.LineBreak, fileInfo.JsonEscape, fileInfo.PrettyPrint, flags)
	case cmd.LTSV:
		return "", encodeLTSV(ctx, fp, view, fileInfo.LineBreak, fileInfo.Encoding)
	case cmd.GFM, cmd.ORG, cmd.TEXT:
		return encodeText(ctx, fp, view, fileInfo.Format, fileInfo.LineBreak, fileInfo.NoHeader, fileInfo.Encoding, flags)
	case cmd.TSV:
		fileInfo.Delimiter = '\t'
		fallthrough
	default: // cmd.CSV
		return "", encodeCSV(ctx, fp, view, fileInfo.Delimiter, fileInfo.LineBreak, fileInfo.NoHeader, fileInfo.Encoding, fileInfo.EncloseAll)
	}
}

func bareValues(ctx context.Context, view *View) ([]string, [][]value.Primary, error) {
	header := view.Header.TableColumnNames()
	records := make([][]value.Primary, 0, view.RecordLen())
	for _, record := range view.RecordSet {
		if ctx.Err() != nil {
			return nil, nil, NewContextIsDone(ctx.Err().Error())
		}

		row := make([]value.Primary, 0, view.FieldLen())
		for _, cell := range record {
			row = append(row, cell.Value())
		}
		records = append(records, row)
	}
	return header, records, nil
}

func encodeCSV(ctx context.Context, fp io.Writer, view *View, delimiter rune, lineBreak text.LineBreak, withoutHeader bool, encoding text.Encoding, encloseAll bool) error {
	header, records, err := bareValues(ctx, view)
	if err != nil {
		return err
	}

	w, err := csv.NewWriter(fp, lineBreak, encoding)
	if err != nil {
		return err
	}
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
		if ctx.Err() != nil {
			return NewContextIsDone(ctx.Err().Error())
		}

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
	return w.Flush()
}

func encodeFixedLengthFormat(ctx context.Context, fp io.Writer, view *View, positions []int, lineBreak text.LineBreak, withoutHeader bool, encoding text.Encoding, singleLine bool) error {
	header, records, err := bareValues(ctx, view)
	if err != nil {
		return err
	}

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
			if ctx.Err() != nil {
				return NewContextIsDone(ctx.Err().Error())
			}

			fields := make([]fixedlen.Field, 0, len(record))
			for _, v := range record {
				str, _, a := ConvertFieldContents(v, false)
				fields = append(fields, fixedlen.NewField(str, a))
			}
			fieldList = append(fieldList, fields)
			m.Measure(fields)
		}

		positions = m.GeneratePositions()
		w, err := fixedlen.NewWriter(fp, positions, lineBreak, encoding)
		if err != nil {
			return err
		}
		w.InsertSpace = true
		for _, fields := range fieldList {
			if ctx.Err() != nil {
				return NewContextIsDone(ctx.Err().Error())
			}

			if err := w.Write(fields); err != nil {
				return err
			}
		}
		err = w.Flush()

	} else {
		w, err := fixedlen.NewWriter(fp, positions, lineBreak, encoding)
		if err != nil {
			return err
		}
		w.SingleLine = singleLine

		fields := make([]fixedlen.Field, len(header))

		if !withoutHeader && !singleLine {
			for i, v := range header {
				fields[i] = fixedlen.NewField(v, text.NotAligned)
			}
			if err := w.Write(fields); err != nil {
				return err
			}
		}

		for _, record := range records {
			if ctx.Err() != nil {
				return NewContextIsDone(ctx.Err().Error())
			}

			for i, v := range record {
				str, _, a := ConvertFieldContents(v, false)
				fields[i] = fixedlen.NewField(str, a)
			}
			if err := w.Write(fields); err != nil {
				return err
			}
		}
		err = w.Flush()
	}
	return err
}

func encodeJson(ctx context.Context, fp io.Writer, view *View, lineBreak text.LineBreak, escapeType txjson.EscapeType, prettyPrint bool, flags *cmd.Flags) error {
	header, records, err := bareValues(ctx, view)
	if err != nil {
		return err
	}

	data, err := json.ConvertTableValueToJsonStructure(header, records)
	if err != nil {
		return errors.New(fmt.Sprintf("encoding to json failed: %s", err.Error()))
	}

	e := txjson.NewEncoder()
	e.EscapeType = escapeType
	e.LineBreak = lineBreak
	e.PrettyPrint = prettyPrint
	if prettyPrint && flags.Color {
		e.Palette = cmd.GetPalette()
	}

	s := e.Encode(data)
	if e.Palette != nil {
		e.Palette.Enable()
	}

	w := bufio.NewWriter(fp)
	if _, err := w.WriteString(s); err != nil {
		return err
	}
	return w.Flush()
}

func encodeText(ctx context.Context, fp io.Writer, view *View, format cmd.Format, lineBreak text.LineBreak, withoutHeader bool, encoding text.Encoding, flags *cmd.Flags) (string, error) {
	header, records, err := bareValues(ctx, view)
	if err != nil {
		return "", err
	}

	isPlainTable := false

	var tableFormat = table.PlainTable
	switch format {
	case cmd.GFM:
		tableFormat = table.GFMTable
	case cmd.ORG:
		tableFormat = table.OrgTable
	default:
		if len(header) < 1 {
			return "Empty Fields", NewEmptyResultSetError()
		}
		if len(records) < 1 {
			return "Empty RecordSet", NewEmptyResultSetError()
		}
		isPlainTable = true
	}

	e := table.NewEncoder(tableFormat, len(records))
	e.LineBreak = lineBreak
	e.EastAsianEncoding = flags.EastAsianEncoding
	e.CountDiacriticalSign = flags.CountDiacriticalSign
	e.CountFormatCode = flags.CountFormatCode
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

	var textStrBuf bytes.Buffer
	var textLineBuf bytes.Buffer
	for i, record := range records {
		if ctx.Err() != nil {
			return "", NewContextIsDone(ctx.Err().Error())
		}

		rfields := make([]table.Field, 0, len(header))
		for _, v := range record {
			str, effect, align := ConvertFieldContents(v, isPlainTable)
			if format == cmd.TEXT {
				textStrBuf.Reset()
				textLineBuf.Reset()

				runes := []rune(str)
				pos := 0
				for {
					if len(runes) <= pos {
						if 0 < textLineBuf.Len() {
							textStrBuf.WriteString(palette.Render(effect, textLineBuf.String()))
						}
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
						if 0 < textLineBuf.Len() {
							textStrBuf.WriteString(palette.Render(effect, textLineBuf.String()))
						}
						textStrBuf.WriteByte('\n')
						textLineBuf.Reset()
					default:
						textLineBuf.WriteRune(r)
					}

					pos++
				}
				str = textStrBuf.String()
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
		return "", err
	}
	w := bufio.NewWriter(fp)
	if _, err := w.WriteString(s); err != nil {
		return "", err
	}
	return "", w.Flush()
}

func encodeLTSV(ctx context.Context, fp io.Writer, view *View, lineBreak text.LineBreak, encoding text.Encoding) error {
	header, records, err := bareValues(ctx, view)
	if err != nil {
		return err
	}

	w, err := ltsv.NewWriter(fp, header, lineBreak, encoding)
	if err != nil {
		return err
	}

	fields := make([]string, len(header))
	for _, record := range records {
		if ctx.Err() != nil {
			return NewContextIsDone(ctx.Err().Error())
		}

		for i, v := range record {
			fields[i], _, _ = ConvertFieldContents(v, false)
		}
		if err := w.Write(fields); err != nil {
			return err
		}
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
