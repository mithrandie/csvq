package action

import (
	"errors"

	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/query"
)

func ShowFields(proc *query.Procedure, filename string) error {
	defer func() {
		if err := query.ReleaseResourcesWithErrors(); err != nil {
			query.WriteToStderrWithLineBreak(err.Error())
		}
	}()

	statements := []parser.Statement{
		parser.ShowFields{
			Type:  parser.Identifier{Literal: "FIELDS"},
			Table: parser.Identifier{Literal: filename},
		},
	}

	_, err := proc.Execute(statements)
	if appErr, ok := err.(query.AppError); ok {
		err = errors.New(appErr.ErrorMessage())
	}

	return err
}
