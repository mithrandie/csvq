package action

import (
	"os"
	"os/signal"

	"github.com/mithrandie/csvq/lib/query"
)

func SetSignalHandler() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill)

	go func() {
		<-ch
		if err := query.Rollback(nil, nil); err != nil {
			query.WriteToStderrWithLineBreak(err.Error())
		}
		if err := query.ReleaseResourcesWithErrors(); err != nil {
			query.WriteToStderrWithLineBreak(err.Error())
		}
		if query.Terminal != nil {
			query.Terminal.Teardown()
		}
		os.Exit(-1)
	}()
}
