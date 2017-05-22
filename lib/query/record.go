package query

import "github.com/mithrandie/csvq/lib/parser"

type Record []Cell

func NewRecord(values []parser.Primary) Record {
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

func (record Record) GroupLen() int {
	return record[0].Len()
}
