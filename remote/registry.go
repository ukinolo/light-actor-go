package remote

import (
	"light-actor-go/actor"
	"sync"
)

type Registry struct {
	mapping map[string]actor.PID
	mu      sync.RWMutex
}

func NewRegistry() *Registry {
	return &Registry{mapping: make(map[string]actor.PID)}
}

func (r *Registry) Add(name string, pid actor.PID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.mapping[name] = pid
	return nil
}

func (r *Registry) Find(name string) actor.PID {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.mapping[name]
}
