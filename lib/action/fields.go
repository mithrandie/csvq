package action

import (
	"errors"

	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/query"
)

func ShowFields(filename string) error {
	SetSignalHandler()

	defer func() {
		query.ReleaseResources()
	}()

	query.UpdateWaitTimeout()

	statements := []parser.Statement{
		parser.ShowFields{
			Table: parser.Identifier{Literal: filename},
		},
	}

	proc := query.NewProcedure()
	_, err := proc.Execute(statements)
	if appErr, ok := err.(query.AppError); ok {
		err = errors.New(appErr.ErrorMessage())
	}

	return err
}
