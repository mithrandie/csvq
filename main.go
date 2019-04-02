package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
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

		return Exit(query.NewIncorrectCommandUsageError(err.Error()))
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
			Action: commandAction(func(ctx context.Context, c *cli.Context, proc *query.Processor) error {
				if c.NArg() != 1 {
					return query.NewIncorrectCommandUsageError("table is not specified")
				}

				table := c.Args().First()

				return action.ShowFields(ctx, proc, table)
			}),
			OnUsageError: func(c *cli.Context, err error, isSubcommand bool) error {
				return Exit(query.NewIncorrectCommandUsageError(err.Error()))
			},
		},
		{
			Name:      "calc",
			Usage:     "Calculate a value from stdin",
			ArgsUsage: "\"expression\"",
			Action: commandAction(func(ctx context.Context, c *cli.Context, proc *query.Processor) error {
				if c.NArg() != 1 {
					return query.NewIncorrectCommandUsageError("expression is empty")
				}

				expr := c.Args().First()
				return action.Calc(ctx, proc, expr)
			}),
			OnUsageError: func(c *cli.Context, err error, isSubcommand bool) error {
				return Exit(query.NewIncorrectCommandUsageError(err.Error()))
			},
		},
		{
			Name:      "syntax",
			Usage:     "Print syntax",
			ArgsUsage: "[search_word ...]",
			Action: commandAction(func(ctx context.Context, c *cli.Context, proc *query.Processor) error {
				words := append([]string{c.Args().First()}, c.Args().Tail()...)
				return action.Syntax(ctx, proc, words)
			}),
			OnUsageError: func(c *cli.Context, err error, isSubcommand bool) error {
				return Exit(query.NewIncorrectCommandUsageError(err.Error()))
			},
		},
		{
			Name:  "check-update",
			Usage: "Check for updates",
			Action: func(c *cli.Context) error {
				return Exit(action.CheckUpdate())
			},
		},
	}

	app.Action = commandAction(func(ctx context.Context, c *cli.Context, proc *query.Processor) error {
		queryString, path, err := readQuery(c, proc.Tx)
		if err != nil {
			return err
		}

		if len(queryString) < 1 {
			err = action.LaunchInteractiveShell(ctx, proc)
		} else {
			err = action.Run(ctx, proc, queryString, path, c.GlobalString("out"))
		}

		return err
	})

	if err := app.Run(os.Args); err != nil {
		log.Fatalln(err.Error())
	}
}

func commandAction(fn func(ctx context.Context, c *cli.Context, proc *query.Processor) error) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		proc, err := generateProcessor(ctx, c)
		if err != nil {
			return Exit(err)
		}
		defer func() {
			if e := proc.AutoRollback(); e != nil {
				proc.LogError(e.Error())
			}
			if err := proc.ReleaseResourcesWithErrors(); err != nil {
				proc.LogError(err.Error())
			}
		}()

		// Handle signals
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, action.Signals...)

		go func() {
			sig := <-ch
			proc.LogWarn(fmt.Sprintf("signal received: %s", sig.String()), false)
			cancel()
		}()

		// Run pre-load commands
		if err := runPreloadCommands(ctx, proc); err != nil {
			return Exit(err)
		}

		// Overwrite Flags with Command Options
		if err := overwriteFlags(c, proc.Tx); err != nil {
			return Exit(query.NewIncorrectCommandUsageError(err.Error()))
		}

		err = fn(ctx, c, proc)
		if _, ok := err.(*query.ContextIsDone); ok {
			err = nil
		}
		return Exit(err)
	}
}

func generateProcessor(ctx context.Context, c *cli.Context) (*query.Processor, error) {
	color.UseEffect = false

	defaultWaitTimeout := file.DefaultWaitTimeout
	if c.IsSet("wait-timeout") {
		if d, err := time.ParseDuration(strconv.FormatFloat(c.GlobalFloat64("wait-timeout"), 'f', -1, 64) + "s"); err == nil {
			defaultWaitTimeout = d
		}
	}

	session := query.NewSession()
	tx, err := query.NewTransaction(ctx, defaultWaitTimeout, file.DefaultRetryDelay, session)
	if err != nil {
		return nil, err
	}

	return query.NewProcessor(tx), nil
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

func runPreloadCommands(ctx context.Context, proc *query.Processor) (err error) {
	files := cmd.GetSpecialFilePath(cmd.PreloadCommandFileName)
	for _, fpath := range files {
		if !file.Exists(fpath) {
			continue
		}

		statements, err := query.LoadStatementsFromFile(ctx, proc.Tx, parser.Identifier{Literal: fpath})
		if err != nil {
			return err
		}

		if _, err := proc.Execute(ctx, statements); err != nil {
			return err
		}
	}
	return nil
}

func readQuery(c *cli.Context, tx *query.Transaction) (queryString string, path string, err error) {
	if c.IsSet("source") && 0 < len(c.GlobalString("source")) {
		path = c.GlobalString("source")

		queryString, err = query.LoadContentsFromFile(context.Background(), tx, parser.Identifier{Literal: path})
	} else {
		switch c.NArg() {
		case 0: //Launch interactive shell
		case 1:
			queryString = c.Args().First()
		default:
			err = Exit(query.NewIncorrectCommandUsageError("multiple queries or statements were passed"))
		}
	}
	return
}

func Exit(err error) error {
	if err == nil {
		return nil
	}
	if exit, ok := err.(*query.ForcedExit); ok && exit.Code() == 0 {
		return nil
	}

	code := query.ReturnCodeApplicationError
	message := err.Error()

	if apperr, ok := err.(query.Error); ok {
		code = apperr.Code()
	}
	return cli.NewExitError(cmd.Error(message), code)
}
