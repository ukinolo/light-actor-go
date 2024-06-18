package actor

import (
	"context"
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
	envelope    Envelope
	state       actorState
	behavior    *Behavior
	children    []PID
}

// NewActorContext creates and initializes a new actorContext
func NewActorContext(ctx context.Context, actorSystem *ActorSystem, props *ActorProps) *ActorContext {
	context := new(ActorContext)
	context.ctx = ctx
	context.props = props
	context.state = actorStart
	context.actorSystem = actorSystem
	context.behavior = NewBehavior() // Initialize behavior
	return context
}

// Adds envelope to the current actor context
func (ctx *ActorContext) AddEnvelope(envelope Envelope) {
	ctx.envelope = envelope
}

// Spawns child actor
func (ctx *ActorContext) SpawnActor(actor Actor) (PID, error) {
	id, err := NewPID()
	if err != nil {
		return PID{}, err
	}
	ctx.children = append(ctx.children, id)
	return id, nil
}

// Send message
func (ctx *ActorContext) Send(message interface{}, reciever PID) {
	sendEnvelope := NewEnvelope(message, reciever)
	ctx.actorSystem.Send(*sendEnvelope)
}

func (ctx *ActorContext) Message() interface{} {
	return ctx.envelope.message
}

// Context returns the context attached to the actorContext
func (ctx *ActorContext) Context() context.Context {
	return ctx.ctx
}

// Message returns the current message being processed
func (ctx *ActorContext) Envelope() Envelope {
	return ctx.envelope
}

// State returns the current state of the actor
func (ctx *ActorContext) State() actorState {
	return ctx.state
}

// Become sets the actor's behavior to a new one
func (ctx *ActorContext) Become(receive ReceiveFunc) {
	ctx.behavior.Become(receive)
}

// BecomeStacked adds a new behavior on top of the current behavior stack
func (ctx *ActorContext) BecomeStacked(receive ReceiveFunc) {
	ctx.behavior.BecomeStacked(receive)
}

// UnbecomeStacked removes the top behavior from the behavior stack
func (ctx *ActorContext) UnbecomeStacked() {
	ctx.behavior.UnbecomeStacked()
}
