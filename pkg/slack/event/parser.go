package event

import (
	"github.com/dpb587/slack-delegate-bot/pkg/message"
	"github.com/dpb587/slack-delegate-bot/pkg/slack"
	"github.com/dpb587/slack-delegate-bot/pkg/slack/slackutil"
	"github.com/slack-go/slack/slackevents"
)

type Parser struct {
	userLookup *slack.UserLookup
}

func NewParser(userLookup *slack.UserLookup) *Parser {
	return &Parser{
		userLookup: userLookup,
	}
}

func (m *Parser) ParseAppMention(raw slackevents.EventsAPIEvent, e slackevents.AppMentionEvent) (message.Message, bool, error) {
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

	msg = slackutil.ParseMessageForChannelReference(msg, isSelf)

	return msg, true, nil
}

func (m *Parser) ParseMessage(raw slackevents.EventsAPIEvent, e slackevents.MessageEvent) (message.Message, bool, error) {
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
		msg = slackutil.ParseMessageForAnyChannelReference(msg)

		if msg.TargetChannelTeamID == msg.ChannelTeamID && msg.TargetChannelID == msg.ChannelID {
			// but if no channel mentioned in the dm, ignore them
			// TODO give a help link; move to responder?
			return message.Message{}, false, nil
		}
	} else if !slackutil.CheckMessageForMention(msg, isSelf) {
		// assume channel-style needing a reference; guess not
		// TODO give a help link; move to responder?
		return message.Message{}, false, nil
	} else {
		// check for contextual channel
		msg = slackutil.ParseMessageForChannelReference(msg, isSelf)

		if msg.TargetChannelTeamID == msg.ChannelTeamID && msg.TargetChannelID == msg.ChannelID {
			// cannot detect channel from assumed-mpim/non-channel messages
			// TODO give a help link; move to responder?
			return message.Message{}, false, nil
		}
	}

	// TODO attachments? how?

	return msg, true, nil
}

func (m *Parser) selfID(appID string) func(string) bool {
	return func(userID string) bool {
		isSelf, err := m.userLookup.IsAppBot(appID, userID)
		if err != nil {
			// TODO log
			return false
		}

		return isSelf
	}
}
