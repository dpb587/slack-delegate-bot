package cmd

import (
	"github.com/dpb587/slack-delegate-bot/cmd/delegatebot/opts"
	"github.com/dpb587/slack-delegate-bot/pkg/slack"
	"github.com/dpb587/slack-delegate-bot/pkg/slack/rtm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	slackapi "github.com/slack-go/slack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type RTMCmd struct {
	*opts.Root `no-flags:"true"`

	SlackLogLevel zapcore.Level `long:"slack-log-level" description:"Log level for Slack client" env:"SLACK_LOG_LEVEL"`
	SlackToken    string        `long:"slack-token" description:"Slack Bot OAuth API token" env:"SLACK_TOKEN"`
}

func (c *RTMCmd) Execute(_ []string) error {
	api := c.slackAPI()

	h, err := c.GetDelegator()
	if err != nil {
		// TODO
		panic(err)
	}

	p := rtm.NewService(api, slack.NewResponder(api, h), c.GetLogger())

	return p.Run()
}

func (c *RTMCmd) slackAPI() *slackapi.Client {
	var apiOpts []slackapi.Option

	if c.SlackLogLevel == zapcore.DebugLevel {
		ll, _ := zap.NewStdLogAt(c.Root.GetLogger(), zapcore.DebugLevel)

		apiOpts = append(
			apiOpts,
			slackapi.OptionDebug(true),
			slackapi.OptionLog(ll),
		)
	}

	return slackapi.New(c.SlackToken, apiOpts...)
}
