package query

import (
	"errors"
	"fmt"
	"github.com/mithrandie/csvq/lib/color"
	"github.com/mithrandie/csvq/lib/text"
	"github.com/mithrandie/ternary"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/file"
	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/value"
)

const (
	DelimiterFlag      = "@@DELIMITER"
	EncodingFlag       = "@@ENCODING"
	LineBreakFlag      = "@@LINE_BREAK"
	TimezoneFlag       = "@@TIMEZONE"
	RepositoryFlag     = "@@REPOSITORY"
	DatetimeFormatFlag = "@@DATETIME_FORMAT"
	NoHeaderFlag       = "@@NO_HEADER"
	WithoutNullFlag    = "@@WITHOUT_NULL"
	WaitTimeoutFlag    = "@@WAIT_TIMEOUT"
	WriteEncodingFlag  = "@@WRITE_ENCODING"
	FormatFlag         = "@@FORMAT"
	WriteDelimiterFlag = "@@WRITE_DELIMITER"
	WithoutHeaderFlag  = "@@WITHOUT_HEADER"
	PrettyPrintFlag    = "@@PRETTY_PRINT"
	ColorFlag          = "@@COLOR"
	QuietFlag          = "@@QUIET"
	CPUFlag            = "@@CPU"
	StatsFlag          = "@@STATS"
)

var flagList = []string{
	DelimiterFlag,
	EncodingFlag,
	LineBreakFlag,
	TimezoneFlag,
	RepositoryFlag,
	DatetimeFormatFlag,
	NoHeaderFlag,
	WithoutNullFlag,
	WaitTimeoutFlag,
	WriteEncodingFlag,
	FormatFlag,
	WriteDelimiterFlag,
	WithoutHeaderFlag,
	PrettyPrintFlag,
	ColorFlag,
	QuietFlag,
	CPUFlag,
	StatsFlag,
}

type ObjectStatus int

const (
	ObjectFixed ObjectStatus = iota
	ObjectCreated
	ObjectUpdated
)

func Print(expr parser.Print, filter *Filter) (string, error) {
	p, err := filter.Evaluate(expr.Value)
	if err != nil {
		return "", err
	}
	return p.String(), err
}

func Printf(expr parser.Printf, filter *Filter) (string, error) {
	var format string
	formatValue, err := filter.Evaluate(expr.Format)
	if err != nil {
		return "", err
	}
	formatString := value.ToString(formatValue)
	if !value.IsNull(formatString) {
		format = formatString.(value.String).Raw()
	}

	args := make([]value.Primary, len(expr.Values))
	for i, v := range expr.Values {
		p, err := filter.Evaluate(v)
		if err != nil {
			return "", err
		}
		args[i] = p
	}

	message, err := FormatString(format, args)
	if err != nil {
		return "", NewPrintfReplaceValueLengthError(expr, err.(AppError).ErrorMessage())
	}
	return message, nil
}

func Source(expr parser.Source, filter *Filter) ([]parser.Statement, error) {
	var fpath string

	if ident, ok := expr.FilePath.(parser.Identifier); ok {
		fpath = ident.Literal
	} else {
		p, err := filter.Evaluate(expr.FilePath)
		if err != nil {
			return nil, err
		}
		s := value.ToString(p)
		if value.IsNull(s) {
			return nil, NewSourceInvalidArgumentError(expr, expr.FilePath)
		}
		fpath = s.(value.String).Raw()
	}

	stat, err := os.Stat(fpath)
	if err != nil {
		return nil, NewSourceFileNotExistError(expr, fpath)
	}
	if stat.IsDir() {
		return nil, NewSourceFileUnableToReadError(expr, fpath)
	}

	fp, err := file.OpenToRead(fpath)
	if err != nil {
		return nil, NewReadFileError(expr, err.Error())
	}
	defer file.Close(fp)

	buf, err := ioutil.ReadAll(fp)
	if err != nil {
		return nil, NewReadFileError(expr, err.Error())
	}
	input := string(buf)

	statements, err := parser.Parse(input, fpath)
	if err != nil {
		syntaxErr := err.(*parser.SyntaxError)
		err = NewSyntaxError(syntaxErr.Message, syntaxErr.Line, syntaxErr.Char, syntaxErr.SourceFile)
	}
	return statements, err
}

