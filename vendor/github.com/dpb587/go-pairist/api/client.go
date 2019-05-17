package api

type Client interface {
	GetTeamCurrent(team string) (*TeamHistorical, error)
	GetTeamLists(team string) (*TeamLists, error)
}
