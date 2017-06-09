package action

import (
	"fmt"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/output"
	"github.com/mithrandie/csvq/lib/query"
)

func Write(input string) error {
	results, err := query.Execute(input)
	if err != nil {
		return err
	}

	flags := cmd.GetFlags()
	var out string

	defer func() {
		if 0 < len(out) {
			output.ToStdout(out)
		}
	}()

	writeDelimiter := flags.WriteDelimiter
	writeEncoding := flags.WriteEncoding
	withoutHeader := flags.WithoutHeader
	format := flags.Format

	for _, result := range results {
		if result.View != nil {
			switch result.Type {
			case query.SELECT:
				flags.WriteDelimiter = writeDelimiter
				flags.WriteEncoding = writeEncoding
				flags.WithoutHeader = withoutHeader
				flags.Format = format
			default:
				flags.WriteDelimiter = result.View.FileInfo.Delimiter
				flags.WriteEncoding = flags.Encoding
				flags.WithoutHeader = flags.NoHeader
				flags.Format = cmd.CSV
			}

			s, err := output.EncodeView(result.View)
			if err != nil {
				return err
			}

			switch result.Type {
			case query.SELECT:
				// Do Nothing
			case query.CREATE_TABLE:
				if err = output.Create(result.View.FileInfo.Path, s); err != nil {
					return err
				}
			default:
				if 0 < result.Count {
					if err = output.Update(result.View.FileInfo.Path, s); err != nil {
						return err
					}
				}
			}

			switch result.Type {
			case query.INSERT:
				out += fmt.Sprintf("%d record(s) inserted on %q\n", result.Count, result.View.FileInfo.Path)
			case query.UPDATE:
				if 0 < result.Count {
					out += fmt.Sprintf("%d record(s) updated on %q\n", result.Count, result.View.FileInfo.Path)
				} else {
					out += fmt.Sprintf("no record updated on %q\n", result.View.FileInfo.Path)
				}
			case query.DELETE:
				if 0 < result.Count {
					out += fmt.Sprintf("%d record(s) deleted on %q\n", result.Count, result.View.FileInfo.Path)
				} else {
					out += fmt.Sprintf("no record deleted on %q\n", result.View.FileInfo.Path)
				}
			case query.CREATE_TABLE:
				out += fmt.Sprintf("file %q is created\n", result.View.FileInfo.Path)
			case query.ADD_COLUMNS:
				out += fmt.Sprintf("%d field(s) added on %q\n", result.Count, result.View.FileInfo.Path)
			case query.DROP_COLUMNS:
				out += fmt.Sprintf("%d field(s) dropped on %q\n", result.Count, result.View.FileInfo.Path)
			case query.RENAME_COLUMN:
				out += fmt.Sprintf("%d field(s) renamed on %q\n", result.Count, result.View.FileInfo.Path)
			default:
				out += s
			}
		}

		if 0 < len(result.Log) {
			out += result.Log + "\n"
		}
	}

	if 0 < len(flags.OutFile) {
		if err = output.Create(flags.OutFile, out); err != nil {
			return err
		}
		out = ""
	}

	return nil
}
