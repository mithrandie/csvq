package action

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
	"unicode"

	"github.com/mithrandie/csvq/lib/cmd"
	csvqfile "github.com/mithrandie/csvq/lib/file"
	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/query"

	"github.com/mithrandie/go-file/v2"
)

func Run(ctx context.Context, proc *query.Processor, input string, sourceFile string, outfile string) error {
	start := time.Now()

	defer func() {
		showStats(ctx, proc, start)
	}()

	statements, _, err := parser.Parse(input, sourceFile, false, proc.Tx.Flags.AnsiQuotes)
	if err != nil {
		return query.NewSyntaxError(err.(*parser.SyntaxError))
	}

	if 0 < len(outfile) {
		if abs, err := filepath.Abs(outfile); err == nil {
			outfile = abs
		}
		if csvqfile.Exists(outfile) {
			return query.NewFileAlreadyExistError(parser.Identifier{Literal: outfile})
		}

		fp, err := file.Create(outfile)
		if err != nil {
			return query.NewIOError(nil, err.Error())
		}
		defer func() {
			if info, err := fp.Stat(); err == nil && info.Size() < 1 {
				if err = os.Remove(outfile); err != nil {
					proc.LogError(err.Error())
				}
			}
			if err = fp.Close(); err != nil {
				proc.LogError(err.Error())
			}
		}()
		proc.Tx.Session.SetOutFile(fp)
	}

	proc.Tx.AutoCommit = true
	_, err = proc.Execute(ctx, statements)
	return err
}

func LaunchInteractiveShell(ctx context.Context, proc *query.Processor) error {
	if proc.Tx.Session.CanReadStdin {
		return query.NewIncorrectCommandUsageError("input from pipe or redirection cannot be used in interactive shell")
	}

	if err := proc.Tx.Session.SetStdin(query.GetStdinForREPL()); err != nil {
		return query.NewIOError(nil, err.Error())
	}

	term, err := query.NewTerminal(ctx, proc.ReferenceScope)
	if err != nil {
		return query.ConvertLoadConfigurationError(err)
	}
	proc.Tx.Session.SetTerminal(term)
	defer func() {
		if e := proc.Tx.Session.Terminal().Teardown(); e != nil {
			proc.LogError(e.Error())
		}
		proc.Tx.Session.SetTerminal(nil)
	}()

	StartUpMessage := "" +
		"csvq interactive shell\n" +
		"Press Ctrl+D or execute \"EXIT;\" to terminate this shell.\n\n"
	proc.Log(StartUpMessage, false)

	lines := make([]string, 0)

	for {
		if ctx.Err() != nil {
			err = query.ConvertContextError(ctx.Err())
			break
		}

		proc.Tx.Session.Terminal().UpdateCompleter()
		line, e := proc.Tx.Session.Terminal().ReadLine()
		if e != nil {
			if e == io.EOF {
				break
			}
			return query.NewIOError(nil, e.Error())
		}

		line = strings.TrimRightFunc(line, unicode.IsSpace)

		if len(lines) < 1 && len(line) < 1 {
			continue
		}

		if 0 < len(line) && line[len(line)-1] == '\\' {
			lines = append(lines, line[:len(line)-1])
			proc.Tx.Session.Terminal().SetContinuousPrompt(ctx)
			continue
		}

		lines = append(lines, line)

		saveLines := make([]string, 0, len(lines))
		for _, l := range lines {
			s := strings.TrimSpace(l)
			if len(s) < 1 {
				continue
			}
			saveLines = append(saveLines, s)
		}

		saveQuery := strings.Join(saveLines, " ")
		if len(saveQuery) < 1 || saveQuery == ";" {
			lines = lines[:0]
			proc.Tx.Session.Terminal().SetPrompt(ctx)
			continue
		}
		if e := proc.Tx.Session.Terminal().SaveHistory(saveQuery); e != nil {
			proc.LogError(e.Error())
		}

		statements, _, e := parser.Parse(strings.Join(lines, "\n"), "", false, proc.Tx.Flags.AnsiQuotes)
		if e != nil {
			if e = query.NewSyntaxError(e.(*parser.SyntaxError)); e != nil {
				proc.LogError(e.Error())
			}
			lines = lines[:0]
			proc.Tx.Session.Terminal().SetPrompt(ctx)
			continue
		}

		flow, e := proc.Execute(ctx, statements)
		if e != nil {
			if ex, ok := e.(*query.ForcedExit); ok {
				err = ex
				break
			} else {
				proc.LogError(e.Error())
				lines = lines[:0]
				proc.Tx.Session.Terminal().SetPrompt(ctx)
				continue
			}
		}

		if flow == query.Exit {
			break
		}

		lines = lines[:0]
		proc.Tx.Session.Terminal().SetPrompt(ctx)
	}

	return err
}

func showStats(ctx context.Context, proc *query.Processor, start time.Time) {
	if ctx.Err() != nil {
		return
	}

	if !proc.Tx.Flags.Stats {
		return
	}
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	exectime := cmd.FormatNumber(time.Since(start).Seconds(), 6, ".", ",", "")
	talloc := cmd.FormatNumber(float64(mem.TotalAlloc), 0, ".", ",", "")
	sys := cmd.FormatNumber(float64(mem.HeapSys), 0, ".", ",", "")
	mallocs := cmd.FormatNumber(float64(mem.Mallocs), 0, ".", ",", "")
	frees := cmd.FormatNumber(float64(mem.Frees), 0, ".", ",", "")

	width := len(exectime)
	for _, v := range []string{talloc, sys, mallocs, frees} {
		if width < len(v) {
			width = len(v)
		}
	}
	width = width + 1

	w := query.NewObjectWriter(proc.Tx)
	w.WriteColor(" TotalTime:", cmd.LableEffect)
	w.WriteSpaces(width - len(exectime))
	w.WriteWithoutLineBreak(exectime + " seconds")
	w.NewLine()
	w.WriteColor("TotalAlloc:", cmd.LableEffect)
	w.WriteSpaces(width - len(talloc))
	w.WriteWithoutLineBreak(talloc + " bytes")
	w.NewLine()
	w.WriteColor("   HeapSys:", cmd.LableEffect)
	w.WriteSpaces(width - len(sys))
	w.WriteWithoutLineBreak(sys + " bytes")
	w.NewLine()
	w.WriteColor("   Mallocs:", cmd.LableEffect)
	w.WriteSpaces(width - len(mallocs))
	w.WriteWithoutLineBreak(mallocs + " objects")
	w.NewLine()
	w.WriteColor("     Frees:", cmd.LableEffect)
	w.WriteSpaces(width - len(frees))
	w.WriteWithoutLineBreak(frees + " objects")
	w.NewLine()
	w.NewLine()

	w.Title1 = "Resource Statistics"

	proc.Log("\n"+w.String(), false)
}
