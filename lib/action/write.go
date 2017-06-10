package action

import (
	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/query"
)

func Write(input string) error {
	var out string
	var err error

	defer func() {
		if 0 < len(out) {
			cmd.ToStdout(out)
		}
	}()

	out, err = query.Execute(input)
	if err != nil {
		return err
	}

	flags := cmd.GetFlags()

	if 0 < len(flags.OutFile) {
		if err = cmd.CreateFile(flags.OutFile, out); err != nil {
			return err
		}
		out = ""
	}

	return nil
}
