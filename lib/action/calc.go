package action

import (
	"context"
	"errors"
	"strings"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/query"
)

func Calc(ctx context.Context, proc *query.Processor, expr string) error {
	_ = proc.Tx.SetFlag(cmd.NoHeaderFlag, true)
	q := "SELECT " + expr + " FROM STDIN"

	program, _, err := parser.Parse(q, "", false, proc.Tx.Flags.AnsiQuotes)
	if err != nil {
		e := err.(*parser.SyntaxError)
		e.SourceFile = ""
		e.Line = 0
		e.Char = 0
		e.Message = "syntax error"
		return query.NewSyntaxError(e)
	}
	selectEntity, _ := program[0].(parser.SelectQuery).SelectEntity.(parser.SelectEntity)

	scope := query.NewReferenceScope(proc.Tx)
	queryScope := scope.CreateNode()

	view, err := query.LoadView(ctx, queryScope, selectEntity.FromClause.(parser.FromClause).Tables, false, false)
	if err != nil {
		if appErr, ok := err.(query.Error); ok {
			err = errors.New(appErr.Message())
		}
		return err
	}

	clause := selectEntity.SelectClause.(parser.SelectClause)

	recordScope := scope.CreateScopeForRecordEvaluation(view, 0)
	values := make([]string, len(clause.Fields))
	for i, v := range clause.Fields {
		field := v.(parser.Field)
		p, err := query.Evaluate(ctx, recordScope, field.Object)
		if err != nil {
			if appErr, ok := err.(query.Error); ok {
				err = errors.New(appErr.Message())
			}
			return err
		}
		values[i], _, _ = query.ConvertFieldContents(p, true)
	}

	return proc.Tx.Session.WriteToStdout(strings.Join(values, string(proc.Tx.Flags.ExportOptions.Delimiter)))
}
