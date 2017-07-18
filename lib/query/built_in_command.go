package query

import (
	"errors"
	"fmt"
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

	p := parser.PrimaryToString(args[0])
	if parser.IsNull(p) {
		return "", nil
	}
	format := []rune(p.(parser.String).Value())
	str := []rune{}

	escaped := false
	placeholderOrder := 0
	for _, r := range format {
		if escaped {
			switch r {
			case 's':
				placeholderOrder++
				if len(args) <= placeholderOrder {
					return "", errors.New(fmt.Sprintf("print format %q: number of replace values does not match", string(format)))
				}
				str = append(str, []rune(args[placeholderOrder].String())...)
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

	if placeholderOrder < len(args)-1 {
		return "", errors.New(fmt.Sprintf("print format %q: number of replace values does not match", string(format)))
	}

	return string(str), nil
}

func Source(expr parser.Source) ([]parser.Statement, error) {
	stat, err := os.Stat(expr.FilePath)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("source file %q does not exist", expr.FilePath))
	}
	if stat.IsDir() {
		return nil, errors.New(fmt.Sprintf("source file %q must be a readable file", expr.FilePath))
	}

	fp, err := os.Open(expr.FilePath)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	buf, err := ioutil.ReadAll(fp)
	if err != nil {
		return nil, err
	}
	input := string(buf)

	return parser.Parse(input)
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
		return errors.New(fmt.Sprintf("invalid flag name: %s", expr.Name))
	}
	if parser.IsNull(p) {
		return errors.New(fmt.Sprintf("invalid flag value: %s = %s", expr.Name, expr.Value))
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
	return err
}
