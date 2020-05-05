package delegate

import "github.com/dpb587/slack-delegate-bot/pkg/message"

type Literal struct {
	Text string
}

var _ message.Delegate = &Literal{}

func (i Literal) String() string {
	return i.Text
}
