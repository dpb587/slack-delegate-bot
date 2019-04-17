package fileloader

import (
	"io/ioutil"
	"path/filepath"

	"github.com/dpb587/slack-alias-bot/conditions"
	"github.com/dpb587/slack-alias-bot/handler"
	"github.com/dpb587/slack-alias-bot/handlers/multiple"
	"github.com/dpb587/slack-alias-bot/handlers/single"
	"github.com/dpb587/slack-alias-bot/interrupts"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

type Loader struct {
	interruptsFactory interrupts.Factory
	conditionsFactory conditions.Factory
}

func New(interruptsFactory interrupts.Factory, conditionsFactory conditions.Factory) *Loader {
	return &Loader{
		interruptsFactory: interruptsFactory,
		conditionsFactory: conditionsFactory,
	}
}

type Options struct {
	When []interface{} `yaml:"when"`
	Then []interface{} `yaml:"then"`
	With WithOptions   `yaml:"with"`
}

type WithOptions struct {
	EmptyResponse    string `yaml:"empty_response"`
	ResponseTemplate string `yaml:"response_template"`
}

func (l Loader) Load(paths []string) (handler.Handler, error) {
	var handlers []handler.Handler

	paths, err := l.squashPaths(paths)
	if err != nil {
		return nil, errors.Wrap(err, "preparing paths")
	}

	for _, path := range paths {
		h, err := l.loadPath(path)
		if err != nil {
			return nil, errors.Wrapf(err, "loading %s", path)
		}

		handlers = append(handlers, h)
	}

	return multiple.Handler{Handlers: handlers}, nil
}

func (l Loader) loadPath(path string) (handler.Handler, error) {
	h := single.Handler{}

	pathBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "reading")
	}

	var parsed Options

	err = yaml.Unmarshal(pathBytes, &parsed)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshalling")
	}

	if parsed.When != nil {
		when, err := l.conditionsFactory.Create("and", parsed.When)
		if err != nil {
			return nil, errors.Wrap(err, "building when")
		}

		h.Condition = when
	}

	then, err := l.interruptsFactory.Create("union", parsed.Then)
	if err != nil {
		return nil, errors.Wrap(err, "building then")
	}

	h.Interrupt = then

	h.Options = handler.Options{
		EmptyResponse:    parsed.With.EmptyResponse,
		ResponseTemplate: parsed.With.ResponseTemplate,
	}

	return h, nil
}

func (l Loader) squashPaths(paths []string) ([]string, error) {
	var squashed []string

	for _, path := range paths {
		globbed, err := filepath.Glob(path)
		if err != nil {
			return nil, errors.Wrapf(err, "globbing %s", path)
		}

		squashed = append(squashed, globbed...)
	}

	return squashed, nil
}
