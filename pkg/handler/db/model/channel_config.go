package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type ChannelConfig struct {
	ID             uuid.UUID `gorm:"primary"`
	TeamID         string    `gorm:"unique:revision"` // TODO unique_index not working on sqlite
	ChannelID      string    `gorm:"unique:revision"`
	RevisionNum    int       `gorm:"unique:revision"`
	RevisionLatest bool

	UpdatedAt     time.Time
	UpdatedByID   string
	UpdatedByName string

	Config string
}

func (m *ChannelConfig) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("ID", uuid.New().String())
}
