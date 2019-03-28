package query

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/file"
	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/syntax"
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

const (
	ReloadConfig = "CONFIG"
)

const (
	ShowTables     = "TABLES"
	ShowViews      = "VIEWS"
	ShowCursors    = "CURSORS"
	ShowFunctions  = "FUNCTIONS"
	ShowStatements = "STATEMENTS"
	ShowFlags      = "FLAGS"
	ShowEnv        = "ENV"
	ShowRuninfo    = "RUNINFO"
)

var ShowObjectList = []string{
	ShowTables,
	ShowViews,
	ShowCursors,
	ShowFunctions,
	ShowStatements,
	ShowFlags,
	ShowEnv,
	ShowRuninfo,
}

func Echo(ctx context.Context, filter *Filter, expr parser.Echo) (string, error) {
	p, err := filter.Evaluate(ctx, expr.Value)
	if err != nil {
		return "", err
	}

	return NewStringFormatter().Format("%s", []value.Primary{p})
}

func Print(ctx context.Context, filter *Filter, expr parser.Print) (string, error) {
	p, err := filter.Evaluate(ctx, expr.Value)
	if err != nil {
		return "", err
	}
	return p.String(), err
}

func Printf(ctx context.Context, filter *Filter, expr parser.Printf) (string, error) {
	var format string
	formatValue, err := filter.Evaluate(ctx, expr.Format)
	if err != nil {
		return "", err
	}
	formatString := value.ToString(formatValue)
	if !value.IsNull(formatString) {
		format = formatString.(value.String).Raw()
	}

	args := make([]value.Primary, len(expr.Values))
	for i, v := range expr.Values {
		p, err := filter.Evaluate(ctx, v)
		if err != nil {
			return "", err
		}
		args[i] = p
	}

	message, err := NewStringFormatter().Format(format, args)
	if err != nil {
		return "", NewReplaceValueLengthError(expr, err.(Error).ErrorMessage())
	}
	return message, nil
}

