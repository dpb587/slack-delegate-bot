package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type ChannelConfig struct {
	ID             uuid.UUID `gorm:"primary_key"`
	TeamID         string    `gorm:"unique_index:channel_config_revision"`
	ChannelID      string    `gorm:"unique_index:channel_config_revision"`
	RevisionNum    int       `gorm:"unique_index:channel_config_revision"`
	RevisionLatest bool

	UpdatedAt     time.Time
	UpdatedByID   string
	UpdatedByName string

	Config        string
	ConfigSecrets string
}

func (m *ChannelConfig) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("ID", uuid.New().String())
}
