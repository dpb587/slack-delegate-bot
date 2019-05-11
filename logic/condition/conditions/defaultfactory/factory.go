package defaultfactory

import (
	"fmt"

	"github.com/dpb587/slack-delegate-bot/logic/condition"
	"github.com/dpb587/slack-delegate-bot/logic/condition/conditions"
	boolandfactory "github.com/dpb587/slack-delegate-bot/logic/condition/conditions/booland/factory"
	boolnotfactory "github.com/dpb587/slack-delegate-bot/logic/condition/conditions/boolnot/factory"
	boolorfactory "github.com/dpb587/slack-delegate-bot/logic/condition/conditions/boolor/factory"
	datefactory "github.com/dpb587/slack-delegate-bot/logic/condition/conditions/date/factory"
	dayfactory "github.com/dpb587/slack-delegate-bot/logic/condition/conditions/day/factory"
	embedfactory "github.com/dpb587/slack-delegate-bot/logic/condition/conditions/embed/factory"
	hoursfactory "github.com/dpb587/slack-delegate-bot/logic/condition/conditions/hours/factory"
	targetfactory "github.com/dpb587/slack-delegate-bot/logic/condition/conditions/target/factory"
)

type factory struct {
	factory map[string]conditions.Factory
}

var _ conditions.Factory = &factory{}

func New() conditions.Factory {
	f := &factory{
		factory: map[string]conditions.Factory{},
	}

	f.factory["and"] = boolandfactory.New(f)
	f.factory["not"] = boolnotfactory.New(f)
	f.factory["or"] = boolorfactory.New(f)
	f.factory["date"] = datefactory.New()
	f.factory["day"] = dayfactory.New()
	f.factory["hours"] = hoursfactory.New()
	f.factory["target"] = targetfactory.New()
	f.factory["embed"] = embedfactory.New(f)

	return f
}

func (f *factory) Create(name string, options interface{}) (condition.Condition, error) {
	ff, known := f.factory[name]
	if !known {
		return nil, fmt.Errorf("unsupported condition: %s", name)
	}

	return ff.Create(name, options)
}
