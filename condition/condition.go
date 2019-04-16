package condition

import "github.com/dpb587/slack-alias-bot/message"

type Condition interface {
	Evaluate(message.Message) (bool, error)
}
