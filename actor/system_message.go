package actor

type SystemMessageType int

const (
	SystemMessageStart SystemMessageType = iota
	SystemMessageStop
)

type SystemMessage struct {
	Type SystemMessageType
}
