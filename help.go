package main

var appHHelpTemplate = `NAME:
   {{.Name}}{{if .Usage}} - {{.Usage}}{{end}}

     https://mithrandie.github.io/csvq

USAGE:
   {{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}} {{if .VisibleFlags}}[options]{{end}}{{if .Commands}} [subcommand]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}{{if .Version}}{{if not .HideVersion}}

VERSION:
   Version {{.Version}}{{end}}{{end}}{{if .Description}}

DESCRIPTION:
   {{.Description}}{{end}}{{if len .Authors}}

AUTHOR{{with $length := len .Authors}}{{if ne 1 $length}}S{{end}}{{end}}:
   {{range $index, $author := .Authors}}{{if $index}}
   {{end}}{{$author}}{{end}}{{end}}{{if .VisibleCommands}}

SUBCOMMANDS:{{range .VisibleCategories}}{{if .Name}}
   {{.Name}}:{{end}}{{range .VisibleCommands}}
     {{join .Names ", "}}{{"\t"}}{{.Usage}}{{end}}{{end}}{{end}}{{if .VisibleFlags}}

OPTIONS:
   {{range $index, $option := .VisibleFlags}}{{if $index}}
   {{end}}{{$option}}{{end}}

PARAMETERS:
   Timezone
       Local | UTC
   Import Format
       CSV | TSV | FIXED | JSON | LTSV
   Export Format
       CSV | TSV | FIXED | JSON | LTSV | GFM | ORG | TEXT
   Import Character Encodings
       AUTO | UTF8 | UTF8M | UTF16 | UTF16BE | UTF16LE | UTF16BEM | UTF16LEM | SJIS
   Export Character Encodings
       UTF8 | UTF8M | UTF16 | UTF16BE | UTF16LE | UTF16BEM | UTF16LEM | SJIS
   Line Break
       CRLF | CR | LF
   JSON Escape Type
       BACKSLASH | HEX | HEXALL{{end}}{{if .Copyright}}

COPYRIGHT:
   {{.Copyright}}{{end}}

`

var commandHelpTemplate = `NAME:
   {{.HelpName}} - {{.Usage}}

USAGE:
   csvq [options] {{.Name}}{{if .VisibleFlags}} [subcommand options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}
{{if .Category}}
CATEGORY:
   {{.Category}}{{end}}{{if .Description}}

DESCRIPTION:
   {{.Description}}{{end}}{{if .VisibleFlags}}
OPTIONS:
   {{range .VisibleFlags}}{{.}}
   {{end}}{{end}}
`
