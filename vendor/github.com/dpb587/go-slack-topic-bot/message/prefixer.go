package message

import (
	"fmt"
)

type prefixer struct {
	prefix  string
	message Messager
}

func Prefix(prefix string, message Messager) Messager {
	return prefixer{
		prefix:  prefix,
		message: message,
	}
}

var _ Messager = &prefixer{}

func (m prefixer) Message() (string, error) {
	msg, err := m.message.Message()
	if err != nil {
		return "", err
	}

	if msg == "" {
		return "", nil
	}

	return fmt.Sprintf("%s%s", m.prefix, msg), nil
}
