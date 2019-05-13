package defaultfactory

import (
	"fmt"

	"github.com/dpb587/slack-delegate-bot/pkg/condition"
	"github.com/dpb587/slack-delegate-bot/pkg/condition/conditions"
	boolandfactory "github.com/dpb587/slack-delegate-bot/pkg/condition/conditions/booland/factory"
	boolnotfactory "github.com/dpb587/slack-delegate-bot/pkg/condition/conditions/boolnot/factory"
	boolorfactory "github.com/dpb587/slack-delegate-bot/pkg/condition/conditions/boolor/factory"
	datefactory "github.com/dpb587/slack-delegate-bot/pkg/condition/conditions/date/factory"
	dayfactory "github.com/dpb587/slack-delegate-bot/pkg/condition/conditions/day/factory"
	hoursfactory "github.com/dpb587/slack-delegate-bot/pkg/condition/conditions/hours/factory"
	targetfactory "github.com/dpb587/slack-delegate-bot/pkg/condition/conditions/target/factory"
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

	return f
}

func (f *factory) Create(name string, options interface{}) (condition.Condition, error) {
	ff, known := f.factory[name]
	if !known {
		return nil, fmt.Errorf("unsupported condition: %s", name)
	}

	return ff.Create(name, options)
}
