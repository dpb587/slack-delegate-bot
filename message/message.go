package message

import "time"

type OriginType string

const ChannelOriginType OriginType = "channel"
const DirectMessageOriginType OriginType = "dm"

type Message struct {
	OriginType      OriginType
	Origin          string
	InterruptTarget string
	Timestamp       time.Time
	Text            string

	response *Response
}

func (m *Message) SetResponse(response Response) {
	m.response = &response
}

func (m *Message) GetResponse() *Response {
	return m.response
}

type Response struct {
	Text string
}