func SetFlag(expr parser.SetFlag) error {
	var p value.Primary

	switch strings.ToUpper(expr.Name) {
	case DelimiterFlag, EncodingFlag, LineBreakFlag, TimezoneFlag, RepositoryFlag, DatetimeFormatFlag,
		WriteEncodingFlag, FormatFlag, WriteDelimiterFlag:
		p = value.ToString(expr.Value)
	case NoHeaderFlag, WithoutNullFlag, ColorFlag, StatsFlag, WithoutHeaderFlag, PrettyPrintFlag, QuietFlag:
		p = value.ToBoolean(expr.Value)
	case WaitTimeoutFlag:
		p = value.ToFloat(expr.Value)
	case CPUFlag:
		p = value.ToInteger(expr.Value)
	default:
		return NewInvalidFlagNameError(expr, expr.Name)
	}
	if value.IsNull(p) {
		return NewInvalidFlagValueError(expr)
	}

	var err error

	switch strings.ToUpper(expr.Name) {
	case DelimiterFlag:
		err = cmd.SetDelimiter(p.(value.String).Raw())
	case EncodingFlag:
		err = cmd.SetEncoding(p.(value.String).Raw())
	case LineBreakFlag:
		err = cmd.SetLineBreak(p.(value.String).Raw())
	case TimezoneFlag:
		err = cmd.SetLocation(p.(value.String).Raw())
	case RepositoryFlag:
		err = cmd.SetRepository(p.(value.String).Raw())
	case DatetimeFormatFlag:
		cmd.SetDatetimeFormat(p.(value.String).Raw())
	case NoHeaderFlag:
		cmd.SetNoHeader(p.(value.Boolean).Raw())
	case WithoutNullFlag:
		cmd.SetWithoutNull(p.(value.Boolean).Raw())
	case WaitTimeoutFlag:
		cmd.SetWaitTimeout(p.(value.Float).Raw())
	case WriteEncodingFlag:
		err = cmd.SetWriteEncoding(p.(value.String).Raw())
	case FormatFlag:
		err = cmd.SetFormat(p.(value.String).Raw())
	case WriteDelimiterFlag:
		err = cmd.SetWriteDelimiter(p.(value.String).Raw())
	case WithoutHeaderFlag:
		cmd.SetWithoutHeader(p.(value.Boolean).Raw())
	case PrettyPrintFlag:
		cmd.SetPrettyPrint(p.(value.Boolean).Raw())
	case ColorFlag:
		cmd.SetColor(p.(value.Boolean).Raw())
	case QuietFlag:
		cmd.SetQuiet(p.(value.Boolean).Raw())
	case CPUFlag:
		cmd.SetCPU(int(p.(value.Integer).Raw()))
	case StatsFlag:
		cmd.SetStats(p.(value.Boolean).Raw())
	}

	if err != nil {
		return NewInvalidFlagValueError(expr)
	}

	return nil
}

func ShowFlag(expr parser.ShowFlag) (string, error) {
	s, err := showFlag(expr.Name)
	if err != nil {
		return s, NewInvalidFlagNameError(expr, expr.Name)
	}

	palette := color.NewPalette()
	palette.Enable()

	flag := strings.ToUpper(expr.Name)
	style := color.StringStyle
	switch flag {
	case WaitTimeoutFlag, CPUFlag:
		style = color.NumberStyle
	case NoHeaderFlag, WithoutNullFlag, WithoutHeaderFlag, PrettyPrintFlag, ColorFlag, QuietFlag, StatsFlag:
		style = color.BooleanStyle
	default:
		if s == "(not set)" {
			style = color.NullStyle
		}
	}

	return " " + palette.Color(flag+":", color.FieldLableStyle) + " " + palette.Color(s, style), nil
}

func showFlag(flag string) (string, error) {
	var s string

	flags := cmd.GetFlags()

	switch strings.ToUpper(flag) {
	case DelimiterFlag:
		s = "'" + cmd.EscapeString(string(flags.Delimiter)) + "'"
	case EncodingFlag:
		s = flags.Encoding.String()
	case LineBreakFlag:
		s = flags.LineBreak.String()
	case TimezoneFlag:
		s = flags.Location
	case RepositoryFlag:
		s = flags.Repository
	case DatetimeFormatFlag:
		if len(flags.DatetimeFormat) < 1 {
			s = "(not set)"
		} else {
			s = flags.DatetimeFormat
		}
	case NoHeaderFlag:
		s = strconv.FormatBool(flags.NoHeader)
	case WithoutNullFlag:
		s = strconv.FormatBool(flags.WithoutNull)
	case WaitTimeoutFlag:
		s = value.Float64ToStr(flags.WaitTimeout)
	case WriteEncodingFlag:
		s = flags.WriteEncoding.String()
	case FormatFlag:
		s = flags.Format.String()
	case WriteDelimiterFlag:
		s = "'" + cmd.EscapeString(string(flags.WriteDelimiter)) + "'"
	case WithoutHeaderFlag:
		s = strconv.FormatBool(flags.WithoutHeader)
	case PrettyPrintFlag:
		s = strconv.FormatBool(flags.PrettyPrint)
	case ColorFlag:
		s = strconv.FormatBool(flags.Color)
	case QuietFlag:
		s = strconv.FormatBool(flags.Quiet)
	case CPUFlag:
		s = strconv.Itoa(flags.CPU)
	case StatsFlag:
		s = strconv.FormatBool(flags.Stats)
	default:
		return s, errors.New("invalid flag name")
	}

	return s, nil
}

