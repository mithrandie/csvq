package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/mithrandie/csvq/lib/action"
	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/file"
	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/query"

	"github.com/mithrandie/go-text/color"
	"github.com/urfave/cli"
)

func main() {
	var proc *query.Processor
	action.CurrentVersion, _ = action.ParseVersion(query.Version)
	if action.CurrentVersion != nil {
		query.Version = action.CurrentVersion.String()
	}

	cli.AppHelpTemplate = appHHelpTemplate
	cli.CommandHelpTemplate = commandHelpTemplate

	app := cli.NewApp()

	app.Name = "csvq"
	app.Usage = "SQL-like query language for csv"
	app.ArgsUsage = "[\"query\"|argument]"
	app.Version = query.Version

	app.OnUsageError = func(c *cli.Context, err error, isSubcommand bool) error {
		if isSubcommand {
			return err
		}

		return NewExitError(fmt.Sprintf("Incorrect Usage: %s", err.Error()), 1)
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "repository, r",
			Usage: "directory `PATH` where files are located",
		},
		cli.StringFlag{
			Name:  "timezone, z",
			Value: "Local",
			Usage: "default timezone. \"Local\", \"UTC\" or a timezone name",
		},
		cli.StringFlag{
			Name:  "datetime-format, t",
			Usage: "datetime format to parse strings",
		},
		cli.Float64Flag{
			Name:  "wait-timeout, w",
			Value: 10,
			Usage: "limit of the waiting time in seconds to wait for locked files to be released",
		},
		cli.StringFlag{
			Name:  "source, s",
			Usage: "load query or statements from `FILE`",
		},
		cli.StringFlag{
			Name:  "import-format, i",
			Value: "CSV",
			Usage: "default format to load files. one of: CSV|TSV|FIXED|JSON|LTSV",
		},
		cli.StringFlag{
			Name:  "delimiter, d",
			Value: ",",
			Usage: "field delimiter for CSV",
		},
		cli.StringFlag{
			Name:  "delimiter-positions, m",
			Usage: "delimiter positions for FIXED",
		},
		cli.StringFlag{
			Name:  "json-query, j",
			Usage: "`QUERY` for JSON",
		},
		cli.StringFlag{
			Name:  "encoding, e",
			Value: "UTF8",
			Usage: "file encoding. one of: UTF8|UTF8M|SJIS",
		},
		cli.BoolFlag{
			Name:  "no-header, n",
			Usage: "import the first line as a record",
		},
		cli.BoolFlag{
			Name:  "without-null, a",
			Usage: "parse empty fields as empty strings",
		},
		cli.StringFlag{
			Name:  "out, o",
			Usage: "export result sets of select queries to `FILE`",
		},
		cli.StringFlag{
			Name:  "format, f",
			Value: "TEXT",
			Usage: "format of query results. one of: CSV|TSV|FIXED|JSON|LTSV|GFM|ORG|TEXT",
		},
		cli.StringFlag{
			Name:  "write-encoding, E",
			Value: "UTF8",
			Usage: "character encoding of query results. one of: UTF8|UTF8M|SJIS",
		},
		cli.StringFlag{
			Name:  "write-delimiter, D",
			Value: ",",
			Usage: "field delimiter for CSV in query results",
		},
		cli.StringFlag{
			Name:  "write-delimiter-positions, M",
			Usage: "delimiter positions for FIXED in query results",
		},
		cli.BoolFlag{
			Name:  "without-header, N",
			Usage: "export result sets of select queries without the header line",
		},
		cli.StringFlag{
			Name:  "line-break, l",
			Value: "LF",
			Usage: "line break in query results. one of: CRLF|LF|CR",
		},
		cli.BoolFlag{
			Name:  "enclose-all, Q",
			Usage: "enclose all string values in CSV and TSV",
		},
		cli.StringFlag{
			Name:  "json-escape, J",
			Value: "BACKSLASH",
			Usage: "JSON escape type. one of: BACKSLASH|HEX|HEXALL",
		},
		cli.BoolFlag{
			Name:  "pretty-print, P",
			Usage: "make JSON output easier to read in query results",
		},
		cli.BoolFlag{
			Name:  "east-asian-encoding, W",
			Usage: "count ambiguous characters as fullwidth",
		},
		cli.BoolFlag{
			Name:  "count-diacritical-sign, S",
			Usage: "count diacritical signs as halfwidth",
		},
		cli.BoolFlag{
			Name:  "count-format-code, A",
			Usage: "count format characters and zero-width spaces as halfwidth",
		},
		cli.BoolFlag{
			Name:  "color, c",
			Usage: "use ANSI color escape sequences",
		},
		cli.BoolFlag{
			Name:  "quiet, q",
			Usage: "suppress operation log output",
		},
		cli.IntFlag{
			Name:  "cpu, p",
			Value: cmd.GetDefaultNumberOfCPU(),
			Usage: "hint for the number of cpu cores to be used",
		},
		cli.BoolFlag{
			Name:  "stats, x",
			Usage: "show execution time and memory statistics",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:      "fields",
			Usage:     "Show fields in a file",
			ArgsUsage: "CSV_FILE_PATH",
			Action: func(c *cli.Context) error {
				if c.NArg() != 1 {
					return NewExitError("table is not specified", 1)
				}

				table := c.Args().First()

				err := action.ShowFields(proc, table)
				if err != nil {
					return NewExitError(err.Error(), 1)
				}

				return nil
			},
			OnUsageError: func(c *cli.Context, err error, isSubcommand bool) error {
				return NewExitError(fmt.Sprintf("Incorrect Usage: %s", err.Error()), 1)
			},
		},
		{
			Name:      "calc",
			Usage:     "Calculate a value from stdin",
			ArgsUsage: "\"expression\"",
			Action: func(c *cli.Context) error {
				if c.NArg() != 1 {
					return NewExitError("expression is empty", 1)
				}

				expr := c.Args().First()
				err := action.Calc(proc, expr)
				if err != nil {
					return NewExitError(err.Error(), 1)
				}

				return nil
			},
			OnUsageError: func(c *cli.Context, err error, isSubcommand bool) error {
				return NewExitError(fmt.Sprintf("Incorrect Usage: %s", err.Error()), 1)
			},
		},
		{
			Name:      "syntax",
			Usage:     "Print syntax",
			ArgsUsage: "[search_word ...]",
			Action: func(c *cli.Context) error {
				words := append([]string{c.Args().First()}, c.Args().Tail()...)
				err := action.Syntax(proc, words)
				if err != nil {
					return NewExitError(err.Error(), 1)
				}

				return nil
			},
			OnUsageError: func(c *cli.Context, err error, isSubcommand bool) error {
				return NewExitError(fmt.Sprintf("Incorrect Usage: %s", err.Error()), 1)
			},
		},
		{
			Name:  "check-update",
			Usage: "Check for updates",
			Action: func(c *cli.Context) error {
				err := action.CheckUpdate(proc)
				if err != nil {
					return NewExitError(err.Error(), 1)
				}

				return nil
			},
		},
	}

	app.Before = func(c *cli.Context) error {
		color.UseEffect = false

		defaultWaitTimeout := file.DefaultWaitTimeout
		if c.IsSet("wait-timeout") {
			if d, err := time.ParseDuration(strconv.FormatFloat(c.GlobalFloat64("wait-timeout"), 'f', -1, 64) + "s"); err == nil {
				defaultWaitTimeout = d
			}
		}

		session := query.NewSession()
		tx, err := query.NewTransaction(context.Background(), defaultWaitTimeout, file.DefaultRetryDelay, session)
		if err != nil {
			return NewExitError(err.Error(), 1)
		}

		proc = query.NewProcessor(tx)

		action.SetSignalHandler(proc)

		// Run pre-load commands
		if err := runPreloadCommands(proc); err != nil {
			return NewExitError(err.Error(), 1)
		}

		// Overwrite Flags with Command Options
		if err := overwriteFlags(c, tx); err != nil {
			return NewExitError(err.Error(), 1)
		}
		return nil
	}

	app.Action = func(c *cli.Context) error {
		queryString, path, err := readQuery(c, proc.Tx)
		if err != nil {
			return NewExitError(err.Error(), 1)
		}

		if len(queryString) < 1 {
			err = action.LaunchInteractiveShell(proc)
		} else {
			err = action.Run(proc, queryString, path, c.GlobalString("out"))
		}

		if err != nil {
			code := 1
			if apperr, ok := err.(query.Error); ok {
				code = apperr.GetCode()
			} else if ex, ok := err.(*query.ForcedExit); ok {
				code = ex.GetCode()
			}
			return NewExitError(err.Error(), code)
		}

		return nil
	}

	if err := app.Run(os.Args); err != nil {
		println(err.Error())
	}
}

