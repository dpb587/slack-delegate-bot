package delegate

import "fmt"

type UserGroup struct {
	ID    string
	Alias string
}

var _ Delegate = &UserGroup{}

func (i UserGroup) String() string {
	return fmt.Sprintf("<!subteam^%s|@%s>", i.ID, i.Alias)
}
