package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type TeamConfig struct {
	ID             uuid.UUID `gorm:"primary"`
	TeamID         string    `gorm:"unique:revision"` // TODO unique_index not working on sqlite
	RevisionNum    int       `gorm:"unique:revision"`
	RevisionLatest bool

	UpdatedAt     time.Time
	UpdatedByID   string
	UpdatedByName string

	HelpText      string
	DefaultConfig string
}

func (m *TeamConfig) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("ID", uuid.New().String())
}
