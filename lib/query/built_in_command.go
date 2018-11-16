package query

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/file"
	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/value"

	"github.com/mithrandie/go-text"
	"github.com/mithrandie/go-text/color"
	"github.com/mithrandie/go-text/fixedlen"
	"github.com/mithrandie/ternary"
)

type ObjectStatus int

const (
	ObjectFixed ObjectStatus = iota
	ObjectCreated
	ObjectUpdated
)

const IgnoredFlagPrefix = "(ignored) "

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

	f := NewStringFormatter()
	message, err := f.Format(format, args)
	if err != nil {
		return "", NewReplaceValueLengthError(expr, err.(AppError).ErrorMessage())
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
			return nil, NewSourceInvalidFilePathError(expr, expr.FilePath)
		}
		fpath = s.(value.String).Raw()
	}
	if len(fpath) < 1 {
		return nil, NewSourceInvalidFilePathError(expr, expr.FilePath)
	}
	if abs, err := filepath.Abs(fpath); err == nil {
		fpath = abs
	}

	stat, err := os.Stat(fpath)
	if err != nil {
		return nil, NewSourceFileNotExistError(expr, fpath)
	}
	if stat.IsDir() {
		return nil, NewSourceFileUnableToReadError(expr, fpath)
	}

	h, err := file.NewHandlerForRead(fpath)
	if err != nil {
		return nil, NewReadFileError(expr, err.Error())
	}
	defer h.Close()

	buf, err := ioutil.ReadAll(h.FileForRead())
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

func ParseExecuteStatements(expr parser.Execute, filter *Filter) ([]parser.Statement, error) {
	var input string
	stmt, err := filter.Evaluate(expr.Statements)
	if err != nil {
		return nil, err
	}
	stmt = value.ToString(stmt)
	if !value.IsNull(stmt) {
		input = stmt.(value.String).Raw()
	}

	args := make([]value.Primary, len(expr.Values))
	for i, v := range expr.Values {
		p, err := filter.Evaluate(v)
		if err != nil {
			return nil, err
		}
		args[i] = p
	}

	f := NewStringFormatter()
	input, err = f.Format(input, args)
	if err != nil {
		return nil, NewReplaceValueLengthError(expr, err.(AppError).ErrorMessage())
	}
	statements, err := parser.Parse(input, fmt.Sprintf("(L:%d C:%d) EXECUTE", expr.Line(), expr.Char()))
	if err != nil {
		syntaxErr := err.(*parser.SyntaxError)
		err = NewSyntaxError(syntaxErr.Message, syntaxErr.Line, syntaxErr.Char, syntaxErr.SourceFile)
	}
	return statements, err
}

