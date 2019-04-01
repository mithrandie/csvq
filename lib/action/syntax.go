package action

import (
	"context"
	"errors"

	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/query"
)

func Syntax(proc *query.Processor, words []string) error {
	defer func() {
		if err := proc.ReleaseResourcesWithErrors(); err != nil {
			proc.LogError(err.Error())
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

	_, err := proc.Execute(context.Background(), statements)
	if appErr, ok := err.(query.Error); ok {
		err = errors.New(appErr.Message())
	}

	return err
}
