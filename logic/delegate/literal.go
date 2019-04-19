package delegate

type Literal struct {
	Text string
}

var _ Delegate = &Literal{}

func (i Literal) String() string {
	return i.Text
}
