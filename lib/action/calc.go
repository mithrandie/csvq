package action

import (
	"strconv"
	"strings"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/query"
)

func Calc(expr string) error {
	cmd.SetNoHeader(true)

	q := "select " + expr + " from stdin"

	parser.SetDebugLevel(0, true)
	program, err := parser.Parse(q)
	if err != nil {
		return err
	}
	selectEntity, _ := program[0].(parser.SelectQuery).SelectEntity.(parser.SelectEntity)

	view := query.NewView()
	err = view.Load(selectEntity.FromClause.(parser.FromClause), nil)
	if err != nil {
		return err
	}

	clause := selectEntity.SelectClause.(parser.SelectClause)

	var filter query.Filter = []query.FilterRecord{{View: view, RecordIndex: 0}}
	values := make([]string, len(clause.Fields))
	for i, v := range clause.Fields {
		field := v.(parser.Field)
		p, err := filter.Evaluate(field.Object)
		if err != nil {
			return err
		}
		values[i] = formatCalcResult(p)
	}

	flags := cmd.GetFlags()
	delimiter := flags.Delimiter
	if delimiter == cmd.UNDEF {
		delimiter = ','
	}

	cmd.ToStdout(strings.Join(values, ","))
	return nil
}

func formatCalcResult(p parser.Primary) string {
	var s string

	switch p.(type) {
	case parser.String:
		s = p.(parser.String).Value()
	case parser.Integer:
		s = parser.Int64ToStr(p.(parser.Integer).Value())
	case parser.Float:
		s = parser.Float64ToStr(p.(parser.Float).Value())
	case parser.Boolean:
		s = strconv.FormatBool(p.(parser.Boolean).Bool())
	case parser.Ternary:
		s = p.(parser.Ternary).String()
	case parser.Datetime:
		s = p.(parser.Datetime).Format()
	case parser.Null:
		s = "null"
	}

	return s
}
