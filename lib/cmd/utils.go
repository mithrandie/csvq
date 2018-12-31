package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/mithrandie/go-text"
	txjson "github.com/mithrandie/go-text/json"
)

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

func EscapeIdentifier(s string) string {
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
		case '`':
			buf.WriteString("\\`")
		case '\\':
			buf.WriteString("\\\\")
		default:
			buf.WriteRune(r)
		}
	}
	return buf.String()
}

func UnescapeIdentifier(s string) string {
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
			case '`', '\\':
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

func QuoteString(s string) string {
	return "\"" + EscapeString(s) + "\""
}

func QuoteIdentifier(s string) string {
	return "`" + EscapeIdentifier(s) + "`"
}

func VariableSymbol(s string) string {
	return VariableSign + s
}

func FlagSymbol(s string) string {
	return FlagSign + s
}

func EnvironmentVariableSymbol(s string) string {
	if MustBeEnclosed(s) {
		s = QuoteIdentifier(s)
	}
	return EnvironmentVariableSign + s
}

func EnclosedEnvironmentVariableSymbol(s string) string {
	return EnvironmentVariableSign + QuoteIdentifier(s)
}

func MustBeEnclosed(s string) bool {
	runes := []rune(s)

	if runes[0] != '_' && !unicode.IsLetter(runes[0]) {
		return true
	}

	for i := 1; i < len(runes); i++ {
		if s[i] != '_' && !unicode.IsLetter(runes[i]) && !unicode.IsDigit(runes[i]) {
			return true
		}
	}
	return false
}

func RuntimeInformationSymbol(s string) string {
	return RuntimeInformationSign + s
}

func FormatInt(i int, thousandsSeparator string) string {
	return FormatNumber(float64(i), 0, ".", thousandsSeparator, "")
}

func FormatNumber(f float64, precision int, decimalPoint string, thousandsSeparator string, decimalSeparator string) string {
	s := strconv.FormatFloat(f, 'f', precision, 64)

	parts := strings.Split(s, ".")
	intPart := parts[0]
	decPart := ""
	if 1 < len(parts) {
		decPart = parts[1]
	}

	intPlaces := make([]string, 0, (len(intPart)/3)+1)
	intLen := len(intPart)
	for i := intLen / 3; i >= 0; i-- {
		end := intLen - i*3
		if end == 0 {
			continue
		}

		start := intLen - (i+1)*3
		if start < 0 {
			start = 0
		}
		intPlaces = append(intPlaces, intPart[start:end])
	}

	decPlaces := make([]string, 0, (len(decPart)/3)+1)
	for i := 0; i < len(decPart); i = i + 3 {
		end := i + 3
		if len(decPart) < end {
			end = len(decPart)
		}
		decPlaces = append(decPlaces, decPart[i:end])
	}

	formatted := strings.Join(intPlaces, thousandsSeparator)
	if 0 < len(decPlaces) {
		formatted = formatted + decimalPoint + strings.Join(decPlaces, decimalSeparator)
	}

	return formatted
}

func IsReadableFromPipeOrRedirection() bool {
	fi, err := os.Stdin.Stat()
	if err == nil && (fi.Mode()&os.ModeNamedPipe != 0 || 0 < fi.Size()) {
		return true
	}
	return false
}

func ParseEncoding(s string) (text.Encoding, error) {
	var encoding text.Encoding
	switch strings.ToUpper(s) {
	case "UTF8":
		encoding = text.UTF8
	case "SJIS":
		encoding = text.SJIS
	default:
		return text.UTF8, errors.New("encoding must be one of UTF8|SJIS")
	}
	return encoding, nil
}

func ParseLineBreak(s string) (text.LineBreak, error) {
	var lb text.LineBreak
	switch strings.ToUpper(s) {
	case "CRLF":
		lb = text.CRLF
	case "CR":
		lb = text.CR
	case "LF":
		lb = text.LF
	default:
		return lb, errors.New("line-break must be one of CRLF|LF|CR")
	}
	return lb, nil
}

func ParseDelimiter(s string, delimiter rune, delimiterPositions []int, delimitAutomatically bool) (rune, []int, bool, error) {
	s = UnescapeString(s)
	strLen := utf8.RuneCountInString(s)

	if strLen < 1 {
		return delimiter, delimiterPositions, delimitAutomatically, errors.New(fmt.Sprintf("delimiter must be one character, %q or JSON array of integers", DelimiteAutomatically))
	}

	if strLen == 1 {
		delimiter = []rune(s)[0]
		delimitAutomatically = false
		delimiterPositions = nil
	} else {
		if strings.EqualFold(DelimiteAutomatically, s) {
			delimiterPositions = nil
			delimitAutomatically = true
		} else {
			var positions []int
			err := json.Unmarshal([]byte(s), &positions)
			if err != nil {
				return delimiter, delimiterPositions, delimitAutomatically, errors.New(fmt.Sprintf("delimiter must be one character, %q or JSON array of integers", DelimiteAutomatically))
			}
			delimiterPositions = positions
			delimitAutomatically = false
		}
	}
	return delimiter, delimiterPositions, delimitAutomatically, nil
}

func ParseFormat(s string, et txjson.EscapeType) (Format, txjson.EscapeType, error) {
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
	case "LTSV":
		fm = LTSV
	case "GFM":
		fm = GFM
	case "ORG":
		fm = ORG
	case "TEXT":
		fm = TEXT
	case "JSONH":
		fm = JSON
		et = txjson.HexDigits
	case "JSONA":
		fm = JSON
		et = txjson.AllWithHexDigits
	default:
		return fm, et, errors.New("format must be one of CSV|TSV|FIXED|JSON|LTSV|GFM|ORG|TEXT")
	}
	return fm, et, nil
}

func ParseJsonEscapeType(s string) (txjson.EscapeType, error) {
	var escape txjson.EscapeType
	switch strings.ToUpper(s) {
	case "BACKSLASH":
		escape = txjson.Backslash
	case "HEX":
		escape = txjson.HexDigits
	case "HEXALL":
		escape = txjson.AllWithHexDigits
	default:
		return escape, errors.New("json-escape must be one of BACKSLASH|HEX|HEXALL")
	}
	return escape, nil
}

func AppendStrIfNotExist(list []string, elem string) []string {
	if len(elem) < 1 {
		return list
	}
	for _, v := range list {
		if elem == v {
			return list
		}
	}
	return append(list, elem)
}

func TextWidth(s string) int {
	return text.Width(s, GetFlags().EastAsianEncoding, GetFlags().CountDiacriticalSign, GetFlags().CountFormatCode)
}

func RuneWidth(r rune) int {
	return text.RuneWidth(r, GetFlags().EastAsianEncoding, GetFlags().CountDiacriticalSign, GetFlags().CountFormatCode)
}
