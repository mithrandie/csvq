package value

type Field []byte

func NewField(s string) Field {
	return []byte(s)
}

func (f Field) ToPrimary() Primary {
	if f == nil {
		return NewNull()
	} else {
		return NewString(string(f))
	}
}
