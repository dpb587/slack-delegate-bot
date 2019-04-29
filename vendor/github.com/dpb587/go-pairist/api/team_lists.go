package api

import (
	"fmt"
)

type TeamLists map[string]TeamList

type TeamList struct {
	Items TeamListItems `json:"items,omitempty"`
	Title string        `json:"title,omitempty"`
}

type TeamListItems map[string]TeamListItem

type TeamListItem struct {
	Checked bool   `json:"checked,omitempty"`
	Title   string `json:"title,omitempty"`
}

func (c *Client) GetTeamLists(team string) (*TeamLists, error) {
	var res TeamLists

	err := c.Get(fmt.Sprintf("teams/%s/lists.json", team), &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
