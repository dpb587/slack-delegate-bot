package slash

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/slack-go/slack"
)

type SyncProcessor struct {
	handler Handler
}

var _ Processor = &SyncProcessor{}

func NewSyncProcessor(handler Handler) Processor {
	return &SyncProcessor{
		handler: handler,
	}
}

func (p *SyncProcessor) Process(since time.Time, event string, payload []byte) error {
	switch event {
	case "/interrupt":
		return p.processInterruptCommand(since, payload)
	}

	return fmt.Errorf("unexpected slash command: %v", event)
}

func (p *SyncProcessor) processInterruptCommand(since time.Time, payload []byte) error {
	req, _ := http.NewRequest(http.MethodPost, "https:slack.local", ioutil.NopCloser(bytes.NewReader(payload)))
	req.Header.Set("content-type", "application/x-www-form-urlencoded")

	cmd, err := slack.SlashCommandParse(req)
	if err != nil {
		return errors.Wrap(err, "parsing payload")
	}

	done, err := p.handler.Handle(cmd)
	if err != nil {
		return errors.Wrap(err, "handling command")
	}

	if !done {
		// TODO respond with confusion? last handler, maybe
	}

	return nil
}
