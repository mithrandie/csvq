package cmd

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strings"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func GetReader(r io.Reader, enc Encoding) io.Reader {
	if enc == SJIS {
		return transform.NewReader(r, japanese.ShiftJIS.NewDecoder())
	}
	return bufio.NewReader(r)
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

	places := []string{}
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