func ShowObjects(expr parser.ShowObjects, filter *Filter) (string, error) {
	var s string

	switch strings.ToUpper(expr.Type.Literal) {
	case "TABLES":
		keys := ViewCache.SortedKeys()

		if len(keys) < 1 {
			s = color.Warn("No table is loaded")
		} else {
			w := text.NewObjectWriter()

			createdFiles, updatedFiles := UncommittedFiles()

			for _, key := range keys {
				fields := ViewCache[key].Header.TableColumnNames()
				info := ViewCache[key].FileInfo

				if _, ok := createdFiles[info.Path]; ok {
					w.WriteColor("*Created* ", color.EmphasisStyle)
				} else if _, ok := updatedFiles[info.Path]; ok {
					w.WriteColor("*Updated* ", color.EmphasisStyle)
				}
				w.WriteColorWithoutLineBreak(info.Path, color.ObjectStyle)
				writeFields(w, fields)

				w.NewLine()
				writeTableAttribute(w, info)
				w.ClearBlock()
				w.NewLine()
			}

			uncommitted := len(createdFiles) + len(updatedFiles)

			w.Title1 = "Loaded Tables"
			if 0 < uncommitted {
				w.Title2 = fmt.Sprintf("(Uncommitted: %s)", FormatCount(uncommitted, "Table"))
				w.Title2Style = color.EmphasisStyle
			}
			s = "\n" + w.String()
		}
	case "VIEWS":
		views := filter.TempViews.All()

		if len(views) < 1 {
			s = color.Warn("No view is declared")
		} else {
			keys := views.SortedKeys()

			w := text.NewObjectWriter()

			updatedViews := UncommittedTempViews()

			for _, key := range keys {
				fields := views[key].Header.TableColumnNames()
				info := views[key].FileInfo

				if _, ok := updatedViews[info.Path]; ok {
					w.WriteColor("*Updated* ", color.EmphasisStyle)
				}
				w.WriteColorWithoutLineBreak(info.Path, color.ObjectStyle)
				writeFields(w, fields)
				w.ClearBlock()
				w.NewLine()
			}

			uncommitted := len(updatedViews)

			w.Title1 = "Views"
			if 0 < uncommitted {
				w.Title2 = fmt.Sprintf("(Uncommitted: %s)", FormatCount(uncommitted, "View"))
				w.Title2Style = color.EmphasisStyle
			}
			s = "\n" + w.String()
		}
	case "CURSORS":
		cursors := filter.Cursors.All()
		if len(cursors) < 1 {
			s = color.Warn("No cursor is declared")
		} else {
			keys := cursors.SortedKeys()

			w := text.NewObjectWriter()

			for _, key := range keys {
				cur := cursors[key]
				isOpen := cur.IsOpen()

				w.WriteColor(cur.name, color.ObjectStyle)
				w.BeginBlock()

				w.NewLine()
				w.WriteColorWithoutLineBreak("Status: ", color.FieldLableStyle)
				if isOpen == ternary.TRUE {
					nor, _ := cur.Count()
					inRange, _ := cur.IsInRange()
					position, _ := cur.Pointer()

					norStr := cmd.HumarizeNumber(strconv.Itoa(nor))

					w.WriteColorWithoutLineBreak("Open", color.TernaryStyle)
					w.WriteColorWithoutLineBreak("    Number of Rows: ", color.FieldLableStyle)
					w.WriteColorWithoutLineBreak(norStr, color.NumberStyle)
					w.WriteSpaces(10 - len(norStr))
					w.WriteColorWithoutLineBreak("Pointer: ", color.FieldLableStyle)
					switch inRange {
					case ternary.FALSE:
						w.WriteColorWithoutLineBreak("Out of Range", color.TernaryStyle)
					case ternary.UNKNOWN:
						w.WriteColorWithoutLineBreak(inRange.String(), color.TernaryStyle)
					default:
						w.WriteColorWithoutLineBreak(cmd.HumarizeNumber(strconv.Itoa(position)), color.NumberStyle)
					}
				} else {
					w.WriteColorWithoutLineBreak("Closed", color.TernaryStyle)
				}

				w.NewLine()
				w.WriteColorWithoutLineBreak("Query: ", color.FieldLableStyle)
				w.WriteColorWithoutLineBreak(cur.query.String(), color.IdentifierStyle)

				w.ClearBlock()
				w.NewLine()
			}
			w.Title1 = "Cursors"
			s = "\n" + w.String()
		}
	case "FUNCTIONS":
		scalas, aggs := filter.Functions.All()
		if len(scalas) < 1 && len(aggs) < 1 {
			s = color.Warn("No function is declared")
		} else {
			if 0 < len(scalas) {
				w := text.NewObjectWriter()
				writeFunctions(w, scalas)
				w.Title1 = "Scala Functions"
				s += "\n" + w.String()
			}
			if 0 < len(aggs) {
				w := text.NewObjectWriter()
				writeFunctions(w, aggs)
				w.Title1 = "Aggregate Functions"
				s += "\n" + w.String()
			}
		}
	case "FLAGS":
		w := text.NewObjectWriter()
		for _, flag := range flagList {
			s, _ := showFlag(flag)
			style := color.StringStyle
			switch flag {
			case WaitTimeoutFlag, CPUFlag:
				style = color.NumberStyle
			case NoHeaderFlag, WithoutNullFlag, WithoutHeaderFlag, PrettyPrintFlag, ColorFlag, QuietFlag, StatsFlag:
				style = color.BooleanStyle
			default:
				if s == "(not set)" {
					style = color.NullStyle
				}
			}
			w.WriteSpaces(18 - len(flag))
			w.WriteColorWithoutLineBreak(flag, color.FieldLableStyle)
			w.WriteColorWithoutLineBreak(":", color.FieldLableStyle)
			w.WriteSpaces(1)
			w.WriteColorWithoutLineBreak(s, style)
			w.NewLine()
		}
		w.Title1 = "Flags"
		s = "\n" + w.String()
	default:
		return "", NewShowInvalidObjectTypeError(expr, expr.Type.String())
	}

	return s, nil
}

