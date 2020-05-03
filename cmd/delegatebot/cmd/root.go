package cmd

import "github.com/dpb587/slack-delegate-bot/cmd/delegatebot/opts"

type Root struct {
	*opts.Root

	Validate *ValidateCmd `command:"validate" description:"For validating configuration"`
	Simulate *SimulateCmd `command:"simulate" description:"For simulating an incoming message"`
	API      *APICmd      `command:"api" description:"Run HTTP API server"`
}
