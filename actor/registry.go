package actor

type Registry struct {
	maping map[PID](chan Envelope) //Stores mailbox chanels
}

func NewRegistry() *Registry {
	return &Registry{maping: make(map[PID](chan Envelope))}
}

func (r *Registry) Add(pid PID, ch chan Envelope) error {
	r.maping[pid] = ch
	return nil
}

func (r *Registry) Find(pid PID) chan Envelope {
	return r.maping[pid]
}
