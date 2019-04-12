package message

type literal struct {
	message string
}

func Literal(message string) Messager {
	return &literal{
		message: message,
	}
}

var _ Messager = &literal{}

func (m literal) Message() (string, error) {
	return m.message, nil
}
