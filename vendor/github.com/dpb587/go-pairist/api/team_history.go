package api

import (
	"fmt"
)

type TeamHistorical struct {
	Entities map[string]TeamHistoricalEntity `json:"entities,omitempty"`
	Lanes    map[string]TeamHistoricalLane   `json:"lanes,omitempty"`
}

type TeamHistoricalEntity struct {
	Color     string   `json:"color,omitempty"`
	Icon      string   `json:"icon,omitempty"`
	Location  string   `json:"location,omitempty"`
	Name      string   `json:"name,omitempty"`
	Picture   string   `json:"picture,omitempty"`
	Tags      []string `json:"tags,omitempty"`
	Type      string   `json:"type,omitempty"`
	UpdatedAt uint     `json:"updatedAt,omitempty"`
}

type TeamHistoricalLane struct {
	Locked    bool `json:"locked,omitempty"`
	SortOrder uint `json:"sortOrder,omitempty"`
}

func (c *Client) GetTeamCurrent(team string) (*TeamHistorical, error) {
	var res TeamHistorical

	err := c.Get(fmt.Sprintf("teams/%s/current.json", team), &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
