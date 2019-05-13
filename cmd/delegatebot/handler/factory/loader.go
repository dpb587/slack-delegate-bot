package factory

import (
	"io/ioutil"
	"path/filepath"
	"sort"

	"github.com/dpb587/slack-delegate-bot/pkg/config"
	"github.com/dpb587/slack-delegate-bot/cmd/delegatebot/handler"
	"github.com/dpb587/slack-delegate-bot/cmd/delegatebot/handler/handlers/multiple"
	"github.com/dpb587/slack-delegate-bot/cmd/delegatebot/handler/handlers/single"
	"github.com/dpb587/slack-delegate-bot/pkg/condition/conditions"
	"github.com/dpb587/slack-delegate-bot/pkg/delegate/delegates"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

type FileLoader struct {
	delegatorsFactory delegates.Factory
	conditionsFactory conditions.Factory
}

func NewFileLoader(delegatorsFactory delegates.Factory, conditionsFactory conditions.Factory) *FileLoader {
	return &FileLoader{
		delegatorsFactory: delegatorsFactory,
		conditionsFactory: conditionsFactory,
	}
}

type ConfigFile struct {
	DelegateBot ConfigFileDelegateBot `yaml:"delegatebot"`
}

type ConfigFileDelegateBot struct {
	Watch    []interface{}                    `yaml:"watch"`
	Delegate interface{}                      `yaml:"delegate"`
	Options  ConfigFileDelegateBotWithOptions `yaml:"options"`
}

type ConfigFileDelegateBotWithOptions struct {
	EmptyMessage string `yaml:"empty_message"`
}

func (l FileLoader) Load(paths []string) (handler.Handler, error) {
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

func (l FileLoader) loadPath(path string) (handler.Handler, error) {
	h := single.Handler{}

	pathBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "reading")
	}

	var parsed ConfigFile

	err = yaml.Unmarshal(pathBytes, &parsed)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshalling")
	}

	if parsed.DelegateBot.Watch != nil {
		watch, err := l.conditionsFactory.Create("or", parsed.DelegateBot.Watch)
		if err != nil {
			return nil, errors.Wrap(err, "building watch")
		}

		h.Condition = watch
	}

	delegateKey, delegateOptions, err := config.KeyValueTuple(parsed.DelegateBot.Delegate)
	if err != nil {
		return nil, errors.Wrap(err, "parsing delegate")
	}

	delegate, err := l.delegatorsFactory.Create(delegateKey, delegateOptions)
	if err != nil {
		return nil, errors.Wrap(err, "building delegate")
	}

	h.Delegator = delegate

	h.Options = handler.Options{
		EmptyMessage: parsed.DelegateBot.Options.EmptyMessage,
	}

	return h, nil
}

func (l FileLoader) squashPaths(paths []string) ([]string, error) {
	var squashed []string

	for _, path := range paths {
		globbed, err := filepath.Glob(path)
		if err != nil {
			return nil, errors.Wrapf(err, "globbing %s", path)
		}

		sort.Strings(globbed)

		squashed = append(squashed, globbed...)
	}

	return squashed, nil
}
