package cmd

import "github.com/dpb587/slack-delegate-bot/cmd/delegatebot/opts"

type Root struct {
	*opts.Root

	Run      *RunCmd      `command:"run" description:"For running the bot"`
	Stats    *StatsCmd    `command:"stats" description:"For generating historical stats"`
	Validate *ValidateCmd `command:"validate" description:"For validating configuration"`
}
