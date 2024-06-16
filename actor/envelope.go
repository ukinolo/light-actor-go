package actor

import "github.com/google/uuid"

type Envelope struct {
	sender   PID
	message  interface{}
	receiver PID
}

func NewEnvelope(sender PID, message interface{}, receiver PID) *Envelope {
	return &Envelope{
		sender:   sender,
		message:  message,
		receiver: receiver,
	}
}

func (e *Envelope) Sender() *PID {
	if e.sender.ID == uuid.Nil {
		return nil
	}
	return &e.sender
}

func (e *Envelope) Receiver() *PID {
	if e.receiver.ID == uuid.Nil {
		return nil
	}
	return &e.receiver
}

func (e *Envelope) Unwrap() (*PID, interface{}, *PID) {
	var sender *PID
	var receiver *PID

	if e.sender.ID != uuid.Nil {
		sender = &e.sender
	} else {
		sender = nil
	}

	if e.receiver.ID == uuid.Nil {
		receiver = &e.receiver
	} else {
		receiver = nil
	}

	return sender, e.message, receiver
}
