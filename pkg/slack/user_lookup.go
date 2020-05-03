package slack

import (
	"sync"

	"github.com/pkg/errors"
	"github.com/slack-go/slack"
)

//go:generate counterfeiter . UserLookupSlackAPI
type UserLookupSlackAPI interface {
	GetUserInfo(string) (*slack.User, error)
}

// TODO interface instead of SlackAPI
type UserLookup struct {
	api UserLookupSlackAPI

	mappedUserApp     map[string]string
	mappedUserAppSync sync.RWMutex
}

func NewUserLookup(api UserLookupSlackAPI) *UserLookup {
	return &UserLookup{
		api:           api,
		mappedUserApp: map[string]string{},
	}
}

func (c *UserLookup) IsAppBot(appID, userID string) (bool, error) {
	c.mappedUserAppSync.RLock()
	appUserID, known := c.mappedUserApp[userID]
	c.mappedUserAppSync.RUnlock()

	if known {
		return appID == appUserID, nil
	}

	userInfo, err := c.api.GetUserInfo(userID)
	if err != nil {
		return false, errors.Wrap(err, "getting user info")
	}

	// weird; why is IsBot = false on test environment?

	var is string

	if userInfo.Profile.ApiAppID == appID {
		is = appID
	}

	// TODO occasional purge of non-apps
	c.mappedUserAppSync.Lock()
	c.mappedUserApp[userID] = is
	c.mappedUserAppSync.Unlock()

	return appID == is, nil
}
