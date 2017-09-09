package action

import (
	"bytes"
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
)

func Run(input string, sourceFile string) error {
	SetSignalHandler()
	start := time.Now()

	defer func() {
		showStats(start, false)
	}()

	query.UpdateWaitTimeout()
	err := query.Execute(input, sourceFile)
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

	SetSignalHandler()
	start := time.Now()
	eof := false

	defer func() {
		query.ReleaseResources()
		showStats(start, eof)
	}()

	var err error

	query.UpdateWaitTimeout()
	proc := query.NewProcedure()
	var buf bytes.Buffer

	term, err := cmd.NewTerminal()
	if err != nil {
		return err
	}
	defer func() {
		term.Restore()
		cmd.Term = nil
	}()
	cmd.Term = term

	for {
		line, e := term.ReadLine()
		if e != nil {
			if e == io.EOF {
				eof = true
				break
			}
			return e
		}

		if buf.Len() < 1 && len(line) < 1 {
			continue
		}

		buf.WriteString(line)
		buf.WriteRune('\n')

		line = strings.TrimRightFunc(line, unicode.IsSpace)
		if 0 < len(line) && line[len(line)-1] != ';' {
			term.SetContinuousPrompt()
			continue
		}

		statements, e := parser.Parse(buf.String(), "")
		if e != nil {
			syntaxErr := e.(*parser.SyntaxError)
			e = query.NewSyntaxError(syntaxErr.Message, syntaxErr.Line, syntaxErr.Char, syntaxErr.SourceFile)
			if werr := term.Write(e.Error() + "\n"); werr != nil {
				return werr
			}
			buf.Reset()
			term.SetPrompt()
			continue
		}

		flow, e := proc.Execute(statements)
		if e != nil {
			if ex, ok := e.(*query.Exit); ok {
				err = ex
				break
			} else {
				if werr := term.Write(e.Error() + "\n"); werr != nil {
					return werr
				}
				buf.Reset()
				term.SetPrompt()
				continue
			}
		}

		if flow == query.EXIT {
			break
		}

		buf.Reset()
		term.SetPrompt()
	}

	createSelectLog()

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

func showStats(start time.Time, lower bool) {
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
		"      Time: %"+w+"[2]s seconds %[1]s"+
			"     Alloc: %"+w+"[3]s bytes %[1]s"+
			"TotalAlloc: %"+w+"[4]s bytes %[1]s"+
			"   HeapSys: %"+w+"[5]s bytes %[1]s"+
			"   Mallocs: %"+w+"[6]s objects %[1]s"+
			"     Frees: %"+w+"[7]s objects %[1]s",
		flags.LineBreak.Value(),
		exectime,
		alloc,
		talloc,
		sys,
		mallocs,
		frees,
	)
	if lower {
		cmd.ToStdout("\n")
	}
	cmd.ToStdout(stats)
}
