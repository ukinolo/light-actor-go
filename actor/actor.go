package actor

type ActorProducer func() Actor

type Actor interface {
	Receive(ctx *ActorContext)
}
