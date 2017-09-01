package query

import (
	"fmt"
	"math"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/value"
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

func Distinguish(list []value.Primary) []value.Primary {
	values := make(map[string]int)
	valueKeys := make([]string, 0, len(list))

	for i, v := range list {
		key := SerializeComparisonKeys([]value.Primary{v})
		if _, ok := values[key]; !ok {
			values[key] = i
			valueKeys = append(valueKeys, key)
		}
	}

	distinguished := make([]value.Primary, len(valueKeys))
	for i, key := range valueKeys {
		distinguished[i] = list[values[key]]
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

func SerializeComparisonKeys(values []value.Primary) string {
	list := make([]string, len(values))

	for i, val := range values {
		list[i] = SerializeKey(val)
	}

	return strings.Join(list, ":")
}

func SerializeKey(val value.Primary) string {
	if value.IsNull(val) {
		return serializeNull()
	} else if in := value.ToInteger(val); !value.IsNull(in) {
		return serializeInteger(in.(value.Integer).Raw())
	} else if f := value.ToFloat(val); !value.IsNull(f) {
		return serializeFlaot(f.(value.Float).Raw())
	} else if dt := value.ToDatetime(val); !value.IsNull(dt) {
		t := dt.(value.Datetime).Raw()
		if t.Nanosecond() > 0 {
			f := float64(t.Unix()) + float64(t.Nanosecond())/1e9
			t2 := value.Float64ToTime(f)
			if t.Equal(t2) {
				return serializeFlaot(f)
			} else {
				return serializeDatetime(t)
			}
		} else {
			return serializeInteger(t.Unix())
		}
	} else if b := value.ToBoolean(val); !value.IsNull(b) {
		return serializeBoolean(b.(value.Boolean).Raw())
	} else if s, ok := val.(value.String); ok {
		return serializeString(s.Raw())
	} else {
		return serializeNull()
	}
}

func serializeNull() string {
	return "[N]"
}

func serializeInteger(i int64) string {
	var b string
	switch i {
	case 0:
		b = "[B]" + strconv.FormatBool(false)
	case 1:
		b = "[B]" + strconv.FormatBool(true)
	}
	return "[I]" + value.Int64ToStr(i) + b
}

func serializeFlaot(f float64) string {
	return "[F]" + value.Float64ToStr(f)
}

func serializeDatetime(t time.Time) string {
	return "[D]" + value.Int64ToStr(t.UnixNano())
}

func serializeDatetimeFromUnixNano(t int64) string {
	return "[D]" + value.Int64ToStr(t)
}

func serializeBoolean(b bool) string {
	var intliteral string
	if b {
		intliteral = "1"
	} else {
		intliteral = "0"
	}
	return "[I]" + intliteral + "[B]" + strconv.FormatBool(b)
}

func serializeString(s string) string {
	return "[S]" + strings.ToUpper(strings.TrimSpace(s))
}

func FormatString(format string, args []value.Primary) (string, error) {
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
					p := value.ToInteger(args[placeholderOrder])
					if !value.IsNull(p) {
						val := float64(p.(value.Integer).Raw())
						sign := numberSign(val, flags)
						i := int64(math.Abs(val))
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
					p := value.ToFloat(args[placeholderOrder])
					if !value.IsNull(p) {
						val := p.(value.Float).Raw()

						var prec float64
						if 0 < len(precision) {
							prec, _ = strconv.ParseFloat(precision, 64)
							val = round(val, prec)
						}
						sign := numberSign(val, flags)
						f := math.Abs(val)
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
					case value.String:
						s = args[placeholderOrder].(value.String).Raw()
					case value.Integer:
						s = args[placeholderOrder].(value.Integer).String()
					case value.Float:
						s = args[placeholderOrder].(value.Float).String()
					case value.Boolean:
						s = args[placeholderOrder].(value.Boolean).String()
					case value.Ternary:
						s = args[placeholderOrder].(value.Ternary).Ternary().String()
					case value.Datetime:
						s = args[placeholderOrder].(value.Datetime).Format(time.RFC3339Nano)
					case value.Null:
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

func IsReadableFromStdin() bool {
	fi, err := os.Stdin.Stat()
	if err == nil && (fi.Mode()&os.ModeNamedPipe != 0 || 0 < fi.Size()) {
		return true
	}
	return false
}

func NumberOfCPU(recordLen int) int {
	num := cmd.GetFlags().CPU
	if 2 < num {
		num = num - 1
	}
	if num <= runtime.NumGoroutine() || recordLen < 150 || recordLen < num {
		num = 1
	}

	return num
}

func RecordRange(cpuIndex int, totalLen int, numberOfCPU int) (int, int) {
	calcLen := totalLen / numberOfCPU

	var start int = cpuIndex * calcLen
	var end int
	if cpuIndex == numberOfCPU-1 {
		end = totalLen
	} else {
		end = (cpuIndex + 1) * calcLen
	}
	return start, end
}
