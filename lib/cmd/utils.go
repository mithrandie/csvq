package cmd

import (
	"bufio"
	"io"
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
