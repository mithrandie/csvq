package query

import (
	"github.com/mithrandie/csvq/lib/value"
)

type Cell []value.Primary

func NewCell(val value.Primary) Cell {
	return []value.Primary{val}
}

func NewGroupCell(values []value.Primary) Cell {
	return values
}

func (cell Cell) Primary() value.Primary {
	return cell[0]
}

func (cell Cell) GroupedPrimary(index int) value.Primary {
	return cell[index]
}

func (cell Cell) Len() int {
	return len(cell)
}
