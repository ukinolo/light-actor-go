package actor

type Mailbox struct {
	actorChan chan Envelope
	queue     []Envelope
}

func NewMailbox(actorChan chan Envelope) *Mailbox {
	m := &Mailbox{
		actorChan: actorChan,
		queue:     make([]Envelope, 0),
	}
	return m
}

func (m *Mailbox) Send(msg Envelope) {
	select {
	case m.actorChan <- msg:
		// Successfully sent the message to the actorChan
	default:
		// The actorChan is full
	}
}

func (m *Mailbox) Receive() Envelope {

	return Envelope{}
}

func (m *Mailbox) Buffer(msg Envelope) {
	m.queue = append(m.queue, msg)
}

func (m *Mailbox) Start() {

}

func (m *Mailbox) Stop() {

}
