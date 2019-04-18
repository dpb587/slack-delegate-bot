package factory

import (
	"fmt"
	"time"

	"github.com/dpb587/slack-delegate-bot/condition"
	"github.com/dpb587/slack-delegate-bot/condition/conditions"
	"github.com/dpb587/slack-delegate-bot/condition/conditions/date"
	"github.com/dpb587/slack-delegate-bot/config"
	"github.com/pkg/errors"
)

type factory struct{}

type Options struct {
	TZ    string   `yaml:"tz"`
	Dates []string `yaml:"dates"`
}

func New() conditions.Factory {
	return &factory{}
}

func (f factory) Create(name string, options interface{}) (condition.Condition, error) {
	if name != "date" {
		return nil, fmt.Errorf("invalid condition: %s", name)
	}

	parsed := Options{}

	err := config.RemarshalYAML(options, &parsed)
	if err != nil {
		return nil, errors.Wrap(err, "remarshalling")
	}

	loc, err := time.LoadLocation(parsed.TZ)
	if err != nil {
		return nil, errors.Wrap(err, "loading timezone")
	}

	return &date.Condition{
		Location: loc,
		Dates:    parsed.Dates,
	}, nil
}
