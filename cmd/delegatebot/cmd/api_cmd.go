package cmd

import (
	nethttp "net/http"
	"time"

	"github.com/dpb587/slack-delegate-bot/cmd/delegatebot/opts"
	zlhttp "github.com/dpb587/slack-delegate-bot/pkg/http"
	"github.com/dpb587/slack-delegate-bot/pkg/slack"
	slackhttp "github.com/dpb587/slack-delegate-bot/pkg/slack/http"
	"github.com/sirupsen/logrus"
	slackapi "github.com/slack-go/slack"
)

type APICmd struct {
	*opts.Root `no-flags:"true"`

	SlackToken         string `long:"slack-token" description:"Slack Bot OAuth API token" env:"SLACK_TOKEN"`
	SlackSigningSecret string `long:"slack-signing-secret" description:"Slack App Signing Secret" env:"SLACK_SIGNING_SECRET"`
}

func (c *APICmd) Execute(_ []string) error {
	http := &nethttp.Server{
		Addr:         ":1234",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	s := zlhttp.NewServer(http)

	services := []zlhttp.Service{
		zlhttp.MetaRuntimeHandler{},
		c.slackService(),
	}

	return s.Run(services...)
}

func (c *APICmd) slackService() zlhttp.Service {
	var apiOpts []slackapi.Option

	if logrus.Level(c.Root.LogLevel) == logrus.DebugLevel {
		apiOpts = append(
			apiOpts,
			slackapi.OptionDebug(true),
			// slackapi.OptionLog(c.Root.GetLogger()),
		)
	}

	api := slackapi.New(c.SlackToken, apiOpts...)

	h, err := c.GetHandler()
	if err != nil {
		// TODO
		panic(err)
	}

	return &slackhttp.Service{
		Messenger:     slack.NewMessenger(api, h),
		SigningSecret: c.SlackSigningSecret,
	}
}
