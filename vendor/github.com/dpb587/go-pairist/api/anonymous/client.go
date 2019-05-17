package anonymous

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dpb587/go-pairist/api"
	"github.com/pkg/errors"
)

var DefaultClient = NewClient(http.DefaultClient, api.DefaultDatabaseURL)

type Client struct {
	client  *http.Client
	baseURL string
}

func NewClient(client *http.Client, baseURL string) *Client {
	return &Client{
		client:  client,
		baseURL: baseURL,
	}
}

func (c *Client) get(path string, data interface{}) error {
	res, err := c.client.Get(fmt.Sprintf("%s/%s", c.baseURL, path))
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

func (c *Client) GetTeamCurrent(team string) (*api.TeamHistorical, error) {
	var res api.TeamHistorical

	err := c.get(fmt.Sprintf("teams/%s/current.json", team), &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) GetTeamLists(team string) (*api.TeamLists, error) {
	var res api.TeamLists

	err := c.get(fmt.Sprintf("teams/%s/lists.json", team), &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
