package event

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/dpb587/slack-delegate-bot/pkg/slack"
	"github.com/pkg/errors"
	"github.com/slack-go/slack/slackevents"
)

type SyncProcessor struct {
	parser    *Parser
	responder *slack.Responder
}

var _ Processor = &SyncProcessor{}

func NewSyncProcessor(parser *Parser, responder *slack.Responder) Processor {
	return &SyncProcessor{
		parser:    parser,
		responder: responder,
	}
}

func (p *SyncProcessor) Process(since time.Time, event string, payload []byte) error {
	switch event {
	case "callback_event":
		return p.processCallbackEvent(since, payload)
	}

	return fmt.Errorf("unexpected event type: %v", event)
}

func (p *SyncProcessor) processCallbackEvent(since time.Time, payload []byte) error {
	event, err := slackevents.ParseEvent(json.RawMessage(payload), slackevents.OptionNoVerifyToken())
	if err != nil {
		return errors.Wrap(err, "parsing event")
	}

	switch inner := event.InnerEvent.Data.(type) {
	case *slackevents.AppMentionEvent:
		msg, reply, err := p.parser.ParseAppMention(event, *inner)
		if err != nil {
			return errors.Wrap(err, "parsing app mention")
		} else if !reply {
			return nil
		}

		return p.responder.ProcessMessage(msg)
	case *slackevents.MessageEvent:
		msg, reply, err := p.parser.ParseMessage(event, *inner)
		if err != nil {
			return errors.Wrap(err, "parsing message")
		} else if !reply {
			return nil
		}

		return p.responder.ProcessMessage(msg)
	}

	return nil
}
