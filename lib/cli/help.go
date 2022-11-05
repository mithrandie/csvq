package cli

var appHHelpTemplate = `Name:
   {{.Name}}{{if .Usage}} - {{.Usage}}{{end}}

     https://mithrandie.github.io/csvq/

Usage:
   {{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}} {{if .VisibleFlags}}[options]{{end}}{{if .Commands}} [subcommand]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}{{if .Version}}{{if not .HideVersion}}

Version:
   Version {{.Version}}{{end}}{{end}}{{if .Description}}

Description:
   {{.Description}}{{end}}{{if len .Authors}}

Author{{with $length := len .Authors}}{{if ne 1 $length}}S{{end}}{{end}}:
   {{range $index, $author := .Authors}}{{if $index}}
   {{end}}{{$author}}{{end}}{{end}}{{if .VisibleCommands}}

Subcommands:{{range .VisibleCategories}}{{if .Name}}
   {{.Name}}:{{end}}{{range .VisibleCommands}}
   {{join .Names ", "}}{{"\t"}}{{.Usage}}{{end}}{{end}}{{end}}{{if .VisibleFlags}}

Options:
   {{range $index, $option := .VisibleFlags}}{{if $index}}
   {{end}}{{$option}}{{end}}

Parameters:
   Timezone
      Local | UTC
   Import Format
      CSV | TSV | FIXED | JSON | JSONL | LTSV
   Export Format
      CSV | TSV | FIXED | JSON | JSONL | LTSV | GFM | ORG | BOX | TEXT
   Import Character Encodings
      AUTO | UTF8 | UTF8M | UTF16 | UTF16BE | UTF16LE | UTF16BEM | UTF16LEM | SJIS
   Export Character Encodings
      UTF8 | UTF8M | UTF16 | UTF16BE | UTF16LE | UTF16BEM | UTF16LEM | SJIS
   Line Break
      CRLF | CR | LF
   JSON Escape Type
      BACKSLASH | HEX | HEXALL{{end}}{{if .Copyright}}

Copyright:
   {{.Copyright}}{{end}}

`

var commandHelpTemplate = `Name:
   {{.HelpName}} - {{.Usage}}

Usage:
   csvq [options] {{.Name}}{{if .VisibleFlags}} [subcommand options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{if .Category}}

Category:
   {{.Category}}{{end}}{{if .Description}}

Description:
   {{.Description}}{{end}}{{if .VisibleFlags}}

Options:
   {{range .VisibleFlags}}{{.}}
   {{end}}{{end}}
`
