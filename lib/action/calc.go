package action

import (
	"context"
	"errors"
	"strings"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/query"
)

func Calc(expr string) error {
	cmd.GetFlags().SetNoHeader(true)

	defer func() {
		if err := query.ReleaseResourcesWithErrors(); err != nil {
			query.LogError(err.Error())
		}
	}()

	ctx := context.Background()

	q := "select " + expr + " from stdin"

	program, err := parser.Parse(q, "")
	if err != nil {
		return errors.New("syntax error")
	}
	selectEntity, _ := program[0].(parser.SelectQuery).SelectEntity.(parser.SelectEntity)

	view := query.NewView()
	err = view.Load(ctx, selectEntity.FromClause.(parser.FromClause), query.NewEmptyFilter().CreateNode())
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
		p, err := filter.Evaluate(ctx, field.Object)
		if err != nil {
			return errors.New("syntax error")
		}
		values[i], _, _ = query.ConvertFieldContents(p, true)
	}

	return query.WriteToStdout(strings.Join(values, string(cmd.GetFlags().WriteDelimiter)))
}
