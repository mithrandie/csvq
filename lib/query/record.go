package query

import (
	"bytes"

	"github.com/mithrandie/csvq/lib/value"
)

type RecordSet []Record

func (r RecordSet) Copy() RecordSet {
	records := make(RecordSet, len(r))
	for i, v := range r {
		records[i] = v.Copy()
	}
	return records
}

type Record []Cell

func NewRecordWithId(internalId int, values []value.Primary) Record {
	record := make(Record, len(values)+1)

	record[0] = NewCell(value.NewInteger(int64(internalId)))

	for i, v := range values {
		record[i+1] = NewCell(v)
	}

	return record
}

func NewRecord(values []value.Primary) Record {
	record := make(Record, len(values))

	for i, v := range values {
		record[i] = NewCell(v)
	}

	return record
}

func NewEmptyRecord(len int) Record {
	record := make(Record, len)
	for i := 0; i < len; i++ {
		record[i] = NewCell(value.NewNull())
	}

	return record
}

func (r Record) GroupLen() int {
	return r[0].Len()
}

func (r Record) Copy() Record {
	record := make(Record, len(r))
	copy(record, r)
	return record

}

func (r Record) SerializeComparisonKeys(buf *bytes.Buffer) {
	for i, cell := range r {
		if 0 < i {
			buf.WriteRune(':')
		}
		SerializeKey(buf, cell.Value())
	}
}

func MergeRecordSetList(list []RecordSet) RecordSet {
	var records RecordSet
	if len(list) == 1 {
		records = list[0]
	} else {
		recordLen := 0
		for _, v := range list {
			recordLen += len(v)
		}
		records = make(RecordSet, 0, recordLen)
		for _, v := range list {
			records = append(records, v...)
		}
	}
	return records
}
