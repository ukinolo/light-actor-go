package envelope

import (
	"light-actor-go/pid"

	"github.com/google/uuid"
)

type Envelope struct {
	Message  interface{}
	receiver pid.PID
}

func NewEnvelope(message interface{}, receiver pid.PID) Envelope {
	return Envelope{
		Message:  message,
		receiver: receiver,
	}
}

func (e *Envelope) Receiver() *pid.PID {
	if e.receiver.ID == uuid.Nil {
		return nil
	}
	return &e.receiver
}

func (e *Envelope) Unwrap() (interface{}, *pid.PID) {

	var receiver *pid.PID

	if e.receiver.ID == uuid.Nil {
		receiver = &e.receiver
	} else {
		receiver = nil
	}

	return e.Message, receiver
}
