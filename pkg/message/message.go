package message

import (
	"time"
)

type OriginType string

// TODO rename to MentionOriginType
const ChannelOriginType OriginType = "channel"
const DirectMessageOriginType OriginType = "dm"

type Message struct {
	ServiceAPI interface{}

	TeamID                string
	OriginUserID          string
	OriginType            OriginType
	Origin                string
	OriginTimestamp       string
	OriginThreadTimestamp string
	InterruptTarget       string
	Timestamp             time.Time
	Text                  string
}
