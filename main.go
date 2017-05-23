package main

import (
	"os"

	"github.com/mithrandie/csvq/lib/cmd"

	"github.com/urfave/cli"
)

var version = "0.0.0"

func main() {
	cli.AppHelpTemplate = appHHelpTemplate
	cli.CommandHelpTemplate = commandHelpTemplate

	app := cli.NewApp()

	app.Name = "csvq"
	app.Usage = "SQL like query language for csv"
	app.Version = version

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "delimiter, d",
			Value: ",",
			Usage: "field delimiter (exam: , for comma, \"\\t\" for tab)",
		},
		cli.StringFlag{
			Name:  "encoding, e",
			Value: "utf8",
			Usage: "file encoding. one of: utf8|sjis",
		},
		cli.StringFlag{
			Name:  "repository, r",
			Usage: "directory path where files are located",
		},
		cli.BoolFlag{
			Name:  "no-header",
			Usage: "import first line as a record",
		},
		cli.BoolFlag{
			Name:  "without-null",
			Usage: "parse empty field as empty string",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "write",
			Usage: "Write output to file",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "write-encoding, E",
					Value: "utf8",
					Usage: "file encoding. one of: utf8|sjis",
				},
				cli.StringFlag{
					Name:  "line-break, l",
					Value: "lf",
					Usage: "line break. one of: crlf|lf|cr",
				},
				cli.StringFlag{
					Name:  "out, o",
					Usage: "write output to `FILE`",
				},
				cli.StringFlag{
					Name:  "format, f",
					Usage: "output format. one of: csv|tsv|json|text",
				},
				cli.BoolFlag{
					Name:  "without-header",
					Usage: "when format is specified as csv or tsv, write without header line",
				},
			},
			Before: func(c *cli.Context) error {
				return setWriteFlags(c)
			},
			Action: func(c *cli.Context) error {
				if c.NArg() != 1 {
					return cli.ShowSubcommandHelp(c)
				}

				q := c.Args().First()

				err := Write(q)
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
		if c.NArg() != 1 {
			return cli.ShowAppHelp(c)
		}

		q := c.Args().First()

		err := Write(q)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		return nil
	}

	app.Run(os.Args)
}

func setGlobalFlags(c *cli.Context) error {
	if err := cmd.SetDelimiter(c.GlobalString("delimiter")); err != nil {
		return err
	}
	if err := cmd.SetEncoding(c.GlobalString("encoding")); err != nil {
		return err
	}
	if err := cmd.SetRepository(c.GlobalString("repository")); err != nil {
		return err
	}
	if err := cmd.SetNoHeader(c.GlobalBool("no-header")); err != nil {
		return err
	}
	if err := cmd.SetWithoutNull(c.GlobalBool("without-null")); err != nil {
		return err
	}
	return nil
}

func setWriteFlags(c *cli.Context) error {
	if err := cmd.SetWriteEncoding(c.String("write-encoding")); err != nil {
		return err
	}
	if err := cmd.SetLineBreak(c.String("line-break")); err != nil {
		return err
	}
	if err := cmd.SetOut(c.String("out")); err != nil {
		return err
	}
	if err := cmd.SetFormat(c.String("format")); err != nil {
		return err
	}
	if err := cmd.SetWithoutHeader(c.Bool("without-header")); err != nil {
		return err
	}
	return nil
}
