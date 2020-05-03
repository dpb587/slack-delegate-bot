package cmd

import (
	"fmt"
	"time"

	"github.com/dpb587/slack-delegate-bot/cmd/delegatebot/message"
	"github.com/dpb587/slack-delegate-bot/cmd/delegatebot/opts"
	"github.com/dpb587/slack-delegate-bot/cmd/delegatebot/service/slack"
	"github.com/dpb587/slack-delegate-bot/pkg/delegate/delegates"
	"github.com/pkg/errors"
	slackapi "github.com/slack-go/slack"
)

type SimulateCmd struct {
	*opts.Root `no-flags:"true"`

	Timestamp string       `long:"timestamp" description:"Timestamp of the message (default: now; format: 2006-01-02T15:04:05Z07:00)"`
	Args      SimulateArgs `positional-args:"true"`
}

type SimulateArgs struct {
	Origin  string `positional-arg-name:"ORIGIN-ID" description:"Channel or DM ID sending the request" required:"true"`
	Message string `positional-arg-name:"MESSAGE" description:"Message sent"`
}

func (c *SimulateCmd) Execute(_ []string) error {
	if c.Timestamp == "" {
		c.Timestamp = time.Now().Format(time.RFC3339)
	}

	ts, err := time.Parse(time.RFC3339, c.Timestamp)
	if err != nil {
		return errors.Wrap(err, "parsing RFC3339 timestamp")
	}

	if c.Args.Message == "" {
		c.Args.Message = "<@U0000000>"
	}

	handler, err := c.Root.GetHandler()
	if err != nil {
		return err
	}

	parser := slack.NewMessageParser(&slackapi.UserDetails{
		ID: "U0000000",
	})

	request, err := parser.ParseMessage(slackapi.Msg{
		Type:      "message",
		User:      "test",
		Channel:   c.Args.Origin,
		Text:      c.Args.Message,
		Timestamp: fmt.Sprintf("%d.0", ts.Unix()),
	})
	if err != nil {
		return errors.Wrap(err, "parsing fake message")
	}

	response, err := handler.Execute(request)
	if err != nil {
		return errors.Wrap(err, "evaluating a response")
	}

	var msg string

	if len(response.Delegates) > 0 {
		msg = delegates.Join(response.Delegates, " ")

		if request.OriginType == message.ChannelOriginType {
			msg = fmt.Sprintf("^ %s", msg)
		}
	} else if response.EmptyMessage != "" {
		msg = response.EmptyMessage
	}

	if msg != "" {
		fmt.Printf("%s\n", msg)
	}

	return nil
}