func readQuery(c *cli.Context, tx *query.Transaction) (queryString string, path string, err error) {
	if c.IsSet("source") && 0 < len(c.GlobalString("source")) {
		path = c.GlobalString("source")
		if abs, err := filepath.Abs(path); err == nil {
			path = abs
		}
		if !file.Exists(path) {
			err = errors.New(fmt.Sprintf("file %q does not exist", path))
			return
		}

		h, e := file.NewHandlerForRead(context.Background(), tx.FileContainer, path, tx.WaitTimeout, tx.RetryDelay)
		if e != nil {
			err = errors.New(fmt.Sprintf("failed to read file: %s", e.Error()))
			return
		}
		defer func() {
			if e := tx.FileContainer.Close(h); e != nil {
				if err == nil {
					err = e
				} else {
					err = errors.New(err.Error() + "\n" + e.Error())
				}
			}
		}()

		buf, e := ioutil.ReadAll(h.FileForRead())
		if e != nil {
			err = errors.New(fmt.Sprintf("failed to read file: %s", e.Error()))
			return
		}

		queryString = string(buf)
	} else {
		if 1 < c.NArg() {
			err = errors.New("multiple queries or statements were passed")
			return
		}
		queryString = c.Args().First()
	}
	return
}

func overwriteFlags(c *cli.Context, tx *query.Transaction) error {
	flags := tx.Flags

	if c.IsSet("color") {
		flags.SetColor(c.GlobalBool("color"))
	}

	if c.IsSet("repository") {
		if err := flags.SetRepository(c.GlobalString("repository")); err != nil {
			return err
		}
	}
	if c.IsSet("timezone") {
		if err := flags.SetLocation(c.String("timezone")); err != nil {
			return err
		}
	}
	if c.IsSet("datetime-format") {
		flags.SetDatetimeFormat(c.GlobalString("datetime-format"))
	}
	if c.IsSet("wait-timeout") {
		tx.UpdateWaitTimeout(c.GlobalFloat64("wait-timeout"), file.DefaultRetryDelay)
	}

	if c.IsSet("import-format") {
		if err := flags.SetImportFormat(c.GlobalString("import-format")); err != nil {
			return err
		}
	}
	if c.IsSet("delimiter") {
		if err := flags.SetDelimiter(c.GlobalString("delimiter")); err != nil {
			return err
		}
	}
	if c.IsSet("delimiter-positions") {
		if err := flags.SetDelimiterPositions(c.GlobalString("delimiter-positions")); err != nil {
			return err
		}
	}
	if c.IsSet("json-query") {
		flags.SetJsonQuery(c.GlobalString("json-query"))
	}
	if c.IsSet("encoding") {
		if err := flags.SetEncoding(c.GlobalString("encoding")); err != nil {
			return err
		}
	}
	if c.IsSet("no-header") {
		flags.SetNoHeader(c.GlobalBool("no-header"))
	}
	if c.IsSet("without-null") {
		flags.SetWithoutNull(c.GlobalBool("without-null"))
	}

	if c.IsSet("format") {
		if err := flags.SetFormat(c.GlobalString("format"), c.GlobalString("out")); err != nil {
			return err
		}
	}
	if c.IsSet("write-encoding") {
		if err := flags.SetWriteEncoding(c.GlobalString("write-encoding")); err != nil {
			return err
		}
	}
	if c.IsSet("write-delimiter") {
		if err := flags.SetWriteDelimiter(c.GlobalString("write-delimiter")); err != nil {
			return err
		}
	}
	if c.IsSet("write-delimiter-positions") {
		if err := flags.SetWriteDelimiterPositions(c.GlobalString("write-delimiter-positions")); err != nil {
			return err
		}
	}
	if c.IsSet("without-header") {
		flags.SetWithoutHeader(c.GlobalBool("without-header"))
	}
	if c.IsSet("line-break") {
		if err := flags.SetLineBreak(c.String("line-break")); err != nil {
			return err
		}
	}
	if c.IsSet("enclose-all") {
		flags.SetEncloseAll(c.GlobalBool("enclose-all"))
	}
	if c.IsSet("json-escape") {
		if err := flags.SetJsonEscape(c.GlobalString("json-escape")); err != nil {
			return err
		}
	}
	if c.IsSet("pretty-print") {
		flags.SetPrettyPrint(c.GlobalBool("pretty-print"))
	}

	if c.IsSet("east-asian-encoding") {
		flags.SetEastAsianEncoding(c.GlobalBool("east-asian-encoding"))
	}
	if c.IsSet("count-diacritical-sign") {
		flags.SetCountDiacriticalSign(c.GlobalBool("count-diacritical-sign"))
	}
	if c.IsSet("count-format-code") {
		flags.SetCountFormatCode(c.GlobalBool("count-format-code"))
	}

	if c.IsSet("quiet") {
		flags.SetQuiet(c.GlobalBool("quiet"))
	}
	if c.IsSet("cpu") {
		flags.SetCPU(c.GlobalInt("cpu"))
	}
	if c.IsSet("stats") {
		flags.SetStats(c.GlobalBool("stats"))
	}

	return nil
}

func runPreloadCommands(proc *query.Processor) (err error) {
	ctx := context.Background()
	files := cmd.GetSpecialFilePath(cmd.PreloadCommandFileName)
	for _, fpath := range files {
		if !file.Exists(fpath) {
			continue
		}

		statements, err := query.LoadStatementsFromFile(ctx, proc.Tx, parser.Source{}, fpath)
		if err != nil {
			if e, ok := err.(*query.ReadFileError); ok {
				err = errors.New(e.ErrorMessage())
			}
			return err
		}

		if _, err := proc.Execute(ctx, statements); err != nil {
			return err
		}
	}
	return nil
}

func NewExitError(message string, code int) *cli.ExitError {
	return cli.NewExitError(cmd.Error(message), code)
}
