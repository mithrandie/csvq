package query

import (
	"strings"

	"github.com/mithrandie/csvq/lib/parser"
)

const INTERNAL_ID_COLUMN = "@__internal_id"

type HeaderField struct {
	Reference  string
	Column     string
	Alias      string
	Number     int
	FromTable  bool
	IsGroupKey bool
}

func (hf HeaderField) Label() string {
	if 0 < len(hf.Alias) {
		return hf.Alias
	}
	return hf.Column
}

type Header []HeaderField

func NewDualHeader() Header {
	h := make([]HeaderField, 1)
	return h
}

func NewHeaderWithId(ref string, words []string) Header {
	h := make([]HeaderField, len(words)+1)

	h[0].Reference = ref
	h[0].Column = INTERNAL_ID_COLUMN

	for i, v := range words {
		h[i+1].Reference = ref
		h[i+1].Column = v
		h[i+1].Number = i + 1
		h[i+1].FromTable = true
	}

	return h
}

func NewHeader(ref string, words []string) Header {
	h := make([]HeaderField, len(words))

	for i, v := range words {
		h[i].Reference = ref
		h[i].Column = v
		h[i].Number = i + 1
		h[i].FromTable = true
	}

	return h
}

func NewEmptyHeader(len int) Header {
	return make([]HeaderField, len)
}

func MergeHeader(h1 Header, h2 Header) Header {
	return append(h1, h2...)
}

func AddHeaderField(h Header, column string, alias string) (header Header, index int) {
	header = append(h, HeaderField{
		Column: column,
		Alias:  alias,
	})
	index = header.Len() - 1
	return
}

func (h Header) Len() int {
	return len(h)
}

func (h Header) TableColumns() []parser.Expression {
	columns := []parser.Expression{}
	for _, f := range h {
		if !f.FromTable {
			continue
		}

		fieldRef := parser.FieldReference{
			Column: parser.Identifier{Literal: f.Column},
		}
		if 0 < len(f.Reference) {
			fieldRef.View = parser.Identifier{Literal: f.Reference}
		}

		columns = append(columns, fieldRef)
	}
	return columns
}

func (h Header) TableColumnNames() []string {
	names := []string{}
	for _, f := range h {
		if !f.FromTable {
			continue
		}
		names = append(names, f.Column)
	}
	return names
}

func (h Header) ContainsObject(obj parser.Expression) (int, error) {
	var fieldRef parser.FieldReference

	if fr, ok := obj.(parser.FieldReference); ok {
		fieldRef = fr
	} else {
		fieldRef = parser.FieldReference{
			Column: parser.Identifier{Literal: obj.String()},
		}
	}
	return h.Contains(fieldRef)
}

func (h Header) ContainsNumber(number parser.ColumnNumber) (int, error) {
	ref := number.View.Literal
	idx := int(number.Number.Value())

	if idx < 1 {
		return -1, NewFieldNumberNotExistError(number)
	}

	for i, f := range h {
		if strings.EqualFold(f.Reference, ref) && f.Number == idx {
			return i, nil
		}
	}
	return -1, NewFieldNumberNotExistError(number)
}

func (h Header) Contains(fieldRef parser.FieldReference) (int, error) {
	var ref string
	if fieldRef.View != nil {
		ref = fieldRef.View.(parser.Identifier).Literal
	}
	column := fieldRef.Column.Literal

	idx := -1

	for i, f := range h {
		if 0 < len(ref) {
			if !strings.EqualFold(f.Reference, ref) || !strings.EqualFold(f.Column, column) {
				continue
			}
		} else {
			if !strings.EqualFold(f.Column, column) && !strings.EqualFold(f.Alias, column) {
				continue
			}
		}

		if -1 < idx {
			return -1, NewFieldAmbiguousError(fieldRef)
		}
		idx = i
	}

	if idx < 0 {
		return -1, NewFieldNotExistError(fieldRef)
	}

	return idx, nil
}

func (h Header) Update(reference string, fields []parser.Expression) error {
	if fields != nil {
		if len(fields) != h.Len() {
			return NewFieldLengthNotMatchError()
		}

		names := make([]string, len(fields))
		for i, v := range fields {
			f, _ := v.(parser.Identifier)
			if InStrSliceWithCaseInsensitive(f.Literal, names) {
				return NewDuplicateFieldNameError(f)
			}
			names[i] = f.Literal
		}
	}

	for i, hf := range h {
		h[i].Reference = reference
		if fields != nil {
			h[i].Column = fields[i].(parser.Identifier).Literal
		} else if 0 < len(hf.Alias) {
			h[i].Column = hf.Alias
		}
		h[i].Alias = ""
	}
	return nil
}

func (h Header) Copy() Header {
	header := make(Header, h.Len())
	copy(header, h)
	return header

}
