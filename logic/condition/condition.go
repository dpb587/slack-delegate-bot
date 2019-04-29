package condition

import "github.com/dpb587/slack-delegate-bot/delegatebot/message"

//go:generate counterfeiter . Condition
type Condition interface {
	Evaluate(message.Message) (bool, error)
}
