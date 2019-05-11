package slack

import (
	"github.com/dpb587/slack-delegate-bot/delegatebot/handler"
	"github.com/nlopes/slack"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
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

				s.logger.Infof("connected: %#+v\n", ev.Info.User)
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

				s.logger.Debugf("received message: %#+v", incoming)

				outgoing, err := s.messageHandler.GetResponse(*incoming, ev)
				if err != nil {
					s.logger.Error(errors.Wrap(err, "getting response"))

					continue
				} else if outgoing == nil {
					continue
				}

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
