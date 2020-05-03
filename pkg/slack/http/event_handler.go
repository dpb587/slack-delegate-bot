package http

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	ourslack "github.com/dpb587/slack-delegate-bot/pkg/slack"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

type EventHandler struct {
	messenger     *ourslack.Messenger
	signingSecret string
}

func NewEventHandler(messenger *ourslack.Messenger, signingSecret string) *EventHandler {
	return &EventHandler{
		messenger:     messenger,
		signingSecret: signingSecret,
	}
}

func (h EventHandler) Accept(c echo.Context) error {
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
		c := ourslack.EventContext{
			AppID:  event.APIAppID,
			TeamID: event.TeamID,
		}

		switch inner := event.InnerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			return h.messenger.HandleAppMention(c, *inner)
		case *slackevents.MessageEvent:
			return h.messenger.HandleMessage(c, *inner)
		}
	}

	// warn to unsubscribe

	return c.NoContent(http.StatusOK)
}
