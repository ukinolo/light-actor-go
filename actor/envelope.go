package actor

import (
	"github.com/google/uuid"
)

type Envelope struct {
	Message  interface{}
	receiver PID
}

func NewEnvelope(message interface{}, receiver PID) Envelope {
	return Envelope{
		Message:  message,
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

	var receiver *PID

	if e.receiver.ID == uuid.Nil {
		receiver = &e.receiver
	} else {
		receiver = nil
	}

	return e.Message, receiver
}
