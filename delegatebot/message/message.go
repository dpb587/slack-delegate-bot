package message

import (
	"time"
)

type OriginType string

const ChannelOriginType OriginType = "channel"
const DirectMessageOriginType OriginType = "dm"

type Message struct {
	OriginType      OriginType
	Origin          string
	InterruptTarget string
	Timestamp       time.Time
	Text            string
}
