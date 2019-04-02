package action

import (
	"context"
	"errors"
	"strings"

	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/query"
)

func Calc(ctx context.Context, proc *query.Processor, expr string) error {
	proc.Tx.Flags.SetNoHeader(true)
	q := "SELECT " + expr + " FROM STDIN"

	program, _, err := parser.Parse(q, "", proc.Tx.Flags.DatetimeFormat, false)
	if err != nil {
		e := err.(*parser.SyntaxError)
		e.SourceFile = ""
		e.Line = 0
		e.Char = 0
		e.Message = "syntax error"
		return query.NewSyntaxError(e)
	}
	selectEntity, _ := program[0].(parser.SelectQuery).SelectEntity.(parser.SelectEntity)

	view := query.NewView(proc.Tx)
	err = view.Load(ctx, query.NewFilter(proc.Tx).CreateNode(), selectEntity.FromClause.(parser.FromClause), false, false)
	if err != nil {
		if appErr, ok := err.(query.Error); ok {
			err = errors.New(appErr.Message())
		}
		return err
	}

	clause := selectEntity.SelectClause.(parser.SelectClause)

	filter := query.NewFilterForRecord(proc.Filter, view, 0)
	values := make([]string, len(clause.Fields))
	for i, v := range clause.Fields {
		field := v.(parser.Field)
		p, err := filter.Evaluate(ctx, field.Object)
		if err != nil {
			if appErr, ok := err.(query.Error); ok {
				err = errors.New(appErr.Message())
			}
			return err
		}
		values[i], _, _ = query.ConvertFieldContents(p, true)
	}

	return proc.Tx.Session.WriteToStdout(strings.Join(values, string(proc.Tx.Flags.WriteDelimiter)))
}
