package defaultfactory

import (
	"fmt"

	"github.com/dpb587/slack-alias-bot/conditions"
	"github.com/dpb587/slack-alias-bot/interrupt"
	"github.com/dpb587/slack-alias-bot/interrupts"
	coalescefactory "github.com/dpb587/slack-alias-bot/interrupts/coalesce/factory"
	conditionalfactory "github.com/dpb587/slack-alias-bot/interrupts/conditional/factory"
	literalfactory "github.com/dpb587/slack-alias-bot/interrupts/literal/factory"
	pairistfactory "github.com/dpb587/slack-alias-bot/interrupts/pairist/factory"
	unionfactory "github.com/dpb587/slack-alias-bot/interrupts/union/factory"
	userfactory "github.com/dpb587/slack-alias-bot/interrupts/user/factory"
	usergroupfactory "github.com/dpb587/slack-alias-bot/interrupts/usergroup/factory"
)

type factory struct {
	factory map[string]interrupts.Factory
}

var _ interrupts.Factory = &factory{}

func New(conditionsFactory conditions.Factory) interrupts.Factory {
	f := &factory{
		factory: map[string]interrupts.Factory{},
	}

	f.factory["coalesce"] = coalescefactory.New(f)
	f.factory["if"] = conditionalfactory.New(f, conditionsFactory)
	f.factory["literal"] = literalfactory.New()
	f.factory["pairist"] = pairistfactory.New()
	f.factory["union"] = unionfactory.New(f)
	f.factory["user"] = userfactory.New()
	f.factory["usergroup"] = usergroupfactory.New()

	return f
}

func (f *factory) Create(name string, options interface{}) (interrupt.Interrupt, error) {
	ff, known := f.factory[name]
	if !known {
		return nil, fmt.Errorf("unsupported interrupt: %s", name)
	}

	return ff.Create(name, options)
}
