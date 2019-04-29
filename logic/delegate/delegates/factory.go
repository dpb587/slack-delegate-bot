package delegates

import "github.com/dpb587/slack-delegate-bot/logic/delegate"

type Factory interface {
	Create(name string, options interface{}) (delegate.Delegator, error)
}
