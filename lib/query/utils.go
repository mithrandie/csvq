package query

import (
	"fmt"
	"strconv"
	"strings"
	"time"

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

func RecordRange(cpuIndex int, totalLen int, numberOfCPU int) (int, int) {
	calcLen := totalLen / numberOfCPU

	var start = cpuIndex * calcLen
	var end int
	if cpuIndex == numberOfCPU-1 {
		end = totalLen
	} else {
		end = (cpuIndex + 1) * calcLen
	}
	return start, end
}
