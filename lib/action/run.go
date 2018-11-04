package action

import (
	"errors"
	"fmt"
	"github.com/mithrandie/csvq/lib/color"
	"io"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/query"
)

func Run(input string, sourceFile string) error {
	SetSignalHandler()
	start := time.Now()

	defer func() {
		query.ReleaseResources()
		showStats(start)
	}()

	statements, err := parser.Parse(input, sourceFile)
	if err != nil {
		syntaxErr := err.(*parser.SyntaxError)
		return query.NewSyntaxError(syntaxErr.Message, syntaxErr.Line, syntaxErr.Char, syntaxErr.SourceFile)
	}

	proc := query.NewProcedure()
	flow, err := proc.Execute(statements)

	if err == nil && flow == query.Terminate {
		err = query.Commit(nil, proc.Filter)
	}

	if err != nil {
		return err
	}

	createSelectLog()

	return nil
}

func LaunchInteractiveShell() error {
	if cmd.IsReadableFromPipeOrRedirection() {
		return errors.New("input from pipe or redirection cannot be used in interactive shell")
	}

	cmd.SetOut("") // Ignore --out option

	SetSignalHandler()

	defer func() {
		query.ReleaseResources()
	}()

	var err error

	term, err := cmd.NewTerminal()
	if err != nil {
		return err
	}
	cmd.Terminal = term
	defer func() {
		cmd.Terminal.Teardown()
		cmd.Terminal = nil
	}()

	StartUpMessage := "" +
		"csvq interactive shell\n" +
		"Press Ctrl+D or execute \"EXIT;\" to terminate this shell.\n\n"
	if werr := cmd.Terminal.Write(StartUpMessage); werr != nil {
		return werr
	}

	proc := query.NewProcedure()
	lines := make([]string, 0)

	for {
		line, e := cmd.Terminal.ReadLine()
		if e != nil {
			if e == io.EOF {
				break
			}
			return e
		}

		line = strings.TrimRightFunc(line, unicode.IsSpace)

		if len(lines) < 1 && len(line) < 1 {
			continue
		}

		if 0 < len(line) && line[len(line)-1] == '\\' {
			lines = append(lines, line[:len(line)-1])
			cmd.Terminal.SetContinuousPrompt()
			continue
		}

		lines = append(lines, line)

		if len(line) < 1 || line[len(line)-1] != ';' {
			cmd.Terminal.SetContinuousPrompt()
			continue
		}

		cmd.Terminal.SaveHistory(strings.Join(lines, " "))

		statements, e := parser.Parse(strings.Join(lines, "\n"), "")
		if e != nil {
			syntaxErr := e.(*parser.SyntaxError)
			e = query.NewSyntaxError(syntaxErr.Message, syntaxErr.Line, syntaxErr.Char, syntaxErr.SourceFile)
			if werr := cmd.Terminal.Write(color.Error(e.Error()) + "\n"); werr != nil {
				return werr
			}
			lines = lines[:0]
			cmd.Terminal.SetPrompt()
			continue
		}

		flow, e := proc.Execute(statements)
		if e != nil {
			if ex, ok := e.(*query.ForcedExit); ok {
				err = ex
				break
			} else {
				if werr := cmd.Terminal.Write(color.Error(e.Error()) + "\n"); werr != nil {
					return werr
				}
				lines = lines[:0]
				cmd.Terminal.SetPrompt()
				continue
			}
		}

		if flow == query.Exit {
			break
		}

		lines = lines[:0]
		cmd.Terminal.SetPrompt()
	}

	return err
}

func createSelectLog() error {
	flags := cmd.GetFlags()
	selectLog := query.ReadSelectLog()
	if 0 < len(flags.OutFile) && 0 < len(selectLog) {
		if err := cmd.CreateFile(flags.OutFile, selectLog); err != nil {
			return err
		}
	}
	return nil
}

func showStats(start time.Time) {
	flags := cmd.GetFlags()
	if !flags.Stats {
		return
	}
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	exectime := cmd.HumarizeNumber(fmt.Sprintf("%f", time.Since(start).Seconds()))
	alloc := cmd.HumarizeNumber(fmt.Sprintf("%v", mem.Alloc))
	talloc := cmd.HumarizeNumber(fmt.Sprintf("%v", mem.TotalAlloc))
	sys := cmd.HumarizeNumber(fmt.Sprintf("%v", mem.HeapSys))
	mallocs := cmd.HumarizeNumber(fmt.Sprintf("%v", mem.Mallocs))
	frees := cmd.HumarizeNumber(fmt.Sprintf("%v", mem.Frees))

	width := len(exectime)
	for _, v := range []string{alloc, talloc, sys, mallocs, frees} {
		if width < len(v) {
			width = len(v)
		}
	}
	w := strconv.Itoa(width)

	stats := fmt.Sprintf(
		""+
			color.BlueB("   TotalTime: ")+"%"+w+"[2]s seconds %[1]s"+
			color.BlueB("       Alloc: ")+"%"+w+"[3]s bytes %[1]s"+
			color.BlueB("  TotalAlloc: ")+"%"+w+"[4]s bytes %[1]s"+
			color.BlueB("     HeapSys: ")+"%"+w+"[5]s bytes %[1]s"+
			color.BlueB("     Mallocs: ")+"%"+w+"[6]s objects %[1]s"+
			color.BlueB("       Frees: ")+"%"+w+"[7]s objects %[1]s",
		"\n",
		exectime,
		alloc,
		talloc,
		sys,
		mallocs,
		frees,
	)
	cmd.ToStdout(stats)
}
