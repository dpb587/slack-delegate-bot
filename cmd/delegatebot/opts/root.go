package opts

import (
	"os"

	"github.com/dpb587/slack-delegate-bot/cmd/delegatebot/handler"
	handlerfactory "github.com/dpb587/slack-delegate-bot/cmd/delegatebot/handler/factory"
	"github.com/dpb587/slack-delegate-bot/cmd/delegatebot/args"
	"github.com/dpb587/slack-delegate-bot/cmd/delegatebot/service/http"
	"github.com/dpb587/slack-delegate-bot/cmd/delegatebot/service/slack"
	conditionsfactory "github.com/dpb587/slack-delegate-bot/pkg/condition/conditions/defaultfactory"
	interruptsfactory "github.com/dpb587/slack-delegate-bot/pkg/delegate/delegates/defaultfactory"
	slackapi "github.com/nlopes/slack"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Root struct {
	SlackToken   string `long:"slack-token" description:"Slack API Token" env:"SLACK_TOKEN"`
	slackAPI     *slackapi.Client
	slackService *slack.Service

	httpService *http.Service

	Configs []string `long:"config" description:"Path to configuration files"`
	handler handler.Handler

	LogLevel args.LogLevel `long:"log-level" description:"Show additional levels of log messages" env:"LOG_LEVEL" default:"INFO"`
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

		loader := handlerfactory.NewFileLoader(interrupts, conditions)
		h, err := loader.Load(r.Configs)
		if err != nil {
			return nil, errors.Wrap(err, "loading configs")
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

func (r *Root) GetHTTPService() (*http.Service, error) {
	if r.httpService == nil {
		r.httpService = http.New(r.GetLogger().WithField("service", "http"))
	}

	return r.httpService, nil
}
