package slash

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dpb587/slack-delegate-bot/pkg/handler"
	"github.com/dpb587/slack-delegate-bot/pkg/message"
	"github.com/pkg/errors"
	"github.com/slack-go/slack"
)

type ShowHandler struct {
	handler    handler.Handler
	serviceAPI interface{}
}

var _ Handler = HelpHandler{}

func NewShowHandler(h handler.Handler, serviceAPI interface{}) Handler {
	return &ShowHandler{
		handler:    h,
		serviceAPI: serviceAPI,
	}
}

func (h ShowHandler) UsageHint() string {
	return "show"
}

func (h ShowHandler) ShortDescription() string {
	return "show the current delegates for this channel"
}

func (h ShowHandler) Handle(cmd slack.SlashCommand) (bool, error) {
	if cmd.Text != "show" {
		return false, nil
	}

	now := time.Now()
	msg := message.Message{
		ServiceAPI:      h.serviceAPI,
		TeamID:          cmd.TeamID,
		OriginUserID:    cmd.UserID,
		OriginType:      message.DirectMessageOriginType,
		Origin:          cmd.ChannelID,
		OriginTimestamp: fmt.Sprintf("%d.0", now.UTC()),
		InterruptTarget: cmd.ChannelID,
		Timestamp:       now,
		Text:            "slash-command",
	}

	response, err := h.handler.Execute(&msg)
	if err != nil {
		return false, errors.Wrap(err, "evaluating a response")
	}

	var responseMessage string

	if len(response.Delegates) > 0 {
		responseMessage = "Here are the current interrupt details for this channel:"

		for _, d := range response.Delegates {
			responseMessage = fmt.Sprintf("%s\n- %s", responseMessage, d)
		}
	} else if response.EmptyMessage != "" {
		responseMessage = fmt.Sprintf(
			"Looks like there is not an interrupt available for this channel:\n> %s",
			strings.ReplaceAll(response.EmptyMessage, "\n", "\n> "),
		)
	} else {
		responseMessage = "Looks like there is no interrupt available for this channel."
	}

	bodyBuf, err := json.Marshal(map[string]string{
		"text":          responseMessage,
		"response_type": "ephemeral",
	})
	if err != nil {
		return false, errors.Wrap(err, "marshalling response")
	}

	res, err := http.DefaultClient.Post(
		cmd.ResponseURL,
		"application/json",
		bytes.NewReader(bodyBuf),
	)
	if err != nil {
		return false, errors.Wrap(err, "posting response")
	} else if res.StatusCode != 200 {
		// TODO valid
		return false, errors.Wrap(err, "status code")
	}

	return true, nil
}
