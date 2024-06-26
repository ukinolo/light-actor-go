package actor

type mailboxState int32

const (
	mailboxRunning mailboxState = iota
	mailboxStopping
	mailboxForceShutdown
	mailboxGraceShutdown
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
		state:       mailboxGraceShutdown,
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

func (m *Mailbox) start() {
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
				m.handleEnvelope(&envelope)
				if m.state == mailboxForceShutdown {
					m.actorChan <- NewEnvelope(forcefulShutdown{}, PID{})
					m.state = mailboxForceShutdown
					m.delete()
					return
				}
				if m.state == mailboxGraceShutdown {
					m.delete()
					return
				}
				m.buffer(envelope)
			}
		}
		newEnvelope = <-m.mailboxChan
		m.handleEnvelope(&newEnvelope)
		if m.state == mailboxForceShutdown {
			m.actorChan <- NewEnvelope(forcefulShutdown{}, PID{})
			m.delete()
			return
		}
		if m.state == mailboxGraceShutdown {
			m.delete()
			return
		}
		haveReady = true
	}
}

func (m *Mailbox) GetChan() chan Envelope {
	return m.mailboxChan
}

func (m *Mailbox) handleEnvelope(envelope *Envelope) {
	//fmt.Printf("Prispela mi je poruka %v koja je tipa %T\n", envelope.Message, envelope.Message)
	switch (*envelope).Message.(type) {
	case forcefulShutdown:
		m.state = mailboxForceShutdown
	case gracefulShutdown:
		m.state = mailboxStopping
	case startingActor:
		(*envelope).Message = ActorStarted{}
	case closeMailbox:
		m.state = mailboxGraceShutdown
	default:
		return
	}
}

func (m *Mailbox) delete() {
	close(m.actorChan)
	clear(m.queue)
}

// func (m *Mailbox) allowedMessage(envelope Envelope) bool {
// 	switch m.state {
// 	case mailboxRunning:
// 		return true
// 	case mailboxForceShutdown:
// 		return false
// 	case mailboxGraceShutdown:
// 		return true
// 	case mailboxStopping:
// 		switch envelope.Message.(type){
// 		case childTerminated:
// 			return true
// 		case closeMailbox:
// 			return true
// 		case
// 		default:
// 			return false
// 		}
// 	}
// }
