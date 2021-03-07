package value

import (
	"bytes"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/mithrandie/csvq/lib/cmd"

	"github.com/mithrandie/ternary"
)

var DatetimeFormats = NewDatetimeFormatMap()

type DatetimeFormatMap struct {
	m *sync.Map
}

func NewDatetimeFormatMap() *DatetimeFormatMap {
	return &DatetimeFormatMap{
		m: &sync.Map{},
	}
}

func (dfmap *DatetimeFormatMap) Store(key string, value string) {
	dfmap.m.Store(key, value)
}

func (dfmap *DatetimeFormatMap) Load(key string) (string, bool) {
	v, ok := dfmap.m.Load(key)
	if ok {
		return v.(string), ok
	}
	return "", ok
}

func (dfmap DatetimeFormatMap) Get(s string) string {
	if f, ok := dfmap.Load(s); ok {
		return f
	}
	f := ConvertDatetimeFormat(s)
	dfmap.Store(s, f)
	return f
}

func StrToTime(s string, formats []string, location *time.Location) (time.Time, bool) {
	s = cmd.TrimSpace(s)

	for _, format := range formats {
		if t, e := time.ParseInLocation(DatetimeFormats.Get(format), s, location); e == nil {
			return t, true
		}
	}

	if 8 <= len(s) && '0' <= s[0] && s[0] <= '9' {
		switch {
		case s[4] == '-':
			if len(s) < 10 {
				if t, e := time.ParseInLocation("2006-1-2", s, location); e == nil {
					return t, true
				}
			} else if len(s) == 10 {
				if t, e := time.ParseInLocation("2006-01-02", s, location); e == nil {
					return t, true
				}
			} else if s[10] == 'T' {
				if s[len(s)-6] == '+' || s[len(s)-6] == '-' || s[len(s)-1] == 'Z' {
					if t, e := time.Parse(time.RFC3339Nano, s); e == nil {
						return t, true
					}
				} else {
					if t, e := time.ParseInLocation("2006-01-02T15:04:05.999999999", s, location); e == nil {
						return t, true
					}
				}
			} else if s[10] == ' ' {
				if t, e := time.ParseInLocation("2006-01-02 15:04:05.999999999", s, location); e == nil {
					return t, true
				} else if t, e := time.Parse("2006-01-02 15:04:05.999999999 Z07:00", s); e == nil {
					return t, true
				} else if t, e := time.Parse("2006-01-02 15:04:05.999999999 -0700", s); e == nil {
					return t, true
				} else if t, e := time.Parse("2006-01-02 15:04:05.999999999 MST", s); e == nil {
					return t, true
				}
			} else {
				if t, e := time.ParseInLocation("2006-1-2 15:04:05.999999999", s, location); e == nil {
					return t, true
				} else if t, e := time.Parse("2006-1-2 15:04:05.999999999 Z07:00", s); e == nil {
					return t, true
				} else if t, e := time.Parse("2006-1-2 15:04:05.999999999 -0700", s); e == nil {
					return t, true
				} else if t, e := time.Parse("2006-1-2 15:04:05.999999999 MST", s); e == nil {
					return t, true
				}
			}
		case s[4] == '/':
			if len(s) < 10 {
				if t, e := time.ParseInLocation("2006/1/2", s, location); e == nil {
					return t, true
				}
			} else if len(s) == 10 {
				if t, e := time.ParseInLocation("2006/01/02", s, location); e == nil {
					return t, true
				}
			} else if s[10] == ' ' {
				if t, e := time.ParseInLocation("2006/01/02 15:04:05.999999999", s, location); e == nil {
					return t, true
				} else if t, e := time.Parse("2006/01/02 15:04:05.999999999 Z07:00", s); e == nil {
					return t, true
				} else if t, e := time.Parse("2006/01/02 15:04:05.999999999 -0700", s); e == nil {
					return t, true
				} else if t, e := time.Parse("2006/01/02 15:04:05.999999999 MST", s); e == nil {
					return t, true
				}
			} else {
				if t, e := time.ParseInLocation("2006/1/2 15:04:05.999999999", s, location); e == nil {
					return t, true
				} else if t, e := time.Parse("2006/1/2 15:04:05.999999999 Z07:00", s); e == nil {
					return t, true
				} else if t, e := time.Parse("2006/1/2 15:04:05.999999999 -0700", s); e == nil {
					return t, true
				} else if t, e := time.Parse("2006/1/2 15:04:05.999999999 MST", s); e == nil {
					return t, true
				}
			}
		default:
			if t, e := time.Parse(time.RFC822, s); e == nil {
				return t, true
			} else if t, e := time.Parse(time.RFC822Z, s); e == nil {
				return t, true
			}
		}
	}
	return time.Time{}, false
}

