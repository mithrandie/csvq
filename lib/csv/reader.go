package csv

import (
	"bufio"
	"errors"
	"fmt"
	"io"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/parser"
)

var EOF error = io.EOF

type Reader struct {
	Delimiter   rune
	WithoutNull bool

	reader *bufio.Reader

	line   int
	column int

	FieldsPerRecords int

	LineBreak cmd.LineBreak
}

func NewReader(r io.Reader) *Reader {
	return &Reader{
		Delimiter:        ',',
		WithoutNull:      false,
		reader:           bufio.NewReader(r),
		line:             1,
		column:           0,
		FieldsPerRecords: 0,
	}
}

func (r *Reader) newError(s string) error {
	return errors.New(fmt.Sprintf("line %d, column %d: %s", r.line, r.column, s))
}

func (r *Reader) ReadHeader() ([]string, error) {
	record, err := r.parseRecord(true)
	if err != nil {
		return nil, err
	}

	header := make([]string, len(record))
	for i, v := range record {
		header[i] = v.(parser.String).Value()
	}
	return header, nil
}

func (r *Reader) Read() ([]parser.Primary, error) {
	record, err := r.parseRecord(r.WithoutNull)
	if err != nil {
		return nil, err
	}
	return record, nil
}

func (r *Reader) ReadAll() ([][]parser.Primary, error) {
	records := [][]parser.Primary{}

	for {
		record, err := r.Read()
		if err == EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	return records, nil
}

func (r *Reader) parseRecord(withoutNull bool) ([]parser.Primary, error) {
	var record []parser.Primary
	if r.FieldsPerRecords < 1 {
		record = []parser.Primary{}
	} else {
		record = make([]parser.Primary, r.FieldsPerRecords)
	}

	fieldIndex := 0
	for {
		if 0 < r.FieldsPerRecords && r.FieldsPerRecords <= fieldIndex {
			return nil, r.newError("wrong number of fields in line")
		}

		field, eol, err := r.parseField(withoutNull)
		if err == EOF {
			if fieldIndex < 1 && (parser.IsNull(field) || len(field.(parser.String).Value()) < 1) {
				return nil, EOF
			}
		} else if err != nil {
			return nil, err
		}

		if r.FieldsPerRecords < 1 {
			record = append(record, field)
		} else {
			record[fieldIndex] = field
		}

		if eol {
			break
		}

		fieldIndex++
	}

	if r.FieldsPerRecords < 1 {
		r.FieldsPerRecords = fieldIndex + 1
	}

	if r.FieldsPerRecords != fieldIndex+1 {
		return nil, r.newError("wrong number of fields in line")
	}

	return record, nil
}

func (r *Reader) parseField(stringOnly bool) (parser.Primary, bool, error) {
	var eof error
	eol := false

	quoted := false
	escaped := false

	field := []rune{}

Read:
	for {
		r1, err := r.readRune()
		if err == io.EOF {
			if !escaped && quoted {
				return nil, eol, r.newError("extraneous \" in field")
			}

			eof = EOF
			eol = true
			break
		}
		if err != nil {
			return nil, eol, err
		}

		if escaped {
			switch r1 {
			case '"':
				escaped = false
				field = append(field, r1)
				continue
			case r.Delimiter:
				break Read
			case '\n':
				eol = true
				break Read
			default:
				r.column--
				return nil, eol, r.newError("unexpected \" in field")
			}
		}

		if quoted {
			if r1 == '"' {
				escaped = true
				continue
			}
			field = append(field, r1)
			continue
		}

		switch r1 {
		case '\n':
			eol = true
			break Read
		case r.Delimiter:
			break Read
		case '"':
			if len(field) < 1 {
				quoted = true
			} else {
				field = append(field, r1)
			}
		default:
			field = append(field, r1)
		}
	}

	s := string(field)

	if !stringOnly && len(s) < 1 && !quoted {
		return parser.NewNull(), eol, eof
	}

	return parser.NewString(s), eol, eof
}

func (r *Reader) readRune() (rune, error) {
	r1, _, err := r.reader.ReadRune()
	r.column++

	if err != nil {
		return r1, err
	}

	if r.isNewLine(r1) {
		r.line++
		r.column = 0
		return '\n', nil
	}
	return r1, nil
}

func (r *Reader) isNewLine(r1 rune) bool {
	switch r1 {
	case '\r':
		r2, _, _ := r.reader.ReadRune()
		if r2 == '\n' {
			if r.LineBreak == "" {
				r.LineBreak = cmd.CRLF
			}
			return true
		}
		r.reader.UnreadRune()
		if r.LineBreak == "" {
			r.LineBreak = cmd.CR
		}
		return true
	case '\n':
		if r.LineBreak == "" {
			r.LineBreak = cmd.LF
		}
		return true
	}
	return false
}
