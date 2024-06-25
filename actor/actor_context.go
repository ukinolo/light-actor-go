package actor

// Define the actorState constants
type actorState int32

const (
	actorStart actorState = iota
	actorStop
	actorStopping
)

// ActorContext holds the state and context of an actor
type ActorContext struct {
	actorSystem *ActorSystem
	props       *ActorProps
	envelope    Envelope
	state       actorState
	children    map[PID]bool
	self        PID
}

// NewActorContext creates and initializes a new actorContext
func NewActorContext(actorSystem *ActorSystem, props *ActorProps, self PID) *ActorContext {
	context := new(ActorContext)
	context.props = props
	context.state = actorStart
	context.actorSystem = actorSystem
	context.self = self
	context.children = make(map[PID]bool)
	return context
}

// Adds envelope to the current actor context
func (ctx *ActorContext) AddEnvelope(envelope Envelope) {
	ctx.envelope = envelope
}

// Spawns child actor
func (ctx *ActorContext) SpawnActor(actor Actor, props ...ActorProps) (PID, error) {
	prop := ConfigureActorProps(props...)
	prop.AddParent(&ctx.self)

	pid, err := ctx.actorSystem.SpawnActor(actor, *prop)
	if err != nil {
		return PID{}, err
	}

	ctx.children[pid] = true

	return pid, nil
}

// Send message
func (ctx *ActorContext) Send(message interface{}, receiver PID) {
	sendEnvelope := NewEnvelope(message, receiver)
	ctx.actorSystem.Send(sendEnvelope)
}

func (ctx *ActorContext) Message() interface{} {
	return ctx.envelope.Message
}

// Message returns the current message being processed
func (ctx *ActorContext) Envelope() Envelope {
	return ctx.envelope
}

// State returns the current state of the actor
func (ctx *ActorContext) State() actorState {
	return ctx.state
}

func (ctx *ActorContext) Self() PID {
	return ctx.self
}

func (ctx *ActorContext) GetChildren() map[PID]bool {
	return ctx.children
}
