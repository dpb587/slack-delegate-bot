package slack

import (
	"fmt"

	"github.com/dpb587/slack-delegate-bot/cmd/delegatebot/handler"
	"github.com/dpb587/slack-delegate-bot/pkg/delegate/delegates"
	"github.com/dpb587/slack-delegate-bot/pkg/message"
	"github.com/dpb587/slack-delegate-bot/pkg/slack/slackutil"
	"github.com/pkg/errors"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

type Messenger struct {
	api     *slack.Client
	handler handler.Handler
	cache   *Cache
}

func NewMessenger(api *slack.Client, handler handler.Handler) *Messenger {
	return &Messenger{
		api:     api,
		handler: handler,
		cache:   NewCache(api),
	}
}

func (m *Messenger) HandleAppMention(c EventContext, e slackevents.AppMentionEvent) error {
	isSelf := m.selfID(c.AppID)

	if isSelf(e.User) {
		return nil
	}

	msg := message.Message{
		TeamID:                c.TeamID,
		OriginUserID:          e.User,
		Origin:                e.Channel,
		OriginTimestamp:       e.TimeStamp,
		OriginThreadTimestamp: e.ThreadTimeStamp,
		OriginType:            message.ChannelOriginType,
		InterruptTarget:       e.Channel,
		Timestamp:             slackutil.MustConvertTimestamp(e.TimeStamp),
		Text:                  e.Text,
	}

	// TODO attachments?

	msg = ParseMessageForChannelReference(msg, isSelf)

	return m.ProcessMessage(msg)
}

func (m *Messenger) HandleMessage(c EventContext, e slackevents.MessageEvent) error {
	isSelf := m.selfID(c.AppID)

	if isSelf(e.User) {
		return nil
	}

	msg := message.Message{
		TeamID:                c.TeamID,
		OriginUserID:          e.User,
		Origin:                e.Channel,
		OriginTimestamp:       e.TimeStamp,
		OriginThreadTimestamp: e.ThreadTimeStamp,
		OriginType:            message.ChannelOriginType,
		InterruptTarget:       e.Channel,
		Timestamp:             slackutil.MustConvertTimestamp(e.TimeStamp),
		Text:                  e.Text,
	}

	if e.ChannelType == "im" {
		msg.OriginType = message.DirectMessageOriginType

		// assume they mention a channel for the interrupt
		msg = ParseMessageForAnyChannelReference(msg)

		if msg.InterruptTarget == msg.Origin {
			// but if no channel mentioned in the dm, ignore them
			// TODO give a help link
			return nil
		}
	} else if !CheckMessageForMention(msg, isSelf) {
		// assume channel-style needing a reference; guess not
		msg = ParseMessageForChannelReference(msg, isSelf)

		return nil
	}

	// TODO attachments? how?

	return m.ProcessMessage(msg)
}

func (m *Messenger) ProcessMessage(msg message.Message) error {
	msg.ServiceAPI = m.api

	response, err := m.handler.Execute(&msg)
	if err != nil {
		return errors.Wrap(err, "finding delegate")
	}

	var responseText string

	if len(response.Delegates) > 0 {
		responseText = delegates.Join(response.Delegates, " ")

		if msg.OriginType == message.ChannelOriginType {
			responseText = fmt.Sprintf("^ %s", responseText)
		}
	} else if response.EmptyMessage != "" {
		responseText = response.EmptyMessage
	}

	opts := []slack.MsgOption{
		slack.MsgOptionText(responseText, false),
	}

	if v := msg.OriginThreadTimestamp; v != "" {
		// always stay in thread if one is started
		opts = append(opts, slack.MsgOptionTS(v))
	} else if msg.OriginType == message.ChannelOriginType {
		// always use a thread if in a channel
		opts = append(opts, slack.MsgOptionTS(msg.OriginTimestamp))
	}

	_, _, err = m.api.PostMessage(msg.Origin, opts...)
	if err != nil {
		return errors.Wrap(err, "posting message")
	}

	return nil
}

func (m *Messenger) selfID(appID string) func(string) bool {
	return func(userID string) bool {
		isSelf, err := m.cache.IsAppBot(appID, userID)
		if err != nil {
			// TODO log
			return false
		}

		return isSelf
	}
}
