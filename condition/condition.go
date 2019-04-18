package condition

import "github.com/dpb587/slack-delegate-bot/message"

type Condition interface {
	Evaluate(message.Message) (bool, error)
}
