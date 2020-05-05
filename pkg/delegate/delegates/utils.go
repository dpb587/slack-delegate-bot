package delegates

import (
	"strings"

	"github.com/dpb587/slack-delegate-bot/pkg/message"
)

func Join(delegates []message.Delegate, sep string) string {
	var str []string

	for _, i := range delegates {
		str = append(str, i.String())
	}

	return strings.Join(str, sep)
}
