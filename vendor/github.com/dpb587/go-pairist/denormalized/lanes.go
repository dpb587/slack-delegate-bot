package denormalized

type Entity struct {
	Color     string
	Icon      string
	Name      string
	Picture   string
	UpdatedAt uint
}

type Lane struct {
	ID     string
	People []Entity
	Roles  []Entity
	Tracks []Entity
}

type Lanes []Lane

func (l Lanes) ByRole(name string) Lanes {
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

func (l Lanes) ByTrack(name string) Lanes {
	var res []Lane

	for _, r := range l {
		for _, b := range r.Tracks {
			if b.Name == name {
				res = append(res, r)

				break
			}
		}
	}

	return res
}
