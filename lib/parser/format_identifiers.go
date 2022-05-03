package parser

import (
	"path/filepath"
	"strings"
	"time"

	"github.com/mithrandie/csvq/lib/value"
)

func FormatTableName(s string) string {
	if len(s) < 1 {
		return ""
	}
	return strings.TrimSuffix(filepath.Base(s), filepath.Ext(s))
}

func FormatFieldIdentifier(e QueryExpression) string {
	if pt, ok := e.(PrimitiveType); ok {
		if s, ok := pt.Value.(*value.String); ok {
			return s.Raw()
		}
		if dt, ok := pt.Value.(*value.Datetime); ok {
			return dt.Format(time.RFC3339Nano)
		}
		return pt.Value.String()
	}
	if fr, ok := e.(FieldReference); ok {
		if col, ok := fr.Column.(Identifier); ok {
			return col.Literal
		}
	}
	return e.String()
}
