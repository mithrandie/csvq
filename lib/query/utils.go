package query

import (
	"bytes"
	"fmt"
	"strings"
	"sync"
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

func Distinguish(list []value.Primary, flags *cmd.Flags) []value.Primary {
	values := make(map[string]int)
	valueKeys := make([]string, 0, len(list))

	buf := GetComparisonKeysBuf()

	for i, v := range list {
		buf.Reset()
		SerializeComparisonKeys(buf, []value.Primary{v}, flags)
		key := buf.String()
		if _, ok := values[key]; !ok {
			values[key] = i
			valueKeys = append(valueKeys, key)
		}
	}

	PutComparisonkeysBuf(buf)

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

var comparisonKeysBufPool = &sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

func GetComparisonKeysBuf() *bytes.Buffer {
	buf := comparisonKeysBufPool.Get().(*bytes.Buffer)
	buf.Reset()
	return buf
}

func PutComparisonkeysBuf(buf *bytes.Buffer) {
	comparisonKeysBufPool.Put(buf)
}

func SerializeComparisonKeys(buf *bytes.Buffer, values []value.Primary, flags *cmd.Flags) {
	for i, val := range values {
		if 0 < i {
			buf.WriteString(":")
		}
		SerializeKey(buf, val, flags)
	}
}

func SerializeKey(buf *bytes.Buffer, val value.Primary, flags *cmd.Flags) {
	if value.IsNull(val) {
		serializeNull(buf)
	} else if in := value.ToInteger(val); !value.IsNull(in) {
		serializeInteger(buf, in.(*value.Integer).Raw())
	} else if f := value.ToFloat(val); !value.IsNull(f) {
		serializeFlaot(buf, f.(*value.Float).Raw())
	} else if dt := value.ToDatetime(val, flags.DatetimeFormat); !value.IsNull(dt) {
		t := dt.(*value.Datetime).Raw()
		if t.Nanosecond() > 0 {
			f := float64(t.Unix()) + float64(t.Nanosecond())/1e9
			t2 := value.Float64ToTime(f)
			if t.Equal(t2) {
				serializeFlaot(buf, f)
			} else {
				serializeDatetime(buf, t)
			}
		} else {
			serializeInteger(buf, t.Unix())
		}
	} else if b := value.ToBoolean(val); !value.IsNull(b) {
		serializeBoolean(buf, b.(*value.Boolean).Raw())
	} else if s, ok := val.(*value.String); ok {
		serializeString(buf, s.Raw())
	} else {
		serializeNull(buf)
	}
}

func serializeNull(buf *bytes.Buffer) {
	buf.WriteString("[N]")
}

func serializeInteger(buf *bytes.Buffer, i int64) {
	buf.WriteString("[I]")
	buf.WriteString(value.Int64ToStr(i))
	switch i {
	case 0:
		buf.WriteString("[B]")
		buf.WriteString("false")
	case 1:
		buf.WriteString("[B]")
		buf.WriteString("true")
	}
}

func serializeFlaot(buf *bytes.Buffer, f float64) {
	buf.WriteString("[F]")
	buf.WriteString(value.Float64ToStr(f))
}

func serializeDatetime(buf *bytes.Buffer, t time.Time) {
	buf.WriteString("[D]")
	buf.WriteString(value.Int64ToStr(t.UnixNano()))
}

func serializeDatetimeFromUnixNano(buf *bytes.Buffer, t int64) {
	buf.WriteString("[D]")
	buf.WriteString(value.Int64ToStr(t))
}

func serializeBoolean(buf *bytes.Buffer, b bool) {
	buf.WriteString("[I]")
	if b {
		buf.WriteString("1")
	} else {
		buf.WriteString("0")
	}
	buf.WriteString("[B]")
	if b {
		buf.WriteString("true")
	} else {
		buf.WriteString("false")
	}
}

func serializeString(buf *bytes.Buffer, s string) {
	buf.WriteString("[S]")
	buf.WriteString(strings.ToUpper(strings.TrimSpace(s)))
}
