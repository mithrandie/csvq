package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/mithrandie/csvq/lib/action"
	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/file"
	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/query"

	"github.com/urfave/cli"
)

func main() {
	action.CurrentVersion, _ = action.ParseVersion(query.Version)
	if !action.CurrentVersion.IsEmpty() {
		query.Version = action.CurrentVersion.String()
	}

	cli.AppHelpTemplate = appHHelpTemplate
	cli.CommandHelpTemplate = commandHelpTemplate

	app := cli.NewApp()

	app.Name = "csvq"
	app.Usage = "SQL-like query language for csv"
	app.ArgsUsage = "[query|argument]"
	app.Version = query.Version
	app.OnUsageError = onUsageError
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "repository, r",
			Usage: "directory `PATH` where files are located",
		},
		cli.StringFlag{
			Name:  "timezone, z",
			Value: "Local",
			Usage: "default timezone",
		},
		cli.StringFlag{
			Name:  "datetime-format, t",
			Usage: "datetime format to parse strings",
		},
		cli.BoolFlag{
			Name:  "ansi-quotes, k",
			Usage: "use double quotation mark as identifier enclosure",
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
			Usage: "default format to load files",
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
			Value: "AUTO",
			Usage: "file encoding",
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
		cli.BoolFlag{
			Name:  "strip-ending-line-break, T",
			Usage: "strip line break from the end of files and query results",
		},
		cli.StringFlag{
			Name:  "format, f",
			Value: "TEXT",
			Usage: "format of query results",
		},
		cli.StringFlag{
			Name:  "write-encoding, E",
			Value: "UTF8",
			Usage: "character encoding of query results",
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
			Usage: "line break in query results",
		},
		cli.BoolFlag{
			Name:  "enclose-all, Q",
			Usage: "enclose all string values in CSV and TSV",
		},
		cli.StringFlag{
			Name:  "json-escape, J",
			Value: "BACKSLASH",
			Usage: "JSON escape type",
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
			Name:  "limit-recursion",
			Value: 1000,
			Usage: "maximum number of iterations for recursive queries",
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
			ArgsUsage: "DATA_FILE_PATH",
			Action: commandAction(func(ctx context.Context, c *cli.Context, proc *query.Processor) error {
				if 1 != c.NArg() {
					return query.NewIncorrectCommandUsageError("fields subcommand takes exactly 1 argument")
				}
				table := c.Args().First()
				return action.ShowFields(ctx, proc, table)
			}),
		},
		{
			Name:      "calc",
			Usage:     "Calculate a value from stdin",
			ArgsUsage: "expression",
			Action: commandAction(func(ctx context.Context, c *cli.Context, proc *query.Processor) error {
				if !proc.Tx.Session.CanReadStdin {
					return query.NewIncorrectCommandUsageError(query.ErrMsgStdinEmpty)
				}
				if 1 != c.NArg() {
					return query.NewIncorrectCommandUsageError("calc subcommand takes exactly 1 argument")
				}
				expr := c.Args().First()
				return action.Calc(ctx, proc, expr)
			}),
		},
		{
			Name:      "syntax",
			Usage:     "Print syntax",
			ArgsUsage: "[search_word ...]",
			Action: commandAction(func(ctx context.Context, c *cli.Context, proc *query.Processor) error {
				words := append([]string{c.Args().First()}, c.Args().Tail()...)
				return action.Syntax(ctx, proc, words)
			}),
		},
		{
			Name:      "check-update",
			Usage:     "Check for updates",
			ArgsUsage: " ",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "include-pre-release",
					Usage: "check including pre-release version",
				},
			},
			Action: commandAction(func(ctx context.Context, c *cli.Context, proc *query.Processor) error {
				if 0 < c.NArg() {
					return query.NewIncorrectCommandUsageError("check-update subcommand takes no argument")
				}

				includePreRelease := false
				if c.IsSet("include-pre-release") {
					includePreRelease = c.Bool("include-pre-release")
				}
				return action.CheckUpdate(includePreRelease)
			}),
		},
	}

	for i := range app.Commands {
		app.Commands[i].OnUsageError = onUsageError
	}

	app.Action = commandAction(func(ctx context.Context, c *cli.Context, proc *query.Processor) error {
		queryString, path, err := readQuery(ctx, c, proc.Tx)
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

func onUsageError(c *cli.Context, err error, isSubcommand bool) error {
	if isSubcommand {
		if e := cli.ShowCommandHelp(c, c.Command.Name); e != nil {
			println(e.Error())
		}
	}
	if _, ok := err.(*query.IncorrectCommandUsageError); !ok {
		err = query.NewIncorrectCommandUsageError(err.Error())
	}
	return Exit(err, nil)
}

func commandAction(fn func(ctx context.Context, c *cli.Context, proc *query.Processor) error) func(c *cli.Context) error {
	return func(c *cli.Context) (err error) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		session := query.NewSession()
		tx, e := query.NewTransaction(ctx, file.DefaultWaitTimeout, file.DefaultRetryDelay, session)
		if e != nil {
			return Exit(e, nil)
		}

		proc := query.NewProcessor(tx)
		defer func() {
			if e := proc.AutoRollback(); e != nil {
				proc.LogError(e.Error())
			}
			if e := proc.ReleaseResourcesWithErrors(); e != nil {
				proc.LogError(e.Error())
			}

			if err != nil {
				if _, ok := err.(*query.IncorrectCommandUsageError); ok {
					err = onUsageError(c, err, 0 < len(c.Command.Name))
				} else {
					err = Exit(err, proc.Tx)
				}
			}
		}()

		// Handle signals
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, action.Signals...)
		var signalReceived error

		go func() {
			sig := <-ch
			signalReceived = query.NewSignalReceived(sig)
			cancel()
		}()

		// Run pre-load commands
		if err = runPreloadCommands(ctx, proc); err != nil {
			return
		}

		// Overwrite Flags with Command Options
		if err = overwriteFlags(c, proc.Tx); err != nil {
			return
		}

		err = fn(ctx, c, proc)
		if signalReceived != nil {
			err = signalReceived
		}
		return
	}
}

