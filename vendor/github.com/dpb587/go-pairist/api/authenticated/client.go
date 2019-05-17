package authenticated

import (
	"context"
	"fmt"

	"firebase.google.com/go"
	"firebase.google.com/go/db"
	"github.com/dpb587/go-pairist/api"
	"google.golang.org/api/option"
)

type Client struct {
	ctx    context.Context
	client *db.Client
}

func CreateConfigClient(ctx context.Context, config *firebase.Config, opts ...option.ClientOption) (*Client, error) {
	app, err := firebase.NewApp(
		ctx,
		config,
		opts...,
	)
	if err != nil {
		return nil, err
	}

	client, err := app.Database(ctx)
	if err != nil {
		return nil, err
	}

	return &Client{
		ctx:    ctx,
		client: client,
	}, nil
}

func CreateClient(apiKey, teamName, teamPassword string) (*Client, error) {
	ctx := context.Background()

	tokenSource, err := CreateTokenSource(ctx, apiKey, teamName, teamPassword)
	if err != nil {
		return nil, err
	}

	config := &firebase.Config{
		DatabaseURL: api.DefaultDatabaseURL,
	}

	return CreateConfigClient(ctx, config, option.WithTokenSource(tokenSource))
}

func (c *Client) GetTeamCurrent(team string) (*api.TeamHistorical, error) {
	var res api.TeamHistorical

	ref := c.client.NewRef(fmt.Sprintf("teams/%s/current", team))

	err := ref.Get(c.ctx, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) GetTeamLists(team string) (*api.TeamLists, error) {
	var res api.TeamLists

	ref := c.client.NewRef(fmt.Sprintf("teams/%s/lists", team))

	err := ref.Get(c.ctx, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
