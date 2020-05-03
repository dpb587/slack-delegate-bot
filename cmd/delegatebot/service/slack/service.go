package slack

import (
	"time"

	"github.com/dpb587/slack-delegate-bot/cmd/delegatebot/handler"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

type Service struct {
	handler handler.Handler
	logger  logrus.FieldLogger

	messageParser  *MessageParser
	messageHandler *MessageHandler

	api *slack.Client
	rtm *slack.RTM
}

func New(api *slack.Client, interruptHandler handler.Handler, logger logrus.FieldLogger) *Service {
	rtm := api.NewRTM()

	return &Service{
		api:            api,
		rtm:            rtm,
		messageHandler: NewMessageHandler(rtm, interruptHandler),
		logger:         logger,
	}
}

func (s *Service) Run() error {
	go s.rtm.ManageConnection()

	for {
		select {
		case msg := <-s.rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.ConnectedEvent:
				s.messageParser = NewMessageParser(ev.Info.User)

				s.logger.Infof("connected as %s (%s)", ev.Info.User.Name, ev.Info.User.ID)
			case *slack.MessageEvent:
				if s.messageParser == nil {
					// we assign messageParser only after we're connected
					s.logger.Errorf("received message, but no parser is available")

					continue
				}

				incoming, err := s.messageParser.ParseMessage(ev.Msg)
				if err != nil {
					s.logger.Error(errors.Wrap(err, "parsing request"))

					continue
				} else if incoming == nil {
					continue
				}

				s.logger.WithFields(logrus.Fields{
					"recv.origin":           incoming.Origin,
					"recv.origin_type":      incoming.OriginType,
					"recv.interrupt_target": incoming.InterruptTarget,
					"recv.text":             incoming.Text,
					"recv.timestamp":        incoming.Timestamp.Format(time.RFC3339),
				}).Debug("received message")

				outgoing, err := s.messageHandler.GetResponse(*incoming, ev)
				if err != nil {
					s.logger.Error(errors.Wrap(err, "getting response"))

					continue
				} else if outgoing == nil {
					continue
				}

				s.logger.WithFields(logrus.Fields{
					"send.channel":          outgoing.Channel,
					"send.text":             outgoing.Text,
					"send.thread_timestamp": outgoing.ThreadTimestamp,
				}).Debug("sending message")

				s.rtm.SendMessage(outgoing)
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
