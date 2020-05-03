package slack

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/slack-go/slack/slackevents"
)

type SyncProcessor struct {
	eventParser *EventParser
	responder   *Responder
}

var _ Processor = &SyncProcessor{}

func NewSyncProcessor(eventParser *EventParser, responder *Responder) Processor {
	return &SyncProcessor{
		eventParser: eventParser,
		responder:   responder,
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
		msg, reply, err := p.eventParser.ParseAppMention(event, *inner)
		if err != nil {
			return errors.Wrap(err, "parsing app mention")
		} else if !reply {
			return nil
		}

		return p.responder.ProcessMessage(msg)
	case *slackevents.MessageEvent:
		msg, reply, err := p.eventParser.ParseMessage(event, *inner)
		if err != nil {
			return errors.Wrap(err, "parsing message")
		} else if !reply {
			return nil
		}

		return p.responder.ProcessMessage(msg)
	}

	return nil
}
