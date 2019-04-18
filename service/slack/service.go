package slack

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/dpb587/slack-delegate-bot/handler"
	"github.com/dpb587/slack-delegate-bot/message"
	"github.com/nlopes/slack"
	"github.com/sirupsen/logrus"
)

type Service struct {
	token   string
	handler handler.Handler
	logger  logrus.FieldLogger

	api  *slack.Client
	self *slack.UserDetails
}

func New(token string, handler_ handler.Handler, logger logrus.FieldLogger) *Service {
	return &Service{
		token:   token,
		handler: handler_,
		logger:  logger,
	}
}

func (s *Service) API() *slack.Client {
	if s.api == nil {
		s.api = slack.New(s.token) // , slack.OptionDebug(true)) // TODO , slack.OptionLog(s.logger))
	}

	return s.api
}

func (s *Service) Run() error {
	rtm := s.API().NewRTM()
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

					err := s.handler.Apply(&intmsg)
					if err != nil {
						s.logger.Errorf("failed to apply handler: %v", err)

						continue
					}

					response := intmsg.GetResponse()
					if response == nil {
						s.logger.Debugf("no response found")

						continue
					}

					outgoing := rtm.NewOutgoingMessage(response.Text, ev.Msg.Channel)

					if intmsg.OriginType == message.ChannelOriginType {
						outgoing.Text = strings.TrimPrefix(outgoing.Text, "^ ") // TODO hacky
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
		intmsg.OriginType = message.DirectMessageOriginType
		intmsg.InterruptTarget = intmsg.Text

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
