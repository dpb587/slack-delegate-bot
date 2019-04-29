package cmd

import "github.com/dpb587/slack-delegate-bot/delegatebot/main/opts"

type Root struct {
	*opts.Root

	Run      *RunCmd      `command:"run" description:"For running the bot"`
	Validate *ValidateCmd `command:"validate" description:"For validating configuration"`
}
