package main

import (
	"os"

	"github.com/dpb587/slack-delegate-bot/cmd/delegatebot/cmd"
	"github.com/dpb587/slack-delegate-bot/cmd/delegatebot/opts"
	"github.com/jessevdk/go-flags"
)

func main() {
	opts := &opts.Root{}
	main := cmd.Root{
		Root: opts,
		API: &cmd.APICmd{
			Root: opts,
		},
		Validate: &cmd.ValidateCmd{
			Root: opts,
		},
		Simulate: &cmd.SimulateCmd{
			Root: opts,
		},
	}

	var parser = flags.NewParser(&main, flags.Default)

	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}
}
