package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type TeamConfig struct {
	ID             uuid.UUID `gorm:"primary_key"`
	TeamID         string    `gorm:"unique_index:team_config_revision"`
	RevisionNum    int       `gorm:"unique_index:team_config_revision"`
	RevisionLatest bool

	UpdatedAt     time.Time
	UpdatedByID   string
	UpdatedByName string

	HelpText             string
	DefaultConfig        string
	DefaultConfigSecrets string
}

func (m *TeamConfig) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("ID", uuid.New().String())
}
