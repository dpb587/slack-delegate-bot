package slack

import (
	"sync"

	"github.com/pkg/errors"
	"github.com/slack-go/slack"
)

type Cache struct {
	api *slack.Client

	mappedUserApp     map[string]string
	mappedUserAppSync sync.RWMutex
}

func NewCache(api *slack.Client) *Cache {
	return &Cache{
		api:           api,
		mappedUserApp: map[string]string{},
	}
}

func (c *Cache) IsAppBot(appID, userID string) (bool, error) {
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

	c.mappedUserAppSync.Lock()
	c.mappedUserApp[userID] = is
	c.mappedUserAppSync.Unlock()

	return appID == is, nil
}