func writeTableAttribute(w *text.ObjectWriter, info *FileInfo) {
	w.WriteColor("Format: ", color.FieldLableStyle)
	w.WriteWithoutLineBreak(info.Format.String())

	w.WriteSpaces(8 - text.StringWidth(info.Format.String()))
	switch info.Format {
	case cmd.CSV:
		w.WriteColorWithoutLineBreak("Delimiter: ", color.FieldLableStyle)
		w.WriteWithoutLineBreak("'" + cmd.EscapeString(string(info.Delimiter)) + "'")
	case cmd.FIXED:
		w.WriteColorWithoutLineBreak("Delimiter Positions: ", color.FieldLableStyle)
		w.WriteWithoutLineBreak(info.DelimiterPositions.String())
	case cmd.JSON:
		w.WriteColorWithoutLineBreak("Query: ", color.FieldLableStyle)
		if len(info.JsonQuery) < 1 {
			w.WriteColorWithoutLineBreak("(empty)", color.NullStyle)
		} else {
			w.WriteWithoutLineBreak(info.JsonQuery)
		}
	}

	w.NewLine()

	w.WriteColor("Encoding: ", color.FieldLableStyle)
	w.WriteWithoutLineBreak(info.Encoding.String())

	w.WriteSpaces(6 - (text.StringWidth(info.Encoding.String())))
	w.WriteColorWithoutLineBreak("LineBreak: ", color.FieldLableStyle)
	w.WriteWithoutLineBreak(info.LineBreak.String())

	if info.Format == cmd.JSON {
		w.WriteSpaces(6 - (text.StringWidth(info.LineBreak.String())))
		w.WriteColorWithoutLineBreak("Pretty Print: ", color.FieldLableStyle)
		w.WriteWithoutLineBreak(strconv.FormatBool(info.PrettyPrint))
	} else {
		w.WriteSpaces(6 - (text.StringWidth(info.LineBreak.String())))
		w.WriteColorWithoutLineBreak("Header: ", color.FieldLableStyle)
		w.WriteWithoutLineBreak(strconv.FormatBool(!info.NoHeader))
	}
}

