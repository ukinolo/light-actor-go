package actor

import (
	"sync"
)

type Registry struct {
	mapping map[PID]chan Envelope // Stores mailbox channels
	mu      sync.RWMutex          // Mutex to make the registry thread-safe
}

func NewRegistry() *Registry {
	return &Registry{mapping: make(map[PID]chan Envelope)}
}

func (r *Registry) Add(pid PID, ch chan Envelope) error {
	r.mu.Lock()         // Lock the mutex for writing
	defer r.mu.Unlock() // Ensure the mutex is unlocked after the operation
	r.mapping[pid] = ch
	return nil
}

func (r *Registry) Find(pid PID) chan Envelope {
	r.mu.RLock()         // Lock the mutex for reading
	defer r.mu.RUnlock() // Ensure the mutex is unlocked after the operation
	if ch, ok := r.mapping[pid]; ok {
		return ch
	}
	return nil
}
