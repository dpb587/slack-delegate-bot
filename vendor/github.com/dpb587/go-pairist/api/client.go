package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

var DefaultClient = NewClient(http.DefaultClient, DefaultDatabaseURL, nil)

type Client struct {
	client  *http.Client
	baseURL string
	auth    *Auth
}

func NewClient(client *http.Client, baseURL string, auth *Auth) *Client {
	return &Client{
		client:  client,
		baseURL: baseURL,
		auth:    auth,
	}
}

func (c *Client) get(path string, data interface{}) error {
	uri := fmt.Sprintf("%s/%s", c.baseURL, path)
	if c.auth != nil {
		token, err := c.auth.IDToken()
		if err != nil {
			return err
		}

		uri = fmt.Sprintf("%s?auth=%s", uri, token)
	}

	res, err := c.client.Get(uri)
	if err != nil {
		return err
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.Wrap(err, "reading response body")
	}

	err = json.Unmarshal(bytes, data)
	if err != nil {
		return errors.Wrap(err, "unmarshalling response")
	}

	return nil
}

func (c *Client) GetTeamCurrent(team string) (*TeamHistorical, error) {
	var res TeamHistorical

	err := c.get(fmt.Sprintf("teams/%s/current.json", team), &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) GetTeamLists(team string) (*TeamLists, error) {
	var res TeamLists

	err := c.get(fmt.Sprintf("teams/%s/lists.json", team), &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
