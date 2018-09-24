package cmd

const (
	TerminalPrompt           string = "csvq > "
	TerminalContinuousPrompt string = "     > "
)

type VirtualTerminal interface {
	ReadLine() (string, error)
	Write(string) error
	SetPrompt()
	SetContinuousPrompt()
	SaveHistory(string)
	Teardown()
}
