package query

import (
	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/ternary"
)

type Records []Record

func (r Records) Contains(record Record) bool {
	for _, v := range r {
		if v.IsEqualTo(record) {
			return true
		}
	}
	return false
}

type Record []Cell

func NewRecord(internalId int, values []parser.Primary) Record {
	record := make(Record, len(values)+1)

	record[0] = NewCell(parser.NewInteger(int64(internalId)))

	for i, v := range values {
		record[i+1] = NewCell(v)
	}

	return record
}

func NewRecordWithoutId(values []parser.Primary) Record {
	record := make(Record, len(values))

	for i, v := range values {
		record[i] = NewCell(v)
	}

	return record
}

func NewEmptyRecord(len int) Record {
	record := make(Record, len)
	for i := 0; i < len; i++ {
		record[i] = NewCell(parser.NewNull())
	}

	return record
}

func (r Record) IsEqualTo(record Record) bool {
	if len(r) != len(record) {
		return false
	}

	for i, cell := range r {
		if EquivalentTo(cell.Primary(), record[i].Primary()) != ternary.TRUE {
			return false
		}
	}

	return true
}

func (r Record) GroupLen() int {
	return r[0].Len()
}

func (r Record) Copy() Record {
	record := make(Record, len(r))
	copy(record, r)
	return record

}
