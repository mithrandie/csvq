package query

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/parser"
)

func Print(expr parser.Print, filter Filter) (string, error) {
	p, err := filter.Evaluate(expr.Value)
	if err != nil {
		return "", err
	}
	return p.String(), err
}

func Printf(expr parser.Printf, filter Filter) (string, error) {
	args := make([]parser.Primary, len(expr.Values))
	for i, v := range expr.Values {
		p, err := filter.Evaluate(v)
		if err != nil {
			return "", err
		}
		args[i] = p
	}

	format := expr.Format
	str := []rune{}

	escaped := false
	placeholderOrder := 0
	for _, r := range format {
		if escaped {
			switch r {
			case 's':
				if len(args) <= placeholderOrder {
					return "", NewPrintfReplaceValueLengthError(expr)
				}
				str = append(str, []rune(args[placeholderOrder].String())...)
				placeholderOrder++
			case '%':
				str = append(str, r)
			default:
				str = append(str, '%', r)
			}
			escaped = false
			continue
		}

		if r == '%' {
			escaped = true
			continue
		}

		str = append(str, r)
	}
	if escaped {
		str = append(str, '%')
	}

	if placeholderOrder < len(args) {
		return "", NewPrintfReplaceValueLengthError(expr)
	}

	return string(str), nil
}

func Source(expr parser.Source, filter Filter) ([]parser.Statement, error) {
	p, err := filter.Evaluate(expr.FilePath)
	if err != nil {
		return nil, err
	}
	s := parser.PrimaryToString(p)
	if parser.IsNull(s) {
		return nil, NewSourceInvalidArgumentError(expr, expr.FilePath)
	}
	fpath := s.(parser.String).Value()

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

	var p parser.Primary

	switch strings.ToUpper(expr.Name) {
	case "@@DELIMITER", "@@ENCODING", "@@LINE_BREAK", "@@REPOSITORY", "@@DATETIME_FORMAT":
		p = parser.PrimaryToString(expr.Value)
	case "@@NO_HEADER", "@@WITHOUT_NULL":
		p = parser.PrimaryToBoolean(expr.Value)
	default:
		return NewInvalidFlagNameError(expr)
	}
	if parser.IsNull(p) {
		return NewInvalidFlagValueError(expr)
	}

	switch strings.ToUpper(expr.Name) {
	case "@@DELIMITER":
		err = cmd.SetDelimiter(p.(parser.String).Value())
	case "@@ENCODING":
		err = cmd.SetEncoding(p.(parser.String).Value())
	case "@@LINE_BREAK":
		err = cmd.SetLineBreak(p.(parser.String).Value())
	case "@@REPOSITORY":
		err = cmd.SetRepository(p.(parser.String).Value())
	case "@@DATETIME_FORMAT":
		cmd.SetDatetimeFormat(p.(parser.String).Value())
	case "@@NO_HEADER":
		cmd.SetNoHeader(p.(parser.Boolean).Value())
	case "@@WITHOUT_NULL":
		cmd.SetWithoutNull(p.(parser.Boolean).Value())
	}

	if err != nil {
		return NewInvalidFlagValueError(expr)
	}

	return nil
}
