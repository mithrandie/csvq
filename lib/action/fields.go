package action

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/csv"
	"github.com/mithrandie/csvq/lib/output"
	"github.com/mithrandie/csvq/lib/query"
)

func ShowFields(input string) error {
	fields, err := readFields(input)
	if err != nil {
		return err
	}

	out := formatFields(fields)
	output.Create("", out)
	return nil
}

func readFields(filename string) ([]string, error) {
	flags := cmd.GetFlags()

	fileInfo, err := query.NewFileInfo(filename, flags.Repository, flags.Delimiter)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(fileInfo.Path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := cmd.GetReader(f, flags.Encoding)

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
