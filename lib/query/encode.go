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

func EncodeView(ctx context.Context, fp io.Writer, view *View, fileInfo *FileInfo, tx *Transaction) (string, error) {
	switch fileInfo.Format {
	case cmd.FIXED:
		return "", encodeFixedLengthFormat(ctx, fp, view, fileInfo.DelimiterPositions, fileInfo.LineBreak, fileInfo.NoHeader, fileInfo.Encoding, fileInfo.SingleLine)
	case cmd.JSON:
		return "", encodeJson(ctx, fp, view, fileInfo.LineBreak, fileInfo.JsonEscape, fileInfo.PrettyPrint, tx)
	case cmd.LTSV:
		return "", encodeLTSV(ctx, fp, view, fileInfo.LineBreak, fileInfo.Encoding)
	case cmd.GFM, cmd.ORG, cmd.TEXT:
		return encodeText(ctx, fp, view, fileInfo.Format, fileInfo.LineBreak, fileInfo.NoHeader, fileInfo.Encoding, tx)
	case cmd.TSV:
		fileInfo.Delimiter = '\t'
		fallthrough
	default: // cmd.CSV
		return "", encodeCSV(ctx, fp, view, fileInfo.Delimiter, fileInfo.LineBreak, fileInfo.NoHeader, fileInfo.Encoding, fileInfo.EncloseAll)
	}
}

func encodeCSV(ctx context.Context, fp io.Writer, view *View, delimiter rune, lineBreak text.LineBreak, withoutHeader bool, encoding text.Encoding, encloseAll bool) error {
	w, err := csv.NewWriter(fp, lineBreak, encoding)
	if err != nil {
		return err
	}
	w.Delimiter = delimiter

	fields := make([]csv.Field, view.FieldLen())

	if !withoutHeader {
		for i := range view.Header {
			fields[i] = csv.NewField(view.Header[i].Column, encloseAll)
		}
		if err := w.Write(fields); err != nil {
			return err
		}
	}

	for i := range view.RecordSet {
		if ctx.Err() != nil {
			return ConvertContextError(ctx.Err())
		}

		for j := range view.RecordSet[i] {
			str, effect, _ := ConvertFieldContents(view.RecordSet[i][j].Value(), false)
			quote := false
			if encloseAll && (effect == cmd.StringEffect || effect == cmd.DatetimeEffect) {
				quote = true
			}
			fields[j] = csv.NewField(str, quote)
		}
		if err := w.Write(fields); err != nil {
			return err
		}
	}
	return w.Flush()
}

