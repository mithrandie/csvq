package action

import (
	"context"
	"errors"

	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/query"
)

func ShowFields(proc *query.Processor, filename string) error {
	defer func() {
		if err := proc.ReleaseResourcesWithErrors(); err != nil {
			proc.LogError(err.Error())
		}
	}()

	statements := []parser.Statement{
		parser.ShowFields{
			Type:  parser.Identifier{Literal: "FIELDS"},
			Table: parser.Identifier{Literal: filename},
		},
	}

	_, err := proc.Execute(context.Background(), statements)
	if appErr, ok := err.(query.Error); ok {
		err = errors.New(appErr.ErrorMessage())
	}

	return err
}
