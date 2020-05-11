package slash

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dpb587/slack-delegate-bot/pkg/delegate"
	"github.com/dpb587/slack-delegate-bot/pkg/message"
	"github.com/pkg/errors"
	"github.com/slack-go/slack"
)

type ShowHandler struct {
	delegator  delegate.Delegator
	serviceAPI interface{}
}

var _ Handler = HelpHandler{}

func NewShowHandler(h delegate.Delegator, serviceAPI interface{}) Handler {
	return &ShowHandler{
		delegator:  h,
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
		ServiceAPI:          h.serviceAPI,
		Delegator:           h.delegator,
		ChannelTeamID:       cmd.TeamID,
		ChannelID:           cmd.ChannelID,
		UserTeamID:          cmd.TeamID,
		UserID:              cmd.UserID,
		RawTimestamp:        fmt.Sprintf("%d.0", now.UTC().Unix()),
		TargetChannelTeamID: cmd.TeamID,
		TargetChannelID:     cmd.ChannelID,
		RawText:             "slash-command",
		Type:                message.DirectMessageMessageType,
		Time:                now,
	}

	dd, err := h.delegator.Delegate(msg)
	if err != nil {
		return false, errors.Wrap(err, "finding delegates")
	}

	var responseMessage string

	if len(dd) == 0 {
		responseMessage = "Looks like there is no interrupt available for this channel."
	} else {
		responseMessage = "Here are the current interrupt details for this channel:"

		for _, d := range dd {
			responseMessage = fmt.Sprintf("%s\n- %s", responseMessage, d)
		}
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
