package factory

import (
	"fmt"
	"io/ioutil"

	"github.com/dpb587/slack-delegate-bot/config"
	"github.com/dpb587/slack-delegate-bot/logic/condition"
	"github.com/dpb587/slack-delegate-bot/logic/condition/conditions"
	"github.com/dpb587/slack-delegate-bot/logic/condition/conditions/boolor"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

type factory struct {
	factory conditions.Factory
}

type Options struct {
	Source string `yaml:"source"`
	Path   string `yaml:"path"`
}

func New(ff conditions.Factory) conditions.Factory {
	return &factory{
		factory: ff,
	}
}

func (f factory) Create(name string, options interface{}) (condition.Condition, error) {
	if name != "embed" {
		return nil, fmt.Errorf("invalid condition: %s", name)
	}

	var parsed Options

	err := config.RemarshalYAML(options, &parsed)
	if err != nil {
		return nil, errors.Wrap(err, "remarshalling")
	}

	bytes, err := ioutil.ReadFile(parsed.Source)
	if err != nil {
		return nil, errors.Wrap(err, "loading embedded file")
	}

	var embedMap map[string]interface{}

	err = yaml.Unmarshal(bytes, &embedMap)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshalling embedded file")
	}

	parsedEmbed, found := embedMap[parsed.Path]
	if !found {
		return nil, fmt.Errorf("unable to find path in embedded file: %s", parsed.Path)
	}

	key, value, err := config.KeyValueTuple(parsedEmbed)
	if err != nil {
		return nil, errors.Wrapf(err, "parsing condition %d", parsedEmbed)
	}

	var ccds []condition.Condition

	condition, err := f.factory.Create(key, value)
	if err != nil {
		return nil, errors.Wrapf(err, "creating condition %d", parsedEmbed)
	}

	ccds = append(ccds, condition)

	return &boolor.Condition{
		Conditions: ccds,
	}, nil
}
