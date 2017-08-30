package parser

import (
	"errors"
	"math"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/ternary"
)

type ConvertedDatetimeFormatMap map[string]string

func (m ConvertedDatetimeFormatMap) Get(s string) string {
	if f, ok := m[s]; ok {
		return f
	}
	f := ConvertDatetimeFormat(s)
	m[s] = f
	return f
}

var DatetimeFormats = ConvertedDatetimeFormatMap{}

func StrToTime(s string) (time.Time, error) {
	s = strings.TrimSpace(s)

	flags := cmd.GetFlags()
	if 0 < len(flags.DatetimeFormat) {
		if t, e := time.ParseInLocation(DatetimeFormats.Get(flags.DatetimeFormat), s, cmd.GetLocation()); e == nil {
			return t, nil
		}
	}

	if 8 <= len(s) && '0' <= s[0] && s[0] <= '9' {
		switch {
		case s[4] == '-':
			if len(s) < 10 {
				if t, e := time.ParseInLocation("2006-1-2", s, cmd.GetLocation()); e == nil {
					return t, nil
				}
			} else if len(s) == 10 {
				if t, e := time.ParseInLocation("2006-01-02", s, cmd.GetLocation()); e == nil {
					return t, nil
				}
			} else if s[10] == 'T' {
				if s[len(s)-3] == '+' || s[len(s)-6] == '-' || s[len(s)-1] == 'Z' {
					if t, e := time.Parse(time.RFC3339Nano, s); e == nil {
						return t, nil
					}
				} else {
					if t, e := time.ParseInLocation("2006-01-02T15:04:05.999999999", s, cmd.GetLocation()); e == nil {
						return t, nil
					}
				}
			} else if s[10] == ' ' {
				if t, e := time.ParseInLocation("2006-01-02 15:04:05.999999999", s, cmd.GetLocation()); e == nil {
					return t, nil
				} else if t, e := time.Parse("2006-01-02 15:04:05.999999999 -07:00", s); e == nil {
					return t, nil
				} else if t, e := time.Parse("2006-01-02 15:04:05.999999999 -0700", s); e == nil {
					return t, nil
				} else if t, e := time.Parse("2006-01-02 15:04:05.999999999 MST", s); e == nil {
					return t, nil
				}
			} else {
				if t, e := time.ParseInLocation("2006-1-2 15:04:05.999999999", s, cmd.GetLocation()); e == nil {
					return t, nil
				} else if t, e := time.Parse("2006-1-2 15:04:05.999999999 -07:00", s); e == nil {
					return t, nil
				} else if t, e := time.Parse("2006-1-2 15:04:05.999999999 -0700", s); e == nil {
					return t, nil
				} else if t, e := time.Parse("2006-1-2 15:04:05.999999999 MST", s); e == nil {
					return t, nil
				}
			}
		case s[4] == '/':
			if len(s) < 10 {
				if t, e := time.ParseInLocation("2006/1/2", s, cmd.GetLocation()); e == nil {
					return t, nil
				}
			} else if len(s) == 10 {
				if t, e := time.ParseInLocation("2006/01/02", s, cmd.GetLocation()); e == nil {
					return t, nil
				}
			} else if s[10] == ' ' {
				if t, e := time.ParseInLocation("2006/01/02 15:04:05.999999999", s, cmd.GetLocation()); e == nil {
					return t, nil
				} else if t, e := time.Parse("2006/01/02 15:04:05.999999999 Z07:00", s); e == nil {
					return t, nil
				} else if t, e := time.Parse("2006/01/02 15:04:05.999999999 -0700", s); e == nil {
					return t, nil
				} else if t, e := time.Parse("2006/01/02 15:04:05.999999999 MST", s); e == nil {
					return t, nil
				}
			} else {
				if t, e := time.ParseInLocation("2006/1/2 15:04:05.999999999", s, cmd.GetLocation()); e == nil {
					return t, nil
				} else if t, e := time.Parse("2006/1/2 15:04:05.999999999 Z07:00", s); e == nil {
					return t, nil
				} else if t, e := time.Parse("2006/1/2 15:04:05.999999999 -0700", s); e == nil {
					return t, nil
				} else if t, e := time.Parse("2006/1/2 15:04:05.999999999 MST", s); e == nil {
					return t, nil
				}
			}
		default:
			if t, e := time.Parse(time.RFC822, s); e == nil {
				return t, nil
			} else if t, e := time.Parse(time.RFC822Z, s); e == nil {
				return t, nil
			}
		}
	}
	return time.Time{}, errors.New("conversion failed")
}

