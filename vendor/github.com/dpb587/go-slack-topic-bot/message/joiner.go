package message

import (
	"strings"

	"github.com/pkg/errors"
)

type joiner struct {
	delimiter string
	messages  []Messager
}

var _ Messager = &joiner{}

func Join(delimiter string, templates ...Messager) Messager {
	return joiner{
		delimiter: delimiter,
		messages:  templates,
	}
}

func (m joiner) Message() (string, error) {
	var msgs []string

	for tplIdx, tpl := range m.messages {
		msg, err := tpl.Message()
		if err != nil {
			return "", errors.Wrapf(err, "template %d", tplIdx)
		}

		if msg == "" {
			continue
		}

		msgs = append(msgs, msg)
	}

	return strings.Join(msgs, m.delimiter), nil
}
