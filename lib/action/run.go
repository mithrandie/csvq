package action

import (
	"errors"
	"fmt"
	"io"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/query"

	"github.com/mithrandie/go-text/color"
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

	exectime := cmd.FormatNumber(time.Since(start).Seconds(), 6, ".", ",", "")
	alloc := cmd.FormatNumber(float64(mem.Alloc), 0, ".", ",", "")
	talloc := cmd.FormatNumber(float64(mem.TotalAlloc), 0, ".", ",", "")
	sys := cmd.FormatNumber(float64(mem.HeapSys), 0, ".", ",", "")
	mallocs := cmd.FormatNumber(float64(mem.Mallocs), 0, ".", ",", "")
	frees := cmd.FormatNumber(float64(mem.Frees), 0, ".", ",", "")

	width := len(exectime)
	for _, v := range []string{alloc, talloc, sys, mallocs, frees} {
		if width < len(v) {
			width = len(v)
		}
	}
	w := strconv.Itoa(width)

	palette := cmd.GetPalette()

	stats := fmt.Sprintf(
		""+
			palette.Render(cmd.LableEffect, "   TotalTime: ")+"%"+w+"[2]s seconds %[1]s"+
			palette.Render(cmd.LableEffect, "       Alloc: ")+"%"+w+"[3]s bytes %[1]s"+
			palette.Render(cmd.LableEffect, "  TotalAlloc: ")+"%"+w+"[4]s bytes %[1]s"+
			palette.Render(cmd.LableEffect, "     HeapSys: ")+"%"+w+"[5]s bytes %[1]s"+
			palette.Render(cmd.LableEffect, "     Mallocs: ")+"%"+w+"[6]s objects %[1]s"+
			palette.Render(cmd.LableEffect, "       Frees: ")+"%"+w+"[7]s objects %[1]s",
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
