package action

import (
	"errors"
	"fmt"
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

	"github.com/mithrandie/go-file"
)

func Run(proc *query.Procedure, input string, sourceFile string, outfile string) error {
	start := time.Now()

	defer func() {
		if err := query.ReleaseResourcesWithErrors(); err != nil {
			cmd.WriteToStdErr(err.Error() + "\n")
		}
		showStats(start)
	}()

	statements, err := parser.Parse(input, sourceFile)
	if err != nil {
		syntaxErr := err.(*parser.SyntaxError)
		return query.NewSyntaxError(syntaxErr.Message, syntaxErr.Line, syntaxErr.Char, syntaxErr.SourceFile)
	}

	if 0 < len(outfile) {
		if abs, err := filepath.Abs(outfile); err == nil {
			outfile = abs
		}
		if csvqfile.Exists(outfile) {
			return errors.New(fmt.Sprintf("file %s already exists", outfile))
		}

		fp, err := file.Create(outfile)
		if err != nil {
			return errors.New(fmt.Sprintf("failed to create file: %s", err.Error()))
		}
		defer func() {
			if info, err := fp.Stat(); err == nil && info.Size() < 1 {
				os.Remove(outfile)
			}
			fp.Close()
		}()
		query.OutFile = fp
	}

	flow, err := proc.Execute(statements)

	if err == nil && flow == query.Terminate {
		if e := query.Commit(nil, proc.Filter); e != nil {
			cmd.WriteToStdErr(e.Error() + "\n")
		}
	} else {
		if e := query.Rollback(nil, proc.Filter); e != nil {
			cmd.WriteToStdErr(e.Error() + "\n")
		}
	}

	return err
}

func LaunchInteractiveShell(proc *query.Procedure) error {
	if cmd.IsReadableFromPipeOrRedirection() {
		return errors.New("input from pipe or redirection cannot be used in interactive shell")
	}

	defer func() {
		if err := query.ReleaseResourcesWithErrors(); err != nil {
			cmd.WriteToStdErr(err.Error() + "\n")
		}
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
			cmd.Terminal.SetPrompt()
			continue
		}
		cmd.Terminal.SaveHistory(saveQuery)

		statements, e := parser.Parse(strings.Join(lines, "\n"), "")
		if e != nil {
			syntaxErr := e.(*parser.SyntaxError)
			e = query.NewSyntaxError(syntaxErr.Message, syntaxErr.Line, syntaxErr.Char, syntaxErr.SourceFile)
			if werr := cmd.Terminal.WriteError(cmd.Error(e.Error()) + "\n"); werr != nil {
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
				if werr := cmd.Terminal.WriteError(cmd.Error(e.Error()) + "\n"); werr != nil {
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

	if e := query.Rollback(nil, proc.Filter); e != nil {
		cmd.WriteToStdErr(e.Error() + "\n")
	}

	return err
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
	width = width + 1

	w := cmd.NewObjectWriter()
	w.WriteColor(" TotalTime:", cmd.LableEffect)
	w.WriteSpaces(width - len(exectime))
	w.WriteWithoutLineBreak(exectime + " seconds")
	w.NewLine()
	w.WriteColor("     Alloc:", cmd.LableEffect)
	w.WriteSpaces(width - len(alloc))
	w.WriteWithoutLineBreak(alloc + " bytes")
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

	cmd.WriteToStdout("\n" + w.String())
}
