package message

type coalescer struct {
	messages []Messager
}

func Coalesce(messages ...Messager) Messager {
	return coalescer{
		messages: messages,
	}
}

var _ Messager = &coalescer{}

func (m coalescer) Message() (string, error) {
	for _, mc := range m.messages {
		msg, err := mc.Message()
		if err != nil {
			return "", err
		} else if msg != "" {
			return msg, nil
		}
	}

	return "", nil
}
