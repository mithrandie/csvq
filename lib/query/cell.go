package query

import (
	"github.com/mithrandie/csvq/lib/parser"
)

type Cell []parser.Primary

func NewCell(value parser.Primary) Cell {
	return []parser.Primary{value}
}

func NewGroupCell(values []parser.Primary) Cell {
	return values
}

func (cell Cell) Primary() parser.Primary {
	return cell[0]
}

func (cell Cell) GroupedPrimary(index int) parser.Primary {
	return cell[index]
}

func (cell Cell) Len() int {
	return len(cell)
}
