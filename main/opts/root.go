package opts

import (
	"os"

	conditionsfactory "github.com/dpb587/slack-delegate-bot/condition/conditions/defaultfactory"
	"github.com/dpb587/slack-delegate-bot/handler"
	"github.com/dpb587/slack-delegate-bot/handler/fileloader"
	interruptsfactory "github.com/dpb587/slack-delegate-bot/interrupt/interrupts/defaultfactory"
	"github.com/dpb587/slack-delegate-bot/main/args"
	"github.com/dpb587/slack-delegate-bot/service/slack"
	slackapi "github.com/nlopes/slack"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Root struct {
	SlackToken   string `long:"slack-token" description:"Slack API Token" env:"SLACK_TOKEN"`
	slackAPI     *slackapi.Client
	slackService *slack.Service

	Handlers []string `long:"handler" description:"Path to handler configuration"`
	handler  handler.Handler

	LogLevel args.LogLevel `long:"log-level" description:"Show additional levels of log messages" env:"LOG_LEVEL" default:"ERROR"`
	logger   logrus.FieldLogger
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

func (r *Root) GetHandler() (handler.Handler, error) {
	if r.handler == nil {
		conditions := conditionsfactory.New()
		interrupts := interruptsfactory.New(conditions, r.GetSlackAPI())

		loader := fileloader.New(interrupts, conditions)
		h, err := loader.Load(r.Handlers)
		if err != nil {
			return nil, errors.Wrap(err, "loading handlers")
		}

		r.handler = h
	}

	return r.handler, nil
}

func (r *Root) GetSlackAPI() *slackapi.Client {
	if r.slackAPI == nil {
		r.slackAPI = slackapi.New(r.SlackToken) // , slack.OptionDebug(true)) // TODO , slack.OptionLog(s.logger))
	}

	return r.slackAPI
}

func (r *Root) GetSlackService() (*slack.Service, error) {
	if r.slackService == nil {
		handler, err := r.GetHandler()
		if err != nil {
			return nil, errors.Wrap(err, "getting handler")
		}

		r.slackService = slack.New(r.GetSlackAPI(), handler, r.GetLogger().WithField("service", "slack"))
	}

	return r.slackService, nil
}
