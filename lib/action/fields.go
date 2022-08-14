package action

import (
	"context"
	"os"

	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/query"
)

func ShowFields(ctx context.Context, proc *query.Processor, filename string) error {
	if len(filename) < 1 {
		e := parser.NewSyntaxError("file name is empty", parser.Token{})
		return query.NewSyntaxError(e.(*parser.SyntaxError))
	}

	var filePath parser.QueryExpression = parser.Identifier{Literal: filename}
	path, err := query.CreateFilePath(filePath.(parser.Identifier), proc.Tx.Flags.Repository)
	if err != nil {
		return query.NewIOError(filePath, err.Error())
	}

	if _, err = os.Stat(path); err != nil {
		statements, _, err := parser.Parse("SELECT 1 FROM "+filename, "", false, proc.Tx.Flags.AnsiQuotes)
		if err != nil {
			return query.NewFileNotExistError(filePath)
		}

		q := statements[0].(parser.SelectQuery)
		filePath = q.SelectEntity.(parser.SelectEntity).FromClause.(parser.FromClause).Tables[0].(parser.Table).Object
		filePath.ClearBaseExpr()
	}

	statements := []parser.Statement{
		parser.ShowFields{
			Type:  parser.Identifier{Literal: "FIELDS"},
			Table: filePath,
		},
	}

	execCtx := context.WithValue(ctx, "CallFromSubcommand", true)
	_, err = proc.Execute(execCtx, statements)
	return err
}
