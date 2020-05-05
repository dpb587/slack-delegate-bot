package delegate

import (
	"fmt"

	"github.com/dpb587/slack-delegate-bot/pkg/message"
)

type UserGroup struct {
	ID    string
	Alias string
}

var _ message.Delegate = &UserGroup{}

func (i UserGroup) String() string {
	return fmt.Sprintf("<!subteam^%s|@%s>", i.ID, i.Alias)
}
