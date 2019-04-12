package message

type conditionaler struct {
	conditional func() bool
	message     Messager
}

func Conditional(conditional func() bool, message Messager) Messager {
	return conditionaler{
		conditional: conditional,
		message:     message,
	}
}

var _ Messager = &conditionaler{}

func (m conditionaler) Message() (string, error) {
	if !m.conditional() {
		return "", nil
	}

	return m.message.Message()
}
