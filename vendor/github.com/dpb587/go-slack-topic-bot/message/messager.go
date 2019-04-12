package message

type Messager interface {
	Message() (string, error)
}