func SetFlag(expr parser.SetFlag, filter *Filter) (string, error) {
	var p value.Primary
	var err error

	if ident, ok := expr.Value.(parser.Identifier); ok {
		p = value.NewString(ident.Literal)
	} else {
		p, err = filter.Evaluate(expr.Value)
		if err != nil {
			return "", err
		}
	}

	switch strings.ToUpper(expr.Name) {
	case cmd.RepositoryFlag, cmd.TimezoneFlag, cmd.DatetimeFormatFlag, cmd.DelimiterFlag, cmd.JsonQuery, cmd.EncodingFlag,
		cmd.WriteEncodingFlag, cmd.FormatFlag, cmd.WriteDelimiterFlag, cmd.LineBreakFlag:
		p = value.ToString(p)
	case cmd.NoHeaderFlag, cmd.WithoutNullFlag, cmd.WithoutHeaderFlag, cmd.EncloseAll, cmd.PrettyPrintFlag,
		cmd.EastAsianEncoding, cmd.CountDiacriticalSign, cmd.CountFormatCode, cmd.ColorFlag, cmd.QuietFlag, cmd.StatsFlag:
		p = value.ToBoolean(p)
	case cmd.WaitTimeoutFlag:
		p = value.ToFloat(p)
	case cmd.CPUFlag:
		p = value.ToInteger(p)
	default:
		return "", NewInvalidFlagNameError(expr, expr.Name)
	}
	if value.IsNull(p) {
		return "", NewFlagValueNotAllowedFormatError(expr)
	}

	switch strings.ToUpper(expr.Name) {
	case cmd.RepositoryFlag:
		err = cmd.SetRepository(p.(value.String).Raw())
	case cmd.TimezoneFlag:
		err = cmd.SetLocation(p.(value.String).Raw())
	case cmd.DatetimeFormatFlag:
		cmd.SetDatetimeFormat(p.(value.String).Raw())
	case cmd.WaitTimeoutFlag:
		cmd.SetWaitTimeout(p.(value.Float).Raw())
	case cmd.DelimiterFlag:
		err = cmd.SetDelimiter(p.(value.String).Raw())
	case cmd.JsonQuery:
		cmd.SetJsonQuery(p.(value.String).Raw())
	case cmd.EncodingFlag:
		err = cmd.SetEncoding(p.(value.String).Raw())
	case cmd.NoHeaderFlag:
		cmd.SetNoHeader(p.(value.Boolean).Raw())
	case cmd.WithoutNullFlag:
		cmd.SetWithoutNull(p.(value.Boolean).Raw())
	case cmd.FormatFlag:
		err = cmd.SetFormat(p.(value.String).Raw())
	case cmd.WriteEncodingFlag:
		err = cmd.SetWriteEncoding(p.(value.String).Raw())
	case cmd.WriteDelimiterFlag:
		err = cmd.SetWriteDelimiter(p.(value.String).Raw())
	case cmd.WithoutHeaderFlag:
		cmd.SetWithoutHeader(p.(value.Boolean).Raw())
	case cmd.LineBreakFlag:
		err = cmd.SetLineBreak(p.(value.String).Raw())
	case cmd.EncloseAll:
		cmd.SetEncloseAll(p.(value.Boolean).Raw())
	case cmd.PrettyPrintFlag:
		cmd.SetPrettyPrint(p.(value.Boolean).Raw())
	case cmd.EastAsianEncoding:
		cmd.SetEastAsianEncoding(p.(value.Boolean).Raw())
	case cmd.CountDiacriticalSign:
		cmd.SetCountDiacriticalSign(p.(value.Boolean).Raw())
	case cmd.CountFormatCode:
		cmd.SetCountFormatCode(p.(value.Boolean).Raw())
	case cmd.ColorFlag:
		cmd.SetColor(p.(value.Boolean).Raw())
	case cmd.QuietFlag:
		cmd.SetQuiet(p.(value.Boolean).Raw())
	case cmd.CPUFlag:
		cmd.SetCPU(int(p.(value.Integer).Raw()))
	case cmd.StatsFlag:
		cmd.SetStats(p.(value.Boolean).Raw())
	}

	if err != nil {
		return "", NewInvalidFlagValueError(expr, err.Error())
	}

	s, _ := showFlag(expr.Name)
	return " " + cmd.GetPalette().Render(cmd.LableEffect, cmd.FlagSymbol(strings.ToUpper(expr.Name)+":")) + " " + s, nil
}

func ShowFlag(expr parser.ShowFlag) (string, error) {
	s, err := showFlag(expr.Name)
	if err != nil {
		return s, NewInvalidFlagNameError(expr, expr.Name)
	}

	return " " + cmd.GetPalette().Render(cmd.LableEffect, cmd.FlagSymbol(strings.ToUpper(expr.Name)+":")) + " " + s, nil
}

