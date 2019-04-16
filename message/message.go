package message

import "time"

type Message struct {
	Origin          string
	InterruptTarget string
	Timestamp       time.Time
	Text            string
}
