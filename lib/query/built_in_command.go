package query

import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/mithrandie/csvq/lib/cmd"
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
	case "@@NO_HEADER", "@@WITHOUT_NULL":
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
	default:
		return s, NewInvalidFlagNameError(expr, expr.Name)
	}

	return s, nil
}