func showFlag(flag string) (string, error) {
	var s string

	flags := cmd.GetFlags()

	switch strings.ToUpper(flag) {
	case cmd.RepositoryFlag:
		s = cmd.GetPalette().Render(cmd.StringEffect, flags.Repository)
	case cmd.TimezoneFlag:
		s = cmd.GetPalette().Render(cmd.StringEffect, flags.Location)
	case cmd.DatetimeFormatFlag:
		if len(flags.DatetimeFormat) < 1 {
			s = cmd.GetPalette().Render(cmd.NullEffect, "(not set)")
		} else {
			s = cmd.GetPalette().Render(cmd.StringEffect, flags.DatetimeFormat)
		}
	case cmd.WaitTimeoutFlag:
		s = cmd.GetPalette().Render(cmd.NumberEffect, value.Float64ToStr(flags.WaitTimeout))
	case cmd.DelimiterFlag:
		d := "'" + cmd.EscapeString(string(flags.Delimiter)) + "'"
		p := fixedlen.DelimiterPositions(flags.DelimiterPositions).String()

		switch flags.ImportFormat() {
		case cmd.CSV, cmd.TSV:
			s = cmd.GetPalette().Render(cmd.StringEffect, d) + cmd.GetPalette().Render(cmd.LableEffect, " | ") + cmd.GetPalette().Render(cmd.NullEffect, p)
		case cmd.FIXED:
			s = cmd.GetPalette().Render(cmd.NullEffect, d) + cmd.GetPalette().Render(cmd.LableEffect, " | ") + cmd.GetPalette().Render(cmd.StringEffect, p)
		default:
			s = cmd.GetPalette().Render(cmd.NullEffect, IgnoredFlagPrefix+d+" | "+p)
		}
	case cmd.JsonQuery:
		q := flags.JsonQuery
		if len(q) < 1 {
			q = "(empty)"
		}

		switch flags.ImportFormat() {
		case cmd.JSON:
			s = cmd.GetPalette().Render(cmd.StringEffect, q)
		default:
			s = cmd.GetPalette().Render(cmd.NullEffect, IgnoredFlagPrefix+q)
		}
	case cmd.EncodingFlag:
		s = cmd.GetPalette().Render(cmd.StringEffect, flags.Encoding.String())
	case cmd.NoHeaderFlag:
		s = cmd.GetPalette().Render(cmd.BooleanEffect, strconv.FormatBool(flags.NoHeader))
	case cmd.WithoutNullFlag:
		s = cmd.GetPalette().Render(cmd.BooleanEffect, strconv.FormatBool(flags.WithoutNull))
	case cmd.FormatFlag:
		s = cmd.GetPalette().Render(cmd.StringEffect, flags.Format.String())
	case cmd.WriteEncodingFlag:
		switch flags.Format {
		case cmd.JSON, cmd.JSONH, cmd.JSONA:
			s = cmd.GetPalette().Render(cmd.NullEffect, IgnoredFlagPrefix+flags.WriteEncoding.String())
		default:
			s = cmd.GetPalette().Render(cmd.StringEffect, flags.WriteEncoding.String())
		}
	case cmd.WriteDelimiterFlag:
		d := "'" + cmd.EscapeString(string(flags.WriteDelimiter)) + "'"
		p := fixedlen.DelimiterPositions(flags.WriteDelimiterPositions).String()
		switch flags.Format {
		case cmd.CSV:
			s = cmd.GetPalette().Render(cmd.StringEffect, d) + cmd.GetPalette().Render(cmd.LableEffect, " | ") + cmd.GetPalette().Render(cmd.NullEffect, p)
		case cmd.FIXED:
			s = cmd.GetPalette().Render(cmd.NullEffect, d) + cmd.GetPalette().Render(cmd.LableEffect, " | ") + cmd.GetPalette().Render(cmd.StringEffect, p)
		default:
			s = cmd.GetPalette().Render(cmd.NullEffect, IgnoredFlagPrefix+d+" | "+p)
		}
	case cmd.WithoutHeaderFlag:
		s = strconv.FormatBool(flags.WithoutHeader)
		switch flags.Format {
		case cmd.CSV, cmd.TSV, cmd.FIXED, cmd.GFM, cmd.ORG:
			s = cmd.GetPalette().Render(cmd.BooleanEffect, s)
		default:
			s = cmd.GetPalette().Render(cmd.NullEffect, IgnoredFlagPrefix+s)
		}
	case cmd.LineBreakFlag:
		s = cmd.GetPalette().Render(cmd.StringEffect, flags.LineBreak.String())
	case cmd.EncloseAll:
		s = strconv.FormatBool(flags.EncloseAll)
		switch flags.Format {
		case cmd.CSV, cmd.TSV:
			s = cmd.GetPalette().Render(cmd.BooleanEffect, s)
		default:
			s = cmd.GetPalette().Render(cmd.NullEffect, IgnoredFlagPrefix+s)
		}
	case cmd.PrettyPrintFlag:
		s = strconv.FormatBool(flags.PrettyPrint)
		switch flags.Format {
		case cmd.JSON, cmd.JSONH, cmd.JSONA:
			s = cmd.GetPalette().Render(cmd.BooleanEffect, s)
		default:
			s = cmd.GetPalette().Render(cmd.NullEffect, IgnoredFlagPrefix+s)
		}
	case cmd.EastAsianEncoding:
		s = strconv.FormatBool(flags.EastAsianEncoding)
		switch flags.Format {
		case cmd.GFM, cmd.ORG, cmd.TEXT:
			s = cmd.GetPalette().Render(cmd.BooleanEffect, s)
		default:
			s = cmd.GetPalette().Render(cmd.NullEffect, IgnoredFlagPrefix+s)
		}
	case cmd.CountDiacriticalSign:
		s = strconv.FormatBool(flags.CountDiacriticalSign)
		switch flags.Format {
		case cmd.GFM, cmd.ORG, cmd.TEXT:
			s = cmd.GetPalette().Render(cmd.BooleanEffect, s)
		default:
			s = cmd.GetPalette().Render(cmd.NullEffect, IgnoredFlagPrefix+s)
		}
	case cmd.CountFormatCode:
		s = strconv.FormatBool(flags.CountFormatCode)
		switch flags.Format {
		case cmd.GFM, cmd.ORG, cmd.TEXT:
			s = cmd.GetPalette().Render(cmd.BooleanEffect, s)
		default:
			s = cmd.GetPalette().Render(cmd.NullEffect, IgnoredFlagPrefix+s)
		}
	case cmd.ColorFlag:
		s = cmd.GetPalette().Render(cmd.BooleanEffect, strconv.FormatBool(flags.Color))
	case cmd.QuietFlag:
		s = cmd.GetPalette().Render(cmd.BooleanEffect, strconv.FormatBool(flags.Quiet))
	case cmd.CPUFlag:
		s = cmd.GetPalette().Render(cmd.NumberEffect, strconv.Itoa(flags.CPU))
	case cmd.StatsFlag:
		s = cmd.GetPalette().Render(cmd.BooleanEffect, strconv.FormatBool(flags.Stats))
	default:
		return s, errors.New("invalid flag name")
	}

	return s, nil
}

