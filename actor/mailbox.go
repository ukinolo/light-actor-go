package actor

type Mailbox struct {
	actorChan   chan Envelope
	mailboxChan chan Envelope
	queue       []Envelope
}

func NewMailbox(actorChan chan Envelope) *Mailbox {
	m := &Mailbox{
		actorChan:   actorChan,
		mailboxChan: make(chan Envelope),
		queue:       make([]Envelope, 0),
	}
	return m
}

func (m *Mailbox) buffer(msg Envelope) {
	m.queue = append(m.queue, msg)
}

func (m *Mailbox) getEnvelope() Envelope {
	defer func() {
		m.queue = m.queue[1:]
	}()
	return m.queue[0]
}

func (m *Mailbox) Start() {
	var newEnvelope Envelope
	var haveReady bool = false
	fn := func() Envelope {
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

func (m *Mailbox) GetChan() chan Envelope {
	return m.mailboxChan
}
