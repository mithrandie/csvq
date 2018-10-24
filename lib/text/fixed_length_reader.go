package text

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/value"
	"io"
	"strconv"
	"strings"
	"unicode"
)

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

type FixedLengthReader struct {
	DelimiterPositions []int
	WithoutNull        bool
	Encoding           cmd.Encoding

	reader *bufio.Reader

	lineBuf   bytes.Buffer
	fieldBuf  bytes.Buffer
	spacesBuf bytes.Buffer
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

func (r *FixedLengthReader) Read() ([]Field, error) {
	return r.parseRecord(r.WithoutNull)
}

func (r *FixedLengthReader) ReadAll() ([][]Field, error) {
	records := make([][]Field, 0, 100)

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

func (r *FixedLengthReader) parseRecord(withoutNull bool) ([]Field, error) {
	r.lineBuf.Reset()

	for {
		line, isPrefix, err := r.reader.ReadLine()
		if err != nil {
			return nil, err
		}

		r.lineBuf.Write(line)
		if !isPrefix {
			break
		}
	}

	record := make([]Field, 0, len(r.DelimiterPositions)+1)
	recordPos := 0
	delimiterPos := 0

	for _, endPos := range r.DelimiterPositions {
		if endPos < 0 || endPos <= delimiterPos {
			return nil, errors.New(fmt.Sprintf("invalid delimiter position: %s", FormatIntSlice(r.DelimiterPositions)))
		}
		delimiterPos = endPos

		r.fieldBuf.Reset()
		r.spacesBuf.Reset()
		var trimChar bool

		for recordPos < delimiterPos {
			c, s, err := r.lineBuf.ReadRune()
			if err != nil {
				if err == io.EOF {
					break
				}
				return nil, err
			}

			trimChar = false
			if unicode.IsSpace(c) {
				trimChar = true
				if 0 < r.fieldBuf.Len() {
					r.spacesBuf.WriteRune(c)
				}
			}

			if !trimChar {
				if 0 < r.spacesBuf.Len() {
					r.fieldBuf.Write(r.spacesBuf.Bytes())
					r.spacesBuf.Reset()
				}
				r.fieldBuf.WriteRune(c)
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
		}

		if r.fieldBuf.Len() < 1 && !withoutNull {
			record = append(record, nil)
		} else {
			field := make([]byte, r.fieldBuf.Len())
			copy(field, r.fieldBuf.Bytes())
			record = append(record, field)
		}
	}

	return record, nil
}

func FormatIntSlice(list []int) string {
	slist := make([]string, 0, len(list))
	for _, v := range list {
		slist = append(slist, strconv.Itoa(v))
	}
	return "[" + strings.Join(slist, ", ") + "]"
}