func overwriteFlags(c *cli.Context, tx *query.Transaction) error {
	if c.GlobalIsSet("repository") {
		if err := tx.SetFlag(cmd.RepositoryFlag, c.GlobalString("repository")); err != nil {
			return query.NewIncorrectCommandUsageError(err.Error())
		}
	}
	if c.GlobalIsSet("timezone") {
		if err := tx.SetFlag(cmd.TimezoneFlag, c.GlobalString("timezone")); err != nil {
			return query.NewIncorrectCommandUsageError(err.Error())
		}
	}
	if c.GlobalIsSet("datetime-format") {
		_ = tx.SetFlag(cmd.DatetimeFormatFlag, c.GlobalString("datetime-format"))
	}
	if c.GlobalIsSet("ansi-quotes") {
		_ = tx.SetFlag(cmd.AnsiQuotesFlag, c.GlobalBool("ansi-quotes"))
	}

	if c.GlobalIsSet("wait-timeout") {
		_ = tx.SetFlag(cmd.WaitTimeoutFlag, c.GlobalFloat64("wait-timeout"))
	}
	if c.GlobalIsSet("color") {
		_ = tx.SetFlag(cmd.ColorFlag, c.GlobalBool("color"))
	}

	if c.GlobalIsSet("import-format") {
		if err := tx.SetFlag(cmd.ImportFormatFlag, c.GlobalString("import-format")); err != nil {
			return query.NewIncorrectCommandUsageError(err.Error())
		}
	}
	if c.GlobalIsSet("delimiter") {
		if err := tx.SetFlag(cmd.DelimiterFlag, c.GlobalString("delimiter")); err != nil {
			return query.NewIncorrectCommandUsageError(err.Error())
		}
	}
	if c.GlobalIsSet("delimiter-positions") {
		if err := tx.SetFlag(cmd.DelimiterPositionsFlag, c.GlobalString("delimiter-positions")); err != nil {
			return query.NewIncorrectCommandUsageError(err.Error())
		}
	}
	if c.GlobalIsSet("json-query") {
		_ = tx.SetFlag(cmd.JsonQueryFlag, c.GlobalString("json-query"))
	}
	if c.GlobalIsSet("encoding") {
		if err := tx.SetFlag(cmd.EncodingFlag, c.GlobalString("encoding")); err != nil {
			return query.NewIncorrectCommandUsageError(err.Error())
		}
	}
	if c.GlobalIsSet("no-header") {
		_ = tx.SetFlag(cmd.NoHeaderFlag, c.GlobalBool("no-header"))
	}
	if c.GlobalIsSet("without-null") {
		_ = tx.SetFlag(cmd.WithoutNullFlag, c.GlobalBool("without-null"))
	}

	if c.GlobalIsSet("strip-ending-line-break") {
		_ = tx.SetFlag(cmd.StripEndingLineBreakFlag, c.GlobalBool("strip-ending-line-break"))
	}
	if c.GlobalIsSet("format") {
		if err := tx.SetFormatFlag(c.GlobalString("format"), c.GlobalString("out")); err != nil {
			return query.NewIncorrectCommandUsageError(err.Error())
		}
	}
	if c.GlobalIsSet("write-encoding") {
		if err := tx.SetFlag(cmd.ExportEncodingFlag, c.GlobalString("write-encoding")); err != nil {
			return query.NewIncorrectCommandUsageError(err.Error())
		}
	}
	if c.GlobalIsSet("write-delimiter") {
		if err := tx.SetFlag(cmd.ExportDelimiterFlag, c.GlobalString("write-delimiter")); err != nil {
			return query.NewIncorrectCommandUsageError(err.Error())
		}
	}
	if c.GlobalIsSet("write-delimiter-positions") {
		if err := tx.SetFlag(cmd.ExportDelimiterPositionsFlag, c.GlobalString("write-delimiter-positions")); err != nil {
			return query.NewIncorrectCommandUsageError(err.Error())
		}
	}
	if c.GlobalIsSet("without-header") {
		_ = tx.SetFlag(cmd.WithoutHeaderFlag, c.GlobalBool("without-header"))
	}
	if c.GlobalIsSet("line-break") {
		if err := tx.SetFlag(cmd.LineBreakFlag, c.GlobalString("line-break")); err != nil {
			return query.NewIncorrectCommandUsageError(err.Error())
		}
	}
	if c.GlobalIsSet("enclose-all") {
		_ = tx.SetFlag(cmd.EncloseAllFlag, c.GlobalBool("enclose-all"))
	}
	if c.GlobalIsSet("json-escape") {
		if err := tx.SetFlag(cmd.JsonEscapeFlag, c.GlobalString("json-escape")); err != nil {
			return query.NewIncorrectCommandUsageError(err.Error())
		}
	}
	if c.GlobalIsSet("pretty-print") {
		_ = tx.SetFlag(cmd.PrettyPrintFlag, c.GlobalBool("pretty-print"))
	}

	if c.GlobalIsSet("east-asian-encoding") {
		_ = tx.SetFlag(cmd.EastAsianEncodingFlag, c.GlobalBool("east-asian-encoding"))
	}
	if c.GlobalIsSet("count-diacritical-sign") {
		_ = tx.SetFlag(cmd.CountDiacriticalSignFlag, c.GlobalBool("count-diacritical-sign"))
	}
	if c.GlobalIsSet("count-format-code") {
		_ = tx.SetFlag(cmd.CountFormatCodeFlag, c.GlobalBool("count-format-code"))
	}

	if c.GlobalIsSet("quiet") {
		_ = tx.SetFlag(cmd.QuietFlag, c.GlobalBool("quiet"))
	}
	if c.GlobalIsSet("limit-recursion") {
		_ = tx.SetFlag(cmd.LimitRecursion, c.GlobalInt64("limit-recursion"))
	}
	if c.GlobalIsSet("cpu") {
		_ = tx.SetFlag(cmd.CPUFlag, c.GlobalInt64("cpu"))
	}
	if c.GlobalIsSet("stats") {
		_ = tx.SetFlag(cmd.StatsFlag, c.GlobalBool("stats"))
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

func readQuery(ctx context.Context, c *cli.Context, tx *query.Transaction) (queryString string, path string, err error) {
	if c.GlobalIsSet("source") && 0 < len(c.GlobalString("source")) {
		if 0 < c.NArg() {
			err = query.NewIncorrectCommandUsageError("no argument can be passed when \"--source\" option is specified")
		} else {
			path = c.GlobalString("source")
			queryString, err = query.LoadContentsFromFile(ctx, tx, parser.Identifier{Literal: path})
		}
	} else {
		switch c.NArg() {
		case 0:
			// Launch interactive shell
		case 1:
			queryString = c.Args().First()
		default:
			err = query.NewIncorrectCommandUsageError("csvq command takes exactly 1 argument")
		}
	}
	return
}

func Exit(err error, tx *query.Transaction) error {
	if err == nil {
		return nil
	}
	if exit, ok := err.(*query.ForcedExit); ok && exit.Code() == 0 {
		return nil
	}

	code := query.ReturnCodeApplicationError
	message := err.Error()
	if tx != nil {
		message = tx.Error(message)
	}

	if apperr, ok := err.(query.Error); ok {
		code = apperr.Code()
	}

	return cli.NewExitError(message, code)
}
