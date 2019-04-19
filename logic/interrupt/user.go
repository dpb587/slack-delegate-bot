package interrupt

import "fmt"

type User struct {
	ID string
}

var _ Interruptible = &User{}

func (i User) String() string {
	return fmt.Sprintf("<@%s>", i.ID)
}
