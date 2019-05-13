package defaultfactory

import (
	"fmt"

	"github.com/dpb587/slack-delegate-bot/pkg/condition/conditions"
	"github.com/dpb587/slack-delegate-bot/pkg/delegate"
	"github.com/dpb587/slack-delegate-bot/pkg/delegate/delegates"
	coalescefactory "github.com/dpb587/slack-delegate-bot/pkg/delegate/delegates/coalesce/factory"
	conditionalfactory "github.com/dpb587/slack-delegate-bot/pkg/delegate/delegates/conditional/factory"
	literalfactory "github.com/dpb587/slack-delegate-bot/pkg/delegate/delegates/literal/factory"
	literalmapfactory "github.com/dpb587/slack-delegate-bot/pkg/delegate/delegates/literalmap/factory"
	pairistfactory "github.com/dpb587/slack-delegate-bot/pkg/delegate/delegates/pairist/factory"
	topiclookupfactory "github.com/dpb587/slack-delegate-bot/pkg/delegate/delegates/topiclookup/factory"
	unionfactory "github.com/dpb587/slack-delegate-bot/pkg/delegate/delegates/union/factory"
	userfactory "github.com/dpb587/slack-delegate-bot/pkg/delegate/delegates/user/factory"
	usergroupfactory "github.com/dpb587/slack-delegate-bot/pkg/delegate/delegates/usergroup/factory"
	"github.com/nlopes/slack"
)

type factory struct {
	factory map[string]delegates.Factory
}

var _ delegates.Factory = &factory{}

func New(conditionsFactory conditions.Factory, slackAPI *slack.Client) delegates.Factory {
	f := &factory{
		factory: map[string]delegates.Factory{},
	}

	f.factory["coalesce"] = coalescefactory.New(f)
	f.factory["if"] = conditionalfactory.New(f, conditionsFactory)
	f.factory["literal"] = literalfactory.New()
	f.factory["literalmap"] = literalmapfactory.New(f)
	f.factory["pairist"] = pairistfactory.New()
	f.factory["topiclookup"] = topiclookupfactory.New(slackAPI)
	f.factory["union"] = unionfactory.New(f)
	f.factory["user"] = userfactory.New()
	f.factory["usergroup"] = usergroupfactory.New()

	return f
}

func (f *factory) Create(name string, options interface{}) (delegate.Delegator, error) {
	ff, known := f.factory[name]
	if !known {
		return nil, fmt.Errorf("unsupported delegate: %s", name)
	}

	return ff.Create(name, options)
}
