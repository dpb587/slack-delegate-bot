package conditions

import "github.com/dpb587/slack-alias-bot/condition"

type Factory interface {
	Create(name string, config interface{}) (condition.Condition, error)
}