func encodeFixedLengthFormat(ctx context.Context, fp io.Writer, view *View, positions []int, lineBreak text.LineBreak, withoutHeader bool, encoding text.Encoding, singleLine bool) error {
	if positions == nil {
		m := fixedlen.NewMeasure()
		m.Encoding = encoding

		var fieldList [][]fixedlen.Field = nil
		var recordStartPos = 0
		var fieldLen = view.FieldLen()

		if withoutHeader {
			fieldList = make([][]fixedlen.Field, view.RecordLen())
		} else {
			fieldList = make([][]fixedlen.Field, view.RecordLen()+1)
			recordStartPos = 1

			fields := make([]fixedlen.Field, fieldLen)
			for i := range view.Header {
				fields[i] = fixedlen.NewField(view.Header[i].Column, text.NotAligned)
			}
			fieldList[0] = fields
			m.Measure(fields)
		}

		for i := range view.RecordSet {
			if ctx.Err() != nil {
				return ConvertContextError(ctx.Err())
			}

			fields := make([]fixedlen.Field, fieldLen)
			for j := range view.RecordSet[i] {
				str, _, a := ConvertFieldContents(view.RecordSet[i][j].Value(), false)
				fields[j] = fixedlen.NewField(str, a)
			}
			fieldList[i+recordStartPos] = fields
			m.Measure(fields)
		}

		positions = m.GeneratePositions()
		w, err := fixedlen.NewWriter(fp, positions, lineBreak, encoding)
		if err != nil {
			return err
		}
		w.InsertSpace = true
		for i := range fieldList {
			if ctx.Err() != nil {
				return ConvertContextError(ctx.Err())
			}

			if err := w.Write(fieldList[i]); err != nil {
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

		fields := make([]fixedlen.Field, view.FieldLen())

		if !withoutHeader && !singleLine {
			for i := range view.Header {
				fields[i] = fixedlen.NewField(view.Header[i].Column, text.NotAligned)
			}
			if err := w.Write(fields); err != nil {
				return err
			}
		}

		for i := range view.RecordSet {
			if ctx.Err() != nil {
				return ConvertContextError(ctx.Err())
			}

			for j := range view.RecordSet[i] {
				str, _, a := ConvertFieldContents(view.RecordSet[i][j].Value(), false)
				fields[j] = fixedlen.NewField(str, a)
			}
			if err := w.Write(fields); err != nil {
				return err
			}
		}
		if err = w.Flush(); err != nil {
			return err
		}
	}
	return nil
}

func encodeJson(ctx context.Context, fp io.Writer, view *View, lineBreak text.LineBreak, escapeType txjson.EscapeType, prettyPrint bool, tx *Transaction) error {
	header := view.Header.TableColumnNames()
	records := make([][]value.Primary, view.RecordLen())
	for i := range view.RecordSet {
		if ctx.Err() != nil {
			return ConvertContextError(ctx.Err())
		}

		row := make([]value.Primary, view.FieldLen())
		for j := range view.RecordSet[i] {
			row[j] = view.RecordSet[i][j].Value()
		}
		records[i] = row
	}

	data, err := json.ConvertTableValueToJsonStructure(ctx, header, records)
	if err != nil {
		if ctx.Err() != nil {
			return ConvertContextError(ctx.Err())
		}
		return errors.New(fmt.Sprintf("encoding to json failed: %s", err.Error()))
	}

	e := txjson.NewEncoder()
	e.EscapeType = escapeType
	e.LineBreak = lineBreak
	e.PrettyPrint = prettyPrint
	if prettyPrint && tx.Flags.Color {
		e.Palette = tx.Palette
	}
	defer tx.UseColor(tx.Flags.Color)

	s := e.Encode(data)

	w := bufio.NewWriter(fp)
	if _, err := w.WriteString(s); err != nil {
		return err
	}
	return w.Flush()
}

func encodeText(ctx context.Context, fp io.Writer, view *View, format cmd.Format, lineBreak text.LineBreak, withoutHeader bool, encoding text.Encoding, tx *Transaction) (string, error) {
	isPlainTable := false

	var tableFormat = table.PlainTable
	switch format {
	case cmd.GFM:
		tableFormat = table.GFMTable
	case cmd.ORG:
		tableFormat = table.OrgTable
	default:
		if view.FieldLen() < 1 {
			return "Empty Fields", NewEmptyResultSetError()
		}
		if view.RecordLen() < 1 {
			return "Empty RecordSet", NewEmptyResultSetError()
		}
		isPlainTable = true
	}

	e := table.NewEncoder(tableFormat, view.RecordLen())
	e.LineBreak = lineBreak
	e.EastAsianEncoding = tx.Flags.EastAsianEncoding
	e.CountDiacriticalSign = tx.Flags.CountDiacriticalSign
	e.CountFormatCode = tx.Flags.CountFormatCode
	e.WithoutHeader = withoutHeader
	e.Encoding = encoding

	fieldLen := view.FieldLen()

	if !withoutHeader {
		hfields := make([]table.Field, fieldLen)
		for i := range view.Header {
			hfields[i] = table.NewField(view.Header[i].Column, text.Centering)
		}
		e.SetHeader(hfields)
	}

	aligns := make([]text.FieldAlignment, fieldLen)

	var textStrBuf bytes.Buffer
	var textLineBuf bytes.Buffer
	for i := range view.RecordSet {
		if ctx.Err() != nil {
			return "", ConvertContextError(ctx.Err())
		}

		rfields := make([]table.Field, fieldLen)
		for j := range view.RecordSet[i] {
			str, effect, align := ConvertFieldContents(view.RecordSet[i][j].Value(), isPlainTable)
			if format == cmd.TEXT {
				textStrBuf.Reset()
				textLineBuf.Reset()

				runes := []rune(str)
				pos := 0
				for {
					if len(runes) <= pos {
						if 0 < textLineBuf.Len() {
							textStrBuf.WriteString(tx.Palette.Render(effect, textLineBuf.String()))
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
							textStrBuf.WriteString(tx.Palette.Render(effect, textLineBuf.String()))
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
			rfields[j] = table.NewField(str, align)

			if i == 0 {
				aligns[j] = align
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
	hfields := make([]string, view.FieldLen())
	for i := range view.Header {
		hfields[i] = view.Header[i].Column
	}

	w, err := ltsv.NewWriter(fp, hfields, lineBreak, encoding)
	if err != nil {
		return err
	}

	fields := make([]string, view.FieldLen())
	for i := range view.RecordSet {
		if ctx.Err() != nil {
			return ConvertContextError(ctx.Err())
		}

		for j := range view.RecordSet[i] {
			fields[j], _, _ = ConvertFieldContents(view.RecordSet[i][j].Value(), false)
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
	case *value.String:
		s = val.(*value.String).Raw()
		effect = cmd.StringEffect
	case *value.Integer:
		s = val.(*value.Integer).String()
		effect = cmd.NumberEffect
		align = text.RightAligned
	case *value.Float:
		s = val.(*value.Float).String()
		effect = cmd.NumberEffect
		align = text.RightAligned
	case *value.Boolean:
		s = val.(*value.Boolean).String()
		effect = cmd.BooleanEffect
		align = text.Centering
	case *value.Ternary:
		t := val.(*value.Ternary)
		if forTextTable {
			s = t.Ternary().String()
			effect = cmd.TernaryEffect
			align = text.Centering
		} else if t.Ternary() != ternary.UNKNOWN {
			s = strconv.FormatBool(t.Ternary().ParseBool())
			effect = cmd.BooleanEffect
			align = text.Centering
		}
	case *value.Datetime:
		s = val.(*value.Datetime).Format(time.RFC3339Nano)
		effect = cmd.DatetimeEffect
	case *value.Null:
		if forTextTable {
			s = "NULL"
			effect = cmd.NullEffect
			align = text.Centering
		}
	}

	return s, effect, align
}
