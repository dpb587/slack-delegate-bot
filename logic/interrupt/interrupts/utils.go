package interrupts

import (
	"strings"

	"github.com/dpb587/slack-delegate-bot/logic/interrupt"
)

func Join(interrupts []interrupt.Interruptible, sep string) string {
	var str []string

	for _, i := range interrupts {
		str = append(str, i.String())
	}

	return strings.Join(str, sep)
}
