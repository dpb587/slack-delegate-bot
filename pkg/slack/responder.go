package slack

import (
	"fmt"

	"github.com/dpb587/slack-delegate-bot/pkg/delegate/delegates"
	"github.com/dpb587/slack-delegate-bot/pkg/handler"
	"github.com/dpb587/slack-delegate-bot/pkg/message"
	"github.com/pkg/errors"
	"github.com/slack-go/slack"
)

//go:generate counterfeiter . ResponderSlackAPI
type ResponderSlackAPI interface {
	PostMessage(string, ...slack.MsgOption) (string, string, error)
}

type Responder struct {
	api     ResponderSlackAPI
	handler handler.Handler
}

func NewResponder(api ResponderSlackAPI, handler handler.Handler) *Responder {
	return &Responder{
		api:     api,
		handler: handler,
	}
}

func (m *Responder) ProcessMessage(msg message.Message) error {
	msg.ServiceAPI = m.api

	response, err := m.handler.Execute(msg)
	if err != nil {
		return errors.Wrap(err, "finding delegate")
	}

	var responseText string

	if len(response.Delegates) > 0 {
		responseText = delegates.Join(response.Delegates, " ")

		if msg.Type == message.ChannelMessageType {
			responseText = fmt.Sprintf("^ %s", responseText)
		}
	} else if response.EmptyMessage != "" {
		responseText = response.EmptyMessage
	} else {
		return nil
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
