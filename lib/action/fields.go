package action

import (
	"fmt"
	"strconv"
	"strings"

	"errors"
	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/csv"
	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/query"
	"github.com/mithrandie/go-file"
)

func ShowFields(input string) error {
	SetSignalHandler()

	fields, err := readFields(input)
	if err != nil {
		return err
	}

	out := formatFields(fields)
	cmd.ToStdout(out)
	return nil
}

func readFields(filename string) ([]string, error) {
	flags := cmd.GetFlags()

	query.UpdateWaitTimeout()
	fileInfo, err := query.NewFileInfo(parser.Identifier{Literal: filename}, flags.Repository, flags.Delimiter)
	if err != nil {
		if appErr, ok := err.(query.AppError); ok {
			return nil, errors.New(appErr.ErrorMessage())
		}
		return nil, errors.New(err.Error())
	}

	fp, err := file.OpenToRead(fileInfo.Path)
	if err != nil {
		return nil, err
	}
	defer file.Close(fp)

	r := cmd.GetReader(fp, flags.Encoding)

	reader := csv.NewReader(r)
	reader.Delimiter = fileInfo.Delimiter
	reader.WithoutNull = flags.WithoutNull

	fields, err := reader.ReadHeader()
	if err != nil {
		return nil, err
	}

	return fields, nil
}

func formatFields(fields []string) string {
	l := len(fields)
	digits := len(strconv.Itoa(l))
	formatted := make([]string, l)

	for i, field := range fields {
		idxstr := strconv.Itoa(i + 1)
		formatted[i] = fmt.Sprintf("%"+strconv.Itoa(digits)+"s. %s", idxstr, field)
	}

	return strings.Join(formatted, "\n") + "\n"
}
