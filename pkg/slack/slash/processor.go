package slash

import "time"

type Processor interface {
	Process(since time.Time, event string, payload []byte) error
}
