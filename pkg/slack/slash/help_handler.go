package slash

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"github.com/slack-go/slack"
)

type HelpHandler struct {
	handlers *Handlers
	baseURL  string
}

var _ Handler = HelpHandler{}

func NewHelpHandler(handlers *Handlers, baseURL string) Handler {
	return &HelpHandler{
		handlers: handlers,
		baseURL:  baseURL,
	}
}

func (h HelpHandler) UsageHint() string {
	return "help"
}

func (h HelpHandler) ShortDescription() string {
	return "show available inline commands"
}

func (h HelpHandler) Handle(cmd slack.SlashCommand) (bool, error) {
	if cmd.Text != "" && cmd.Text != "help" {
		return false, nil
	}

	helps := []string{}

	for _, handler := range *h.handlers {
		usage := handler.UsageHint()
		if usage == "" {
			continue
		}

		shortDescription := handler.ShortDescription()
		if shortDescription != "" {
			shortDescription = fmt.Sprintf(" – %s", shortDescription)
		}

		helps = append(helps, fmt.Sprintf("_%s_%s", usage, shortDescription))
	}

	managementURL := fmt.Sprintf("%s/team/%s/channel/%s/", h.baseURL, cmd.TeamID, cmd.ChannelID)

	responseMessage := fmt.Sprintf(
		"Hello! I understand the following commands, and you can find more details at <%s|%s>.\n\n%s",
		managementURL,
		strings.SplitN(managementURL, "/", 4)[2],
		strings.Join(helps, "\n"),
	)

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
