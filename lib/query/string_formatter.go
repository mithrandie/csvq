package query

import (
	"bytes"
	"fmt"
	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/value"
	"reflect"
	"strconv"
	"time"
)

const EOF = -1

type StringFormatter struct {
	format    []rune
	formatPos int
	offset    int
	values    []value.Primary

	buf bytes.Buffer

	err error
}

func NewStringFormatter() *StringFormatter {
	return &StringFormatter{}
}

func (f *StringFormatter) Format(format string, values []value.Primary) (string, error) {
	f.format = []rune(format)
	f.values = values
	f.formatPos = 0
	f.buf.Reset()

	var placeholder bytes.Buffer
	placeholderOrder := 0

	for {
		ch := f.next()
		if ch == EOF {
			break
		}

		if ch != '%' {
			f.buf.WriteRune(ch)
			continue
		}

		ch = f.next()
		if ch == '%' {
			f.buf.WriteRune(ch)
			continue
		}

		if len(values) <= placeholderOrder {
			return "", NewFormatStringLengthNotMatchError()
		}

		placeholder.Reset()
		placeholder.WriteRune('%')

		precision := -1

		if f.isFlag(ch) {
			placeholder.WriteRune(ch)
			ch = f.next()
		}

		if f.isDecimal(ch) {
			f.offset = 1
			f.scanDecimal()
			placeholder.WriteString(f.literal())
			ch = f.next()
		}

		if ch == '.' {
			ch = f.next()
			if f.isDecimal(ch) {
				f.offset = 1
				f.scanDecimal()
				precision = f.integer()
				ch = f.next()
			} else {
				precision = 0
			}
		}

		switch ch {
		case 's', 'q', 'i', 'T':
			placeholder.WriteRune('s')
		default:
			if -1 < precision {
				placeholder.WriteRune('.')
				placeholder.WriteString(strconv.Itoa(precision))
			}

			placeholder.WriteRune(ch)
		}

		var s string

		switch ch {
		case 'b', 'o', 'd', 'x', 'X':
			var val int64 = 0

			p := value.ToInteger(values[placeholderOrder])
			if !value.IsNull(p) {
				val = p.(value.Integer).Raw()
			}

			s = fmt.Sprintf(placeholder.String(), val)
		case 'e', 'E', 'f':
			var val float64 = 0

			p := value.ToFloat(values[placeholderOrder])
			if !value.IsNull(p) {
				val = p.(value.Float).Raw()
			}

			s = fmt.Sprintf(placeholder.String(), val)
		case 's', 'q', 'i':
			switch values[placeholderOrder].(type) {
			case value.String:
				s = values[placeholderOrder].(value.String).Raw()
			case value.Datetime:
				s = values[placeholderOrder].(value.Datetime).Format(time.RFC3339Nano)
			default:
				s = values[placeholderOrder].String()
			}

			if -1 < precision {
				s = s[:precision]
			}

			switch ch {
			case 'q':
				s = cmd.QuoteString(s)
			case 'i':
				s = cmd.QuoteIdentifier(s)
			}

			s = fmt.Sprintf(placeholder.String(), s)
		case 'T':
			s = reflect.TypeOf(values[placeholderOrder]).Name()
			if -1 < precision {
				s = s[:precision]
			}

			s = fmt.Sprintf(placeholder.String(), s)
		case EOF:
			return "", NewFormatUnexpectedTerminationError()
		default:
			return "", NewUnknownFormatPlaceholderError(ch)
		}

		f.buf.WriteString(s)

		placeholderOrder++
	}

	if placeholderOrder < len(values) {
		return "", NewFormatStringLengthNotMatchError()
	}

	return f.buf.String(), nil
}

func (f *StringFormatter) runes() []rune {
	return f.format[(f.formatPos - f.offset):f.formatPos]
}

func (f *StringFormatter) literal() string {
	return string(f.runes())
}

func (f *StringFormatter) integer() int {
	i, _ := strconv.Atoi(string(f.runes()))
	return i
}

func (f *StringFormatter) peek() rune {
	if len(f.format) <= f.formatPos {
		return EOF
	}
	return f.format[f.formatPos]
}

func (f *StringFormatter) next() rune {
	ch := f.peek()
	if ch == EOF {
		return ch
	}

	f.formatPos++
	f.offset++
	return ch
}

func (f *StringFormatter) isFlag(ch rune) bool {
	switch ch {
	case '+', '-', ' ', '0':
		return true
	default:
		return false
	}
}

func (f *StringFormatter) isDecimal(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

func (f *StringFormatter) scanDecimal() {
	for f.isDecimal(f.peek()) {
		f.next()
	}
}
