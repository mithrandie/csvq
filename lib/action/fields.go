package action

import (
	"context"

	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/query"
)

func ShowFields(ctx context.Context, proc *query.Processor, filename string) error {
	statements := []parser.Statement{
		parser.ShowFields{
			Type:  parser.Identifier{Literal: "FIELDS"},
			Table: parser.Identifier{Literal: filename},
		},
	}

	_, err := proc.Execute(ctx, statements)
	return err
}
