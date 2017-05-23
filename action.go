package main

import (
	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/output"
	"github.com/mithrandie/csvq/lib/query"
)

func Write(input string) error {
	results, err := query.Execute(input)
	if err != nil {
		return err
	}

	flags := cmd.GetFlags()

	for _, result := range results {
		out := output.Encode(flags.Format, result)
		output.Write(flags.OutFile, out)
	}
	return nil
}
