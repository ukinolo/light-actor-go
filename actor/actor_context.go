package actor

import (
	"context"
	"light-actor-go/envelope"
	"light-actor-go/pid"
)

// Define the actorState constants
type actorState int32

const (
	actorStart actorState = iota
	actorStop
)

// ActorContext holds the state and context of an actor
type ActorContext struct {
	actorSystem *ActorSystem
	ctx         context.Context
	props       *ActorProps
	envelope    envelope.Envelope
	state       actorState
	children    []pid.PID
	self        pid.PID
}

// NewActorContext creates and initializes a new actorContext
func NewActorContext(ctx context.Context, actorSystem *ActorSystem, props *ActorProps, self pid.PID) *ActorContext {
	context := new(ActorContext)
	context.ctx = ctx
	context.props = props
	context.state = actorStart
	context.actorSystem = actorSystem
	context.self = self
	return context
}

// Adds envelope to the current actor context
func (ctx *ActorContext) AddEnvelope(envelope envelope.Envelope) {
	ctx.envelope = envelope
}

// Spawns child actor
func (ctx *ActorContext) SpawnActor(actor Actor, props ...ActorProps) (pid.PID, error) {
	prop := ConfigureActorProps(props...)
	prop.AddParent(&ctx.self)

	id, err := pid.NewPID()
	if err != nil {
		return pid.PID{}, err
	}
	ctx.children = append(ctx.children, id)
	return id, nil
}

// Send message
func (ctx *ActorContext) Send(message interface{}, reciever pid.PID) {
	sendEnvelope := envelope.NewEnvelope(message, reciever)
	ctx.actorSystem.Send(sendEnvelope)
}

func (ctx *ActorContext) Message() interface{} {
	return ctx.envelope.Message
}

// Context returns the context attached to the actorContext
func (ctx *ActorContext) Context() context.Context {
	return ctx.ctx
}

// Message returns the current message being processed
func (ctx *ActorContext) Envelope() envelope.Envelope {
	return ctx.envelope
}

// State returns the current state of the actor
func (ctx *ActorContext) State() actorState {
	return ctx.state
}

func (ctx *ActorContext) Self() pid.PID {
	return ctx.self
}
