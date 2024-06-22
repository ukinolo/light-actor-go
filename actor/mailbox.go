package actor

import "light-actor-go/envelope"

type Mailbox struct {
	actorChan   chan envelope.Envelope
	mailboxChan chan envelope.Envelope
	queue       []envelope.Envelope
}

func NewMailbox(actorChan chan envelope.Envelope) *Mailbox {
	m := &Mailbox{
		actorChan:   actorChan,
		mailboxChan: make(chan envelope.Envelope),
		queue:       make([]envelope.Envelope, 0),
	}
	return m
}

func (m *Mailbox) buffer(msg envelope.Envelope) {
	m.queue = append(m.queue, msg)
}

func (m *Mailbox) getEnvelope() envelope.Envelope {
	defer func() {
		m.queue = m.queue[1:]
	}()
	return m.queue[0]
}

func (m *Mailbox) Start() {
	var newEnvelope envelope.Envelope
	var haveReady bool = false
	fn := func() envelope.Envelope {
		if !haveReady {
			haveReady = true
			if len(m.queue) > 0 {
				newEnvelope = m.getEnvelope()
			} else {
				newEnvelope = <-m.mailboxChan
			}
		}
		return newEnvelope
	}
	for {
		select {
		case m.actorChan <- fn():
			haveReady = false
		case envelope := <-m.mailboxChan:
			m.buffer(envelope)
		}
	}
}

func (m *Mailbox) GetChan() chan envelope.Envelope {
	return m.mailboxChan
}
