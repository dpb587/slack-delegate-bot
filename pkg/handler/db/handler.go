package db

import (
	"github.com/dpb587/slack-delegate-bot/pkg/handler"
	"github.com/dpb587/slack-delegate-bot/pkg/handler/db/model"
	"github.com/dpb587/slack-delegate-bot/pkg/handler/yaml"
	"github.com/dpb587/slack-delegate-bot/pkg/message"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type Handler struct {
	db     *gorm.DB
	parser *yaml.Parser
}

func NewHandler(db *gorm.DB, parser *yaml.Parser) handler.Handler {
	return &Handler{
		db:     db,
		parser: parser,
	}
}

func (h *Handler) Execute(msg message.Message) (message.MessageResponse, error) {
	var config model.ChannelConfig

	err := h.db.Model(config).
		Where("team_id = ? AND channel_id = ?", msg.TargetChannelTeamID, msg.TargetChannelID).
		Where("revision_latest = ?", true).
		First(&config).
		Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return message.MessageResponse{}, errors.Wrap(err, "loading channel config")
	}

	if config.Config == "" {
		return h.executeTeam(msg)
	}

	return h.executeWithConfig(msg, config.Config)
}

func (h *Handler) executeTeam(msg message.Message) (message.MessageResponse, error) {
	var config model.TeamConfig

	err := h.db.Model(config).
		Where("team_id = ?", msg.TargetChannelTeamID).
		Where("revision_latest = ?", true).
		First(&config).
		Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return message.MessageResponse{}, errors.Wrap(err, "loading channel config")
	}

	if config.DefaultConfig == "" {
		return message.MessageResponse{}, nil
	}

	return h.executeWithConfig(msg, config.DefaultConfig)
}

func (h *Handler) executeWithConfig(msg message.Message, config string) (message.MessageResponse, error) {
	configHandler, err := h.parser.Parse([]byte(config))
	if err != nil {
		return message.MessageResponse{}, errors.Wrap(err, "parsing config")
	}

	return configHandler.Execute(msg)
}
