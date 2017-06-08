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
	s = strings.Replace(s, "\\a", "\a", -1)
	s = strings.Replace(s, "\\b", "\b", -1)
	s = strings.Replace(s, "\\f", "\f", -1)
	s = strings.Replace(s, "\\n", "\n", -1)
	s = strings.Replace(s, "\\r", "\r", -1)
	s = strings.Replace(s, "\\t", "\t", -1)
	s = strings.Replace(s, "\\v", "\v", -1)
	s = strings.Replace(s, "\\\"", "\"", -1)
	s = strings.Replace(s, "\\\\", "\\", -1)
	return s
}
