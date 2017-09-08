package action

import (
	"os"
	"os/signal"

	"github.com/mithrandie/csvq/lib/query"
)

func SetSignalHandler() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	go func() {
		<-ch
		query.ReleaseResources()
		os.Exit(-1)
	}()
}
