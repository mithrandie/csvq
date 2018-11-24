package excmd

const EOF = -(iota + 1)

type ElementType int

const (
	FixedString         ElementType = 1
	Variable            ElementType = 2
	EnvironmentVariable ElementType = 3
	RuntimeInformation  ElementType = 4
	CsvqExpression      ElementType = 5
)
