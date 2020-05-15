package rtm

import (
	"fmt"
	"strings"

	"github.com/dpb587/slack-delegate-bot/pkg/message"
	"github.com/dpb587/slack-delegate-bot/pkg/slack/slackutil"
	"github.com/slack-go/slack"
)

type Parser struct {
	teamID string
	userID string
}

func NewParser(teamID, userID string) *Parser {
	return &Parser{
		teamID: teamID,
		userID: userID,
	}
}

func (p *Parser) ParseMessage(msg slack.Msg) (message.Message, bool, error) {
	if msg.Type != "message" {
		return message.Message{}, false, nil
	} else if msg.SubType == "message_deleted" {
		// no sense responding to deleted message notifications
		return message.Message{}, false, nil
	} else if msg.SubType == "group_topic" || strings.Contains(msg.Text, "set the channel topic: ") {
		// no sense responding to a reference in the topic
		// trivia: slack doesn't support topic threads, but still allows bots to
		// respond which means you get mentioned, but the browser app doesn't
		// render the thread in New Threads so you can't mark it as read unless you
		// use the mobile app (which happens to show it as -1 replies).
		return message.Message{}, false, nil
	} else if p.isSelf(msg.User) {
		// avoid accidentally talking to ourselves into a recursive DoS
		return message.Message{}, false, nil
	}

	incoming := message.Message{
		UserTeamID: p.teamID, // TODO incorrect? i.e. shared channels
		UserID:     msg.User,

		ChannelTeamID: p.teamID,
		ChannelID:     msg.Channel,

		TargetChannelTeamID: p.teamID,
		TargetChannelID:     msg.Channel,

		RawText:            msg.Text,
		RawTimestamp:       msg.Timestamp,
		RawThreadTimestamp: msg.ThreadTimestamp,

		Type: message.ChannelMessageType,
		Time: slackutil.MustConvertTimestamp(msg.Timestamp),
	}

	// include attachments
	for _, attachment := range msg.Attachments {
		if attachment.Fallback == "" {
			continue
		}

		incoming.RawText = fmt.Sprintf("%s\n\n---\n\n%s", incoming.RawText, attachment.Fallback)
	}

	if msg.Channel[0] == 'D' { // TODO better way to detect if this is our bot DM?
		incoming.Type = message.DirectMessageMessageType

		incoming = slackutil.ParseMessageForAnyChannelReference(incoming)

		return incoming, true, nil
	} else if !slackutil.CheckMessageForMention(incoming, p.isSelf) {
		return message.Message{}, false, nil
	}

	incoming = slackutil.ParseMessageForChannelReference(incoming, p.isSelf)

	return incoming, true, nil
}

func (p *Parser) isSelf(userID string) bool {
	return p.userID == userID
}
