package actor

import (
	"light-actor-go/envelope"
	"light-actor-go/pid"
)

type Registry struct {
	maping map[pid.PID](chan envelope.Envelope) //Stores mailbox chanels
}

func NewRegistry() *Registry {
	return &Registry{maping: make(map[pid.PID](chan envelope.Envelope))}
}

func (r *Registry) Add(pid pid.PID, ch chan envelope.Envelope) error {
	r.maping[pid] = ch
	return nil
}

func (r *Registry) Find(pid pid.PID) chan envelope.Envelope {
	return r.maping[pid]
}
