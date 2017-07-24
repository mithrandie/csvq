package action

import (
	"errors"
	"strconv"
	"strings"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/query"
)

func Calc(expr string) error {
	cmd.SetNoHeader(true)

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

	filter := query.Filter{
		Records: []query.FilterRecord{
			{
				View:        view,
				RecordIndex: 0,
			},
		},
	}
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

func formatCalcResult(p parser.Primary) string {
	var s string

	switch p.(type) {
	case parser.String:
		s = strings.TrimSpace(p.(parser.String).Value())
	case parser.Integer:
		s = parser.Int64ToStr(p.(parser.Integer).Value())
	case parser.Float:
		s = parser.Float64ToStr(p.(parser.Float).Value())
	case parser.Boolean:
		s = strconv.FormatBool(p.(parser.Boolean).Value())
	case parser.Ternary:
		s = p.(parser.Ternary).String()
	case parser.Datetime:
		s = p.(parser.Datetime).Format()
	case parser.Null:
		s = "null"
	}

	return s
}
