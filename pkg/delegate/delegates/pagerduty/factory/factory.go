package factory

import (
	"fmt"
	"os"
	"strings"

	pagerdutyapi "github.com/PagerDuty/go-pagerduty"
	"github.com/dpb587/slack-delegate-bot/pkg/config"
	"github.com/dpb587/slack-delegate-bot/pkg/delegate"
	"github.com/dpb587/slack-delegate-bot/pkg/delegate/delegates"
	"github.com/dpb587/slack-delegate-bot/pkg/delegate/delegates/pagerduty"
	"github.com/pkg/errors"
)

type factory struct{}

type Options struct {
	APIKey           string `yaml:"api_key"`
	EscalationPolicy string `yaml:"escalation_policy"`
	EscalationLevel  *uint  `yaml:"escalation_level"`
}

func New() delegates.Factory {
	return &factory{}
}

func (f factory) Create(name string, options interface{}) (delegate.Delegator, error) {
	if name != "pagerduty" {
		return nil, fmt.Errorf("unsupported delegate: %s", name)
	}

	parsed := Options{}

	err := config.RemarshalYAML(options, &parsed)
	if err != nil {
		return nil, errors.Wrap(err, "remarshalling")
	}

	if parsed.EscalationLevel == nil {
		var r uint = 1
		parsed.EscalationLevel = &r
	}

	apiKey := parsed.APIKey

	if strings.HasPrefix(apiKey, "$") && len(apiKey) > 1 {
		apiKey = os.Getenv(apiKey[1:])
	}

	return &pagerduty.Delegator{
		Client:           pagerdutyapi.NewClient(apiKey),
		EscalationPolicy: parsed.EscalationPolicy,
		EscalationLevel:  *parsed.EscalationLevel,
	}, nil
}
