package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

var DefaultClient = NewClient(http.DefaultClient, "https://pairist-9de4d.firebaseio.com")

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

func (c *Client) Get(path string, data interface{}) error {
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
