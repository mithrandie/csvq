package query

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/csv"
	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/value"

	"github.com/mithrandie/go-file"
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
	case "@@DELIMITER", "@@ENCODING", "@@LINE_BREAK", "@@REPOSITORY", "@@DATETIME_FORMAT":
		p = value.ToString(expr.Value)
	case "@@WAIT_TIMEOUT":
		p = value.ToFloat(expr.Value)
	case "@@NO_HEADER", "@@WITHOUT_NULL", "@@STATS":
		p = value.ToBoolean(expr.Value)
	default:
		return NewInvalidFlagNameError(expr, expr.Name)
	}
	if value.IsNull(p) {
		return NewInvalidFlagValueError(expr)
	}

	var err error

	switch strings.ToUpper(expr.Name) {
	case "@@DELIMITER":
		err = cmd.SetDelimiter(p.(value.String).Raw())
	case "@@ENCODING":
		err = cmd.SetEncoding(p.(value.String).Raw())
	case "@@LINE_BREAK":
		err = cmd.SetLineBreak(p.(value.String).Raw())
	case "@@REPOSITORY":
		err = cmd.SetRepository(p.(value.String).Raw())
	case "@@DATETIME_FORMAT":
		cmd.SetDatetimeFormat(p.(value.String).Raw())
	case "@@WAIT_TIMEOUT":
		err = cmd.SetWaitTimeout(value.Float64ToStr(p.(value.Float).Raw()))
		if err == nil {
			UpdateWaitTimeout()
		}
	case "@@NO_HEADER":
		cmd.SetNoHeader(p.(value.Boolean).Raw())
	case "@@WITHOUT_NULL":
		cmd.SetWithoutNull(p.(value.Boolean).Raw())
	case "@@STATS":
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
	case "@@DELIMITER":
		if flags.Delimiter == cmd.UNDEF {
			s = "(not set)"
		} else {
			s = "'" + cmd.EscapeString(string(flags.Delimiter)) + "'"
		}
	case "@@ENCODING":
		s = flags.Encoding.String()
	case "@@LINE_BREAK":
		s = flags.LineBreak.String()
	case "@@REPOSITORY":
		s = flags.Repository
	case "@@DATETIME_FORMAT":
		if len(flags.DatetimeFormat) < 1 {
			s = "(not set)"
		} else {
			s = flags.DatetimeFormat
		}
	case "@@WAIT_TIMEOUT":
		s = value.Float64ToStr(flags.WaitTimeout)
	case "@@NO_HEADER":
		s = strconv.FormatBool(flags.NoHeader)
	case "@@WITHOUT_NULL":
		s = strconv.FormatBool(flags.WithoutNull)
	case "@@STATS":
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
			case strings.ToUpper(cmd.CSV_EXT), strings.ToUpper(cmd.TSV_EXT):
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
			s = fmt.Sprintf("Repository %q is empty", repository)
		} else {
			if 0 < len(filePaths) {
				s += "\n"
				s += fmt.Sprintf("    Tables in %s\n", repository)
				s += strings.Repeat("-", len(repository)+18) + "\n"
				s += strings.Join(filePaths, "\n")
				s += "\n"
			}
			if 0 < len(cachedPaths) {
				s += "\n"
				s += "    Tables in other directories\n"
				s += "-----------------------------------\n"
				s += strings.Join(cachedPaths, "\n")
				s += "\n"
			}
		}
	case parser.VIEWS:
		views := filter.TempViews.List()
		if len(views) < 1 {
			s = "No view is declared"
		} else {
			s += "\n"
			s += "    Views\n"
			s += "-------------\n"
			s += strings.Join(views, "\n")
			s += "\n"
		}
	case parser.CURSORS:
		cursors := filter.Cursors.List()
		if len(cursors) < 1 {
			s = "No cursor is declared"
		} else {
			s += "\n"
			s += "    Cursors\n"
			s += "---------------\n"
			s += strings.Join(cursors, "\n")
			s += "\n"
		}
	case parser.FUNCTIONS:
		scalas, aggs := filter.Functions.List()
		if len(scalas) < 1 && len(aggs) < 1 {
			s = "No function is declared"
		} else {
			if 0 < len(scalas) {
				s += "\n"
				s += "    Scala Functions\n"
				s += "-----------------------\n"
				s += strings.Join(scalas, "\n")
				s += "\n"
			}
			if 0 < len(aggs) {
				s += "\n"
				s += "    Aggregate Functions\n"
				s += "---------------------------\n"
				s += strings.Join(aggs, "\n")
				s += "\n"
			}
		}
	}

	return s, nil
}

func ShowFields(expr parser.ShowFields, filter *Filter) (string, error) {
	var fields []string

	if filter.TempViews.Exists(expr.Table.Literal) {
		view, _ := filter.TempViews.Get(expr.Table)
		fields = view.Header.TableColumnNames()
	} else {
		flags := cmd.GetFlags()

		fileInfo, err := NewFileInfoForCreate(expr.Table, flags.Repository, flags.Delimiter)
		if err != nil {
			return "", err
		}

		if ViewCache.Exists(fileInfo.Path) {
			pathIdent := parser.Identifier{Literal: fileInfo.Path}
			view, _ := ViewCache.Get(pathIdent)
			fields = view.Header.TableColumnNames()
		} else {
			fileInfo, err = NewFileInfo(expr.Table, flags.Repository, flags.Delimiter)
			if err != nil {
				return "", err
			}

			if ViewCache.Exists(fileInfo.Path) {
				pathIdent := parser.Identifier{Literal: fileInfo.Path}
				view, _ := ViewCache.Get(pathIdent)
				fields = view.Header.TableColumnNames()
			} else {
				if !FileLocks.CanRead(fileInfo.Path) {
					return "", NewFileLockTimeoutError(expr.Table, fileInfo.Path)
				}
				fp, err := file.OpenToReadWithTimeout(fileInfo.Path)
				if err != nil {
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

	var s string
	s += "\n"
	s += fmt.Sprintf("    Fields in %s\n", expr.Table.Literal)
	s += strings.Repeat("-", len(expr.Table.Literal)+18) + "\n"
	s += formatFields(fields)
	s += "\n"

	return s, nil
}

func formatFields(fields []string) string {
	l := len(fields)
	digits := len(strconv.Itoa(l))
	formatted := make([]string, l)

	for i, field := range fields {
		idxstr := strconv.Itoa(i + 1)
		formatted[i] = fmt.Sprintf("%"+strconv.Itoa(digits)+"s. %s", idxstr, field)
	}

	return strings.Join(formatted, "\n")
}
