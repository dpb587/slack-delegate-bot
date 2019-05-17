package api

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