func writeFields(w *text.ObjectWriter, fields []string) {
	w.BeginBlock()
	w.NewLine()
	w.WriteColor("Fields: ", color.FieldLableStyle)
	w.BeginSubBlock()
	lastIdx := len(fields) - 1
	for i, f := range fields {
		escaped := cmd.EscapeString(f)
		if i < lastIdx && !w.FitInLine(escaped+", ") {
			w.NewLine()
		}
		w.WriteColor(escaped, color.AttributeStyle)
		if i < lastIdx {
			w.WriteWithoutLineBreak(", ")
		}
	}
	w.EndSubBlock()
}

func writeFunctions(w *text.ObjectWriter, funcs UserDefinedFunctionMap) {
	keys := funcs.SortedKeys()

	for _, key := range keys {
		fn := funcs[key]

		w.WriteColor(fn.Name.String(), color.ObjectStyle)
		w.WriteWithoutLineBreak(" (")

		if fn.IsAggregate {
			w.WriteColorWithoutLineBreak(fn.Cursor.String(), color.IdentifierStyle)
			if 0 < len(fn.Parameters) {
				w.WriteWithoutLineBreak(", ")
			}
		}

		for i, p := range fn.Parameters {
			if 0 < i {
				w.WriteWithoutLineBreak(", ")
			}
			if def, ok := fn.Defaults[p.String()]; ok {
				w.WriteColorWithoutLineBreak(p.String(), color.AttributeStyle)
				w.WriteWithoutLineBreak(" = ")
				w.WriteColorWithoutLineBreak(def.String(), color.ValueStyle)
			} else {
				w.WriteColorWithoutLineBreak(p.String(), color.AttributeStyle)
			}
		}

		w.WriteWithoutLineBreak(")")
		w.ClearBlock()
		w.NewLine()
	}
}

func ShowFields(expr parser.ShowFields, filter *Filter) (string, error) {
	if !strings.EqualFold(expr.Type.Literal, "FIELDS") {
		return "", NewShowInvalidObjectTypeError(expr, expr.Type.Literal)
	}

	var status = ObjectFixed

	view := NewView()
	err := view.LoadFromTableIdentifier(expr.Table, filter.CreateNode())
	if err != nil {
		return "", err
	}

	if view.FileInfo.IsTemporary {
		updatedViews := UncommittedTempViews()

		if _, ok := updatedViews[view.FileInfo.Path]; ok {
			status = ObjectUpdated
		}
	} else {
		createdViews, updatedView := UncommittedFiles()
		if _, ok := createdViews[view.FileInfo.Path]; ok {
			status = ObjectCreated
		} else if _, ok := updatedView[view.FileInfo.Path]; ok {
			status = ObjectUpdated
		}
	}

	w := text.NewObjectWriter()
	w.WriteColorWithoutLineBreak("Type: ", color.FieldLableStyle)
	if view.FileInfo.IsTemporary {
		w.WriteWithoutLineBreak("View")
	} else {
		w.WriteWithoutLineBreak("Table")
		w.NewLine()
		w.WriteColorWithoutLineBreak("Path: ", color.FieldLableStyle)
		w.WriteColorWithoutLineBreak(view.FileInfo.Path, color.ObjectStyle)
		w.NewLine()
		writeTableAttribute(w, view.FileInfo)
	}

	w.NewLine()
	w.WriteColorWithoutLineBreak("Status: ", color.FieldLableStyle)
	switch status {
	case ObjectCreated:
		w.WriteColorWithoutLineBreak("Created", color.EmphasisStyle)
	case ObjectUpdated:
		w.WriteColorWithoutLineBreak("Updated", color.EmphasisStyle)
	default:
		w.WriteWithoutLineBreak("Fixed")
	}

	w.NewLine()
	writeFieldList(w, view.Header.TableColumnNames())

	w.Title1 = "Fields in"
	w.Title2 = expr.Table.Literal
	w.Title2Style = color.IdentifierStyle
	return "\n" + w.String(), nil
}

func writeFieldList(w *text.ObjectWriter, fields []string) {
	l := len(fields)
	digits := len(strconv.Itoa(l))
	fieldNumbers := make([]string, 0, l)
	for i := 0; i < l; i++ {
		idxstr := strconv.Itoa(i + 1)
		fieldNumbers = append(fieldNumbers, strings.Repeat(" ", digits-len(idxstr))+idxstr)
	}

	w.WriteColorWithoutLineBreak("Fields:", color.FieldLableStyle)
	w.NewLine()
	w.WriteSpaces(2)
	w.BeginSubBlock()
	for i := 0; i < l; i++ {
		w.WriteColor(fieldNumbers[i], color.NumberStyle)
		w.Write(".")
		w.WriteSpaces(1)
		w.WriteColorWithoutLineBreak(fields[i], color.AttributeStyle)
		w.NewLine()
	}
}