func ShowObjects(expr parser.ShowObjects, filter *Filter) (string, error) {
	var s string

	w := cmd.NewObjectWriter()

	switch strings.ToUpper(expr.Type.Literal) {
	case "TABLES":
		keys := ViewCache.SortedKeys()

		if len(keys) < 1 {
			s = color.Warn("No table is loaded")
		} else {
			createdFiles, updatedFiles := UncommittedFiles()

			for _, key := range keys {
				fields := ViewCache[key].Header.TableColumnNames()
				info := ViewCache[key].FileInfo

				if _, ok := createdFiles[info.Path]; ok {
					w.WriteColor("*Created* ", cmd.EmphasisEffect)
				} else if _, ok := updatedFiles[info.Path]; ok {
					w.WriteColor("*Updated* ", cmd.EmphasisEffect)
				}
				w.WriteColorWithoutLineBreak(info.Path, cmd.ObjectEffect)
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
				w.Title2Effect = cmd.EmphasisEffect
			}
			s = "\n" + w.String()
		}
	case "VIEWS":
		views := filter.TempViews.All()

		if len(views) < 1 {
			s = color.Warn("No view is declared")
		} else {
			keys := views.SortedKeys()

			updatedViews := UncommittedTempViews()

			for _, key := range keys {
				fields := views[key].Header.TableColumnNames()
				info := views[key].FileInfo

				if _, ok := updatedViews[info.Path]; ok {
					w.WriteColor("*Updated* ", cmd.EmphasisEffect)
				}
				w.WriteColorWithoutLineBreak(info.Path, cmd.ObjectEffect)
				writeFields(w, fields)
				w.ClearBlock()
				w.NewLine()
			}

			uncommitted := len(updatedViews)

			w.Title1 = "Views"
			if 0 < uncommitted {
				w.Title2 = fmt.Sprintf("(Uncommitted: %s)", FormatCount(uncommitted, "View"))
				w.Title2Effect = cmd.EmphasisEffect
			}
			s = "\n" + w.String()
		}
	case "CURSORS":
		cursors := filter.Cursors.All()
		if len(cursors) < 1 {
			s = color.Warn("No cursor is declared")
		} else {
			keys := cursors.SortedKeys()

			for _, key := range keys {
				cur := cursors[key]
				isOpen := cur.IsOpen()

				w.WriteColor(cur.name, cmd.ObjectEffect)
				w.BeginBlock()

				w.NewLine()
				w.WriteColorWithoutLineBreak("Status: ", cmd.LableEffect)
				if isOpen == ternary.TRUE {
					nor, _ := cur.Count()
					inRange, _ := cur.IsInRange()
					position, _ := cur.Pointer()

					norStr := cmd.FormatInt(nor, ",")

					w.WriteColorWithoutLineBreak("Open", cmd.TernaryEffect)
					w.WriteColorWithoutLineBreak("    Number of Rows: ", cmd.LableEffect)
					w.WriteColorWithoutLineBreak(norStr, cmd.NumberEffect)
					w.WriteSpaces(10 - len(norStr))
					w.WriteColorWithoutLineBreak("Pointer: ", cmd.LableEffect)
					switch inRange {
					case ternary.FALSE:
						w.WriteColorWithoutLineBreak("Out of Range", cmd.TernaryEffect)
					case ternary.UNKNOWN:
						w.WriteColorWithoutLineBreak(inRange.String(), cmd.TernaryEffect)
					default:
						w.WriteColorWithoutLineBreak(cmd.FormatInt(position, ","), cmd.NumberEffect)
					}
				} else {
					w.WriteColorWithoutLineBreak("Closed", cmd.TernaryEffect)
				}

				w.NewLine()
				w.WriteColorWithoutLineBreak("Query: ", cmd.LableEffect)
				w.WriteColorWithoutLineBreak(cur.query.String(), cmd.IdentifierEffect)

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
				w.Clear()
				writeFunctions(w, scalas)
				w.Title1 = "Scala Functions"
				s += "\n" + w.String()
			}
			if 0 < len(aggs) {
				w.Clear()
				writeFunctions(w, aggs)
				w.Title1 = "Aggregate Functions"
				s += "\n" + w.String()
			}
		}
	case "FLAGS":
		for _, flag := range cmd.FlagList {
			symbol := cmd.FlagSymbol(flag)
			s, _ := showFlag(flag)
			w.WriteSpaces(24 - len(symbol))
			w.WriteColorWithoutLineBreak(symbol, cmd.LableEffect)
			w.WriteColorWithoutLineBreak(":", cmd.LableEffect)
			w.WriteSpaces(1)
			w.WriteWithoutLineBreak(s)
			w.NewLine()
		}
		w.Title1 = "Flags"
		s = "\n" + w.String()
	default:
		return "", NewShowInvalidObjectTypeError(expr, expr.Type.String())
	}

	return s, nil
}

