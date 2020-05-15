package rtm

import (
	"encoding/json"
	"fmt"

	ourslack "github.com/dpb587/slack-delegate-bot/pkg/slack"
	"github.com/pkg/errors"
	"github.com/slack-go/slack"
	"go.uber.org/zap"
)

type Service struct {
	parser    *Parser
	rtm       *slack.RTM
	responder *ourslack.Responder
	logger    *zap.Logger
}

func NewService(api *slack.Client, responder *ourslack.Responder, logger *zap.Logger) *Service {
	return &Service{
		rtm:       api.NewRTM(),
		responder: responder,
		logger:    logger,
	}
}

func (s *Service) Run() error {
	go s.rtm.ManageConnection()

	for {
		select {
		case msg := <-s.rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.ConnectedEvent:
				s.parser = NewParser(ev.Info.Team.ID, ev.Info.User.ID)

				s.logger.Info(fmt.Sprintf("connected as %s (%s)", ev.Info.User.Name, ev.Info.User.ID))
			case *slack.MessageEvent:
				if s.parser == nil {
					// we assign parser only after we're connected
					s.logger.Warn("ignoring received message because no parser is ready")

					continue
				}

				msg, reply, err := s.parser.ParseMessage(ev.Msg)
				if err != nil {
					s.logger.Error("unable to parse message", zap.Error(err))

					continue
				} else if !reply {
					continue
				}

				msgBuf, err := json.Marshal(ev.Msg)
				if err != nil {
					s.logger.Error("failed to dump incoming message for debugging", zap.Error(err))
				}

				s.logger.Debug("received rtm message", zap.ByteString("payload", msgBuf))

				err = s.responder.ProcessMessage(msg)
				if err != nil {
					s.logger.Error("failed to process message", zap.Error(err))
				}
			case *slack.RTMError:
				s.logger.Error("rtm error occurred", zap.Error(ev))
			case *slack.InvalidAuthEvent:
				return errors.New("invalid credentials")
			}
		}
	}
}
