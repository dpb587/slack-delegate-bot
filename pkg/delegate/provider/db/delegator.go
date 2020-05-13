package db

import (
	"github.com/dpb587/slack-delegate-bot/pkg/configutil"
	"github.com/dpb587/slack-delegate-bot/pkg/delegate"
	"github.com/dpb587/slack-delegate-bot/pkg/delegate/provider/db/model"
	ouryaml "github.com/dpb587/slack-delegate-bot/pkg/delegate/provider/yaml"
	"github.com/dpb587/slack-delegate-bot/pkg/message"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Delegator struct {
	db     *gorm.DB
	parser *ouryaml.Parser
}

var _ delegate.Delegator = &Delegator{}

func NewDelegator(db *gorm.DB, parser *ouryaml.Parser) delegate.Delegator {
	return &Delegator{
		db:     db,
		parser: parser,
	}
}

func (h *Delegator) Delegate(msg message.Message) ([]message.Delegate, error) {
	var config model.ChannelConfig

	err := h.db.Model(config).
		Where("team_id = ? AND channel_id = ?", msg.TargetChannelTeamID, msg.TargetChannelID).
		Where("revision_latest = ?", true).
		First(&config).
		Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, errors.Wrap(err, "loading channel config")
	}

	if config.Config == "" {
		return h.delegateTeam(msg)
	}

	return h.delegateWithConfig(msg, config.Config, config.ConfigSecrets)
}

func (h *Delegator) delegateTeam(msg message.Message) ([]message.Delegate, error) {
	var config model.TeamConfig

	err := h.db.Model(config).
		Where("team_id = ?", msg.TargetChannelTeamID).
		Where("revision_latest = ?", true).
		First(&config).
		Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, errors.Wrap(err, "loading channel config")
	}

	if config.DefaultConfig == "" {
		return nil, nil
	}

	return h.delegateWithConfig(msg, config.DefaultConfig, config.DefaultConfigSecrets)
}

func (h *Delegator) delegateWithConfig(msg message.Message, rawConfig, rawConfigSecrets string) ([]message.Delegate, error) {
	var secrets map[string]interface{}

	if rawConfigSecrets != "" {
		err := yaml.Unmarshal([]byte(rawConfigSecrets), &secrets)
		if err != nil {
			return nil, errors.Wrap(err, "parsing secrets")
		}
	}

	desanitizedConfig, _, err := configutil.DesanitizeConfig(rawConfig, secrets)
	if err != nil {
		return nil, errors.Wrap(err, "reinjecting secrets")
	}

	d, err := h.parser.Parse([]byte(desanitizedConfig))
	if err != nil {
		return nil, errors.Wrap(err, "parsing config")
	}

	return d.Delegate(msg)
}
