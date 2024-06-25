package actor

import (
	"sync"
)

type Registry struct {
	mapping map[PID]chan Envelope // Stores mailbox channels
	mu      sync.RWMutex
}

func NewRegistry() *Registry {
	return &Registry{mapping: make(map[PID]chan Envelope)}
}

func (r *Registry) Add(pid PID, ch chan Envelope) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.mapping[pid] = ch
	return nil
}

func (r *Registry) Find(pid PID) chan Envelope {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.mapping[pid]
}

func (r *Registry) Delete(pid PID) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.mapping, pid)
}
