package parser

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/mithrandie/csvq/lib/ternary"
)

func StrToTime(s string) (time.Time, error) {
	if t, e := time.ParseInLocation(DATETIME_FORMAT, s, time.Local); e == nil {
		return t, nil
	} else if t, e := time.Parse(DATETIME_FORMAT+" MST", s); e == nil {
		return t, nil
	} else if t, e := time.Parse(time.RFC822, s); e == nil {
		return t, nil
	} else if t, e := time.Parse(time.RFC822Z, s); e == nil {
		return t, nil
	} else if t, e := time.Parse(time.RFC3339Nano, s); e == nil {
		return t, nil
	}
	return time.Time{}, errors.New(fmt.Sprintf("%q does not match with datetime format", s))
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
	case Ternary:
		if p.(Ternary).Ternary() != ternary.UNKNOWN {
			return NewBoolean(p.(Ternary).Bool())
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
