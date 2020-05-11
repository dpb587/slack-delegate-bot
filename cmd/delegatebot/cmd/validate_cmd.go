package cmd

import (
	"github.com/dpb587/slack-delegate-bot/cmd/delegatebot/opts"
)

type ValidateCmd struct {
	*opts.Root `no-flags:"true"`
}

func (c *ValidateCmd) Execute(_ []string) error {
	_, err := c.Root.GetDelegator()

	return err
}
