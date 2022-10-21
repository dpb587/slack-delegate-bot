package cmd

import (
	"fmt"
	nethttp "net/http"
	"time"

	"github.com/dpb587/slack-delegate-bot/cmd/delegatebot/opts"
	zlhttp "github.com/dpb587/slack-delegate-bot/pkg/http"
	"github.com/dpb587/slack-delegate-bot/pkg/slack"
	slackevent "github.com/dpb587/slack-delegate-bot/pkg/slack/event"
	slackhttp "github.com/dpb587/slack-delegate-bot/pkg/slack/http"
	"github.com/dpb587/slack-delegate-bot/pkg/slack/slash"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	slackapi "github.com/slack-go/slack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type APICmd struct {
	*opts.Root `no-flags:"true"`

	BindHost    string `long:"bind-host" description:"Bind host/IP" env:"BINDING" default:"0.0.0.0"`
	BindPort    int    `long:"bind-port" description:"Bind port" env:"PORT" default:"8080"`
	ExternalURL string `long:"external-url" description:"Public URL" env:"HTTP_EXTERNAL_URL"`

	SlackToken         string `long:"slack-token" description:"Slack Bot OAuth API token" env:"SLACK_TOKEN"`
	SlackSigningSecret string `long:"slack-signing-secret" description:"Slack App Signing Secret" env:"SLACK_SIGNING_SECRET"`
}

func (c *APICmd) Execute(_ []string) error {
	http := &nethttp.Server{
		Addr:         fmt.Sprintf("%s:%d", c.BindHost, c.BindPort),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	s := zlhttp.NewServer(http, c.GetLogger())

	services := []zlhttp.Service{
		zlhttp.MetaRuntimeHandler{},
		c.slackService(),
	}

	return s.Run(services...)
}

func (c *APICmd) slackService() zlhttp.Service {
	var apiOpts []slackapi.Option

	if c.Root.LogLevel == zapcore.DebugLevel {
		ll, _ := zap.NewStdLogAt(c.Root.GetLogger(), zapcore.DebugLevel)

		apiOpts = append(
			apiOpts,
			slackapi.OptionDebug(true),
			slackapi.OptionLog(ll),
		)
	}

	api := slackapi.New(c.SlackToken, apiOpts...)

	h, err := c.GetDelegator()
	if err != nil {
		// TODO
		panic(err)
	}

	processor := slackevent.NewSyncProcessor(
		slackevent.NewParser(slack.NewUserLookup(api)),
		slack.NewResponder(api, h),
	)

	slashHandler := slash.Handlers{
		slash.NewShowHandler(h, api),
	}

	slashHandler = append(slashHandler, slash.NewHelpHandler(
		&slashHandler,
		c.ExternalURL,
	))

	slashProcessor := slash.NewSyncProcessor(slashHandler)

	return &slackhttp.Service{
		EventProcessor: processor,
		SlashProcessor: slashProcessor,
		SigningSecret:  c.SlackSigningSecret,
		Logger:         c.GetLogger(),
	}
}
