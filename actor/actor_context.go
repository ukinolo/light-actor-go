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
	actor       Actor
	ctx         context.Context
	props       *Props
	envelope    Envelope
	state       actorState
	behavior    *Behavior // Add behavior as an attribute
}

// NewActorContext creates and initializes a new actorContext
func NewActorContext(ctx context.Context, actorSystem *ActorSystem, props *Props) *ActorContext {
	context := new(ActorContext)
	context.ctx = ctx
	context.props = props
	context.state = actorStart
	context.actorSystem = actorSystem
	context.behavior = NewBehavior() // Initialize behavior
	return context
}

// Context returns the context attached to the actorContext
func (ctx *ActorContext) Context() context.Context {
	return ctx.ctx
}

// Actor returns the actor associated with the actorContext
func (ctx *ActorContext) Actor() Actor {
	return ctx.actor
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

// Receive processes a message using the current behavior
func (ctx *ActorContext) Receive() {
	ctx.behavior.Receive(*ctx)
}
