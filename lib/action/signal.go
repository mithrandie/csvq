package action

import (
	"os"
	"os/signal"

	"github.com/mithrandie/csvq/lib/cmd"

	"github.com/mithrandie/csvq/lib/query"
)

func SetSignalHandler() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill)

	go func() {
		<-ch
		if e := query.Rollback(nil, nil); e != nil {
			cmd.WriteToStdErr(e.Error() + "\n")
		}
		if errs := query.ReleaseResourcesWithErrors(); errs != nil {
			for _, err := range errs {
				cmd.WriteToStdErr(err.Error() + "\n")
			}
		}
		os.Exit(-1)
	}()
}
