package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"

	"github.com/mithrandie/csvq/lib/action"
	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/query"

	"github.com/mithrandie/go-text/color"
	"github.com/urfave/cli"
)

var version = "v1.5.4"

func main() {
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

	defaultCPU := runtime.NumCPU() / 2
	if defaultCPU < 1 {
		defaultCPU = 1
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
			Usage: "field delimiter for csv, or delimiter positions for fixed-length format",
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
			Usage: "format of query results. one of: CSV|TSV|FIXED|JSON|JSONH|JSONA|GFM|ORG|TEXT",
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
			Name:  "pretty-print, P",
			Usage: "make JSON output easier to read in query results",
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

				err := action.ShowFields(table)
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
		err := setFlags(c)
		if err != nil {
			return NewExitError(err.Error(), 1)
		}
		return nil
	}

	app.Action = func(c *cli.Context) error {
		queryString, err := readQuery(c)
		if err != nil {
			return NewExitError(err.Error(), 1)
		}

		if len(queryString) < 1 {
			err = action.LaunchInteractiveShell()
		} else {
			err = action.Run(queryString, cmd.GetFlags().Source)
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

func readQuery(c *cli.Context) (string, error) {
	var queryString string

	flags := cmd.GetFlags()
	if 0 < len(flags.Source) {
		fp, err := os.Open(flags.Source)
		if err != nil {
			return queryString, err
		}
		defer fp.Close()

		buf, err := ioutil.ReadAll(fp)
		if err != nil {
			return queryString, err
		}
		queryString = string(buf)

	} else {
		if 1 < c.NArg() {
			return queryString, errors.New("multiple queries or statements were passed")
		}
		queryString = c.Args().First()
	}

	return queryString, nil
}

func setFlags(c *cli.Context) error {
	cmd.SetColor(c.GlobalBool("color"))

	if err := cmd.SetRepository(c.GlobalString("repository")); err != nil {
		return err
	}
	if err := cmd.SetLocation(c.String("timezone")); err != nil {
		return err
	}
	cmd.SetDatetimeFormat(c.GlobalString("datetime-format"))
	cmd.SetWaitTimeout(c.GlobalFloat64("wait-timeout"))

	if err := cmd.SetSource(c.GlobalString("source")); err != nil {
		return err
	}

	if err := cmd.SetDelimiter(c.GlobalString("delimiter")); err != nil {
		return err
	}
	cmd.SetJsonQuery(c.GlobalString("json-query"))
	if err := cmd.SetEncoding(c.GlobalString("encoding")); err != nil {
		return err
	}
	cmd.SetNoHeader(c.GlobalBool("no-header"))
	cmd.SetWithoutNull(c.GlobalBool("without-null"))

	if err := cmd.SetOut(c.GlobalString("out")); err != nil {
		return err
	}
	if err := cmd.SetFormat(c.GlobalString("format")); err != nil {
		return err
	}
	if err := cmd.SetWriteEncoding(c.GlobalString("write-encoding")); err != nil {
		return err
	}
	if err := cmd.SetWriteDelimiter(c.GlobalString("write-delimiter")); err != nil {
		return err
	}
	cmd.SetWithoutHeader(c.GlobalBool("without-header"))
	if err := cmd.SetLineBreak(c.String("line-break")); err != nil {
		return err
	}
	cmd.SetPrettyPrint(c.GlobalBool("pretty-print"))

	cmd.SetQuiet(c.GlobalBool("quiet"))
	cmd.SetCPU(c.GlobalInt("cpu"))
	cmd.SetStats(c.GlobalBool("stats"))

	return nil
}

func NewExitError(message string, code int) *cli.ExitError {
	return cli.NewExitError(color.Error(message), code)
}
