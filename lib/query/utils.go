package query

import (
	"fmt"
	"os"
	"strings"

	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/ternary"
)

func InIntSlice(i int, list []int) bool {
	for _, v := range list {
		if i == v {
			return true
		}
	}
	return false
}

func InStrSliceWithCaseInsensitive(s string, list []string) bool {
	for _, v := range list {
		if strings.EqualFold(s, v) {
			return true
		}
	}
	return false
}

func Distinguish(list []parser.Primary) []parser.Primary {
	var in = func(list []parser.Primary, item parser.Primary) bool {
		for _, v := range list {
			if EquivalentTo(item, v) == ternary.TRUE {
				return true
			}
		}
		return false
	}

	distinguished := []parser.Primary{}
	for _, v := range list {
		if !in(distinguished, v) {
			distinguished = append(distinguished, v)
		}
	}
	return distinguished
}

func FormatCount(i int, obj string) string {
	var s string
	if i == 0 {
		s = fmt.Sprintf("no %s", obj)
	} else if i == 1 {
		s = fmt.Sprintf("%d %s", i, obj)
	} else {
		s = fmt.Sprintf("%d %ss", i, obj)
	}
	return s
}

func IsReadableFromStdin() bool {
	fi, err := os.Stdin.Stat()
	if err == nil && (fi.Mode()&os.ModeNamedPipe != 0 || 0 < fi.Size()) {
		return true
	}
	return false
}

func FormatString(format string, args []parser.Primary) (string, error) {
	str := []rune{}

	escaped := false
	placeholderOrder := 0
	for _, r := range format {
		if escaped {
			switch r {
			case 's':
				if len(args) <= placeholderOrder {
					return "", NewFormatStringLengthNotMatchError()
				}
				str = append(str, []rune(args[placeholderOrder].String())...)
				placeholderOrder++
			case '%':
				str = append(str, r)
			default:
				str = append(str, '%', r)
			}
			escaped = false
			continue
		}

		if r == '%' {
			escaped = true
			continue
		}

		str = append(str, r)
	}
	if escaped {
		str = append(str, '%')
	}

	if placeholderOrder < len(args) {
		return "", NewFormatStringLengthNotMatchError()
	}

	return string(str), nil
}
