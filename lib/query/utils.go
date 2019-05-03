package query

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/mithrandie/csvq/lib/cmd"

	"github.com/mithrandie/csvq/lib/value"
)

const LimitToUseUintSlicePool = 20

type UintPool struct {
	limitToUseSlice int
	m               map[uint]bool
	values          []uint
}

func NewUintPool(initCap int, limitToUseSlice int) *UintPool {
	return &UintPool{
		limitToUseSlice: limitToUseSlice,
		m:               make(map[uint]bool, initCap),
		values:          make([]uint, 0, initCap),
	}
}

func (c *UintPool) Exists(val uint) bool {
	if c.limitToUseSlice <= len(c.values) {
		_, ok := c.m[val]
		return ok
	}

	for i := range c.values {
		if val == c.values[i] {
			return true
		}
	}
	return false
}

func (c *UintPool) Add(val uint) {
	c.m[val] = true
	c.values = append(c.values, val)
}

func (c *UintPool) Range(fn func(idx int, value uint) error) error {
	var err error
	for i := range c.values {
		if err = fn(i, c.values[i]); err != nil {
			break
		}
	}
	return err
}

func (c *UintPool) Len() int {
	return len(c.values)
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
	values := make(map[string]int, 40)
	valueKeys := make([]string, 0, 40)

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
		return &bytes.Buffer{}
	},
}

func GetComparisonKeysBuf() *bytes.Buffer {
	buf := comparisonKeysBufPool.Get().(*bytes.Buffer)
	return buf
}

func PutComparisonkeysBuf(buf *bytes.Buffer) {
	buf.Reset()
	comparisonKeysBufPool.Put(buf)
}

func SerializeComparisonKeys(buf *bytes.Buffer, values []value.Primary, flags *cmd.Flags) {
	for i, val := range values {
		if 0 < i {
			buf.WriteByte(58)
		}
		SerializeKey(buf, val, flags)
	}
}

func SerializeKey(buf *bytes.Buffer, val value.Primary, flags *cmd.Flags) {
	if value.IsNull(val) {
		serializeNull(buf)
	} else if in := value.ToInteger(val); !value.IsNull(in) {
		serializeInteger(buf, in.(*value.Integer).Raw())
		value.Discard(in)
	} else if f := value.ToFloat(val); !value.IsNull(f) {
		serializeFloat(buf, f.(*value.Float).Raw())
		value.Discard(f)
	} else if dt := value.ToDatetime(val, flags.DatetimeFormat); !value.IsNull(dt) {
		serializeDatetime(buf, dt.(*value.Datetime).Raw())
		value.Discard(dt)
	} else if b := value.ToBoolean(val); !value.IsNull(b) {
		if b.(*value.Boolean).Raw() {
			serializeInteger(buf, 1)
		} else {
			serializeInteger(buf, 0)
		}
	} else if s, ok := val.(*value.String); ok {
		serializeString(buf, s.Raw())
	} else {
		serializeNull(buf)
	}
}

func serializeNull(buf *bytes.Buffer) {
	buf.Write([]byte{91, 78, 93})
}

func serializeInteger(buf *bytes.Buffer, i int64) {
	if i == 0 {
		buf.Write([]byte{91, 73, 93, 48})
	} else if i == 1 {
		buf.Write([]byte{91, 73, 93, 49})
	} else {
		buf.Write([]byte{91, 73, 93})
		_ = binary.Write(buf, binary.LittleEndian, i)
	}
}

func serializeFloat(buf *bytes.Buffer, f float64) {
	buf.Write([]byte{91, 70, 93})
	_ = binary.Write(buf, binary.LittleEndian, f)
}

func serializeDatetime(buf *bytes.Buffer, t time.Time) {
	serializeDatetimeFromUnixNano(buf, t.UnixNano())
}

func serializeDatetimeFromUnixNano(buf *bytes.Buffer, t int64) {
	buf.Write([]byte{91, 68, 93})
	_ = binary.Write(buf, binary.LittleEndian, t)
}

func serializeString(buf *bytes.Buffer, s string) {
	buf.Write([]byte{91, 83, 93})
	buf.WriteString(strings.ToUpper(cmd.TrimSpace(s)))
}
