package actor

type Producer func() Receiver

type Receiver interface {
	Recieve(ctx actorContext)
}
