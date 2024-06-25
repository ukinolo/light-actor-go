package actor

type SystemMessageType int

const (
	SystemMessageStart SystemMessageType = iota
	SystemMessageStop
	SystemMessageGracefulStop
	SystemMessageChildTerminated
	DeleteMailbox
)

type SystemMessage struct {
	Type   SystemMessageType
	Extras interface{}
}
