package slackutil

import (
	"regexp"

	"github.com/dpb587/slack-delegate-bot/pkg/message"
)

var reMention = regexp.MustCompile(`<@([^>]+)>`)
var reChannel = regexp.MustCompile(`<#([^|>]+)(\|([^>]*))?>`)
var reChannelMention = regexp.MustCompile(`<#([^|>]+)(\|([^>]*))?>\s+<@([^>]+)>`)

func ParseMessageForAnyChannelReference(msg message.Message) message.Message {
	match := reChannel.FindStringSubmatch(msg.RawText)
	if len(match) == 0 {
		return msg
	}

	msg.TargetChannelID = match[1]

	return msg
}

func ParseMessageForChannelReference(msg message.Message, isSelf func(string) bool) message.Message {
	for _, match := range reChannelMention.FindAllStringSubmatch(msg.RawText, 32) {
		if isSelf(match[4]) {
			msg.TargetChannelID = match[1]

			break
		}
	}

	return msg
}

func CheckMessageForMention(msg message.Message, isSelf func(string) bool) bool {
	for _, match := range reMention.FindAllStringSubmatch(msg.RawText, 32) {
		if isSelf(match[1]) {
			return true
		}
	}

	return false
}
