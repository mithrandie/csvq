package action

import (
	"context"
	"errors"
	"strings"

	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/query"
)

func Calc(proc *query.Processor, expr string) error {
	proc.Tx.Flags.SetNoHeader(true)

	defer func() {
		if err := proc.ReleaseResourcesWithErrors(); err != nil {
			proc.LogError(err.Error())
		}
	}()

	ctx := context.Background()

	q := "select " + expr + " from stdin"

	program, err := parser.Parse(q, "", proc.Tx.Flags.DatetimeFormat)
	if err != nil {
		return errors.New("syntax error")
	}
	selectEntity, _ := program[0].(parser.SelectQuery).SelectEntity.(parser.SelectEntity)

	view := query.NewView(proc.Tx)
	err = view.Load(ctx, query.NewFilter(proc.Tx).CreateNode(), selectEntity.FromClause.(parser.FromClause))
	if err != nil {
		if appErr, ok := err.(query.Error); ok {
			return errors.New(appErr.ErrorMessage())
		}
		return err
	}

	clause := selectEntity.SelectClause.(parser.SelectClause)

	filter := query.NewFilterForRecord(query.NewFilter(proc.Tx), view, 0)
	values := make([]string, len(clause.Fields))
	for i, v := range clause.Fields {
		field := v.(parser.Field)
		p, err := filter.Evaluate(ctx, field.Object)
		if err != nil {
			return errors.New("syntax error")
		}
		values[i], _, _ = query.ConvertFieldContents(p, true)
	}

	return proc.Tx.Session.WriteToStdout(strings.Join(values, string(proc.Tx.Flags.WriteDelimiter)))
}
