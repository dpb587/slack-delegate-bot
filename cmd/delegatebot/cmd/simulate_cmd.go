package cmd

import (
	"fmt"
	"time"

	"github.com/dpb587/slack-delegate-bot/cmd/delegatebot/opts"
	"github.com/dpb587/slack-delegate-bot/pkg/delegate/delegates"
	"github.com/dpb587/slack-delegate-bot/pkg/message"
	"github.com/pkg/errors"
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

	delegator, err := c.Root.GetDelegator()
	if err != nil {
		return err
	}

	msg := message.Message{
		ChannelTeamID:       "T1234567",
		ChannelID:           "D1234567",
		UserTeamID:          "T1234567",
		UserID:              "U1234567",
		TargetChannelTeamID: "T1234567",
		TargetChannelID:     c.Args.Origin,
		RawTimestamp:        fmt.Sprintf("%d.0", ts.Unix()),
		RawText:             c.Args.Message,
		Type:                message.DirectMessageMessageType,
	}

	dd, err := delegator.Delegate(msg)
	if err != nil {
		return errors.Wrap(err, "evaluating a response")
	}

	if len(dd) == 0 {
		return nil
	}

	fmt.Printf("%s\n", delegates.Join(dd, " "))

	return nil
}
