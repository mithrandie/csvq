package query

import "github.com/mithrandie/csvq/lib/parser"

const INTERNAL_ID_FIELD = "@__internal_id"

type HeaderField struct {
	Reference  string
	Column     string
	Alias      string
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

func NewHeader(ref string, words []string) Header {
	h := make([]HeaderField, len(words)+1)

	h[0].Reference = ref
	h[0].Column = INTERNAL_ID_FIELD

	for i, v := range words {
		h[i+1].Reference = ref
		h[i+1].Column = v
		h[i+1].FromTable = true
	}

	return h
}

func NewHeaderWithoutId(ref string, words []string) Header {
	h := make([]HeaderField, len(words))

	for i, v := range words {
		h[i].Reference = ref
		h[i].Column = v
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

func AddHeaderField(h Header, alias string) (header Header, index int) {
	header = append(h, HeaderField{
		Alias: alias,
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

		var lit string
		if 0 < len(f.Reference) {
			lit = f.Reference + "." + f.Column
		} else {
			lit = f.Column
		}
		columns = append(columns, parser.Identifier{Literal: lit})
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

func (h Header) Contains(ref string, column string) (int, error) {
	identifier := column
	if 0 < len(ref) {
		identifier = ref + "." + column
	}

	idx := -1

	for i, f := range h {
		if 0 < len(ref) {
			if f.Reference != ref || f.Column != column {
				continue
			}
		} else {
			if f.Column != column && f.Alias != column {
				continue
			}
		}

		if -1 < idx {
			return -1, h.newError(identifier, ErrFieldAmbiguous)
		}
		idx = i
	}

	if idx < 0 {
		return -1, h.newError(identifier, ErrFieldNotExist)
	}

	return idx, nil
}

func (h Header) newError(identifier string, err error) error {
	return &IdentificationError{
		Identifier: identifier,
		Err:        err,
	}
}

func (h Header) Copy() Header {
	header := make(Header, h.Len())
	copy(header, h)
	return header

}
