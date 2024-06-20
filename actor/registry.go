package actor

import (
	"light-actor-go/envelope"
	"light-actor-go/pid"
	"sync"
)

type Registry struct {
	mapping map[pid.PID]chan envelope.Envelope // Stores mailbox channels
	mu      sync.RWMutex                       // Mutex to make the registry thread-safe
}

func NewRegistry() *Registry {
	return &Registry{mapping: make(map[pid.PID]chan envelope.Envelope)}
}

func (r *Registry) Add(pid pid.PID, ch chan envelope.Envelope) error {
	r.mu.Lock()         // Lock the mutex for writing
	defer r.mu.Unlock() // Ensure the mutex is unlocked after the operation
	r.mapping[pid] = ch
	return nil
}

func (r *Registry) Find(pid pid.PID) chan envelope.Envelope {
	r.mu.RLock()         // Lock the mutex for reading
	defer r.mu.RUnlock() // Ensure the mutex is unlocked after the operation
	if ch, ok := r.mapping[pid]; ok {
		return ch
	}
	return nil
}
