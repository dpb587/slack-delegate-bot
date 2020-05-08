package http

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/dpb587/slack-delegate-bot/pkg/slack/slash"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/slack-go/slack"
)

type SlashHandler struct {
	processor     slash.Processor
	signingSecret string
}

func NewSlashHandler(processor slash.Processor, signingSecret string) *SlashHandler {
	return &SlashHandler{
		processor:     processor,
		signingSecret: signingSecret,
	}
}

func (h SlashHandler) Accept(c echo.Context) error {
	at := time.Now()

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

	// TODO reconsider+refactor
	req := c.Request()
	req.Body = ioutil.NopCloser(bytes.NewReader(body))

	cmd, err := slack.SlashCommandParse(req)
	if err != nil {
		return errors.Wrap(err, "parsing incoming command")
	}

	err = h.processor.Process(at, cmd.Command, body)
	if err != nil {
		return errors.Wrap(err, "processing")
	}

	return c.NoContent(http.StatusOK)
}
