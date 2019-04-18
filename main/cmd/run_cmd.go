package cmd

import (
	"github.com/dpb587/slack-delegate-bot/main/opts"
	"github.com/pkg/errors"
)

type RunCmd struct {
	*opts.Root `no-flags:"true"`
}

func (c *RunCmd) Execute(_ []string) error {
	svc, err := c.Root.GetSlackService()
	if err != nil {
		return errors.Wrap(err, "loading slack service")
	}

	return svc.Run()
}
