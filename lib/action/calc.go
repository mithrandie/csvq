package action

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/query"
	"github.com/mithrandie/csvq/lib/value"
)

func Calc(expr string) error {
	cmd.SetNoHeader(true)

	SetSignalHandler()

	defer func() {
		query.ReleaseResources()
	}()

	q := "select " + expr + " from stdin"

	program, err := parser.Parse(q, "")
	if err != nil {
		return errors.New("syntax error")
	}
	selectEntity, _ := program[0].(parser.SelectQuery).SelectEntity.(parser.SelectEntity)

	view := query.NewView()
	err = view.Load(selectEntity.FromClause.(parser.FromClause), query.NewEmptyFilter().CreateNode())
	if err != nil {
		if appErr, ok := err.(query.AppError); ok {
			return errors.New(appErr.ErrorMessage())
		}
		return err
	}

	clause := selectEntity.SelectClause.(parser.SelectClause)

	filter := query.NewFilterForRecord(view, 0, query.NewEmptyFilter())
	values := make([]string, len(clause.Fields))
	for i, v := range clause.Fields {
		field := v.(parser.Field)
		p, err := filter.Evaluate(field.Object)
		if err != nil {
			return errors.New("syntax error")
		}
		values[i] = formatCalcResult(p)
	}

	flags := cmd.GetFlags()
	delimiter := flags.Delimiter
	if delimiter == cmd.UNDEF {
		delimiter = ','
	}

	cmd.ToStdout(strings.Join(values, string(delimiter)))
	return nil
}

func formatCalcResult(p value.Primary) string {
	var s string

	switch p.(type) {
	case value.String:
		s = strings.TrimSpace(p.(value.String).Raw())
	case value.Integer:
		s = value.Int64ToStr(p.(value.Integer).Raw())
	case value.Float:
		s = value.Float64ToStr(p.(value.Float).Raw())
	case value.Boolean:
		s = strconv.FormatBool(p.(value.Boolean).Raw())
	case value.Ternary:
		s = p.(value.Ternary).String()
	case value.Datetime:
		s = p.(value.Datetime).Format(time.RFC3339Nano)
	case value.Null:
		s = "null"
	}

	return s
}
