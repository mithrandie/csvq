package query

import (
	"strings"

	"github.com/mithrandie/csvq/lib/value"
)

type Records []Record

func (r Records) Copy() Records {
	records := make(Records, len(r))
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

func (r Record) SerializeComparisonKeys() string {
	list := make([]string, len(r))

	for i, cell := range r {
		list[i] = SerializeKey(cell.Primary())
	}

	return strings.Join(list, ":")
}

func MergeRecord(r1 Record, r2 Record) Record {
	r := make(Record, len(r1)+len(r2))
	for i, v := range r1 {
		r[i] = v
	}
	for i, v := range r2 {
		r[i+len(r1)] = v
	}
	return r
}

func MergeRecordsList(list []Records) Records {
	var records Records
	if len(list) == 1 {
		records = list[0]
	} else {
		recordLen := 0
		for _, v := range list {
			recordLen += len(v)
		}
		records = make(Records, recordLen)
		idx := 0
		for _, v := range list {
			for _, r := range v {
				records[idx] = r
				idx++
			}
		}
	}
	return records
}
