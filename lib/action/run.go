package action

import (
	"fmt"
	"runtime"
	"strconv"
	"time"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/query"
)

func Run(input string, sourceFile string) error {
	var log string
	var selectLog string
	var err error

	defer func() {
		if 0 < len(log) {
			cmd.ToStdout(log)
		}
	}()

	flags := cmd.GetFlags()

	var start time.Time
	if flags.Stats {
		start = time.Now()
	}

	log, selectLog, err = query.Execute(input, sourceFile)
	if err != nil {
		return err
	}

	if 0 < len(flags.OutFile) && 0 < len(selectLog) {
		if err = cmd.CreateFile(flags.OutFile, selectLog); err != nil {
			return err
		}
	}

	if flags.Stats {
		var mem runtime.MemStats
		runtime.ReadMemStats(&mem)

		exectime := cmd.HumarizeNumber(fmt.Sprintf("%f", time.Now().Sub(start).Seconds()))
		alloc := cmd.HumarizeNumber(fmt.Sprintf("%v", mem.Alloc))
		talloc := cmd.HumarizeNumber(fmt.Sprintf("%v", mem.TotalAlloc))
		sys := cmd.HumarizeNumber(fmt.Sprintf("%v", mem.HeapSys))
		released := cmd.HumarizeNumber(fmt.Sprintf("%v", mem.HeapReleased))

		width := len(exectime)
		mslen := len(sys)
		if width < mslen {
			width = mslen
		}
		w := strconv.Itoa(width)

		log += fmt.Sprintf(""+
			"        Time: %"+w+"[2]s seconds %[1]s"+
			"       Alloc: %"+w+"[3]s bytes %[1]s"+
			"  TotalAlloc: %"+w+"[4]s bytes %[1]s"+
			"     HeapSys: %"+w+"[5]s bytes %[1]s"+
			"HeapReleased: %"+w+"[6]s bytes %[1]s",
			flags.LineBreak.Value(),
			exectime,
			alloc,
			talloc,
			sys,
			released)
	}

	return nil
}
