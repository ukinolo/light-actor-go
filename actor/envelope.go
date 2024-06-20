package actor

import (
	"github.com/google/uuid"
)

type Envelope struct {
	message  interface{}
	receiver PID
}

func NewEnvelope(message interface{}, receiver PID) *Envelope {
	return &Envelope{
		message:  message,
		receiver: receiver,
	}
}

func (e *Envelope) Receiver() *PID {
	if e.receiver.ID == uuid.Nil {
		return nil
	}
	return &e.receiver
}

func (e *Envelope) Unwrap() (interface{}, *PID) {
	return e.message, &e.receiver
}
