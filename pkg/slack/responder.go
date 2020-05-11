package slack

import (
	"fmt"

	"github.com/dpb587/slack-delegate-bot/pkg/delegate"
	"github.com/dpb587/slack-delegate-bot/pkg/delegate/delegates"
	"github.com/dpb587/slack-delegate-bot/pkg/message"
	"github.com/pkg/errors"
	"github.com/slack-go/slack"
)

//go:generate counterfeiter . ResponderSlackAPI
type ResponderSlackAPI interface {
	PostMessage(string, ...slack.MsgOption) (string, string, error)
}

type Responder struct {
	api       ResponderSlackAPI
	delegator delegate.Delegator
}

func NewResponder(api ResponderSlackAPI, delegator delegate.Delegator) *Responder {
	return &Responder{
		api:       api,
		delegator: delegator,
	}
}

func (m *Responder) ProcessMessage(msg message.Message) error {
	msg.ServiceAPI = m.api
	msg.Delegator = m.delegator

	dd, err := m.delegator.Delegate(msg)
	if err != nil {
		return errors.Wrap(err, "finding delegate")
	}

	if len(dd) == 0 {
		return nil
	}

	responseText := delegates.Join(dd, " ")

	if msg.Type == message.ChannelMessageType {
		responseText = fmt.Sprintf("^ %s", responseText)
	}

	opts := []slack.MsgOption{
		slack.MsgOptionText(responseText, false),
	}

	if v := msg.RawThreadTimestamp; v != "" {
		// always stay in thread if one is started
		opts = append(opts, slack.MsgOptionTS(v))
	} else if msg.Type == message.ChannelMessageType {
		// always use a thread if in a channel
		opts = append(opts, slack.MsgOptionTS(msg.RawTimestamp))
	}

	_, _, err = m.api.PostMessage(msg.ChannelID, opts...)
	if err != nil {
		return errors.Wrap(err, "posting message")
	}

	return nil
}
