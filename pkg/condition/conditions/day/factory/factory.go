package factory

import (
	"fmt"
	"time"

	"github.com/dpb587/slack-delegate-bot/pkg/condition"
	"github.com/dpb587/slack-delegate-bot/pkg/condition/conditions"
	"github.com/dpb587/slack-delegate-bot/pkg/condition/conditions/day"
	"github.com/dpb587/slack-delegate-bot/pkg/configutil"
	"github.com/pkg/errors"
)

type factory struct{}

type Options struct {
	TZ   string   `yaml:"tz"`
	Days []string `yaml:"days"`
}

func New() conditions.Factory {
	return &factory{}
}

func (f factory) Create(name string, options interface{}) (condition.Condition, error) {
	if name != "day" {
		return nil, fmt.Errorf("invalid condition: %s", name)
	}

	parsed := Options{}

	err := configutil.RemarshalYAML(options, &parsed)
	if err != nil {
		return nil, errors.Wrap(err, "remarshalling")
	}

	loc, err := time.LoadLocation(parsed.TZ)
	if err != nil {
		return nil, errors.Wrap(err, "loading timezone")
	}

	return &day.Condition{
		Location: loc,
		Days:     parsed.Days,
	}, nil
}
