package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"github.com/mithrandie/csvq/lib/action"
	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/file"
	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/query"

	"github.com/mithrandie/go-text/color"
	"github.com/urfave/cli"
)

var version = "v1.6.6"

func main() {
	var proc *query.Procedure

	var defaultCPU = runtime.NumCPU() / 2
	if defaultCPU < 1 {
		defaultCPU = 1
	}
	query.Version = version

	cli.AppHelpTemplate = appHHelpTemplate
	cli.CommandHelpTemplate = commandHelpTemplate

	app := cli.NewApp()

	app.Name = "csvq"
	app.Usage = "SQL like query language for csv"
	app.ArgsUsage = "[\"query\"|argument]"
	app.Version = version

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
			Usage: "default timezone. \"Local\", \"UTC\" or a timezone name(e.g. \"America/Los_Angeles\")",
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
			Name:  "delimiter, d",
			Value: ",",
			Usage: "field delimiter for CSV, or delimiter positions for Fixed-Length Format",
		},
		cli.StringFlag{
			Name:  "json-query, j",
			Usage: "`QUERY` for JSON data passed from standard input",
		},
		cli.StringFlag{
			Name:  "encoding, e",
			Value: "UTF8",
			Usage: "file encoding. one of: UTF8|SJIS",
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
			Usage: "export query results and logs to `FILE`",
		},
		cli.StringFlag{
			Name:  "format, f",
			Value: "TEXT",
			Usage: "format of query results. one of: CSV|TSV|FIXED|JSON|GFM|ORG|TEXT|JSONH|JSONA",
		},
		cli.StringFlag{
			Name:  "write-encoding, E",
			Value: "UTF8",
			Usage: "character encoding of query results. one of: UTF8|SJIS",
		},
		cli.StringFlag{
			Name:  "write-delimiter, D",
			Value: ",",
			Usage: "field delimiter or delimiter positions in query results",
		},
		cli.BoolFlag{
			Name:  "without-header, N",
			Usage: "write without the header line in query results",
		},
		cli.StringFlag{
			Name:  "line-break, l",
			Value: "LF",
			Usage: "line break in query results. one of: CRLF|LF|CR",
		},
		cli.BoolFlag{
			Name:  "enclose-all, Q",
			Usage: "enclose all string values in CSV",
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
			Value: defaultCPU,
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
				err := action.Calc(expr)
				if err != nil {
					return NewExitError(err.Error(), 1)
				}

				return nil
			},
			OnUsageError: func(c *cli.Context, err error, isSubcommand bool) error {
				return NewExitError(fmt.Sprintf("Incorrect Usage: %s", err.Error()), 1)
			},
		},
	}

	app.Before = func(c *cli.Context) error {
		action.SetSignalHandler()
		color.UseEffect = false

		proc = query.NewProcedure()

		// Init Single Objects
		if _, err := cmd.GetEnvironment(); err != nil {
			return NewExitError(err.Error(), 1)
		}
		if _, err := cmd.GetPalette(); err != nil {
			return NewExitError(err.Error(), 1)
		}
		cmd.GetFlags()

		// Run pre-load commands
		if err := runPreloadCommands(proc); err != nil {
			return NewExitError(err.Error(), 1)
		}

		// Overwrite Flags with Command Options
		if err := overwriteFlags(c); err != nil {
			return NewExitError(err.Error(), 1)
		}
		return nil
	}

	app.Action = func(c *cli.Context) error {
		queryString, path, err := readQuery(c)
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
			if apperr, ok := err.(query.AppError); ok {
				code = apperr.GetCode()
			} else if ex, ok := err.(*query.ForcedExit); ok {
				code = ex.GetCode()
			}
			return NewExitError(err.Error(), code)
		}

		return nil
	}

	app.Run(os.Args)
}

func readQuery(c *cli.Context) (string, string, error) {
	var queryString string
	var path string

	if c.IsSet("source") && 0 < len(c.GlobalString("source")) {
		path = c.GlobalString("source")
		if abs, err := filepath.Abs(path); err == nil {
			path = abs
		}
		if !file.Exists(path) {
			return queryString, path, errors.New(fmt.Sprintf("file %q does not exist", path))
		}
		h, err := file.NewHandlerForRead(path)
		if err != nil {
			return queryString, path, errors.New(fmt.Sprintf("failed to read file: %s", err.Error()))
		}
		defer h.Close()

		buf, err := ioutil.ReadAll(h.FileForRead())
		if err != nil {
			return queryString, path, errors.New(fmt.Sprintf("failed to read file: %s", err.Error()))
		}

		queryString = string(buf)
	} else {
		if 1 < c.NArg() {
			return queryString, path, errors.New("multiple queries or statements were passed")
		}
		queryString = c.Args().First()
	}

	return queryString, path, nil
}

func overwriteFlags(c *cli.Context) error {
	flags := cmd.GetFlags()
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
		flags.SetWaitTimeout(c.GlobalFloat64("wait-timeout"))
	}

	if c.IsSet("delimiter") {
		if err := flags.SetDelimiter(c.GlobalString("delimiter")); err != nil {
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

func runPreloadCommands(proc *query.Procedure) error {
	handlers := make([]*file.Handler, 0, 4)
	defer func() {
		for _, h := range handlers {
			h.Close()
		}
	}()

	files := cmd.GetSpecialFilePath(cmd.PreloadCommandFileName)
	for _, fpath := range files {
		if !file.Exists(fpath) {
			continue
		}

		statements, err := query.LoadStatementsFromFile(parser.Source{}, fpath)
		if err != nil {
			if e, ok := err.(*query.ReadFileError); ok {
				err = errors.New(e.ErrorMessage())
			}
			return err
		}

		if _, err := proc.Execute(statements); err != nil {
			return err
		}
	}
	return nil
}

func NewExitError(message string, code int) *cli.ExitError {
	return cli.NewExitError(cmd.Error(message), code)
}
