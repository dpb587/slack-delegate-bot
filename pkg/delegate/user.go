package delegate

import "fmt"

type User struct {
	ID string
}

var _ Delegate = &User{}

func (i User) String() string {
	return fmt.Sprintf("<@%s>", i.ID)
}
