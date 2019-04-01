package action

import (
	"context"

	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/query"
)

func Syntax(proc *query.Processor, words []string) error {
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
	return err
}
