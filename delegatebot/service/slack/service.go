package slack

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/dpb587/slack-delegate-bot/delegatebot/handler"
	"github.com/dpb587/slack-delegate-bot/logic/delegate/delegates"
	"github.com/dpb587/slack-delegate-bot/delegatebot/message"
	"github.com/nlopes/slack"
	"github.com/sirupsen/logrus"
)

type Service struct {
	handler handler.Handler
	logger  logrus.FieldLogger

	api  *slack.Client
	self *slack.UserDetails
}

func New(api *slack.Client, handler_ handler.Handler, logger logrus.FieldLogger) *Service {
	return &Service{
		api:     api,
		handler: handler_,
		logger:  logger,
	}
}

func (s *Service) Run() error {
	rtm := s.api.NewRTM()
	go rtm.ManageConnection()

	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.ConnectedEvent:
				s.self = ev.Info.User

				s.logger.Infof("connected: %#+v\n", s.self)

			case *slack.MessageEvent:
				if ev.Msg.Type == "message" && ev.Msg.SubType != "message_deleted" {
					intmsg, useful := s.buildMessage(ev.Msg)
					if !useful {
						continue
					}

					s.logger.Debugf("received message: %#+v", intmsg)

					response, err := s.handler.Execute(&intmsg)
					if err != nil {
						s.logger.Errorf("failed to apply handler: %v", err)

						continue
					}

					var msg string

					if len(response.Delegates) > 0 {
						msg = delegates.Join(response.Delegates, " ")

						if intmsg.OriginType == message.ChannelOriginType {
							msg = fmt.Sprintf("^ %s", msg)
						}
					} else if response.EmptyMessage != "" {
						msg = response.EmptyMessage
					}

					if msg == "" {
						s.logger.Debugf("no response")

						continue
					}

					outgoing := rtm.NewOutgoingMessage(msg, ev.Msg.Channel)

					if intmsg.OriginType == message.ChannelOriginType {
						outgoing.ThreadTimestamp = ev.Msg.Timestamp
					}

					rtm.SendMessage(outgoing)
				}
			case *slack.RTMError:
				s.logger.Warnf("RTM: %s", ev.Error())

			case *slack.InvalidAuthEvent:
				err := errors.New("invalid credentials")

				s.logger.Errorf("auth: %s", err)

				return err
			}
		}
	}
}

func (s *Service) buildMessage(msg slack.Msg) (message.Message, bool) {
	intmsg := message.Message{
		Origin:          msg.Channel,
		OriginType:      message.ChannelOriginType,
		InterruptTarget: msg.Channel,
		Timestamp:       time.Now(), // TODO
		Text:            msg.Text,
	}

	if msg.Channel[0] == 'D' { // TODO better way to detect if this is our bot DM?
		re, err := regexp.Compile(`<#([^|]+)|([^>]+)>`)
		if err != nil {
			panic(err)
		}

		matches := re.FindStringSubmatch(msg.Text)
		if len(matches) > 0 {
			intmsg.InterruptTarget = matches[1]
		}

		intmsg.OriginType = message.DirectMessageOriginType

		return intmsg, true
	} else if strings.Contains(msg.Text, fmt.Sprintf("<@%s>", s.self.ID)) {
		re, err := regexp.Compile(fmt.Sprintf(`<#([^|]+)|([^>]+)>\s+<@%s>`, regexp.QuoteMeta(s.self.ID)))
		if err != nil {
			panic(err)
		}

		matches := re.FindStringSubmatch(msg.Text)
		if len(matches) > 0 {
			intmsg.InterruptTarget = matches[1]
		}

		return intmsg, true
	}

	return message.Message{}, false
}
