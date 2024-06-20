package actor

import (
	"light-actor-go/pid"

	"github.com/google/uuid"
)

type Envelope struct {
	message  interface{}
	receiver pid.PID
}


func NewEnvelope(message interface{}, receiver PID) Envelope {
	return Envelope{
		message:  message,
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
	return e.message, &e.receiver
}