func writeTableAttribute(w *cmd.ObjectWriter, info *FileInfo) {
	w.WriteColor("Format: ", cmd.LableEffect)
	w.WriteWithoutLineBreak(info.Format.String())

	w.WriteSpaces(8 - cmd.TextWidth(info.Format.String()))
	switch info.Format {
	case cmd.CSV:
		w.WriteColorWithoutLineBreak("Delimiter: ", cmd.LableEffect)
		w.WriteWithoutLineBreak("'" + cmd.EscapeString(string(info.Delimiter)) + "'")
	case cmd.TSV:
		w.WriteColorWithoutLineBreak("Delimiter: ", cmd.LableEffect)
		w.WriteColorWithoutLineBreak("'\\t'", cmd.NullEffect)
	case cmd.FIXED:
		w.WriteColorWithoutLineBreak("Delimiter Positions: ", cmd.LableEffect)
		w.WriteWithoutLineBreak(info.DelimiterPositions.String())
	case cmd.JSON, cmd.JSONH, cmd.JSONA:
		w.WriteColorWithoutLineBreak("Query: ", cmd.LableEffect)
		if len(info.JsonQuery) < 1 {
			w.WriteColorWithoutLineBreak("(empty)", cmd.NullEffect)
		} else {
			w.WriteColorWithoutLineBreak(info.JsonQuery, cmd.NullEffect)
		}
	}

	switch info.Format {
	case cmd.CSV, cmd.TSV:
		w.WriteSpaces(4 - (cmd.TextWidth(cmd.EscapeString(string(info.Delimiter)))))
		w.WriteColorWithoutLineBreak("Enclose All: ", cmd.LableEffect)
		w.WriteWithoutLineBreak(strconv.FormatBool(info.EncloseAll))
	}

	w.NewLine()

	w.WriteColor("Encoding: ", cmd.LableEffect)
	switch info.Format {
	case cmd.JSON, cmd.JSONH, cmd.JSONA:
		w.WriteColorWithoutLineBreak(text.UTF8.String(), cmd.NullEffect)
	default:
		w.WriteWithoutLineBreak(info.Encoding.String())
	}

	w.WriteSpaces(6 - (cmd.TextWidth(info.Encoding.String())))
	w.WriteColorWithoutLineBreak("LineBreak: ", cmd.LableEffect)
	w.WriteWithoutLineBreak(info.LineBreak.String())

	switch info.Format {
	case cmd.JSON, cmd.JSONH, cmd.JSONA:
		w.WriteSpaces(6 - (cmd.TextWidth(info.LineBreak.String())))
		w.WriteColorWithoutLineBreak("Pretty Print: ", cmd.LableEffect)
		w.WriteWithoutLineBreak(strconv.FormatBool(info.PrettyPrint))
	case cmd.CSV, cmd.TSV, cmd.FIXED, cmd.GFM, cmd.ORG:
		w.WriteSpaces(6 - (cmd.TextWidth(info.LineBreak.String())))
		w.WriteColorWithoutLineBreak("Header: ", cmd.LableEffect)
		w.WriteWithoutLineBreak(strconv.FormatBool(!info.NoHeader))
	}
}

