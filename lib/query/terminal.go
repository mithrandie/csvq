package query

const (
	TerminalPrompt           string = "csvq > "
	TerminalContinuousPrompt string = "     > "
)

type VirtualTerminal interface {
	ReadLine() (string, error)
	Write(string) error
	WriteError(string) error
	SetPrompt()
	SetContinuousPrompt()
	SaveHistory(string)
	Teardown()
	GetSize() (int, int, error)
	RestoreRawMode() error
	RestoreOriginalMode() error
}
