package query

import "github.com/mithrandie/csvq/lib/parser"

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

func (record Record) GroupLen() int {
	return record[0].Len()
}
