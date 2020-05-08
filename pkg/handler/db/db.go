package db

import (
	"fmt"

	"github.com/dpb587/slack-delegate-bot/pkg/handler/db/model"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

func OpenDB(adapter, profile string) (*gorm.DB, error) {
	if adapter == "mysql2" {
		profile = fmt.Sprintf("%s?charset=utf8&&parseTime=true", profile) // TODO ampersand
	}

	db, err := gorm.Open(adapter, profile)
	if err != nil {
		panic(errors.Wrap(err, "opening database"))
	}

	if err := db.AutoMigrate(&model.TeamConfig{}).Error; err != nil {
		return nil, errors.Wrap(err, "auto-migrating TeamConfig")
	}

	if err := db.AutoMigrate(&model.ChannelConfig{}).Error; err != nil {
		return nil, errors.Wrap(err, "auto-migrating ChannelConfig")
	}

	return db, nil
}
