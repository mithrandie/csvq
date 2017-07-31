package main

import (
	"errors"
	"io/ioutil"
	"os"

	"github.com/mithrandie/csvq/lib/action"
	"github.com/mithrandie/csvq/lib/cmd"

	"github.com/urfave/cli"
)

var version = "v0.4.11"

func main() {
	cli.AppHelpTemplate = appHHelpTemplate
	cli.CommandHelpTemplate = commandHelpTemplate

	app := cli.NewApp()

	app.Name = "csvq"
	app.Usage = "SQL like query language for csv"
	app.ArgsUsage = "[\"query\"|\"statements\"]"
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
		cli.BoolFlag{
			Name:  "no-header, n",
			Usage: "import the first line as a record",
		},
		cli.BoolFlag{
			Name:  "without-null, a",
			Usage: "parse empty fields as empty strings",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:      "write",
			Usage:     "Write output to a file",
			ArgsUsage: "[\"query\"|\"statements\"]",

			Flags: []cli.Flag{
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
			},
			Before: func(c *cli.Context) error {
				return setWriteFlags(c)
			},
			Action: func(c *cli.Context) error {
				query, err := readQuery(c)
				if err != nil {
					cli.ShowCommandHelp(c, "write")
					return cli.NewExitError(err.Error(), 1)
				}

				err = action.Write(query, cmd.GetFlags().Source)
				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}

				return nil
			},
		},
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
		return setGlobalFlags(c)
	}

	app.Action = func(c *cli.Context) error {
		query, err := readQuery(c)
		if err != nil {
			cli.ShowAppHelp(c)
			return cli.NewExitError(err.Error(), 1)
		}

		err = action.Write(query, cmd.GetFlags().Source)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		return nil
	}

	app.Run(os.Args)
}

func readQuery(c *cli.Context) (string, error) {
	var query string

	flags := cmd.GetFlags()
	if 0 < len(flags.Source) {
		fp, err := os.Open(flags.Source)
		if err != nil {
			return query, err
		}
		defer fp.Close()

		buf, err := ioutil.ReadAll(fp)
		if err != nil {
			return query, err
		}
		query = string(buf)

	} else {
		if c.NArg() != 1 {
			return query, errors.New("query is empty")
		}
		query = c.Args().First()
	}

	return query, nil
}

func setGlobalFlags(c *cli.Context) error {
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
	cmd.SetNoHeader(c.GlobalBool("no-header"))
	cmd.SetWithoutNull(c.GlobalBool("without-null"))
	return nil
}

func setWriteFlags(c *cli.Context) error {
	if err := cmd.SetWriteEncoding(c.String("write-encoding")); err != nil {
		return err
	}
	if err := cmd.SetOut(c.String("out")); err != nil {
		return err
	}
	if err := cmd.SetFormat(c.String("format")); err != nil {
		return err
	}
	if err := cmd.SetWriteDelimiter(c.String("write-delimiter")); err != nil {
		return err
	}
	cmd.SetWithoutHeader(c.Bool("without-header"))
	return nil
}