func writeFields(w *cmd.ObjectWriter, fields []string) {
	w.BeginBlock()
	w.NewLine()
	w.WriteColor("Fields: ", cmd.LableEffect)
	w.BeginSubBlock()
	lastIdx := len(fields) - 1
	for i, f := range fields {
		escaped := cmd.EscapeString(f)
		if i < lastIdx && !w.FitInLine(escaped+", ") {
			w.NewLine()
		}
		w.WriteColor(escaped, cmd.AttributeEffect)
		if i < lastIdx {
			w.WriteWithoutLineBreak(", ")
		}
	}
	w.EndSubBlock()
}

func writeFunctions(w *cmd.ObjectWriter, funcs UserDefinedFunctionMap) {
	keys := funcs.SortedKeys()

	for _, key := range keys {
		fn := funcs[key]

		w.WriteColor(fn.Name.String(), cmd.ObjectEffect)
		w.WriteWithoutLineBreak(" (")

		if fn.IsAggregate {
			w.WriteColorWithoutLineBreak(fn.Cursor.String(), cmd.IdentifierEffect)
			if 0 < len(fn.Parameters) {
				w.WriteWithoutLineBreak(", ")
			}
		}

		for i, p := range fn.Parameters {
			if 0 < i {
				w.WriteWithoutLineBreak(", ")
			}
			if def, ok := fn.Defaults[p.Name]; ok {
				w.WriteColorWithoutLineBreak(p.String(), cmd.AttributeEffect)
				w.WriteWithoutLineBreak(" = ")
				w.WriteColorWithoutLineBreak(def.String(), cmd.ValueEffect)
			} else {
				w.WriteColorWithoutLineBreak(p.String(), cmd.AttributeEffect)
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

	w := cmd.NewObjectWriter()
	w.WriteColorWithoutLineBreak("Type: ", cmd.LableEffect)
	if view.FileInfo.IsTemporary {
		w.WriteWithoutLineBreak("View")
	} else {
		w.WriteWithoutLineBreak("Table")
		w.NewLine()
		w.WriteColorWithoutLineBreak("Path: ", cmd.LableEffect)
		w.WriteColorWithoutLineBreak(view.FileInfo.Path, cmd.ObjectEffect)
		w.NewLine()
		writeTableAttribute(w, view.FileInfo)
	}

	w.NewLine()
	w.WriteColorWithoutLineBreak("Status: ", cmd.LableEffect)
	switch status {
	case ObjectCreated:
		w.WriteColorWithoutLineBreak("Created", cmd.EmphasisEffect)
	case ObjectUpdated:
		w.WriteColorWithoutLineBreak("Updated", cmd.EmphasisEffect)
	default:
		w.WriteWithoutLineBreak("Fixed")
	}

	w.NewLine()
	writeFieldList(w, view.Header.TableColumnNames())

	w.Title1 = "Fields in"
	w.Title2 = expr.Table.Literal
	w.Title2Effect = cmd.IdentifierEffect
	return "\n" + w.String(), nil
}

func writeFieldList(w *cmd.ObjectWriter, fields []string) {
	l := len(fields)
	digits := len(strconv.Itoa(l))
	fieldNumbers := make([]string, 0, l)
	for i := 0; i < l; i++ {
		idxstr := strconv.Itoa(i + 1)
		fieldNumbers = append(fieldNumbers, strings.Repeat(" ", digits-len(idxstr))+idxstr)
	}

	w.WriteColorWithoutLineBreak("Fields:", cmd.LableEffect)
	w.NewLine()
	w.WriteSpaces(2)
	w.BeginSubBlock()
	for i := 0; i < l; i++ {
		w.WriteColor(fieldNumbers[i], cmd.NumberEffect)
		w.Write(".")
		w.WriteSpaces(1)
		w.WriteColorWithoutLineBreak(fields[i], cmd.AttributeEffect)
		w.NewLine()
	}
}
