package cmd

const (
	TERMINAL_PROMPT            string = "csvq > "
	TERMINAL_CONTINUOUS_PROMPT string = "     > "
)

type VirtualTerminal interface {
	ReadLine() (string, error)
	Write(string) error
	SetPrompt()
	SetContinuousPrompt()
	SaveHistory(string)
	Teardown()
}
