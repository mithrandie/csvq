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
		prefix := "@__PT:"
		switch pt.Value.(type) {
		case *value.String:
			prefix = prefix + "S"
		case *value.Integer:
			prefix = prefix + "I"
		case *value.Float:
			prefix = prefix + "F"
		case *value.Boolean:
			prefix = prefix + "B"
		case *value.Ternary:
			prefix = prefix + "T"
		case *value.Datetime:
			prefix = prefix + "D"
		case *value.Null:
			prefix = prefix + "N"
		}
		return prefix + ":" + FormatFieldLabel(e)
	}
	if fr, ok := e.(FieldReference); ok {
		if col, ok := fr.Column.(Identifier); ok {
			return "@__IDENT:" + col.Literal
		}
	}
	return e.String()
}

func FormatFieldLabel(e QueryExpression) string {
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
