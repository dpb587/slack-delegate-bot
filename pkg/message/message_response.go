package message

type MessageResponse struct {
	Delegates    []Delegate
	EmptyMessage string
}

func (mr MessageResponse) IsUnset() bool {
	return len(mr.Delegates) == 0 && len(mr.EmptyMessage) == 0
}
