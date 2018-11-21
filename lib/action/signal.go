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
		if err := query.Rollback(nil, nil); err != nil {
			cmd.WriteToStderrWithLineBreak(err.Error())
		}
		if err := query.ReleaseResourcesWithErrors(); err != nil {
			cmd.WriteToStderrWithLineBreak(err.Error())
		}
		if cmd.Terminal != nil {
			cmd.Terminal.Teardown()
		}
		os.Exit(-1)
	}()
}
