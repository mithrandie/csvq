package parser

import (
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/ternary"
)

func StrToTime(s string) (time.Time, error) {
	flags := cmd.GetFlags()
	if 0 < len(flags.DatetimeFormat) {
		if t, e := time.Parse(ConvertDatetimeFormat(flags.DatetimeFormat), s); e == nil {
			return t, nil
		}
	}

	if t, e := time.Parse(time.RFC3339Nano, s); e == nil {
		return t, nil
	} else if t, e := time.ParseInLocation(DATETIME_FORMAT, s, cmd.GetLocation()); e == nil {
		return t, nil
	} else if t, e := time.Parse(DATETIME_FORMAT+" MST", s); e == nil {
		return t, nil
	} else if t, e := time.ParseInLocation("2006-01-02", s, cmd.GetLocation()); e == nil {
		return t, nil
	} else if t, e := time.ParseInLocation("2006/01/02 15:04:05.999999999", s, cmd.GetLocation()); e == nil {
		return t, nil
	} else if t, e := time.Parse("2006/01/02 15:04:05.999999999 MST", s); e == nil {
		return t, nil
	} else if t, e := time.ParseInLocation("2006/01/02", s, cmd.GetLocation()); e == nil {
		return t, nil
	} else if t, e := time.Parse(time.RFC822, s); e == nil {
		return t, nil
	} else if t, e := time.Parse(time.RFC822Z, s); e == nil {
		return t, nil
	}
	return time.Time{}, errors.New(fmt.Sprintf("%q does not match with datetime format", s))
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
	s := Float64ToStr(f)
	if i, e := strconv.ParseInt(s, 10, 64); e == nil {
		return NewInteger(i)
	}
	return NewFloat(f)
}

func PrimaryToInteger(p Primary) Primary {
	switch p.(type) {
	case Integer:
		return p
	case Float:
		if i, e := strconv.ParseInt(p.(Float).literal, 10, 64); e == nil {
			return NewInteger(i)
		}
	case String:
		if i, e := strconv.ParseInt(p.(String).Value(), 10, 64); e == nil {
			return NewInteger(i)
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
		if f, e := strconv.ParseFloat(p.(String).Value(), 64); e == nil {
			return NewFloat(f)
		}
	}

	return NewNull()
}

func PrimaryToDatetime(p Primary) Primary {
	switch p.(type) {
	case Integer:
		dt := time.Unix(p.(Integer).Value(), 0)
		return NewDatetime(dt)
	case Float:
		dt := float64ToTime(p.(Float).Value())
		return NewDatetime(dt)
	case Datetime:
		return p
	case String:
		if f, e := strconv.ParseFloat(p.(String).Value(), 64); e == nil {
			dt := float64ToTime(f)
			return NewDatetime(dt)
		}

		if dt, e := StrToTime(p.(String).Value()); e == nil {
			return NewDatetime(dt)
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
		if b, e := strconv.ParseBool(p.(String).Value()); e == nil {
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

func float64ToTime(f float64) time.Time {
	s := Float64ToStr(f)
	ns := strings.Split(s, ".")
	sec, _ := strconv.ParseInt(ns[0], 10, 64)
	var nsec int64
	if 1 < len(ns) {
		nsec, _ = strconv.ParseInt(ns[1]+strings.Repeat("0", 9-len(ns[1])), 10, 64)
	}
	return time.Unix(sec, nsec)
}

func FormatTableName(s string) string {
	return strings.TrimSuffix(path.Base(s), filepath.Ext(s))
}