func ConvertDatetimeFormat(format string) string {
	runes := []rune(format)
	dtfmt := []rune{}

	escaped := false
	for _, r := range runes {
		if !escaped {
			switch r {
			case '%':
				escaped = true
			default:
				dtfmt = append(dtfmt, r)
			}
			continue
		}

		switch r {
		case 'a':
			dtfmt = append(dtfmt, []rune("Mon")...)
		case 'b':
			dtfmt = append(dtfmt, []rune("Jan")...)
		case 'c':
			dtfmt = append(dtfmt, []rune("1")...)
		case 'd':
			dtfmt = append(dtfmt, []rune("02")...)
		case 'E':
			dtfmt = append(dtfmt, []rune("_2")...)
		case 'e':
			dtfmt = append(dtfmt, []rune("2")...)
		case 'F':
			dtfmt = append(dtfmt, []rune(".999999")...)
		case 'f':
			dtfmt = append(dtfmt, []rune(".000000")...)
		case 'H':
			dtfmt = append(dtfmt, []rune("15")...)
		case 'h':
			dtfmt = append(dtfmt, []rune("03")...)
		case 'i':
			dtfmt = append(dtfmt, []rune("04")...)
		case 'l':
			dtfmt = append(dtfmt, []rune("3")...)
		case 'M':
			dtfmt = append(dtfmt, []rune("January")...)
		case 'm':
			dtfmt = append(dtfmt, []rune("01")...)
		case 'N':
			dtfmt = append(dtfmt, []rune(".999999999")...)
		case 'n':
			dtfmt = append(dtfmt, []rune(".000000000")...)
		case 'p':
			dtfmt = append(dtfmt, []rune("PM")...)
		case 'r':
			dtfmt = append(dtfmt, []rune("03:04:05 PM")...)
		case 's':
			dtfmt = append(dtfmt, []rune("05")...)
		case 'T':
			dtfmt = append(dtfmt, []rune("15:04:05")...)
		case 'W':
			dtfmt = append(dtfmt, []rune("Monday")...)
		case 'Y':
			dtfmt = append(dtfmt, []rune("2006")...)
		case 'y':
			dtfmt = append(dtfmt, []rune("06")...)
		case 'Z':
			dtfmt = append(dtfmt, []rune("Z07:00")...)
		case 'z':
			dtfmt = append(dtfmt, []rune("MST")...)
		default:
			dtfmt = append(dtfmt, r)
		}
		escaped = false
	}

	return string(dtfmt)
}

func Int64ToStr(i int64) string {
	return strconv.FormatInt(i, 10)
}

