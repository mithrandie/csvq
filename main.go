package main

import (
	"errors"
	"io/ioutil"
	"os"

	"github.com/mithrandie/csvq/lib/action"
	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/query"

	"github.com/urfave/cli"
)

var version = "v0.7.10"

func main() {
	cli.AppHelpTemplate = appHHelpTemplate
	cli.CommandHelpTemplate = commandHelpTemplate

	app := cli.NewApp()

	app.Name = "csvq"
	app.Usage = "SQL like query language for csv"
	app.ArgsUsage = "[\"query\"|\"statements\"|argument]"
	app.Version = version

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "delimiter, d",
			Usage: "field delimiter. Default is \",\" for csv files, \"\\t\" for tsv files.",
		},
		cli.StringFlag{
			Name:  "encoding, e",
			Value: "UTF8",
			Usage: "file encoding. one of: UTF8|SJIS",
		},
		cli.StringFlag{
			Name:  "line-break, l",
			Value: "LF",
			Usage: "line break. one of: CRLF|LF|CR",
		},
		cli.StringFlag{
			Name:  "timezone, z",
			Value: "Local",
			Usage: "default timezone. \"Local\", \"UTC\" or a timezone name(e.g. \"America/Los_Angeles\")",
		},
		cli.StringFlag{
			Name:  "repository, r",
			Usage: "directory path where files are located",
		},
		cli.StringFlag{
			Name:  "source, s",
			Usage: "load query from `FILE`",
		},
		cli.StringFlag{
			Name:  "datetime-format, t",
			Usage: "set datetime format to parse strings",
		},
		cli.StringFlag{
			Name:  "wait-timeout, w",
			Value: "10",
			Usage: "limit of the waiting time in seconds to wait for locked files to be released",
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
			Name:  "write-encoding, E",
			Value: "UTF8",
			Usage: "file encoding. one of: UTF8|SJIS",
		},
		cli.StringFlag{
			Name:  "out, o",
			Usage: "write output to `FILE`",
		},
		cli.StringFlag{
			Name:  "format, f",
			Usage: "output format. one of: CSV|TSV|JSON|TEXT",
		},
		cli.StringFlag{
			Name:  "write-delimiter, D",
			Usage: "field delimiter for CSV",
		},
		cli.BoolFlag{
			Name:  "without-header, N",
			Usage: "when the file format is specified as CSV or TSV, write without the header line",
		},
		cli.BoolFlag{
			Name:  "quiet, q",
			Usage: "suppress operation log output",
		},
		cli.IntFlag{
			Name:  "cpu, p",
			Usage: "hint for the number of cpu cores to be used. 1 - number of cpu cores",
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
					cli.ShowCommandHelp(c, "fields")
					return cli.NewExitError("table is not specified", 1)
				}

				table := c.Args().First()

				err := action.ShowFields(table)
				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}

				return nil
			},
		},
		{
			Name:      "calc",
			Usage:     "Calculate a value from stdin",
			ArgsUsage: "\"expression\"",
			Action: func(c *cli.Context) error {
				if c.NArg() != 1 {
					cli.ShowCommandHelp(c, "calc")
					return cli.NewExitError("expression is empty", 1)
				}

				expr := c.Args().First()
				err := action.Calc(expr)
				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}

				return nil
			},
		},
	}

	app.Before = func(c *cli.Context) error {
		err := setFlags(c)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		return nil
	}

	app.Action = func(c *cli.Context) error {
		queryString, err := readQuery(c)
		if err != nil {
			cli.ShowAppHelp(c)
			return cli.NewExitError(err.Error(), 1)
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
			} else if ex, ok := err.(*query.Exit); ok {
				code = ex.GetCode()
			}
			return cli.NewExitError(err.Error(), code)
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
	if err := cmd.SetDelimiter(c.GlobalString("delimiter")); err != nil {
		return err
	}
	if err := cmd.SetEncoding(c.GlobalString("encoding")); err != nil {
		return err
	}
	if err := cmd.SetLineBreak(c.String("line-break")); err != nil {
		return err
	}
	if err := cmd.SetLocation(c.String("timezone")); err != nil {
		return err
	}
	if err := cmd.SetRepository(c.GlobalString("repository")); err != nil {
		return err
	}
	if err := cmd.SetSource(c.GlobalString("source")); err != nil {
		return err
	}
	cmd.SetDatetimeFormat(c.GlobalString("datetime-format"))
	if err := cmd.SetWaitTimeout(c.GlobalString("wait-timeout")); err != nil {
		return err
	}
	cmd.SetNoHeader(c.GlobalBool("no-header"))
	cmd.SetWithoutNull(c.GlobalBool("without-null"))

	if err := cmd.SetWriteEncoding(c.GlobalString("write-encoding")); err != nil {
		return err
	}
	if err := cmd.SetOut(c.GlobalString("out")); err != nil {
		return err
	}
	if err := cmd.SetFormat(c.GlobalString("format")); err != nil {
		return err
	}
	if err := cmd.SetWriteDelimiter(c.GlobalString("write-delimiter")); err != nil {
		return err
	}
	cmd.SetWithoutHeader(c.GlobalBool("without-header"))

	cmd.SetQuiet(c.GlobalBool("quiet"))
	cmd.SetCPU(c.GlobalInt("cpu"))
	cmd.SetStats(c.GlobalBool("stats"))

	return nil
}
