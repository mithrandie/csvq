package query

import (
	"fmt"
	"math"
	"os"
	"reflect"
	"strconv"
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

func InRuneSlice(r rune, list []rune) bool {
	for _, v := range list {
		if r == v {
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
	var pad = func(s string, length int, flags []rune) string {
		if length <= len(s) {
			return s
		}

		padchar := " "
		if InRuneSlice('0', flags) {
			padchar = "0"
		}
		padstr := strings.Repeat(padchar, length-len(s))
		if InRuneSlice('-', flags) {
			s = s + padstr
		} else {
			s = padstr + s
		}
		return s
	}

	var numberSign = func(value float64, flags []rune) string {
		sign := ""
		if value < 0 {
			sign = "-"
		} else {
			switch {
			case InRuneSlice('+', flags):
				sign = "+"
			case InRuneSlice(' ', flags):
				sign = " "
			}
		}
		return sign
	}

	str := []rune{}

	escaped := false
	placeholderOrder := 0
	flags := []rune{}
	var length string
	var precision string
	var isPrecision bool
	for _, r := range format {
		if escaped {
			if isPrecision {
				if '0' <= r && r <= '9' {
					precision += string(r)
					continue
				} else {
					isPrecision = false
				}
			}

			if 0 < len(length) && '0' <= r && r <= '9' {
				length += string(r)
				continue
			}

			switch r {
			case '+', '-', ' ', '0':
				flags = append(flags, r)
				continue
			case '1', '2', '3', '4', '5', '6', '7', '8', '9':
				length = string(r)
				continue
			case '.':
				isPrecision = true
				continue
			case 'b', 'o', 'd', 'x', 'X', 'e', 'E', 'f', 's', 'q', 'T':
				if len(args) <= placeholderOrder {
					return "", NewFormatStringLengthNotMatchError()
				}

				switch r {
				case 'b', 'o', 'd', 'x', 'X':
					p := parser.PrimaryToInteger(args[placeholderOrder])
					if !parser.IsNull(p) {
						value := float64(p.(parser.Integer).Value())
						sign := numberSign(value, flags)
						i := int64(math.Abs(value))
						var s string
						switch r {
						case 'b':
							s = strconv.FormatInt(i, 2)
						case 'o':
							s = strconv.FormatInt(i, 8)
						case 'd':
							s = strconv.FormatInt(i, 10)
						case 'x':
							s = strconv.FormatInt(i, 16)
						case 'X':
							s = strings.ToUpper(strconv.FormatInt(i, 16))
						}
						l, _ := strconv.Atoi(length)
						s = sign + pad(s, l-len(sign), flags)
						str = append(str, []rune(s)...)
					}
				case 'e', 'E', 'f':
					p := parser.PrimaryToFloat(args[placeholderOrder])
					if !parser.IsNull(p) {
						value := p.(parser.Float).Value()

						var prec float64
						if 0 < len(precision) {
							prec, _ = strconv.ParseFloat(precision, 64)
							value = round(value, prec)
						}
						sign := numberSign(value, flags)
						f := math.Abs(value)
						s := strconv.FormatFloat(f, byte(r), -1, 64)

						if 0 < prec {
							parts := strings.Split(s, ".")
							intpart := parts[0]
							var dec string
							var en string
							if len(parts) < 2 {
								dec = ""
							} else {
								dec = parts[1]
							}
							if r != 'f' {
								if 0 < len(dec) {
									enidx := strings.Index(dec, string(r))
									en = dec[enidx:]
									dec = dec[:enidx]
								} else {
									enidx := strings.Index(intpart, string(r))
									en = intpart[enidx:]
									intpart = intpart[:enidx]
								}
							}
							if len(dec) < int(prec) {
								dec = dec + strings.Repeat("0", int(prec)-len(dec))
								s = intpart + "." + dec + en
							}
						}

						l, _ := strconv.Atoi(length)
						s = sign + pad(s, l-len(sign), flags)
						str = append(str, []rune(s)...)
					}
				case 's':
					var s string
					switch args[placeholderOrder].(type) {
					case parser.String:
						s = args[placeholderOrder].(parser.String).Value()
					case parser.Integer:
						s = parser.Int64ToStr(args[placeholderOrder].(parser.Integer).Value())
					case parser.Float:
						s = parser.Float64ToStr(args[placeholderOrder].(parser.Float).Value())
					case parser.Boolean:
						s = strconv.FormatBool(args[placeholderOrder].(parser.Boolean).Value())
					case parser.Ternary:
						s = args[placeholderOrder].(parser.Ternary).Ternary().String()
					case parser.Datetime:
						s = args[placeholderOrder].(parser.Datetime).Format()
					case parser.Null:
						s = "NULL"
					}
					l, _ := strconv.Atoi(length)
					s = pad(s, l, flags)
					str = append(str, []rune(s)...)
				case 'q':
					str = append(str, []rune(args[placeholderOrder].String())...)
				case 'T':
					str = append(str, []rune(reflect.TypeOf(args[placeholderOrder]).Name())...)
				}

				placeholderOrder++
			case '%':
				str = append(str, r)
			default:
				str = append(str, '%', r)
			}

			escaped = false
			flags = []rune{}
			length = ""
			precision = ""
			isPrecision = false
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
