package slack

import "time"

type Processor interface {
	Process(since time.Time, event string, payload []byte) error
}
