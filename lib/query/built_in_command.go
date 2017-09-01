package query

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/value"
)

func Print(expr parser.Print, filter *Filter) (string, error) {
	p, err := filter.Evaluate(expr.Value)
	if err != nil {
		return "", err
	}
	return p.String(), err
}

func Printf(expr parser.Printf, filter *Filter) (string, error) {
	args := make([]value.Primary, len(expr.Values))
	for i, v := range expr.Values {
		p, err := filter.Evaluate(v)
		if err != nil {
			return "", err
		}
		args[i] = p
	}

	message, err := FormatString(expr.Format, args)
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
	s := value.PrimaryToString(p)
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

	fp, err := os.Open(fpath)
	if err != nil {
		return nil, NewReadFileError(expr, err.Error())
	}
	defer fp.Close()

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
	var err error

	var p value.Primary

	switch strings.ToUpper(expr.Name) {
	case "@@DELIMITER", "@@ENCODING", "@@LINE_BREAK", "@@REPOSITORY", "@@DATETIME_FORMAT":
		p = value.PrimaryToString(expr.Value)
	case "@@NO_HEADER", "@@WITHOUT_NULL":
		p = value.PrimaryToBoolean(expr.Value)
	default:
		return NewInvalidFlagNameError(expr)
	}
	if value.IsNull(p) {
		return NewInvalidFlagValueError(expr)
	}

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
