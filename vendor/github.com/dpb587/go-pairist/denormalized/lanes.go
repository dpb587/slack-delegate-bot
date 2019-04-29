package denormalized

import (
	"github.com/dpb587/go-pairist/api"
)

type Lane struct {
	ID     string
	People []Entity
	Roles  []Entity
	Tracks []Entity
}

type Lanes []Lane

func (l Lanes) ByRole(name string) []Lane {
	var res []Lane

	for _, r := range l {
		for _, b := range r.Roles {
			if b.Name == name {
				res = append(res, r)

				break
			}
		}
	}

	return res
}

type Entity struct {
	Color     string
	Icon      string
	Name      string
	Picture   string
	UpdatedAt uint
}

func BuildLanes(historical *api.TeamHistorical) Lanes {
	var lanes Lanes

	for laneID := range historical.Lanes {
		lane := Lane{}

		for _, entity := range historical.Entities {
			if entity.Location != laneID {
				continue
			}

			denormalizedEntity := Entity{
				Color:     entity.Color,
				Icon:      entity.Icon,
				Name:      entity.Name,
				Picture:   entity.Picture,
				UpdatedAt: entity.UpdatedAt,
			}

			switch entity.Type {
			case "person":
				lane.People = append(lane.People, denormalizedEntity)
			case "role":
				lane.Roles = append(lane.Roles, denormalizedEntity)
			case "track":
				lane.Tracks = append(lane.Tracks, denormalizedEntity)
			}
		}

		lanes = append(lanes, lane)
	}

	return lanes
}
