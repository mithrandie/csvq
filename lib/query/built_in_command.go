package query

import (
	"fmt"
	"github.com/mithrandie/csvq/lib/color"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/csv"
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
	WaitTimeoutFlag    = "@@WAIT_TIMEOUT"
	NoHeaderFlag       = "@@NO_HEADER"
	WithoutNullFlag    = "@@WITHOUT_NULL"
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
	p, err := filter.Evaluate(expr.FilePath)
	if err != nil {
		return nil, err
	}
	s := value.ToString(p)
	if value.IsNull(s) {
		return nil, NewSourceInvalidArgumentError(expr, expr.FilePath)
	}
	fpath := s.(value.String).Raw()

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
	case WaitTimeoutFlag:
		cmd.SetWaitTimeout(p.(value.Float).Raw())
	case NoHeaderFlag:
		cmd.SetNoHeader(p.(value.Boolean).Raw())
	case WithoutNullFlag:
		cmd.SetWithoutNull(p.(value.Boolean).Raw())
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
	var s string

	flags := cmd.GetFlags()

	switch strings.ToUpper(expr.Name) {
	case DelimiterFlag:
		if flags.Delimiter == cmd.UNDEF {
			s = "(not set)"
		} else {
			s = "'" + cmd.EscapeString(string(flags.Delimiter)) + "'"
		}
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
	case WaitTimeoutFlag:
		s = value.Float64ToStr(flags.WaitTimeout)
	case NoHeaderFlag:
		s = strconv.FormatBool(flags.NoHeader)
	case WithoutNullFlag:
		s = strconv.FormatBool(flags.WithoutNull)
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
		return s, NewInvalidFlagNameError(expr, expr.Name)
	}

	return s, nil
}

func ShowObjects(expr parser.ShowObjects, filter *Filter) (string, error) {
	var s string

	switch expr.Type {
	case parser.TABLES:
		repository := cmd.GetFlags().Repository

		files, err := ioutil.ReadDir(repository)
		if err != nil {
			return "", err
		}

		var filePaths []string
		var absPaths []string
		for _, f := range files {
			if f.IsDir() {
				continue
			}

			ext := filepath.Ext(f.Name())
			switch strings.ToUpper(ext) {
			case strings.ToUpper(cmd.CsvExt), strings.ToUpper(cmd.TsvExt):
				absPath := filepath.Join(repository, f.Name())
				if !filepath.IsAbs(absPath) {
					p, err := filepath.Abs(absPath)
					if err != nil {
						return "", err
					}
					absPath = p
				}
				filePaths = append(filePaths, f.Name())
				absPaths = append(absPaths, absPath)
			}
		}
		sort.Strings(filePaths)

		var cachedPaths []string
		for _, v := range ViewCache {
			cachedPath := v.FileInfo.Path
			if !InStrSlice(cachedPath, absPaths) {
				cachedPaths = append(cachedPaths, cachedPath)
			}
		}
		sort.Strings(cachedPaths)

		if len(filePaths) < 1 && len(cachedPaths) < 1 {
			s = color.Warn(fmt.Sprintf("Repository %q is empty", repository))
		} else {
			if 0 < len(filePaths) {
				s += formatHeader("Tables in ", repository) + strings.Join(filePaths, "\n") + "\n"
			}
			if 0 < len(cachedPaths) {
				s += formatHeader("Tables in other directories", "") + strings.Join(cachedPaths, "\n") + "\n"
			}
		}
	case parser.VIEWS:
		views := filter.TempViews.List()
		if len(views) < 1 {
			s = color.Warn("No view is declared")
		} else {
			s = formatHeader("Views", "") + strings.Join(views, "\n") + "\n"
		}
	case parser.CURSORS:
		cursors := filter.Cursors.List()
		if len(cursors) < 1 {
			s = color.Warn("No cursor is declared")
		} else {
			s = formatHeader("Cursors", "") + strings.Join(cursors, "\n") + "\n"
		}
	case parser.FUNCTIONS:
		scalas, aggs := filter.Functions.List()
		if len(scalas) < 1 && len(aggs) < 1 {
			s = color.Warn("No function is declared")
		} else {
			if 0 < len(scalas) {
				s += formatHeader("Scala Functions", "") + strings.Join(scalas, "\n") + "\n"
			}
			if 0 < len(aggs) {
				s += formatHeader("Aggregate Functions", "") + strings.Join(aggs, "\n") + "\n"
			}
		}
	}

	return s, nil
}

func ShowFields(expr parser.ShowFields, filter *Filter) (string, error) {
	var fields []string

	if filter.TempViews.Exists(expr.Table.Literal) {
		header, _ := filter.TempViews.GetHeader(expr.Table)
		fields = header.TableColumnNames()
	} else {
		flags := cmd.GetFlags()

		fileInfo, err := NewFileInfoForCreate(expr.Table, flags.Repository, flags.Delimiter)
		if err != nil {
			return "", err
		}

		if ViewCache.Exists(fileInfo.Path) {
			pathIdent := parser.Identifier{Literal: fileInfo.Path}
			header, _ := ViewCache.GetHeader(pathIdent)
			fields = header.TableColumnNames()
		} else {
			fileInfo, err = NewFileInfo(expr.Table, flags.Repository, flags.Delimiter, cmd.CSV)
			if err != nil {
				return "", err
			}

			if ViewCache.Exists(fileInfo.Path) {
				pathIdent := parser.Identifier{Literal: fileInfo.Path}
				header, _ := ViewCache.GetHeader(pathIdent)
				fields = header.TableColumnNames()
			} else {
				fp, err := file.OpenToRead(fileInfo.Path)
				if err != nil {
					if _, ok := err.(*file.TimeoutError); ok {
						NewFileLockTimeoutError(expr.Table, fileInfo.Path)
					}
					return "", NewReadFileError(expr.Table, err.Error())
				}
				defer file.Close(fp)

				r := cmd.GetReader(fp, flags.Encoding)
				reader := csv.NewReader(r)
				reader.Delimiter = fileInfo.Delimiter
				reader.WithoutNull = flags.WithoutNull

				header, err := reader.ReadHeader()
				if err != nil && err != csv.EOF {
					return "", err
				}
				fields = header
			}
		}
	}

	s := formatHeader("Fields in ", expr.Table.Literal) + formatFields(fields)

	return s, nil
}

func formatFields(fields []string) string {
	l := len(fields)
	digits := len(strconv.Itoa(l))
	formatted := make([]string, l)

	for i, field := range fields {
		idxstr := strconv.Itoa(i + 1)
		formatted[i] = color.MagentaB(fmt.Sprintf("%"+strconv.Itoa(digits)+"s", idxstr)) + ". " + color.Cyan(field)
	}

	return strings.Join(formatted, "\n") + "\n"
}

func formatHeader(title string, colorItem string) string {
	colorItem = color.CyanB(colorItem)

	return "\n    " +
		title + colorItem + "\n" +
		strings.Repeat("-", len(title)+len(color.StripEscapeSequence(colorItem))+8) + "\n"
}
