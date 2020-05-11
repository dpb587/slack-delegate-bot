package lookup

import (
	"fmt"

	"github.com/dpb587/slack-delegate-bot/pkg/delegate"
	"github.com/dpb587/slack-delegate-bot/pkg/message"
	"github.com/pkg/errors"
)

type Delegator struct {
	Channel string
}

var _ delegate.Delegator = &Delegator{}

func (i Delegator) Delegate(msg message.Message) ([]message.Delegate, error) {
	if msg.TargetChannelID == i.Channel {
		// no use
		return nil, nil
	} else if msg.RecursionDepth >= 3 {
		// no more
		return nil, fmt.Errorf("maximum recursion depth reached: %d", msg.RecursionDepth)
	} else if msg.Delegator == nil {
		return nil, errors.New("no delegator available from message context")
	}

	newmsg := message.Message{
		ServiceAPI:     msg.ServiceAPI,
		Delegator:      msg.Delegator,
		RecursionDepth: msg.RecursionDepth + 1, // changed

		UserTeamID:          msg.UserTeamID,
		UserID:              msg.UserID,
		ChannelTeamID:       msg.ChannelTeamID,
		ChannelID:           msg.ChannelID,
		TargetChannelTeamID: msg.TargetChannelTeamID,
		TargetChannelID:     i.Channel, // changed
		RawText:             msg.RawText,
		RawTimestamp:        msg.RawTimestamp,
		RawThreadTimestamp:  msg.RawThreadTimestamp,
		Time:                msg.Time,
		Type:                msg.Type,
	}

	d, ok := newmsg.Delegator.(delegate.Delegator)
	if !ok {
		return nil, fmt.Errorf("expected type delegate.Delegator for recursion: got %T", newmsg.Delegator)
	}

	res, err := d.Delegate(newmsg)
	if err != nil {
		return nil, errors.Wrapf(err, "recursing lookup for %s", i.Channel)
	}

	return res, nil
}
