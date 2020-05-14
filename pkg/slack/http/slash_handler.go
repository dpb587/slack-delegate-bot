package http

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/dpb587/slack-delegate-bot/pkg/slack/slash"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/slack-go/slack"
	"go.uber.org/zap"
)

type SlashHandler struct {
	processor     slash.Processor
	signingSecret string
	logger        *zap.Logger
}

func NewSlashHandler(processor slash.Processor, signingSecret string, logger *zap.Logger) *SlashHandler {
	return &SlashHandler{
		processor:     processor,
		signingSecret: signingSecret,
		logger:        logger,
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

	if err = verifier.Ensure(); err != nil {
		h.logger.Debug("received unverified slack slash command", zap.ByteString("payload", body))
		h.logger.Warn("unable to verify incoming slack slash command", zap.Error(err))

		return c.String(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
	}

	h.logger.Debug("received slack slash event", zap.ByteString("payload", body))

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
