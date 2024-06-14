package actor

import "github.com/google/uuid"

type Envelope struct {
	sender   uuid.UUID
	message  interface{}
	receiver uuid.UUID
}

func NewEnvelope(sender uuid.UUID, message interface{}, receiver uuid.UUID) *Envelope {
	return &Envelope{
		sender:   sender,
		message:  message,
		receiver: receiver,
	}
}

func (e *Envelope) Sender() *uuid.UUID {
	if e.sender == uuid.Nil {
		return nil
	}
	return &e.sender
}

func (e *Envelope) Receiver() *uuid.UUID {
	if e.receiver == uuid.Nil {
		return nil
	}
	return &e.receiver
}

func (e *Envelope) Unwrap() (*uuid.UUID, interface{}, *uuid.UUID) {
	var sender *uuid.UUID
	var receiver *uuid.UUID

	if e.sender != uuid.Nil {
		sender = &e.sender
	} else {
		sender = nil
	}

	if e.receiver == uuid.Nil {
		receiver = &e.receiver
	} else {
		receiver = nil
	}

	return sender, e.message, receiver
}
