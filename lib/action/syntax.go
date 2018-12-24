package action

import (
	"errors"

	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/query"
)

func Syntax(proc *query.Procedure, words []string) error {
	defer func() {
		if err := query.ReleaseResourcesWithErrors(); err != nil {
			query.LogError(err.Error())
		}
	}()

	keywords := make([]parser.QueryExpression, 0, len(words))
	for _, w := range words {
		keywords = append(keywords, parser.NewStringValue(w))
	}

	statements := []parser.Statement{
		parser.Syntax{
			Keywords: keywords,
		},
	}

	_, err := proc.Execute(statements)
	if appErr, ok := err.(query.AppError); ok {
		err = errors.New(appErr.ErrorMessage())
	}

	return err
}
