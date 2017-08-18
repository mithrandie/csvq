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

func (r Records) Copy() Records {
	records := make(Records, len(r))
	for i, v := range r {
		records[i] = v.Copy()
	}
	return records
}

type Record []Cell

func NewRecordWithId(internalId int, values []parser.Primary) Record {
	record := make(Record, len(values)+1)

	record[0] = NewCell(parser.NewInteger(int64(internalId)))

	for i, v := range values {
		record[i+1] = NewCell(v)
	}

	return record
}

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

func (r Record) Match(record Record, indices []int) bool {
	for _, i := range indices {
		if EquivalentTo(r[i].Primary(), record[i].Primary()) != ternary.TRUE {
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
