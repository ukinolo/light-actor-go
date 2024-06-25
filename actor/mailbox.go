package actor

type mailboxState int32

const (
	mailboxRunning mailboxState = iota
	mailboxStopping
	mailboxStoped
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
		state:       mailboxStoped,
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
	fn := func() Envelope {
		if !haveReady {
			haveReady = true
			if len(m.queue) > 0 {
				newEnvelope = m.getEnvelope()
			} else {
				newEnvelope = <-m.mailboxChan
			}
		}
		m.handleEnvelope(&newEnvelope)
		return newEnvelope
	}
	for {
		select {
		case m.actorChan <- fn():
			haveReady = false
			if m.state == mailboxStoped {
				return
			}
		case envelope := <-m.mailboxChan:
			m.handleEnvelope(&envelope)
			if m.state == mailboxStoped {
				m.actorChan <- NewEnvelope(forcefulShutdown{}, PID{})
				return
			}
			if m.state == mailboxStopping {
				//Handle gracefull shutdown??
			}
			m.buffer(envelope)
		}
	}
}

func (m *Mailbox) GetChan() chan Envelope {
	return m.mailboxChan
}

func (m *Mailbox) handleEnvelope(envelope *Envelope) {
	switch (*envelope).Message.(type) {
	case forcefulShutdown:
		m.state = mailboxStoped
	case gracefulShutdown:
		m.state = mailboxStopping
	case startingActor:
		(*envelope).Message = ActorStarted{}
	case stopingActor:
		(*envelope).Message = ActorStoped{}
	default:
		return
	}
}
