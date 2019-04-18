package main

import (
	"os"

	"github.com/dpb587/slack-delegate-bot/main/cmd"
	"github.com/dpb587/slack-delegate-bot/main/opts"
	"github.com/jessevdk/go-flags"
)

func main() {
	opts := &opts.Root{}
	main := cmd.Root{
		Root: opts,
		Run: &cmd.RunCmd{
			Root: opts,
		},
		Validate: &cmd.ValidateCmd{
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
