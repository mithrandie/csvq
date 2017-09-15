package query

import (
	"bytes"
	"fmt"
	"math"
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

func InStrSlice(s string, list []string) bool {
	for _, v := range list {
		if s == v {
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
	var pad = func(buf *bytes.Buffer, s string, sign []byte, length int, flags []rune) {
		padlen := length - len(sign) - len(s)
		if padlen < 1 {
			buf.Write(sign)
			buf.WriteString(s)
			return
		}

		var padchar byte = ' '
		if InRuneSlice('0', flags) {
			padchar = '0'
		}
		padstr := bytes.Repeat([]byte{padchar}, padlen)
		if InRuneSlice('-', flags) {
			buf.Write(sign)
			buf.WriteString(s)
			buf.Write(padstr)
		} else {
			if padchar == ' ' {
				buf.Write(padstr)
				buf.Write(sign)
			} else {
				buf.Write(sign)
				buf.Write(padstr)
			}
			buf.WriteString(s)
		}
		return
	}

	var numberSign = func(value float64, flags []rune) []byte {
		sign := make([]byte, 0, 1)
		if value < 0 {
			sign = append(sign, '-')
		} else {
			switch {
			case InRuneSlice('+', flags):
				sign = append(sign, '+')
			case InRuneSlice(' ', flags):
				sign = append(sign, ' ')
			}
		}
		return sign
	}

	var buf bytes.Buffer

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
						i := p.(value.Integer).Raw()
						val := float64(i)
						sign := numberSign(val, flags)
						if i < 0 {
							i = i * -1
						}
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
						pad(&buf, s, sign, l, flags)
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
						pad(&buf, s, sign, l, flags)
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
					pad(&buf, s, []byte{}, l, flags)
				case 'q':
					buf.WriteString(args[placeholderOrder].String())
				case 'T':
					buf.WriteString(reflect.TypeOf(args[placeholderOrder]).Name())
				}

				placeholderOrder++
			case '%':
				buf.WriteRune(r)
			default:
				buf.WriteRune('%')
				buf.WriteRune(r)
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

		buf.WriteRune(r)
	}
	if escaped {
		buf.WriteRune('%')
	}

	if placeholderOrder < len(args) {
		return "", NewFormatStringLengthNotMatchError()
	}

	return buf.String(), nil
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
