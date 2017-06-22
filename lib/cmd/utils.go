package cmd

import (
	"bufio"
	"io"

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
	unescaped := []rune{}

	escaped := false
	for _, r := range runes {
		if escaped {
			switch r {
			case 'a':
				unescaped = append(unescaped, '\a')
			case 'b':
				unescaped = append(unescaped, '\b')
			case 'f':
				unescaped = append(unescaped, '\f')
			case 'n':
				unescaped = append(unescaped, '\n')
			case 'r':
				unescaped = append(unescaped, '\r')
			case 't':
				unescaped = append(unescaped, '\t')
			case 'v':
				unescaped = append(unescaped, '\v')
			case '"':
				unescaped = append(unescaped, '"')
			case '\\':
				unescaped = append(unescaped, '\\')
			default:
				unescaped = append(unescaped, '\\', r)
			}
			escaped = false
			continue
		}

		if r == '\\' {
			escaped = true
			continue
		}

		unescaped = append(unescaped, r)
	}
	if escaped {
		unescaped = append(unescaped, '\\')
	}

	return string(unescaped)
}
