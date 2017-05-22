package main

import (
	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/query"
	"github.com/mithrandie/csvq/lib/stdout"
)

func Write(input string) error {
	results, err := query.Execute(input)
	if err != nil {
		return err
	}

	flags := cmd.GetFlags()

	for _, result := range results {
		switch result.Statement {
		case query.SELECT:
			if result.Count < 1 {
				stdout.Write("Empty\n")
			} else {
				var out string

				switch flags.Format {
				case cmd.STDOUT:
					out = stdout.Format(result.View)
				}

				switch flags.OutFile {
				case "":
					stdout.Write(out)
				}
			}
		}
	}
	return nil
}