func Float64ToStr(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func Float64ToPrimary(f float64) Primary {
	if math.Remainder(f, 1) == 0 {
		return NewInteger(int64(f))
	}
	return NewFloat(f)
}

func PrimaryToInteger(p Primary) Primary {
	switch p.(type) {
	case Integer:
		return p
	case Float:
		f := p.(Float).Value()
		if math.Remainder(f, 1) == 0 {
			return NewInteger(int64(f))
		}
	case String:
		s := strings.TrimSpace(p.(String).Value())
		if maybeNumber(s) {
			if i, e := strconv.ParseInt(s, 10, 64); e == nil {
				return NewInteger(i)
			}
			if f, e := strconv.ParseFloat(s, 64); e == nil {
				if math.Remainder(f, 1) == 0 {
					return NewInteger(int64(f))
				}
			}
		}
	}

	return NewNull()
}

func PrimaryToFloat(p Primary) Primary {
	switch p.(type) {
	case Integer:
		return NewFloat(float64(p.(Integer).Value()))
	case Float:
		return p
	case String:
		s := strings.TrimSpace(p.(String).Value())
		if maybeNumber(s) {
			if f, e := strconv.ParseFloat(p.(String).Value(), 64); e == nil {
				return NewFloat(f)
			}
		}
	}

	return NewNull()
}

func maybeNumber(s string) bool {
	slen := len(s)
	if 1 < slen && (s[0] == '-' || s[0] == '+') && '0' <= s[1] && s[1] <= '9' {
		return true
	}
	if 0 < slen && '0' <= s[0] && s[0] <= '9' {
		if 8 <= slen {
			if s[4] == '-' && (s[6] == '-' || s[7] == '-') {
				return false
			}
			if s[4] == '/' && (s[6] == '/' || s[7] == '/') {
				return false
			}
			if s[2] == ' ' {
				return false
			}
		}
		return true
	}
	return false
}

func PrimaryToDatetime(p Primary) Primary {
	switch p.(type) {
	case Integer:
		dt := time.Unix(p.(Integer).Value(), 0)
		return NewDatetime(dt)
	case Float:
		dt := Float64ToTime(p.(Float).Value())
		return NewDatetime(dt)
	case Datetime:
		return p
	case String:
		s := strings.TrimSpace(p.(String).Value())
		if dt, e := StrToTime(s); e == nil {
			return NewDatetime(dt)
		}
		if maybeNumber(s) {
			if i, e := strconv.ParseInt(s, 10, 64); e == nil {
				dt := time.Unix(i, 0)
				return NewDatetime(dt)
			}
			if f, e := strconv.ParseFloat(s, 64); e == nil {
				dt := Float64ToTime(f)
				return NewDatetime(dt)
			}
		}
	}

	return NewNull()
}

func PrimaryToBoolean(p Primary) Primary {
	switch p.(type) {
	case Boolean:
		return p
	case Integer, Float, Ternary:
		if p.Ternary() != ternary.UNKNOWN {
			return NewBoolean(p.Ternary().BoolValue())
		}
	case String:
		s := strings.TrimSpace(p.(String).Value())
		if b, e := strconv.ParseBool(s); e == nil {
			return NewBoolean(b)
		}
	}
	return NewNull()
}

func PrimaryToString(p Primary) Primary {
	switch p.(type) {
	case String:
		return p
	case Integer:
		return NewString(Int64ToStr(p.(Integer).Value()))
	case Float:
		return NewString(Float64ToStr(p.(Float).Value()))
	}
	return NewNull()
}

func Float64ToTime(f float64) time.Time {
	s := Float64ToStr(f)
	ns := strings.Split(s, ".")
	sec, _ := strconv.ParseInt(ns[0], 10, 64)
	var nsec int64
	if 1 < len(ns) {
		if 9 < len(ns[1]) {
			ns[1] = ns[1][:9]
		}
		nsec, _ = strconv.ParseInt(ns[1]+strings.Repeat("0", 9-len(ns[1])), 10, 64)
	}
	return time.Unix(sec, nsec)
}

func FormatTableName(s string) string {
	return strings.TrimSuffix(filepath.Base(s), filepath.Ext(s))
}

func FieldIdentifier(e QueryExpression) string {
	if pt, ok := e.(PrimitiveType); ok {
		if s, ok := pt.Value.(String); ok {
			return s.Value()
		}
		if dt, ok := pt.Value.(Datetime); ok {
			return dt.Format(time.RFC3339Nano)
		}
		return pt.Value.String()
	}
	if fr, ok := e.(FieldReference); ok {
		return fr.Column.Literal
	}
	return e.String()
}
