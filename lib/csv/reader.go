package csv

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/value"
)

var EOF = io.EOF

type Field []byte

func NewField(s string) Field {
	return []byte(s)
}

func (f Field) ToPrimary() value.Primary {
	if f == nil {
		return value.NewNull()
	} else {
		return value.NewString(string(f))
	}
}

type Reader struct {
	Delimiter   rune
	WithoutNull bool

	reader *bufio.Reader
	line   int
	column int

	recordBuf     bytes.Buffer
	fieldStartPos []int
	fieldQuoted   []bool

	FieldsPerRecord int

	LineBreak cmd.LineBreak
}

func NewReader(r io.Reader) *Reader {
	return &Reader{
		Delimiter:       ',',
		WithoutNull:     false,
		reader:          bufio.NewReader(r),
		line:            1,
		column:          0,
		FieldsPerRecord: 0,
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
		header[i] = string(v)
	}
	return header, nil
}

func (r *Reader) Read() ([]Field, error) {
	record, err := r.parseRecord(r.WithoutNull)
	if err != nil {
		return nil, err
	}
	return record, nil
}

func (r *Reader) ReadAll() ([][]Field, error) {
	records := make([][]Field, 0)

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

func (r *Reader) parseRecord(withoutNull bool) ([]Field, error) {
	r.recordBuf.Reset()
	r.fieldStartPos = r.fieldStartPos[:0]
	r.fieldQuoted = r.fieldQuoted[:0]

	fieldIndex := 0
	fieldPosition := 0
	for {
		if 0 < r.FieldsPerRecord && r.FieldsPerRecord <= fieldIndex {
			return nil, r.newError("wrong number of fields in line")
		}

		fieldPosition = r.recordBuf.Len()
		quoted, eol, err := r.parseField()

		if err != nil {
			if err == EOF {
				if fieldIndex < 1 && r.recordBuf.Len() < 1 {
					return nil, EOF
				}
			} else {
				return nil, err
			}
		}

		if eol && fieldIndex < 1 && r.recordBuf.Len() < 1 {
			continue
		}

		r.fieldStartPos = append(r.fieldStartPos, fieldPosition)
		r.fieldQuoted = append(r.fieldQuoted, quoted)
		fieldIndex++

		if eol {
			break
		}
	}

	if r.FieldsPerRecord < 1 {
		r.FieldsPerRecord = fieldIndex
	} else if fieldIndex < r.FieldsPerRecord {
		r.line--
		return nil, r.newError("wrong number of fields in line")
	}

	record := make([]Field, 0, r.FieldsPerRecord)
	recordStr := make([]byte, r.recordBuf.Len())
	copy(recordStr, r.recordBuf.Bytes())
	for i, pos := range r.fieldStartPos {
		var endPos int
		if i == len(r.fieldStartPos)-1 {
			endPos = r.recordBuf.Len()
		} else {
			endPos = r.fieldStartPos[i+1]
		}

		if !withoutNull && pos == endPos && !r.fieldQuoted[i] {
			record = append(record, nil)
		} else {
			record = append(record, recordStr[pos:endPos])
		}
	}

	return record, nil
}

func (r *Reader) parseField() (bool, bool, error) {
	var eof error
	eol := false
	startPos := r.recordBuf.Len()

	quoted := false
	escaped := false

	var lineBreak cmd.LineBreak

Read:
	for {
		lineBreak = ""

		r1, _, err := r.reader.ReadRune()
		r.column++

		if err != nil {
			if err == io.EOF {
				if !escaped && quoted {
					return quoted, eol, r.newError("extraneous \" in field")
				}
				eol = true
			}
			return quoted, eol, err
		}

		switch r1 {
		case '\r':
			r2, _, _ := r.reader.ReadRune()
			if r2 == '\n' {
				lineBreak = cmd.CRLF
			} else {
				r.reader.UnreadRune()
				lineBreak = cmd.CR
			}
			r1 = '\n'
		case '\n':
			lineBreak = cmd.LF
		}
		if r1 == '\n' {
			r.line++
			r.column = 0
		}

		if quoted {
			if escaped {
				switch r1 {
				case '"':
					escaped = false
					r.recordBuf.WriteRune(r1)
					continue
				case r.Delimiter:
					break Read
				case '\n':
					if r.LineBreak == "" {
						r.LineBreak = lineBreak
					}
					eol = true
					break Read
				default:
					r.column--
					return quoted, eol, r.newError("unexpected \" in field")
				}
			}

			switch r1 {
			case '"':
				escaped = true
			case '\n':
				r.recordBuf.WriteString(lineBreak.Value())
			default:
				r.recordBuf.WriteRune(r1)
			}
			continue
		}

		switch r1 {
		case '\n':
			if r.LineBreak == "" {
				r.LineBreak = lineBreak
			}
			eol = true
			break Read
		case r.Delimiter:
			break Read
		case '"':
			if startPos == r.recordBuf.Len() {
				quoted = true
			} else {
				r.recordBuf.WriteRune(r1)
			}
		default:
			r.recordBuf.WriteRune(r1)
		}
	}

	return quoted, eol, eof
}
