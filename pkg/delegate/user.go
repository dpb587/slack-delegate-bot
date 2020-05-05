package delegate

import (
	"fmt"

	"github.com/dpb587/slack-delegate-bot/pkg/message"
)

type User struct {
	ID string
}

var _ message.Delegate = &User{}

func (i User) String() string {
	return fmt.Sprintf("<@%s>", i.ID)
}
