package cli

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/mithrandie/csvq/lib/action"
	"github.com/mithrandie/csvq/lib/file"
	"github.com/mithrandie/csvq/lib/option"
	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/query"

	"github.com/urfave/cli/v2"
)

func Run() {
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
	app.Authors = []*cli.Author{{Name: "Yuki et al."}}
	app.OnUsageError = onUsageError
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "repository",
			Aliases: []string{"r"},
			Usage:   "directory `PATH` where files are located",
		},
		&cli.StringFlag{
			Name:    "timezone",
			Aliases: []string{"z"},
			Value:   "Local",
			Usage:   "default timezone",
		},
		&cli.StringFlag{
			Name:    "datetime-format",
			Aliases: []string{"t"},
			Usage:   "datetime format to parse strings",
		},
		&cli.BoolFlag{
			Name:    "ansi-quotes",
			Aliases: []string{"k"},
			Usage:   "use double quotation mark as identifier enclosure",
		},
		&cli.BoolFlag{
			Name:    "strict-equal",
			Aliases: []string{"g"},
			Usage:   "compare strictly equal or not in DISTINCT, GROUP BY and ORDER BY",
		},
		&cli.Float64Flag{
			Name:    "wait-timeout",
			Aliases: []string{"w"},
			Value:   10,
			Usage:   "maximum time in seconds to wait for locked files to be released",
		},
		&cli.StringFlag{
			Name:    "source",
			Aliases: []string{"s"},
			Usage:   "load query or statements from `FILE`",
		},
		&cli.StringFlag{
			Name:    "import-format",
			Aliases: []string{"i"},
			Value:   "CSV",
			Usage:   "default format to load files",
		},
		&cli.StringFlag{
			Name:    "delimiter",
			Aliases: []string{"d"},
			Value:   ",",
			Usage:   "field delimiter for CSV",
		},
		&cli.BoolFlag{
			Name:  "allow-uneven-fields",
			Usage: "allow loading CSV files with uneven field length",
		},
		&cli.StringFlag{
			Name:    "delimiter-positions",
			Aliases: []string{"m"},
			Usage:   "delimiter positions for FIXED",
		},
		&cli.StringFlag{
			Name:    "json-query",
			Aliases: []string{"j"},
			Usage:   "`QUERY` for JSON",
		},
		&cli.StringFlag{
			Name:    "encoding",
			Aliases: []string{"e"},
			Value:   "AUTO",
			Usage:   "file encoding",
		},
		&cli.BoolFlag{
			Name:    "no-header",
			Aliases: []string{"n"},
			Usage:   "import the first line as a record",
		},
		&cli.BoolFlag{
			Name:    "without-null",
			Aliases: []string{"a"},
			Usage:   "parse empty fields as empty strings",
		},
		&cli.StringFlag{
			Name:    "out",
			Aliases: []string{"o"},
			Usage:   "export result sets of select queries to `FILE`",
		},
		&cli.BoolFlag{
			Name:    "strip-ending-line-break",
			Aliases: []string{"T"},
			Usage:   "strip line break from the end of files and query results",
		},
		&cli.StringFlag{
			Name:    "format",
			Aliases: []string{"f"},
			Usage:   "format of query results. (default: \"CSV\" for output to pipe, \"TEXT\" otherwise)",
		},
		&cli.StringFlag{
			Name:    "write-encoding",
			Aliases: []string{"E"},
			Value:   "UTF8",
			Usage:   "character encoding of query results",
		},
		&cli.StringFlag{
			Name:    "write-delimiter",
			Aliases: []string{"D"},
			Value:   ",",
			Usage:   "field delimiter for CSV in query results",
		},
		&cli.StringFlag{
			Name:    "write-delimiter-positions",
			Aliases: []string{"M"},
			Usage:   "delimiter positions for FIXED in query results",
		},
		&cli.BoolFlag{
			Name:    "without-header",
			Aliases: []string{"N"},
			Usage:   "export result sets of select queries without the header line",
		},
		&cli.StringFlag{
			Name:    "line-break",
			Aliases: []string{"l"},
			Value:   "LF",
			Usage:   "line break in query results",
		},
		&cli.BoolFlag{
			Name:    "enclose-all",
			Aliases: []string{"Q"},
			Usage:   "enclose all string values in CSV and TSV",
		},
		&cli.StringFlag{
			Name:    "json-escape",
			Aliases: []string{"J"},
			Value:   "BACKSLASH",
			Usage:   "JSON escape type",
		},
		&cli.BoolFlag{
			Name:    "pretty-print",
			Aliases: []string{"P"},
			Usage:   "make JSON output easier to read in query results",
		},
		&cli.BoolFlag{
			Name:    "scientific-notation",
			Aliases: []string{"SN"},
			Usage:   "use scientific notation for large exponents in output",
		},
		&cli.BoolFlag{
			Name:    "east-asian-encoding",
			Aliases: []string{"W"},
			Usage:   "count ambiguous characters as fullwidth",
		},
		&cli.BoolFlag{
			Name:    "count-diacritical-sign",
			Aliases: []string{"S"},
			Usage:   "count diacritical signs as halfwidth",
		},
		&cli.BoolFlag{
			Name:    "count-format-code",
			Aliases: []string{"A"},
			Usage:   "count format characters and zero-width spaces as halfwidth",
		},
		&cli.BoolFlag{
			Name:    "color",
			Aliases: []string{"c"},
			Usage:   "use ANSI color escape sequences",
		},
		&cli.BoolFlag{
			Name:    "quiet",
			Aliases: []string{"q"},
			Usage:   "suppress operation log output",
		},
		&cli.IntFlag{
			Name:  "limit-recursion",
			Value: 1000,
			Usage: "maximum number of iterations for recursive queries",
		},
		&cli.IntFlag{
			Name:    "cpu",
			Aliases: []string{"p"},
			Value:   option.GetDefaultNumberOfCPU(),
			Usage:   "hint for the number of cpu cores to be used",
		},
		&cli.BoolFlag{
			Name:    "stats",
			Aliases: []string{"x"},
			Usage:   "show execution time and memory statistics",
		},
	}

	app.Commands = []*cli.Command{
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
				&cli.BoolFlag{
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
			err = action.Run(ctx, proc, queryString, path, c.String("out"))
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
	if c.IsSet("repository") {
		if err := tx.SetFlag(option.RepositoryFlag, c.String("repository")); err != nil {
			return query.NewIncorrectCommandUsageError(err.Error())
		}
	}
	if c.IsSet("timezone") {
		if err := tx.SetFlag(option.TimezoneFlag, c.String("timezone")); err != nil {
			return query.NewIncorrectCommandUsageError(err.Error())
		}
	}
	if c.IsSet("datetime-format") {
		_ = tx.SetFlag(option.DatetimeFormatFlag, c.String("datetime-format"))
	}
	if c.IsSet("ansi-quotes") {
		_ = tx.SetFlag(option.AnsiQuotesFlag, c.Bool("ansi-quotes"))
	}
	if c.IsSet("strict-equal") {
		_ = tx.SetFlag(option.StrictEqualFlag, c.Bool("strict-equal"))
	}

	if c.IsSet("wait-timeout") {
		_ = tx.SetFlag(option.WaitTimeoutFlag, c.Float64("wait-timeout"))
	}
	if c.IsSet("color") {
		_ = tx.SetFlag(option.ColorFlag, c.Bool("color"))
	}

	if c.IsSet("import-format") {
		if err := tx.SetFlag(option.ImportFormatFlag, c.String("import-format")); err != nil {
			return query.NewIncorrectCommandUsageError(err.Error())
		}
	}
	if c.IsSet("delimiter") {
		if err := tx.SetFlag(option.DelimiterFlag, c.String("delimiter")); err != nil {
			return query.NewIncorrectCommandUsageError(err.Error())
		}
	}
	if c.IsSet("allow-uneven-fields") {
		_ = tx.SetFlag(option.AllowUnevenFieldsFlag, c.Bool("allow-uneven-fields"))
	}
	if c.IsSet("delimiter-positions") {
		if err := tx.SetFlag(option.DelimiterPositionsFlag, c.String("delimiter-positions")); err != nil {
			return query.NewIncorrectCommandUsageError(err.Error())
		}
	}
	if c.IsSet("json-query") {
		_ = tx.SetFlag(option.JsonQueryFlag, c.String("json-query"))
	}
	if c.IsSet("encoding") {
		if err := tx.SetFlag(option.EncodingFlag, c.String("encoding")); err != nil {
			return query.NewIncorrectCommandUsageError(err.Error())
		}
	}
	if c.IsSet("no-header") {
		_ = tx.SetFlag(option.NoHeaderFlag, c.Bool("no-header"))
	}
	if c.IsSet("without-null") {
		_ = tx.SetFlag(option.WithoutNullFlag, c.Bool("without-null"))
	}

	if c.IsSet("strip-ending-line-break") {
		_ = tx.SetFlag(option.StripEndingLineBreakFlag, c.Bool("strip-ending-line-break"))
	}

	if err := tx.SetFormatFlag(c.String("format"), c.String("out")); err != nil {
		return query.NewIncorrectCommandUsageError(err.Error())
	}

	if c.IsSet("write-encoding") {
		if err := tx.SetFlag(option.ExportEncodingFlag, c.String("write-encoding")); err != nil {
			return query.NewIncorrectCommandUsageError(err.Error())
		}
	}
	if c.IsSet("write-delimiter") {
		if err := tx.SetFlag(option.ExportDelimiterFlag, c.String("write-delimiter")); err != nil {
			return query.NewIncorrectCommandUsageError(err.Error())
		}
	}
	if c.IsSet("write-delimiter-positions") {
		if err := tx.SetFlag(option.ExportDelimiterPositionsFlag, c.String("write-delimiter-positions")); err != nil {
			return query.NewIncorrectCommandUsageError(err.Error())
		}
	}
	if c.IsSet("without-header") {
		_ = tx.SetFlag(option.WithoutHeaderFlag, c.Bool("without-header"))
	}
	if c.IsSet("line-break") {
		if err := tx.SetFlag(option.LineBreakFlag, c.String("line-break")); err != nil {
			return query.NewIncorrectCommandUsageError(err.Error())
		}
	}
	if c.IsSet("enclose-all") {
		_ = tx.SetFlag(option.EncloseAllFlag, c.Bool("enclose-all"))
	}
	if c.IsSet("json-escape") {
		if err := tx.SetFlag(option.JsonEscapeFlag, c.String("json-escape")); err != nil {
			return query.NewIncorrectCommandUsageError(err.Error())
		}
	}
	if c.IsSet("pretty-print") {
		_ = tx.SetFlag(option.PrettyPrintFlag, c.Bool("pretty-print"))
	}
	if c.IsSet("scientific-notation") {
		_ = tx.SetFlag(option.ScientificNotationFlag, c.Bool("scientific-notation"))
	}

	if c.IsSet("east-asian-encoding") {
		_ = tx.SetFlag(option.EastAsianEncodingFlag, c.Bool("east-asian-encoding"))
	}
	if c.IsSet("count-diacritical-sign") {
		_ = tx.SetFlag(option.CountDiacriticalSignFlag, c.Bool("count-diacritical-sign"))
	}
	if c.IsSet("count-format-code") {
		_ = tx.SetFlag(option.CountFormatCodeFlag, c.Bool("count-format-code"))
	}

	if c.IsSet("quiet") {
		_ = tx.SetFlag(option.QuietFlag, c.Bool("quiet"))
	}
	if c.IsSet("limit-recursion") {
		_ = tx.SetFlag(option.LimitRecursion, c.Int64("limit-recursion"))
	}
	if c.IsSet("cpu") {
		_ = tx.SetFlag(option.CPUFlag, c.Int64("cpu"))
	}
	if c.IsSet("stats") {
		_ = tx.SetFlag(option.StatsFlag, c.Bool("stats"))
	}

	return nil
}

func runPreloadCommands(ctx context.Context, proc *query.Processor) (err error) {
	files := option.GetSpecialFilePath(option.PreloadCommandFileName)
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
	if c.IsSet("source") && 0 < len(c.String("source")) {
		if 0 < c.NArg() {
			err = query.NewIncorrectCommandUsageError("no argument can be passed when \"--source\" option is specified")
		} else {
			path = c.String("source")
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

	return cli.Exit(message, code)
}
