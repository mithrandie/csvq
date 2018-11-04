package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func GetReader(r io.Reader, enc Encoding) io.Reader {
	if enc == SJIS {
		return transform.NewReader(r, japanese.ShiftJIS.NewDecoder())
	}
	return r
}

func EscapeString(s string) string {
	runes := []rune(s)
	var buf bytes.Buffer

	for _, r := range runes {
		switch r {
		case '\a':
			buf.WriteString("\\a")
		case '\b':
			buf.WriteString("\\b")
		case '\f':
			buf.WriteString("\\f")
		case '\n':
			buf.WriteString("\\n")
		case '\r':
			buf.WriteString("\\r")
		case '\t':
			buf.WriteString("\\t")
		case '\v':
			buf.WriteString("\\v")
		case '"':
			buf.WriteString("\\\"")
		case '\'':
			buf.WriteString("\\'")
		case '\\':
			buf.WriteString("\\\\")
		default:
			buf.WriteRune(r)
		}
	}
	return buf.String()
}

func UnescapeString(s string) string {
	runes := []rune(s)
	var buf bytes.Buffer

	escaped := false
	for _, r := range runes {
		if escaped {
			switch r {
			case 'a':
				buf.WriteRune('\a')
			case 'b':
				buf.WriteRune('\b')
			case 'f':
				buf.WriteRune('\f')
			case 'n':
				buf.WriteRune('\n')
			case 'r':
				buf.WriteRune('\r')
			case 't':
				buf.WriteRune('\t')
			case 'v':
				buf.WriteRune('\v')
			case '"', '\'', '\\':
				buf.WriteRune(r)
			default:
				buf.WriteRune('\\')
				buf.WriteRune(r)
			}
			escaped = false
			continue
		}

		if r == '\\' {
			escaped = true
			continue
		}

		buf.WriteRune(r)
	}
	if escaped {
		buf.WriteRune('\\')
	}

	return buf.String()
}

func HumarizeNumber(s string) string {
	parts := strings.Split(s, ".")
	intPart := parts[0]
	decPart := ""
	if 1 < len(parts) {
		decPart = "." + parts[1]
	}

	places := make([]string, 0)
	slen := len(intPart)
	for i := slen / 3; i >= 0; i-- {
		end := slen - i*3
		if end == 0 {
			continue
		}

		start := slen - (i+1)*3
		if start < 0 {
			start = 0
		}
		places = append(places, intPart[start:end])
	}

	return strings.Join(places, ",") + decPart
}

func IsReadableFromPipeOrRedirection() bool {
	fi, err := os.Stdin.Stat()
	if err == nil && (fi.Mode()&os.ModeNamedPipe != 0 || 0 < fi.Size()) {
		return true
	}
	return false
}

func ParseEncoding(s string) (Encoding, error) {
	var encoding Encoding
	switch strings.ToUpper(s) {
	case "UTF8":
		encoding = UTF8
	case "SJIS":
		encoding = SJIS
	default:
		return UTF8, errors.New("encoding must be one of UTF8|SJIS")
	}
	return encoding, nil
}

func ParseLineBreak(s string) (LineBreak, error) {
	var lb LineBreak
	switch strings.ToUpper(s) {
	case "CRLF":
		lb = CRLF
	case "CR":
		lb = CR
	case "LF":
		lb = LF
	default:
		return lb, errors.New("line-break must be one of CRLF|LF|CR")
	}
	return lb, nil
}

func ParseDelimiter(s string, delimiter rune, delimiterPositions []int, delimitAutomatically bool) (rune, []int, bool, error) {
	s = UnescapeString(s)
	strLen := utf8.RuneCountInString(s)

	if strLen < 1 {
		return delimiter, delimiterPositions, delimitAutomatically, errors.New("delimiter must be one character, \"SPACES\" or JSON array of integers")
	}

	if strLen == 1 {
		delimiter = []rune(s)[0]
		delimitAutomatically = false
		delimiterPositions = nil
	} else {
		if strings.EqualFold("SPACES", s) {
			delimiterPositions = nil
			delimitAutomatically = true
		} else {
			var positions []int
			err := json.Unmarshal([]byte(s), &positions)
			if err != nil {
				return delimiter, delimiterPositions, delimitAutomatically, errors.New("delimiter must be one character, \"SPACES\" or JSON array of integers")
			}
			delimiterPositions = positions
			delimitAutomatically = false
		}
	}
	return delimiter, delimiterPositions, delimitAutomatically, nil
}

func ParseFormat(s string) (Format, error) {
	var fm Format
	switch strings.ToUpper(s) {
	case "CSV":
		fm = CSV
	case "TSV":
		fm = TSV
	case "FIXED":
		fm = FIXED
	case "JSON":
		fm = JSON
	case "JSONH":
		fm = JSONH
	case "JSONA":
		fm = JSONA
	case "GFM":
		fm = GFM
	case "ORG":
		fm = ORG
	case "TEXT":
		fm = TEXT
	default:
		return fm, errors.New("format must be one of CSV|TSV|FIXED|JSON|JSONH|JSONA|GFM|ORG|TEXT")
	}
	return fm, nil
}
