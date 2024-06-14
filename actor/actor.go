package actor

type ActorProducer func() Actor

type Actor interface {
	Recieve(ctx actorContext)
}
