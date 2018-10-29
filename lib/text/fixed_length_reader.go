package text

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/value"
	"io"
)

type FixedLengthReader struct {
	DelimiterPositions DelimiterPositions
	WithoutNull        bool
	Encoding           cmd.Encoding

	reader *bufio.Reader
	buf    bytes.Buffer

	LineBreak cmd.LineBreak
}

func NewFixedLengthReader(r io.Reader, positions []int) *FixedLengthReader {
	return &FixedLengthReader{
		DelimiterPositions: positions,
		WithoutNull:        false,
		reader:             bufio.NewReader(r),
	}
}

func (r *FixedLengthReader) ReadHeader() ([]string, error) {
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

func (r *FixedLengthReader) Read() ([]value.Field, error) {
	return r.parseRecord(r.WithoutNull)
}

func (r *FixedLengthReader) ReadAll() ([][]value.Field, error) {
	records := make([][]value.Field, 0, 100)

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	return records, nil
}

func (r *FixedLengthReader) parseRecord(withoutNull bool) ([]value.Field, error) {
	record := make([]value.Field, 0, len(r.DelimiterPositions))
	recordPos := 0
	delimiterPos := 0

	var lineBreak cmd.LineBreak
	lineEnd := false

	for _, endPos := range r.DelimiterPositions {
		if endPos < 0 || endPos <= delimiterPos {
			return nil, errors.New(fmt.Sprintf("invalid delimiter position: %s", r.DelimiterPositions))
		}
		delimiterPos = endPos

		r.buf.Reset()
		for !lineEnd && recordPos < delimiterPos {
			c, s, err := r.reader.ReadRune()

			if err != nil {
				if err != io.EOF || recordPos < 1 {
					return nil, err
				}
				lineEnd = true
				continue
			} else {
				switch c {
				case '\r':
					c2, _, _ := r.reader.ReadRune()
					if c2 == '\n' {
						lineBreak = cmd.CRLF
					} else {
						r.reader.UnreadRune()
						lineBreak = cmd.CR
					}
					c = '\n'
				case '\n':
					lineBreak = cmd.LF
				}
				if c == '\n' {
					lineEnd = true
					continue
				}
			}

			switch r.Encoding {
			case cmd.SJIS:
				recordPos = recordPos + SJISCharByteSize(c)
			default:
				recordPos = recordPos + s
			}

			if delimiterPos < recordPos {
				return nil, errors.New("cannot delimit lines at the position of byte array of a character")
			}

			r.buf.WriteRune(c)
		}

		b := r.buf.Bytes()
		b = bytes.TrimSpace(b)

		if len(b) < 1 && !withoutNull {
			record = append(record, nil)
		} else {
			field := make([]byte, len(b))
			copy(field, b)
			record = append(record, field)
		}
	}

	if !lineEnd {
		for {
			c, _, err := r.reader.ReadRune()
			if err != nil {
				if err != io.EOF || recordPos < 1 {
					return nil, err
				}
				break
			}
			switch c {
			case '\r':
				c2, _, _ := r.reader.ReadRune()
				if c2 == '\n' {
					lineBreak = cmd.CRLF
				} else {
					r.reader.UnreadRune()
					lineBreak = cmd.CR
				}
				c = '\n'
			case '\n':
				lineBreak = cmd.LF
			}
			if c == '\n' {
				lineEnd = true
				break
			}
			recordPos++
		}
	}

	if r.LineBreak == "" {
		r.LineBreak = lineBreak
	}

	return record, nil
}
