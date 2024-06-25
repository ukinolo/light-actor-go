package actor

type mailboxState int32

const (
	mailboxRunning mailboxState = iota
	mailboxStopping
)

type Mailbox struct {
	actorChan   chan Envelope
	mailboxChan chan Envelope
	queue       []Envelope
	state       mailboxState
}

func NewMailbox(actorChan chan Envelope) *Mailbox {
	m := &Mailbox{
		actorChan:   actorChan,
		mailboxChan: make(chan Envelope),
		queue:       make([]Envelope, 0),
		state:       mailboxStopping,
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
	m.state = mailboxRunning
	var newEnvelope Envelope
	var haveReady bool = false
	for {
		for haveReady {
			select {
			case m.actorChan <- newEnvelope:
				if len(m.queue) > 0 {
					newEnvelope = m.getEnvelope()
				} else {
					haveReady = false
				}
			case envelope := <-m.mailboxChan:
				switch msg := envelope.Message.(type) {
				case SystemMessage:
					if msg.Type == DeleteMailbox {
						m.delete()
						return
					}
				}
				m.buffer(envelope)
			}
		}
		newEnvelope = <-m.mailboxChan
		switch msg := newEnvelope.Message.(type) {
		case SystemMessage:
			if msg.Type == DeleteMailbox {
				m.delete()
				return
			}
		}
		haveReady = true
	}
}

func (m *Mailbox) GetChan() chan Envelope {
	return m.mailboxChan
}

func (m *Mailbox) delete() {
	close(m.actorChan)
	clear(m.queue)
}
