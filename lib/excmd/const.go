package excmd

const EOF = -(iota + 1)

type NodeType int

const (
	FixedString         NodeType = 1
	Variable            NodeType = 2
	EnvironmentVariable NodeType = 3
	RuntimeInformation  NodeType = 4
	CsvqExpression      NodeType = 5
)
