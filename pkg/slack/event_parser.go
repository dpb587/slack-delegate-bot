package slack

import (
	"github.com/dpb587/slack-delegate-bot/pkg/message"
	"github.com/dpb587/slack-delegate-bot/pkg/slack/slackutil"
	"github.com/slack-go/slack/slackevents"
)

type EventParser struct {
	userLookup *UserLookup
}

func NewEventParser(userLookup *UserLookup) *EventParser {
	return &EventParser{
		userLookup: userLookup,
	}
}

func (m *EventParser) ParseAppMention(raw slackevents.EventsAPIEvent, e slackevents.AppMentionEvent) (message.Message, bool, error) {
	isSelf := m.selfID(raw.APIAppID)

	if isSelf(e.User) {
		return message.Message{}, false, nil
	}

	msg := message.Message{
		ChannelTeamID:       raw.TeamID,
		ChannelID:           e.Channel,
		UserID:              e.User,
		UserTeamID:          e.UserTeam,
		RawTimestamp:        e.TimeStamp,
		RawThreadTimestamp:  e.ThreadTimeStamp,
		RawText:             e.Text,
		TargetChannelTeamID: raw.TeamID,
		TargetChannelID:     e.Channel,
		Type:                message.ChannelMessageType,
		Time:                slackutil.MustConvertTimestamp(e.TimeStamp),
	}

	// TODO attachments?

	msg = ParseMessageForChannelReference(msg, isSelf)

	return msg, true, nil
}

func (m *EventParser) ParseMessage(raw slackevents.EventsAPIEvent, e slackevents.MessageEvent) (message.Message, bool, error) {
	isSelf := m.selfID(raw.APIAppID)

	if isSelf(e.User) {
		return message.Message{}, false, nil
	}

	msg := message.Message{
		ChannelTeamID:       raw.TeamID,
		ChannelID:           e.Channel,
		UserID:              e.User,
		UserTeamID:          e.UserTeam,
		TargetChannelTeamID: raw.TeamID,
		TargetChannelID:     e.Channel,
		RawTimestamp:        e.TimeStamp,
		RawThreadTimestamp:  e.ThreadTimeStamp,
		RawText:             e.Text,
		Type:                message.ChannelMessageType,
		Time:                slackutil.MustConvertTimestamp(e.TimeStamp),
	}

	if e.ChannelType == "im" {
		msg.Type = message.DirectMessageMessageType

		// assume they mention a channel for the interrupt
		msg = ParseMessageForAnyChannelReference(msg)

		if msg.TargetChannelTeamID == msg.ChannelTeamID && msg.TargetChannelID == msg.ChannelID {
			// but if no channel mentioned in the dm, ignore them
			// TODO give a help link; move to responder?
			return message.Message{}, false, nil
		}
	} else if !CheckMessageForMention(msg, isSelf) {
		// assume channel-style needing a reference; guess not
		// TODO give a help link; move to responder?
		return message.Message{}, false, nil
	} else {
		// check for contextual channel
		msg = ParseMessageForChannelReference(msg, isSelf)

		if msg.TargetChannelTeamID == msg.ChannelTeamID && msg.TargetChannelID == msg.ChannelID {
			// cannot detect channel from assumed-mpim/non-channel messages
			// TODO give a help link; move to responder?
			return message.Message{}, false, nil
		}
	}

	// TODO attachments? how?

	return msg, true, nil
}

func (m *EventParser) selfID(appID string) func(string) bool {
	return func(userID string) bool {
		isSelf, err := m.userLookup.IsAppBot(appID, userID)
		if err != nil {
			// TODO log
			return false
		}

		return isSelf
	}
}
