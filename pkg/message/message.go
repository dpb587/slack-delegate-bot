package message

import (
	"time"
)

type MessageType string

// TODO rename to MentionMessageType
const ChannelMessageType MessageType = "channel"
const DirectMessageMessageType MessageType = "dm"

type Message struct {
	// TODO these in separate context?
	ServiceAPI     interface{}
	Delegator      interface{}
	RecursionDepth int

	UserTeamID string
	UserID     string

	ChannelTeamID string
	ChannelID     string

	TargetChannelTeamID string
	TargetChannelID     string

	RawText            string
	RawTimestamp       string
	RawThreadTimestamp string

	Time time.Time
	Type MessageType
}
