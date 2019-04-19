package interrupt

type Literal struct {
	Text string
}

var _ Interruptible = &Literal{}

func (i Literal) String() string {
	return i.Text
}
