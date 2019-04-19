package literalmap

import (
	"github.com/dpb587/slack-delegate-bot/logic/interrupt"
	"github.com/dpb587/slack-delegate-bot/message"
)

type Interrupt struct {
	From       interrupt.Interrupt
	Users      map[string]string
	Usergroups map[string]string
}

var _ interrupt.Interrupt = &Interrupt{}

func (i Interrupt) Lookup(m message.Message) ([]interrupt.Interruptible, error) {
	inner, err := i.From.Lookup(m)
	if err != nil {
		return nil, err
	}

	var res []interrupt.Interruptible

	for _, innerInterrupt := range inner {
		literalInterrupt, ok := innerInterrupt.(interrupt.Literal)
		if !ok {
			res = append(res, innerInterrupt)

			continue
		}

		var newres interrupt.Interruptible = literalInterrupt

		if mapped, found := i.Users[literalInterrupt.Text]; found {
			newres = interrupt.User{ID: mapped}
		} else if mapped, found := i.Usergroups[literalInterrupt.Text]; found {
			newres = interrupt.UserGroup{ID: mapped}
		}

		res = append(res, newres)
	}

	return res, nil
}
