package opts

import (
	"os"

	"github.com/dpb587/slack-delegate-bot/cmd/delegatebot/args"
	conditionsfactory "github.com/dpb587/slack-delegate-bot/pkg/condition/conditions/defaultfactory"
	interruptsfactory "github.com/dpb587/slack-delegate-bot/pkg/delegate/delegates/defaultfactory"
	"github.com/dpb587/slack-delegate-bot/pkg/handler"
	"github.com/dpb587/slack-delegate-bot/pkg/handler/fs"
	"github.com/dpb587/slack-delegate-bot/pkg/handler/yaml"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Root struct {
	Configs []string `long:"config" description:"Path to configuration files"`
	handler handler.Handler

	LogLevel args.LogLevel `long:"log-level" description:"Show additional levels of log messages" env:"LOG_LEVEL" default:"INFO"`
	logger   logrus.FieldLogger

	parser *yaml.Parser
}

func (r *Root) GetLogger() logrus.FieldLogger {
	if r.logger == nil {
		var logger = logrus.New()
		logger.Out = os.Stderr
		logger.Formatter = &logrus.JSONFormatter{}
		logger.Level = logrus.Level(r.LogLevel)

		r.logger = logger
	}

	return r.logger
}

func (r *Root) GetParser() *yaml.Parser {
	if r.parser == nil {
		conditions := conditionsfactory.New()
		interrupts := interruptsfactory.New(conditions)

		r.parser = yaml.NewParser(interrupts, conditions)
	}

	return r.parser
}

func (r *Root) GetHandler() (handler.Handler, error) {
	if r.handler == nil {
		h, err := fs.BuildHandler(r.GetParser(), r.Configs...)
		if err != nil {
			return nil, errors.Wrap(err, "loading configs")
		}

		r.handler = h
	}

	return r.handler, nil
}
