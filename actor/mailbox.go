package actor

import (
	"sync"
)

type Mailbox struct {
	actorChan   chan Envelope
	mailboxChan chan Envelope
	queue       []Envelope
	queueMutex  sync.Mutex
	running     bool
	stopChan    chan struct{}
}

// Creates a new Mailbox instance.
func NewMailbox(actorChan chan Envelope, mailboxChan chan Envelope) *Mailbox {
	return &Mailbox{
		actorChan:   actorChan,
		mailboxChan: mailboxChan,
		queue:       make([]Envelope, 0),
		stopChan:    make(chan struct{}),
	}
}

// Attempts to send a message to the actorChan, or buffers it if the channel is full.
func (m *Mailbox) Send(msg Envelope) {
	select {
	case m.actorChan <- msg:
		// Successfully sent the message to the actorChan
	default:
		// The actorChan is full, buffer the message
		m.Buffer(msg)
	}
}

// Adds a message to the mailbox queue.
func (m *Mailbox) Buffer(msg Envelope) {
	m.queueMutex.Lock()
	defer m.queueMutex.Unlock()
	m.queue = append(m.queue, msg)
}

// Begins the mailbox's message dispatching loop.
func (m *Mailbox) Start() {
	m.running = true
	go func() {
		for {
			select {
			case <-m.stopChan:
				m.running = false
				return
			case msg := <-m.mailboxChan:
				m.Receive(msg)
			default:
				m.queueMutex.Lock()
				if len(m.queue) > 0 {
					m.dispatchMessages()
				}
				m.queueMutex.Unlock()
			}
		}
	}()
}

// Stop halts the mailbox's message dispatching loop.
func (m *Mailbox) Stop() {
	if m.running {
		close(m.stopChan)
	}
}

// Receive enqueues a message from the mailbox channel into the queue.
func (m *Mailbox) Receive(msg Envelope) {
	m.queueMutex.Lock()
	defer m.queueMutex.Unlock()
	m.queue = append(m.queue, msg)
	m.dispatchMessages()
}

// Attempts to send all buffered messages to the actorChan.
func (m *Mailbox) dispatchMessages() {
	for len(m.queue) > 0 {
		select {
		case m.actorChan <- m.queue[0]:
			// Successfully sent the message to the actorChan, remove it from the queue
			m.queue = m.queue[1:]
		default:
			// The actorChan is full, stop attempting to send
			return
		}
	}
}