func ConvertDatetimeFormat(format string) string {
	runes := []rune(format)
	var buf bytes.Buffer

	escaped := false
	for _, r := range runes {
		if !escaped {
			switch r {
			case '%':
				escaped = true
			default:
				buf.WriteRune(r)
			}
			continue
		}

		switch r {
		case 'a':
			buf.WriteString("Mon")
		case 'b':
			buf.WriteString("Jan")
		case 'c':
			buf.WriteString("1")
		case 'd':
			buf.WriteString("02")
		case 'E':
			buf.WriteString("_2")
		case 'e':
			buf.WriteString("2")
		case 'F':
			buf.WriteString(".999999")
		case 'f':
			buf.WriteString(".000000")
		case 'H':
			buf.WriteString("15")
		case 'h':
			buf.WriteString("03")
		case 'i':
			buf.WriteString("04")
		case 'l':
			buf.WriteString("3")
		case 'M':
			buf.WriteString("January")
		case 'm':
			buf.WriteString("01")
		case 'N':
			buf.WriteString(".999999999")
		case 'n':
			buf.WriteString(".000000000")
		case 'p':
			buf.WriteString("PM")
		case 'r':
			buf.WriteString("03:04:05 PM")
		case 's':
			buf.WriteString("05")
		case 'T':
			buf.WriteString("15:04:05")
		case 'W':
			buf.WriteString("Monday")
		case 'Y':
			buf.WriteString("2006")
		case 'y':
			buf.WriteString("06")
		case 'Z':
			buf.WriteString("Z07:00")
		case 'z':
			buf.WriteString("MST")
		default:
			buf.WriteRune(r)
		}
		escaped = false
	}

	return buf.String()
}

func Float64ToTime(f float64, location *time.Location) time.Time {
	s := Float64ToStr(f)
	pointIdx := strings.Index(s, ".")
	if -1 < pointIdx {
		afterPLen := len(s) - 1 - pointIdx
		if 9 < afterPLen {
			s = s[:pointIdx+10]
			afterPLen = 9
		}
		s = s[:pointIdx] + s[pointIdx+1:] + strings.Repeat("0", 9-afterPLen)
	} else {
		s = s + strings.Repeat("0", 9)
	}
	nsec, _ := strconv.ParseInt(s, 10, 64)
	return TimeFromUnixTime(0, nsec, location)
}

func Int64ToStr(i int64) string {
	return strconv.FormatInt(i, 10)
}

func Float64ToStr(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func ParseFloat64(f float64) Primary {
	if math.Remainder(f, 1) == 0 {
		return NewInteger(int64(f))
	}
	return NewFloat(f)
}

func ToInteger(p Primary) Primary {
	switch p.(type) {
	case *Integer:
		return NewInteger(p.(*Integer).Raw())
	case *Float:
		f := p.(*Float).Raw()
		if math.Remainder(f, 1) == 0 {
			return NewInteger(int64(f))
		}
	case *String:
		s := cmd.TrimSpace(p.(*String).Raw())
		if MaybeInteger(s) {
			if i, e := strconv.ParseInt(s, 10, 64); e == nil {
				return NewInteger(i)
			}
		}
		if MaybeNumber(s) {
			if f, e := strconv.ParseFloat(s, 64); e == nil {
				if math.Remainder(f, 1) == 0 {
					return NewInteger(int64(f))
				}
			}
		}
	}

	return NewNull()
}

func ToFloat(p Primary) Primary {
	switch p.(type) {
	case *Integer:
		return NewFloat(float64(p.(*Integer).Raw()))
	case *Float:
		return NewFloat(p.(*Float).Raw())
	case *String:
		s := cmd.TrimSpace(p.(*String).Raw())
		if MaybeNumber(s) {
			if f, e := strconv.ParseFloat(p.(*String).Raw(), 64); e == nil {
				return NewFloat(f)
			}
		}
	}

	return NewNull()
}

func MaybeInteger(s string) bool {
	if len(s) < 1 {
		return false
	}

	start := 0
	if s[start] == '+' || s[start] == '-' {
		start++
	}

	for i := start; i < len(s); i++ {
		if !isDecimal(s[i]) {
			return false
		}
	}
	return true
}

func MaybeNumber(s string) bool {
	if len(s) < 1 {
		return false
	}

	start := 0
	if s[start] == '+' || s[start] == '-' {
		start++
	}
	if len(s) < start+1 {
		return false
	}

	pointExists := false
	eExists := false
	for i := start; i < len(s); i++ {
		if !pointExists && s[i] == '.' {
			if len(s) < i+2 {
				return false
			}
			pointExists = true
			continue
		}

		if !eExists && (s[i] == 'e' || s[i] == 'E') {
			i++
			if len(s) < i+2 {
				return false
			}
			if s[i] != '+' && s[i] != '-' {
				return false
			}
			eExists = true
			continue
		}

		if isDecimal(s[i]) {
			continue
		}

		return false
	}

	return true
}

func isDecimal(b byte) bool {
	return '0' <= b && b <= '9'
}

func ToDatetime(p Primary, formats []string, location *time.Location) Primary {
	switch p.(type) {
	case *Datetime:
		return NewDatetime(p.(*Datetime).Raw())
	case *String:
		if dt, ok := StrToTime(p.(*String).Raw(), formats, location); ok {
			return NewDatetime(dt)
		}
	}

	return NewNull()
}

func TimeFromUnixTime(sec int64, nano int64, location *time.Location) time.Time {
	return time.Unix(sec, nano).In(location)
}

func ToBoolean(p Primary) Primary {
	switch p.(type) {
	case *Boolean:
		return NewBoolean(p.(*Boolean).Raw())
	case *String, *Integer, *Float, *Ternary:
		if p.Ternary() != ternary.UNKNOWN {
			return NewBoolean(p.Ternary().ParseBool())
		}
	}
	return NewNull()
}

func ToString(p Primary) Primary {
	switch p.(type) {
	case *String:
		return NewString(p.(*String).Raw())
	case *Integer:
		return NewString(Int64ToStr(p.(*Integer).Raw()))
	case *Float:
		return NewString(Float64ToStr(p.(*Float).Raw()))
	}
	return NewNull()
}
