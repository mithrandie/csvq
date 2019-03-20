package query

type Discard struct {
}

func NewDiscard() *Discard {
	return &Discard{}
}

func (d Discard) Write(p []byte) (int, error) {
	return len(p), nil
}

func (d Discard) Close() error {
	return nil
}
