package http

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	ourslack "github.com/dpb587/slack-delegate-bot/pkg/slack"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

type EventHandler struct {
	processor     ourslack.Processor
	signingSecret string
}

func NewEventHandler(processor ourslack.Processor, signingSecret string) *EventHandler {
	return &EventHandler{
		processor:     processor,
		signingSecret: signingSecret,
	}
}

func (h EventHandler) Accept(c echo.Context) error {
	at := time.Now()

	if c.Request().Header.Get("content-type") != "application/json" {
		return c.String(http.StatusUnsupportedMediaType, http.StatusText(http.StatusUnsupportedMediaType))
	}

	verifier, err := slack.NewSecretsVerifier(c.Request().Header, h.signingSecret)
	if err != nil {
		return errors.Wrap(err, "building secrets verifier")
	}

	body, err := ioutil.ReadAll(io.TeeReader(c.Request().Body, &verifier))
	if err != nil {
		return errors.Wrap(err, "reading body")
	}

	fmt.Printf("%s\n", body) // TODO log.DEBUG

	if err = verifier.Ensure(); err != nil {
		// TODO log err
		return c.String(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
	}

	event, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionNoVerifyToken())
	if err != nil {
		return errors.Wrap(err, "parsing event")
	}

	switch event.Type {
	case slackevents.URLVerification:
		var r *slackevents.ChallengeResponse

		err := json.Unmarshal([]byte(body), &r)
		if err != nil {
			return errors.Wrap(err, "unmarshalling verification payload")
		}

		return c.String(http.StatusOK, r.Challenge)
	case slackevents.CallbackEvent:
		err := h.processor.Process(at, "callback_event", body)
		if err != nil {
			return errors.Wrap(err, "processing event")
		}

		return c.NoContent(http.StatusAccepted)
	}

	// TODO warn to add handling of or unsubscribe from event

	return c.NoContent(http.StatusOK)
}