func Source(ctx context.Context, filter *Filter, expr parser.Source) ([]parser.Statement, error) {
	var fpath string

	if ident, ok := expr.FilePath.(parser.Identifier); ok {
		fpath = ident.Literal
	} else {
		p, err := filter.Evaluate(ctx, expr.FilePath)
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

	return LoadStatementsFromFile(ctx, filter.tx, expr, fpath)
}

func LoadStatementsFromFile(ctx context.Context, tx *Transaction, expr parser.Source, fpath string) (statements []parser.Statement, err error) {
	if !filepath.IsAbs(fpath) {
		if abs, err := filepath.Abs(fpath); err == nil {
			fpath = abs
		}
	}

	if !file.Exists(fpath) {
		return nil, NewFileNotExistError(expr.FilePath)
	}

	h, err := file.NewHandlerForRead(ctx, tx.FileContainer, fpath, tx.WaitTimeout, tx.RetryDelay)
	if err != nil {
		return nil, NewReadFileError(expr, err.Error())
	}
	defer func() {
		err = AppendCompositeError(err, tx.FileContainer.Close(h))
	}()

	buf, err := ioutil.ReadAll(h.File())
	if err != nil {
		return nil, NewReadFileError(expr, err.Error())
	}
	input := string(buf)

	statements, _, err = parser.Parse(input, fpath, tx.Flags.DatetimeFormat, false)
	if err != nil {
		err = NewSyntaxError(err.(*parser.SyntaxError))
	}
	return statements, err
}

func ParseExecuteStatements(ctx context.Context, filter *Filter, expr parser.Execute) ([]parser.Statement, error) {
	var input string
	stmt, err := filter.Evaluate(ctx, expr.Statements)
	if err != nil {
		return nil, err
	}
	stmt = value.ToString(stmt)
	if !value.IsNull(stmt) {
		input = stmt.(value.String).Raw()
	}

	args := make([]value.Primary, len(expr.Values))
	for i, v := range expr.Values {
		p, err := filter.Evaluate(ctx, v)
		if err != nil {
			return nil, err
		}
		args[i] = p
	}

	input, err = NewStringFormatter().Format(input, args)
	if err != nil {
		return nil, NewReplaceValueLengthError(expr, err.(Error).ErrorMessage())
	}
	statements, _, err := parser.Parse(input, fmt.Sprintf("(L:%d C:%d) EXECUTE", expr.Line(), expr.Char()), filter.tx.Flags.DatetimeFormat, false)
	if err != nil {
		err = NewSyntaxError(err.(*parser.SyntaxError))
	}
	return statements, err
}

func SetFlag(ctx context.Context, filter *Filter, expr parser.SetFlag) error {
	var p value.Primary
	var err error

	if ident, ok := expr.Value.(parser.Identifier); ok {
		p = value.NewString(ident.Literal)
	} else {
		p, err = filter.Evaluate(ctx, expr.Value)
		if err != nil {
			return err
		}
	}

	switch strings.ToUpper(expr.Name) {
	case cmd.RepositoryFlag, cmd.TimezoneFlag, cmd.DatetimeFormatFlag,
		cmd.ImportFormatFlag, cmd.DelimiterFlag, cmd.DelimiterPositionsFlag, cmd.JsonQueryFlag, cmd.EncodingFlag,
		cmd.WriteEncodingFlag, cmd.FormatFlag, cmd.WriteDelimiterFlag, cmd.WriteDelimiterPositionsFlag, cmd.LineBreakFlag, cmd.JsonEscape:
		p = value.ToString(p)
	case cmd.NoHeaderFlag, cmd.WithoutNullFlag, cmd.WithoutHeaderFlag, cmd.EncloseAll, cmd.PrettyPrintFlag,
		cmd.EastAsianEncodingFlag, cmd.CountDiacriticalSignFlag, cmd.CountFormatCodeFlag, cmd.ColorFlag, cmd.QuietFlag, cmd.StatsFlag:
		p = value.ToBoolean(p)
	case cmd.WaitTimeoutFlag:
		p = value.ToFloat(p)
	case cmd.CPUFlag:
		p = value.ToInteger(p)
	default:
		return NewInvalidFlagNameError(expr, expr.Name)
	}
	if value.IsNull(p) {
		return NewFlagValueNotAllowedFormatError(expr)
	}

	switch strings.ToUpper(expr.Name) {
	case cmd.RepositoryFlag:
		err = filter.tx.Flags.SetRepository(p.(value.String).Raw())
	case cmd.TimezoneFlag:
		err = filter.tx.Flags.SetLocation(p.(value.String).Raw())
	case cmd.DatetimeFormatFlag:
		filter.tx.Flags.SetDatetimeFormat(p.(value.String).Raw())
	case cmd.WaitTimeoutFlag:
		filter.tx.UpdateWaitTimeout(p.(value.Float).Raw(), file.DefaultRetryDelay)
	case cmd.ImportFormatFlag:
		err = filter.tx.Flags.SetImportFormat(p.(value.String).Raw())
	case cmd.DelimiterFlag:
		err = filter.tx.Flags.SetDelimiter(p.(value.String).Raw())
	case cmd.DelimiterPositionsFlag:
		err = filter.tx.Flags.SetDelimiterPositions(p.(value.String).Raw())
	case cmd.JsonQueryFlag:
		filter.tx.Flags.SetJsonQuery(p.(value.String).Raw())
	case cmd.EncodingFlag:
		err = filter.tx.Flags.SetEncoding(p.(value.String).Raw())
	case cmd.NoHeaderFlag:
		filter.tx.Flags.SetNoHeader(p.(value.Boolean).Raw())
	case cmd.WithoutNullFlag:
		filter.tx.Flags.SetWithoutNull(p.(value.Boolean).Raw())
	case cmd.FormatFlag:
		err = filter.tx.Flags.SetFormat(p.(value.String).Raw(), "")
	case cmd.WriteEncodingFlag:
		err = filter.tx.Flags.SetWriteEncoding(p.(value.String).Raw())
	case cmd.WriteDelimiterFlag:
		err = filter.tx.Flags.SetWriteDelimiter(p.(value.String).Raw())
	case cmd.WriteDelimiterPositionsFlag:
		err = filter.tx.Flags.SetWriteDelimiterPositions(p.(value.String).Raw())
	case cmd.WithoutHeaderFlag:
		filter.tx.Flags.SetWithoutHeader(p.(value.Boolean).Raw())
	case cmd.LineBreakFlag:
		err = filter.tx.Flags.SetLineBreak(p.(value.String).Raw())
	case cmd.EncloseAll:
		filter.tx.Flags.SetEncloseAll(p.(value.Boolean).Raw())
	case cmd.JsonEscape:
		err = filter.tx.Flags.SetJsonEscape(p.(value.String).Raw())
	case cmd.PrettyPrintFlag:
		filter.tx.Flags.SetPrettyPrint(p.(value.Boolean).Raw())
	case cmd.EastAsianEncodingFlag:
		filter.tx.Flags.SetEastAsianEncoding(p.(value.Boolean).Raw())
	case cmd.CountDiacriticalSignFlag:
		filter.tx.Flags.SetCountDiacriticalSign(p.(value.Boolean).Raw())
	case cmd.CountFormatCodeFlag:
		filter.tx.Flags.SetCountFormatCode(p.(value.Boolean).Raw())
	case cmd.ColorFlag:
		filter.tx.Flags.SetColor(p.(value.Boolean).Raw())
	case cmd.QuietFlag:
		filter.tx.Flags.SetQuiet(p.(value.Boolean).Raw())
	case cmd.CPUFlag:
		filter.tx.Flags.SetCPU(int(p.(value.Integer).Raw()))
	case cmd.StatsFlag:
		filter.tx.Flags.SetStats(p.(value.Boolean).Raw())
	}

	if err != nil {
		return NewInvalidFlagValueError(expr, err.Error())
	}
	return nil
}

func AddFlagElement(ctx context.Context, filter *Filter, expr parser.AddFlagElement) error {
	switch strings.ToUpper(expr.Name) {
	case cmd.DatetimeFormatFlag:
		e := parser.SetFlag{
			BaseExpr: expr.GetBaseExpr(),
			Name:     expr.Name,
			Value:    expr.Value,
		}
		return SetFlag(ctx, filter, e)
	case cmd.RepositoryFlag, cmd.TimezoneFlag, cmd.DelimiterFlag, cmd.JsonQueryFlag, cmd.EncodingFlag,
		cmd.WriteEncodingFlag, cmd.FormatFlag, cmd.WriteDelimiterFlag, cmd.LineBreakFlag, cmd.JsonEscape,
		cmd.NoHeaderFlag, cmd.WithoutNullFlag, cmd.WithoutHeaderFlag, cmd.EncloseAll, cmd.PrettyPrintFlag,
		cmd.EastAsianEncodingFlag, cmd.CountDiacriticalSignFlag, cmd.CountFormatCodeFlag, cmd.ColorFlag, cmd.QuietFlag, cmd.StatsFlag,
		cmd.WaitTimeoutFlag,
		cmd.CPUFlag:

		return NewAddFlagNotSupportedNameError(expr)
	default:
		return NewInvalidFlagNameError(expr, expr.Name)
	}
}

func RemoveFlagElement(ctx context.Context, filter *Filter, expr parser.RemoveFlagElement) error {
	var p value.Primary
	var err error

	p, err = filter.Evaluate(ctx, expr.Value)
	if err != nil {
		return err
	}

	switch strings.ToUpper(expr.Name) {
	case cmd.DatetimeFormatFlag:
		if i := value.ToInteger(p); !value.IsNull(i) {
			idx := int(i.(value.Integer).Raw())
			if -1 < idx && idx < len(filter.tx.Flags.DatetimeFormat) {
				filter.tx.Flags.DatetimeFormat = append(filter.tx.Flags.DatetimeFormat[:idx], filter.tx.Flags.DatetimeFormat[idx+1:]...)
			}

		} else if s := value.ToString(p); !value.IsNull(s) {
			val := s.(value.String).Raw()
			formats := make([]string, 0, len(filter.tx.Flags.DatetimeFormat))
			for _, v := range filter.tx.Flags.DatetimeFormat {
				if val != v {
					formats = append(formats, v)
				}
			}
			filter.tx.Flags.DatetimeFormat = formats
		} else {
			return NewInvalidFlagValueToBeRemovedError(expr)
		}
	case cmd.RepositoryFlag, cmd.TimezoneFlag,
		cmd.ImportFormatFlag, cmd.DelimiterFlag, cmd.DelimiterPositionsFlag, cmd.JsonQueryFlag, cmd.EncodingFlag,
		cmd.WriteEncodingFlag, cmd.FormatFlag, cmd.WriteDelimiterFlag, cmd.WriteDelimiterPositionsFlag, cmd.LineBreakFlag, cmd.JsonEscape,
		cmd.NoHeaderFlag, cmd.WithoutNullFlag, cmd.WithoutHeaderFlag, cmd.EncloseAll, cmd.PrettyPrintFlag,
		cmd.EastAsianEncodingFlag, cmd.CountDiacriticalSignFlag, cmd.CountFormatCodeFlag, cmd.ColorFlag, cmd.QuietFlag, cmd.StatsFlag,
		cmd.WaitTimeoutFlag,
		cmd.CPUFlag:

		return NewRemoveFlagNotSupportedNameError(expr)
	default:
		return NewInvalidFlagNameError(expr, expr.Name)
	}

	return nil
}

func ShowFlag(flags *cmd.Flags, expr parser.ShowFlag) (string, error) {
	s, err := showFlag(flags, expr.Name)
	if err != nil {
		return s, NewInvalidFlagNameError(expr, expr.Name)
	}

	palette := cmd.GetPalette()
	return palette.Render(cmd.LableEffect, cmd.FlagSymbol(strings.ToUpper(expr.Name)+":")) + " " + s, nil
}

func showFlag(flags *cmd.Flags, flag string) (string, error) {
	var s string

	palette := cmd.GetPalette()

	switch strings.ToUpper(flag) {
	case cmd.RepositoryFlag:
		if len(flags.Repository) < 1 {
			wd, _ := os.Getwd()
			s = palette.Render(cmd.NullEffect, fmt.Sprintf("(current dir: %s)", wd))
		} else {
			s = palette.Render(cmd.StringEffect, flags.Repository)
		}
	case cmd.TimezoneFlag:
		s = palette.Render(cmd.StringEffect, flags.Location)
	case cmd.DatetimeFormatFlag:
		if len(flags.DatetimeFormat) < 1 {
			s = palette.Render(cmd.NullEffect, "(not set)")
		} else {
			list := make([]string, 0, len(flags.DatetimeFormat))
			for _, f := range flags.DatetimeFormat {
				list = append(list, "\""+f+"\"")
			}
			s = palette.Render(cmd.StringEffect, "["+strings.Join(list, ", ")+"]")
		}
	case cmd.WaitTimeoutFlag:
		s = palette.Render(cmd.NumberEffect, value.Float64ToStr(flags.WaitTimeout))
	case cmd.ImportFormatFlag:
		s = palette.Render(cmd.StringEffect, flags.ImportFormat.String())
	case cmd.DelimiterFlag:
		s = palette.Render(cmd.StringEffect, "'"+cmd.EscapeString(string(flags.Delimiter))+"'")
	case cmd.DelimiterPositionsFlag:
		p := fixedlen.DelimiterPositions(flags.DelimiterPositions).String()
		if flags.SingleLine {
			p = "S" + p
		}
		s = palette.Render(cmd.StringEffect, p)
	case cmd.JsonQueryFlag:
		if len(flags.JsonQuery) < 1 {
			s = palette.Render(cmd.NullEffect, "(empty)")
		} else {
			s = palette.Render(cmd.StringEffect, flags.JsonQuery)
		}
	case cmd.EncodingFlag:
		s = palette.Render(cmd.StringEffect, flags.Encoding.String())
	case cmd.NoHeaderFlag:
		s = palette.Render(cmd.BooleanEffect, strconv.FormatBool(flags.NoHeader))
	case cmd.WithoutNullFlag:
		s = palette.Render(cmd.BooleanEffect, strconv.FormatBool(flags.WithoutNull))
	case cmd.FormatFlag:
		s = palette.Render(cmd.StringEffect, flags.Format.String())
	case cmd.WriteEncodingFlag:
		switch flags.Format {
		case cmd.JSON:
			s = palette.Render(cmd.NullEffect, IgnoredFlagPrefix+flags.WriteEncoding.String())
		default:
			s = palette.Render(cmd.StringEffect, flags.WriteEncoding.String())
		}
	case cmd.WriteDelimiterFlag:
		s = "'" + cmd.EscapeString(string(flags.WriteDelimiter)) + "'"
		switch flags.Format {
		case cmd.CSV:
			s = palette.Render(cmd.StringEffect, s)
		default:
			s = palette.Render(cmd.NullEffect, IgnoredFlagPrefix+s)
		}
	case cmd.WriteDelimiterPositionsFlag:
		s = fixedlen.DelimiterPositions(flags.WriteDelimiterPositions).String()
		if flags.WriteAsSingleLine {
			s = "S" + s
		}
		switch flags.Format {
		case cmd.FIXED:
			s = palette.Render(cmd.StringEffect, s)
		default:
			s = palette.Render(cmd.NullEffect, IgnoredFlagPrefix+s)
		}
	case cmd.WithoutHeaderFlag:
		s = strconv.FormatBool(flags.WithoutHeader)
		switch flags.Format {
		case cmd.CSV, cmd.TSV, cmd.FIXED, cmd.GFM, cmd.ORG:
			if flags.Format == cmd.FIXED && flags.WriteAsSingleLine {
				s = palette.Render(cmd.NullEffect, IgnoredFlagPrefix+s)
			} else {
				s = palette.Render(cmd.BooleanEffect, s)
			}
		default:
			s = palette.Render(cmd.NullEffect, IgnoredFlagPrefix+s)
		}
	case cmd.LineBreakFlag:
		if flags.Format == cmd.FIXED && flags.WriteAsSingleLine {
			s = palette.Render(cmd.NullEffect, IgnoredFlagPrefix+flags.LineBreak.String())
		} else {
			s = palette.Render(cmd.StringEffect, flags.LineBreak.String())
		}
	case cmd.EncloseAll:
		s = strconv.FormatBool(flags.EncloseAll)
		switch flags.Format {
		case cmd.CSV, cmd.TSV:
			s = palette.Render(cmd.BooleanEffect, s)
		default:
			s = palette.Render(cmd.NullEffect, IgnoredFlagPrefix+s)
		}
	case cmd.JsonEscape:
		s = cmd.JsonEscapeTypeToString(flags.JsonEscape)
		switch flags.Format {
		case cmd.JSON:
			s = palette.Render(cmd.StringEffect, s)
		default:
			s = palette.Render(cmd.NullEffect, IgnoredFlagPrefix+s)
		}
	case cmd.PrettyPrintFlag:
		s = strconv.FormatBool(flags.PrettyPrint)
		switch flags.Format {
		case cmd.JSON:
			s = palette.Render(cmd.BooleanEffect, s)
		default:
			s = palette.Render(cmd.NullEffect, IgnoredFlagPrefix+s)
		}
	case cmd.EastAsianEncodingFlag:
		s = strconv.FormatBool(flags.EastAsianEncoding)
		switch flags.Format {
		case cmd.GFM, cmd.ORG, cmd.TEXT:
			s = palette.Render(cmd.BooleanEffect, s)
		default:
			s = palette.Render(cmd.NullEffect, IgnoredFlagPrefix+s)
		}
	case cmd.CountDiacriticalSignFlag:
		s = strconv.FormatBool(flags.CountDiacriticalSign)
		switch flags.Format {
		case cmd.GFM, cmd.ORG, cmd.TEXT:
			s = palette.Render(cmd.BooleanEffect, s)
		default:
			s = palette.Render(cmd.NullEffect, IgnoredFlagPrefix+s)
		}
	case cmd.CountFormatCodeFlag:
		s = strconv.FormatBool(flags.CountFormatCode)
		switch flags.Format {
		case cmd.GFM, cmd.ORG, cmd.TEXT:
			s = palette.Render(cmd.BooleanEffect, s)
		default:
			s = palette.Render(cmd.NullEffect, IgnoredFlagPrefix+s)
		}
	case cmd.ColorFlag:
		s = palette.Render(cmd.BooleanEffect, strconv.FormatBool(flags.Color))
	case cmd.QuietFlag:
		s = palette.Render(cmd.BooleanEffect, strconv.FormatBool(flags.Quiet))
	case cmd.CPUFlag:
		s = palette.Render(cmd.NumberEffect, strconv.Itoa(flags.CPU))
	case cmd.StatsFlag:
		s = palette.Render(cmd.BooleanEffect, strconv.FormatBool(flags.Stats))
	default:
		return s, errors.New("invalid flag name")
	}

	return s, nil
}

func ShowObjects(filter *Filter, expr parser.ShowObjects) (string, error) {
	var s string

	w := NewObjectWriter(filter.tx)

	switch strings.ToUpper(expr.Type.Literal) {
	case ShowTables:
		keys := filter.tx.cachedViews.SortedKeys()

		if len(keys) < 1 {
			s = cmd.Warn("No table is loaded")
		} else {
			createdFiles, updatedFiles := filter.tx.uncommittedViews.UncommittedFiles()

			for _, key := range keys {
				fields := filter.tx.cachedViews[key].Header.TableColumnNames()
				info := filter.tx.cachedViews[key].FileInfo
				ufpath := strings.ToUpper(info.Path)

				if _, ok := createdFiles[ufpath]; ok {
					w.WriteColor("*Created* ", cmd.EmphasisEffect)
				} else if _, ok := updatedFiles[ufpath]; ok {
					w.WriteColor("*Updated* ", cmd.EmphasisEffect)
				}
				w.WriteColorWithoutLineBreak(info.Path, cmd.ObjectEffect)
				writeFields(w, fields)

				w.NewLine()
				writeTableAttribute(w, filter.tx.Flags, info)
				w.ClearBlock()
				w.NewLine()
			}

			uncommitted := len(createdFiles) + len(updatedFiles)

			w.Title1 = "Loaded Tables"
			if 0 < uncommitted {
				w.Title2 = fmt.Sprintf("(Uncommitted: %s)", FormatCount(uncommitted, "Table"))
				w.Title2Effect = cmd.EmphasisEffect
			}
			s = "\n" + w.String() + "\n"
		}
	case ShowViews:
		views := filter.tempViews.All()

		if len(views) < 1 {
			s = cmd.Warn("No view is declared")
		} else {
			keys := views.SortedKeys()

			updatedViews := filter.tx.uncommittedViews.UncommittedTempViews()

			for _, key := range keys {
				fields := views[key].Header.TableColumnNames()
				info := views[key].FileInfo
				ufpath := strings.ToUpper(info.Path)

				if _, ok := updatedViews[ufpath]; ok {
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
			s = "\n" + w.String() + "\n"
		}
	case ShowCursors:
		cursors := filter.cursors.All()
		if len(cursors) < 1 {
			s = cmd.Warn("No cursor is declared")
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
				if cur.query.SelectEntity != nil {
					w.WriteColor("Query: ", cmd.LableEffect)
					writeQuery(w, cur.query.String())
				} else {
					w.WriteColorWithoutLineBreak("Statement: ", cmd.LableEffect)
					w.WriteColorWithoutLineBreak(cur.statement.String(), cmd.IdentifierEffect)
				}

				w.ClearBlock()
				w.NewLine()
			}
			w.Title1 = "Cursors"
			s = "\n" + w.String() + "\n"
		}
	case ShowFunctions:
		scalas, aggs := filter.functions.All()
		if len(scalas) < 1 && len(aggs) < 1 {
			s = cmd.Warn("No function is declared")
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
				s += "\n" + w.String() + "\n"
			} else {
				s += "\n"
			}
		}
	case ShowStatements:
		if len(filter.tx.PreparedStatements) < 1 {
			s = cmd.Warn("No statement is prepared")
		} else {
			keys := make([]string, 0, len(filter.tx.PreparedStatements))
			for k := range filter.tx.PreparedStatements {
				keys = append(keys, k)
			}
			sort.Strings(keys)

			for _, key := range keys {
				stmt := filter.tx.PreparedStatements[key]

				w.WriteColor(stmt.Name, cmd.ObjectEffect)
				w.BeginBlock()

				w.NewLine()
				w.WriteColorWithoutLineBreak("Placeholder Number: ", cmd.LableEffect)
				w.WriteColorWithoutLineBreak(strconv.Itoa(stmt.HolderNumber), cmd.NumberEffect)
				w.NewLine()
				w.WriteColorWithoutLineBreak("Statement: ", cmd.LableEffect)
				writeQuery(w, stmt.StatementString)

				w.ClearBlock()
				w.NewLine()
			}
			w.Title1 = "Prepared Statements"
			s = "\n" + w.String() + "\n"

		}
	case ShowFlags:
		for _, flag := range cmd.FlagList {
			symbol := cmd.FlagSymbol(flag)
			s, _ := showFlag(filter.tx.Flags, flag)
			w.WriteSpaces(27 - len(symbol))
			w.WriteColorWithoutLineBreak(symbol, cmd.LableEffect)
			w.WriteColorWithoutLineBreak(":", cmd.LableEffect)
			w.WriteSpaces(1)
			w.WriteWithoutLineBreak(s)
			w.NewLine()
		}
		w.Title1 = "Flags"
		s = "\n" + w.String() + "\n"
	case ShowEnv:
		env := os.Environ()
		names := make([]string, 0, len(env))
		vars := make([]string, 0, len(env))
		nameWidth := 0

		for _, e := range env {
			words := strings.Split(e, "=")
			name := string(parser.VariableSign) + string(parser.EnvironmentVariableSign) + words[0]
			if nameWidth < len(name) {
				nameWidth = len(name)
			}

			var val string
			if 1 < len(words) {
				val = strings.Join(words[1:], "=")
			}
			vars = append(vars, val)
			names = append(names, name)
		}

		for i, name := range names {
			w.WriteSpaces(nameWidth - len(name))
			w.WriteColorWithoutLineBreak(name, cmd.LableEffect)
			w.WriteColorWithoutLineBreak(":", cmd.LableEffect)
			w.WriteSpaces(1)
			w.WriteWithoutLineBreak(vars[i])
			w.NewLine()
		}
		w.Title1 = "Environment Variables"
		s = "\n" + w.String() + "\n"
	case ShowRuninfo:
		for _, ri := range RuntimeInformatinList {
			label := string(parser.VariableSign) + string(parser.RuntimeInformationSign) + ri
			p, _ := GetRuntimeInformation(filter.tx, parser.RuntimeInformation{Name: ri})

			w.WriteSpaces(19 - len(label))
			w.WriteColorWithoutLineBreak(label, cmd.LableEffect)
			w.WriteColorWithoutLineBreak(":", cmd.LableEffect)
			w.WriteSpaces(1)
			switch ri {
			case WorkingDirectory, VersionInformation:
				w.WriteColorWithoutLineBreak(p.(value.String).Raw(), cmd.StringEffect)
			case UncommittedInformation:
				w.WriteColorWithoutLineBreak(p.(value.Boolean).String(), cmd.BooleanEffect)
			default:
				w.WriteColorWithoutLineBreak(p.(value.Integer).String(), cmd.NumberEffect)
			}
			w.NewLine()
		}
		w.Title1 = "Runtime Information"
		s = "\n" + w.String() + "\n"
	default:
		return "", NewShowInvalidObjectTypeError(expr, expr.Type.String())
	}

	return s, nil
}

func writeTableAttribute(w *ObjectWriter, flags *cmd.Flags, info *FileInfo) {
	w.WriteColor("Format: ", cmd.LableEffect)
	w.WriteWithoutLineBreak(info.Format.String())

	w.WriteSpaces(9 - cmd.TextWidth(info.Format.String(), flags))
	switch info.Format {
	case cmd.CSV:
		w.WriteColorWithoutLineBreak("Delimiter: ", cmd.LableEffect)
		w.WriteWithoutLineBreak("'" + cmd.EscapeString(string(info.Delimiter)) + "'")
	case cmd.TSV:
		w.WriteColorWithoutLineBreak("Delimiter: ", cmd.LableEffect)
		w.WriteColorWithoutLineBreak("'\\t'", cmd.NullEffect)
	case cmd.FIXED:
		dp := info.DelimiterPositions.String()
		if info.SingleLine {
			dp = "S" + dp
		}

		w.WriteColorWithoutLineBreak("Delimiter Positions: ", cmd.LableEffect)
		w.WriteWithoutLineBreak(dp)
	case cmd.JSON:
		escapeStr := cmd.JsonEscapeTypeToString(info.JsonEscape)
		w.WriteColorWithoutLineBreak("Escape: ", cmd.LableEffect)
		w.WriteWithoutLineBreak(escapeStr)

		spaces := 9 - len(escapeStr)
		if spaces < 2 {
			spaces = 2
		}
		w.WriteSpaces(spaces)

		w.WriteColorWithoutLineBreak("Query: ", cmd.LableEffect)
		if len(info.JsonQuery) < 1 {
			w.WriteColorWithoutLineBreak("(empty)", cmd.NullEffect)
		} else {
			w.WriteColorWithoutLineBreak(info.JsonQuery, cmd.NullEffect)
		}
	}

	switch info.Format {
	case cmd.CSV, cmd.TSV:
		w.WriteSpaces(4 - (cmd.TextWidth(cmd.EscapeString(string(info.Delimiter)), flags)))
		w.WriteColorWithoutLineBreak("Enclose All: ", cmd.LableEffect)
		w.WriteWithoutLineBreak(strconv.FormatBool(info.EncloseAll))
	}

	w.NewLine()

	w.WriteColor("Encoding: ", cmd.LableEffect)
	switch info.Format {
	case cmd.JSON:
		w.WriteColorWithoutLineBreak(text.UTF8.String(), cmd.NullEffect)
	default:
		w.WriteWithoutLineBreak(info.Encoding.String())
	}

	if !(info.Format == cmd.FIXED && info.SingleLine) {
		w.WriteSpaces(7 - (cmd.TextWidth(info.Encoding.String(), flags)))
		w.WriteColorWithoutLineBreak("LineBreak: ", cmd.LableEffect)
		w.WriteWithoutLineBreak(info.LineBreak.String())
	}

	switch info.Format {
	case cmd.JSON:
		w.WriteSpaces(6 - (cmd.TextWidth(info.LineBreak.String(), flags)))
		w.WriteColorWithoutLineBreak("Pretty Print: ", cmd.LableEffect)
		w.WriteWithoutLineBreak(strconv.FormatBool(info.PrettyPrint))
	case cmd.CSV, cmd.TSV, cmd.FIXED, cmd.GFM, cmd.ORG:
		if !(info.Format == cmd.FIXED && info.SingleLine) {
			w.WriteSpaces(6 - (cmd.TextWidth(info.LineBreak.String(), flags)))
			w.WriteColorWithoutLineBreak("Header: ", cmd.LableEffect)
			w.WriteWithoutLineBreak(strconv.FormatBool(!info.NoHeader))
		}
	}
}

func writeFields(w *ObjectWriter, fields []string) {
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

func writeFunctions(w *ObjectWriter, funcs UserDefinedFunctionMap) {
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

func ShowFields(ctx context.Context, filter *Filter, expr parser.ShowFields) (string, error) {
	if !strings.EqualFold(expr.Type.Literal, "FIELDS") {
		return "", NewShowInvalidObjectTypeError(expr, expr.Type.Literal)
	}

	var status = ObjectFixed

	view := NewView(filter.tx)
	err := view.LoadFromTableIdentifier(ctx, filter.CreateNode(), expr.Table, false, false)
	if err != nil {
		return "", err
	}

	if view.FileInfo.IsTemporary {
		updatedViews := filter.tx.uncommittedViews.UncommittedTempViews()
		ufpath := strings.ToUpper(view.FileInfo.Path)

		if _, ok := updatedViews[ufpath]; ok {
			status = ObjectUpdated
		}
	} else {
		createdViews, updatedView := filter.tx.uncommittedViews.UncommittedFiles()
		ufpath := strings.ToUpper(view.FileInfo.Path)

		if _, ok := createdViews[ufpath]; ok {
			status = ObjectCreated
		} else if _, ok := updatedView[ufpath]; ok {
			status = ObjectUpdated
		}
	}

	w := NewObjectWriter(filter.tx)
	w.WriteColorWithoutLineBreak("Type: ", cmd.LableEffect)
	if view.FileInfo.IsTemporary {
		w.WriteWithoutLineBreak("View")
	} else {
		w.WriteWithoutLineBreak("Table")
		w.NewLine()
		w.WriteColorWithoutLineBreak("Path: ", cmd.LableEffect)
		w.WriteColorWithoutLineBreak(view.FileInfo.Path, cmd.ObjectEffect)
		w.NewLine()
		writeTableAttribute(w, filter.tx.Flags, view.FileInfo)
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
	if i, ok := expr.Table.(parser.Identifier); ok {
		w.Title2 = i.Literal
	} else if to, ok := expr.Table.(parser.TableObject); ok {
		w.Title2 = to.Path.Literal
	}
	w.Title2Effect = cmd.IdentifierEffect
	return "\n" + w.String() + "\n", nil
}

func writeFieldList(w *ObjectWriter, fields []string) {
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

func writeQuery(w *ObjectWriter, s string) {
	words := strings.Split(s, " ")

	w.BeginSubBlock()
	for _, v := range words {
		if !w.FitInLine(v + " ") {
			w.NewLine()
		}
		w.Write(v + " ")
	}
	w.EndSubBlock()
}

func SetEnvVar(ctx context.Context, filter *Filter, expr parser.SetEnvVar) error {
	var p value.Primary
	var err error

	if ident, ok := expr.Value.(parser.Identifier); ok {
		p = value.NewString(ident.Literal)
	} else {
		p, err = filter.Evaluate(ctx, expr.Value)
		if err != nil {
			return err
		}
	}

	var val string
	if p = value.ToString(p); !value.IsNull(p) {
		val = p.(value.String).Raw()
	}
	return os.Setenv(expr.EnvVar.Name, val)
}

func UnsetEnvVar(expr parser.UnsetEnvVar) error {
	return os.Unsetenv(expr.EnvVar.Name)
}

func Chdir(ctx context.Context, filter *Filter, expr parser.Chdir) error {
	var dirpath string
	var err error

	if ident, ok := expr.DirPath.(parser.Identifier); ok {
		dirpath = ident.Literal
	} else {
		p, err := filter.Evaluate(ctx, expr.DirPath)
		if err != nil {
			return err
		}
		s := value.ToString(p)
		if value.IsNull(s) {
			return NewInvalidPathError(expr, expr.DirPath.String(), "invalid directory path")
		}
		dirpath = s.(value.String).Raw()
	}

	if err = os.Chdir(dirpath); err != nil {
		if patherr, ok := err.(*os.PathError); ok {
			err = NewInvalidPathError(expr, patherr.Path, patherr.Err.Error())
		}
	}
	return err
}

func Pwd(expr parser.Pwd) (string, error) {
	dirpath, err := os.Getwd()
	if err != nil {
		if patherr, ok := err.(*os.PathError); ok {
			err = NewInvalidPathError(expr, patherr.Path, patherr.Err.Error())
		}
	}
	return dirpath, err
}

func Reload(ctx context.Context, tx *Transaction, expr parser.Reload) error {
	switch strings.ToUpper(expr.Type.Literal) {
	case ReloadConfig:
		if err := tx.Environment.Load(ctx, tx.WaitTimeout, tx.RetryDelay); err != nil {
			return NewLoadConfigurationError(expr, err.Error())
		}

		for _, v := range tx.Environment.DatetimeFormat {
			tx.Flags.DatetimeFormat = cmd.AppendStrIfNotExist(tx.Flags.DatetimeFormat, v)
		}

		palette, err := color.GeneratePalette(tx.Environment.Palette)
		if err != nil {
			return NewLoadConfigurationError(expr, err.Error())
		}
		oldPalette := cmd.GetPalette()
		oldPalette.Merge(palette)

		if tx.Session.Terminal != nil {
			if err := tx.Session.Terminal.ReloadConfig(); err != nil {
				return NewLoadConfigurationError(expr, err.Error())
			}
		}

	default:
		return NewInvalidReloadTypeError(expr, expr.Type.Literal)
	}
	return nil
}

func Syntax(ctx context.Context, filter *Filter, expr parser.Syntax) string {
	keys := make([]string, 0, len(expr.Keywords))
	for _, key := range expr.Keywords {
		var keystr string
		if fr, ok := key.(parser.FieldReference); ok {
			keystr = fr.Column.Literal
		} else {
			if p, err := filter.Evaluate(ctx, key); err == nil {
				if s := value.ToString(p); !value.IsNull(s) {
					keystr = s.(value.String).Raw()
				}
			}
		}

		if 0 < len(keystr) {
			words := strings.Split(strings.TrimSpace(keystr), " ")
			for _, w := range words {
				w = strings.TrimSpace(w)
				if 0 < len(w) {
					keys = append(keys, w)
				}
			}
		}
	}

	store := syntax.NewStore()
	exps := store.Search(keys)

	var p *color.Palette
	if filter.tx.Flags.Color {
		p = cmd.GetPalette()
	}

	w := NewObjectWriter(filter.tx)

	for _, exp := range exps {
		w.WriteColor(exp.Label, cmd.LableEffect)
		w.NewLine()
		if len(exps) < 4 {
			w.BeginBlock()

			if 0 < len(exp.Description.Template) {
				w.WriteWithAutoLineBreak(exp.Description.Format(p))
				w.NewLine()
				w.NewLine()
			}

			for _, def := range exp.Grammar {
				w.Write(def.Name.Format(p))
				w.NewLine()
				w.BeginBlock()
				for i, gram := range def.Group {
					if i == 0 {
						w.Write(": ")
					} else {
						w.Write("| ")
					}
					w.BeginSubBlock()
					w.WriteWithAutoLineBreak(gram.Format(p))
					w.EndSubBlock()
					w.NewLine()
				}

				if 0 < len(def.Description.Template) {
					if 0 < len(def.Group) {
						w.NewLine()
					}
					w.WriteWithAutoLineBreak(def.Description.Format(p))
					w.NewLine()
				}

				w.EndBlock()
				w.NewLine()
			}
			w.EndBlock()
		}

		if 0 < len(exp.Children) && (len(keys) < 1 || strings.EqualFold(exp.Label, strings.Join(keys, " "))) {
			w.BeginBlock()
			for _, child := range exp.Children {
				w.WriteColor(child.Label, cmd.LableEffect)
				w.NewLine()
			}
		}

		w.ClearBlock()
	}

	if len(keys) < 1 {
		w.Title1 = "Contents"
	} else {
		w.Title1 = "Search: " + strings.Join(keys, " ")
	}
	return "\n" + w.String() + "\n"

}
