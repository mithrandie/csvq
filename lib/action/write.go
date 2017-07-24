package action

import (
	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/query"
)

func Write(input string, sourceFile string) error {
	var log string
	var selectLog string
	var err error

	defer func() {
		if 0 < len(log) {
			cmd.ToStdout(log)
		}
	}()

	log, selectLog, err = query.Execute(input, sourceFile)
	if err != nil {
		return err
	}

	flags := cmd.GetFlags()

	if 0 < len(flags.OutFile) && 0 < len(selectLog) {
		if err = cmd.CreateFile(flags.OutFile, selectLog); err != nil {
			return err
		}
	}

	return nil
}
