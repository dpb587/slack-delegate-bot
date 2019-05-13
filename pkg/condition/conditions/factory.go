package conditions

import "github.com/dpb587/slack-delegate-bot/pkg/condition"

type Factory interface {
	Create(name string, config interface{}) (condition.Condition, error)
}
