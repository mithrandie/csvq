package action

import (
	"os"
	"os/signal"

	"github.com/mithrandie/csvq/lib/query"
)

func SetSignalHandler(proc *query.Processor) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill)

	go func() {
		<-ch
		if err := proc.AutoRollback(); err != nil {
			proc.LogError(err.Error())
		}
		if err := proc.ReleaseResourcesWithErrors(); err != nil {
			proc.LogError(err.Error())
		}
		if proc.Tx.Session.Terminal != nil {
			if err := proc.Tx.Session.Terminal.Teardown(); err != nil {
				proc.LogError(err.Error())
			}
		}
		os.Exit(-1)
	}()
}
