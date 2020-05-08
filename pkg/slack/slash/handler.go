package slash

import "github.com/slack-go/slack"

type Handler interface {
	Handle(slack.SlashCommand) (bool, error)

	// TODO separate interface
	UsageHint() string
	ShortDescription() string
}
